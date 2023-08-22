package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Session struct {
	Id                         string `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	Language                   string `json:"language,omitempty" path:"language,omitempty" url:"language,omitempty"`
	LoginToken                 string `json:"login_token,omitempty" path:"login_token,omitempty" url:"login_token,omitempty"`
	LoginTokenDomain           string `json:"login_token_domain,omitempty" path:"login_token_domain,omitempty" url:"login_token_domain,omitempty"`
	MaxDirListingSize          int64  `json:"max_dir_listing_size,omitempty" path:"max_dir_listing_size,omitempty" url:"max_dir_listing_size,omitempty"`
	MultipleRegions            *bool  `json:"multiple_regions,omitempty" path:"multiple_regions,omitempty" url:"multiple_regions,omitempty"`
	ReadOnly                   *bool  `json:"read_only,omitempty" path:"read_only,omitempty" url:"read_only,omitempty"`
	RootPath                   string `json:"root_path,omitempty" path:"root_path,omitempty" url:"root_path,omitempty"`
	SftpInsecureCiphers        *bool  `json:"sftp_insecure_ciphers,omitempty" path:"sftp_insecure_ciphers,omitempty" url:"sftp_insecure_ciphers,omitempty"`
	SiteId                     int64  `json:"site_id,omitempty" path:"site_id,omitempty" url:"site_id,omitempty"`
	SslRequired                *bool  `json:"ssl_required,omitempty" path:"ssl_required,omitempty" url:"ssl_required,omitempty"`
	TlsDisabled                *bool  `json:"tls_disabled,omitempty" path:"tls_disabled,omitempty" url:"tls_disabled,omitempty"`
	TwoFactorSetupNeeded       *bool  `json:"two_factor_setup_needed,omitempty" path:"two_factor_setup_needed,omitempty" url:"two_factor_setup_needed,omitempty"`
	Allowed2faMethodSms        *bool  `json:"allowed_2fa_method_sms,omitempty" path:"allowed_2fa_method_sms,omitempty" url:"allowed_2fa_method_sms,omitempty"`
	Allowed2faMethodTotp       *bool  `json:"allowed_2fa_method_totp,omitempty" path:"allowed_2fa_method_totp,omitempty" url:"allowed_2fa_method_totp,omitempty"`
	Allowed2faMethodU2f        *bool  `json:"allowed_2fa_method_u2f,omitempty" path:"allowed_2fa_method_u2f,omitempty" url:"allowed_2fa_method_u2f,omitempty"`
	Allowed2faMethodWebauthn   *bool  `json:"allowed_2fa_method_webauthn,omitempty" path:"allowed_2fa_method_webauthn,omitempty" url:"allowed_2fa_method_webauthn,omitempty"`
	Allowed2faMethodYubi       *bool  `json:"allowed_2fa_method_yubi,omitempty" path:"allowed_2fa_method_yubi,omitempty" url:"allowed_2fa_method_yubi,omitempty"`
	UseProvidedModifiedAt      *bool  `json:"use_provided_modified_at,omitempty" path:"use_provided_modified_at,omitempty" url:"use_provided_modified_at,omitempty"`
	WindowsModeFtp             *bool  `json:"windows_mode_ftp,omitempty" path:"windows_mode_ftp,omitempty" url:"windows_mode_ftp,omitempty"`
	Username                   string `json:"username,omitempty" path:"username,omitempty" url:"username,omitempty"`
	Password                   string `json:"password,omitempty" path:"password,omitempty" url:"password,omitempty"`
	ChangePassword             string `json:"change_password,omitempty" path:"change_password,omitempty" url:"change_password,omitempty"`
	ChangePasswordConfirmation string `json:"change_password_confirmation,omitempty" path:"change_password_confirmation,omitempty" url:"change_password_confirmation,omitempty"`
	Interface                  string `json:"interface,omitempty" path:"interface,omitempty" url:"interface,omitempty"`
	Locale                     string `json:"locale,omitempty" path:"locale,omitempty" url:"locale,omitempty"`
	NoCookie                   *bool  `json:"no_cookie,omitempty" path:"no_cookie,omitempty" url:"no_cookie,omitempty"`
	OauthProvider              string `json:"oauth_provider,omitempty" path:"oauth_provider,omitempty" url:"oauth_provider,omitempty"`
	OauthCode                  string `json:"oauth_code,omitempty" path:"oauth_code,omitempty" url:"oauth_code,omitempty"`
	OauthState                 string `json:"oauth_state,omitempty" path:"oauth_state,omitempty" url:"oauth_state,omitempty"`
	Otp                        string `json:"otp,omitempty" path:"otp,omitempty" url:"otp,omitempty"`
	PartialSessionId           string `json:"partial_session_id,omitempty" path:"partial_session_id,omitempty" url:"partial_session_id,omitempty"`
}

func (s Session) Identifier() interface{} {
	return s.Id
}

type SessionCollection []Session

type SessionCreateParams struct {
	Username                   string `url:"username,omitempty" required:"false" json:"username,omitempty" path:"username"`
	Password                   string `url:"password,omitempty" required:"false" json:"password,omitempty" path:"password"`
	ChangePassword             string `url:"change_password,omitempty" required:"false" json:"change_password,omitempty" path:"change_password"`
	ChangePasswordConfirmation string `url:"change_password_confirmation,omitempty" required:"false" json:"change_password_confirmation,omitempty" path:"change_password_confirmation"`
	Interface                  string `url:"interface,omitempty" required:"false" json:"interface,omitempty" path:"interface"`
	Locale                     string `url:"locale,omitempty" required:"false" json:"locale,omitempty" path:"locale"`
	NoCookie                   *bool  `url:"no_cookie,omitempty" required:"false" json:"no_cookie,omitempty" path:"no_cookie"`
	OauthProvider              string `url:"oauth_provider,omitempty" required:"false" json:"oauth_provider,omitempty" path:"oauth_provider"`
	OauthCode                  string `url:"oauth_code,omitempty" required:"false" json:"oauth_code,omitempty" path:"oauth_code"`
	OauthState                 string `url:"oauth_state,omitempty" required:"false" json:"oauth_state,omitempty" path:"oauth_state"`
	Otp                        string `url:"otp,omitempty" required:"false" json:"otp,omitempty" path:"otp"`
	PartialSessionId           string `url:"partial_session_id,omitempty" required:"false" json:"partial_session_id,omitempty" path:"partial_session_id"`
}

type SessionForgotResetParams struct {
	Code            string `url:"code,omitempty" required:"true" json:"code,omitempty" path:"code"`
	Password        string `url:"password,omitempty" required:"true" json:"password,omitempty" path:"password"`
	ConfirmPassword string `url:"confirm_password,omitempty" required:"false" json:"confirm_password,omitempty" path:"confirm_password"`
	Interface       string `url:"interface,omitempty" required:"false" json:"interface,omitempty" path:"interface"`
	Locale          string `url:"locale,omitempty" required:"false" json:"locale,omitempty" path:"locale"`
	Otp             string `url:"otp,omitempty" required:"false" json:"otp,omitempty" path:"otp"`
}

type SessionForgotValidateParams struct {
	Code string `url:"code,omitempty" required:"true" json:"code,omitempty" path:"code"`
}

type SessionForgotParams struct {
	Email           string `url:"email,omitempty" required:"false" json:"email,omitempty" path:"email"`
	Username        string `url:"username,omitempty" required:"false" json:"username,omitempty" path:"username"`
	UsernameOrEmail string `url:"username_or_email,omitempty" required:"false" json:"username_or_email,omitempty" path:"username_or_email"`
}

type SessionPairingKeyParams struct {
	Key string `url:"-,omitempty" required:"false" json:"-,omitempty" path:"key"`
}

type SessionOauthParams struct {
	Provider string `url:"provider,omitempty" required:"true" json:"provider,omitempty" path:"provider"`
	State    string `url:"state,omitempty" required:"false" json:"state,omitempty" path:"state"`
}

func (s *Session) UnmarshalJSON(data []byte) error {
	type session Session
	var v session
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*s = Session(v)
	return nil
}

func (s *SessionCollection) UnmarshalJSON(data []byte) error {
	type sessions SessionCollection
	var v sessions
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*s = SessionCollection(v)
	return nil
}

func (s *SessionCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*s))
	for i, v := range *s {
		ret[i] = v
	}

	return &ret
}
