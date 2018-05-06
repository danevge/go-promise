package go_promise

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
	"time"
)

const (
	testStr1 = "aaa_a-1"
	testStr2 = "bbb_b-2"
	testStr3 = "ccc_c-3"
	testStr4 = "ddd_d-4"
	testStr5 = "eee_e-3"
	testErr1 = "ups 1"
	testErr2 = "ups 2"
	testErr3 = "ups 2"
)

func TestSuccessAlonePromise(t *testing.T) {

	value, err := NewPromise(func(d interface{}) interface{} { return testStr1 }).
		Get()

	assert.Equal(t, testStr1, value)
	assert.NoError(t, err)
}

func TestAddSimpleThen(t *testing.T) {

	value, err := NewPromise(func(d interface{}) interface{} { return testStr1 }).
		Then(func(d interface{}) interface{} { return testStr2 }).
		Get()

	assert.Equal(t, testStr2, value)
	assert.NoError(t, err)
}

func TestErrorAlonePromiseWithoutCatch(t *testing.T) {

	value, err := NewPromise(func(d interface{}) interface{} { return fmt.Errorf(testErr1) }).
		Get()

	assert.Equal(t, nil, value)
	assert.Errorf(t, err, testErr1)
}

func TestErrorAlonePromiseWithCatch(t *testing.T) {

	value, err := NewPromise(func(d interface{}) interface{} { return fmt.Errorf(testErr1) }).
		Catch(func(err error) interface{} { return fmt.Errorf(testErr2) }).
		Get()

	assert.Equal(t, nil, value)
	assert.Errorf(t, err, testErr2)
}

func TestDoubleThen(t *testing.T) {

	parent := NewPromise(func(d interface{}) interface{} { return testStr1 })
	value, err := parent.Get()
	value1, err1 := parent.Then(func(d interface{}) interface{} { return testStr2 }).Get()
	value2, err2 := parent.Then(func(d interface{}) interface{} { return testStr3 }).Get()

	assert.Equal(t, testStr1, value)
	assert.Equal(t, testStr2, value1)
	assert.Equal(t, testStr3, value2)
	assert.NoError(t, err)
	assert.NoError(t, err1)
	assert.NoError(t, err2)
}

func TestAddSimpleThenAndCatchByData(t *testing.T) {

	value, err := NewPromise(func(d interface{}) interface{} { return testStr1 }).
		ThenAndCatch(
		func(d interface{}) interface{} { return testStr2 },
		func(err error) interface{} { return fmt.Errorf(testErr2) }).
		Get()

	assert.Equal(t, testStr2, value)
	assert.NoError(t, err)
}

func TestAddSimpleThenAndCatchByError(t *testing.T) {

	value, err := NewPromise(func(d interface{}) interface{} { return fmt.Errorf(testErr1) }).
		ThenAndCatch(
		func(d interface{}) interface{} { return testStr2 },
		func(err error) interface{} { return fmt.Errorf(testErr2) }).
		Get()

	assert.Equal(t, nil, value)
	assert.Errorf(t, err, testErr2)
}

func TestUseParentPromiseByData(t *testing.T) {

	value, err := NewPromise(func(d interface{}) interface{} { return testStr1 }).
		Then(func(d interface{}) interface{} { return d.(string) + testStr2 }).
		Get()

	assert.Equal(t, testStr1+testStr2, value)
	assert.NoError(t, err)
}

func TestUseParentPromiseByError(t *testing.T) {

	value, err := NewPromise(func(d interface{}) interface{} { return fmt.Errorf(testErr1) }).
		ThenAndCatch(
		func(d interface{}) interface{} { return testStr2 },
		func(err error) interface{} { return err }).
		Get()

	assert.Equal(t, nil, value)
	assert.Errorf(t, err, testErr1)
}

func TestUseParentData(t *testing.T) {

	value, err := NewPromise(func(d interface{}) interface{} { return testStr1 }).
		Then(func(d interface{}) interface{} { return d.(string) + testStr2 }).
		Get()

	assert.Equal(t, testStr1+testStr2, value)
	assert.NoError(t, err)
}

func TestAlonePromiseResolveErrorToValue(t *testing.T) {

	value, err := NewPromise(func(d interface{}) interface{} { return fmt.Errorf(testErr1) }).
		Catch(func(err2 error) interface{} { return testStr3 }).
		Get()

	assert.Equal(t, testStr3, value)
	assert.NoError(t, err)
}

func TestErrorDoublePromiseByNotCatch(t *testing.T) {

	parent := NewPromise(func(d interface{}) interface{} { return fmt.Errorf(testErr1) })
	value, err := parent.Get()
	value1, err1 := parent.Then(func(d interface{}) interface{} { return testStr2 }).Get()
	value2, err2 := parent.Then(func(d interface{}) interface{} { return testStr3 }).Get()

	assert.Equal(t, nil, value)
	assert.Equal(t, nil, value1)
	assert.Equal(t, nil, value2)
	assert.Errorf(t, err, testErr1)
	assert.Errorf(t, err1, testErr1)
	assert.Errorf(t, err2, testErr1)
}

func TestErrorManyPromiseByNotCatch(t *testing.T) {

	value, err := NewPromise(func(d interface{}) interface{} { return fmt.Errorf(testErr1) }).
		Then(func(d interface{}) interface{} { return testStr2 }).
		Then(func(d interface{}) interface{} { return testStr3 }).
		Then(func(d interface{}) interface{} { return testStr4 }).
		Then(func(d interface{}) interface{} { return testStr5 }).
		Get()

	assert.Equal(t, nil, value)
	assert.Errorf(t, err, testErr1)
}

func TestChildResolveErrorToValue(t *testing.T) {

	value, err := NewPromise(func(d interface{}) interface{} { return fmt.Errorf(testErr1) }).
		Then(func(d interface{}) interface{} { return testStr2 }).
		Catch(func(err error) interface{} { return testStr3 }).
		Get()

	assert.Equal(t, testStr3, value)
	assert.NoError(t, err)
}

func TestThenForProcessedPromiseWithUseOldData(t *testing.T) {

	parent := NewPromise(func(d interface{}) interface{} { return testStr1 })
	value, err := parent.Then(func(d interface{}) interface{} { return testStr2 }).Get()

	assert.Equal(t, testStr2, value)
	assert.NoError(t, err)

	value, err = parent.Then(func(d interface{}) interface{} { return d.(string) + testStr3 }).Get()

	assert.Equal(t, testStr1+testStr3, value)
	assert.NoError(t, err)
}

func TestAlonePromiseReturnNewPromise(t *testing.T) {
	value, err := NewPromise(func(d interface{}) interface{} {
		return NewPromise(func(d interface{}) interface{} {
			return testStr1
		})
	}).Get()

	assert.Equal(t, testStr1, value)
	assert.NoError(t, err)
}

func TestThenNewPromise(t *testing.T) {
	value, err := NewPromise(func(d interface{}) interface{} { return testStr1 }).
		Then(func(d interface{}) interface{} {
		return NewPromise(func(d interface{}) interface{} { return testStr2 })
	}).Get()

	assert.Equal(t, testStr2, value)
	assert.NoError(t, err)
}

func TestCatchNewPromise(t *testing.T) {
	value, err := NewPromise(func(d interface{}) interface{} { return testStr1 }).
		Then(func(d interface{}) interface{} { return fmt.Errorf(testErr1) }).
		Catch(func(err error) interface{} {
		return NewPromise(func(d interface{}) interface{} { return testStr3 })
	}).Get()

	assert.Equal(t, testStr3, value)
	assert.NoError(t, err)
}

func TestTimeoutError(t *testing.T) {
	value, err := NewPromise(func(d interface{}) interface{} {
		time.Sleep(time.Second)
		return testStr1
	}).GetWithTimeout(100 * time.Millisecond)

	assert.Equal(t, nil, value)
	assert.Errorf(t, err, "timeout error")
}

func TestTimeout(t *testing.T) {
	value, err := NewPromise(func(d interface{}) interface{} {
		time.Sleep(400 * time.Millisecond)
		return testStr1
	}).GetWithTimeout(500 * time.Millisecond)

	assert.Equal(t, testStr1, value)
	assert.NoError(t, err)
}

func TestTimeoutErrorInNewPromise(t *testing.T) {

	value, err := NewPromise(func(d interface{}) interface{} {
		return NewPromise(func(d interface{}) interface{} {
			time.Sleep(500 * time.Millisecond)
			return testStr1
		})
	}).Get()

	assert.Equal(t, nil, value)
	assert.Errorf(t, err, "timeout error")
}

func TestResolveErrorToValueAndUseValue(t *testing.T) {

	value, err := NewPromise(func(d interface{}) interface{} { return fmt.Errorf(testErr1) }).
		Then(func(d interface{}) interface{} { return testStr2 }).
		Catch(func(err error) interface{} { return testStr3 }).
		Then(func(d interface{}) interface{} { return d.(string) + testStr5 }).
		Get()

	assert.Equal(t, testStr3+testStr5, value)
	assert.NoError(t, err)
}
