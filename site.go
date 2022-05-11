package files_sdk

import (
	"encoding/json"
	"io"
	"time"
)

type Site struct {
	Name                                 string          `json:"name,omitempty"`
	Allowed2faMethodSms                  *bool           `json:"allowed_2fa_method_sms,omitempty"`
	Allowed2faMethodTotp                 *bool           `json:"allowed_2fa_method_totp,omitempty"`
	Allowed2faMethodU2f                  *bool           `json:"allowed_2fa_method_u2f,omitempty"`
	Allowed2faMethodWebauthn             *bool           `json:"allowed_2fa_method_webauthn,omitempty"`
	Allowed2faMethodYubi                 *bool           `json:"allowed_2fa_method_yubi,omitempty"`
	Allowed2faMethodBypassForFtpSftpDav  *bool           `json:"allowed_2fa_method_bypass_for_ftp_sftp_dav,omitempty"`
	AdminUserId                          int64           `json:"admin_user_id,omitempty"`
	AllowBundleNames                     *bool           `json:"allow_bundle_names,omitempty"`
	AllowedCountries                     string          `json:"allowed_countries,omitempty"`
	AllowedIps                           string          `json:"allowed_ips,omitempty"`
	AskAboutOverwrites                   *bool           `json:"ask_about_overwrites,omitempty"`
	BundleExpiration                     int64           `json:"bundle_expiration,omitempty"`
	BundlePasswordRequired               *bool           `json:"bundle_password_required,omitempty"`
	BundleRequireShareRecipient          *bool           `json:"bundle_require_share_recipient,omitempty"`
	BundleWatermarkAttachment            Image           `json:"bundle_watermark_attachment,omitempty"`
	BundleWatermarkValue                 json.RawMessage `json:"bundle_watermark_value,omitempty"`
	Color2Left                           string          `json:"color2_left,omitempty"`
	Color2Link                           string          `json:"color2_link,omitempty"`
	Color2Text                           string          `json:"color2_text,omitempty"`
	Color2Top                            string          `json:"color2_top,omitempty"`
	Color2TopText                        string          `json:"color2_top_text,omitempty"`
	ContactName                          string          `json:"contact_name,omitempty"`
	CreatedAt                            time.Time       `json:"created_at,omitempty"`
	Currency                             string          `json:"currency,omitempty"`
	CustomNamespace                      *bool           `json:"custom_namespace,omitempty"`
	DaysToRetainBackups                  int64           `json:"days_to_retain_backups,omitempty"`
	DefaultTimeZone                      string          `json:"default_time_zone,omitempty"`
	DesktopApp                           *bool           `json:"desktop_app,omitempty"`
	DesktopAppSessionIpPinning           *bool           `json:"desktop_app_session_ip_pinning,omitempty"`
	DesktopAppSessionLifetime            int64           `json:"desktop_app_session_lifetime,omitempty"`
	MobileApp                            *bool           `json:"mobile_app,omitempty"`
	MobileAppSessionIpPinning            *bool           `json:"mobile_app_session_ip_pinning,omitempty"`
	MobileAppSessionLifetime             int64           `json:"mobile_app_session_lifetime,omitempty"`
	DisallowedCountries                  string          `json:"disallowed_countries,omitempty"`
	DisableNotifications                 *bool           `json:"disable_notifications,omitempty"`
	DisablePasswordReset                 *bool           `json:"disable_password_reset,omitempty"`
	Domain                               string          `json:"domain,omitempty"`
	DomainHstsHeader                     *bool           `json:"domain_hsts_header,omitempty"`
	DomainLetsencryptChain               string          `json:"domain_letsencrypt_chain,omitempty"`
	Email                                string          `json:"email,omitempty"`
	FtpEnabled                           *bool           `json:"ftp_enabled,omitempty"`
	ReplyToEmail                         string          `json:"reply_to_email,omitempty"`
	NonSsoGroupsAllowed                  *bool           `json:"non_sso_groups_allowed,omitempty"`
	NonSsoUsersAllowed                   *bool           `json:"non_sso_users_allowed,omitempty"`
	FolderPermissionsGroupsOnly          *bool           `json:"folder_permissions_groups_only,omitempty"`
	Hipaa                                *bool           `json:"hipaa,omitempty"`
	Icon128                              Image           `json:"icon128,omitempty"`
	Icon16                               Image           `json:"icon16,omitempty"`
	Icon32                               Image           `json:"icon32,omitempty"`
	Icon48                               Image           `json:"icon48,omitempty"`
	ImmutableFilesSetAt                  time.Time       `json:"immutable_files_set_at,omitempty"`
	IncludePasswordInWelcomeEmail        *bool           `json:"include_password_in_welcome_email,omitempty"`
	Language                             string          `json:"language,omitempty"`
	LdapBaseDn                           string          `json:"ldap_base_dn,omitempty"`
	LdapDomain                           string          `json:"ldap_domain,omitempty"`
	LdapEnabled                          *bool           `json:"ldap_enabled,omitempty"`
	LdapGroupAction                      string          `json:"ldap_group_action,omitempty"`
	LdapGroupExclusion                   string          `json:"ldap_group_exclusion,omitempty"`
	LdapGroupInclusion                   string          `json:"ldap_group_inclusion,omitempty"`
	LdapHost                             string          `json:"ldap_host,omitempty"`
	LdapHost2                            string          `json:"ldap_host_2,omitempty"`
	LdapHost3                            string          `json:"ldap_host_3,omitempty"`
	LdapPort                             int64           `json:"ldap_port,omitempty"`
	LdapSecure                           *bool           `json:"ldap_secure,omitempty"`
	LdapType                             string          `json:"ldap_type,omitempty"`
	LdapUserAction                       string          `json:"ldap_user_action,omitempty"`
	LdapUserIncludeGroups                string          `json:"ldap_user_include_groups,omitempty"`
	LdapUsername                         string          `json:"ldap_username,omitempty"`
	LdapUsernameField                    string          `json:"ldap_username_field,omitempty"`
	LoginHelpText                        string          `json:"login_help_text,omitempty"`
	Logo                                 Image           `json:"logo,omitempty"`
	MaxPriorPasswords                    int64           `json:"max_prior_passwords,omitempty"`
	NextBillingAmount                    float32         `json:"next_billing_amount,omitempty"`
	NextBillingDate                      string          `json:"next_billing_date,omitempty"`
	OfficeIntegrationAvailable           *bool           `json:"office_integration_available,omitempty"`
	OncehubLink                          string          `json:"oncehub_link,omitempty"`
	OptOutGlobal                         *bool           `json:"opt_out_global,omitempty"`
	OverageNotifiedAt                    time.Time       `json:"overage_notified_at,omitempty"`
	OverageNotify                        *bool           `json:"overage_notify,omitempty"`
	Overdue                              *bool           `json:"overdue,omitempty"`
	PasswordMinLength                    int64           `json:"password_min_length,omitempty"`
	PasswordRequireLetter                *bool           `json:"password_require_letter,omitempty"`
	PasswordRequireMixed                 *bool           `json:"password_require_mixed,omitempty"`
	PasswordRequireNumber                *bool           `json:"password_require_number,omitempty"`
	PasswordRequireSpecial               *bool           `json:"password_require_special,omitempty"`
	PasswordRequireUnbreached            *bool           `json:"password_require_unbreached,omitempty"`
	PasswordRequirementsApplyToBundles   *bool           `json:"password_requirements_apply_to_bundles,omitempty"`
	PasswordValidityDays                 int64           `json:"password_validity_days,omitempty"`
	Phone                                string          `json:"phone,omitempty"`
	Require2fa                           *bool           `json:"require_2fa,omitempty"`
	Require2faStopTime                   time.Time       `json:"require_2fa_stop_time,omitempty"`
	Require2faUserType                   string          `json:"require_2fa_user_type,omitempty"`
	Session                              Session         `json:"session,omitempty"`
	SessionPinnedByIp                    *bool           `json:"session_pinned_by_ip,omitempty"`
	SftpEnabled                          *bool           `json:"sftp_enabled,omitempty"`
	SftpInsecureCiphers                  *bool           `json:"sftp_insecure_ciphers,omitempty"`
	SftpUserRootEnabled                  *bool           `json:"sftp_user_root_enabled,omitempty"`
	SharingEnabled                       *bool           `json:"sharing_enabled,omitempty"`
	ShowRequestAccessLink                *bool           `json:"show_request_access_link,omitempty"`
	SiteFooter                           string          `json:"site_footer,omitempty"`
	SiteHeader                           string          `json:"site_header,omitempty"`
	SmtpAddress                          string          `json:"smtp_address,omitempty"`
	SmtpAuthentication                   string          `json:"smtp_authentication,omitempty"`
	SmtpFrom                             string          `json:"smtp_from,omitempty"`
	SmtpPort                             int64           `json:"smtp_port,omitempty"`
	SmtpUsername                         string          `json:"smtp_username,omitempty"`
	SessionExpiry                        float32         `json:"session_expiry,omitempty"`
	SslRequired                          *bool           `json:"ssl_required,omitempty"`
	Subdomain                            string          `json:"subdomain,omitempty"`
	SwitchToPlanDate                     time.Time       `json:"switch_to_plan_date,omitempty"`
	TlsDisabled                          *bool           `json:"tls_disabled,omitempty"`
	TrialDaysLeft                        int64           `json:"trial_days_left,omitempty"`
	TrialUntil                           time.Time       `json:"trial_until,omitempty"`
	UpdatedAt                            time.Time       `json:"updated_at,omitempty"`
	UseProvidedModifiedAt                *bool           `json:"use_provided_modified_at,omitempty"`
	User                                 User            `json:"user,omitempty"`
	UserLockout                          *bool           `json:"user_lockout,omitempty"`
	UserLockoutLockPeriod                int64           `json:"user_lockout_lock_period,omitempty"`
	UserLockoutTries                     int64           `json:"user_lockout_tries,omitempty"`
	UserLockoutWithin                    int64           `json:"user_lockout_within,omitempty"`
	UserRequestsEnabled                  *bool           `json:"user_requests_enabled,omitempty"`
	WelcomeCustomText                    string          `json:"welcome_custom_text,omitempty"`
	WelcomeEmailCc                       string          `json:"welcome_email_cc,omitempty"`
	WelcomeEmailSubject                  string          `json:"welcome_email_subject,omitempty"`
	WelcomeEmailEnabled                  *bool           `json:"welcome_email_enabled,omitempty"`
	WelcomeScreen                        string          `json:"welcome_screen,omitempty"`
	WindowsModeFtp                       *bool           `json:"windows_mode_ftp,omitempty"`
	DisableUsersFromInactivityPeriodDays int64           `json:"disable_users_from_inactivity_period_days,omitempty"`
}

type SiteCollection []Site

type SiteUpdateParams struct {
	Name                                 string    `url:"name,omitempty" required:"false" json:"name,omitempty"`
	Subdomain                            string    `url:"subdomain,omitempty" required:"false" json:"subdomain,omitempty"`
	Domain                               string    `url:"domain,omitempty" required:"false" json:"domain,omitempty"`
	DomainHstsHeader                     *bool     `url:"domain_hsts_header,omitempty" required:"false" json:"domain_hsts_header,omitempty"`
	DomainLetsencryptChain               string    `url:"domain_letsencrypt_chain,omitempty" required:"false" json:"domain_letsencrypt_chain,omitempty"`
	Email                                string    `url:"email,omitempty" required:"false" json:"email,omitempty"`
	ReplyToEmail                         string    `url:"reply_to_email,omitempty" required:"false" json:"reply_to_email,omitempty"`
	AllowBundleNames                     *bool     `url:"allow_bundle_names,omitempty" required:"false" json:"allow_bundle_names,omitempty"`
	BundleExpiration                     int64     `url:"bundle_expiration,omitempty" required:"false" json:"bundle_expiration,omitempty"`
	OverageNotify                        *bool     `url:"overage_notify,omitempty" required:"false" json:"overage_notify,omitempty"`
	WelcomeEmailEnabled                  *bool     `url:"welcome_email_enabled,omitempty" required:"false" json:"welcome_email_enabled,omitempty"`
	AskAboutOverwrites                   *bool     `url:"ask_about_overwrites,omitempty" required:"false" json:"ask_about_overwrites,omitempty"`
	ShowRequestAccessLink                *bool     `url:"show_request_access_link,omitempty" required:"false" json:"show_request_access_link,omitempty"`
	WelcomeEmailCc                       string    `url:"welcome_email_cc,omitempty" required:"false" json:"welcome_email_cc,omitempty"`
	WelcomeEmailSubject                  string    `url:"welcome_email_subject,omitempty" required:"false" json:"welcome_email_subject,omitempty"`
	WelcomeCustomText                    string    `url:"welcome_custom_text,omitempty" required:"false" json:"welcome_custom_text,omitempty"`
	Language                             string    `url:"language,omitempty" required:"false" json:"language,omitempty"`
	WindowsModeFtp                       *bool     `url:"windows_mode_ftp,omitempty" required:"false" json:"windows_mode_ftp,omitempty"`
	DefaultTimeZone                      string    `url:"default_time_zone,omitempty" required:"false" json:"default_time_zone,omitempty"`
	DesktopApp                           *bool     `url:"desktop_app,omitempty" required:"false" json:"desktop_app,omitempty"`
	DesktopAppSessionIpPinning           *bool     `url:"desktop_app_session_ip_pinning,omitempty" required:"false" json:"desktop_app_session_ip_pinning,omitempty"`
	DesktopAppSessionLifetime            int64     `url:"desktop_app_session_lifetime,omitempty" required:"false" json:"desktop_app_session_lifetime,omitempty"`
	MobileApp                            *bool     `url:"mobile_app,omitempty" required:"false" json:"mobile_app,omitempty"`
	MobileAppSessionIpPinning            *bool     `url:"mobile_app_session_ip_pinning,omitempty" required:"false" json:"mobile_app_session_ip_pinning,omitempty"`
	MobileAppSessionLifetime             int64     `url:"mobile_app_session_lifetime,omitempty" required:"false" json:"mobile_app_session_lifetime,omitempty"`
	FolderPermissionsGroupsOnly          *bool     `url:"folder_permissions_groups_only,omitempty" required:"false" json:"folder_permissions_groups_only,omitempty"`
	WelcomeScreen                        string    `url:"welcome_screen,omitempty" required:"false" json:"welcome_screen,omitempty"`
	OfficeIntegrationAvailable           *bool     `url:"office_integration_available,omitempty" required:"false" json:"office_integration_available,omitempty"`
	SessionExpiry                        float32   `url:"session_expiry,omitempty" required:"false" json:"session_expiry,omitempty"`
	SslRequired                          *bool     `url:"ssl_required,omitempty" required:"false" json:"ssl_required,omitempty"`
	TlsDisabled                          *bool     `url:"tls_disabled,omitempty" required:"false" json:"tls_disabled,omitempty"`
	SftpInsecureCiphers                  *bool     `url:"sftp_insecure_ciphers,omitempty" required:"false" json:"sftp_insecure_ciphers,omitempty"`
	UserLockout                          *bool     `url:"user_lockout,omitempty" required:"false" json:"user_lockout,omitempty"`
	UserLockoutTries                     int64     `url:"user_lockout_tries,omitempty" required:"false" json:"user_lockout_tries,omitempty"`
	UserLockoutWithin                    int64     `url:"user_lockout_within,omitempty" required:"false" json:"user_lockout_within,omitempty"`
	UserLockoutLockPeriod                int64     `url:"user_lockout_lock_period,omitempty" required:"false" json:"user_lockout_lock_period,omitempty"`
	IncludePasswordInWelcomeEmail        *bool     `url:"include_password_in_welcome_email,omitempty" required:"false" json:"include_password_in_welcome_email,omitempty"`
	AllowedCountries                     string    `url:"allowed_countries,omitempty" required:"false" json:"allowed_countries,omitempty"`
	AllowedIps                           string    `url:"allowed_ips,omitempty" required:"false" json:"allowed_ips,omitempty"`
	DisallowedCountries                  string    `url:"disallowed_countries,omitempty" required:"false" json:"disallowed_countries,omitempty"`
	DaysToRetainBackups                  int64     `url:"days_to_retain_backups,omitempty" required:"false" json:"days_to_retain_backups,omitempty"`
	MaxPriorPasswords                    int64     `url:"max_prior_passwords,omitempty" required:"false" json:"max_prior_passwords,omitempty"`
	PasswordValidityDays                 int64     `url:"password_validity_days,omitempty" required:"false" json:"password_validity_days,omitempty"`
	PasswordMinLength                    int64     `url:"password_min_length,omitempty" required:"false" json:"password_min_length,omitempty"`
	PasswordRequireLetter                *bool     `url:"password_require_letter,omitempty" required:"false" json:"password_require_letter,omitempty"`
	PasswordRequireMixed                 *bool     `url:"password_require_mixed,omitempty" required:"false" json:"password_require_mixed,omitempty"`
	PasswordRequireSpecial               *bool     `url:"password_require_special,omitempty" required:"false" json:"password_require_special,omitempty"`
	PasswordRequireNumber                *bool     `url:"password_require_number,omitempty" required:"false" json:"password_require_number,omitempty"`
	PasswordRequireUnbreached            *bool     `url:"password_require_unbreached,omitempty" required:"false" json:"password_require_unbreached,omitempty"`
	SftpUserRootEnabled                  *bool     `url:"sftp_user_root_enabled,omitempty" required:"false" json:"sftp_user_root_enabled,omitempty"`
	DisablePasswordReset                 *bool     `url:"disable_password_reset,omitempty" required:"false" json:"disable_password_reset,omitempty"`
	ImmutableFiles                       *bool     `url:"immutable_files,omitempty" required:"false" json:"immutable_files,omitempty"`
	SessionPinnedByIp                    *bool     `url:"session_pinned_by_ip,omitempty" required:"false" json:"session_pinned_by_ip,omitempty"`
	BundlePasswordRequired               *bool     `url:"bundle_password_required,omitempty" required:"false" json:"bundle_password_required,omitempty"`
	BundleRequireShareRecipient          *bool     `url:"bundle_require_share_recipient,omitempty" required:"false" json:"bundle_require_share_recipient,omitempty"`
	PasswordRequirementsApplyToBundles   *bool     `url:"password_requirements_apply_to_bundles,omitempty" required:"false" json:"password_requirements_apply_to_bundles,omitempty"`
	OptOutGlobal                         *bool     `url:"opt_out_global,omitempty" required:"false" json:"opt_out_global,omitempty"`
	UseProvidedModifiedAt                *bool     `url:"use_provided_modified_at,omitempty" required:"false" json:"use_provided_modified_at,omitempty"`
	CustomNamespace                      *bool     `url:"custom_namespace,omitempty" required:"false" json:"custom_namespace,omitempty"`
	DisableUsersFromInactivityPeriodDays int64     `url:"disable_users_from_inactivity_period_days,omitempty" required:"false" json:"disable_users_from_inactivity_period_days,omitempty"`
	NonSsoGroupsAllowed                  *bool     `url:"non_sso_groups_allowed,omitempty" required:"false" json:"non_sso_groups_allowed,omitempty"`
	NonSsoUsersAllowed                   *bool     `url:"non_sso_users_allowed,omitempty" required:"false" json:"non_sso_users_allowed,omitempty"`
	SharingEnabled                       *bool     `url:"sharing_enabled,omitempty" required:"false" json:"sharing_enabled,omitempty"`
	UserRequestsEnabled                  *bool     `url:"user_requests_enabled,omitempty" required:"false" json:"user_requests_enabled,omitempty"`
	FtpEnabled                           *bool     `url:"ftp_enabled,omitempty" required:"false" json:"ftp_enabled,omitempty"`
	SftpEnabled                          *bool     `url:"sftp_enabled,omitempty" required:"false" json:"sftp_enabled,omitempty"`
	Allowed2faMethodSms                  *bool     `url:"allowed_2fa_method_sms,omitempty" required:"false" json:"allowed_2fa_method_sms,omitempty"`
	Allowed2faMethodU2f                  *bool     `url:"allowed_2fa_method_u2f,omitempty" required:"false" json:"allowed_2fa_method_u2f,omitempty"`
	Allowed2faMethodTotp                 *bool     `url:"allowed_2fa_method_totp,omitempty" required:"false" json:"allowed_2fa_method_totp,omitempty"`
	Allowed2faMethodWebauthn             *bool     `url:"allowed_2fa_method_webauthn,omitempty" required:"false" json:"allowed_2fa_method_webauthn,omitempty"`
	Allowed2faMethodYubi                 *bool     `url:"allowed_2fa_method_yubi,omitempty" required:"false" json:"allowed_2fa_method_yubi,omitempty"`
	Allowed2faMethodBypassForFtpSftpDav  *bool     `url:"allowed_2fa_method_bypass_for_ftp_sftp_dav,omitempty" required:"false" json:"allowed_2fa_method_bypass_for_ftp_sftp_dav,omitempty"`
	Require2fa                           *bool     `url:"require_2fa,omitempty" required:"false" json:"require_2fa,omitempty"`
	Require2faUserType                   string    `url:"require_2fa_user_type,omitempty" required:"false" json:"require_2fa_user_type,omitempty"`
	Color2Top                            string    `url:"color2_top,omitempty" required:"false" json:"color2_top,omitempty"`
	Color2Left                           string    `url:"color2_left,omitempty" required:"false" json:"color2_left,omitempty"`
	Color2Link                           string    `url:"color2_link,omitempty" required:"false" json:"color2_link,omitempty"`
	Color2Text                           string    `url:"color2_text,omitempty" required:"false" json:"color2_text,omitempty"`
	Color2TopText                        string    `url:"color2_top_text,omitempty" required:"false" json:"color2_top_text,omitempty"`
	SiteHeader                           string    `url:"site_header,omitempty" required:"false" json:"site_header,omitempty"`
	SiteFooter                           string    `url:"site_footer,omitempty" required:"false" json:"site_footer,omitempty"`
	LoginHelpText                        string    `url:"login_help_text,omitempty" required:"false" json:"login_help_text,omitempty"`
	SmtpAddress                          string    `url:"smtp_address,omitempty" required:"false" json:"smtp_address,omitempty"`
	SmtpAuthentication                   string    `url:"smtp_authentication,omitempty" required:"false" json:"smtp_authentication,omitempty"`
	SmtpFrom                             string    `url:"smtp_from,omitempty" required:"false" json:"smtp_from,omitempty"`
	SmtpUsername                         string    `url:"smtp_username,omitempty" required:"false" json:"smtp_username,omitempty"`
	SmtpPort                             int64     `url:"smtp_port,omitempty" required:"false" json:"smtp_port,omitempty"`
	LdapEnabled                          *bool     `url:"ldap_enabled,omitempty" required:"false" json:"ldap_enabled,omitempty"`
	LdapType                             string    `url:"ldap_type,omitempty" required:"false" json:"ldap_type,omitempty"`
	LdapHost                             string    `url:"ldap_host,omitempty" required:"false" json:"ldap_host,omitempty"`
	LdapHost2                            string    `url:"ldap_host_2,omitempty" required:"false" json:"ldap_host_2,omitempty"`
	LdapHost3                            string    `url:"ldap_host_3,omitempty" required:"false" json:"ldap_host_3,omitempty"`
	LdapPort                             int64     `url:"ldap_port,omitempty" required:"false" json:"ldap_port,omitempty"`
	LdapSecure                           *bool     `url:"ldap_secure,omitempty" required:"false" json:"ldap_secure,omitempty"`
	LdapUsername                         string    `url:"ldap_username,omitempty" required:"false" json:"ldap_username,omitempty"`
	LdapUsernameField                    string    `url:"ldap_username_field,omitempty" required:"false" json:"ldap_username_field,omitempty"`
	LdapDomain                           string    `url:"ldap_domain,omitempty" required:"false" json:"ldap_domain,omitempty"`
	LdapUserAction                       string    `url:"ldap_user_action,omitempty" required:"false" json:"ldap_user_action,omitempty"`
	LdapGroupAction                      string    `url:"ldap_group_action,omitempty" required:"false" json:"ldap_group_action,omitempty"`
	LdapUserIncludeGroups                string    `url:"ldap_user_include_groups,omitempty" required:"false" json:"ldap_user_include_groups,omitempty"`
	LdapGroupExclusion                   string    `url:"ldap_group_exclusion,omitempty" required:"false" json:"ldap_group_exclusion,omitempty"`
	LdapGroupInclusion                   string    `url:"ldap_group_inclusion,omitempty" required:"false" json:"ldap_group_inclusion,omitempty"`
	LdapBaseDn                           string    `url:"ldap_base_dn,omitempty" required:"false" json:"ldap_base_dn,omitempty"`
	Icon16File                           io.Writer `url:"icon16_file,omitempty" required:"false" json:"icon16_file,omitempty"`
	Icon16Delete                         *bool     `url:"icon16_delete,omitempty" required:"false" json:"icon16_delete,omitempty"`
	Icon32File                           io.Writer `url:"icon32_file,omitempty" required:"false" json:"icon32_file,omitempty"`
	Icon32Delete                         *bool     `url:"icon32_delete,omitempty" required:"false" json:"icon32_delete,omitempty"`
	Icon48File                           io.Writer `url:"icon48_file,omitempty" required:"false" json:"icon48_file,omitempty"`
	Icon48Delete                         *bool     `url:"icon48_delete,omitempty" required:"false" json:"icon48_delete,omitempty"`
	Icon128File                          io.Writer `url:"icon128_file,omitempty" required:"false" json:"icon128_file,omitempty"`
	Icon128Delete                        *bool     `url:"icon128_delete,omitempty" required:"false" json:"icon128_delete,omitempty"`
	LogoFile                             io.Writer `url:"logo_file,omitempty" required:"false" json:"logo_file,omitempty"`
	LogoDelete                           *bool     `url:"logo_delete,omitempty" required:"false" json:"logo_delete,omitempty"`
	BundleWatermarkAttachmentFile        io.Writer `url:"bundle_watermark_attachment_file,omitempty" required:"false" json:"bundle_watermark_attachment_file,omitempty"`
	BundleWatermarkAttachmentDelete      *bool     `url:"bundle_watermark_attachment_delete,omitempty" required:"false" json:"bundle_watermark_attachment_delete,omitempty"`
	Disable2faWithDelay                  *bool     `url:"disable_2fa_with_delay,omitempty" required:"false" json:"disable_2fa_with_delay,omitempty"`
	LdapPasswordChange                   string    `url:"ldap_password_change,omitempty" required:"false" json:"ldap_password_change,omitempty"`
	LdapPasswordChangeConfirmation       string    `url:"ldap_password_change_confirmation,omitempty" required:"false" json:"ldap_password_change_confirmation,omitempty"`
	SmtpPassword                         string    `url:"smtp_password,omitempty" required:"false" json:"smtp_password,omitempty"`
}

func (s *Site) UnmarshalJSON(data []byte) error {
	type site Site
	var v site
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*s = Site(v)
	return nil
}

func (s *SiteCollection) UnmarshalJSON(data []byte) error {
	type sites []Site
	var v sites
	if err := json.Unmarshal(data, &v); err != nil {
		return err
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
