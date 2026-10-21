package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// PublicHostingRequestLog is a Files.com API resource.
type PublicHostingRequestLog struct {
	// Start Time of Action. Deprecrated: Use created_at.
	Timestamp *time.Time `json:"timestamp,omitempty" path:"timestamp,omitempty" url:"timestamp,omitempty"`
	// IP Address of Public Hosting Client.
	RemoteIp string `json:"remote_ip,omitempty" path:"remote_ip,omitempty" url:"remote_ip,omitempty"`
	// IP Address of Public Hosting Server.
	ServerIp string `json:"server_ip,omitempty" path:"server_ip,omitempty" url:"server_ip,omitempty"`
	// HTTP Request Hostname.
	Hostname string `json:"hostname,omitempty" path:"hostname,omitempty" url:"hostname,omitempty"`
	// HTTP Request Path. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.
	Path string `json:"path,omitempty" path:"path,omitempty" url:"path,omitempty"`
	// HTTP Response Code.
	ResponseCode int64 `json:"responseCode,omitempty" path:"responseCode,omitempty" url:"responseCode,omitempty"`
	// Whether SFTP Action was successful.
	Success *bool `json:"success,omitempty" path:"success,omitempty" url:"success,omitempty"`
	// Duration (in milliseconds).
	DurationMs int64 `json:"duration_ms,omitempty" path:"duration_ms,omitempty" url:"duration_ms,omitempty"`
	// Start Time of Action.
	CreatedAt *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	// The number of bytes transferred for file downloads.
	BytesTransferred int64 `json:"bytes_transferred,omitempty" path:"bytes_transferred,omitempty" url:"bytes_transferred,omitempty"`
	// Method of the HTTP call.
	HttpMethod string `json:"http_method,omitempty" path:"http_method,omitempty" url:"http_method,omitempty"`
}

// Identifier returns the resource path.
func (p PublicHostingRequestLog) Identifier() interface{} {
	return p.Path
}

// PublicHostingRequestLogCollection is a list of PublicHostingRequestLog resources.
type PublicHostingRequestLogCollection []PublicHostingRequestLog

// PublicHostingRequestLogListParams contains the request parameters for GET /public_hosting_request_logs.
type PublicHostingRequestLogListParams struct {
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `path`, `remote_ip`, `success` or `created_at`. Valid field combinations are `[ path ]`, `[ remote_ip ]`, `[ success ]`, `[ created_at ]`, `[ path, remote_ip ]`, `[ path, success ]`, `[ path, created_at ]`, `[ remote_ip, success ]`, `[ remote_ip, created_at ]`, `[ success, created_at ]`, `[ path, remote_ip, success ]`, `[ path, remote_ip, created_at ]`, `[ path, success, created_at ]`, `[ remote_ip, success, created_at ]` or `[ path, remote_ip, success, created_at ]`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// If set, return records where the specified field is greater than the supplied value. Valid fields are `created_at`. Valid field combinations are `[ path ]`, `[ remote_ip ]`, `[ success ]`, `[ created_at ]`, `[ path, remote_ip ]`, `[ path, success ]`, `[ path, created_at ]`, `[ remote_ip, success ]`, `[ remote_ip, created_at ]`, `[ success, created_at ]`, `[ path, remote_ip, success ]`, `[ path, remote_ip, created_at ]`, `[ path, success, created_at ]`, `[ remote_ip, success, created_at ]` or `[ path, remote_ip, success, created_at ]`.
	FilterGt interface{} `url:"filter_gt,omitempty" json:"filter_gt,omitempty" path:"filter_gt"`
	// If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `created_at`. Valid field combinations are `[ path ]`, `[ remote_ip ]`, `[ success ]`, `[ created_at ]`, `[ path, remote_ip ]`, `[ path, success ]`, `[ path, created_at ]`, `[ remote_ip, success ]`, `[ remote_ip, created_at ]`, `[ success, created_at ]`, `[ path, remote_ip, success ]`, `[ path, remote_ip, created_at ]`, `[ path, success, created_at ]`, `[ remote_ip, success, created_at ]` or `[ path, remote_ip, success, created_at ]`.
	FilterGteq interface{} `url:"filter_gteq,omitempty" json:"filter_gteq,omitempty" path:"filter_gteq"`
	// If set, return records where the specified field is prefixed by the supplied value. Valid fields are `path`. Valid field combinations are `[ path ]`, `[ remote_ip ]`, `[ success ]`, `[ created_at ]`, `[ path, remote_ip ]`, `[ path, success ]`, `[ path, created_at ]`, `[ remote_ip, success ]`, `[ remote_ip, created_at ]`, `[ success, created_at ]`, `[ path, remote_ip, success ]`, `[ path, remote_ip, created_at ]`, `[ path, success, created_at ]`, `[ remote_ip, success, created_at ]` or `[ path, remote_ip, success, created_at ]`.
	FilterPrefix interface{} `url:"filter_prefix,omitempty" json:"filter_prefix,omitempty" path:"filter_prefix"`
	// If set, return records where the specified field is less than the supplied value. Valid fields are `created_at`. Valid field combinations are `[ path ]`, `[ remote_ip ]`, `[ success ]`, `[ created_at ]`, `[ path, remote_ip ]`, `[ path, success ]`, `[ path, created_at ]`, `[ remote_ip, success ]`, `[ remote_ip, created_at ]`, `[ success, created_at ]`, `[ path, remote_ip, success ]`, `[ path, remote_ip, created_at ]`, `[ path, success, created_at ]`, `[ remote_ip, success, created_at ]` or `[ path, remote_ip, success, created_at ]`.
	FilterLt interface{} `url:"filter_lt,omitempty" json:"filter_lt,omitempty" path:"filter_lt"`
	// If set, return records where the specified field is less than or equal the supplied value. Valid fields are `created_at`. Valid field combinations are `[ path ]`, `[ remote_ip ]`, `[ success ]`, `[ created_at ]`, `[ path, remote_ip ]`, `[ path, success ]`, `[ path, created_at ]`, `[ remote_ip, success ]`, `[ remote_ip, created_at ]`, `[ success, created_at ]`, `[ path, remote_ip, success ]`, `[ path, remote_ip, created_at ]`, `[ path, success, created_at ]`, `[ remote_ip, success, created_at ]` or `[ path, remote_ip, success, created_at ]`.
	FilterLteq interface{} `url:"filter_lteq,omitempty" json:"filter_lteq,omitempty" path:"filter_lteq"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (p *PublicHostingRequestLog) UnmarshalJSON(data []byte) error {
	type publicHostingRequestLog PublicHostingRequestLog
	var v publicHostingRequestLog
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*p = PublicHostingRequestLog(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (p *PublicHostingRequestLogCollection) UnmarshalJSON(data []byte) error {
	type publicHostingRequestLogs PublicHostingRequestLogCollection
	var v publicHostingRequestLogs
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*p = PublicHostingRequestLogCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (p *PublicHostingRequestLogCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*p))
	for i, v := range *p {
		ret[i] = v
	}

	return &ret
}
