package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// Workspace is a Files.com API resource.
type Workspace struct {
	// Workspace ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Workspace name.
	Name string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
}

// Identifier returns the resource ID.
func (w Workspace) Identifier() interface{} {
	return w.Id
}

// WorkspaceCollection is a list of Workspace resources.
type WorkspaceCollection []Workspace

// WorkspaceListParams contains the request parameters for GET /workspaces.
type WorkspaceListParams struct {
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `name`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `name`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// If set, return records where the specified field is prefixed by the supplied value. Valid fields are `name`.
	FilterPrefix interface{} `url:"filter_prefix,omitempty" json:"filter_prefix,omitempty" path:"filter_prefix"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// WorkspaceFindParams contains the request parameters for GET /workspaces/{id}.
type WorkspaceFindParams struct {
	// Workspace ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// WorkspaceCreateParams contains the request parameters for POST /workspaces.
type WorkspaceCreateParams struct {
	// Workspace name.
	Name string `url:"name,omitempty" json:"name,omitempty" path:"name"`
}

// WorkspaceUpdateParams contains the request parameters for PATCH /workspaces/{id}.
type WorkspaceUpdateParams struct {
	// Workspace ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// Workspace name.
	Name string `url:"name,omitempty" json:"name,omitempty" path:"name"`
}

// WorkspaceDeleteParams contains the request parameters for DELETE /workspaces/{id}.
type WorkspaceDeleteParams struct {
	// Workspace ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (w *Workspace) UnmarshalJSON(data []byte) error {
	type workspace Workspace
	var v workspace
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*w = Workspace(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (w *WorkspaceCollection) UnmarshalJSON(data []byte) error {
	type workspaces WorkspaceCollection
	var v workspaces
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*w = WorkspaceCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (w *WorkspaceCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*w))
	for i, v := range *w {
		ret[i] = v
	}

	return &ret
}
