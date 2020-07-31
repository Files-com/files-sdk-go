package files_sdk

import (
	"encoding/json"
)

type Errors struct {
	Fields   []string `json:"fields,omitempty"`
	Messages []string `json:"messages,omitempty"`
}

type ErrorsCollection []Errors

func (e *Errors) UnmarshalJSON(data []byte) error {
	type errors Errors
	var v errors
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*e = Errors(v)
	return nil
}

func (e *ErrorsCollection) UnmarshalJSON(data []byte) error {
	type errorss []Errors
	var v errorss
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*e = ErrorsCollection(v)
	return nil
}
