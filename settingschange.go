package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// SettingsChange is a Files.com API resource.
type SettingsChange struct {
	// The API key id responsible for this change.
	ApiKeyId int64 `json:"api_key_id,omitempty" path:"api_key_id,omitempty" url:"api_key_id,omitempty"`
	// Markdown-formatted change messages.
	Changes []string `json:"changes,omitempty" path:"changes,omitempty" url:"changes,omitempty"`
	// The time this change was made.
	CreatedAt *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	// The user id responsible for this change.
	UserId int64 `json:"user_id,omitempty" path:"user_id,omitempty" url:"user_id,omitempty"`
	// true if this change was performed by Files.com support.
	UserIsFilesSupport *bool `json:"user_is_files_support,omitempty" path:"user_is_files_support,omitempty" url:"user_is_files_support,omitempty"`
	// true if this change was performed by a user on a parent site.
	UserIsFromParentSite *bool `json:"user_is_from_parent_site,omitempty" path:"user_is_from_parent_site,omitempty" url:"user_is_from_parent_site,omitempty"`
	// The username of the user responsible for this change.
	Username string `json:"username,omitempty" path:"username,omitempty" url:"username,omitempty"`
}

// Identifier no path or id

// SettingsChangeCollection is a list of SettingsChange resources.
type SettingsChangeCollection []SettingsChange

// SettingsChangeListParams contains the request parameters for GET /settings_changes.
type SettingsChangeListParams struct {
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `created_at`, `api_key_id` or `user_id`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `api_key_id` and `user_id`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (s *SettingsChange) UnmarshalJSON(data []byte) error {
	type settingsChange SettingsChange
	var v settingsChange
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*s = SettingsChange(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (s *SettingsChangeCollection) UnmarshalJSON(data []byte) error {
	type settingsChanges SettingsChangeCollection
	var v settingsChanges
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*s = SettingsChangeCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (s *SettingsChangeCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*s))
	for i, v := range *s {
		ret[i] = v
	}

	return &ret
}
