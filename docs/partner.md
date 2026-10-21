# Partner

`files_sdk.Partner` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/partner`, named `partner`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/partner"
)
```

## Operations

### List

`GET /partners`

List Partners

```go
func (c *Client) List(params files_sdk.PartnerListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.PartnerListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.Partner` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.PartnerListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id` and `name`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `workspace_id`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := partner.List(files_sdk.PartnerListParams{})
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

`GET /partners/{id}`

Show Partner

```go
func (c *Client) Find(params files_sdk.PartnerFindParams, opts ...files_sdk.RequestResponseOption) (partner files_sdk.Partner, err error)
func Find(params files_sdk.PartnerFindParams, opts ...files_sdk.RequestResponseOption) (partner files_sdk.Partner, err error)
```

Parameters, `files_sdk.PartnerFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Partner ID. |

Example:

```go
result, err := partner.Find(files_sdk.PartnerFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Create

`POST /partners`

Create Partner

```go
func (c *Client) Create(params files_sdk.PartnerCreateParams, opts ...files_sdk.RequestResponseOption) (partner files_sdk.Partner, err error)
func Create(params files_sdk.PartnerCreateParams, opts ...files_sdk.RequestResponseOption) (partner files_sdk.Partner, err error)
```

Parameters, `files_sdk.PartnerCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `AiAssistantPersonalityId` | `int64` | `ai_assistant_personality_id` | no | AI Assistant Personality ID assigned to this Partner, if any. Users in the Partner inherit it unless a direct per-user assignment overrides it. |
| `AllowedIps` | `string` | `allowed_ips` | no | A list of allowed IPs for this Partner. Newline delimited. Partner User IP access is allowed when the IP matches the Partner, User, or Site allowed IP lists. |
| `AllowBypassing2faPolicies` | `*bool` | `allow_bypassing_2fa_policies` | no | Allow Partner Admins to change Two-Factor Authentication requirements for Partner Users. |
| `AllowCredentialChanges` | `*bool` | `allow_credential_changes` | no | Allow Partner Admins to change or reset credentials for users belonging to this Partner. |
| `AllowProvidingGpgKeys` | `*bool` | `allow_providing_gpg_keys` | no | Allow Partner Admins to provide GPG keys. |
| `AllowUserCreation` | `*bool` | `allow_user_creation` | no | Allow Partner Admins to create users. |
| `CcEmailsToResponsibleParty` | `*bool` | `cc_emails_to_responsible_party` | no | When `true`, emails sent to Partner users are copied to the responsible User or Group. |
| `Notes` | `string` | `notes` | no | Notes about this Partner. |
| `PartnerChannelTemplateId` | `int64` | `partner_channel_template_id` | no | ID of the Partner Channel Template assigned to this Partner. |
| `ResponsibleGroupId` | `int64` | `responsible_group_id` | no | ID of the Group responsible for this Partner. |
| `ResponsibleUserId` | `int64` | `responsible_user_id` | no | ID of the User responsible for this Partner. |
| `ShowPartnerChannelHomePage` | `*bool` | `show_partner_channel_home_page` | no | Show Partner users a simplified home page built from this Partner's Channels. |
| `Tags` | `string` | `tags` | no | Comma-separated list of Tags for this Partner. Tags are used for other features, such as UserLifecycleRules, which can target specific tags. Tags must only contain lowercase letters, numbers, and hyphens. |
| `Name` | `string` | `name` | yes | The name of the Partner. |
| `RootFolder` | `string` | `root_folder` | yes | The root folder path for this Partner. |
| `WorkspaceId` | `int64` | `workspace_id` | no | ID of the Workspace associated with this Partner. |

Example:

```go
result, err := partner.Create(files_sdk.PartnerCreateParams{
	Name:       "Acme Corp",
	RootFolder: "/AcmeCorp",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Update

`PATCH /partners/{id}`

Update Partner

```go
func (c *Client) Update(params files_sdk.PartnerUpdateParams, opts ...files_sdk.RequestResponseOption) (partner files_sdk.Partner, err error)
func Update(params files_sdk.PartnerUpdateParams, opts ...files_sdk.RequestResponseOption) (partner files_sdk.Partner, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (partner files_sdk.Partner, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (partner files_sdk.Partner, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.PartnerUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Partner ID. |
| `AiAssistantPersonalityId` | `int64` | `ai_assistant_personality_id` | no | AI Assistant Personality ID assigned to this Partner, if any. Users in the Partner inherit it unless a direct per-user assignment overrides it. |
| `AllowedIps` | `string` | `allowed_ips` | no | A list of allowed IPs for this Partner. Newline delimited. Partner User IP access is allowed when the IP matches the Partner, User, or Site allowed IP lists. |
| `AllowBypassing2faPolicies` | `*bool` | `allow_bypassing_2fa_policies` | no | Allow Partner Admins to change Two-Factor Authentication requirements for Partner Users. |
| `AllowCredentialChanges` | `*bool` | `allow_credential_changes` | no | Allow Partner Admins to change or reset credentials for users belonging to this Partner. |
| `AllowProvidingGpgKeys` | `*bool` | `allow_providing_gpg_keys` | no | Allow Partner Admins to provide GPG keys. |
| `AllowUserCreation` | `*bool` | `allow_user_creation` | no | Allow Partner Admins to create users. |
| `CcEmailsToResponsibleParty` | `*bool` | `cc_emails_to_responsible_party` | no | When `true`, emails sent to Partner users are copied to the responsible User or Group. |
| `Notes` | `string` | `notes` | no | Notes about this Partner. |
| `PartnerChannelTemplateId` | `int64` | `partner_channel_template_id` | no | ID of the Partner Channel Template assigned to this Partner. |
| `ResponsibleGroupId` | `int64` | `responsible_group_id` | no | ID of the Group responsible for this Partner. |
| `ResponsibleUserId` | `int64` | `responsible_user_id` | no | ID of the User responsible for this Partner. |
| `ShowPartnerChannelHomePage` | `*bool` | `show_partner_channel_home_page` | no | Show Partner users a simplified home page built from this Partner's Channels. |
| `Tags` | `string` | `tags` | no | Comma-separated list of Tags for this Partner. Tags are used for other features, such as UserLifecycleRules, which can target specific tags. Tags must only contain lowercase letters, numbers, and hyphens. |
| `Name` | `string` | `name` | no | The name of the Partner. |
| `PartnershipRole` | `files_sdk.PartnerPartnershipRoleEnum` | `partnership_role` | no | This site's role for this Partner in Connected Sites relationships. `host` is a Partner this site configured. `guest` is a Partner created by approving another site's connection request; it has no root folder and cannot hold users, permissions, or Partner Channels, or host a connection. `host_and_guest` is a configured Partner that is also the guest side of a connection. Promote a `guest` Partner by setting this to `host_and_guest` together with a `root_folder`. Values: `host`, `guest`, `host_and_guest`. |
| `RootFolder` | `string` | `root_folder` | no | The root folder path for this Partner. |

Example:

```go
result, err := partner.Update(files_sdk.PartnerUpdateParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := partner.UpdateWithMap(map[string]interface{}{
	"id": 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /partners/{id}`

Delete Partner

```go
func (c *Client) Delete(params files_sdk.PartnerDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.PartnerDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.PartnerDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Partner ID. |

Example:

```go
err := partner.Delete(files_sdk.PartnerDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.Partner` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `AllowBypassing2faPolicies` | `*bool` | `allow_bypassing_2fa_policies` | Allow Partner Admins to change Two-Factor Authentication requirements for Partner Users. |
| `AllowedIps` | `string` | `allowed_ips` | A list of allowed IPs for this Partner. Newline delimited. Partner User IP access is allowed when the IP matches the Partner, User, or Site allowed IP lists. |
| `AllowCredentialChanges` | `*bool` | `allow_credential_changes` | Allow Partner Admins to change or reset credentials for users belonging to this Partner. |
| `AllowProvidingGpgKeys` | `*bool` | `allow_providing_gpg_keys` | Allow Partner Admins to provide GPG keys. |
| `AllowUserCreation` | `*bool` | `allow_user_creation` | Allow Partner Admins to create users. |
| `CcEmailsToResponsibleParty` | `*bool` | `cc_emails_to_responsible_party` | When `true`, emails sent to Partner users are copied to the responsible User or Group. |
| `Connections` | `[]files_sdk.PartnerConnection` | `connections` | Approved Connected Sites relationships for this Partner, in both directions. Empty when this Partner has no connections. Read-only. |
| `Id` | `int64` | `id` | The unique ID of the Partner. |
| `AiAssistantPersonalityId` | `int64` | `ai_assistant_personality_id` | AI Assistant Personality ID assigned to this Partner, if any. Users in the Partner inherit it unless a direct per-user assignment overrides it. |
| `WorkspaceId` | `int64` | `workspace_id` | ID of the Workspace associated with this Partner. |
| `Name` | `string` | `name` | The name of the Partner. |
| `Notes` | `string` | `notes` | Notes about this Partner. |
| `PartnerAdminIds` | `[]int64` | `partner_admin_ids` | Array of User IDs that are Partner Admins for this Partner. |
| `PartnerChannelTemplateId` | `int64` | `partner_channel_template_id` | ID of the Partner Channel Template assigned to this Partner. |
| `PartnershipRole` | `string` | `partnership_role` | This site's role for this Partner in Connected Sites relationships. `host` is a Partner this site configured. `guest` is a Partner created by approving another site's connection request; it has no root folder and cannot hold users, permissions, or Partner Channels, or host a connection. `host_and_guest` is a configured Partner that is also the guest side of a connection. Promote a `guest` Partner by setting this to `host_and_guest` together with a `root_folder`. |
| `ResponsibleGroupId` | `int64` | `responsible_group_id` | ID of the Group responsible for this Partner. |
| `ResponsibleUserId` | `int64` | `responsible_user_id` | ID of the User responsible for this Partner. |
| `RootFolder` | `string` | `root_folder` | The root folder path for this Partner. |
| `ShowPartnerChannelHomePage` | `*bool` | `show_partner_channel_home_page` | Show Partner users a simplified home page built from this Partner's Channels. |
| `Tags` | `string` | `tags` | Comma-separated list of Tags for this Partner. Tags are used for other features, such as UserLifecycleRules, which can target specific tags. Tags must only contain lowercase letters, numbers, and hyphens. |
| `UserIds` | `[]int64` | `user_ids` | Array of User IDs that belong to this Partner. |
