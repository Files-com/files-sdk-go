package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// ScimLog is a Files.com API resource.
type ScimLog struct {
	Id               int64  `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	CreatedAt        string `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	RequestPath      string `json:"request_path,omitempty" path:"request_path,omitempty" url:"request_path,omitempty"`
	RequestMethod    string `json:"request_method,omitempty" path:"request_method,omitempty" url:"request_method,omitempty"`
	HttpResponseCode string `json:"http_response_code,omitempty" path:"http_response_code,omitempty" url:"http_response_code,omitempty"`
	UserAgent        string `json:"user_agent,omitempty" path:"user_agent,omitempty" url:"user_agent,omitempty"`
	RequestJson      string `json:"request_json,omitempty" path:"request_json,omitempty" url:"request_json,omitempty"`
	ResponseJson     string `json:"response_json,omitempty" path:"response_json,omitempty" url:"response_json,omitempty"`
}

// Identifier returns the resource ID.
func (s ScimLog) Identifier() interface{} {
	return s.Id
}

// ScimLogCollection is a list of ScimLog resources.
type ScimLogCollection []ScimLog

// ScimLogListParams contains the request parameters for GET /scim_logs.
type ScimLogListParams struct {
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	ListParams
}

// ScimLogFindParams contains the request parameters for GET /scim_logs/{id}.
type ScimLogFindParams struct {
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
