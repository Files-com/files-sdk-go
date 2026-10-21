package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// ScimLog is a Files.com API resource.
type ScimLog struct {
	// The unique ID of this SCIM request.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// The date and time when this SCIM request occurred.
	CreatedAt string `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	// The path portion of the URL requested.
	RequestPath string `json:"request_path,omitempty" path:"request_path,omitempty" url:"request_path,omitempty"`
	// The HTTP method used for this request.
	RequestMethod string `json:"request_method,omitempty" path:"request_method,omitempty" url:"request_method,omitempty"`
	// The HTTP response code returned for this request.
	HttpResponseCode string `json:"http_response_code,omitempty" path:"http_response_code,omitempty" url:"http_response_code,omitempty"`
	// The User-Agent header sent with the request.
	UserAgent string `json:"user_agent,omitempty" path:"user_agent,omitempty" url:"user_agent,omitempty"`
	// The JSON payload sent with the request.
	RequestJson string `json:"request_json,omitempty" path:"request_json,omitempty" url:"request_json,omitempty"`
	// The JSON payload returned in the response.
	ResponseJson string `json:"response_json,omitempty" path:"response_json,omitempty" url:"response_json,omitempty"`
}

// Identifier returns the resource ID.
func (s ScimLog) Identifier() interface{} {
	return s.Id
}

// ScimLogCollection is a list of ScimLog resources.
type ScimLogCollection []ScimLog

// ScimLogListParams contains the request parameters for GET /scim_logs.
type ScimLogListParams struct {
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `created_at`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// ScimLogFindParams contains the request parameters for GET /scim_logs/{id}.
type ScimLogFindParams struct {
	// Scim Log ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (s *ScimLog) UnmarshalJSON(data []byte) error {
	type scimLog ScimLog
	var v scimLog
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*s = ScimLog(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (s *ScimLogCollection) UnmarshalJSON(data []byte) error {
	type scimLogs ScimLogCollection
	var v scimLogs
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*s = ScimLogCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (s *ScimLogCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*s))
	for i, v := range *s {
		ret[i] = v
	}

	return &ret
}
