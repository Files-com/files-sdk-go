# Workspace

`files_sdk.Workspace` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/workspace`, named `workspace`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/workspace"
)
```

## Operations

### List

`GET /workspaces`

List Workspaces

```go
func (c *Client) List(params files_sdk.WorkspaceListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.WorkspaceListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.Workspace` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.WorkspaceListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `name`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `name`. |
| `FilterPrefix` | `interface{}` | `filter_prefix` | no | If set, return records where the specified field is prefixed by the supplied value. Valid fields are `name`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := workspace.List(files_sdk.WorkspaceListParams{})
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

`GET /workspaces/{id}`

Show Workspace

```go
func (c *Client) Find(params files_sdk.WorkspaceFindParams, opts ...files_sdk.RequestResponseOption) (workspace files_sdk.Workspace, err error)
func Find(params files_sdk.WorkspaceFindParams, opts ...files_sdk.RequestResponseOption) (workspace files_sdk.Workspace, err error)
```

Parameters, `files_sdk.WorkspaceFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Workspace ID. |

Example:

```go
result, err := workspace.Find(files_sdk.WorkspaceFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Create

`POST /workspaces`

Create Workspace

```go
func (c *Client) Create(params files_sdk.WorkspaceCreateParams, opts ...files_sdk.RequestResponseOption) (workspace files_sdk.Workspace, err error)
func Create(params files_sdk.WorkspaceCreateParams, opts ...files_sdk.RequestResponseOption) (workspace files_sdk.Workspace, err error)
```

Parameters, `files_sdk.WorkspaceCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Name` | `string` | `name` | no | Workspace name |

Example:

```go
result, err := workspace.Create(files_sdk.WorkspaceCreateParams{})
if err != nil {
	return err
}
fmt.Println(result)
```

### Update

`PATCH /workspaces/{id}`

Update Workspace

```go
func (c *Client) Update(params files_sdk.WorkspaceUpdateParams, opts ...files_sdk.RequestResponseOption) (workspace files_sdk.Workspace, err error)
func Update(params files_sdk.WorkspaceUpdateParams, opts ...files_sdk.RequestResponseOption) (workspace files_sdk.Workspace, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (workspace files_sdk.Workspace, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (workspace files_sdk.Workspace, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.WorkspaceUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Workspace ID. |
| `Name` | `string` | `name` | no | Workspace name |

Example:

```go
result, err := workspace.Update(files_sdk.WorkspaceUpdateParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := workspace.UpdateWithMap(map[string]interface{}{
	"id": 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /workspaces/{id}`

Delete Workspace

```go
func (c *Client) Delete(params files_sdk.WorkspaceDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.WorkspaceDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.WorkspaceDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Workspace ID. |

Example:

```go
err := workspace.Delete(files_sdk.WorkspaceDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.Workspace` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Workspace ID |
| `Name` | `string` | `name` | Workspace name |
