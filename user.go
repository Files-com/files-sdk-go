package files_sdk

import (
	"encoding/json"
	"io"
	"time"

	lib "github.com/Files-com/files-sdk-go/lib"
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
	Company                    string    `json:"company,omitempty"`
	Notes                      string    `json:"notes,omitempty"`
	NotificationDailySendTime  int       `json:"notification_daily_send_time,omitempty"`
	OfficeIntegrationEnabled   *bool     `json:"office_integration_enabled,omitempty"`
	PasswordSetAt              time.Time `json:"password_set_at,omitempty"`
	PasswordValidityDays       int       `json:"password_validity_days,omitempty"`
	PublicKeysCount            int       `json:"public_keys_count,omitempty"`
	ReceiveAdminAlerts         *bool     `json:"receive_admin_alerts,omitempty"`
	Require2fa                 *bool     `json:"require_2fa,omitempty"`
	Active2fa                  *bool     `json:"active_2fa,omitempty"`
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
	Page       int             `url:"page,omitempty" required:"false"`
	PerPage    int             `url:"per_page,omitempty" required:"false"`
	Action     string          `url:"action,omitempty" required:"false"`
	Cursor     string          `url:"cursor,omitempty" required:"false"`
	SortBy     json.RawMessage `url:"sort_by,omitempty" required:"false"`
	Filter     json.RawMessage `url:"filter,omitempty" required:"false"`
	FilterGt   json.RawMessage `url:"filter_gt,omitempty" required:"false"`
	FilterGteq json.RawMessage `url:"filter_gteq,omitempty" required:"false"`
	FilterLike json.RawMessage `url:"filter_like,omitempty" required:"false"`
	FilterLt   json.RawMessage `url:"filter_lt,omitempty" required:"false"`
	FilterLteq json.RawMessage `url:"filter_lteq,omitempty" required:"false"`
	Ids        string          `url:"ids,omitempty" required:"false"`
	QParam     QParam          `url:"q,omitempty" required:"false"`
	Search     string          `url:"search,omitempty" required:"false"`
	lib.ListParams
}

type UserFindParams struct {
	Id int64 `url:"-,omitempty" required:"true"`
}

type UserCreateParams struct {
	AvatarFile                 io.Writer `url:"avatar_file,omitempty" required:"false"`
	AvatarDelete               *bool     `url:"avatar_delete,omitempty" required:"false"`
	ChangePassword             string    `url:"change_password,omitempty" required:"false"`
	ChangePasswordConfirmation string    `url:"change_password_confirmation,omitempty" required:"false"`
	Email                      string    `url:"email,omitempty" required:"false"`
	GrantPermission            string    `url:"grant_permission,omitempty" required:"false"`
	GroupId                    int64     `url:"group_id,omitempty" required:"false"`
	GroupIds                   string    `url:"group_ids,omitempty" required:"false"`
	Password                   string    `url:"password,omitempty" required:"false"`
	PasswordConfirmation       string    `url:"password_confirmation,omitempty" required:"false"`
	AnnouncementsRead          *bool     `url:"announcements_read,omitempty" required:"false"`
	AllowedIps                 string    `url:"allowed_ips,omitempty" required:"false"`
	AttachmentsPermission      *bool     `url:"attachments_permission,omitempty" required:"false"`
	AuthenticateUntil          time.Time `url:"authenticate_until,omitempty" required:"false"`
	AuthenticationMethod       string    `url:"authentication_method,omitempty" required:"false"`
	BillingPermission          *bool     `url:"billing_permission,omitempty" required:"false"`
	BypassInactiveDisable      *bool     `url:"bypass_inactive_disable,omitempty" required:"false"`
	BypassSiteAllowedIps       *bool     `url:"bypass_site_allowed_ips,omitempty" required:"false"`
	DavPermission              *bool     `url:"dav_permission,omitempty" required:"false"`
	Disabled                   *bool     `url:"disabled,omitempty" required:"false"`
	FtpPermission              *bool     `url:"ftp_permission,omitempty" required:"false"`
	HeaderText                 string    `url:"header_text,omitempty" required:"false"`
	Language                   string    `url:"language,omitempty" required:"false"`
	NotificationDailySendTime  int       `url:"notification_daily_send_time,omitempty" required:"false"`
	Name                       string    `url:"name,omitempty" required:"false"`
	Company                    string    `url:"company,omitempty" required:"false"`
	Notes                      string    `url:"notes,omitempty" required:"false"`
	OfficeIntegrationEnabled   *bool     `url:"office_integration_enabled,omitempty" required:"false"`
	PasswordValidityDays       int       `url:"password_validity_days,omitempty" required:"false"`
	ReceiveAdminAlerts         *bool     `url:"receive_admin_alerts,omitempty" required:"false"`
	RequirePasswordChange      *bool     `url:"require_password_change,omitempty" required:"false"`
	RestapiPermission          *bool     `url:"restapi_permission,omitempty" required:"false"`
	SelfManaged                *bool     `url:"self_managed,omitempty" required:"false"`
	SftpPermission             *bool     `url:"sftp_permission,omitempty" required:"false"`
	SiteAdmin                  *bool     `url:"site_admin,omitempty" required:"false"`
	SkipWelcomeScreen          *bool     `url:"skip_welcome_screen,omitempty" required:"false"`
	SslRequired                string    `url:"ssl_required,omitempty" required:"false"`
	SsoStrategyId              int64     `url:"sso_strategy_id,omitempty" required:"false"`
	SubscribeToNewsletter      *bool     `url:"subscribe_to_newsletter,omitempty" required:"false"`
	TimeZone                   string    `url:"time_zone,omitempty" required:"false"`
	UserRoot                   string    `url:"user_root,omitempty" required:"false"`
	Username                   string    `url:"username,omitempty" required:"false"`
}

// Unlock user who has been locked out due to failed logins
type UserUnlockParams struct {
	Id int64 `url:"-,omitempty" required:"true"`
}

// Resend user welcome email
type UserResendWelcomeEmailParams struct {
	Id int64 `url:"-,omitempty" required:"true"`
}

// Trigger 2FA Reset process for user who has lost access to their existing 2FA methods
type UserUser2faResetParams struct {
	Id int64 `url:"-,omitempty" required:"true"`
}

type UserUpdateParams struct {
	Id                         int64     `url:"-,omitempty" required:"true"`
	AvatarFile                 io.Writer `url:"avatar_file,omitempty" required:"false"`
	AvatarDelete               *bool     `url:"avatar_delete,omitempty" required:"false"`
	ChangePassword             string    `url:"change_password,omitempty" required:"false"`
	ChangePasswordConfirmation string    `url:"change_password_confirmation,omitempty" required:"false"`
	Email                      string    `url:"email,omitempty" required:"false"`
	GrantPermission            string    `url:"grant_permission,omitempty" required:"false"`
	GroupId                    int64     `url:"group_id,omitempty" required:"false"`
	GroupIds                   string    `url:"group_ids,omitempty" required:"false"`
	Password                   string    `url:"password,omitempty" required:"false"`
	PasswordConfirmation       string    `url:"password_confirmation,omitempty" required:"false"`
	AnnouncementsRead          *bool     `url:"announcements_read,omitempty" required:"false"`
	AllowedIps                 string    `url:"allowed_ips,omitempty" required:"false"`
	AttachmentsPermission      *bool     `url:"attachments_permission,omitempty" required:"false"`
	AuthenticateUntil          time.Time `url:"authenticate_until,omitempty" required:"false"`
	AuthenticationMethod       string    `url:"authentication_method,omitempty" required:"false"`
	BillingPermission          *bool     `url:"billing_permission,omitempty" required:"false"`
	BypassInactiveDisable      *bool     `url:"bypass_inactive_disable,omitempty" required:"false"`
	BypassSiteAllowedIps       *bool     `url:"bypass_site_allowed_ips,omitempty" required:"false"`
	DavPermission              *bool     `url:"dav_permission,omitempty" required:"false"`
	Disabled                   *bool     `url:"disabled,omitempty" required:"false"`
	FtpPermission              *bool     `url:"ftp_permission,omitempty" required:"false"`
	HeaderText                 string    `url:"header_text,omitempty" required:"false"`
	Language                   string    `url:"language,omitempty" required:"false"`
	NotificationDailySendTime  int       `url:"notification_daily_send_time,omitempty" required:"false"`
	Name                       string    `url:"name,omitempty" required:"false"`
	Company                    string    `url:"company,omitempty" required:"false"`
	Notes                      string    `url:"notes,omitempty" required:"false"`
	OfficeIntegrationEnabled   *bool     `url:"office_integration_enabled,omitempty" required:"false"`
	PasswordValidityDays       int       `url:"password_validity_days,omitempty" required:"false"`
	ReceiveAdminAlerts         *bool     `url:"receive_admin_alerts,omitempty" required:"false"`
	RequirePasswordChange      *bool     `url:"require_password_change,omitempty" required:"false"`
	RestapiPermission          *bool     `url:"restapi_permission,omitempty" required:"false"`
	SelfManaged                *bool     `url:"self_managed,omitempty" required:"false"`
	SftpPermission             *bool     `url:"sftp_permission,omitempty" required:"false"`
	SiteAdmin                  *bool     `url:"site_admin,omitempty" required:"false"`
	SkipWelcomeScreen          *bool     `url:"skip_welcome_screen,omitempty" required:"false"`
	SslRequired                string    `url:"ssl_required,omitempty" required:"false"`
	SsoStrategyId              int64     `url:"sso_strategy_id,omitempty" required:"false"`
	SubscribeToNewsletter      *bool     `url:"subscribe_to_newsletter,omitempty" required:"false"`
	TimeZone                   string    `url:"time_zone,omitempty" required:"false"`
	UserRoot                   string    `url:"user_root,omitempty" required:"false"`
	Username                   string    `url:"username,omitempty" required:"false"`
}

type UserDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"true"`
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
