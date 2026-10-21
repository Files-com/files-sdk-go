package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// HistoryExport is a Files.com API resource.
type HistoryExport struct {
	// Workspace of the export. 0 represents the default workspace. A null value means a site-wide export.
	WorkspaceId int64 `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	// History Export ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Version of the history for the export.
	HistoryVersion string `json:"history_version,omitempty" path:"history_version,omitempty" url:"history_version,omitempty"`
	// Start date/time of export range.
	StartAt *time.Time `json:"start_at,omitempty" path:"start_at,omitempty" url:"start_at,omitempty"`
	// End date/time of export range.
	EndAt *time.Time `json:"end_at,omitempty" path:"end_at,omitempty" url:"end_at,omitempty"`
	// Status of export. Will be: `building`, `ready`, or `failed`.
	Status string `json:"status,omitempty" path:"status,omitempty" url:"status,omitempty"`
	// Filter results by this this action type. Valid values: `create`, `read`, `update`, `destroy`, `move`, `login`, `failedlogin`, `copy`, `user_create`, `user_update`, `user_destroy`, `group_create`, `group_update`, `group_destroy`, `permission_create`, `permission_destroy`, `api_key_create`, `api_key_update`, `api_key_destroy`, `archived_delete`.
	QueryAction string `json:"query_action,omitempty" path:"query_action,omitempty" url:"query_action,omitempty"`
	// Filter results by this this interface type. Valid values: `web`, `ftp`, `robot`, `jsapi`, `webdesktopapi`, `sftp`, `dav`, `desktop`, `restapi`, `scim`, `office`, `mobile`, `as2`, `inbound_email`, `remote`, `inbound_s3`.
	QueryInterface string `json:"query_interface,omitempty" path:"query_interface,omitempty" url:"query_interface,omitempty"`
	// Return results that are actions performed by the user indicated by this User ID.
	QueryUserId string `json:"query_user_id,omitempty" path:"query_user_id,omitempty" url:"query_user_id,omitempty"`
	// Return results that are file actions related to the file indicated by this File ID.
	QueryFileId string `json:"query_file_id,omitempty" path:"query_file_id,omitempty" url:"query_file_id,omitempty"`
	// Return results that are file actions inside the parent folder specified by this folder ID.
	QueryParentId string `json:"query_parent_id,omitempty" path:"query_parent_id,omitempty" url:"query_parent_id,omitempty"`
	// Return results that are file actions related to paths matching this pattern.
	QueryPath string `json:"query_path,omitempty" path:"query_path,omitempty" url:"query_path,omitempty"`
	// Return results that are file actions related to files or folders inside folder paths matching this pattern.
	QueryFolder string `json:"query_folder,omitempty" path:"query_folder,omitempty" url:"query_folder,omitempty"`
	// Return results that are file moves originating from paths matching this pattern.
	QuerySrc string `json:"query_src,omitempty" path:"query_src,omitempty" url:"query_src,omitempty"`
	// Return results that are file moves with paths matching this pattern as destination.
	QueryDestination string `json:"query_destination,omitempty" path:"query_destination,omitempty" url:"query_destination,omitempty"`
	// Filter results by this IP address.
	QueryIp string `json:"query_ip,omitempty" path:"query_ip,omitempty" url:"query_ip,omitempty"`
	// Filter results by this username.
	QueryUsername string `json:"query_username,omitempty" path:"query_username,omitempty" url:"query_username,omitempty"`
	// If searching for Histories about login failures, this parameter restricts results to failures of this specific type. Valid values: `expired_trial`, `account_overdue`, `locked_out`, `ip_mismatch`, `password_mismatch`, `site_mismatch`, `username_not_found`, `none`, `no_ftp_permission`, `no_web_permission`, `no_directory`, `errno_enoent`, `no_sftp_permission`, `no_dav_permission`, `no_restapi_permission`, `key_mismatch`, `region_mismatch`, `expired_access`, `desktop_ip_mismatch`, `desktop_api_key_not_used_quickly_enough`, `disabled`, `country_mismatch`, `insecure_ftp`, `insecure_cipher`, `rate_limited`, `no_s3_compatible_endpoint_permission`, `mobile_ip_mismatch`, `mobile_api_key_not_used_quickly_enough`, `desktop_app_disabled`, `mobile_app_disabled`.
	QueryFailureType string `json:"query_failure_type,omitempty" path:"query_failure_type,omitempty" url:"query_failure_type,omitempty"`
	// If searching for Histories about specific objects (such as Users, or API Keys), this parameter restricts results to objects that match this ID.
	QueryTargetId string `json:"query_target_id,omitempty" path:"query_target_id,omitempty" url:"query_target_id,omitempty"`
	// If searching for Histories about Users, Groups or other objects with names, this parameter restricts results to objects with this name/username.
	QueryTargetName string `json:"query_target_name,omitempty" path:"query_target_name,omitempty" url:"query_target_name,omitempty"`
	// If searching for Histories about Permissions, this parameter restricts results to permissions of this level.
	QueryTargetPermission string `json:"query_target_permission,omitempty" path:"query_target_permission,omitempty" url:"query_target_permission,omitempty"`
	// If searching for Histories about API keys, this parameter restricts results to API keys created by/for this user ID.
	QueryTargetUserId string `json:"query_target_user_id,omitempty" path:"query_target_user_id,omitempty" url:"query_target_user_id,omitempty"`
	// If searching for Histories about API keys, this parameter restricts results to API keys created by/for this username.
	QueryTargetUsername string `json:"query_target_username,omitempty" path:"query_target_username,omitempty" url:"query_target_username,omitempty"`
	// If searching for Histories about API keys, this parameter restricts results to API keys associated with this platform.
	QueryTargetPlatform string `json:"query_target_platform,omitempty" path:"query_target_platform,omitempty" url:"query_target_platform,omitempty"`
	// If searching for Histories about API keys, this parameter restricts results to API keys with this permission set.
	QueryTargetPermissionSet string `json:"query_target_permission_set,omitempty" path:"query_target_permission_set,omitempty" url:"query_target_permission_set,omitempty"`
	// If `status` is `ready`, this will be a URL where all the results can be downloaded at once as a CSV.
	ResultsUrl string `json:"results_url,omitempty" path:"results_url,omitempty" url:"results_url,omitempty"`
	// User ID. Provide a value of `0` to operate the current session's user.
	UserId int64 `json:"user_id,omitempty" path:"user_id,omitempty" url:"user_id,omitempty"`
}

// Identifier returns the resource ID.
func (h HistoryExport) Identifier() interface{} {
	return h.Id
}

// HistoryExportCollection is a list of HistoryExport resources.
type HistoryExportCollection []HistoryExport

// HistoryExportFindParams contains the request parameters for GET /history_exports/{id}.
type HistoryExportFindParams struct {
	// History Export ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// HistoryExportCreateParams contains the request parameters for POST /history_exports.
type HistoryExportCreateParams struct {
	// User ID. Provide a value of `0` to operate the current session's user.
	UserId int64 `url:"user_id,omitempty" json:"user_id,omitempty" path:"user_id"`
	// Start date/time of export range.
	StartAt *time.Time `url:"start_at,omitempty" json:"start_at,omitempty" path:"start_at"`
	// End date/time of export range.
	EndAt *time.Time `url:"end_at,omitempty" json:"end_at,omitempty" path:"end_at"`
	// Filter results by this this action type. Valid values: `create`, `read`, `update`, `destroy`, `move`, `login`, `failedlogin`, `copy`, `user_create`, `user_update`, `user_destroy`, `group_create`, `group_update`, `group_destroy`, `permission_create`, `permission_destroy`, `api_key_create`, `api_key_update`, `api_key_destroy`, `archived_delete`.
	QueryAction string `url:"query_action,omitempty" json:"query_action,omitempty" path:"query_action"`
	// Filter results by this this interface type. Valid values: `web`, `ftp`, `robot`, `jsapi`, `webdesktopapi`, `sftp`, `dav`, `desktop`, `restapi`, `scim`, `office`, `mobile`, `as2`, `inbound_email`, `remote`, `inbound_s3`.
	QueryInterface string `url:"query_interface,omitempty" json:"query_interface,omitempty" path:"query_interface"`
	// Return results that are actions performed by the user indicated by this User ID.
	QueryUserId string `url:"query_user_id,omitempty" json:"query_user_id,omitempty" path:"query_user_id"`
	// Return results that are file actions related to the file indicated by this File ID.
	QueryFileId string `url:"query_file_id,omitempty" json:"query_file_id,omitempty" path:"query_file_id"`
	// Return results that are file actions inside the parent folder specified by this folder ID.
	QueryParentId string `url:"query_parent_id,omitempty" json:"query_parent_id,omitempty" path:"query_parent_id"`
	// Return results that are file actions related to paths matching this pattern.
	QueryPath string `url:"query_path,omitempty" json:"query_path,omitempty" path:"query_path"`
	// Return results that are file actions related to files or folders inside folder paths matching this pattern.
	QueryFolder string `url:"query_folder,omitempty" json:"query_folder,omitempty" path:"query_folder"`
	// Return results that are file moves originating from paths matching this pattern.
	QuerySrc string `url:"query_src,omitempty" json:"query_src,omitempty" path:"query_src"`
	// Return results that are file moves with paths matching this pattern as destination.
	QueryDestination string `url:"query_destination,omitempty" json:"query_destination,omitempty" path:"query_destination"`
	// Filter results by this IP address.
	QueryIp string `url:"query_ip,omitempty" json:"query_ip,omitempty" path:"query_ip"`
	// Filter results by this username.
	QueryUsername string `url:"query_username,omitempty" json:"query_username,omitempty" path:"query_username"`
	// If searching for Histories about login failures, this parameter restricts results to failures of this specific type. Valid values: `expired_trial`, `account_overdue`, `locked_out`, `ip_mismatch`, `password_mismatch`, `site_mismatch`, `username_not_found`, `none`, `no_ftp_permission`, `no_web_permission`, `no_directory`, `errno_enoent`, `no_sftp_permission`, `no_dav_permission`, `no_restapi_permission`, `key_mismatch`, `region_mismatch`, `expired_access`, `desktop_ip_mismatch`, `desktop_api_key_not_used_quickly_enough`, `disabled`, `country_mismatch`, `insecure_ftp`, `insecure_cipher`, `rate_limited`, `no_s3_compatible_endpoint_permission`, `mobile_ip_mismatch`, `mobile_api_key_not_used_quickly_enough`, `desktop_app_disabled`, `mobile_app_disabled`.
	QueryFailureType string `url:"query_failure_type,omitempty" json:"query_failure_type,omitempty" path:"query_failure_type"`
	// If searching for Histories about specific objects (such as Users, or API Keys), this parameter restricts results to objects that match this ID.
	QueryTargetId string `url:"query_target_id,omitempty" json:"query_target_id,omitempty" path:"query_target_id"`
	// If searching for Histories about Users, Groups or other objects with names, this parameter restricts results to objects with this name/username.
	QueryTargetName string `url:"query_target_name,omitempty" json:"query_target_name,omitempty" path:"query_target_name"`
	// If searching for Histories about Permissions, this parameter restricts results to permissions of this level.
	QueryTargetPermission string `url:"query_target_permission,omitempty" json:"query_target_permission,omitempty" path:"query_target_permission"`
	// If searching for Histories about API keys, this parameter restricts results to API keys created by/for this user ID.
	QueryTargetUserId string `url:"query_target_user_id,omitempty" json:"query_target_user_id,omitempty" path:"query_target_user_id"`
	// If searching for Histories about API keys, this parameter restricts results to API keys created by/for this username.
	QueryTargetUsername string `url:"query_target_username,omitempty" json:"query_target_username,omitempty" path:"query_target_username"`
	// If searching for Histories about API keys, this parameter restricts results to API keys associated with this platform.
	QueryTargetPlatform string `url:"query_target_platform,omitempty" json:"query_target_platform,omitempty" path:"query_target_platform"`
	// If searching for Histories about API keys, this parameter restricts results to API keys with this permission set.
	QueryTargetPermissionSet string `url:"query_target_permission_set,omitempty" json:"query_target_permission_set,omitempty" path:"query_target_permission_set"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (h *HistoryExport) UnmarshalJSON(data []byte) error {
	type historyExport HistoryExport
	var v historyExport
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*h = HistoryExport(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (h *HistoryExportCollection) UnmarshalJSON(data []byte) error {
	type historyExports HistoryExportCollection
	var v historyExports
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*h = HistoryExportCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (h *HistoryExportCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*h))
	for i, v := range *h {
		ret[i] = v
	}

	return &ret
}
