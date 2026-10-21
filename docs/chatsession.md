# ChatSession

`files_sdk.ChatSession` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/chatsession`, named `chat_session`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	chat_session "github.com/Files-com/files-sdk-go/v3/chatsession"
)
```

## Operations

### List

`GET /chat_sessions`

List Chat Sessions

```go
func (c *Client) List(params files_sdk.ChatSessionListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.ChatSessionListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.ChatSession` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.ChatSessionListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `id`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `ai_task_id`, `user_id` or `workspace_id`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := chat_session.List(files_sdk.ChatSessionListParams{})
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

`GET /chat_sessions/{id}`

Show Chat Session

```go
func (c *Client) Find(params files_sdk.ChatSessionFindParams, opts ...files_sdk.RequestResponseOption) (chatSession files_sdk.ChatSession, err error)
func Find(params files_sdk.ChatSessionFindParams, opts ...files_sdk.RequestResponseOption) (chatSession files_sdk.ChatSession, err error)
```

Parameters, `files_sdk.ChatSessionFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `string` | `id`, in the path | yes | Chat Session ID. |

Example:

```go
result, err := chat_session.Find(files_sdk.ChatSessionFindParams{
	Id: "id",
})
if err != nil {
	return err
}
fmt.Println(result)
```

## Model

`files_sdk.ChatSession` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `string` | `id` | Chat Session ID. |
| `Title` | `string` | `title` | Short AI-generated chat title. |
| `UserId` | `int64` | `user_id` | User ID. |
| `AiTaskId` | `int64` | `ai_task_id` | AI Task ID. Present when the conversation was started by an AI Task. |
| `WorkspaceId` | `int64` | `workspace_id` | Workspace ID. `0` means the default workspace. |
| `LastActiveAt` | `*time.Time` | `last_active_at` | Most recent chat activity date/time. |
| `CreatedAt` | `*time.Time` | `created_at` | Chat session creation date/time. |
| `Messages` | `[]files_sdk.ChatMessage` | `messages` | Visible conversation messages in this chat session. For performance reasons, this is not provided when listing chat sessions. |
