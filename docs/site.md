# Site

`files_sdk.Site` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/site`, named `site`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/site"
)
```

## Operations

### Get

`GET /site`

Show Site Settings

```go
func (c *Client) Get(opts ...files_sdk.RequestResponseOption) (site files_sdk.Site, err error)
func Get(opts ...files_sdk.RequestResponseOption) (site files_sdk.Site, err error)
```

Example:

```go
result, err := site.Get()
if err != nil {
	return err
}
fmt.Println(result)
```

### GetUsage

`GET /site/usage`

Get the most recent usage snapshot (usage data for billing purposes) for a Site

```go
func (c *Client) GetUsage(opts ...files_sdk.RequestResponseOption) (usageSnapshot files_sdk.UsageSnapshot, err error)
func GetUsage(opts ...files_sdk.RequestResponseOption) (usageSnapshot files_sdk.UsageSnapshot, err error)
```

Example:

```go
result, err := site.GetUsage()
if err != nil {
	return err
}
fmt.Println(result)
```

### Update

`PATCH /site`

Update Site Settings

```go
func (c *Client) Update(params files_sdk.SiteUpdateParams, opts ...files_sdk.RequestResponseOption) (site files_sdk.Site, err error)
func Update(params files_sdk.SiteUpdateParams, opts ...files_sdk.RequestResponseOption) (site files_sdk.Site, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (site files_sdk.Site, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (site files_sdk.Site, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.SiteUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Name` | `string` | `name` | no | Site name |
| `Subdomain` | `string` | `subdomain` | no | Site subdomain |
| `Domain` | `string` | `domain` | no | Custom domain |
| `DomainHstsHeader` | `*bool` | `domain_hsts_header` | no | Send HSTS (HTTP Strict Transport Security) header when visitors access the site via a custom domain? |
| `DomainLetsencryptChain` | `string` | `domain_letsencrypt_chain` | no | Letsencrypt chain to use when registering SSL Certificate for domain. No longer used as of 2026. |
| `Email` | `string` | `email` | no | Main email for this site |
| `ReplyToEmail` | `string` | `reply_to_email` | no | Reply-to email for this site |
| `AllowBundleNames` | `*bool` | `allow_bundle_names` | no | Are manual Bundle names allowed? |
| `BundleExpiration` | `int64` | `bundle_expiration` | no | Site-wide Bundle expiration in days |
| `WelcomeEmailEnabled` | `*bool` | `welcome_email_enabled` | no | Will the welcome email be sent to new users? |
| `AskAboutOverwrites` | `*bool` | `ask_about_overwrites` | no | If false, rename conflicting files instead of asking for overwrite confirmation. Only applies to web interface. |
| `ShowRequestAccessLink` | `*bool` | `show_request_access_link` | no | Show request access link for users without access? Currently unused. |
| `AlwaysMkdirParents` | `*bool` | `always_mkdir_parents` | no | Create parent directories if they do not exist during uploads? This is primarily used to work around broken upload clients that assume servers will perform this step. |
| `WelcomeEmailCc` | `string` | `welcome_email_cc` | no | Include this email in welcome emails if enabled |
| `WelcomeEmailSubject` | `string` | `welcome_email_subject` | no | Include this email subject in welcome emails if enabled |
| `WelcomeCustomText` | `string` | `welcome_custom_text` | no | Custom text send in user welcome email |
| `Language` | `string` | `language` | no | Site default language |
| `WindowsModeFtp` | `*bool` | `windows_mode_ftp` | no | Does FTP user Windows emulation mode? |
| `DefaultTimeZone` | `string` | `default_time_zone` | no | Site default time zone |
| `DesktopApp` | `*bool` | `desktop_app` | no | Is the desktop app enabled? |
| `DesktopAppSessionIpPinning` | `*bool` | `desktop_app_session_ip_pinning` | no | Is desktop app session IP pinning enabled? |
| `DesktopAppSessionLifetime` | `int64` | `desktop_app_session_lifetime` | no | Desktop app session lifetime (in hours) |
| `MobileApp` | `*bool` | `mobile_app` | no | Is the mobile app enabled? |
| `MobileAppSessionIpPinning` | `*bool` | `mobile_app_session_ip_pinning` | no | Is mobile app session IP pinning enabled? |
| `MobileAppSessionLifetime` | `int64` | `mobile_app_session_lifetime` | no | Mobile app session lifetime (in hours) |
| `FolderPermissionsGroupsOnly` | `*bool` | `folder_permissions_groups_only` | no | If true, permissions for this site must be bound to a group (not a user). |
| `WelcomeScreen` | `string` | `welcome_screen` | no | Does the welcome screen appear? |
| `OfficeIntegrationAvailable` | `*bool` | `office_integration_available` | no | If true, allows users to use a document editing integration. |
| `OfficeIntegrationType` | `string` | `office_integration_type` | no | Which document editing integration to support. Files.com Editor or Microsoft Office for the Web. |
| `PinAllRemoteServersToSiteRegion` | `*bool` | `pin_all_remote_servers_to_site_region` | no | If true, we will ensure that all internal communications with any remote server are made through the primary region of the site. This setting overrides individual remote server settings. |
| `MotdText` | `string` | `motd_text` | no | A message to show users when they connect via FTP or SFTP. |
| `MotdUseForFtp` | `*bool` | `motd_use_for_ftp` | no | Show message to users connecting via FTP |
| `MotdUseForSftp` | `*bool` | `motd_use_for_sftp` | no | Show message to users connecting via SFTP |
| `LeftNavigationVisibility` | `interface{}` | `left_navigation_visibility` | no | Visibility settings for account navigation |
| `DisableAllAiFeatures` | `*bool` | `disable_all_ai_features` | no | If true, all AI features are disabled for this site. |
| `AiFeatureAvailability` | `interface{}` | `ai_feature_availability` | no | Availability settings for AI features. Each feature requires the site_admins, workspace_admins, folder_admins, and all_users keys. Optional selected_group_members defaults to false; when true, members of at least one group in group_ids get access regardless of the other options. All availability options are additive: any enabled option matching the user grants access. Optional group_ids is an array of integer IDs of active groups on this site, from any workspace. Omitted or empty group_ids grants no access through selected_group_members and does not affect other options. Disabling all AI features overrides these settings. |
| `McpDcrEnabled` | `*bool` | `mcp_dcr_enabled` | no | Is OAuth DCR (dynamic client registration) for MCP enabled? |
| `AdditionalTextFileTypes` | `[]string` | `additional_text_file_types` | no | Additional extensions that are considered text files |
| `BundleRequireNote` | `*bool` | `bundle_require_note` | no | Do Bundles require internal notes? |
| `BundleSendSharedReceipts` | `*bool` | `bundle_send_shared_receipts` | no | Do Bundle creators receive receipts of invitations? |
| `BundlesDefaultOwnedByPrimaryGroup` | `*bool` | `bundles_default_owned_by_primary_group` | no | If true, new Share Links created by a user with a primary group will default to that group as owner. |
| `CalculateFileChecksumsCrc32` | `*bool` | `calculate_file_checksums_crc32` | no | Calculate CRC32 checksums for files? |
| `CalculateFileChecksumsMd5` | `*bool` | `calculate_file_checksums_md5` | no | Calculate MD5 checksums for files? |
| `CalculateFileChecksumsSha1` | `*bool` | `calculate_file_checksums_sha1` | no | Calculate SHA1 checksums for files? |
| `CalculateFileChecksumsSha256` | `*bool` | `calculate_file_checksums_sha256` | no | Calculate SHA256 checksums for files? |
| `LegacyChecksumsMode` | `*bool` | `legacy_checksums_mode` | no | Use legacy checksums mode? |
| `MigrateRemoteServerSyncToSync` | `*bool` | `migrate_remote_server_sync_to_sync` | no | If true, we will migrate all remote server syncs to the new Sync model. |
| `As2MessageRetentionDays` | `int64` | `as2_message_retention_days` | no | Number of days to retain AS2 messages (incoming and outgoing). |
| `FileActivityEventsEnabled` | `*bool` | `file_activity_events_enabled` | no | If true, file operations produce file_action.activity.v1 EventChannel events in all workspaces. Defaults to false; subscription filters and workspace access still control delivery. |
| `UsernameDisplay` | `string` | `username_display` | no | How usernames are displayed in the web UI. Can be `username_only`, `full_name_only`, `full_name_username`, `full_name_company`, or `full_name_username_company`. |
| `SessionExpiryMinutes` | `int64` | `session_expiry_minutes` | no | Session expiry in minutes |
| `SslRequired` | `*bool` | `ssl_required` | no | Is SSL required? Disabling this is insecure. |
| `SftpInsecureCiphers` | `*bool` | `sftp_insecure_ciphers` | no | If true, we will allow weak and known insecure ciphers to be used for SFTP connections. Enabling this setting severely weakens the security of your site and it is not recommend, except as a last resort for compatibility. |
| `SftpInsecureDiffieHellman` | `*bool` | `sftp_insecure_diffie_hellman` | no | If true, we will allow weak Diffie Hellman parameters to be used within ciphers for SFTP that are otherwise on our secure list. This has the effect of making the cipher weaker than our normal threshold for security, but is required to support certain legacy or broken SSH and MFT clients. Enabling this weakens security, but not nearly as much as enabling the full `sftp_insecure_ciphers` option. |
| `DisableFilesCertificateGeneration` | `*bool` | `disable_files_certificate_generation` | no | If set, Files.com will not set the CAA records required to generate future SSL certificates for this domain. |
| `Fedramp` | `*bool` | `fedramp` | no | Are FedRAMP security restrictions enabled for this site? |
| `UserLockout` | `*bool` | `user_lockout` | no | Will users be locked out after incorrect login attempts? |
| `UserLockoutTries` | `int64` | `user_lockout_tries` | no | Number of login tries within `user_lockout_within` hours before users are locked out |
| `UserLockoutWithin` | `int64` | `user_lockout_within` | no | Number of hours for user lockout window |
| `UserLockoutLockPeriod` | `int64` | `user_lockout_lock_period` | no | How many hours to lock user out for failed password? |
| `IncludePasswordInWelcomeEmail` | `*bool` | `include_password_in_welcome_email` | no | Include password in emails to new users? |
| `AllowedCountries` | `string` | `allowed_countries` | no | Comma separated list of allowed Country codes |
| `AllowedIps` | `string` | `allowed_ips` | no | List of allowed IP addresses |
| `AllowUserLevel2faOverride` | `*bool` | `allow_user_level_2fa_override` | no | Allow the site-wide two-factor authentication requirement to be overriden on a per-user-basis? |
| `AllowUserLevelAllowedIpOverride` | `*bool` | `allow_user_level_allowed_ip_override` | no | Allow the site-wide allowed IP restriction to be overriden on a per-user-basis? |
| `AllowUserLevelSslOverride` | `*bool` | `allow_user_level_ssl_override` | no | Allow the site-wide FTP SSL requirement to be overriden on a per-user-basis? |
| `EnforceParentSiteUserIpRestrictions` | `*bool` | `enforce_parent_site_user_ip_restrictions` | no | Require Parent Site users entering this Child Site to match its allowed IP addresses in addition to their existing IP restrictions. Applies to all Workspaces. Defaults to false. |
| `DisallowedCountries` | `string` | `disallowed_countries` | no | Comma separated list of disallowed Country codes |
| `DaysToRetainBackups` | `int64` | `days_to_retain_backups` | no | Number of days to keep deleted files |
| `MaxPriorPasswords` | `int64` | `max_prior_passwords` | no | Number of prior passwords to disallow |
| `PasswordValidityDays` | `int64` | `password_validity_days` | no | Number of days password is valid |
| `PasswordMinLength` | `int64` | `password_min_length` | no | Shortest password length for users |
| `PasswordRequireLetter` | `*bool` | `password_require_letter` | no | Require a letter in passwords? |
| `PasswordRequireMixed` | `*bool` | `password_require_mixed` | no | Require lower and upper case letters in passwords? |
| `PasswordRequireSpecial` | `*bool` | `password_require_special` | no | Require special characters in password? |
| `PasswordRequireNumber` | `*bool` | `password_require_number` | no | Require a number in passwords? |
| `PasswordRequireUnbreached` | `*bool` | `password_require_unbreached` | no | Require passwords that have not been previously breached? (see https://haveibeenpwned.com/) |
| `RequireLogoutFromBundlesAndInboxes` | `*bool` | `require_logout_from_bundles_and_inboxes` | no | If true, we will hide the 'Remember Me' box on Inbox and Bundle registration pages, requiring that the user logout and log back in every time they visit the page. |
| `DavUserRootEnabled` | `*bool` | `dav_user_root_enabled` | no | Use user FTP roots also for WebDAV? |
| `SftpUserRootEnabled` | `*bool` | `sftp_user_root_enabled` | no | Use user FTP roots also for SFTP? |
| `DisablePasswordReset` | `*bool` | `disable_password_reset` | no | Is password reset disabled? |
| `ImmutableFiles` | `*bool` | `immutable_files` | no | Are files protected from modification? |
| `BundleNotFoundMessage` | `string` | `bundle_not_found_message` | no | Custom error message to show when bundle is not found. |
| `BundlePasswordRequired` | `*bool` | `bundle_password_required` | no | Do Bundles require password protection? |
| `BundleRequireRegistration` | `*bool` | `bundle_require_registration` | no | Do Bundles require registration? |
| `BundleRequireShareRecipient` | `*bool` | `bundle_require_share_recipient` | no | Do Bundles require recipients for sharing? |
| `BundleSendOneTimePasswordToRecipientAtRegistration` | `*bool` | `bundle_send_one_time_password_to_recipient_at_registration` | no | If true, new Share Links must send a one-time password to the recipient when they register. Requires bundle_require_share_recipient and cannot be enabled with bundle_password_required. |
| `BundleRegistrationNotifications` | `string` | `bundle_registration_notifications` | no | Do Bundle owners receive registration notification? |
| `BundleActivityNotifications` | `string` | `bundle_activity_notifications` | no | Do Bundle owners receive activity notifications? |
| `BundleUploadReceiptNotifications` | `string` | `bundle_upload_receipt_notifications` | no | Do Bundle uploaders receive upload confirmation notifications? |
| `DocumentEditsInBundleAllowed` | `*bool` | `document_edits_in_bundle_allowed` | no | If true, allow public viewers of Bundles with full permissions to use document editing integrations. |
| `PasswordRequirementsApplyToBundles` | `*bool` | `password_requirements_apply_to_bundles` | no | Require bundles' passwords, and passwords for other items (inboxes, public shares, etc.) to conform to the same requirements as users' passwords? |
| `PreventRootPermissionsForNonSiteAdmins` | `*bool` | `prevent_root_permissions_for_non_site_admins` | no | If true, we will prevent non-administrators from receiving any permissions directly on the root folder. This is commonly used to prevent the accidental application of permissions. |
| `RestrictRootFolderBehaviorsToSiteAdmins` | `*bool` | `restrict_root_folder_behaviors_to_site_admins` | no | If true, only site admins may create, modify, or delete any behavior at the site root, or a skip that would disable one. |
| `RootFolderBehaviorsApplyToWorkspaces` | `*bool` | `root_folder_behaviors_apply_to_workspaces` | no | If true, supported protective behaviors at the site root also apply within named workspaces. Requires restrict_root_folder_behaviors_to_site_admins to be enabled. |
| `OptOutGlobal` | `*bool` | `opt_out_global` | no | Use servers in the USA only? |
| `UseProvidedModifiedAt` | `*bool` | `use_provided_modified_at` | no | Allow uploaders to set `provided_modified_at` for uploaded files? |
| `CustomNamespace` | `*bool` | `custom_namespace` | no | Is this site using a custom namespace for users? |
| `NonSsoGroupsAllowed` | `*bool` | `non_sso_groups_allowed` | no | If true, groups can be manually created / modified / deleted by Site Admins. Otherwise, groups can only be managed via your SSO provider. |
| `NonSsoUsersAllowed` | `*bool` | `non_sso_users_allowed` | no | If true, users can be manually created / modified / deleted by Site Admins. Otherwise, users can only be managed via your SSO provider. |
| `SharingEnabled` | `*bool` | `sharing_enabled` | no | Allow bundle creation |
| `SnapshotSharingEnabled` | `*bool` | `snapshot_sharing_enabled` | no | Allow snapshot share links creation |
| `UserRequestsEnabled` | `*bool` | `user_requests_enabled` | no | Enable User Requests feature |
| `UserRequestsNotifyAdmins` | `*bool` | `user_requests_notify_admins` | no | Send email to site admins when a user request is received? |
| `DavEnabled` | `*bool` | `dav_enabled` | no | Is WebDAV enabled? |
| `FilesComRemoteServerEnabled` | `*bool` | `files_com_remote_server_enabled` | no | Allow other Files.com sites to use this site's API keys for native Files.com Remote Server connections? Defaults to true and applies to all keys and workspaces on this site. When false, new pairings and access through existing connections are rejected without revoking keys or deleting connections. Re-enabling permits access again with usable keys. Does not disable this site's connections to other sites, ordinary API access, or Connected Sites. |
| `FtpEnabled` | `*bool` | `ftp_enabled` | no | Is FTP enabled? |
| `S3CompatibleEndpointEnabled` | `*bool` | `s3_compatible_endpoint_enabled` | no | Is the S3-compatible endpoint enabled for all users and workspaces on this site? Defaults to true. When false, user and group S3 permissions do not allow access. |
| `SftpEnabled` | `*bool` | `sftp_enabled` | no | Is SFTP enabled? |
| `SftpFinalizePartialUploads` | `*bool` | `sftp_finalize_partial_uploads` | no | Finalize partial SFTP uploads from interrupted connections? Default: true. |
| `UsersCanCreateApiKeys` | `*bool` | `users_can_create_api_keys` | no | Allow users to create their own API keys? |
| `UsersCanCreateSshKeys` | `*bool` | `users_can_create_ssh_keys` | no | Allow users to create their own SSH keys? |
| `ShowUserNotificationsLogInLink` | `*bool` | `show_user_notifications_log_in_link` | no | Show log in link in user notifications? |
| `SftpHostKeyType` | `string` | `sftp_host_key_type` | no | Sftp Host Key Type |
| `ActiveSftpHostKeyId` | `int64` | `active_sftp_host_key_id` | no | Id of the currently selected custom SFTP Host Key |
| `ActiveSftpHostKeyIds` | `[]int64` | `active_sftp_host_key_ids` | no | Ids of the selected custom SFTP Host Keys |
| `ProtocolAccessGroupsOnly` | `*bool` | `protocol_access_groups_only` | no | If true, protocol access permissions on users will be ignored, and only protocol access permissions set on Groups will be honored. Make sure that your current user is a member of a group with API permission when changing this value to avoid locking yourself out of your site. |
| `RevokeBundleAccessOnDisableOrDelete` | `*bool` | `revoke_bundle_access_on_disable_or_delete` | no | Auto-removes bundles for disabled/deleted users and enforces bundle expiry within user access period. |
| `GroupAdminsCanAddUsers` | `*bool` | `group_admins_can_add_users` | no | Allow group admins to create users in their groups |
| `GroupAdminsCanManageGroupMemberships` | `*bool` | `group_admins_can_manage_group_memberships` | no | Allow group admins to add or remove existing users in their groups |
| `GroupAdminsCanDeleteUsers` | `*bool` | `group_admins_can_delete_users` | no | Allow group admins to delete users in their groups |
| `GroupAdminsCanEnableDisableUsers` | `*bool` | `group_admins_can_enable_disable_users` | no | Allow group admins to enable or disable users in their groups |
| `GroupAdminsCanModifyUsers` | `*bool` | `group_admins_can_modify_users` | no | Allow group admins to modify users in their groups |
| `GroupAdminsCanBypassUserLifecycleRules` | `*bool` | `group_admins_can_bypass_user_lifecycle_rules` | no | Allow group admins to exempt users in their groups from lifecycle rules |
| `GroupAdminsCanResetPasswords` | `*bool` | `group_admins_can_reset_passwords` | no | Allow group admins to reset passwords for users in their groups |
| `GroupAdminsCanSetUserPassword` | `*bool` | `group_admins_can_set_user_password` | no | Allow group admins to set password authentication method |
| `ChildSiteFolderAdminsCanSeeParentSiteDirectory` | `*bool` | `child_site_folder_admins_can_see_parent_site_directory` | no | If true, administrators on this site's Child Sites, including Folder Admins, can list this site's Default Workspace users and groups and grant them folder permissions and notifications on Child Site paths. |
| `WorkspaceFolderAdminsCanSeeDefaultWorkspaceDirectory` | `*bool` | `workspace_folder_admins_can_see_default_workspace_directory` | no | If true, Workspace Administrators and Folder Admins in a Custom Workspace can list the Default Workspace users and groups that have permissions in that Workspace and choose them as Email Notification recipients. |
| `BundleRecipientBlacklistFreeEmailDomains` | `*bool` | `bundle_recipient_blacklist_free_email_domains` | no | Disallow free email domains for Bundle/Inbox recipients? |
| `BundleRecipientBlacklistDomains` | `[]string` | `bundle_recipient_blacklist_domains` | no | List of email domains to disallow when entering a Bundle/Inbox recipients |
| `AdminsBypassLockedSubfolders` | `*bool` | `admins_bypass_locked_subfolders` | no | Allow admins to bypass the locked subfolders setting. |
| `Allowed2faMethodSms` | `*bool` | `allowed_2fa_method_sms` | no | Is SMS two factor authentication allowed? |
| `Allowed2faMethodTotp` | `*bool` | `allowed_2fa_method_totp` | no | Is TOTP two factor authentication allowed? |
| `Allowed2faMethodWebauthn` | `*bool` | `allowed_2fa_method_webauthn` | no | Is WebAuthn two factor authentication allowed? |
| `Allowed2faMethodYubi` | `*bool` | `allowed_2fa_method_yubi` | no | Is yubikey two factor authentication allowed? |
| `Allowed2faMethodEmail` | `*bool` | `allowed_2fa_method_email` | no | Is OTP via email two factor authentication allowed? |
| `Allowed2faMethodStatic` | `*bool` | `allowed_2fa_method_static` | no | Is OTP via static codes for two factor authentication allowed? |
| `Allowed2faMethodBypassForFtpSftpDav` | `*bool` | `allowed_2fa_method_bypass_for_ftp_sftp_dav` | no | Are users allowed to configure their two factor authentication to be bypassed for FTP/SFTP/WebDAV? |
| `Require2fa` | `*bool` | `require_2fa` | no | Require two-factor authentication for all users? |
| `Require2faExemptAllSsoUsers` | `*bool` | `require_2fa_exempt_all_sso_users` | no | If true, SSO users using the default user-level two-factor authentication setting are exempt from the site-wide two-factor authentication requirement. |
| `Require2faUserType` | `string` | `require_2fa_user_type` | no | What type of user is required to use two-factor authentication (when require_2fa is set to `true` for this site)? |
| `Color2Top` | `string` | `color2_top` | no | Top bar background color |
| `Color2Left` | `string` | `color2_left` | no | Page link and button color |
| `Color2Link` | `string` | `color2_link` | no | Top bar link color |
| `Color2Text` | `string` | `color2_text` | no | Page link and button color |
| `Color2TopText` | `string` | `color2_top_text` | no | Top bar text color |
| `SiteHeader` | `string` | `site_header` | no | Custom site header text for authenticated pages |
| `SiteFooter` | `string` | `site_footer` | no | Custom site footer text for authenticated pages |
| `SitePublicHeader` | `string` | `site_public_header` | no | Custom site header text for public pages |
| `SitePublicFooter` | `string` | `site_public_footer` | no | Custom site footer text for public pages |
| `LoginHelpText` | `string` | `login_help_text` | no | Login help text |
| `UseDedicatedIpsForSmtp` | `*bool` | `use_dedicated_ips_for_smtp` | no | If using custom SMTP, should we use dedicated IPs to deliver emails? |
| `EmailFooterCustomText` | `string` | `email_footer_custom_text` | no | Custom footer text for system-generated emails (as Markdown). Supports standard strftime date/time patterns like %Y (4-digit year), %m (month), %d (day). |
| `SmtpAddress` | `string` | `smtp_address` | no | SMTP server hostname or IP |
| `SmtpAuthentication` | `string` | `smtp_authentication` | no | SMTP server authentication type |
| `SmtpFrom` | `string` | `smtp_from` | no | From address to use when mailing through custom SMTP |
| `SmtpUsername` | `string` | `smtp_username` | no | SMTP server username |
| `SmtpPort` | `int64` | `smtp_port` | no | SMTP server port |
| `SmtpSsl` | `string` | `smtp_ssl` | no | Custom SMTP encryption mode: if_available (default) uses STARTTLS when offered and otherwise sends credentials and messages unencrypted. A certificate-verified STARTTLS connection automatically changes if_available to require unless smtp_ssl is managed by a parent policy. require requires STARTTLS before authentication; require_implicit uses TLS from connection start; never disables TLS. TLS verifies the server certificate against smtp_address. Encryption is never automatically downgraded. |
| `LdapEnabled` | `*bool` | `ldap_enabled` | no | Main LDAP setting: is LDAP enabled? |
| `LdapType` | `string` | `ldap_type` | no | LDAP type |
| `LdapHost` | `string` | `ldap_host` | no | LDAP host |
| `LdapHost2` | `string` | `ldap_host_2` | no | LDAP backup host |
| `LdapHost3` | `string` | `ldap_host_3` | no | LDAP backup host |
| `LdapPort` | `int64` | `ldap_port` | no | LDAP port |
| `LdapSecure` | `*bool` | `ldap_secure` | no | Use secure LDAP? |
| `LdapServerCertificate` | `string` | `ldap_server_certificate` | no | How to validate the LDAP server certificate. `require_match` validates the certificate chain and hostname; `allow_any` disables certificate validation. |
| `LdapUsername` | `string` | `ldap_username` | no | Username for signing in to LDAP server. |
| `LdapUsernameField` | `string` | `ldap_username_field` | no | LDAP username field |
| `LdapDomain` | `string` | `ldap_domain` | no | Domain name that will be appended to usernames |
| `LdapUserAction` | `string` | `ldap_user_action` | no | Should we sync users from LDAP server? |
| `LdapGroupAction` | `string` | `ldap_group_action` | no | Should we sync groups from LDAP server? |
| `LdapUserIncludeGroups` | `string` | `ldap_user_include_groups` | no | Comma or newline separated list of group names (with optional wildcards) - if provided, only users in these groups will be added or synced. |
| `LdapGroupExclusion` | `string` | `ldap_group_exclusion` | no | Comma or newline separated list of group names (with optional wildcards) to exclude when syncing. |
| `LdapGroupInclusion` | `string` | `ldap_group_inclusion` | no | Comma or newline separated list of group names (with optional wildcards) to include when syncing. |
| `LdapBaseDn` | `string` | `ldap_base_dn` | no | Base DN for looking up users in LDAP server |
| `UploadsViaEmailAuthentication` | `*bool` | `uploads_via_email_authentication` | no | Require email authentication, virus, and spam checks for incoming emails to Inboxes and Incoming Email Automations in every Workspace on this site? |
| `BundleWatermarkValue` | `interface{}` | `bundle_watermark_value` | no | Preview watermark settings applied to all bundle items. Uses the same keys as Behavior.value |
| `Icon16File` | `io.Writer` | `icon16_file` | no |  |
| `Icon16FileAttachment` | `*lib.JSONAttachment` | `icon16_file` | no | Sends file content as icon16_file; see lib.JSONAttachment. Leave Icon16File nil when you set it. |
| `Icon16Delete` | `*bool` | `icon16_delete` | no | If true, will delete the file stored in icon16 |
| `Icon32File` | `io.Writer` | `icon32_file` | no |  |
| `Icon32FileAttachment` | `*lib.JSONAttachment` | `icon32_file` | no | Sends file content as icon32_file; see lib.JSONAttachment. Leave Icon32File nil when you set it. |
| `Icon32Delete` | `*bool` | `icon32_delete` | no | If true, will delete the file stored in icon32 |
| `Icon48File` | `io.Writer` | `icon48_file` | no |  |
| `Icon48FileAttachment` | `*lib.JSONAttachment` | `icon48_file` | no | Sends file content as icon48_file; see lib.JSONAttachment. Leave Icon48File nil when you set it. |
| `Icon48Delete` | `*bool` | `icon48_delete` | no | If true, will delete the file stored in icon48 |
| `Icon128File` | `io.Writer` | `icon128_file` | no |  |
| `Icon128FileAttachment` | `*lib.JSONAttachment` | `icon128_file` | no | Sends file content as icon128_file; see lib.JSONAttachment. Leave Icon128File nil when you set it. |
| `Icon128Delete` | `*bool` | `icon128_delete` | no | If true, will delete the file stored in icon128 |
| `LogoFile` | `io.Writer` | `logo_file` | no |  |
| `LogoFileAttachment` | `*lib.JSONAttachment` | `logo_file` | no | Sends file content as logo_file; see lib.JSONAttachment. Leave LogoFile nil when you set it. |
| `LogoDelete` | `*bool` | `logo_delete` | no | If true, will delete the file stored in logo |
| `BundleWatermarkAttachmentFile` | `io.Writer` | `bundle_watermark_attachment_file` | no |  |
| `BundleWatermarkAttachmentFileAttachment` | `*lib.JSONAttachment` | `bundle_watermark_attachment_file` | no | Sends file content as bundle_watermark_attachment_file; see lib.JSONAttachment. Leave BundleWatermarkAttachmentFile nil when you set it. |
| `BundleWatermarkAttachmentDelete` | `*bool` | `bundle_watermark_attachment_delete` | no | If true, will delete the file stored in bundle_watermark_attachment |
| `LoginPageBackgroundImageFile` | `io.Writer` | `login_page_background_image_file` | no |  |
| `LoginPageBackgroundImageFileAttachment` | `*lib.JSONAttachment` | `login_page_background_image_file` | no | Sends file content as login_page_background_image_file; see lib.JSONAttachment. Leave LoginPageBackgroundImageFile nil when you set it. |
| `LoginPageBackgroundImageDelete` | `*bool` | `login_page_background_image_delete` | no | If true, will delete the file stored in login_page_background_image |
| `Disable2faWithDelay` | `*bool` | `disable_2fa_with_delay` | no | If set to true, we will begin the process of disabling 2FA on this site. |
| `LdapPasswordChange` | `string` | `ldap_password_change` | no | New LDAP password. |
| `LdapPasswordChangeConfirmation` | `string` | `ldap_password_change_confirmation` | no | Confirm new LDAP password. |
| `RedirectOldSubdomain` | `*bool` | `redirect_old_subdomain` | no | If true, and if changing the site subdomain, then create a redirect from the previous Files.com subdomain to the new Files.com subdomain. |
| `SmtpPassword` | `string` | `smtp_password` | no | Password for SMTP server. |

Example:

```go
result, err := site.Update(files_sdk.SiteUpdateParams{})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := site.UpdateWithMap(map[string]interface{}{})
if err != nil {
	return err
}
fmt.Println(result)
```

## Model

`files_sdk.Site` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Site Id |
| `Name` | `string` | `name` | Site name |
| `AdditionalTextFileTypes` | `[]string` | `additional_text_file_types` | Additional extensions that are considered text files |
| `AiFeatureAvailability` | `interface{}` | `ai_feature_availability` | Availability settings for AI features. Each feature requires the site_admins, workspace_admins, folder_admins, and all_users keys. Optional selected_group_members defaults to false; when true, members of at least one group in group_ids get access regardless of the other options. All availability options are additive: any enabled option matching the user grants access. Optional group_ids is an array of integer IDs of active groups on this site, from any workspace. Omitted or empty group_ids grants no access through selected_group_members and does not affect other options. Disabling all AI features overrides these settings. |
| `Allowed2faMethodSms` | `*bool` | `allowed_2fa_method_sms` | Is SMS two factor authentication allowed? |
| `Allowed2faMethodTotp` | `*bool` | `allowed_2fa_method_totp` | Is TOTP two factor authentication allowed? |
| `Allowed2faMethodWebauthn` | `*bool` | `allowed_2fa_method_webauthn` | Is WebAuthn two factor authentication allowed? |
| `Allowed2faMethodYubi` | `*bool` | `allowed_2fa_method_yubi` | Is yubikey two factor authentication allowed? |
| `Allowed2faMethodEmail` | `*bool` | `allowed_2fa_method_email` | Is OTP via email two factor authentication allowed? |
| `Allowed2faMethodStatic` | `*bool` | `allowed_2fa_method_static` | Is OTP via static codes for two factor authentication allowed? |
| `Allowed2faMethodBypassForFtpSftpDav` | `*bool` | `allowed_2fa_method_bypass_for_ftp_sftp_dav` | Are users allowed to configure their two factor authentication to be bypassed for FTP/SFTP/WebDAV? |
| `AdminUserId` | `int64` | `admin_user_id` | User ID for the main site administrator |
| `AdminsBypassLockedSubfolders` | `*bool` | `admins_bypass_locked_subfolders` | Allow admins to bypass the locked subfolders setting. |
| `AllowBundleNames` | `*bool` | `allow_bundle_names` | Are manual Bundle names allowed? |
| `AllowUserLevel2faOverride` | `*bool` | `allow_user_level_2fa_override` | Allow the site-wide two-factor authentication requirement to be overriden on a per-user-basis? |
| `AllowUserLevelAllowedIpOverride` | `*bool` | `allow_user_level_allowed_ip_override` | Allow the site-wide allowed IP restriction to be overriden on a per-user-basis? |
| `AllowUserLevelSslOverride` | `*bool` | `allow_user_level_ssl_override` | Allow the site-wide FTP SSL requirement to be overriden on a per-user-basis? |
| `AllowedCountries` | `string` | `allowed_countries` | Comma separated list of allowed Country codes |
| `AllowedIps` | `string` | `allowed_ips` | List of allowed IP addresses |
| `EnforceParentSiteUserIpRestrictions` | `*bool` | `enforce_parent_site_user_ip_restrictions` | Require Parent Site users entering this Child Site to match its allowed IP addresses in addition to their existing IP restrictions. Applies to all Workspaces. Defaults to false. |
| `AlwaysMkdirParents` | `*bool` | `always_mkdir_parents` | Create parent directories if they do not exist during uploads? This is primarily used to work around broken upload clients that assume servers will perform this step. |
| `As2MessageRetentionDays` | `int64` | `as2_message_retention_days` | Number of days to retain AS2 messages (incoming and outgoing). |
| `AskAboutOverwrites` | `*bool` | `ask_about_overwrites` | If false, rename conflicting files instead of asking for overwrite confirmation. Only applies to web interface. |
| `BundleActivityNotifications` | `string` | `bundle_activity_notifications` | Do Bundle owners receive activity notifications? |
| `BundleExpiration` | `int64` | `bundle_expiration` | Site-wide Bundle expiration in days |
| `BundleNotFoundMessage` | `string` | `bundle_not_found_message` | Custom error message to show when bundle is not found. |
| `BundlePasswordRequired` | `*bool` | `bundle_password_required` | Do Bundles require password protection? |
| `BundlesDefaultOwnedByPrimaryGroup` | `*bool` | `bundles_default_owned_by_primary_group` | If true, new Share Links created by a user with a primary group will default to that group as owner. |
| `BundleRecipientBlacklistDomains` | `[]string` | `bundle_recipient_blacklist_domains` | List of email domains to disallow when entering a Bundle/Inbox recipients |
| `BundleRecipientBlacklistFreeEmailDomains` | `*bool` | `bundle_recipient_blacklist_free_email_domains` | Disallow free email domains for Bundle/Inbox recipients? |
| `BundleRegistrationNotifications` | `string` | `bundle_registration_notifications` | Do Bundle owners receive registration notification? |
| `BundleRequireRegistration` | `*bool` | `bundle_require_registration` | Do Bundles require registration? |
| `BundleRequireShareRecipient` | `*bool` | `bundle_require_share_recipient` | Do Bundles require recipients for sharing? |
| `BundleSendOneTimePasswordToRecipientAtRegistration` | `*bool` | `bundle_send_one_time_password_to_recipient_at_registration` | If true, new Share Links must send a one-time password to the recipient when they register. Requires bundle_require_share_recipient and cannot be enabled with bundle_password_required. |
| `BundleRequireNote` | `*bool` | `bundle_require_note` | Do Bundles require internal notes? |
| `BundleSendSharedReceipts` | `*bool` | `bundle_send_shared_receipts` | Do Bundle creators receive receipts of invitations? |
| `BundleUploadReceiptNotifications` | `string` | `bundle_upload_receipt_notifications` | Do Bundle uploaders receive upload confirmation notifications? |
| `BundleWatermarkAttachment` | `files_sdk.Image` | `bundle_watermark_attachment` | Preview watermark image applied to all bundle items. |
| `BundleWatermarkValue` | `interface{}` | `bundle_watermark_value` | Preview watermark settings applied to all bundle items. Uses the same keys as Behavior.value |
| `CalculateFileChecksumsCrc32` | `*bool` | `calculate_file_checksums_crc32` | Calculate CRC32 checksums for files? |
| `CalculateFileChecksumsMd5` | `*bool` | `calculate_file_checksums_md5` | Calculate MD5 checksums for files? |
| `CalculateFileChecksumsSha1` | `*bool` | `calculate_file_checksums_sha1` | Calculate SHA1 checksums for files? |
| `CalculateFileChecksumsSha256` | `*bool` | `calculate_file_checksums_sha256` | Calculate SHA256 checksums for files? |
| `ChildSiteFolderAdminsCanSeeParentSiteDirectory` | `*bool` | `child_site_folder_admins_can_see_parent_site_directory` | If true, administrators on this site's Child Sites, including Folder Admins, can list this site's Default Workspace users and groups and grant them folder permissions and notifications on Child Site paths. |
| `WorkspaceFolderAdminsCanSeeDefaultWorkspaceDirectory` | `*bool` | `workspace_folder_admins_can_see_default_workspace_directory` | If true, Workspace Administrators and Folder Admins in a Custom Workspace can list the Default Workspace users and groups that have permissions in that Workspace and choose them as Email Notification recipients. |
| `UploadsViaEmailAuthentication` | `*bool` | `uploads_via_email_authentication` | Require email authentication, virus, and spam checks for incoming emails to Inboxes and Incoming Email Automations in every Workspace on this site? |
| `Color2Left` | `string` | `color2_left` | Page link and button color |
| `Color2Link` | `string` | `color2_link` | Top bar link color |
| `Color2Text` | `string` | `color2_text` | Page link and button color |
| `Color2Top` | `string` | `color2_top` | Top bar background color |
| `Color2TopText` | `string` | `color2_top_text` | Top bar text color |
| `ContactName` | `string` | `contact_name` | Site main contact name |
| `CreatedAt` | `*time.Time` | `created_at` | Time this site was created |
| `Currency` | `string` | `currency` | Preferred currency |
| `CustomNamespace` | `*bool` | `custom_namespace` | Is this site using a custom namespace for users? |
| `DavEnabled` | `*bool` | `dav_enabled` | Is WebDAV enabled? |
| `DavUserRootEnabled` | `*bool` | `dav_user_root_enabled` | Use user FTP roots also for WebDAV? |
| `DaysToRetainBackups` | `int64` | `days_to_retain_backups` | Number of days to keep deleted files |
| `DocumentEditsInBundleAllowed` | `*bool` | `document_edits_in_bundle_allowed` | If true, allow public viewers of Bundles with full permissions to use document editing integrations. |
| `DefaultTimeZone` | `string` | `default_time_zone` | Site default time zone |
| `DesktopApp` | `*bool` | `desktop_app` | Is the desktop app enabled? |
| `DesktopAppSessionIpPinning` | `*bool` | `desktop_app_session_ip_pinning` | Is desktop app session IP pinning enabled? |
| `DesktopAppSessionLifetime` | `int64` | `desktop_app_session_lifetime` | Desktop app session lifetime (in hours) |
| `LegacyChecksumsMode` | `*bool` | `legacy_checksums_mode` | Use legacy checksums mode? |
| `MigrateRemoteServerSyncToSync` | `*bool` | `migrate_remote_server_sync_to_sync` | If true, we will migrate all remote server syncs to the new Sync model. |
| `McpDcrEnabled` | `*bool` | `mcp_dcr_enabled` | Is OAuth DCR (dynamic client registration) for MCP enabled? |
| `MobileApp` | `*bool` | `mobile_app` | Is the mobile app enabled? |
| `MobileAppSessionIpPinning` | `*bool` | `mobile_app_session_ip_pinning` | Is mobile app session IP pinning enabled? |
| `MobileAppSessionLifetime` | `int64` | `mobile_app_session_lifetime` | Mobile app session lifetime (in hours) |
| `DisallowedCountries` | `string` | `disallowed_countries` | Comma separated list of disallowed Country codes |
| `DisableAllAiFeatures` | `*bool` | `disable_all_ai_features` | If true, all AI features are disabled for this site. |
| `DisableFilesCertificateGeneration` | `*bool` | `disable_files_certificate_generation` | If set, Files.com will not set the CAA records required to generate future SSL certificates for this domain. |
| `DisableNotifications` | `*bool` | `disable_notifications` | Are notifications disabled? |
| `FileActivityEventsEnabled` | `*bool` | `file_activity_events_enabled` | If true, file operations produce file_action.activity.v1 EventChannel events in all workspaces. Defaults to false; subscription filters and workspace access still control delivery. |
| `DisablePasswordReset` | `*bool` | `disable_password_reset` | Is password reset disabled? |
| `Domain` | `string` | `domain` | Custom domain |
| `DomainHstsHeader` | `*bool` | `domain_hsts_header` | Send HSTS (HTTP Strict Transport Security) header when visitors access the site via a custom domain? |
| `DomainLetsencryptChain` | `string` | `domain_letsencrypt_chain` | Letsencrypt chain to use when registering SSL Certificate for domain. No longer used as of 2026. |
| `Email` | `string` | `email` | Main email for this site |
| `Fedramp` | `*bool` | `fedramp` | Are FedRAMP security restrictions enabled for this site? |
| `FilesComRemoteServerEnabled` | `*bool` | `files_com_remote_server_enabled` | Allow other Files.com sites to use this site's API keys for native Files.com Remote Server connections? Defaults to true and applies to all keys and workspaces on this site. When false, new pairings and access through existing connections are rejected without revoking keys or deleting connections. Re-enabling permits access again with usable keys. Does not disable this site's connections to other sites, ordinary API access, or Connected Sites. |
| `FtpEnabled` | `*bool` | `ftp_enabled` | Is FTP enabled? |
| `ReplyToEmail` | `string` | `reply_to_email` | Reply-to email for this site |
| `NonSsoGroupsAllowed` | `*bool` | `non_sso_groups_allowed` | If true, groups can be manually created / modified / deleted by Site Admins. Otherwise, groups can only be managed via your SSO provider. |
| `NonSsoUsersAllowed` | `*bool` | `non_sso_users_allowed` | If true, users can be manually created / modified / deleted by Site Admins. Otherwise, users can only be managed via your SSO provider. |
| `FolderPermissionsGroupsOnly` | `*bool` | `folder_permissions_groups_only` | If true, permissions for this site must be bound to a group (not a user). |
| `GroupAdminsCanAddUsers` | `*bool` | `group_admins_can_add_users` | Allow group admins to create users in their groups |
| `GroupAdminsCanManageGroupMemberships` | `*bool` | `group_admins_can_manage_group_memberships` | Allow group admins to add or remove existing users in their groups |
| `GroupAdminsCanDeleteUsers` | `*bool` | `group_admins_can_delete_users` | Allow group admins to delete users in their groups |
| `GroupAdminsCanEnableDisableUsers` | `*bool` | `group_admins_can_enable_disable_users` | Allow group admins to enable or disable users in their groups |
| `GroupAdminsCanModifyUsers` | `*bool` | `group_admins_can_modify_users` | Allow group admins to modify users in their groups |
| `GroupAdminsCanBypassUserLifecycleRules` | `*bool` | `group_admins_can_bypass_user_lifecycle_rules` | Allow group admins to exempt users in their groups from lifecycle rules |
| `GroupAdminsCanResetPasswords` | `*bool` | `group_admins_can_reset_passwords` | Allow group admins to reset passwords for users in their groups |
| `GroupAdminsCanSetUserPassword` | `*bool` | `group_admins_can_set_user_password` | Allow group admins to set password authentication method |
| `Hipaa` | `*bool` | `hipaa` | Is there a signed HIPAA BAA between Files.com and this site? |
| `Icon128` | `files_sdk.Image` | `icon128` | Branded icon 128x128 |
| `Icon16` | `files_sdk.Image` | `icon16` | Branded icon 16x16 |
| `Icon32` | `files_sdk.Image` | `icon32` | Branded icon 32x32 |
| `Icon48` | `files_sdk.Image` | `icon48` | Branded icon 48x48 |
| `ImmutableFilesSetAt` | `*time.Time` | `immutable_files_set_at` | Can files be modified? |
| `IncludePasswordInWelcomeEmail` | `*bool` | `include_password_in_welcome_email` | Include password in emails to new users? |
| `Language` | `string` | `language` | Site default language |
| `LdapBaseDn` | `string` | `ldap_base_dn` | Base DN for looking up users in LDAP server |
| `LdapDomain` | `string` | `ldap_domain` | Domain name that will be appended to usernames |
| `LdapEnabled` | `*bool` | `ldap_enabled` | Main LDAP setting: is LDAP enabled? |
| `LdapGroupAction` | `string` | `ldap_group_action` | Should we sync groups from LDAP server? |
| `LdapGroupExclusion` | `string` | `ldap_group_exclusion` | Comma or newline separated list of group names (with optional wildcards) to exclude when syncing. |
| `LdapGroupInclusion` | `string` | `ldap_group_inclusion` | Comma or newline separated list of group names (with optional wildcards) to include when syncing. |
| `LdapHost` | `string` | `ldap_host` | LDAP host |
| `LdapHost2` | `string` | `ldap_host_2` | LDAP backup host |
| `LdapHost3` | `string` | `ldap_host_3` | LDAP backup host |
| `LdapPort` | `int64` | `ldap_port` | LDAP port |
| `LdapSecure` | `*bool` | `ldap_secure` | Use secure LDAP? |
| `LdapServerCertificate` | `string` | `ldap_server_certificate` | How to validate the LDAP server certificate. `require_match` validates the certificate chain and hostname; `allow_any` disables certificate validation. |
| `LdapType` | `string` | `ldap_type` | LDAP type |
| `LdapUserAction` | `string` | `ldap_user_action` | Should we sync users from LDAP server? |
| `LdapUserIncludeGroups` | `string` | `ldap_user_include_groups` | Comma or newline separated list of group names (with optional wildcards) - if provided, only users in these groups will be added or synced. |
| `LdapUsername` | `string` | `ldap_username` | Username for signing in to LDAP server. |
| `LdapUsernameField` | `string` | `ldap_username_field` | LDAP username field |
| `LoginHelpText` | `string` | `login_help_text` | Login help text |
| `Logo` | `files_sdk.Image` | `logo` | Branded logo |
| `LoginPageBackgroundImage` | `files_sdk.Image` | `login_page_background_image` | Branded login page background |
| `MaxPriorPasswords` | `int64` | `max_prior_passwords` | Number of prior passwords to disallow |
| `ManagedSiteSettings` | `interface{}` | `managed_site_settings` | Site settings managed by the parent site. Each key identifies a managed setting, with a null value. Non-secret effective values are returned in their own Site settings fields. |
| `MotdText` | `string` | `motd_text` | A message to show users when they connect via FTP or SFTP. |
| `MotdUseForFtp` | `*bool` | `motd_use_for_ftp` | Show message to users connecting via FTP |
| `MotdUseForSftp` | `*bool` | `motd_use_for_sftp` | Show message to users connecting via SFTP |
| `NextBillingAmount` | `string` | `next_billing_amount` | Next billing amount |
| `NextBillingDate` | `string` | `next_billing_date` | Next billing date |
| `OfficeIntegrationAvailable` | `*bool` | `office_integration_available` | If true, allows users to use a document editing integration. |
| `OfficeIntegrationType` | `string` | `office_integration_type` | Which document editing integration to support. Files.com Editor or Microsoft Office for the Web. |
| `OncehubLink` | `string` | `oncehub_link` | Link to scheduling a meeting with our Sales team |
| `OptOutGlobal` | `*bool` | `opt_out_global` | Use servers in the USA only? |
| `Overdue` | `*bool` | `overdue` | Is this site's billing overdue? |
| `PasswordMinLength` | `int64` | `password_min_length` | Shortest password length for users |
| `PasswordRequireLetter` | `*bool` | `password_require_letter` | Require a letter in passwords? |
| `PasswordRequireMixed` | `*bool` | `password_require_mixed` | Require lower and upper case letters in passwords? |
| `PasswordRequireNumber` | `*bool` | `password_require_number` | Require a number in passwords? |
| `PasswordRequireSpecial` | `*bool` | `password_require_special` | Require special characters in password? |
| `PasswordRequireUnbreached` | `*bool` | `password_require_unbreached` | Require passwords that have not been previously breached? (see https://haveibeenpwned.com/) |
| `PasswordRequirementsApplyToBundles` | `*bool` | `password_requirements_apply_to_bundles` | Require bundles' passwords, and passwords for other items (inboxes, public shares, etc.) to conform to the same requirements as users' passwords? |
| `PasswordValidityDays` | `int64` | `password_validity_days` | Number of days password is valid |
| `Phone` | `string` | `phone` | Site phone number |
| `PinAllRemoteServersToSiteRegion` | `*bool` | `pin_all_remote_servers_to_site_region` | If true, we will ensure that all internal communications with any remote server are made through the primary region of the site. This setting overrides individual remote server settings. |
| `PreventRootPermissionsForNonSiteAdmins` | `*bool` | `prevent_root_permissions_for_non_site_admins` | If true, we will prevent non-administrators from receiving any permissions directly on the root folder. This is commonly used to prevent the accidental application of permissions. |
| `ProtocolAccessGroupsOnly` | `*bool` | `protocol_access_groups_only` | If true, protocol access permissions on users will be ignored, and only protocol access permissions set on Groups will be honored. Make sure that your current user is a member of a group with API permission when changing this value to avoid locking yourself out of your site. |
| `Require2fa` | `*bool` | `require_2fa` | Require two-factor authentication for all users? |
| `RestrictRootFolderBehaviorsToSiteAdmins` | `*bool` | `restrict_root_folder_behaviors_to_site_admins` | If true, only site admins may create, modify, or delete any behavior at the site root, or a skip that would disable one. |
| `RootFolderBehaviorsApplyToWorkspaces` | `*bool` | `root_folder_behaviors_apply_to_workspaces` | If true, supported protective behaviors at the site root also apply within named workspaces. Requires restrict_root_folder_behaviors_to_site_admins to be enabled. |
| `Require2faExemptAllSsoUsers` | `*bool` | `require_2fa_exempt_all_sso_users` | If true, SSO users using the default user-level two-factor authentication setting are exempt from the site-wide two-factor authentication requirement. |
| `Require2faStopTime` | `*time.Time` | `require_2fa_stop_time` | If set, requirement for two-factor authentication has been scheduled to end on this date-time. |
| `RevokeBundleAccessOnDisableOrDelete` | `*bool` | `revoke_bundle_access_on_disable_or_delete` | Auto-removes bundles for disabled/deleted users and enforces bundle expiry within user access period. |
| `Require2faUserType` | `string` | `require_2fa_user_type` | What type of user is required to use two-factor authentication (when require_2fa is set to `true` for this site)? |
| `RequireLogoutFromBundlesAndInboxes` | `*bool` | `require_logout_from_bundles_and_inboxes` | If true, we will hide the 'Remember Me' box on Inbox and Bundle registration pages, requiring that the user logout and log back in every time they visit the page. |
| `Session` | `files_sdk.Session` | `session` | Current session |
| `S3CompatibleEndpointEnabled` | `*bool` | `s3_compatible_endpoint_enabled` | Is the S3-compatible endpoint enabled for all users and workspaces on this site? Defaults to true. When false, user and group S3 permissions do not allow access. |
| `SftpEnabled` | `*bool` | `sftp_enabled` | Is SFTP enabled? |
| `SftpFinalizePartialUploads` | `*bool` | `sftp_finalize_partial_uploads` | Finalize partial SFTP uploads from interrupted connections? Default: true. |
| `SftpHostKeyType` | `string` | `sftp_host_key_type` | Sftp Host Key Type |
| `ActiveSftpHostKeyId` | `int64` | `active_sftp_host_key_id` | Id of the currently selected custom SFTP Host Key |
| `ActiveSftpHostKeyIds` | `[]int64` | `active_sftp_host_key_ids` | Ids of the selected custom SFTP Host Keys |
| `SftpInsecureCiphers` | `*bool` | `sftp_insecure_ciphers` | If true, we will allow weak and known insecure ciphers to be used for SFTP connections. Enabling this setting severely weakens the security of your site and it is not recommend, except as a last resort for compatibility. |
| `SftpInsecureDiffieHellman` | `*bool` | `sftp_insecure_diffie_hellman` | If true, we will allow weak Diffie Hellman parameters to be used within ciphers for SFTP that are otherwise on our secure list. This has the effect of making the cipher weaker than our normal threshold for security, but is required to support certain legacy or broken SSH and MFT clients. Enabling this weakens security, but not nearly as much as enabling the full `sftp_insecure_ciphers` option. |
| `SftpUserRootEnabled` | `*bool` | `sftp_user_root_enabled` | Use user FTP roots also for SFTP? |
| `SharingEnabled` | `*bool` | `sharing_enabled` | Allow bundle creation |
| `ShowUserNotificationsLogInLink` | `*bool` | `show_user_notifications_log_in_link` | Show log in link in user notifications? |
| `ShowRequestAccessLink` | `*bool` | `show_request_access_link` | Show request access link for users without access? Currently unused. |
| `SiteFooter` | `string` | `site_footer` | Custom site footer text for authenticated pages |
| `SiteHeader` | `string` | `site_header` | Custom site header text for authenticated pages |
| `SitePublicFooter` | `string` | `site_public_footer` | Custom site footer text for public pages |
| `SitePublicHeader` | `string` | `site_public_header` | Custom site header text for public pages |
| `SmtpAddress` | `string` | `smtp_address` | SMTP server hostname or IP |
| `SmtpAuthentication` | `string` | `smtp_authentication` | SMTP server authentication type |
| `SmtpFrom` | `string` | `smtp_from` | From address to use when mailing through custom SMTP |
| `SmtpPort` | `int64` | `smtp_port` | SMTP server port |
| `SmtpSsl` | `string` | `smtp_ssl` | Custom SMTP encryption mode: if_available (default) uses STARTTLS when offered and otherwise sends credentials and messages unencrypted. A certificate-verified STARTTLS connection automatically changes if_available to require unless smtp_ssl is managed by a parent policy. require requires STARTTLS before authentication; require_implicit uses TLS from connection start; never disables TLS. TLS verifies the server certificate against smtp_address. Encryption is never automatically downgraded. |
| `SmtpUsername` | `string` | `smtp_username` | SMTP server username |
| `SessionExpiryMinutes` | `int64` | `session_expiry_minutes` | Session expiry in minutes |
| `SnapshotSharingEnabled` | `*bool` | `snapshot_sharing_enabled` | Allow snapshot share links creation |
| `SslRequired` | `*bool` | `ssl_required` | Is SSL required? Disabling this is insecure. |
| `Subdomain` | `string` | `subdomain` | Site subdomain |
| `SwitchToPlanDate` | `*time.Time` | `switch_to_plan_date` | If switching plans, when does the new plan take effect? |
| `TrialDaysLeft` | `int64` | `trial_days_left` | Number of days left in trial |
| `TrialUntil` | `*time.Time` | `trial_until` | When does this Site trial expire? |
| `UseDedicatedIpsForSmtp` | `*bool` | `use_dedicated_ips_for_smtp` | If using custom SMTP, should we use dedicated IPs to deliver emails? |
| `UseProvidedModifiedAt` | `*bool` | `use_provided_modified_at` | Allow uploaders to set `provided_modified_at` for uploaded files? |
| `User` | `files_sdk.User` | `user` | User of current session |
| `UserLockout` | `*bool` | `user_lockout` | Will users be locked out after incorrect login attempts? |
| `UserLockoutLockPeriod` | `int64` | `user_lockout_lock_period` | How many hours to lock user out for failed password? |
| `UserLockoutTries` | `int64` | `user_lockout_tries` | Number of login tries within `user_lockout_within` hours before users are locked out |
| `UserLockoutWithin` | `int64` | `user_lockout_within` | Number of hours for user lockout window |
| `UserRequestsEnabled` | `*bool` | `user_requests_enabled` | Enable User Requests feature |
| `UserRequestsNotifyAdmins` | `*bool` | `user_requests_notify_admins` | Send email to site admins when a user request is received? |
| `UsersCanCreateApiKeys` | `*bool` | `users_can_create_api_keys` | Allow users to create their own API keys? |
| `UsersCanCreateSshKeys` | `*bool` | `users_can_create_ssh_keys` | Allow users to create their own SSH keys? |
| `UsernameDisplay` | `string` | `username_display` | How usernames are displayed in the web UI. Can be `username_only`, `full_name_only`, `full_name_username`, `full_name_company`, or `full_name_username_company`. |
| `WelcomeCustomText` | `string` | `welcome_custom_text` | Custom text send in user welcome email |
| `EmailFooterCustomText` | `string` | `email_footer_custom_text` | Custom footer text for system-generated emails (as Markdown). Supports standard strftime date/time patterns like %Y (4-digit year), %m (month), %d (day). |
| `WelcomeEmailCc` | `string` | `welcome_email_cc` | Include this email in welcome emails if enabled |
| `WelcomeEmailSubject` | `string` | `welcome_email_subject` | Include this email subject in welcome emails if enabled |
| `WelcomeEmailEnabled` | `*bool` | `welcome_email_enabled` | Will the welcome email be sent to new users? |
| `WelcomeScreen` | `string` | `welcome_screen` | Does the welcome screen appear? |
| `WindowsModeFtp` | `*bool` | `windows_mode_ftp` | Does FTP user Windows emulation mode? |
