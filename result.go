package go_promise

// result wrapper
type result struct {
	value      interface{}
	err        error
	promise    *Promise
	resultType resultType
}

func (r *result) copy() *result {

	return &result{
		resultType: r.resultType,
		value:      r.value,
		err:        r.err,
		// promise don't can be copy
	}
}

type resultType int

const (
	_       resultType = iota
	ERROR
	VALUE
	PROMISE
)

/*
	This method simplifies API

	Example, method onReject and onSuccess can use:
		return err
		return 5
		return promise

	This approach like original JS API
 */
func resolve(data interface{}) *result {

	switch data.(type) {
	case nil:
		return &result{
			resultType: VALUE,
			value:      nil,
		}
	case error:
		return &result{
			resultType: ERROR,
			err:        data.(error),
		}
	case *Promise:
		return &result{
			resultType: PROMISE,
			promise:    data.(*Promise),
		}
	case Promise:
		promise := data.(Promise)
		return &result{
			resultType: PROMISE,
			promise:    &promise,
		}
	default:
		return &result{
			resultType: VALUE,
			value:      data,
		}
	}
}
