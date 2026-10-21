# UsageSnapshot

`files_sdk.UsageSnapshot` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/usagesnapshot`, named `usage_snapshot`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	usage_snapshot "github.com/Files-com/files-sdk-go/v3/usagesnapshot"
)
```

## Operations

### List

`GET /usage_snapshots`

List Usage Snapshots

```go
func (c *Client) List(params files_sdk.UsageSnapshotListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.UsageSnapshotListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.UsageSnapshot` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.UsageSnapshotListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := usage_snapshot.List(files_sdk.UsageSnapshotListParams{})
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

`files_sdk.UsageSnapshot` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Usage snapshot ID |
| `StartAt` | `*time.Time` | `start_at` | Usage snapshot start date/time |
| `EndAt` | `*time.Time` | `end_at` | Usage snapshot end date/time |
| `HighWaterUserCount` | `int64` | `high_water_user_count` | Highest user count number in time period |
| `TransformationCreditsUsage` | `int64` | `transformation_credits_usage` | Total Transformation and AI credits used in this usage period. Updated daily. Does not include credits consumed today. |
| `CurrentStorage` | `int64` | `current_storage` | Current total Storage Usage GB as of end date (not necessarily high water mark, which is used for billing) |
| `HighWaterStorage` | `int64` | `high_water_storage` | Highest Storage Usage GB recorded in time period (used for billing) |
| `RootStorage` | `int64` | `root_storage` | Storage Usage for root folder as of end date (not necessarily high water mark, which is used for billing) |
| `DeletedFilesCountedInMinimum` | `int64` | `deleted_files_counted_in_minimum` | Storage Usage for files that are deleted but uploaded within last 30 days as of end date (not necessarily high water mark, which is used for billing) |
| `DeletedFilesStorage` | `int64` | `deleted_files_storage` | Storage Usage for files that are deleted but retained as backups as of end date (not necessarily high water mark, which is used for billing) |
| `TotalBillableUsage` | `int64` | `total_billable_usage` | Storage + Transfer Usage - Total Billable amount |
| `TotalBillableTransferUsage` | `int64` | `total_billable_transfer_usage` | Transfer usage for period - Total Billable amount |
| `BytesSent` | `int64` | `bytes_sent` | Transfer Usage for period - Outbound GB from Files Native Storage |
| `SyncBytesReceived` | `int64` | `sync_bytes_received` | Transfer Usage for period - Inbound GB to Remote Servers (Sync/Mount) |
| `SyncBytesSent` | `int64` | `sync_bytes_sent` | Transfer Usage for period - Outbound GB from Remote Servers (Sync/Mount) |
| `UsageByTopLevelDir` | `[]files_sdk.UsageByTopLevelDir` | `usage_by_top_level_dir` | Storage Usage - map of root folders to their usage as of end date (not necessarily high water mark, which is used for billing) |
