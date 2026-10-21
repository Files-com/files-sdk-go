# GpgKey

`files_sdk.GpgKey` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/gpgkey`, named `gpg_key`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	gpg_key "github.com/Files-com/files-sdk-go/v3/gpgkey"
)
```

## Operations

### List

`GET /gpg_keys`

List GPG Keys

```go
func (c *Client) List(params files_sdk.GpgKeyListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.GpgKeyListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.GpgKey` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.GpgKeyListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `UserId` | `int64` | `user_id` | no | User ID. Provide a value of `0` to operate the current session's user. |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id`, `name` or `expires_at`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `workspace_id`, `partner_id` or `expires_at`. Valid field combinations are `[ workspace_id, expires_at ]`. |
| `FilterGt` | `interface{}` | `filter_gt` | no | If set, return records where the specified field is greater than the supplied value. Valid fields are `expires_at`. |
| `FilterGteq` | `interface{}` | `filter_gteq` | no | If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `expires_at`. |
| `FilterLt` | `interface{}` | `filter_lt` | no | If set, return records where the specified field is less than the supplied value. Valid fields are `expires_at`. |
| `FilterLteq` | `interface{}` | `filter_lteq` | no | If set, return records where the specified field is less than or equal the supplied value. Valid fields are `expires_at`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := gpg_key.List(files_sdk.GpgKeyListParams{})
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

`GET /gpg_keys/{id}`

Show GPG Key

```go
func (c *Client) Find(params files_sdk.GpgKeyFindParams, opts ...files_sdk.RequestResponseOption) (gpgKey files_sdk.GpgKey, err error)
func Find(params files_sdk.GpgKeyFindParams, opts ...files_sdk.RequestResponseOption) (gpgKey files_sdk.GpgKey, err error)
```

Parameters, `files_sdk.GpgKeyFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Gpg Key ID. |

Example:

```go
result, err := gpg_key.Find(files_sdk.GpgKeyFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Create

`POST /gpg_keys`

Create GPG Key

```go
func (c *Client) Create(params files_sdk.GpgKeyCreateParams, opts ...files_sdk.RequestResponseOption) (gpgKey files_sdk.GpgKey, err error)
func Create(params files_sdk.GpgKeyCreateParams, opts ...files_sdk.RequestResponseOption) (gpgKey files_sdk.GpgKey, err error)
```

Parameters, `files_sdk.GpgKeyCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `UserId` | `int64` | `user_id` | no | User ID. Provide a value of `0` to operate the current session's user. |
| `PartnerId` | `int64` | `partner_id` | no | Partner ID who owns this GPG Key, if applicable. |
| `PublicKey` | `string` | `public_key` | no | The GPG public key |
| `PrivateKey` | `string` | `private_key` | no | The GPG private key |
| `PrivateKeyPassword` | `string` | `private_key_password` | no | The GPG private key password |
| `Name` | `string` | `name` | yes | GPG key name. |
| `WorkspaceId` | `int64` | `workspace_id` | no | Workspace ID (0 for default workspace). |
| `GenerateExpiresAt` | `*time.Time` | `generate_expires_at` | no | Expiration date of the key. Used for the generation of the key. Will be ignored if `generate_keypair` is false. |
| `GenerateKeypair` | `*bool` | `generate_keypair` | no | If true, generate a new GPG key pair. Can not be used with `public_key`/`private_key` |
| `GenerateFullName` | `string` | `generate_full_name` | no | Full name of the key owner. Used for the generation of the key. Will be ignored if `generate_keypair` is false. |
| `GenerateEmail` | `string` | `generate_email` | no | Email address of the key owner. Used for the generation of the key. Will be ignored if `generate_keypair` is false. |

Example:

```go
result, err := gpg_key.Create(files_sdk.GpgKeyCreateParams{
	Name: "key name",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Update

`PATCH /gpg_keys/{id}`

Update GPG Key

```go
func (c *Client) Update(params files_sdk.GpgKeyUpdateParams, opts ...files_sdk.RequestResponseOption) (gpgKey files_sdk.GpgKey, err error)
func Update(params files_sdk.GpgKeyUpdateParams, opts ...files_sdk.RequestResponseOption) (gpgKey files_sdk.GpgKey, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (gpgKey files_sdk.GpgKey, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (gpgKey files_sdk.GpgKey, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.GpgKeyUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Gpg Key ID. |
| `PartnerId` | `int64` | `partner_id` | no | Partner ID who owns this GPG Key, if applicable. |
| `PublicKey` | `string` | `public_key` | no | The GPG public key |
| `PrivateKey` | `string` | `private_key` | no | The GPG private key |
| `PrivateKeyPassword` | `string` | `private_key_password` | no | The GPG private key password |
| `Name` | `string` | `name` | no | GPG key name. |

Example:

```go
result, err := gpg_key.Update(files_sdk.GpgKeyUpdateParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := gpg_key.UpdateWithMap(map[string]interface{}{
	"id": 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /gpg_keys/{id}`

Delete GPG Key

```go
func (c *Client) Delete(params files_sdk.GpgKeyDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.GpgKeyDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.GpgKeyDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Gpg Key ID. |

Example:

```go
err := gpg_key.Delete(files_sdk.GpgKeyDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.GpgKey` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | GPG key ID. |
| `WorkspaceId` | `int64` | `workspace_id` | Workspace ID (0 for default workspace). |
| `ExpiresAt` | `*time.Time` | `expires_at` | GPG key expiration date. |
| `Name` | `string` | `name` | GPG key name. |
| `PartnerId` | `int64` | `partner_id` | Partner ID who owns this GPG Key, if applicable. |
| `PartnerName` | `string` | `partner_name` | Name of the Partner who owns this GPG Key, if applicable. |
| `UserId` | `int64` | `user_id` | User ID who owns this GPG Key, if applicable. |
| `PublicKeyMd5` | `string` | `public_key_md5` | MD5 hash of the GPG public key |
| `PrivateKeyMd5` | `string` | `private_key_md5` | MD5 hash of the GPG private key. |
| `GeneratedPublicKey` | `string` | `generated_public_key` | GPG public key |
| `GeneratedPrivateKey` | `string` | `generated_private_key` | GPG private key. |
| `PrivateKeyPasswordMd5` | `string` | `private_key_password_md5` | GPG private key password. Only required for password protected keys. |
| `PublicKey` | `string` | `public_key` | The GPG public key |
| `PrivateKey` | `string` | `private_key` | The GPG private key |
| `PrivateKeyPassword` | `string` | `private_key_password` | The GPG private key password |
| `GenerateExpiresAt` | `string` | `generate_expires_at` | Expiration date of the key. Used for the generation of the key. Will be ignored if `generate_keypair` is false. |
| `GenerateKeypair` | `*bool` | `generate_keypair` | If true, generate a new GPG key pair. Can not be used with `public_key`/`private_key` |
| `GenerateFullName` | `string` | `generate_full_name` | Full name of the key owner. Used for the generation of the key. Will be ignored if `generate_keypair` is false. |
| `GenerateEmail` | `string` | `generate_email` | Email address of the key owner. Used for the generation of the key. Will be ignored if `generate_keypair` is false. |
