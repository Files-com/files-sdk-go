package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Settings struct {
	ImageRegex                             string                 `json:"image_regex,omitempty" path:"image_regex,omitempty" url:"image_regex,omitempty"`
	VideoRegex                             string                 `json:"video_regex,omitempty" path:"video_regex,omitempty" url:"video_regex,omitempty"`
	AudioRegex                             string                 `json:"audio_regex,omitempty" path:"audio_regex,omitempty" url:"audio_regex,omitempty"`
	PdfRegex                               string                 `json:"pdf_regex,omitempty" path:"pdf_regex,omitempty" url:"pdf_regex,omitempty"`
	CurrentLanguage                        string                 `json:"current_language,omitempty" path:"current_language,omitempty" url:"current_language,omitempty"`
	CurrentTime                            *time.Time             `json:"current_time,omitempty" path:"current_time,omitempty" url:"current_time,omitempty"`
	LinodeRegions                          []string               `json:"linode_regions,omitempty" path:"linode_regions,omitempty" url:"linode_regions,omitempty"`
	PrimarySubDomainBase                   string                 `json:"primary_sub_domain_base,omitempty" path:"primary_sub_domain_base,omitempty" url:"primary_sub_domain_base,omitempty"`
	ReadOnly                               string                 `json:"read_only,omitempty" path:"read_only,omitempty" url:"read_only,omitempty"`
	Reauth                                 *bool                  `json:"reauth,omitempty" path:"reauth,omitempty" url:"reauth,omitempty"`
	Regions                                []string               `json:"regions,omitempty" path:"regions,omitempty" url:"regions,omitempty"`
	S3Regions                              []string               `json:"s3_regions,omitempty" path:"s3_regions,omitempty" url:"s3_regions,omitempty"`
	SalesTaxRegions                        []string               `json:"sales_tax_regions,omitempty" path:"sales_tax_regions,omitempty" url:"sales_tax_regions,omitempty"`
	SessionLanguage                        string                 `json:"session_language,omitempty" path:"session_language,omitempty" url:"session_language,omitempty"`
	TabConfig                              string                 `json:"tab_config,omitempty" path:"tab_config,omitempty" url:"tab_config,omitempty"`
	BetaFeatures                           *bool                  `json:"beta_features,omitempty" path:"beta_features,omitempty" url:"beta_features,omitempty"`
	BetaFeature2                           *bool                  `json:"beta_feature2,omitempty" path:"beta_feature2,omitempty" url:"beta_feature2,omitempty"`
	BetaFeature3                           *bool                  `json:"beta_feature3,omitempty" path:"beta_feature3,omitempty" url:"beta_feature3,omitempty"`
	Color2Left                             string                 `json:"color2_left,omitempty" path:"color2_left,omitempty" url:"color2_left,omitempty"`
	Color2Link                             string                 `json:"color2_link,omitempty" path:"color2_link,omitempty" url:"color2_link,omitempty"`
	Color2Text                             string                 `json:"color2_text,omitempty" path:"color2_text,omitempty" url:"color2_text,omitempty"`
	Color2Top                              string                 `json:"color2_top,omitempty" path:"color2_top,omitempty" url:"color2_top,omitempty"`
	Color2TopText                          string                 `json:"color2_top_text,omitempty" path:"color2_top_text,omitempty" url:"color2_top_text,omitempty"`
	Domain                                 string                 `json:"domain,omitempty" path:"domain,omitempty" url:"domain,omitempty"`
	DisablePasswordReset                   *bool                  `json:"disable_password_reset,omitempty" path:"disable_password_reset,omitempty" url:"disable_password_reset,omitempty"`
	LoginHelpText                          string                 `json:"login_help_text,omitempty" path:"login_help_text,omitempty" url:"login_help_text,omitempty"`
	LoginHelpTextMarkdown                  string                 `json:"login_help_text_markdown,omitempty" path:"login_help_text_markdown,omitempty" url:"login_help_text_markdown,omitempty"`
	SiteName                               string                 `json:"site_name,omitempty" path:"site_name,omitempty" url:"site_name,omitempty"`
	OfficeIntegrationType                  string                 `json:"office_integration_type,omitempty" path:"office_integration_type,omitempty" url:"office_integration_type,omitempty"`
	OfficeIntegrationEnabled               *bool                  `json:"office_integration_enabled,omitempty" path:"office_integration_enabled,omitempty" url:"office_integration_enabled,omitempty"`
	OfficeIntegrationHost                  string                 `json:"office_integration_host,omitempty" path:"office_integration_host,omitempty" url:"office_integration_host,omitempty"`
	OncehubLink                            string                 `json:"oncehub_link,omitempty" path:"oncehub_link,omitempty" url:"oncehub_link,omitempty"`
	OfficeIntegrationAvailable             *bool                  `json:"office_integration_available,omitempty" path:"office_integration_available,omitempty" url:"office_integration_available,omitempty"`
	RequireLogoutFromBundlesAndInboxes     *bool                  `json:"require_logout_from_bundles_and_inboxes,omitempty" path:"require_logout_from_bundles_and_inboxes,omitempty" url:"require_logout_from_bundles_and_inboxes,omitempty"`
	ShowRequestAccessLink                  *bool                  `json:"show_request_access_link,omitempty" path:"show_request_access_link,omitempty" url:"show_request_access_link,omitempty"`
	SiteFooter                             string                 `json:"site_footer,omitempty" path:"site_footer,omitempty" url:"site_footer,omitempty"`
	SiteHeader                             string                 `json:"site_header,omitempty" path:"site_header,omitempty" url:"site_header,omitempty"`
	SiteFooterMarkdown                     string                 `json:"site_footer_markdown,omitempty" path:"site_footer_markdown,omitempty" url:"site_footer_markdown,omitempty"`
	SiteHeaderMarkdown                     string                 `json:"site_header_markdown,omitempty" path:"site_header_markdown,omitempty" url:"site_header_markdown,omitempty"`
	SiteLanguage                           string                 `json:"site_language,omitempty" path:"site_language,omitempty" url:"site_language,omitempty"`
	SsoStrategies                          []string               `json:"sso_strategies,omitempty" path:"sso_strategies,omitempty" url:"sso_strategies,omitempty"`
	Subdomain                              string                 `json:"subdomain,omitempty" path:"subdomain,omitempty" url:"subdomain,omitempty"`
	UseProvidedModifiedAt                  *bool                  `json:"use_provided_modified_at,omitempty" path:"use_provided_modified_at,omitempty" url:"use_provided_modified_at,omitempty"`
	UserRequestsEnabled                    *bool                  `json:"user_requests_enabled,omitempty" path:"user_requests_enabled,omitempty" url:"user_requests_enabled,omitempty"`
	WelcomeScreen                          string                 `json:"welcome_screen,omitempty" path:"welcome_screen,omitempty" url:"welcome_screen,omitempty"`
	Icon128                                Image                  `json:"icon128,omitempty" path:"icon128,omitempty" url:"icon128,omitempty"`
	Icon16                                 Image                  `json:"icon16,omitempty" path:"icon16,omitempty" url:"icon16,omitempty"`
	Icon32                                 Image                  `json:"icon32,omitempty" path:"icon32,omitempty" url:"icon32,omitempty"`
	Icon48                                 Image                  `json:"icon48,omitempty" path:"icon48,omitempty" url:"icon48,omitempty"`
	Logo                                   Image                  `json:"logo,omitempty" path:"logo,omitempty" url:"logo,omitempty"`
	LogoThumbnail                          Image                  `json:"logo_thumbnail,omitempty" path:"logo_thumbnail,omitempty" url:"logo_thumbnail,omitempty"`
	AttachmentsPermission                  *bool                  `json:"attachments_permission,omitempty" path:"attachments_permission,omitempty" url:"attachments_permission,omitempty"`
	AuthenticationMethod                   string                 `json:"authentication_method,omitempty" path:"authentication_method,omitempty" url:"authentication_method,omitempty"`
	AvatarUrl                              string                 `json:"avatar_url,omitempty" path:"avatar_url,omitempty" url:"avatar_url,omitempty"`
	BillingPermission                      *bool                  `json:"billing_permission,omitempty" path:"billing_permission,omitempty" url:"billing_permission,omitempty"`
	CachedPermissions                      []string               `json:"cached_permissions,omitempty" path:"cached_permissions,omitempty" url:"cached_permissions,omitempty"`
	CanAdminSomewhere                      *bool                  `json:"can_admin_somewhere,omitempty" path:"can_admin_somewhere,omitempty" url:"can_admin_somewhere,omitempty"`
	CanBundleSomewhere                     *bool                  `json:"can_bundle_somewhere,omitempty" path:"can_bundle_somewhere,omitempty" url:"can_bundle_somewhere,omitempty"`
	CanWriteSomewhere                      *bool                  `json:"can_write_somewhere,omitempty" path:"can_write_somewhere,omitempty" url:"can_write_somewhere,omitempty"`
	DavPermission                          *bool                  `json:"dav_permission,omitempty" path:"dav_permission,omitempty" url:"dav_permission,omitempty"`
	DaysRemainingUntilPasswordExpire       int64                  `json:"days_remaining_until_password_expire,omitempty" path:"days_remaining_until_password_expire,omitempty" url:"days_remaining_until_password_expire,omitempty"`
	Email                                  string                 `json:"email,omitempty" path:"email,omitempty" url:"email,omitempty"`
	FtpPermission                          *bool                  `json:"ftp_permission,omitempty" path:"ftp_permission,omitempty" url:"ftp_permission,omitempty"`
	GroupAdmin                             *bool                  `json:"group_admin,omitempty" path:"group_admin,omitempty" url:"group_admin,omitempty"`
	HeaderText                             string                 `json:"header_text,omitempty" path:"header_text,omitempty" url:"header_text,omitempty"`
	Id                                     int64                  `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	LastReadAnnouncementsAt                *time.Time             `json:"last_read_announcements_at,omitempty" path:"last_read_announcements_at,omitempty" url:"last_read_announcements_at,omitempty"`
	Name                                   string                 `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	NotificationDailySendTime              int64                  `json:"notification_daily_send_time,omitempty" path:"notification_daily_send_time,omitempty" url:"notification_daily_send_time,omitempty"`
	SelfManaged                            *bool                  `json:"self_managed,omitempty" path:"self_managed,omitempty" url:"self_managed,omitempty"`
	SftpPermission                         *bool                  `json:"sftp_permission,omitempty" path:"sftp_permission,omitempty" url:"sftp_permission,omitempty"`
	SiteAdmin                              *bool                  `json:"site_admin,omitempty" path:"site_admin,omitempty" url:"site_admin,omitempty"`
	SkipWelcomeScreen                      *bool                  `json:"skip_welcome_screen,omitempty" path:"skip_welcome_screen,omitempty" url:"skip_welcome_screen,omitempty"`
	ExternallyManaged                      *bool                  `json:"externally_managed,omitempty" path:"externally_managed,omitempty" url:"externally_managed,omitempty"`
	TimeZone                               string                 `json:"time_zone,omitempty" path:"time_zone,omitempty" url:"time_zone,omitempty"`
	TypeOf2fa                              string                 `json:"type_of_2fa,omitempty" path:"type_of_2fa,omitempty" url:"type_of_2fa,omitempty"`
	Reauth2fa                              string                 `json:"reauth_2fa,omitempty" path:"reauth_2fa,omitempty" url:"reauth_2fa,omitempty"`
	UserId                                 int64                  `json:"user_id,omitempty" path:"user_id,omitempty" url:"user_id,omitempty"`
	UserLanguage                           string                 `json:"user_language,omitempty" path:"user_language,omitempty" url:"user_language,omitempty"`
	UserRoot                               string                 `json:"user_root,omitempty" path:"user_root,omitempty" url:"user_root,omitempty"`
	Username                               string                 `json:"username,omitempty" path:"username,omitempty" url:"username,omitempty"`
	WebRoot                                string                 `json:"web_root,omitempty" path:"web_root,omitempty" url:"web_root,omitempty"`
	Allowed2faMethodBypassForFtpSftpDav    *bool                  `json:"allowed_2fa_method_bypass_for_ftp_sftp_dav,omitempty" path:"allowed_2fa_method_bypass_for_ftp_sftp_dav,omitempty" url:"allowed_2fa_method_bypass_for_ftp_sftp_dav,omitempty"`
	Allowed2faMethodSms                    *bool                  `json:"allowed_2fa_method_sms,omitempty" path:"allowed_2fa_method_sms,omitempty" url:"allowed_2fa_method_sms,omitempty"`
	Allowed2faMethodTotp                   *bool                  `json:"allowed_2fa_method_totp,omitempty" path:"allowed_2fa_method_totp,omitempty" url:"allowed_2fa_method_totp,omitempty"`
	Allowed2faMethodU2f                    *bool                  `json:"allowed_2fa_method_u2f,omitempty" path:"allowed_2fa_method_u2f,omitempty" url:"allowed_2fa_method_u2f,omitempty"`
	Allowed2faMethodWebauthn               *bool                  `json:"allowed_2fa_method_webauthn,omitempty" path:"allowed_2fa_method_webauthn,omitempty" url:"allowed_2fa_method_webauthn,omitempty"`
	Allowed2faMethodYubi                   *bool                  `json:"allowed_2fa_method_yubi,omitempty" path:"allowed_2fa_method_yubi,omitempty" url:"allowed_2fa_method_yubi,omitempty"`
	AllowBundleNames                       *bool                  `json:"allow_bundle_names,omitempty" path:"allow_bundle_names,omitempty" url:"allow_bundle_names,omitempty"`
	BundleActivityNotifications            string                 `json:"bundle_activity_notifications,omitempty" path:"bundle_activity_notifications,omitempty" url:"bundle_activity_notifications,omitempty"`
	BundleRegistrationNotifications        string                 `json:"bundle_registration_notifications,omitempty" path:"bundle_registration_notifications,omitempty" url:"bundle_registration_notifications,omitempty"`
	BundleRequireRegistration              *bool                  `json:"bundle_require_registration,omitempty" path:"bundle_require_registration,omitempty" url:"bundle_require_registration,omitempty"`
	BundleUploadReceiptNotifications       string                 `json:"bundle_upload_receipt_notifications,omitempty" path:"bundle_upload_receipt_notifications,omitempty" url:"bundle_upload_receipt_notifications,omitempty"`
	BundleExpiration                       int64                  `json:"bundle_expiration,omitempty" path:"bundle_expiration,omitempty" url:"bundle_expiration,omitempty"`
	BundlePasswordRequired                 *bool                  `json:"bundle_password_required,omitempty" path:"bundle_password_required,omitempty" url:"bundle_password_required,omitempty"`
	BundleRequireShareRecipient            *bool                  `json:"bundle_require_share_recipient,omitempty" path:"bundle_require_share_recipient,omitempty" url:"bundle_require_share_recipient,omitempty"`
	ChildSiteCountForPlan                  int64                  `json:"child_site_count_for_plan,omitempty" path:"child_site_count_for_plan,omitempty" url:"child_site_count_for_plan,omitempty"`
	DesktopApp                             *bool                  `json:"desktop_app,omitempty" path:"desktop_app,omitempty" url:"desktop_app,omitempty"`
	FeatureBundleEca                       *bool                  `json:"feature_bundle_eca,omitempty" path:"feature_bundle_eca,omitempty" url:"feature_bundle_eca,omitempty"`
	FeatureBundlePower                     *bool                  `json:"feature_bundle_power,omitempty" path:"feature_bundle_power,omitempty" url:"feature_bundle_power,omitempty"`
	FeatureBundlePremier                   *bool                  `json:"feature_bundle_premier,omitempty" path:"feature_bundle_premier,omitempty" url:"feature_bundle_premier,omitempty"`
	FolderPermissionsGroupsOnly            *bool                  `json:"folder_permissions_groups_only,omitempty" path:"folder_permissions_groups_only,omitempty" url:"folder_permissions_groups_only,omitempty"`
	GroupAdminsCanSetUserPassword          *bool                  `json:"group_admins_can_set_user_password,omitempty" path:"group_admins_can_set_user_password,omitempty" url:"group_admins_can_set_user_password,omitempty"`
	HasAccount                             *bool                  `json:"has_account,omitempty" path:"has_account,omitempty" url:"has_account,omitempty"`
	HideBilling                            *bool                  `json:"hide_billing,omitempty" path:"hide_billing,omitempty" url:"hide_billing,omitempty"`
	HighUsersCount                         *bool                  `json:"high_users_count,omitempty" path:"high_users_count,omitempty" url:"high_users_count,omitempty"`
	HistoryUnavailable                     *bool                  `json:"history_unavailable,omitempty" path:"history_unavailable,omitempty" url:"history_unavailable,omitempty"`
	ImmutableFiles                         *bool                  `json:"immutable_files,omitempty" path:"immutable_files,omitempty" url:"immutable_files,omitempty"`
	IntersitialPage                        *bool                  `json:"intersitial_page,omitempty" path:"intersitial_page,omitempty" url:"intersitial_page,omitempty"`
	LeftNavigationVisibility               map[string]interface{} `json:"left_navigation_visibility,omitempty" path:"left_navigation_visibility,omitempty" url:"left_navigation_visibility,omitempty"`
	MinRemoteSyncInterval                  int64                  `json:"min_remote_sync_interval,omitempty" path:"min_remote_sync_interval,omitempty" url:"min_remote_sync_interval,omitempty"`
	NonSsoGroupsAllowed                    *bool                  `json:"non_sso_groups_allowed,omitempty" path:"non_sso_groups_allowed,omitempty" url:"non_sso_groups_allowed,omitempty"`
	NonSsoUsersAllowed                     *bool                  `json:"non_sso_users_allowed,omitempty" path:"non_sso_users_allowed,omitempty" url:"non_sso_users_allowed,omitempty"`
	Overdue                                *bool                  `json:"overdue,omitempty" path:"overdue,omitempty" url:"overdue,omitempty"`
	SiteUnavailable                        *bool                  `json:"site_unavailable,omitempty" path:"site_unavailable,omitempty" url:"site_unavailable,omitempty"`
	PasswordMinLength                      int64                  `json:"password_min_length,omitempty" path:"password_min_length,omitempty" url:"password_min_length,omitempty"`
	PasswordRequireLetter                  *bool                  `json:"password_require_letter,omitempty" path:"password_require_letter,omitempty" url:"password_require_letter,omitempty"`
	PasswordRequireMixed                   *bool                  `json:"password_require_mixed,omitempty" path:"password_require_mixed,omitempty" url:"password_require_mixed,omitempty"`
	PasswordRequireNumber                  *bool                  `json:"password_require_number,omitempty" path:"password_require_number,omitempty" url:"password_require_number,omitempty"`
	PasswordRequireSpecial                 *bool                  `json:"password_require_special,omitempty" path:"password_require_special,omitempty" url:"password_require_special,omitempty"`
	PasswordRequireUnbreached              *bool                  `json:"password_require_unbreached,omitempty" path:"password_require_unbreached,omitempty" url:"password_require_unbreached,omitempty"`
	PasswordRequirementsApplyToBundles     *bool                  `json:"password_requirements_apply_to_bundles,omitempty" path:"password_requirements_apply_to_bundles,omitempty" url:"password_requirements_apply_to_bundles,omitempty"`
	PlanAs2Included                        *bool                  `json:"plan_as2_included,omitempty" path:"plan_as2_included,omitempty" url:"plan_as2_included,omitempty"`
	PreventRootPermissionsForNonSiteAdmins *bool                  `json:"prevent_root_permissions_for_non_site_admins,omitempty" path:"prevent_root_permissions_for_non_site_admins,omitempty" url:"prevent_root_permissions_for_non_site_admins,omitempty"`
	PublicUrl                              string                 `json:"public_url,omitempty" path:"public_url,omitempty" url:"public_url,omitempty"`
	PublicSharingAllowed                   *bool                  `json:"public_sharing_allowed,omitempty" path:"public_sharing_allowed,omitempty" url:"public_sharing_allowed,omitempty"`
	Require2fa                             *bool                  `json:"require_2fa,omitempty" path:"require_2fa,omitempty" url:"require_2fa,omitempty"`
	RootRegion                             *bool                  `json:"root_region,omitempty" path:"root_region,omitempty" url:"root_region,omitempty"`
	SharingEnabled                         *bool                  `json:"sharing_enabled,omitempty" path:"sharing_enabled,omitempty" url:"sharing_enabled,omitempty"`
	StagingSiteCountForPlan                int64                  `json:"staging_site_count_for_plan,omitempty" path:"staging_site_count_for_plan,omitempty" url:"staging_site_count_for_plan,omitempty"`
	TrialFlaggedAsDuplicate                *bool                  `json:"trial_flagged_as_duplicate,omitempty" path:"trial_flagged_as_duplicate,omitempty" url:"trial_flagged_as_duplicate,omitempty"`
	TrialDaysLeft                          int64                  `json:"trial_days_left,omitempty" path:"trial_days_left,omitempty" url:"trial_days_left,omitempty"`
	TrialLocked                            *bool                  `json:"trial_locked,omitempty" path:"trial_locked,omitempty" url:"trial_locked,omitempty"`
	TrialUntil                             *time.Time             `json:"trial_until,omitempty" path:"trial_until,omitempty" url:"trial_until,omitempty"`
	UsageIncluded                          int64                  `json:"usage_included,omitempty" path:"usage_included,omitempty" url:"usage_included,omitempty"`
	UsersCount                             int64                  `json:"users_count,omitempty" path:"users_count,omitempty" url:"users_count,omitempty"`
}

func (s Settings) Identifier() interface{} {
	return s.Id
}

type SettingsCollection []Settings

func (s *Settings) UnmarshalJSON(data []byte) error {
	type settings Settings
	var v settings
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*s = Settings(v)
	return nil
}

func (s *SettingsCollection) UnmarshalJSON(data []byte) error {
	type settingss SettingsCollection
	var v settingss
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*s = SettingsCollection(v)
	return nil
}

func (s *SettingsCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*s))
	for i, v := range *s {
		ret[i] = v
	}

	return &ret
}
