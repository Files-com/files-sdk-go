package files_sdk

import (
	"encoding/json"
)

type Status struct {
	Code          int64  `json:"code,omitempty"`
	Message       string `json:"message,omitempty"`
	Status        string `json:"status,omitempty"`
	Data          string `json:"data,omitempty"`
	Errors        string `json:"errors,omitempty"`
	ClickwrapId   int64  `json:"clickwrap_id,omitempty"`
	ClickwrapBody string `json:"clickwrap_body,omitempty"`
}

type StatusCollection []Status

func (s *Status) UnmarshalJSON(data []byte) error {
	type status Status
	var v status
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*s = Status(v)
	return nil
}

func (s *StatusCollection) UnmarshalJSON(data []byte) error {
	type statuss []Status
	var v statuss
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*s = StatusCollection(v)
	return nil
}

func (s *StatusCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*s))
	for i, v := range *s {
		ret[i] = v
	}

	return &ret
}
