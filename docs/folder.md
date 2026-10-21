# Folder

`files_sdk.Folder` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/folder`, named `folder`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/folder"
)
```

## Operations

### ListFor

`GET /folders/{path}`

List Folders by Path

```go
func (c *Client) ListFor(params files_sdk.FolderListForParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func ListFor(params files_sdk.FolderListForParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.File` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.FolderListForParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Path` | `string` | `path`, in the path | yes | Path to operate on. |
| `PreviewSize` | `string` | `preview_size` | no | Request a preview size. Can be `small` (default), `large`, `xlarge`, or `pdf`. |
| `SortBy` | `interface{}` | `sort_by` | no | Search by field and direction. Valid fields are `path`, `size`, `modified_at_datetime`, `provided_modified_at`. Valid directions are `asc` and `desc`. Defaults to `{"path":"asc"}`. |
| `Search` | `string` | `search` | no | If specified, will search the folders/files list by name. Ignores text before last `/`. This is the same API used by the search bar in the web UI when running 'Search This Folder'. Search results are a best effort, not real time, and not guaranteed to perfectly match the latest folder listing. Results may be truncated if more than 1,000 possible matches exist. This field should only be used for ad-hoc (human) searching, and not as part of an automated process. |
| `SearchCustomMetadataKey` | `string` | `search_custom_metadata_key` | no | If provided, the search string in `search` will search for files where this custom metadata key matches the value sent in `search`. Set this to `*` to allow any metadata key to match the value sent in `search`. |
| `SearchAll` | `*bool` | `search_all` | no | Search entire site? If true, we will search the entire site. Do not provide a path when using this parameter. This is the same API used by the search bar in the web UI when running 'Search All Files'. Search results are a best effort, not real time, and not guaranteed to match every file. This field should only be used for ad-hoc (human) searching, and not as part of an automated process. |
| `WithPreviews` | `*bool` | `with_previews` | no | Include file previews? |
| `WithPriorityColor` | `*bool` | `with_priority_color` | no | Include file priority color information? |
| `Type` | `string` | `type` | no | Type of objects to return. Can be `folder` or `file`. |
| `ModifiedAtDatetime` | `*time.Time` | `modified_at_datetime` | no | If provided, will only return files/folders modified after this time. Can be used only in combination with `type` filter. |
| `ConcurrencyManager` | `lib.ConcurrencyManagerWithSubWorker` | not sent | no | Its tags send nothing to the API. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := folder.ListFor(files_sdk.FolderListForParams{
	Path: "path",
})
if err != nil {
	return err
}
for item, err := range it.All() {
	if err != nil {
		return err
	}
	fmt.Println(item)
}
```

### Create

`POST /folders/{path}`

Create Folder

```go
func (c *Client) Create(params files_sdk.FolderCreateParams, opts ...files_sdk.RequestResponseOption) (file files_sdk.File, err error)
func Create(params files_sdk.FolderCreateParams, opts ...files_sdk.RequestResponseOption) (file files_sdk.File, err error)
```

Parameters, `files_sdk.FolderCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Path` | `string` | `path`, in the path | yes | Path to operate on. |
| `MkdirParents` | `*bool` | `mkdir_parents` | no | Create parent directories if they do not exist? |
| `ProvidedMtime` | `*time.Time` | `provided_mtime` | no | User provided modification time. |

Example:

```go
result, err := folder.Create(files_sdk.FolderCreateParams{
	Path: "path",
})
if err != nil {
	return err
}
fmt.Println(result)
```

## Model

`files_sdk.Folder` has these fields from the API schema:

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
| `MkdirParents` | `*bool` | `mkdir_parents` | Create parent directories if they do not exist? |
