package files_sdk

import (
	"encoding/json"
	"io"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// Site is a Files.com API resource.
type Site struct {
	// Site Id.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Site name.
	Name string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	// Additional extensions that are considered text files.
	AdditionalTextFileTypes []string `json:"additional_text_file_types,omitempty" path:"additional_text_file_types,omitempty" url:"additional_text_file_types,omitempty"`
	// Availability settings for AI features. Each feature requires the site_admins, workspace_admins, folder_admins, and all_users keys. Optional selected_group_members defaults to false; when true, members of at least one group in group_ids get access regardless of the other options. All availability options are additive: any enabled option matching the user grants access. Optional group_ids is an array of integer IDs of active groups on this site, from any workspace. Omitted or empty group_ids grants no access through selected_group_members and does not affect other options. Disabling all AI features overrides these settings.
	AiFeatureAvailability interface{} `json:"ai_feature_availability,omitempty" path:"ai_feature_availability,omitempty" url:"ai_feature_availability,omitempty"`
	// Is SMS two factor authentication allowed?
	Allowed2faMethodSms *bool `json:"allowed_2fa_method_sms,omitempty" path:"allowed_2fa_method_sms,omitempty" url:"allowed_2fa_method_sms,omitempty"`
	// Is TOTP two factor authentication allowed?
	Allowed2faMethodTotp *bool `json:"allowed_2fa_method_totp,omitempty" path:"allowed_2fa_method_totp,omitempty" url:"allowed_2fa_method_totp,omitempty"`
	// Is WebAuthn two factor authentication allowed?
	Allowed2faMethodWebauthn *bool `json:"allowed_2fa_method_webauthn,omitempty" path:"allowed_2fa_method_webauthn,omitempty" url:"allowed_2fa_method_webauthn,omitempty"`
	// Is yubikey two factor authentication allowed?
	Allowed2faMethodYubi *bool `json:"allowed_2fa_method_yubi,omitempty" path:"allowed_2fa_method_yubi,omitempty" url:"allowed_2fa_method_yubi,omitempty"`
	// Is OTP via email two factor authentication allowed?
	Allowed2faMethodEmail *bool `json:"allowed_2fa_method_email,omitempty" path:"allowed_2fa_method_email,omitempty" url:"allowed_2fa_method_email,omitempty"`
	// Is OTP via static codes for two factor authentication allowed?
	Allowed2faMethodStatic *bool `json:"allowed_2fa_method_static,omitempty" path:"allowed_2fa_method_static,omitempty" url:"allowed_2fa_method_static,omitempty"`
	// Are users allowed to configure their two factor authentication to be bypassed for FTP/SFTP/WebDAV?
	Allowed2faMethodBypassForFtpSftpDav *bool `json:"allowed_2fa_method_bypass_for_ftp_sftp_dav,omitempty" path:"allowed_2fa_method_bypass_for_ftp_sftp_dav,omitempty" url:"allowed_2fa_method_bypass_for_ftp_sftp_dav,omitempty"`
	// User ID for the main site administrator.
	AdminUserId int64 `json:"admin_user_id,omitempty" path:"admin_user_id,omitempty" url:"admin_user_id,omitempty"`
	// Allow admins to bypass the locked subfolders setting.
	AdminsBypassLockedSubfolders *bool `json:"admins_bypass_locked_subfolders,omitempty" path:"admins_bypass_locked_subfolders,omitempty" url:"admins_bypass_locked_subfolders,omitempty"`
	// Are manual Bundle names allowed?
	AllowBundleNames *bool `json:"allow_bundle_names,omitempty" path:"allow_bundle_names,omitempty" url:"allow_bundle_names,omitempty"`
	// Allow the site-wide two-factor authentication requirement to be overriden on a per-user-basis?
	AllowUserLevel2faOverride *bool `json:"allow_user_level_2fa_override,omitempty" path:"allow_user_level_2fa_override,omitempty" url:"allow_user_level_2fa_override,omitempty"`
	// Allow the site-wide allowed IP restriction to be overriden on a per-user-basis?
	AllowUserLevelAllowedIpOverride *bool `json:"allow_user_level_allowed_ip_override,omitempty" path:"allow_user_level_allowed_ip_override,omitempty" url:"allow_user_level_allowed_ip_override,omitempty"`
	// Allow the site-wide FTP SSL requirement to be overriden on a per-user-basis?
	AllowUserLevelSslOverride *bool `json:"allow_user_level_ssl_override,omitempty" path:"allow_user_level_ssl_override,omitempty" url:"allow_user_level_ssl_override,omitempty"`
	// Comma separated list of allowed Country codes.
	AllowedCountries string `json:"allowed_countries,omitempty" path:"allowed_countries,omitempty" url:"allowed_countries,omitempty"`
	// List of allowed IP addresses.
	AllowedIps string `json:"allowed_ips,omitempty" path:"allowed_ips,omitempty" url:"allowed_ips,omitempty"`
	// Require Parent Site users entering this Child Site to match its allowed IP addresses in addition to their existing IP restrictions. Applies to all Workspaces. Defaults to false.
	EnforceParentSiteUserIpRestrictions *bool `json:"enforce_parent_site_user_ip_restrictions,omitempty" path:"enforce_parent_site_user_ip_restrictions,omitempty" url:"enforce_parent_site_user_ip_restrictions,omitempty"`
	// Create parent directories if they do not exist during uploads? This is primarily used to work around broken upload clients that assume servers will perform this step.
	AlwaysMkdirParents *bool `json:"always_mkdir_parents,omitempty" path:"always_mkdir_parents,omitempty" url:"always_mkdir_parents,omitempty"`
	// Number of days to retain AS2 messages (incoming and outgoing).
	As2MessageRetentionDays int64 `json:"as2_message_retention_days,omitempty" path:"as2_message_retention_days,omitempty" url:"as2_message_retention_days,omitempty"`
	// If false, rename conflicting files instead of asking for overwrite confirmation. Only applies to web interface.
	AskAboutOverwrites *bool `json:"ask_about_overwrites,omitempty" path:"ask_about_overwrites,omitempty" url:"ask_about_overwrites,omitempty"`
	// Do Bundle owners receive activity notifications?
	BundleActivityNotifications string `json:"bundle_activity_notifications,omitempty" path:"bundle_activity_notifications,omitempty" url:"bundle_activity_notifications,omitempty"`
	// Site-wide Bundle expiration in days.
	BundleExpiration int64 `json:"bundle_expiration,omitempty" path:"bundle_expiration,omitempty" url:"bundle_expiration,omitempty"`
	// Custom error message to show when bundle is not found.
	BundleNotFoundMessage string `json:"bundle_not_found_message,omitempty" path:"bundle_not_found_message,omitempty" url:"bundle_not_found_message,omitempty"`
	// Do Bundles require password protection?
	BundlePasswordRequired *bool `json:"bundle_password_required,omitempty" path:"bundle_password_required,omitempty" url:"bundle_password_required,omitempty"`
	// If true, new Share Links created by a user with a primary group will default to that group as owner.
	BundlesDefaultOwnedByPrimaryGroup *bool `json:"bundles_default_owned_by_primary_group,omitempty" path:"bundles_default_owned_by_primary_group,omitempty" url:"bundles_default_owned_by_primary_group,omitempty"`
	// List of email domains to disallow when entering a Bundle/Inbox recipients.
	BundleRecipientBlacklistDomains []string `json:"bundle_recipient_blacklist_domains,omitempty" path:"bundle_recipient_blacklist_domains,omitempty" url:"bundle_recipient_blacklist_domains,omitempty"`
	// Disallow free email domains for Bundle/Inbox recipients?
	BundleRecipientBlacklistFreeEmailDomains *bool `json:"bundle_recipient_blacklist_free_email_domains,omitempty" path:"bundle_recipient_blacklist_free_email_domains,omitempty" url:"bundle_recipient_blacklist_free_email_domains,omitempty"`
	// Do Bundle owners receive registration notification?
	BundleRegistrationNotifications string `json:"bundle_registration_notifications,omitempty" path:"bundle_registration_notifications,omitempty" url:"bundle_registration_notifications,omitempty"`
	// Do Bundles require registration?
	BundleRequireRegistration *bool `json:"bundle_require_registration,omitempty" path:"bundle_require_registration,omitempty" url:"bundle_require_registration,omitempty"`
	// Do Bundles require recipients for sharing?
	BundleRequireShareRecipient *bool `json:"bundle_require_share_recipient,omitempty" path:"bundle_require_share_recipient,omitempty" url:"bundle_require_share_recipient,omitempty"`
	// If true, new Share Links must send a one-time password to the recipient when they register. Requires bundle_require_share_recipient and cannot be enabled with bundle_password_required.
	BundleSendOneTimePasswordToRecipientAtRegistration *bool `json:"bundle_send_one_time_password_to_recipient_at_registration,omitempty" path:"bundle_send_one_time_password_to_recipient_at_registration,omitempty" url:"bundle_send_one_time_password_to_recipient_at_registration,omitempty"`
	// Do Bundles require internal notes?
	BundleRequireNote *bool `json:"bundle_require_note,omitempty" path:"bundle_require_note,omitempty" url:"bundle_require_note,omitempty"`
	// Do Bundle creators receive receipts of invitations?
	BundleSendSharedReceipts *bool `json:"bundle_send_shared_receipts,omitempty" path:"bundle_send_shared_receipts,omitempty" url:"bundle_send_shared_receipts,omitempty"`
	// Do Bundle uploaders receive upload confirmation notifications?
	BundleUploadReceiptNotifications string `json:"bundle_upload_receipt_notifications,omitempty" path:"bundle_upload_receipt_notifications,omitempty" url:"bundle_upload_receipt_notifications,omitempty"`
	// Preview watermark image applied to all bundle items.
	BundleWatermarkAttachment Image `json:"bundle_watermark_attachment,omitempty" path:"bundle_watermark_attachment,omitempty" url:"bundle_watermark_attachment,omitempty"`
	// Preview watermark settings applied to all bundle items. Uses the same keys as Behavior.value.
	BundleWatermarkValue interface{} `json:"bundle_watermark_value,omitempty" path:"bundle_watermark_value,omitempty" url:"bundle_watermark_value,omitempty"`
	// Calculate CRC32 checksums for files?
	CalculateFileChecksumsCrc32 *bool `json:"calculate_file_checksums_crc32,omitempty" path:"calculate_file_checksums_crc32,omitempty" url:"calculate_file_checksums_crc32,omitempty"`
	// Calculate MD5 checksums for files?
	CalculateFileChecksumsMd5 *bool `json:"calculate_file_checksums_md5,omitempty" path:"calculate_file_checksums_md5,omitempty" url:"calculate_file_checksums_md5,omitempty"`
	// Calculate SHA1 checksums for files?
	CalculateFileChecksumsSha1 *bool `json:"calculate_file_checksums_sha1,omitempty" path:"calculate_file_checksums_sha1,omitempty" url:"calculate_file_checksums_sha1,omitempty"`
	// Calculate SHA256 checksums for files?
	CalculateFileChecksumsSha256 *bool `json:"calculate_file_checksums_sha256,omitempty" path:"calculate_file_checksums_sha256,omitempty" url:"calculate_file_checksums_sha256,omitempty"`
	// If true, administrators on this site's Child Sites, including Folder Admins, can list this site's Default Workspace users and groups and grant them folder permissions and notifications on Child Site paths.
	ChildSiteFolderAdminsCanSeeParentSiteDirectory *bool `json:"child_site_folder_admins_can_see_parent_site_directory,omitempty" path:"child_site_folder_admins_can_see_parent_site_directory,omitempty" url:"child_site_folder_admins_can_see_parent_site_directory,omitempty"`
	// If true, Workspace Administrators and Folder Admins in a Custom Workspace can list the Default Workspace users and groups that have permissions in that Workspace and choose them as Email Notification recipients.
	WorkspaceFolderAdminsCanSeeDefaultWorkspaceDirectory *bool `json:"workspace_folder_admins_can_see_default_workspace_directory,omitempty" path:"workspace_folder_admins_can_see_default_workspace_directory,omitempty" url:"workspace_folder_admins_can_see_default_workspace_directory,omitempty"`
	// Require email authentication, virus, and spam checks for incoming emails to Inboxes and Incoming Email Automations in every Workspace on this site?
	UploadsViaEmailAuthentication *bool `json:"uploads_via_email_authentication,omitempty" path:"uploads_via_email_authentication,omitempty" url:"uploads_via_email_authentication,omitempty"`
	// Page link and button color.
	Color2Left string `json:"color2_left,omitempty" path:"color2_left,omitempty" url:"color2_left,omitempty"`
	// Top bar link color.
	Color2Link string `json:"color2_link,omitempty" path:"color2_link,omitempty" url:"color2_link,omitempty"`
	// Page link and button color.
	Color2Text string `json:"color2_text,omitempty" path:"color2_text,omitempty" url:"color2_text,omitempty"`
	// Top bar background color.
	Color2Top string `json:"color2_top,omitempty" path:"color2_top,omitempty" url:"color2_top,omitempty"`
	// Top bar text color.
	Color2TopText string `json:"color2_top_text,omitempty" path:"color2_top_text,omitempty" url:"color2_top_text,omitempty"`
	// Site main contact name.
	ContactName string `json:"contact_name,omitempty" path:"contact_name,omitempty" url:"contact_name,omitempty"`
	// Time this site was created.
	CreatedAt *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	// Preferred currency.
	Currency string `json:"currency,omitempty" path:"currency,omitempty" url:"currency,omitempty"`
	// Is this site using a custom namespace for users?
	CustomNamespace *bool `json:"custom_namespace,omitempty" path:"custom_namespace,omitempty" url:"custom_namespace,omitempty"`
	// Is WebDAV enabled?
	DavEnabled *bool `json:"dav_enabled,omitempty" path:"dav_enabled,omitempty" url:"dav_enabled,omitempty"`
	// Use user FTP roots also for WebDAV?
	DavUserRootEnabled *bool `json:"dav_user_root_enabled,omitempty" path:"dav_user_root_enabled,omitempty" url:"dav_user_root_enabled,omitempty"`
	// Number of days to keep deleted files.
	DaysToRetainBackups int64 `json:"days_to_retain_backups,omitempty" path:"days_to_retain_backups,omitempty" url:"days_to_retain_backups,omitempty"`
	// If true, allow public viewers of Bundles with full permissions to use document editing integrations.
	DocumentEditsInBundleAllowed *bool `json:"document_edits_in_bundle_allowed,omitempty" path:"document_edits_in_bundle_allowed,omitempty" url:"document_edits_in_bundle_allowed,omitempty"`
	// Site default time zone.
	DefaultTimeZone string `json:"default_time_zone,omitempty" path:"default_time_zone,omitempty" url:"default_time_zone,omitempty"`
	// Is the desktop app enabled?
	DesktopApp *bool `json:"desktop_app,omitempty" path:"desktop_app,omitempty" url:"desktop_app,omitempty"`
	// Is desktop app session IP pinning enabled?
	DesktopAppSessionIpPinning *bool `json:"desktop_app_session_ip_pinning,omitempty" path:"desktop_app_session_ip_pinning,omitempty" url:"desktop_app_session_ip_pinning,omitempty"`
	// Desktop app session lifetime (in hours).
	DesktopAppSessionLifetime int64 `json:"desktop_app_session_lifetime,omitempty" path:"desktop_app_session_lifetime,omitempty" url:"desktop_app_session_lifetime,omitempty"`
	// Use legacy checksums mode?
	LegacyChecksumsMode *bool `json:"legacy_checksums_mode,omitempty" path:"legacy_checksums_mode,omitempty" url:"legacy_checksums_mode,omitempty"`
	// If true, we will migrate all remote server syncs to the new Sync model.
	MigrateRemoteServerSyncToSync *bool `json:"migrate_remote_server_sync_to_sync,omitempty" path:"migrate_remote_server_sync_to_sync,omitempty" url:"migrate_remote_server_sync_to_sync,omitempty"`
	// Is OAuth DCR (dynamic client registration) for MCP enabled?
	McpDcrEnabled *bool `json:"mcp_dcr_enabled,omitempty" path:"mcp_dcr_enabled,omitempty" url:"mcp_dcr_enabled,omitempty"`
	// Is the mobile app enabled?
	MobileApp *bool `json:"mobile_app,omitempty" path:"mobile_app,omitempty" url:"mobile_app,omitempty"`
	// Is mobile app session IP pinning enabled?
	MobileAppSessionIpPinning *bool `json:"mobile_app_session_ip_pinning,omitempty" path:"mobile_app_session_ip_pinning,omitempty" url:"mobile_app_session_ip_pinning,omitempty"`
	// Mobile app session lifetime (in hours).
	MobileAppSessionLifetime int64 `json:"mobile_app_session_lifetime,omitempty" path:"mobile_app_session_lifetime,omitempty" url:"mobile_app_session_lifetime,omitempty"`
	// Comma separated list of disallowed Country codes.
	DisallowedCountries string `json:"disallowed_countries,omitempty" path:"disallowed_countries,omitempty" url:"disallowed_countries,omitempty"`
	// If true, all AI features are disabled for this site.
	DisableAllAiFeatures *bool `json:"disable_all_ai_features,omitempty" path:"disable_all_ai_features,omitempty" url:"disable_all_ai_features,omitempty"`
	// If set, Files.com will not set the CAA records required to generate future SSL certificates for this domain.
	DisableFilesCertificateGeneration *bool `json:"disable_files_certificate_generation,omitempty" path:"disable_files_certificate_generation,omitempty" url:"disable_files_certificate_generation,omitempty"`
	// Are notifications disabled?
	DisableNotifications *bool `json:"disable_notifications,omitempty" path:"disable_notifications,omitempty" url:"disable_notifications,omitempty"`
	// If true, file operations produce file_action.activity.v1 EventChannel events in all workspaces. Defaults to false; subscription filters and workspace access still control delivery.
	FileActivityEventsEnabled *bool `json:"file_activity_events_enabled,omitempty" path:"file_activity_events_enabled,omitempty" url:"file_activity_events_enabled,omitempty"`
	// Is password reset disabled?
	DisablePasswordReset *bool `json:"disable_password_reset,omitempty" path:"disable_password_reset,omitempty" url:"disable_password_reset,omitempty"`
	// Custom domain.
	Domain string `json:"domain,omitempty" path:"domain,omitempty" url:"domain,omitempty"`
	// Send HSTS (HTTP Strict Transport Security) header when visitors access the site via a custom domain?
	DomainHstsHeader *bool `json:"domain_hsts_header,omitempty" path:"domain_hsts_header,omitempty" url:"domain_hsts_header,omitempty"`
	// Letsencrypt chain to use when registering SSL Certificate for domain. No longer used as of 2026.
	DomainLetsencryptChain string `json:"domain_letsencrypt_chain,omitempty" path:"domain_letsencrypt_chain,omitempty" url:"domain_letsencrypt_chain,omitempty"`
	// Main email for this site.
	Email string `json:"email,omitempty" path:"email,omitempty" url:"email,omitempty"`
	// Are FedRAMP security restrictions enabled for this site?
	Fedramp *bool `json:"fedramp,omitempty" path:"fedramp,omitempty" url:"fedramp,omitempty"`
	// Allow other Files.com sites to use this site's API keys for native Files.com Remote Server connections? Defaults to true and applies to all keys and workspaces on this site. When false, new pairings and access through existing connections are rejected without revoking keys or deleting connections. Re-enabling permits access again with usable keys. Does not disable this site's connections to other sites, ordinary API access, or Connected Sites.
	FilesComRemoteServerEnabled *bool `json:"files_com_remote_server_enabled,omitempty" path:"files_com_remote_server_enabled,omitempty" url:"files_com_remote_server_enabled,omitempty"`
	// Is FTP enabled?
	FtpEnabled *bool `json:"ftp_enabled,omitempty" path:"ftp_enabled,omitempty" url:"ftp_enabled,omitempty"`
	// Reply-to email for this site.
	ReplyToEmail string `json:"reply_to_email,omitempty" path:"reply_to_email,omitempty" url:"reply_to_email,omitempty"`
	// If true, groups can be manually created / modified / deleted by Site Admins. Otherwise, groups can only be managed via your SSO provider.
	NonSsoGroupsAllowed *bool `json:"non_sso_groups_allowed,omitempty" path:"non_sso_groups_allowed,omitempty" url:"non_sso_groups_allowed,omitempty"`
	// If true, users can be manually created / modified / deleted by Site Admins. Otherwise, users can only be managed via your SSO provider.
	NonSsoUsersAllowed *bool `json:"non_sso_users_allowed,omitempty" path:"non_sso_users_allowed,omitempty" url:"non_sso_users_allowed,omitempty"`
	// If true, permissions for this site must be bound to a group (not a user).
	FolderPermissionsGroupsOnly *bool `json:"folder_permissions_groups_only,omitempty" path:"folder_permissions_groups_only,omitempty" url:"folder_permissions_groups_only,omitempty"`
	// Allow group admins to create users in their groups.
	GroupAdminsCanAddUsers *bool `json:"group_admins_can_add_users,omitempty" path:"group_admins_can_add_users,omitempty" url:"group_admins_can_add_users,omitempty"`
	// Allow group admins to add or remove existing users in their groups.
	GroupAdminsCanManageGroupMemberships *bool `json:"group_admins_can_manage_group_memberships,omitempty" path:"group_admins_can_manage_group_memberships,omitempty" url:"group_admins_can_manage_group_memberships,omitempty"`
	// Allow group admins to delete users in their groups.
	GroupAdminsCanDeleteUsers *bool `json:"group_admins_can_delete_users,omitempty" path:"group_admins_can_delete_users,omitempty" url:"group_admins_can_delete_users,omitempty"`
	// Allow group admins to enable or disable users in their groups.
	GroupAdminsCanEnableDisableUsers *bool `json:"group_admins_can_enable_disable_users,omitempty" path:"group_admins_can_enable_disable_users,omitempty" url:"group_admins_can_enable_disable_users,omitempty"`
	// Allow group admins to modify users in their groups.
	GroupAdminsCanModifyUsers *bool `json:"group_admins_can_modify_users,omitempty" path:"group_admins_can_modify_users,omitempty" url:"group_admins_can_modify_users,omitempty"`
	// Allow group admins to exempt users in their groups from lifecycle rules.
	GroupAdminsCanBypassUserLifecycleRules *bool `json:"group_admins_can_bypass_user_lifecycle_rules,omitempty" path:"group_admins_can_bypass_user_lifecycle_rules,omitempty" url:"group_admins_can_bypass_user_lifecycle_rules,omitempty"`
	// Allow group admins to reset passwords for users in their groups.
	GroupAdminsCanResetPasswords *bool `json:"group_admins_can_reset_passwords,omitempty" path:"group_admins_can_reset_passwords,omitempty" url:"group_admins_can_reset_passwords,omitempty"`
	// Allow group admins to set password authentication method.
	GroupAdminsCanSetUserPassword *bool `json:"group_admins_can_set_user_password,omitempty" path:"group_admins_can_set_user_password,omitempty" url:"group_admins_can_set_user_password,omitempty"`
	// Is there a signed HIPAA BAA between Files.com and this site?
	Hipaa *bool `json:"hipaa,omitempty" path:"hipaa,omitempty" url:"hipaa,omitempty"`
	// Branded icon 128x128.
	Icon128 Image `json:"icon128,omitempty" path:"icon128,omitempty" url:"icon128,omitempty"`
	// Branded icon 16x16.
	Icon16 Image `json:"icon16,omitempty" path:"icon16,omitempty" url:"icon16,omitempty"`
	// Branded icon 32x32.
	Icon32 Image `json:"icon32,omitempty" path:"icon32,omitempty" url:"icon32,omitempty"`
	// Branded icon 48x48.
	Icon48 Image `json:"icon48,omitempty" path:"icon48,omitempty" url:"icon48,omitempty"`
	// Can files be modified?
	ImmutableFilesSetAt *time.Time `json:"immutable_files_set_at,omitempty" path:"immutable_files_set_at,omitempty" url:"immutable_files_set_at,omitempty"`
	// Include password in emails to new users?
	IncludePasswordInWelcomeEmail *bool `json:"include_password_in_welcome_email,omitempty" path:"include_password_in_welcome_email,omitempty" url:"include_password_in_welcome_email,omitempty"`
	// Site default language.
	Language string `json:"language,omitempty" path:"language,omitempty" url:"language,omitempty"`
	// Base DN for looking up users in LDAP server.
	LdapBaseDn string `json:"ldap_base_dn,omitempty" path:"ldap_base_dn,omitempty" url:"ldap_base_dn,omitempty"`
	// Domain name that will be appended to usernames.
	LdapDomain string `json:"ldap_domain,omitempty" path:"ldap_domain,omitempty" url:"ldap_domain,omitempty"`
	// Main LDAP setting: is LDAP enabled?
	LdapEnabled *bool `json:"ldap_enabled,omitempty" path:"ldap_enabled,omitempty" url:"ldap_enabled,omitempty"`
	// Should we sync groups from LDAP server?
	LdapGroupAction string `json:"ldap_group_action,omitempty" path:"ldap_group_action,omitempty" url:"ldap_group_action,omitempty"`
	// Comma or newline separated list of group names (with optional wildcards) to exclude when syncing.
	LdapGroupExclusion string `json:"ldap_group_exclusion,omitempty" path:"ldap_group_exclusion,omitempty" url:"ldap_group_exclusion,omitempty"`
	// Comma or newline separated list of group names (with optional wildcards) to include when syncing.
	LdapGroupInclusion string `json:"ldap_group_inclusion,omitempty" path:"ldap_group_inclusion,omitempty" url:"ldap_group_inclusion,omitempty"`
	// LDAP host.
	LdapHost string `json:"ldap_host,omitempty" path:"ldap_host,omitempty" url:"ldap_host,omitempty"`
	// LDAP backup host.
	LdapHost2 string `json:"ldap_host_2,omitempty" path:"ldap_host_2,omitempty" url:"ldap_host_2,omitempty"`
	// LDAP backup host.
	LdapHost3 string `json:"ldap_host_3,omitempty" path:"ldap_host_3,omitempty" url:"ldap_host_3,omitempty"`
	// LDAP port.
	LdapPort int64 `json:"ldap_port,omitempty" path:"ldap_port,omitempty" url:"ldap_port,omitempty"`
	// Use secure LDAP?
	LdapSecure *bool `json:"ldap_secure,omitempty" path:"ldap_secure,omitempty" url:"ldap_secure,omitempty"`
	// How to validate the LDAP server certificate. `require_match` validates the certificate chain and hostname; `allow_any` disables certificate validation.
	LdapServerCertificate string `json:"ldap_server_certificate,omitempty" path:"ldap_server_certificate,omitempty" url:"ldap_server_certificate,omitempty"`
	// LDAP type.
	LdapType string `json:"ldap_type,omitempty" path:"ldap_type,omitempty" url:"ldap_type,omitempty"`
	// Should we sync users from LDAP server?
	LdapUserAction string `json:"ldap_user_action,omitempty" path:"ldap_user_action,omitempty" url:"ldap_user_action,omitempty"`
	// Comma or newline separated list of group names (with optional wildcards) - if provided, only users in these groups will be added or synced.
	LdapUserIncludeGroups string `json:"ldap_user_include_groups,omitempty" path:"ldap_user_include_groups,omitempty" url:"ldap_user_include_groups,omitempty"`
	// Username for signing in to LDAP server.
	LdapUsername string `json:"ldap_username,omitempty" path:"ldap_username,omitempty" url:"ldap_username,omitempty"`
	// LDAP username field.
	LdapUsernameField string `json:"ldap_username_field,omitempty" path:"ldap_username_field,omitempty" url:"ldap_username_field,omitempty"`
	// Login help text.
	LoginHelpText string `json:"login_help_text,omitempty" path:"login_help_text,omitempty" url:"login_help_text,omitempty"`
	// Branded logo.
	Logo Image `json:"logo,omitempty" path:"logo,omitempty" url:"logo,omitempty"`
	// Branded login page background.
	LoginPageBackgroundImage Image `json:"login_page_background_image,omitempty" path:"login_page_background_image,omitempty" url:"login_page_background_image,omitempty"`
	// Number of prior passwords to disallow.
	MaxPriorPasswords int64 `json:"max_prior_passwords,omitempty" path:"max_prior_passwords,omitempty" url:"max_prior_passwords,omitempty"`
	// Site settings managed by the parent site. Each key identifies a managed setting, with a null value. Non-secret effective values are returned in their own Site settings fields.
	ManagedSiteSettings interface{} `json:"managed_site_settings,omitempty" path:"managed_site_settings,omitempty" url:"managed_site_settings,omitempty"`
	// A message to show users when they connect via FTP or SFTP.
	MotdText string `json:"motd_text,omitempty" path:"motd_text,omitempty" url:"motd_text,omitempty"`
	// Show message to users connecting via FTP.
	MotdUseForFtp *bool `json:"motd_use_for_ftp,omitempty" path:"motd_use_for_ftp,omitempty" url:"motd_use_for_ftp,omitempty"`
	// Show message to users connecting via SFTP.
	MotdUseForSftp *bool `json:"motd_use_for_sftp,omitempty" path:"motd_use_for_sftp,omitempty" url:"motd_use_for_sftp,omitempty"`
	// Next billing amount.
	NextBillingAmount string `json:"next_billing_amount,omitempty" path:"next_billing_amount,omitempty" url:"next_billing_amount,omitempty"`
	// Next billing date.
	NextBillingDate string `json:"next_billing_date,omitempty" path:"next_billing_date,omitempty" url:"next_billing_date,omitempty"`
	// If true, allows users to use a document editing integration.
	OfficeIntegrationAvailable *bool `json:"office_integration_available,omitempty" path:"office_integration_available,omitempty" url:"office_integration_available,omitempty"`
	// Which document editing integration to support. Files.com Editor or Microsoft Office for the Web.
	OfficeIntegrationType string `json:"office_integration_type,omitempty" path:"office_integration_type,omitempty" url:"office_integration_type,omitempty"`
	// Link to scheduling a meeting with our Sales team.
	OncehubLink string `json:"oncehub_link,omitempty" path:"oncehub_link,omitempty" url:"oncehub_link,omitempty"`
	// Use servers in the USA only?
	OptOutGlobal *bool `json:"opt_out_global,omitempty" path:"opt_out_global,omitempty" url:"opt_out_global,omitempty"`
	// Is this site's billing overdue?
	Overdue *bool `json:"overdue,omitempty" path:"overdue,omitempty" url:"overdue,omitempty"`
	// Shortest password length for users.
	PasswordMinLength int64 `json:"password_min_length,omitempty" path:"password_min_length,omitempty" url:"password_min_length,omitempty"`
	// Require a letter in passwords?
	PasswordRequireLetter *bool `json:"password_require_letter,omitempty" path:"password_require_letter,omitempty" url:"password_require_letter,omitempty"`
	// Require lower and upper case letters in passwords?
	PasswordRequireMixed *bool `json:"password_require_mixed,omitempty" path:"password_require_mixed,omitempty" url:"password_require_mixed,omitempty"`
	// Require a number in passwords?
	PasswordRequireNumber *bool `json:"password_require_number,omitempty" path:"password_require_number,omitempty" url:"password_require_number,omitempty"`
	// Require special characters in password?
	PasswordRequireSpecial *bool `json:"password_require_special,omitempty" path:"password_require_special,omitempty" url:"password_require_special,omitempty"`
	// Require passwords that have not been previously breached? (see https://haveibeenpwned.com/).
	PasswordRequireUnbreached *bool `json:"password_require_unbreached,omitempty" path:"password_require_unbreached,omitempty" url:"password_require_unbreached,omitempty"`
	// Require bundles' passwords, and passwords for other items (inboxes, public shares, etc.) to conform to the same requirements as users' passwords?
	PasswordRequirementsApplyToBundles *bool `json:"password_requirements_apply_to_bundles,omitempty" path:"password_requirements_apply_to_bundles,omitempty" url:"password_requirements_apply_to_bundles,omitempty"`
	// Number of days password is valid.
	PasswordValidityDays int64 `json:"password_validity_days,omitempty" path:"password_validity_days,omitempty" url:"password_validity_days,omitempty"`
	// Site phone number.
	Phone string `json:"phone,omitempty" path:"phone,omitempty" url:"phone,omitempty"`
	// If true, we will ensure that all internal communications with any remote server are made through the primary region of the site. This setting overrides individual remote server settings.
	PinAllRemoteServersToSiteRegion *bool `json:"pin_all_remote_servers_to_site_region,omitempty" path:"pin_all_remote_servers_to_site_region,omitempty" url:"pin_all_remote_servers_to_site_region,omitempty"`
	// If true, we will prevent non-administrators from receiving any permissions directly on the root folder. This is commonly used to prevent the accidental application of permissions.
	PreventRootPermissionsForNonSiteAdmins *bool `json:"prevent_root_permissions_for_non_site_admins,omitempty" path:"prevent_root_permissions_for_non_site_admins,omitempty" url:"prevent_root_permissions_for_non_site_admins,omitempty"`
	// If true, protocol access permissions on users will be ignored, and only protocol access permissions set on Groups will be honored. Make sure that your current user is a member of a group with API permission when changing this value to avoid locking yourself out of your site.
	ProtocolAccessGroupsOnly *bool `json:"protocol_access_groups_only,omitempty" path:"protocol_access_groups_only,omitempty" url:"protocol_access_groups_only,omitempty"`
	// Require two-factor authentication for all users?
	Require2fa *bool `json:"require_2fa,omitempty" path:"require_2fa,omitempty" url:"require_2fa,omitempty"`
	// If true, only site admins may create, modify, or delete any behavior at the site root, or a skip that would disable one.
	RestrictRootFolderBehaviorsToSiteAdmins *bool `json:"restrict_root_folder_behaviors_to_site_admins,omitempty" path:"restrict_root_folder_behaviors_to_site_admins,omitempty" url:"restrict_root_folder_behaviors_to_site_admins,omitempty"`
	// If true, supported protective behaviors at the site root also apply within named workspaces. Requires restrict_root_folder_behaviors_to_site_admins to be enabled.
	RootFolderBehaviorsApplyToWorkspaces *bool `json:"root_folder_behaviors_apply_to_workspaces,omitempty" path:"root_folder_behaviors_apply_to_workspaces,omitempty" url:"root_folder_behaviors_apply_to_workspaces,omitempty"`
	// If true, SSO users using the default user-level two-factor authentication setting are exempt from the site-wide two-factor authentication requirement.
	Require2faExemptAllSsoUsers *bool `json:"require_2fa_exempt_all_sso_users,omitempty" path:"require_2fa_exempt_all_sso_users,omitempty" url:"require_2fa_exempt_all_sso_users,omitempty"`
	// If set, requirement for two-factor authentication has been scheduled to end on this date-time.
	Require2faStopTime *time.Time `json:"require_2fa_stop_time,omitempty" path:"require_2fa_stop_time,omitempty" url:"require_2fa_stop_time,omitempty"`
	// Auto-removes bundles for disabled/deleted users and enforces bundle expiry within user access period.
	RevokeBundleAccessOnDisableOrDelete *bool `json:"revoke_bundle_access_on_disable_or_delete,omitempty" path:"revoke_bundle_access_on_disable_or_delete,omitempty" url:"revoke_bundle_access_on_disable_or_delete,omitempty"`
	// What type of user is required to use two-factor authentication (when require_2fa is set to `true` for this site)?
	Require2faUserType string `json:"require_2fa_user_type,omitempty" path:"require_2fa_user_type,omitempty" url:"require_2fa_user_type,omitempty"`
	// If true, we will hide the 'Remember Me' box on Inbox and Bundle registration pages, requiring that the user logout and log back in every time they visit the page.
	RequireLogoutFromBundlesAndInboxes *bool `json:"require_logout_from_bundles_and_inboxes,omitempty" path:"require_logout_from_bundles_and_inboxes,omitempty" url:"require_logout_from_bundles_and_inboxes,omitempty"`
	// Current session.
	Session Session `json:"session,omitempty" path:"session,omitempty" url:"session,omitempty"`
	// Is the S3-compatible endpoint enabled for all users and workspaces on this site? Defaults to true. When false, user and group S3 permissions do not allow access.
	S3CompatibleEndpointEnabled *bool `json:"s3_compatible_endpoint_enabled,omitempty" path:"s3_compatible_endpoint_enabled,omitempty" url:"s3_compatible_endpoint_enabled,omitempty"`
	// Is SFTP enabled?
	SftpEnabled *bool `json:"sftp_enabled,omitempty" path:"sftp_enabled,omitempty" url:"sftp_enabled,omitempty"`
	// Finalize partial SFTP uploads from interrupted connections? Default: true.
	SftpFinalizePartialUploads *bool `json:"sftp_finalize_partial_uploads,omitempty" path:"sftp_finalize_partial_uploads,omitempty" url:"sftp_finalize_partial_uploads,omitempty"`
	// Sftp Host Key Type.
	SftpHostKeyType string `json:"sftp_host_key_type,omitempty" path:"sftp_host_key_type,omitempty" url:"sftp_host_key_type,omitempty"`
	// Id of the currently selected custom SFTP Host Key.
	ActiveSftpHostKeyId int64 `json:"active_sftp_host_key_id,omitempty" path:"active_sftp_host_key_id,omitempty" url:"active_sftp_host_key_id,omitempty"`
	// Ids of the selected custom SFTP Host Keys.
	ActiveSftpHostKeyIds []int64 `json:"active_sftp_host_key_ids,omitempty" path:"active_sftp_host_key_ids,omitempty" url:"active_sftp_host_key_ids,omitempty"`
	// If true, we will allow weak and known insecure ciphers to be used for SFTP connections. Enabling this setting severely weakens the security of your site and it is not recommend, except as a last resort for compatibility.
	SftpInsecureCiphers *bool `json:"sftp_insecure_ciphers,omitempty" path:"sftp_insecure_ciphers,omitempty" url:"sftp_insecure_ciphers,omitempty"`
	// If true, we will allow weak Diffie Hellman parameters to be used within ciphers for SFTP that are otherwise on our secure list. This has the effect of making the cipher weaker than our normal threshold for security, but is required to support certain legacy or broken SSH and MFT clients. Enabling this weakens security, but not nearly as much as enabling the full `sftp_insecure_ciphers` option.
	SftpInsecureDiffieHellman *bool `json:"sftp_insecure_diffie_hellman,omitempty" path:"sftp_insecure_diffie_hellman,omitempty" url:"sftp_insecure_diffie_hellman,omitempty"`
	// Use user FTP roots also for SFTP?
	SftpUserRootEnabled *bool `json:"sftp_user_root_enabled,omitempty" path:"sftp_user_root_enabled,omitempty" url:"sftp_user_root_enabled,omitempty"`
	// Allow bundle creation.
	SharingEnabled *bool `json:"sharing_enabled,omitempty" path:"sharing_enabled,omitempty" url:"sharing_enabled,omitempty"`
	// Show log in link in user notifications?
	ShowUserNotificationsLogInLink *bool `json:"show_user_notifications_log_in_link,omitempty" path:"show_user_notifications_log_in_link,omitempty" url:"show_user_notifications_log_in_link,omitempty"`
	// Show request access link for users without access? Currently unused.
	ShowRequestAccessLink *bool `json:"show_request_access_link,omitempty" path:"show_request_access_link,omitempty" url:"show_request_access_link,omitempty"`
	// Custom site footer text for authenticated pages.
	SiteFooter string `json:"site_footer,omitempty" path:"site_footer,omitempty" url:"site_footer,omitempty"`
	// Custom site header text for authenticated pages.
	SiteHeader string `json:"site_header,omitempty" path:"site_header,omitempty" url:"site_header,omitempty"`
	// Custom site footer text for public pages.
	SitePublicFooter string `json:"site_public_footer,omitempty" path:"site_public_footer,omitempty" url:"site_public_footer,omitempty"`
	// Custom site header text for public pages.
	SitePublicHeader string `json:"site_public_header,omitempty" path:"site_public_header,omitempty" url:"site_public_header,omitempty"`
	// SMTP server hostname or IP.
	SmtpAddress string `json:"smtp_address,omitempty" path:"smtp_address,omitempty" url:"smtp_address,omitempty"`
	// SMTP server authentication type.
	SmtpAuthentication string `json:"smtp_authentication,omitempty" path:"smtp_authentication,omitempty" url:"smtp_authentication,omitempty"`
	// From address to use when mailing through custom SMTP.
	SmtpFrom string `json:"smtp_from,omitempty" path:"smtp_from,omitempty" url:"smtp_from,omitempty"`
	// SMTP server port.
	SmtpPort int64 `json:"smtp_port,omitempty" path:"smtp_port,omitempty" url:"smtp_port,omitempty"`
	// Custom SMTP encryption mode: if_available (default) uses STARTTLS when offered and otherwise sends credentials and messages unencrypted. A certificate-verified STARTTLS connection automatically changes if_available to require unless smtp_ssl is managed by a parent policy. require requires STARTTLS before authentication; require_implicit uses TLS from connection start; never disables TLS. TLS verifies the server certificate against smtp_address. Encryption is never automatically downgraded.
	SmtpSsl string `json:"smtp_ssl,omitempty" path:"smtp_ssl,omitempty" url:"smtp_ssl,omitempty"`
	// SMTP server username.
	SmtpUsername string `json:"smtp_username,omitempty" path:"smtp_username,omitempty" url:"smtp_username,omitempty"`
	// Session expiry in minutes.
	SessionExpiryMinutes int64 `json:"session_expiry_minutes,omitempty" path:"session_expiry_minutes,omitempty" url:"session_expiry_minutes,omitempty"`
	// Allow snapshot share links creation.
	SnapshotSharingEnabled *bool `json:"snapshot_sharing_enabled,omitempty" path:"snapshot_sharing_enabled,omitempty" url:"snapshot_sharing_enabled,omitempty"`
	// Is SSL required? Disabling this is insecure.
	SslRequired *bool `json:"ssl_required,omitempty" path:"ssl_required,omitempty" url:"ssl_required,omitempty"`
	// Site subdomain.
	Subdomain string `json:"subdomain,omitempty" path:"subdomain,omitempty" url:"subdomain,omitempty"`
	// If switching plans, when does the new plan take effect?
	SwitchToPlanDate *time.Time `json:"switch_to_plan_date,omitempty" path:"switch_to_plan_date,omitempty" url:"switch_to_plan_date,omitempty"`
	// Number of days left in trial.
	TrialDaysLeft int64 `json:"trial_days_left,omitempty" path:"trial_days_left,omitempty" url:"trial_days_left,omitempty"`
	// When does this Site trial expire?
	TrialUntil *time.Time `json:"trial_until,omitempty" path:"trial_until,omitempty" url:"trial_until,omitempty"`
	// If using custom SMTP, should we use dedicated IPs to deliver emails?
	UseDedicatedIpsForSmtp *bool `json:"use_dedicated_ips_for_smtp,omitempty" path:"use_dedicated_ips_for_smtp,omitempty" url:"use_dedicated_ips_for_smtp,omitempty"`
	// Allow uploaders to set `provided_modified_at` for uploaded files?
	UseProvidedModifiedAt *bool `json:"use_provided_modified_at,omitempty" path:"use_provided_modified_at,omitempty" url:"use_provided_modified_at,omitempty"`
	// User of current session.
	User User `json:"user,omitempty" path:"user,omitempty" url:"user,omitempty"`
	// Will users be locked out after incorrect login attempts?
	UserLockout *bool `json:"user_lockout,omitempty" path:"user_lockout,omitempty" url:"user_lockout,omitempty"`
	// How many hours to lock user out for failed password?
	UserLockoutLockPeriod int64 `json:"user_lockout_lock_period,omitempty" path:"user_lockout_lock_period,omitempty" url:"user_lockout_lock_period,omitempty"`
	// Number of login tries within `user_lockout_within` hours before users are locked out.
	UserLockoutTries int64 `json:"user_lockout_tries,omitempty" path:"user_lockout_tries,omitempty" url:"user_lockout_tries,omitempty"`
	// Number of hours for user lockout window.
	UserLockoutWithin int64 `json:"user_lockout_within,omitempty" path:"user_lockout_within,omitempty" url:"user_lockout_within,omitempty"`
	// Enable User Requests feature.
	UserRequestsEnabled *bool `json:"user_requests_enabled,omitempty" path:"user_requests_enabled,omitempty" url:"user_requests_enabled,omitempty"`
	// Send email to site admins when a user request is received?
	UserRequestsNotifyAdmins *bool `json:"user_requests_notify_admins,omitempty" path:"user_requests_notify_admins,omitempty" url:"user_requests_notify_admins,omitempty"`
	// Allow users to create their own API keys?
	UsersCanCreateApiKeys *bool `json:"users_can_create_api_keys,omitempty" path:"users_can_create_api_keys,omitempty" url:"users_can_create_api_keys,omitempty"`
	// Allow users to create their own SSH keys?
	UsersCanCreateSshKeys *bool `json:"users_can_create_ssh_keys,omitempty" path:"users_can_create_ssh_keys,omitempty" url:"users_can_create_ssh_keys,omitempty"`
	// How usernames are displayed in the web UI. Can be `username_only`, `full_name_only`, `full_name_username`, `full_name_company`, or `full_name_username_company`.
	UsernameDisplay string `json:"username_display,omitempty" path:"username_display,omitempty" url:"username_display,omitempty"`
	// Custom text send in user welcome email.
	WelcomeCustomText string `json:"welcome_custom_text,omitempty" path:"welcome_custom_text,omitempty" url:"welcome_custom_text,omitempty"`
	// Custom footer text for system-generated emails (as Markdown). Supports standard strftime date/time patterns like %Y (4-digit year), %m (month), %d (day).
	EmailFooterCustomText string `json:"email_footer_custom_text,omitempty" path:"email_footer_custom_text,omitempty" url:"email_footer_custom_text,omitempty"`
	// Include this email in welcome emails if enabled.
	WelcomeEmailCc string `json:"welcome_email_cc,omitempty" path:"welcome_email_cc,omitempty" url:"welcome_email_cc,omitempty"`
	// Include this email subject in welcome emails if enabled.
	WelcomeEmailSubject string `json:"welcome_email_subject,omitempty" path:"welcome_email_subject,omitempty" url:"welcome_email_subject,omitempty"`
	// Will the welcome email be sent to new users?
	WelcomeEmailEnabled *bool `json:"welcome_email_enabled,omitempty" path:"welcome_email_enabled,omitempty" url:"welcome_email_enabled,omitempty"`
	// Does the welcome screen appear?
	WelcomeScreen string `json:"welcome_screen,omitempty" path:"welcome_screen,omitempty" url:"welcome_screen,omitempty"`
	// Does FTP user Windows emulation mode?
	WindowsModeFtp *bool `json:"windows_mode_ftp,omitempty" path:"windows_mode_ftp,omitempty" url:"windows_mode_ftp,omitempty"`
}

// Identifier returns the resource ID.
func (s Site) Identifier() interface{} {
	return s.Id
}

// SiteCollection is a list of Site resources.
type SiteCollection []Site

// SiteUpdateParams contains the request parameters for PATCH /site.
type SiteUpdateParams struct {
	// Site name.
	Name string `url:"name,omitempty" json:"name,omitempty" path:"name"`
	// Site subdomain.
	Subdomain string `url:"subdomain,omitempty" json:"subdomain,omitempty" path:"subdomain"`
	// Custom domain.
	Domain string `url:"domain,omitempty" json:"domain,omitempty" path:"domain"`
	// Send HSTS (HTTP Strict Transport Security) header when visitors access the site via a custom domain?
	DomainHstsHeader *bool `url:"domain_hsts_header,omitempty" json:"domain_hsts_header,omitempty" path:"domain_hsts_header"`
	// Letsencrypt chain to use when registering SSL Certificate for domain. No longer used as of 2026.
	DomainLetsencryptChain string `url:"domain_letsencrypt_chain,omitempty" json:"domain_letsencrypt_chain,omitempty" path:"domain_letsencrypt_chain"`
	// Main email for this site.
	Email string `url:"email,omitempty" json:"email,omitempty" path:"email"`
	// Reply-to email for this site.
	ReplyToEmail string `url:"reply_to_email,omitempty" json:"reply_to_email,omitempty" path:"reply_to_email"`
	// Are manual Bundle names allowed?
	AllowBundleNames *bool `url:"allow_bundle_names,omitempty" json:"allow_bundle_names,omitempty" path:"allow_bundle_names"`
	// Site-wide Bundle expiration in days.
	BundleExpiration int64 `url:"bundle_expiration,omitempty" json:"bundle_expiration,omitempty" path:"bundle_expiration"`
	// Will the welcome email be sent to new users?
	WelcomeEmailEnabled *bool `url:"welcome_email_enabled,omitempty" json:"welcome_email_enabled,omitempty" path:"welcome_email_enabled"`
	// If false, rename conflicting files instead of asking for overwrite confirmation. Only applies to web interface.
	AskAboutOverwrites *bool `url:"ask_about_overwrites,omitempty" json:"ask_about_overwrites,omitempty" path:"ask_about_overwrites"`
	// Show request access link for users without access? Currently unused.
	ShowRequestAccessLink *bool `url:"show_request_access_link,omitempty" json:"show_request_access_link,omitempty" path:"show_request_access_link"`
	// Create parent directories if they do not exist during uploads? This is primarily used to work around broken upload clients that assume servers will perform this step.
	AlwaysMkdirParents *bool `url:"always_mkdir_parents,omitempty" json:"always_mkdir_parents,omitempty" path:"always_mkdir_parents"`
	// Include this email in welcome emails if enabled.
	WelcomeEmailCc string `url:"welcome_email_cc,omitempty" json:"welcome_email_cc,omitempty" path:"welcome_email_cc"`
	// Include this email subject in welcome emails if enabled.
	WelcomeEmailSubject string `url:"welcome_email_subject,omitempty" json:"welcome_email_subject,omitempty" path:"welcome_email_subject"`
	// Custom text send in user welcome email.
	WelcomeCustomText string `url:"welcome_custom_text,omitempty" json:"welcome_custom_text,omitempty" path:"welcome_custom_text"`
	// Site default language.
	Language string `url:"language,omitempty" json:"language,omitempty" path:"language"`
	// Does FTP user Windows emulation mode?
	WindowsModeFtp *bool `url:"windows_mode_ftp,omitempty" json:"windows_mode_ftp,omitempty" path:"windows_mode_ftp"`
	// Site default time zone.
	DefaultTimeZone string `url:"default_time_zone,omitempty" json:"default_time_zone,omitempty" path:"default_time_zone"`
	// Is the desktop app enabled?
	DesktopApp *bool `url:"desktop_app,omitempty" json:"desktop_app,omitempty" path:"desktop_app"`
	// Is desktop app session IP pinning enabled?
	DesktopAppSessionIpPinning *bool `url:"desktop_app_session_ip_pinning,omitempty" json:"desktop_app_session_ip_pinning,omitempty" path:"desktop_app_session_ip_pinning"`
	// Desktop app session lifetime (in hours).
	DesktopAppSessionLifetime int64 `url:"desktop_app_session_lifetime,omitempty" json:"desktop_app_session_lifetime,omitempty" path:"desktop_app_session_lifetime"`
	// Is the mobile app enabled?
	MobileApp *bool `url:"mobile_app,omitempty" json:"mobile_app,omitempty" path:"mobile_app"`
	// Is mobile app session IP pinning enabled?
	MobileAppSessionIpPinning *bool `url:"mobile_app_session_ip_pinning,omitempty" json:"mobile_app_session_ip_pinning,omitempty" path:"mobile_app_session_ip_pinning"`
	// Mobile app session lifetime (in hours).
	MobileAppSessionLifetime int64 `url:"mobile_app_session_lifetime,omitempty" json:"mobile_app_session_lifetime,omitempty" path:"mobile_app_session_lifetime"`
	// If true, permissions for this site must be bound to a group (not a user).
	FolderPermissionsGroupsOnly *bool `url:"folder_permissions_groups_only,omitempty" json:"folder_permissions_groups_only,omitempty" path:"folder_permissions_groups_only"`
	// Does the welcome screen appear?
	WelcomeScreen string `url:"welcome_screen,omitempty" json:"welcome_screen,omitempty" path:"welcome_screen"`
	// If true, allows users to use a document editing integration.
	OfficeIntegrationAvailable *bool `url:"office_integration_available,omitempty" json:"office_integration_available,omitempty" path:"office_integration_available"`
	// Which document editing integration to support. Files.com Editor or Microsoft Office for the Web.
	OfficeIntegrationType string `url:"office_integration_type,omitempty" json:"office_integration_type,omitempty" path:"office_integration_type"`
	// If true, we will ensure that all internal communications with any remote server are made through the primary region of the site. This setting overrides individual remote server settings.
	PinAllRemoteServersToSiteRegion *bool `url:"pin_all_remote_servers_to_site_region,omitempty" json:"pin_all_remote_servers_to_site_region,omitempty" path:"pin_all_remote_servers_to_site_region"`
	// A message to show users when they connect via FTP or SFTP.
	MotdText string `url:"motd_text,omitempty" json:"motd_text,omitempty" path:"motd_text"`
	// Show message to users connecting via FTP.
	MotdUseForFtp *bool `url:"motd_use_for_ftp,omitempty" json:"motd_use_for_ftp,omitempty" path:"motd_use_for_ftp"`
	// Show message to users connecting via SFTP.
	MotdUseForSftp *bool `url:"motd_use_for_sftp,omitempty" json:"motd_use_for_sftp,omitempty" path:"motd_use_for_sftp"`
	// Visibility settings for account navigation.
	LeftNavigationVisibility interface{} `url:"left_navigation_visibility,omitempty" json:"left_navigation_visibility,omitempty" path:"left_navigation_visibility"`
	// If true, all AI features are disabled for this site.
	DisableAllAiFeatures *bool `url:"disable_all_ai_features,omitempty" json:"disable_all_ai_features,omitempty" path:"disable_all_ai_features"`
	// Availability settings for AI features. Each feature requires the site_admins, workspace_admins, folder_admins, and all_users keys. Optional selected_group_members defaults to false; when true, members of at least one group in group_ids get access regardless of the other options. All availability options are additive: any enabled option matching the user grants access. Optional group_ids is an array of integer IDs of active groups on this site, from any workspace. Omitted or empty group_ids grants no access through selected_group_members and does not affect other options. Disabling all AI features overrides these settings.
	AiFeatureAvailability interface{} `url:"ai_feature_availability,omitempty" json:"ai_feature_availability,omitempty" path:"ai_feature_availability"`
	// Is OAuth DCR (dynamic client registration) for MCP enabled?
	McpDcrEnabled *bool `url:"mcp_dcr_enabled,omitempty" json:"mcp_dcr_enabled,omitempty" path:"mcp_dcr_enabled"`
	// Additional extensions that are considered text files.
	AdditionalTextFileTypes []string `url:"additional_text_file_types,omitempty" json:"additional_text_file_types,omitempty" path:"additional_text_file_types"`
	// Do Bundles require internal notes?
	BundleRequireNote *bool `url:"bundle_require_note,omitempty" json:"bundle_require_note,omitempty" path:"bundle_require_note"`
	// Do Bundle creators receive receipts of invitations?
	BundleSendSharedReceipts *bool `url:"bundle_send_shared_receipts,omitempty" json:"bundle_send_shared_receipts,omitempty" path:"bundle_send_shared_receipts"`
	// If true, new Share Links created by a user with a primary group will default to that group as owner.
	BundlesDefaultOwnedByPrimaryGroup *bool `url:"bundles_default_owned_by_primary_group,omitempty" json:"bundles_default_owned_by_primary_group,omitempty" path:"bundles_default_owned_by_primary_group"`
	// Calculate CRC32 checksums for files?
	CalculateFileChecksumsCrc32 *bool `url:"calculate_file_checksums_crc32,omitempty" json:"calculate_file_checksums_crc32,omitempty" path:"calculate_file_checksums_crc32"`
	// Calculate MD5 checksums for files?
	CalculateFileChecksumsMd5 *bool `url:"calculate_file_checksums_md5,omitempty" json:"calculate_file_checksums_md5,omitempty" path:"calculate_file_checksums_md5"`
	// Calculate SHA1 checksums for files?
	CalculateFileChecksumsSha1 *bool `url:"calculate_file_checksums_sha1,omitempty" json:"calculate_file_checksums_sha1,omitempty" path:"calculate_file_checksums_sha1"`
	// Calculate SHA256 checksums for files?
	CalculateFileChecksumsSha256 *bool `url:"calculate_file_checksums_sha256,omitempty" json:"calculate_file_checksums_sha256,omitempty" path:"calculate_file_checksums_sha256"`
	// Use legacy checksums mode?
	LegacyChecksumsMode *bool `url:"legacy_checksums_mode,omitempty" json:"legacy_checksums_mode,omitempty" path:"legacy_checksums_mode"`
	// If true, we will migrate all remote server syncs to the new Sync model.
	MigrateRemoteServerSyncToSync *bool `url:"migrate_remote_server_sync_to_sync,omitempty" json:"migrate_remote_server_sync_to_sync,omitempty" path:"migrate_remote_server_sync_to_sync"`
	// Number of days to retain AS2 messages (incoming and outgoing).
	As2MessageRetentionDays int64 `url:"as2_message_retention_days,omitempty" json:"as2_message_retention_days,omitempty" path:"as2_message_retention_days"`
	// If true, file operations produce file_action.activity.v1 EventChannel events in all workspaces. Defaults to false; subscription filters and workspace access still control delivery.
	FileActivityEventsEnabled *bool `url:"file_activity_events_enabled,omitempty" json:"file_activity_events_enabled,omitempty" path:"file_activity_events_enabled"`
	// How usernames are displayed in the web UI. Can be `username_only`, `full_name_only`, `full_name_username`, `full_name_company`, or `full_name_username_company`.
	UsernameDisplay string `url:"username_display,omitempty" json:"username_display,omitempty" path:"username_display"`
	// Session expiry in minutes.
	SessionExpiryMinutes int64 `url:"session_expiry_minutes,omitempty" json:"session_expiry_minutes,omitempty" path:"session_expiry_minutes"`
	// Is SSL required? Disabling this is insecure.
	SslRequired *bool `url:"ssl_required,omitempty" json:"ssl_required,omitempty" path:"ssl_required"`
	// If true, we will allow weak and known insecure ciphers to be used for SFTP connections. Enabling this setting severely weakens the security of your site and it is not recommend, except as a last resort for compatibility.
	SftpInsecureCiphers *bool `url:"sftp_insecure_ciphers,omitempty" json:"sftp_insecure_ciphers,omitempty" path:"sftp_insecure_ciphers"`
	// If true, we will allow weak Diffie Hellman parameters to be used within ciphers for SFTP that are otherwise on our secure list. This has the effect of making the cipher weaker than our normal threshold for security, but is required to support certain legacy or broken SSH and MFT clients. Enabling this weakens security, but not nearly as much as enabling the full `sftp_insecure_ciphers` option.
	SftpInsecureDiffieHellman *bool `url:"sftp_insecure_diffie_hellman,omitempty" json:"sftp_insecure_diffie_hellman,omitempty" path:"sftp_insecure_diffie_hellman"`
	// If set, Files.com will not set the CAA records required to generate future SSL certificates for this domain.
	DisableFilesCertificateGeneration *bool `url:"disable_files_certificate_generation,omitempty" json:"disable_files_certificate_generation,omitempty" path:"disable_files_certificate_generation"`
	// Are FedRAMP security restrictions enabled for this site?
	Fedramp *bool `url:"fedramp,omitempty" json:"fedramp,omitempty" path:"fedramp"`
	// Will users be locked out after incorrect login attempts?
	UserLockout *bool `url:"user_lockout,omitempty" json:"user_lockout,omitempty" path:"user_lockout"`
	// Number of login tries within `user_lockout_within` hours before users are locked out.
	UserLockoutTries int64 `url:"user_lockout_tries,omitempty" json:"user_lockout_tries,omitempty" path:"user_lockout_tries"`
	// Number of hours for user lockout window.
	UserLockoutWithin int64 `url:"user_lockout_within,omitempty" json:"user_lockout_within,omitempty" path:"user_lockout_within"`
	// How many hours to lock user out for failed password?
	UserLockoutLockPeriod int64 `url:"user_lockout_lock_period,omitempty" json:"user_lockout_lock_period,omitempty" path:"user_lockout_lock_period"`
	// Include password in emails to new users?
	IncludePasswordInWelcomeEmail *bool `url:"include_password_in_welcome_email,omitempty" json:"include_password_in_welcome_email,omitempty" path:"include_password_in_welcome_email"`
	// Comma separated list of allowed Country codes.
	AllowedCountries string `url:"allowed_countries,omitempty" json:"allowed_countries,omitempty" path:"allowed_countries"`
	// List of allowed IP addresses.
	AllowedIps string `url:"allowed_ips,omitempty" json:"allowed_ips,omitempty" path:"allowed_ips"`
	// Allow the site-wide two-factor authentication requirement to be overriden on a per-user-basis?
	AllowUserLevel2faOverride *bool `url:"allow_user_level_2fa_override,omitempty" json:"allow_user_level_2fa_override,omitempty" path:"allow_user_level_2fa_override"`
	// Allow the site-wide allowed IP restriction to be overriden on a per-user-basis?
	AllowUserLevelAllowedIpOverride *bool `url:"allow_user_level_allowed_ip_override,omitempty" json:"allow_user_level_allowed_ip_override,omitempty" path:"allow_user_level_allowed_ip_override"`
	// Allow the site-wide FTP SSL requirement to be overriden on a per-user-basis?
	AllowUserLevelSslOverride *bool `url:"allow_user_level_ssl_override,omitempty" json:"allow_user_level_ssl_override,omitempty" path:"allow_user_level_ssl_override"`
	// Require Parent Site users entering this Child Site to match its allowed IP addresses in addition to their existing IP restrictions. Applies to all Workspaces. Defaults to false.
	EnforceParentSiteUserIpRestrictions *bool `url:"enforce_parent_site_user_ip_restrictions,omitempty" json:"enforce_parent_site_user_ip_restrictions,omitempty" path:"enforce_parent_site_user_ip_restrictions"`
	// Comma separated list of disallowed Country codes.
	DisallowedCountries string `url:"disallowed_countries,omitempty" json:"disallowed_countries,omitempty" path:"disallowed_countries"`
	// Number of days to keep deleted files.
	DaysToRetainBackups int64 `url:"days_to_retain_backups,omitempty" json:"days_to_retain_backups,omitempty" path:"days_to_retain_backups"`
	// Number of prior passwords to disallow.
	MaxPriorPasswords int64 `url:"max_prior_passwords,omitempty" json:"max_prior_passwords,omitempty" path:"max_prior_passwords"`
	// Number of days password is valid.
	PasswordValidityDays int64 `url:"password_validity_days,omitempty" json:"password_validity_days,omitempty" path:"password_validity_days"`
	// Shortest password length for users.
	PasswordMinLength int64 `url:"password_min_length,omitempty" json:"password_min_length,omitempty" path:"password_min_length"`
	// Require a letter in passwords?
	PasswordRequireLetter *bool `url:"password_require_letter,omitempty" json:"password_require_letter,omitempty" path:"password_require_letter"`
	// Require lower and upper case letters in passwords?
	PasswordRequireMixed *bool `url:"password_require_mixed,omitempty" json:"password_require_mixed,omitempty" path:"password_require_mixed"`
	// Require special characters in password?
	PasswordRequireSpecial *bool `url:"password_require_special,omitempty" json:"password_require_special,omitempty" path:"password_require_special"`
	// Require a number in passwords?
	PasswordRequireNumber *bool `url:"password_require_number,omitempty" json:"password_require_number,omitempty" path:"password_require_number"`
	// Require passwords that have not been previously breached? (see https://haveibeenpwned.com/).
	PasswordRequireUnbreached *bool `url:"password_require_unbreached,omitempty" json:"password_require_unbreached,omitempty" path:"password_require_unbreached"`
	// If true, we will hide the 'Remember Me' box on Inbox and Bundle registration pages, requiring that the user logout and log back in every time they visit the page.
	RequireLogoutFromBundlesAndInboxes *bool `url:"require_logout_from_bundles_and_inboxes,omitempty" json:"require_logout_from_bundles_and_inboxes,omitempty" path:"require_logout_from_bundles_and_inboxes"`
	// Use user FTP roots also for WebDAV?
	DavUserRootEnabled *bool `url:"dav_user_root_enabled,omitempty" json:"dav_user_root_enabled,omitempty" path:"dav_user_root_enabled"`
	// Use user FTP roots also for SFTP?
	SftpUserRootEnabled *bool `url:"sftp_user_root_enabled,omitempty" json:"sftp_user_root_enabled,omitempty" path:"sftp_user_root_enabled"`
	// Is password reset disabled?
	DisablePasswordReset *bool `url:"disable_password_reset,omitempty" json:"disable_password_reset,omitempty" path:"disable_password_reset"`
	// Are files protected from modification?
	ImmutableFiles *bool `url:"immutable_files,omitempty" json:"immutable_files,omitempty" path:"immutable_files"`
	// Custom error message to show when bundle is not found.
	BundleNotFoundMessage string `url:"bundle_not_found_message,omitempty" json:"bundle_not_found_message,omitempty" path:"bundle_not_found_message"`
	// Do Bundles require password protection?
	BundlePasswordRequired *bool `url:"bundle_password_required,omitempty" json:"bundle_password_required,omitempty" path:"bundle_password_required"`
	// Do Bundles require registration?
	BundleRequireRegistration *bool `url:"bundle_require_registration,omitempty" json:"bundle_require_registration,omitempty" path:"bundle_require_registration"`
	// Do Bundles require recipients for sharing?
	BundleRequireShareRecipient *bool `url:"bundle_require_share_recipient,omitempty" json:"bundle_require_share_recipient,omitempty" path:"bundle_require_share_recipient"`
	// If true, new Share Links must send a one-time password to the recipient when they register. Requires bundle_require_share_recipient and cannot be enabled with bundle_password_required.
	BundleSendOneTimePasswordToRecipientAtRegistration *bool `url:"bundle_send_one_time_password_to_recipient_at_registration,omitempty" json:"bundle_send_one_time_password_to_recipient_at_registration,omitempty" path:"bundle_send_one_time_password_to_recipient_at_registration"`
	// Do Bundle owners receive registration notification?
	BundleRegistrationNotifications string `url:"bundle_registration_notifications,omitempty" json:"bundle_registration_notifications,omitempty" path:"bundle_registration_notifications"`
	// Do Bundle owners receive activity notifications?
	BundleActivityNotifications string `url:"bundle_activity_notifications,omitempty" json:"bundle_activity_notifications,omitempty" path:"bundle_activity_notifications"`
	// Do Bundle uploaders receive upload confirmation notifications?
	BundleUploadReceiptNotifications string `url:"bundle_upload_receipt_notifications,omitempty" json:"bundle_upload_receipt_notifications,omitempty" path:"bundle_upload_receipt_notifications"`
	// If true, allow public viewers of Bundles with full permissions to use document editing integrations.
	DocumentEditsInBundleAllowed *bool `url:"document_edits_in_bundle_allowed,omitempty" json:"document_edits_in_bundle_allowed,omitempty" path:"document_edits_in_bundle_allowed"`
	// Require bundles' passwords, and passwords for other items (inboxes, public shares, etc.) to conform to the same requirements as users' passwords?
	PasswordRequirementsApplyToBundles *bool `url:"password_requirements_apply_to_bundles,omitempty" json:"password_requirements_apply_to_bundles,omitempty" path:"password_requirements_apply_to_bundles"`
	// If true, we will prevent non-administrators from receiving any permissions directly on the root folder. This is commonly used to prevent the accidental application of permissions.
	PreventRootPermissionsForNonSiteAdmins *bool `url:"prevent_root_permissions_for_non_site_admins,omitempty" json:"prevent_root_permissions_for_non_site_admins,omitempty" path:"prevent_root_permissions_for_non_site_admins"`
	// If true, only site admins may create, modify, or delete any behavior at the site root, or a skip that would disable one.
	RestrictRootFolderBehaviorsToSiteAdmins *bool `url:"restrict_root_folder_behaviors_to_site_admins,omitempty" json:"restrict_root_folder_behaviors_to_site_admins,omitempty" path:"restrict_root_folder_behaviors_to_site_admins"`
	// If true, supported protective behaviors at the site root also apply within named workspaces. Requires restrict_root_folder_behaviors_to_site_admins to be enabled.
	RootFolderBehaviorsApplyToWorkspaces *bool `url:"root_folder_behaviors_apply_to_workspaces,omitempty" json:"root_folder_behaviors_apply_to_workspaces,omitempty" path:"root_folder_behaviors_apply_to_workspaces"`
	// Use servers in the USA only?
	OptOutGlobal *bool `url:"opt_out_global,omitempty" json:"opt_out_global,omitempty" path:"opt_out_global"`
	// Allow uploaders to set `provided_modified_at` for uploaded files?
	UseProvidedModifiedAt *bool `url:"use_provided_modified_at,omitempty" json:"use_provided_modified_at,omitempty" path:"use_provided_modified_at"`
	// Is this site using a custom namespace for users?
	CustomNamespace *bool `url:"custom_namespace,omitempty" json:"custom_namespace,omitempty" path:"custom_namespace"`
	// If true, groups can be manually created / modified / deleted by Site Admins. Otherwise, groups can only be managed via your SSO provider.
	NonSsoGroupsAllowed *bool `url:"non_sso_groups_allowed,omitempty" json:"non_sso_groups_allowed,omitempty" path:"non_sso_groups_allowed"`
	// If true, users can be manually created / modified / deleted by Site Admins. Otherwise, users can only be managed via your SSO provider.
	NonSsoUsersAllowed *bool `url:"non_sso_users_allowed,omitempty" json:"non_sso_users_allowed,omitempty" path:"non_sso_users_allowed"`
	// Allow bundle creation.
	SharingEnabled *bool `url:"sharing_enabled,omitempty" json:"sharing_enabled,omitempty" path:"sharing_enabled"`
	// Allow snapshot share links creation.
	SnapshotSharingEnabled *bool `url:"snapshot_sharing_enabled,omitempty" json:"snapshot_sharing_enabled,omitempty" path:"snapshot_sharing_enabled"`
	// Enable User Requests feature.
	UserRequestsEnabled *bool `url:"user_requests_enabled,omitempty" json:"user_requests_enabled,omitempty" path:"user_requests_enabled"`
	// Send email to site admins when a user request is received?
	UserRequestsNotifyAdmins *bool `url:"user_requests_notify_admins,omitempty" json:"user_requests_notify_admins,omitempty" path:"user_requests_notify_admins"`
	// Is WebDAV enabled?
	DavEnabled *bool `url:"dav_enabled,omitempty" json:"dav_enabled,omitempty" path:"dav_enabled"`
	// Allow other Files.com sites to use this site's API keys for native Files.com Remote Server connections? Defaults to true and applies to all keys and workspaces on this site. When false, new pairings and access through existing connections are rejected without revoking keys or deleting connections. Re-enabling permits access again with usable keys. Does not disable this site's connections to other sites, ordinary API access, or Connected Sites.
	FilesComRemoteServerEnabled *bool `url:"files_com_remote_server_enabled,omitempty" json:"files_com_remote_server_enabled,omitempty" path:"files_com_remote_server_enabled"`
	// Is FTP enabled?
	FtpEnabled *bool `url:"ftp_enabled,omitempty" json:"ftp_enabled,omitempty" path:"ftp_enabled"`
	// Is the S3-compatible endpoint enabled for all users and workspaces on this site? Defaults to true. When false, user and group S3 permissions do not allow access.
	S3CompatibleEndpointEnabled *bool `url:"s3_compatible_endpoint_enabled,omitempty" json:"s3_compatible_endpoint_enabled,omitempty" path:"s3_compatible_endpoint_enabled"`
	// Is SFTP enabled?
	SftpEnabled *bool `url:"sftp_enabled,omitempty" json:"sftp_enabled,omitempty" path:"sftp_enabled"`
	// Finalize partial SFTP uploads from interrupted connections? Default: true.
	SftpFinalizePartialUploads *bool `url:"sftp_finalize_partial_uploads,omitempty" json:"sftp_finalize_partial_uploads,omitempty" path:"sftp_finalize_partial_uploads"`
	// Allow users to create their own API keys?
	UsersCanCreateApiKeys *bool `url:"users_can_create_api_keys,omitempty" json:"users_can_create_api_keys,omitempty" path:"users_can_create_api_keys"`
	// Allow users to create their own SSH keys?
	UsersCanCreateSshKeys *bool `url:"users_can_create_ssh_keys,omitempty" json:"users_can_create_ssh_keys,omitempty" path:"users_can_create_ssh_keys"`
	// Show log in link in user notifications?
	ShowUserNotificationsLogInLink *bool `url:"show_user_notifications_log_in_link,omitempty" json:"show_user_notifications_log_in_link,omitempty" path:"show_user_notifications_log_in_link"`
	// Sftp Host Key Type.
	SftpHostKeyType string `url:"sftp_host_key_type,omitempty" json:"sftp_host_key_type,omitempty" path:"sftp_host_key_type"`
	// Id of the currently selected custom SFTP Host Key.
	ActiveSftpHostKeyId int64 `url:"active_sftp_host_key_id,omitempty" json:"active_sftp_host_key_id,omitempty" path:"active_sftp_host_key_id"`
	// Ids of the selected custom SFTP Host Keys.
	ActiveSftpHostKeyIds []int64 `url:"active_sftp_host_key_ids,omitempty" json:"active_sftp_host_key_ids,omitempty" path:"active_sftp_host_key_ids"`
	// If true, protocol access permissions on users will be ignored, and only protocol access permissions set on Groups will be honored. Make sure that your current user is a member of a group with API permission when changing this value to avoid locking yourself out of your site.
	ProtocolAccessGroupsOnly *bool `url:"protocol_access_groups_only,omitempty" json:"protocol_access_groups_only,omitempty" path:"protocol_access_groups_only"`
	// Auto-removes bundles for disabled/deleted users and enforces bundle expiry within user access period.
	RevokeBundleAccessOnDisableOrDelete *bool `url:"revoke_bundle_access_on_disable_or_delete,omitempty" json:"revoke_bundle_access_on_disable_or_delete,omitempty" path:"revoke_bundle_access_on_disable_or_delete"`
	// Allow group admins to create users in their groups.
	GroupAdminsCanAddUsers *bool `url:"group_admins_can_add_users,omitempty" json:"group_admins_can_add_users,omitempty" path:"group_admins_can_add_users"`
	// Allow group admins to add or remove existing users in their groups.
	GroupAdminsCanManageGroupMemberships *bool `url:"group_admins_can_manage_group_memberships,omitempty" json:"group_admins_can_manage_group_memberships,omitempty" path:"group_admins_can_manage_group_memberships"`
	// Allow group admins to delete users in their groups.
	GroupAdminsCanDeleteUsers *bool `url:"group_admins_can_delete_users,omitempty" json:"group_admins_can_delete_users,omitempty" path:"group_admins_can_delete_users"`
	// Allow group admins to enable or disable users in their groups.
	GroupAdminsCanEnableDisableUsers *bool `url:"group_admins_can_enable_disable_users,omitempty" json:"group_admins_can_enable_disable_users,omitempty" path:"group_admins_can_enable_disable_users"`
	// Allow group admins to modify users in their groups.
	GroupAdminsCanModifyUsers *bool `url:"group_admins_can_modify_users,omitempty" json:"group_admins_can_modify_users,omitempty" path:"group_admins_can_modify_users"`
	// Allow group admins to exempt users in their groups from lifecycle rules.
	GroupAdminsCanBypassUserLifecycleRules *bool `url:"group_admins_can_bypass_user_lifecycle_rules,omitempty" json:"group_admins_can_bypass_user_lifecycle_rules,omitempty" path:"group_admins_can_bypass_user_lifecycle_rules"`
	// Allow group admins to reset passwords for users in their groups.
	GroupAdminsCanResetPasswords *bool `url:"group_admins_can_reset_passwords,omitempty" json:"group_admins_can_reset_passwords,omitempty" path:"group_admins_can_reset_passwords"`
	// Allow group admins to set password authentication method.
	GroupAdminsCanSetUserPassword *bool `url:"group_admins_can_set_user_password,omitempty" json:"group_admins_can_set_user_password,omitempty" path:"group_admins_can_set_user_password"`
	// If true, administrators on this site's Child Sites, including Folder Admins, can list this site's Default Workspace users and groups and grant them folder permissions and notifications on Child Site paths.
	ChildSiteFolderAdminsCanSeeParentSiteDirectory *bool `url:"child_site_folder_admins_can_see_parent_site_directory,omitempty" json:"child_site_folder_admins_can_see_parent_site_directory,omitempty" path:"child_site_folder_admins_can_see_parent_site_directory"`
	// If true, Workspace Administrators and Folder Admins in a Custom Workspace can list the Default Workspace users and groups that have permissions in that Workspace and choose them as Email Notification recipients.
	WorkspaceFolderAdminsCanSeeDefaultWorkspaceDirectory *bool `url:"workspace_folder_admins_can_see_default_workspace_directory,omitempty" json:"workspace_folder_admins_can_see_default_workspace_directory,omitempty" path:"workspace_folder_admins_can_see_default_workspace_directory"`
	// Disallow free email domains for Bundle/Inbox recipients?
	BundleRecipientBlacklistFreeEmailDomains *bool `url:"bundle_recipient_blacklist_free_email_domains,omitempty" json:"bundle_recipient_blacklist_free_email_domains,omitempty" path:"bundle_recipient_blacklist_free_email_domains"`
	// List of email domains to disallow when entering a Bundle/Inbox recipients.
	BundleRecipientBlacklistDomains []string `url:"bundle_recipient_blacklist_domains,omitempty" json:"bundle_recipient_blacklist_domains,omitempty" path:"bundle_recipient_blacklist_domains"`
	// Allow admins to bypass the locked subfolders setting.
	AdminsBypassLockedSubfolders *bool `url:"admins_bypass_locked_subfolders,omitempty" json:"admins_bypass_locked_subfolders,omitempty" path:"admins_bypass_locked_subfolders"`
	// Is SMS two factor authentication allowed?
	Allowed2faMethodSms *bool `url:"allowed_2fa_method_sms,omitempty" json:"allowed_2fa_method_sms,omitempty" path:"allowed_2fa_method_sms"`
	// Is TOTP two factor authentication allowed?
	Allowed2faMethodTotp *bool `url:"allowed_2fa_method_totp,omitempty" json:"allowed_2fa_method_totp,omitempty" path:"allowed_2fa_method_totp"`
	// Is WebAuthn two factor authentication allowed?
	Allowed2faMethodWebauthn *bool `url:"allowed_2fa_method_webauthn,omitempty" json:"allowed_2fa_method_webauthn,omitempty" path:"allowed_2fa_method_webauthn"`
	// Is yubikey two factor authentication allowed?
	Allowed2faMethodYubi *bool `url:"allowed_2fa_method_yubi,omitempty" json:"allowed_2fa_method_yubi,omitempty" path:"allowed_2fa_method_yubi"`
	// Is OTP via email two factor authentication allowed?
	Allowed2faMethodEmail *bool `url:"allowed_2fa_method_email,omitempty" json:"allowed_2fa_method_email,omitempty" path:"allowed_2fa_method_email"`
	// Is OTP via static codes for two factor authentication allowed?
	Allowed2faMethodStatic *bool `url:"allowed_2fa_method_static,omitempty" json:"allowed_2fa_method_static,omitempty" path:"allowed_2fa_method_static"`
	// Are users allowed to configure their two factor authentication to be bypassed for FTP/SFTP/WebDAV?
	Allowed2faMethodBypassForFtpSftpDav *bool `url:"allowed_2fa_method_bypass_for_ftp_sftp_dav,omitempty" json:"allowed_2fa_method_bypass_for_ftp_sftp_dav,omitempty" path:"allowed_2fa_method_bypass_for_ftp_sftp_dav"`
	// Require two-factor authentication for all users?
	Require2fa *bool `url:"require_2fa,omitempty" json:"require_2fa,omitempty" path:"require_2fa"`
	// If true, SSO users using the default user-level two-factor authentication setting are exempt from the site-wide two-factor authentication requirement.
	Require2faExemptAllSsoUsers *bool `url:"require_2fa_exempt_all_sso_users,omitempty" json:"require_2fa_exempt_all_sso_users,omitempty" path:"require_2fa_exempt_all_sso_users"`
	// What type of user is required to use two-factor authentication (when require_2fa is set to `true` for this site)?
	Require2faUserType string `url:"require_2fa_user_type,omitempty" json:"require_2fa_user_type,omitempty" path:"require_2fa_user_type"`
	// Top bar background color.
	Color2Top string `url:"color2_top,omitempty" json:"color2_top,omitempty" path:"color2_top"`
	// Page link and button color.
	Color2Left string `url:"color2_left,omitempty" json:"color2_left,omitempty" path:"color2_left"`
	// Top bar link color.
	Color2Link string `url:"color2_link,omitempty" json:"color2_link,omitempty" path:"color2_link"`
	// Page link and button color.
	Color2Text string `url:"color2_text,omitempty" json:"color2_text,omitempty" path:"color2_text"`
	// Top bar text color.
	Color2TopText string `url:"color2_top_text,omitempty" json:"color2_top_text,omitempty" path:"color2_top_text"`
	// Custom site header text for authenticated pages.
	SiteHeader string `url:"site_header,omitempty" json:"site_header,omitempty" path:"site_header"`
	// Custom site footer text for authenticated pages.
	SiteFooter string `url:"site_footer,omitempty" json:"site_footer,omitempty" path:"site_footer"`
	// Custom site header text for public pages.
	SitePublicHeader string `url:"site_public_header,omitempty" json:"site_public_header,omitempty" path:"site_public_header"`
	// Custom site footer text for public pages.
	SitePublicFooter string `url:"site_public_footer,omitempty" json:"site_public_footer,omitempty" path:"site_public_footer"`
	// Login help text.
	LoginHelpText string `url:"login_help_text,omitempty" json:"login_help_text,omitempty" path:"login_help_text"`
	// If using custom SMTP, should we use dedicated IPs to deliver emails?
	UseDedicatedIpsForSmtp *bool `url:"use_dedicated_ips_for_smtp,omitempty" json:"use_dedicated_ips_for_smtp,omitempty" path:"use_dedicated_ips_for_smtp"`
	// Custom footer text for system-generated emails (as Markdown). Supports standard strftime date/time patterns like %Y (4-digit year), %m (month), %d (day).
	EmailFooterCustomText string `url:"email_footer_custom_text,omitempty" json:"email_footer_custom_text,omitempty" path:"email_footer_custom_text"`
	// SMTP server hostname or IP.
	SmtpAddress string `url:"smtp_address,omitempty" json:"smtp_address,omitempty" path:"smtp_address"`
	// SMTP server authentication type.
	SmtpAuthentication string `url:"smtp_authentication,omitempty" json:"smtp_authentication,omitempty" path:"smtp_authentication"`
	// From address to use when mailing through custom SMTP.
	SmtpFrom string `url:"smtp_from,omitempty" json:"smtp_from,omitempty" path:"smtp_from"`
	// SMTP server username.
	SmtpUsername string `url:"smtp_username,omitempty" json:"smtp_username,omitempty" path:"smtp_username"`
	// SMTP server port.
	SmtpPort int64 `url:"smtp_port,omitempty" json:"smtp_port,omitempty" path:"smtp_port"`
	// Custom SMTP encryption mode: if_available (default) uses STARTTLS when offered and otherwise sends credentials and messages unencrypted. A certificate-verified STARTTLS connection automatically changes if_available to require unless smtp_ssl is managed by a parent policy. require requires STARTTLS before authentication; require_implicit uses TLS from connection start; never disables TLS. TLS verifies the server certificate against smtp_address. Encryption is never automatically downgraded.
	SmtpSsl string `url:"smtp_ssl,omitempty" json:"smtp_ssl,omitempty" path:"smtp_ssl"`
	// Main LDAP setting: is LDAP enabled?
	LdapEnabled *bool `url:"ldap_enabled,omitempty" json:"ldap_enabled,omitempty" path:"ldap_enabled"`
	// LDAP type.
	LdapType string `url:"ldap_type,omitempty" json:"ldap_type,omitempty" path:"ldap_type"`
	// LDAP host.
	LdapHost string `url:"ldap_host,omitempty" json:"ldap_host,omitempty" path:"ldap_host"`
	// LDAP backup host.
	LdapHost2 string `url:"ldap_host_2,omitempty" json:"ldap_host_2,omitempty" path:"ldap_host_2"`
	// LDAP backup host.
	LdapHost3 string `url:"ldap_host_3,omitempty" json:"ldap_host_3,omitempty" path:"ldap_host_3"`
	// LDAP port.
	LdapPort int64 `url:"ldap_port,omitempty" json:"ldap_port,omitempty" path:"ldap_port"`
	// Use secure LDAP?
	LdapSecure *bool `url:"ldap_secure,omitempty" json:"ldap_secure,omitempty" path:"ldap_secure"`
	// How to validate the LDAP server certificate. `require_match` validates the certificate chain and hostname; `allow_any` disables certificate validation.
	LdapServerCertificate string `url:"ldap_server_certificate,omitempty" json:"ldap_server_certificate,omitempty" path:"ldap_server_certificate"`
	// Username for signing in to LDAP server.
	LdapUsername string `url:"ldap_username,omitempty" json:"ldap_username,omitempty" path:"ldap_username"`
	// LDAP username field.
	LdapUsernameField string `url:"ldap_username_field,omitempty" json:"ldap_username_field,omitempty" path:"ldap_username_field"`
	// Domain name that will be appended to usernames.
	LdapDomain string `url:"ldap_domain,omitempty" json:"ldap_domain,omitempty" path:"ldap_domain"`
	// Should we sync users from LDAP server?
	LdapUserAction string `url:"ldap_user_action,omitempty" json:"ldap_user_action,omitempty" path:"ldap_user_action"`
	// Should we sync groups from LDAP server?
	LdapGroupAction string `url:"ldap_group_action,omitempty" json:"ldap_group_action,omitempty" path:"ldap_group_action"`
	// Comma or newline separated list of group names (with optional wildcards) - if provided, only users in these groups will be added or synced.
	LdapUserIncludeGroups string `url:"ldap_user_include_groups,omitempty" json:"ldap_user_include_groups,omitempty" path:"ldap_user_include_groups"`
	// Comma or newline separated list of group names (with optional wildcards) to exclude when syncing.
	LdapGroupExclusion string `url:"ldap_group_exclusion,omitempty" json:"ldap_group_exclusion,omitempty" path:"ldap_group_exclusion"`
	// Comma or newline separated list of group names (with optional wildcards) to include when syncing.
	LdapGroupInclusion string `url:"ldap_group_inclusion,omitempty" json:"ldap_group_inclusion,omitempty" path:"ldap_group_inclusion"`
	// Base DN for looking up users in LDAP server.
	LdapBaseDn string `url:"ldap_base_dn,omitempty" json:"ldap_base_dn,omitempty" path:"ldap_base_dn"`
	// Require email authentication, virus, and spam checks for incoming emails to Inboxes and Incoming Email Automations in every Workspace on this site?
	UploadsViaEmailAuthentication *bool `url:"uploads_via_email_authentication,omitempty" json:"uploads_via_email_authentication,omitempty" path:"uploads_via_email_authentication"`
	// Preview watermark settings applied to all bundle items. Uses the same keys as Behavior.value.
	BundleWatermarkValue interface{} `url:"bundle_watermark_value,omitempty" json:"bundle_watermark_value,omitempty" path:"bundle_watermark_value"`
	Icon16File           io.Writer   `url:"icon16_file,omitempty" json:"icon16_file,omitempty" path:"icon16_file"`
	// Icon16FileAttachment sends file content as icon16_file; see lib.JSONAttachment. Leave Icon16File nil when you set it.
	Icon16FileAttachment *lib.JSONAttachment `json:"-" url:"-" path:"-" attachment:"icon16_file"`
	// If true, will delete the file stored in icon16.
	Icon16Delete *bool     `url:"icon16_delete,omitempty" json:"icon16_delete,omitempty" path:"icon16_delete"`
	Icon32File   io.Writer `url:"icon32_file,omitempty" json:"icon32_file,omitempty" path:"icon32_file"`
	// Icon32FileAttachment sends file content as icon32_file; see lib.JSONAttachment. Leave Icon32File nil when you set it.
	Icon32FileAttachment *lib.JSONAttachment `json:"-" url:"-" path:"-" attachment:"icon32_file"`
	// If true, will delete the file stored in icon32.
	Icon32Delete *bool     `url:"icon32_delete,omitempty" json:"icon32_delete,omitempty" path:"icon32_delete"`
	Icon48File   io.Writer `url:"icon48_file,omitempty" json:"icon48_file,omitempty" path:"icon48_file"`
	// Icon48FileAttachment sends file content as icon48_file; see lib.JSONAttachment. Leave Icon48File nil when you set it.
	Icon48FileAttachment *lib.JSONAttachment `json:"-" url:"-" path:"-" attachment:"icon48_file"`
	// If true, will delete the file stored in icon48.
	Icon48Delete *bool     `url:"icon48_delete,omitempty" json:"icon48_delete,omitempty" path:"icon48_delete"`
	Icon128File  io.Writer `url:"icon128_file,omitempty" json:"icon128_file,omitempty" path:"icon128_file"`
	// Icon128FileAttachment sends file content as icon128_file; see lib.JSONAttachment. Leave Icon128File nil when you set it.
	Icon128FileAttachment *lib.JSONAttachment `json:"-" url:"-" path:"-" attachment:"icon128_file"`
	// If true, will delete the file stored in icon128.
	Icon128Delete *bool     `url:"icon128_delete,omitempty" json:"icon128_delete,omitempty" path:"icon128_delete"`
	LogoFile      io.Writer `url:"logo_file,omitempty" json:"logo_file,omitempty" path:"logo_file"`
	// LogoFileAttachment sends file content as logo_file; see lib.JSONAttachment. Leave LogoFile nil when you set it.
	LogoFileAttachment *lib.JSONAttachment `json:"-" url:"-" path:"-" attachment:"logo_file"`
	// If true, will delete the file stored in logo.
	LogoDelete                    *bool     `url:"logo_delete,omitempty" json:"logo_delete,omitempty" path:"logo_delete"`
	BundleWatermarkAttachmentFile io.Writer `url:"bundle_watermark_attachment_file,omitempty" json:"bundle_watermark_attachment_file,omitempty" path:"bundle_watermark_attachment_file"`
	// BundleWatermarkAttachmentFileAttachment sends file content as bundle_watermark_attachment_file; see lib.JSONAttachment. Leave BundleWatermarkAttachmentFile nil when you set it.
	BundleWatermarkAttachmentFileAttachment *lib.JSONAttachment `json:"-" url:"-" path:"-" attachment:"bundle_watermark_attachment_file"`
	// If true, will delete the file stored in bundle_watermark_attachment.
	BundleWatermarkAttachmentDelete *bool     `url:"bundle_watermark_attachment_delete,omitempty" json:"bundle_watermark_attachment_delete,omitempty" path:"bundle_watermark_attachment_delete"`
	LoginPageBackgroundImageFile    io.Writer `url:"login_page_background_image_file,omitempty" json:"login_page_background_image_file,omitempty" path:"login_page_background_image_file"`
	// LoginPageBackgroundImageFileAttachment sends file content as login_page_background_image_file; see lib.JSONAttachment. Leave LoginPageBackgroundImageFile nil when you set it.
	LoginPageBackgroundImageFileAttachment *lib.JSONAttachment `json:"-" url:"-" path:"-" attachment:"login_page_background_image_file"`
	// If true, will delete the file stored in login_page_background_image.
	LoginPageBackgroundImageDelete *bool `url:"login_page_background_image_delete,omitempty" json:"login_page_background_image_delete,omitempty" path:"login_page_background_image_delete"`
	// If set to true, we will begin the process of disabling 2FA on this site.
	Disable2faWithDelay *bool `url:"disable_2fa_with_delay,omitempty" json:"disable_2fa_with_delay,omitempty" path:"disable_2fa_with_delay"`
	// New LDAP password.
	LdapPasswordChange string `url:"ldap_password_change,omitempty" json:"ldap_password_change,omitempty" path:"ldap_password_change"`
	// Confirm new LDAP password.
	LdapPasswordChangeConfirmation string `url:"ldap_password_change_confirmation,omitempty" json:"ldap_password_change_confirmation,omitempty" path:"ldap_password_change_confirmation"`
	// If true, and if changing the site subdomain, then create a redirect from the previous Files.com subdomain to the new Files.com subdomain.
	RedirectOldSubdomain *bool `url:"redirect_old_subdomain,omitempty" json:"redirect_old_subdomain,omitempty" path:"redirect_old_subdomain"`
	// Password for SMTP server.
	SmtpPassword string `url:"smtp_password,omitempty" json:"smtp_password,omitempty" path:"smtp_password"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (s *Site) UnmarshalJSON(data []byte) error {
	type site Site
	var v site
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*s = Site(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (s *SiteCollection) UnmarshalJSON(data []byte) error {
	type sites SiteCollection
	var v sites
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*s = SiteCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (s *SiteCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*s))
	for i, v := range *s {
		ret[i] = v
	}

	return &ret
}
