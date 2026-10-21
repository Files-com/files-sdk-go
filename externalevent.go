package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// ExternalEvent is a Files.com API resource.
type ExternalEvent struct {
	// Event ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Type of event being recorded.
	EventType string `json:"event_type,omitempty" path:"event_type,omitempty" url:"event_type,omitempty"`
	// Status of event.
	Status string `json:"status,omitempty" path:"status,omitempty" url:"status,omitempty"`
	// Event body.
	Body string `json:"body,omitempty" path:"body,omitempty" url:"body,omitempty"`
	// External event create date/time.
	CreatedAt *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	// Link to log file.
	BodyUrl string `json:"body_url,omitempty" path:"body_url,omitempty" url:"body_url,omitempty"`
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
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `created_at`, `status` or `event_type`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `created_at` and `status`. Valid field combinations are `[ status, created_at ]`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// If set, return records where the specified field is greater than the supplied value. Valid fields are `created_at`.
	FilterGt interface{} `url:"filter_gt,omitempty" json:"filter_gt,omitempty" path:"filter_gt"`
	// If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `created_at`.
	FilterGteq interface{} `url:"filter_gteq,omitempty" json:"filter_gteq,omitempty" path:"filter_gteq"`
	// If set, return records where the specified field is less than the supplied value. Valid fields are `created_at`.
	FilterLt interface{} `url:"filter_lt,omitempty" json:"filter_lt,omitempty" path:"filter_lt"`
	// If set, return records where the specified field is less than or equal the supplied value. Valid fields are `created_at`.
	FilterLteq interface{} `url:"filter_lteq,omitempty" json:"filter_lteq,omitempty" path:"filter_lteq"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// ExternalEventFindParams contains the request parameters for GET /external_events/{id}.
type ExternalEventFindParams struct {
	// External Event ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// ExternalEventCreateParams contains the request parameters for POST /external_events.
type ExternalEventCreateParams struct {
	// Status of event. Values: success, failure, partial_failure, in_progress, skipped. Required.
	Status ExternalEventStatusEnum `url:"status" json:"status" path:"status"`
	// Event body. Required.
	Body string `url:"body" json:"body" path:"body"`
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
