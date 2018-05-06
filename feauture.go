package go_promise

import (
	"fmt"
	"log"
)

/*
	Create parent promise
 */
func NewPromise(onSuccess func(value interface{}) interface{}) *Promise {

	promise := newPromise(nil)
	promise.onSuccess = onSuccess
	go promise.process(nil)
	return promise
}

/*
	Wait execute all process or error
 */
func All(functions ...func(value interface{}) interface{}) *Promise {

	childs := make([]*Promise, len(functions), len(functions))

	promiseFun := func(d interface{}) interface{} {
		result := make([]interface{}, len(functions), len(functions))
		for i, onSuccess := range functions {
			childs[i] = NewPromise(onSuccess)
		}
		for i, child := range childs {
			value, err := child.Get()
			if err != nil {
				return err
			}
			result[i] = value
		}
		return result
	}

	return NewPromise(promiseFun)
}

/*
	get first success result
 */
func Race(functions ...func(value interface{}) interface{}) *Promise {

	childs := make([]*Promise, len(functions), len(functions))

	promiseFun := func(d interface{}) interface{} {

		result := make(chan *Promise)

		for i, onSuccess := range functions {
			childs[i] = NewPromise(onSuccess)
		}

		for _, child := range childs {
			go func(p *Promise) {
				select {
				case _, ok := <-p.final:
					if !ok {
						result <- p
					}
				}
			}(child)
		}

		for i := 1; i <= len(childs); i++ {

			p := <-result
			log.Printf("%v is first (%v)", p, p.state)
			if p.state == success {
				return p.result.value
			}
		}

		return fmt.Errorf("not success promises")
	}

	promise := NewPromise(promiseFun)
	log.Printf("%v is Race promise", promise)
	return promise
}

/*
	resolve data like JS
 */
func Resolve(d interface{}) *Promise {

	return NewPromise(func(value interface{}) interface{} { return d })
}

/*
	resolve data like JS
 */
func Reject(err error) *Promise {

	return Resolve(err)
}

/*
	Super function resolver for beautiful API

	Use:
		F(func(value interface{}) interface{})
		F(func(d interface{}) {...})
		F(func(){...} interface{})
		F(func(){...})
		F(value)
		F(error)
 */
func F(d interface{}) func(value interface{}) interface{} {

	switch d.(type) {
	case func(value interface{}) interface{}:
		return d.(func(value interface{}) interface{})
	case func(value interface{}):
		return func(value interface{}) interface{} {
			d.(func(value interface{}))(value)
			return nil
		}
	case func():
		return func(value interface{}) interface{} {
			d.(func())()
			return nil
		}
	case func() interface{}:
		return func(value interface{}) interface{} {
			return d.(func() interface{})()
		}
	default:
		return func(value interface{}) interface{} { return d }
	}
}
