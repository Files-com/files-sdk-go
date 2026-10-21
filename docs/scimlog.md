# ScimLog

`files_sdk.ScimLog` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/scimlog`, named `scim_log`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	scim_log "github.com/Files-com/files-sdk-go/v3/scimlog"
)
```

## Operations

### List

`GET /scim_logs`

List Scim Logs

```go
func (c *Client) List(params files_sdk.ScimLogListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.ScimLogListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.ScimLog` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.ScimLogListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `created_at`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := scim_log.List(files_sdk.ScimLogListParams{})
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

`GET /scim_logs/{id}`

Show Scim Log

```go
func (c *Client) Find(params files_sdk.ScimLogFindParams, opts ...files_sdk.RequestResponseOption) (scimLog files_sdk.ScimLog, err error)
func Find(params files_sdk.ScimLogFindParams, opts ...files_sdk.RequestResponseOption) (scimLog files_sdk.ScimLog, err error)
```

Parameters, `files_sdk.ScimLogFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Scim Log ID. |

Example:

```go
result, err := scim_log.Find(files_sdk.ScimLogFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

## Model

`files_sdk.ScimLog` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | The unique ID of this SCIM request. |
| `CreatedAt` | `string` | `created_at` | The date and time when this SCIM request occurred. |
| `RequestPath` | `string` | `request_path` | The path portion of the URL requested. |
| `RequestMethod` | `string` | `request_method` | The HTTP method used for this request. |
| `HttpResponseCode` | `string` | `http_response_code` | The HTTP response code returned for this request. |
| `UserAgent` | `string` | `user_agent` | The User-Agent header sent with the request. |
| `RequestJson` | `string` | `request_json` | The JSON payload sent with the request. |
| `ResponseJson` | `string` | `response_json` | The JSON payload returned in the response. |
