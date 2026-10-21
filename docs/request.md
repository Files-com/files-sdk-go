# Request

`files_sdk.Request` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/request`, named `request`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/request"
)
```

## Operations

### List

`GET /requests`

List Requests

```go
func (c *Client) List(params files_sdk.RequestListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.RequestListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.Request` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.RequestListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are . |
| `Mine` | `*bool` | `mine` | no | Only show requests of the current user? (Defaults to true if current user is not a site admin.) |
| `Path` | `string` | `path` | no | Path to show requests for. If omitted, shows all paths. Send `/` to represent the root directory. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := request.List(files_sdk.RequestListParams{})
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

### GetFolder

`GET /requests/folders/{path}`

List Requests

```go
func (c *Client) GetFolder(params files_sdk.RequestGetFolderParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func GetFolder(params files_sdk.RequestGetFolderParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.Request` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.RequestGetFolderParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are . |
| `Mine` | `*bool` | `mine` | no | Only show requests of the current user? (Defaults to true if current user is not a site admin.) |
| `Path` | `string` | `path`, in the path | yes | Path to show requests for. If omitted, shows all paths. Send `/` to represent the root directory. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := request.GetFolder(files_sdk.RequestGetFolderParams{
	Path: "path",
})
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

### Create

`POST /requests`

Create Request

```go
func (c *Client) Create(params files_sdk.RequestCreateParams, opts ...files_sdk.RequestResponseOption) (request files_sdk.Request, err error)
func Create(params files_sdk.RequestCreateParams, opts ...files_sdk.RequestResponseOption) (request files_sdk.Request, err error)
```

Parameters, `files_sdk.RequestCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Path` | `string` | `path` | yes | Folder path on which to request the file. |
| `Destination` | `string` | `destination` | yes | Destination filename (without extension) to request. |
| `UserIds` | `string` | `user_ids` | no | A list of user IDs to request the file from. If sent as a string, it should be comma-delimited. |
| `GroupIds` | `string` | `group_ids` | no | A list of group IDs to request the file from. If sent as a string, it should be comma-delimited. |

Example:

```go
result, err := request.Create(files_sdk.RequestCreateParams{
	Path:        "path",
	Destination: "destination",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /requests/{id}`

Delete Request

```go
func (c *Client) Delete(params files_sdk.RequestDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.RequestDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.RequestDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Request ID. |

Example:

```go
err := request.Delete(files_sdk.RequestDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.Request` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Request ID |
| `Path` | `string` | `path` | Folder path. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters. |
| `Source` | `string` | `source` | Source filename, if applicable |
| `Destination` | `string` | `destination` | Destination filename |
| `AutomationId` | `int64` | `automation_id` | ID of automation that created request |
| `UserDisplayName` | `string` | `user_display_name` | User making the request (if applicable) |
| `UserIds` | `string` | `user_ids` | A list of user IDs to request the file from. If sent as a string, it should be comma-delimited. |
| `GroupIds` | `string` | `group_ids` | A list of group IDs to request the file from. If sent as a string, it should be comma-delimited. |
