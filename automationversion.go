package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// AutomationVersion is a Files.com API resource.
type AutomationVersion struct {
	// Automation revision ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// ID of the Automation that owns this revision.
	AutomationId int64 `json:"automation_id,omitempty" path:"automation_id,omitempty" url:"automation_id,omitempty"`
	// Workspace that owns the Automation. Zero is the Default Workspace.
	WorkspaceId int64 `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	// Revision number within the Automation.
	Version int64 `json:"version,omitempty" path:"version,omitempty" url:"version,omitempty"`
	// Time the revision was saved.
	CreatedAt *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	// ID of the user who saved the revision, when available.
	UserId int64 `json:"user_id,omitempty" path:"user_id,omitempty" url:"user_id,omitempty"`
	// Username recorded when the revision was saved, when available.
	Username string `json:"username,omitempty" path:"username,omitempty" url:"username,omitempty"`
	// ID of the API key used to save the revision, when available.
	ApiKeyId int64 `json:"api_key_id,omitempty" path:"api_key_id,omitempty" url:"api_key_id,omitempty"`
	// Saved Automation v2 graph definition. Included when retrieving a single revision.
	Definition interface{} `json:"definition,omitempty" path:"definition,omitempty" url:"definition,omitempty"`
}

// Identifier returns the resource ID.
func (a AutomationVersion) Identifier() interface{} {
	return a.Id
}

// AutomationVersionCollection is a list of AutomationVersion resources.
type AutomationVersionCollection []AutomationVersion

// AutomationVersionListParams contains the request parameters for GET /automation_versions.
type AutomationVersionListParams struct {
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are .
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// ID of the Automation that owns this revision. Required.
	AutomationId int64 `url:"automation_id" json:"automation_id" path:"automation_id"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// AutomationVersionFindParams contains the request parameters for GET /automation_versions/{id}.
type AutomationVersionFindParams struct {
	// Automation Version ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// AutomationVersionRestoreParams contains the request parameters for POST /automation_versions/{id}/restore.
//
// Restore an Automation v2 definition as a new current version
type AutomationVersionRestoreParams struct {
	// Automation Version ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (a *AutomationVersion) UnmarshalJSON(data []byte) error {
	type automationVersion AutomationVersion
	var v automationVersion
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*a = AutomationVersion(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (a *AutomationVersionCollection) UnmarshalJSON(data []byte) error {
	type automationVersions AutomationVersionCollection
	var v automationVersions
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*a = AutomationVersionCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (a *AutomationVersionCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*a))
	for i, v := range *a {
		ret[i] = v
	}

	return &ret
}
