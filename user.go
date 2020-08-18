package files_sdk

import (
	"encoding/json"
	lib "github.com/Files-com/files-sdk-go/lib"
	"io"
	"time"
)

type User struct {
	Id                         int64     `json:"id,omitempty"`
	Username                   string    `json:"username,omitempty"`
	AdminGroupIds              []string  `json:"admin_group_ids,omitempty"`
	AllowedIps                 string    `json:"allowed_ips,omitempty"`
	AttachmentsPermission      *bool     `json:"attachments_permission,omitempty"`
	ApiKeysCount               int       `json:"api_keys_count,omitempty"`
	AuthenticateUntil          time.Time `json:"authenticate_until,omitempty"`
	AuthenticationMethod       string    `json:"authentication_method,omitempty"`
	AvatarUrl                  string    `json:"avatar_url,omitempty"`
	BillingPermission          *bool     `json:"billing_permission,omitempty"`
	BypassSiteAllowedIps       *bool     `json:"bypass_site_allowed_ips,omitempty"`
	BypassInactiveDisable      *bool     `json:"bypass_inactive_disable,omitempty"`
	CreatedAt                  time.Time `json:"created_at,omitempty"`
	DavPermission              *bool     `json:"dav_permission,omitempty"`
	Disabled                   *bool     `json:"disabled,omitempty"`
	Email                      string    `json:"email,omitempty"`
	FtpPermission              *bool     `json:"ftp_permission,omitempty"`
	GroupIds                   string    `json:"group_ids,omitempty"`
	HeaderText                 string    `json:"header_text,omitempty"`
	Language                   string    `json:"language,omitempty"`
	LastLoginAt                time.Time `json:"last_login_at,omitempty"`
	LastProtocolCipher         string    `json:"last_protocol_cipher,omitempty"`
	LockoutExpires             time.Time `json:"lockout_expires,omitempty"`
	Name                       string    `json:"name,omitempty"`
	Notes                      string    `json:"notes,omitempty"`
	NotificationDailySendTime  int       `json:"notification_daily_send_time,omitempty"`
	OfficeIntegrationEnabled   *bool     `json:"office_integration_enabled,omitempty"`
	PasswordSetAt              time.Time `json:"password_set_at,omitempty"`
	PasswordValidityDays       int       `json:"password_validity_days,omitempty"`
	PublicKeysCount            int       `json:"public_keys_count,omitempty"`
	ReceiveAdminAlerts         *bool     `json:"receive_admin_alerts,omitempty"`
	Require2fa                 *bool     `json:"require_2fa,omitempty"`
	RequirePasswordChange      *bool     `json:"require_password_change,omitempty"`
	RestapiPermission          *bool     `json:"restapi_permission,omitempty"`
	SelfManaged                *bool     `json:"self_managed,omitempty"`
	SftpPermission             *bool     `json:"sftp_permission,omitempty"`
	SiteAdmin                  *bool     `json:"site_admin,omitempty"`
	SkipWelcomeScreen          *bool     `json:"skip_welcome_screen,omitempty"`
	SslRequired                string    `json:"ssl_required,omitempty"`
	SsoStrategyId              int64     `json:"sso_strategy_id,omitempty"`
	SubscribeToNewsletter      *bool     `json:"subscribe_to_newsletter,omitempty"`
	ExternallyManaged          *bool     `json:"externally_managed,omitempty"`
	TimeZone                   string    `json:"time_zone,omitempty"`
	TypeOf2fa                  string    `json:"type_of_2fa,omitempty"`
	UserRoot                   string    `json:"user_root,omitempty"`
	AvatarFile                 io.Reader `json:"avatar_file,omitempty"`
	AvatarDelete               *bool     `json:"avatar_delete,omitempty"`
	ChangePassword             string    `json:"change_password,omitempty"`
	ChangePasswordConfirmation string    `json:"change_password_confirmation,omitempty"`
	GrantPermission            string    `json:"grant_permission,omitempty"`
	GroupId                    int64     `json:"group_id,omitempty"`
	Password                   string    `json:"password,omitempty"`
	PasswordConfirmation       string    `json:"password_confirmation,omitempty"`
	AnnouncementsRead          *bool     `json:"announcements_read,omitempty"`
}

type UserCollection []User

type QParam struct {
	Username             string `url:"username,omitempty"`
	Email                string `url:"email,omitempty"`
	Notes                string `url:"notes,omitempty"`
	Admin                string `url:"admin,omitempty"`
	AllowedIps           string `url:"allowed_ips,omitempty"`
	PasswordValidityDays string `url:"password_validity_days,omitempty"`
	SslRequired          string `url:"ssl_required,omitempty"`
}

type UserListParams struct {
	Page       int             `url:"page,omitempty"`
	PerPage    int             `url:"per_page,omitempty"`
	Action     string          `url:"action,omitempty"`
	Cursor     string          `url:"cursor,omitempty"`
	SortBy     json.RawMessage `url:"sort_by,omitempty"`
	Filter     json.RawMessage `url:"filter,omitempty"`
	FilterGt   json.RawMessage `url:"filter_gt,omitempty"`
	FilterGteq json.RawMessage `url:"filter_gteq,omitempty"`
	FilterLike json.RawMessage `url:"filter_like,omitempty"`
	FilterLt   json.RawMessage `url:"filter_lt,omitempty"`
	FilterLteq json.RawMessage `url:"filter_lteq,omitempty"`
	Ids        string          `url:"ids,omitempty"`
	QParam     QParam          `url:"q,omitempty"`
	Search     string          `url:"search,omitempty"`
	lib.ListParams
}

type UserFindParams struct {
	Id int64 `url:"-,omitempty"`
}

type UserCreateParams struct {
	AvatarFile                 io.Writer `url:"avatar_file,omitempty"`
	AvatarDelete               *bool     `url:"avatar_delete,omitempty"`
	ChangePassword             string    `url:"change_password,omitempty"`
	ChangePasswordConfirmation string    `url:"change_password_confirmation,omitempty"`
	Email                      string    `url:"email,omitempty"`
	GrantPermission            string    `url:"grant_permission,omitempty"`
	GroupId                    int64     `url:"group_id,omitempty"`
	GroupIds                   string    `url:"group_ids,omitempty"`
	Password                   string    `url:"password,omitempty"`
	PasswordConfirmation       string    `url:"password_confirmation,omitempty"`
	AnnouncementsRead          *bool     `url:"announcements_read,omitempty"`
	AllowedIps                 string    `url:"allowed_ips,omitempty"`
	AttachmentsPermission      *bool     `url:"attachments_permission,omitempty"`
	AuthenticateUntil          string    `url:"authenticate_until,omitempty"`
	AuthenticationMethod       string    `url:"authentication_method,omitempty"`
	BillingPermission          *bool     `url:"billing_permission,omitempty"`
	BypassInactiveDisable      *bool     `url:"bypass_inactive_disable,omitempty"`
	BypassSiteAllowedIps       *bool     `url:"bypass_site_allowed_ips,omitempty"`
	DavPermission              *bool     `url:"dav_permission,omitempty"`
	Disabled                   *bool     `url:"disabled,omitempty"`
	FtpPermission              *bool     `url:"ftp_permission,omitempty"`
	HeaderText                 string    `url:"header_text,omitempty"`
	Language                   string    `url:"language,omitempty"`
	NotificationDailySendTime  int       `url:"notification_daily_send_time,omitempty"`
	Name                       string    `url:"name,omitempty"`
	Notes                      string    `url:"notes,omitempty"`
	OfficeIntegrationEnabled   *bool     `url:"office_integration_enabled,omitempty"`
	PasswordValidityDays       int       `url:"password_validity_days,omitempty"`
	ReceiveAdminAlerts         *bool     `url:"receive_admin_alerts,omitempty"`
	RequirePasswordChange      *bool     `url:"require_password_change,omitempty"`
	RestapiPermission          *bool     `url:"restapi_permission,omitempty"`
	SelfManaged                *bool     `url:"self_managed,omitempty"`
	SftpPermission             *bool     `url:"sftp_permission,omitempty"`
	SiteAdmin                  *bool     `url:"site_admin,omitempty"`
	SkipWelcomeScreen          *bool     `url:"skip_welcome_screen,omitempty"`
	SslRequired                string    `url:"ssl_required,omitempty"`
	SsoStrategyId              int64     `url:"sso_strategy_id,omitempty"`
	SubscribeToNewsletter      *bool     `url:"subscribe_to_newsletter,omitempty"`
	TimeZone                   string    `url:"time_zone,omitempty"`
	UserRoot                   string    `url:"user_root,omitempty"`
	Username                   string    `url:"username,omitempty"`
}

type UserUnlockParams struct {
	Id int64 `url:"-,omitempty"`
}

type UserResendWelcomeEmailParams struct {
	Id int64 `url:"-,omitempty"`
}

type UserUser2faResetParams struct {
	Id int64 `url:"-,omitempty"`
}

type UserUpdateParams struct {
	Id                         int64     `url:"-,omitempty"`
	AvatarFile                 io.Writer `url:"avatar_file,omitempty"`
	AvatarDelete               *bool     `url:"avatar_delete,omitempty"`
	ChangePassword             string    `url:"change_password,omitempty"`
	ChangePasswordConfirmation string    `url:"change_password_confirmation,omitempty"`
	Email                      string    `url:"email,omitempty"`
	GrantPermission            string    `url:"grant_permission,omitempty"`
	GroupId                    int64     `url:"group_id,omitempty"`
	GroupIds                   string    `url:"group_ids,omitempty"`
	Password                   string    `url:"password,omitempty"`
	PasswordConfirmation       string    `url:"password_confirmation,omitempty"`
	AnnouncementsRead          *bool     `url:"announcements_read,omitempty"`
	AllowedIps                 string    `url:"allowed_ips,omitempty"`
	AttachmentsPermission      *bool     `url:"attachments_permission,omitempty"`
	AuthenticateUntil          string    `url:"authenticate_until,omitempty"`
	AuthenticationMethod       string    `url:"authentication_method,omitempty"`
	BillingPermission          *bool     `url:"billing_permission,omitempty"`
	BypassInactiveDisable      *bool     `url:"bypass_inactive_disable,omitempty"`
	BypassSiteAllowedIps       *bool     `url:"bypass_site_allowed_ips,omitempty"`
	DavPermission              *bool     `url:"dav_permission,omitempty"`
	Disabled                   *bool     `url:"disabled,omitempty"`
	FtpPermission              *bool     `url:"ftp_permission,omitempty"`
	HeaderText                 string    `url:"header_text,omitempty"`
	Language                   string    `url:"language,omitempty"`
	NotificationDailySendTime  int       `url:"notification_daily_send_time,omitempty"`
	Name                       string    `url:"name,omitempty"`
	Notes                      string    `url:"notes,omitempty"`
	OfficeIntegrationEnabled   *bool     `url:"office_integration_enabled,omitempty"`
	PasswordValidityDays       int       `url:"password_validity_days,omitempty"`
	ReceiveAdminAlerts         *bool     `url:"receive_admin_alerts,omitempty"`
	RequirePasswordChange      *bool     `url:"require_password_change,omitempty"`
	RestapiPermission          *bool     `url:"restapi_permission,omitempty"`
	SelfManaged                *bool     `url:"self_managed,omitempty"`
	SftpPermission             *bool     `url:"sftp_permission,omitempty"`
	SiteAdmin                  *bool     `url:"site_admin,omitempty"`
	SkipWelcomeScreen          *bool     `url:"skip_welcome_screen,omitempty"`
	SslRequired                string    `url:"ssl_required,omitempty"`
	SsoStrategyId              int64     `url:"sso_strategy_id,omitempty"`
	SubscribeToNewsletter      *bool     `url:"subscribe_to_newsletter,omitempty"`
	TimeZone                   string    `url:"time_zone,omitempty"`
	UserRoot                   string    `url:"user_root,omitempty"`
	Username                   string    `url:"username,omitempty"`
}

type UserDeleteParams struct {
	Id int64 `url:"-,omitempty"`
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
