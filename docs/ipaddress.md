# IpAddress

`files_sdk.IpAddress` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/ipaddress`, named `ip_address`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	ip_address "github.com/Files-com/files-sdk-go/v3/ipaddress"
)
```

## Operations

### List

`GET /ip_addresses`

List IP Addresses associated with the current site

```go
func (c *Client) List(params files_sdk.IpAddressListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.IpAddressListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields each of this listing's `files_sdk.IpAddress` values as an `any` value, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.IpAddressListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := ip_address.List(files_sdk.IpAddressListParams{})
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

### GetSmartfileReserved

`GET /ip_addresses/smartfile-reserved`

List all possible public SmartFile IP addresses

```go
func (c *Client) GetSmartfileReserved(params files_sdk.IpAddressGetSmartfileReservedParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func GetSmartfileReserved(params files_sdk.IpAddressGetSmartfileReservedParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields each of this listing's `files_sdk.PublicIpAddress` values as an `any` value, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.IpAddressGetSmartfileReservedParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := ip_address.GetSmartfileReserved(files_sdk.IpAddressGetSmartfileReservedParams{})
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

### GetExavaultReserved

`GET /ip_addresses/exavault-reserved`

List all possible public ExaVault IP addresses

```go
func (c *Client) GetExavaultReserved(params files_sdk.IpAddressGetExavaultReservedParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func GetExavaultReserved(params files_sdk.IpAddressGetExavaultReservedParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields each of this listing's `files_sdk.PublicIpAddress` values as an `any` value, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.IpAddressGetExavaultReservedParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := ip_address.GetExavaultReserved(files_sdk.IpAddressGetExavaultReservedParams{})
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

### GetReserved

`GET /ip_addresses/reserved`

List all possible public IP addresses

```go
func (c *Client) GetReserved(params files_sdk.IpAddressGetReservedParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func GetReserved(params files_sdk.IpAddressGetReservedParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields each of this listing's `files_sdk.PublicIpAddress` values as an `any` value, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.IpAddressGetReservedParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := ip_address.GetReserved(files_sdk.IpAddressGetReservedParams{})
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

`files_sdk.IpAddress` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `string` | `id` | Unique label for list; used by Zapier and other integrations. |
| `AssociatedWith` | `string` | `associated_with` | The object that this public IP address list is associated with. |
| `GroupId` | `int64` | `group_id` | Group ID |
| `IpAddresses` | `[]string` | `ip_addresses` | A list of IP addresses. |
