# GroupUser

`files_sdk.GroupUser` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/groupuser`, named `group_user`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	group_user "github.com/Files-com/files-sdk-go/v3/groupuser"
)
```

## Operations

### List

`GET /group_users`

List Group Users

```go
func (c *Client) List(params files_sdk.GroupUserListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.GroupUserListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.GroupUser` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.GroupUserListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `GroupId` | `int64` | `group_id` | no | Group ID. If provided, returns memberships of this group. Requires a Site Administrator, a Read-only Administrator, a Workspace Administrator for the group's workspace, or a Group Administrator of this group. |
| `UserId` | `int64` | `user_id` | no | User ID. If provided, will return group_users of this user. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := group_user.List(files_sdk.GroupUserListParams{})
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

`POST /group_users`

Create Group User

```go
func (c *Client) Create(params files_sdk.GroupUserCreateParams, opts ...files_sdk.RequestResponseOption) (groupUser files_sdk.GroupUser, err error)
func Create(params files_sdk.GroupUserCreateParams, opts ...files_sdk.RequestResponseOption) (groupUser files_sdk.GroupUser, err error)
```

Parameters, `files_sdk.GroupUserCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `GroupId` | `int64` | `group_id` | yes | Group ID to add user to. |
| `UserId` | `int64` | `user_id` | yes | User ID to add to group. |
| `Admin` | `*bool` | `admin` | no | Is the user a group administrator? |

Example:

```go
result, err := group_user.Create(files_sdk.GroupUserCreateParams{
	GroupId: 1,
	UserId:  1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Update

`PATCH /group_users/{id}`

Update Group User

```go
func (c *Client) Update(params files_sdk.GroupUserUpdateParams, opts ...files_sdk.RequestResponseOption) (groupUser files_sdk.GroupUser, err error)
func Update(params files_sdk.GroupUserUpdateParams, opts ...files_sdk.RequestResponseOption) (groupUser files_sdk.GroupUser, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (groupUser files_sdk.GroupUser, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (groupUser files_sdk.GroupUser, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.GroupUserUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Group User ID. |
| `GroupId` | `int64` | `group_id` | yes | Group ID to add user to. |
| `UserId` | `int64` | `user_id` | yes | User ID to add to group. |
| `Admin` | `*bool` | `admin` | no | Is the user a group administrator? |

Example:

```go
result, err := group_user.Update(files_sdk.GroupUserUpdateParams{
	Id:      1,
	GroupId: 1,
	UserId:  1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := group_user.UpdateWithMap(map[string]interface{}{
	"id":       1,
	"group_id": 1,
	"user_id":  1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /group_users/{id}`

Delete Group User

```go
func (c *Client) Delete(params files_sdk.GroupUserDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.GroupUserDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.GroupUserDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Group User ID. |
| `GroupId` | `int64` | `group_id` | yes | Group ID from which to remove user. |
| `UserId` | `int64` | `user_id` | yes | User ID to remove from group. |

Example:

```go
err := group_user.Delete(files_sdk.GroupUserDeleteParams{
	Id:      1,
	GroupId: 1,
	UserId:  1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.GroupUser` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `GroupName` | `string` | `group_name` | Group name |
| `GroupId` | `int64` | `group_id` | Group ID |
| `UserId` | `int64` | `user_id` | User ID |
| `Admin` | `*bool` | `admin` | Is this user an administrator of this group? |
| `Username` | `string` | `username` | Username of the user |
| `Id` | `int64` | `id` | Group User ID. |
