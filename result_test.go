package go_promise

import (
	"testing"
	"fmt"
	"github.com/stretchr/testify/assert"
)

func TestResolveResultByError(t *testing.T) {

	result := resolve(fmt.Errorf(testStr1))
	assert.Equal(t, ERROR, result.resultType)
	assert.Errorf(t, result.err, testErr1)
	assert.True(t, nil == result.promise)
	assert.Equal(t, nil, result.value)
}

func TestResolveResultByValue(t *testing.T) {

	result := resolve(testStr1)
	assert.Equal(t, VALUE, result.resultType)
	assert.NoError(t, result.err)
	assert.True(t, nil == result.promise)
	assert.Equal(t, testStr1, result.value)
}

func TestResolveResultByPromise(t *testing.T) {

	result := resolve(*NewPromise(F(testStr1)))
	assert.Equal(t, PROMISE, result.resultType)
	assert.NoError(t, result.err)
	assert.True(t, nil != result.promise)
	assert.Equal(t, nil, result.value)
}

func TestResolveResultByLinkPromise(t *testing.T) {

	result := resolve(NewPromise(F(testStr1)))
	assert.Equal(t, PROMISE, result.resultType)
	assert.NoError(t, result.err)
	assert.True(t, nil != result.promise)
	assert.Equal(t, nil, result.value)
}

func TestResolveResultByNil(t *testing.T) {

	result := resolve(nil)
	assert.Equal(t, VALUE, result.resultType)
	assert.NoError(t, result.err)
	assert.True(t, nil == result.promise)
	assert.Equal(t, nil, result.value)
}
