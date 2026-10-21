# Lock

`files_sdk.Lock` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/lock`, named `lock`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/lock"
)
```

## Operations

### ListFor

`GET /locks/{path}`

List Locks by Path

```go
func (c *Client) ListFor(params files_sdk.LockListForParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func ListFor(params files_sdk.LockListForParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.Lock` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.LockListForParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Path` | `string` | `path`, in the path | yes | Path to operate on. |
| `IncludeChildren` | `*bool` | `include_children` | no | Include locks from children objects? |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := lock.ListFor(files_sdk.LockListForParams{
	Path: "path",
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

### Create

`POST /locks/{path}`

Create Lock

```go
func (c *Client) Create(params files_sdk.LockCreateParams, opts ...files_sdk.RequestResponseOption) (lock files_sdk.Lock, err error)
func Create(params files_sdk.LockCreateParams, opts ...files_sdk.RequestResponseOption) (lock files_sdk.Lock, err error)
```

Parameters, `files_sdk.LockCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Path` | `string` | `path`, in the path | yes | Path |
| `Token` | `string` | `token` | no | Lock token. With expected_token, use the same value to refresh or a different value to replace the existing token. |
| `ExpectedToken` | `string` | `expected_token` | no | Require this existing, unexpired token before refreshing or replacing a lock. Set token to the same value to refresh, or a different value to replace. |
| `AllowAccessByAnyUser` | `*bool` | `allow_access_by_any_user` | no | Can lock be modified by users other than its creator? |
| `Exclusive` | `*bool` | `exclusive` | no | Is lock exclusive? |
| `Recursive` | `*bool` | `recursive` | no | Does lock apply to subfolders? |
| `Owner` | `string` | `owner` | no | Arbitrary descriptive label for the lock. Does not change the lock creator or permissions. |
| `Timeout` | `int64` | `timeout` | no | Lock timeout in seconds |

Example:

```go
result, err := lock.Create(files_sdk.LockCreateParams{
	Path: "locked_file",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /locks/{path}`

Delete Lock

```go
func (c *Client) Delete(params files_sdk.LockDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.LockDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.LockDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Path` | `string` | `path`, in the path | yes | Path |
| `Token` | `string` | `token` | yes | Lock token |

Example:

```go
err := lock.Delete(files_sdk.LockDeleteParams{
	Path:  "locked_file",
	Token: "token",
})
if err != nil {
	return err
}
```

## Model

`files_sdk.Lock` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Path` | `string` | `path` | Path. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters. |
| `Timeout` | `int64` | `timeout` | Lock timeout in seconds |
| `Depth` | `string` | `depth` |  |
| `Recursive` | `*bool` | `recursive` | Does lock apply to subfolders? |
| `Owner` | `string` | `owner` | Arbitrary descriptive label for the lock. Does not change the lock creator or permissions. |
| `Scope` | `string` | `scope` |  |
| `Exclusive` | `*bool` | `exclusive` | Is lock exclusive? |
| `Token` | `string` | `token` | Lock token. Use to release lock. |
| `Type` | `string` | `type` |  |
| `AllowAccessByAnyUser` | `*bool` | `allow_access_by_any_user` | Can lock be modified by users other than its creator? |
| `UserId` | `int64` | `user_id` | Lock creator user ID |
| `Username` | `string` | `username` | Lock creator username |
| `ExpectedToken` | `string` | `expected_token` | Require this existing, unexpired token before refreshing or replacing a lock. Set token to the same value to refresh, or a different value to replace. |
