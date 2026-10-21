# SftpHostKey

`files_sdk.SftpHostKey` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/sftphostkey`, named `sftp_host_key`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	sftp_host_key "github.com/Files-com/files-sdk-go/v3/sftphostkey"
)
```

## Operations

### List

`GET /sftp_host_keys`

List SFTP Host Keys

```go
func (c *Client) List(params files_sdk.SftpHostKeyListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.SftpHostKeyListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.SftpHostKey` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.SftpHostKeyListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := sftp_host_key.List(files_sdk.SftpHostKeyListParams{})
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

`GET /sftp_host_keys/{id}`

Show SFTP Host Key

```go
func (c *Client) Find(params files_sdk.SftpHostKeyFindParams, opts ...files_sdk.RequestResponseOption) (sftpHostKey files_sdk.SftpHostKey, err error)
func Find(params files_sdk.SftpHostKeyFindParams, opts ...files_sdk.RequestResponseOption) (sftpHostKey files_sdk.SftpHostKey, err error)
```

Parameters, `files_sdk.SftpHostKeyFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Sftp Host Key ID. |

Example:

```go
result, err := sftp_host_key.Find(files_sdk.SftpHostKeyFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Create

`POST /sftp_host_keys`

Create SFTP Host Key

```go
func (c *Client) Create(params files_sdk.SftpHostKeyCreateParams, opts ...files_sdk.RequestResponseOption) (sftpHostKey files_sdk.SftpHostKey, err error)
func Create(params files_sdk.SftpHostKeyCreateParams, opts ...files_sdk.RequestResponseOption) (sftpHostKey files_sdk.SftpHostKey, err error)
```

Parameters, `files_sdk.SftpHostKeyCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Active` | `*bool` | `active` | no | If true, use this SFTP Host Key. |
| `CustomDomainId` | `int64` | `custom_domain_id` | no | Custom Domain ID. If set, this key is used only for that Custom Domain. |
| `Name` | `string` | `name` | no | The friendly name of this SFTP Host Key. |
| `PrivateKey` | `string` | `private_key` | no | The private key data. |

Example:

```go
result, err := sftp_host_key.Create(files_sdk.SftpHostKeyCreateParams{})
if err != nil {
	return err
}
fmt.Println(result)
```

### Update

`PATCH /sftp_host_keys/{id}`

Update SFTP Host Key

```go
func (c *Client) Update(params files_sdk.SftpHostKeyUpdateParams, opts ...files_sdk.RequestResponseOption) (sftpHostKey files_sdk.SftpHostKey, err error)
func Update(params files_sdk.SftpHostKeyUpdateParams, opts ...files_sdk.RequestResponseOption) (sftpHostKey files_sdk.SftpHostKey, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (sftpHostKey files_sdk.SftpHostKey, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (sftpHostKey files_sdk.SftpHostKey, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.SftpHostKeyUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Sftp Host Key ID. |
| `Active` | `*bool` | `active` | no | If true, use this SFTP Host Key. |
| `CustomDomainId` | `int64` | `custom_domain_id` | no | Custom Domain ID. If set, this key is used only for that Custom Domain. |
| `Name` | `string` | `name` | no | The friendly name of this SFTP Host Key. |
| `PrivateKey` | `string` | `private_key` | no | The private key data. |

Example:

```go
result, err := sftp_host_key.Update(files_sdk.SftpHostKeyUpdateParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := sftp_host_key.UpdateWithMap(map[string]interface{}{
	"id": 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /sftp_host_keys/{id}`

Delete SFTP Host Key

```go
func (c *Client) Delete(params files_sdk.SftpHostKeyDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.SftpHostKeyDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.SftpHostKeyDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Sftp Host Key ID. |

Example:

```go
err := sftp_host_key.Delete(files_sdk.SftpHostKeyDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.SftpHostKey` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Active` | `*bool` | `active` | If true, use this SFTP Host Key. |
| `CustomDomainId` | `int64` | `custom_domain_id` | Custom Domain ID. If set, this key is used only for that Custom Domain. |
| `Id` | `int64` | `id` | SFTP Host Key ID |
| `Name` | `string` | `name` | The friendly name of this SFTP Host Key. |
| `KeyType` | `string` | `key_type` | SSH key type |
| `FingerprintMd5` | `string` | `fingerprint_md5` | MD5 Fingerprint of the public key |
| `FingerprintSha256` | `string` | `fingerprint_sha256` | SHA256 Fingerprint of the public key |
| `PrivateKey` | `string` | `private_key` | The private key data. |
