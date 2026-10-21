# AccountLineItem

`files_sdk.AccountLineItem` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/accountlineitem`, named `account_line_item`.

This package has no API operations of its own.

## Model

`files_sdk.AccountLineItem` has these fields from the API schema:

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
