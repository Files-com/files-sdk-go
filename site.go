package files_sdk

import (
	"encoding/json"
	"io"
	"time"
)

type Site struct {
	Name                                 string    `json:"name,omitempty"`
	Allowed2faMethodSms                  *bool     `json:"allowed_2fa_method_sms,omitempty"`
	Allowed2faMethodTotp                 *bool     `json:"allowed_2fa_method_totp,omitempty"`
	Allowed2faMethodU2f                  *bool     `json:"allowed_2fa_method_u2f,omitempty"`
	Allowed2faMethodYubi                 *bool     `json:"allowed_2fa_method_yubi,omitempty"`
	AdminUserId                          int64     `json:"admin_user_id,omitempty"`
	AllowBundleNames                     *bool     `json:"allow_bundle_names,omitempty"`
	AllowedCountries                     string    `json:"allowed_countries,omitempty"`
	AllowedIps                           string    `json:"allowed_ips,omitempty"`
	AskAboutOverwrites                   *bool     `json:"ask_about_overwrites,omitempty"`
	BundleExpiration                     int       `json:"bundle_expiration,omitempty"`
	BundlePasswordRequired               *bool     `json:"bundle_password_required,omitempty"`
	Color2Left                           string    `json:"color2_left,omitempty"`
	Color2Link                           string    `json:"color2_link,omitempty"`
	Color2Text                           string    `json:"color2_text,omitempty"`
	Color2Top                            string    `json:"color2_top,omitempty"`
	Color2TopText                        string    `json:"color2_top_text,omitempty"`
	CreatedAt                            time.Time `json:"created_at,omitempty"`
	Currency                             string    `json:"currency,omitempty"`
	CustomNamespace                      *bool     `json:"custom_namespace,omitempty"`
	DaysToRetainBackups                  int       `json:"days_to_retain_backups,omitempty"`
	DefaultTimeZone                      string    `json:"default_time_zone,omitempty"`
	DesktopApp                           *bool     `json:"desktop_app,omitempty"`
	DesktopAppSessionIpPinning           *bool     `json:"desktop_app_session_ip_pinning,omitempty"`
	DesktopAppSessionLifetime            int       `json:"desktop_app_session_lifetime,omitempty"`
	DisallowedCountries                  string    `json:"disallowed_countries,omitempty"`
	DisableNotifications                 *bool     `json:"disable_notifications,omitempty"`
	DisablePasswordReset                 *bool     `json:"disable_password_reset,omitempty"`
	Domain                               string    `json:"domain,omitempty"`
	Email                                string    `json:"email,omitempty"`
	ReplyToEmail                         string    `json:"reply_to_email,omitempty"`
	NonSsoGroupsAllowed                  *bool     `json:"non_sso_groups_allowed,omitempty"`
	NonSsoUsersAllowed                   *bool     `json:"non_sso_users_allowed,omitempty"`
	FolderPermissionsGroupsOnly          *bool     `json:"folder_permissions_groups_only,omitempty"`
	Hipaa                                *bool     `json:"hipaa,omitempty"`
	Icon128                              string    `json:"icon128,omitempty"`
	Icon16                               string    `json:"icon16,omitempty"`
	Icon32                               string    `json:"icon32,omitempty"`
	Icon48                               string    `json:"icon48,omitempty"`
	ImmutableFilesSetAt                  time.Time `json:"immutable_files_set_at,omitempty"`
	IncludePasswordInWelcomeEmail        *bool     `json:"include_password_in_welcome_email,omitempty"`
	Language                             string    `json:"language,omitempty"`
	LdapBaseDn                           string    `json:"ldap_base_dn,omitempty"`
	LdapDomain                           string    `json:"ldap_domain,omitempty"`
	LdapEnabled                          *bool     `json:"ldap_enabled,omitempty"`
	LdapGroupAction                      string    `json:"ldap_group_action,omitempty"`
	LdapGroupExclusion                   string    `json:"ldap_group_exclusion,omitempty"`
	LdapGroupInclusion                   string    `json:"ldap_group_inclusion,omitempty"`
	LdapHost                             string    `json:"ldap_host,omitempty"`
	LdapHost2                            string    `json:"ldap_host_2,omitempty"`
	LdapHost3                            string    `json:"ldap_host_3,omitempty"`
	LdapPort                             int       `json:"ldap_port,omitempty"`
	LdapSecure                           *bool     `json:"ldap_secure,omitempty"`
	LdapType                             string    `json:"ldap_type,omitempty"`
	LdapUserAction                       string    `json:"ldap_user_action,omitempty"`
	LdapUserIncludeGroups                string    `json:"ldap_user_include_groups,omitempty"`
	LdapUsername                         string    `json:"ldap_username,omitempty"`
	LdapUsernameField                    string    `json:"ldap_username_field,omitempty"`
	LoginHelpText                        string    `json:"login_help_text,omitempty"`
	Logo                                 string    `json:"logo,omitempty"`
	MaxPriorPasswords                    int       `json:"max_prior_passwords,omitempty"`
	NextBillingAmount                    float32   `json:"next_billing_amount,omitempty"`
	NextBillingDate                      string    `json:"next_billing_date,omitempty"`
	OfficeIntegrationAvailable           *bool     `json:"office_integration_available,omitempty"`
	OptOutGlobal                         *bool     `json:"opt_out_global,omitempty"`
	OverageNotifiedAt                    time.Time `json:"overage_notified_at,omitempty"`
	OverageNotify                        *bool     `json:"overage_notify,omitempty"`
	Overdue                              *bool     `json:"overdue,omitempty"`
	PasswordMinLength                    int       `json:"password_min_length,omitempty"`
	PasswordRequireLetter                *bool     `json:"password_require_letter,omitempty"`
	PasswordRequireMixed                 *bool     `json:"password_require_mixed,omitempty"`
	PasswordRequireNumber                *bool     `json:"password_require_number,omitempty"`
	PasswordRequireSpecial               *bool     `json:"password_require_special,omitempty"`
	PasswordRequireUnbreached            *bool     `json:"password_require_unbreached,omitempty"`
	PasswordRequirementsApplyToBundles   *bool     `json:"password_requirements_apply_to_bundles,omitempty"`
	PasswordValidityDays                 int       `json:"password_validity_days,omitempty"`
	Phone                                string    `json:"phone,omitempty"`
	Require2fa                           *bool     `json:"require_2fa,omitempty"`
	Require2faStopTime                   time.Time `json:"require_2fa_stop_time,omitempty"`
	Require2faUserType                   string    `json:"require_2fa_user_type,omitempty"`
	Session                              string    `json:"session,omitempty"`
	SessionPinnedByIp                    *bool     `json:"session_pinned_by_ip,omitempty"`
	SftpUserRootEnabled                  *bool     `json:"sftp_user_root_enabled,omitempty"`
	SharingEnabled                       *bool     `json:"sharing_enabled,omitempty"`
	ShowRequestAccessLink                *bool     `json:"show_request_access_link,omitempty"`
	SiteFooter                           string    `json:"site_footer,omitempty"`
	SiteHeader                           string    `json:"site_header,omitempty"`
	SmtpAddress                          string    `json:"smtp_address,omitempty"`
	SmtpAuthentication                   string    `json:"smtp_authentication,omitempty"`
	SmtpFrom                             string    `json:"smtp_from,omitempty"`
	SmtpPort                             int       `json:"smtp_port,omitempty"`
	SmtpUsername                         string    `json:"smtp_username,omitempty"`
	SessionExpiry                        float32   `json:"session_expiry,omitempty"`
	SslRequired                          *bool     `json:"ssl_required,omitempty"`
	Subdomain                            string    `json:"subdomain,omitempty"`
	SwitchToPlanDate                     time.Time `json:"switch_to_plan_date,omitempty"`
	TlsDisabled                          *bool     `json:"tls_disabled,omitempty"`
	TrialDaysLeft                        int       `json:"trial_days_left,omitempty"`
	TrialUntil                           time.Time `json:"trial_until,omitempty"`
	UpdatedAt                            time.Time `json:"updated_at,omitempty"`
	UseProvidedModifiedAt                *bool     `json:"use_provided_modified_at,omitempty"`
	User                                 string    `json:"user,omitempty"`
	UserLockout                          *bool     `json:"user_lockout,omitempty"`
	UserLockoutLockPeriod                int       `json:"user_lockout_lock_period,omitempty"`
	UserLockoutTries                     int       `json:"user_lockout_tries,omitempty"`
	UserLockoutWithin                    int       `json:"user_lockout_within,omitempty"`
	UserRequestsEnabled                  *bool     `json:"user_requests_enabled,omitempty"`
	WelcomeCustomText                    string    `json:"welcome_custom_text,omitempty"`
	WelcomeEmailCc                       string    `json:"welcome_email_cc,omitempty"`
	WelcomeEmailEnabled                  *bool     `json:"welcome_email_enabled,omitempty"`
	WelcomeScreen                        string    `json:"welcome_screen,omitempty"`
	WindowsModeFtp                       *bool     `json:"windows_mode_ftp,omitempty"`
	DisableUsersFromInactivityPeriodDays int       `json:"disable_users_from_inactivity_period_days,omitempty"`
}

type SiteCollection []Site

type SiteUpdateParams struct {
	Name                                 string    `url:"name,omitempty" required:"false"`
	Subdomain                            string    `url:"subdomain,omitempty" required:"false"`
	Domain                               string    `url:"domain,omitempty" required:"false"`
	Email                                string    `url:"email,omitempty" required:"false"`
	ReplyToEmail                         string    `url:"reply_to_email,omitempty" required:"false"`
	AllowBundleNames                     *bool     `url:"allow_bundle_names,omitempty" required:"false"`
	BundleExpiration                     int       `url:"bundle_expiration,omitempty" required:"false"`
	OverageNotify                        *bool     `url:"overage_notify,omitempty" required:"false"`
	WelcomeEmailEnabled                  *bool     `url:"welcome_email_enabled,omitempty" required:"false"`
	AskAboutOverwrites                   *bool     `url:"ask_about_overwrites,omitempty" required:"false"`
	ShowRequestAccessLink                *bool     `url:"show_request_access_link,omitempty" required:"false"`
	WelcomeEmailCc                       string    `url:"welcome_email_cc,omitempty" required:"false"`
	WelcomeCustomText                    string    `url:"welcome_custom_text,omitempty" required:"false"`
	Language                             string    `url:"language,omitempty" required:"false"`
	WindowsModeFtp                       *bool     `url:"windows_mode_ftp,omitempty" required:"false"`
	DefaultTimeZone                      string    `url:"default_time_zone,omitempty" required:"false"`
	DesktopApp                           *bool     `url:"desktop_app,omitempty" required:"false"`
	DesktopAppSessionIpPinning           *bool     `url:"desktop_app_session_ip_pinning,omitempty" required:"false"`
	DesktopAppSessionLifetime            int       `url:"desktop_app_session_lifetime,omitempty" required:"false"`
	FolderPermissionsGroupsOnly          *bool     `url:"folder_permissions_groups_only,omitempty" required:"false"`
	WelcomeScreen                        string    `url:"welcome_screen,omitempty" required:"false"`
	OfficeIntegrationAvailable           *bool     `url:"office_integration_available,omitempty" required:"false"`
	SessionExpiry                        float32   `url:"session_expiry,omitempty" required:"false"`
	SslRequired                          *bool     `url:"ssl_required,omitempty" required:"false"`
	TlsDisabled                          *bool     `url:"tls_disabled,omitempty" required:"false"`
	UserLockout                          *bool     `url:"user_lockout,omitempty" required:"false"`
	UserLockoutTries                     int       `url:"user_lockout_tries,omitempty" required:"false"`
	UserLockoutWithin                    int       `url:"user_lockout_within,omitempty" required:"false"`
	UserLockoutLockPeriod                int       `url:"user_lockout_lock_period,omitempty" required:"false"`
	IncludePasswordInWelcomeEmail        *bool     `url:"include_password_in_welcome_email,omitempty" required:"false"`
	AllowedCountries                     string    `url:"allowed_countries,omitempty" required:"false"`
	AllowedIps                           string    `url:"allowed_ips,omitempty" required:"false"`
	DisallowedCountries                  string    `url:"disallowed_countries,omitempty" required:"false"`
	DaysToRetainBackups                  int       `url:"days_to_retain_backups,omitempty" required:"false"`
	MaxPriorPasswords                    int       `url:"max_prior_passwords,omitempty" required:"false"`
	PasswordValidityDays                 int       `url:"password_validity_days,omitempty" required:"false"`
	PasswordMinLength                    int       `url:"password_min_length,omitempty" required:"false"`
	PasswordRequireLetter                *bool     `url:"password_require_letter,omitempty" required:"false"`
	PasswordRequireMixed                 *bool     `url:"password_require_mixed,omitempty" required:"false"`
	PasswordRequireSpecial               *bool     `url:"password_require_special,omitempty" required:"false"`
	PasswordRequireNumber                *bool     `url:"password_require_number,omitempty" required:"false"`
	PasswordRequireUnbreached            *bool     `url:"password_require_unbreached,omitempty" required:"false"`
	SftpUserRootEnabled                  *bool     `url:"sftp_user_root_enabled,omitempty" required:"false"`
	DisablePasswordReset                 *bool     `url:"disable_password_reset,omitempty" required:"false"`
	ImmutableFiles                       *bool     `url:"immutable_files,omitempty" required:"false"`
	SessionPinnedByIp                    *bool     `url:"session_pinned_by_ip,omitempty" required:"false"`
	BundlePasswordRequired               *bool     `url:"bundle_password_required,omitempty" required:"false"`
	PasswordRequirementsApplyToBundles   *bool     `url:"password_requirements_apply_to_bundles,omitempty" required:"false"`
	OptOutGlobal                         *bool     `url:"opt_out_global,omitempty" required:"false"`
	UseProvidedModifiedAt                *bool     `url:"use_provided_modified_at,omitempty" required:"false"`
	CustomNamespace                      *bool     `url:"custom_namespace,omitempty" required:"false"`
	DisableUsersFromInactivityPeriodDays int       `url:"disable_users_from_inactivity_period_days,omitempty" required:"false"`
	NonSsoGroupsAllowed                  *bool     `url:"non_sso_groups_allowed,omitempty" required:"false"`
	NonSsoUsersAllowed                   *bool     `url:"non_sso_users_allowed,omitempty" required:"false"`
	SharingEnabled                       *bool     `url:"sharing_enabled,omitempty" required:"false"`
	UserRequestsEnabled                  *bool     `url:"user_requests_enabled,omitempty" required:"false"`
	Allowed2faMethodSms                  *bool     `url:"allowed_2fa_method_sms,omitempty" required:"false"`
	Allowed2faMethodU2f                  *bool     `url:"allowed_2fa_method_u2f,omitempty" required:"false"`
	Allowed2faMethodTotp                 *bool     `url:"allowed_2fa_method_totp,omitempty" required:"false"`
	Allowed2faMethodYubi                 *bool     `url:"allowed_2fa_method_yubi,omitempty" required:"false"`
	Require2fa                           *bool     `url:"require_2fa,omitempty" required:"false"`
	Require2faUserType                   string    `url:"require_2fa_user_type,omitempty" required:"false"`
	Color2Top                            string    `url:"color2_top,omitempty" required:"false"`
	Color2Left                           string    `url:"color2_left,omitempty" required:"false"`
	Color2Link                           string    `url:"color2_link,omitempty" required:"false"`
	Color2Text                           string    `url:"color2_text,omitempty" required:"false"`
	Color2TopText                        string    `url:"color2_top_text,omitempty" required:"false"`
	SiteHeader                           string    `url:"site_header,omitempty" required:"false"`
	SiteFooter                           string    `url:"site_footer,omitempty" required:"false"`
	LoginHelpText                        string    `url:"login_help_text,omitempty" required:"false"`
	SmtpAddress                          string    `url:"smtp_address,omitempty" required:"false"`
	SmtpAuthentication                   string    `url:"smtp_authentication,omitempty" required:"false"`
	SmtpFrom                             string    `url:"smtp_from,omitempty" required:"false"`
	SmtpUsername                         string    `url:"smtp_username,omitempty" required:"false"`
	SmtpPort                             int       `url:"smtp_port,omitempty" required:"false"`
	LdapEnabled                          *bool     `url:"ldap_enabled,omitempty" required:"false"`
	LdapType                             string    `url:"ldap_type,omitempty" required:"false"`
	LdapHost                             string    `url:"ldap_host,omitempty" required:"false"`
	LdapHost2                            string    `url:"ldap_host_2,omitempty" required:"false"`
	LdapHost3                            string    `url:"ldap_host_3,omitempty" required:"false"`
	LdapPort                             int       `url:"ldap_port,omitempty" required:"false"`
	LdapSecure                           *bool     `url:"ldap_secure,omitempty" required:"false"`
	LdapUsername                         string    `url:"ldap_username,omitempty" required:"false"`
	LdapUsernameField                    string    `url:"ldap_username_field,omitempty" required:"false"`
	LdapDomain                           string    `url:"ldap_domain,omitempty" required:"false"`
	LdapUserAction                       string    `url:"ldap_user_action,omitempty" required:"false"`
	LdapGroupAction                      string    `url:"ldap_group_action,omitempty" required:"false"`
	LdapUserIncludeGroups                string    `url:"ldap_user_include_groups,omitempty" required:"false"`
	LdapGroupExclusion                   string    `url:"ldap_group_exclusion,omitempty" required:"false"`
	LdapGroupInclusion                   string    `url:"ldap_group_inclusion,omitempty" required:"false"`
	LdapBaseDn                           string    `url:"ldap_base_dn,omitempty" required:"false"`
	Icon16File                           io.Writer `url:"icon16_file,omitempty" required:"false"`
	Icon16Delete                         *bool     `url:"icon16_delete,omitempty" required:"false"`
	Icon32File                           io.Writer `url:"icon32_file,omitempty" required:"false"`
	Icon32Delete                         *bool     `url:"icon32_delete,omitempty" required:"false"`
	Icon48File                           io.Writer `url:"icon48_file,omitempty" required:"false"`
	Icon48Delete                         *bool     `url:"icon48_delete,omitempty" required:"false"`
	Icon128File                          io.Writer `url:"icon128_file,omitempty" required:"false"`
	Icon128Delete                        *bool     `url:"icon128_delete,omitempty" required:"false"`
	LogoFile                             io.Writer `url:"logo_file,omitempty" required:"false"`
	LogoDelete                           *bool     `url:"logo_delete,omitempty" required:"false"`
	Disable2faWithDelay                  *bool     `url:"disable_2fa_with_delay,omitempty" required:"false"`
	LdapPasswordChange                   string    `url:"ldap_password_change,omitempty" required:"false"`
	LdapPasswordChangeConfirmation       string    `url:"ldap_password_change_confirmation,omitempty" required:"false"`
	SmtpPassword                         string    `url:"smtp_password,omitempty" required:"false"`
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
