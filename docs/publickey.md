# PublicKey

`files_sdk.PublicKey` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/publickey`, named `public_key`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	public_key "github.com/Files-com/files-sdk-go/v3/publickey"
)
```

## Operations

### List

`GET /public_keys`

List Public Keys

```go
func (c *Client) List(params files_sdk.PublicKeyListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.PublicKeyListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.PublicKey` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.PublicKeyListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `UserId` | `int64` | `user_id` | no | User ID. Provide a value of `0` to operate the current session's user. |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id`, `user_id`, `title` or `created_at`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `created_at` and `workspace_id`. Valid field combinations are `[ workspace_id, created_at ]`. |
| `FilterGt` | `interface{}` | `filter_gt` | no | If set, return records where the specified field is greater than the supplied value. Valid fields are `created_at`. |
| `FilterGteq` | `interface{}` | `filter_gteq` | no | If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `created_at`. |
| `FilterLt` | `interface{}` | `filter_lt` | no | If set, return records where the specified field is less than the supplied value. Valid fields are `created_at`. |
| `FilterLteq` | `interface{}` | `filter_lteq` | no | If set, return records where the specified field is less than or equal the supplied value. Valid fields are `created_at`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := public_key.List(files_sdk.PublicKeyListParams{})
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

`GET /public_keys/{id}`

Show Public Key

```go
func (c *Client) Find(params files_sdk.PublicKeyFindParams, opts ...files_sdk.RequestResponseOption) (publicKey files_sdk.PublicKey, err error)
func Find(params files_sdk.PublicKeyFindParams, opts ...files_sdk.RequestResponseOption) (publicKey files_sdk.PublicKey, err error)
```

Parameters, `files_sdk.PublicKeyFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Public Key ID. |

Example:

```go
result, err := public_key.Find(files_sdk.PublicKeyFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Create

`POST /public_keys`

Create Public Key

```go
func (c *Client) Create(params files_sdk.PublicKeyCreateParams, opts ...files_sdk.RequestResponseOption) (publicKey files_sdk.PublicKey, err error)
func Create(params files_sdk.PublicKeyCreateParams, opts ...files_sdk.RequestResponseOption) (publicKey files_sdk.PublicKey, err error)
```

Parameters, `files_sdk.PublicKeyCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `UserId` | `int64` | `user_id` | no | User ID. Provide a value of `0` to operate the current session's user. |
| `Title` | `string` | `title` | yes | Internal reference for key. |
| `PublicKey` | `string` | `public_key` | no | Actual contents of SSH key. |
| `GenerateKeypair` | `*bool` | `generate_keypair` | no | If true, generate a new SSH key pair. Can not be used with `public_key` |
| `GeneratePrivateKeyPassword` | `string` | `generate_private_key_password` | no | Password for the private key. Used for the generation of the key. Will be ignored if `generate_keypair` is false. |
| `GenerateAlgorithm` | `string` | `generate_algorithm` | no | Type of key to generate. One of rsa, dsa, ecdsa, ed25519. Used for the generation of the key. Will be ignored if `generate_keypair` is false. |
| `GenerateLength` | `int64` | `generate_length` | no | Length of key to generate. If algorithm is ecdsa, this is the signature size. Used for the generation of the key. Will be ignored if `generate_keypair` is false. |

Example:

```go
result, err := public_key.Create(files_sdk.PublicKeyCreateParams{
	Title: "My Main Key",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Update

`PATCH /public_keys/{id}`

Update Public Key

```go
func (c *Client) Update(params files_sdk.PublicKeyUpdateParams, opts ...files_sdk.RequestResponseOption) (publicKey files_sdk.PublicKey, err error)
func Update(params files_sdk.PublicKeyUpdateParams, opts ...files_sdk.RequestResponseOption) (publicKey files_sdk.PublicKey, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (publicKey files_sdk.PublicKey, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (publicKey files_sdk.PublicKey, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.PublicKeyUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Public Key ID. |
| `Title` | `string` | `title` | yes | Internal reference for key. |

Example:

```go
result, err := public_key.Update(files_sdk.PublicKeyUpdateParams{
	Id:    1,
	Title: "My Main Key",
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := public_key.UpdateWithMap(map[string]interface{}{
	"id":    1,
	"title": "My Main Key",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /public_keys/{id}`

Delete Public Key

```go
func (c *Client) Delete(params files_sdk.PublicKeyDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.PublicKeyDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.PublicKeyDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Public Key ID. |

Example:

```go
err := public_key.Delete(files_sdk.PublicKeyDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.PublicKey` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Public key ID |
| `WorkspaceId` | `int64` | `workspace_id` | Workspace ID (0 for default workspace). |
| `Title` | `string` | `title` | Public key title |
| `CreatedAt` | `*time.Time` | `created_at` | Public key created at date/time |
| `ExpiresAt` | `*time.Time` | `expires_at` | Public key expiration date/time |
| `Expired` | `*bool` | `expired` | Is this public key expired? |
| `Fingerprint` | `string` | `fingerprint` | Public key fingerprint (MD5) |
| `FingerprintSha256` | `string` | `fingerprint_sha256` | Public key fingerprint (SHA256) |
| `Status` | `string` | `status` | Only returned when generating keys. Can be invalid, not_generated, generating, complete |
| `LastLoginAt` | `*time.Time` | `last_login_at` | Key's most recent login time via SFTP |
| `GeneratedPrivateKey` | `string` | `generated_private_key` | Only returned when generating keys. Private key generated for the user. |
| `GeneratedPublicKey` | `string` | `generated_public_key` | Only returned when generating keys. Public key generated for the user. |
| `Username` | `string` | `username` | Username of the user this public key is associated with |
| `UserId` | `int64` | `user_id` | User ID this public key is associated with |
| `PublicKey` | `string` | `public_key` | Actual contents of SSH key. |
| `GenerateKeypair` | `*bool` | `generate_keypair` | If true, generate a new SSH key pair. Can not be used with `public_key` |
| `GeneratePrivateKeyPassword` | `string` | `generate_private_key_password` | Password for the private key. Used for the generation of the key. Will be ignored if `generate_keypair` is false. |
| `GenerateAlgorithm` | `string` | `generate_algorithm` | Type of key to generate. One of rsa, dsa, ecdsa, ed25519. Used for the generation of the key. Will be ignored if `generate_keypair` is false. |
| `GenerateLength` | `int64` | `generate_length` | Length of key to generate. If algorithm is ecdsa, this is the signature size. Used for the generation of the key. Will be ignored if `generate_keypair` is false. |
