package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// Request is a Files.com API resource.
type Request struct {
	// Request ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Folder path. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.
	Path string `json:"path,omitempty" path:"path,omitempty" url:"path,omitempty"`
	// Source filename, if applicable.
	Source string `json:"source,omitempty" path:"source,omitempty" url:"source,omitempty"`
	// Destination filename.
	Destination string `json:"destination,omitempty" path:"destination,omitempty" url:"destination,omitempty"`
	// ID of automation that created request.
	AutomationId int64 `json:"automation_id,omitempty" path:"automation_id,omitempty" url:"automation_id,omitempty"`
	// User making the request (if applicable).
	UserDisplayName string `json:"user_display_name,omitempty" path:"user_display_name,omitempty" url:"user_display_name,omitempty"`
	// A list of user IDs to request the file from. If sent as a string, it should be comma-delimited.
	UserIds string `json:"user_ids,omitempty" path:"user_ids,omitempty" url:"user_ids,omitempty"`
	// A list of group IDs to request the file from. If sent as a string, it should be comma-delimited.
	GroupIds string `json:"group_ids,omitempty" path:"group_ids,omitempty" url:"group_ids,omitempty"`
}

// Identifier returns the resource ID.
func (r Request) Identifier() interface{} {
	return r.Id
}

// RequestCollection is a list of Request resources.
type RequestCollection []Request

// RequestListParams contains the request parameters for GET /requests.
type RequestListParams struct {
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are .
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// Only show requests of the current user? (Defaults to true if current user is not a site admin.).
	Mine *bool `url:"mine,omitempty" json:"mine,omitempty" path:"mine"`
	// Path to show requests for. If omitted, shows all paths. Send `/` to represent the root directory.
	Path string `url:"path,omitempty" json:"path,omitempty" path:"path"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// RequestGetFolderParams contains the request parameters for GET /requests/folders/{path}.
type RequestGetFolderParams struct {
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are .
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// Only show requests of the current user? (Defaults to true if current user is not a site admin.).
	Mine *bool `url:"mine,omitempty" json:"mine,omitempty" path:"mine"`
	// Path to show requests for. If omitted, shows all paths. Send `/` to represent the root directory. Required; sent in the request path.
	Path string `url:"-,omitempty" json:"-,omitempty" path:"path"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// RequestCreateParams contains the request parameters for POST /requests.
type RequestCreateParams struct {
	// Folder path on which to request the file. Required.
	Path string `url:"path" json:"path" path:"path"`
	// Destination filename (without extension) to request. Required.
	Destination string `url:"destination" json:"destination" path:"destination"`
	// A list of user IDs to request the file from. If sent as a string, it should be comma-delimited.
	UserIds string `url:"user_ids,omitempty" json:"user_ids,omitempty" path:"user_ids"`
	// A list of group IDs to request the file from. If sent as a string, it should be comma-delimited.
	GroupIds string `url:"group_ids,omitempty" json:"group_ids,omitempty" path:"group_ids"`
}

// RequestDeleteParams contains the request parameters for DELETE /requests/{id}.
type RequestDeleteParams struct {
	// Request ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (r *Request) UnmarshalJSON(data []byte) error {
	type request Request
	var v request
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*r = Request(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (r *RequestCollection) UnmarshalJSON(data []byte) error {
	type requests RequestCollection
	var v requests
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*r = RequestCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (r *RequestCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*r))
	for i, v := range *r {
		ret[i] = v
	}

	return &ret
}
