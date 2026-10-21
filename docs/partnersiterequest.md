# PartnerSiteRequest

`files_sdk.PartnerSiteRequest` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/partnersiterequest`, named `partner_site_request`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	partner_site_request "github.com/Files-com/files-sdk-go/v3/partnersiterequest"
)
```

## Operations

### List

`GET /partner_site_requests`

List Partner Site Requests

```go
func (c *Client) List(params files_sdk.PartnerSiteRequestListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.PartnerSiteRequestListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.PartnerSiteRequest` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.PartnerSiteRequestListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `host_partner_id`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `host_partner_id`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := partner_site_request.List(files_sdk.PartnerSiteRequestListParams{})
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

### FindByPairingKey

`GET /partner_site_requests/find_by_pairing_key`

Find partner site request by pairing key

```go
func (c *Client) FindByPairingKey(params files_sdk.PartnerSiteRequestFindByPairingKeyParams, opts ...files_sdk.RequestResponseOption) (err error)
func FindByPairingKey(params files_sdk.PartnerSiteRequestFindByPairingKeyParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.PartnerSiteRequestFindByPairingKeyParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `PairingKey` | `string` | `pairing_key` | yes | Pairing key for the partner site request |

Example:

```go
err := partner_site_request.FindByPairingKey(files_sdk.PartnerSiteRequestFindByPairingKeyParams{
	PairingKey: "pairing_key",
})
if err != nil {
	return err
}
```

### Create

`POST /partner_site_requests`

Create Partner Site Request

```go
func (c *Client) Create(params files_sdk.PartnerSiteRequestCreateParams, opts ...files_sdk.RequestResponseOption) (partnerSiteRequest files_sdk.PartnerSiteRequest, err error)
func Create(params files_sdk.PartnerSiteRequestCreateParams, opts ...files_sdk.RequestResponseOption) (partnerSiteRequest files_sdk.PartnerSiteRequest, err error)
```

Parameters, `files_sdk.PartnerSiteRequestCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `HostPartnerId` | `int64` | `host_partner_id` | yes | Host Partner ID to link with |
| `GuestSiteUrl` | `string` | `guest_site_url` | yes | Guest Site URL to link to |

Example:

```go
result, err := partner_site_request.Create(files_sdk.PartnerSiteRequestCreateParams{
	HostPartnerId: 1,
	GuestSiteUrl:  "guest_site_url",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Reject

`POST /partner_site_requests/reject`

Reject partner site request

```go
func (c *Client) Reject(params files_sdk.PartnerSiteRequestRejectParams, opts ...files_sdk.RequestResponseOption) (err error)
func Reject(params files_sdk.PartnerSiteRequestRejectParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.PartnerSiteRequestRejectParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `PairingKey` | `string` | `pairing_key` | yes | Pairing key for the partner site request |

Example:

```go
err := partner_site_request.Reject(files_sdk.PartnerSiteRequestRejectParams{
	PairingKey: "pairing_key",
})
if err != nil {
	return err
}
```

### Approve

`POST /partner_site_requests/approve`

Approve partner site request

```go
func (c *Client) Approve(params files_sdk.PartnerSiteRequestApproveParams, opts ...files_sdk.RequestResponseOption) (err error)
func Approve(params files_sdk.PartnerSiteRequestApproveParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.PartnerSiteRequestApproveParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `PairingKey` | `string` | `pairing_key` | yes | Pairing key for the partner site request |
| `PartnerId` | `int64` | `partner_id` | no | ID of an existing Partner on this site, with the host role, that represents the requesting organization. The connection binds to that Partner and makes it host_and_guest. When omitted, a guest Partner named after the host site is created. |

Example:

```go
err := partner_site_request.Approve(files_sdk.PartnerSiteRequestApproveParams{
	PairingKey: "pairing_key",
})
if err != nil {
	return err
}
```

### Delete

`DELETE /partner_site_requests/{id}`

Delete Partner Site Request

```go
func (c *Client) Delete(params files_sdk.PartnerSiteRequestDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.PartnerSiteRequestDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.PartnerSiteRequestDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Partner Site Request ID. |

Example:

```go
err := partner_site_request.Delete(files_sdk.PartnerSiteRequestDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.PartnerSiteRequest` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Partner Site Request ID |
| `HostPartnerId` | `int64` | `host_partner_id` | Host Partner ID |
| `GuestSiteUrl` | `string` | `guest_site_url` | Guest Site URL |
| `Status` | `string` | `status` | Request status (pending, approved, rejected) |
| `HostSiteName` | `string` | `host_site_name` | Host Site Name |
| `PairingKey` | `string` | `pairing_key` | Pairing key used to approve this request on the Guest Site |
| `CreatedAt` | `*time.Time` | `created_at` | Request creation date/time |
| `UpdatedAt` | `*time.Time` | `updated_at` | Request last updated date/time |
