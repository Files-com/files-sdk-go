# EventTarget

`files_sdk.EventTarget` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/eventtarget`, named `event_target`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	event_target "github.com/Files-com/files-sdk-go/v3/eventtarget"
)
```

## Operations

### List

`GET /event_targets`

List Event Targets

```go
func (c *Client) List(params files_sdk.EventTargetListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.EventTargetListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.EventTarget` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.EventTargetListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `name`, `enabled` or `workspace_id`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `enabled`, `target_type` or `workspace_id`. Valid field combinations are `[ enabled, target_type ]`, `[ workspace_id, enabled ]` or `[ workspace_id, enabled, target_type ]`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := event_target.List(files_sdk.EventTargetListParams{})
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

`GET /event_targets/{id}`

Show Event Target

```go
func (c *Client) Find(params files_sdk.EventTargetFindParams, opts ...files_sdk.RequestResponseOption) (eventTarget files_sdk.EventTarget, err error)
func Find(params files_sdk.EventTargetFindParams, opts ...files_sdk.RequestResponseOption) (eventTarget files_sdk.EventTarget, err error)
```

Parameters, `files_sdk.EventTargetFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Event Target ID. |

Example:

```go
result, err := event_target.Find(files_sdk.EventTargetFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Create

`POST /event_targets`

Create Event Target

```go
func (c *Client) Create(params files_sdk.EventTargetCreateParams, opts ...files_sdk.RequestResponseOption) (eventTarget files_sdk.EventTarget, err error)
func Create(params files_sdk.EventTargetCreateParams, opts ...files_sdk.RequestResponseOption) (eventTarget files_sdk.EventTarget, err error)
```

Parameters, `files_sdk.EventTargetCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Name` | `string` | `name` | yes | Event Target name. |
| `WorkspaceId` | `int64` | `workspace_id` | no | Workspace ID. 0 means the default workspace or site-wide. |
| `ApplyToAllWorkspaces` | `*bool` | `apply_to_all_workspaces` | no | If true, this default-workspace target can receive events from all workspaces. |
| `Enabled` | `*bool` | `enabled` | no | Whether this Event Target can receive events. |
| `Config` | `interface{}` | `config` | yes | Event Target configuration. Folder targets accept path and format (json or csv). |
| `DeliveryPolicy` | `interface{}` | `delivery_policy` | no | Event Target delivery policy. Email and folder targets support batch_interval in seconds, between 600 and 86400. Slack and Teams support delivery_mode (immediate or subfolder_summary), summary_path (empty string means root), and summary_delay (15–300 seconds, default 60). The summary applies to matching file activity under summary_path; other matching events retain immediate delivery. |
| `TargetType` | `files_sdk.EventTargetTargetTypeEnum` | `target_type` | yes | Event Target type. Values: `email`, `webhook`, `slack_webhook`, `teams_webhook`, `amazon_sns`, `google_pubsub`, `folder`. |

Example:

```go
// Set this required value before the call.
var config interface{}
result, err := event_target.Create(files_sdk.EventTargetCreateParams{
	Name:       "example",
	Config:     config,
	TargetType: "email",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Update

`PATCH /event_targets/{id}`

Update Event Target

```go
func (c *Client) Update(params files_sdk.EventTargetUpdateParams, opts ...files_sdk.RequestResponseOption) (eventTarget files_sdk.EventTarget, err error)
func Update(params files_sdk.EventTargetUpdateParams, opts ...files_sdk.RequestResponseOption) (eventTarget files_sdk.EventTarget, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (eventTarget files_sdk.EventTarget, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (eventTarget files_sdk.EventTarget, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.EventTargetUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Event Target ID. |
| `Name` | `string` | `name` | no | Event Target name. |
| `WorkspaceId` | `int64` | `workspace_id` | no | Workspace ID. 0 means the default workspace or site-wide. |
| `ApplyToAllWorkspaces` | `*bool` | `apply_to_all_workspaces` | no | If true, this default-workspace target can receive events from all workspaces. |
| `Enabled` | `*bool` | `enabled` | no | Whether this Event Target can receive events. |
| `Config` | `interface{}` | `config` | no | Event Target configuration. Folder targets accept path and format (json or csv). |
| `DeliveryPolicy` | `interface{}` | `delivery_policy` | no | Event Target delivery policy. Email and folder targets support batch_interval in seconds, between 600 and 86400. Slack and Teams support delivery_mode (immediate or subfolder_summary), summary_path (empty string means root), and summary_delay (15–300 seconds, default 60). The summary applies to matching file activity under summary_path; other matching events retain immediate delivery. |

Example:

```go
result, err := event_target.Update(files_sdk.EventTargetUpdateParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := event_target.UpdateWithMap(map[string]interface{}{
	"id": 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /event_targets/{id}`

Delete Event Target

```go
func (c *Client) Delete(params files_sdk.EventTargetDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.EventTargetDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.EventTargetDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Event Target ID. |

Example:

```go
err := event_target.Delete(files_sdk.EventTargetDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.EventTarget` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Event Target ID |
| `Name` | `string` | `name` | Event Target name. |
| `TargetType` | `string` | `target_type` | Event Target type. |
| `WorkspaceId` | `int64` | `workspace_id` | Workspace ID. 0 means the default workspace or site-wide. |
| `ApplyToAllWorkspaces` | `*bool` | `apply_to_all_workspaces` | If true, this default-workspace target can receive events from all workspaces. |
| `Enabled` | `*bool` | `enabled` | Whether this Event Target can receive events. |
| `Config` | `interface{}` | `config` | Event Target configuration. Folder targets accept path and format (json or csv). |
| `DeliveryPolicy` | `interface{}` | `delivery_policy` | Event Target delivery policy. Email and folder targets support batch_interval in seconds, between 600 and 86400. Slack and Teams support delivery_mode (immediate or subfolder_summary), summary_path (empty string means root), and summary_delay (15–300 seconds, default 60). The summary applies to matching file activity under summary_path; other matching events retain immediate delivery. |
| `CreatedAt` | `*time.Time` | `created_at` | Event Target create date/time. |
| `UpdatedAt` | `*time.Time` | `updated_at` | Event Target update date/time. |
