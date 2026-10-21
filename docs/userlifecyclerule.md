# UserLifecycleRule

`files_sdk.UserLifecycleRule` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/userlifecyclerule`, named `user_lifecycle_rule`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	user_lifecycle_rule "github.com/Files-com/files-sdk-go/v3/userlifecyclerule"
)
```

## Operations

### List

`GET /user_lifecycle_rules`

List User Lifecycle Rules

```go
func (c *Client) List(params files_sdk.UserLifecycleRuleListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.UserLifecycleRuleListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.UserLifecycleRule` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.UserLifecycleRuleListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `site_id` and `workspace_id`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `workspace_id`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := user_lifecycle_rule.List(files_sdk.UserLifecycleRuleListParams{})
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

`GET /user_lifecycle_rules/{id}`

Show User Lifecycle Rule

```go
func (c *Client) Find(params files_sdk.UserLifecycleRuleFindParams, opts ...files_sdk.RequestResponseOption) (userLifecycleRule files_sdk.UserLifecycleRule, err error)
func Find(params files_sdk.UserLifecycleRuleFindParams, opts ...files_sdk.RequestResponseOption) (userLifecycleRule files_sdk.UserLifecycleRule, err error)
```

Parameters, `files_sdk.UserLifecycleRuleFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | User Lifecycle Rule ID. |

Example:

```go
result, err := user_lifecycle_rule.Find(files_sdk.UserLifecycleRuleFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Create

`POST /user_lifecycle_rules`

Create User Lifecycle Rule

```go
func (c *Client) Create(params files_sdk.UserLifecycleRuleCreateParams, opts ...files_sdk.RequestResponseOption) (userLifecycleRule files_sdk.UserLifecycleRule, err error)
func Create(params files_sdk.UserLifecycleRuleCreateParams, opts ...files_sdk.RequestResponseOption) (userLifecycleRule files_sdk.UserLifecycleRule, err error)
```

Parameters, `files_sdk.UserLifecycleRuleCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Action` | `files_sdk.UserLifecycleRuleActionEnum` | `action` | no | Action to take on inactive users (disable or delete) Values: `disable`, `delete`. |
| `ApplyToAllWorkspaces` | `*bool` | `apply_to_all_workspaces` | no | If true, a Default Workspace rule also applies to users in all Custom Workspaces. Can only be enabled when `workspace_id` is `0`. |
| `AuthenticationMethod` | `files_sdk.UserLifecycleRuleAuthenticationMethodEnum` | `authentication_method` | no | User authentication method for which the rule will apply. Use `all_non_sso` to target every non-SSO authentication method with one rule. Values: `all`, `password`, `sso`, `none`, `email_signup`, `password_with_imported_hash`, `password_and_ssh_key`, `all_non_sso`. |
| `GroupIds` | `[]int64` | `group_ids` | no | Array of Group IDs to which the rule applies. If empty or not set, the rule applies to all users. |
| `InactivityDays` | `int64` | `inactivity_days` | no | Number of days of inactivity before the rule applies |
| `IncludeSiteAdmins` | `*bool` | `include_site_admins` | no | If true, the rule includes Site Administrators, who always belong to the Default Workspace. Can only be enabled when `workspace_id` is `0`. |
| `IncludeFolderAdmins` | `*bool` | `include_folder_admins` | no | If true, the rule will apply to folder admins. |
| `Name` | `string` | `name` | no | User Lifecycle Rule name |
| `NotifyUsers` | `*bool` | `notify_users` | no | If true, users will be emailed before the rule disables or deletes them. |
| `PartnerTag` | `string` | `partner_tag` | no | If provided, only users belonging to Partners with this tag at the Partner level will be affected by the rule. Tags must only contain lowercase letters, numbers, and hyphens. |
| `UserState` | `files_sdk.UserLifecycleRuleUserStateEnum` | `user_state` | no | State of the users to apply the rule to (inactive or disabled) Values: `inactive`, `disabled`. |
| `UserTag` | `string` | `user_tag` | no | If provided, only users with this tag will be affected by the rule. Tags must only contain lowercase letters, numbers, and hyphens. |
| `WorkspaceId` | `int64` | `workspace_id` | no | Workspace whose users the rule applies to. `0` means the Default Workspace. A Custom Workspace rule applies only to users who belong to that Workspace, regardless of access granted to other users. |

Example:

```go
result, err := user_lifecycle_rule.Create(files_sdk.UserLifecycleRuleCreateParams{})
if err != nil {
	return err
}
fmt.Println(result)
```

### Update

`PATCH /user_lifecycle_rules/{id}`

Update User Lifecycle Rule

```go
func (c *Client) Update(params files_sdk.UserLifecycleRuleUpdateParams, opts ...files_sdk.RequestResponseOption) (userLifecycleRule files_sdk.UserLifecycleRule, err error)
func Update(params files_sdk.UserLifecycleRuleUpdateParams, opts ...files_sdk.RequestResponseOption) (userLifecycleRule files_sdk.UserLifecycleRule, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (userLifecycleRule files_sdk.UserLifecycleRule, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (userLifecycleRule files_sdk.UserLifecycleRule, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.UserLifecycleRuleUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | User Lifecycle Rule ID. |
| `Action` | `files_sdk.UserLifecycleRuleActionEnum` | `action` | no | Action to take on inactive users (disable or delete) Values: `disable`, `delete`. |
| `ApplyToAllWorkspaces` | `*bool` | `apply_to_all_workspaces` | no | If true, a Default Workspace rule also applies to users in all Custom Workspaces. Can only be enabled when `workspace_id` is `0`. |
| `AuthenticationMethod` | `files_sdk.UserLifecycleRuleAuthenticationMethodEnum` | `authentication_method` | no | User authentication method for which the rule will apply. Use `all_non_sso` to target every non-SSO authentication method with one rule. Values: `all`, `password`, `sso`, `none`, `email_signup`, `password_with_imported_hash`, `password_and_ssh_key`, `all_non_sso`. |
| `GroupIds` | `[]int64` | `group_ids` | no | Array of Group IDs to which the rule applies. If empty or not set, the rule applies to all users. |
| `InactivityDays` | `int64` | `inactivity_days` | no | Number of days of inactivity before the rule applies |
| `IncludeSiteAdmins` | `*bool` | `include_site_admins` | no | If true, the rule includes Site Administrators, who always belong to the Default Workspace. Can only be enabled when `workspace_id` is `0`. |
| `IncludeFolderAdmins` | `*bool` | `include_folder_admins` | no | If true, the rule will apply to folder admins. |
| `Name` | `string` | `name` | no | User Lifecycle Rule name |
| `NotifyUsers` | `*bool` | `notify_users` | no | If true, users will be emailed before the rule disables or deletes them. |
| `PartnerTag` | `string` | `partner_tag` | no | If provided, only users belonging to Partners with this tag at the Partner level will be affected by the rule. Tags must only contain lowercase letters, numbers, and hyphens. |
| `UserState` | `files_sdk.UserLifecycleRuleUserStateEnum` | `user_state` | no | State of the users to apply the rule to (inactive or disabled) Values: `inactive`, `disabled`. |
| `UserTag` | `string` | `user_tag` | no | If provided, only users with this tag will be affected by the rule. Tags must only contain lowercase letters, numbers, and hyphens. |
| `WorkspaceId` | `int64` | `workspace_id` | no | Workspace whose users the rule applies to. `0` means the Default Workspace. A Custom Workspace rule applies only to users who belong to that Workspace, regardless of access granted to other users. |

Example:

```go
result, err := user_lifecycle_rule.Update(files_sdk.UserLifecycleRuleUpdateParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := user_lifecycle_rule.UpdateWithMap(map[string]interface{}{
	"id": 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /user_lifecycle_rules/{id}`

Delete User Lifecycle Rule

```go
func (c *Client) Delete(params files_sdk.UserLifecycleRuleDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.UserLifecycleRuleDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.UserLifecycleRuleDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | User Lifecycle Rule ID. |

Example:

```go
err := user_lifecycle_rule.Delete(files_sdk.UserLifecycleRuleDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.UserLifecycleRule` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | User Lifecycle Rule ID |
| `AuthenticationMethod` | `string` | `authentication_method` | User authentication method for which the rule will apply. Use `all_non_sso` to target every non-SSO authentication method with one rule. |
| `GroupIds` | `[]int64` | `group_ids` | Array of Group IDs to which the rule applies. If empty or not set, the rule applies to all users. |
| `Action` | `string` | `action` | Action to take on inactive users (disable or delete) |
| `InactivityDays` | `int64` | `inactivity_days` | Number of days of inactivity before the rule applies |
| `IncludeFolderAdmins` | `*bool` | `include_folder_admins` | If true, the rule will apply to folder admins. |
| `IncludeSiteAdmins` | `*bool` | `include_site_admins` | If true, the rule includes Site Administrators, who always belong to the Default Workspace. Can only be enabled when `workspace_id` is `0`. |
| `ApplyToAllWorkspaces` | `*bool` | `apply_to_all_workspaces` | If true, a Default Workspace rule also applies to users in all Custom Workspaces. Can only be enabled when `workspace_id` is `0`. |
| `Name` | `string` | `name` | User Lifecycle Rule name |
| `NotifyUsers` | `*bool` | `notify_users` | If true, users will be emailed before the rule disables or deletes them. |
| `PartnerTag` | `string` | `partner_tag` | If provided, only users belonging to Partners with this tag at the Partner level will be affected by the rule. Tags must only contain lowercase letters, numbers, and hyphens. |
| `SiteId` | `int64` | `site_id` | Site ID |
| `WorkspaceId` | `int64` | `workspace_id` | Workspace whose users the rule applies to. `0` means the Default Workspace. A Custom Workspace rule applies only to users who belong to that Workspace, regardless of access granted to other users. |
| `UserState` | `string` | `user_state` | State of the users to apply the rule to (inactive or disabled) |
| `UserTag` | `string` | `user_tag` | If provided, only users with this tag will be affected by the rule. Tags must only contain lowercase letters, numbers, and hyphens. |
