package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// ExavaultApiRequestLog is a Files.com API resource.
type ExavaultApiRequestLog struct {
	// Start Time of Action. Deprecrated: Use created_at.
	Timestamp *time.Time `json:"timestamp,omitempty" path:"timestamp,omitempty" url:"timestamp,omitempty"`
	// Name of API Endpoint.
	Endpoint string `json:"endpoint,omitempty" path:"endpoint,omitempty" url:"endpoint,omitempty"`
	// Exavault API Version.
	Version float64 `json:"version,omitempty" path:"version,omitempty" url:"version,omitempty"`
	// IP of requesting client.
	RequestIp string `json:"request_ip,omitempty" path:"request_ip,omitempty" url:"request_ip,omitempty"`
	// HTTP Method.
	RequestMethod string `json:"request_method,omitempty" path:"request_method,omitempty" url:"request_method,omitempty"`
	// Error type, if applicable.
	ErrorType string `json:"error_type,omitempty" path:"error_type,omitempty" url:"error_type,omitempty"`
	// Error message, if applicable.
	ErrorMessage string `json:"error_message,omitempty" path:"error_message,omitempty" url:"error_message,omitempty"`
	// User-Agent.
	UserAgent string `json:"user_agent,omitempty" path:"user_agent,omitempty" url:"user_agent,omitempty"`
	// HTTP Response Code.
	ResponseCode int64 `json:"response_code,omitempty" path:"response_code,omitempty" url:"response_code,omitempty"`
	// `false` if HTTP Response Code is 4xx or 5xx.
	Success *bool `json:"success,omitempty" path:"success,omitempty" url:"success,omitempty"`
	// Duration (in milliseconds).
	DurationMs int64 `json:"duration_ms,omitempty" path:"duration_ms,omitempty" url:"duration_ms,omitempty"`
	// Start Time of Action.
	CreatedAt *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
}

// Identifier no path or id

// ExavaultApiRequestLogCollection is a list of ExavaultApiRequestLog resources.
type ExavaultApiRequestLogCollection []ExavaultApiRequestLog

// ExavaultApiRequestLogListParams contains the request parameters for GET /exavault_api_request_logs.
type ExavaultApiRequestLogListParams struct {
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `request_ip`, `request_method`, `success` or `created_at`. Valid field combinations are `[ request_ip ]`, `[ request_method ]`, `[ success ]`, `[ created_at ]`, `[ request_ip, request_method ]`, `[ request_ip, success ]`, `[ request_ip, created_at ]`, `[ request_method, success ]`, `[ request_method, created_at ]`, `[ success, created_at ]`, `[ request_ip, request_method, success ]`, `[ request_ip, request_method, created_at ]`, `[ request_ip, success, created_at ]`, `[ request_method, success, created_at ]` or `[ request_ip, request_method, success, created_at ]`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// If set, return records where the specified field is greater than the supplied value. Valid fields are `created_at`. Valid field combinations are `[ request_ip ]`, `[ request_method ]`, `[ success ]`, `[ created_at ]`, `[ request_ip, request_method ]`, `[ request_ip, success ]`, `[ request_ip, created_at ]`, `[ request_method, success ]`, `[ request_method, created_at ]`, `[ success, created_at ]`, `[ request_ip, request_method, success ]`, `[ request_ip, request_method, created_at ]`, `[ request_ip, success, created_at ]`, `[ request_method, success, created_at ]` or `[ request_ip, request_method, success, created_at ]`.
	FilterGt interface{} `url:"filter_gt,omitempty" json:"filter_gt,omitempty" path:"filter_gt"`
	// If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `created_at`. Valid field combinations are `[ request_ip ]`, `[ request_method ]`, `[ success ]`, `[ created_at ]`, `[ request_ip, request_method ]`, `[ request_ip, success ]`, `[ request_ip, created_at ]`, `[ request_method, success ]`, `[ request_method, created_at ]`, `[ success, created_at ]`, `[ request_ip, request_method, success ]`, `[ request_ip, request_method, created_at ]`, `[ request_ip, success, created_at ]`, `[ request_method, success, created_at ]` or `[ request_ip, request_method, success, created_at ]`.
	FilterGteq interface{} `url:"filter_gteq,omitempty" json:"filter_gteq,omitempty" path:"filter_gteq"`
	// If set, return records where the specified field is prefixed by the supplied value. Valid fields are `request_ip` and `request_method`. Valid field combinations are `[ request_ip ]`, `[ request_method ]`, `[ success ]`, `[ created_at ]`, `[ request_ip, request_method ]`, `[ request_ip, success ]`, `[ request_ip, created_at ]`, `[ request_method, success ]`, `[ request_method, created_at ]`, `[ success, created_at ]`, `[ request_ip, request_method, success ]`, `[ request_ip, request_method, created_at ]`, `[ request_ip, success, created_at ]`, `[ request_method, success, created_at ]` or `[ request_ip, request_method, success, created_at ]`.
	FilterPrefix interface{} `url:"filter_prefix,omitempty" json:"filter_prefix,omitempty" path:"filter_prefix"`
	// If set, return records where the specified field is less than the supplied value. Valid fields are `created_at`. Valid field combinations are `[ request_ip ]`, `[ request_method ]`, `[ success ]`, `[ created_at ]`, `[ request_ip, request_method ]`, `[ request_ip, success ]`, `[ request_ip, created_at ]`, `[ request_method, success ]`, `[ request_method, created_at ]`, `[ success, created_at ]`, `[ request_ip, request_method, success ]`, `[ request_ip, request_method, created_at ]`, `[ request_ip, success, created_at ]`, `[ request_method, success, created_at ]` or `[ request_ip, request_method, success, created_at ]`.
	FilterLt interface{} `url:"filter_lt,omitempty" json:"filter_lt,omitempty" path:"filter_lt"`
	// If set, return records where the specified field is less than or equal the supplied value. Valid fields are `created_at`. Valid field combinations are `[ request_ip ]`, `[ request_method ]`, `[ success ]`, `[ created_at ]`, `[ request_ip, request_method ]`, `[ request_ip, success ]`, `[ request_ip, created_at ]`, `[ request_method, success ]`, `[ request_method, created_at ]`, `[ success, created_at ]`, `[ request_ip, request_method, success ]`, `[ request_ip, request_method, created_at ]`, `[ request_ip, success, created_at ]`, `[ request_method, success, created_at ]` or `[ request_ip, request_method, success, created_at ]`.
	FilterLteq interface{} `url:"filter_lteq,omitempty" json:"filter_lteq,omitempty" path:"filter_lteq"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (e *ExavaultApiRequestLog) UnmarshalJSON(data []byte) error {
	type exavaultApiRequestLog ExavaultApiRequestLog
	var v exavaultApiRequestLog
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*e = ExavaultApiRequestLog(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (e *ExavaultApiRequestLogCollection) UnmarshalJSON(data []byte) error {
	type exavaultApiRequestLogs ExavaultApiRequestLogCollection
	var v exavaultApiRequestLogs
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*e = ExavaultApiRequestLogCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (e *ExavaultApiRequestLogCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*e))
	for i, v := range *e {
		ret[i] = v
	}

	return &ret
}
