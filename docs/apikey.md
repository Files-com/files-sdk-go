# ApiKey

`files_sdk.ApiKey` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/apikey`, named `api_key`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	api_key "github.com/Files-com/files-sdk-go/v3/apikey"
)
```

## Operations

### List

`GET /api_keys`

List API Keys

```go
func (c *Client) List(params files_sdk.ApiKeyListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.ApiKeyListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.ApiKey` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.ApiKeyListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `UserId` | `int64` | `user_id` | no | User ID. Provide a value of `0` to operate the current session's user. |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `site_id` and `workspace_id`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `aws_style_credentials` and `expires_at`. |
| `FilterGt` | `interface{}` | `filter_gt` | no | If set, return records where the specified field is greater than the supplied value. Valid fields are `expires_at`. |
| `FilterGteq` | `interface{}` | `filter_gteq` | no | If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `expires_at`. |
| `FilterLt` | `interface{}` | `filter_lt` | no | If set, return records where the specified field is less than the supplied value. Valid fields are `expires_at`. |
| `FilterLteq` | `interface{}` | `filter_lteq` | no | If set, return records where the specified field is less than or equal the supplied value. Valid fields are `expires_at`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := api_key.List(files_sdk.ApiKeyListParams{})
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

### FindCurrent

`GET /api_key`

Show information about current API key. (Requires current API connection to be using an API key.)

```go
func (c *Client) FindCurrent(opts ...files_sdk.RequestResponseOption) (apiKey files_sdk.ApiKey, err error)
func FindCurrent(opts ...files_sdk.RequestResponseOption) (apiKey files_sdk.ApiKey, err error)
```

Example:

```go
result, err := api_key.FindCurrent()
if err != nil {
	return err
}
fmt.Println(result)
```

### Find

`GET /api_keys/{id}`

Show API Key

```go
func (c *Client) Find(params files_sdk.ApiKeyFindParams, opts ...files_sdk.RequestResponseOption) (apiKey files_sdk.ApiKey, err error)
func Find(params files_sdk.ApiKeyFindParams, opts ...files_sdk.RequestResponseOption) (apiKey files_sdk.ApiKey, err error)
```

Parameters, `files_sdk.ApiKeyFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Api Key ID. |

Example:

```go
result, err := api_key.Find(files_sdk.ApiKeyFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Create

`POST /api_keys`

Create API Key

```go
func (c *Client) Create(params files_sdk.ApiKeyCreateParams, opts ...files_sdk.RequestResponseOption) (apiKey files_sdk.ApiKey, err error)
func Create(params files_sdk.ApiKeyCreateParams, opts ...files_sdk.RequestResponseOption) (apiKey files_sdk.ApiKey, err error)
```

Parameters, `files_sdk.ApiKeyCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `UserId` | `int64` | `user_id` | no | User ID. Provide a value of `0` to operate the current session's user. |
| `Description` | `string` | `description` | no | User-supplied description of API key. |
| `ExpiresAt` | `*time.Time` | `expires_at` | no | API Key expiration date |
| `Name` | `string` | `name` | yes | Internal name for the API Key. For your use. |
| `AwsStyleCredentials` | `*bool` | `aws_style_credentials` | no | If `true`, this API key will be usable with AWS-compatible endpoints, such as our Inbound S3-compatible endpoint. |
| `Path` | `string` | `path` | no | Restricts the file and folder operations made with this key, meaning the files, folders, and file_actions endpoints, to the specified folder and its descendants, including copy and move destinations. For `GET /file_migrations/{id}`, a key with a path can read only migrations it started. Other endpoints do not apply the path restriction; use the `files_only` permission set to confine a key to file operations and their supporting lookups. Does not grant access beyond the owning user's permissions. Optional except for `office_integration` keys, which require a path the owning user can read. |
| `PermissionSet` | `files_sdk.ApiKeyPermissionSetEnum` | `permission_set` | no | Permissions for this API Key. Keys with the `desktop_app` permission set only have the ability to do the functions provided in our Desktop App (File and Share Link operations). Keys with the `mobile_app` permission set can only perform the functions provided by our Mobile App (File and Share Link operations). Keys with the `office_integration` permission set are auto generated, and automatically expire, to allow users to interact with office integration platforms. Keys with the `files_only` permission set can use only the files, folders, and file_actions endpoints, excluding the GPG Encrypt and GPG Decrypt file actions, where they perform file operations as a full-access file user in the key's workspace scope, along with `GET /file_migrations/{id}` and `GET /api_key`. On the migration lookup, any key with a `path` can read only migrations it started. Without a `path`, user-owned keys retain user-scoped access; site-wide Workspace-bound keys can read only migrations whose source and destination are both in their Workspace. Every `files_only` key is Workspace-bound, including Workspace `0`, as is a full-access key in a named Workspace. A full-access site-wide key in the default Workspace retains Site-wide migration access. Migrations outside the caller's visibility return `not-found`. Keys with `files_only` cannot use site admin, workspace admin, folder admin, group admin, partner admin, or billing privileges from the owning user. Both GPG file actions and every other endpoint deny them with `not-authorized/api-key-only-for-file-operations`. Folder-configured GPG encryption and decryption still apply to files uploaded with a `files_only` key. Values: `none`, `full`, `desktop_app`, `sync_app`, `office_integration`, `mobile_app`, `files_only`. |
| `WorkspaceId` | `int64` | `workspace_id` | no | Workspace ID for this API Key. `0` means the default workspace. |

Example:

```go
result, err := api_key.Create(files_sdk.ApiKeyCreateParams{
	Name: "My Main API Key",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### UpdateCurrent

`PATCH /api_key`

Update current API key. (Requires current API connection to be using an API key.)

```go
func (c *Client) UpdateCurrent(params files_sdk.ApiKeyUpdateCurrentParams, opts ...files_sdk.RequestResponseOption) (apiKey files_sdk.ApiKey, err error)
func UpdateCurrent(params files_sdk.ApiKeyUpdateCurrentParams, opts ...files_sdk.RequestResponseOption) (apiKey files_sdk.ApiKey, err error)
func (c *Client) UpdateCurrentWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (apiKey files_sdk.ApiKey, err error)
func UpdateCurrentWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (apiKey files_sdk.ApiKey, err error)
```

`UpdateCurrentWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.ApiKeyUpdateCurrentParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `ExpiresAt` | `*time.Time` | `expires_at` | no | API Key expiration date |
| `Name` | `string` | `name` | no | Internal name for the API Key. For your use. |
| `PermissionSet` | `files_sdk.ApiKeyPermissionSetEnum` | `permission_set` | no | Permissions for this API Key. Keys with the `desktop_app` permission set only have the ability to do the functions provided in our Desktop App (File and Share Link operations). Keys with the `mobile_app` permission set can only perform the functions provided by our Mobile App (File and Share Link operations). Keys with the `office_integration` permission set are auto generated, and automatically expire, to allow users to interact with office integration platforms. Keys with the `files_only` permission set can use only the files, folders, and file_actions endpoints, excluding the GPG Encrypt and GPG Decrypt file actions, where they perform file operations as a full-access file user in the key's workspace scope, along with `GET /file_migrations/{id}` and `GET /api_key`. On the migration lookup, any key with a `path` can read only migrations it started. Without a `path`, user-owned keys retain user-scoped access; site-wide Workspace-bound keys can read only migrations whose source and destination are both in their Workspace. Every `files_only` key is Workspace-bound, including Workspace `0`, as is a full-access key in a named Workspace. A full-access site-wide key in the default Workspace retains Site-wide migration access. Migrations outside the caller's visibility return `not-found`. Keys with `files_only` cannot use site admin, workspace admin, folder admin, group admin, partner admin, or billing privileges from the owning user. Both GPG file actions and every other endpoint deny them with `not-authorized/api-key-only-for-file-operations`. Folder-configured GPG encryption and decryption still apply to files uploaded with a `files_only` key. Values: `none`, `full`, `desktop_app`, `sync_app`, `office_integration`, `mobile_app`, `files_only`. |

Example:

```go
result, err := api_key.UpdateCurrent(files_sdk.ApiKeyUpdateCurrentParams{})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := api_key.UpdateCurrentWithMap(map[string]interface{}{})
if err != nil {
	return err
}
fmt.Println(result)
```

### Update

`PATCH /api_keys/{id}`

Update API Key

```go
func (c *Client) Update(params files_sdk.ApiKeyUpdateParams, opts ...files_sdk.RequestResponseOption) (apiKey files_sdk.ApiKey, err error)
func Update(params files_sdk.ApiKeyUpdateParams, opts ...files_sdk.RequestResponseOption) (apiKey files_sdk.ApiKey, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (apiKey files_sdk.ApiKey, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (apiKey files_sdk.ApiKey, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.ApiKeyUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Api Key ID. |
| `Description` | `string` | `description` | no | User-supplied description of API key. |
| `ExpiresAt` | `*time.Time` | `expires_at` | no | API Key expiration date |
| `Name` | `string` | `name` | no | Internal name for the API Key. For your use. |

Example:

```go
result, err := api_key.Update(files_sdk.ApiKeyUpdateParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := api_key.UpdateWithMap(map[string]interface{}{
	"id": 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### DeleteCurrent

`DELETE /api_key`

Delete current API key. (Requires current API connection to be using an API key.)

```go
func (c *Client) DeleteCurrent(opts ...files_sdk.RequestResponseOption) (err error)
func DeleteCurrent(opts ...files_sdk.RequestResponseOption) (err error)
```

Example:

```go
err := api_key.DeleteCurrent()
if err != nil {
	return err
}
```

### Delete

`DELETE /api_keys/{id}`

Delete API Key

```go
func (c *Client) Delete(params files_sdk.ApiKeyDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.ApiKeyDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.ApiKeyDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Api Key ID. |

Example:

```go
err := api_key.Delete(files_sdk.ApiKeyDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.ApiKey` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | API Key ID |
| `DescriptiveLabel` | `string` | `descriptive_label` | Unique label that describes this API key. Useful for external systems where you may have API keys from multiple accounts and want a human-readable label for each key. |
| `Description` | `string` | `description` | User-supplied description of API key. |
| `CreatedAt` | `*time.Time` | `created_at` | Time which API Key was created |
| `ExpiresAt` | `*time.Time` | `expires_at` | API Key expiration date |
| `Key` | `string` | `key` | API Key actual key string |
| `AwsStyleCredentials` | `*bool` | `aws_style_credentials` | If `true`, this API key will be usable with AWS-compatible endpoints, such as our Inbound S3-compatible endpoint. |
| `AwsAccessKeyId` | `string` | `aws_access_key_id` | AWS Access Key ID to use with AWS-compatible endpoints, such as our Inbound S3-compatible endpoint. |
| `AwsSecretKey` | `string` | `aws_secret_key` | AWS Secret Key to use with AWS-compatible endpoints, such as our Inbound S3-compatible endpoint. |
| `LastUseAt` | `*time.Time` | `last_use_at` | API Key last used - note this value is only updated once per 3 hour period, so the 'actual' time of last use may be up to 3 hours later than this timestamp. |
| `Name` | `string` | `name` | Internal name for the API Key. For your use. |
| `PermissionSet` | `string` | `permission_set` | Permissions for this API Key. Keys with the `desktop_app` permission set only have the ability to do the functions provided in our Desktop App (File and Share Link operations). Keys with the `mobile_app` permission set can only perform the functions provided by our Mobile App (File and Share Link operations). Keys with the `office_integration` permission set are auto generated, and automatically expire, to allow users to interact with office integration platforms. Keys with the `files_only` permission set can use only the files, folders, and file_actions endpoints, excluding the GPG Encrypt and GPG Decrypt file actions, where they perform file operations as a full-access file user in the key's workspace scope, along with `GET /file_migrations/{id}` and `GET /api_key`. On the migration lookup, any key with a `path` can read only migrations it started. Without a `path`, user-owned keys retain user-scoped access; site-wide Workspace-bound keys can read only migrations whose source and destination are both in their Workspace. Every `files_only` key is Workspace-bound, including Workspace `0`, as is a full-access key in a named Workspace. A full-access site-wide key in the default Workspace retains Site-wide migration access. Migrations outside the caller's visibility return `not-found`. Keys with `files_only` cannot use site admin, workspace admin, folder admin, group admin, partner admin, or billing privileges from the owning user. Both GPG file actions and every other endpoint deny them with `not-authorized/api-key-only-for-file-operations`. Folder-configured GPG encryption and decryption still apply to files uploaded with a `files_only` key. |
| `Platform` | `string` | `platform` | If this API key represents a Desktop app, what platform was it created on? |
| `SiteId` | `int64` | `site_id` | Site ID |
| `SiteName` | `string` | `site_name` | Site Name |
| `Url` | `string` | `url` | URL for API host. |
| `UserId` | `int64` | `user_id` | User ID for the owner of this API Key. May be blank for Site-wide API Keys. |
| `WorkspaceId` | `int64` | `workspace_id` | Workspace ID for this API Key. `0` means the default workspace. |
| `Path` | `string` | `path` | Restricts the file and folder operations made with this key, meaning the files, folders, and file_actions endpoints, to the specified folder and its descendants, including copy and move destinations. For `GET /file_migrations/{id}`, a key with a path can read only migrations it started. Other endpoints do not apply the path restriction; use the `files_only` permission set to confine a key to file operations and their supporting lookups. Does not grant access beyond the owning user's permissions. Optional except for `office_integration` keys, which require a path the owning user can read. |
