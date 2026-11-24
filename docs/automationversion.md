# AutomationVersion

`files_sdk.AutomationVersion` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/automationversion`, named `automation_version`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	automation_version "github.com/Files-com/files-sdk-go/v3/automationversion"
)
```

## Operations

### List

`GET /automation_versions`

List Automation Versions

```go
func (c *Client) List(params files_sdk.AutomationVersionListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.AutomationVersionListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.AutomationVersion` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.AutomationVersionListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are . |
| `AutomationId` | `int64` | `automation_id` | yes | ID of the Automation that owns this revision. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := automation_version.List(files_sdk.AutomationVersionListParams{
	AutomationId: 1,
})
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

`GET /automation_versions/{id}`

Show Automation Version

```go
func (c *Client) Find(params files_sdk.AutomationVersionFindParams, opts ...files_sdk.RequestResponseOption) (automationVersion files_sdk.AutomationVersion, err error)
func Find(params files_sdk.AutomationVersionFindParams, opts ...files_sdk.RequestResponseOption) (automationVersion files_sdk.AutomationVersion, err error)
```

Parameters, `files_sdk.AutomationVersionFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Automation Version ID. |

Example:

```go
result, err := automation_version.Find(files_sdk.AutomationVersionFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Restore

`POST /automation_versions/{id}/restore`

Restore an Automation v2 definition as a new current version

```go
func (c *Client) Restore(params files_sdk.AutomationVersionRestoreParams, opts ...files_sdk.RequestResponseOption) (automation files_sdk.Automation, err error)
func Restore(params files_sdk.AutomationVersionRestoreParams, opts ...files_sdk.RequestResponseOption) (automation files_sdk.Automation, err error)
```

Parameters, `files_sdk.AutomationVersionRestoreParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Automation Version ID. |

Example:

```go
result, err := automation_version.Restore(files_sdk.AutomationVersionRestoreParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

## Model

`files_sdk.AutomationVersion` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Automation revision ID. |
| `AutomationId` | `int64` | `automation_id` | ID of the Automation that owns this revision. |
| `WorkspaceId` | `int64` | `workspace_id` | Workspace that owns the Automation. Zero is the Default Workspace. |
| `Version` | `int64` | `version` | Revision number within the Automation. |
| `CreatedAt` | `*time.Time` | `created_at` | Time the revision was saved. |
| `UserId` | `int64` | `user_id` | ID of the user who saved the revision, when available. |
| `Username` | `string` | `username` | Username recorded when the revision was saved, when available. |
| `ApiKeyId` | `int64` | `api_key_id` | ID of the API key used to save the revision, when available. |
| `Definition` | `interface{}` | `definition` | Saved Automation v2 graph definition. Included when retrieving a single revision. |
