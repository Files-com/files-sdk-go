package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// Status is a Files.com API resource.
type Status struct {
	// Status HTTP code.
	Code int64 `json:"code,omitempty" path:"code,omitempty" url:"code,omitempty"`
	// Error message.
	Message string `json:"message,omitempty" path:"message,omitempty" url:"message,omitempty"`
	// Status message.
	Status string `json:"status,omitempty" path:"status,omitempty" url:"status,omitempty"`
	// Additional data.
	Data Auto `json:"data,omitempty" path:"data,omitempty" url:"data,omitempty"`
	// A list of api errors.
	Errors []Errors `json:"errors,omitempty" path:"errors,omitempty" url:"errors,omitempty"`
	// Required Clickwrap id.
	ClickwrapId int64 `json:"clickwrap_id,omitempty" path:"clickwrap_id,omitempty" url:"clickwrap_id,omitempty"`
	// Required Clickwrap body.
	ClickwrapBody string `json:"clickwrap_body,omitempty" path:"clickwrap_body,omitempty" url:"clickwrap_body,omitempty"`
}

// Identifier no path or id

// StatusCollection is a list of Status resources.
type StatusCollection []Status

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (s *Status) UnmarshalJSON(data []byte) error {
	type status Status
	var v status
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*s = Status(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (s *StatusCollection) UnmarshalJSON(data []byte) error {
	type statuss StatusCollection
	var v statuss
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*s = StatusCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (s *StatusCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*s))
	for i, v := range *s {
		ret[i] = v
	}

	return &ret
}
