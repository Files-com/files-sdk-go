package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type As2Station struct {
	Id                   int64  `json:"id,omitempty"`
	Name                 string `json:"name,omitempty"`
	Uri                  string `json:"uri,omitempty"`
	Domain               string `json:"domain,omitempty"`
	PublicCertificate    string `json:"public_certificate,omitempty"`
	PublicCertificateMd5 string `json:"public_certificate_md5,omitempty"`
	PrivateKeyMd5        string `json:"private_key_md5,omitempty"`
	PrivateKey           string `json:"private_key,omitempty"`
}

type As2StationCollection []As2Station

type As2StationListParams struct {
	Cursor  string `url:"cursor,omitempty" required:"false"`
	PerPage int64  `url:"per_page,omitempty" required:"false"`
	lib.ListParams
}

type As2StationFindParams struct {
	Id int64 `url:"-,omitempty" required:"true"`
}

type As2StationCreateParams struct {
	Name              string `url:"name,omitempty" required:"true"`
	Domain            string `url:"domain,omitempty" required:"true"`
	Uri               string `url:"uri,omitempty" required:"true"`
	PublicCertificate string `url:"public_certificate,omitempty" required:"true"`
	PrivateKey        string `url:"private_key,omitempty" required:"true"`
}

type As2StationUpdateParams struct {
	Id                int64  `url:"-,omitempty" required:"true"`
	Name              string `url:"name,omitempty" required:"false"`
	Domain            string `url:"domain,omitempty" required:"false"`
	Uri               string `url:"uri,omitempty" required:"false"`
	PublicCertificate string `url:"public_certificate,omitempty" required:"false"`
	PrivateKey        string `url:"private_key,omitempty" required:"false"`
}

type As2StationDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"true"`
}

func (a *As2Station) UnmarshalJSON(data []byte) error {
	type as2Station As2Station
	var v as2Station
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*a = As2Station(v)
	return nil
}

func (a *As2StationCollection) UnmarshalJSON(data []byte) error {
	type as2Stations []As2Station
	var v as2Stations
	if err := json.Unmarshal(data, &v); err != nil {
		return err
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
