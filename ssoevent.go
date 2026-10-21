package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// SsoEvent is a Files.com API resource.
type SsoEvent struct {
	// Event ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Type of SSO event being recorded.
	EventType string `json:"event_type,omitempty" path:"event_type,omitempty" url:"event_type,omitempty"`
	// Status of event.
	Status string `json:"status,omitempty" path:"status,omitempty" url:"status,omitempty"`
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
	// Username on Files.com for the SSO login attempt.
	Username string `json:"username,omitempty" path:"username,omitempty" url:"username,omitempty"`
	// Identity Provider UID for the SSO login attempt.
	IdpUid string `json:"idp_uid,omitempty" path:"idp_uid,omitempty" url:"idp_uid,omitempty"`
	// SSO provider for the SSO login attempt.
	Provider string `json:"provider,omitempty" path:"provider,omitempty" url:"provider,omitempty"`
	// SSO provider label for the SSO login attempt.
	ProviderLabel string `json:"provider_label,omitempty" path:"provider_label,omitempty" url:"provider_label,omitempty"`
	// IP address for the SSO login attempt.
	Ip string `json:"ip,omitempty" path:"ip,omitempty" url:"ip,omitempty"`
	// Region for the SSO login attempt.
	Region string `json:"region,omitempty" path:"region,omitempty" url:"region,omitempty"`
}

// Identifier returns the resource ID.
func (s SsoEvent) Identifier() interface{} {
	return s.Id
}

// SsoEventCollection is a list of SsoEvent resources.
type SsoEventCollection []SsoEvent

// SsoEventListParams contains the request parameters for GET /sso_events.
type SsoEventListParams struct {
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `created_at`, `event_type`, `status` or `user_id`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `created_at`, `event_type`, `idp_uid`, `ip`, `provider`, `status`, `user_id` or `username`. Valid field combinations are `[ event_type, created_at ]`, `[ idp_uid, created_at ]`, `[ ip, created_at ]`, `[ provider, created_at ]`, `[ status, created_at ]`, `[ user_id, created_at ]`, `[ username, created_at ]`, `[ event_type, status ]`, `[ idp_uid, status ]`, `[ ip, status ]`, `[ provider, status ]`, `[ user_id, status ]`, `[ username, status ]`, `[ event_type, status, created_at ]`, `[ idp_uid, status, created_at ]`, `[ ip, status, created_at ]`, `[ provider, status, created_at ]`, `[ user_id, status, created_at ]` or `[ username, status, created_at ]`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// If set, return records where the specified field is greater than the supplied value. Valid fields are `created_at`.
	FilterGt interface{} `url:"filter_gt,omitempty" json:"filter_gt,omitempty" path:"filter_gt"`
	// If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `created_at`.
	FilterGteq interface{} `url:"filter_gteq,omitempty" json:"filter_gteq,omitempty" path:"filter_gteq"`
	// If set, return records where the specified field is prefixed by the supplied value. Valid fields are `idp_uid`, `ip`, `provider` or `username`.
	FilterPrefix interface{} `url:"filter_prefix,omitempty" json:"filter_prefix,omitempty" path:"filter_prefix"`
	// If set, return records where the specified field is less than the supplied value. Valid fields are `created_at`.
	FilterLt interface{} `url:"filter_lt,omitempty" json:"filter_lt,omitempty" path:"filter_lt"`
	// If set, return records where the specified field is less than or equal the supplied value. Valid fields are `created_at`.
	FilterLteq interface{} `url:"filter_lteq,omitempty" json:"filter_lteq,omitempty" path:"filter_lteq"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// SsoEventFindParams contains the request parameters for GET /sso_events/{id}.
type SsoEventFindParams struct {
	// Sso Event ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (s *SsoEvent) UnmarshalJSON(data []byte) error {
	type ssoEvent SsoEvent
	var v ssoEvent
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*s = SsoEvent(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (s *SsoEventCollection) UnmarshalJSON(data []byte) error {
	type ssoEvents SsoEventCollection
	var v ssoEvents
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*s = SsoEventCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (s *SsoEventCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*s))
	for i, v := range *s {
		ret[i] = v
	}

	return &ret
}
