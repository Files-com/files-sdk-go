package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// ApiKey is a Files.com API resource.
type ApiKey struct {
	// API Key ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Unique label that describes this API key. Useful for external systems where you may have API keys from multiple accounts and want a human-readable label for each key.
	DescriptiveLabel string `json:"descriptive_label,omitempty" path:"descriptive_label,omitempty" url:"descriptive_label,omitempty"`
	// User-supplied description of API key.
	Description string `json:"description,omitempty" path:"description,omitempty" url:"description,omitempty"`
	// Time which API Key was created.
	CreatedAt *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	// API Key expiration date.
	ExpiresAt *time.Time `json:"expires_at,omitempty" path:"expires_at,omitempty" url:"expires_at,omitempty"`
	// API Key actual key string.
	Key string `json:"key,omitempty" path:"key,omitempty" url:"key,omitempty"`
	// If `true`, this API key will be usable with AWS-compatible endpoints, such as our Inbound S3-compatible endpoint.
	AwsStyleCredentials *bool `json:"aws_style_credentials,omitempty" path:"aws_style_credentials,omitempty" url:"aws_style_credentials,omitempty"`
	// AWS Access Key ID to use with AWS-compatible endpoints, such as our Inbound S3-compatible endpoint.
	AwsAccessKeyId string `json:"aws_access_key_id,omitempty" path:"aws_access_key_id,omitempty" url:"aws_access_key_id,omitempty"`
	// AWS Secret Key to use with AWS-compatible endpoints, such as our Inbound S3-compatible endpoint.
	AwsSecretKey string `json:"aws_secret_key,omitempty" path:"aws_secret_key,omitempty" url:"aws_secret_key,omitempty"`
	// API Key last used - note this value is only updated once per 3 hour period, so the 'actual' time of last use may be up to 3 hours later than this timestamp.
	LastUseAt *time.Time `json:"last_use_at,omitempty" path:"last_use_at,omitempty" url:"last_use_at,omitempty"`
	// Internal name for the API Key. For your use.
	Name string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	// Permissions for this API Key. Keys with the `desktop_app` permission set only have the ability to do the functions provided in our Desktop App (File and Share Link operations). Keys with the `mobile_app` permission set can only perform the functions provided by our Mobile App (File and Share Link operations). Keys with the `office_integration` permission set are auto generated, and automatically expire, to allow users to interact with office integration platforms. Keys with the `files_only` permission set can use only the files, folders, and file_actions endpoints, excluding the GPG Encrypt and GPG Decrypt file actions, where they perform file operations as a full-access file user in the key's workspace scope, along with `GET /file_migrations/{id}` and `GET /api_key`. On the migration lookup, any key with a `path` can read only migrations it started. Without a `path`, user-owned keys retain user-scoped access; site-wide Workspace-bound keys can read only migrations whose source and destination are both in their Workspace. Every `files_only` key is Workspace-bound, including Workspace `0`, as is a full-access key in a named Workspace. A full-access site-wide key in the default Workspace retains Site-wide migration access. Migrations outside the caller's visibility return `not-found`. Keys with `files_only` cannot use site admin, workspace admin, folder admin, group admin, partner admin, or billing privileges from the owning user. Both GPG file actions and every other endpoint deny them with `not-authorized/api-key-only-for-file-operations`. Folder-configured GPG encryption and decryption still apply to files uploaded with a `files_only` key.
	PermissionSet string `json:"permission_set,omitempty" path:"permission_set,omitempty" url:"permission_set,omitempty"`
	// If this API key represents a Desktop app, what platform was it created on?
	Platform string `json:"platform,omitempty" path:"platform,omitempty" url:"platform,omitempty"`
	// Site ID.
	SiteId int64 `json:"site_id,omitempty" path:"site_id,omitempty" url:"site_id,omitempty"`
	// Site Name.
	SiteName string `json:"site_name,omitempty" path:"site_name,omitempty" url:"site_name,omitempty"`
	// URL for API host.
	Url string `json:"url,omitempty" path:"url,omitempty" url:"url,omitempty"`
	// User ID for the owner of this API Key. May be blank for Site-wide API Keys.
	UserId int64 `json:"user_id,omitempty" path:"user_id,omitempty" url:"user_id,omitempty"`
	// Workspace ID for this API Key. `0` means the default workspace.
	WorkspaceId int64 `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	// Restricts the file and folder operations made with this key, meaning the files, folders, and file_actions endpoints, to the specified folder and its descendants, including copy and move destinations. For `GET /file_migrations/{id}`, a key with a path can read only migrations it started. Does not grant access beyond the owning user's permissions. Optional for `files_only` keys and required for `office_integration` keys, which need a path the owning user can read. Keys with any other permission set cannot have a path.
	Path string `json:"path,omitempty" path:"path,omitempty" url:"path,omitempty"`
}

// Identifier returns the resource ID.
func (a ApiKey) Identifier() interface{} {
	return a.Id
}

// ApiKeyCollection is a list of ApiKey resources.
type ApiKeyCollection []ApiKey

// ApiKeyPermissionSetEnum is a string value for permission_set.
// Enum lists the values documented by the API.
type ApiKeyPermissionSetEnum string

// String returns the API parameter value.
func (u ApiKeyPermissionSetEnum) String() string {
	return string(u)
}

// Enum returns the documented values keyed by their API strings.
func (u ApiKeyPermissionSetEnum) Enum() map[string]ApiKeyPermissionSetEnum {
	return map[string]ApiKeyPermissionSetEnum{
		"none":               ApiKeyPermissionSetEnum("none"),
		"full":               ApiKeyPermissionSetEnum("full"),
		"desktop_app":        ApiKeyPermissionSetEnum("desktop_app"),
		"sync_app":           ApiKeyPermissionSetEnum("sync_app"),
		"office_integration": ApiKeyPermissionSetEnum("office_integration"),
		"mobile_app":         ApiKeyPermissionSetEnum("mobile_app"),
		"files_only":         ApiKeyPermissionSetEnum("files_only"),
	}
}

// ApiKeyListParams contains the request parameters for GET /api_keys.
type ApiKeyListParams struct {
	// User ID. Provide a value of `0` to operate the current session's user.
	UserId int64 `url:"user_id,omitempty" json:"user_id,omitempty" path:"user_id"`
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `site_id` and `workspace_id`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `aws_style_credentials` and `expires_at`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// If set, return records where the specified field is greater than the supplied value. Valid fields are `expires_at`.
	FilterGt interface{} `url:"filter_gt,omitempty" json:"filter_gt,omitempty" path:"filter_gt"`
	// If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `expires_at`.
	FilterGteq interface{} `url:"filter_gteq,omitempty" json:"filter_gteq,omitempty" path:"filter_gteq"`
	// If set, return records where the specified field is less than the supplied value. Valid fields are `expires_at`.
	FilterLt interface{} `url:"filter_lt,omitempty" json:"filter_lt,omitempty" path:"filter_lt"`
	// If set, return records where the specified field is less than or equal the supplied value. Valid fields are `expires_at`.
	FilterLteq interface{} `url:"filter_lteq,omitempty" json:"filter_lteq,omitempty" path:"filter_lteq"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// ApiKeyFindParams contains the request parameters for GET /api_keys/{id}.
type ApiKeyFindParams struct {
	// Api Key ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// ApiKeyCreateParams contains the request parameters for POST /api_keys.
type ApiKeyCreateParams struct {
	// User ID. Provide a value of `0` to operate the current session's user.
	UserId int64 `url:"user_id,omitempty" json:"user_id,omitempty" path:"user_id"`
	// User-supplied description of API key.
	Description string `url:"description,omitempty" json:"description,omitempty" path:"description"`
	// API Key expiration date.
	ExpiresAt *time.Time `url:"expires_at,omitempty" json:"expires_at,omitempty" path:"expires_at"`
	// Internal name for the API Key. For your use. Required.
	Name string `url:"name" json:"name" path:"name"`
	// If `true`, this API key will be usable with AWS-compatible endpoints, such as our Inbound S3-compatible endpoint.
	AwsStyleCredentials *bool `url:"aws_style_credentials,omitempty" json:"aws_style_credentials,omitempty" path:"aws_style_credentials"`
	// Restricts the file and folder operations made with this key, meaning the files, folders, and file_actions endpoints, to the specified folder and its descendants, including copy and move destinations. For `GET /file_migrations/{id}`, a key with a path can read only migrations it started. Does not grant access beyond the owning user's permissions. Optional for `files_only` keys and required for `office_integration` keys, which need a path the owning user can read. Keys with any other permission set cannot have a path.
	Path string `url:"path,omitempty" json:"path,omitempty" path:"path"`
	// Permissions for this API Key. Keys with the `desktop_app` permission set only have the ability to do the functions provided in our Desktop App (File and Share Link operations). Keys with the `mobile_app` permission set can only perform the functions provided by our Mobile App (File and Share Link operations). Keys with the `office_integration` permission set are auto generated, and automatically expire, to allow users to interact with office integration platforms. Keys with the `files_only` permission set can use only the files, folders, and file_actions endpoints, excluding the GPG Encrypt and GPG Decrypt file actions, where they perform file operations as a full-access file user in the key's workspace scope, along with `GET /file_migrations/{id}` and `GET /api_key`. On the migration lookup, any key with a `path` can read only migrations it started. Without a `path`, user-owned keys retain user-scoped access; site-wide Workspace-bound keys can read only migrations whose source and destination are both in their Workspace. Every `files_only` key is Workspace-bound, including Workspace `0`, as is a full-access key in a named Workspace. A full-access site-wide key in the default Workspace retains Site-wide migration access. Migrations outside the caller's visibility return `not-found`. Keys with `files_only` cannot use site admin, workspace admin, folder admin, group admin, partner admin, or billing privileges from the owning user. Both GPG file actions and every other endpoint deny them with `not-authorized/api-key-only-for-file-operations`. Folder-configured GPG encryption and decryption still apply to files uploaded with a `files_only` key. Values: none, full, desktop_app, sync_app, office_integration, mobile_app, files_only.
	PermissionSet ApiKeyPermissionSetEnum `url:"permission_set,omitempty" json:"permission_set,omitempty" path:"permission_set"`
	// Workspace ID for this API Key. `0` means the default workspace.
	WorkspaceId int64 `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
}

// ApiKeyUpdateCurrentParams contains the request parameters for PATCH /api_key.
type ApiKeyUpdateCurrentParams struct {
	// API Key expiration date.
	ExpiresAt *time.Time `url:"expires_at,omitempty" json:"expires_at,omitempty" path:"expires_at"`
	// Internal name for the API Key. For your use.
	Name string `url:"name,omitempty" json:"name,omitempty" path:"name"`
	// Permissions for this API Key. Keys with the `desktop_app` permission set only have the ability to do the functions provided in our Desktop App (File and Share Link operations). Keys with the `mobile_app` permission set can only perform the functions provided by our Mobile App (File and Share Link operations). Keys with the `office_integration` permission set are auto generated, and automatically expire, to allow users to interact with office integration platforms. Keys with the `files_only` permission set can use only the files, folders, and file_actions endpoints, excluding the GPG Encrypt and GPG Decrypt file actions, where they perform file operations as a full-access file user in the key's workspace scope, along with `GET /file_migrations/{id}` and `GET /api_key`. On the migration lookup, any key with a `path` can read only migrations it started. Without a `path`, user-owned keys retain user-scoped access; site-wide Workspace-bound keys can read only migrations whose source and destination are both in their Workspace. Every `files_only` key is Workspace-bound, including Workspace `0`, as is a full-access key in a named Workspace. A full-access site-wide key in the default Workspace retains Site-wide migration access. Migrations outside the caller's visibility return `not-found`. Keys with `files_only` cannot use site admin, workspace admin, folder admin, group admin, partner admin, or billing privileges from the owning user. Both GPG file actions and every other endpoint deny them with `not-authorized/api-key-only-for-file-operations`. Folder-configured GPG encryption and decryption still apply to files uploaded with a `files_only` key. Values: none, full, desktop_app, sync_app, office_integration, mobile_app, files_only.
	PermissionSet ApiKeyPermissionSetEnum `url:"permission_set,omitempty" json:"permission_set,omitempty" path:"permission_set"`
}

// ApiKeyUpdateParams contains the request parameters for PATCH /api_keys/{id}.
type ApiKeyUpdateParams struct {
	// Api Key ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// User-supplied description of API key.
	Description string `url:"description,omitempty" json:"description,omitempty" path:"description"`
	// API Key expiration date.
	ExpiresAt *time.Time `url:"expires_at,omitempty" json:"expires_at,omitempty" path:"expires_at"`
	// Internal name for the API Key. For your use.
	Name string `url:"name,omitempty" json:"name,omitempty" path:"name"`
}

// ApiKeyDeleteParams contains the request parameters for DELETE /api_keys/{id}.
type ApiKeyDeleteParams struct {
	// Api Key ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (a *ApiKey) UnmarshalJSON(data []byte) error {
	type apiKey ApiKey
	var v apiKey
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*a = ApiKey(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (a *ApiKeyCollection) UnmarshalJSON(data []byte) error {
	type apiKeys ApiKeyCollection
	var v apiKeys
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*a = ApiKeyCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (a *ApiKeyCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*a))
	for i, v := range *a {
		ret[i] = v
	}

	return &ret
}
