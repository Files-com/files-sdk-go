# HistoryExport

`files_sdk.HistoryExport` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/historyexport`, named `history_export`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	history_export "github.com/Files-com/files-sdk-go/v3/historyexport"
)
```

## Operations

### Find

`GET /history_exports/{id}`

Show History Export

```go
func (c *Client) Find(params files_sdk.HistoryExportFindParams, opts ...files_sdk.RequestResponseOption) (historyExport files_sdk.HistoryExport, err error)
func Find(params files_sdk.HistoryExportFindParams, opts ...files_sdk.RequestResponseOption) (historyExport files_sdk.HistoryExport, err error)
```

Parameters, `files_sdk.HistoryExportFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | History Export ID. |

Example:

```go
result, err := history_export.Find(files_sdk.HistoryExportFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Create

`POST /history_exports`

Create History Export

```go
func (c *Client) Create(params files_sdk.HistoryExportCreateParams, opts ...files_sdk.RequestResponseOption) (historyExport files_sdk.HistoryExport, err error)
func Create(params files_sdk.HistoryExportCreateParams, opts ...files_sdk.RequestResponseOption) (historyExport files_sdk.HistoryExport, err error)
```

Parameters, `files_sdk.HistoryExportCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `UserId` | `int64` | `user_id` | no | User ID. Provide a value of `0` to operate the current session's user. |
| `StartAt` | `*time.Time` | `start_at` | no | Start date/time of export range. |
| `EndAt` | `*time.Time` | `end_at` | no | End date/time of export range. |
| `QueryAction` | `string` | `query_action` | no | Filter results by this this action type. Valid values: `create`, `read`, `update`, `destroy`, `move`, `login`, `failedlogin`, `copy`, `user_create`, `user_update`, `user_destroy`, `group_create`, `group_update`, `group_destroy`, `permission_create`, `permission_destroy`, `api_key_create`, `api_key_update`, `api_key_destroy`, `archived_delete` |
| `QueryInterface` | `string` | `query_interface` | no | Filter results by this this interface type. Valid values: `web`, `ftp`, `robot`, `jsapi`, `webdesktopapi`, `sftp`, `dav`, `desktop`, `restapi`, `scim`, `office`, `mobile`, `as2`, `inbound_email`, `remote`, `inbound_s3` |
| `QueryUserId` | `string` | `query_user_id` | no | Return results that are actions performed by the user indicated by this User ID |
| `QueryFileId` | `string` | `query_file_id` | no | Return results that are file actions related to the file indicated by this File ID |
| `QueryParentId` | `string` | `query_parent_id` | no | Return results that are file actions inside the parent folder specified by this folder ID |
| `QueryPath` | `string` | `query_path` | no | Return results that are file actions related to paths matching this pattern. |
| `QueryFolder` | `string` | `query_folder` | no | Return results that are file actions related to files or folders inside folder paths matching this pattern. |
| `QuerySrc` | `string` | `query_src` | no | Return results that are file moves originating from paths matching this pattern. |
| `QueryDestination` | `string` | `query_destination` | no | Return results that are file moves with paths matching this pattern as destination. |
| `QueryIp` | `string` | `query_ip` | no | Filter results by this IP address. |
| `QueryUsername` | `string` | `query_username` | no | Filter results by this username. |
| `QueryFailureType` | `string` | `query_failure_type` | no | If searching for Histories about login failures, this parameter restricts results to failures of this specific type. Valid values: `expired_trial`, `account_overdue`, `locked_out`, `ip_mismatch`, `password_mismatch`, `site_mismatch`, `username_not_found`, `none`, `no_ftp_permission`, `no_web_permission`, `no_directory`, `errno_enoent`, `no_sftp_permission`, `no_dav_permission`, `no_restapi_permission`, `key_mismatch`, `region_mismatch`, `expired_access`, `desktop_ip_mismatch`, `desktop_api_key_not_used_quickly_enough`, `disabled`, `country_mismatch`, `insecure_ftp`, `insecure_cipher`, `rate_limited`, `no_s3_compatible_endpoint_permission`, `mobile_ip_mismatch`, `mobile_api_key_not_used_quickly_enough`, `desktop_app_disabled`, `mobile_app_disabled` |
| `QueryTargetId` | `string` | `query_target_id` | no | If searching for Histories about specific objects (such as Users, or API Keys), this parameter restricts results to objects that match this ID. |
| `QueryTargetName` | `string` | `query_target_name` | no | If searching for Histories about Users, Groups or other objects with names, this parameter restricts results to objects with this name/username. |
| `QueryTargetPermission` | `string` | `query_target_permission` | no | If searching for Histories about Permissions, this parameter restricts results to permissions of this level. |
| `QueryTargetUserId` | `string` | `query_target_user_id` | no | If searching for Histories about API keys, this parameter restricts results to API keys created by/for this user ID. |
| `QueryTargetUsername` | `string` | `query_target_username` | no | If searching for Histories about API keys, this parameter restricts results to API keys created by/for this username. |
| `QueryTargetPlatform` | `string` | `query_target_platform` | no | If searching for Histories about API keys, this parameter restricts results to API keys associated with this platform. |
| `QueryTargetPermissionSet` | `string` | `query_target_permission_set` | no | If searching for Histories about API keys, this parameter restricts results to API keys with this permission set. |

Example:

```go
result, err := history_export.Create(files_sdk.HistoryExportCreateParams{})
if err != nil {
	return err
}
fmt.Println(result)
```

## Model

`files_sdk.HistoryExport` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `WorkspaceId` | `int64` | `workspace_id` | Workspace of the export. 0 represents the default workspace. A null value means a site-wide export. |
| `Id` | `int64` | `id` | History Export ID |
| `HistoryVersion` | `string` | `history_version` | Version of the history for the export. |
| `StartAt` | `*time.Time` | `start_at` | Start date/time of export range. |
| `EndAt` | `*time.Time` | `end_at` | End date/time of export range. |
| `Status` | `string` | `status` | Status of export. Will be: `building`, `ready`, or `failed` |
| `QueryAction` | `string` | `query_action` | Filter results by this this action type. Valid values: `create`, `read`, `update`, `destroy`, `move`, `login`, `failedlogin`, `copy`, `user_create`, `user_update`, `user_destroy`, `group_create`, `group_update`, `group_destroy`, `permission_create`, `permission_destroy`, `api_key_create`, `api_key_update`, `api_key_destroy`, `archived_delete` |
| `QueryInterface` | `string` | `query_interface` | Filter results by this this interface type. Valid values: `web`, `ftp`, `robot`, `jsapi`, `webdesktopapi`, `sftp`, `dav`, `desktop`, `restapi`, `scim`, `office`, `mobile`, `as2`, `inbound_email`, `remote`, `inbound_s3` |
| `QueryUserId` | `string` | `query_user_id` | Return results that are actions performed by the user indicated by this User ID |
| `QueryFileId` | `string` | `query_file_id` | Return results that are file actions related to the file indicated by this File ID |
| `QueryParentId` | `string` | `query_parent_id` | Return results that are file actions inside the parent folder specified by this folder ID |
| `QueryPath` | `string` | `query_path` | Return results that are file actions related to paths matching this pattern. |
| `QueryFolder` | `string` | `query_folder` | Return results that are file actions related to files or folders inside folder paths matching this pattern. |
| `QuerySrc` | `string` | `query_src` | Return results that are file moves originating from paths matching this pattern. |
| `QueryDestination` | `string` | `query_destination` | Return results that are file moves with paths matching this pattern as destination. |
| `QueryIp` | `string` | `query_ip` | Filter results by this IP address. |
| `QueryUsername` | `string` | `query_username` | Filter results by this username. |
| `QueryFailureType` | `string` | `query_failure_type` | If searching for Histories about login failures, this parameter restricts results to failures of this specific type. Valid values: `expired_trial`, `account_overdue`, `locked_out`, `ip_mismatch`, `password_mismatch`, `site_mismatch`, `username_not_found`, `none`, `no_ftp_permission`, `no_web_permission`, `no_directory`, `errno_enoent`, `no_sftp_permission`, `no_dav_permission`, `no_restapi_permission`, `key_mismatch`, `region_mismatch`, `expired_access`, `desktop_ip_mismatch`, `desktop_api_key_not_used_quickly_enough`, `disabled`, `country_mismatch`, `insecure_ftp`, `insecure_cipher`, `rate_limited`, `no_s3_compatible_endpoint_permission`, `mobile_ip_mismatch`, `mobile_api_key_not_used_quickly_enough`, `desktop_app_disabled`, `mobile_app_disabled` |
| `QueryTargetId` | `string` | `query_target_id` | If searching for Histories about specific objects (such as Users, or API Keys), this parameter restricts results to objects that match this ID. |
| `QueryTargetName` | `string` | `query_target_name` | If searching for Histories about Users, Groups or other objects with names, this parameter restricts results to objects with this name/username. |
| `QueryTargetPermission` | `string` | `query_target_permission` | If searching for Histories about Permissions, this parameter restricts results to permissions of this level. |
| `QueryTargetUserId` | `string` | `query_target_user_id` | If searching for Histories about API keys, this parameter restricts results to API keys created by/for this user ID. |
| `QueryTargetUsername` | `string` | `query_target_username` | If searching for Histories about API keys, this parameter restricts results to API keys created by/for this username. |
| `QueryTargetPlatform` | `string` | `query_target_platform` | If searching for Histories about API keys, this parameter restricts results to API keys associated with this platform. |
| `QueryTargetPermissionSet` | `string` | `query_target_permission_set` | If searching for Histories about API keys, this parameter restricts results to API keys with this permission set. |
| `ResultsUrl` | `string` | `results_url` | If `status` is `ready`, this will be a URL where all the results can be downloaded at once as a CSV. |
| `UserId` | `int64` | `user_id` | User ID. Provide a value of `0` to operate the current session's user. |
