# Restore

`files_sdk.Restore` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/restore`, named `restore`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/restore"
)
```

## Operations

### List

`GET /restores`

List Restores

```go
func (c *Client) List(params files_sdk.RestoreListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.RestoreListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.Restore` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.RestoreListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `restoration_type`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := restore.List(files_sdk.RestoreListParams{})
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

### Create

`POST /restores`

Create Restore

```go
func (c *Client) Create(params files_sdk.RestoreCreateParams, opts ...files_sdk.RequestResponseOption) (restore files_sdk.Restore, err error)
func Create(params files_sdk.RestoreCreateParams, opts ...files_sdk.RequestResponseOption) (restore files_sdk.Restore, err error)
```

Parameters, `files_sdk.RestoreCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `EarliestDate` | `*time.Time` | `earliest_date` | yes | Restore files or users deleted on or after this date/time. Don't set this earlier than you need. Can not be greater than 365 days prior to the restore request. |
| `Prefix` | `string` | `prefix` | no | Prefix of the files/folders to restore, or a case-insensitive username prefix for a user restore. A trailing slash selects that folder and its contents only; without it, the prefix matches any path beginning with that text. Do not use a leading slash. To restore all deleted items of the selected restoration type within the selected site or workspace scope, specify an empty string (`''`) or omit the field. |
| `RestorationType` | `files_sdk.RestoreRestorationTypeEnum` | `restoration_type` | no | Type of restoration to perform. `files` restores deleted filesystem items. `users` restores deleted users and associated authentication records, plus permissions when `restore_deleted_permissions=true`. A 5-minute grace period allows for batch deletion processes: associated records deleted within the 5 minutes before the user are included, even if removed independently. Records deleted earlier, or before `earliest_date`, remain deleted. Values: `files`, `users`. |
| `RestoreDeletedPermissions` | `*bool` | `restore_deleted_permissions` | no | If true, a user restore restores permissions deleted along with the selected users, including permissions removed independently within the 5 minutes before user deletion. Permissions deleted earlier, or before `earliest_date`, remain deleted. An in-place file restore restores user and group permissions on returned folders and their subfolders that were deleted from the time of the folder's deletion through 24 hours afterward, including permissions removed independently within that window. File restores into a new folder do not restore permissions. |
| `RestoreInPlace` | `*bool` | `restore_in_place` | no | If true, we will restore the files in place (into their original paths). If false, we will create a new restoration folder in the root and restore files there. |
| `UpdateTimestamps` | `*bool` | `update_timestamps` | no | If true, we will update the last modified timestamp of restored files to today's date. If false, we might trigger File Expiration to delete the file again. |
| `WorkspaceId` | `int64` | `workspace_id` | no | Workspace ID for a workspace-scoped restore. `0` means the default site-wide scope. |

Example:

```go
// Set this required value before the call.
var earliestDate *time.Time
result, err := restore.Create(files_sdk.RestoreCreateParams{
	EarliestDate: earliestDate,
})
if err != nil {
	return err
}
fmt.Println(result)
```

## Model

`files_sdk.Restore` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `EarliestDate` | `*time.Time` | `earliest_date` | Restore files or users deleted on or after this date/time. Don't set this earlier than you need. Can not be greater than 365 days prior to the restore request. |
| `Id` | `int64` | `id` | Restore Record ID. |
| `DirsRestored` | `int64` | `dirs_restored` | Number of directories that were successfully restored. |
| `DirsErrored` | `int64` | `dirs_errored` | Number of directories that were not able to be restored. |
| `DirsTotal` | `int64` | `dirs_total` | Total number of directories processed. |
| `FilesRestored` | `int64` | `files_restored` | Number of files successfully restored. |
| `FilesErrored` | `int64` | `files_errored` | Number of files that were not able to be restored. |
| `FilesTotal` | `int64` | `files_total` | Total number of files processed. |
| `Prefix` | `string` | `prefix` | Prefix of the files/folders to restore, or a case-insensitive username prefix for a user restore. A trailing slash selects that folder and its contents only; without it, the prefix matches any path beginning with that text. Do not use a leading slash. To restore all deleted items of the selected restoration type within the selected site or workspace scope, specify an empty string (`''`) or omit the field. |
| `RestorationType` | `string` | `restoration_type` | Type of restoration to perform. `files` restores deleted filesystem items. `users` restores deleted users and associated authentication records, plus permissions when `restore_deleted_permissions=true`. A 5-minute grace period allows for batch deletion processes: associated records deleted within the 5 minutes before the user are included, even if removed independently. Records deleted earlier, or before `earliest_date`, remain deleted. |
| `RestoreInPlace` | `*bool` | `restore_in_place` | If true, we will restore the files in place (into their original paths). If false, we will create a new restoration folder in the root and restore files there. |
| `RestoreDeletedPermissions` | `*bool` | `restore_deleted_permissions` | If true, a user restore restores permissions deleted along with the selected users, including permissions removed independently within the 5 minutes before user deletion. Permissions deleted earlier, or before `earliest_date`, remain deleted. An in-place file restore restores user and group permissions on returned folders and their subfolders that were deleted from the time of the folder's deletion through 24 hours afterward, including permissions removed independently within that window. File restores into a new folder do not restore permissions. |
| `UsersRestored` | `int64` | `users_restored` | Number of users successfully restored (only present for `restoration_type=users`). |
| `UsersErrored` | `int64` | `users_errored` | Number of users that failed to restore (only present for `restoration_type=users`). |
| `UsersTotal` | `int64` | `users_total` | Total number of users processed (only present for `restoration_type=users`). |
| `ApiKeysRestored` | `int64` | `api_keys_restored` | Number of API keys restored (only present for `restoration_type=users`). |
| `PublicKeysRestored` | `int64` | `public_keys_restored` | Number of public keys restored (only present for `restoration_type=users`). |
| `TwoFactorAuthenticationMethodsRestored` | `int64` | `two_factor_authentication_methods_restored` | Number of two factor authentication methods restored (only present for `restoration_type=users`). |
| `Status` | `string` | `status` | Status of the restoration process. |
| `UpdateTimestamps` | `*bool` | `update_timestamps` | If true, we will update the last modified timestamp of restored files to today's date. If false, we might trigger File Expiration to delete the file again. |
| `WorkspaceId` | `int64` | `workspace_id` | Workspace ID for a workspace-scoped restore. `0` means the default site-wide scope. |
| `ErrorMessages` | `[]string` | `error_messages` | Error messages received while restoring files and/or directories. Only present if there were errors. |
