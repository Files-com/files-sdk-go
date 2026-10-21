# FileComment

`files_sdk.FileComment` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/filecomment`, named `file_comment`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	file_comment "github.com/Files-com/files-sdk-go/v3/filecomment"
)
```

## Operations

### ListFor

`GET /file_comments/files/{path}`

List File Comments by Path

```go
func (c *Client) ListFor(params files_sdk.FileCommentListForParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func ListFor(params files_sdk.FileCommentListForParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.FileComment` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.FileCommentListForParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Path` | `string` | `path`, in the path | yes | Path to operate on. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := file_comment.ListFor(files_sdk.FileCommentListForParams{
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

`POST /file_comments`

Create File Comment

```go
func (c *Client) Create(params files_sdk.FileCommentCreateParams, opts ...files_sdk.RequestResponseOption) (fileComment files_sdk.FileComment, err error)
func Create(params files_sdk.FileCommentCreateParams, opts ...files_sdk.RequestResponseOption) (fileComment files_sdk.FileComment, err error)
```

Parameters, `files_sdk.FileCommentCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Body` | `string` | `body` | yes | Comment body. |
| `Path` | `string` | `path` | yes | File path. |

Example:

```go
result, err := file_comment.Create(files_sdk.FileCommentCreateParams{
	Body: "What a great file!",
	Path: "path",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Update

`PATCH /file_comments/{id}`

Update File Comment

```go
func (c *Client) Update(params files_sdk.FileCommentUpdateParams, opts ...files_sdk.RequestResponseOption) (fileComment files_sdk.FileComment, err error)
func Update(params files_sdk.FileCommentUpdateParams, opts ...files_sdk.RequestResponseOption) (fileComment files_sdk.FileComment, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (fileComment files_sdk.FileComment, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (fileComment files_sdk.FileComment, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.FileCommentUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | File Comment ID. |
| `Body` | `string` | `body` | yes | Comment body. |

Example:

```go
result, err := file_comment.Update(files_sdk.FileCommentUpdateParams{
	Id:   1,
	Body: "What a great file!",
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := file_comment.UpdateWithMap(map[string]interface{}{
	"id":   1,
	"body": "What a great file!",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /file_comments/{id}`

Delete File Comment

```go
func (c *Client) Delete(params files_sdk.FileCommentDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.FileCommentDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.FileCommentDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | File Comment ID. |

Example:

```go
err := file_comment.Delete(files_sdk.FileCommentDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.FileComment` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | File Comment ID |
| `Body` | `string` | `body` | Comment body. |
| `Reactions` | `[]files_sdk.FileCommentReaction` | `reactions` | Reactions to this comment. |
| `Path` | `string` | `path` | File path. |
