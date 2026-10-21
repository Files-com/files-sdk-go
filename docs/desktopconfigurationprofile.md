# DesktopConfigurationProfile

`files_sdk.DesktopConfigurationProfile` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/desktopconfigurationprofile`, named `desktop_configuration_profile`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	desktop_configuration_profile "github.com/Files-com/files-sdk-go/v3/desktopconfigurationprofile"
)
```

## Operations

### List

`GET /desktop_configuration_profiles`

List Desktop Configuration Profiles

```go
func (c *Client) List(params files_sdk.DesktopConfigurationProfileListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.DesktopConfigurationProfileListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.DesktopConfigurationProfile` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.DesktopConfigurationProfileListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id` and `name`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `workspace_id`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := desktop_configuration_profile.List(files_sdk.DesktopConfigurationProfileListParams{})
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

`GET /desktop_configuration_profiles/{id}`

Show Desktop Configuration Profile

```go
func (c *Client) Find(params files_sdk.DesktopConfigurationProfileFindParams, opts ...files_sdk.RequestResponseOption) (desktopConfigurationProfile files_sdk.DesktopConfigurationProfile, err error)
func Find(params files_sdk.DesktopConfigurationProfileFindParams, opts ...files_sdk.RequestResponseOption) (desktopConfigurationProfile files_sdk.DesktopConfigurationProfile, err error)
```

Parameters, `files_sdk.DesktopConfigurationProfileFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Desktop Configuration Profile ID. |

Example:

```go
result, err := desktop_configuration_profile.Find(files_sdk.DesktopConfigurationProfileFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Create

`POST /desktop_configuration_profiles`

Create Desktop Configuration Profile

```go
func (c *Client) Create(params files_sdk.DesktopConfigurationProfileCreateParams, opts ...files_sdk.RequestResponseOption) (desktopConfigurationProfile files_sdk.DesktopConfigurationProfile, err error)
func Create(params files_sdk.DesktopConfigurationProfileCreateParams, opts ...files_sdk.RequestResponseOption) (desktopConfigurationProfile files_sdk.DesktopConfigurationProfile, err error)
```

Parameters, `files_sdk.DesktopConfigurationProfileCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Name` | `string` | `name` | yes | Profile name |
| `MountMappings` | `interface{}` | `mount_mappings` | yes | Mount point mappings for the desktop app. Keys must be a single uppercase Windows drive letter other than A, B, or C, and values are Files.com paths to mount there. |
| `WorkspaceId` | `int64` | `workspace_id` | no | Workspace ID |
| `UseForAllUsers` | `*bool` | `use_for_all_users` | no | Whether this profile applies to all users in the Workspace by default |
| `DisableDriveMounting` | `*bool` | `disable_drive_mounting` | no | Whether the desktop app should hide drive mounting, prevent new drive mounts, and unmount active drive mounts for users with this profile |

Example:

```go
// Set this required value before the call.
var mountMappings interface{}
result, err := desktop_configuration_profile.Create(files_sdk.DesktopConfigurationProfileCreateParams{
	Name:          "North America Desktop Profile",
	MountMappings: mountMappings,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Update

`PATCH /desktop_configuration_profiles/{id}`

Update Desktop Configuration Profile

```go
func (c *Client) Update(params files_sdk.DesktopConfigurationProfileUpdateParams, opts ...files_sdk.RequestResponseOption) (desktopConfigurationProfile files_sdk.DesktopConfigurationProfile, err error)
func Update(params files_sdk.DesktopConfigurationProfileUpdateParams, opts ...files_sdk.RequestResponseOption) (desktopConfigurationProfile files_sdk.DesktopConfigurationProfile, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (desktopConfigurationProfile files_sdk.DesktopConfigurationProfile, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (desktopConfigurationProfile files_sdk.DesktopConfigurationProfile, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.DesktopConfigurationProfileUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Desktop Configuration Profile ID. |
| `Name` | `string` | `name` | no | Profile name |
| `WorkspaceId` | `int64` | `workspace_id` | no | Workspace ID |
| `MountMappings` | `interface{}` | `mount_mappings` | no | Mount point mappings for the desktop app. Keys must be a single uppercase Windows drive letter other than A, B, or C, and values are Files.com paths to mount there. |
| `UseForAllUsers` | `*bool` | `use_for_all_users` | no | Whether this profile applies to all users in the Workspace by default |
| `DisableDriveMounting` | `*bool` | `disable_drive_mounting` | no | Whether the desktop app should hide drive mounting, prevent new drive mounts, and unmount active drive mounts for users with this profile |

Example:

```go
result, err := desktop_configuration_profile.Update(files_sdk.DesktopConfigurationProfileUpdateParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := desktop_configuration_profile.UpdateWithMap(map[string]interface{}{
	"id": 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /desktop_configuration_profiles/{id}`

Delete Desktop Configuration Profile

```go
func (c *Client) Delete(params files_sdk.DesktopConfigurationProfileDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.DesktopConfigurationProfileDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.DesktopConfigurationProfileDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Desktop Configuration Profile ID. |

Example:

```go
err := desktop_configuration_profile.Delete(files_sdk.DesktopConfigurationProfileDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.DesktopConfigurationProfile` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Desktop Configuration Profile ID |
| `Name` | `string` | `name` | Profile name |
| `WorkspaceId` | `int64` | `workspace_id` | Workspace ID |
| `UseForAllUsers` | `*bool` | `use_for_all_users` | Whether this profile applies to all users in the Workspace by default |
| `DisableDriveMounting` | `*bool` | `disable_drive_mounting` | Whether the desktop app should hide drive mounting, prevent new drive mounts, and unmount active drive mounts for users with this profile |
| `MountMappings` | `interface{}` | `mount_mappings` | Mount point mappings for the desktop app. Keys must be a single uppercase Windows drive letter other than A, B, or C, and values are Files.com paths to mount there. |
