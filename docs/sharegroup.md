# ShareGroup

`files_sdk.ShareGroup` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/sharegroup`, named `share_group`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	share_group "github.com/Files-com/files-sdk-go/v3/sharegroup"
)
```

## Operations

### List

`GET /share_groups`

List Share Groups

```go
func (c *Client) List(params files_sdk.ShareGroupListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.ShareGroupListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.ShareGroup` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.ShareGroupListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `UserId` | `int64` | `user_id` | no | User ID. Provide a value of `0` to operate the current session's user. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := share_group.List(files_sdk.ShareGroupListParams{})
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

`GET /share_groups/{id}`

Show Share Group

```go
func (c *Client) Find(params files_sdk.ShareGroupFindParams, opts ...files_sdk.RequestResponseOption) (shareGroup files_sdk.ShareGroup, err error)
func Find(params files_sdk.ShareGroupFindParams, opts ...files_sdk.RequestResponseOption) (shareGroup files_sdk.ShareGroup, err error)
```

Parameters, `files_sdk.ShareGroupFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Share Group ID. |

Example:

```go
result, err := share_group.Find(files_sdk.ShareGroupFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Create

`POST /share_groups`

Create Share Group

```go
func (c *Client) Create(params files_sdk.ShareGroupCreateParams, opts ...files_sdk.RequestResponseOption) (shareGroup files_sdk.ShareGroup, err error)
func Create(params files_sdk.ShareGroupCreateParams, opts ...files_sdk.RequestResponseOption) (shareGroup files_sdk.ShareGroup, err error)
```

Parameters, `files_sdk.ShareGroupCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `UserId` | `int64` | `user_id` | no | User ID. Provide a value of `0` to operate the current session's user. |
| `Notes` | `string` | `notes` | no | Additional notes of the share group |
| `Name` | `string` | `name` | yes | Name of the share group |
| `Members` | `[]map[string]interface{}` | `members` | yes | A list of share group members. |

Example:

```go
// Set this required value before the call.
var members []map[string]interface{}
result, err := share_group.Create(files_sdk.ShareGroupCreateParams{
	Name:    "Test group 1",
	Members: members,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Update

`PATCH /share_groups/{id}`

Update Share Group

```go
func (c *Client) Update(params files_sdk.ShareGroupUpdateParams, opts ...files_sdk.RequestResponseOption) (shareGroup files_sdk.ShareGroup, err error)
func Update(params files_sdk.ShareGroupUpdateParams, opts ...files_sdk.RequestResponseOption) (shareGroup files_sdk.ShareGroup, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (shareGroup files_sdk.ShareGroup, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (shareGroup files_sdk.ShareGroup, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.ShareGroupUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Share Group ID. |
| `Notes` | `string` | `notes` | no | Additional notes of the share group |
| `Name` | `string` | `name` | no | Name of the share group |
| `Members` | `[]map[string]interface{}` | `members` | no | A list of share group members. |

Example:

```go
result, err := share_group.Update(files_sdk.ShareGroupUpdateParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := share_group.UpdateWithMap(map[string]interface{}{
	"id": 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /share_groups/{id}`

Delete Share Group

```go
func (c *Client) Delete(params files_sdk.ShareGroupDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.ShareGroupDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.ShareGroupDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Share Group ID. |

Example:

```go
err := share_group.Delete(files_sdk.ShareGroupDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.ShareGroup` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Share Group ID |
| `Name` | `string` | `name` | Name of the share group |
| `Notes` | `string` | `notes` | Additional notes of the share group |
| `UserId` | `int64` | `user_id` | Owner User ID |
| `Members` | `[]files_sdk.ShareGroupMember` | `members` | A list of share group members |
