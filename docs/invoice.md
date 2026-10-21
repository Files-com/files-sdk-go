# Invoice

`files_sdk.Invoice` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/invoice`, named `invoice`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/invoice"
)
```

## Operations

### List

`GET /invoices`

List Invoices

```go
func (c *Client) List(params files_sdk.InvoiceListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.InvoiceListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.AccountLineItem` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.InvoiceListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := invoice.List(files_sdk.InvoiceListParams{})
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

`GET /invoices/{id}`

Show Invoice

```go
func (c *Client) Find(params files_sdk.InvoiceFindParams, opts ...files_sdk.RequestResponseOption) (accountLineItem files_sdk.AccountLineItem, err error)
func Find(params files_sdk.InvoiceFindParams, opts ...files_sdk.RequestResponseOption) (accountLineItem files_sdk.AccountLineItem, err error)
```

Parameters, `files_sdk.InvoiceFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Invoice ID. |

Example:

```go
result, err := invoice.Find(files_sdk.InvoiceFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

## Model

`files_sdk.Invoice` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Line item Id |
| `Amount` | `string` | `amount` | Line item amount |
| `Balance` | `string` | `balance` | Line item balance |
| `CreatedAt` | `*time.Time` | `created_at` | Line item created at |
| `Currency` | `string` | `currency` | Line item currency |
| `DownloadUri` | `string` | `download_uri` | Line item download uri |
| `InvoiceLineItems` | `[]files_sdk.InvoiceLineItem` | `invoice_line_items` | Associated invoice line items |
| `Method` | `string` | `method` | Line item payment method |
| `PaymentLineItems` | `[]files_sdk.PaymentLineItem` | `payment_line_items` | Associated payment line items |
| `PaymentReversedAt` | `*time.Time` | `payment_reversed_at` | Date/time payment was reversed if applicable |
| `PaymentType` | `string` | `payment_type` | Type of payment if applicable |
| `SiteName` | `string` | `site_name` | Site name this line item is for |
| `Type` | `string` | `type` | Type of line item, either payment or invoice |
