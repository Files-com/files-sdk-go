# PaymentLineItem

`files_sdk.PaymentLineItem` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/paymentlineitem`, named `payment_line_item`.

This package has no API operations of its own.

## Model

`files_sdk.PaymentLineItem` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Amount` | `string` | `amount` | Payment line item amount |
| `CreatedAt` | `*time.Time` | `created_at` | Payment line item created at date/time |
| `InvoiceId` | `int64` | `invoice_id` | Invoice ID |
| `PaymentId` | `int64` | `payment_id` | Payment ID |
