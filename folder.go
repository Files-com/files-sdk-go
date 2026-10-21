package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// Folder is a Files.com API resource.
type Folder struct {
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
	// Create parent directories if they do not exist?
	MkdirParents *bool `json:"mkdir_parents,omitempty" path:"mkdir_parents,omitempty" url:"mkdir_parents,omitempty"`
}

// Identifier returns the resource path.
func (f Folder) Identifier() interface{} {
	return f.Path
}

// FolderCollection is a list of Folder resources.
type FolderCollection []Folder

// FolderListForParams contains the request parameters for GET /folders/{path}.
type FolderListForParams struct {
	// Path to operate on. Required; sent in the request path.
	Path string `url:"-,omitempty" json:"-,omitempty" path:"path"`
	// Request a preview size. Can be `small` (default), `large`, `xlarge`, or `pdf`.
	PreviewSize string `url:"preview_size,omitempty" json:"preview_size,omitempty" path:"preview_size"`
	// Search by field and direction. Valid fields are `path`, `size`, `modified_at_datetime`, `provided_modified_at`. Valid directions are `asc` and `desc`. Defaults to `{"path":"asc"}`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If specified, will search the folders/files list by name. Ignores text before last `/`. This is the same API used by the search bar in the web UI when running 'Search This Folder'. Search results are a best effort, not real time, and not guaranteed to perfectly match the latest folder listing. Results may be truncated if more than 1,000 possible matches exist. This field should only be used for ad-hoc (human) searching, and not as part of an automated process.
	Search string `url:"search,omitempty" json:"search,omitempty" path:"search"`
	// If provided, the search string in `search` will search for files where this custom metadata key matches the value sent in `search`. Set this to `*` to allow any metadata key to match the value sent in `search`.
	SearchCustomMetadataKey string `url:"search_custom_metadata_key,omitempty" json:"search_custom_metadata_key,omitempty" path:"search_custom_metadata_key"`
	// Search entire site? If true, we will search the entire site. Do not provide a path when using this parameter. This is the same API used by the search bar in the web UI when running 'Search All Files'. Search results are a best effort, not real time, and not guaranteed to match every file. This field should only be used for ad-hoc (human) searching, and not as part of an automated process.
	SearchAll *bool `url:"search_all,omitempty" json:"search_all,omitempty" path:"search_all"`
	// Include file previews?
	WithPreviews *bool `url:"with_previews,omitempty" json:"with_previews,omitempty" path:"with_previews"`
	// Include file priority color information?
	WithPriorityColor *bool `url:"with_priority_color,omitempty" json:"with_priority_color,omitempty" path:"with_priority_color"`
	// Type of objects to return. Can be `folder` or `file`.
	Type string `url:"type,omitempty" json:"type,omitempty" path:"type"`
	// If provided, will only return files/folders modified after this time. Can be used only in combination with `type` filter.
	ModifiedAtDatetime *time.Time `url:"modified_at_datetime,omitempty" json:"modified_at_datetime,omitempty" path:"modified_at_datetime"`
	// Not sent to the API.
	ConcurrencyManager lib.ConcurrencyManagerWithSubWorker `url:"-" required:"false" json:"-"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// FolderCreateParams contains the request parameters for POST /folders/{path}.
type FolderCreateParams struct {
	// Path to operate on. Required; sent in the request path.
	Path string `url:"-,omitempty" json:"-,omitempty" path:"path"`
	// Create parent directories if they do not exist?
	MkdirParents *bool `url:"mkdir_parents,omitempty" json:"mkdir_parents,omitempty" path:"mkdir_parents"`
	// User provided modification time.
	ProvidedMtime *time.Time `url:"provided_mtime,omitempty" json:"provided_mtime,omitempty" path:"provided_mtime"`
}

// ToFile copies the fields shared with File through their JSON representation.
func (f *Folder) ToFile() (File, error) {
	bodyBytes, err := json.Marshal(f)
	if err != nil {
		return File{}, err
	}
	file := File{}
	file.UnmarshalJSON(bodyBytes)
	return file, nil
}

// IsDir reports whether the API resource type is "directory".
func (f Folder) IsDir() bool {
	return f.Type == "directory"
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (f *Folder) UnmarshalJSON(data []byte) error {
	type folder Folder
	var v folder
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*f = Folder(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (f *FolderCollection) UnmarshalJSON(data []byte) error {
	type folders FolderCollection
	var v folders
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*f = FolderCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (f *FolderCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*f))
	for i, v := range *f {
		ret[i] = v
	}

	return &ret
}
