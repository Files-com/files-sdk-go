# FileCommentReaction

`files_sdk.FileCommentReaction` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/filecommentreaction`, named `file_comment_reaction`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	file_comment_reaction "github.com/Files-com/files-sdk-go/v3/filecommentreaction"
)
```

## Operations

### Create

`POST /file_comment_reactions`

Create File Comment Reaction

```go
func (c *Client) Create(params files_sdk.FileCommentReactionCreateParams, opts ...files_sdk.RequestResponseOption) (fileCommentReaction files_sdk.FileCommentReaction, err error)
func Create(params files_sdk.FileCommentReactionCreateParams, opts ...files_sdk.RequestResponseOption) (fileCommentReaction files_sdk.FileCommentReaction, err error)
```

Parameters, `files_sdk.FileCommentReactionCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `UserId` | `int64` | `user_id` | no | User ID. Provide a value of `0` to operate the current session's user. |
| `FileCommentId` | `int64` | `file_comment_id` | yes | ID of file comment to attach reaction to. |
| `Emoji` | `string` | `emoji` | yes | Emoji to react with. |

Example:

```go
result, err := file_comment_reaction.Create(files_sdk.FileCommentReactionCreateParams{
	FileCommentId: 1,
	Emoji:         "emoji",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /file_comment_reactions/{id}`

Delete File Comment Reaction

```go
func (c *Client) Delete(params files_sdk.FileCommentReactionDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.FileCommentReactionDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.FileCommentReactionDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | File Comment Reaction ID. |

Example:

```go
err := file_comment_reaction.Delete(files_sdk.FileCommentReactionDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.FileCommentReaction` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Reaction ID |
| `Emoji` | `string` | `emoji` | Emoji used in the reaction. |
| `UserId` | `int64` | `user_id` | User ID. Provide a value of `0` to operate the current session's user. |
| `FileCommentId` | `int64` | `file_comment_id` | ID of file comment to attach reaction to. |
