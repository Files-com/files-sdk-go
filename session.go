package files_sdk

import (
	"encoding/json"
)

type Session struct {
	Id                       string `json:"id,omitempty"`
	Language                 string `json:"language,omitempty"`
	LoginToken               string `json:"login_token,omitempty"`
	LoginTokenDomain         string `json:"login_token_domain,omitempty"`
	MaxDirListingSize        int64  `json:"max_dir_listing_size,omitempty"`
	MultipleRegions          *bool  `json:"multiple_regions,omitempty"`
	ReadOnly                 *bool  `json:"read_only,omitempty"`
	RootPath                 string `json:"root_path,omitempty"`
	SftpInsecureCiphers      *bool  `json:"sftp_insecure_ciphers,omitempty"`
	SiteId                   int64  `json:"site_id,omitempty"`
	SslRequired              *bool  `json:"ssl_required,omitempty"`
	TlsDisabled              *bool  `json:"tls_disabled,omitempty"`
	TwoFactorSetupNeeded     *bool  `json:"two_factor_setup_needed,omitempty"`
	Allowed2faMethodSms      *bool  `json:"allowed_2fa_method_sms,omitempty"`
	Allowed2faMethodTotp     *bool  `json:"allowed_2fa_method_totp,omitempty"`
	Allowed2faMethodU2f      *bool  `json:"allowed_2fa_method_u2f,omitempty"`
	Allowed2faMethodWebauthn *bool  `json:"allowed_2fa_method_webauthn,omitempty"`
	Allowed2faMethodYubi     *bool  `json:"allowed_2fa_method_yubi,omitempty"`
	UseProvidedModifiedAt    *bool  `json:"use_provided_modified_at,omitempty"`
	WindowsModeFtp           *bool  `json:"windows_mode_ftp,omitempty"`
	Username                 string `json:"username,omitempty"`
	Password                 string `json:"password,omitempty"`
	Otp                      string `json:"otp,omitempty"`
	PartialSessionId         string `json:"partial_session_id,omitempty"`
}

type SessionCollection []Session

type SessionCreateParams struct {
	Username         string `url:"username,omitempty" required:"false" json:"username,omitempty"`
	Password         string `url:"password,omitempty" required:"false" json:"password,omitempty"`
	Otp              string `url:"otp,omitempty" required:"false" json:"otp,omitempty"`
	PartialSessionId string `url:"partial_session_id,omitempty" required:"false" json:"partial_session_id,omitempty"`
}

func (s *Session) UnmarshalJSON(data []byte) error {
	type session Session
	var v session
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*s = Session(v)
	return nil
}

func (s *SessionCollection) UnmarshalJSON(data []byte) error {
	type sessions []Session
	var v sessions
	if err := json.Unmarshal(data, &v); err != nil {
		return err
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
