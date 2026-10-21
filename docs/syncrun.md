# SyncRun

`files_sdk.SyncRun` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/syncrun`, named `sync_run`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	sync_run "github.com/Files-com/files-sdk-go/v3/syncrun"
)
```

## Operations

### List

`GET /sync_runs`

List Sync Runs

```go
func (c *Client) List(params files_sdk.SyncRunListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.SyncRunListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.SyncRun` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.SyncRunListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `UserId` | `int64` | `user_id` | no | User ID. Provide a value of `0` to operate the current session's user. |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `site_id`, `workspace_id`, `sync_id`, `created_at` or `status`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `created_at`, `status`, `dry_run`, `workspace_id`, `src_remote_server_type`, `dest_remote_server_type` or `sync_id`. Valid field combinations are `[ status, created_at ]`, `[ workspace_id, created_at ]`, `[ src_remote_server_type, created_at ]`, `[ dest_remote_server_type, created_at ]`, `[ sync_id, created_at ]`, `[ workspace_id, status ]`, `[ src_remote_server_type, status ]`, `[ dest_remote_server_type, status ]`, `[ sync_id, status ]`, `[ workspace_id, src_remote_server_type ]`, `[ workspace_id, dest_remote_server_type ]`, `[ workspace_id, sync_id ]`, `[ workspace_id, status, created_at ]`, `[ src_remote_server_type, status, created_at ]`, `[ dest_remote_server_type, status, created_at ]`, `[ sync_id, status, created_at ]`, `[ workspace_id, src_remote_server_type, created_at ]`, `[ workspace_id, dest_remote_server_type, created_at ]`, `[ workspace_id, sync_id, created_at ]`, `[ workspace_id, src_remote_server_type, status ]`, `[ workspace_id, dest_remote_server_type, status ]`, `[ workspace_id, sync_id, status ]`, `[ workspace_id, src_remote_server_type, status, created_at ]`, `[ workspace_id, dest_remote_server_type, status, created_at ]` or `[ workspace_id, sync_id, status, created_at ]`. |
| `FilterGt` | `interface{}` | `filter_gt` | no | If set, return records where the specified field is greater than the supplied value. Valid fields are `created_at`. |
| `FilterGteq` | `interface{}` | `filter_gteq` | no | If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `created_at`. |
| `FilterLt` | `interface{}` | `filter_lt` | no | If set, return records where the specified field is less than the supplied value. Valid fields are `created_at`. |
| `FilterLteq` | `interface{}` | `filter_lteq` | no | If set, return records where the specified field is less than or equal the supplied value. Valid fields are `created_at`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := sync_run.List(files_sdk.SyncRunListParams{})
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

### Find

`GET /sync_runs/{id}`

Show Sync Run

```go
func (c *Client) Find(params files_sdk.SyncRunFindParams, opts ...files_sdk.RequestResponseOption) (syncRun files_sdk.SyncRun, err error)
func Find(params files_sdk.SyncRunFindParams, opts ...files_sdk.RequestResponseOption) (syncRun files_sdk.SyncRun, err error)
```

Parameters, `files_sdk.SyncRunFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Sync Run ID. |

Example:

```go
result, err := sync_run.Find(files_sdk.SyncRunFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

## Model

`files_sdk.SyncRun` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | SyncRun ID |
| `Body` | `string` | `body` | Log or summary body for this run |
| `BytesSynced` | `int64` | `bytes_synced` | Total bytes synced in this run |
| `ComparedFiles` | `int64` | `compared_files` | Number of files compared |
| `ComparedFolders` | `int64` | `compared_folders` | Number of folders compared |
| `CompletedAt` | `*time.Time` | `completed_at` | When this run was completed |
| `CreatedAt` | `*time.Time` | `created_at` | When this run was created |
| `DestRemoteServerType` | `string` | `dest_remote_server_type` | Destination remote server type, if any |
| `DryRun` | `*bool` | `dry_run` | Whether this run was a dry run (no actual changes made) |
| `ErroredFiles` | `int64` | `errored_files` | Number of files that errored |
| `EstimatedBytesCount` | `int64` | `estimated_bytes_count` | Estimated bytes count for this run |
| `EventErrors` | `[]string` | `event_errors` | Array of errors encountered during the run |
| `LogUrl` | `string` | `log_url` | Link to external log file. |
| `Runtime` | `float64` | `runtime` | Total runtime in seconds |
| `SiteId` | `int64` | `site_id` | Site ID |
| `WorkspaceId` | `int64` | `workspace_id` | Workspace ID |
| `SrcRemoteServerType` | `string` | `src_remote_server_type` | Source remote server type, if any |
| `Status` | `string` | `status` | Status of the sync run (success, failure, partial_failure, in_progress, skipped) |
| `SuccessfulFiles` | `int64` | `successful_files` | Number of files successfully synced |
| `SyncId` | `int64` | `sync_id` | ID of the Sync this run belongs to |
| `SyncName` | `string` | `sync_name` | Name of the Sync this run belongs to |
| `UpdatedAt` | `*time.Time` | `updated_at` | When this run was last updated |
| `LiveTransfers` | `[]files_sdk.SyncRunLiveTransfer` | `live_transfers` | Array of in-progress file transfers with progress data. Only present when the sync run status is in_progress. |
