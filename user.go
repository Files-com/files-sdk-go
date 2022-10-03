package files_sdk

import (
	"encoding/json"
	"io"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type User struct {
	Id                         int64      `json:"id,omitempty" path:"id"`
	Username                   string     `json:"username,omitempty" path:"username"`
	AdminGroupIds              []int64    `json:"admin_group_ids,omitempty" path:"admin_group_ids"`
	AllowedIps                 string     `json:"allowed_ips,omitempty" path:"allowed_ips"`
	AttachmentsPermission      *bool      `json:"attachments_permission,omitempty" path:"attachments_permission"`
	ApiKeysCount               int64      `json:"api_keys_count,omitempty" path:"api_keys_count"`
	AuthenticateUntil          *time.Time `json:"authenticate_until,omitempty" path:"authenticate_until"`
	AuthenticationMethod       string     `json:"authentication_method,omitempty" path:"authentication_method"`
	AvatarUrl                  string     `json:"avatar_url,omitempty" path:"avatar_url"`
	BillingPermission          *bool      `json:"billing_permission,omitempty" path:"billing_permission"`
	BypassSiteAllowedIps       *bool      `json:"bypass_site_allowed_ips,omitempty" path:"bypass_site_allowed_ips"`
	BypassInactiveDisable      *bool      `json:"bypass_inactive_disable,omitempty" path:"bypass_inactive_disable"`
	CreatedAt                  *time.Time `json:"created_at,omitempty" path:"created_at"`
	DavPermission              *bool      `json:"dav_permission,omitempty" path:"dav_permission"`
	Disabled                   *bool      `json:"disabled,omitempty" path:"disabled"`
	Email                      string     `json:"email,omitempty" path:"email"`
	FirstLoginAt               *time.Time `json:"first_login_at,omitempty" path:"first_login_at"`
	FtpPermission              *bool      `json:"ftp_permission,omitempty" path:"ftp_permission"`
	GroupIds                   string     `json:"group_ids,omitempty" path:"group_ids"`
	HeaderText                 string     `json:"header_text,omitempty" path:"header_text"`
	Language                   string     `json:"language,omitempty" path:"language"`
	LastLoginAt                *time.Time `json:"last_login_at,omitempty" path:"last_login_at"`
	LastWebLoginAt             *time.Time `json:"last_web_login_at,omitempty" path:"last_web_login_at"`
	LastFtpLoginAt             *time.Time `json:"last_ftp_login_at,omitempty" path:"last_ftp_login_at"`
	LastSftpLoginAt            *time.Time `json:"last_sftp_login_at,omitempty" path:"last_sftp_login_at"`
	LastDavLoginAt             *time.Time `json:"last_dav_login_at,omitempty" path:"last_dav_login_at"`
	LastDesktopLoginAt         *time.Time `json:"last_desktop_login_at,omitempty" path:"last_desktop_login_at"`
	LastRestapiLoginAt         *time.Time `json:"last_restapi_login_at,omitempty" path:"last_restapi_login_at"`
	LastApiUseAt               *time.Time `json:"last_api_use_at,omitempty" path:"last_api_use_at"`
	LastActiveAt               *time.Time `json:"last_active_at,omitempty" path:"last_active_at"`
	LastProtocolCipher         string     `json:"last_protocol_cipher,omitempty" path:"last_protocol_cipher"`
	LockoutExpires             *time.Time `json:"lockout_expires,omitempty" path:"lockout_expires"`
	Name                       string     `json:"name,omitempty" path:"name"`
	Company                    string     `json:"company,omitempty" path:"company"`
	Notes                      string     `json:"notes,omitempty" path:"notes"`
	NotificationDailySendTime  int64      `json:"notification_daily_send_time,omitempty" path:"notification_daily_send_time"`
	OfficeIntegrationEnabled   *bool      `json:"office_integration_enabled,omitempty" path:"office_integration_enabled"`
	PasswordSetAt              *time.Time `json:"password_set_at,omitempty" path:"password_set_at"`
	PasswordValidityDays       int64      `json:"password_validity_days,omitempty" path:"password_validity_days"`
	PublicKeysCount            int64      `json:"public_keys_count,omitempty" path:"public_keys_count"`
	ReceiveAdminAlerts         *bool      `json:"receive_admin_alerts,omitempty" path:"receive_admin_alerts"`
	Require2fa                 string     `json:"require_2fa,omitempty" path:"require_2fa"`
	Active2fa                  *bool      `json:"active_2fa,omitempty" path:"active_2fa"`
	RequirePasswordChange      *bool      `json:"require_password_change,omitempty" path:"require_password_change"`
	PasswordExpired            *bool      `json:"password_expired,omitempty" path:"password_expired"`
	RestapiPermission          *bool      `json:"restapi_permission,omitempty" path:"restapi_permission"`
	SelfManaged                *bool      `json:"self_managed,omitempty" path:"self_managed"`
	SftpPermission             *bool      `json:"sftp_permission,omitempty" path:"sftp_permission"`
	SiteAdmin                  *bool      `json:"site_admin,omitempty" path:"site_admin"`
	SkipWelcomeScreen          *bool      `json:"skip_welcome_screen,omitempty" path:"skip_welcome_screen"`
	SslRequired                string     `json:"ssl_required,omitempty" path:"ssl_required"`
	SsoStrategyId              int64      `json:"sso_strategy_id,omitempty" path:"sso_strategy_id"`
	SubscribeToNewsletter      *bool      `json:"subscribe_to_newsletter,omitempty" path:"subscribe_to_newsletter"`
	ExternallyManaged          *bool      `json:"externally_managed,omitempty" path:"externally_managed"`
	TimeZone                   string     `json:"time_zone,omitempty" path:"time_zone"`
	TypeOf2fa                  string     `json:"type_of_2fa,omitempty" path:"type_of_2fa"`
	UpdatedAt                  *time.Time `json:"updated_at,omitempty" path:"updated_at"`
	UserRoot                   string     `json:"user_root,omitempty" path:"user_root"`
	AvatarFile                 io.Reader  `json:"avatar_file,omitempty" path:"avatar_file"`
	AvatarDelete               *bool      `json:"avatar_delete,omitempty" path:"avatar_delete"`
	ChangePassword             string     `json:"change_password,omitempty" path:"change_password"`
	ChangePasswordConfirmation string     `json:"change_password_confirmation,omitempty" path:"change_password_confirmation"`
	GrantPermission            string     `json:"grant_permission,omitempty" path:"grant_permission"`
	GroupId                    int64      `json:"group_id,omitempty" path:"group_id"`
	ImportedPasswordHash       string     `json:"imported_password_hash,omitempty" path:"imported_password_hash"`
	Password                   string     `json:"password,omitempty" path:"password"`
	PasswordConfirmation       string     `json:"password_confirmation,omitempty" path:"password_confirmation"`
	AnnouncementsRead          *bool      `json:"announcements_read,omitempty" path:"announcements_read"`
}

type UserCollection []User

type QParam struct {
	Username             string `url:"username,omitempty" json:"username,omitempty" path:"username"`
	Email                string `url:"email,omitempty" json:"email,omitempty" path:"email"`
	Notes                string `url:"notes,omitempty" json:"notes,omitempty" path:"notes"`
	Admin                string `url:"admin,omitempty" json:"admin,omitempty" path:"admin"`
	AllowedIps           string `url:"allowed_ips,omitempty" json:"allowed_ips,omitempty" path:"allowed_ips"`
	PasswordValidityDays string `url:"password_validity_days,omitempty" json:"password_validity_days,omitempty" path:"password_validity_days"`
	SslRequired          string `url:"ssl_required,omitempty" json:"ssl_required,omitempty" path:"ssl_required"`
}

type UserAuthenticationMethodEnum string

func (u UserAuthenticationMethodEnum) String() string {
	return string(u)
}

func (u UserAuthenticationMethodEnum) Enum() map[string]UserAuthenticationMethodEnum {
	return map[string]UserAuthenticationMethodEnum{
		"password":                    UserAuthenticationMethodEnum("password"),
		"unused_former_ldap":          UserAuthenticationMethodEnum("unused_former_ldap"),
		"sso":                         UserAuthenticationMethodEnum("sso"),
		"none":                        UserAuthenticationMethodEnum("none"),
		"email_signup":                UserAuthenticationMethodEnum("email_signup"),
		"password_with_imported_hash": UserAuthenticationMethodEnum("password_with_imported_hash"),
	}
}

type UserSslRequiredEnum string

func (u UserSslRequiredEnum) String() string {
	return string(u)
}

func (u UserSslRequiredEnum) Enum() map[string]UserSslRequiredEnum {
	return map[string]UserSslRequiredEnum{
		"use_system_setting": UserSslRequiredEnum("use_system_setting"),
		"always_require":     UserSslRequiredEnum("always_require"),
		"never_require":      UserSslRequiredEnum("never_require"),
	}
}

type UserRequire2faEnum string

func (u UserRequire2faEnum) String() string {
	return string(u)
}

func (u UserRequire2faEnum) Enum() map[string]UserRequire2faEnum {
	return map[string]UserRequire2faEnum{
		"use_system_setting": UserRequire2faEnum("use_system_setting"),
		"always_require":     UserRequire2faEnum("always_require"),
		"never_require":      UserRequire2faEnum("never_require"),
	}
}

type UserListParams struct {
	SortBy     json.RawMessage `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty" path:"sort_by"`
	Filter     json.RawMessage `url:"filter,omitempty" required:"false" json:"filter,omitempty" path:"filter"`
	FilterGt   json.RawMessage `url:"filter_gt,omitempty" required:"false" json:"filter_gt,omitempty" path:"filter_gt"`
	FilterGteq json.RawMessage `url:"filter_gteq,omitempty" required:"false" json:"filter_gteq,omitempty" path:"filter_gteq"`
	FilterLike json.RawMessage `url:"filter_like,omitempty" required:"false" json:"filter_like,omitempty" path:"filter_like"`
	FilterLt   json.RawMessage `url:"filter_lt,omitempty" required:"false" json:"filter_lt,omitempty" path:"filter_lt"`
	FilterLteq json.RawMessage `url:"filter_lteq,omitempty" required:"false" json:"filter_lteq,omitempty" path:"filter_lteq"`
	Ids        string          `url:"ids,omitempty" required:"false" json:"ids,omitempty" path:"ids"`
	QParam     QParam          `url:"q,omitempty" required:"false" json:"q,omitempty" path:"q"`
	Search     string          `url:"search,omitempty" required:"false" json:"search,omitempty" path:"search"`
	lib.ListParams
}

type UserFindParams struct {
	Id int64 `url:"-,omitempty" required:"true" json:"-,omitempty" path:"id"`
}

type UserCreateParams struct {
	AvatarFile                 io.Writer                    `url:"avatar_file,omitempty" required:"false" json:"avatar_file,omitempty" path:"avatar_file"`
	AvatarDelete               *bool                        `url:"avatar_delete,omitempty" required:"false" json:"avatar_delete,omitempty" path:"avatar_delete"`
	ChangePassword             string                       `url:"change_password,omitempty" required:"false" json:"change_password,omitempty" path:"change_password"`
	ChangePasswordConfirmation string                       `url:"change_password_confirmation,omitempty" required:"false" json:"change_password_confirmation,omitempty" path:"change_password_confirmation"`
	Email                      string                       `url:"email,omitempty" required:"false" json:"email,omitempty" path:"email"`
	GrantPermission            string                       `url:"grant_permission,omitempty" required:"false" json:"grant_permission,omitempty" path:"grant_permission"`
	GroupId                    int64                        `url:"group_id,omitempty" required:"false" json:"group_id,omitempty" path:"group_id"`
	GroupIds                   string                       `url:"group_ids,omitempty" required:"false" json:"group_ids,omitempty" path:"group_ids"`
	ImportedPasswordHash       string                       `url:"imported_password_hash,omitempty" required:"false" json:"imported_password_hash,omitempty" path:"imported_password_hash"`
	Password                   string                       `url:"password,omitempty" required:"false" json:"password,omitempty" path:"password"`
	PasswordConfirmation       string                       `url:"password_confirmation,omitempty" required:"false" json:"password_confirmation,omitempty" path:"password_confirmation"`
	AnnouncementsRead          *bool                        `url:"announcements_read,omitempty" required:"false" json:"announcements_read,omitempty" path:"announcements_read"`
	AllowedIps                 string                       `url:"allowed_ips,omitempty" required:"false" json:"allowed_ips,omitempty" path:"allowed_ips"`
	AttachmentsPermission      *bool                        `url:"attachments_permission,omitempty" required:"false" json:"attachments_permission,omitempty" path:"attachments_permission"`
	AuthenticateUntil          *time.Time                   `url:"authenticate_until,omitempty" required:"false" json:"authenticate_until,omitempty" path:"authenticate_until"`
	AuthenticationMethod       UserAuthenticationMethodEnum `url:"authentication_method,omitempty" required:"false" json:"authentication_method,omitempty" path:"authentication_method"`
	BillingPermission          *bool                        `url:"billing_permission,omitempty" required:"false" json:"billing_permission,omitempty" path:"billing_permission"`
	BypassInactiveDisable      *bool                        `url:"bypass_inactive_disable,omitempty" required:"false" json:"bypass_inactive_disable,omitempty" path:"bypass_inactive_disable"`
	BypassSiteAllowedIps       *bool                        `url:"bypass_site_allowed_ips,omitempty" required:"false" json:"bypass_site_allowed_ips,omitempty" path:"bypass_site_allowed_ips"`
	DavPermission              *bool                        `url:"dav_permission,omitempty" required:"false" json:"dav_permission,omitempty" path:"dav_permission"`
	Disabled                   *bool                        `url:"disabled,omitempty" required:"false" json:"disabled,omitempty" path:"disabled"`
	FtpPermission              *bool                        `url:"ftp_permission,omitempty" required:"false" json:"ftp_permission,omitempty" path:"ftp_permission"`
	HeaderText                 string                       `url:"header_text,omitempty" required:"false" json:"header_text,omitempty" path:"header_text"`
	Language                   string                       `url:"language,omitempty" required:"false" json:"language,omitempty" path:"language"`
	NotificationDailySendTime  int64                        `url:"notification_daily_send_time,omitempty" required:"false" json:"notification_daily_send_time,omitempty" path:"notification_daily_send_time"`
	Name                       string                       `url:"name,omitempty" required:"false" json:"name,omitempty" path:"name"`
	Company                    string                       `url:"company,omitempty" required:"false" json:"company,omitempty" path:"company"`
	Notes                      string                       `url:"notes,omitempty" required:"false" json:"notes,omitempty" path:"notes"`
	OfficeIntegrationEnabled   *bool                        `url:"office_integration_enabled,omitempty" required:"false" json:"office_integration_enabled,omitempty" path:"office_integration_enabled"`
	PasswordValidityDays       int64                        `url:"password_validity_days,omitempty" required:"false" json:"password_validity_days,omitempty" path:"password_validity_days"`
	ReceiveAdminAlerts         *bool                        `url:"receive_admin_alerts,omitempty" required:"false" json:"receive_admin_alerts,omitempty" path:"receive_admin_alerts"`
	RequirePasswordChange      *bool                        `url:"require_password_change,omitempty" required:"false" json:"require_password_change,omitempty" path:"require_password_change"`
	RestapiPermission          *bool                        `url:"restapi_permission,omitempty" required:"false" json:"restapi_permission,omitempty" path:"restapi_permission"`
	SelfManaged                *bool                        `url:"self_managed,omitempty" required:"false" json:"self_managed,omitempty" path:"self_managed"`
	SftpPermission             *bool                        `url:"sftp_permission,omitempty" required:"false" json:"sftp_permission,omitempty" path:"sftp_permission"`
	SiteAdmin                  *bool                        `url:"site_admin,omitempty" required:"false" json:"site_admin,omitempty" path:"site_admin"`
	SkipWelcomeScreen          *bool                        `url:"skip_welcome_screen,omitempty" required:"false" json:"skip_welcome_screen,omitempty" path:"skip_welcome_screen"`
	SslRequired                UserSslRequiredEnum          `url:"ssl_required,omitempty" required:"false" json:"ssl_required,omitempty" path:"ssl_required"`
	SsoStrategyId              int64                        `url:"sso_strategy_id,omitempty" required:"false" json:"sso_strategy_id,omitempty" path:"sso_strategy_id"`
	SubscribeToNewsletter      *bool                        `url:"subscribe_to_newsletter,omitempty" required:"false" json:"subscribe_to_newsletter,omitempty" path:"subscribe_to_newsletter"`
	Require2fa                 UserRequire2faEnum           `url:"require_2fa,omitempty" required:"false" json:"require_2fa,omitempty" path:"require_2fa"`
	TimeZone                   string                       `url:"time_zone,omitempty" required:"false" json:"time_zone,omitempty" path:"time_zone"`
	UserRoot                   string                       `url:"user_root,omitempty" required:"false" json:"user_root,omitempty" path:"user_root"`
	Username                   string                       `url:"username,omitempty" required:"false" json:"username,omitempty" path:"username"`
}

// Unlock user who has been locked out due to failed logins
type UserUnlockParams struct {
	Id int64 `url:"-,omitempty" required:"true" json:"-,omitempty" path:"id"`
}

// Resend user welcome email
type UserResendWelcomeEmailParams struct {
	Id int64 `url:"-,omitempty" required:"true" json:"-,omitempty" path:"id"`
}

// Trigger 2FA Reset process for user who has lost access to their existing 2FA methods
type UserUser2faResetParams struct {
	Id int64 `url:"-,omitempty" required:"true" json:"-,omitempty" path:"id"`
}

type UserUpdateParams struct {
	Id                         int64                        `url:"-,omitempty" required:"true" json:"-,omitempty" path:"id"`
	AvatarFile                 io.Writer                    `url:"avatar_file,omitempty" required:"false" json:"avatar_file,omitempty" path:"avatar_file"`
	AvatarDelete               *bool                        `url:"avatar_delete,omitempty" required:"false" json:"avatar_delete,omitempty" path:"avatar_delete"`
	ChangePassword             string                       `url:"change_password,omitempty" required:"false" json:"change_password,omitempty" path:"change_password"`
	ChangePasswordConfirmation string                       `url:"change_password_confirmation,omitempty" required:"false" json:"change_password_confirmation,omitempty" path:"change_password_confirmation"`
	Email                      string                       `url:"email,omitempty" required:"false" json:"email,omitempty" path:"email"`
	GrantPermission            string                       `url:"grant_permission,omitempty" required:"false" json:"grant_permission,omitempty" path:"grant_permission"`
	GroupId                    int64                        `url:"group_id,omitempty" required:"false" json:"group_id,omitempty" path:"group_id"`
	GroupIds                   string                       `url:"group_ids,omitempty" required:"false" json:"group_ids,omitempty" path:"group_ids"`
	ImportedPasswordHash       string                       `url:"imported_password_hash,omitempty" required:"false" json:"imported_password_hash,omitempty" path:"imported_password_hash"`
	Password                   string                       `url:"password,omitempty" required:"false" json:"password,omitempty" path:"password"`
	PasswordConfirmation       string                       `url:"password_confirmation,omitempty" required:"false" json:"password_confirmation,omitempty" path:"password_confirmation"`
	AnnouncementsRead          *bool                        `url:"announcements_read,omitempty" required:"false" json:"announcements_read,omitempty" path:"announcements_read"`
	AllowedIps                 string                       `url:"allowed_ips,omitempty" required:"false" json:"allowed_ips,omitempty" path:"allowed_ips"`
	AttachmentsPermission      *bool                        `url:"attachments_permission,omitempty" required:"false" json:"attachments_permission,omitempty" path:"attachments_permission"`
	AuthenticateUntil          *time.Time                   `url:"authenticate_until,omitempty" required:"false" json:"authenticate_until,omitempty" path:"authenticate_until"`
	AuthenticationMethod       UserAuthenticationMethodEnum `url:"authentication_method,omitempty" required:"false" json:"authentication_method,omitempty" path:"authentication_method"`
	BillingPermission          *bool                        `url:"billing_permission,omitempty" required:"false" json:"billing_permission,omitempty" path:"billing_permission"`
	BypassInactiveDisable      *bool                        `url:"bypass_inactive_disable,omitempty" required:"false" json:"bypass_inactive_disable,omitempty" path:"bypass_inactive_disable"`
	BypassSiteAllowedIps       *bool                        `url:"bypass_site_allowed_ips,omitempty" required:"false" json:"bypass_site_allowed_ips,omitempty" path:"bypass_site_allowed_ips"`
	DavPermission              *bool                        `url:"dav_permission,omitempty" required:"false" json:"dav_permission,omitempty" path:"dav_permission"`
	Disabled                   *bool                        `url:"disabled,omitempty" required:"false" json:"disabled,omitempty" path:"disabled"`
	FtpPermission              *bool                        `url:"ftp_permission,omitempty" required:"false" json:"ftp_permission,omitempty" path:"ftp_permission"`
	HeaderText                 string                       `url:"header_text,omitempty" required:"false" json:"header_text,omitempty" path:"header_text"`
	Language                   string                       `url:"language,omitempty" required:"false" json:"language,omitempty" path:"language"`
	NotificationDailySendTime  int64                        `url:"notification_daily_send_time,omitempty" required:"false" json:"notification_daily_send_time,omitempty" path:"notification_daily_send_time"`
	Name                       string                       `url:"name,omitempty" required:"false" json:"name,omitempty" path:"name"`
	Company                    string                       `url:"company,omitempty" required:"false" json:"company,omitempty" path:"company"`
	Notes                      string                       `url:"notes,omitempty" required:"false" json:"notes,omitempty" path:"notes"`
	OfficeIntegrationEnabled   *bool                        `url:"office_integration_enabled,omitempty" required:"false" json:"office_integration_enabled,omitempty" path:"office_integration_enabled"`
	PasswordValidityDays       int64                        `url:"password_validity_days,omitempty" required:"false" json:"password_validity_days,omitempty" path:"password_validity_days"`
	ReceiveAdminAlerts         *bool                        `url:"receive_admin_alerts,omitempty" required:"false" json:"receive_admin_alerts,omitempty" path:"receive_admin_alerts"`
	RequirePasswordChange      *bool                        `url:"require_password_change,omitempty" required:"false" json:"require_password_change,omitempty" path:"require_password_change"`
	RestapiPermission          *bool                        `url:"restapi_permission,omitempty" required:"false" json:"restapi_permission,omitempty" path:"restapi_permission"`
	SelfManaged                *bool                        `url:"self_managed,omitempty" required:"false" json:"self_managed,omitempty" path:"self_managed"`
	SftpPermission             *bool                        `url:"sftp_permission,omitempty" required:"false" json:"sftp_permission,omitempty" path:"sftp_permission"`
	SiteAdmin                  *bool                        `url:"site_admin,omitempty" required:"false" json:"site_admin,omitempty" path:"site_admin"`
	SkipWelcomeScreen          *bool                        `url:"skip_welcome_screen,omitempty" required:"false" json:"skip_welcome_screen,omitempty" path:"skip_welcome_screen"`
	SslRequired                UserSslRequiredEnum          `url:"ssl_required,omitempty" required:"false" json:"ssl_required,omitempty" path:"ssl_required"`
	SsoStrategyId              int64                        `url:"sso_strategy_id,omitempty" required:"false" json:"sso_strategy_id,omitempty" path:"sso_strategy_id"`
	SubscribeToNewsletter      *bool                        `url:"subscribe_to_newsletter,omitempty" required:"false" json:"subscribe_to_newsletter,omitempty" path:"subscribe_to_newsletter"`
	Require2fa                 UserRequire2faEnum           `url:"require_2fa,omitempty" required:"false" json:"require_2fa,omitempty" path:"require_2fa"`
	TimeZone                   string                       `url:"time_zone,omitempty" required:"false" json:"time_zone,omitempty" path:"time_zone"`
	UserRoot                   string                       `url:"user_root,omitempty" required:"false" json:"user_root,omitempty" path:"user_root"`
	Username                   string                       `url:"username,omitempty" required:"false" json:"username,omitempty" path:"username"`
}

type UserDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"true" json:"-,omitempty" path:"id"`
}

func (u *User) UnmarshalJSON(data []byte) error {
	type user User
	var v user
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*u = User(v)
	return nil
}

func (u *UserCollection) UnmarshalJSON(data []byte) error {
	type users UserCollection
	var v users
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*u = UserCollection(v)
	return nil
}

func (u *UserCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*u))
	for i, v := range *u {
		ret[i] = v
	}

	return &ret
}
