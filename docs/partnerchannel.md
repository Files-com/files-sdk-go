# PartnerChannel

`files_sdk.PartnerChannel` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/partnerchannel`, named `partner_channel`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	partner_channel "github.com/Files-com/files-sdk-go/v3/partnerchannel"
)
```

## Operations

### List

`GET /partner_channels`

List Partner Channels

```go
func (c *Client) List(params files_sdk.PartnerChannelListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.PartnerChannelListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.PartnerChannel` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.PartnerChannelListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id`, `path` or `partner_id`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `partner_id` and `workspace_id`. Valid field combinations are `[ workspace_id, partner_id ]`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := partner_channel.List(files_sdk.PartnerChannelListParams{})
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

`GET /partner_channels/{id}`

Show Partner Channel

```go
func (c *Client) Find(params files_sdk.PartnerChannelFindParams, opts ...files_sdk.RequestResponseOption) (partnerChannel files_sdk.PartnerChannel, err error)
func Find(params files_sdk.PartnerChannelFindParams, opts ...files_sdk.RequestResponseOption) (partnerChannel files_sdk.PartnerChannel, err error)
```

Parameters, `files_sdk.PartnerChannelFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Partner Channel ID. |

Example:

```go
result, err := partner_channel.Find(files_sdk.PartnerChannelFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Create

`POST /partner_channels`

Create Partner Channel

```go
func (c *Client) Create(params files_sdk.PartnerChannelCreateParams, opts ...files_sdk.RequestResponseOption) (partnerChannel files_sdk.PartnerChannel, err error)
func Create(params files_sdk.PartnerChannelCreateParams, opts ...files_sdk.RequestResponseOption) (partnerChannel files_sdk.PartnerChannel, err error)
```

Parameters, `files_sdk.PartnerChannelCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Direction` | `files_sdk.PartnerChannelDirectionEnum` | `direction` | no | Channel directions. `two_way` enables both directions, `to_partner` enables outgoing downloads, and `from_partner` enables incoming uploads. Values: `two_way`, `to_partner`, `from_partner`. |
| `UseChannelRoot` | `*bool` | `use_channel_root` | no | Use the Channel folder directly for a one-way exchange. Defaults to false. Cannot be changed after creation. Folder name overrides must be blank when enabled, and the Channel must remain one-way. |
| `FromPartnerFolderName` | `string` | `from_partner_folder_name` | no | Optional Channel-level from-Partner folder name override. |
| `FromPartnerManagedFolderPaths` | `[]string` | `from_partner_managed_folder_paths` | no | Managed folder paths inside the from-Partner folder. |
| `FromPartnerRoutePath` | `string` | `from_partner_route_path` | no | Optional route path for files uploaded by the Partner. |
| `ToPartnerFolderName` | `string` | `to_partner_folder_name` | no | Optional Channel-level to-Partner folder name override. |
| `ToPartnerManagedFolderPaths` | `[]string` | `to_partner_managed_folder_paths` | no | Managed folder paths inside the to-Partner folder. |
| `ToPartnerRoutePath` | `string` | `to_partner_route_path` | no | Optional route path for files delivered to the Partner. |
| `PartnerId` | `int64` | `partner_id` | yes | ID of the Partner this Channel belongs to. |
| `Path` | `string` | `path` | yes | Channel path relative to the Partner root folder. |
| `WorkspaceId` | `int64` | `workspace_id` | no | ID of the Workspace associated with this Partner Channel. |

Example:

```go
result, err := partner_channel.Create(files_sdk.PartnerChannelCreateParams{
	PartnerId: 1,
	Path:      "claims/medical",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Update

`PATCH /partner_channels/{id}`

Update Partner Channel

```go
func (c *Client) Update(params files_sdk.PartnerChannelUpdateParams, opts ...files_sdk.RequestResponseOption) (partnerChannel files_sdk.PartnerChannel, err error)
func Update(params files_sdk.PartnerChannelUpdateParams, opts ...files_sdk.RequestResponseOption) (partnerChannel files_sdk.PartnerChannel, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (partnerChannel files_sdk.PartnerChannel, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (partnerChannel files_sdk.PartnerChannel, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.PartnerChannelUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Partner Channel ID. |
| `Direction` | `files_sdk.PartnerChannelDirectionEnum` | `direction` | no | Channel directions. `two_way` enables both directions, `to_partner` enables outgoing downloads, and `from_partner` enables incoming uploads. Values: `two_way`, `to_partner`, `from_partner`. |
| `UseChannelRoot` | `*bool` | `use_channel_root` | no | Use the Channel folder directly for a one-way exchange. Defaults to false. Cannot be changed after creation. Folder name overrides must be blank when enabled, and the Channel must remain one-way. |
| `FromPartnerFolderName` | `string` | `from_partner_folder_name` | no | Optional Channel-level from-Partner folder name override. |
| `FromPartnerManagedFolderPaths` | `[]string` | `from_partner_managed_folder_paths` | no | Managed folder paths inside the from-Partner folder. |
| `FromPartnerRoutePath` | `string` | `from_partner_route_path` | no | Optional route path for files uploaded by the Partner. |
| `ToPartnerFolderName` | `string` | `to_partner_folder_name` | no | Optional Channel-level to-Partner folder name override. |
| `ToPartnerManagedFolderPaths` | `[]string` | `to_partner_managed_folder_paths` | no | Managed folder paths inside the to-Partner folder. |
| `ToPartnerRoutePath` | `string` | `to_partner_route_path` | no | Optional route path for files delivered to the Partner. |
| `Path` | `string` | `path` | no | Channel path relative to the Partner root folder. |

Example:

```go
result, err := partner_channel.Update(files_sdk.PartnerChannelUpdateParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := partner_channel.UpdateWithMap(map[string]interface{}{
	"id": 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /partner_channels/{id}`

Delete Partner Channel

```go
func (c *Client) Delete(params files_sdk.PartnerChannelDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.PartnerChannelDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.PartnerChannelDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Partner Channel ID. |

Example:

```go
err := partner_channel.Delete(files_sdk.PartnerChannelDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.PartnerChannel` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | The unique ID of the Partner Channel. |
| `WorkspaceId` | `int64` | `workspace_id` | ID of the Workspace associated with this Partner Channel. |
| `Direction` | `string` | `direction` | Channel directions. `two_way` enables both directions, `to_partner` enables outgoing downloads, and `from_partner` enables incoming uploads. |
| `UseChannelRoot` | `*bool` | `use_channel_root` | Use the Channel folder directly for a one-way exchange. Defaults to false. Cannot be changed after creation. Folder name overrides must be blank when enabled, and the Channel must remain one-way. |
| `PartnerId` | `int64` | `partner_id` | ID of the Partner this Channel belongs to. |
| `PartnerChannelTemplateId` | `int64` | `partner_channel_template_id` | ID of the Partner Channel Template that manages this Channel, if any. |
| `Path` | `string` | `path` | Channel path relative to the Partner root folder. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters. |
| `ToPartnerFolderName` | `string` | `to_partner_folder_name` | Optional Channel-level to-Partner folder name override. |
| `FromPartnerFolderName` | `string` | `from_partner_folder_name` | Optional Channel-level from-Partner folder name override. |
| `FromPartnerRoutePath` | `string` | `from_partner_route_path` | Optional route path for files uploaded by the Partner. |
| `ToPartnerRoutePath` | `string` | `to_partner_route_path` | Optional route path for files delivered to the Partner. |
| `ToPartnerManagedFolderPaths` | `[]string` | `to_partner_managed_folder_paths` | Managed folder paths inside the to-Partner folder. |
| `FromPartnerManagedFolderPaths` | `[]string` | `from_partner_managed_folder_paths` | Managed folder paths inside the from-Partner folder. |
| `EffectiveToPartnerFolderName` | `string` | `effective_to_partner_folder_name` | Resolved to-Partner subfolder name. Null when the direction is disabled or uses the Channel folder directly. |
| `EffectiveFromPartnerFolderName` | `string` | `effective_from_partner_folder_name` | Resolved from-Partner subfolder name. Null when the direction is disabled or uses the Channel folder directly. |
| `ChannelPath` | `string` | `channel_path` | Resolved Channel folder path. |
| `ToPartnerFolderPath` | `string` | `to_partner_folder_path` | Resolved to-Partner folder path. |
| `FromPartnerFolderPath` | `string` | `from_partner_folder_path` | Resolved from-Partner folder path. |
