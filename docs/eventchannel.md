# EventChannel

`files_sdk.EventChannel` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/eventchannel`, named `event_channel`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	event_channel "github.com/Files-com/files-sdk-go/v3/eventchannel"
)
```

## Operations

### List

`GET /event_channels`

List Event Channels

```go
func (c *Client) List(params files_sdk.EventChannelListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.EventChannelListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.EventChannel` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.EventChannelListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `name`, `enabled`, `default_channel` or `workspace_id`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `enabled`, `default_channel` or `workspace_id`. Valid field combinations are `[ workspace_id, enabled ]` and `[ workspace_id, default_channel ]`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := event_channel.List(files_sdk.EventChannelListParams{})
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

`GET /event_channels/{id}`

Show Event Channel

```go
func (c *Client) Find(params files_sdk.EventChannelFindParams, opts ...files_sdk.RequestResponseOption) (eventChannel files_sdk.EventChannel, err error)
func Find(params files_sdk.EventChannelFindParams, opts ...files_sdk.RequestResponseOption) (eventChannel files_sdk.EventChannel, err error)
```

Parameters, `files_sdk.EventChannelFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Event Channel ID. |

Example:

```go
result, err := event_channel.Find(files_sdk.EventChannelFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Create

`POST /event_channels`

Create Event Channel

```go
func (c *Client) Create(params files_sdk.EventChannelCreateParams, opts ...files_sdk.RequestResponseOption) (eventChannel files_sdk.EventChannel, err error)
func Create(params files_sdk.EventChannelCreateParams, opts ...files_sdk.RequestResponseOption) (eventChannel files_sdk.EventChannel, err error)
```

Parameters, `files_sdk.EventChannelCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Name` | `string` | `name` | yes | Event Channel name. |
| `WorkspaceId` | `int64` | `workspace_id` | no | Workspace ID. 0 means the default workspace. |
| `Description` | `string` | `description` | no | Event Channel description. |
| `Enabled` | `*bool` | `enabled` | no | Whether this Event Channel can dispatch events. |
| `DefaultChannel` | `*bool` | `default_channel` | no | Whether this Event Channel is the default destination for newly published events. |

Example:

```go
result, err := event_channel.Create(files_sdk.EventChannelCreateParams{
	Name: "example",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Update

`PATCH /event_channels/{id}`

Update Event Channel

```go
func (c *Client) Update(params files_sdk.EventChannelUpdateParams, opts ...files_sdk.RequestResponseOption) (eventChannel files_sdk.EventChannel, err error)
func Update(params files_sdk.EventChannelUpdateParams, opts ...files_sdk.RequestResponseOption) (eventChannel files_sdk.EventChannel, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (eventChannel files_sdk.EventChannel, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (eventChannel files_sdk.EventChannel, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.EventChannelUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Event Channel ID. |
| `Name` | `string` | `name` | no | Event Channel name. |
| `WorkspaceId` | `int64` | `workspace_id` | no | Workspace ID. 0 means the default workspace. |
| `Description` | `string` | `description` | no | Event Channel description. |
| `Enabled` | `*bool` | `enabled` | no | Whether this Event Channel can dispatch events. |
| `DefaultChannel` | `*bool` | `default_channel` | no | Whether this Event Channel is the default destination for newly published events. |

Example:

```go
result, err := event_channel.Update(files_sdk.EventChannelUpdateParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := event_channel.UpdateWithMap(map[string]interface{}{
	"id": 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /event_channels/{id}`

Delete Event Channel

```go
func (c *Client) Delete(params files_sdk.EventChannelDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.EventChannelDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.EventChannelDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Event Channel ID. |

Example:

```go
err := event_channel.Delete(files_sdk.EventChannelDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.EventChannel` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Event Channel ID |
| `Name` | `string` | `name` | Event Channel name. |
| `WorkspaceId` | `int64` | `workspace_id` | Workspace ID. 0 means the default workspace. |
| `Description` | `string` | `description` | Event Channel description. |
| `Enabled` | `*bool` | `enabled` | Whether this Event Channel can dispatch events. |
| `DefaultChannel` | `*bool` | `default_channel` | Whether this Event Channel is the default destination for newly published events. |
| `CreatedAt` | `*time.Time` | `created_at` | Event Channel create date/time. |
| `UpdatedAt` | `*time.Time` | `updated_at` | Event Channel update date/time. |
