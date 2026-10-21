# ExternalEvent

`files_sdk.ExternalEvent` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/externalevent`, named `external_event`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	external_event "github.com/Files-com/files-sdk-go/v3/externalevent"
)
```

## Operations

### List

`GET /external_events`

List External Events

```go
func (c *Client) List(params files_sdk.ExternalEventListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.ExternalEventListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.ExternalEvent` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.ExternalEventListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `created_at`, `status` or `event_type`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `created_at` and `status`. Valid field combinations are `[ status, created_at ]`. |
| `FilterGt` | `interface{}` | `filter_gt` | no | If set, return records where the specified field is greater than the supplied value. Valid fields are `created_at`. |
| `FilterGteq` | `interface{}` | `filter_gteq` | no | If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `created_at`. |
| `FilterLt` | `interface{}` | `filter_lt` | no | If set, return records where the specified field is less than the supplied value. Valid fields are `created_at`. |
| `FilterLteq` | `interface{}` | `filter_lteq` | no | If set, return records where the specified field is less than or equal the supplied value. Valid fields are `created_at`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := external_event.List(files_sdk.ExternalEventListParams{})
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

`GET /external_events/{id}`

Show External Event

```go
func (c *Client) Find(params files_sdk.ExternalEventFindParams, opts ...files_sdk.RequestResponseOption) (externalEvent files_sdk.ExternalEvent, err error)
func Find(params files_sdk.ExternalEventFindParams, opts ...files_sdk.RequestResponseOption) (externalEvent files_sdk.ExternalEvent, err error)
```

Parameters, `files_sdk.ExternalEventFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | External Event ID. |

Example:

```go
result, err := external_event.Find(files_sdk.ExternalEventFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Create

`POST /external_events`

Create External Event

```go
func (c *Client) Create(params files_sdk.ExternalEventCreateParams, opts ...files_sdk.RequestResponseOption) (externalEvent files_sdk.ExternalEvent, err error)
func Create(params files_sdk.ExternalEventCreateParams, opts ...files_sdk.RequestResponseOption) (externalEvent files_sdk.ExternalEvent, err error)
```

Parameters, `files_sdk.ExternalEventCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Status` | `files_sdk.ExternalEventStatusEnum` | `status` | yes | Status of event. Values: `success`, `failure`, `partial_failure`, `in_progress`, `skipped`. |
| `Body` | `string` | `body` | yes | Event body |

Example:

```go
result, err := external_event.Create(files_sdk.ExternalEventCreateParams{
	Status: "success",
	Body:   "example",
})
if err != nil {
	return err
}
fmt.Println(result)
```

## Model

`files_sdk.ExternalEvent` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Event ID |
| `EventType` | `string` | `event_type` | Type of event being recorded. |
| `Status` | `string` | `status` | Status of event. |
| `Body` | `string` | `body` | Event body |
| `CreatedAt` | `*time.Time` | `created_at` | External event create date/time |
| `BodyUrl` | `string` | `body_url` | Link to log file. |
