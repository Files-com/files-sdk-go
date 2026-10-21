# UserCipherUse

`files_sdk.UserCipherUse` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/usercipheruse`, named `user_cipher_use`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	user_cipher_use "github.com/Files-com/files-sdk-go/v3/usercipheruse"
)
```

## Operations

### List

`GET /user_cipher_uses`

List User Cipher Uses

```go
func (c *Client) List(params files_sdk.UserCipherUseListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.UserCipherUseListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.UserCipherUse` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.UserCipherUseListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `UserId` | `int64` | `user_id` | no | User ID. If provided, will return uses for this user. |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `updated_at`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `insecure` and `updated_at`. Valid field combinations are `[ insecure, updated_at ]`. |
| `FilterGt` | `interface{}` | `filter_gt` | no | If set, return records where the specified field is greater than the supplied value. Valid fields are `updated_at`. |
| `FilterGteq` | `interface{}` | `filter_gteq` | no | If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `updated_at`. |
| `FilterLt` | `interface{}` | `filter_lt` | no | If set, return records where the specified field is less than the supplied value. Valid fields are `updated_at`. |
| `FilterLteq` | `interface{}` | `filter_lteq` | no | If set, return records where the specified field is less than or equal the supplied value. Valid fields are `updated_at`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := user_cipher_use.List(files_sdk.UserCipherUseListParams{})
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

`files_sdk.UserCipherUse` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | UserCipherUse ID |
| `UserId` | `int64` | `user_id` | ID of the user who performed this access |
| `Username` | `string` | `username` | Username of the user who performed this access |
| `ProtocolCipher` | `string` | `protocol_cipher` | The protocol and cipher employed |
| `CreatedAt` | `*time.Time` | `created_at` | The earliest recorded use of this combination of interface and protocol and cipher (for this user) |
| `Insecure` | `*bool` | `insecure` | Is this cipher considered insecure? |
| `Interface` | `string` | `interface` | The interface accessed |
| `UpdatedAt` | `*time.Time` | `updated_at` | The most recent use of this combination of interface and protocol and cipher (for this user) |
