# Session

`files_sdk.Session` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/session`, named `session`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/session"
)
```

## Operations

### Create

`POST /sessions`

Create user session (log in)

```go
func (c *Client) Create(params files_sdk.SessionCreateParams, opts ...files_sdk.RequestResponseOption) (session files_sdk.Session, err error)
func Create(params files_sdk.SessionCreateParams, opts ...files_sdk.RequestResponseOption) (session files_sdk.Session, err error)
```

Parameters, `files_sdk.SessionCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Username` | `string` | `username` | no | Username to sign in as |
| `Password` | `string` | `password` | no | Password for sign in |
| `Otp` | `string` | `otp` | no | If this user has a 2FA device, provide its OTP or code here. |
| `PartialSessionId` | `string` | `partial_session_id` | no | Identifier for a partially-completed login |

Example:

```go
result, err := session.Create(files_sdk.SessionCreateParams{})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /sessions`

Delete user session (log out)

```go
func (c *Client) Delete(opts ...files_sdk.RequestResponseOption) (err error)
func Delete(opts ...files_sdk.RequestResponseOption) (err error)
```

Example:

```go
err := session.Delete()
if err != nil {
	return err
}
```

## Model

`files_sdk.Session` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `string` | `id` | Session ID |
| `Language` | `string` | `language` | Session language |
| `ReadOnly` | `*bool` | `read_only` | Is this session read only? |
| `SftpInsecureCiphers` | `*bool` | `sftp_insecure_ciphers` | Are insecure SFTP ciphers allowed for this user? (If this is set to true, the site administrator has signaled that it is ok to use less secure SSH ciphers for this user.) |
| `Username` | `string` | `username` | Username to sign in as |
| `Password` | `string` | `password` | Password for sign in |
| `Otp` | `string` | `otp` | If this user has a 2FA device, provide its OTP or code here. |
| `PartialSessionId` | `string` | `partial_session_id` | Identifier for a partially-completed login |
