package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// UserSecurityEvent is a Files.com API resource.
type UserSecurityEvent struct {
	Id          int64      `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	EventType   string     `json:"event_type,omitempty" path:"event_type,omitempty" url:"event_type,omitempty"`
	Body        string     `json:"body,omitempty" path:"body,omitempty" url:"body,omitempty"`
	EventErrors []string   `json:"event_errors,omitempty" path:"event_errors,omitempty" url:"event_errors,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	BodyUrl     string     `json:"body_url,omitempty" path:"body_url,omitempty" url:"body_url,omitempty"`
	UserId      int64      `json:"user_id,omitempty" path:"user_id,omitempty" url:"user_id,omitempty"`
}

// Identifier returns the resource ID.
func (u UserSecurityEvent) Identifier() interface{} {
	return u.Id
}

// UserSecurityEventCollection is a list of UserSecurityEvent resources.
type UserSecurityEventCollection []UserSecurityEvent

// UserSecurityEventListParams contains the request parameters for GET /user_security_events.
type UserSecurityEventListParams struct {
	SortBy     interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	Filter     interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	FilterGt   interface{} `url:"filter_gt,omitempty" json:"filter_gt,omitempty" path:"filter_gt"`
	FilterGteq interface{} `url:"filter_gteq,omitempty" json:"filter_gteq,omitempty" path:"filter_gteq"`
	FilterLt   interface{} `url:"filter_lt,omitempty" json:"filter_lt,omitempty" path:"filter_lt"`
	FilterLteq interface{} `url:"filter_lteq,omitempty" json:"filter_lteq,omitempty" path:"filter_lteq"`
	ListParams
}

// UserSecurityEventFindParams contains the request parameters for GET /user_security_events/{id}.
type UserSecurityEventFindParams struct {
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
