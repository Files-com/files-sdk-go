# UserSftpClientUse

`files_sdk.UserSftpClientUse` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/usersftpclientuse`, named `user_sftp_client_use`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	user_sftp_client_use "github.com/Files-com/files-sdk-go/v3/usersftpclientuse"
)
```

## Operations

### List

`GET /user_sftp_client_uses`

List User SFTP Client Uses

```go
func (c *Client) List(params files_sdk.UserSftpClientUseListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.UserSftpClientUseListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.UserSftpClientUse` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.UserSftpClientUseListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `UserId` | `int64` | `user_id` | no | User ID. If provided, will return uses for this user. |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `created_at`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := user_sftp_client_use.List(files_sdk.UserSftpClientUseListParams{})
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

## Model

`files_sdk.UserSftpClientUse` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | UserSftpClientUse ID |
| `SftpClient` | `string` | `sftp_client` | The SFTP client used |
| `CreatedAt` | `*time.Time` | `created_at` | The earliest recorded use of this SFTP client (for this user) |
| `UpdatedAt` | `*time.Time` | `updated_at` | The most recent use of this SFTP client (for this user) |
| `UserId` | `int64` | `user_id` | ID of the user who performed this access |
