# Style

`files_sdk.Style` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/style`, named `style`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/style"
)
```

## Operations

### Find

`GET /styles/{path}`

Show Style

```go
func (c *Client) Find(params files_sdk.StyleFindParams, opts ...files_sdk.RequestResponseOption) (style files_sdk.Style, err error)
func Find(params files_sdk.StyleFindParams, opts ...files_sdk.RequestResponseOption) (style files_sdk.Style, err error)
```

Parameters, `files_sdk.StyleFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Path` | `string` | `path`, in the path | yes | Style path. |

Example:

```go
result, err := style.Find(files_sdk.StyleFindParams{
	Path: "path",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Update

`PATCH /styles/{path}`

Update Style

```go
func (c *Client) Update(params files_sdk.StyleUpdateParams, opts ...files_sdk.RequestResponseOption) (style files_sdk.Style, err error)
func Update(params files_sdk.StyleUpdateParams, opts ...files_sdk.RequestResponseOption) (style files_sdk.Style, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (style files_sdk.Style, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (style files_sdk.Style, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.StyleUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Path` | `string` | `path`, in the path | yes | Style path. |
| `File` | `io.Writer` | `file` | no | Logo for custom branding. Required when creating a new style. |
| `FileAttachment` | `*lib.JSONAttachment` | `file` | no | Sends file content as file; see lib.JSONAttachment. Leave File nil when you set it. |
| `LogoClickHref` | `string` | `logo_click_href` | no | URL to open when a public visitor clicks the logo. |

Example:

```go
result, err := style.Update(files_sdk.StyleUpdateParams{
	Path: "path",
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := style.UpdateWithMap(map[string]interface{}{
	"path": "path",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /styles/{path}`

Delete Style

```go
func (c *Client) Delete(params files_sdk.StyleDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.StyleDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.StyleDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Path` | `string` | `path`, in the path | yes | Style path. |

Example:

```go
err := style.Delete(files_sdk.StyleDeleteParams{
	Path: "path",
})
if err != nil {
	return err
}
```

## Model

`files_sdk.Style` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Style ID |
| `Path` | `string` | `path` | Folder path. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters. |
| `Logo` | `files_sdk.Image` | `logo` | Logo |
| `LogoClickHref` | `string` | `logo_click_href` | URL to open when a public visitor clicks the logo |
| `Thumbnail` | `files_sdk.Image` | `thumbnail` | Logo thumbnail |
| `File` | `io.Reader` | `file` | Logo for custom branding. Required when creating a new style. |
