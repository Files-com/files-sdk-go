# As2Station

`files_sdk.As2Station` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/as2station`, named `as2_station`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	as2_station "github.com/Files-com/files-sdk-go/v3/as2station"
)
```

## Operations

### List

`GET /as2_stations`

List AS2 Stations

```go
func (c *Client) List(params files_sdk.As2StationListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.As2StationListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.As2Station` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.As2StationListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id` and `name`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `workspace_id`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := as2_station.List(files_sdk.As2StationListParams{})
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

`GET /as2_stations/{id}`

Show AS2 Station

```go
func (c *Client) Find(params files_sdk.As2StationFindParams, opts ...files_sdk.RequestResponseOption) (as2Station files_sdk.As2Station, err error)
func Find(params files_sdk.As2StationFindParams, opts ...files_sdk.RequestResponseOption) (as2Station files_sdk.As2Station, err error)
```

Parameters, `files_sdk.As2StationFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | As2 Station ID. |

Example:

```go
result, err := as2_station.Find(files_sdk.As2StationFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Create

`POST /as2_stations`

Create AS2 Station

```go
func (c *Client) Create(params files_sdk.As2StationCreateParams, opts ...files_sdk.RequestResponseOption) (as2Station files_sdk.As2Station, err error)
func Create(params files_sdk.As2StationCreateParams, opts ...files_sdk.RequestResponseOption) (as2Station files_sdk.As2Station, err error)
```

Parameters, `files_sdk.As2StationCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Name` | `string` | `name` | yes | The station's formal AS2 name. |
| `WorkspaceId` | `int64` | `workspace_id` | no | ID of the Workspace associated with this AS2 Station. |
| `PublicCertificate` | `string` | `public_certificate` | no | Public certificate used for message security. |
| `PrivateKey` | `string` | `private_key` | no | PEM-encoded private key matching public_certificate. |
| `PrivateKeyPassword` | `string` | `private_key_password` | no | Password for the PEM-encoded private key. |
| `Pkcs12` | `string` | `pkcs12` | no | Base64-encoded PKCS#12 (.pfx or .p12) file containing the identity's certificate and private key. Provide this instead of public_certificate and private_key. |
| `Pkcs12Password` | `string` | `pkcs12_password` | no | Password for pkcs12. The file and password are used only for import; the extracted certificate and private key are stored as PEM. |

Example:

```go
result, err := as2_station.Create(files_sdk.As2StationCreateParams{
	Name: "AS2 Station Name",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Update

`PATCH /as2_stations/{id}`

Update AS2 Station

```go
func (c *Client) Update(params files_sdk.As2StationUpdateParams, opts ...files_sdk.RequestResponseOption) (as2Station files_sdk.As2Station, err error)
func Update(params files_sdk.As2StationUpdateParams, opts ...files_sdk.RequestResponseOption) (as2Station files_sdk.As2Station, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (as2Station files_sdk.As2Station, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (as2Station files_sdk.As2Station, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.As2StationUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | As2 Station ID. |
| `Name` | `string` | `name` | no | The station's formal AS2 name. |
| `PublicCertificate` | `string` | `public_certificate` | no | Public certificate used for message security. |
| `PrivateKey` | `string` | `private_key` | no | PEM-encoded private key matching public_certificate. |
| `PrivateKeyPassword` | `string` | `private_key_password` | no | Password for the PEM-encoded private key. |
| `Pkcs12` | `string` | `pkcs12` | no | Base64-encoded PKCS#12 (.pfx or .p12) file containing the identity's replacement certificate and private key. Provide this instead of public_certificate and private_key. |
| `Pkcs12Password` | `string` | `pkcs12_password` | no | Password for pkcs12. The file and password are used only for import; the extracted certificate and private key are stored as PEM. |

Example:

```go
result, err := as2_station.Update(files_sdk.As2StationUpdateParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := as2_station.UpdateWithMap(map[string]interface{}{
	"id": 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /as2_stations/{id}`

Delete AS2 Station

```go
func (c *Client) Delete(params files_sdk.As2StationDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.As2StationDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.As2StationDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | As2 Station ID. |

Example:

```go
err := as2_station.Delete(files_sdk.As2StationDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.As2Station` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Id of the AS2 Station. |
| `WorkspaceId` | `int64` | `workspace_id` | ID of the Workspace associated with this AS2 Station. |
| `Name` | `string` | `name` | The station's formal AS2 name. |
| `Uri` | `string` | `uri` | Public URI for sending AS2 message to. |
| `Domain` | `string` | `domain` | The station's AS2 domain name. |
| `HexPublicCertificateSerial` | `string` | `hex_public_certificate_serial` | Serial of public certificate used for message security in hex format. |
| `PublicCertificateMd5` | `string` | `public_certificate_md5` | MD5 hash of public certificate used for message security. |
| `PublicCertificate` | `string` | `public_certificate` | Public certificate used for message security. |
| `PrivateKeyMd5` | `string` | `private_key_md5` | MD5 hash of private key used for message security. |
| `PublicCertificateSubject` | `string` | `public_certificate_subject` | Subject of public certificate used for message security. |
| `PublicCertificateIssuer` | `string` | `public_certificate_issuer` | Issuer of public certificate used for message security. |
| `PublicCertificateSerial` | `string` | `public_certificate_serial` | Serial of public certificate used for message security. |
| `PublicCertificateNotBefore` | `string` | `public_certificate_not_before` | Not before value of public certificate used for message security. |
| `PublicCertificateNotAfter` | `string` | `public_certificate_not_after` | Not after value of public certificate used for message security. |
| `PrivateKeyPasswordMd5` | `string` | `private_key_password_md5` | MD5 hash of private key password used for message security. |
| `PrivateKey` | `string` | `private_key` | PEM-encoded private key matching public_certificate. |
| `PrivateKeyPassword` | `string` | `private_key_password` | Password for the PEM-encoded private key. |
| `Pkcs12` | `string` | `pkcs12` | Base64-encoded PKCS#12 (.pfx or .p12) file containing the identity's certificate and private key. Provide this instead of public_certificate and private_key. |
| `Pkcs12Password` | `string` | `pkcs12_password` | Password for pkcs12. The file and password are used only for import; the extracted certificate and private key are stored as PEM. |
