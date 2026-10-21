# UsageByTopLevelDir

`files_sdk.UsageByTopLevelDir` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/usagebytopleveldir`, named `usage_by_top_level_dir`.

This package has no API operations of its own.

## Model

`files_sdk.UsageByTopLevelDir` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Dir` | `string` | `dir` | Directory name |
| `Size` | `int64` | `size` | Usage |
| `Count` | `int64` | `count` | File count |
