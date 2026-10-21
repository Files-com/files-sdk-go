package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// File is a Files.com API resource.
type File struct {
	// File/Folder path. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.
	Path string `json:"path,omitempty" path:"path,omitempty" url:"path,omitempty"`
	// User ID of the User who created the file/folder.
	CreatedById int64 `json:"created_by_id,omitempty" path:"created_by_id,omitempty" url:"created_by_id,omitempty"`
	// ID of the API key that created the file/folder.
	CreatedByApiKeyId int64 `json:"created_by_api_key_id,omitempty" path:"created_by_api_key_id,omitempty" url:"created_by_api_key_id,omitempty"`
	// ID of the AS2 Incoming Message that created the file/folder.
	CreatedByAs2IncomingMessageId int64 `json:"created_by_as2_incoming_message_id,omitempty" path:"created_by_as2_incoming_message_id,omitempty" url:"created_by_as2_incoming_message_id,omitempty"`
	// ID of the Automation that created the file/folder.
	CreatedByAutomationId int64 `json:"created_by_automation_id,omitempty" path:"created_by_automation_id,omitempty" url:"created_by_automation_id,omitempty"`
	// ID of the Bundle Registration that created the file/folder.
	CreatedByBundleRegistrationId int64 `json:"created_by_bundle_registration_id,omitempty" path:"created_by_bundle_registration_id,omitempty" url:"created_by_bundle_registration_id,omitempty"`
	// ID of the Inbox that created the file/folder.
	CreatedByInboxId int64 `json:"created_by_inbox_id,omitempty" path:"created_by_inbox_id,omitempty" url:"created_by_inbox_id,omitempty"`
	// ID of the Remote Server that created the file/folder.
	CreatedByRemoteServerId int64 `json:"created_by_remote_server_id,omitempty" path:"created_by_remote_server_id,omitempty" url:"created_by_remote_server_id,omitempty"`
	// ID of the Sync that created the file/folder.
	CreatedBySyncId int64 `json:"created_by_sync_id,omitempty" path:"created_by_sync_id,omitempty" url:"created_by_sync_id,omitempty"`
	// Custom metadata map of keys and values. Limited to 32 keys, 256 characters per key and 1024 characters per value.
	CustomMetadata interface{} `json:"custom_metadata,omitempty" path:"custom_metadata,omitempty" url:"custom_metadata,omitempty"`
	// File/Folder display name.
	DisplayName string `json:"display_name,omitempty" path:"display_name,omitempty" url:"display_name,omitempty"`
	// Type: `directory` or `file`.
	Type string `json:"type,omitempty" path:"type,omitempty" url:"type,omitempty"`
	// File/Folder size.
	Size int64 `json:"size,omitempty" path:"size,omitempty" url:"size,omitempty"`
	// File created date/time.
	CreatedAt *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	// User ID of the User who last modified the file/folder.
	LastModifiedById int64 `json:"last_modified_by_id,omitempty" path:"last_modified_by_id,omitempty" url:"last_modified_by_id,omitempty"`
	// ID of the API key that last modified the file/folder.
	LastModifiedByApiKeyId int64 `json:"last_modified_by_api_key_id,omitempty" path:"last_modified_by_api_key_id,omitempty" url:"last_modified_by_api_key_id,omitempty"`
	// ID of the Automation that last modified the file/folder.
	LastModifiedByAutomationId int64 `json:"last_modified_by_automation_id,omitempty" path:"last_modified_by_automation_id,omitempty" url:"last_modified_by_automation_id,omitempty"`
	// ID of the Bundle Registration that last modified the file/folder.
	LastModifiedByBundleRegistrationId int64 `json:"last_modified_by_bundle_registration_id,omitempty" path:"last_modified_by_bundle_registration_id,omitempty" url:"last_modified_by_bundle_registration_id,omitempty"`
	// ID of the Remote Server that last modified the file/folder.
	LastModifiedByRemoteServerId int64 `json:"last_modified_by_remote_server_id,omitempty" path:"last_modified_by_remote_server_id,omitempty" url:"last_modified_by_remote_server_id,omitempty"`
	// ID of the Sync that last modified the file/folder.
	LastModifiedBySyncId int64 `json:"last_modified_by_sync_id,omitempty" path:"last_modified_by_sync_id,omitempty" url:"last_modified_by_sync_id,omitempty"`
	// File last modified date/time, according to the server. This is the timestamp of the last Files.com operation of the file, regardless of what modified timestamp was sent.
	Mtime *time.Time `json:"mtime,omitempty" path:"mtime,omitempty" url:"mtime,omitempty"`
	// File last modified date/time, according to the client who set it. Files.com allows desktop, FTP, SFTP, and WebDAV clients to set modified at times. This allows Desktop<->Cloud syncing to preserve modified at times.
	ProvidedMtime *time.Time `json:"provided_mtime,omitempty" path:"provided_mtime,omitempty" url:"provided_mtime,omitempty"`
	// File CRC32 checksum. This is sometimes delayed, so if you get a blank response, wait and try again.
	Crc32 string `json:"crc32,omitempty" path:"crc32,omitempty" url:"crc32,omitempty"`
	// File MD5 checksum. This is sometimes delayed, so if you get a blank response, wait and try again.
	Md5 string `json:"md5,omitempty" path:"md5,omitempty" url:"md5,omitempty"`
	// File SHA1 checksum. This is sometimes delayed, so if you get a blank response, wait and try again.
	Sha1 string `json:"sha1,omitempty" path:"sha1,omitempty" url:"sha1,omitempty"`
	// File SHA256 checksum. This is sometimes delayed, so if you get a blank response, wait and try again.
	Sha256 string `json:"sha256,omitempty" path:"sha256,omitempty" url:"sha256,omitempty"`
	// MIME Type. This is determined by the filename extension and is not stored separately internally.
	MimeType string `json:"mime_type,omitempty" path:"mime_type,omitempty" url:"mime_type,omitempty"`
	// Region location.
	Region string `json:"region,omitempty" path:"region,omitempty" url:"region,omitempty"`
	// A short string representing the current user's permissions. Can be `r` (Read),`w` (Write),`d` (Delete), `l` (List) or any combination.
	Permissions string `json:"permissions,omitempty" path:"permissions,omitempty" url:"permissions,omitempty"`
	// Are subfolders locked and unable to be modified?
	SubfoldersLocked *bool `json:"subfolders_locked?,omitempty" path:"subfolders_locked?,omitempty" url:"subfolders_locked?,omitempty"`
	// Is this folder locked and unable to be modified?
	IsLocked *bool `json:"is_locked,omitempty" path:"is_locked,omitempty" url:"is_locked,omitempty"`
	// Link to download file. Provided only in response to a download request.
	DownloadUri string `json:"download_uri,omitempty" path:"download_uri,omitempty" url:"download_uri,omitempty"`
	// Optional direct connection information for direct Agent transfer attempts.
	DirectConnectionInfo DirectConnectionInfo `json:"direct_connection_info,omitempty" path:"direct_connection_info,omitempty" url:"direct_connection_info,omitempty"`
	// Bookmark/priority color of file/folder.
	PriorityColor string `json:"priority_color,omitempty" path:"priority_color,omitempty" url:"priority_color,omitempty"`
	// File preview ID.
	PreviewId int64 `json:"preview_id,omitempty" path:"preview_id,omitempty" url:"preview_id,omitempty"`
	// File preview.
	Preview Preview `json:"preview,omitempty" path:"preview,omitempty" url:"preview,omitempty"`
	// The action to perform. Can be `append`, `attachment`, `end`, `upload`, `put`, or may not exist.
	Action string `json:"action,omitempty" path:"action,omitempty" url:"action,omitempty"`
	// Length of file.
	Length int64 `json:"length,omitempty" path:"length,omitempty" url:"length,omitempty"`
	// Create parent directories if they do not exist?
	MkdirParents *bool `json:"mkdir_parents,omitempty" path:"mkdir_parents,omitempty" url:"mkdir_parents,omitempty"`
	// Part if uploading a part.
	Part int64 `json:"part,omitempty" path:"part,omitempty" url:"part,omitempty"`
	// How many parts to fetch?
	Parts int64  `json:"parts,omitempty" path:"parts,omitempty" url:"parts,omitempty"`
	Ref   string `json:"ref,omitempty" path:"ref,omitempty" url:"ref,omitempty"`
	// File byte offset to restart from.
	Restart int64 `json:"restart,omitempty" path:"restart,omitempty" url:"restart,omitempty"`
	// If copying a folder, also copy supported behaviors, email notification subscriptions, and per-folder branding to the destination folder tree?
	CopyBehaviors *bool `json:"copy_behaviors,omitempty" path:"copy_behaviors,omitempty" url:"copy_behaviors,omitempty"`
	// If copying folder, copy just the structure?
	Structure string `json:"structure,omitempty" path:"structure,omitempty" url:"structure,omitempty"`
	// Allow file rename instead of overwrite?
	WithRename *bool `json:"with_rename,omitempty" path:"with_rename,omitempty" url:"with_rename,omitempty"`
	// If true, and the path refers to a destination not stored on Files.com (such as a remote server mount), the upload will be uploaded first to Files.com before being sent to the remote server mount. This can allow clients to upload using parallel parts to a remote server destination that does not offer parallel parts support natively.
	BufferedUpload *bool `json:"buffered_upload,omitempty" path:"buffered_upload,omitempty" url:"buffered_upload,omitempty"`
	// Include optional direct connection information for a direct Agent transfer attempt?
	WithDirectConnectionInfo *bool `json:"with_direct_connection_info,omitempty" path:"with_direct_connection_info,omitempty" url:"with_direct_connection_info,omitempty"`
}

// Identifier returns the resource path.
func (f File) Identifier() interface{} {
	return f.Path
}

// FileCollection is a list of File resources.
type FileCollection []File

// EtagsParam contains the nested request parameters named by its fields.
type EtagsParam struct {
	// etag identifier. Required in FileCreateParams.
	Etag string `url:"etag,omitempty" json:"etag,omitempty" path:"etag"`
	// Part number. Required in FileCreateParams.
	Part string `url:"part,omitempty" json:"part,omitempty" path:"part"`
}

// FileDownloadParams contains the request parameters for GET /files/{path}.
//
// Download File
type FileDownloadParams struct {
	// Path to operate on. Required; sent in the request path.
	Path string `url:"-,omitempty" json:"-,omitempty" path:"path"`
	// Can be blank, `redirect` or `stat`. If set to `stat`, we will return file information but without a download URL, and without logging a download. If set to `redirect` we will serve a 302 redirect directly to the file. This is used for integrations with Zapier, and is not recommended for most integrations.
	Action string `url:"action,omitempty" json:"action,omitempty" path:"action"`
	// Request a preview size. Can be `small` (default), `large`, `xlarge`, or `pdf`.
	PreviewSize string `url:"preview_size,omitempty" json:"preview_size,omitempty" path:"preview_size"`
	// Include file preview information?
	WithPreviews *bool `url:"with_previews,omitempty" json:"with_previews,omitempty" path:"with_previews"`
	// Include file priority color information?
	WithPriorityColor *bool `url:"with_priority_color,omitempty" json:"with_priority_color,omitempty" path:"with_priority_color"`
	// Include optional direct connection information for a direct Agent transfer attempt?
	WithDirectConnectionInfo *bool `url:"with_direct_connection_info,omitempty" json:"with_direct_connection_info,omitempty" path:"with_direct_connection_info"`
	// Not sent to the API.
	File File `url:"-" required:"false" json:"-"`
}

// FileCreateParams contains the request parameters for POST /files/{path}.
type FileCreateParams struct {
	// Path to operate on. Required; sent in the request path.
	Path string `url:"-,omitempty" json:"-,omitempty" path:"path"`
	// The action to perform. Can be `append`, `attachment`, `end`, `upload`, `put`, or may not exist.
	Action string `url:"action,omitempty" json:"action,omitempty" path:"action"`
	// Custom metadata map to save when `action=end` completes the upload. Replaces existing metadata; an empty map clears it. No separate metadata-edit permission is required. Supported on native files and configured remote mounts, excluding remote server automount paths. Limited to 32 keys, 256 characters per key and 1024 characters per value.
	CustomMetadata interface{} `url:"custom_metadata,omitempty" json:"custom_metadata,omitempty" path:"custom_metadata"`
	// Nested parameters etags[etag] and etags[part].
	// etags[etag]: etag identifier. Required.
	// etags[part]: Part number. Required.
	EtagsParam []EtagsParam `url:"etags,omitempty" json:"etags,omitempty" path:"etags"`
	// Length of file.
	Length int64 `url:"length,omitempty" json:"length,omitempty" path:"length"`
	// Create parent directories if they do not exist?
	MkdirParents *bool `url:"mkdir_parents,omitempty" json:"mkdir_parents,omitempty" path:"mkdir_parents"`
	// Part if uploading a part.
	Part int64 `url:"part,omitempty" json:"part,omitempty" path:"part"`
	// How many parts to fetch?
	Parts int64 `url:"parts,omitempty" json:"parts,omitempty" path:"parts"`
	// User provided modification time.
	ProvidedMtime *time.Time `url:"provided_mtime,omitempty" json:"provided_mtime,omitempty" path:"provided_mtime"`
	Ref           string     `url:"ref,omitempty" json:"ref,omitempty" path:"ref"`
	// File byte offset to restart from.
	Restart int64 `url:"restart,omitempty" json:"restart,omitempty" path:"restart"`
	// Size of file.
	Size int64 `url:"size,omitempty" json:"size,omitempty" path:"size"`
	// If copying a folder, also copy supported behaviors, email notification subscriptions, and per-folder branding to the destination folder tree?
	CopyBehaviors *bool `url:"copy_behaviors,omitempty" json:"copy_behaviors,omitempty" path:"copy_behaviors"`
	// If copying folder, copy just the structure?
	Structure string `url:"structure,omitempty" json:"structure,omitempty" path:"structure"`
	// Allow file rename instead of overwrite?
	WithRename *bool `url:"with_rename,omitempty" json:"with_rename,omitempty" path:"with_rename"`
	// If true, and the path refers to a destination not stored on Files.com (such as a remote server mount), the upload will be uploaded first to Files.com before being sent to the remote server mount. This can allow clients to upload using parallel parts to a remote server destination that does not offer parallel parts support natively.
	BufferedUpload *bool `url:"buffered_upload,omitempty" json:"buffered_upload,omitempty" path:"buffered_upload"`
	// Include optional direct connection information for a direct Agent transfer attempt?
	WithDirectConnectionInfo *bool `url:"with_direct_connection_info,omitempty" json:"with_direct_connection_info,omitempty" path:"with_direct_connection_info"`
	// Sent to the API as action_attributes.
	ActionAttributes map[string]any `url:"action_attributes,omitempty" json:"action_attributes,omitempty" path:"action_attributes"`
}

// FileUpdateParams contains the request parameters for PATCH /files/{path}.
type FileUpdateParams struct {
	// Path to operate on. Required; sent in the request path.
	Path string `url:"-,omitempty" json:"-,omitempty" path:"path"`
	// Custom metadata map of keys and values. Limited to 32 keys, 256 characters per key and 1024 characters per value.
	CustomMetadata interface{} `url:"custom_metadata,omitempty" json:"custom_metadata,omitempty" path:"custom_metadata"`
	// Modified time of file.
	ProvidedMtime *time.Time `url:"provided_mtime,omitempty" json:"provided_mtime,omitempty" path:"provided_mtime"`
	// Priority/Bookmark color of file.
	PriorityColor string `url:"priority_color,omitempty" json:"priority_color,omitempty" path:"priority_color"`
}

// FileDeleteParams contains the request parameters for DELETE /files/{path}.
type FileDeleteParams struct {
	// Path to operate on. Required; sent in the request path.
	Path string `url:"-,omitempty" json:"-,omitempty" path:"path"`
	// If true, will recursively delete folders. Otherwise, will error on non-empty folders.
	Recursive *bool `url:"recursive,omitempty" json:"recursive,omitempty" path:"recursive"`
}

// FileFindParams contains the request parameters for GET /file_actions/metadata/{path}.
type FileFindParams struct {
	// Path to operate on. Required; sent in the request path.
	Path string `url:"-,omitempty" json:"-,omitempty" path:"path"`
	// Request a preview size. Can be `small` (default), `large`, `xlarge`, or `pdf`.
	PreviewSize string `url:"preview_size,omitempty" json:"preview_size,omitempty" path:"preview_size"`
	// Include file preview information?
	WithPreviews *bool `url:"with_previews,omitempty" json:"with_previews,omitempty" path:"with_previews"`
	// Include file priority color information?
	WithPriorityColor *bool `url:"with_priority_color,omitempty" json:"with_priority_color,omitempty" path:"with_priority_color"`
}

// FileZipListContentsParams contains the request parameters for GET /file_actions/zip_list/{path}.
//
// List the contents of a ZIP file
type FileZipListContentsParams struct {
	// Path to operate on. Required; sent in the request path.
	Path string `url:"-,omitempty" json:"-,omitempty" path:"path"`
}

// FileBeginDeleteParams contains the request parameters for POST /file_actions/delete/{path}.
//
// Delete File/Folder and return its File Action.  A pending or processing deletion continues in the background after this returns; check its File Migration for completion.  Waiting for the File Migration does not cancel the deletion
type FileBeginDeleteParams struct {
	// Path to operate on. Required; sent in the request path.
	Path string `url:"-,omitempty" json:"-,omitempty" path:"path"`
	// If true, will recursively delete folders. Otherwise, will error on non-empty folders.
	Recursive *bool `url:"recursive,omitempty" json:"recursive,omitempty" path:"recursive"`
}

// FileCopyParams contains the request parameters for POST /file_actions/copy/{path}.
//
// Copy File/Folder
type FileCopyParams struct {
	// Path to operate on. Required; sent in the request path.
	Path string `url:"-,omitempty" json:"-,omitempty" path:"path"`
	// Copy destination path. Required.
	Destination string `url:"destination" json:"destination" path:"destination"`
	// If copying a folder, also copy supported behaviors, email notification subscriptions, and per-folder branding to the destination folder tree?
	CopyBehaviors *bool `url:"copy_behaviors,omitempty" json:"copy_behaviors,omitempty" path:"copy_behaviors"`
	// Copy structure only?
	Structure *bool `url:"structure,omitempty" json:"structure,omitempty" path:"structure"`
	// Overwrite existing file(s) in the destination?
	Overwrite *bool `url:"overwrite,omitempty" json:"overwrite,omitempty" path:"overwrite"`
}

// FileMoveParams contains the request parameters for POST /file_actions/move/{path}.
//
// Move File/Folder
type FileMoveParams struct {
	// Path to operate on. Required; sent in the request path.
	Path string `url:"-,omitempty" json:"-,omitempty" path:"path"`
	// Move destination path. Required.
	Destination string `url:"destination" json:"destination" path:"destination"`
	// Overwrite existing file(s) in the destination?
	Overwrite *bool `url:"overwrite,omitempty" json:"overwrite,omitempty" path:"overwrite"`
}

// FileTransformParams contains the request parameters for POST /file_actions/transform/{path}.
//
// Transform a file and save the output to a destination path
type FileTransformParams struct {
	// Path to operate on. Required; sent in the request path.
	Path string `url:"-,omitempty" json:"-,omitempty" path:"path"`
	// Destination file path for the transformed output. Required.
	Destination string `url:"destination" json:"destination" path:"destination"`
	// Transform type. Supported values are `image_convert`, `document_convert`, and `files_transform_script_execute`. Required.
	TransformType string `url:"transform_type" json:"transform_type" path:"transform_type"`
	// Destination format to create. Required.
	TargetFormat string `url:"target_format" json:"target_format" path:"target_format"`
	// Files TransformScript source. Required when transform_type is `files_transform_script_execute`.
	Script string `url:"script,omitempty" json:"script,omitempty" path:"script"`
	// Maximum output width for image_convert.
	Width int64 `url:"width,omitempty" json:"width,omitempty" path:"width"`
	// Maximum output height for image_convert.
	Height int64 `url:"height,omitempty" json:"height,omitempty" path:"height"`
	// Overwrite existing file in the destination?
	Overwrite *bool `url:"overwrite,omitempty" json:"overwrite,omitempty" path:"overwrite"`
}

// FileGpgDecryptParams contains the request parameters for POST /file_actions/gpg_decrypt/{path}.
//
// Decrypt a GPG-encrypted file and save it to a destination path
type FileGpgDecryptParams struct {
	// Path to operate on. Required; sent in the request path.
	Path string `url:"-,omitempty" json:"-,omitempty" path:"path"`
	// Destination file path for the decrypted file. Required.
	Destination string `url:"destination" json:"destination" path:"destination"`
	// GPG Key IDs to decrypt with. If omitted, every accessible private GPG key in the source workspace is used.
	GpgKeyIds []int64 `url:"gpg_key_ids,omitempty" json:"gpg_key_ids,omitempty" path:"gpg_key_ids"`
	// Partner ID whose GPG keys should be used for decryption.
	GpgKeyPartnerId int64 `url:"gpg_key_partner_id,omitempty" json:"gpg_key_partner_id,omitempty" path:"gpg_key_partner_id"`
	// Use every accessible private GPG key in the source workspace for decryption.
	UseAllPrivateKeys *bool `url:"use_all_private_keys,omitempty" json:"use_all_private_keys,omitempty" path:"use_all_private_keys"`
	// Ignore errors from the MDC (modification detection code) check.
	IgnoreMdcError *bool `url:"ignore_mdc_error,omitempty" json:"ignore_mdc_error,omitempty" path:"ignore_mdc_error"`
	// Overwrite existing file in the destination?
	Overwrite *bool `url:"overwrite,omitempty" json:"overwrite,omitempty" path:"overwrite"`
}

// FileGpgEncryptParams contains the request parameters for POST /file_actions/gpg_encrypt/{path}.
//
// Encrypt a file with GPG and save it to a destination path
type FileGpgEncryptParams struct {
	// Path to operate on. Required; sent in the request path.
	Path string `url:"-,omitempty" json:"-,omitempty" path:"path"`
	// Destination file path for the encrypted file. Required.
	Destination string `url:"destination" json:"destination" path:"destination"`
	// GPG Key IDs to encrypt with.
	GpgKeyIds []int64 `url:"gpg_key_ids,omitempty" json:"gpg_key_ids,omitempty" path:"gpg_key_ids"`
	// Partner ID whose GPG keys should be used for encryption.
	GpgKeyPartnerId int64 `url:"gpg_key_partner_id,omitempty" json:"gpg_key_partner_id,omitempty" path:"gpg_key_partner_id"`
	// Optional GPG Key ID to sign with.
	SigningKeyId int64 `url:"signing_key_id,omitempty" json:"signing_key_id,omitempty" path:"signing_key_id"`
	// Output ASCII-armored encrypted data.
	Armor *bool `url:"armor,omitempty" json:"armor,omitempty" path:"armor"`
	// Overwrite existing file in the destination?
	Overwrite *bool `url:"overwrite,omitempty" json:"overwrite,omitempty" path:"overwrite"`
}

// FileUnzipParams contains the request parameters for POST /file_actions/unzip.
//
// Extract a ZIP file to a destination folder
type FileUnzipParams struct {
	// ZIP file path to extract. Required.
	Path string `url:"path" json:"path" path:"path"`
	// Destination folder path for extracted files. Required.
	Destination string `url:"destination" json:"destination" path:"destination"`
	// Optional single entry filename to extract.
	Filename string `url:"filename,omitempty" json:"filename,omitempty" path:"filename"`
	// Overwrite existing files in the destination?
	Overwrite *bool `url:"overwrite,omitempty" json:"overwrite,omitempty" path:"overwrite"`
}

// FileZipParams contains the request parameters for POST /file_actions/zip.
type FileZipParams struct {
	// Paths to include in the ZIP. Required.
	Paths []string `url:"paths" json:"paths" path:"paths"`
	// Destination file path for the ZIP. Required.
	Destination string `url:"destination" json:"destination" path:"destination"`
	// Overwrite existing file in the destination?
	Overwrite *bool `url:"overwrite,omitempty" json:"overwrite,omitempty" path:"overwrite"`
}

// FileBeginUploadParams contains the request parameters for POST /file_actions/begin_upload/{path}.
//
// Begin File Upload
type FileBeginUploadParams struct {
	// Path to operate on. Required; sent in the request path.
	Path string `url:"-,omitempty" json:"-,omitempty" path:"path"`
	// Create parent directories if they do not exist?
	MkdirParents *bool `url:"mkdir_parents,omitempty" json:"mkdir_parents,omitempty" path:"mkdir_parents"`
	// Part if uploading a part.
	Part int64 `url:"part,omitempty" json:"part,omitempty" path:"part"`
	// How many parts to fetch?
	Parts int64  `url:"parts,omitempty" json:"parts,omitempty" path:"parts"`
	Ref   string `url:"ref,omitempty" json:"ref,omitempty" path:"ref"`
	// File byte offset to restart from.
	Restart int64 `url:"restart,omitempty" json:"restart,omitempty" path:"restart"`
	// Total bytes of file being uploaded (include bytes being retained if appending/restarting).
	Size int64 `url:"size,omitempty" json:"size,omitempty" path:"size"`
	// Allow file rename instead of overwrite?
	WithRename *bool `url:"with_rename,omitempty" json:"with_rename,omitempty" path:"with_rename"`
	// If true, and the path refers to a destination not stored on Files.com (such as a remote server mount), the upload will be uploaded first to Files.com before being sent to the remote server mount. This can allow clients to upload using parallel parts to a remote server destination that does not offer parallel parts support natively.
	BufferedUpload *bool `url:"buffered_upload,omitempty" json:"buffered_upload,omitempty" path:"buffered_upload"`
	// Include optional direct connection information for a direct Agent transfer attempt?
	WithDirectConnectionInfo *bool `url:"with_direct_connection_info,omitempty" json:"with_direct_connection_info,omitempty" path:"with_direct_connection_info"`
}

// ToFolder copies the fields shared with Folder through their JSON representation.
// It does not check whether f is a directory.
func (f File) ToFolder() (Folder, error) {
	bodyBytes, err := json.Marshal(f)
	if err != nil {
		return Folder{}, err
	}
	folder := Folder{}
	folder.UnmarshalJSON(bodyBytes)
	return folder, nil
}

// String returns the Files.com path.
func (f File) String() string {
	return f.Path
}

// Iterable reports whether f is a directory whose children can be listed.
func (f File) Iterable() bool {
	return f.IsDir()
}

// IsDir reports whether the API resource type is "directory".
func (f File) IsDir() bool {
	return f.Type == "directory"
}

// ModTime returns ProvidedMtime when present, then Mtime, or the zero time
// when neither timestamp is available.
func (f File) ModTime() time.Time {
	if f.ProvidedMtime != nil {
		return *f.ProvidedMtime
	}
	if f.Mtime != nil {
		return *f.Mtime
	}
	return time.Time{}
}

// CreationTime returns CreatedAt, or the zero time when it is unavailable.
func (f File) CreationTime() time.Time {
	if f.CreatedAt != nil {
		return *f.CreatedAt
	}
	return time.Time{}
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (f *File) UnmarshalJSON(data []byte) error {
	type file File
	var v file
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*f = File(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (f *FileCollection) UnmarshalJSON(data []byte) error {
	type files FileCollection
	var v files
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*f = FileCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (f *FileCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*f))
	for i, v := range *f {
		ret[i] = v
	}

	return &ret
}
