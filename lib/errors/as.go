package errors

import "errors"

func As[T any](err, target error) (t T, ok bool) {
	for err != nil {
		if e, ok := err.(T); ok {
			return e, true
		}
		if errors.Is(err, target) {
			return t, false
		}
		err = errors.Unwrap(err)
	}
	return t, false
}
