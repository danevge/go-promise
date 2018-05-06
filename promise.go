package go_promise

import (
	"log"
	"time"
	"fmt"
)

const defaultTimeout = 250 * time.Millisecond

type TimeoutError error
type state int

const (
	_        state = iota
	pending
	success
	rejected
)

/*
	id - param for logging
	state - promise state
	result - end result all process for promise
	onSuccess - main function
	onReject - resolve error function
	final - broadcast about finalize all process about build end result
 */
type Promise struct {
	id        string
	state     state
	result    *result
	onSuccess func(value interface{}) interface{}
	onReject  func(err error) interface{}
	final     chan bool
}

/*
	Add new promise with handler for current promise

	JS example: promise.then( result => { ... });
 */
func (p *Promise) Then(onSuccess func(value interface{}) interface{}) *Promise {

	promise := newPromise(p)
	promise.onSuccess = onSuccess
	p.add(promise)
	return promise
}

/*
	Add new promise with handler and catch handler for current promise

	JS example:
					promise.then(
						result => ...,
						error => ...
					);
 */
func (p *Promise) ThenAndCatch(onSuccess func(value interface{}) interface{},
	onRejected func(err error) interface{}) *Promise {

	promise := newPromise(p)
	promise.onSuccess = onSuccess
	promise.onReject = onRejected
	p.add(promise)
	return promise
}

/*
	Add catch handler for current promise

	JS example: promise.catch(error => { ... });
 */
func (p *Promise) Catch(onRejected func(err error) interface{}) *Promise {

	p.onReject = onRejected
	return p
}

/*
	adaptation for backend
	return current value or timeout error
 */
func (p *Promise) Get() (interface{}, error) {

	return p.GetWithTimeout(defaultTimeout)
}

/*
	get result with custom timeout
 */
func (p *Promise) GetWithTimeout(timeout time.Duration) (interface{}, error) {

	select {
	case _, ok := <-p.final:
		if !ok {
			return p.result.value, p.result.err
		}
	case <-time.After(timeout):
		return nil, fmt.Errorf("timeout error")
	}
	return nil, TimeoutError(fmt.Errorf("timeout error"))
}

func (p *Promise) String() string {
	return fmt.Sprintf("Promise[id: %v; state: %v]", p.id, p.state)
}

/*
	default like JS
 */
func defaultOnRejected(err error) interface{} {

	return err
}

/*
	default like JS
 */
func defaultOnSuccess(value interface{}) interface{} {

	return value
}

func newPromise(parent *Promise) *Promise {

	oldId := ""
	if parent != nil {
		oldId = parent.id
	}
	return &Promise{
		id:        id(oldId),
		state:     pending,
		onSuccess: defaultOnSuccess,
		onReject:  defaultOnRejected,
		final:     make(chan bool, 1),
	}
}

func (p *Promise) process(oldResult *result) {

	log.Printf("%v - process", p)
	if oldResult == nil {
		// first promise and new promise in process line
		p.result = resolve(p.onSuccess(nil))
	} else {
		// let's see result previous promise
		switch oldResult.resultType {
		case ERROR:
			p.result = oldResult.copy()
		case VALUE:
			p.result = resolve(p.onSuccess(oldResult.value))
		case PROMISE:
			panic("promise result type 'PROMISE' is error!")
		default:
			panic("promise result type is undefined!")
		}
	}
	log.Printf("%v - promise is calculated %v", p, p.result)

	p.postProcess()
}

func (p *Promise) postProcess() {

	log.Printf("%v - post process", p)
	switch p.result.resultType {
	case ERROR:

		p.result = resolve(p.onReject(p.result.err))
		if p.result.resultType == ERROR {
			p.finalize(rejected)
			break
		}
		log.Printf("%v - resolve error, new result %v", p, p.result)
		p.postProcess()
	case PROMISE:
		p.processNewPromise()
	case VALUE:
		p.finalize(success)
	default:
		panic("promise result type is undefined!")
	}
}

func (p *Promise) processNewPromise() {
	newP := p.result.promise
	log.Printf("%v - wait result new promise", p)
	_, err := newP.Get()

	if _, ok := err.(TimeoutError); ok {
		log.Printf("%v - new promise %v fail by timeout", p, newP)
		p.result = &result{
			resultType: ERROR,
			err:        err,
		}
		return
	}

	log.Printf("%v -  change result from %v to %v", p, p.result, newP.result)
	p.result = newP.result.copy()
	p.postProcess()
}

/*
	Using close channel for send broadcast for all process
 */
func (p *Promise) finalize(state state) {
	log.Printf("%v - finalize to %v", p, state)
	p.state = state
	close(p.final)
}

/*
	Start or freeze start new promise for current promise
 */
func (p *Promise) add(child *Promise) {

	if p.state != pending {
		log.Printf("%v - start now", child)
		go child.process(p.result)
		return
	}
	log.Printf("%v - wait", child)
	go func() {
		_, ok := <-p.final
		if !ok {
			log.Printf("%v - freeze start", child)
			go child.process(p.result)
		}
	}()
}
