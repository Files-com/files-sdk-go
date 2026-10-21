# Status

`files_sdk.Status` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/status`, named `status`.

This package has no API operations of its own.

## Model

`files_sdk.Status` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Code` | `int64` | `code` | Status HTTP code |
| `Message` | `string` | `message` | Error message |
| `Status` | `string` | `status` | Status message |
| `Data` | `files_sdk.Auto` | `data` | Additional data |
| `Errors` | `[]files_sdk.Errors` | `errors` | A list of api errors |
| `ClickwrapId` | `int64` | `clickwrap_id` | Required Clickwrap id |
| `ClickwrapBody` | `string` | `clickwrap_body` | Required Clickwrap body |
