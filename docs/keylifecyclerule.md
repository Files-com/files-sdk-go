# KeyLifecycleRule

`files_sdk.KeyLifecycleRule` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/keylifecyclerule`, named `key_lifecycle_rule`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	key_lifecycle_rule "github.com/Files-com/files-sdk-go/v3/keylifecyclerule"
)
```

## Operations

### List

`GET /key_lifecycle_rules`

List Key Lifecycle Rules

```go
func (c *Client) List(params files_sdk.KeyLifecycleRuleListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.KeyLifecycleRuleListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.KeyLifecycleRule` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.KeyLifecycleRuleListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id` and `key_type`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `key_type` and `workspace_id`. Valid field combinations are `[ workspace_id, key_type ]`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := key_lifecycle_rule.List(files_sdk.KeyLifecycleRuleListParams{})
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

`GET /key_lifecycle_rules/{id}`

Show Key Lifecycle Rule

```go
func (c *Client) Find(params files_sdk.KeyLifecycleRuleFindParams, opts ...files_sdk.RequestResponseOption) (keyLifecycleRule files_sdk.KeyLifecycleRule, err error)
func Find(params files_sdk.KeyLifecycleRuleFindParams, opts ...files_sdk.RequestResponseOption) (keyLifecycleRule files_sdk.KeyLifecycleRule, err error)
```

Parameters, `files_sdk.KeyLifecycleRuleFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Key Lifecycle Rule ID. |

Example:

```go
result, err := key_lifecycle_rule.Find(files_sdk.KeyLifecycleRuleFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Create

`POST /key_lifecycle_rules`

Create Key Lifecycle Rule

```go
func (c *Client) Create(params files_sdk.KeyLifecycleRuleCreateParams, opts ...files_sdk.RequestResponseOption) (keyLifecycleRule files_sdk.KeyLifecycleRule, err error)
func Create(params files_sdk.KeyLifecycleRuleCreateParams, opts ...files_sdk.RequestResponseOption) (keyLifecycleRule files_sdk.KeyLifecycleRule, err error)
```

Parameters, `files_sdk.KeyLifecycleRuleCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `ApplyToAllWorkspaces` | `*bool` | `apply_to_all_workspaces` | no | If true, a default-workspace rule also applies to keys in all workspaces. |
| `ExpirationDays` | `int64` | `expiration_days` | no | Number of days after creation before an SSH key expires. Applies only to SSH keys. |
| `KeyType` | `files_sdk.KeyLifecycleRuleKeyTypeEnum` | `key_type` | no | Key type for which the rule will apply (gpg, ssh, or api). Values: `gpg`, `ssh`, `api`. |
| `InactivityDays` | `int64` | `inactivity_days` | no | Number of days of inactivity before the rule applies. |
| `Name` | `string` | `name` | no | Key Lifecycle Rule name |
| `WorkspaceId` | `int64` | `workspace_id` | no | Workspace ID. `0` means the default workspace. |

Example:

```go
result, err := key_lifecycle_rule.Create(files_sdk.KeyLifecycleRuleCreateParams{})
if err != nil {
	return err
}
fmt.Println(result)
```

### Update

`PATCH /key_lifecycle_rules/{id}`

Update Key Lifecycle Rule

```go
func (c *Client) Update(params files_sdk.KeyLifecycleRuleUpdateParams, opts ...files_sdk.RequestResponseOption) (keyLifecycleRule files_sdk.KeyLifecycleRule, err error)
func Update(params files_sdk.KeyLifecycleRuleUpdateParams, opts ...files_sdk.RequestResponseOption) (keyLifecycleRule files_sdk.KeyLifecycleRule, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (keyLifecycleRule files_sdk.KeyLifecycleRule, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (keyLifecycleRule files_sdk.KeyLifecycleRule, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.KeyLifecycleRuleUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Key Lifecycle Rule ID. |
| `ApplyToAllWorkspaces` | `*bool` | `apply_to_all_workspaces` | no | If true, a default-workspace rule also applies to keys in all workspaces. |
| `ExpirationDays` | `int64` | `expiration_days` | no | Number of days after creation before an SSH key expires. Applies only to SSH keys. |
| `KeyType` | `files_sdk.KeyLifecycleRuleKeyTypeEnum` | `key_type` | no | Key type for which the rule will apply (gpg, ssh, or api). Values: `gpg`, `ssh`, `api`. |
| `InactivityDays` | `int64` | `inactivity_days` | no | Number of days of inactivity before the rule applies. |
| `Name` | `string` | `name` | no | Key Lifecycle Rule name |
| `WorkspaceId` | `int64` | `workspace_id` | no | Workspace ID. `0` means the default workspace. |

Example:

```go
result, err := key_lifecycle_rule.Update(files_sdk.KeyLifecycleRuleUpdateParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := key_lifecycle_rule.UpdateWithMap(map[string]interface{}{
	"id": 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /key_lifecycle_rules/{id}`

Delete Key Lifecycle Rule

```go
func (c *Client) Delete(params files_sdk.KeyLifecycleRuleDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.KeyLifecycleRuleDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.KeyLifecycleRuleDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Key Lifecycle Rule ID. |

Example:

```go
err := key_lifecycle_rule.Delete(files_sdk.KeyLifecycleRuleDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.KeyLifecycleRule` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Key Lifecycle Rule ID |
| `KeyType` | `string` | `key_type` | Key type for which the rule will apply (gpg, ssh, or api). |
| `InactivityDays` | `int64` | `inactivity_days` | Number of days of inactivity before the rule applies. |
| `ExpirationDays` | `int64` | `expiration_days` | Number of days after creation before an SSH key expires. Applies only to SSH keys. |
| `ApplyToAllWorkspaces` | `*bool` | `apply_to_all_workspaces` | If true, a default-workspace rule also applies to keys in all workspaces. |
| `Name` | `string` | `name` | Key Lifecycle Rule name |
| `WorkspaceId` | `int64` | `workspace_id` | Workspace ID. `0` means the default workspace. |
