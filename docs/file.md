# File

`files_sdk.File` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/file`, named `file`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file"
)
```

## Operations

### Download

The API schema lists this operation as `GET /files/{path}`. The package provides these
functions and methods for it:

```go
func (c *Client) Get(Path string, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error)

func Get(Path string, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error)

// FileStats checks for current file size and modification time.
// File{}.Size and File{}.Mtime are not always up to date. This calls HEAD on File{}.DownloadUri to get the latest info.
// Some Download URLs won't support HEAD. In this case the size is reported as UntrustedSizeValue. The size can be known post download
// using Client{}.DownloadRequestStatus. This applies to the remote mount types FTP, SFTP, and WebDAV.
func (c *Client) FileStats(file files_sdk.File, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error)

func (c *Client) DownloadRequestStatus(fileDownloadUrl string, downloadRequestId string, opts ...files_sdk.RequestResponseOption) (files_sdk.ResponseError, error)

func (c *Client) DownloadUri(params files_sdk.FileDownloadParams, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error)

func (c *Client) Download(params files_sdk.FileDownloadParams, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error)

func Download(params files_sdk.FileDownloadParams, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error)
```

Parameters, `files_sdk.FileDownloadParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Path` | `string` | `path`, in the path | yes | Path to operate on. |
| `Action` | `string` | `action` | no | Can be blank, `redirect` or `stat`. If set to `stat`, we will return file information but without a download URL, and without logging a download. If set to `redirect` we will serve a 302 redirect directly to the file. This is used for integrations with Zapier, and is not recommended for most integrations. |
| `PreviewSize` | `string` | `preview_size` | no | Request a preview size. Can be `small` (default), `large`, `xlarge`, or `pdf`. |
| `WithPreviews` | `*bool` | `with_previews` | no | Include file preview information? |
| `WithPriorityColor` | `*bool` | `with_priority_color` | no | Include file priority color information? |
| `WithDirectConnectionInfo` | `*bool` | `with_direct_connection_info` | no | Include optional direct connection information for a direct Agent transfer attempt? |
| `File` | `files_sdk.File` | not sent | no | Its tags send nothing to the API. |

Example:

```go
result, err := file.Download(files_sdk.FileDownloadParams{
	Path: "path",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Create

`POST /files/{path}`

Upload File

```go
func (c *Client) Create(params files_sdk.FileCreateParams, opts ...files_sdk.RequestResponseOption) (file files_sdk.File, err error)
func Create(params files_sdk.FileCreateParams, opts ...files_sdk.RequestResponseOption) (file files_sdk.File, err error)
```

Parameters, `files_sdk.FileCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Path` | `string` | `path`, in the path | yes | Path to operate on. |
| `Action` | `string` | `action` | no | The action to perform. Can be `append`, `attachment`, `end`, `upload`, `put`, or may not exist |
| `CustomMetadata` | `interface{}` | `custom_metadata` | no | Custom metadata map to save when `action=end` completes the upload. Replaces existing metadata; an empty map clears it. No separate metadata-edit permission is required. Supported on native files and configured remote mounts, excluding remote server automount paths. Limited to 32 keys, 256 characters per key and 1024 characters per value. |
| `EtagsParam` | `[]files_sdk.EtagsParam` | `etags` | no | Nested parameters: see `files_sdk.EtagsParam`. |
| `Length` | `int64` | `length` | no | Length of file. |
| `MkdirParents` | `*bool` | `mkdir_parents` | no | Create parent directories if they do not exist? |
| `Part` | `int64` | `part` | no | Part if uploading a part. |
| `Parts` | `int64` | `parts` | no | How many parts to fetch? |
| `ProvidedMtime` | `*time.Time` | `provided_mtime` | no | User provided modification time. |
| `Ref` | `string` | `ref` | no |  |
| `Restart` | `int64` | `restart` | no | File byte offset to restart from. |
| `Size` | `int64` | `size` | no | Size of file. |
| `CopyBehaviors` | `*bool` | `copy_behaviors` | no | If copying a folder, also copy supported behaviors, email notification subscriptions, and per-folder branding to the destination folder tree? |
| `Structure` | `string` | `structure` | no | If copying folder, copy just the structure? |
| `WithRename` | `*bool` | `with_rename` | no | Allow file rename instead of overwrite? |
| `BufferedUpload` | `*bool` | `buffered_upload` | no | If true, and the path refers to a destination not stored on Files.com (such as a remote server mount), the upload will be uploaded first to Files.com before being sent to the remote server mount. This can allow clients to upload using parallel parts to a remote server destination that does not offer parallel parts support natively. |
| `WithDirectConnectionInfo` | `*bool` | `with_direct_connection_info` | no | Include optional direct connection information for a direct Agent transfer attempt? |
| `ActionAttributes` | `map[string]any` | `action_attributes` | no |  |

Example:

```go
result, err := file.Create(files_sdk.FileCreateParams{
	Path: "path",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Update

`PATCH /files/{path}`

Update File/Folder Metadata

```go
func (c *Client) Update(params files_sdk.FileUpdateParams, opts ...files_sdk.RequestResponseOption) (file files_sdk.File, err error)
func Update(params files_sdk.FileUpdateParams, opts ...files_sdk.RequestResponseOption) (file files_sdk.File, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (file files_sdk.File, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (file files_sdk.File, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.FileUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Path` | `string` | `path`, in the path | yes | Path to operate on. |
| `CustomMetadata` | `interface{}` | `custom_metadata` | no | Custom metadata map of keys and values. Limited to 32 keys, 256 characters per key and 1024 characters per value. |
| `ProvidedMtime` | `*time.Time` | `provided_mtime` | no | Modified time of file. |
| `PriorityColor` | `string` | `priority_color` | no | Priority/Bookmark color of file. |

Example:

```go
result, err := file.Update(files_sdk.FileUpdateParams{
	Path: "path",
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := file.UpdateWithMap(map[string]interface{}{
	"path": "path",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /files/{path}`

Delete File/Folder

```go
func (c *Client) Delete(params files_sdk.FileDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.FileDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.FileDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Path` | `string` | `path`, in the path | yes | Path to operate on. |
| `Recursive` | `*bool` | `recursive` | no | If true, will recursively delete folders. Otherwise, will error on non-empty folders. |

Example:

```go
err := file.Delete(files_sdk.FileDeleteParams{
	Path: "path",
})
if err != nil {
	return err
}
```

### Find

`GET /file_actions/metadata/{path}`

Find File/Folder by Path

```go
func (c *Client) Find(params files_sdk.FileFindParams, opts ...files_sdk.RequestResponseOption) (file files_sdk.File, err error)
func Find(params files_sdk.FileFindParams, opts ...files_sdk.RequestResponseOption) (file files_sdk.File, err error)
```

Parameters, `files_sdk.FileFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Path` | `string` | `path`, in the path | yes | Path to operate on. |
| `PreviewSize` | `string` | `preview_size` | no | Request a preview size. Can be `small` (default), `large`, `xlarge`, or `pdf`. |
| `WithPreviews` | `*bool` | `with_previews` | no | Include file preview information? |
| `WithPriorityColor` | `*bool` | `with_priority_color` | no | Include file priority color information? |

Example:

```go
result, err := file.Find(files_sdk.FileFindParams{
	Path: "path",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### ZipListContents

`GET /file_actions/zip_list/{path}`

List the contents of a ZIP file

```go
func (c *Client) ZipListContents(params files_sdk.FileZipListContentsParams, opts ...files_sdk.RequestResponseOption) (zipListEntryCollection files_sdk.ZipListEntryCollection, err error)
func ZipListContents(params files_sdk.FileZipListContentsParams, opts ...files_sdk.RequestResponseOption) (zipListEntryCollection files_sdk.ZipListEntryCollection, err error)
```

Parameters, `files_sdk.FileZipListContentsParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Path` | `string` | `path`, in the path | yes | Path to operate on. |

Example:

```go
result, err := file.ZipListContents(files_sdk.FileZipListContentsParams{
	Path: "path",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### BeginDelete

`POST /file_actions/delete/{path}`

Delete File/Folder and return its File Action. A pending or processing deletion continues in the background after this returns; check its File Migration for completion. Waiting for the File Migration does not cancel the deletion

```go
func (c *Client) BeginDelete(params files_sdk.FileBeginDeleteParams, opts ...files_sdk.RequestResponseOption) (fileAction files_sdk.FileAction, err error)
func BeginDelete(params files_sdk.FileBeginDeleteParams, opts ...files_sdk.RequestResponseOption) (fileAction files_sdk.FileAction, err error)
```

Parameters, `files_sdk.FileBeginDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Path` | `string` | `path`, in the path | yes | Path to operate on. |
| `Recursive` | `*bool` | `recursive` | no | If true, will recursively delete folders. Otherwise, will error on non-empty folders. |

Example:

```go
result, err := file.BeginDelete(files_sdk.FileBeginDeleteParams{
	Path: "path",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Copy

`POST /file_actions/copy/{path}`

Copy File/Folder

```go
func (c *Client) Copy(params files_sdk.FileCopyParams, opts ...files_sdk.RequestResponseOption) (fileAction files_sdk.FileAction, err error)
func Copy(params files_sdk.FileCopyParams, opts ...files_sdk.RequestResponseOption) (fileAction files_sdk.FileAction, err error)
```

Parameters, `files_sdk.FileCopyParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Path` | `string` | `path`, in the path | yes | Path to operate on. |
| `Destination` | `string` | `destination` | yes | Copy destination path. |
| `CopyBehaviors` | `*bool` | `copy_behaviors` | no | If copying a folder, also copy supported behaviors, email notification subscriptions, and per-folder branding to the destination folder tree? |
| `Structure` | `*bool` | `structure` | no | Copy structure only? |
| `Overwrite` | `*bool` | `overwrite` | no | Overwrite existing file(s) in the destination? |

Example:

```go
result, err := file.Copy(files_sdk.FileCopyParams{
	Path:        "path",
	Destination: "destination",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Move

`POST /file_actions/move/{path}`

Move File/Folder

```go
func (c *Client) Move(params files_sdk.FileMoveParams, opts ...files_sdk.RequestResponseOption) (fileAction files_sdk.FileAction, err error)
func Move(params files_sdk.FileMoveParams, opts ...files_sdk.RequestResponseOption) (fileAction files_sdk.FileAction, err error)
```

Parameters, `files_sdk.FileMoveParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Path` | `string` | `path`, in the path | yes | Path to operate on. |
| `Destination` | `string` | `destination` | yes | Move destination path. |
| `Overwrite` | `*bool` | `overwrite` | no | Overwrite existing file(s) in the destination? |

Example:

```go
result, err := file.Move(files_sdk.FileMoveParams{
	Path:        "path",
	Destination: "destination",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Transform

`POST /file_actions/transform/{path}`

Transform a file and save the output to a destination path

```go
func (c *Client) Transform(params files_sdk.FileTransformParams, opts ...files_sdk.RequestResponseOption) (fileAction files_sdk.FileAction, err error)
func Transform(params files_sdk.FileTransformParams, opts ...files_sdk.RequestResponseOption) (fileAction files_sdk.FileAction, err error)
```

Parameters, `files_sdk.FileTransformParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Path` | `string` | `path`, in the path | yes | Path to operate on. |
| `Destination` | `string` | `destination` | yes | Destination file path for the transformed output. |
| `TransformType` | `string` | `transform_type` | yes | Transform type. Supported values are `image_convert`, `document_convert`, and `files_transform_script_execute`. |
| `TargetFormat` | `string` | `target_format` | yes | Destination format to create. |
| `Script` | `string` | `script` | no | Files TransformScript source. Required when transform_type is `files_transform_script_execute`. |
| `Width` | `int64` | `width` | no | Maximum output width for image_convert. |
| `Height` | `int64` | `height` | no | Maximum output height for image_convert. |
| `Overwrite` | `*bool` | `overwrite` | no | Overwrite existing file in the destination? |

Example:

```go
result, err := file.Transform(files_sdk.FileTransformParams{
	Path:          "path",
	Destination:   "destination",
	TransformType: "transform_type",
	TargetFormat:  "target_format",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### GpgDecrypt

`POST /file_actions/gpg_decrypt/{path}`

Decrypt a GPG-encrypted file and save it to a destination path

```go
func (c *Client) GpgDecrypt(params files_sdk.FileGpgDecryptParams, opts ...files_sdk.RequestResponseOption) (fileAction files_sdk.FileAction, err error)
func GpgDecrypt(params files_sdk.FileGpgDecryptParams, opts ...files_sdk.RequestResponseOption) (fileAction files_sdk.FileAction, err error)
```

Parameters, `files_sdk.FileGpgDecryptParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Path` | `string` | `path`, in the path | yes | Path to operate on. |
| `Destination` | `string` | `destination` | yes | Destination file path for the decrypted file. |
| `GpgKeyIds` | `[]int64` | `gpg_key_ids` | no | GPG Key IDs to decrypt with. If omitted, every accessible private GPG key in the source workspace is used. |
| `GpgKeyPartnerId` | `int64` | `gpg_key_partner_id` | no | Partner ID whose GPG keys should be used for decryption. |
| `UseAllPrivateKeys` | `*bool` | `use_all_private_keys` | no | Use every accessible private GPG key in the source workspace for decryption. |
| `IgnoreMdcError` | `*bool` | `ignore_mdc_error` | no | Ignore errors from the MDC (modification detection code) check. |
| `Overwrite` | `*bool` | `overwrite` | no | Overwrite existing file in the destination? |

Example:

```go
result, err := file.GpgDecrypt(files_sdk.FileGpgDecryptParams{
	Path:        "path",
	Destination: "destination",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### GpgEncrypt

`POST /file_actions/gpg_encrypt/{path}`

Encrypt a file with GPG and save it to a destination path

```go
func (c *Client) GpgEncrypt(params files_sdk.FileGpgEncryptParams, opts ...files_sdk.RequestResponseOption) (fileAction files_sdk.FileAction, err error)
func GpgEncrypt(params files_sdk.FileGpgEncryptParams, opts ...files_sdk.RequestResponseOption) (fileAction files_sdk.FileAction, err error)
```

Parameters, `files_sdk.FileGpgEncryptParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Path` | `string` | `path`, in the path | yes | Path to operate on. |
| `Destination` | `string` | `destination` | yes | Destination file path for the encrypted file. |
| `GpgKeyIds` | `[]int64` | `gpg_key_ids` | no | GPG Key IDs to encrypt with. |
| `GpgKeyPartnerId` | `int64` | `gpg_key_partner_id` | no | Partner ID whose GPG keys should be used for encryption. |
| `SigningKeyId` | `int64` | `signing_key_id` | no | Optional GPG Key ID to sign with. |
| `Armor` | `*bool` | `armor` | no | Output ASCII-armored encrypted data. |
| `Overwrite` | `*bool` | `overwrite` | no | Overwrite existing file in the destination? |

Example:

```go
result, err := file.GpgEncrypt(files_sdk.FileGpgEncryptParams{
	Path:        "path",
	Destination: "destination",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Unzip

`POST /file_actions/unzip`

Extract a ZIP file to a destination folder

```go
func (c *Client) Unzip(params files_sdk.FileUnzipParams, opts ...files_sdk.RequestResponseOption) (fileAction files_sdk.FileAction, err error)
func Unzip(params files_sdk.FileUnzipParams, opts ...files_sdk.RequestResponseOption) (fileAction files_sdk.FileAction, err error)
```

Parameters, `files_sdk.FileUnzipParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Path` | `string` | `path` | yes | ZIP file path to extract. |
| `Destination` | `string` | `destination` | yes | Destination folder path for extracted files. |
| `Filename` | `string` | `filename` | no | Optional single entry filename to extract. |
| `Overwrite` | `*bool` | `overwrite` | no | Overwrite existing files in the destination? |

Example:

```go
result, err := file.Unzip(files_sdk.FileUnzipParams{
	Path:        "path",
	Destination: "destination",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Zip

`POST /file_actions/zip`

Create a ZIP from one or more paths and save it to a destination path

```go
func (c *Client) Zip(params files_sdk.FileZipParams, opts ...files_sdk.RequestResponseOption) (fileAction files_sdk.FileAction, err error)
func Zip(params files_sdk.FileZipParams, opts ...files_sdk.RequestResponseOption) (fileAction files_sdk.FileAction, err error)
```

Parameters, `files_sdk.FileZipParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Paths` | `[]string` | `paths` | yes | Paths to include in the ZIP. |
| `Destination` | `string` | `destination` | yes | Destination file path for the ZIP. |
| `Overwrite` | `*bool` | `overwrite` | no | Overwrite existing file in the destination? |

Example:

```go
// Set this required value before the call.
var paths []string
result, err := file.Zip(files_sdk.FileZipParams{
	Paths:       paths,
	Destination: "destination",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### BeginUpload

`POST /file_actions/begin_upload/{path}`

Begin File Upload

```go
func (c *Client) BeginUpload(params files_sdk.FileBeginUploadParams, opts ...files_sdk.RequestResponseOption) (fileUploadPartCollection files_sdk.FileUploadPartCollection, err error)
func BeginUpload(params files_sdk.FileBeginUploadParams, opts ...files_sdk.RequestResponseOption) (fileUploadPartCollection files_sdk.FileUploadPartCollection, err error)
```

Parameters, `files_sdk.FileBeginUploadParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Path` | `string` | `path`, in the path | yes | Path to operate on. |
| `MkdirParents` | `*bool` | `mkdir_parents` | no | Create parent directories if they do not exist? |
| `Part` | `int64` | `part` | no | Part if uploading a part. |
| `Parts` | `int64` | `parts` | no | How many parts to fetch? |
| `Ref` | `string` | `ref` | no |  |
| `Restart` | `int64` | `restart` | no | File byte offset to restart from. |
| `Size` | `int64` | `size` | no | Total bytes of file being uploaded (include bytes being retained if appending/restarting). |
| `WithRename` | `*bool` | `with_rename` | no | Allow file rename instead of overwrite? |
| `BufferedUpload` | `*bool` | `buffered_upload` | no | If true, and the path refers to a destination not stored on Files.com (such as a remote server mount), the upload will be uploaded first to Files.com before being sent to the remote server mount. This can allow clients to upload using parallel parts to a remote server destination that does not offer parallel parts support natively. |
| `WithDirectConnectionInfo` | `*bool` | `with_direct_connection_info` | no | Include optional direct connection information for a direct Agent transfer attempt? |

Example:

```go
result, err := file.BeginUpload(files_sdk.FileBeginUploadParams{
	Path: "path",
})
if err != nil {
	return err
}
fmt.Println(result)
```

## Other functions and methods

The package also provides these functions and methods:

```go
func (i Iter) LoadResource(identifier interface{}, opts ...files_sdk.RequestResponseOption) (interface{}, error)

func (i Iter) Iterate(identifier interface{}, opts ...files_sdk.RequestResponseOption) (files_sdk.IterI, error)

func (c *Client) ListFor(params files_sdk.FolderListForParams, opts ...files_sdk.RequestResponseOption) (Iter, error)

func (c *Client) CreateFolder(params files_sdk.FolderCreateParams, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error)

func (r RecursiveItem) Err() error

func (c *Client) ListForRecursive(params files_sdk.FolderListForParams, opts ...files_sdk.RequestResponseOption) (files_sdk.TypedIterI[RecursiveItem], error)
```

## Model

`files_sdk.File` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Path` | `string` | `path` | File/Folder path. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters. |
| `CreatedById` | `int64` | `created_by_id` | User ID of the User who created the file/folder |
| `CreatedByApiKeyId` | `int64` | `created_by_api_key_id` | ID of the API key that created the file/folder |
| `CreatedByAs2IncomingMessageId` | `int64` | `created_by_as2_incoming_message_id` | ID of the AS2 Incoming Message that created the file/folder |
| `CreatedByAutomationId` | `int64` | `created_by_automation_id` | ID of the Automation that created the file/folder |
| `CreatedByBundleRegistrationId` | `int64` | `created_by_bundle_registration_id` | ID of the Bundle Registration that created the file/folder |
| `CreatedByInboxId` | `int64` | `created_by_inbox_id` | ID of the Inbox that created the file/folder |
| `CreatedByRemoteServerId` | `int64` | `created_by_remote_server_id` | ID of the Remote Server that created the file/folder |
| `CreatedBySyncId` | `int64` | `created_by_sync_id` | ID of the Sync that created the file/folder |
| `CustomMetadata` | `interface{}` | `custom_metadata` | Custom metadata map of keys and values. Limited to 32 keys, 256 characters per key and 1024 characters per value. |
| `DisplayName` | `string` | `display_name` | File/Folder display name |
| `Type` | `string` | `type` | Type: `directory` or `file`. |
| `Size` | `int64` | `size` | File/Folder size |
| `CreatedAt` | `*time.Time` | `created_at` | File created date/time |
| `LastModifiedById` | `int64` | `last_modified_by_id` | User ID of the User who last modified the file/folder |
| `LastModifiedByApiKeyId` | `int64` | `last_modified_by_api_key_id` | ID of the API key that last modified the file/folder |
| `LastModifiedByAutomationId` | `int64` | `last_modified_by_automation_id` | ID of the Automation that last modified the file/folder |
| `LastModifiedByBundleRegistrationId` | `int64` | `last_modified_by_bundle_registration_id` | ID of the Bundle Registration that last modified the file/folder |
| `LastModifiedByRemoteServerId` | `int64` | `last_modified_by_remote_server_id` | ID of the Remote Server that last modified the file/folder |
| `LastModifiedBySyncId` | `int64` | `last_modified_by_sync_id` | ID of the Sync that last modified the file/folder |
| `Mtime` | `*time.Time` | `mtime` | File last modified date/time, according to the server. This is the timestamp of the last Files.com operation of the file, regardless of what modified timestamp was sent. |
| `ProvidedMtime` | `*time.Time` | `provided_mtime` | File last modified date/time, according to the client who set it. Files.com allows desktop, FTP, SFTP, and WebDAV clients to set modified at times. This allows Desktop<->Cloud syncing to preserve modified at times. |
| `Crc32` | `string` | `crc32` | File CRC32 checksum. This is sometimes delayed, so if you get a blank response, wait and try again. |
| `Md5` | `string` | `md5` | File MD5 checksum. This is sometimes delayed, so if you get a blank response, wait and try again. |
| `Sha1` | `string` | `sha1` | File SHA1 checksum. This is sometimes delayed, so if you get a blank response, wait and try again. |
| `Sha256` | `string` | `sha256` | File SHA256 checksum. This is sometimes delayed, so if you get a blank response, wait and try again. |
| `MimeType` | `string` | `mime_type` | MIME Type. This is determined by the filename extension and is not stored separately internally. |
| `Region` | `string` | `region` | Region location |
| `Permissions` | `string` | `permissions` | A short string representing the current user's permissions. Can be `r` (Read),`w` (Write),`d` (Delete), `l` (List) or any combination |
| `SubfoldersLocked` | `*bool` | `subfolders_locked?` | Are subfolders locked and unable to be modified? |
| `IsLocked` | `*bool` | `is_locked` | Is this folder locked and unable to be modified? |
| `DownloadUri` | `string` | `download_uri` | Link to download file. Provided only in response to a download request. |
| `DirectConnectionInfo` | `files_sdk.DirectConnectionInfo` | `direct_connection_info` | Optional direct connection information for direct Agent transfer attempts |
| `PriorityColor` | `string` | `priority_color` | Bookmark/priority color of file/folder |
| `PreviewId` | `int64` | `preview_id` | File preview ID |
| `Preview` | `files_sdk.Preview` | `preview` | File preview |
| `Action` | `string` | `action` | The action to perform. Can be `append`, `attachment`, `end`, `upload`, `put`, or may not exist |
| `Length` | `int64` | `length` | Length of file. |
| `MkdirParents` | `*bool` | `mkdir_parents` | Create parent directories if they do not exist? |
| `Part` | `int64` | `part` | Part if uploading a part. |
| `Parts` | `int64` | `parts` | How many parts to fetch? |
| `Ref` | `string` | `ref` |  |
| `Restart` | `int64` | `restart` | File byte offset to restart from. |
| `CopyBehaviors` | `*bool` | `copy_behaviors` | If copying a folder, also copy supported behaviors, email notification subscriptions, and per-folder branding to the destination folder tree? |
| `Structure` | `string` | `structure` | If copying folder, copy just the structure? |
| `WithRename` | `*bool` | `with_rename` | Allow file rename instead of overwrite? |
| `BufferedUpload` | `*bool` | `buffered_upload` | If true, and the path refers to a destination not stored on Files.com (such as a remote server mount), the upload will be uploaded first to Files.com before being sent to the remote server mount. This can allow clients to upload using parallel parts to a remote server destination that does not offer parallel parts support natively. |
| `WithDirectConnectionInfo` | `*bool` | `with_direct_connection_info` | Include optional direct connection information for a direct Agent transfer attempt? |
