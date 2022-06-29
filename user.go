package files_sdk

import (
	"encoding/json"
	"io"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type User struct {
	Id                         int64      `json:"id,omitempty"`
	Username                   string     `json:"username,omitempty"`
	AdminGroupIds              []int64    `json:"admin_group_ids,omitempty"`
	AllowedIps                 string     `json:"allowed_ips,omitempty"`
	AttachmentsPermission      *bool      `json:"attachments_permission,omitempty"`
	ApiKeysCount               int64      `json:"api_keys_count,omitempty"`
	AuthenticateUntil          *time.Time `json:"authenticate_until,omitempty"`
	AuthenticationMethod       string     `json:"authentication_method,omitempty"`
	AvatarUrl                  string     `json:"avatar_url,omitempty"`
	BillingPermission          *bool      `json:"billing_permission,omitempty"`
	BypassSiteAllowedIps       *bool      `json:"bypass_site_allowed_ips,omitempty"`
	BypassInactiveDisable      *bool      `json:"bypass_inactive_disable,omitempty"`
	CreatedAt                  *time.Time `json:"created_at,omitempty"`
	DavPermission              *bool      `json:"dav_permission,omitempty"`
	Disabled                   *bool      `json:"disabled,omitempty"`
	Email                      string     `json:"email,omitempty"`
	FtpPermission              *bool      `json:"ftp_permission,omitempty"`
	GroupIds                   string     `json:"group_ids,omitempty"`
	HeaderText                 string     `json:"header_text,omitempty"`
	Language                   string     `json:"language,omitempty"`
	LastLoginAt                *time.Time `json:"last_login_at,omitempty"`
	LastProtocolCipher         string     `json:"last_protocol_cipher,omitempty"`
	LockoutExpires             *time.Time `json:"lockout_expires,omitempty"`
	Name                       string     `json:"name,omitempty"`
	Company                    string     `json:"company,omitempty"`
	Notes                      string     `json:"notes,omitempty"`
	NotificationDailySendTime  int64      `json:"notification_daily_send_time,omitempty"`
	OfficeIntegrationEnabled   *bool      `json:"office_integration_enabled,omitempty"`
	PasswordSetAt              *time.Time `json:"password_set_at,omitempty"`
	PasswordValidityDays       int64      `json:"password_validity_days,omitempty"`
	PublicKeysCount            int64      `json:"public_keys_count,omitempty"`
	ReceiveAdminAlerts         *bool      `json:"receive_admin_alerts,omitempty"`
	Require2fa                 string     `json:"require_2fa,omitempty"`
	Active2fa                  *bool      `json:"active_2fa,omitempty"`
	RequirePasswordChange      *bool      `json:"require_password_change,omitempty"`
	PasswordExpired            *bool      `json:"password_expired,omitempty"`
	RestapiPermission          *bool      `json:"restapi_permission,omitempty"`
	SelfManaged                *bool      `json:"self_managed,omitempty"`
	SftpPermission             *bool      `json:"sftp_permission,omitempty"`
	SiteAdmin                  *bool      `json:"site_admin,omitempty"`
	SkipWelcomeScreen          *bool      `json:"skip_welcome_screen,omitempty"`
	SslRequired                string     `json:"ssl_required,omitempty"`
	SsoStrategyId              int64      `json:"sso_strategy_id,omitempty"`
	SubscribeToNewsletter      *bool      `json:"subscribe_to_newsletter,omitempty"`
	ExternallyManaged          *bool      `json:"externally_managed,omitempty"`
	TimeZone                   string     `json:"time_zone,omitempty"`
	TypeOf2fa                  string     `json:"type_of_2fa,omitempty"`
	UserRoot                   string     `json:"user_root,omitempty"`
	AvatarFile                 io.Reader  `json:"avatar_file,omitempty"`
	AvatarDelete               *bool      `json:"avatar_delete,omitempty"`
	ChangePassword             string     `json:"change_password,omitempty"`
	ChangePasswordConfirmation string     `json:"change_password_confirmation,omitempty"`
	GrantPermission            string     `json:"grant_permission,omitempty"`
	GroupId                    int64      `json:"group_id,omitempty"`
	ImportedPasswordHash       string     `json:"imported_password_hash,omitempty"`
	Password                   string     `json:"password,omitempty"`
	PasswordConfirmation       string     `json:"password_confirmation,omitempty"`
	AnnouncementsRead          *bool      `json:"announcements_read,omitempty"`
}

type UserCollection []User

type QParam struct {
	Username             string `url:"username,omitempty" json:"username,omitempty"`
	Email                string `url:"email,omitempty" json:"email,omitempty"`
	Notes                string `url:"notes,omitempty" json:"notes,omitempty"`
	Admin                string `url:"admin,omitempty" json:"admin,omitempty"`
	AllowedIps           string `url:"allowed_ips,omitempty" json:"allowed_ips,omitempty"`
	PasswordValidityDays string `url:"password_validity_days,omitempty" json:"password_validity_days,omitempty"`
	SslRequired          string `url:"ssl_required,omitempty" json:"ssl_required,omitempty"`
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
	Cursor     string          `url:"cursor,omitempty" required:"false" json:"cursor,omitempty"`
	PerPage    int64           `url:"per_page,omitempty" required:"false" json:"per_page,omitempty"`
	SortBy     json.RawMessage `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty"`
	Filter     json.RawMessage `url:"filter,omitempty" required:"false" json:"filter,omitempty"`
	FilterGt   json.RawMessage `url:"filter_gt,omitempty" required:"false" json:"filter_gt,omitempty"`
	FilterGteq json.RawMessage `url:"filter_gteq,omitempty" required:"false" json:"filter_gteq,omitempty"`
	FilterLike json.RawMessage `url:"filter_like,omitempty" required:"false" json:"filter_like,omitempty"`
	FilterLt   json.RawMessage `url:"filter_lt,omitempty" required:"false" json:"filter_lt,omitempty"`
	FilterLteq json.RawMessage `url:"filter_lteq,omitempty" required:"false" json:"filter_lteq,omitempty"`
	Ids        string          `url:"ids,omitempty" required:"false" json:"ids,omitempty"`
	QParam     QParam          `url:"q,omitempty" required:"false" json:"q,omitempty"`
	Search     string          `url:"search,omitempty" required:"false" json:"search,omitempty"`
	lib.ListParams
}

type UserFindParams struct {
	Id int64 `url:"-,omitempty" required:"true" json:"-,omitempty"`
}

type UserCreateParams struct {
	AvatarFile                 io.Writer                    `url:"avatar_file,omitempty" required:"false" json:"avatar_file,omitempty"`
	AvatarDelete               *bool                        `url:"avatar_delete,omitempty" required:"false" json:"avatar_delete,omitempty"`
	ChangePassword             string                       `url:"change_password,omitempty" required:"false" json:"change_password,omitempty"`
	ChangePasswordConfirmation string                       `url:"change_password_confirmation,omitempty" required:"false" json:"change_password_confirmation,omitempty"`
	Email                      string                       `url:"email,omitempty" required:"false" json:"email,omitempty"`
	GrantPermission            string                       `url:"grant_permission,omitempty" required:"false" json:"grant_permission,omitempty"`
	GroupId                    int64                        `url:"group_id,omitempty" required:"false" json:"group_id,omitempty"`
	GroupIds                   string                       `url:"group_ids,omitempty" required:"false" json:"group_ids,omitempty"`
	ImportedPasswordHash       string                       `url:"imported_password_hash,omitempty" required:"false" json:"imported_password_hash,omitempty"`
	Password                   string                       `url:"password,omitempty" required:"false" json:"password,omitempty"`
	PasswordConfirmation       string                       `url:"password_confirmation,omitempty" required:"false" json:"password_confirmation,omitempty"`
	AnnouncementsRead          *bool                        `url:"announcements_read,omitempty" required:"false" json:"announcements_read,omitempty"`
	AllowedIps                 string                       `url:"allowed_ips,omitempty" required:"false" json:"allowed_ips,omitempty"`
	AttachmentsPermission      *bool                        `url:"attachments_permission,omitempty" required:"false" json:"attachments_permission,omitempty"`
	AuthenticateUntil          *time.Time                   `url:"authenticate_until,omitempty" required:"false" json:"authenticate_until,omitempty"`
	AuthenticationMethod       UserAuthenticationMethodEnum `url:"authentication_method,omitempty" required:"false" json:"authentication_method,omitempty"`
	BillingPermission          *bool                        `url:"billing_permission,omitempty" required:"false" json:"billing_permission,omitempty"`
	BypassInactiveDisable      *bool                        `url:"bypass_inactive_disable,omitempty" required:"false" json:"bypass_inactive_disable,omitempty"`
	BypassSiteAllowedIps       *bool                        `url:"bypass_site_allowed_ips,omitempty" required:"false" json:"bypass_site_allowed_ips,omitempty"`
	DavPermission              *bool                        `url:"dav_permission,omitempty" required:"false" json:"dav_permission,omitempty"`
	Disabled                   *bool                        `url:"disabled,omitempty" required:"false" json:"disabled,omitempty"`
	FtpPermission              *bool                        `url:"ftp_permission,omitempty" required:"false" json:"ftp_permission,omitempty"`
	HeaderText                 string                       `url:"header_text,omitempty" required:"false" json:"header_text,omitempty"`
	Language                   string                       `url:"language,omitempty" required:"false" json:"language,omitempty"`
	NotificationDailySendTime  int64                        `url:"notification_daily_send_time,omitempty" required:"false" json:"notification_daily_send_time,omitempty"`
	Name                       string                       `url:"name,omitempty" required:"false" json:"name,omitempty"`
	Company                    string                       `url:"company,omitempty" required:"false" json:"company,omitempty"`
	Notes                      string                       `url:"notes,omitempty" required:"false" json:"notes,omitempty"`
	OfficeIntegrationEnabled   *bool                        `url:"office_integration_enabled,omitempty" required:"false" json:"office_integration_enabled,omitempty"`
	PasswordValidityDays       int64                        `url:"password_validity_days,omitempty" required:"false" json:"password_validity_days,omitempty"`
	ReceiveAdminAlerts         *bool                        `url:"receive_admin_alerts,omitempty" required:"false" json:"receive_admin_alerts,omitempty"`
	RequirePasswordChange      *bool                        `url:"require_password_change,omitempty" required:"false" json:"require_password_change,omitempty"`
	RestapiPermission          *bool                        `url:"restapi_permission,omitempty" required:"false" json:"restapi_permission,omitempty"`
	SelfManaged                *bool                        `url:"self_managed,omitempty" required:"false" json:"self_managed,omitempty"`
	SftpPermission             *bool                        `url:"sftp_permission,omitempty" required:"false" json:"sftp_permission,omitempty"`
	SiteAdmin                  *bool                        `url:"site_admin,omitempty" required:"false" json:"site_admin,omitempty"`
	SkipWelcomeScreen          *bool                        `url:"skip_welcome_screen,omitempty" required:"false" json:"skip_welcome_screen,omitempty"`
	SslRequired                UserSslRequiredEnum          `url:"ssl_required,omitempty" required:"false" json:"ssl_required,omitempty"`
	SsoStrategyId              int64                        `url:"sso_strategy_id,omitempty" required:"false" json:"sso_strategy_id,omitempty"`
	SubscribeToNewsletter      *bool                        `url:"subscribe_to_newsletter,omitempty" required:"false" json:"subscribe_to_newsletter,omitempty"`
	Require2fa                 UserRequire2faEnum           `url:"require_2fa,omitempty" required:"false" json:"require_2fa,omitempty"`
	TimeZone                   string                       `url:"time_zone,omitempty" required:"false" json:"time_zone,omitempty"`
	UserRoot                   string                       `url:"user_root,omitempty" required:"false" json:"user_root,omitempty"`
	Username                   string                       `url:"username,omitempty" required:"false" json:"username,omitempty"`
}

// Unlock user who has been locked out due to failed logins
type UserUnlockParams struct {
	Id int64 `url:"-,omitempty" required:"true" json:"-,omitempty"`
}

// Resend user welcome email
type UserResendWelcomeEmailParams struct {
	Id int64 `url:"-,omitempty" required:"true" json:"-,omitempty"`
}

// Trigger 2FA Reset process for user who has lost access to their existing 2FA methods
type UserUser2faResetParams struct {
	Id int64 `url:"-,omitempty" required:"true" json:"-,omitempty"`
}

type UserUpdateParams struct {
	Id                         int64                        `url:"-,omitempty" required:"true" json:"-,omitempty"`
	AvatarFile                 io.Writer                    `url:"avatar_file,omitempty" required:"false" json:"avatar_file,omitempty"`
	AvatarDelete               *bool                        `url:"avatar_delete,omitempty" required:"false" json:"avatar_delete,omitempty"`
	ChangePassword             string                       `url:"change_password,omitempty" required:"false" json:"change_password,omitempty"`
	ChangePasswordConfirmation string                       `url:"change_password_confirmation,omitempty" required:"false" json:"change_password_confirmation,omitempty"`
	Email                      string                       `url:"email,omitempty" required:"false" json:"email,omitempty"`
	GrantPermission            string                       `url:"grant_permission,omitempty" required:"false" json:"grant_permission,omitempty"`
	GroupId                    int64                        `url:"group_id,omitempty" required:"false" json:"group_id,omitempty"`
	GroupIds                   string                       `url:"group_ids,omitempty" required:"false" json:"group_ids,omitempty"`
	ImportedPasswordHash       string                       `url:"imported_password_hash,omitempty" required:"false" json:"imported_password_hash,omitempty"`
	Password                   string                       `url:"password,omitempty" required:"false" json:"password,omitempty"`
	PasswordConfirmation       string                       `url:"password_confirmation,omitempty" required:"false" json:"password_confirmation,omitempty"`
	AnnouncementsRead          *bool                        `url:"announcements_read,omitempty" required:"false" json:"announcements_read,omitempty"`
	AllowedIps                 string                       `url:"allowed_ips,omitempty" required:"false" json:"allowed_ips,omitempty"`
	AttachmentsPermission      *bool                        `url:"attachments_permission,omitempty" required:"false" json:"attachments_permission,omitempty"`
	AuthenticateUntil          *time.Time                   `url:"authenticate_until,omitempty" required:"false" json:"authenticate_until,omitempty"`
	AuthenticationMethod       UserAuthenticationMethodEnum `url:"authentication_method,omitempty" required:"false" json:"authentication_method,omitempty"`
	BillingPermission          *bool                        `url:"billing_permission,omitempty" required:"false" json:"billing_permission,omitempty"`
	BypassInactiveDisable      *bool                        `url:"bypass_inactive_disable,omitempty" required:"false" json:"bypass_inactive_disable,omitempty"`
	BypassSiteAllowedIps       *bool                        `url:"bypass_site_allowed_ips,omitempty" required:"false" json:"bypass_site_allowed_ips,omitempty"`
	DavPermission              *bool                        `url:"dav_permission,omitempty" required:"false" json:"dav_permission,omitempty"`
	Disabled                   *bool                        `url:"disabled,omitempty" required:"false" json:"disabled,omitempty"`
	FtpPermission              *bool                        `url:"ftp_permission,omitempty" required:"false" json:"ftp_permission,omitempty"`
	HeaderText                 string                       `url:"header_text,omitempty" required:"false" json:"header_text,omitempty"`
	Language                   string                       `url:"language,omitempty" required:"false" json:"language,omitempty"`
	NotificationDailySendTime  int64                        `url:"notification_daily_send_time,omitempty" required:"false" json:"notification_daily_send_time,omitempty"`
	Name                       string                       `url:"name,omitempty" required:"false" json:"name,omitempty"`
	Company                    string                       `url:"company,omitempty" required:"false" json:"company,omitempty"`
	Notes                      string                       `url:"notes,omitempty" required:"false" json:"notes,omitempty"`
	OfficeIntegrationEnabled   *bool                        `url:"office_integration_enabled,omitempty" required:"false" json:"office_integration_enabled,omitempty"`
	PasswordValidityDays       int64                        `url:"password_validity_days,omitempty" required:"false" json:"password_validity_days,omitempty"`
	ReceiveAdminAlerts         *bool                        `url:"receive_admin_alerts,omitempty" required:"false" json:"receive_admin_alerts,omitempty"`
	RequirePasswordChange      *bool                        `url:"require_password_change,omitempty" required:"false" json:"require_password_change,omitempty"`
	RestapiPermission          *bool                        `url:"restapi_permission,omitempty" required:"false" json:"restapi_permission,omitempty"`
	SelfManaged                *bool                        `url:"self_managed,omitempty" required:"false" json:"self_managed,omitempty"`
	SftpPermission             *bool                        `url:"sftp_permission,omitempty" required:"false" json:"sftp_permission,omitempty"`
	SiteAdmin                  *bool                        `url:"site_admin,omitempty" required:"false" json:"site_admin,omitempty"`
	SkipWelcomeScreen          *bool                        `url:"skip_welcome_screen,omitempty" required:"false" json:"skip_welcome_screen,omitempty"`
	SslRequired                UserSslRequiredEnum          `url:"ssl_required,omitempty" required:"false" json:"ssl_required,omitempty"`
	SsoStrategyId              int64                        `url:"sso_strategy_id,omitempty" required:"false" json:"sso_strategy_id,omitempty"`
	SubscribeToNewsletter      *bool                        `url:"subscribe_to_newsletter,omitempty" required:"false" json:"subscribe_to_newsletter,omitempty"`
	Require2fa                 UserRequire2faEnum           `url:"require_2fa,omitempty" required:"false" json:"require_2fa,omitempty"`
	TimeZone                   string                       `url:"time_zone,omitempty" required:"false" json:"time_zone,omitempty"`
	UserRoot                   string                       `url:"user_root,omitempty" required:"false" json:"user_root,omitempty"`
	Username                   string                       `url:"username,omitempty" required:"false" json:"username,omitempty"`
}

type UserDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"true" json:"-,omitempty"`
}

func (u *User) UnmarshalJSON(data []byte) error {
	type user User
	var v user
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*u = User(v)
	return nil
}

func (u *UserCollection) UnmarshalJSON(data []byte) error {
	type users []User
	var v users
	if err := json.Unmarshal(data, &v); err != nil {
		return err
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
