package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

type As2Partner struct {
	Id                         int64  `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	As2StationId               int64  `json:"as2_station_id,omitempty" path:"as2_station_id,omitempty" url:"as2_station_id,omitempty"`
	Name                       string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	Uri                        string `json:"uri,omitempty" path:"uri,omitempty" url:"uri,omitempty"`
	ServerCertificate          string `json:"server_certificate,omitempty" path:"server_certificate,omitempty" url:"server_certificate,omitempty"`
	HttpAuthUsername           string `json:"http_auth_username,omitempty" path:"http_auth_username,omitempty" url:"http_auth_username,omitempty"`
	MdnValidationLevel         string `json:"mdn_validation_level,omitempty" path:"mdn_validation_level,omitempty" url:"mdn_validation_level,omitempty"`
	EnableDedicatedIps         *bool  `json:"enable_dedicated_ips,omitempty" path:"enable_dedicated_ips,omitempty" url:"enable_dedicated_ips,omitempty"`
	HexPublicCertificateSerial string `json:"hex_public_certificate_serial,omitempty" path:"hex_public_certificate_serial,omitempty" url:"hex_public_certificate_serial,omitempty"`
	PublicCertificateMd5       string `json:"public_certificate_md5,omitempty" path:"public_certificate_md5,omitempty" url:"public_certificate_md5,omitempty"`
	PublicCertificateSubject   string `json:"public_certificate_subject,omitempty" path:"public_certificate_subject,omitempty" url:"public_certificate_subject,omitempty"`
	PublicCertificateIssuer    string `json:"public_certificate_issuer,omitempty" path:"public_certificate_issuer,omitempty" url:"public_certificate_issuer,omitempty"`
	PublicCertificateSerial    string `json:"public_certificate_serial,omitempty" path:"public_certificate_serial,omitempty" url:"public_certificate_serial,omitempty"`
	PublicCertificateNotBefore string `json:"public_certificate_not_before,omitempty" path:"public_certificate_not_before,omitempty" url:"public_certificate_not_before,omitempty"`
	PublicCertificateNotAfter  string `json:"public_certificate_not_after,omitempty" path:"public_certificate_not_after,omitempty" url:"public_certificate_not_after,omitempty"`
	HttpAuthPassword           string `json:"http_auth_password,omitempty" path:"http_auth_password,omitempty" url:"http_auth_password,omitempty"`
	PublicCertificate          string `json:"public_certificate,omitempty" path:"public_certificate,omitempty" url:"public_certificate,omitempty"`
}

func (a As2Partner) Identifier() interface{} {
	return a.Id
}

type As2PartnerCollection []As2Partner

type As2PartnerMdnValidationLevelEnum string

func (u As2PartnerMdnValidationLevelEnum) String() string {
	return string(u)
}

func (u As2PartnerMdnValidationLevelEnum) Enum() map[string]As2PartnerMdnValidationLevelEnum {
	return map[string]As2PartnerMdnValidationLevelEnum{
		"none":   As2PartnerMdnValidationLevelEnum("none"),
		"weak":   As2PartnerMdnValidationLevelEnum("weak"),
		"normal": As2PartnerMdnValidationLevelEnum("normal"),
		"strict": As2PartnerMdnValidationLevelEnum("strict"),
	}
}

type As2PartnerServerCertificateEnum string

func (u As2PartnerServerCertificateEnum) String() string {
	return string(u)
}

func (u As2PartnerServerCertificateEnum) Enum() map[string]As2PartnerServerCertificateEnum {
	return map[string]As2PartnerServerCertificateEnum{
		"require_match": As2PartnerServerCertificateEnum("require_match"),
		"allow_any":     As2PartnerServerCertificateEnum("allow_any"),
	}
}

type As2PartnerListParams struct {
	ListParams
}

type As2PartnerFindParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
}

type As2PartnerCreateParams struct {
	EnableDedicatedIps *bool                            `url:"enable_dedicated_ips,omitempty" required:"false" json:"enable_dedicated_ips,omitempty" path:"enable_dedicated_ips"`
	HttpAuthUsername   string                           `url:"http_auth_username,omitempty" required:"false" json:"http_auth_username,omitempty" path:"http_auth_username"`
	HttpAuthPassword   string                           `url:"http_auth_password,omitempty" required:"false" json:"http_auth_password,omitempty" path:"http_auth_password"`
	MdnValidationLevel As2PartnerMdnValidationLevelEnum `url:"mdn_validation_level,omitempty" required:"false" json:"mdn_validation_level,omitempty" path:"mdn_validation_level"`
	ServerCertificate  As2PartnerServerCertificateEnum  `url:"server_certificate,omitempty" required:"false" json:"server_certificate,omitempty" path:"server_certificate"`
	As2StationId       int64                            `url:"as2_station_id,omitempty" required:"true" json:"as2_station_id,omitempty" path:"as2_station_id"`
	Name               string                           `url:"name,omitempty" required:"true" json:"name,omitempty" path:"name"`
	Uri                string                           `url:"uri,omitempty" required:"true" json:"uri,omitempty" path:"uri"`
	PublicCertificate  string                           `url:"public_certificate,omitempty" required:"true" json:"public_certificate,omitempty" path:"public_certificate"`
}

type As2PartnerUpdateParams struct {
	Id                 int64                            `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
	EnableDedicatedIps *bool                            `url:"enable_dedicated_ips,omitempty" required:"false" json:"enable_dedicated_ips,omitempty" path:"enable_dedicated_ips"`
	HttpAuthUsername   string                           `url:"http_auth_username,omitempty" required:"false" json:"http_auth_username,omitempty" path:"http_auth_username"`
	HttpAuthPassword   string                           `url:"http_auth_password,omitempty" required:"false" json:"http_auth_password,omitempty" path:"http_auth_password"`
	MdnValidationLevel As2PartnerMdnValidationLevelEnum `url:"mdn_validation_level,omitempty" required:"false" json:"mdn_validation_level,omitempty" path:"mdn_validation_level"`
	ServerCertificate  As2PartnerServerCertificateEnum  `url:"server_certificate,omitempty" required:"false" json:"server_certificate,omitempty" path:"server_certificate"`
	Name               string                           `url:"name,omitempty" required:"false" json:"name,omitempty" path:"name"`
	Uri                string                           `url:"uri,omitempty" required:"false" json:"uri,omitempty" path:"uri"`
	PublicCertificate  string                           `url:"public_certificate,omitempty" required:"false" json:"public_certificate,omitempty" path:"public_certificate"`
}

type As2PartnerDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
}

func (a *As2Partner) UnmarshalJSON(data []byte) error {
	type as2Partner As2Partner
	var v as2Partner
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*a = As2Partner(v)
	return nil
}

func (a *As2PartnerCollection) UnmarshalJSON(data []byte) error {
	type as2Partners As2PartnerCollection
	var v as2Partners
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*a = As2PartnerCollection(v)
	return nil
}

func (a *As2PartnerCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*a))
	for i, v := range *a {
		ret[i] = v
	}

	return &ret
}
