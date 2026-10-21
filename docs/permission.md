# Permission

`files_sdk.Permission` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/permission`, named `permission`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/permission"
)
```

## Operations

### List

`GET /permissions`

List Permissions

```go
func (c *Client) List(params files_sdk.PermissionListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.PermissionListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.Permission` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.PermissionListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `site_id`, `group_id`, `path`, `user_id`, `partner_id` or `id`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `path`, `group_id`, `partner_id` or `user_id`. Valid field combinations are `[ group_id, path ]`, `[ partner_id, path ]`, `[ user_id, path ]`, `[ user_id, group_id ]`, `[ user_id, group_id, path ]`, `[ user_id, group_id, partner_id ]` or `[ user_id, group_id, partner_id, path ]`. |
| `FilterPrefix` | `interface{}` | `filter_prefix` | no | If set, return records where the specified field is prefixed by the supplied value. Valid fields are `path`. |
| `Path` | `string` | `path` | no | Permission path. If provided, will scope all permissions(including upward) to this path. |
| `IncludeGroups` | `*bool` | `include_groups` | no | If searching by user or group, also include user's permissions that are inherited from its groups? |
| `GroupId` | `string` | `group_id` | no |  |
| `PartnerId` | `string` | `partner_id` | no |  |
| `UserId` | `string` | `user_id` | no |  |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := permission.List(files_sdk.PermissionListParams{})
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

`POST /permissions`

Create Permission

```go
func (c *Client) Create(params files_sdk.PermissionCreateParams, opts ...files_sdk.RequestResponseOption) (permission files_sdk.Permission, err error)
func Create(params files_sdk.PermissionCreateParams, opts ...files_sdk.RequestResponseOption) (permission files_sdk.Permission, err error)
```

Parameters, `files_sdk.PermissionCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Path` | `string` | `path` | yes | Folder path |
| `GroupId` | `int64` | `group_id` | no | Group ID. Provide `group_name` or `group_id` |
| `GroupIds` | `string` | `group_ids` | no | Group IDs when the Permission requires membership in every listed group. If sent as a string, it should be comma-delimited. |
| `Permission` | `string` | `permission` | no | Permission type. Can be `admin`, `full`, `readonly`, `writeonly`, `list`, or `history` |
| `Recursive` | `*bool` | `recursive` | no | Apply to subfolders recursively? Must be true for `admin` Permissions. |
| `PartnerId` | `int64` | `partner_id` | no | Partner ID if this Permission belongs to a partner. |
| `UserId` | `int64` | `user_id` | no | User ID. Provide `username` or `user_id` |
| `Username` | `string` | `username` | no | User username. Provide `username` or `user_id` |
| `GroupName` | `string` | `group_name` | no | Group name. Provide `group_name` or `group_id` |
| `SiteId` | `int64` | `site_id` | no | Site ID. If not provided, will default to current site. Used when creating a permission for a child site. |

Example:

```go
result, err := permission.Create(files_sdk.PermissionCreateParams{
	Path: "path",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /permissions/{id}`

Delete Permission

```go
func (c *Client) Delete(params files_sdk.PermissionDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.PermissionDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.PermissionDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Permission ID. |

Example:

```go
err := permission.Delete(files_sdk.PermissionDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.Permission` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Permission ID |
| `Path` | `string` | `path` | Folder path. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters. |
| `UserId` | `int64` | `user_id` | User ID |
| `Username` | `string` | `username` | Username (if applicable) |
| `GroupId` | `int64` | `group_id` | Group ID |
| `GroupName` | `string` | `group_name` | Group name (if applicable) |
| `GroupIds` | `[]int64` | `group_ids` | Group IDs when this Permission requires membership in every listed group. |
| `GroupNames` | `[]string` | `group_names` | Group names when this permission requires multiple groups |
| `PartnerId` | `int64` | `partner_id` | Partner ID (if applicable) |
| `PartnerName` | `string` | `partner_name` | Partner name (if applicable) |
| `Permission` | `string` | `permission` | Permission type. See the table referenced in the documentation for an explanation of each permission. |
| `Recursive` | `*bool` | `recursive` | Whether this Permission applies to subfolders. Must be true for `admin` Permissions. |
| `SiteId` | `int64` | `site_id` | Site ID |
