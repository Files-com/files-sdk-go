# Errors

`files_sdk.Errors` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/errors`, named `errors`.

This package has no API operations of its own.

## Model

`files_sdk.Errors` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Fields` | `[]string` | `fields` | A list of fields where errors occur |
| `Messages` | `[]string` | `messages` | A list of error messages |
