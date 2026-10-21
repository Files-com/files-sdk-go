# Preview

`files_sdk.Preview` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/preview`, named `preview`.

This package has no API operations of its own.

## Model

`files_sdk.Preview` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Preview ID |
| `Status` | `string` | `status` | Preview status. Can be invalid, not_generated, generating, complete, or file_too_large |
| `DownloadUri` | `string` | `download_uri` | Link to download preview |
| `Type` | `string` | `type` | Preview type. Can be image, pdf, pdf_native, video, audio, or text |
| `Size` | `string` | `size` | Preview size |
