package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// History is a Files.com API resource.
type History struct {
	// Action ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Path. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.
	Path string `json:"path,omitempty" path:"path,omitempty" url:"path,omitempty"`
	// Action occurrence date/time.
	When *time.Time `json:"when,omitempty" path:"when,omitempty" url:"when,omitempty"`
	// The destination path for this action, if applicable.
	Destination string `json:"destination,omitempty" path:"destination,omitempty" url:"destination,omitempty"`
	// Friendly displayed output.
	Display string `json:"display,omitempty" path:"display,omitempty" url:"display,omitempty"`
	// IP Address that performed this action.
	Ip string `json:"ip,omitempty" path:"ip,omitempty" url:"ip,omitempty"`
	// The source path for this action, if applicable.
	Source string `json:"source,omitempty" path:"source,omitempty" url:"source,omitempty"`
	// Targets.
	Targets interface{} `json:"targets,omitempty" path:"targets,omitempty" url:"targets,omitempty"`
	// User ID.
	UserId int64 `json:"user_id,omitempty" path:"user_id,omitempty" url:"user_id,omitempty"`
	// Username.
	Username string `json:"username,omitempty" path:"username,omitempty" url:"username,omitempty"`
	// true if this change was performed by a user on a parent site.
	UserIsFromParentSite *bool `json:"user_is_from_parent_site,omitempty" path:"user_is_from_parent_site,omitempty" url:"user_is_from_parent_site,omitempty"`
	// Type of action.
	Action string `json:"action,omitempty" path:"action,omitempty" url:"action,omitempty"`
	// Failure type. If action was a user login or session failure, why did it fail?
	FailureType string `json:"failure_type,omitempty" path:"failure_type,omitempty" url:"failure_type,omitempty"`
	// Interface on which this action occurred.
	Interface string `json:"interface,omitempty" path:"interface,omitempty" url:"interface,omitempty"`
}

// Identifier returns the resource ID.
func (h History) Identifier() interface{} {
	return h.Id
}

// HistoryCollection is a list of History resources.
type HistoryCollection []History

// HistoryListForFileParams contains the request parameters for GET /history/files/{path}.
type HistoryListForFileParams struct {
	// Leave blank or set to a date/time to filter earlier entries.
	StartAt *time.Time `url:"start_at,omitempty" json:"start_at,omitempty" path:"start_at"`
	// Leave blank or set to a date/time to filter later entries.
	EndAt *time.Time `url:"end_at,omitempty" json:"end_at,omitempty" path:"end_at"`
	// Display format. Leave blank or set to `full` or `parent`.
	Display string `url:"display,omitempty" json:"display,omitempty" path:"display"`
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `created_at`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// Path to operate on. Required; sent in the request path.
	Path string `url:"-,omitempty" json:"-,omitempty" path:"path"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// HistoryListForFolderParams contains the request parameters for GET /history/folders/{path}.
type HistoryListForFolderParams struct {
	// Leave blank or set to a date/time to filter earlier entries.
	StartAt *time.Time `url:"start_at,omitempty" json:"start_at,omitempty" path:"start_at"`
	// Leave blank or set to a date/time to filter later entries.
	EndAt *time.Time `url:"end_at,omitempty" json:"end_at,omitempty" path:"end_at"`
	// Display format. Leave blank or set to `full` or `parent`.
	Display string `url:"display,omitempty" json:"display,omitempty" path:"display"`
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `created_at`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// Path to operate on. Required; sent in the request path.
	Path string `url:"-,omitempty" json:"-,omitempty" path:"path"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// HistoryListForUserParams contains the request parameters for GET /history/users/{user_id}.
type HistoryListForUserParams struct {
	// Leave blank or set to a date/time to filter earlier entries.
	StartAt *time.Time `url:"start_at,omitempty" json:"start_at,omitempty" path:"start_at"`
	// Leave blank or set to a date/time to filter later entries.
	EndAt *time.Time `url:"end_at,omitempty" json:"end_at,omitempty" path:"end_at"`
	// Display format. Leave blank or set to `full` or `parent`.
	Display string `url:"display,omitempty" json:"display,omitempty" path:"display"`
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `created_at`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// User ID. Required; sent in the request path.
	UserId int64 `url:"-,omitempty" json:"-,omitempty" path:"user_id"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// HistoryListLoginsParams contains the request parameters for GET /history/login.
type HistoryListLoginsParams struct {
	// Leave blank or set to a date/time to filter earlier entries.
	StartAt *time.Time `url:"start_at,omitempty" json:"start_at,omitempty" path:"start_at"`
	// Leave blank or set to a date/time to filter later entries.
	EndAt *time.Time `url:"end_at,omitempty" json:"end_at,omitempty" path:"end_at"`
	// Display format. Leave blank or set to `full` or `parent`.
	Display string `url:"display,omitempty" json:"display,omitempty" path:"display"`
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `created_at`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// HistoryListParams contains the request parameters for GET /history.
type HistoryListParams struct {
	// Leave blank or set to a date/time to filter earlier entries.
	StartAt *time.Time `url:"start_at,omitempty" json:"start_at,omitempty" path:"start_at"`
	// Leave blank or set to a date/time to filter later entries.
	EndAt *time.Time `url:"end_at,omitempty" json:"end_at,omitempty" path:"end_at"`
	// Display format. Leave blank or set to `full` or `parent`.
	Display string `url:"display,omitempty" json:"display,omitempty" path:"display"`
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `created_at`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `user_id`, `folder` or `path`. Valid field combinations are `[ ]`, `[ path ]`, `[ path ]` or `[ path ]`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// If set, return records where the specified field is prefixed by the supplied value. Valid fields are `path`.
	FilterPrefix interface{} `url:"filter_prefix,omitempty" json:"filter_prefix,omitempty" path:"filter_prefix"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (h *History) UnmarshalJSON(data []byte) error {
	type history History
	var v history
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*h = History(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (h *HistoryCollection) UnmarshalJSON(data []byte) error {
	type historys HistoryCollection
	var v historys
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*h = HistoryCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (h *HistoryCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*h))
	for i, v := range *h {
		ret[i] = v
	}

	return &ret
}
