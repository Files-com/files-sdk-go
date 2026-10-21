# MetadataCategory

`files_sdk.MetadataCategory` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/metadatacategory`, named `metadata_category`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	metadata_category "github.com/Files-com/files-sdk-go/v3/metadatacategory"
)
```

## Operations

### List

`GET /metadata_categories`

List Metadata Categories

```go
func (c *Client) List(params files_sdk.MetadataCategoryListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.MetadataCategoryListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.MetadataCategory` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.MetadataCategoryListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are . |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := metadata_category.List(files_sdk.MetadataCategoryListParams{})
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

`GET /metadata_categories/{id}`

Show Metadata Category

```go
func (c *Client) Find(params files_sdk.MetadataCategoryFindParams, opts ...files_sdk.RequestResponseOption) (metadataCategory files_sdk.MetadataCategory, err error)
func Find(params files_sdk.MetadataCategoryFindParams, opts ...files_sdk.RequestResponseOption) (metadataCategory files_sdk.MetadataCategory, err error)
```

Parameters, `files_sdk.MetadataCategoryFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Metadata Category ID. |

Example:

```go
result, err := metadata_category.Find(files_sdk.MetadataCategoryFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### ListFor

`GET /metadata_categories/list_by_path/{path}`

List Metadata Categories by Path

```go
func (c *Client) ListFor(params files_sdk.MetadataCategoryListForParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func ListFor(params files_sdk.MetadataCategoryListForParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.MetadataCategory` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.MetadataCategoryListForParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Path` | `string` | `path`, in the path | yes | Path to operate on. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := metadata_category.ListFor(files_sdk.MetadataCategoryListForParams{
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

`POST /metadata_categories`

Create Metadata Category

```go
func (c *Client) Create(params files_sdk.MetadataCategoryCreateParams, opts ...files_sdk.RequestResponseOption) (metadataCategory files_sdk.MetadataCategory, err error)
func Create(params files_sdk.MetadataCategoryCreateParams, opts ...files_sdk.RequestResponseOption) (metadataCategory files_sdk.MetadataCategory, err error)
```

Parameters, `files_sdk.MetadataCategoryCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Name` | `string` | `name` | yes | Name of the metadata category. |
| `DefaultColumns` | `[]string` | `default_columns` | no | Metadata keys that should appear as columns in the UI by default. |

Example:

```go
result, err := metadata_category.Create(files_sdk.MetadataCategoryCreateParams{
	Name: "Approval Workflow",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Update

`PATCH /metadata_categories/{id}`

Update Metadata Category

```go
func (c *Client) Update(params files_sdk.MetadataCategoryUpdateParams, opts ...files_sdk.RequestResponseOption) (metadataCategory files_sdk.MetadataCategory, err error)
func Update(params files_sdk.MetadataCategoryUpdateParams, opts ...files_sdk.RequestResponseOption) (metadataCategory files_sdk.MetadataCategory, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (metadataCategory files_sdk.MetadataCategory, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (metadataCategory files_sdk.MetadataCategory, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.MetadataCategoryUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Metadata Category ID. |
| `Name` | `string` | `name` | no | Name of the metadata category. |
| `DefaultColumns` | `[]string` | `default_columns` | no | Metadata keys that should appear as columns in the UI by default. |

Example:

```go
result, err := metadata_category.Update(files_sdk.MetadataCategoryUpdateParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := metadata_category.UpdateWithMap(map[string]interface{}{
	"id": 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /metadata_categories/{id}`

Delete Metadata Category

```go
func (c *Client) Delete(params files_sdk.MetadataCategoryDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.MetadataCategoryDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.MetadataCategoryDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Metadata Category ID. |

Example:

```go
err := metadata_category.Delete(files_sdk.MetadataCategoryDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.MetadataCategory` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Metadata Category ID |
| `Name` | `string` | `name` | Name of the metadata category. |
| `Definitions` | `map[string][]string` | `definitions` | Map of key names to arrays of allowed values. An empty array means free-form text. |
| `DefaultColumns` | `[]string` | `default_columns` | Metadata keys that should appear as columns in the UI by default. |
