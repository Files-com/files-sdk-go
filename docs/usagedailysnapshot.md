# UsageDailySnapshot

`files_sdk.UsageDailySnapshot` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/usagedailysnapshot`, named `usage_daily_snapshot`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	usage_daily_snapshot "github.com/Files-com/files-sdk-go/v3/usagedailysnapshot"
)
```

## Operations

### List

`GET /usage_daily_snapshots`

List Usage Daily Snapshots

```go
func (c *Client) List(params files_sdk.UsageDailySnapshotListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.UsageDailySnapshotListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.UsageDailySnapshot` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.UsageDailySnapshotListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `date`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `date` and `usage_snapshot_id`. Valid field combinations are `[ usage_snapshot_id, date ]`. |
| `FilterGt` | `interface{}` | `filter_gt` | no | If set, return records where the specified field is greater than the supplied value. Valid fields are `date`. |
| `FilterGteq` | `interface{}` | `filter_gteq` | no | If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `date`. |
| `FilterLt` | `interface{}` | `filter_lt` | no | If set, return records where the specified field is less than the supplied value. Valid fields are `date`. |
| `FilterLteq` | `interface{}` | `filter_lteq` | no | If set, return records where the specified field is less than or equal the supplied value. Valid fields are `date`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := usage_daily_snapshot.List(files_sdk.UsageDailySnapshotListParams{})
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

`files_sdk.UsageDailySnapshot` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | ID of the usage record |
| `Date` | `*calendar.Date` | `date` | The date of this usage record |
| `ApiUsageAvailable` | `*bool` | `api_usage_available` | True if the API usage fields `read_api_usage` and `write_api_usage` can be relied upon. If this is false, we suggest hiding that value from any UI. |
| `ReadApiUsage` | `int64` | `read_api_usage` | Read API Calls used on this day. Note: only updated for days before the current day. |
| `WriteApiUsage` | `int64` | `write_api_usage` | Write API Calls used on this day. Note: only updated for days before the current day. |
| `TransformationCreditsUsage` | `int64` | `transformation_credits_usage` | Transformation and AI credits used on this day. Note: only updated for days before the current day. |
| `UserCount` | `int64` | `user_count` | Number of billable users as of this day. |
| `CurrentStorage` | `int64` | `current_storage` | GB of Files Native Storage used on this day. |
| `DeletedFilesStorage` | `int64` | `deleted_files_storage` | GB of Files Native Storage used on this day for files that have been deleted and are stored as backups. |
| `DeletedFilesCountedInMinimum` | `int64` | `deleted_files_counted_in_minimum` | GB of Files Native Storage used on this day for files that have been permanently deleted but were uploaded less than 30 days ago, and are still billable. |
| `RootStorage` | `int64` | `root_storage` | GB of Files Native Storage used for the root folder. Included here because this value will not be part of `usage_by_top_level_dir` |
| `UsageByTopLevelDir` | `[]files_sdk.UsageByTopLevelDir` | `usage_by_top_level_dir` | Usage broken down by each top-level folder |
