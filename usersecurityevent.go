package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// UserSecurityEvent is a Files.com API resource.
type UserSecurityEvent struct {
	// Event ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Type of user security event being recorded.
	EventType string `json:"event_type,omitempty" path:"event_type,omitempty" url:"event_type,omitempty"`
	// Event body.
	Body string `json:"body,omitempty" path:"body,omitempty" url:"body,omitempty"`
	// Event errors.
	EventErrors []string `json:"event_errors,omitempty" path:"event_errors,omitempty" url:"event_errors,omitempty"`
	// Event create date/time.
	CreatedAt *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	// Link to log file.
	BodyUrl string `json:"body_url,omitempty" path:"body_url,omitempty" url:"body_url,omitempty"`
	// User ID.
	UserId int64 `json:"user_id,omitempty" path:"user_id,omitempty" url:"user_id,omitempty"`
}

// Identifier returns the resource ID.
func (u UserSecurityEvent) Identifier() interface{} {
	return u.Id
}

// UserSecurityEventCollection is a list of UserSecurityEvent resources.
type UserSecurityEventCollection []UserSecurityEvent

// UserSecurityEventListParams contains the request parameters for GET /user_security_events.
type UserSecurityEventListParams struct {
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `created_at` and `user_id`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `created_at` and `user_id`. Valid field combinations are `[ user_id, created_at ]`.
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

// UserSecurityEventFindParams contains the request parameters for GET /user_security_events/{id}.
type UserSecurityEventFindParams struct {
	// User Security Event ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (u *UserSecurityEvent) UnmarshalJSON(data []byte) error {
	type userSecurityEvent UserSecurityEvent
	var v userSecurityEvent
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*u = UserSecurityEvent(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (u *UserSecurityEventCollection) UnmarshalJSON(data []byte) error {
	type userSecurityEvents UserSecurityEventCollection
	var v userSecurityEvents
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*u = UserSecurityEventCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (u *UserSecurityEventCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*u))
	for i, v := range *u {
		ret[i] = v
	}

	return &ret
}
