# PartnerChannelTemplate

`files_sdk.PartnerChannelTemplate` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/partnerchanneltemplate`, named `partner_channel_template`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	partner_channel_template "github.com/Files-com/files-sdk-go/v3/partnerchanneltemplate"
)
```

## Operations

### List

`GET /partner_channel_templates`

List Partner Channel Templates

```go
func (c *Client) List(params files_sdk.PartnerChannelTemplateListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.PartnerChannelTemplateListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.PartnerChannelTemplate` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.PartnerChannelTemplateListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id` and `name`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `workspace_id`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := partner_channel_template.List(files_sdk.PartnerChannelTemplateListParams{})
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

`GET /partner_channel_templates/{id}`

Show Partner Channel Template

```go
func (c *Client) Find(params files_sdk.PartnerChannelTemplateFindParams, opts ...files_sdk.RequestResponseOption) (partnerChannelTemplate files_sdk.PartnerChannelTemplate, err error)
func Find(params files_sdk.PartnerChannelTemplateFindParams, opts ...files_sdk.RequestResponseOption) (partnerChannelTemplate files_sdk.PartnerChannelTemplate, err error)
```

Parameters, `files_sdk.PartnerChannelTemplateFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Partner Channel Template ID. |

Example:

```go
result, err := partner_channel_template.Find(files_sdk.PartnerChannelTemplateFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Create

`POST /partner_channel_templates`

Create Partner Channel Template

```go
func (c *Client) Create(params files_sdk.PartnerChannelTemplateCreateParams, opts ...files_sdk.RequestResponseOption) (partnerChannelTemplate files_sdk.PartnerChannelTemplate, err error)
func Create(params files_sdk.PartnerChannelTemplateCreateParams, opts ...files_sdk.RequestResponseOption) (partnerChannelTemplate files_sdk.PartnerChannelTemplate, err error)
```

Parameters, `files_sdk.PartnerChannelTemplateCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Direction` | `files_sdk.PartnerChannelTemplateDirectionEnum` | `direction` | no | Channel directions. `two_way` enables both directions, `to_partner` enables outgoing downloads, and `from_partner` enables incoming uploads. Values: `two_way`, `to_partner`, `from_partner`. |
| `UseChannelRoot` | `*bool` | `use_channel_root` | no | Use the Channel folder directly for a one-way exchange. Defaults to false. Cannot be changed after creation. Folder name overrides must be blank when enabled, and the Template must remain one-way. |
| `FromPartnerFolderName` | `string` | `from_partner_folder_name` | no | Optional Channel-level from-Partner folder name override. |
| `FromPartnerManagedFolderPaths` | `[]string` | `from_partner_managed_folder_paths` | no | Managed folder paths inside the from-Partner folder. |
| `FromPartnerRoutePathPattern` | `string` | `from_partner_route_path_pattern` | no | Optional route path pattern for files uploaded by the Partner. Supports {{partner_name}}. |
| `ToPartnerFolderName` | `string` | `to_partner_folder_name` | no | Optional Channel-level to-Partner folder name override. |
| `ToPartnerManagedFolderPaths` | `[]string` | `to_partner_managed_folder_paths` | no | Managed folder paths inside the to-Partner folder. |
| `ToPartnerRoutePathPattern` | `string` | `to_partner_route_path_pattern` | no | Optional route path pattern for files delivered to the Partner. Supports {{partner_name}}. |
| `Name` | `string` | `name` | yes | The name of the Partner Channel Template. |
| `Path` | `string` | `path` | yes | Channel path relative to the Partner root folder. |
| `WorkspaceId` | `int64` | `workspace_id` | no | ID of the Workspace associated with this Partner Channel Template. |

Example:

```go
result, err := partner_channel_template.Create(files_sdk.PartnerChannelTemplateCreateParams{
	Name: "Claims Template",
	Path: "claims/medical",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Update

`PATCH /partner_channel_templates/{id}`

Update Partner Channel Template

```go
func (c *Client) Update(params files_sdk.PartnerChannelTemplateUpdateParams, opts ...files_sdk.RequestResponseOption) (partnerChannelTemplate files_sdk.PartnerChannelTemplate, err error)
func Update(params files_sdk.PartnerChannelTemplateUpdateParams, opts ...files_sdk.RequestResponseOption) (partnerChannelTemplate files_sdk.PartnerChannelTemplate, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (partnerChannelTemplate files_sdk.PartnerChannelTemplate, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (partnerChannelTemplate files_sdk.PartnerChannelTemplate, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.PartnerChannelTemplateUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Partner Channel Template ID. |
| `Direction` | `files_sdk.PartnerChannelTemplateDirectionEnum` | `direction` | no | Channel directions. `two_way` enables both directions, `to_partner` enables outgoing downloads, and `from_partner` enables incoming uploads. Values: `two_way`, `to_partner`, `from_partner`. |
| `UseChannelRoot` | `*bool` | `use_channel_root` | no | Use the Channel folder directly for a one-way exchange. Defaults to false. Cannot be changed after creation. Folder name overrides must be blank when enabled, and the Template must remain one-way. |
| `FromPartnerFolderName` | `string` | `from_partner_folder_name` | no | Optional Channel-level from-Partner folder name override. |
| `FromPartnerManagedFolderPaths` | `[]string` | `from_partner_managed_folder_paths` | no | Managed folder paths inside the from-Partner folder. |
| `FromPartnerRoutePathPattern` | `string` | `from_partner_route_path_pattern` | no | Optional route path pattern for files uploaded by the Partner. Supports {{partner_name}}. |
| `ToPartnerFolderName` | `string` | `to_partner_folder_name` | no | Optional Channel-level to-Partner folder name override. |
| `ToPartnerManagedFolderPaths` | `[]string` | `to_partner_managed_folder_paths` | no | Managed folder paths inside the to-Partner folder. |
| `ToPartnerRoutePathPattern` | `string` | `to_partner_route_path_pattern` | no | Optional route path pattern for files delivered to the Partner. Supports {{partner_name}}. |
| `Name` | `string` | `name` | no | The name of the Partner Channel Template. |
| `Path` | `string` | `path` | no | Channel path relative to the Partner root folder. |

Example:

```go
result, err := partner_channel_template.Update(files_sdk.PartnerChannelTemplateUpdateParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := partner_channel_template.UpdateWithMap(map[string]interface{}{
	"id": 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /partner_channel_templates/{id}`

Delete Partner Channel Template

```go
func (c *Client) Delete(params files_sdk.PartnerChannelTemplateDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.PartnerChannelTemplateDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.PartnerChannelTemplateDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Partner Channel Template ID. |

Example:

```go
err := partner_channel_template.Delete(files_sdk.PartnerChannelTemplateDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.PartnerChannelTemplate` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | The unique ID of the Partner Channel Template. |
| `WorkspaceId` | `int64` | `workspace_id` | ID of the Workspace associated with this Partner Channel Template. |
| `Direction` | `string` | `direction` | Channel directions. `two_way` enables both directions, `to_partner` enables outgoing downloads, and `from_partner` enables incoming uploads. |
| `UseChannelRoot` | `*bool` | `use_channel_root` | Use the Channel folder directly for a one-way exchange. Defaults to false. Cannot be changed after creation. Folder name overrides must be blank when enabled, and the Template must remain one-way. |
| `Name` | `string` | `name` | The name of the Partner Channel Template. |
| `Path` | `string` | `path` | Channel path relative to the Partner root folder. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters. |
| `ToPartnerFolderName` | `string` | `to_partner_folder_name` | Optional Channel-level to-Partner folder name override. |
| `FromPartnerFolderName` | `string` | `from_partner_folder_name` | Optional Channel-level from-Partner folder name override. |
| `FromPartnerRoutePathPattern` | `string` | `from_partner_route_path_pattern` | Optional route path pattern for files uploaded by the Partner. Supports {{partner_name}}. |
| `ToPartnerRoutePathPattern` | `string` | `to_partner_route_path_pattern` | Optional route path pattern for files delivered to the Partner. Supports {{partner_name}}. |
| `ToPartnerManagedFolderPaths` | `[]string` | `to_partner_managed_folder_paths` | Managed folder paths inside the to-Partner folder. |
| `FromPartnerManagedFolderPaths` | `[]string` | `from_partner_managed_folder_paths` | Managed folder paths inside the from-Partner folder. |
| `EffectiveToPartnerFolderName` | `string` | `effective_to_partner_folder_name` | Resolved to-Partner subfolder name. Null when the direction is disabled or uses the Channel folder directly. |
| `EffectiveFromPartnerFolderName` | `string` | `effective_from_partner_folder_name` | Resolved from-Partner subfolder name. Null when the direction is disabled or uses the Channel folder directly. |
