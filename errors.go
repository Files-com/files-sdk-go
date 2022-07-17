package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Errors struct {
	Fields   []string `json:"fields,omitempty" path:"fields"`
	Messages []string `json:"messages,omitempty" path:"messages"`
}

type ErrorsCollection []Errors

func (e *Errors) UnmarshalJSON(data []byte) error {
	type errors Errors
	var v errors
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*e = Errors(v)
	return nil
}

func (e *ErrorsCollection) UnmarshalJSON(data []byte) error {
	type errorss ErrorsCollection
	var v errorss
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*e = ErrorsCollection(v)
	return nil
}

func (e *ErrorsCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*e))
	for i, v := range *e {
		ret[i] = v
	}

	return &ret
}
