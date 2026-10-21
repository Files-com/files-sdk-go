package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// HistoryExportResult is a Files.com API resource.
type HistoryExportResult struct {
	// Action ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// When the action happened.
	CreatedAt int64 `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	// When the action happened, in ISO8601 format.
	CreatedAtIso8601 string `json:"created_at_iso8601,omitempty" path:"created_at_iso8601,omitempty" url:"created_at_iso8601,omitempty"`
	// User ID.
	UserId int64 `json:"user_id,omitempty" path:"user_id,omitempty" url:"user_id,omitempty"`
	// File ID related to the action.
	FileId int64 `json:"file_id,omitempty" path:"file_id,omitempty" url:"file_id,omitempty"`
	// ID of the parent folder.
	ParentId int64 `json:"parent_id,omitempty" path:"parent_id,omitempty" url:"parent_id,omitempty"`
	// Path of the related action. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.
	Path string `json:"path,omitempty" path:"path,omitempty" url:"path,omitempty"`
	// Folder in which the action occurred.
	Folder string `json:"folder,omitempty" path:"folder,omitempty" url:"folder,omitempty"`
	// File move originated from this path.
	Src string `json:"src,omitempty" path:"src,omitempty" url:"src,omitempty"`
	// File moved to this destination folder.
	Destination string `json:"destination,omitempty" path:"destination,omitempty" url:"destination,omitempty"`
	// Client IP that performed the action.
	Ip string `json:"ip,omitempty" path:"ip,omitempty" url:"ip,omitempty"`
	// Username of the user that performed the action.
	Username string `json:"username,omitempty" path:"username,omitempty" url:"username,omitempty"`
	// true if this change was performed by a user on a parent site.
	UserIsFromParentSite *bool `json:"user_is_from_parent_site,omitempty" path:"user_is_from_parent_site,omitempty" url:"user_is_from_parent_site,omitempty"`
	// What action was taken. Valid values: `create`, `read`, `update`, `destroy`, `move`, `login`, `failedlogin`, `copy`, `user_create`, `user_update`, `user_destroy`, `group_create`, `group_update`, `group_destroy`, `permission_create`, `permission_destroy`, `api_key_create`, `api_key_update`, `api_key_destroy`, `archived_delete`.
	Action string `json:"action,omitempty" path:"action,omitempty" url:"action,omitempty"`
	// The type of login failure, if applicable. Valid values: `expired_trial`, `account_overdue`, `locked_out`, `ip_mismatch`, `password_mismatch`, `site_mismatch`, `username_not_found`, `none`, `no_ftp_permission`, `no_web_permission`, `no_directory`, `errno_enoent`, `no_sftp_permission`, `no_dav_permission`, `no_restapi_permission`, `key_mismatch`, `region_mismatch`, `expired_access`, `desktop_ip_mismatch`, `desktop_api_key_not_used_quickly_enough`, `disabled`, `country_mismatch`, `insecure_ftp`, `insecure_cipher`, `rate_limited`, `no_s3_compatible_endpoint_permission`, `mobile_ip_mismatch`, `mobile_api_key_not_used_quickly_enough`, `desktop_app_disabled`, `mobile_app_disabled`.
	FailureType string `json:"failure_type,omitempty" path:"failure_type,omitempty" url:"failure_type,omitempty"`
	// Interface through which the action was taken. Valid values: `web`, `ftp`, `robot`, `jsapi`, `webdesktopapi`, `sftp`, `dav`, `desktop`, `restapi`, `scim`, `office`, `mobile`, `as2`, `inbound_email`, `remote`, `inbound_s3`.
	Interface string `json:"interface,omitempty" path:"interface,omitempty" url:"interface,omitempty"`
	// ID of the object (such as Users, or API Keys) on which the action was taken.
	TargetId int64 `json:"target_id,omitempty" path:"target_id,omitempty" url:"target_id,omitempty"`
	// Name of the User, Group or other object with a name related to this action.
	TargetName string `json:"target_name,omitempty" path:"target_name,omitempty" url:"target_name,omitempty"`
	// Permission level of the action.
	TargetPermission string `json:"target_permission,omitempty" path:"target_permission,omitempty" url:"target_permission,omitempty"`
	// Whether or not the action was recursive.
	TargetRecursive *bool `json:"target_recursive,omitempty" path:"target_recursive,omitempty" url:"target_recursive,omitempty"`
	// If searching for Histories about API keys, this is when the API key will expire. Represented as a Unix timestamp.
	TargetExpiresAt int64 `json:"target_expires_at,omitempty" path:"target_expires_at,omitempty" url:"target_expires_at,omitempty"`
	// If searching for Histories about API keys, this is when the API key will expire. Represented in ISO8601 format.
	TargetExpiresAtIso8601 string `json:"target_expires_at_iso8601,omitempty" path:"target_expires_at_iso8601,omitempty" url:"target_expires_at_iso8601,omitempty"`
	// If searching for Histories about API keys, this represents the permission set of the associated API key.
	TargetPermissionSet string `json:"target_permission_set,omitempty" path:"target_permission_set,omitempty" url:"target_permission_set,omitempty"`
	// If searching for Histories about API keys, this is the platform on which the action was taken.
	TargetPlatform string `json:"target_platform,omitempty" path:"target_platform,omitempty" url:"target_platform,omitempty"`
	// If searching for Histories about API keys, this is the username on which the action was taken.
	TargetUsername string `json:"target_username,omitempty" path:"target_username,omitempty" url:"target_username,omitempty"`
	// If searching for Histories about API keys, this is the User ID on which the action was taken.
	TargetUserId int64 `json:"target_user_id,omitempty" path:"target_user_id,omitempty" url:"target_user_id,omitempty"`
}

// Identifier returns the resource ID.
func (h HistoryExportResult) Identifier() interface{} {
	return h.Id
}

// HistoryExportResultCollection is a list of HistoryExportResult resources.
type HistoryExportResultCollection []HistoryExportResult

// HistoryExportResultListParams contains the request parameters for GET /history_export_results.
type HistoryExportResultListParams struct {
	// User ID. Provide a value of `0` to operate the current session's user.
	UserId int64 `url:"user_id,omitempty" json:"user_id,omitempty" path:"user_id"`
	// ID of the associated history export. Required.
	HistoryExportId int64 `url:"history_export_id" json:"history_export_id" path:"history_export_id"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (h *HistoryExportResult) UnmarshalJSON(data []byte) error {
	type historyExportResult HistoryExportResult
	var v historyExportResult
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*h = HistoryExportResult(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (h *HistoryExportResultCollection) UnmarshalJSON(data []byte) error {
	type historyExportResults HistoryExportResultCollection
	var v historyExportResults
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*h = HistoryExportResultCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (h *HistoryExportResultCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*h))
	for i, v := range *h {
		ret[i] = v
	}

	return &ret
}
