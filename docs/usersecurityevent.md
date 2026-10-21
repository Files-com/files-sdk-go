# UserSecurityEvent

`files_sdk.UserSecurityEvent` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/usersecurityevent`, named `user_security_event`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	user_security_event "github.com/Files-com/files-sdk-go/v3/usersecurityevent"
)
```

## Operations

### List

`GET /user_security_events`

List User Security Events

```go
func (c *Client) List(params files_sdk.UserSecurityEventListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.UserSecurityEventListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.UserSecurityEvent` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.UserSecurityEventListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `created_at` and `user_id`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `created_at` and `user_id`. Valid field combinations are `[ user_id, created_at ]`. |
| `FilterGt` | `interface{}` | `filter_gt` | no | If set, return records where the specified field is greater than the supplied value. Valid fields are `created_at`. |
| `FilterGteq` | `interface{}` | `filter_gteq` | no | If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `created_at`. |
| `FilterLt` | `interface{}` | `filter_lt` | no | If set, return records where the specified field is less than the supplied value. Valid fields are `created_at`. |
| `FilterLteq` | `interface{}` | `filter_lteq` | no | If set, return records where the specified field is less than or equal the supplied value. Valid fields are `created_at`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := user_security_event.List(files_sdk.UserSecurityEventListParams{})
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

`GET /user_security_events/{id}`

Show User Security Event

```go
func (c *Client) Find(params files_sdk.UserSecurityEventFindParams, opts ...files_sdk.RequestResponseOption) (userSecurityEvent files_sdk.UserSecurityEvent, err error)
func Find(params files_sdk.UserSecurityEventFindParams, opts ...files_sdk.RequestResponseOption) (userSecurityEvent files_sdk.UserSecurityEvent, err error)
```

Parameters, `files_sdk.UserSecurityEventFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | User Security Event ID. |

Example:

```go
result, err := user_security_event.Find(files_sdk.UserSecurityEventFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

## Model

`files_sdk.UserSecurityEvent` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Event ID |
| `EventType` | `string` | `event_type` | Type of user security event being recorded. |
| `Body` | `string` | `body` | Event body. |
| `EventErrors` | `[]string` | `event_errors` | Event errors. |
| `CreatedAt` | `*time.Time` | `created_at` | Event create date/time. |
| `BodyUrl` | `string` | `body_url` | Link to log file. |
| `UserId` | `int64` | `user_id` | User ID. |
