package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type TwoFactorAuthenticationMethod struct {
	Id                          int64      `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	Name                        string     `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	MethodType                  string     `json:"method_type,omitempty" path:"method_type,omitempty" url:"method_type,omitempty"`
	PhoneNumber                 string     `json:"phone_number,omitempty" path:"phone_number,omitempty" url:"phone_number,omitempty"`
	PhoneNumberCountry          string     `json:"phone_number_country,omitempty" path:"phone_number_country,omitempty" url:"phone_number_country,omitempty"`
	PhoneNumberNationalFormat   string     `json:"phone_number_national_format,omitempty" path:"phone_number_national_format,omitempty" url:"phone_number_national_format,omitempty"`
	SetupExpired                *bool      `json:"setup_expired,omitempty" path:"setup_expired,omitempty" url:"setup_expired,omitempty"`
	SetupComplete               *bool      `json:"setup_complete,omitempty" path:"setup_complete,omitempty" url:"setup_complete,omitempty"`
	SetupExpiresAt              *time.Time `json:"setup_expires_at,omitempty" path:"setup_expires_at,omitempty" url:"setup_expires_at,omitempty"`
	TotpProvisioningUri         string     `json:"totp_provisioning_uri,omitempty" path:"totp_provisioning_uri,omitempty" url:"totp_provisioning_uri,omitempty"`
	U2fAppId                    string     `json:"u2f_app_id,omitempty" path:"u2f_app_id,omitempty" url:"u2f_app_id,omitempty"`
	U2fRegistrationRequests     []string   `json:"u2f_registration_requests,omitempty" path:"u2f_registration_requests,omitempty" url:"u2f_registration_requests,omitempty"`
	WebauthnRegistrationOptions []string   `json:"webauthn_registration_options,omitempty" path:"webauthn_registration_options,omitempty" url:"webauthn_registration_options,omitempty"`
	BypassForFtpSftpDav         *bool      `json:"bypass_for_ftp_sftp_dav,omitempty" path:"bypass_for_ftp_sftp_dav,omitempty" url:"bypass_for_ftp_sftp_dav,omitempty"`
	Otp                         string     `json:"otp,omitempty" path:"otp,omitempty" url:"otp,omitempty"`
}

func (t TwoFactorAuthenticationMethod) Identifier() interface{} {
	return t.Id
}

type TwoFactorAuthenticationMethodCollection []TwoFactorAuthenticationMethod

type TwoFactorAuthenticationMethodMethodTypeEnum string

func (u TwoFactorAuthenticationMethodMethodTypeEnum) String() string {
	return string(u)
}

func (u TwoFactorAuthenticationMethodMethodTypeEnum) Enum() map[string]TwoFactorAuthenticationMethodMethodTypeEnum {
	return map[string]TwoFactorAuthenticationMethodMethodTypeEnum{
		"totp":     TwoFactorAuthenticationMethodMethodTypeEnum("totp"),
		"yubi":     TwoFactorAuthenticationMethodMethodTypeEnum("yubi"),
		"sms":      TwoFactorAuthenticationMethodMethodTypeEnum("sms"),
		"u2f":      TwoFactorAuthenticationMethodMethodTypeEnum("u2f"),
		"webauthn": TwoFactorAuthenticationMethodMethodTypeEnum("webauthn"),
	}
}

type TwoFactorAuthenticationMethodGetParams struct {
	Action string `url:"action,omitempty" required:"false" json:"action,omitempty" path:"action"`
	ListParams
}

type TwoFactorAuthenticationMethodCreateParams struct {
	MethodType          TwoFactorAuthenticationMethodMethodTypeEnum `url:"method_type,omitempty" required:"true" json:"method_type,omitempty" path:"method_type"`
	Name                string                                      `url:"name,omitempty" required:"false" json:"name,omitempty" path:"name"`
	PhoneNumber         string                                      `url:"phone_number,omitempty" required:"false" json:"phone_number,omitempty" path:"phone_number"`
	BypassForFtpSftpDav *bool                                       `url:"bypass_for_ftp_sftp_dav,omitempty" required:"false" json:"bypass_for_ftp_sftp_dav,omitempty" path:"bypass_for_ftp_sftp_dav"`
}

type TwoFactorAuthenticationMethodSendCodeParams struct {
	U2fOnly *bool `url:"u2f_only,omitempty" required:"false" json:"u2f_only,omitempty" path:"u2f_only"`
}

type TwoFactorAuthenticationMethodUpdateParams struct {
	Id                  int64  `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
	Otp                 string `url:"otp,omitempty" required:"false" json:"otp,omitempty" path:"otp"`
	Name                string `url:"name,omitempty" required:"false" json:"name,omitempty" path:"name"`
	BypassForFtpSftpDav *bool  `url:"bypass_for_ftp_sftp_dav,omitempty" required:"false" json:"bypass_for_ftp_sftp_dav,omitempty" path:"bypass_for_ftp_sftp_dav"`
}

type TwoFactorAuthenticationMethodDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
}

func (t *TwoFactorAuthenticationMethod) UnmarshalJSON(data []byte) error {
	type twoFactorAuthenticationMethod TwoFactorAuthenticationMethod
	var v twoFactorAuthenticationMethod
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*t = TwoFactorAuthenticationMethod(v)
	return nil
}

func (t *TwoFactorAuthenticationMethodCollection) UnmarshalJSON(data []byte) error {
	type twoFactorAuthenticationMethods TwoFactorAuthenticationMethodCollection
	var v twoFactorAuthenticationMethods
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*t = TwoFactorAuthenticationMethodCollection(v)
	return nil
}

func (t *TwoFactorAuthenticationMethodCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*t))
	for i, v := range *t {
		ret[i] = v
	}

	return &ret
}
