package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Certificate struct {
	Id                         int64      `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	Name                       string     `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	Certificate                string     `json:"certificate,omitempty" path:"certificate,omitempty" url:"certificate,omitempty"`
	CreatedAt                  *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	DisplayStatus              string     `json:"display_status,omitempty" path:"display_status,omitempty" url:"display_status,omitempty"`
	Domains                    []string   `json:"domains,omitempty" path:"domains,omitempty" url:"domains,omitempty"`
	ExpiresAt                  *time.Time `json:"expires_at,omitempty" path:"expires_at,omitempty" url:"expires_at,omitempty"`
	BrickManaged               *bool      `json:"brick_managed,omitempty" path:"brick_managed,omitempty" url:"brick_managed,omitempty"`
	Intermediates              string     `json:"intermediates,omitempty" path:"intermediates,omitempty" url:"intermediates,omitempty"`
	IpAddresses                []string   `json:"ip_addresses,omitempty" path:"ip_addresses,omitempty" url:"ip_addresses,omitempty"`
	Issuer                     string     `json:"issuer,omitempty" path:"issuer,omitempty" url:"issuer,omitempty"`
	KeyType                    string     `json:"key_type,omitempty" path:"key_type,omitempty" url:"key_type,omitempty"`
	Request                    string     `json:"request,omitempty" path:"request,omitempty" url:"request,omitempty"`
	Status                     string     `json:"status,omitempty" path:"status,omitempty" url:"status,omitempty"`
	Subject                    string     `json:"subject,omitempty" path:"subject,omitempty" url:"subject,omitempty"`
	UpdatedAt                  *time.Time `json:"updated_at,omitempty" path:"updated_at,omitempty" url:"updated_at,omitempty"`
	CertificateDomain          string     `json:"certificate_domain,omitempty" path:"certificate_domain,omitempty" url:"certificate_domain,omitempty"`
	CertificateCountry         string     `json:"certificate_country,omitempty" path:"certificate_country,omitempty" url:"certificate_country,omitempty"`
	CertificateStateOrProvince string     `json:"certificate_state_or_province,omitempty" path:"certificate_state_or_province,omitempty" url:"certificate_state_or_province,omitempty"`
	CertificateCityOrLocale    string     `json:"certificate_city_or_locale,omitempty" path:"certificate_city_or_locale,omitempty" url:"certificate_city_or_locale,omitempty"`
	CertificateCompanyName     string     `json:"certificate_company_name,omitempty" path:"certificate_company_name,omitempty" url:"certificate_company_name,omitempty"`
	CsrOu1                     string     `json:"csr_ou1,omitempty" path:"csr_ou1,omitempty" url:"csr_ou1,omitempty"`
	CsrOu2                     string     `json:"csr_ou2,omitempty" path:"csr_ou2,omitempty" url:"csr_ou2,omitempty"`
	CsrOu3                     string     `json:"csr_ou3,omitempty" path:"csr_ou3,omitempty" url:"csr_ou3,omitempty"`
	CertificateEmailAddress    string     `json:"certificate_email_address,omitempty" path:"certificate_email_address,omitempty" url:"certificate_email_address,omitempty"`
	PrivateKey                 string     `json:"private_key,omitempty" path:"private_key,omitempty" url:"private_key,omitempty"`
	Password                   string     `json:"password,omitempty" path:"password,omitempty" url:"password,omitempty"`
}

func (c Certificate) Identifier() interface{} {
	return c.Id
}

type CertificateCollection []Certificate

type CertificateListParams struct {
	Action string `url:"action,omitempty" required:"false" json:"action,omitempty" path:"action"`
	ListParams
}

type CertificateFindParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
}

type CertificateCreateParams struct {
	Name                       string `url:"name,omitempty" required:"true" json:"name,omitempty" path:"name"`
	CertificateDomain          string `url:"certificate_domain,omitempty" required:"false" json:"certificate_domain,omitempty" path:"certificate_domain"`
	CertificateCountry         string `url:"certificate_country,omitempty" required:"false" json:"certificate_country,omitempty" path:"certificate_country"`
	CertificateStateOrProvince string `url:"certificate_state_or_province,omitempty" required:"false" json:"certificate_state_or_province,omitempty" path:"certificate_state_or_province"`
	CertificateCityOrLocale    string `url:"certificate_city_or_locale,omitempty" required:"false" json:"certificate_city_or_locale,omitempty" path:"certificate_city_or_locale"`
	CertificateCompanyName     string `url:"certificate_company_name,omitempty" required:"false" json:"certificate_company_name,omitempty" path:"certificate_company_name"`
	CsrOu1                     string `url:"csr_ou1,omitempty" required:"false" json:"csr_ou1,omitempty" path:"csr_ou1"`
	CsrOu2                     string `url:"csr_ou2,omitempty" required:"false" json:"csr_ou2,omitempty" path:"csr_ou2"`
	CsrOu3                     string `url:"csr_ou3,omitempty" required:"false" json:"csr_ou3,omitempty" path:"csr_ou3"`
	CertificateEmailAddress    string `url:"certificate_email_address,omitempty" required:"false" json:"certificate_email_address,omitempty" path:"certificate_email_address"`
	KeyType                    string `url:"key_type,omitempty" required:"false" json:"key_type,omitempty" path:"key_type"`
	Certificate                string `url:"certificate,omitempty" required:"false" json:"certificate,omitempty" path:"certificate"`
	PrivateKey                 string `url:"private_key,omitempty" required:"false" json:"private_key,omitempty" path:"private_key"`
	Password                   string `url:"password,omitempty" required:"false" json:"password,omitempty" path:"password"`
	Intermediates              string `url:"intermediates,omitempty" required:"false" json:"intermediates,omitempty" path:"intermediates"`
}

// Deactivate SSL Certificate
type CertificateDeactivateParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
}

// Activate SSL Certificate
type CertificateActivateParams struct {
	Id          int64  `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
	ReplaceCert string `url:"replace_cert,omitempty" required:"false" json:"replace_cert,omitempty" path:"replace_cert"`
}

type CertificateUpdateParams struct {
	Id            int64  `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
	Name          string `url:"name,omitempty" required:"false" json:"name,omitempty" path:"name"`
	Intermediates string `url:"intermediates,omitempty" required:"false" json:"intermediates,omitempty" path:"intermediates"`
	Certificate   string `url:"certificate,omitempty" required:"false" json:"certificate,omitempty" path:"certificate"`
}

type CertificateDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
}

func (c *Certificate) UnmarshalJSON(data []byte) error {
	type certificate Certificate
	var v certificate
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*c = Certificate(v)
	return nil
}

func (c *CertificateCollection) UnmarshalJSON(data []byte) error {
	type certificates CertificateCollection
	var v certificates
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*c = CertificateCollection(v)
	return nil
}

func (c *CertificateCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*c))
	for i, v := range *c {
		ret[i] = v
	}

	return &ret
}
