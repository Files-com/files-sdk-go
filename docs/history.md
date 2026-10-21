# History

`files_sdk.History` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/history`, named `history`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/history"
)
```

## Operations

### ListForFile

`GET /history/files/{path}`

List history for specific file

```go
func (c *Client) ListForFile(params files_sdk.HistoryListForFileParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func ListForFile(params files_sdk.HistoryListForFileParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.Action` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.HistoryListForFileParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `StartAt` | `*time.Time` | `start_at` | no | Leave blank or set to a date/time to filter earlier entries. |
| `EndAt` | `*time.Time` | `end_at` | no | Leave blank or set to a date/time to filter later entries. |
| `Display` | `string` | `display` | no | Display format. Leave blank or set to `full` or `parent`. |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `created_at`. |
| `Path` | `string` | `path`, in the path | yes | Path to operate on. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := history.ListForFile(files_sdk.HistoryListForFileParams{
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

### ListForFolder

`GET /history/folders/{path}`

List history for specific folder

```go
func (c *Client) ListForFolder(params files_sdk.HistoryListForFolderParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func ListForFolder(params files_sdk.HistoryListForFolderParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.Action` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.HistoryListForFolderParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `StartAt` | `*time.Time` | `start_at` | no | Leave blank or set to a date/time to filter earlier entries. |
| `EndAt` | `*time.Time` | `end_at` | no | Leave blank or set to a date/time to filter later entries. |
| `Display` | `string` | `display` | no | Display format. Leave blank or set to `full` or `parent`. |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `created_at`. |
| `Path` | `string` | `path`, in the path | yes | Path to operate on. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := history.ListForFolder(files_sdk.HistoryListForFolderParams{
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

### ListForUser

`GET /history/users/{user_id}`

List history for specific user

```go
func (c *Client) ListForUser(params files_sdk.HistoryListForUserParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func ListForUser(params files_sdk.HistoryListForUserParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.Action` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.HistoryListForUserParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `StartAt` | `*time.Time` | `start_at` | no | Leave blank or set to a date/time to filter earlier entries. |
| `EndAt` | `*time.Time` | `end_at` | no | Leave blank or set to a date/time to filter later entries. |
| `Display` | `string` | `display` | no | Display format. Leave blank or set to `full` or `parent`. |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `created_at`. |
| `UserId` | `int64` | `user_id`, in the path | yes | User ID. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := history.ListForUser(files_sdk.HistoryListForUserParams{
	UserId: 1,
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

### ListLogins

`GET /history/login`

List site login history

```go
func (c *Client) ListLogins(params files_sdk.HistoryListLoginsParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func ListLogins(params files_sdk.HistoryListLoginsParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.Action` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.HistoryListLoginsParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `StartAt` | `*time.Time` | `start_at` | no | Leave blank or set to a date/time to filter earlier entries. |
| `EndAt` | `*time.Time` | `end_at` | no | Leave blank or set to a date/time to filter later entries. |
| `Display` | `string` | `display` | no | Display format. Leave blank or set to `full` or `parent`. |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `created_at`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := history.ListLogins(files_sdk.HistoryListLoginsParams{})
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

### List

`GET /history`

List site full action history

```go
func (c *Client) List(params files_sdk.HistoryListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.HistoryListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.Action` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.HistoryListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `StartAt` | `*time.Time` | `start_at` | no | Leave blank or set to a date/time to filter earlier entries. |
| `EndAt` | `*time.Time` | `end_at` | no | Leave blank or set to a date/time to filter later entries. |
| `Display` | `string` | `display` | no | Display format. Leave blank or set to `full` or `parent`. |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `created_at`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `user_id`, `folder` or `path`. Valid field combinations are `[ ]`, `[ path ]`, `[ path ]` or `[ path ]`. |
| `FilterPrefix` | `interface{}` | `filter_prefix` | no | If set, return records where the specified field is prefixed by the supplied value. Valid fields are `path`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := history.List(files_sdk.HistoryListParams{})
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

## Model

`files_sdk.History` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Action ID |
| `Path` | `string` | `path` | Path. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters. |
| `When` | `*time.Time` | `when` | Action occurrence date/time |
| `Destination` | `string` | `destination` | The destination path for this action, if applicable |
| `Display` | `string` | `display` | Friendly displayed output |
| `Ip` | `string` | `ip` | IP Address that performed this action |
| `Source` | `string` | `source` | The source path for this action, if applicable |
| `Targets` | `interface{}` | `targets` | Targets |
| `UserId` | `int64` | `user_id` | User ID |
| `Username` | `string` | `username` | Username |
| `UserIsFromParentSite` | `*bool` | `user_is_from_parent_site` | true if this change was performed by a user on a parent site. |
| `Action` | `string` | `action` | Type of action |
| `FailureType` | `string` | `failure_type` | Failure type. If action was a user login or session failure, why did it fail? |
| `Interface` | `string` | `interface` | Interface on which this action occurred. |
