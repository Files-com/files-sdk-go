package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type As2Partner struct {
	Id                         int64  `json:"id,omitempty"`
	As2StationId               int64  `json:"as2_station_id,omitempty"`
	Name                       string `json:"name,omitempty"`
	Uri                        string `json:"uri,omitempty"`
	PublicCertificateMd5       string `json:"public_certificate_md5,omitempty"`
	PublicCertificateSubject   string `json:"public_certificate_subject,omitempty"`
	PublicCertificateIssuer    string `json:"public_certificate_issuer,omitempty"`
	PublicCertificateSerial    string `json:"public_certificate_serial,omitempty"`
	PublicCertificateNotBefore string `json:"public_certificate_not_before,omitempty"`
	PublicCertificateNotAfter  string `json:"public_certificate_not_after,omitempty"`
	PublicCertificate          string `json:"public_certificate,omitempty"`
}

type As2PartnerCollection []As2Partner

type As2PartnerListParams struct {
	Cursor  string `url:"cursor,omitempty" required:"false"`
	PerPage int64  `url:"per_page,omitempty" required:"false"`
	lib.ListParams
}

type As2PartnerFindParams struct {
	Id int64 `url:"-,omitempty" required:"true"`
}

type As2PartnerCreateParams struct {
	Name              string `url:"name,omitempty" required:"true"`
	Uri               string `url:"uri,omitempty" required:"true"`
	PublicCertificate string `url:"public_certificate,omitempty" required:"true"`
	As2StationId      int64  `url:"as2_station_id,omitempty" required:"true"`
}

type As2PartnerUpdateParams struct {
	Id                int64  `url:"-,omitempty" required:"true"`
	Name              string `url:"name,omitempty" required:"false"`
	Uri               string `url:"uri,omitempty" required:"false"`
	PublicCertificate string `url:"public_certificate,omitempty" required:"false"`
}

type As2PartnerDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"true"`
}

func (a *As2Partner) UnmarshalJSON(data []byte) error {
	type as2Partner As2Partner
	var v as2Partner
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*a = As2Partner(v)
	return nil
}

func (a *As2PartnerCollection) UnmarshalJSON(data []byte) error {
	type as2Partners []As2Partner
	var v as2Partners
	if err := json.Unmarshal(data, &v); err != nil {
		return err
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
