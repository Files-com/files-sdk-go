package lib

import (
	"encoding/json"
)

type ErrorWithOriginalResponse struct {
	error
	originalResponse interface{}
}

func (u ErrorWithOriginalResponse) OriginalResponse() interface{} {
	return u.originalResponse
}

func (u ErrorWithOriginalResponse) ProcessError(data []byte, err error, t interface{}) error {
	unmarshalError, ok := err.(*json.UnmarshalTypeError)

	if ok {
		ignoreErr := json.Unmarshal(data, &t)
		if ignoreErr == nil {
			return ErrorWithOriginalResponse{error: unmarshalError, originalResponse: t}
		}
	}
	errorWithOriginalResponse, ok := err.(ErrorWithOriginalResponse)
	if ok {
		ignoreErr := json.Unmarshal(data, &t)
		if ignoreErr == nil {
			return ErrorWithOriginalResponse{error: errorWithOriginalResponse.error, originalResponse: t}
		}
	}
	return err
}
