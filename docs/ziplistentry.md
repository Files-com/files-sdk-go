# ZipListEntry

`files_sdk.ZipListEntry` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/ziplistentry`, named `zip_list_entry`.

This package has no API operations of its own.

## Model

`files_sdk.ZipListEntry` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Path` | `string` | `path` | Entry path inside the ZIP. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters. |
| `Size` | `int64` | `size` | Uncompressed size in bytes. |
