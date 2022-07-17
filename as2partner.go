package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type As2Partner struct {
	Id                         int64  `json:"id,omitempty" path:"id"`
	As2StationId               int64  `json:"as2_station_id,omitempty" path:"as2_station_id"`
	Name                       string `json:"name,omitempty" path:"name"`
	Uri                        string `json:"uri,omitempty" path:"uri"`
	ServerCertificate          string `json:"server_certificate,omitempty" path:"server_certificate"`
	HexPublicCertificateSerial string `json:"hex_public_certificate_serial,omitempty" path:"hex_public_certificate_serial"`
	PublicCertificateMd5       string `json:"public_certificate_md5,omitempty" path:"public_certificate_md5"`
	PublicCertificateSubject   string `json:"public_certificate_subject,omitempty" path:"public_certificate_subject"`
	PublicCertificateIssuer    string `json:"public_certificate_issuer,omitempty" path:"public_certificate_issuer"`
	PublicCertificateSerial    string `json:"public_certificate_serial,omitempty" path:"public_certificate_serial"`
	PublicCertificateNotBefore string `json:"public_certificate_not_before,omitempty" path:"public_certificate_not_before"`
	PublicCertificateNotAfter  string `json:"public_certificate_not_after,omitempty" path:"public_certificate_not_after"`
	PublicCertificate          string `json:"public_certificate,omitempty" path:"public_certificate"`
}

type As2PartnerCollection []As2Partner

type As2PartnerListParams struct {
	lib.ListParams
}

type As2PartnerFindParams struct {
	Id int64 `url:"-,omitempty" required:"true" json:"-,omitempty" path:"id"`
}

type As2PartnerCreateParams struct {
	Name              string `url:"name,omitempty" required:"true" json:"name,omitempty" path:"name"`
	Uri               string `url:"uri,omitempty" required:"true" json:"uri,omitempty" path:"uri"`
	PublicCertificate string `url:"public_certificate,omitempty" required:"true" json:"public_certificate,omitempty" path:"public_certificate"`
	As2StationId      int64  `url:"as2_station_id,omitempty" required:"true" json:"as2_station_id,omitempty" path:"as2_station_id"`
	ServerCertificate string `url:"server_certificate,omitempty" required:"false" json:"server_certificate,omitempty" path:"server_certificate"`
}

type As2PartnerUpdateParams struct {
	Id                int64  `url:"-,omitempty" required:"true" json:"-,omitempty" path:"id"`
	Name              string `url:"name,omitempty" required:"false" json:"name,omitempty" path:"name"`
	Uri               string `url:"uri,omitempty" required:"false" json:"uri,omitempty" path:"uri"`
	ServerCertificate string `url:"server_certificate,omitempty" required:"false" json:"server_certificate,omitempty" path:"server_certificate"`
	PublicCertificate string `url:"public_certificate,omitempty" required:"false" json:"public_certificate,omitempty" path:"public_certificate"`
}

type As2PartnerDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"true" json:"-,omitempty" path:"id"`
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
