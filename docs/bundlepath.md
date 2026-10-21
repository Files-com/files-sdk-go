# BundlePath

`files_sdk.BundlePath` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/bundlepath`, named `bundle_path`.

This package has no API operations of its own.

## Model

`files_sdk.BundlePath` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Recursive` | `*bool` | `recursive` | Allow access to subfolders content? |
| `Path` | `string` | `path` | The path to the resource relative to filesystem. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters. |
