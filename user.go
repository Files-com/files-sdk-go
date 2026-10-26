package files_sdk

import (
	"encoding/json"
	"io"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// User is a Files.com API resource.
type User struct {
	// User ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// User's username.
	Username string `json:"username,omitempty" path:"username,omitempty" url:"username,omitempty"`
	// List of group IDs of which this user is an administrator.
	AdminGroupIds []int64 `json:"admin_group_ids,omitempty" path:"admin_group_ids,omitempty" url:"admin_group_ids,omitempty"`
	// A list of allowed IPs if applicable. Newline delimited.
	AllowedIps string `json:"allowed_ips,omitempty" path:"allowed_ips,omitempty" url:"allowed_ips,omitempty"`
	// If `true`, the user can user create Bundles (aka Share Links). Use the bundle permission instead.
	AttachmentsPermission *bool `json:"attachments_permission,omitempty" path:"attachments_permission,omitempty" url:"attachments_permission,omitempty"`
	// Number of API keys associated with this user.
	ApiKeysCount int64 `json:"api_keys_count,omitempty" path:"api_keys_count,omitempty" url:"api_keys_count,omitempty"`
	// Scheduled Date/Time at which user will be deactivated.
	AuthenticateUntil *time.Time `json:"authenticate_until,omitempty" path:"authenticate_until,omitempty" url:"authenticate_until,omitempty"`
	// How is this user authenticated?
	AuthenticationMethod string `json:"authentication_method,omitempty" path:"authentication_method,omitempty" url:"authentication_method,omitempty"`
	// URL holding the user's avatar.
	AvatarUrl string `json:"avatar_url,omitempty" path:"avatar_url,omitempty" url:"avatar_url,omitempty"`
	// Is this a billable user record?
	Billable *bool `json:"billable,omitempty" path:"billable,omitempty" url:"billable,omitempty"`
	// Allow this user to perform operations on the account, payments, and invoices?
	BillingPermission *bool `json:"billing_permission,omitempty" path:"billing_permission,omitempty" url:"billing_permission,omitempty"`
	// Allow this user to skip site-wide IP blacklists?
	BypassSiteAllowedIps *bool `json:"bypass_site_allowed_ips,omitempty" path:"bypass_site_allowed_ips,omitempty" url:"bypass_site_allowed_ips,omitempty"`
	// Exempt this user from user lifecycle rules?
	BypassUserLifecycleRules *bool `json:"bypass_user_lifecycle_rules,omitempty" path:"bypass_user_lifecycle_rules,omitempty" url:"bypass_user_lifecycle_rules,omitempty"`
	// When this user was created.
	CreatedAt *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	// Can the user connect with WebDAV?
	DavPermission *bool `json:"dav_permission,omitempty" path:"dav_permission,omitempty" url:"dav_permission,omitempty"`
	// Is user disabled? Disabled users cannot log in, and do not count for billing purposes. Users can be automatically disabled after an inactivity period via a Site setting or schedule to be deactivated after specific date.
	Disabled *bool `json:"disabled,omitempty" path:"disabled,omitempty" url:"disabled,omitempty"`
	// Computed property that returns true if user disabled or expired or inactive.
	DisabledExpiredOrInactive *bool `json:"disabled_expired_or_inactive,omitempty" path:"disabled_expired_or_inactive,omitempty" url:"disabled_expired_or_inactive,omitempty"`
	// AI Assistant Personality ID assigned directly to this user, if any.
	AiAssistantPersonalityId int64 `json:"ai_assistant_personality_id,omitempty" path:"ai_assistant_personality_id,omitempty" url:"ai_assistant_personality_id,omitempty"`
	// Desktop Configuration Profile ID assigned directly to this user, if any.
	DesktopConfigurationProfileId int64 `json:"desktop_configuration_profile_id,omitempty" path:"desktop_configuration_profile_id,omitempty" url:"desktop_configuration_profile_id,omitempty"`
	// Integration Centric Profile ID assigned directly to this user, if any.
	IntegrationCentricProfileId int64 `json:"integration_centric_profile_id,omitempty" path:"integration_centric_profile_id,omitempty" url:"integration_centric_profile_id,omitempty"`
	// User email address.
	Email string `json:"email,omitempty" path:"email,omitempty" url:"email,omitempty"`
	// File system layout.
	FilesystemLayout string `json:"filesystem_layout,omitempty" path:"filesystem_layout,omitempty" url:"filesystem_layout,omitempty"`
	// User's first login time.
	FirstLoginAt *time.Time `json:"first_login_at,omitempty" path:"first_login_at,omitempty" url:"first_login_at,omitempty"`
	// Can the user access with FTP/FTPS?
	FtpPermission *bool `json:"ftp_permission,omitempty" path:"ftp_permission,omitempty" url:"ftp_permission,omitempty"`
	// Comma-separated list of group IDs of which this user is a member.
	GroupIds string `json:"group_ids,omitempty" path:"group_ids,omitempty" url:"group_ids,omitempty"`
	// Text to display to the user in the header of the UI.
	HeaderText string `json:"header_text,omitempty" path:"header_text,omitempty" url:"header_text,omitempty"`
	// Preferred language.
	Language string `json:"language,omitempty" path:"language,omitempty" url:"language,omitempty"`
	// User's most recent login time via any protocol.
	LastLoginAt *time.Time `json:"last_login_at,omitempty" path:"last_login_at,omitempty" url:"last_login_at,omitempty"`
	// User's most recent login time via web.
	LastWebLoginAt *time.Time `json:"last_web_login_at,omitempty" path:"last_web_login_at,omitempty" url:"last_web_login_at,omitempty"`
	// User's most recent login time via FTP.
	LastFtpLoginAt *time.Time `json:"last_ftp_login_at,omitempty" path:"last_ftp_login_at,omitempty" url:"last_ftp_login_at,omitempty"`
	// User's most recent login time via SFTP.
	LastSftpLoginAt *time.Time `json:"last_sftp_login_at,omitempty" path:"last_sftp_login_at,omitempty" url:"last_sftp_login_at,omitempty"`
	// User's most recent login time via WebDAV.
	LastDavLoginAt *time.Time `json:"last_dav_login_at,omitempty" path:"last_dav_login_at,omitempty" url:"last_dav_login_at,omitempty"`
	// User's most recent login time via Desktop app.
	LastDesktopLoginAt *time.Time `json:"last_desktop_login_at,omitempty" path:"last_desktop_login_at,omitempty" url:"last_desktop_login_at,omitempty"`
	// User's most recent login time via Rest API.
	LastRestapiLoginAt *time.Time `json:"last_restapi_login_at,omitempty" path:"last_restapi_login_at,omitempty" url:"last_restapi_login_at,omitempty"`
	// User's most recent API use time.
	LastApiUseAt *time.Time `json:"last_api_use_at,omitempty" path:"last_api_use_at,omitempty" url:"last_api_use_at,omitempty"`
	// User's most recent activity time, which is the latest of most recent login, most recent API use, enablement, or creation.
	LastActiveAt *time.Time `json:"last_active_at,omitempty" path:"last_active_at,omitempty" url:"last_active_at,omitempty"`
	// The most recent protocol and cipher used.
	LastProtocolCipher string `json:"last_protocol_cipher,omitempty" path:"last_protocol_cipher,omitempty" url:"last_protocol_cipher,omitempty"`
	// Time in the future that the user will no longer be locked out if applicable.
	LockoutExpires *time.Time `json:"lockout_expires,omitempty" path:"lockout_expires,omitempty" url:"lockout_expires,omitempty"`
	// User's full name.
	Name string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	// User's company.
	Company string `json:"company,omitempty" path:"company,omitempty" url:"company,omitempty"`
	// Any internal notes on the user.
	Notes string `json:"notes,omitempty" path:"notes,omitempty" url:"notes,omitempty"`
	// Hour of the day at which daily notifications should be sent. Can be in range 0 to 23.
	NotificationDailySendTime int64 `json:"notification_daily_send_time,omitempty" path:"notification_daily_send_time,omitempty" url:"notification_daily_send_time,omitempty"`
	// Enable integration with Office for the web?
	OfficeIntegrationEnabled *bool `json:"office_integration_enabled,omitempty" path:"office_integration_enabled,omitempty" url:"office_integration_enabled,omitempty"`
	// Is this user a Partner administrator?
	PartnerAdmin *bool `json:"partner_admin,omitempty" path:"partner_admin,omitempty" url:"partner_admin,omitempty"`
	// Partner ID if this user belongs to a Partner.
	PartnerId int64 `json:"partner_id,omitempty" path:"partner_id,omitempty" url:"partner_id,omitempty"`
	// Name of the Partner if this user belongs to a Partner.
	PartnerName string `json:"partner_name,omitempty" path:"partner_name,omitempty" url:"partner_name,omitempty"`
	// Last time the user's password was set.
	PasswordSetAt *time.Time `json:"password_set_at,omitempty" path:"password_set_at,omitempty" url:"password_set_at,omitempty"`
	// Number of days to allow user to use the same password.
	PasswordValidityDays int64 `json:"password_validity_days,omitempty" path:"password_validity_days,omitempty" url:"password_validity_days,omitempty"`
	// Primary group ID for Group Admin scoping.
	PrimaryGroupId int64 `json:"primary_group_id,omitempty" path:"primary_group_id,omitempty" url:"primary_group_id,omitempty"`
	// Number of public keys associated with this user.
	PublicKeysCount int64 `json:"public_keys_count,omitempty" path:"public_keys_count,omitempty" url:"public_keys_count,omitempty"`
	// Deprecated. Use notify_on_all_site_warnings and granular failure notification preferences instead.
	ReceiveAdminAlerts *bool `json:"receive_admin_alerts,omitempty" path:"receive_admin_alerts,omitempty" url:"receive_admin_alerts,omitempty"`
	// Should the user receive site warnings via email?
	NotifyOnAllSiteWarnings *bool `json:"notify_on_all_site_warnings,omitempty" path:"notify_on_all_site_warnings,omitempty" url:"notify_on_all_site_warnings,omitempty"`
	// Should the user receive sso/scim/ldap configuration/sync failures via email?
	NotifyOnAllSsoFailures *bool `json:"notify_on_all_sso_failures,omitempty" path:"notify_on_all_sso_failures,omitempty" url:"notify_on_all_sso_failures,omitempty"`
	// Should the user receive user security events via email?
	NotifyOnAllUserSecurityEvents *bool `json:"notify_on_all_user_security_events,omitempty" path:"notify_on_all_user_security_events,omitempty" url:"notify_on_all_user_security_events,omitempty"`
	// Should the user receive pending work failures via email?
	NotifyOnAllPendingWorkFailures *bool `json:"notify_on_all_pending_work_failures,omitempty" path:"notify_on_all_pending_work_failures,omitempty" url:"notify_on_all_pending_work_failures,omitempty"`
	// Should the user receive siem failures via email?
	NotifyOnAllSiemHttpDestinationFailures *bool `json:"notify_on_all_siem_http_destination_failures,omitempty" path:"notify_on_all_siem_http_destination_failures,omitempty" url:"notify_on_all_siem_http_destination_failures,omitempty"`
	// Should the user receive sync failures via email?
	NotifyOnAllSyncFailures *bool `json:"notify_on_all_sync_failures,omitempty" path:"notify_on_all_sync_failures,omitempty" url:"notify_on_all_sync_failures,omitempty"`
	// Should the user receive automation failures via email?
	NotifyOnAllAutomationFailures *bool `json:"notify_on_all_automation_failures,omitempty" path:"notify_on_all_automation_failures,omitempty" url:"notify_on_all_automation_failures,omitempty"`
	// Should the user receive expectation failures and misses via email?
	NotifyOnAllExpectationFailures *bool `json:"notify_on_all_expectation_failures,omitempty" path:"notify_on_all_expectation_failures,omitempty" url:"notify_on_all_expectation_failures,omitempty"`
	// 2FA required setting. `use_system_setting` uses the site-wide setting, including SSO exemptions. `always_require` and `never_require` override the site-wide setting when user-level overrides are allowed.
	Require2fa string `json:"require_2fa,omitempty" path:"require_2fa,omitempty" url:"require_2fa,omitempty"`
	// Require user to login by specified date otherwise it will be disabled.
	RequireLoginBy *time.Time `json:"require_login_by,omitempty" path:"require_login_by,omitempty" url:"require_login_by,omitempty"`
	// Is 2fa active for the user?
	Active2fa *bool `json:"active_2fa,omitempty" path:"active_2fa,omitempty" url:"active_2fa,omitempty"`
	// Is a password change required upon next user login?
	RequirePasswordChange *bool `json:"require_password_change,omitempty" path:"require_password_change,omitempty" url:"require_password_change,omitempty"`
	// Is user's password expired?
	PasswordExpired *bool `json:"password_expired,omitempty" path:"password_expired,omitempty" url:"password_expired,omitempty"`
	// ID of the internal Group responsible for this Partner User, overriding the Partner default.
	ResponsibleGroupId int64 `json:"responsible_group_id,omitempty" path:"responsible_group_id,omitempty" url:"responsible_group_id,omitempty"`
	// ID of the internal User responsible for this Partner User, overriding the Partner default.
	ResponsibleUserId int64 `json:"responsible_user_id,omitempty" path:"responsible_user_id,omitempty" url:"responsible_user_id,omitempty"`
	// Is the user an allowed to view all (non-billing) site configuration for this site?
	ReadonlySiteAdmin *bool `json:"readonly_site_admin,omitempty" path:"readonly_site_admin,omitempty" url:"readonly_site_admin,omitempty"`
	// Can this user access the Web app, Desktop app, SDKs, or REST API? (All of these tools use the API internally, so this is one unified permission set.).
	RestapiPermission *bool `json:"restapi_permission,omitempty" path:"restapi_permission,omitempty" url:"restapi_permission,omitempty"`
	// Can the user access the S3-compatible endpoint? Defaults to true.
	S3CompatibleEndpointPermission *bool `json:"s3_compatible_endpoint_permission,omitempty" path:"s3_compatible_endpoint_permission,omitempty" url:"s3_compatible_endpoint_permission,omitempty"`
	// Does this user manage it's own credentials or is it a shared/bot user?
	SelfManaged *bool `json:"self_managed,omitempty" path:"self_managed,omitempty" url:"self_managed,omitempty"`
	// Can the user access with SFTP?
	SftpPermission *bool `json:"sftp_permission,omitempty" path:"sftp_permission,omitempty" url:"sftp_permission,omitempty"`
	// Is the user an administrator for this site?
	SiteAdmin *bool `json:"site_admin,omitempty" path:"site_admin,omitempty" url:"site_admin,omitempty"`
	// Whether the user is an administrator of their own Custom Workspace. Does not reflect administration granted through Permissions.
	WorkspaceAdmin *bool `json:"workspace_admin,omitempty" path:"workspace_admin,omitempty" url:"workspace_admin,omitempty"`
	// Site ID.
	SiteId int64 `json:"site_id,omitempty" path:"site_id,omitempty" url:"site_id,omitempty"`
	// ID of the Workspace the user belongs to. 0 is the Default Workspace.
	WorkspaceId int64 `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	// Workspace ID the user should land in by default when more than one Workspace is available.
	DefaultWorkspaceId int64 `json:"default_workspace_id,omitempty" path:"default_workspace_id,omitempty" url:"default_workspace_id,omitempty"`
	// Skip Welcome page in the UI?
	SkipWelcomeScreen *bool `json:"skip_welcome_screen,omitempty" path:"skip_welcome_screen,omitempty" url:"skip_welcome_screen,omitempty"`
	// SSL required setting.
	SslRequired string `json:"ssl_required,omitempty" path:"ssl_required,omitempty" url:"ssl_required,omitempty"`
	// SSO (Single Sign On) strategy ID for the user, if applicable.
	SsoStrategyId int64 `json:"sso_strategy_id,omitempty" path:"sso_strategy_id,omitempty" url:"sso_strategy_id,omitempty"`
	// Is this user managed by a SsoStrategy?
	ExternallyManaged *bool `json:"externally_managed,omitempty" path:"externally_managed,omitempty" url:"externally_managed,omitempty"`
	// Comma-separated list of Tags for this user. Tags are used for other features, such as UserLifecycleRules, which can target specific tags. Tags must only contain lowercase letters, numbers, and hyphens.
	Tags string `json:"tags,omitempty" path:"tags,omitempty" url:"tags,omitempty"`
	// User time zone.
	TimeZone string `json:"time_zone,omitempty" path:"time_zone,omitempty" url:"time_zone,omitempty"`
	// Type(s) of 2FA methods in use, for programmatic use. Will be either `sms`, `totp`, `webauthn`, `yubi`, `email`, or multiple values sorted alphabetically and joined by an underscore. Does not specify whether user has more than one of a given method.
	TypeOf2fa string `json:"type_of_2fa,omitempty" path:"type_of_2fa,omitempty" url:"type_of_2fa,omitempty"`
	// Type(s) of 2FA methods in use, formatted for displaying in the UI. Unlike `type_of_2fa`, this value will make clear when a user has more than 1 of the same type of method.
	TypeOf2faForDisplay string `json:"type_of_2fa_for_display,omitempty" path:"type_of_2fa_for_display,omitempty" url:"type_of_2fa_for_display,omitempty"`
	// If filesystem layout is user_root, this path is the root path the user is fixed to for all interfaces. If the filesystem layout is site_root or partner_root, this acts as a root folder only for FTP and SFTP (SFTP applicability also requires a site-wide setting to be set). For partner_root layout, this path is relative to the Partner root folder for all callers and blank opts out of an additional protocol root. In this situation, this path is not applied to the API, Desktop, or Web interface.
	UserRoot string `json:"user_root,omitempty" path:"user_root,omitempty" url:"user_root,omitempty"`
	// Home folder for FTP/SFTP. For users with the partner_root filesystem layout, this path is relative to the Partner root folder. In all other cases, it is an absolute path. Only applies to FTP and SFTP, and not any other interface.
	UserHome string `json:"user_home,omitempty" path:"user_home,omitempty" url:"user_home,omitempty"`
	// Number of days remaining until password expires.
	DaysRemainingUntilPasswordExpire int64 `json:"days_remaining_until_password_expire,omitempty" path:"days_remaining_until_password_expire,omitempty" url:"days_remaining_until_password_expire,omitempty"`
	// Password expiration datetime.
	PasswordExpireAt *time.Time `json:"password_expire_at,omitempty" path:"password_expire_at,omitempty" url:"password_expire_at,omitempty"`
	// An image file for your user avatar.
	AvatarFile io.Reader `json:"avatar_file,omitempty" path:"avatar_file,omitempty" url:"avatar_file,omitempty"`
	// If true, the avatar will be deleted.
	AvatarDelete *bool `json:"avatar_delete,omitempty" path:"avatar_delete,omitempty" url:"avatar_delete,omitempty"`
	// Used for changing a password on an existing user.
	ChangePassword string `json:"change_password,omitempty" path:"change_password,omitempty" url:"change_password,omitempty"`
	// Optional, but if provided, we will ensure that it matches the value sent in `change_password`.
	ChangePasswordConfirmation string `json:"change_password_confirmation,omitempty" path:"change_password_confirmation,omitempty" url:"change_password_confirmation,omitempty"`
	// Permission to grant on the User Root upon creation, or on the Workspace root when the User Root is blank. For the Default Workspace, a blank User Root uses the site root. `admin` on a Custom Workspace's root grants Workspace Administrator access; on a subfolder it grants Folder Admin access. Can be blank or `admin`, `full`, `read`, `write`, `list`, `read+write`, `list+write`, or `history`.
	GrantPermission string `json:"grant_permission,omitempty" path:"grant_permission,omitempty" url:"grant_permission,omitempty"`
	// Group ID to associate this user with.
	GroupId int64 `json:"group_id,omitempty" path:"group_id,omitempty" url:"group_id,omitempty"`
	// Pre-calculated hash of the user's password. If supplied, this will be used to authenticate the user on first login. Supported hash methods are MD5, SHA1, and SHA256.
	ImportedPasswordHash string `json:"imported_password_hash,omitempty" path:"imported_password_hash,omitempty" url:"imported_password_hash,omitempty"`
	// User password.
	Password string `json:"password,omitempty" path:"password,omitempty" url:"password,omitempty"`
	// Optional, but if provided, we will ensure that it matches the value sent in `password`.
	PasswordConfirmation string `json:"password_confirmation,omitempty" path:"password_confirmation,omitempty" url:"password_confirmation,omitempty"`
	// Signifies that the user has read all the announcements in the UI.
	AnnouncementsRead *bool `json:"announcements_read,omitempty" path:"announcements_read,omitempty" url:"announcements_read,omitempty"`
	// If true when changing authentication_method from `password` to `sso`, remove all two-factor methods. Ignored in all other cases.
	Clear2fa *bool `json:"clear_2fa,omitempty" path:"clear_2fa,omitempty" url:"clear_2fa,omitempty"`
	// Required when assigning a Partner to an existing non-Partner user. If true, convert the user by assigning the partner_id provided.
	ConvertToPartnerUser *bool `json:"convert_to_partner_user,omitempty" path:"convert_to_partner_user,omitempty" url:"convert_to_partner_user,omitempty"`
}

// Identifier returns the resource ID.
func (u User) Identifier() interface{} {
	return u.Id
}

// UserCollection is a list of User resources.
type UserCollection []User

// UserAuthenticationMethodEnum is a string value for authentication_method.
// Enum lists the values documented by the API.
type UserAuthenticationMethodEnum string

// String returns the API parameter value.
func (u UserAuthenticationMethodEnum) String() string {
	return string(u)
}

// Enum returns the documented values keyed by their API strings.
func (u UserAuthenticationMethodEnum) Enum() map[string]UserAuthenticationMethodEnum {
	return map[string]UserAuthenticationMethodEnum{
		"password":                    UserAuthenticationMethodEnum("password"),
		"sso":                         UserAuthenticationMethodEnum("sso"),
		"none":                        UserAuthenticationMethodEnum("none"),
		"email_signup":                UserAuthenticationMethodEnum("email_signup"),
		"password_with_imported_hash": UserAuthenticationMethodEnum("password_with_imported_hash"),
		"password_and_ssh_key":        UserAuthenticationMethodEnum("password_and_ssh_key"),
	}
}

// UserFilesystemLayoutEnum is a string value for filesystem_layout.
// Enum lists the values documented by the API.
type UserFilesystemLayoutEnum string

// String returns the API parameter value.
func (u UserFilesystemLayoutEnum) String() string {
	return string(u)
}

// Enum returns the documented values keyed by their API strings.
func (u UserFilesystemLayoutEnum) Enum() map[string]UserFilesystemLayoutEnum {
	return map[string]UserFilesystemLayoutEnum{
		"site_root":           UserFilesystemLayoutEnum("site_root"),
		"user_root":           UserFilesystemLayoutEnum("user_root"),
		"partner_root":        UserFilesystemLayoutEnum("partner_root"),
		"integration_centric": UserFilesystemLayoutEnum("integration_centric"),
		"workspace_root":      UserFilesystemLayoutEnum("workspace_root"),
	}
}

// UserSslRequiredEnum is a string value for ssl_required.
// Enum lists the values documented by the API.
type UserSslRequiredEnum string

// String returns the API parameter value.
func (u UserSslRequiredEnum) String() string {
	return string(u)
}

// Enum returns the documented values keyed by their API strings.
func (u UserSslRequiredEnum) Enum() map[string]UserSslRequiredEnum {
	return map[string]UserSslRequiredEnum{
		"use_system_setting": UserSslRequiredEnum("use_system_setting"),
		"always_require":     UserSslRequiredEnum("always_require"),
		"never_require":      UserSslRequiredEnum("never_require"),
	}
}

// UserRequire2faEnum is a string value for require_2fa.
// Enum lists the values documented by the API.
type UserRequire2faEnum string

// String returns the API parameter value.
func (u UserRequire2faEnum) String() string {
	return string(u)
}

// Enum returns the documented values keyed by their API strings.
func (u UserRequire2faEnum) Enum() map[string]UserRequire2faEnum {
	return map[string]UserRequire2faEnum{
		"use_system_setting": UserRequire2faEnum("use_system_setting"),
		"always_require":     UserRequire2faEnum("always_require"),
		"never_require":      UserRequire2faEnum("never_require"),
	}
}

// UserListParams contains the request parameters for GET /users.
type UserListParams struct {
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `site_id`, `workspace_id`, `company`, `name`, `disabled`, `authenticate_until`, `username`, `email`, `site_admin`, `last_desktop_login_at`, `last_login_at`, `password_validity_days` or `ssl_required`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `username`, `name`, `email`, `company`, `site_admin`, `password_validity_days`, `ssl_required`, `last_login_at`, `authenticate_until`, `not_site_admin`, `disabled`, `partner_id`, `primary_group_id` or `workspace_id`. Valid field combinations are `[ site_admin, username ]`, `[ not_site_admin, username ]`, `[ workspace_id, username ]`, `[ company, name ]`, `[ workspace_id, name ]`, `[ workspace_id, email ]`, `[ workspace_id, company ]`, `[ workspace_id, site_admin ]`, `[ workspace_id, not_site_admin ]`, `[ workspace_id, disabled ]`, `[ workspace_id, partner_id ]`, `[ workspace_id, site_admin, username ]`, `[ workspace_id, not_site_admin, username ]`, `[ workspace_id, disabled, username ]`, `[ workspace_id, partner_id, username ]` or `[ workspace_id, company, name ]`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// If set, return records where the specified field is greater than the supplied value. Valid fields are `password_validity_days`, `last_login_at` or `authenticate_until`.
	FilterGt interface{} `url:"filter_gt,omitempty" json:"filter_gt,omitempty" path:"filter_gt"`
	// If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `password_validity_days`, `last_login_at` or `authenticate_until`.
	FilterGteq interface{} `url:"filter_gteq,omitempty" json:"filter_gteq,omitempty" path:"filter_gteq"`
	// If set, return records where the specified field is prefixed by the supplied value. Valid fields are `username`, `name`, `email` or `company`. Valid field combinations are `[ company, name ]`.
	FilterPrefix interface{} `url:"filter_prefix,omitempty" json:"filter_prefix,omitempty" path:"filter_prefix"`
	// If set, return records where the specified field is less than the supplied value. Valid fields are `password_validity_days`, `last_login_at` or `authenticate_until`.
	FilterLt interface{} `url:"filter_lt,omitempty" json:"filter_lt,omitempty" path:"filter_lt"`
	// If set, return records where the specified field is less than or equal the supplied value. Valid fields are `password_validity_days`, `last_login_at` or `authenticate_until`.
	FilterLteq interface{} `url:"filter_lteq,omitempty" json:"filter_lteq,omitempty" path:"filter_lteq"`
	// comma-separated list of User IDs.
	Ids string `url:"ids,omitempty" json:"ids,omitempty" path:"ids"`
	// Include users from the parent site.
	IncludeParentSiteUsers *bool `url:"include_parent_site_users,omitempty" json:"include_parent_site_users,omitempty" path:"include_parent_site_users"`
	// When scoped to a Custom Workspace, include Default Workspace users who have permissions in that Workspace. Available to Site Administrators, and to Workspace Administrators and Folder Admins in that Workspace when the site shares its Default Workspace directory with Workspaces.
	IncludeDefaultWorkspaceUsers *bool `url:"include_default_workspace_users,omitempty" json:"include_default_workspace_users,omitempty" path:"include_default_workspace_users"`
	// Searches for partial matches of name, username, or email.
	Search string `url:"search,omitempty" json:"search,omitempty" path:"search"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// UserFindParams contains the request parameters for GET /users/{id}.
type UserFindParams struct {
	// User ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UserCreateParams contains the request parameters for POST /users.
type UserCreateParams struct {
	// An image file for your user avatar.
	AvatarFile io.Writer `url:"avatar_file,omitempty" json:"avatar_file,omitempty" path:"avatar_file"`
	// AvatarFileAttachment sends file content as avatar_file; see lib.JSONAttachment. Leave AvatarFile nil when you set it.
	AvatarFileAttachment *lib.JSONAttachment `json:"-" url:"-" path:"-" attachment:"avatar_file"`
	// If true, the avatar will be deleted.
	AvatarDelete *bool `url:"avatar_delete,omitempty" json:"avatar_delete,omitempty" path:"avatar_delete"`
	// Used for changing a password on an existing user.
	ChangePassword string `url:"change_password,omitempty" json:"change_password,omitempty" path:"change_password"`
	// Optional, but if provided, we will ensure that it matches the value sent in `change_password`.
	ChangePasswordConfirmation string `url:"change_password_confirmation,omitempty" json:"change_password_confirmation,omitempty" path:"change_password_confirmation"`
	// User's email.
	Email string `url:"email,omitempty" json:"email,omitempty" path:"email"`
	// Permission to grant on the User Root upon creation, or on the Workspace root when the User Root is blank. For the Default Workspace, a blank User Root uses the site root. `admin` on a Custom Workspace's root grants Workspace Administrator access; on a subfolder it grants Folder Admin access. Can be blank or `admin`, `full`, `read`, `write`, `list`, `read+write`, `list+write`, or `history`.
	GrantPermission string `url:"grant_permission,omitempty" json:"grant_permission,omitempty" path:"grant_permission"`
	// Group ID to associate this user with.
	GroupId int64 `url:"group_id,omitempty" json:"group_id,omitempty" path:"group_id"`
	// A list of group ids to associate this user with. Comma delimited.
	GroupIds string `url:"group_ids,omitempty" json:"group_ids,omitempty" path:"group_ids"`
	// Pre-calculated hash of the user's password. If supplied, this will be used to authenticate the user on first login. Supported hash methods are MD5, SHA1, and SHA256.
	ImportedPasswordHash string `url:"imported_password_hash,omitempty" json:"imported_password_hash,omitempty" path:"imported_password_hash"`
	// User password.
	Password string `url:"password,omitempty" json:"password,omitempty" path:"password"`
	// Optional, but if provided, we will ensure that it matches the value sent in `password`.
	PasswordConfirmation string `url:"password_confirmation,omitempty" json:"password_confirmation,omitempty" path:"password_confirmation"`
	// Signifies that the user has read all the announcements in the UI.
	AnnouncementsRead *bool `url:"announcements_read,omitempty" json:"announcements_read,omitempty" path:"announcements_read"`
	// AI Assistant Personality ID assigned directly to this user, if any.
	AiAssistantPersonalityId int64 `url:"ai_assistant_personality_id,omitempty" json:"ai_assistant_personality_id,omitempty" path:"ai_assistant_personality_id"`
	// A list of allowed IPs if applicable. Newline delimited.
	AllowedIps string `url:"allowed_ips,omitempty" json:"allowed_ips,omitempty" path:"allowed_ips"`
	// DEPRECATED: If `true`, the user can user create Bundles (aka Share Links). Use the bundle permission instead.
	AttachmentsPermission *bool `url:"attachments_permission,omitempty" json:"attachments_permission,omitempty" path:"attachments_permission"`
	// Scheduled Date/Time at which user will be deactivated.
	AuthenticateUntil *time.Time `url:"authenticate_until,omitempty" json:"authenticate_until,omitempty" path:"authenticate_until"`
	// How is this user authenticated? Values: password, sso, none, email_signup, password_with_imported_hash, password_and_ssh_key.
	AuthenticationMethod UserAuthenticationMethodEnum `url:"authentication_method,omitempty" json:"authentication_method,omitempty" path:"authentication_method"`
	// Allow this user to perform operations on the account, payments, and invoices?
	BillingPermission *bool `url:"billing_permission,omitempty" json:"billing_permission,omitempty" path:"billing_permission"`
	// Exempt this user from user lifecycle rules?
	BypassUserLifecycleRules *bool `url:"bypass_user_lifecycle_rules,omitempty" json:"bypass_user_lifecycle_rules,omitempty" path:"bypass_user_lifecycle_rules"`
	// Allow this user to skip site-wide IP blacklists?
	BypassSiteAllowedIps *bool `url:"bypass_site_allowed_ips,omitempty" json:"bypass_site_allowed_ips,omitempty" path:"bypass_site_allowed_ips"`
	// Can the user connect with WebDAV?
	DavPermission *bool `url:"dav_permission,omitempty" json:"dav_permission,omitempty" path:"dav_permission"`
	// Desktop Configuration Profile ID assigned directly to this user, if any.
	DesktopConfigurationProfileId int64 `url:"desktop_configuration_profile_id,omitempty" json:"desktop_configuration_profile_id,omitempty" path:"desktop_configuration_profile_id"`
	// Workspace ID the user should land in by default when more than one Workspace is available.
	DefaultWorkspaceId int64 `url:"default_workspace_id,omitempty" json:"default_workspace_id,omitempty" path:"default_workspace_id"`
	// Is user disabled? Disabled users cannot log in, and do not count for billing purposes. Users can be automatically disabled after an inactivity period via a Site setting or schedule to be deactivated after specific date.
	Disabled *bool `url:"disabled,omitempty" json:"disabled,omitempty" path:"disabled"`
	// File system layout. Values: site_root, user_root, partner_root, integration_centric, workspace_root.
	FilesystemLayout UserFilesystemLayoutEnum `url:"filesystem_layout,omitempty" json:"filesystem_layout,omitempty" path:"filesystem_layout"`
	// Can the user access with FTP/FTPS?
	FtpPermission *bool `url:"ftp_permission,omitempty" json:"ftp_permission,omitempty" path:"ftp_permission"`
	// Text to display to the user in the header of the UI.
	HeaderText string `url:"header_text,omitempty" json:"header_text,omitempty" path:"header_text"`
	// Integration Centric Profile ID assigned directly to this user, if any.
	IntegrationCentricProfileId int64 `url:"integration_centric_profile_id,omitempty" json:"integration_centric_profile_id,omitempty" path:"integration_centric_profile_id"`
	// Preferred language.
	Language string `url:"language,omitempty" json:"language,omitempty" path:"language"`
	// Hour of the day at which daily notifications should be sent. Can be in range 0 to 23.
	NotificationDailySendTime int64 `url:"notification_daily_send_time,omitempty" json:"notification_daily_send_time,omitempty" path:"notification_daily_send_time"`
	// User's full name.
	Name string `url:"name,omitempty" json:"name,omitempty" path:"name"`
	// User's company.
	Company string `url:"company,omitempty" json:"company,omitempty" path:"company"`
	// Any internal notes on the user.
	Notes string `url:"notes,omitempty" json:"notes,omitempty" path:"notes"`
	// Enable integration with Office for the web?
	OfficeIntegrationEnabled *bool `url:"office_integration_enabled,omitempty" json:"office_integration_enabled,omitempty" path:"office_integration_enabled"`
	// Is this user a Partner administrator?
	PartnerAdmin *bool `url:"partner_admin,omitempty" json:"partner_admin,omitempty" path:"partner_admin"`
	// Partner ID if this user belongs to a Partner.
	PartnerId int64 `url:"partner_id,omitempty" json:"partner_id,omitempty" path:"partner_id"`
	// Number of days to allow user to use the same password.
	PasswordValidityDays int64 `url:"password_validity_days,omitempty" json:"password_validity_days,omitempty" path:"password_validity_days"`
	// Primary group ID for Group Admin scoping.
	PrimaryGroupId int64 `url:"primary_group_id,omitempty" json:"primary_group_id,omitempty" path:"primary_group_id"`
	// Is the user an allowed to view all (non-billing) site configuration for this site?
	ReadonlySiteAdmin *bool `url:"readonly_site_admin,omitempty" json:"readonly_site_admin,omitempty" path:"readonly_site_admin"`
	// Deprecated. Use notify_on_all_site_warnings and granular failure notification preferences instead.
	ReceiveAdminAlerts *bool `url:"receive_admin_alerts,omitempty" json:"receive_admin_alerts,omitempty" path:"receive_admin_alerts"`
	// Should the user receive site warnings via email?
	NotifyOnAllSiteWarnings *bool `url:"notify_on_all_site_warnings,omitempty" json:"notify_on_all_site_warnings,omitempty" path:"notify_on_all_site_warnings"`
	// Should the user receive sso/scim/ldap configuration/sync failures via email?
	NotifyOnAllSsoFailures *bool `url:"notify_on_all_sso_failures,omitempty" json:"notify_on_all_sso_failures,omitempty" path:"notify_on_all_sso_failures"`
	// Should the user receive user security events via email?
	NotifyOnAllUserSecurityEvents *bool `url:"notify_on_all_user_security_events,omitempty" json:"notify_on_all_user_security_events,omitempty" path:"notify_on_all_user_security_events"`
	// Should the user receive pending work failures via email?
	NotifyOnAllPendingWorkFailures *bool `url:"notify_on_all_pending_work_failures,omitempty" json:"notify_on_all_pending_work_failures,omitempty" path:"notify_on_all_pending_work_failures"`
	// Should the user receive siem failures via email?
	NotifyOnAllSiemHttpDestinationFailures *bool `url:"notify_on_all_siem_http_destination_failures,omitempty" json:"notify_on_all_siem_http_destination_failures,omitempty" path:"notify_on_all_siem_http_destination_failures"`
	// Should the user receive sync failures via email?
	NotifyOnAllSyncFailures *bool `url:"notify_on_all_sync_failures,omitempty" json:"notify_on_all_sync_failures,omitempty" path:"notify_on_all_sync_failures"`
	// Should the user receive automation failures via email?
	NotifyOnAllAutomationFailures *bool `url:"notify_on_all_automation_failures,omitempty" json:"notify_on_all_automation_failures,omitempty" path:"notify_on_all_automation_failures"`
	// Should the user receive expectation failures and misses via email?
	NotifyOnAllExpectationFailures *bool `url:"notify_on_all_expectation_failures,omitempty" json:"notify_on_all_expectation_failures,omitempty" path:"notify_on_all_expectation_failures"`
	// Require user to login by specified date otherwise it will be disabled.
	RequireLoginBy *time.Time `url:"require_login_by,omitempty" json:"require_login_by,omitempty" path:"require_login_by"`
	// Is a password change required upon next user login?
	RequirePasswordChange *bool `url:"require_password_change,omitempty" json:"require_password_change,omitempty" path:"require_password_change"`
	// ID of the internal Group responsible for this Partner User, overriding the Partner default.
	ResponsibleGroupId int64 `url:"responsible_group_id,omitempty" json:"responsible_group_id,omitempty" path:"responsible_group_id"`
	// ID of the internal User responsible for this Partner User, overriding the Partner default.
	ResponsibleUserId int64 `url:"responsible_user_id,omitempty" json:"responsible_user_id,omitempty" path:"responsible_user_id"`
	// Can this user access the Web app, Desktop app, SDKs, or REST API? (All of these tools use the API internally, so this is one unified permission set.).
	RestapiPermission *bool `url:"restapi_permission,omitempty" json:"restapi_permission,omitempty" path:"restapi_permission"`
	// Can the user access the S3-compatible endpoint? Defaults to true.
	S3CompatibleEndpointPermission *bool `url:"s3_compatible_endpoint_permission,omitempty" json:"s3_compatible_endpoint_permission,omitempty" path:"s3_compatible_endpoint_permission"`
	// Does this user manage it's own credentials or is it a shared/bot user?
	SelfManaged *bool `url:"self_managed,omitempty" json:"self_managed,omitempty" path:"self_managed"`
	// Can the user access with SFTP?
	SftpPermission *bool `url:"sftp_permission,omitempty" json:"sftp_permission,omitempty" path:"sftp_permission"`
	// Is the user an administrator for this site?
	SiteAdmin *bool `url:"site_admin,omitempty" json:"site_admin,omitempty" path:"site_admin"`
	// Skip Welcome page in the UI?
	SkipWelcomeScreen *bool `url:"skip_welcome_screen,omitempty" json:"skip_welcome_screen,omitempty" path:"skip_welcome_screen"`
	// SSL required setting. Values: use_system_setting, always_require, never_require.
	SslRequired UserSslRequiredEnum `url:"ssl_required,omitempty" json:"ssl_required,omitempty" path:"ssl_required"`
	// SSO (Single Sign On) strategy ID for the user, if applicable.
	SsoStrategyId int64 `url:"sso_strategy_id,omitempty" json:"sso_strategy_id,omitempty" path:"sso_strategy_id"`
	// 2FA required setting. `use_system_setting` uses the site-wide setting, including SSO exemptions. `always_require` and `never_require` override the site-wide setting when user-level overrides are allowed. Values: use_system_setting, always_require, never_require.
	Require2fa UserRequire2faEnum `url:"require_2fa,omitempty" json:"require_2fa,omitempty" path:"require_2fa"`
	// Comma-separated list of Tags for this user. Tags are used for other features, such as UserLifecycleRules, which can target specific tags. Tags must only contain lowercase letters, numbers, and hyphens.
	Tags string `url:"tags,omitempty" json:"tags,omitempty" path:"tags"`
	// User time zone.
	TimeZone string `url:"time_zone,omitempty" json:"time_zone,omitempty" path:"time_zone"`
	// If filesystem layout is user_root, this path is the root path the user is fixed to for all interfaces. If the filesystem layout is site_root or partner_root, this acts as a root folder only for FTP and SFTP (SFTP applicability also requires a site-wide setting to be set). For partner_root layout, this path is relative to the Partner root folder for all callers and blank opts out of an additional protocol root. In this situation, this path is not applied to the API, Desktop, or Web interface.
	UserRoot string `url:"user_root,omitempty" json:"user_root,omitempty" path:"user_root"`
	// Home folder for FTP/SFTP. For users with the partner_root filesystem layout, this path is relative to the Partner root folder. In all other cases, it is an absolute path. Only applies to FTP and SFTP, and not any other interface.
	UserHome string `url:"user_home,omitempty" json:"user_home,omitempty" path:"user_home"`
	// Whether the user is an administrator of their own Custom Workspace. Does not reflect administration granted through Permissions.
	WorkspaceAdmin *bool `url:"workspace_admin,omitempty" json:"workspace_admin,omitempty" path:"workspace_admin"`
	// User's username. Required.
	Username string `url:"username" json:"username" path:"username"`
	// ID of the Workspace the user belongs to. 0 is the Default Workspace.
	WorkspaceId int64 `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
}

// UserUnlockParams contains the request parameters for POST /users/{id}/unlock.
//
// Unlock user who has been locked out due to failed logins
type UserUnlockParams struct {
	// User ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UserResendWelcomeEmailParams contains the request parameters for POST /users/{id}/resend_welcome_email.
//
// Resend user welcome email
type UserResendWelcomeEmailParams struct {
	// User ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UserUser2faResetParams contains the request parameters for POST /users/{id}/2fa/reset.
//
// Trigger 2FA Reset process for user who has lost access to their existing 2FA methods
type UserUser2faResetParams struct {
	// User ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UserUpdateParams contains the request parameters for PATCH /users/{id}.
type UserUpdateParams struct {
	// User ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// An image file for your user avatar.
	AvatarFile io.Writer `url:"avatar_file,omitempty" json:"avatar_file,omitempty" path:"avatar_file"`
	// AvatarFileAttachment sends file content as avatar_file; see lib.JSONAttachment. Leave AvatarFile nil when you set it.
	AvatarFileAttachment *lib.JSONAttachment `json:"-" url:"-" path:"-" attachment:"avatar_file"`
	// If true, the avatar will be deleted.
	AvatarDelete *bool `url:"avatar_delete,omitempty" json:"avatar_delete,omitempty" path:"avatar_delete"`
	// Used for changing a password on an existing user.
	ChangePassword string `url:"change_password,omitempty" json:"change_password,omitempty" path:"change_password"`
	// Optional, but if provided, we will ensure that it matches the value sent in `change_password`.
	ChangePasswordConfirmation string `url:"change_password_confirmation,omitempty" json:"change_password_confirmation,omitempty" path:"change_password_confirmation"`
	// User's email.
	Email string `url:"email,omitempty" json:"email,omitempty" path:"email"`
	// Permission to grant on the User Root upon creation, or on the Workspace root when the User Root is blank. For the Default Workspace, a blank User Root uses the site root. `admin` on a Custom Workspace's root grants Workspace Administrator access; on a subfolder it grants Folder Admin access. Can be blank or `admin`, `full`, `read`, `write`, `list`, `read+write`, `list+write`, or `history`.
	GrantPermission string `url:"grant_permission,omitempty" json:"grant_permission,omitempty" path:"grant_permission"`
	// Group ID to associate this user with.
	GroupId int64 `url:"group_id,omitempty" json:"group_id,omitempty" path:"group_id"`
	// A list of group ids to associate this user with. Comma delimited.
	GroupIds string `url:"group_ids,omitempty" json:"group_ids,omitempty" path:"group_ids"`
	// Pre-calculated hash of the user's password. If supplied, this will be used to authenticate the user on first login. Supported hash methods are MD5, SHA1, and SHA256.
	ImportedPasswordHash string `url:"imported_password_hash,omitempty" json:"imported_password_hash,omitempty" path:"imported_password_hash"`
	// User password.
	Password string `url:"password,omitempty" json:"password,omitempty" path:"password"`
	// Optional, but if provided, we will ensure that it matches the value sent in `password`.
	PasswordConfirmation string `url:"password_confirmation,omitempty" json:"password_confirmation,omitempty" path:"password_confirmation"`
	// Signifies that the user has read all the announcements in the UI.
	AnnouncementsRead *bool `url:"announcements_read,omitempty" json:"announcements_read,omitempty" path:"announcements_read"`
	// AI Assistant Personality ID assigned directly to this user, if any.
	AiAssistantPersonalityId int64 `url:"ai_assistant_personality_id,omitempty" json:"ai_assistant_personality_id,omitempty" path:"ai_assistant_personality_id"`
	// A list of allowed IPs if applicable. Newline delimited.
	AllowedIps string `url:"allowed_ips,omitempty" json:"allowed_ips,omitempty" path:"allowed_ips"`
	// DEPRECATED: If `true`, the user can user create Bundles (aka Share Links). Use the bundle permission instead.
	AttachmentsPermission *bool `url:"attachments_permission,omitempty" json:"attachments_permission,omitempty" path:"attachments_permission"`
	// Scheduled Date/Time at which user will be deactivated.
	AuthenticateUntil *time.Time `url:"authenticate_until,omitempty" json:"authenticate_until,omitempty" path:"authenticate_until"`
	// How is this user authenticated? Values: password, sso, none, email_signup, password_with_imported_hash, password_and_ssh_key.
	AuthenticationMethod UserAuthenticationMethodEnum `url:"authentication_method,omitempty" json:"authentication_method,omitempty" path:"authentication_method"`
	// Allow this user to perform operations on the account, payments, and invoices?
	BillingPermission *bool `url:"billing_permission,omitempty" json:"billing_permission,omitempty" path:"billing_permission"`
	// Exempt this user from user lifecycle rules?
	BypassUserLifecycleRules *bool `url:"bypass_user_lifecycle_rules,omitempty" json:"bypass_user_lifecycle_rules,omitempty" path:"bypass_user_lifecycle_rules"`
	// Allow this user to skip site-wide IP blacklists?
	BypassSiteAllowedIps *bool `url:"bypass_site_allowed_ips,omitempty" json:"bypass_site_allowed_ips,omitempty" path:"bypass_site_allowed_ips"`
	// Can the user connect with WebDAV?
	DavPermission *bool `url:"dav_permission,omitempty" json:"dav_permission,omitempty" path:"dav_permission"`
	// Desktop Configuration Profile ID assigned directly to this user, if any.
	DesktopConfigurationProfileId int64 `url:"desktop_configuration_profile_id,omitempty" json:"desktop_configuration_profile_id,omitempty" path:"desktop_configuration_profile_id"`
	// Workspace ID the user should land in by default when more than one Workspace is available.
	DefaultWorkspaceId int64 `url:"default_workspace_id,omitempty" json:"default_workspace_id,omitempty" path:"default_workspace_id"`
	// Is user disabled? Disabled users cannot log in, and do not count for billing purposes. Users can be automatically disabled after an inactivity period via a Site setting or schedule to be deactivated after specific date.
	Disabled *bool `url:"disabled,omitempty" json:"disabled,omitempty" path:"disabled"`
	// File system layout. Values: site_root, user_root, partner_root, integration_centric, workspace_root.
	FilesystemLayout UserFilesystemLayoutEnum `url:"filesystem_layout,omitempty" json:"filesystem_layout,omitempty" path:"filesystem_layout"`
	// Can the user access with FTP/FTPS?
	FtpPermission *bool `url:"ftp_permission,omitempty" json:"ftp_permission,omitempty" path:"ftp_permission"`
	// Text to display to the user in the header of the UI.
	HeaderText string `url:"header_text,omitempty" json:"header_text,omitempty" path:"header_text"`
	// Integration Centric Profile ID assigned directly to this user, if any.
	IntegrationCentricProfileId int64 `url:"integration_centric_profile_id,omitempty" json:"integration_centric_profile_id,omitempty" path:"integration_centric_profile_id"`
	// Preferred language.
	Language string `url:"language,omitempty" json:"language,omitempty" path:"language"`
	// Hour of the day at which daily notifications should be sent. Can be in range 0 to 23.
	NotificationDailySendTime int64 `url:"notification_daily_send_time,omitempty" json:"notification_daily_send_time,omitempty" path:"notification_daily_send_time"`
	// User's full name.
	Name string `url:"name,omitempty" json:"name,omitempty" path:"name"`
	// User's company.
	Company string `url:"company,omitempty" json:"company,omitempty" path:"company"`
	// Any internal notes on the user.
	Notes string `url:"notes,omitempty" json:"notes,omitempty" path:"notes"`
	// Enable integration with Office for the web?
	OfficeIntegrationEnabled *bool `url:"office_integration_enabled,omitempty" json:"office_integration_enabled,omitempty" path:"office_integration_enabled"`
	// Is this user a Partner administrator?
	PartnerAdmin *bool `url:"partner_admin,omitempty" json:"partner_admin,omitempty" path:"partner_admin"`
	// Partner ID if this user belongs to a Partner.
	PartnerId int64 `url:"partner_id,omitempty" json:"partner_id,omitempty" path:"partner_id"`
	// Number of days to allow user to use the same password.
	PasswordValidityDays int64 `url:"password_validity_days,omitempty" json:"password_validity_days,omitempty" path:"password_validity_days"`
	// Primary group ID for Group Admin scoping.
	PrimaryGroupId int64 `url:"primary_group_id,omitempty" json:"primary_group_id,omitempty" path:"primary_group_id"`
	// Is the user an allowed to view all (non-billing) site configuration for this site?
	ReadonlySiteAdmin *bool `url:"readonly_site_admin,omitempty" json:"readonly_site_admin,omitempty" path:"readonly_site_admin"`
	// Deprecated. Use notify_on_all_site_warnings and granular failure notification preferences instead.
	ReceiveAdminAlerts *bool `url:"receive_admin_alerts,omitempty" json:"receive_admin_alerts,omitempty" path:"receive_admin_alerts"`
	// Should the user receive site warnings via email?
	NotifyOnAllSiteWarnings *bool `url:"notify_on_all_site_warnings,omitempty" json:"notify_on_all_site_warnings,omitempty" path:"notify_on_all_site_warnings"`
	// Should the user receive sso/scim/ldap configuration/sync failures via email?
	NotifyOnAllSsoFailures *bool `url:"notify_on_all_sso_failures,omitempty" json:"notify_on_all_sso_failures,omitempty" path:"notify_on_all_sso_failures"`
	// Should the user receive user security events via email?
	NotifyOnAllUserSecurityEvents *bool `url:"notify_on_all_user_security_events,omitempty" json:"notify_on_all_user_security_events,omitempty" path:"notify_on_all_user_security_events"`
	// Should the user receive pending work failures via email?
	NotifyOnAllPendingWorkFailures *bool `url:"notify_on_all_pending_work_failures,omitempty" json:"notify_on_all_pending_work_failures,omitempty" path:"notify_on_all_pending_work_failures"`
	// Should the user receive siem failures via email?
	NotifyOnAllSiemHttpDestinationFailures *bool `url:"notify_on_all_siem_http_destination_failures,omitempty" json:"notify_on_all_siem_http_destination_failures,omitempty" path:"notify_on_all_siem_http_destination_failures"`
	// Should the user receive sync failures via email?
	NotifyOnAllSyncFailures *bool `url:"notify_on_all_sync_failures,omitempty" json:"notify_on_all_sync_failures,omitempty" path:"notify_on_all_sync_failures"`
	// Should the user receive automation failures via email?
	NotifyOnAllAutomationFailures *bool `url:"notify_on_all_automation_failures,omitempty" json:"notify_on_all_automation_failures,omitempty" path:"notify_on_all_automation_failures"`
	// Should the user receive expectation failures and misses via email?
	NotifyOnAllExpectationFailures *bool `url:"notify_on_all_expectation_failures,omitempty" json:"notify_on_all_expectation_failures,omitempty" path:"notify_on_all_expectation_failures"`
	// Require user to login by specified date otherwise it will be disabled.
	RequireLoginBy *time.Time `url:"require_login_by,omitempty" json:"require_login_by,omitempty" path:"require_login_by"`
	// Is a password change required upon next user login?
	RequirePasswordChange *bool `url:"require_password_change,omitempty" json:"require_password_change,omitempty" path:"require_password_change"`
	// ID of the internal Group responsible for this Partner User, overriding the Partner default.
	ResponsibleGroupId int64 `url:"responsible_group_id,omitempty" json:"responsible_group_id,omitempty" path:"responsible_group_id"`
	// ID of the internal User responsible for this Partner User, overriding the Partner default.
	ResponsibleUserId int64 `url:"responsible_user_id,omitempty" json:"responsible_user_id,omitempty" path:"responsible_user_id"`
	// Can this user access the Web app, Desktop app, SDKs, or REST API? (All of these tools use the API internally, so this is one unified permission set.).
	RestapiPermission *bool `url:"restapi_permission,omitempty" json:"restapi_permission,omitempty" path:"restapi_permission"`
	// Can the user access the S3-compatible endpoint? Defaults to true.
	S3CompatibleEndpointPermission *bool `url:"s3_compatible_endpoint_permission,omitempty" json:"s3_compatible_endpoint_permission,omitempty" path:"s3_compatible_endpoint_permission"`
	// Does this user manage it's own credentials or is it a shared/bot user?
	SelfManaged *bool `url:"self_managed,omitempty" json:"self_managed,omitempty" path:"self_managed"`
	// Can the user access with SFTP?
	SftpPermission *bool `url:"sftp_permission,omitempty" json:"sftp_permission,omitempty" path:"sftp_permission"`
	// Is the user an administrator for this site?
	SiteAdmin *bool `url:"site_admin,omitempty" json:"site_admin,omitempty" path:"site_admin"`
	// Skip Welcome page in the UI?
	SkipWelcomeScreen *bool `url:"skip_welcome_screen,omitempty" json:"skip_welcome_screen,omitempty" path:"skip_welcome_screen"`
	// SSL required setting. Values: use_system_setting, always_require, never_require.
	SslRequired UserSslRequiredEnum `url:"ssl_required,omitempty" json:"ssl_required,omitempty" path:"ssl_required"`
	// SSO (Single Sign On) strategy ID for the user, if applicable.
	SsoStrategyId int64 `url:"sso_strategy_id,omitempty" json:"sso_strategy_id,omitempty" path:"sso_strategy_id"`
	// 2FA required setting. `use_system_setting` uses the site-wide setting, including SSO exemptions. `always_require` and `never_require` override the site-wide setting when user-level overrides are allowed. Values: use_system_setting, always_require, never_require.
	Require2fa UserRequire2faEnum `url:"require_2fa,omitempty" json:"require_2fa,omitempty" path:"require_2fa"`
	// Comma-separated list of Tags for this user. Tags are used for other features, such as UserLifecycleRules, which can target specific tags. Tags must only contain lowercase letters, numbers, and hyphens.
	Tags string `url:"tags,omitempty" json:"tags,omitempty" path:"tags"`
	// User time zone.
	TimeZone string `url:"time_zone,omitempty" json:"time_zone,omitempty" path:"time_zone"`
	// If filesystem layout is user_root, this path is the root path the user is fixed to for all interfaces. If the filesystem layout is site_root or partner_root, this acts as a root folder only for FTP and SFTP (SFTP applicability also requires a site-wide setting to be set). For partner_root layout, this path is relative to the Partner root folder for all callers and blank opts out of an additional protocol root. In this situation, this path is not applied to the API, Desktop, or Web interface.
	UserRoot string `url:"user_root,omitempty" json:"user_root,omitempty" path:"user_root"`
	// Home folder for FTP/SFTP. For users with the partner_root filesystem layout, this path is relative to the Partner root folder. In all other cases, it is an absolute path. Only applies to FTP and SFTP, and not any other interface.
	UserHome string `url:"user_home,omitempty" json:"user_home,omitempty" path:"user_home"`
	// Whether the user is an administrator of their own Custom Workspace. Does not reflect administration granted through Permissions.
	WorkspaceAdmin *bool `url:"workspace_admin,omitempty" json:"workspace_admin,omitempty" path:"workspace_admin"`
	// User's username.
	Username string `url:"username,omitempty" json:"username,omitempty" path:"username"`
	// Workspace ID. Only Site Administrators can change this field. Values supplied by Workspace Administrators, Group Administrators, or other non-Site Administrators using `/user` are ignored.
	WorkspaceId int64 `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
	// If true when changing authentication_method from `password` to `sso`, remove all two-factor methods. Ignored in all other cases.
	Clear2fa *bool `url:"clear_2fa,omitempty" json:"clear_2fa,omitempty" path:"clear_2fa"`
	// Required when assigning a Partner to an existing non-Partner user. If true, convert the user by assigning the partner_id provided.
	ConvertToPartnerUser *bool `url:"convert_to_partner_user,omitempty" json:"convert_to_partner_user,omitempty" path:"convert_to_partner_user"`
}

// UserDeleteParams contains the request parameters for DELETE /users/{id}.
type UserDeleteParams struct {
	// User ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// Provide a User ID here to transfer ownership of certain resources such as Automations and Share Links (Bundles) to that new user.
	NewOwnerId int64 `url:"new_owner_id,omitempty" json:"new_owner_id,omitempty" path:"new_owner_id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (u *User) UnmarshalJSON(data []byte) error {
	type user User
	var v user
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*u = User(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (u *UserCollection) UnmarshalJSON(data []byte) error {
	type users UserCollection
	var v users
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*u = UserCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (u *UserCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*u))
	for i, v := range *u {
		ret[i] = v
	}

	return &ret
}
