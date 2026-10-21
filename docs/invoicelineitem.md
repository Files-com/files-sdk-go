# InvoiceLineItem

`files_sdk.InvoiceLineItem` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/invoicelineitem`, named `invoice_line_item`.

This package has no API operations of its own.

## Model

`files_sdk.InvoiceLineItem` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Invoice Line item Id |
| `Amount` | `string` | `amount` | Invoice line item amount |
| `CreatedAt` | `*time.Time` | `created_at` | Invoice line item created at date/time |
| `Description` | `string` | `description` | Invoice line item description |
| `Type` | `string` | `type` | Invoice line item type |
| `ServiceEndAt` | `*time.Time` | `service_end_at` | Invoice line item service end date/time |
| `ServiceStartAt` | `*time.Time` | `service_start_at` | Invoice line item service start date/time |
| `Plan` | `string` | `plan` | Plan name |
| `Site` | `string` | `site` | Site name |
| `PrepaidBytes` | `int64` | `prepaid_bytes` | Prepaid bytes purchased for this invoice line item |
| `PrepaidBytesExpireAt` | `*time.Time` | `prepaid_bytes_expire_at` | When the prepaid bytes expire |
| `PrepaidBytesUsed` | `int64` | `prepaid_bytes_used` | Total prepaid bytes used for this invoice line item |
| `PrepaidBytesAvailable` | `int64` | `prepaid_bytes_available` | Available prepaid bytes for this invoice line item |
