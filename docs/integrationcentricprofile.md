# IntegrationCentricProfile

`files_sdk.IntegrationCentricProfile` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/integrationcentricprofile`, named `integration_centric_profile`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	integration_centric_profile "github.com/Files-com/files-sdk-go/v3/integrationcentricprofile"
)
```

## Operations

### List

`GET /integration_centric_profiles`

List Integration Centric Profiles

```go
func (c *Client) List(params files_sdk.IntegrationCentricProfileListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.IntegrationCentricProfileListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.IntegrationCentricProfile` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.IntegrationCentricProfileListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id` and `name`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `workspace_id`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := integration_centric_profile.List(files_sdk.IntegrationCentricProfileListParams{})
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

`GET /integration_centric_profiles/{id}`

Show Integration Centric Profile

```go
func (c *Client) Find(params files_sdk.IntegrationCentricProfileFindParams, opts ...files_sdk.RequestResponseOption) (integrationCentricProfile files_sdk.IntegrationCentricProfile, err error)
func Find(params files_sdk.IntegrationCentricProfileFindParams, opts ...files_sdk.RequestResponseOption) (integrationCentricProfile files_sdk.IntegrationCentricProfile, err error)
```

Parameters, `files_sdk.IntegrationCentricProfileFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Integration Centric Profile ID. |

Example:

```go
result, err := integration_centric_profile.Find(files_sdk.IntegrationCentricProfileFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Create

`POST /integration_centric_profiles`

Create Integration Centric Profile

```go
func (c *Client) Create(params files_sdk.IntegrationCentricProfileCreateParams, opts ...files_sdk.RequestResponseOption) (integrationCentricProfile files_sdk.IntegrationCentricProfile, err error)
func Create(params files_sdk.IntegrationCentricProfileCreateParams, opts ...files_sdk.RequestResponseOption) (integrationCentricProfile files_sdk.IntegrationCentricProfile, err error)
```

Parameters, `files_sdk.IntegrationCentricProfileCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Name` | `string` | `name` | yes | Profile name |
| `ExpectedRemoteServers` | `[]map[string]interface{}` | `expected_remote_servers` | yes | Remote Server integrations the user is expected to add and connect. Each entry requires `server_type` and may include a display `name`. |
| `WorkspaceId` | `int64` | `workspace_id` | no | Workspace ID |
| `UseForAllUsers` | `*bool` | `use_for_all_users` | no | Whether this profile applies to all users in the Workspace by default |

Example:

```go
// Set this required value before the call.
var expectedRemoteServers []map[string]interface{}
result, err := integration_centric_profile.Create(files_sdk.IntegrationCentricProfileCreateParams{
	Name:                  "Business Systems Onboarding",
	ExpectedRemoteServers: expectedRemoteServers,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Update

`PATCH /integration_centric_profiles/{id}`

Update Integration Centric Profile

```go
func (c *Client) Update(params files_sdk.IntegrationCentricProfileUpdateParams, opts ...files_sdk.RequestResponseOption) (integrationCentricProfile files_sdk.IntegrationCentricProfile, err error)
func Update(params files_sdk.IntegrationCentricProfileUpdateParams, opts ...files_sdk.RequestResponseOption) (integrationCentricProfile files_sdk.IntegrationCentricProfile, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (integrationCentricProfile files_sdk.IntegrationCentricProfile, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (integrationCentricProfile files_sdk.IntegrationCentricProfile, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.IntegrationCentricProfileUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Integration Centric Profile ID. |
| `Name` | `string` | `name` | no | Profile name |
| `WorkspaceId` | `int64` | `workspace_id` | no | Workspace ID |
| `ExpectedRemoteServers` | `[]map[string]interface{}` | `expected_remote_servers` | no | Remote Server integrations the user is expected to add and connect. Each entry requires `server_type` and may include a display `name`. |
| `UseForAllUsers` | `*bool` | `use_for_all_users` | no | Whether this profile applies to all users in the Workspace by default |

Example:

```go
result, err := integration_centric_profile.Update(files_sdk.IntegrationCentricProfileUpdateParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := integration_centric_profile.UpdateWithMap(map[string]interface{}{
	"id": 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /integration_centric_profiles/{id}`

Delete Integration Centric Profile

```go
func (c *Client) Delete(params files_sdk.IntegrationCentricProfileDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.IntegrationCentricProfileDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.IntegrationCentricProfileDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Integration Centric Profile ID. |

Example:

```go
err := integration_centric_profile.Delete(files_sdk.IntegrationCentricProfileDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.IntegrationCentricProfile` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Integration Centric Profile ID |
| `Name` | `string` | `name` | Profile name |
| `WorkspaceId` | `int64` | `workspace_id` | Workspace ID |
| `UseForAllUsers` | `*bool` | `use_for_all_users` | Whether this profile applies to all users in the Workspace by default |
| `ExpectedRemoteServers` | `[]map[string]interface{}` | `expected_remote_servers` | Remote Server integrations the user is expected to add and connect. Each entry requires `server_type` and may include a display `name`. |
