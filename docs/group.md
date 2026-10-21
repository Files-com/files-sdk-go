# Group

`files_sdk.Group` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/group`, named `group`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/group"
)
```

## Operations

### List

`GET /groups`

List Groups

```go
func (c *Client) List(params files_sdk.GroupListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.GroupListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.Group` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.GroupListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `site_id`, `workspace_id` or `name`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `name` and `workspace_id`. Valid field combinations are `[ workspace_id, name ]`. |
| `FilterPrefix` | `interface{}` | `filter_prefix` | no | If set, return records where the specified field is prefixed by the supplied value. Valid fields are `name`. |
| `Ids` | `string` | `ids` | no | Comma-separated list of group ids to include in results. |
| `IncludeParentSiteGroups` | `*bool` | `include_parent_site_groups` | no | Include groups from the parent site. |
| `IncludeDefaultWorkspaceGroups` | `*bool` | `include_default_workspace_groups` | no | When scoped to a Custom Workspace, include Default Workspace groups that have permissions in that Workspace. Available to Site Administrators, and to Workspace Administrators and Folder Admins in that Workspace when the site shares its Default Workspace directory with Workspaces. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := group.List(files_sdk.GroupListParams{})
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

`GET /groups/{id}`

Show Group

```go
func (c *Client) Find(params files_sdk.GroupFindParams, opts ...files_sdk.RequestResponseOption) (group files_sdk.Group, err error)
func Find(params files_sdk.GroupFindParams, opts ...files_sdk.RequestResponseOption) (group files_sdk.Group, err error)
```

Parameters, `files_sdk.GroupFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Group ID. |

Example:

```go
result, err := group.Find(files_sdk.GroupFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Create

`POST /groups`

Create Group

```go
func (c *Client) Create(params files_sdk.GroupCreateParams, opts ...files_sdk.RequestResponseOption) (group files_sdk.Group, err error)
func Create(params files_sdk.GroupCreateParams, opts ...files_sdk.RequestResponseOption) (group files_sdk.Group, err error)
```

Parameters, `files_sdk.GroupCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Notes` | `string` | `notes` | no | Group notes. |
| `UserIds` | `string` | `user_ids` | no | A list of user ids. If sent as a string, should be comma-delimited. |
| `AdminIds` | `string` | `admin_ids` | no | A list of group admin user ids. If sent as a string, should be comma-delimited. |
| `AiAssistantPersonalityId` | `int64` | `ai_assistant_personality_id` | no | AI Assistant Personality ID assigned to this Group, if any. Users in the Group inherit it unless a direct per-user or Partner assignment overrides it. |
| `FtpPermission` | `*bool` | `ftp_permission` | no | If true, users in this group can use FTP to login. This will override a false value of `ftp_permission` on the user level. |
| `SftpPermission` | `*bool` | `sftp_permission` | no | If true, users in this group can use SFTP to login. This will override a false value of `sftp_permission` on the user level. |
| `DavPermission` | `*bool` | `dav_permission` | no | If true, users in this group can use WebDAV to login. This will override a false value of `dav_permission` on the user level. |
| `RestapiPermission` | `*bool` | `restapi_permission` | no | If true, users in this group can use the REST API to login. This will override a false value of `restapi_permission` on the user level. |
| `S3CompatibleEndpointPermission` | `*bool` | `s3_compatible_endpoint_permission` | no | If true, users in this group can access the S3-compatible endpoint. This will override a false value of `s3_compatible_endpoint_permission` on the user level. Defaults to false. |
| `DesktopConfigurationProfileId` | `int64` | `desktop_configuration_profile_id` | no | Desktop Configuration Profile ID assigned to this Group, if any. Users in the Group inherit it unless a direct per-user assignment overrides it. |
| `IntegrationCentricProfileId` | `int64` | `integration_centric_profile_id` | no | Integration Centric Profile ID assigned to this Group, if any. Users in the Group inherit it unless a direct per-user assignment overrides it. |
| `AllowedIps` | `string` | `allowed_ips` | no | A list of allowed IPs if applicable. Newline delimited |
| `Name` | `string` | `name` | yes | Group name. |
| `WorkspaceId` | `int64` | `workspace_id` | no | Workspace ID |

Example:

```go
result, err := group.Create(files_sdk.GroupCreateParams{
	Name: "name",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Update

`PATCH /groups/{id}`

Update Group

```go
func (c *Client) Update(params files_sdk.GroupUpdateParams, opts ...files_sdk.RequestResponseOption) (group files_sdk.Group, err error)
func Update(params files_sdk.GroupUpdateParams, opts ...files_sdk.RequestResponseOption) (group files_sdk.Group, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (group files_sdk.Group, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (group files_sdk.Group, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.GroupUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Group ID. |
| `Notes` | `string` | `notes` | no | Group notes. |
| `UserIds` | `string` | `user_ids` | no | A list of user ids. If sent as a string, should be comma-delimited. |
| `AdminIds` | `string` | `admin_ids` | no | A list of group admin user ids. If sent as a string, should be comma-delimited. |
| `AiAssistantPersonalityId` | `int64` | `ai_assistant_personality_id` | no | AI Assistant Personality ID assigned to this Group, if any. Users in the Group inherit it unless a direct per-user or Partner assignment overrides it. |
| `FtpPermission` | `*bool` | `ftp_permission` | no | If true, users in this group can use FTP to login. This will override a false value of `ftp_permission` on the user level. |
| `SftpPermission` | `*bool` | `sftp_permission` | no | If true, users in this group can use SFTP to login. This will override a false value of `sftp_permission` on the user level. |
| `DavPermission` | `*bool` | `dav_permission` | no | If true, users in this group can use WebDAV to login. This will override a false value of `dav_permission` on the user level. |
| `RestapiPermission` | `*bool` | `restapi_permission` | no | If true, users in this group can use the REST API to login. This will override a false value of `restapi_permission` on the user level. |
| `S3CompatibleEndpointPermission` | `*bool` | `s3_compatible_endpoint_permission` | no | If true, users in this group can access the S3-compatible endpoint. This will override a false value of `s3_compatible_endpoint_permission` on the user level. Defaults to false. |
| `DesktopConfigurationProfileId` | `int64` | `desktop_configuration_profile_id` | no | Desktop Configuration Profile ID assigned to this Group, if any. Users in the Group inherit it unless a direct per-user assignment overrides it. |
| `IntegrationCentricProfileId` | `int64` | `integration_centric_profile_id` | no | Integration Centric Profile ID assigned to this Group, if any. Users in the Group inherit it unless a direct per-user assignment overrides it. |
| `AllowedIps` | `string` | `allowed_ips` | no | A list of allowed IPs if applicable. Newline delimited |
| `Name` | `string` | `name` | no | Group name. |

Example:

```go
result, err := group.Update(files_sdk.GroupUpdateParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := group.UpdateWithMap(map[string]interface{}{
	"id": 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /groups/{id}`

Delete Group

```go
func (c *Client) Delete(params files_sdk.GroupDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.GroupDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.GroupDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Group ID. |

Example:

```go
err := group.Delete(files_sdk.GroupDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.Group` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Group ID |
| `Name` | `string` | `name` | Group name |
| `AllowedIps` | `string` | `allowed_ips` | A list of allowed IPs if applicable. Newline delimited |
| `AdminIds` | `string` | `admin_ids` | Comma-delimited list of user IDs who are group administrators (separated by commas) |
| `Notes` | `string` | `notes` | Notes about this group |
| `UserIds` | `string` | `user_ids` | Comma-delimited list of user IDs who belong to this group (separated by commas) |
| `Usernames` | `string` | `usernames` | Comma-delimited list of usernames who belong to this group (separated by commas) |
| `AiAssistantPersonalityId` | `int64` | `ai_assistant_personality_id` | AI Assistant Personality ID assigned to this Group, if any. Users in the Group inherit it unless a direct per-user or Partner assignment overrides it. |
| `FtpPermission` | `*bool` | `ftp_permission` | If true, users in this group can use FTP to login. This will override a false value of `ftp_permission` on the user level. |
| `SftpPermission` | `*bool` | `sftp_permission` | If true, users in this group can use SFTP to login. This will override a false value of `sftp_permission` on the user level. |
| `DavPermission` | `*bool` | `dav_permission` | If true, users in this group can use WebDAV to login. This will override a false value of `dav_permission` on the user level. |
| `RestapiPermission` | `*bool` | `restapi_permission` | If true, users in this group can use the REST API to login. This will override a false value of `restapi_permission` on the user level. |
| `S3CompatibleEndpointPermission` | `*bool` | `s3_compatible_endpoint_permission` | If true, users in this group can access the S3-compatible endpoint. This will override a false value of `s3_compatible_endpoint_permission` on the user level. Defaults to false. |
| `DesktopConfigurationProfileId` | `int64` | `desktop_configuration_profile_id` | Desktop Configuration Profile ID assigned to this Group, if any. Users in the Group inherit it unless a direct per-user assignment overrides it. |
| `IntegrationCentricProfileId` | `int64` | `integration_centric_profile_id` | Integration Centric Profile ID assigned to this Group, if any. Users in the Group inherit it unless a direct per-user assignment overrides it. |
| `SiteId` | `int64` | `site_id` | Site ID |
| `WorkspaceId` | `int64` | `workspace_id` | Workspace ID |
