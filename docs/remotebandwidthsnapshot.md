# RemoteBandwidthSnapshot

`files_sdk.RemoteBandwidthSnapshot` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/remotebandwidthsnapshot`, named `remote_bandwidth_snapshot`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	remote_bandwidth_snapshot "github.com/Files-com/files-sdk-go/v3/remotebandwidthsnapshot"
)
```

## Operations

### List

`GET /remote_bandwidth_snapshots`

List Remote Bandwidth Snapshots

```go
func (c *Client) List(params files_sdk.RemoteBandwidthSnapshotListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.RemoteBandwidthSnapshotListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.RemoteBandwidthSnapshot` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.RemoteBandwidthSnapshotListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `logged_at`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `logged_at`. |
| `FilterGt` | `interface{}` | `filter_gt` | no | If set, return records where the specified field is greater than the supplied value. Valid fields are `logged_at`. |
| `FilterGteq` | `interface{}` | `filter_gteq` | no | If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `logged_at`. |
| `FilterLt` | `interface{}` | `filter_lt` | no | If set, return records where the specified field is less than the supplied value. Valid fields are `logged_at`. |
| `FilterLteq` | `interface{}` | `filter_lteq` | no | If set, return records where the specified field is less than or equal the supplied value. Valid fields are `logged_at`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := remote_bandwidth_snapshot.List(files_sdk.RemoteBandwidthSnapshotListParams{})
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

`files_sdk.RemoteBandwidthSnapshot` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Site bandwidth ID |
| `SyncBytesReceived` | `int64` | `sync_bytes_received` | Site sync bandwidth report bytes received |
| `SyncBytesSent` | `int64` | `sync_bytes_sent` | Site sync bandwidth report bytes sent |
| `LoggedAt` | `*time.Time` | `logged_at` | Time the site bandwidth report was logged |
| `RemoteServerId` | `int64` | `remote_server_id` | ID of related Remote Server |
