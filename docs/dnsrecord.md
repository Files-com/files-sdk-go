# DnsRecord

`files_sdk.DnsRecord` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/dnsrecord`, named `dns_record`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	dns_record "github.com/Files-com/files-sdk-go/v3/dnsrecord"
)
```

## Operations

### List

`GET /dns_records`

Show Site DNS Configuration

```go
func (c *Client) List(params files_sdk.DnsRecordListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.DnsRecordListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.DnsRecord` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.DnsRecordListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := dns_record.List(files_sdk.DnsRecordListParams{})
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

`files_sdk.DnsRecord` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `string` | `id` | Unique label for DNS record; used by Zapier and other integrations. |
| `Domain` | `string` | `domain` | DNS record domain name |
| `Rrtype` | `string` | `rrtype` | DNS record type |
| `Value` | `string` | `value` | DNS record value |
