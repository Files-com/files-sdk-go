package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type As2Station struct {
	Id                         int64  `json:"id,omitempty" path:"id"`
	Name                       string `json:"name,omitempty" path:"name"`
	Uri                        string `json:"uri,omitempty" path:"uri"`
	Domain                     string `json:"domain,omitempty" path:"domain"`
	HexPublicCertificateSerial string `json:"hex_public_certificate_serial,omitempty" path:"hex_public_certificate_serial"`
	PublicCertificateMd5       string `json:"public_certificate_md5,omitempty" path:"public_certificate_md5"`
	PrivateKeyMd5              string `json:"private_key_md5,omitempty" path:"private_key_md5"`
	PublicCertificateSubject   string `json:"public_certificate_subject,omitempty" path:"public_certificate_subject"`
	PublicCertificateIssuer    string `json:"public_certificate_issuer,omitempty" path:"public_certificate_issuer"`
	PublicCertificateSerial    string `json:"public_certificate_serial,omitempty" path:"public_certificate_serial"`
	PublicCertificateNotBefore string `json:"public_certificate_not_before,omitempty" path:"public_certificate_not_before"`
	PublicCertificateNotAfter  string `json:"public_certificate_not_after,omitempty" path:"public_certificate_not_after"`
	PrivateKeyPasswordMd5      string `json:"private_key_password_md5,omitempty" path:"private_key_password_md5"`
	PublicCertificate          string `json:"public_certificate,omitempty" path:"public_certificate"`
	PrivateKey                 string `json:"private_key,omitempty" path:"private_key"`
	PrivateKeyPassword         string `json:"private_key_password,omitempty" path:"private_key_password"`
}

type As2StationCollection []As2Station

type As2StationListParams struct {
	lib.ListParams
}

type As2StationFindParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
}

type As2StationCreateParams struct {
	Name               string `url:"name,omitempty" required:"true" json:"name,omitempty" path:"name"`
	PublicCertificate  string `url:"public_certificate,omitempty" required:"true" json:"public_certificate,omitempty" path:"public_certificate"`
	PrivateKey         string `url:"private_key,omitempty" required:"true" json:"private_key,omitempty" path:"private_key"`
	PrivateKeyPassword string `url:"private_key_password,omitempty" required:"false" json:"private_key_password,omitempty" path:"private_key_password"`
}

type As2StationUpdateParams struct {
	Id                 int64  `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
	Name               string `url:"name,omitempty" required:"false" json:"name,omitempty" path:"name"`
	PublicCertificate  string `url:"public_certificate,omitempty" required:"false" json:"public_certificate,omitempty" path:"public_certificate"`
	PrivateKey         string `url:"private_key,omitempty" required:"false" json:"private_key,omitempty" path:"private_key"`
	PrivateKeyPassword string `url:"private_key_password,omitempty" required:"false" json:"private_key_password,omitempty" path:"private_key_password"`
}

type As2StationDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
}

func (a *As2Station) UnmarshalJSON(data []byte) error {
	type as2Station As2Station
	var v as2Station
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*a = As2Station(v)
	return nil
}

func (a *As2StationCollection) UnmarshalJSON(data []byte) error {
	type as2Stations As2StationCollection
	var v as2Stations
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*a = As2StationCollection(v)
	return nil
}

func (a *As2StationCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*a))
	for i, v := range *a {
		ret[i] = v
	}

	return &ret
}
