# As2Partner

`files_sdk.As2Partner` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/as2partner`, named `as2_partner`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	as2_partner "github.com/Files-com/files-sdk-go/v3/as2partner"
)
```

## Operations

### List

`GET /as2_partners`

List AS2 Partners

```go
func (c *Client) List(params files_sdk.As2PartnerListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.As2PartnerListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.As2Partner` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.As2PartnerListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id`, `as2_station_id` or `name`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `as2_station_id` and `workspace_id`. Valid field combinations are `[ workspace_id, as2_station_id ]`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := as2_partner.List(files_sdk.As2PartnerListParams{})
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

`GET /as2_partners/{id}`

Show AS2 Partner

```go
func (c *Client) Find(params files_sdk.As2PartnerFindParams, opts ...files_sdk.RequestResponseOption) (as2Partner files_sdk.As2Partner, err error)
func Find(params files_sdk.As2PartnerFindParams, opts ...files_sdk.RequestResponseOption) (as2Partner files_sdk.As2Partner, err error)
```

Parameters, `files_sdk.As2PartnerFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | As2 Partner ID. |

Example:

```go
result, err := as2_partner.Find(files_sdk.As2PartnerFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Create

`POST /as2_partners`

Create AS2 Partner

```go
func (c *Client) Create(params files_sdk.As2PartnerCreateParams, opts ...files_sdk.RequestResponseOption) (as2Partner files_sdk.As2Partner, err error)
func Create(params files_sdk.As2PartnerCreateParams, opts ...files_sdk.RequestResponseOption) (as2Partner files_sdk.As2Partner, err error)
```

Parameters, `files_sdk.As2PartnerCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `EnableDedicatedIps` | `*bool` | `enable_dedicated_ips` | no | If `true`, we will use your site's dedicated IPs for all outbound connections to this AS2 Partner. |
| `HttpAuthUsername` | `string` | `http_auth_username` | no | Username to send to server for HTTP Authentication. |
| `HttpAuthPassword` | `string` | `http_auth_password` | no | Password to send to server for HTTP Authentication. |
| `MdnValidationLevel` | `files_sdk.As2PartnerMdnValidationLevelEnum` | `mdn_validation_level` | no | How should Files.com evaluate message transfer success based on a partner's MDN response? This setting does not affect MDN storage; all MDNs received from a partner are always stored. `none`: MDN is stored for informational purposes only, a successful HTTPS transfer is a successful AS2 transfer. `weak`: Inspect the MDN for MIC and Disposition only. `normal`: `weak` plus validate MDN signature matches body, `strict`: `normal` but do not allow signatures from self-signed or incorrectly purposed certificates. `auto`: Automatically set the correct value for this setting based on next mdn received. Values: `none`, `weak`, `normal`, `strict`, `auto`. |
| `SignatureValidationLevel` | `files_sdk.As2PartnerSignatureValidationLevelEnum` | `signature_validation_level` | no | Should Files.com require signatures on incoming AS2 messages? `normal`: require that incoming messages are signed with a valid matching signature. `none`: Unsigned incoming messages are allowed. `auto`: A one-time starting state. The next message received is always accepted and switches this to `normal` if it carries a valid signature, or `none` if it does not. Values: `normal`, `none`, `auto`. |
| `ServerCertificate` | `files_sdk.As2PartnerServerCertificateEnum` | `server_certificate` | no | Should we require that the remote HTTP server have a valid SSL Certificate for HTTPS? (This only applies to Outgoing AS2 message from Files.com to a Partner.) Values: `require_match`, `allow_any`. |
| `DefaultMimeType` | `string` | `default_mime_type` | no | Default mime type of the file attached to the encrypted message |
| `AdditionalHttpHeaders` | `interface{}` | `additional_http_headers` | no | Additional HTTP Headers for outgoing message sent to this partner. |
| `As2StationId` | `int64` | `as2_station_id` | yes | ID of the AS2 Station associated with this partner. |
| `Name` | `string` | `name` | yes | The partner's formal AS2 name. |
| `Uri` | `string` | `uri` | yes | Public URI where we will send the AS2 messages (via HTTP/HTTPS). |
| `PublicCertificate` | `string` | `public_certificate` | no | Public certificate for AS2 Partner. Note: This is the certificate for AS2 message security, not a certificate used for HTTPS authentication. |
| `Pkcs12` | `string` | `pkcs12` | no | Base64-encoded PKCS#12 (.pfx or .p12) file containing the partner's public certificate. Provide this instead of public_certificate. Any private key in the file is discarded. A public-only file must contain exactly one certificate. |
| `Pkcs12Password` | `string` | `pkcs12_password` | no | Password for pkcs12. The file and password are used only for import; the extracted public certificate is stored as PEM. |

Example:

```go
result, err := as2_partner.Create(files_sdk.As2PartnerCreateParams{
	As2StationId: 1,
	Name:         "AS2 Partner Name",
	Uri:          "example",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Update

`PATCH /as2_partners/{id}`

Update AS2 Partner

```go
func (c *Client) Update(params files_sdk.As2PartnerUpdateParams, opts ...files_sdk.RequestResponseOption) (as2Partner files_sdk.As2Partner, err error)
func Update(params files_sdk.As2PartnerUpdateParams, opts ...files_sdk.RequestResponseOption) (as2Partner files_sdk.As2Partner, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (as2Partner files_sdk.As2Partner, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (as2Partner files_sdk.As2Partner, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.As2PartnerUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | As2 Partner ID. |
| `EnableDedicatedIps` | `*bool` | `enable_dedicated_ips` | no | If `true`, we will use your site's dedicated IPs for all outbound connections to this AS2 Partner. |
| `HttpAuthUsername` | `string` | `http_auth_username` | no | Username to send to server for HTTP Authentication. |
| `HttpAuthPassword` | `string` | `http_auth_password` | no | Password to send to server for HTTP Authentication. |
| `MdnValidationLevel` | `files_sdk.As2PartnerMdnValidationLevelEnum` | `mdn_validation_level` | no | How should Files.com evaluate message transfer success based on a partner's MDN response? This setting does not affect MDN storage; all MDNs received from a partner are always stored. `none`: MDN is stored for informational purposes only, a successful HTTPS transfer is a successful AS2 transfer. `weak`: Inspect the MDN for MIC and Disposition only. `normal`: `weak` plus validate MDN signature matches body, `strict`: `normal` but do not allow signatures from self-signed or incorrectly purposed certificates. `auto`: Automatically set the correct value for this setting based on next mdn received. Values: `none`, `weak`, `normal`, `strict`, `auto`. |
| `SignatureValidationLevel` | `files_sdk.As2PartnerSignatureValidationLevelEnum` | `signature_validation_level` | no | Should Files.com require signatures on incoming AS2 messages? `normal`: require that incoming messages are signed with a valid matching signature. `none`: Unsigned incoming messages are allowed. `auto`: A one-time starting state. The next message received is always accepted and switches this to `normal` if it carries a valid signature, or `none` if it does not. Values: `normal`, `none`, `auto`. |
| `ServerCertificate` | `files_sdk.As2PartnerServerCertificateEnum` | `server_certificate` | no | Should we require that the remote HTTP server have a valid SSL Certificate for HTTPS? (This only applies to Outgoing AS2 message from Files.com to a Partner.) Values: `require_match`, `allow_any`. |
| `DefaultMimeType` | `string` | `default_mime_type` | no | Default mime type of the file attached to the encrypted message |
| `AdditionalHttpHeaders` | `interface{}` | `additional_http_headers` | no | Additional HTTP Headers for outgoing message sent to this partner. |
| `Name` | `string` | `name` | no | The partner's formal AS2 name. |
| `Uri` | `string` | `uri` | no | Public URI where we will send the AS2 messages (via HTTP/HTTPS). |
| `PublicCertificate` | `string` | `public_certificate` | no | Public certificate for AS2 Partner. Note: This is the certificate for AS2 message security, not a certificate used for HTTPS authentication. |
| `Pkcs12` | `string` | `pkcs12` | no | Base64-encoded PKCS#12 (.pfx or .p12) file containing the partner's public certificate. Provide this instead of public_certificate. Any private key in the file is discarded. A public-only file must contain exactly one certificate. |
| `Pkcs12Password` | `string` | `pkcs12_password` | no | Password for pkcs12. The file and password are used only for import; the extracted public certificate is stored as PEM. |

Example:

```go
result, err := as2_partner.Update(files_sdk.As2PartnerUpdateParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := as2_partner.UpdateWithMap(map[string]interface{}{
	"id": 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /as2_partners/{id}`

Delete AS2 Partner

```go
func (c *Client) Delete(params files_sdk.As2PartnerDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.As2PartnerDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.As2PartnerDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | As2 Partner ID. |

Example:

```go
err := as2_partner.Delete(files_sdk.As2PartnerDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.As2Partner` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | ID of the AS2 Partner. |
| `WorkspaceId` | `int64` | `workspace_id` | ID of the Workspace associated with this AS2 Partner. |
| `As2StationId` | `int64` | `as2_station_id` | ID of the AS2 Station associated with this partner. |
| `Name` | `string` | `name` | The partner's formal AS2 name. |
| `Uri` | `string` | `uri` | Public URI where we will send the AS2 messages (via HTTP/HTTPS). |
| `ServerCertificate` | `string` | `server_certificate` | Should we require that the remote HTTP server have a valid SSL Certificate for HTTPS? (This only applies to Outgoing AS2 message from Files.com to a Partner.) |
| `HttpAuthUsername` | `string` | `http_auth_username` | Username to send to server for HTTP Authentication. |
| `AdditionalHttpHeaders` | `interface{}` | `additional_http_headers` | Additional HTTP Headers for outgoing message sent to this partner. |
| `DefaultMimeType` | `string` | `default_mime_type` | Default mime type of the file attached to the encrypted message |
| `MdnValidationLevel` | `string` | `mdn_validation_level` | How should Files.com evaluate message transfer success based on a partner's MDN response? This setting does not affect MDN storage; all MDNs received from a partner are always stored. `none`: MDN is stored for informational purposes only, a successful HTTPS transfer is a successful AS2 transfer. `weak`: Inspect the MDN for MIC and Disposition only. `normal`: `weak` plus validate MDN signature matches body, `strict`: `normal` but do not allow signatures from self-signed or incorrectly purposed certificates. `auto`: Automatically set the correct value for this setting based on next mdn received. |
| `SignatureValidationLevel` | `string` | `signature_validation_level` | Should Files.com require signatures on incoming AS2 messages? `normal`: require that incoming messages are signed with a valid matching signature. `none`: Unsigned incoming messages are allowed. `auto`: A one-time starting state. The next message received is always accepted and switches this to `normal` if it carries a valid signature, or `none` if it does not. |
| `EnableDedicatedIps` | `*bool` | `enable_dedicated_ips` | If `true`, we will use your site's dedicated IPs for all outbound connections to this AS2 Partner. |
| `HexPublicCertificateSerial` | `string` | `hex_public_certificate_serial` | Serial of public certificate used for message security in hex format. |
| `PublicCertificate` | `string` | `public_certificate` | Public certificate used for message security. |
| `PublicCertificateMd5` | `string` | `public_certificate_md5` | MD5 hash of public certificate used for message security. |
| `PublicCertificateSubject` | `string` | `public_certificate_subject` | Subject of public certificate used for message security. |
| `PublicCertificateIssuer` | `string` | `public_certificate_issuer` | Issuer of public certificate used for message security. |
| `PublicCertificateSerial` | `string` | `public_certificate_serial` | Serial of public certificate used for message security. |
| `PublicCertificateNotBefore` | `string` | `public_certificate_not_before` | Not before value of public certificate used for message security. |
| `PublicCertificateNotAfter` | `string` | `public_certificate_not_after` | Not after value of public certificate used for message security. |
| `HttpAuthPassword` | `string` | `http_auth_password` | Password to send to server for HTTP Authentication. |
| `Pkcs12` | `string` | `pkcs12` | Base64-encoded PKCS#12 (.pfx or .p12) file containing the partner's public certificate. Provide this instead of public_certificate. Any private key in the file is discarded. A public-only file must contain exactly one certificate. |
| `Pkcs12Password` | `string` | `pkcs12_password` | Password for pkcs12. The file and password are used only for import; the extracted public certificate is stored as PEM. |
