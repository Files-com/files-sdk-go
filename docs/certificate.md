# Certificate

`files_sdk.Certificate` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/certificate`, named `certificate`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/certificate"
)
```

## Operations

### List

`GET /certificates`

List Certificates

```go
func (c *Client) List(params files_sdk.CertificateListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.CertificateListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.Certificate` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.CertificateListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := certificate.List(files_sdk.CertificateListParams{})
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

`GET /certificates/{id}`

Show Certificate

```go
func (c *Client) Find(params files_sdk.CertificateFindParams, opts ...files_sdk.RequestResponseOption) (certificate files_sdk.Certificate, err error)
func Find(params files_sdk.CertificateFindParams, opts ...files_sdk.RequestResponseOption) (certificate files_sdk.Certificate, err error)
```

Parameters, `files_sdk.CertificateFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Certificate ID. |

Example:

```go
result, err := certificate.Find(files_sdk.CertificateFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Create

`POST /certificates`

Create CSR or import existing SSL Certificate

```go
func (c *Client) Create(params files_sdk.CertificateCreateParams, opts ...files_sdk.RequestResponseOption) (certificate files_sdk.Certificate, err error)
func Create(params files_sdk.CertificateCreateParams, opts ...files_sdk.RequestResponseOption) (certificate files_sdk.Certificate, err error)
```

Parameters, `files_sdk.CertificateCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Name` | `string` | `name` | yes | Internal name of the SSL certificate. |
| `CertificateDomain` | `string` | `certificate_domain` | no | Domain for certificate. |
| `CertificateExtraDomains` | `[]string` | `certificate_extra_domains` | no | Additional domains for the certificate CSR. |
| `CertificateCountry` | `string` | `certificate_country` | no | Country. |
| `CertificateStateOrProvince` | `string` | `certificate_state_or_province` | no | State or province. |
| `CertificateCityOrLocale` | `string` | `certificate_city_or_locale` | no | City or locale. |
| `CertificateCompanyName` | `string` | `certificate_company_name` | no | Company name. |
| `CsrOu1` | `string` | `csr_ou1` | no | Department name or organization unit 1 |
| `CsrOu2` | `string` | `csr_ou2` | no | Department name or organization unit 2 |
| `CsrOu3` | `string` | `csr_ou3` | no | Department name or organization unit 3 |
| `CertificateEmailAddress` | `string` | `certificate_email_address` | no | Email address for the certificate owner. |
| `KeyType` | `string` | `key_type` | no | Any supported key type. Defaults to `RSA-4096`. |
| `Certificate` | `string` | `certificate` | no | The certificate. PEM or PKCS7 formats accepted. |
| `PrivateKey` | `string` | `private_key` | no | Certificate private key. |
| `Password` | `string` | `password` | no | Certificate password. |
| `Intermediates` | `string` | `intermediates` | no | Intermediate certificates. PEM or PKCS7 formats accepted. |

Example:

```go
result, err := certificate.Create(files_sdk.CertificateCreateParams{
	Name: "Name",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Activate

`POST /certificates/{id}/activate`

Activate SSL Certificate

```go
func (c *Client) Activate(params files_sdk.CertificateActivateParams, opts ...files_sdk.RequestResponseOption) (err error)
func Activate(params files_sdk.CertificateActivateParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.CertificateActivateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Certificate ID. |
| `ReplaceCert` | `string` | `replace_cert` | no | Leave blank or set to `any` to replace any currently active certificate, if applicable. Set to `new_ips` to activate this certificate as an additional certificate on your Site by allocating a new set of Dedicated IPs (may require a Plan upgrade). Set to the ID of a currently active certificate to replace that certificate on its set of dedicated IPs. A certificate belonging to a Custom Domain activates only for that domain: leave blank, use `any`, or specify that domain's current certificate ID. |

Example:

```go
err := certificate.Activate(files_sdk.CertificateActivateParams{
	Id: 1,
})
if err != nil {
	return err
}
```

### Deactivate

`POST /certificates/{id}/deactivate`

Deactivate SSL Certificate

```go
func (c *Client) Deactivate(params files_sdk.CertificateDeactivateParams, opts ...files_sdk.RequestResponseOption) (err error)
func Deactivate(params files_sdk.CertificateDeactivateParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.CertificateDeactivateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Certificate ID. |

Example:

```go
err := certificate.Deactivate(files_sdk.CertificateDeactivateParams{
	Id: 1,
})
if err != nil {
	return err
}
```

### ReplaceWithFiles

`POST /certificates/{id}/replace_with_files_certificate`

Replace this active custom SSL Certificate for the Primary Custom Domain with a Files.com SSL Certificate. The Primary Custom Domain's CNAME record must point at the Site's subdomain. The Files.com certificate is issued and activated in the background

```go
func (c *Client) ReplaceWithFiles(params files_sdk.CertificateReplaceWithFilesParams, opts ...files_sdk.RequestResponseOption) (err error)
func ReplaceWithFiles(params files_sdk.CertificateReplaceWithFilesParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.CertificateReplaceWithFilesParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Certificate ID. |

Example:

```go
err := certificate.ReplaceWithFiles(files_sdk.CertificateReplaceWithFilesParams{
	Id: 1,
})
if err != nil {
	return err
}
```

### Update

`PATCH /certificates/{id}`

Update Certificate

```go
func (c *Client) Update(params files_sdk.CertificateUpdateParams, opts ...files_sdk.RequestResponseOption) (certificate files_sdk.Certificate, err error)
func Update(params files_sdk.CertificateUpdateParams, opts ...files_sdk.RequestResponseOption) (certificate files_sdk.Certificate, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (certificate files_sdk.Certificate, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (certificate files_sdk.Certificate, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.CertificateUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Certificate ID. |
| `Name` | `string` | `name` | no | Internal certificate name. |
| `Intermediates` | `string` | `intermediates` | no | Any intermediate certificates. PEM or PKCS7 formats accepted. |
| `Certificate` | `string` | `certificate` | no | The certificate. PEM or PKCS7 formats accepted. |

Example:

```go
result, err := certificate.Update(files_sdk.CertificateUpdateParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := certificate.UpdateWithMap(map[string]interface{}{
	"id": 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /certificates/{id}`

Delete Certificate

```go
func (c *Client) Delete(params files_sdk.CertificateDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.CertificateDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.CertificateDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Certificate ID. |

Example:

```go
err := certificate.Delete(files_sdk.CertificateDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.Certificate` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Certificate ID |
| `Name` | `string` | `name` | Descriptive name of certificate |
| `CustomDomainId` | `int64` | `custom_domain_id` | Custom Domain that owns this certificate, including while the certificate is inactive. |
| `Certificate` | `string` | `certificate` | Full text of SSL certificate |
| `CreatedAt` | `*time.Time` | `created_at` | Certificate created at |
| `DisplayStatus` | `string` | `display_status` | Certificate status. (One of `Request Generated`, `New`, `Active`, `Active/Expired`, `Expired`, or `Available`) |
| `Domains` | `[]string` | `domains` | Domains on this certificate other than files.com domains |
| `ExpiresAt` | `*time.Time` | `expires_at` | Certificate expiration date/time |
| `BrickManaged` | `*bool` | `brick_managed` | Is this certificate automatically managed and renewed by Files.com? |
| `Intermediates` | `string` | `intermediates` | Intermediate certificates |
| `IpAddresses` | `[]string` | `ip_addresses` | A list of IP addresses associated with this SSL Certificate |
| `Issuer` | `string` | `issuer` | X509 issuer |
| `KeyType` | `string` | `key_type` | Type of key |
| `Request` | `string` | `request` | Certificate signing request text |
| `Status` | `string` | `status` | Certificate status (Request Generated, New, Active, Active/Expired, Expired, or Available) |
| `Subject` | `string` | `subject` | X509 Subject name |
| `CertificateDomain` | `string` | `certificate_domain` | Domain for certificate. |
| `CertificateExtraDomains` | `[]string` | `certificate_extra_domains` | Additional domains for the certificate CSR. |
| `CertificateCountry` | `string` | `certificate_country` | Country. |
| `CertificateStateOrProvince` | `string` | `certificate_state_or_province` | State or province. |
| `CertificateCityOrLocale` | `string` | `certificate_city_or_locale` | City or locale. |
| `CertificateCompanyName` | `string` | `certificate_company_name` | Company name. |
| `CsrOu1` | `string` | `csr_ou1` | Department name or organization unit 1 |
| `CsrOu2` | `string` | `csr_ou2` | Department name or organization unit 2 |
| `CsrOu3` | `string` | `csr_ou3` | Department name or organization unit 3 |
| `CertificateEmailAddress` | `string` | `certificate_email_address` | Email address for the certificate owner. |
| `PrivateKey` | `string` | `private_key` | Certificate private key. |
| `Password` | `string` | `password` | Certificate password. |
