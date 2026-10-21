# FileUploadPart

`files_sdk.FileUploadPart` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/fileuploadpart`, named `file_upload_part`.

This package has no API operations of its own.

## Model

`files_sdk.FileUploadPart` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Send` | `interface{}` | `send` | Content-Type and File to send |
| `Action` | `string` | `action` | Type of upload |
| `AskAboutOverwrites` | `*bool` | `ask_about_overwrites` | If `true`, this file exists and you may wish to ask the user for overwrite confirmation |
| `AvailableParts` | `int64` | `available_parts` | Number of parts in the upload |
| `Expires` | `string` | `expires` | Date/time of when this Upload part expires and the URL cannot be used any more |
| `Headers` | `interface{}` | `headers` | Additional upload headers to provide as part of the upload |
| `HttpMethod` | `string` | `http_method` | HTTP Method to use for uploading the part, usually `PUT` |
| `NextPartsize` | `int64` | `next_partsize` | Size in bytes for this part |
| `ParallelParts` | `*bool` | `parallel_parts` | If `true`, multiple parts may be uploaded in parallel. If `false`, be sure to only upload one part at a time, in order. |
| `RetryParts` | `*bool` | `retry_parts` | If `true`, parts may be retried. If `false`, a part cannot be retried and the upload should be restarted. |
| `Parameters` | `interface{}` | `parameters` | Additional HTTP parameters to send with the upload |
| `PartNumber` | `int64` | `part_number` | Number of this upload part |
| `Partsize` | `int64` | `partsize` | Size in bytes for the next upload part |
| `Path` | `string` | `path` | New file path. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters. |
| `Ref` | `string` | `ref` | Reference name for this upload part |
| `UploadUri` | `string` | `upload_uri` | URI to upload this part to |
| `DirectConnectionInfo` | `files_sdk.DirectConnectionInfo` | `direct_connection_info` | Optional direct connection information for direct Agent transfer attempts |
