package go_promise

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
	"time"
)

func TestAll(t *testing.T) {

	value, err := All(
		func(d interface{}) interface{} { return testStr1 },
		func(d interface{}) interface{} { return testStr2 },
		func(d interface{}) interface{} { return testStr3 },
		func(d interface{}) interface{} { return testStr4 },
		func(d interface{}) interface{} { return testStr5 },
	).Get()

	testArray := []interface{}{testStr1, testStr2, testStr3, testStr4, testStr5}
	assert.Equal(t, testArray, value)
	assert.NoError(t, err)
}

func TestAllByError(t *testing.T) {

	value, err := All(
		func(d interface{}) interface{} { return testStr1 },
		func(d interface{}) interface{} { return testStr2 },
		func(d interface{}) interface{} { return fmt.Errorf(testErr1) },
		func(d interface{}) interface{} { return testStr4 },
		func(d interface{}) interface{} { return testStr5 },
	).Get()

	assert.Equal(t, nil, value)
	assert.Errorf(t, err, testErr1)
}

func TestAllAndUseResult(t *testing.T) {

	value, err := All(
		func(d interface{}) interface{} { return testStr1 },
		func(d interface{}) interface{} { return testStr2 },
		func(d interface{}) interface{} { return testStr3 },
		func(d interface{}) interface{} { return testStr4 },
		func(d interface{}) interface{} { return testStr5 },
	).Then(func(d interface{}) interface{} {
		result := ""
		array := d.([]interface{})
		for _, str := range array {
			result += str.(string)
		}
		return result
	}).Get()

	assert.Equal(t, testStr1+testStr2+testStr3+testStr4+testStr5, value)
	assert.NoError(t, err)
}

func TestResolve(t *testing.T) {

	value, err := Resolve(testStr1).Get()
	assert.Equal(t, testStr1, value)
	assert.NoError(t, err)
}

func TestReject(t *testing.T) {

	value, err := Reject(fmt.Errorf(testErr1)).Get()
	assert.Equal(t, nil, value)
	assert.Errorf(t, err, testErr1)
}

func TestRaceByRaceCondition(t *testing.T) {

	value, err := Race(
		func(d interface{}) interface{} {
			time.Sleep(600 * time.Millisecond)
			return testStr1
		},
		func(d interface{}) interface{} {
			time.Sleep(200 * time.Millisecond)
			return testStr2
		},
		func(d interface{}) interface{} {
			time.Sleep(800 * time.Millisecond)
			return testStr3
		},
	).GetWithTimeout(2 * time.Second)

	assert.Equal(t, testStr2, value)
	assert.NoError(t, err)
}

func TestRaceByError(t *testing.T) {

	value, err := Race(
		func(d interface{}) interface{} {
			time.Sleep(600 * time.Millisecond)
			return testStr1
		},
		func(d interface{}) interface{} {
			time.Sleep(200 * time.Millisecond)
			return fmt.Errorf(testErr1)
		},
		func(d interface{}) interface{} {
			time.Sleep(800 * time.Millisecond)
			return testStr3
		},
	).GetWithTimeout(2 * time.Second)

	assert.Equal(t, testStr1, value)
	assert.NoError(t, err)
}

func TestRaceByAllError(t *testing.T) {

	value, err := Race(
		func(d interface{}) interface{} {
			time.Sleep(600 * time.Millisecond)
			return fmt.Errorf(testErr1)
		},
		func(d interface{}) interface{} {
			time.Sleep(200 * time.Millisecond)
			return fmt.Errorf(testErr2)
		},
		func(d interface{}) interface{} {
			time.Sleep(800 * time.Millisecond)
			return fmt.Errorf(testErr3)
		},
	).GetWithTimeout(2 * time.Second)

	assert.Equal(t, nil, value)
	assert.Errorf(t, err, "not success promises")
}

func TestRaceByTimeout(t *testing.T) {

	value, err := Race(
		func(d interface{}) interface{} {
			time.Sleep(500 * time.Millisecond)
			return fmt.Errorf(testErr1)
		},
		func(d interface{}) interface{} {
			time.Sleep(400 * time.Millisecond)
			return fmt.Errorf(testErr2)
		},
		func(d interface{}) interface{} {
			time.Sleep(800 * time.Millisecond)
			return fmt.Errorf(testErr3)
		},
	).Get()

	assert.Equal(t, nil, value)
	assert.Errorf(t, err, "timeout error")
}

func TestResolveNormFunc(t *testing.T) {

	value, err := NewPromise(F(func(d interface{}) interface{} { return testStr1 })).Get()
	assert.Equal(t, testStr1, value)
	assert.NoError(t, err)
}

func TestResolveFuncWithoutInputParams(t *testing.T) {

	value, err := NewPromise(F(func() interface{} { return testStr1 })).Get()
	assert.Equal(t, testStr1, value)
	assert.NoError(t, err)
}

func TestResolveFuncWithoutOutputParam(t *testing.T) {

	a := "not a"
	value, err := NewPromise(F(func(d interface{}) { a = "a" })).Get()
	assert.Equal(t, "a", a)
	assert.Equal(t, nil, value)
	assert.NoError(t, err)
}

func TestResolveSimpleFunc(t *testing.T) {

	a := "not a"
	value, err := NewPromise(F(func() { a = "a" })).Get()
	assert.Equal(t, "a", a)
	assert.Equal(t, nil, value)
	assert.NoError(t, err)
}


func TestResolveError(t *testing.T) {

	value, err := NewPromise(F(fmt.Errorf(testErr1))).Get()
	assert.Equal(t, nil, value)
	assert.Errorf(t, err, testErr1)
}


func TestResolveValue(t *testing.T) {

	value, err := NewPromise(F(testStr1)).Get()
	assert.Equal(t, testStr1, value)
	assert.NoError(t, err)
}
