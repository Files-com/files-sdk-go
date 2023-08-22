package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type StagingSite struct {
	Name                                   string                 `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	Allowed2faMethodSms                    *bool                  `json:"allowed_2fa_method_sms,omitempty" path:"allowed_2fa_method_sms,omitempty" url:"allowed_2fa_method_sms,omitempty"`
	Allowed2faMethodTotp                   *bool                  `json:"allowed_2fa_method_totp,omitempty" path:"allowed_2fa_method_totp,omitempty" url:"allowed_2fa_method_totp,omitempty"`
	Allowed2faMethodU2f                    *bool                  `json:"allowed_2fa_method_u2f,omitempty" path:"allowed_2fa_method_u2f,omitempty" url:"allowed_2fa_method_u2f,omitempty"`
	Allowed2faMethodWebauthn               *bool                  `json:"allowed_2fa_method_webauthn,omitempty" path:"allowed_2fa_method_webauthn,omitempty" url:"allowed_2fa_method_webauthn,omitempty"`
	Allowed2faMethodYubi                   *bool                  `json:"allowed_2fa_method_yubi,omitempty" path:"allowed_2fa_method_yubi,omitempty" url:"allowed_2fa_method_yubi,omitempty"`
	Allowed2faMethodBypassForFtpSftpDav    *bool                  `json:"allowed_2fa_method_bypass_for_ftp_sftp_dav,omitempty" path:"allowed_2fa_method_bypass_for_ftp_sftp_dav,omitempty" url:"allowed_2fa_method_bypass_for_ftp_sftp_dav,omitempty"`
	AdminUserId                            int64                  `json:"admin_user_id,omitempty" path:"admin_user_id,omitempty" url:"admin_user_id,omitempty"`
	AllowBundleNames                       *bool                  `json:"allow_bundle_names,omitempty" path:"allow_bundle_names,omitempty" url:"allow_bundle_names,omitempty"`
	AllowedCountries                       string                 `json:"allowed_countries,omitempty" path:"allowed_countries,omitempty" url:"allowed_countries,omitempty"`
	AllowedIps                             string                 `json:"allowed_ips,omitempty" path:"allowed_ips,omitempty" url:"allowed_ips,omitempty"`
	AskAboutOverwrites                     *bool                  `json:"ask_about_overwrites,omitempty" path:"ask_about_overwrites,omitempty" url:"ask_about_overwrites,omitempty"`
	BundleActivityNotifications            string                 `json:"bundle_activity_notifications,omitempty" path:"bundle_activity_notifications,omitempty" url:"bundle_activity_notifications,omitempty"`
	BundleExpiration                       int64                  `json:"bundle_expiration,omitempty" path:"bundle_expiration,omitempty" url:"bundle_expiration,omitempty"`
	BundlePasswordRequired                 *bool                  `json:"bundle_password_required,omitempty" path:"bundle_password_required,omitempty" url:"bundle_password_required,omitempty"`
	BundleRegistrationNotifications        string                 `json:"bundle_registration_notifications,omitempty" path:"bundle_registration_notifications,omitempty" url:"bundle_registration_notifications,omitempty"`
	BundleRequireRegistration              *bool                  `json:"bundle_require_registration,omitempty" path:"bundle_require_registration,omitempty" url:"bundle_require_registration,omitempty"`
	BundleRequireShareRecipient            *bool                  `json:"bundle_require_share_recipient,omitempty" path:"bundle_require_share_recipient,omitempty" url:"bundle_require_share_recipient,omitempty"`
	BundleUploadReceiptNotifications       string                 `json:"bundle_upload_receipt_notifications,omitempty" path:"bundle_upload_receipt_notifications,omitempty" url:"bundle_upload_receipt_notifications,omitempty"`
	BundleWatermarkAttachment              Image                  `json:"bundle_watermark_attachment,omitempty" path:"bundle_watermark_attachment,omitempty" url:"bundle_watermark_attachment,omitempty"`
	BundleWatermarkValue                   map[string]interface{} `json:"bundle_watermark_value,omitempty" path:"bundle_watermark_value,omitempty" url:"bundle_watermark_value,omitempty"`
	UploadsViaEmailAuthentication          *bool                  `json:"uploads_via_email_authentication,omitempty" path:"uploads_via_email_authentication,omitempty" url:"uploads_via_email_authentication,omitempty"`
	Color2Left                             string                 `json:"color2_left,omitempty" path:"color2_left,omitempty" url:"color2_left,omitempty"`
	Color2Link                             string                 `json:"color2_link,omitempty" path:"color2_link,omitempty" url:"color2_link,omitempty"`
	Color2Text                             string                 `json:"color2_text,omitempty" path:"color2_text,omitempty" url:"color2_text,omitempty"`
	Color2Top                              string                 `json:"color2_top,omitempty" path:"color2_top,omitempty" url:"color2_top,omitempty"`
	Color2TopText                          string                 `json:"color2_top_text,omitempty" path:"color2_top_text,omitempty" url:"color2_top_text,omitempty"`
	ContactName                            string                 `json:"contact_name,omitempty" path:"contact_name,omitempty" url:"contact_name,omitempty"`
	CreatedAt                              *time.Time             `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	Currency                               string                 `json:"currency,omitempty" path:"currency,omitempty" url:"currency,omitempty"`
	CustomNamespace                        *bool                  `json:"custom_namespace,omitempty" path:"custom_namespace,omitempty" url:"custom_namespace,omitempty"`
	DaysToRetainBackups                    int64                  `json:"days_to_retain_backups,omitempty" path:"days_to_retain_backups,omitempty" url:"days_to_retain_backups,omitempty"`
	DefaultTimeZone                        string                 `json:"default_time_zone,omitempty" path:"default_time_zone,omitempty" url:"default_time_zone,omitempty"`
	DesktopApp                             *bool                  `json:"desktop_app,omitempty" path:"desktop_app,omitempty" url:"desktop_app,omitempty"`
	DesktopAppSessionIpPinning             *bool                  `json:"desktop_app_session_ip_pinning,omitempty" path:"desktop_app_session_ip_pinning,omitempty" url:"desktop_app_session_ip_pinning,omitempty"`
	DesktopAppSessionLifetime              int64                  `json:"desktop_app_session_lifetime,omitempty" path:"desktop_app_session_lifetime,omitempty" url:"desktop_app_session_lifetime,omitempty"`
	MobileApp                              *bool                  `json:"mobile_app,omitempty" path:"mobile_app,omitempty" url:"mobile_app,omitempty"`
	MobileAppSessionIpPinning              *bool                  `json:"mobile_app_session_ip_pinning,omitempty" path:"mobile_app_session_ip_pinning,omitempty" url:"mobile_app_session_ip_pinning,omitempty"`
	MobileAppSessionLifetime               int64                  `json:"mobile_app_session_lifetime,omitempty" path:"mobile_app_session_lifetime,omitempty" url:"mobile_app_session_lifetime,omitempty"`
	DisallowedCountries                    string                 `json:"disallowed_countries,omitempty" path:"disallowed_countries,omitempty" url:"disallowed_countries,omitempty"`
	DisableFilesCertificateGeneration      *bool                  `json:"disable_files_certificate_generation,omitempty" path:"disable_files_certificate_generation,omitempty" url:"disable_files_certificate_generation,omitempty"`
	DisableNotifications                   *bool                  `json:"disable_notifications,omitempty" path:"disable_notifications,omitempty" url:"disable_notifications,omitempty"`
	DisablePasswordReset                   *bool                  `json:"disable_password_reset,omitempty" path:"disable_password_reset,omitempty" url:"disable_password_reset,omitempty"`
	Domain                                 string                 `json:"domain,omitempty" path:"domain,omitempty" url:"domain,omitempty"`
	DomainHstsHeader                       *bool                  `json:"domain_hsts_header,omitempty" path:"domain_hsts_header,omitempty" url:"domain_hsts_header,omitempty"`
	DomainLetsencryptChain                 string                 `json:"domain_letsencrypt_chain,omitempty" path:"domain_letsencrypt_chain,omitempty" url:"domain_letsencrypt_chain,omitempty"`
	Email                                  string                 `json:"email,omitempty" path:"email,omitempty" url:"email,omitempty"`
	FtpEnabled                             *bool                  `json:"ftp_enabled,omitempty" path:"ftp_enabled,omitempty" url:"ftp_enabled,omitempty"`
	ReplyToEmail                           string                 `json:"reply_to_email,omitempty" path:"reply_to_email,omitempty" url:"reply_to_email,omitempty"`
	NonSsoGroupsAllowed                    *bool                  `json:"non_sso_groups_allowed,omitempty" path:"non_sso_groups_allowed,omitempty" url:"non_sso_groups_allowed,omitempty"`
	NonSsoUsersAllowed                     *bool                  `json:"non_sso_users_allowed,omitempty" path:"non_sso_users_allowed,omitempty" url:"non_sso_users_allowed,omitempty"`
	FolderPermissionsGroupsOnly            *bool                  `json:"folder_permissions_groups_only,omitempty" path:"folder_permissions_groups_only,omitempty" url:"folder_permissions_groups_only,omitempty"`
	Hipaa                                  *bool                  `json:"hipaa,omitempty" path:"hipaa,omitempty" url:"hipaa,omitempty"`
	Icon128                                Image                  `json:"icon128,omitempty" path:"icon128,omitempty" url:"icon128,omitempty"`
	Icon16                                 Image                  `json:"icon16,omitempty" path:"icon16,omitempty" url:"icon16,omitempty"`
	Icon32                                 Image                  `json:"icon32,omitempty" path:"icon32,omitempty" url:"icon32,omitempty"`
	Icon48                                 Image                  `json:"icon48,omitempty" path:"icon48,omitempty" url:"icon48,omitempty"`
	ImmutableFilesSetAt                    *time.Time             `json:"immutable_files_set_at,omitempty" path:"immutable_files_set_at,omitempty" url:"immutable_files_set_at,omitempty"`
	IncludePasswordInWelcomeEmail          *bool                  `json:"include_password_in_welcome_email,omitempty" path:"include_password_in_welcome_email,omitempty" url:"include_password_in_welcome_email,omitempty"`
	Language                               string                 `json:"language,omitempty" path:"language,omitempty" url:"language,omitempty"`
	LdapBaseDn                             string                 `json:"ldap_base_dn,omitempty" path:"ldap_base_dn,omitempty" url:"ldap_base_dn,omitempty"`
	LdapDomain                             string                 `json:"ldap_domain,omitempty" path:"ldap_domain,omitempty" url:"ldap_domain,omitempty"`
	LdapEnabled                            *bool                  `json:"ldap_enabled,omitempty" path:"ldap_enabled,omitempty" url:"ldap_enabled,omitempty"`
	LdapGroupAction                        string                 `json:"ldap_group_action,omitempty" path:"ldap_group_action,omitempty" url:"ldap_group_action,omitempty"`
	LdapGroupExclusion                     string                 `json:"ldap_group_exclusion,omitempty" path:"ldap_group_exclusion,omitempty" url:"ldap_group_exclusion,omitempty"`
	LdapGroupInclusion                     string                 `json:"ldap_group_inclusion,omitempty" path:"ldap_group_inclusion,omitempty" url:"ldap_group_inclusion,omitempty"`
	LdapHost                               string                 `json:"ldap_host,omitempty" path:"ldap_host,omitempty" url:"ldap_host,omitempty"`
	LdapHost2                              string                 `json:"ldap_host_2,omitempty" path:"ldap_host_2,omitempty" url:"ldap_host_2,omitempty"`
	LdapHost3                              string                 `json:"ldap_host_3,omitempty" path:"ldap_host_3,omitempty" url:"ldap_host_3,omitempty"`
	LdapPort                               int64                  `json:"ldap_port,omitempty" path:"ldap_port,omitempty" url:"ldap_port,omitempty"`
	LdapSecure                             *bool                  `json:"ldap_secure,omitempty" path:"ldap_secure,omitempty" url:"ldap_secure,omitempty"`
	LdapType                               string                 `json:"ldap_type,omitempty" path:"ldap_type,omitempty" url:"ldap_type,omitempty"`
	LdapUserAction                         string                 `json:"ldap_user_action,omitempty" path:"ldap_user_action,omitempty" url:"ldap_user_action,omitempty"`
	LdapUserIncludeGroups                  string                 `json:"ldap_user_include_groups,omitempty" path:"ldap_user_include_groups,omitempty" url:"ldap_user_include_groups,omitempty"`
	LdapUsername                           string                 `json:"ldap_username,omitempty" path:"ldap_username,omitempty" url:"ldap_username,omitempty"`
	LdapUsernameField                      string                 `json:"ldap_username_field,omitempty" path:"ldap_username_field,omitempty" url:"ldap_username_field,omitempty"`
	LoginHelpText                          string                 `json:"login_help_text,omitempty" path:"login_help_text,omitempty" url:"login_help_text,omitempty"`
	Logo                                   Image                  `json:"logo,omitempty" path:"logo,omitempty" url:"logo,omitempty"`
	MaxPriorPasswords                      int64                  `json:"max_prior_passwords,omitempty" path:"max_prior_passwords,omitempty" url:"max_prior_passwords,omitempty"`
	MotdText                               string                 `json:"motd_text,omitempty" path:"motd_text,omitempty" url:"motd_text,omitempty"`
	MotdUseForFtp                          *bool                  `json:"motd_use_for_ftp,omitempty" path:"motd_use_for_ftp,omitempty" url:"motd_use_for_ftp,omitempty"`
	MotdUseForSftp                         *bool                  `json:"motd_use_for_sftp,omitempty" path:"motd_use_for_sftp,omitempty" url:"motd_use_for_sftp,omitempty"`
	NextBillingAmount                      string                 `json:"next_billing_amount,omitempty" path:"next_billing_amount,omitempty" url:"next_billing_amount,omitempty"`
	NextBillingDate                        string                 `json:"next_billing_date,omitempty" path:"next_billing_date,omitempty" url:"next_billing_date,omitempty"`
	OfficeIntegrationAvailable             *bool                  `json:"office_integration_available,omitempty" path:"office_integration_available,omitempty" url:"office_integration_available,omitempty"`
	OfficeIntegrationType                  string                 `json:"office_integration_type,omitempty" path:"office_integration_type,omitempty" url:"office_integration_type,omitempty"`
	OncehubLink                            string                 `json:"oncehub_link,omitempty" path:"oncehub_link,omitempty" url:"oncehub_link,omitempty"`
	OptOutGlobal                           *bool                  `json:"opt_out_global,omitempty" path:"opt_out_global,omitempty" url:"opt_out_global,omitempty"`
	Overdue                                *bool                  `json:"overdue,omitempty" path:"overdue,omitempty" url:"overdue,omitempty"`
	PasswordMinLength                      int64                  `json:"password_min_length,omitempty" path:"password_min_length,omitempty" url:"password_min_length,omitempty"`
	PasswordRequireLetter                  *bool                  `json:"password_require_letter,omitempty" path:"password_require_letter,omitempty" url:"password_require_letter,omitempty"`
	PasswordRequireMixed                   *bool                  `json:"password_require_mixed,omitempty" path:"password_require_mixed,omitempty" url:"password_require_mixed,omitempty"`
	PasswordRequireNumber                  *bool                  `json:"password_require_number,omitempty" path:"password_require_number,omitempty" url:"password_require_number,omitempty"`
	PasswordRequireSpecial                 *bool                  `json:"password_require_special,omitempty" path:"password_require_special,omitempty" url:"password_require_special,omitempty"`
	PasswordRequireUnbreached              *bool                  `json:"password_require_unbreached,omitempty" path:"password_require_unbreached,omitempty" url:"password_require_unbreached,omitempty"`
	PasswordRequirementsApplyToBundles     *bool                  `json:"password_requirements_apply_to_bundles,omitempty" path:"password_requirements_apply_to_bundles,omitempty" url:"password_requirements_apply_to_bundles,omitempty"`
	PasswordValidityDays                   int64                  `json:"password_validity_days,omitempty" path:"password_validity_days,omitempty" url:"password_validity_days,omitempty"`
	Phone                                  string                 `json:"phone,omitempty" path:"phone,omitempty" url:"phone,omitempty"`
	PinAllRemoteServersToSiteRegion        *bool                  `json:"pin_all_remote_servers_to_site_region,omitempty" path:"pin_all_remote_servers_to_site_region,omitempty" url:"pin_all_remote_servers_to_site_region,omitempty"`
	PreventRootPermissionsForNonSiteAdmins *bool                  `json:"prevent_root_permissions_for_non_site_admins,omitempty" path:"prevent_root_permissions_for_non_site_admins,omitempty" url:"prevent_root_permissions_for_non_site_admins,omitempty"`
	Require2fa                             *bool                  `json:"require_2fa,omitempty" path:"require_2fa,omitempty" url:"require_2fa,omitempty"`
	Require2faStopTime                     *time.Time             `json:"require_2fa_stop_time,omitempty" path:"require_2fa_stop_time,omitempty" url:"require_2fa_stop_time,omitempty"`
	Require2faUserType                     string                 `json:"require_2fa_user_type,omitempty" path:"require_2fa_user_type,omitempty" url:"require_2fa_user_type,omitempty"`
	RequireLogoutFromBundlesAndInboxes     *bool                  `json:"require_logout_from_bundles_and_inboxes,omitempty" path:"require_logout_from_bundles_and_inboxes,omitempty" url:"require_logout_from_bundles_and_inboxes,omitempty"`
	Session                                Session                `json:"session,omitempty" path:"session,omitempty" url:"session,omitempty"`
	SessionPinnedByIp                      *bool                  `json:"session_pinned_by_ip,omitempty" path:"session_pinned_by_ip,omitempty" url:"session_pinned_by_ip,omitempty"`
	SftpEnabled                            *bool                  `json:"sftp_enabled,omitempty" path:"sftp_enabled,omitempty" url:"sftp_enabled,omitempty"`
	SftpHostKeyType                        string                 `json:"sftp_host_key_type,omitempty" path:"sftp_host_key_type,omitempty" url:"sftp_host_key_type,omitempty"`
	ActiveSftpHostKeyId                    int64                  `json:"active_sftp_host_key_id,omitempty" path:"active_sftp_host_key_id,omitempty" url:"active_sftp_host_key_id,omitempty"`
	SftpInsecureCiphers                    *bool                  `json:"sftp_insecure_ciphers,omitempty" path:"sftp_insecure_ciphers,omitempty" url:"sftp_insecure_ciphers,omitempty"`
	SftpUserRootEnabled                    *bool                  `json:"sftp_user_root_enabled,omitempty" path:"sftp_user_root_enabled,omitempty" url:"sftp_user_root_enabled,omitempty"`
	SharingEnabled                         *bool                  `json:"sharing_enabled,omitempty" path:"sharing_enabled,omitempty" url:"sharing_enabled,omitempty"`
	ShowRequestAccessLink                  *bool                  `json:"show_request_access_link,omitempty" path:"show_request_access_link,omitempty" url:"show_request_access_link,omitempty"`
	SiteFooter                             string                 `json:"site_footer,omitempty" path:"site_footer,omitempty" url:"site_footer,omitempty"`
	SiteHeader                             string                 `json:"site_header,omitempty" path:"site_header,omitempty" url:"site_header,omitempty"`
	SmtpAddress                            string                 `json:"smtp_address,omitempty" path:"smtp_address,omitempty" url:"smtp_address,omitempty"`
	SmtpAuthentication                     string                 `json:"smtp_authentication,omitempty" path:"smtp_authentication,omitempty" url:"smtp_authentication,omitempty"`
	SmtpFrom                               string                 `json:"smtp_from,omitempty" path:"smtp_from,omitempty" url:"smtp_from,omitempty"`
	SmtpPort                               int64                  `json:"smtp_port,omitempty" path:"smtp_port,omitempty" url:"smtp_port,omitempty"`
	SmtpUsername                           string                 `json:"smtp_username,omitempty" path:"smtp_username,omitempty" url:"smtp_username,omitempty"`
	SessionExpiry                          string                 `json:"session_expiry,omitempty" path:"session_expiry,omitempty" url:"session_expiry,omitempty"`
	SessionExpiryMinutes                   int64                  `json:"session_expiry_minutes,omitempty" path:"session_expiry_minutes,omitempty" url:"session_expiry_minutes,omitempty"`
	SslRequired                            *bool                  `json:"ssl_required,omitempty" path:"ssl_required,omitempty" url:"ssl_required,omitempty"`
	Subdomain                              string                 `json:"subdomain,omitempty" path:"subdomain,omitempty" url:"subdomain,omitempty"`
	SwitchToPlanDate                       *time.Time             `json:"switch_to_plan_date,omitempty" path:"switch_to_plan_date,omitempty" url:"switch_to_plan_date,omitempty"`
	TlsDisabled                            *bool                  `json:"tls_disabled,omitempty" path:"tls_disabled,omitempty" url:"tls_disabled,omitempty"`
	TrialDaysLeft                          int64                  `json:"trial_days_left,omitempty" path:"trial_days_left,omitempty" url:"trial_days_left,omitempty"`
	TrialUntil                             *time.Time             `json:"trial_until,omitempty" path:"trial_until,omitempty" url:"trial_until,omitempty"`
	UpdatedAt                              *time.Time             `json:"updated_at,omitempty" path:"updated_at,omitempty" url:"updated_at,omitempty"`
	UseProvidedModifiedAt                  *bool                  `json:"use_provided_modified_at,omitempty" path:"use_provided_modified_at,omitempty" url:"use_provided_modified_at,omitempty"`
	User                                   User                   `json:"user,omitempty" path:"user,omitempty" url:"user,omitempty"`
	UserLockout                            *bool                  `json:"user_lockout,omitempty" path:"user_lockout,omitempty" url:"user_lockout,omitempty"`
	UserLockoutLockPeriod                  int64                  `json:"user_lockout_lock_period,omitempty" path:"user_lockout_lock_period,omitempty" url:"user_lockout_lock_period,omitempty"`
	UserLockoutTries                       int64                  `json:"user_lockout_tries,omitempty" path:"user_lockout_tries,omitempty" url:"user_lockout_tries,omitempty"`
	UserLockoutWithin                      int64                  `json:"user_lockout_within,omitempty" path:"user_lockout_within,omitempty" url:"user_lockout_within,omitempty"`
	UserRequestsEnabled                    *bool                  `json:"user_requests_enabled,omitempty" path:"user_requests_enabled,omitempty" url:"user_requests_enabled,omitempty"`
	UserRequestsNotifyAdmins               *bool                  `json:"user_requests_notify_admins,omitempty" path:"user_requests_notify_admins,omitempty" url:"user_requests_notify_admins,omitempty"`
	WelcomeCustomText                      string                 `json:"welcome_custom_text,omitempty" path:"welcome_custom_text,omitempty" url:"welcome_custom_text,omitempty"`
	WelcomeEmailCc                         string                 `json:"welcome_email_cc,omitempty" path:"welcome_email_cc,omitempty" url:"welcome_email_cc,omitempty"`
	WelcomeEmailSubject                    string                 `json:"welcome_email_subject,omitempty" path:"welcome_email_subject,omitempty" url:"welcome_email_subject,omitempty"`
	WelcomeEmailEnabled                    *bool                  `json:"welcome_email_enabled,omitempty" path:"welcome_email_enabled,omitempty" url:"welcome_email_enabled,omitempty"`
	WelcomeScreen                          string                 `json:"welcome_screen,omitempty" path:"welcome_screen,omitempty" url:"welcome_screen,omitempty"`
	WindowsModeFtp                         *bool                  `json:"windows_mode_ftp,omitempty" path:"windows_mode_ftp,omitempty" url:"windows_mode_ftp,omitempty"`
	DisableUsersFromInactivityPeriodDays   int64                  `json:"disable_users_from_inactivity_period_days,omitempty" path:"disable_users_from_inactivity_period_days,omitempty" url:"disable_users_from_inactivity_period_days,omitempty"`
	GroupAdminsCanSetUserPassword          *bool                  `json:"group_admins_can_set_user_password,omitempty" path:"group_admins_can_set_user_password,omitempty" url:"group_admins_can_set_user_password,omitempty"`
	Username                               string                 `json:"username,omitempty" path:"username,omitempty" url:"username,omitempty"`
	Password                               string                 `json:"password,omitempty" path:"password,omitempty" url:"password,omitempty"`
}

// Identifier no path or id

type StagingSiteCollection []StagingSite

type StagingSiteListParams struct {
	Action string `url:"action,omitempty" required:"false" json:"action,omitempty" path:"action"`
	ListParams
}

type StagingSiteCreateParams struct {
	Name     string `url:"name,omitempty" required:"true" json:"name,omitempty" path:"name"`
	Email    string `url:"email,omitempty" required:"true" json:"email,omitempty" path:"email"`
	Username string `url:"username,omitempty" required:"false" json:"username,omitempty" path:"username"`
	Password string `url:"password,omitempty" required:"false" json:"password,omitempty" path:"password"`
}

func (s *StagingSite) UnmarshalJSON(data []byte) error {
	type stagingSite StagingSite
	var v stagingSite
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*s = StagingSite(v)
	return nil
}

func (s *StagingSiteCollection) UnmarshalJSON(data []byte) error {
	type stagingSites StagingSiteCollection
	var v stagingSites
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*s = StagingSiteCollection(v)
	return nil
}

func (s *StagingSiteCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*s))
	for i, v := range *s {
		ret[i] = v
	}

	return &ret
}
