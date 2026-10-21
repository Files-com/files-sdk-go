package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// Restore is a Files.com API resource.
type Restore struct {
	// Restore files or users deleted on or after this date/time. Don't set this earlier than you need. Can not be greater than 365 days prior to the restore request.
	EarliestDate *time.Time `json:"earliest_date,omitempty" path:"earliest_date,omitempty" url:"earliest_date,omitempty"`
	// Restore Record ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Number of directories that were successfully restored.
	DirsRestored int64 `json:"dirs_restored,omitempty" path:"dirs_restored,omitempty" url:"dirs_restored,omitempty"`
	// Number of directories that were not able to be restored.
	DirsErrored int64 `json:"dirs_errored,omitempty" path:"dirs_errored,omitempty" url:"dirs_errored,omitempty"`
	// Total number of directories processed.
	DirsTotal int64 `json:"dirs_total,omitempty" path:"dirs_total,omitempty" url:"dirs_total,omitempty"`
	// Number of files successfully restored.
	FilesRestored int64 `json:"files_restored,omitempty" path:"files_restored,omitempty" url:"files_restored,omitempty"`
	// Number of files that were not able to be restored.
	FilesErrored int64 `json:"files_errored,omitempty" path:"files_errored,omitempty" url:"files_errored,omitempty"`
	// Total number of files processed.
	FilesTotal int64 `json:"files_total,omitempty" path:"files_total,omitempty" url:"files_total,omitempty"`
	// Prefix of the files/folders to restore, or a case-insensitive username prefix for a user restore. A trailing slash selects that folder and its contents only; without it, the prefix matches any path beginning with that text. Do not use a leading slash. To restore all deleted items of the selected restoration type within the selected site or workspace scope, specify an empty string (`''`) or omit the field.
	Prefix string `json:"prefix,omitempty" path:"prefix,omitempty" url:"prefix,omitempty"`
	// Type of restoration to perform. `files` restores deleted filesystem items. `users` restores deleted users and associated authentication records, plus permissions when `restore_deleted_permissions=true`. A 5-minute grace period allows for batch deletion processes: associated records deleted within the 5 minutes before the user are included, even if removed independently. Records deleted earlier, or before `earliest_date`, remain deleted.
	RestorationType string `json:"restoration_type,omitempty" path:"restoration_type,omitempty" url:"restoration_type,omitempty"`
	// If true, we will restore the files in place (into their original paths). If false, we will create a new restoration folder in the root and restore files there.
	RestoreInPlace *bool `json:"restore_in_place,omitempty" path:"restore_in_place,omitempty" url:"restore_in_place,omitempty"`
	// If true, a user restore restores permissions deleted along with the selected users, including permissions removed independently within the 5 minutes before user deletion. Permissions deleted earlier, or before `earliest_date`, remain deleted. An in-place file restore restores user and group permissions on returned folders and their subfolders that were deleted from the time of the folder's deletion through 24 hours afterward, including permissions removed independently within that window. File restores into a new folder do not restore permissions.
	RestoreDeletedPermissions *bool `json:"restore_deleted_permissions,omitempty" path:"restore_deleted_permissions,omitempty" url:"restore_deleted_permissions,omitempty"`
	// Number of users successfully restored (only present for `restoration_type=users`).
	UsersRestored int64 `json:"users_restored,omitempty" path:"users_restored,omitempty" url:"users_restored,omitempty"`
	// Number of users that failed to restore (only present for `restoration_type=users`).
	UsersErrored int64 `json:"users_errored,omitempty" path:"users_errored,omitempty" url:"users_errored,omitempty"`
	// Total number of users processed (only present for `restoration_type=users`).
	UsersTotal int64 `json:"users_total,omitempty" path:"users_total,omitempty" url:"users_total,omitempty"`
	// Number of API keys restored (only present for `restoration_type=users`).
	ApiKeysRestored int64 `json:"api_keys_restored,omitempty" path:"api_keys_restored,omitempty" url:"api_keys_restored,omitempty"`
	// Number of public keys restored (only present for `restoration_type=users`).
	PublicKeysRestored int64 `json:"public_keys_restored,omitempty" path:"public_keys_restored,omitempty" url:"public_keys_restored,omitempty"`
	// Number of two factor authentication methods restored (only present for `restoration_type=users`).
	TwoFactorAuthenticationMethodsRestored int64 `json:"two_factor_authentication_methods_restored,omitempty" path:"two_factor_authentication_methods_restored,omitempty" url:"two_factor_authentication_methods_restored,omitempty"`
	// Status of the restoration process.
	Status string `json:"status,omitempty" path:"status,omitempty" url:"status,omitempty"`
	// If true, we will update the last modified timestamp of restored files to today's date. If false, we might trigger File Expiration to delete the file again.
	UpdateTimestamps *bool `json:"update_timestamps,omitempty" path:"update_timestamps,omitempty" url:"update_timestamps,omitempty"`
	// Workspace ID for a workspace-scoped restore. `0` means the default site-wide scope.
	WorkspaceId int64 `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	// Error messages received while restoring files and/or directories. Only present if there were errors.
	ErrorMessages []string `json:"error_messages,omitempty" path:"error_messages,omitempty" url:"error_messages,omitempty"`
}

// Identifier returns the resource ID.
func (r Restore) Identifier() interface{} {
	return r.Id
}

// RestoreCollection is a list of Restore resources.
type RestoreCollection []Restore

// RestoreRestorationTypeEnum is a string value for restoration_type.
// Enum lists the values documented by the API.
type RestoreRestorationTypeEnum string

// String returns the API parameter value.
func (u RestoreRestorationTypeEnum) String() string {
	return string(u)
}

// Enum returns the documented values keyed by their API strings.
func (u RestoreRestorationTypeEnum) Enum() map[string]RestoreRestorationTypeEnum {
	return map[string]RestoreRestorationTypeEnum{
		"files": RestoreRestorationTypeEnum("files"),
		"users": RestoreRestorationTypeEnum("users"),
	}
}

// RestoreListParams contains the request parameters for GET /restores.
type RestoreListParams struct {
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `restoration_type`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// RestoreCreateParams contains the request parameters for POST /restores.
type RestoreCreateParams struct {
	// Restore files or users deleted on or after this date/time. Don't set this earlier than you need. Can not be greater than 365 days prior to the restore request. Required.
	EarliestDate *time.Time `url:"earliest_date" json:"earliest_date" path:"earliest_date"`
	// Prefix of the files/folders to restore, or a case-insensitive username prefix for a user restore. A trailing slash selects that folder and its contents only; without it, the prefix matches any path beginning with that text. Do not use a leading slash. To restore all deleted items of the selected restoration type within the selected site or workspace scope, specify an empty string (`''`) or omit the field.
	Prefix string `url:"prefix,omitempty" json:"prefix,omitempty" path:"prefix"`
	// Type of restoration to perform. `files` restores deleted filesystem items. `users` restores deleted users and associated authentication records, plus permissions when `restore_deleted_permissions=true`. A 5-minute grace period allows for batch deletion processes: associated records deleted within the 5 minutes before the user are included, even if removed independently. Records deleted earlier, or before `earliest_date`, remain deleted. Values: files, users.
	RestorationType RestoreRestorationTypeEnum `url:"restoration_type,omitempty" json:"restoration_type,omitempty" path:"restoration_type"`
	// If true, a user restore restores permissions deleted along with the selected users, including permissions removed independently within the 5 minutes before user deletion. Permissions deleted earlier, or before `earliest_date`, remain deleted. An in-place file restore restores user and group permissions on returned folders and their subfolders that were deleted from the time of the folder's deletion through 24 hours afterward, including permissions removed independently within that window. File restores into a new folder do not restore permissions.
	RestoreDeletedPermissions *bool `url:"restore_deleted_permissions,omitempty" json:"restore_deleted_permissions,omitempty" path:"restore_deleted_permissions"`
	// If true, we will restore the files in place (into their original paths). If false, we will create a new restoration folder in the root and restore files there.
	RestoreInPlace *bool `url:"restore_in_place,omitempty" json:"restore_in_place,omitempty" path:"restore_in_place"`
	// If true, we will update the last modified timestamp of restored files to today's date. If false, we might trigger File Expiration to delete the file again.
	UpdateTimestamps *bool `url:"update_timestamps,omitempty" json:"update_timestamps,omitempty" path:"update_timestamps"`
	// Workspace ID for a workspace-scoped restore. `0` means the default site-wide scope.
	WorkspaceId int64 `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (r *Restore) UnmarshalJSON(data []byte) error {
	type restore Restore
	var v restore
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*r = Restore(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (r *RestoreCollection) UnmarshalJSON(data []byte) error {
	type restores RestoreCollection
	var v restores
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*r = RestoreCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (r *RestoreCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*r))
	for i, v := range *r {
		ret[i] = v
	}

	return &ret
}
