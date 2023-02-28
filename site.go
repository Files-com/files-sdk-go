package files_sdk

import (
	"encoding/json"
	"io"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Site struct {
	Name                                 string          `json:"name,omitempty" path:"name"`
	Allowed2faMethodSms                  *bool           `json:"allowed_2fa_method_sms,omitempty" path:"allowed_2fa_method_sms"`
	Allowed2faMethodTotp                 *bool           `json:"allowed_2fa_method_totp,omitempty" path:"allowed_2fa_method_totp"`
	Allowed2faMethodU2f                  *bool           `json:"allowed_2fa_method_u2f,omitempty" path:"allowed_2fa_method_u2f"`
	Allowed2faMethodWebauthn             *bool           `json:"allowed_2fa_method_webauthn,omitempty" path:"allowed_2fa_method_webauthn"`
	Allowed2faMethodYubi                 *bool           `json:"allowed_2fa_method_yubi,omitempty" path:"allowed_2fa_method_yubi"`
	Allowed2faMethodBypassForFtpSftpDav  *bool           `json:"allowed_2fa_method_bypass_for_ftp_sftp_dav,omitempty" path:"allowed_2fa_method_bypass_for_ftp_sftp_dav"`
	AdminUserId                          int64           `json:"admin_user_id,omitempty" path:"admin_user_id"`
	AllowBundleNames                     *bool           `json:"allow_bundle_names,omitempty" path:"allow_bundle_names"`
	AllowedCountries                     string          `json:"allowed_countries,omitempty" path:"allowed_countries"`
	AllowedIps                           string          `json:"allowed_ips,omitempty" path:"allowed_ips"`
	AskAboutOverwrites                   *bool           `json:"ask_about_overwrites,omitempty" path:"ask_about_overwrites"`
	BundleActivityNotifications          string          `json:"bundle_activity_notifications,omitempty" path:"bundle_activity_notifications"`
	BundleExpiration                     int64           `json:"bundle_expiration,omitempty" path:"bundle_expiration"`
	BundlePasswordRequired               *bool           `json:"bundle_password_required,omitempty" path:"bundle_password_required"`
	BundleRegistrationNotifications      string          `json:"bundle_registration_notifications,omitempty" path:"bundle_registration_notifications"`
	BundleRequireShareRecipient          *bool           `json:"bundle_require_share_recipient,omitempty" path:"bundle_require_share_recipient"`
	BundleUploadReceiptNotifications     string          `json:"bundle_upload_receipt_notifications,omitempty" path:"bundle_upload_receipt_notifications"`
	BundleWatermarkAttachment            Image           `json:"bundle_watermark_attachment,omitempty" path:"bundle_watermark_attachment"`
	BundleWatermarkValue                 json.RawMessage `json:"bundle_watermark_value,omitempty" path:"bundle_watermark_value"`
	UploadsViaEmailAuthentication        *bool           `json:"uploads_via_email_authentication,omitempty" path:"uploads_via_email_authentication"`
	Color2Left                           string          `json:"color2_left,omitempty" path:"color2_left"`
	Color2Link                           string          `json:"color2_link,omitempty" path:"color2_link"`
	Color2Text                           string          `json:"color2_text,omitempty" path:"color2_text"`
	Color2Top                            string          `json:"color2_top,omitempty" path:"color2_top"`
	Color2TopText                        string          `json:"color2_top_text,omitempty" path:"color2_top_text"`
	ContactName                          string          `json:"contact_name,omitempty" path:"contact_name"`
	CreatedAt                            *time.Time      `json:"created_at,omitempty" path:"created_at"`
	Currency                             string          `json:"currency,omitempty" path:"currency"`
	CustomNamespace                      *bool           `json:"custom_namespace,omitempty" path:"custom_namespace"`
	DaysToRetainBackups                  int64           `json:"days_to_retain_backups,omitempty" path:"days_to_retain_backups"`
	DefaultTimeZone                      string          `json:"default_time_zone,omitempty" path:"default_time_zone"`
	DesktopApp                           *bool           `json:"desktop_app,omitempty" path:"desktop_app"`
	DesktopAppSessionIpPinning           *bool           `json:"desktop_app_session_ip_pinning,omitempty" path:"desktop_app_session_ip_pinning"`
	DesktopAppSessionLifetime            int64           `json:"desktop_app_session_lifetime,omitempty" path:"desktop_app_session_lifetime"`
	MobileApp                            *bool           `json:"mobile_app,omitempty" path:"mobile_app"`
	MobileAppSessionIpPinning            *bool           `json:"mobile_app_session_ip_pinning,omitempty" path:"mobile_app_session_ip_pinning"`
	MobileAppSessionLifetime             int64           `json:"mobile_app_session_lifetime,omitempty" path:"mobile_app_session_lifetime"`
	DisallowedCountries                  string          `json:"disallowed_countries,omitempty" path:"disallowed_countries"`
	DisableFilesCertificateGeneration    *bool           `json:"disable_files_certificate_generation,omitempty" path:"disable_files_certificate_generation"`
	DisableNotifications                 *bool           `json:"disable_notifications,omitempty" path:"disable_notifications"`
	DisablePasswordReset                 *bool           `json:"disable_password_reset,omitempty" path:"disable_password_reset"`
	Domain                               string          `json:"domain,omitempty" path:"domain"`
	DomainHstsHeader                     *bool           `json:"domain_hsts_header,omitempty" path:"domain_hsts_header"`
	DomainLetsencryptChain               string          `json:"domain_letsencrypt_chain,omitempty" path:"domain_letsencrypt_chain"`
	Email                                string          `json:"email,omitempty" path:"email"`
	FtpEnabled                           *bool           `json:"ftp_enabled,omitempty" path:"ftp_enabled"`
	ReplyToEmail                         string          `json:"reply_to_email,omitempty" path:"reply_to_email"`
	NonSsoGroupsAllowed                  *bool           `json:"non_sso_groups_allowed,omitempty" path:"non_sso_groups_allowed"`
	NonSsoUsersAllowed                   *bool           `json:"non_sso_users_allowed,omitempty" path:"non_sso_users_allowed"`
	FolderPermissionsGroupsOnly          *bool           `json:"folder_permissions_groups_only,omitempty" path:"folder_permissions_groups_only"`
	Hipaa                                *bool           `json:"hipaa,omitempty" path:"hipaa"`
	Icon128                              Image           `json:"icon128,omitempty" path:"icon128"`
	Icon16                               Image           `json:"icon16,omitempty" path:"icon16"`
	Icon32                               Image           `json:"icon32,omitempty" path:"icon32"`
	Icon48                               Image           `json:"icon48,omitempty" path:"icon48"`
	ImmutableFilesSetAt                  *time.Time      `json:"immutable_files_set_at,omitempty" path:"immutable_files_set_at"`
	IncludePasswordInWelcomeEmail        *bool           `json:"include_password_in_welcome_email,omitempty" path:"include_password_in_welcome_email"`
	Language                             string          `json:"language,omitempty" path:"language"`
	LdapBaseDn                           string          `json:"ldap_base_dn,omitempty" path:"ldap_base_dn"`
	LdapDomain                           string          `json:"ldap_domain,omitempty" path:"ldap_domain"`
	LdapEnabled                          *bool           `json:"ldap_enabled,omitempty" path:"ldap_enabled"`
	LdapGroupAction                      string          `json:"ldap_group_action,omitempty" path:"ldap_group_action"`
	LdapGroupExclusion                   string          `json:"ldap_group_exclusion,omitempty" path:"ldap_group_exclusion"`
	LdapGroupInclusion                   string          `json:"ldap_group_inclusion,omitempty" path:"ldap_group_inclusion"`
	LdapHost                             string          `json:"ldap_host,omitempty" path:"ldap_host"`
	LdapHost2                            string          `json:"ldap_host_2,omitempty" path:"ldap_host_2"`
	LdapHost3                            string          `json:"ldap_host_3,omitempty" path:"ldap_host_3"`
	LdapPort                             int64           `json:"ldap_port,omitempty" path:"ldap_port"`
	LdapSecure                           *bool           `json:"ldap_secure,omitempty" path:"ldap_secure"`
	LdapType                             string          `json:"ldap_type,omitempty" path:"ldap_type"`
	LdapUserAction                       string          `json:"ldap_user_action,omitempty" path:"ldap_user_action"`
	LdapUserIncludeGroups                string          `json:"ldap_user_include_groups,omitempty" path:"ldap_user_include_groups"`
	LdapUsername                         string          `json:"ldap_username,omitempty" path:"ldap_username"`
	LdapUsernameField                    string          `json:"ldap_username_field,omitempty" path:"ldap_username_field"`
	LoginHelpText                        string          `json:"login_help_text,omitempty" path:"login_help_text"`
	Logo                                 Image           `json:"logo,omitempty" path:"logo"`
	MaxPriorPasswords                    int64           `json:"max_prior_passwords,omitempty" path:"max_prior_passwords"`
	MotdText                             string          `json:"motd_text,omitempty" path:"motd_text"`
	MotdUseForFtp                        *bool           `json:"motd_use_for_ftp,omitempty" path:"motd_use_for_ftp"`
	MotdUseForSftp                       *bool           `json:"motd_use_for_sftp,omitempty" path:"motd_use_for_sftp"`
	NextBillingAmount                    string          `json:"next_billing_amount,omitempty" path:"next_billing_amount"`
	NextBillingDate                      string          `json:"next_billing_date,omitempty" path:"next_billing_date"`
	OfficeIntegrationAvailable           *bool           `json:"office_integration_available,omitempty" path:"office_integration_available"`
	OfficeIntegrationType                string          `json:"office_integration_type,omitempty" path:"office_integration_type"`
	OncehubLink                          string          `json:"oncehub_link,omitempty" path:"oncehub_link"`
	OptOutGlobal                         *bool           `json:"opt_out_global,omitempty" path:"opt_out_global"`
	OverageNotifiedAt                    *time.Time      `json:"overage_notified_at,omitempty" path:"overage_notified_at"`
	OverageNotify                        *bool           `json:"overage_notify,omitempty" path:"overage_notify"`
	Overdue                              *bool           `json:"overdue,omitempty" path:"overdue"`
	PasswordMinLength                    int64           `json:"password_min_length,omitempty" path:"password_min_length"`
	PasswordRequireLetter                *bool           `json:"password_require_letter,omitempty" path:"password_require_letter"`
	PasswordRequireMixed                 *bool           `json:"password_require_mixed,omitempty" path:"password_require_mixed"`
	PasswordRequireNumber                *bool           `json:"password_require_number,omitempty" path:"password_require_number"`
	PasswordRequireSpecial               *bool           `json:"password_require_special,omitempty" path:"password_require_special"`
	PasswordRequireUnbreached            *bool           `json:"password_require_unbreached,omitempty" path:"password_require_unbreached"`
	PasswordRequirementsApplyToBundles   *bool           `json:"password_requirements_apply_to_bundles,omitempty" path:"password_requirements_apply_to_bundles"`
	PasswordValidityDays                 int64           `json:"password_validity_days,omitempty" path:"password_validity_days"`
	Phone                                string          `json:"phone,omitempty" path:"phone"`
	PinAllRemoteServersToSiteRegion      *bool           `json:"pin_all_remote_servers_to_site_region,omitempty" path:"pin_all_remote_servers_to_site_region"`
	Require2fa                           *bool           `json:"require_2fa,omitempty" path:"require_2fa"`
	Require2faStopTime                   *time.Time      `json:"require_2fa_stop_time,omitempty" path:"require_2fa_stop_time"`
	Require2faUserType                   string          `json:"require_2fa_user_type,omitempty" path:"require_2fa_user_type"`
	Session                              Session         `json:"session,omitempty" path:"session"`
	SessionPinnedByIp                    *bool           `json:"session_pinned_by_ip,omitempty" path:"session_pinned_by_ip"`
	SftpEnabled                          *bool           `json:"sftp_enabled,omitempty" path:"sftp_enabled"`
	SftpHostKeyType                      string          `json:"sftp_host_key_type,omitempty" path:"sftp_host_key_type"`
	ActiveSftpHostKeyId                  int64           `json:"active_sftp_host_key_id,omitempty" path:"active_sftp_host_key_id"`
	SftpInsecureCiphers                  *bool           `json:"sftp_insecure_ciphers,omitempty" path:"sftp_insecure_ciphers"`
	SftpUserRootEnabled                  *bool           `json:"sftp_user_root_enabled,omitempty" path:"sftp_user_root_enabled"`
	SharingEnabled                       *bool           `json:"sharing_enabled,omitempty" path:"sharing_enabled"`
	ShowRequestAccessLink                *bool           `json:"show_request_access_link,omitempty" path:"show_request_access_link"`
	SiteFooter                           string          `json:"site_footer,omitempty" path:"site_footer"`
	SiteHeader                           string          `json:"site_header,omitempty" path:"site_header"`
	SmtpAddress                          string          `json:"smtp_address,omitempty" path:"smtp_address"`
	SmtpAuthentication                   string          `json:"smtp_authentication,omitempty" path:"smtp_authentication"`
	SmtpFrom                             string          `json:"smtp_from,omitempty" path:"smtp_from"`
	SmtpPort                             int64           `json:"smtp_port,omitempty" path:"smtp_port"`
	SmtpUsername                         string          `json:"smtp_username,omitempty" path:"smtp_username"`
	SessionExpiry                        string          `json:"session_expiry,omitempty" path:"session_expiry"`
	SessionExpiryMinutes                 int64           `json:"session_expiry_minutes,omitempty" path:"session_expiry_minutes"`
	SslRequired                          *bool           `json:"ssl_required,omitempty" path:"ssl_required"`
	Subdomain                            string          `json:"subdomain,omitempty" path:"subdomain"`
	SwitchToPlanDate                     *time.Time      `json:"switch_to_plan_date,omitempty" path:"switch_to_plan_date"`
	TlsDisabled                          *bool           `json:"tls_disabled,omitempty" path:"tls_disabled"`
	TrialDaysLeft                        int64           `json:"trial_days_left,omitempty" path:"trial_days_left"`
	TrialUntil                           *time.Time      `json:"trial_until,omitempty" path:"trial_until"`
	UpdatedAt                            *time.Time      `json:"updated_at,omitempty" path:"updated_at"`
	UseProvidedModifiedAt                *bool           `json:"use_provided_modified_at,omitempty" path:"use_provided_modified_at"`
	User                                 User            `json:"user,omitempty" path:"user"`
	UserLockout                          *bool           `json:"user_lockout,omitempty" path:"user_lockout"`
	UserLockoutLockPeriod                int64           `json:"user_lockout_lock_period,omitempty" path:"user_lockout_lock_period"`
	UserLockoutTries                     int64           `json:"user_lockout_tries,omitempty" path:"user_lockout_tries"`
	UserLockoutWithin                    int64           `json:"user_lockout_within,omitempty" path:"user_lockout_within"`
	UserRequestsEnabled                  *bool           `json:"user_requests_enabled,omitempty" path:"user_requests_enabled"`
	UserRequestsNotifyAdmins             *bool           `json:"user_requests_notify_admins,omitempty" path:"user_requests_notify_admins"`
	WelcomeCustomText                    string          `json:"welcome_custom_text,omitempty" path:"welcome_custom_text"`
	WelcomeEmailCc                       string          `json:"welcome_email_cc,omitempty" path:"welcome_email_cc"`
	WelcomeEmailSubject                  string          `json:"welcome_email_subject,omitempty" path:"welcome_email_subject"`
	WelcomeEmailEnabled                  *bool           `json:"welcome_email_enabled,omitempty" path:"welcome_email_enabled"`
	WelcomeScreen                        string          `json:"welcome_screen,omitempty" path:"welcome_screen"`
	WindowsModeFtp                       *bool           `json:"windows_mode_ftp,omitempty" path:"windows_mode_ftp"`
	DisableUsersFromInactivityPeriodDays int64           `json:"disable_users_from_inactivity_period_days,omitempty" path:"disable_users_from_inactivity_period_days"`
}

type SiteCollection []Site

type SiteUpdateParams struct {
	Name                                 string          `url:"name,omitempty" required:"false" json:"name,omitempty" path:"name"`
	Subdomain                            string          `url:"subdomain,omitempty" required:"false" json:"subdomain,omitempty" path:"subdomain"`
	Domain                               string          `url:"domain,omitempty" required:"false" json:"domain,omitempty" path:"domain"`
	DomainHstsHeader                     *bool           `url:"domain_hsts_header,omitempty" required:"false" json:"domain_hsts_header,omitempty" path:"domain_hsts_header"`
	DomainLetsencryptChain               string          `url:"domain_letsencrypt_chain,omitempty" required:"false" json:"domain_letsencrypt_chain,omitempty" path:"domain_letsencrypt_chain"`
	Email                                string          `url:"email,omitempty" required:"false" json:"email,omitempty" path:"email"`
	ReplyToEmail                         string          `url:"reply_to_email,omitempty" required:"false" json:"reply_to_email,omitempty" path:"reply_to_email"`
	AllowBundleNames                     *bool           `url:"allow_bundle_names,omitempty" required:"false" json:"allow_bundle_names,omitempty" path:"allow_bundle_names"`
	BundleExpiration                     int64           `url:"bundle_expiration,omitempty" required:"false" json:"bundle_expiration,omitempty" path:"bundle_expiration"`
	OverageNotify                        *bool           `url:"overage_notify,omitempty" required:"false" json:"overage_notify,omitempty" path:"overage_notify"`
	WelcomeEmailEnabled                  *bool           `url:"welcome_email_enabled,omitempty" required:"false" json:"welcome_email_enabled,omitempty" path:"welcome_email_enabled"`
	AskAboutOverwrites                   *bool           `url:"ask_about_overwrites,omitempty" required:"false" json:"ask_about_overwrites,omitempty" path:"ask_about_overwrites"`
	ShowRequestAccessLink                *bool           `url:"show_request_access_link,omitempty" required:"false" json:"show_request_access_link,omitempty" path:"show_request_access_link"`
	WelcomeEmailCc                       string          `url:"welcome_email_cc,omitempty" required:"false" json:"welcome_email_cc,omitempty" path:"welcome_email_cc"`
	WelcomeEmailSubject                  string          `url:"welcome_email_subject,omitempty" required:"false" json:"welcome_email_subject,omitempty" path:"welcome_email_subject"`
	WelcomeCustomText                    string          `url:"welcome_custom_text,omitempty" required:"false" json:"welcome_custom_text,omitempty" path:"welcome_custom_text"`
	Language                             string          `url:"language,omitempty" required:"false" json:"language,omitempty" path:"language"`
	WindowsModeFtp                       *bool           `url:"windows_mode_ftp,omitempty" required:"false" json:"windows_mode_ftp,omitempty" path:"windows_mode_ftp"`
	DefaultTimeZone                      string          `url:"default_time_zone,omitempty" required:"false" json:"default_time_zone,omitempty" path:"default_time_zone"`
	DesktopApp                           *bool           `url:"desktop_app,omitempty" required:"false" json:"desktop_app,omitempty" path:"desktop_app"`
	DesktopAppSessionIpPinning           *bool           `url:"desktop_app_session_ip_pinning,omitempty" required:"false" json:"desktop_app_session_ip_pinning,omitempty" path:"desktop_app_session_ip_pinning"`
	DesktopAppSessionLifetime            int64           `url:"desktop_app_session_lifetime,omitempty" required:"false" json:"desktop_app_session_lifetime,omitempty" path:"desktop_app_session_lifetime"`
	MobileApp                            *bool           `url:"mobile_app,omitempty" required:"false" json:"mobile_app,omitempty" path:"mobile_app"`
	MobileAppSessionIpPinning            *bool           `url:"mobile_app_session_ip_pinning,omitempty" required:"false" json:"mobile_app_session_ip_pinning,omitempty" path:"mobile_app_session_ip_pinning"`
	MobileAppSessionLifetime             int64           `url:"mobile_app_session_lifetime,omitempty" required:"false" json:"mobile_app_session_lifetime,omitempty" path:"mobile_app_session_lifetime"`
	FolderPermissionsGroupsOnly          *bool           `url:"folder_permissions_groups_only,omitempty" required:"false" json:"folder_permissions_groups_only,omitempty" path:"folder_permissions_groups_only"`
	WelcomeScreen                        string          `url:"welcome_screen,omitempty" required:"false" json:"welcome_screen,omitempty" path:"welcome_screen"`
	OfficeIntegrationAvailable           *bool           `url:"office_integration_available,omitempty" required:"false" json:"office_integration_available,omitempty" path:"office_integration_available"`
	OfficeIntegrationType                string          `url:"office_integration_type,omitempty" required:"false" json:"office_integration_type,omitempty" path:"office_integration_type"`
	PinAllRemoteServersToSiteRegion      *bool           `url:"pin_all_remote_servers_to_site_region,omitempty" required:"false" json:"pin_all_remote_servers_to_site_region,omitempty" path:"pin_all_remote_servers_to_site_region"`
	MotdText                             string          `url:"motd_text,omitempty" required:"false" json:"motd_text,omitempty" path:"motd_text"`
	MotdUseForFtp                        *bool           `url:"motd_use_for_ftp,omitempty" required:"false" json:"motd_use_for_ftp,omitempty" path:"motd_use_for_ftp"`
	MotdUseForSftp                       *bool           `url:"motd_use_for_sftp,omitempty" required:"false" json:"motd_use_for_sftp,omitempty" path:"motd_use_for_sftp"`
	SessionExpiry                        string          `url:"session_expiry,omitempty" required:"false" json:"session_expiry,omitempty" path:"session_expiry"`
	SslRequired                          *bool           `url:"ssl_required,omitempty" required:"false" json:"ssl_required,omitempty" path:"ssl_required"`
	TlsDisabled                          *bool           `url:"tls_disabled,omitempty" required:"false" json:"tls_disabled,omitempty" path:"tls_disabled"`
	SftpInsecureCiphers                  *bool           `url:"sftp_insecure_ciphers,omitempty" required:"false" json:"sftp_insecure_ciphers,omitempty" path:"sftp_insecure_ciphers"`
	DisableFilesCertificateGeneration    *bool           `url:"disable_files_certificate_generation,omitempty" required:"false" json:"disable_files_certificate_generation,omitempty" path:"disable_files_certificate_generation"`
	UserLockout                          *bool           `url:"user_lockout,omitempty" required:"false" json:"user_lockout,omitempty" path:"user_lockout"`
	UserLockoutTries                     int64           `url:"user_lockout_tries,omitempty" required:"false" json:"user_lockout_tries,omitempty" path:"user_lockout_tries"`
	UserLockoutWithin                    int64           `url:"user_lockout_within,omitempty" required:"false" json:"user_lockout_within,omitempty" path:"user_lockout_within"`
	UserLockoutLockPeriod                int64           `url:"user_lockout_lock_period,omitempty" required:"false" json:"user_lockout_lock_period,omitempty" path:"user_lockout_lock_period"`
	IncludePasswordInWelcomeEmail        *bool           `url:"include_password_in_welcome_email,omitempty" required:"false" json:"include_password_in_welcome_email,omitempty" path:"include_password_in_welcome_email"`
	AllowedCountries                     string          `url:"allowed_countries,omitempty" required:"false" json:"allowed_countries,omitempty" path:"allowed_countries"`
	AllowedIps                           string          `url:"allowed_ips,omitempty" required:"false" json:"allowed_ips,omitempty" path:"allowed_ips"`
	DisallowedCountries                  string          `url:"disallowed_countries,omitempty" required:"false" json:"disallowed_countries,omitempty" path:"disallowed_countries"`
	DaysToRetainBackups                  int64           `url:"days_to_retain_backups,omitempty" required:"false" json:"days_to_retain_backups,omitempty" path:"days_to_retain_backups"`
	MaxPriorPasswords                    int64           `url:"max_prior_passwords,omitempty" required:"false" json:"max_prior_passwords,omitempty" path:"max_prior_passwords"`
	PasswordValidityDays                 int64           `url:"password_validity_days,omitempty" required:"false" json:"password_validity_days,omitempty" path:"password_validity_days"`
	PasswordMinLength                    int64           `url:"password_min_length,omitempty" required:"false" json:"password_min_length,omitempty" path:"password_min_length"`
	PasswordRequireLetter                *bool           `url:"password_require_letter,omitempty" required:"false" json:"password_require_letter,omitempty" path:"password_require_letter"`
	PasswordRequireMixed                 *bool           `url:"password_require_mixed,omitempty" required:"false" json:"password_require_mixed,omitempty" path:"password_require_mixed"`
	PasswordRequireSpecial               *bool           `url:"password_require_special,omitempty" required:"false" json:"password_require_special,omitempty" path:"password_require_special"`
	PasswordRequireNumber                *bool           `url:"password_require_number,omitempty" required:"false" json:"password_require_number,omitempty" path:"password_require_number"`
	PasswordRequireUnbreached            *bool           `url:"password_require_unbreached,omitempty" required:"false" json:"password_require_unbreached,omitempty" path:"password_require_unbreached"`
	SftpUserRootEnabled                  *bool           `url:"sftp_user_root_enabled,omitempty" required:"false" json:"sftp_user_root_enabled,omitempty" path:"sftp_user_root_enabled"`
	DisablePasswordReset                 *bool           `url:"disable_password_reset,omitempty" required:"false" json:"disable_password_reset,omitempty" path:"disable_password_reset"`
	ImmutableFiles                       *bool           `url:"immutable_files,omitempty" required:"false" json:"immutable_files,omitempty" path:"immutable_files"`
	SessionPinnedByIp                    *bool           `url:"session_pinned_by_ip,omitempty" required:"false" json:"session_pinned_by_ip,omitempty" path:"session_pinned_by_ip"`
	BundlePasswordRequired               *bool           `url:"bundle_password_required,omitempty" required:"false" json:"bundle_password_required,omitempty" path:"bundle_password_required"`
	BundleRequireShareRecipient          *bool           `url:"bundle_require_share_recipient,omitempty" required:"false" json:"bundle_require_share_recipient,omitempty" path:"bundle_require_share_recipient"`
	BundleRegistrationNotifications      string          `url:"bundle_registration_notifications,omitempty" required:"false" json:"bundle_registration_notifications,omitempty" path:"bundle_registration_notifications"`
	BundleActivityNotifications          string          `url:"bundle_activity_notifications,omitempty" required:"false" json:"bundle_activity_notifications,omitempty" path:"bundle_activity_notifications"`
	BundleUploadReceiptNotifications     string          `url:"bundle_upload_receipt_notifications,omitempty" required:"false" json:"bundle_upload_receipt_notifications,omitempty" path:"bundle_upload_receipt_notifications"`
	PasswordRequirementsApplyToBundles   *bool           `url:"password_requirements_apply_to_bundles,omitempty" required:"false" json:"password_requirements_apply_to_bundles,omitempty" path:"password_requirements_apply_to_bundles"`
	OptOutGlobal                         *bool           `url:"opt_out_global,omitempty" required:"false" json:"opt_out_global,omitempty" path:"opt_out_global"`
	UseProvidedModifiedAt                *bool           `url:"use_provided_modified_at,omitempty" required:"false" json:"use_provided_modified_at,omitempty" path:"use_provided_modified_at"`
	CustomNamespace                      *bool           `url:"custom_namespace,omitempty" required:"false" json:"custom_namespace,omitempty" path:"custom_namespace"`
	DisableUsersFromInactivityPeriodDays int64           `url:"disable_users_from_inactivity_period_days,omitempty" required:"false" json:"disable_users_from_inactivity_period_days,omitempty" path:"disable_users_from_inactivity_period_days"`
	NonSsoGroupsAllowed                  *bool           `url:"non_sso_groups_allowed,omitempty" required:"false" json:"non_sso_groups_allowed,omitempty" path:"non_sso_groups_allowed"`
	NonSsoUsersAllowed                   *bool           `url:"non_sso_users_allowed,omitempty" required:"false" json:"non_sso_users_allowed,omitempty" path:"non_sso_users_allowed"`
	SharingEnabled                       *bool           `url:"sharing_enabled,omitempty" required:"false" json:"sharing_enabled,omitempty" path:"sharing_enabled"`
	UserRequestsEnabled                  *bool           `url:"user_requests_enabled,omitempty" required:"false" json:"user_requests_enabled,omitempty" path:"user_requests_enabled"`
	UserRequestsNotifyAdmins             *bool           `url:"user_requests_notify_admins,omitempty" required:"false" json:"user_requests_notify_admins,omitempty" path:"user_requests_notify_admins"`
	FtpEnabled                           *bool           `url:"ftp_enabled,omitempty" required:"false" json:"ftp_enabled,omitempty" path:"ftp_enabled"`
	SftpEnabled                          *bool           `url:"sftp_enabled,omitempty" required:"false" json:"sftp_enabled,omitempty" path:"sftp_enabled"`
	SftpHostKeyType                      string          `url:"sftp_host_key_type,omitempty" required:"false" json:"sftp_host_key_type,omitempty" path:"sftp_host_key_type"`
	ActiveSftpHostKeyId                  int64           `url:"active_sftp_host_key_id,omitempty" required:"false" json:"active_sftp_host_key_id,omitempty" path:"active_sftp_host_key_id"`
	BundleWatermarkValue                 json.RawMessage `url:"bundle_watermark_value,omitempty" required:"false" json:"bundle_watermark_value,omitempty" path:"bundle_watermark_value"`
	Allowed2faMethodSms                  *bool           `url:"allowed_2fa_method_sms,omitempty" required:"false" json:"allowed_2fa_method_sms,omitempty" path:"allowed_2fa_method_sms"`
	Allowed2faMethodU2f                  *bool           `url:"allowed_2fa_method_u2f,omitempty" required:"false" json:"allowed_2fa_method_u2f,omitempty" path:"allowed_2fa_method_u2f"`
	Allowed2faMethodTotp                 *bool           `url:"allowed_2fa_method_totp,omitempty" required:"false" json:"allowed_2fa_method_totp,omitempty" path:"allowed_2fa_method_totp"`
	Allowed2faMethodWebauthn             *bool           `url:"allowed_2fa_method_webauthn,omitempty" required:"false" json:"allowed_2fa_method_webauthn,omitempty" path:"allowed_2fa_method_webauthn"`
	Allowed2faMethodYubi                 *bool           `url:"allowed_2fa_method_yubi,omitempty" required:"false" json:"allowed_2fa_method_yubi,omitempty" path:"allowed_2fa_method_yubi"`
	Allowed2faMethodBypassForFtpSftpDav  *bool           `url:"allowed_2fa_method_bypass_for_ftp_sftp_dav,omitempty" required:"false" json:"allowed_2fa_method_bypass_for_ftp_sftp_dav,omitempty" path:"allowed_2fa_method_bypass_for_ftp_sftp_dav"`
	Require2fa                           *bool           `url:"require_2fa,omitempty" required:"false" json:"require_2fa,omitempty" path:"require_2fa"`
	Require2faUserType                   string          `url:"require_2fa_user_type,omitempty" required:"false" json:"require_2fa_user_type,omitempty" path:"require_2fa_user_type"`
	Color2Top                            string          `url:"color2_top,omitempty" required:"false" json:"color2_top,omitempty" path:"color2_top"`
	Color2Left                           string          `url:"color2_left,omitempty" required:"false" json:"color2_left,omitempty" path:"color2_left"`
	Color2Link                           string          `url:"color2_link,omitempty" required:"false" json:"color2_link,omitempty" path:"color2_link"`
	Color2Text                           string          `url:"color2_text,omitempty" required:"false" json:"color2_text,omitempty" path:"color2_text"`
	Color2TopText                        string          `url:"color2_top_text,omitempty" required:"false" json:"color2_top_text,omitempty" path:"color2_top_text"`
	SiteHeader                           string          `url:"site_header,omitempty" required:"false" json:"site_header,omitempty" path:"site_header"`
	SiteFooter                           string          `url:"site_footer,omitempty" required:"false" json:"site_footer,omitempty" path:"site_footer"`
	LoginHelpText                        string          `url:"login_help_text,omitempty" required:"false" json:"login_help_text,omitempty" path:"login_help_text"`
	SmtpAddress                          string          `url:"smtp_address,omitempty" required:"false" json:"smtp_address,omitempty" path:"smtp_address"`
	SmtpAuthentication                   string          `url:"smtp_authentication,omitempty" required:"false" json:"smtp_authentication,omitempty" path:"smtp_authentication"`
	SmtpFrom                             string          `url:"smtp_from,omitempty" required:"false" json:"smtp_from,omitempty" path:"smtp_from"`
	SmtpUsername                         string          `url:"smtp_username,omitempty" required:"false" json:"smtp_username,omitempty" path:"smtp_username"`
	SmtpPort                             int64           `url:"smtp_port,omitempty" required:"false" json:"smtp_port,omitempty" path:"smtp_port"`
	LdapEnabled                          *bool           `url:"ldap_enabled,omitempty" required:"false" json:"ldap_enabled,omitempty" path:"ldap_enabled"`
	LdapType                             string          `url:"ldap_type,omitempty" required:"false" json:"ldap_type,omitempty" path:"ldap_type"`
	LdapHost                             string          `url:"ldap_host,omitempty" required:"false" json:"ldap_host,omitempty" path:"ldap_host"`
	LdapHost2                            string          `url:"ldap_host_2,omitempty" required:"false" json:"ldap_host_2,omitempty" path:"ldap_host_2"`
	LdapHost3                            string          `url:"ldap_host_3,omitempty" required:"false" json:"ldap_host_3,omitempty" path:"ldap_host_3"`
	LdapPort                             int64           `url:"ldap_port,omitempty" required:"false" json:"ldap_port,omitempty" path:"ldap_port"`
	LdapSecure                           *bool           `url:"ldap_secure,omitempty" required:"false" json:"ldap_secure,omitempty" path:"ldap_secure"`
	LdapUsername                         string          `url:"ldap_username,omitempty" required:"false" json:"ldap_username,omitempty" path:"ldap_username"`
	LdapUsernameField                    string          `url:"ldap_username_field,omitempty" required:"false" json:"ldap_username_field,omitempty" path:"ldap_username_field"`
	LdapDomain                           string          `url:"ldap_domain,omitempty" required:"false" json:"ldap_domain,omitempty" path:"ldap_domain"`
	LdapUserAction                       string          `url:"ldap_user_action,omitempty" required:"false" json:"ldap_user_action,omitempty" path:"ldap_user_action"`
	LdapGroupAction                      string          `url:"ldap_group_action,omitempty" required:"false" json:"ldap_group_action,omitempty" path:"ldap_group_action"`
	LdapUserIncludeGroups                string          `url:"ldap_user_include_groups,omitempty" required:"false" json:"ldap_user_include_groups,omitempty" path:"ldap_user_include_groups"`
	LdapGroupExclusion                   string          `url:"ldap_group_exclusion,omitempty" required:"false" json:"ldap_group_exclusion,omitempty" path:"ldap_group_exclusion"`
	LdapGroupInclusion                   string          `url:"ldap_group_inclusion,omitempty" required:"false" json:"ldap_group_inclusion,omitempty" path:"ldap_group_inclusion"`
	LdapBaseDn                           string          `url:"ldap_base_dn,omitempty" required:"false" json:"ldap_base_dn,omitempty" path:"ldap_base_dn"`
	UploadsViaEmailAuthentication        *bool           `url:"uploads_via_email_authentication,omitempty" required:"false" json:"uploads_via_email_authentication,omitempty" path:"uploads_via_email_authentication"`
	Icon16File                           io.Writer       `url:"icon16_file,omitempty" required:"false" json:"icon16_file,omitempty" path:"icon16_file"`
	Icon16Delete                         *bool           `url:"icon16_delete,omitempty" required:"false" json:"icon16_delete,omitempty" path:"icon16_delete"`
	Icon32File                           io.Writer       `url:"icon32_file,omitempty" required:"false" json:"icon32_file,omitempty" path:"icon32_file"`
	Icon32Delete                         *bool           `url:"icon32_delete,omitempty" required:"false" json:"icon32_delete,omitempty" path:"icon32_delete"`
	Icon48File                           io.Writer       `url:"icon48_file,omitempty" required:"false" json:"icon48_file,omitempty" path:"icon48_file"`
	Icon48Delete                         *bool           `url:"icon48_delete,omitempty" required:"false" json:"icon48_delete,omitempty" path:"icon48_delete"`
	Icon128File                          io.Writer       `url:"icon128_file,omitempty" required:"false" json:"icon128_file,omitempty" path:"icon128_file"`
	Icon128Delete                        *bool           `url:"icon128_delete,omitempty" required:"false" json:"icon128_delete,omitempty" path:"icon128_delete"`
	LogoFile                             io.Writer       `url:"logo_file,omitempty" required:"false" json:"logo_file,omitempty" path:"logo_file"`
	LogoDelete                           *bool           `url:"logo_delete,omitempty" required:"false" json:"logo_delete,omitempty" path:"logo_delete"`
	BundleWatermarkAttachmentFile        io.Writer       `url:"bundle_watermark_attachment_file,omitempty" required:"false" json:"bundle_watermark_attachment_file,omitempty" path:"bundle_watermark_attachment_file"`
	BundleWatermarkAttachmentDelete      *bool           `url:"bundle_watermark_attachment_delete,omitempty" required:"false" json:"bundle_watermark_attachment_delete,omitempty" path:"bundle_watermark_attachment_delete"`
	Disable2faWithDelay                  *bool           `url:"disable_2fa_with_delay,omitempty" required:"false" json:"disable_2fa_with_delay,omitempty" path:"disable_2fa_with_delay"`
	LdapPasswordChange                   string          `url:"ldap_password_change,omitempty" required:"false" json:"ldap_password_change,omitempty" path:"ldap_password_change"`
	LdapPasswordChangeConfirmation       string          `url:"ldap_password_change_confirmation,omitempty" required:"false" json:"ldap_password_change_confirmation,omitempty" path:"ldap_password_change_confirmation"`
	SmtpPassword                         string          `url:"smtp_password,omitempty" required:"false" json:"smtp_password,omitempty" path:"smtp_password"`
	SessionExpiryMinutes                 int64           `url:"session_expiry_minutes,omitempty" required:"false" json:"session_expiry_minutes,omitempty" path:"session_expiry_minutes"`
}

func (s *Site) UnmarshalJSON(data []byte) error {
	type site Site
	var v site
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*s = Site(v)
	return nil
}

func (s *SiteCollection) UnmarshalJSON(data []byte) error {
	type sites SiteCollection
	var v sites
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*s = SiteCollection(v)
	return nil
}

func (s *SiteCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*s))
	for i, v := range *s {
		ret[i] = v
	}

	return &ret
}
