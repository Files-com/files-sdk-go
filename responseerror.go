package files_sdk

import (
	"encoding/json"
	"fmt"
	"strings"
)

type ResponseError struct {
	ErrorMessage string   `json:"error"`
	HttpCode     string   `json:"http-code"`
	Errors       []string `json:"errors"`
}

func (e ResponseError) Error() string {
	if len(e.Errors) == 0 {
		return fmt.Sprintf("%v - http-code: %v", e.ErrorMessage, e.HttpCode)
	}
	return fmt.Sprintf(strings.Join(e.Errors, "\n"))
}

func (e ResponseError) IsNil() bool {
	return e.ErrorMessage == "" && len(e.Errors) == 0
}

func (e *ResponseError) UnmarshalJSON(data []byte) error {
	type re ResponseError
	var v re
	json.Unmarshal(data, &v)
	*e = ResponseError(v)
	if !e.IsNil() {
		return e
	}
	return nil
}
