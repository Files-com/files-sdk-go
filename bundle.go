package files_sdk

import (
	"encoding/json"
	"io"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// Bundle is a Files.com API resource.
type Bundle struct {
	// Bundle code. This code forms the end part of the Public URL.
	Code string `json:"code,omitempty" path:"code,omitempty" url:"code,omitempty"`
	// Page link and button color.
	ColorLeft string `json:"color_left,omitempty" path:"color_left,omitempty" url:"color_left,omitempty"`
	// Top bar link color.
	ColorLink string `json:"color_link,omitempty" path:"color_link,omitempty" url:"color_link,omitempty"`
	// Page link and button color.
	ColorText string `json:"color_text,omitempty" path:"color_text,omitempty" url:"color_text,omitempty"`
	// Top bar background color.
	ColorTop string `json:"color_top,omitempty" path:"color_top,omitempty" url:"color_top,omitempty"`
	// Top bar text color.
	ColorTopText string `json:"color_top_text,omitempty" path:"color_top_text,omitempty" url:"color_top_text,omitempty"`
	// Public URL of Share Link.
	Url string `json:"url,omitempty" path:"url,omitempty" url:"url,omitempty"`
	// Public description.
	Description string `json:"description,omitempty" path:"description,omitempty" url:"description,omitempty"`
	// Explicit Bundle expiration date/time. If not set, the site-wide expiration setting may apply.
	ExpiresAt *time.Time `json:"expires_at,omitempty" path:"expires_at,omitempty" url:"expires_at,omitempty"`
	// Is this bundle password protected?
	PasswordProtected *bool `json:"password_protected,omitempty" path:"password_protected,omitempty" url:"password_protected,omitempty"`
	// Permissions that apply to Folders in this Share Link.
	Permissions string `json:"permissions,omitempty" path:"permissions,omitempty" url:"permissions,omitempty"`
	PreviewOnly *bool  `json:"preview_only,omitempty" path:"preview_only,omitempty" url:"preview_only,omitempty"`
	// Show a registration page that captures the downloader's name and email address?
	RequireRegistration *bool `json:"require_registration,omitempty" path:"require_registration,omitempty" url:"require_registration,omitempty"`
	// Only allow access to recipients who have explicitly received the share via an email sent through the Files.com UI?
	RequireShareRecipient *bool `json:"require_share_recipient,omitempty" path:"require_share_recipient,omitempty" url:"require_share_recipient,omitempty"`
	// If true, we will hide the 'Remember Me' box on the Bundle registration page, requiring that the user logout and log back in every time they visit the page.
	RequireLogout *bool `json:"require_logout,omitempty" path:"require_logout,omitempty" url:"require_logout,omitempty"`
	// Legal text that must be agreed to prior to accessing Bundle.
	ClickwrapBody string `json:"clickwrap_body,omitempty" path:"clickwrap_body,omitempty" url:"clickwrap_body,omitempty"`
	// Custom Form to use.
	FormFieldSet FormFieldSet `json:"form_field_set,omitempty" path:"form_field_set,omitempty" url:"form_field_set,omitempty"`
	// BundleRegistrations can be saved without providing name?
	SkipName *bool `json:"skip_name,omitempty" path:"skip_name,omitempty" url:"skip_name,omitempty"`
	// BundleRegistrations can be saved without providing email?
	SkipEmail *bool `json:"skip_email,omitempty" path:"skip_email,omitempty" url:"skip_email,omitempty"`
	// Date when share will start to be accessible. If `nil` access granted right after create.
	StartAccessOnDate *time.Time `json:"start_access_on_date,omitempty" path:"start_access_on_date,omitempty" url:"start_access_on_date,omitempty"`
	// BundleRegistrations can be saved without providing company?
	SkipCompany *bool `json:"skip_company,omitempty" path:"skip_company,omitempty" url:"skip_company,omitempty"`
	// Bundle ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// If true, this Share Link bypasses site-wide expiration rules. Only site admins may set this.
	BypassesSiteExpirationRules *bool `json:"bypasses_site_expiration_rules,omitempty" path:"bypasses_site_expiration_rules,omitempty" url:"bypasses_site_expiration_rules,omitempty"`
	// Bundle created at date/time.
	CreatedAt *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	// Indicates if the bundle has been deleted.
	Deleted *bool `json:"deleted,omitempty" path:"deleted,omitempty" url:"deleted,omitempty"`
	// Bundle deleted at date/time.
	DeletedAt *time.Time `json:"deleted_at,omitempty" path:"deleted_at,omitempty" url:"deleted_at,omitempty"`
	// Do not create subfolders for files uploaded to this share. Note: there are subtle security pitfalls with allowing anonymous uploads from multiple users to live in the same folder. We strongly discourage use of this option unless absolutely required.
	DontSeparateSubmissionsByFolder *bool `json:"dont_separate_submissions_by_folder,omitempty" path:"dont_separate_submissions_by_folder,omitempty" url:"dont_separate_submissions_by_folder,omitempty"`
	// Read-only expiration date/time, using the explicit expiration or the site-wide setting when applicable. Null when the Share Link does not expire.
	EffectiveExpiresAt *time.Time `json:"effective_expires_at,omitempty" path:"effective_expires_at,omitempty" url:"effective_expires_at,omitempty"`
	// Maximum number of times bundle can be accessed.
	MaxUses int64 `json:"max_uses,omitempty" path:"max_uses,omitempty" url:"max_uses,omitempty"`
	// Internal name for identifying this Share Link.
	InternalName string `json:"internal_name,omitempty" path:"internal_name,omitempty" url:"internal_name,omitempty"`
	// Bundle internal note.
	Note string `json:"note,omitempty" path:"note,omitempty" url:"note,omitempty"`
	// Template for creating submission subfolders. Can use the uploader's name, email address, ip, company, `strftime` directives, and any custom form data.
	PathTemplate string `json:"path_template,omitempty" path:"path_template,omitempty" url:"path_template,omitempty"`
	// Timezone to use when rendering timestamps in path templates.
	PathTemplateTimeZone string `json:"path_template_time_zone,omitempty" path:"path_template_time_zone,omitempty" url:"path_template_time_zone,omitempty"`
	// Send delivery receipt to the uploader. Note: For writable share only.
	SendEmailReceiptToUploader *bool `json:"send_email_receipt_to_uploader,omitempty" path:"send_email_receipt_to_uploader,omitempty" url:"send_email_receipt_to_uploader,omitempty"`
	// ID of the snapshot containing this bundle's contents.
	SnapshotId int64 `json:"snapshot_id,omitempty" path:"snapshot_id,omitempty" url:"snapshot_id,omitempty"`
	// Bundle creator user ID.
	UserId int64 `json:"user_id,omitempty" path:"user_id,omitempty" url:"user_id,omitempty"`
	// Bundle creator username.
	Username string `json:"username,omitempty" path:"username,omitempty" url:"username,omitempty"`
	// Owning group ID. If set, members of this group can view, edit, and share this Share Link.
	GroupId int64 `json:"group_id,omitempty" path:"group_id,omitempty" url:"group_id,omitempty"`
	// ID of the clickwrap to use with this bundle.
	ClickwrapId int64 `json:"clickwrap_id,omitempty" path:"clickwrap_id,omitempty" url:"clickwrap_id,omitempty"`
	// ID of the associated inbox, if available.
	InboxId int64 `json:"inbox_id,omitempty" path:"inbox_id,omitempty" url:"inbox_id,omitempty"`
	// Preview watermark image applied to all bundle items.
	WatermarkAttachment Image `json:"watermark_attachment,omitempty" path:"watermark_attachment,omitempty" url:"watermark_attachment,omitempty"`
	// Preview watermark settings applied to all bundle items. Uses the same keys as Behavior.value.
	WatermarkValue interface{} `json:"watermark_value,omitempty" path:"watermark_value,omitempty" url:"watermark_value,omitempty"`
	// If true, require_share_recipient bundles will send a one-time password to the recipient when they register. Cannot be enabled if the bundle has a password set.
	SendOneTimePasswordToRecipientAtRegistration *bool `json:"send_one_time_password_to_recipient_at_registration,omitempty" path:"send_one_time_password_to_recipient_at_registration,omitempty" url:"send_one_time_password_to_recipient_at_registration,omitempty"`
	// Workspace ID. `0` means the default workspace.
	WorkspaceId int64 `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	// Does this bundle have an associated inbox?
	HasInbox *bool `json:"has_inbox,omitempty" path:"has_inbox,omitempty" url:"has_inbox,omitempty"`
	// Should folder uploads be prevented?
	DontAllowFoldersInUploads *bool `json:"dont_allow_folders_in_uploads,omitempty" path:"dont_allow_folders_in_uploads,omitempty" url:"dont_allow_folders_in_uploads,omitempty"`
	// Upload slots requested by the associated Inbox. Each slot contains a name used as its label and destination subfolder name.
	RequestedUploadSlots []map[string]interface{} `json:"requested_upload_slots,omitempty" path:"requested_upload_slots,omitempty" url:"requested_upload_slots,omitempty"`
	// A list of paths in this bundle. For performance reasons, this is not provided when listing bundles.
	Paths []string `json:"paths,omitempty" path:"paths,omitempty" url:"paths,omitempty"`
	// A list of bundlepaths in this bundle. For performance reasons, this is not provided when listing bundles.
	Bundlepaths []BundlePath `json:"bundlepaths,omitempty" path:"bundlepaths,omitempty" url:"bundlepaths,omitempty"`
	// Password for this bundle.
	Password string `json:"password,omitempty" path:"password,omitempty" url:"password,omitempty"`
	// Id of Form Field Set to use with this bundle.
	FormFieldSetId int64 `json:"form_field_set_id,omitempty" path:"form_field_set_id,omitempty" url:"form_field_set_id,omitempty"`
	// If true, create a snapshot of this bundle's contents.
	CreateSnapshot *bool `json:"create_snapshot,omitempty" path:"create_snapshot,omitempty" url:"create_snapshot,omitempty"`
	// If true, finalize the snapshot of this bundle's contents. Note that `create_snapshot` must also be true.
	FinalizeSnapshot *bool `json:"finalize_snapshot,omitempty" path:"finalize_snapshot,omitempty" url:"finalize_snapshot,omitempty"`
	// Preview watermark image applied to all bundle items.
	WatermarkAttachmentFile io.Reader `json:"watermark_attachment_file,omitempty" path:"watermark_attachment_file,omitempty" url:"watermark_attachment_file,omitempty"`
	// If true, will delete the file stored in watermark_attachment.
	WatermarkAttachmentDelete *bool `json:"watermark_attachment_delete,omitempty" path:"watermark_attachment_delete,omitempty" url:"watermark_attachment_delete,omitempty"`
}

// Identifier returns the resource ID.
func (b Bundle) Identifier() interface{} {
	return b.Id
}

// BundleCollection is a list of Bundle resources.
type BundleCollection []Bundle

// BundlePermissionsEnum is a string value for permissions.
// Enum lists the values documented by the API.
type BundlePermissionsEnum string

// String returns the API parameter value.
func (u BundlePermissionsEnum) String() string {
	return string(u)
}

// Enum returns the documented values keyed by their API strings.
func (u BundlePermissionsEnum) Enum() map[string]BundlePermissionsEnum {
	return map[string]BundlePermissionsEnum{
		"read":         BundlePermissionsEnum("read"),
		"write":        BundlePermissionsEnum("write"),
		"read_write":   BundlePermissionsEnum("read_write"),
		"full":         BundlePermissionsEnum("full"),
		"none":         BundlePermissionsEnum("none"),
		"preview_only": BundlePermissionsEnum("preview_only"),
	}
}

// BundleListParams contains the request parameters for GET /bundles.
type BundleListParams struct {
	// User ID. Provide a value of `0` to operate the current session's user.
	UserId int64 `url:"user_id,omitempty" json:"user_id,omitempty" path:"user_id"`
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `expires_at`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `created_at`, `expires_at`, `code`, `group_id`, `user_id` or `bypasses_site_expiration_rules`. Valid field combinations are `[ group_id, expires_at ]` and `[ user_id, expires_at ]`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// If set, return records where the specified field is greater than the supplied value. Valid fields are `created_at` and `expires_at`.
	FilterGt interface{} `url:"filter_gt,omitempty" json:"filter_gt,omitempty" path:"filter_gt"`
	// If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `created_at` and `expires_at`.
	FilterGteq interface{} `url:"filter_gteq,omitempty" json:"filter_gteq,omitempty" path:"filter_gteq"`
	// If set, return records where the specified field is prefixed by the supplied value. Valid fields are `code`.
	FilterPrefix interface{} `url:"filter_prefix,omitempty" json:"filter_prefix,omitempty" path:"filter_prefix"`
	// If set, return records where the specified field is less than the supplied value. Valid fields are `created_at` and `expires_at`.
	FilterLt interface{} `url:"filter_lt,omitempty" json:"filter_lt,omitempty" path:"filter_lt"`
	// If set, return records where the specified field is less than or equal the supplied value. Valid fields are `created_at` and `expires_at`.
	FilterLteq interface{} `url:"filter_lteq,omitempty" json:"filter_lteq,omitempty" path:"filter_lteq"`
	// If true, only list deleted Share Links.
	Deleted *bool `url:"deleted,omitempty" json:"deleted,omitempty" path:"deleted"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// BundleFindParams contains the request parameters for GET /bundles/{id}.
type BundleFindParams struct {
	// Bundle ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// If true, show a deleted Share Link.
	Deleted *bool `url:"deleted,omitempty" json:"deleted,omitempty" path:"deleted"`
}

// BundleCreateParams contains the request parameters for POST /bundles.
type BundleCreateParams struct {
	// User ID. Provide a value of `0` to operate the current session's user.
	UserId int64 `url:"user_id,omitempty" json:"user_id,omitempty" path:"user_id"`
	// A list of paths to include in this bundle. Required.
	Paths []string `url:"paths" json:"paths" path:"paths"`
	// Password for this bundle.
	Password string `url:"password,omitempty" json:"password,omitempty" path:"password"`
	// If true, this Share Link bypasses site-wide expiration rules. Only site admins may set this.
	BypassesSiteExpirationRules *bool `url:"bypasses_site_expiration_rules,omitempty" json:"bypasses_site_expiration_rules,omitempty" path:"bypasses_site_expiration_rules"`
	// Id of Form Field Set to use with this bundle.
	FormFieldSetId int64 `url:"form_field_set_id,omitempty" json:"form_field_set_id,omitempty" path:"form_field_set_id"`
	// If true, create a snapshot of this bundle's contents.
	CreateSnapshot *bool `url:"create_snapshot,omitempty" json:"create_snapshot,omitempty" path:"create_snapshot"`
	// Do not create subfolders for files uploaded to this share. Note: there are subtle security pitfalls with allowing anonymous uploads from multiple users to live in the same folder. We strongly discourage use of this option unless absolutely required.
	DontSeparateSubmissionsByFolder *bool `url:"dont_separate_submissions_by_folder,omitempty" json:"dont_separate_submissions_by_folder,omitempty" path:"dont_separate_submissions_by_folder"`
	// Explicit Bundle expiration date/time. If not set, the site-wide expiration setting may apply.
	ExpiresAt *time.Time `url:"expires_at,omitempty" json:"expires_at,omitempty" path:"expires_at"`
	// If true, finalize the snapshot of this bundle's contents. Note that `create_snapshot` must also be true.
	FinalizeSnapshot *bool `url:"finalize_snapshot,omitempty" json:"finalize_snapshot,omitempty" path:"finalize_snapshot"`
	// Maximum number of times bundle can be accessed.
	MaxUses int64 `url:"max_uses,omitempty" json:"max_uses,omitempty" path:"max_uses"`
	// Owning group ID. If set, members of this group can view, edit, and share this Share Link.
	GroupId int64 `url:"group_id,omitempty" json:"group_id,omitempty" path:"group_id"`
	// Internal name for identifying this Share Link.
	InternalName string `url:"internal_name,omitempty" json:"internal_name,omitempty" path:"internal_name"`
	// Public description.
	Description string `url:"description,omitempty" json:"description,omitempty" path:"description"`
	// Bundle internal note.
	Note string `url:"note,omitempty" json:"note,omitempty" path:"note"`
	// Bundle code. This code forms the end part of the Public URL.
	Code string `url:"code,omitempty" json:"code,omitempty" path:"code"`
	// Template for creating submission subfolders. Can use the uploader's name, email address, ip, company, `strftime` directives, and any custom form data.
	PathTemplate string `url:"path_template,omitempty" json:"path_template,omitempty" path:"path_template"`
	// Timezone to use when rendering timestamps in path templates.
	PathTemplateTimeZone string `url:"path_template_time_zone,omitempty" json:"path_template_time_zone,omitempty" path:"path_template_time_zone"`
	// Permissions that apply to Folders in this Share Link. Values: read, write, read_write, full, none, preview_only.
	Permissions BundlePermissionsEnum `url:"permissions,omitempty" json:"permissions,omitempty" path:"permissions"`
	// Show a registration page that captures the downloader's name and email address?
	RequireRegistration *bool `url:"require_registration,omitempty" json:"require_registration,omitempty" path:"require_registration"`
	// ID of the clickwrap to use with this bundle.
	ClickwrapId int64 `url:"clickwrap_id,omitempty" json:"clickwrap_id,omitempty" path:"clickwrap_id"`
	// ID of the associated inbox, if available.
	InboxId int64 `url:"inbox_id,omitempty" json:"inbox_id,omitempty" path:"inbox_id"`
	// Only allow access to recipients who have explicitly received the share via an email sent through the Files.com UI?
	RequireShareRecipient *bool `url:"require_share_recipient,omitempty" json:"require_share_recipient,omitempty" path:"require_share_recipient"`
	// If true, require_share_recipient bundles will send a one-time password to the recipient when they register. Cannot be enabled if the bundle has a password set.
	SendOneTimePasswordToRecipientAtRegistration *bool `url:"send_one_time_password_to_recipient_at_registration,omitempty" json:"send_one_time_password_to_recipient_at_registration,omitempty" path:"send_one_time_password_to_recipient_at_registration"`
	// Send delivery receipt to the uploader. Note: For writable share only.
	SendEmailReceiptToUploader *bool `url:"send_email_receipt_to_uploader,omitempty" json:"send_email_receipt_to_uploader,omitempty" path:"send_email_receipt_to_uploader"`
	// BundleRegistrations can be saved without providing email?
	SkipEmail *bool `url:"skip_email,omitempty" json:"skip_email,omitempty" path:"skip_email"`
	// BundleRegistrations can be saved without providing name?
	SkipName *bool `url:"skip_name,omitempty" json:"skip_name,omitempty" path:"skip_name"`
	// BundleRegistrations can be saved without providing company?
	SkipCompany *bool `url:"skip_company,omitempty" json:"skip_company,omitempty" path:"skip_company"`
	// Date when share will start to be accessible. If `nil` access granted right after create.
	StartAccessOnDate *time.Time `url:"start_access_on_date,omitempty" json:"start_access_on_date,omitempty" path:"start_access_on_date"`
	// ID of the snapshot containing this bundle's contents.
	SnapshotId int64 `url:"snapshot_id,omitempty" json:"snapshot_id,omitempty" path:"snapshot_id"`
	// Workspace ID. `0` means the default workspace.
	WorkspaceId int64 `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
	// Preview watermark image applied to all bundle items.
	WatermarkAttachmentFile io.Writer `url:"watermark_attachment_file,omitempty" json:"watermark_attachment_file,omitempty" path:"watermark_attachment_file"`
	// WatermarkAttachmentFileAttachment sends file content as watermark_attachment_file; see lib.JSONAttachment. Leave WatermarkAttachmentFile nil when you set it.
	WatermarkAttachmentFileAttachment *lib.JSONAttachment `json:"-" url:"-" path:"-" attachment:"watermark_attachment_file"`
	// Preview watermark settings applied to all bundle items. Uses the same keys as Behavior.value.
	WatermarkValue interface{} `url:"watermark_value,omitempty" json:"watermark_value,omitempty" path:"watermark_value"`
}

// BundleShareParams contains the request parameters for POST /bundles/{id}/share.
//
// Send email(s) with a link to bundle
type BundleShareParams struct {
	// Bundle ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// A list of email addresses to share this bundle with. Required unless `recipients` is used.
	To []string `url:"to,omitempty" json:"to,omitempty" path:"to"`
	// Note to include in email.
	Note string `url:"note,omitempty" json:"note,omitempty" path:"note"`
	// A list of recipients to share this bundle with. Required unless `to` is used.
	Recipients []map[string]interface{} `url:"recipients,omitempty" json:"recipients,omitempty" path:"recipients"`
}

// BundleUpdateParams contains the request parameters for PATCH /bundles/{id}.
type BundleUpdateParams struct {
	// Bundle ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// A list of paths to include in this bundle.
	Paths []string `url:"paths,omitempty" json:"paths,omitempty" path:"paths"`
	// Password for this bundle.
	Password string `url:"password,omitempty" json:"password,omitempty" path:"password"`
	// If true, this Share Link bypasses site-wide expiration rules. Only site admins may set this.
	BypassesSiteExpirationRules *bool `url:"bypasses_site_expiration_rules,omitempty" json:"bypasses_site_expiration_rules,omitempty" path:"bypasses_site_expiration_rules"`
	// Id of Form Field Set to use with this bundle.
	FormFieldSetId int64 `url:"form_field_set_id,omitempty" json:"form_field_set_id,omitempty" path:"form_field_set_id"`
	// ID of the clickwrap to use with this bundle.
	ClickwrapId int64 `url:"clickwrap_id,omitempty" json:"clickwrap_id,omitempty" path:"clickwrap_id"`
	// Bundle code. This code forms the end part of the Public URL.
	Code string `url:"code,omitempty" json:"code,omitempty" path:"code"`
	// If true, create a snapshot of this bundle's contents.
	CreateSnapshot *bool `url:"create_snapshot,omitempty" json:"create_snapshot,omitempty" path:"create_snapshot"`
	// Public description.
	Description string `url:"description,omitempty" json:"description,omitempty" path:"description"`
	// Do not create subfolders for files uploaded to this share. Note: there are subtle security pitfalls with allowing anonymous uploads from multiple users to live in the same folder. We strongly discourage use of this option unless absolutely required.
	DontSeparateSubmissionsByFolder *bool `url:"dont_separate_submissions_by_folder,omitempty" json:"dont_separate_submissions_by_folder,omitempty" path:"dont_separate_submissions_by_folder"`
	// Explicit Bundle expiration date/time. If not set, the site-wide expiration setting may apply.
	ExpiresAt *time.Time `url:"expires_at,omitempty" json:"expires_at,omitempty" path:"expires_at"`
	// If true, finalize the snapshot of this bundle's contents. Note that `create_snapshot` must also be true.
	FinalizeSnapshot *bool `url:"finalize_snapshot,omitempty" json:"finalize_snapshot,omitempty" path:"finalize_snapshot"`
	// ID of the associated inbox, if available.
	InboxId int64 `url:"inbox_id,omitempty" json:"inbox_id,omitempty" path:"inbox_id"`
	// Maximum number of times bundle can be accessed.
	MaxUses int64 `url:"max_uses,omitempty" json:"max_uses,omitempty" path:"max_uses"`
	// Owning group ID. If set, members of this group can view, edit, and share this Share Link.
	GroupId int64 `url:"group_id,omitempty" json:"group_id,omitempty" path:"group_id"`
	// Internal name for identifying this Share Link.
	InternalName string `url:"internal_name,omitempty" json:"internal_name,omitempty" path:"internal_name"`
	// Bundle internal note.
	Note string `url:"note,omitempty" json:"note,omitempty" path:"note"`
	// Template for creating submission subfolders. Can use the uploader's name, email address, ip, company, `strftime` directives, and any custom form data.
	PathTemplate string `url:"path_template,omitempty" json:"path_template,omitempty" path:"path_template"`
	// Timezone to use when rendering timestamps in path templates.
	PathTemplateTimeZone string `url:"path_template_time_zone,omitempty" json:"path_template_time_zone,omitempty" path:"path_template_time_zone"`
	// Permissions that apply to Folders in this Share Link. Values: read, write, read_write, full, none, preview_only.
	Permissions BundlePermissionsEnum `url:"permissions,omitempty" json:"permissions,omitempty" path:"permissions"`
	// Show a registration page that captures the downloader's name and email address?
	RequireRegistration *bool `url:"require_registration,omitempty" json:"require_registration,omitempty" path:"require_registration"`
	// Only allow access to recipients who have explicitly received the share via an email sent through the Files.com UI?
	RequireShareRecipient *bool `url:"require_share_recipient,omitempty" json:"require_share_recipient,omitempty" path:"require_share_recipient"`
	// If true, require_share_recipient bundles will send a one-time password to the recipient when they register. Cannot be enabled if the bundle has a password set.
	SendOneTimePasswordToRecipientAtRegistration *bool `url:"send_one_time_password_to_recipient_at_registration,omitempty" json:"send_one_time_password_to_recipient_at_registration,omitempty" path:"send_one_time_password_to_recipient_at_registration"`
	// Send delivery receipt to the uploader. Note: For writable share only.
	SendEmailReceiptToUploader *bool `url:"send_email_receipt_to_uploader,omitempty" json:"send_email_receipt_to_uploader,omitempty" path:"send_email_receipt_to_uploader"`
	// BundleRegistrations can be saved without providing company?
	SkipCompany *bool `url:"skip_company,omitempty" json:"skip_company,omitempty" path:"skip_company"`
	// Date when share will start to be accessible. If `nil` access granted right after create.
	StartAccessOnDate *time.Time `url:"start_access_on_date,omitempty" json:"start_access_on_date,omitempty" path:"start_access_on_date"`
	// BundleRegistrations can be saved without providing email?
	SkipEmail *bool `url:"skip_email,omitempty" json:"skip_email,omitempty" path:"skip_email"`
	// BundleRegistrations can be saved without providing name?
	SkipName *bool `url:"skip_name,omitempty" json:"skip_name,omitempty" path:"skip_name"`
	// The owning user id. Only site admins can set this.
	UserId int64 `url:"user_id,omitempty" json:"user_id,omitempty" path:"user_id"`
	// If true, will delete the file stored in watermark_attachment.
	WatermarkAttachmentDelete *bool `url:"watermark_attachment_delete,omitempty" json:"watermark_attachment_delete,omitempty" path:"watermark_attachment_delete"`
	// Preview watermark image applied to all bundle items.
	WatermarkAttachmentFile io.Writer `url:"watermark_attachment_file,omitempty" json:"watermark_attachment_file,omitempty" path:"watermark_attachment_file"`
	// WatermarkAttachmentFileAttachment sends file content as watermark_attachment_file; see lib.JSONAttachment. Leave WatermarkAttachmentFile nil when you set it.
	WatermarkAttachmentFileAttachment *lib.JSONAttachment `json:"-" url:"-" path:"-" attachment:"watermark_attachment_file"`
	// Preview watermark settings applied to all bundle items. Uses the same keys as Behavior.value.
	WatermarkValue interface{} `url:"watermark_value,omitempty" json:"watermark_value,omitempty" path:"watermark_value"`
	// Workspace ID. `0` means the default workspace.
	WorkspaceId int64 `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
}

// BundleDeleteParams contains the request parameters for DELETE /bundles/{id}.
type BundleDeleteParams struct {
	// Bundle ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (b *Bundle) UnmarshalJSON(data []byte) error {
	type bundle Bundle
	var v bundle
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*b = Bundle(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (b *BundleCollection) UnmarshalJSON(data []byte) error {
	type bundles BundleCollection
	var v bundles
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*b = BundleCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (b *BundleCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*b))
	for i, v := range *b {
		ret[i] = v
	}

	return &ret
}
