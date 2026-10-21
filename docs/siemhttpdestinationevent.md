# SiemHttpDestinationEvent

`files_sdk.SiemHttpDestinationEvent` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/siemhttpdestinationevent`, named `siem_http_destination_event`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	siem_http_destination_event "github.com/Files-com/files-sdk-go/v3/siemhttpdestinationevent"
)
```

## Operations

### List

`GET /siem_http_destination_events`

List SIEM HTTP Destination Events

```go
func (c *Client) List(params files_sdk.SiemHttpDestinationEventListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.SiemHttpDestinationEventListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.SiemHttpDestinationEvent` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.SiemHttpDestinationEventListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `created_at`, `status` or `siem_http_destination_id`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `created_at`, `siem_http_destination_id` or `status`. Valid field combinations are `[ siem_http_destination_id, created_at ]`, `[ status, created_at ]`, `[ siem_http_destination_id, status ]` or `[ siem_http_destination_id, status, created_at ]`. |
| `FilterGt` | `interface{}` | `filter_gt` | no | If set, return records where the specified field is greater than the supplied value. Valid fields are `created_at`. |
| `FilterGteq` | `interface{}` | `filter_gteq` | no | If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `created_at`. |
| `FilterLt` | `interface{}` | `filter_lt` | no | If set, return records where the specified field is less than the supplied value. Valid fields are `created_at`. |
| `FilterLteq` | `interface{}` | `filter_lteq` | no | If set, return records where the specified field is less than or equal the supplied value. Valid fields are `created_at`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := siem_http_destination_event.List(files_sdk.SiemHttpDestinationEventListParams{})
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

`GET /siem_http_destination_events/{id}`

Show SIEM HTTP Destination Event

```go
func (c *Client) Find(params files_sdk.SiemHttpDestinationEventFindParams, opts ...files_sdk.RequestResponseOption) (siemHttpDestinationEvent files_sdk.SiemHttpDestinationEvent, err error)
func Find(params files_sdk.SiemHttpDestinationEventFindParams, opts ...files_sdk.RequestResponseOption) (siemHttpDestinationEvent files_sdk.SiemHttpDestinationEvent, err error)
```

Parameters, `files_sdk.SiemHttpDestinationEventFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Siem Http Destination Event ID. |

Example:

```go
result, err := siem_http_destination_event.Find(files_sdk.SiemHttpDestinationEventFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

## Model

`files_sdk.SiemHttpDestinationEvent` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Event ID |
| `EventType` | `string` | `event_type` | Type of SIEM event being recorded. |
| `Status` | `string` | `status` | Status of event. |
| `Body` | `string` | `body` | Event body. |
| `EventErrors` | `[]string` | `event_errors` | Event errors. |
| `CreatedAt` | `*time.Time` | `created_at` | Event create date/time. |
| `BodyUrl` | `string` | `body_url` | Link to log file. |
| `SiemHttpDestinationId` | `int64` | `siem_http_destination_id` | SIEM ID. |
