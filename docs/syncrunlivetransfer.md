# SyncRunLiveTransfer

`files_sdk.SyncRunLiveTransfer` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/syncrunlivetransfer`, named `sync_run_live_transfer`.

This package has no API operations of its own.

## Model

`files_sdk.SyncRunLiveTransfer` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Path` | `string` | `path` | The file path being transferred. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters. |
| `Status` | `string` | `status` | Status of this individual transfer |
| `BytesCopied` | `int64` | `bytes_copied` | Bytes transferred so far |
| `BytesTotal` | `int64` | `bytes_total` | Total bytes of the file being transferred |
| `Percentage` | `float64` | `percentage` | Transfer progress from 0.0 to 1.0 |
| `Eta` | `string` | `eta` | Estimated time remaining (human-readable) |
| `StartedAt` | `string` | `started_at` | When this individual transfer started |
