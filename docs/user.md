# User

`files_sdk.User` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/user`, named `user`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/user"
)
```

## Operations

### List

`GET /users`

List Users

```go
func (c *Client) List(params files_sdk.UserListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.UserListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.User` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.UserListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `site_id`, `workspace_id`, `company`, `name`, `disabled`, `authenticate_until`, `username`, `email`, `site_admin`, `last_desktop_login_at`, `last_login_at`, `password_validity_days` or `ssl_required`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `username`, `name`, `email`, `company`, `site_admin`, `password_validity_days`, `ssl_required`, `last_login_at`, `authenticate_until`, `not_site_admin`, `disabled`, `partner_id`, `primary_group_id` or `workspace_id`. Valid field combinations are `[ site_admin, username ]`, `[ not_site_admin, username ]`, `[ workspace_id, username ]`, `[ company, name ]`, `[ workspace_id, name ]`, `[ workspace_id, email ]`, `[ workspace_id, company ]`, `[ workspace_id, site_admin ]`, `[ workspace_id, not_site_admin ]`, `[ workspace_id, disabled ]`, `[ workspace_id, partner_id ]`, `[ workspace_id, site_admin, username ]`, `[ workspace_id, not_site_admin, username ]`, `[ workspace_id, disabled, username ]`, `[ workspace_id, partner_id, username ]` or `[ workspace_id, company, name ]`. |
| `FilterGt` | `interface{}` | `filter_gt` | no | If set, return records where the specified field is greater than the supplied value. Valid fields are `password_validity_days`, `last_login_at` or `authenticate_until`. |
| `FilterGteq` | `interface{}` | `filter_gteq` | no | If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `password_validity_days`, `last_login_at` or `authenticate_until`. |
| `FilterPrefix` | `interface{}` | `filter_prefix` | no | If set, return records where the specified field is prefixed by the supplied value. Valid fields are `username`, `name`, `email` or `company`. Valid field combinations are `[ company, name ]`. |
| `FilterLt` | `interface{}` | `filter_lt` | no | If set, return records where the specified field is less than the supplied value. Valid fields are `password_validity_days`, `last_login_at` or `authenticate_until`. |
| `FilterLteq` | `interface{}` | `filter_lteq` | no | If set, return records where the specified field is less than or equal the supplied value. Valid fields are `password_validity_days`, `last_login_at` or `authenticate_until`. |
| `Ids` | `string` | `ids` | no | comma-separated list of User IDs |
| `IncludeParentSiteUsers` | `*bool` | `include_parent_site_users` | no | Include users from the parent site. |
| `IncludeDefaultWorkspaceUsers` | `*bool` | `include_default_workspace_users` | no | When scoped to a Custom Workspace, include Default Workspace users who have permissions in that Workspace. Available to Site Administrators, and to Workspace Administrators and Folder Admins in that Workspace when the site shares its Default Workspace directory with Workspaces. |
| `Search` | `string` | `search` | no | Searches for partial matches of name, username, or email. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := user.List(files_sdk.UserListParams{})
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

`GET /users/{id}`

Show User

```go
func (c *Client) Find(params files_sdk.UserFindParams, opts ...files_sdk.RequestResponseOption) (user files_sdk.User, err error)
func Find(params files_sdk.UserFindParams, opts ...files_sdk.RequestResponseOption) (user files_sdk.User, err error)
```

Parameters, `files_sdk.UserFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | User ID. |

Example:

```go
result, err := user.Find(files_sdk.UserFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Create

`POST /users`

Create User

```go
func (c *Client) Create(params files_sdk.UserCreateParams, opts ...files_sdk.RequestResponseOption) (user files_sdk.User, err error)
func Create(params files_sdk.UserCreateParams, opts ...files_sdk.RequestResponseOption) (user files_sdk.User, err error)
```

Parameters, `files_sdk.UserCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `AvatarFile` | `io.Writer` | `avatar_file` | no | An image file for your user avatar. |
| `AvatarFileAttachment` | `*lib.JSONAttachment` | `avatar_file` | no | Sends file content as avatar_file; see lib.JSONAttachment. Leave AvatarFile nil when you set it. |
| `AvatarDelete` | `*bool` | `avatar_delete` | no | If true, the avatar will be deleted. |
| `ChangePassword` | `string` | `change_password` | no | Used for changing a password on an existing user. |
| `ChangePasswordConfirmation` | `string` | `change_password_confirmation` | no | Optional, but if provided, we will ensure that it matches the value sent in `change_password`. |
| `Email` | `string` | `email` | no | User's email. |
| `GrantPermission` | `string` | `grant_permission` | no | Permission to grant on the User Root upon creation, or on the Workspace root when the User Root is blank. For the Default Workspace, a blank User Root uses the site root. `admin` on a Custom Workspace's root grants Workspace Administrator access; on a subfolder it grants Folder Admin access. Can be blank or `admin`, `full`, `read`, `write`, `list`, `read+write`, `list+write`, or `history`. |
| `GroupId` | `int64` | `group_id` | no | Group ID to associate this user with. |
| `GroupIds` | `string` | `group_ids` | no | A list of group ids to associate this user with. Comma delimited. |
| `ImportedPasswordHash` | `string` | `imported_password_hash` | no | Pre-calculated hash of the user's password. If supplied, this will be used to authenticate the user on first login. Supported hash methods are MD5, SHA1, and SHA256. |
| `Password` | `string` | `password` | no | User password. |
| `PasswordConfirmation` | `string` | `password_confirmation` | no | Optional, but if provided, we will ensure that it matches the value sent in `password`. |
| `AnnouncementsRead` | `*bool` | `announcements_read` | no | Signifies that the user has read all the announcements in the UI. |
| `AiAssistantPersonalityId` | `int64` | `ai_assistant_personality_id` | no | AI Assistant Personality ID assigned directly to this user, if any. |
| `AllowedIps` | `string` | `allowed_ips` | no | A list of allowed IPs if applicable. Newline delimited |
| `AttachmentsPermission` | `*bool` | `attachments_permission` | no | DEPRECATED: If `true`, the user can user create Bundles (aka Share Links). Use the bundle permission instead. |
| `AuthenticateUntil` | `*time.Time` | `authenticate_until` | no | Scheduled Date/Time at which user will be deactivated |
| `AuthenticationMethod` | `files_sdk.UserAuthenticationMethodEnum` | `authentication_method` | no | How is this user authenticated? Values: `password`, `sso`, `none`, `email_signup`, `password_with_imported_hash`, `password_and_ssh_key`. |
| `BillingPermission` | `*bool` | `billing_permission` | no | Allow this user to perform operations on the account, payments, and invoices? |
| `BypassUserLifecycleRules` | `*bool` | `bypass_user_lifecycle_rules` | no | Exempt this user from user lifecycle rules? |
| `BypassSiteAllowedIps` | `*bool` | `bypass_site_allowed_ips` | no | Allow this user to skip site-wide IP blacklists? |
| `DavPermission` | `*bool` | `dav_permission` | no | Can the user connect with WebDAV? |
| `DesktopConfigurationProfileId` | `int64` | `desktop_configuration_profile_id` | no | Desktop Configuration Profile ID assigned directly to this user, if any. |
| `DefaultWorkspaceId` | `int64` | `default_workspace_id` | no | Workspace ID the user should land in by default when more than one Workspace is available. |
| `Disabled` | `*bool` | `disabled` | no | Is user disabled? Disabled users cannot log in, and do not count for billing purposes. Users can be automatically disabled after an inactivity period via a Site setting or schedule to be deactivated after specific date. |
| `FilesystemLayout` | `files_sdk.UserFilesystemLayoutEnum` | `filesystem_layout` | no | File system layout Values: `site_root`, `user_root`, `partner_root`, `integration_centric`, `workspace_root`. |
| `FtpPermission` | `*bool` | `ftp_permission` | no | Can the user access with FTP/FTPS? |
| `HeaderText` | `string` | `header_text` | no | Text to display to the user in the header of the UI |
| `IntegrationCentricProfileId` | `int64` | `integration_centric_profile_id` | no | Integration Centric Profile ID assigned directly to this user, if any. |
| `Language` | `string` | `language` | no | Preferred language |
| `NotificationDailySendTime` | `int64` | `notification_daily_send_time` | no | Hour of the day at which daily notifications should be sent. Can be in range 0 to 23 |
| `Name` | `string` | `name` | no | User's full name |
| `Company` | `string` | `company` | no | User's company |
| `Notes` | `string` | `notes` | no | Any internal notes on the user |
| `OfficeIntegrationEnabled` | `*bool` | `office_integration_enabled` | no | Enable integration with Office for the web? |
| `PartnerAdmin` | `*bool` | `partner_admin` | no | Is this user a Partner administrator? |
| `PartnerId` | `int64` | `partner_id` | no | Partner ID if this user belongs to a Partner |
| `PasswordValidityDays` | `int64` | `password_validity_days` | no | Number of days to allow user to use the same password |
| `PrimaryGroupId` | `int64` | `primary_group_id` | no | Primary group ID for Group Admin scoping |
| `ReadonlySiteAdmin` | `*bool` | `readonly_site_admin` | no | Is the user an allowed to view all (non-billing) site configuration for this site? |
| `ReceiveAdminAlerts` | `*bool` | `receive_admin_alerts` | no | Deprecated. Use notify_on_all_site_warnings and granular failure notification preferences instead. |
| `NotifyOnAllSiteWarnings` | `*bool` | `notify_on_all_site_warnings` | no | Should the user receive site warnings via email? |
| `NotifyOnAllSsoFailures` | `*bool` | `notify_on_all_sso_failures` | no | Should the user receive sso/scim/ldap configuration/sync failures via email? |
| `NotifyOnAllUserSecurityEvents` | `*bool` | `notify_on_all_user_security_events` | no | Should the user receive user security events via email? |
| `NotifyOnAllPendingWorkFailures` | `*bool` | `notify_on_all_pending_work_failures` | no | Should the user receive pending work failures via email? |
| `NotifyOnAllSiemHttpDestinationFailures` | `*bool` | `notify_on_all_siem_http_destination_failures` | no | Should the user receive siem failures via email? |
| `NotifyOnAllSyncFailures` | `*bool` | `notify_on_all_sync_failures` | no | Should the user receive sync failures via email? |
| `NotifyOnAllAutomationFailures` | `*bool` | `notify_on_all_automation_failures` | no | Should the user receive automation failures via email? |
| `NotifyOnAllExpectationFailures` | `*bool` | `notify_on_all_expectation_failures` | no | Should the user receive expectation failures and misses via email? |
| `RequireLoginBy` | `*time.Time` | `require_login_by` | no | Require user to login by specified date otherwise it will be disabled. |
| `RequirePasswordChange` | `*bool` | `require_password_change` | no | Is a password change required upon next user login? |
| `ResponsibleGroupId` | `int64` | `responsible_group_id` | no | ID of the internal Group responsible for this Partner User, overriding the Partner default. |
| `ResponsibleUserId` | `int64` | `responsible_user_id` | no | ID of the internal User responsible for this Partner User, overriding the Partner default. |
| `RestapiPermission` | `*bool` | `restapi_permission` | no | Can this user access the Web app, Desktop app, SDKs, or REST API? (All of these tools use the API internally, so this is one unified permission set.) |
| `S3CompatibleEndpointPermission` | `*bool` | `s3_compatible_endpoint_permission` | no | Can the user access the S3-compatible endpoint? Defaults to true. |
| `SelfManaged` | `*bool` | `self_managed` | no | Does this user manage it's own credentials or is it a shared/bot user? |
| `SftpPermission` | `*bool` | `sftp_permission` | no | Can the user access with SFTP? |
| `SiteAdmin` | `*bool` | `site_admin` | no | Is the user an administrator for this site? |
| `SkipWelcomeScreen` | `*bool` | `skip_welcome_screen` | no | Skip Welcome page in the UI? |
| `SslRequired` | `files_sdk.UserSslRequiredEnum` | `ssl_required` | no | SSL required setting Values: `use_system_setting`, `always_require`, `never_require`. |
| `SsoStrategyId` | `int64` | `sso_strategy_id` | no | SSO (Single Sign On) strategy ID for the user, if applicable. |
| `Require2fa` | `files_sdk.UserRequire2faEnum` | `require_2fa` | no | 2FA required setting. `use_system_setting` uses the site-wide setting, including SSO exemptions. `always_require` and `never_require` override the site-wide setting when user-level overrides are allowed. Values: `use_system_setting`, `always_require`, `never_require`. |
| `Tags` | `string` | `tags` | no | Comma-separated list of Tags for this user. Tags are used for other features, such as UserLifecycleRules, which can target specific tags. Tags must only contain lowercase letters, numbers, and hyphens. |
| `TimeZone` | `string` | `time_zone` | no | User time zone |
| `UserRoot` | `string` | `user_root` | no | If filesystem layout is user_root, this path is the root path the user is fixed to for all interfaces. If the filesystem layout is site_root or partner_root, this acts as a root folder only for FTP and SFTP (SFTP applicability also requires a site-wide setting to be set). For partner_root layout, this path is relative to the Partner root folder for all callers and blank opts out of an additional protocol root. In this situation, this path is not applied to the API, Desktop, or Web interface. |
| `UserHome` | `string` | `user_home` | no | Home folder for FTP/SFTP. For users with the partner_root filesystem layout, this path is relative to the Partner root folder. In all other cases, it is an absolute path. Only applies to FTP and SFTP, and not any other interface. |
| `WorkspaceAdmin` | `*bool` | `workspace_admin` | no | Whether the user is an administrator of their own Custom Workspace. Does not reflect administration granted through Permissions. |
| `Username` | `string` | `username` | yes | User's username |
| `WorkspaceId` | `int64` | `workspace_id` | no | ID of the Workspace the user belongs to. 0 is the Default Workspace. |

Example:

```go
result, err := user.Create(files_sdk.UserCreateParams{
	Username: "user",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Unlock

`POST /users/{id}/unlock`

Unlock user who has been locked out due to failed logins

```go
func (c *Client) Unlock(params files_sdk.UserUnlockParams, opts ...files_sdk.RequestResponseOption) (err error)
func Unlock(params files_sdk.UserUnlockParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.UserUnlockParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | User ID. |

Example:

```go
err := user.Unlock(files_sdk.UserUnlockParams{
	Id: 1,
})
if err != nil {
	return err
}
```

### ResendWelcomeEmail

`POST /users/{id}/resend_welcome_email`

Resend user welcome email

```go
func (c *Client) ResendWelcomeEmail(params files_sdk.UserResendWelcomeEmailParams, opts ...files_sdk.RequestResponseOption) (err error)
func ResendWelcomeEmail(params files_sdk.UserResendWelcomeEmailParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.UserResendWelcomeEmailParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | User ID. |

Example:

```go
err := user.ResendWelcomeEmail(files_sdk.UserResendWelcomeEmailParams{
	Id: 1,
})
if err != nil {
	return err
}
```

### User2faReset

`POST /users/{id}/2fa/reset`

Trigger 2FA Reset process for user who has lost access to their existing 2FA methods

```go
func (c *Client) User2faReset(params files_sdk.UserUser2faResetParams, opts ...files_sdk.RequestResponseOption) (err error)
func User2faReset(params files_sdk.UserUser2faResetParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.UserUser2faResetParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | User ID. |

Example:

```go
err := user.User2faReset(files_sdk.UserUser2faResetParams{
	Id: 1,
})
if err != nil {
	return err
}
```

### Update

`PATCH /users/{id}`

Update User

```go
func (c *Client) Update(params files_sdk.UserUpdateParams, opts ...files_sdk.RequestResponseOption) (user files_sdk.User, err error)
func Update(params files_sdk.UserUpdateParams, opts ...files_sdk.RequestResponseOption) (user files_sdk.User, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (user files_sdk.User, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (user files_sdk.User, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.UserUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | User ID. |
| `AvatarFile` | `io.Writer` | `avatar_file` | no | An image file for your user avatar. |
| `AvatarFileAttachment` | `*lib.JSONAttachment` | `avatar_file` | no | Sends file content as avatar_file; see lib.JSONAttachment. Leave AvatarFile nil when you set it. |
| `AvatarDelete` | `*bool` | `avatar_delete` | no | If true, the avatar will be deleted. |
| `ChangePassword` | `string` | `change_password` | no | Used for changing a password on an existing user. |
| `ChangePasswordConfirmation` | `string` | `change_password_confirmation` | no | Optional, but if provided, we will ensure that it matches the value sent in `change_password`. |
| `Email` | `string` | `email` | no | User's email. |
| `GrantPermission` | `string` | `grant_permission` | no | Permission to grant on the User Root upon creation, or on the Workspace root when the User Root is blank. For the Default Workspace, a blank User Root uses the site root. `admin` on a Custom Workspace's root grants Workspace Administrator access; on a subfolder it grants Folder Admin access. Can be blank or `admin`, `full`, `read`, `write`, `list`, `read+write`, `list+write`, or `history`. |
| `GroupId` | `int64` | `group_id` | no | Group ID to associate this user with. |
| `GroupIds` | `string` | `group_ids` | no | A list of group ids to associate this user with. Comma delimited. |
| `ImportedPasswordHash` | `string` | `imported_password_hash` | no | Pre-calculated hash of the user's password. If supplied, this will be used to authenticate the user on first login. Supported hash methods are MD5, SHA1, and SHA256. |
| `Password` | `string` | `password` | no | User password. |
| `PasswordConfirmation` | `string` | `password_confirmation` | no | Optional, but if provided, we will ensure that it matches the value sent in `password`. |
| `AnnouncementsRead` | `*bool` | `announcements_read` | no | Signifies that the user has read all the announcements in the UI. |
| `AiAssistantPersonalityId` | `int64` | `ai_assistant_personality_id` | no | AI Assistant Personality ID assigned directly to this user, if any. |
| `AllowedIps` | `string` | `allowed_ips` | no | A list of allowed IPs if applicable. Newline delimited |
| `AttachmentsPermission` | `*bool` | `attachments_permission` | no | DEPRECATED: If `true`, the user can user create Bundles (aka Share Links). Use the bundle permission instead. |
| `AuthenticateUntil` | `*time.Time` | `authenticate_until` | no | Scheduled Date/Time at which user will be deactivated |
| `AuthenticationMethod` | `files_sdk.UserAuthenticationMethodEnum` | `authentication_method` | no | How is this user authenticated? Values: `password`, `sso`, `none`, `email_signup`, `password_with_imported_hash`, `password_and_ssh_key`. |
| `BillingPermission` | `*bool` | `billing_permission` | no | Allow this user to perform operations on the account, payments, and invoices? |
| `BypassUserLifecycleRules` | `*bool` | `bypass_user_lifecycle_rules` | no | Exempt this user from user lifecycle rules? |
| `BypassSiteAllowedIps` | `*bool` | `bypass_site_allowed_ips` | no | Allow this user to skip site-wide IP blacklists? |
| `DavPermission` | `*bool` | `dav_permission` | no | Can the user connect with WebDAV? |
| `DesktopConfigurationProfileId` | `int64` | `desktop_configuration_profile_id` | no | Desktop Configuration Profile ID assigned directly to this user, if any. |
| `DefaultWorkspaceId` | `int64` | `default_workspace_id` | no | Workspace ID the user should land in by default when more than one Workspace is available. |
| `Disabled` | `*bool` | `disabled` | no | Is user disabled? Disabled users cannot log in, and do not count for billing purposes. Users can be automatically disabled after an inactivity period via a Site setting or schedule to be deactivated after specific date. |
| `FilesystemLayout` | `files_sdk.UserFilesystemLayoutEnum` | `filesystem_layout` | no | File system layout Values: `site_root`, `user_root`, `partner_root`, `integration_centric`, `workspace_root`. |
| `FtpPermission` | `*bool` | `ftp_permission` | no | Can the user access with FTP/FTPS? |
| `HeaderText` | `string` | `header_text` | no | Text to display to the user in the header of the UI |
| `IntegrationCentricProfileId` | `int64` | `integration_centric_profile_id` | no | Integration Centric Profile ID assigned directly to this user, if any. |
| `Language` | `string` | `language` | no | Preferred language |
| `NotificationDailySendTime` | `int64` | `notification_daily_send_time` | no | Hour of the day at which daily notifications should be sent. Can be in range 0 to 23 |
| `Name` | `string` | `name` | no | User's full name |
| `Company` | `string` | `company` | no | User's company |
| `Notes` | `string` | `notes` | no | Any internal notes on the user |
| `OfficeIntegrationEnabled` | `*bool` | `office_integration_enabled` | no | Enable integration with Office for the web? |
| `PartnerAdmin` | `*bool` | `partner_admin` | no | Is this user a Partner administrator? |
| `PartnerId` | `int64` | `partner_id` | no | Partner ID if this user belongs to a Partner |
| `PasswordValidityDays` | `int64` | `password_validity_days` | no | Number of days to allow user to use the same password |
| `PrimaryGroupId` | `int64` | `primary_group_id` | no | Primary group ID for Group Admin scoping |
| `ReadonlySiteAdmin` | `*bool` | `readonly_site_admin` | no | Is the user an allowed to view all (non-billing) site configuration for this site? |
| `ReceiveAdminAlerts` | `*bool` | `receive_admin_alerts` | no | Deprecated. Use notify_on_all_site_warnings and granular failure notification preferences instead. |
| `NotifyOnAllSiteWarnings` | `*bool` | `notify_on_all_site_warnings` | no | Should the user receive site warnings via email? |
| `NotifyOnAllSsoFailures` | `*bool` | `notify_on_all_sso_failures` | no | Should the user receive sso/scim/ldap configuration/sync failures via email? |
| `NotifyOnAllUserSecurityEvents` | `*bool` | `notify_on_all_user_security_events` | no | Should the user receive user security events via email? |
| `NotifyOnAllPendingWorkFailures` | `*bool` | `notify_on_all_pending_work_failures` | no | Should the user receive pending work failures via email? |
| `NotifyOnAllSiemHttpDestinationFailures` | `*bool` | `notify_on_all_siem_http_destination_failures` | no | Should the user receive siem failures via email? |
| `NotifyOnAllSyncFailures` | `*bool` | `notify_on_all_sync_failures` | no | Should the user receive sync failures via email? |
| `NotifyOnAllAutomationFailures` | `*bool` | `notify_on_all_automation_failures` | no | Should the user receive automation failures via email? |
| `NotifyOnAllExpectationFailures` | `*bool` | `notify_on_all_expectation_failures` | no | Should the user receive expectation failures and misses via email? |
| `RequireLoginBy` | `*time.Time` | `require_login_by` | no | Require user to login by specified date otherwise it will be disabled. |
| `RequirePasswordChange` | `*bool` | `require_password_change` | no | Is a password change required upon next user login? |
| `ResponsibleGroupId` | `int64` | `responsible_group_id` | no | ID of the internal Group responsible for this Partner User, overriding the Partner default. |
| `ResponsibleUserId` | `int64` | `responsible_user_id` | no | ID of the internal User responsible for this Partner User, overriding the Partner default. |
| `RestapiPermission` | `*bool` | `restapi_permission` | no | Can this user access the Web app, Desktop app, SDKs, or REST API? (All of these tools use the API internally, so this is one unified permission set.) |
| `S3CompatibleEndpointPermission` | `*bool` | `s3_compatible_endpoint_permission` | no | Can the user access the S3-compatible endpoint? Defaults to true. |
| `SelfManaged` | `*bool` | `self_managed` | no | Does this user manage it's own credentials or is it a shared/bot user? |
| `SftpPermission` | `*bool` | `sftp_permission` | no | Can the user access with SFTP? |
| `SiteAdmin` | `*bool` | `site_admin` | no | Is the user an administrator for this site? |
| `SkipWelcomeScreen` | `*bool` | `skip_welcome_screen` | no | Skip Welcome page in the UI? |
| `SslRequired` | `files_sdk.UserSslRequiredEnum` | `ssl_required` | no | SSL required setting Values: `use_system_setting`, `always_require`, `never_require`. |
| `SsoStrategyId` | `int64` | `sso_strategy_id` | no | SSO (Single Sign On) strategy ID for the user, if applicable. |
| `Require2fa` | `files_sdk.UserRequire2faEnum` | `require_2fa` | no | 2FA required setting. `use_system_setting` uses the site-wide setting, including SSO exemptions. `always_require` and `never_require` override the site-wide setting when user-level overrides are allowed. Values: `use_system_setting`, `always_require`, `never_require`. |
| `Tags` | `string` | `tags` | no | Comma-separated list of Tags for this user. Tags are used for other features, such as UserLifecycleRules, which can target specific tags. Tags must only contain lowercase letters, numbers, and hyphens. |
| `TimeZone` | `string` | `time_zone` | no | User time zone |
| `UserRoot` | `string` | `user_root` | no | If filesystem layout is user_root, this path is the root path the user is fixed to for all interfaces. If the filesystem layout is site_root or partner_root, this acts as a root folder only for FTP and SFTP (SFTP applicability also requires a site-wide setting to be set). For partner_root layout, this path is relative to the Partner root folder for all callers and blank opts out of an additional protocol root. In this situation, this path is not applied to the API, Desktop, or Web interface. |
| `UserHome` | `string` | `user_home` | no | Home folder for FTP/SFTP. For users with the partner_root filesystem layout, this path is relative to the Partner root folder. In all other cases, it is an absolute path. Only applies to FTP and SFTP, and not any other interface. |
| `WorkspaceAdmin` | `*bool` | `workspace_admin` | no | Whether the user is an administrator of their own Custom Workspace. Does not reflect administration granted through Permissions. |
| `Username` | `string` | `username` | no | User's username |
| `WorkspaceId` | `int64` | `workspace_id` | no | Workspace ID. Only Site Administrators can change this field. Values supplied by Workspace Administrators, Group Administrators, or other non-Site Administrators using `/user` are ignored. |
| `Clear2fa` | `*bool` | `clear_2fa` | no | If true when changing authentication_method from `password` to `sso`, remove all two-factor methods. Ignored in all other cases. |
| `ConvertToPartnerUser` | `*bool` | `convert_to_partner_user` | no | Required when assigning a Partner to an existing non-Partner user. If true, convert the user by assigning the partner_id provided. |

Example:

```go
result, err := user.Update(files_sdk.UserUpdateParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := user.UpdateWithMap(map[string]interface{}{
	"id": 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /users/{id}`

Delete User

```go
func (c *Client) Delete(params files_sdk.UserDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.UserDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.UserDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | User ID. |
| `NewOwnerId` | `int64` | `new_owner_id` | no | Provide a User ID here to transfer ownership of certain resources such as Automations and Share Links (Bundles) to that new user. |

Example:

```go
err := user.Delete(files_sdk.UserDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.User` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | User ID |
| `Username` | `string` | `username` | User's username |
| `AdminGroupIds` | `[]int64` | `admin_group_ids` | List of group IDs of which this user is an administrator |
| `AllowedIps` | `string` | `allowed_ips` | A list of allowed IPs if applicable. Newline delimited |
| `AttachmentsPermission` | `*bool` | `attachments_permission` | If `true`, the user can user create Bundles (aka Share Links). Use the bundle permission instead. |
| `ApiKeysCount` | `int64` | `api_keys_count` | Number of API keys associated with this user |
| `AuthenticateUntil` | `*time.Time` | `authenticate_until` | Scheduled Date/Time at which user will be deactivated |
| `AuthenticationMethod` | `string` | `authentication_method` | How is this user authenticated? |
| `AvatarUrl` | `string` | `avatar_url` | URL holding the user's avatar |
| `Billable` | `*bool` | `billable` | Is this a billable user record? |
| `BillingPermission` | `*bool` | `billing_permission` | Allow this user to perform operations on the account, payments, and invoices? |
| `BypassSiteAllowedIps` | `*bool` | `bypass_site_allowed_ips` | Allow this user to skip site-wide IP blacklists? |
| `BypassUserLifecycleRules` | `*bool` | `bypass_user_lifecycle_rules` | Exempt this user from user lifecycle rules? |
| `CreatedAt` | `*time.Time` | `created_at` | When this user was created |
| `DavPermission` | `*bool` | `dav_permission` | Can the user connect with WebDAV? |
| `Disabled` | `*bool` | `disabled` | Is user disabled? Disabled users cannot log in, and do not count for billing purposes. Users can be automatically disabled after an inactivity period via a Site setting or schedule to be deactivated after specific date. |
| `DisabledExpiredOrInactive` | `*bool` | `disabled_expired_or_inactive` | Computed property that returns true if user disabled or expired or inactive. |
| `AiAssistantPersonalityId` | `int64` | `ai_assistant_personality_id` | AI Assistant Personality ID assigned directly to this user, if any. |
| `DesktopConfigurationProfileId` | `int64` | `desktop_configuration_profile_id` | Desktop Configuration Profile ID assigned directly to this user, if any. |
| `IntegrationCentricProfileId` | `int64` | `integration_centric_profile_id` | Integration Centric Profile ID assigned directly to this user, if any. |
| `Email` | `string` | `email` | User email address |
| `FilesystemLayout` | `string` | `filesystem_layout` | File system layout |
| `FirstLoginAt` | `*time.Time` | `first_login_at` | User's first login time |
| `FtpPermission` | `*bool` | `ftp_permission` | Can the user access with FTP/FTPS? |
| `GroupIds` | `string` | `group_ids` | Comma-separated list of group IDs of which this user is a member |
| `HeaderText` | `string` | `header_text` | Text to display to the user in the header of the UI |
| `Language` | `string` | `language` | Preferred language |
| `LastLoginAt` | `*time.Time` | `last_login_at` | User's most recent login time via any protocol |
| `LastWebLoginAt` | `*time.Time` | `last_web_login_at` | User's most recent login time via web |
| `LastFtpLoginAt` | `*time.Time` | `last_ftp_login_at` | User's most recent login time via FTP |
| `LastSftpLoginAt` | `*time.Time` | `last_sftp_login_at` | User's most recent login time via SFTP |
| `LastDavLoginAt` | `*time.Time` | `last_dav_login_at` | User's most recent login time via WebDAV |
| `LastDesktopLoginAt` | `*time.Time` | `last_desktop_login_at` | User's most recent login time via Desktop app |
| `LastRestapiLoginAt` | `*time.Time` | `last_restapi_login_at` | User's most recent login time via Rest API |
| `LastApiUseAt` | `*time.Time` | `last_api_use_at` | User's most recent API use time |
| `LastActiveAt` | `*time.Time` | `last_active_at` | User's most recent activity time, which is the latest of most recent login, most recent API use, enablement, or creation |
| `LastProtocolCipher` | `string` | `last_protocol_cipher` | The most recent protocol and cipher used |
| `LockoutExpires` | `*time.Time` | `lockout_expires` | Time in the future that the user will no longer be locked out if applicable |
| `Name` | `string` | `name` | User's full name |
| `Company` | `string` | `company` | User's company |
| `Notes` | `string` | `notes` | Any internal notes on the user |
| `NotificationDailySendTime` | `int64` | `notification_daily_send_time` | Hour of the day at which daily notifications should be sent. Can be in range 0 to 23 |
| `OfficeIntegrationEnabled` | `*bool` | `office_integration_enabled` | Enable integration with Office for the web? |
| `PartnerAdmin` | `*bool` | `partner_admin` | Is this user a Partner administrator? |
| `PartnerId` | `int64` | `partner_id` | Partner ID if this user belongs to a Partner |
| `PartnerName` | `string` | `partner_name` | Name of the Partner if this user belongs to a Partner |
| `PasswordSetAt` | `*time.Time` | `password_set_at` | Last time the user's password was set |
| `PasswordValidityDays` | `int64` | `password_validity_days` | Number of days to allow user to use the same password |
| `PrimaryGroupId` | `int64` | `primary_group_id` | Primary group ID for Group Admin scoping |
| `PublicKeysCount` | `int64` | `public_keys_count` | Number of public keys associated with this user |
| `ReceiveAdminAlerts` | `*bool` | `receive_admin_alerts` | Deprecated. Use notify_on_all_site_warnings and granular failure notification preferences instead. |
| `NotifyOnAllSiteWarnings` | `*bool` | `notify_on_all_site_warnings` | Should the user receive site warnings via email? |
| `NotifyOnAllSsoFailures` | `*bool` | `notify_on_all_sso_failures` | Should the user receive sso/scim/ldap configuration/sync failures via email? |
| `NotifyOnAllUserSecurityEvents` | `*bool` | `notify_on_all_user_security_events` | Should the user receive user security events via email? |
| `NotifyOnAllPendingWorkFailures` | `*bool` | `notify_on_all_pending_work_failures` | Should the user receive pending work failures via email? |
| `NotifyOnAllSiemHttpDestinationFailures` | `*bool` | `notify_on_all_siem_http_destination_failures` | Should the user receive siem failures via email? |
| `NotifyOnAllSyncFailures` | `*bool` | `notify_on_all_sync_failures` | Should the user receive sync failures via email? |
| `NotifyOnAllAutomationFailures` | `*bool` | `notify_on_all_automation_failures` | Should the user receive automation failures via email? |
| `NotifyOnAllExpectationFailures` | `*bool` | `notify_on_all_expectation_failures` | Should the user receive expectation failures and misses via email? |
| `Require2fa` | `string` | `require_2fa` | 2FA required setting. `use_system_setting` uses the site-wide setting, including SSO exemptions. `always_require` and `never_require` override the site-wide setting when user-level overrides are allowed. |
| `RequireLoginBy` | `*time.Time` | `require_login_by` | Require user to login by specified date otherwise it will be disabled. |
| `Active2fa` | `*bool` | `active_2fa` | Is 2fa active for the user? |
| `RequirePasswordChange` | `*bool` | `require_password_change` | Is a password change required upon next user login? |
| `PasswordExpired` | `*bool` | `password_expired` | Is user's password expired? |
| `ResponsibleGroupId` | `int64` | `responsible_group_id` | ID of the internal Group responsible for this Partner User, overriding the Partner default. |
| `ResponsibleUserId` | `int64` | `responsible_user_id` | ID of the internal User responsible for this Partner User, overriding the Partner default. |
| `ReadonlySiteAdmin` | `*bool` | `readonly_site_admin` | Is the user an allowed to view all (non-billing) site configuration for this site? |
| `RestapiPermission` | `*bool` | `restapi_permission` | Can this user access the Web app, Desktop app, SDKs, or REST API? (All of these tools use the API internally, so this is one unified permission set.) |
| `S3CompatibleEndpointPermission` | `*bool` | `s3_compatible_endpoint_permission` | Can the user access the S3-compatible endpoint? Defaults to true. |
| `SelfManaged` | `*bool` | `self_managed` | Does this user manage it's own credentials or is it a shared/bot user? |
| `SftpPermission` | `*bool` | `sftp_permission` | Can the user access with SFTP? |
| `SiteAdmin` | `*bool` | `site_admin` | Is the user an administrator for this site? |
| `WorkspaceAdmin` | `*bool` | `workspace_admin` | Whether the user is an administrator of their own Custom Workspace. Does not reflect administration granted through Permissions. |
| `SiteId` | `int64` | `site_id` | Site ID |
| `WorkspaceId` | `int64` | `workspace_id` | ID of the Workspace the user belongs to. 0 is the Default Workspace. |
| `DefaultWorkspaceId` | `int64` | `default_workspace_id` | Workspace ID the user should land in by default when more than one Workspace is available. |
| `SkipWelcomeScreen` | `*bool` | `skip_welcome_screen` | Skip Welcome page in the UI? |
| `SslRequired` | `string` | `ssl_required` | SSL required setting |
| `SsoStrategyId` | `int64` | `sso_strategy_id` | SSO (Single Sign On) strategy ID for the user, if applicable. |
| `ExternallyManaged` | `*bool` | `externally_managed` | Is this user managed by a SsoStrategy? |
| `Tags` | `string` | `tags` | Comma-separated list of Tags for this user. Tags are used for other features, such as UserLifecycleRules, which can target specific tags. Tags must only contain lowercase letters, numbers, and hyphens. |
| `TimeZone` | `string` | `time_zone` | User time zone |
| `TypeOf2fa` | `string` | `type_of_2fa` | Type(s) of 2FA methods in use, for programmatic use. Will be either `sms`, `totp`, `webauthn`, `yubi`, `email`, or multiple values sorted alphabetically and joined by an underscore. Does not specify whether user has more than one of a given method. |
| `TypeOf2faForDisplay` | `string` | `type_of_2fa_for_display` | Type(s) of 2FA methods in use, formatted for displaying in the UI. Unlike `type_of_2fa`, this value will make clear when a user has more than 1 of the same type of method. |
| `UserRoot` | `string` | `user_root` | If filesystem layout is user_root, this path is the root path the user is fixed to for all interfaces. If the filesystem layout is site_root or partner_root, this acts as a root folder only for FTP and SFTP (SFTP applicability also requires a site-wide setting to be set). For partner_root layout, this path is relative to the Partner root folder for all callers and blank opts out of an additional protocol root. In this situation, this path is not applied to the API, Desktop, or Web interface. |
| `UserHome` | `string` | `user_home` | Home folder for FTP/SFTP. For users with the partner_root filesystem layout, this path is relative to the Partner root folder. In all other cases, it is an absolute path. Only applies to FTP and SFTP, and not any other interface. |
| `DaysRemainingUntilPasswordExpire` | `int64` | `days_remaining_until_password_expire` | Number of days remaining until password expires |
| `PasswordExpireAt` | `*time.Time` | `password_expire_at` | Password expiration datetime |
| `AvatarFile` | `io.Reader` | `avatar_file` | An image file for your user avatar. |
| `AvatarDelete` | `*bool` | `avatar_delete` | If true, the avatar will be deleted. |
| `ChangePassword` | `string` | `change_password` | Used for changing a password on an existing user. |
| `ChangePasswordConfirmation` | `string` | `change_password_confirmation` | Optional, but if provided, we will ensure that it matches the value sent in `change_password`. |
| `GrantPermission` | `string` | `grant_permission` | Permission to grant on the User Root upon creation, or on the Workspace root when the User Root is blank. For the Default Workspace, a blank User Root uses the site root. `admin` on a Custom Workspace's root grants Workspace Administrator access; on a subfolder it grants Folder Admin access. Can be blank or `admin`, `full`, `read`, `write`, `list`, `read+write`, `list+write`, or `history`. |
| `GroupId` | `int64` | `group_id` | Group ID to associate this user with. |
| `ImportedPasswordHash` | `string` | `imported_password_hash` | Pre-calculated hash of the user's password. If supplied, this will be used to authenticate the user on first login. Supported hash methods are MD5, SHA1, and SHA256. |
| `Password` | `string` | `password` | User password. |
| `PasswordConfirmation` | `string` | `password_confirmation` | Optional, but if provided, we will ensure that it matches the value sent in `password`. |
| `AnnouncementsRead` | `*bool` | `announcements_read` | Signifies that the user has read all the announcements in the UI. |
| `Clear2fa` | `*bool` | `clear_2fa` | If true when changing authentication_method from `password` to `sso`, remove all two-factor methods. Ignored in all other cases. |
| `ConvertToPartnerUser` | `*bool` | `convert_to_partner_user` | Required when assigning a Partner to an existing non-Partner user. If true, convert the user by assigning the partner_id provided. |
