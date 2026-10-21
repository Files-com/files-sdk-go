# ChildSiteManagementPolicy

`files_sdk.ChildSiteManagementPolicy` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/childsitemanagementpolicy`, named `child_site_management_policy`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	child_site_management_policy "github.com/Files-com/files-sdk-go/v3/childsitemanagementpolicy"
)
```

## Operations

### List

`GET /child_site_management_policies`

List Child Site Management Policies. On a Child Site, returns only the policy covering that site. Its value contains managed setting names with null values

```go
func (c *Client) List(params files_sdk.ChildSiteManagementPolicyListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.ChildSiteManagementPolicyListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.ChildSiteManagementPolicy` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.ChildSiteManagementPolicyListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := child_site_management_policy.List(files_sdk.ChildSiteManagementPolicyListParams{})
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

`GET /child_site_management_policies/{id}`

Show Child Site Management Policy

```go
func (c *Client) Find(params files_sdk.ChildSiteManagementPolicyFindParams, opts ...files_sdk.RequestResponseOption) (childSiteManagementPolicy files_sdk.ChildSiteManagementPolicy, err error)
func Find(params files_sdk.ChildSiteManagementPolicyFindParams, opts ...files_sdk.RequestResponseOption) (childSiteManagementPolicy files_sdk.ChildSiteManagementPolicy, err error)
```

Parameters, `files_sdk.ChildSiteManagementPolicyFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Child Site Management Policy ID. |

Example:

```go
result, err := child_site_management_policy.Find(files_sdk.ChildSiteManagementPolicyFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Create

`POST /child_site_management_policies`

Create Child Site Management Policy

```go
func (c *Client) Create(params files_sdk.ChildSiteManagementPolicyCreateParams, opts ...files_sdk.RequestResponseOption) (childSiteManagementPolicy files_sdk.ChildSiteManagementPolicy, err error)
func Create(params files_sdk.ChildSiteManagementPolicyCreateParams, opts ...files_sdk.RequestResponseOption) (childSiteManagementPolicy files_sdk.ChildSiteManagementPolicy, err error)
```

Parameters, `files_sdk.ChildSiteManagementPolicyCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Value` | `interface{}` | `value` | no | Policy configuration data. Attributes differ by policy type. For more information, refer to the Value Hash section of the developer documentation. |
| `SkipChildSiteIds` | `[]int64` | `skip_child_site_ids` | no | IDs of child sites excluded from this default policy. |
| `ChildSiteIds` | `[]int64` | `child_site_ids` | no | IDs of child sites explicitly assigned to this non-default policy. |
| `DefaultPolicy` | `*bool` | `default_policy` | no | Whether this policy applies to child sites not explicitly assigned to another policy. |
| `PolicyType` | `files_sdk.ChildSiteManagementPolicyPolicyTypeEnum` | `policy_type` | yes | Type of policy. Valid values: `settings`. Values: `settings`. |
| `Name` | `string` | `name` | no | Name for this policy. |
| `Description` | `string` | `description` | no | Description for this policy. |

Example:

```go
result, err := child_site_management_policy.Create(files_sdk.ChildSiteManagementPolicyCreateParams{
	PolicyType: "settings",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Update

`PATCH /child_site_management_policies/{id}`

Update Child Site Management Policy

```go
func (c *Client) Update(params files_sdk.ChildSiteManagementPolicyUpdateParams, opts ...files_sdk.RequestResponseOption) (childSiteManagementPolicy files_sdk.ChildSiteManagementPolicy, err error)
func Update(params files_sdk.ChildSiteManagementPolicyUpdateParams, opts ...files_sdk.RequestResponseOption) (childSiteManagementPolicy files_sdk.ChildSiteManagementPolicy, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (childSiteManagementPolicy files_sdk.ChildSiteManagementPolicy, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (childSiteManagementPolicy files_sdk.ChildSiteManagementPolicy, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.ChildSiteManagementPolicyUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Child Site Management Policy ID. |
| `Value` | `interface{}` | `value` | no | Policy configuration data. Attributes differ by policy type. For more information, refer to the Value Hash section of the developer documentation. |
| `SkipChildSiteIds` | `[]int64` | `skip_child_site_ids` | no | IDs of child sites excluded from this default policy. |
| `ChildSiteIds` | `[]int64` | `child_site_ids` | no | IDs of child sites explicitly assigned to this non-default policy. |
| `DefaultPolicy` | `*bool` | `default_policy` | no | Whether this policy applies to child sites not explicitly assigned to another policy. |
| `PolicyType` | `files_sdk.ChildSiteManagementPolicyPolicyTypeEnum` | `policy_type` | no | Type of policy. Valid values: `settings`. Values: `settings`. |
| `Name` | `string` | `name` | no | Name for this policy. |
| `Description` | `string` | `description` | no | Description for this policy. |

Example:

```go
result, err := child_site_management_policy.Update(files_sdk.ChildSiteManagementPolicyUpdateParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := child_site_management_policy.UpdateWithMap(map[string]interface{}{
	"id": 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /child_site_management_policies/{id}`

Delete Child Site Management Policy

```go
func (c *Client) Delete(params files_sdk.ChildSiteManagementPolicyDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.ChildSiteManagementPolicyDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.ChildSiteManagementPolicyDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Child Site Management Policy ID. |

Example:

```go
err := child_site_management_policy.Delete(files_sdk.ChildSiteManagementPolicyDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.ChildSiteManagementPolicy` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Policy ID. |
| `PolicyType` | `string` | `policy_type` | Type of policy. Valid values: `settings`. |
| `Name` | `string` | `name` | Name for this policy. |
| `Description` | `string` | `description` | Description for this policy. |
| `Value` | `interface{}` | `value` | Policy configuration data. On a Child Site, contains only managed setting names with null values. Settings policies accept site settings plus an optional `folder_behaviors` array for parent-managed root behaviors on child sites. For more information, refer to the Value Hash section of the developer documentation. |
| `AppliedChildSiteIds` | `[]int64` | `applied_child_site_ids` | IDs of child sites that this policy has been applied to. This field is read-only. |
| `SkipChildSiteIds` | `[]int64` | `skip_child_site_ids` | IDs of child sites excluded from this default policy. |
| `ChildSiteIds` | `[]int64` | `child_site_ids` | IDs of child sites explicitly assigned to this non-default policy. |
| `DefaultPolicy` | `*bool` | `default_policy` | Whether this policy applies to child sites not explicitly assigned to another policy. |
| `CreatedAt` | `*time.Time` | `created_at` | When this policy was created. |
| `UpdatedAt` | `*time.Time` | `updated_at` | When this policy was last updated. |
