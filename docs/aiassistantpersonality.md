# AiAssistantPersonality

`files_sdk.AiAssistantPersonality` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/aiassistantpersonality`, named `ai_assistant_personality`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	ai_assistant_personality "github.com/Files-com/files-sdk-go/v3/aiassistantpersonality"
)
```

## Operations

### List

`GET /ai_assistant_personalities`

List Ai Assistant Personalities

```go
func (c *Client) List(params files_sdk.AiAssistantPersonalityListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.AiAssistantPersonalityListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.AiAssistantPersonality` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.AiAssistantPersonalityListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id` and `id`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `workspace_id`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := ai_assistant_personality.List(files_sdk.AiAssistantPersonalityListParams{})
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

`GET /ai_assistant_personalities/{id}`

Show Ai Assistant Personality

```go
func (c *Client) Find(params files_sdk.AiAssistantPersonalityFindParams, opts ...files_sdk.RequestResponseOption) (aiAssistantPersonality files_sdk.AiAssistantPersonality, err error)
func Find(params files_sdk.AiAssistantPersonalityFindParams, opts ...files_sdk.RequestResponseOption) (aiAssistantPersonality files_sdk.AiAssistantPersonality, err error)
```

Parameters, `files_sdk.AiAssistantPersonalityFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Ai Assistant Personality ID. |

Example:

```go
result, err := ai_assistant_personality.Find(files_sdk.AiAssistantPersonalityFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Create

`POST /ai_assistant_personalities`

Create Ai Assistant Personality

```go
func (c *Client) Create(params files_sdk.AiAssistantPersonalityCreateParams, opts ...files_sdk.RequestResponseOption) (aiAssistantPersonality files_sdk.AiAssistantPersonality, err error)
func Create(params files_sdk.AiAssistantPersonalityCreateParams, opts ...files_sdk.RequestResponseOption) (aiAssistantPersonality files_sdk.AiAssistantPersonality, err error)
```

Parameters, `files_sdk.AiAssistantPersonalityCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `ApplyToAllWorkspaces` | `*bool` | `apply_to_all_workspaces` | no | If true, this default-workspace personality can apply to users in all workspaces. |
| `Name` | `string` | `name` | yes | AI Assistant Personality name. |
| `SystemPrompt` | `string` | `system_prompt` | yes | System prompt injected into the in-app AI Assistant. |
| `UseByDefault` | `*bool` | `use_by_default` | no | Whether this personality is the default personality for the Workspace. |
| `WorkspaceId` | `int64` | `workspace_id` | no | Workspace ID. `0` means the default workspace. |

Example:

```go
result, err := ai_assistant_personality.Create(files_sdk.AiAssistantPersonalityCreateParams{
	Name:         "Concise Assistant",
	SystemPrompt: "Respond as a concise operations assistant.",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Update

`PATCH /ai_assistant_personalities/{id}`

Update Ai Assistant Personality

```go
func (c *Client) Update(params files_sdk.AiAssistantPersonalityUpdateParams, opts ...files_sdk.RequestResponseOption) (aiAssistantPersonality files_sdk.AiAssistantPersonality, err error)
func Update(params files_sdk.AiAssistantPersonalityUpdateParams, opts ...files_sdk.RequestResponseOption) (aiAssistantPersonality files_sdk.AiAssistantPersonality, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (aiAssistantPersonality files_sdk.AiAssistantPersonality, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (aiAssistantPersonality files_sdk.AiAssistantPersonality, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.AiAssistantPersonalityUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Ai Assistant Personality ID. |
| `ApplyToAllWorkspaces` | `*bool` | `apply_to_all_workspaces` | no | If true, this default-workspace personality can apply to users in all workspaces. |
| `Name` | `string` | `name` | no | AI Assistant Personality name. |
| `SystemPrompt` | `string` | `system_prompt` | no | System prompt injected into the in-app AI Assistant. |
| `UseByDefault` | `*bool` | `use_by_default` | no | Whether this personality is the default personality for the Workspace. |
| `WorkspaceId` | `int64` | `workspace_id` | no | Workspace ID. `0` means the default workspace. |

Example:

```go
result, err := ai_assistant_personality.Update(files_sdk.AiAssistantPersonalityUpdateParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := ai_assistant_personality.UpdateWithMap(map[string]interface{}{
	"id": 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /ai_assistant_personalities/{id}`

Delete Ai Assistant Personality

```go
func (c *Client) Delete(params files_sdk.AiAssistantPersonalityDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.AiAssistantPersonalityDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.AiAssistantPersonalityDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Ai Assistant Personality ID. |

Example:

```go
err := ai_assistant_personality.Delete(files_sdk.AiAssistantPersonalityDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.AiAssistantPersonality` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | AI Assistant Personality ID. |
| `WorkspaceId` | `int64` | `workspace_id` | Workspace ID. `0` means the default workspace. |
| `Name` | `string` | `name` | AI Assistant Personality name. |
| `SystemPrompt` | `string` | `system_prompt` | System prompt injected into the in-app AI Assistant. |
| `UseByDefault` | `*bool` | `use_by_default` | Whether this personality is the default personality for the Workspace. |
| `ApplyToAllWorkspaces` | `*bool` | `apply_to_all_workspaces` | If true, this default-workspace personality can apply to users in all workspaces. |
| `CreatedAt` | `*time.Time` | `created_at` | Creation time. |
| `UpdatedAt` | `*time.Time` | `updated_at` | Last update time. |
