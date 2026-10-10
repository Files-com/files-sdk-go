package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// ExternalEvent is a Files.com API resource.
type ExternalEvent struct {
	Id        int64      `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	EventType string     `json:"event_type,omitempty" path:"event_type,omitempty" url:"event_type,omitempty"`
	Status    string     `json:"status,omitempty" path:"status,omitempty" url:"status,omitempty"`
	Body      string     `json:"body,omitempty" path:"body,omitempty" url:"body,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	BodyUrl   string     `json:"body_url,omitempty" path:"body_url,omitempty" url:"body_url,omitempty"`
}

// Identifier returns the resource ID.
func (e ExternalEvent) Identifier() interface{} {
	return e.Id
}

// ExternalEventCollection is a list of ExternalEvent resources.
type ExternalEventCollection []ExternalEvent

// ExternalEventStatusEnum is a string value for status.
// Enum lists the values documented by the API.
type ExternalEventStatusEnum string

// String returns the API parameter value.
func (u ExternalEventStatusEnum) String() string {
	return string(u)
}

// Enum returns the documented values keyed by their API strings.
func (u ExternalEventStatusEnum) Enum() map[string]ExternalEventStatusEnum {
	return map[string]ExternalEventStatusEnum{
		"success":         ExternalEventStatusEnum("success"),
		"failure":         ExternalEventStatusEnum("failure"),
		"partial_failure": ExternalEventStatusEnum("partial_failure"),
		"in_progress":     ExternalEventStatusEnum("in_progress"),
		"skipped":         ExternalEventStatusEnum("skipped"),
	}
}

// ExternalEventListParams contains the request parameters for GET /external_events.
type ExternalEventListParams struct {
	SortBy     interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	Filter     interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	FilterGt   interface{} `url:"filter_gt,omitempty" json:"filter_gt,omitempty" path:"filter_gt"`
	FilterGteq interface{} `url:"filter_gteq,omitempty" json:"filter_gteq,omitempty" path:"filter_gteq"`
	FilterLt   interface{} `url:"filter_lt,omitempty" json:"filter_lt,omitempty" path:"filter_lt"`
	FilterLteq interface{} `url:"filter_lteq,omitempty" json:"filter_lteq,omitempty" path:"filter_lteq"`
	ListParams
}

// ExternalEventFindParams contains the request parameters for GET /external_events/{id}.
type ExternalEventFindParams struct {
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// ExternalEventCreateParams contains the request parameters for POST /external_events.
type ExternalEventCreateParams struct {
	Status ExternalEventStatusEnum `url:"status" json:"status" path:"status"`
	Body   string                  `url:"body" json:"body" path:"body"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (e *ExternalEvent) UnmarshalJSON(data []byte) error {
	type externalEvent ExternalEvent
	var v externalEvent
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*e = ExternalEvent(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (e *ExternalEventCollection) UnmarshalJSON(data []byte) error {
	type externalEvents ExternalEventCollection
	var v externalEvents
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*e = ExternalEventCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (e *ExternalEventCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*e))
	for i, v := range *e {
		ret[i] = v
	}

	return &ret
}
