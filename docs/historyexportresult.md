# HistoryExportResult

`files_sdk.HistoryExportResult` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/historyexportresult`, named `history_export_result`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	history_export_result "github.com/Files-com/files-sdk-go/v3/historyexportresult"
)
```

## Operations

### List

`GET /history_export_results`

List History Export Results

```go
func (c *Client) List(params files_sdk.HistoryExportResultListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.HistoryExportResultListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.HistoryExportResult` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.HistoryExportResultListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `UserId` | `int64` | `user_id` | no | User ID. Provide a value of `0` to operate the current session's user. |
| `HistoryExportId` | `int64` | `history_export_id` | yes | ID of the associated history export. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := history_export_result.List(files_sdk.HistoryExportResultListParams{
	HistoryExportId: 1,
})
if err != nil {
	return err
}
for item, err := range it.All() {
	if err != nil {
		return err
	}
	fmt.Println(item)
}
```

## Model

`files_sdk.HistoryExportResult` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Action ID |
| `CreatedAt` | `int64` | `created_at` | When the action happened |
| `CreatedAtIso8601` | `string` | `created_at_iso8601` | When the action happened, in ISO8601 format. |
| `UserId` | `int64` | `user_id` | User ID |
| `FileId` | `int64` | `file_id` | File ID related to the action |
| `ParentId` | `int64` | `parent_id` | ID of the parent folder |
| `Path` | `string` | `path` | Path of the related action. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters. |
| `Folder` | `string` | `folder` | Folder in which the action occurred |
| `Src` | `string` | `src` | File move originated from this path |
| `Destination` | `string` | `destination` | File moved to this destination folder |
| `Ip` | `string` | `ip` | Client IP that performed the action |
| `Username` | `string` | `username` | Username of the user that performed the action |
| `UserIsFromParentSite` | `*bool` | `user_is_from_parent_site` | true if this change was performed by a user on a parent site. |
| `Action` | `string` | `action` | What action was taken. Valid values: `create`, `read`, `update`, `destroy`, `move`, `login`, `failedlogin`, `copy`, `user_create`, `user_update`, `user_destroy`, `group_create`, `group_update`, `group_destroy`, `permission_create`, `permission_destroy`, `api_key_create`, `api_key_update`, `api_key_destroy`, `archived_delete` |
| `FailureType` | `string` | `failure_type` | The type of login failure, if applicable. Valid values: `expired_trial`, `account_overdue`, `locked_out`, `ip_mismatch`, `password_mismatch`, `site_mismatch`, `username_not_found`, `none`, `no_ftp_permission`, `no_web_permission`, `no_directory`, `errno_enoent`, `no_sftp_permission`, `no_dav_permission`, `no_restapi_permission`, `key_mismatch`, `region_mismatch`, `expired_access`, `desktop_ip_mismatch`, `desktop_api_key_not_used_quickly_enough`, `disabled`, `country_mismatch`, `insecure_ftp`, `insecure_cipher`, `rate_limited`, `no_s3_compatible_endpoint_permission`, `mobile_ip_mismatch`, `mobile_api_key_not_used_quickly_enough`, `desktop_app_disabled`, `mobile_app_disabled` |
| `Interface` | `string` | `interface` | Interface through which the action was taken. Valid values: `web`, `ftp`, `robot`, `jsapi`, `webdesktopapi`, `sftp`, `dav`, `desktop`, `restapi`, `scim`, `office`, `mobile`, `as2`, `inbound_email`, `remote`, `inbound_s3` |
| `TargetId` | `int64` | `target_id` | ID of the object (such as Users, or API Keys) on which the action was taken |
| `TargetName` | `string` | `target_name` | Name of the User, Group or other object with a name related to this action |
| `TargetPermission` | `string` | `target_permission` | Permission level of the action |
| `TargetRecursive` | `*bool` | `target_recursive` | Whether or not the action was recursive |
| `TargetExpiresAt` | `int64` | `target_expires_at` | If searching for Histories about API keys, this is when the API key will expire. Represented as a Unix timestamp. |
| `TargetExpiresAtIso8601` | `string` | `target_expires_at_iso8601` | If searching for Histories about API keys, this is when the API key will expire. Represented in ISO8601 format. |
| `TargetPermissionSet` | `string` | `target_permission_set` | If searching for Histories about API keys, this represents the permission set of the associated API key |
| `TargetPlatform` | `string` | `target_platform` | If searching for Histories about API keys, this is the platform on which the action was taken |
| `TargetUsername` | `string` | `target_username` | If searching for Histories about API keys, this is the username on which the action was taken |
| `TargetUserId` | `int64` | `target_user_id` | If searching for Histories about API keys, this is the User ID on which the action was taken |
