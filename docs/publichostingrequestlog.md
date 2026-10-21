# PublicHostingRequestLog

`files_sdk.PublicHostingRequestLog` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/publichostingrequestlog`, named `public_hosting_request_log`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	public_hosting_request_log "github.com/Files-com/files-sdk-go/v3/publichostingrequestlog"
)
```

## Operations

### List

`GET /public_hosting_request_logs`

List Public Hosting Request Logs

```go
func (c *Client) List(params files_sdk.PublicHostingRequestLogListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.PublicHostingRequestLogListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.PublicHostingRequestLog` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.PublicHostingRequestLogListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `path`, `remote_ip`, `success` or `created_at`. Valid field combinations are `[ path ]`, `[ remote_ip ]`, `[ success ]`, `[ created_at ]`, `[ path, remote_ip ]`, `[ path, success ]`, `[ path, created_at ]`, `[ remote_ip, success ]`, `[ remote_ip, created_at ]`, `[ success, created_at ]`, `[ path, remote_ip, success ]`, `[ path, remote_ip, created_at ]`, `[ path, success, created_at ]`, `[ remote_ip, success, created_at ]` or `[ path, remote_ip, success, created_at ]`. |
| `FilterGt` | `interface{}` | `filter_gt` | no | If set, return records where the specified field is greater than the supplied value. Valid fields are `created_at`. Valid field combinations are `[ path ]`, `[ remote_ip ]`, `[ success ]`, `[ created_at ]`, `[ path, remote_ip ]`, `[ path, success ]`, `[ path, created_at ]`, `[ remote_ip, success ]`, `[ remote_ip, created_at ]`, `[ success, created_at ]`, `[ path, remote_ip, success ]`, `[ path, remote_ip, created_at ]`, `[ path, success, created_at ]`, `[ remote_ip, success, created_at ]` or `[ path, remote_ip, success, created_at ]`. |
| `FilterGteq` | `interface{}` | `filter_gteq` | no | If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `created_at`. Valid field combinations are `[ path ]`, `[ remote_ip ]`, `[ success ]`, `[ created_at ]`, `[ path, remote_ip ]`, `[ path, success ]`, `[ path, created_at ]`, `[ remote_ip, success ]`, `[ remote_ip, created_at ]`, `[ success, created_at ]`, `[ path, remote_ip, success ]`, `[ path, remote_ip, created_at ]`, `[ path, success, created_at ]`, `[ remote_ip, success, created_at ]` or `[ path, remote_ip, success, created_at ]`. |
| `FilterPrefix` | `interface{}` | `filter_prefix` | no | If set, return records where the specified field is prefixed by the supplied value. Valid fields are `path`. Valid field combinations are `[ path ]`, `[ remote_ip ]`, `[ success ]`, `[ created_at ]`, `[ path, remote_ip ]`, `[ path, success ]`, `[ path, created_at ]`, `[ remote_ip, success ]`, `[ remote_ip, created_at ]`, `[ success, created_at ]`, `[ path, remote_ip, success ]`, `[ path, remote_ip, created_at ]`, `[ path, success, created_at ]`, `[ remote_ip, success, created_at ]` or `[ path, remote_ip, success, created_at ]`. |
| `FilterLt` | `interface{}` | `filter_lt` | no | If set, return records where the specified field is less than the supplied value. Valid fields are `created_at`. Valid field combinations are `[ path ]`, `[ remote_ip ]`, `[ success ]`, `[ created_at ]`, `[ path, remote_ip ]`, `[ path, success ]`, `[ path, created_at ]`, `[ remote_ip, success ]`, `[ remote_ip, created_at ]`, `[ success, created_at ]`, `[ path, remote_ip, success ]`, `[ path, remote_ip, created_at ]`, `[ path, success, created_at ]`, `[ remote_ip, success, created_at ]` or `[ path, remote_ip, success, created_at ]`. |
| `FilterLteq` | `interface{}` | `filter_lteq` | no | If set, return records where the specified field is less than or equal the supplied value. Valid fields are `created_at`. Valid field combinations are `[ path ]`, `[ remote_ip ]`, `[ success ]`, `[ created_at ]`, `[ path, remote_ip ]`, `[ path, success ]`, `[ path, created_at ]`, `[ remote_ip, success ]`, `[ remote_ip, created_at ]`, `[ success, created_at ]`, `[ path, remote_ip, success ]`, `[ path, remote_ip, created_at ]`, `[ path, success, created_at ]`, `[ remote_ip, success, created_at ]` or `[ path, remote_ip, success, created_at ]`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := public_hosting_request_log.List(files_sdk.PublicHostingRequestLogListParams{})
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

`files_sdk.PublicHostingRequestLog` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Timestamp` | `*time.Time` | `timestamp` | Start Time of Action. Deprecrated: Use created_at. |
| `RemoteIp` | `string` | `remote_ip` | IP Address of Public Hosting Client. |
| `ServerIp` | `string` | `server_ip` | IP Address of Public Hosting Server. |
| `Hostname` | `string` | `hostname` | HTTP Request Hostname. |
| `Path` | `string` | `path` | HTTP Request Path. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters. |
| `ResponseCode` | `int64` | `responseCode` | HTTP Response Code. |
| `Success` | `*bool` | `success` | Whether SFTP Action was successful. |
| `DurationMs` | `int64` | `duration_ms` | Duration (in milliseconds). |
| `CreatedAt` | `*time.Time` | `created_at` | Start Time of Action. |
| `BytesTransferred` | `int64` | `bytes_transferred` | The number of bytes transferred for file downloads. |
| `HttpMethod` | `string` | `http_method` | Method of the HTTP call. |
