# Bundle

`files_sdk.Bundle` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/bundle`, named `bundle`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/bundle"
)
```

## Operations

### List

`GET /bundles`

List Share Links

```go
func (c *Client) List(params files_sdk.BundleListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.BundleListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.Bundle` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.BundleListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `UserId` | `int64` | `user_id` | no | User ID. Provide a value of `0` to operate the current session's user. |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `expires_at`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `created_at`, `expires_at`, `code`, `group_id`, `user_id` or `bypasses_site_expiration_rules`. Valid field combinations are `[ group_id, expires_at ]` and `[ user_id, expires_at ]`. |
| `FilterGt` | `interface{}` | `filter_gt` | no | If set, return records where the specified field is greater than the supplied value. Valid fields are `created_at` and `expires_at`. |
| `FilterGteq` | `interface{}` | `filter_gteq` | no | If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `created_at` and `expires_at`. |
| `FilterPrefix` | `interface{}` | `filter_prefix` | no | If set, return records where the specified field is prefixed by the supplied value. Valid fields are `code`. |
| `FilterLt` | `interface{}` | `filter_lt` | no | If set, return records where the specified field is less than the supplied value. Valid fields are `created_at` and `expires_at`. |
| `FilterLteq` | `interface{}` | `filter_lteq` | no | If set, return records where the specified field is less than or equal the supplied value. Valid fields are `created_at` and `expires_at`. |
| `Deleted` | `*bool` | `deleted` | no | If true, only list deleted Share Links. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := bundle.List(files_sdk.BundleListParams{})
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

### Find

`GET /bundles/{id}`

Show Share Link

```go
func (c *Client) Find(params files_sdk.BundleFindParams, opts ...files_sdk.RequestResponseOption) (bundle files_sdk.Bundle, err error)
func Find(params files_sdk.BundleFindParams, opts ...files_sdk.RequestResponseOption) (bundle files_sdk.Bundle, err error)
```

Parameters, `files_sdk.BundleFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Bundle ID. |
| `Deleted` | `*bool` | `deleted` | no | If true, show a deleted Share Link. |

Example:

```go
result, err := bundle.Find(files_sdk.BundleFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Create

`POST /bundles`

Create Share Link

```go
func (c *Client) Create(params files_sdk.BundleCreateParams, opts ...files_sdk.RequestResponseOption) (bundle files_sdk.Bundle, err error)
func Create(params files_sdk.BundleCreateParams, opts ...files_sdk.RequestResponseOption) (bundle files_sdk.Bundle, err error)
```

Parameters, `files_sdk.BundleCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `UserId` | `int64` | `user_id` | no | User ID. Provide a value of `0` to operate the current session's user. |
| `Paths` | `[]string` | `paths` | yes | A list of paths to include in this bundle. |
| `Password` | `string` | `password` | no | Password for this bundle. |
| `BypassesSiteExpirationRules` | `*bool` | `bypasses_site_expiration_rules` | no | If true, this Share Link bypasses site-wide expiration rules. Only site admins may set this. |
| `FormFieldSetId` | `int64` | `form_field_set_id` | no | Id of Form Field Set to use with this bundle |
| `CreateSnapshot` | `*bool` | `create_snapshot` | no | If true, create a snapshot of this bundle's contents. |
| `DontSeparateSubmissionsByFolder` | `*bool` | `dont_separate_submissions_by_folder` | no | Do not create subfolders for files uploaded to this share. Note: there are subtle security pitfalls with allowing anonymous uploads from multiple users to live in the same folder. We strongly discourage use of this option unless absolutely required. |
| `ExpiresAt` | `*time.Time` | `expires_at` | no | Explicit Bundle expiration date/time. If not set, the site-wide expiration setting may apply. |
| `FinalizeSnapshot` | `*bool` | `finalize_snapshot` | no | If true, finalize the snapshot of this bundle's contents. Note that `create_snapshot` must also be true. |
| `MaxUses` | `int64` | `max_uses` | no | Maximum number of times bundle can be accessed |
| `GroupId` | `int64` | `group_id` | no | Owning group ID. If set, members of this group can view, edit, and share this Share Link. |
| `InternalName` | `string` | `internal_name` | no | Internal name for identifying this Share Link. |
| `Description` | `string` | `description` | no | Public description |
| `Note` | `string` | `note` | no | Bundle internal note |
| `Code` | `string` | `code` | no | Bundle code. This code forms the end part of the Public URL. |
| `PathTemplate` | `string` | `path_template` | no | Template for creating submission subfolders. Can use the uploader's name, email address, ip, company, `strftime` directives, and any custom form data. |
| `PathTemplateTimeZone` | `string` | `path_template_time_zone` | no | Timezone to use when rendering timestamps in path templates. |
| `Permissions` | `files_sdk.BundlePermissionsEnum` | `permissions` | no | Permissions that apply to Folders in this Share Link. Values: `read`, `write`, `read_write`, `full`, `none`, `preview_only`. |
| `RequireRegistration` | `*bool` | `require_registration` | no | Show a registration page that captures the downloader's name and email address? |
| `ClickwrapId` | `int64` | `clickwrap_id` | no | ID of the clickwrap to use with this bundle. |
| `InboxId` | `int64` | `inbox_id` | no | ID of the associated inbox, if available. |
| `RequireShareRecipient` | `*bool` | `require_share_recipient` | no | Only allow access to recipients who have explicitly received the share via an email sent through the Files.com UI? |
| `SendOneTimePasswordToRecipientAtRegistration` | `*bool` | `send_one_time_password_to_recipient_at_registration` | no | If true, require_share_recipient bundles will send a one-time password to the recipient when they register. Cannot be enabled if the bundle has a password set. |
| `SendEmailReceiptToUploader` | `*bool` | `send_email_receipt_to_uploader` | no | Send delivery receipt to the uploader. Note: For writable share only |
| `SkipEmail` | `*bool` | `skip_email` | no | BundleRegistrations can be saved without providing email? |
| `SkipName` | `*bool` | `skip_name` | no | BundleRegistrations can be saved without providing name? |
| `SkipCompany` | `*bool` | `skip_company` | no | BundleRegistrations can be saved without providing company? |
| `StartAccessOnDate` | `*time.Time` | `start_access_on_date` | no | Date when share will start to be accessible. If `nil` access granted right after create. |
| `SnapshotId` | `int64` | `snapshot_id` | no | ID of the snapshot containing this bundle's contents. |
| `WorkspaceId` | `int64` | `workspace_id` | no | Workspace ID. `0` means the default workspace. |
| `WatermarkAttachmentFile` | `io.Writer` | `watermark_attachment_file` | no | Preview watermark image applied to all bundle items. |
| `WatermarkAttachmentFileAttachment` | `*lib.JSONAttachment` | `watermark_attachment_file` | no | Sends file content as watermark_attachment_file; see lib.JSONAttachment. Leave WatermarkAttachmentFile nil when you set it. |
| `WatermarkValue` | `interface{}` | `watermark_value` | no | Preview watermark settings applied to all bundle items. Uses the same keys as Behavior.value |

Example:

```go
result, err := bundle.Create(files_sdk.BundleCreateParams{
	Paths: []string{"file.txt"},
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Share

`POST /bundles/{id}/share`

Send email(s) with a link to bundle

```go
func (c *Client) Share(params files_sdk.BundleShareParams, opts ...files_sdk.RequestResponseOption) (err error)
func Share(params files_sdk.BundleShareParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.BundleShareParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Bundle ID. |
| `To` | `[]string` | `to` | no | A list of email addresses to share this bundle with. Required unless `recipients` is used. |
| `Note` | `string` | `note` | no | Note to include in email. |
| `Recipients` | `[]map[string]interface{}` | `recipients` | no | A list of recipients to share this bundle with. Required unless `to` is used. |

Example:

```go
err := bundle.Share(files_sdk.BundleShareParams{
	Id: 1,
})
if err != nil {
	return err
}
```

### Update

`PATCH /bundles/{id}`

Update Share Link

```go
func (c *Client) Update(params files_sdk.BundleUpdateParams, opts ...files_sdk.RequestResponseOption) (bundle files_sdk.Bundle, err error)
func Update(params files_sdk.BundleUpdateParams, opts ...files_sdk.RequestResponseOption) (bundle files_sdk.Bundle, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (bundle files_sdk.Bundle, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (bundle files_sdk.Bundle, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.BundleUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Bundle ID. |
| `Paths` | `[]string` | `paths` | no | A list of paths to include in this bundle. |
| `Password` | `string` | `password` | no | Password for this bundle. |
| `BypassesSiteExpirationRules` | `*bool` | `bypasses_site_expiration_rules` | no | If true, this Share Link bypasses site-wide expiration rules. Only site admins may set this. |
| `FormFieldSetId` | `int64` | `form_field_set_id` | no | Id of Form Field Set to use with this bundle |
| `ClickwrapId` | `int64` | `clickwrap_id` | no | ID of the clickwrap to use with this bundle. |
| `Code` | `string` | `code` | no | Bundle code. This code forms the end part of the Public URL. |
| `CreateSnapshot` | `*bool` | `create_snapshot` | no | If true, create a snapshot of this bundle's contents. |
| `Description` | `string` | `description` | no | Public description |
| `DontSeparateSubmissionsByFolder` | `*bool` | `dont_separate_submissions_by_folder` | no | Do not create subfolders for files uploaded to this share. Note: there are subtle security pitfalls with allowing anonymous uploads from multiple users to live in the same folder. We strongly discourage use of this option unless absolutely required. |
| `ExpiresAt` | `*time.Time` | `expires_at` | no | Explicit Bundle expiration date/time. If not set, the site-wide expiration setting may apply. |
| `FinalizeSnapshot` | `*bool` | `finalize_snapshot` | no | If true, finalize the snapshot of this bundle's contents. Note that `create_snapshot` must also be true. |
| `InboxId` | `int64` | `inbox_id` | no | ID of the associated inbox, if available. |
| `MaxUses` | `int64` | `max_uses` | no | Maximum number of times bundle can be accessed |
| `GroupId` | `int64` | `group_id` | no | Owning group ID. If set, members of this group can view, edit, and share this Share Link. |
| `InternalName` | `string` | `internal_name` | no | Internal name for identifying this Share Link. |
| `Note` | `string` | `note` | no | Bundle internal note |
| `PathTemplate` | `string` | `path_template` | no | Template for creating submission subfolders. Can use the uploader's name, email address, ip, company, `strftime` directives, and any custom form data. |
| `PathTemplateTimeZone` | `string` | `path_template_time_zone` | no | Timezone to use when rendering timestamps in path templates. |
| `Permissions` | `files_sdk.BundlePermissionsEnum` | `permissions` | no | Permissions that apply to Folders in this Share Link. Values: `read`, `write`, `read_write`, `full`, `none`, `preview_only`. |
| `RequireRegistration` | `*bool` | `require_registration` | no | Show a registration page that captures the downloader's name and email address? |
| `RequireShareRecipient` | `*bool` | `require_share_recipient` | no | Only allow access to recipients who have explicitly received the share via an email sent through the Files.com UI? |
| `SendOneTimePasswordToRecipientAtRegistration` | `*bool` | `send_one_time_password_to_recipient_at_registration` | no | If true, require_share_recipient bundles will send a one-time password to the recipient when they register. Cannot be enabled if the bundle has a password set. |
| `SendEmailReceiptToUploader` | `*bool` | `send_email_receipt_to_uploader` | no | Send delivery receipt to the uploader. Note: For writable share only |
| `SkipCompany` | `*bool` | `skip_company` | no | BundleRegistrations can be saved without providing company? |
| `StartAccessOnDate` | `*time.Time` | `start_access_on_date` | no | Date when share will start to be accessible. If `nil` access granted right after create. |
| `SkipEmail` | `*bool` | `skip_email` | no | BundleRegistrations can be saved without providing email? |
| `SkipName` | `*bool` | `skip_name` | no | BundleRegistrations can be saved without providing name? |
| `UserId` | `int64` | `user_id` | no | The owning user id. Only site admins can set this. |
| `WatermarkAttachmentDelete` | `*bool` | `watermark_attachment_delete` | no | If true, will delete the file stored in watermark_attachment |
| `WatermarkAttachmentFile` | `io.Writer` | `watermark_attachment_file` | no | Preview watermark image applied to all bundle items. |
| `WatermarkAttachmentFileAttachment` | `*lib.JSONAttachment` | `watermark_attachment_file` | no | Sends file content as watermark_attachment_file; see lib.JSONAttachment. Leave WatermarkAttachmentFile nil when you set it. |
| `WatermarkValue` | `interface{}` | `watermark_value` | no | Preview watermark settings applied to all bundle items. Uses the same keys as Behavior.value |
| `WorkspaceId` | `int64` | `workspace_id` | no | Workspace ID. `0` means the default workspace. |

Example:

```go
result, err := bundle.Update(files_sdk.BundleUpdateParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := bundle.UpdateWithMap(map[string]interface{}{
	"id": 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /bundles/{id}`

Delete Share Link

```go
func (c *Client) Delete(params files_sdk.BundleDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.BundleDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.BundleDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Bundle ID. |

Example:

```go
err := bundle.Delete(files_sdk.BundleDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.Bundle` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Code` | `string` | `code` | Bundle code. This code forms the end part of the Public URL. |
| `ColorLeft` | `string` | `color_left` | Page link and button color |
| `ColorLink` | `string` | `color_link` | Top bar link color |
| `ColorText` | `string` | `color_text` | Page link and button color |
| `ColorTop` | `string` | `color_top` | Top bar background color |
| `ColorTopText` | `string` | `color_top_text` | Top bar text color |
| `Url` | `string` | `url` | Public URL of Share Link |
| `Description` | `string` | `description` | Public description |
| `ExpiresAt` | `*time.Time` | `expires_at` | Explicit Bundle expiration date/time. If not set, the site-wide expiration setting may apply. |
| `PasswordProtected` | `*bool` | `password_protected` | Is this bundle password protected? |
| `Permissions` | `string` | `permissions` | Permissions that apply to Folders in this Share Link. |
| `PreviewOnly` | `*bool` | `preview_only` |  |
| `RequireRegistration` | `*bool` | `require_registration` | Show a registration page that captures the downloader's name and email address? |
| `RequireShareRecipient` | `*bool` | `require_share_recipient` | Only allow access to recipients who have explicitly received the share via an email sent through the Files.com UI? |
| `RequireLogout` | `*bool` | `require_logout` | If true, we will hide the 'Remember Me' box on the Bundle registration page, requiring that the user logout and log back in every time they visit the page. |
| `ClickwrapBody` | `string` | `clickwrap_body` | Legal text that must be agreed to prior to accessing Bundle. |
| `FormFieldSet` | `files_sdk.FormFieldSet` | `form_field_set` | Custom Form to use |
| `SkipName` | `*bool` | `skip_name` | BundleRegistrations can be saved without providing name? |
| `SkipEmail` | `*bool` | `skip_email` | BundleRegistrations can be saved without providing email? |
| `StartAccessOnDate` | `*time.Time` | `start_access_on_date` | Date when share will start to be accessible. If `nil` access granted right after create. |
| `SkipCompany` | `*bool` | `skip_company` | BundleRegistrations can be saved without providing company? |
| `Id` | `int64` | `id` | Bundle ID |
| `BypassesSiteExpirationRules` | `*bool` | `bypasses_site_expiration_rules` | If true, this Share Link bypasses site-wide expiration rules. Only site admins may set this. |
| `CreatedAt` | `*time.Time` | `created_at` | Bundle created at date/time |
| `Deleted` | `*bool` | `deleted` | Indicates if the bundle has been deleted. |
| `DeletedAt` | `*time.Time` | `deleted_at` | Bundle deleted at date/time |
| `DontSeparateSubmissionsByFolder` | `*bool` | `dont_separate_submissions_by_folder` | Do not create subfolders for files uploaded to this share. Note: there are subtle security pitfalls with allowing anonymous uploads from multiple users to live in the same folder. We strongly discourage use of this option unless absolutely required. |
| `EffectiveExpiresAt` | `*time.Time` | `effective_expires_at` | Read-only expiration date/time, using the explicit expiration or the site-wide setting when applicable. Null when the Share Link does not expire. |
| `MaxUses` | `int64` | `max_uses` | Maximum number of times bundle can be accessed |
| `InternalName` | `string` | `internal_name` | Internal name for identifying this Share Link. |
| `Note` | `string` | `note` | Bundle internal note |
| `PathTemplate` | `string` | `path_template` | Template for creating submission subfolders. Can use the uploader's name, email address, ip, company, `strftime` directives, and any custom form data. |
| `PathTemplateTimeZone` | `string` | `path_template_time_zone` | Timezone to use when rendering timestamps in path templates. |
| `SendEmailReceiptToUploader` | `*bool` | `send_email_receipt_to_uploader` | Send delivery receipt to the uploader. Note: For writable share only |
| `SnapshotId` | `int64` | `snapshot_id` | ID of the snapshot containing this bundle's contents. |
| `UserId` | `int64` | `user_id` | Bundle creator user ID |
| `Username` | `string` | `username` | Bundle creator username |
| `GroupId` | `int64` | `group_id` | Owning group ID. If set, members of this group can view, edit, and share this Share Link. |
| `ClickwrapId` | `int64` | `clickwrap_id` | ID of the clickwrap to use with this bundle. |
| `InboxId` | `int64` | `inbox_id` | ID of the associated inbox, if available. |
| `WatermarkAttachment` | `files_sdk.Image` | `watermark_attachment` | Preview watermark image applied to all bundle items. |
| `WatermarkValue` | `interface{}` | `watermark_value` | Preview watermark settings applied to all bundle items. Uses the same keys as Behavior.value |
| `SendOneTimePasswordToRecipientAtRegistration` | `*bool` | `send_one_time_password_to_recipient_at_registration` | If true, require_share_recipient bundles will send a one-time password to the recipient when they register. Cannot be enabled if the bundle has a password set. |
| `WorkspaceId` | `int64` | `workspace_id` | Workspace ID. `0` means the default workspace. |
| `HasInbox` | `*bool` | `has_inbox` | Does this bundle have an associated inbox? |
| `DontAllowFoldersInUploads` | `*bool` | `dont_allow_folders_in_uploads` | Should folder uploads be prevented? |
| `RequestedUploadSlots` | `[]map[string]interface{}` | `requested_upload_slots` | Upload slots requested by the associated Inbox. Each slot contains a name used as its label and destination subfolder name. |
| `Paths` | `[]string` | `paths` | A list of paths in this bundle. For performance reasons, this is not provided when listing bundles. |
| `Bundlepaths` | `[]files_sdk.BundlePath` | `bundlepaths` | A list of bundlepaths in this bundle. For performance reasons, this is not provided when listing bundles. |
| `Password` | `string` | `password` | Password for this bundle. |
| `FormFieldSetId` | `int64` | `form_field_set_id` | Id of Form Field Set to use with this bundle |
| `CreateSnapshot` | `*bool` | `create_snapshot` | If true, create a snapshot of this bundle's contents. |
| `FinalizeSnapshot` | `*bool` | `finalize_snapshot` | If true, finalize the snapshot of this bundle's contents. Note that `create_snapshot` must also be true. |
| `WatermarkAttachmentFile` | `io.Reader` | `watermark_attachment_file` | Preview watermark image applied to all bundle items. |
| `WatermarkAttachmentDelete` | `*bool` | `watermark_attachment_delete` | If true, will delete the file stored in watermark_attachment |
