package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// Certificate is a Files.com API resource.
type Certificate struct {
	// Certificate ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Descriptive name of certificate.
	Name string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	// Custom Domain that owns this certificate, including while the certificate is inactive.
	CustomDomainId int64 `json:"custom_domain_id,omitempty" path:"custom_domain_id,omitempty" url:"custom_domain_id,omitempty"`
	// Full text of SSL certificate.
	Certificate string `json:"certificate,omitempty" path:"certificate,omitempty" url:"certificate,omitempty"`
	// Certificate created at.
	CreatedAt *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	// Certificate status. (One of `Request Generated`, `New`, `Active`, `Active/Expired`, `Expired`, or `Available`).
	DisplayStatus string `json:"display_status,omitempty" path:"display_status,omitempty" url:"display_status,omitempty"`
	// Domains on this certificate other than files.com domains.
	Domains []string `json:"domains,omitempty" path:"domains,omitempty" url:"domains,omitempty"`
	// Certificate expiration date/time.
	ExpiresAt *time.Time `json:"expires_at,omitempty" path:"expires_at,omitempty" url:"expires_at,omitempty"`
	// Is this certificate automatically managed and renewed by Files.com?
	BrickManaged *bool `json:"brick_managed,omitempty" path:"brick_managed,omitempty" url:"brick_managed,omitempty"`
	// Intermediate certificates.
	Intermediates string `json:"intermediates,omitempty" path:"intermediates,omitempty" url:"intermediates,omitempty"`
	// A list of IP addresses associated with this SSL Certificate.
	IpAddresses []string `json:"ip_addresses,omitempty" path:"ip_addresses,omitempty" url:"ip_addresses,omitempty"`
	// X509 issuer.
	Issuer string `json:"issuer,omitempty" path:"issuer,omitempty" url:"issuer,omitempty"`
	// Type of key.
	KeyType string `json:"key_type,omitempty" path:"key_type,omitempty" url:"key_type,omitempty"`
	// Certificate signing request text.
	Request string `json:"request,omitempty" path:"request,omitempty" url:"request,omitempty"`
	// Certificate status (Request Generated, New, Active, Active/Expired, Expired, or Available).
	Status string `json:"status,omitempty" path:"status,omitempty" url:"status,omitempty"`
	// X509 Subject name.
	Subject string `json:"subject,omitempty" path:"subject,omitempty" url:"subject,omitempty"`
	// Domain for certificate.
	CertificateDomain string `json:"certificate_domain,omitempty" path:"certificate_domain,omitempty" url:"certificate_domain,omitempty"`
	// Additional domains for the certificate CSR.
	CertificateExtraDomains []string `json:"certificate_extra_domains,omitempty" path:"certificate_extra_domains,omitempty" url:"certificate_extra_domains,omitempty"`
	// Country.
	CertificateCountry string `json:"certificate_country,omitempty" path:"certificate_country,omitempty" url:"certificate_country,omitempty"`
	// State or province.
	CertificateStateOrProvince string `json:"certificate_state_or_province,omitempty" path:"certificate_state_or_province,omitempty" url:"certificate_state_or_province,omitempty"`
	// City or locale.
	CertificateCityOrLocale string `json:"certificate_city_or_locale,omitempty" path:"certificate_city_or_locale,omitempty" url:"certificate_city_or_locale,omitempty"`
	// Company name.
	CertificateCompanyName string `json:"certificate_company_name,omitempty" path:"certificate_company_name,omitempty" url:"certificate_company_name,omitempty"`
	// Department name or organization unit 1.
	CsrOu1 string `json:"csr_ou1,omitempty" path:"csr_ou1,omitempty" url:"csr_ou1,omitempty"`
	// Department name or organization unit 2.
	CsrOu2 string `json:"csr_ou2,omitempty" path:"csr_ou2,omitempty" url:"csr_ou2,omitempty"`
	// Department name or organization unit 3.
	CsrOu3 string `json:"csr_ou3,omitempty" path:"csr_ou3,omitempty" url:"csr_ou3,omitempty"`
	// Email address for the certificate owner.
	CertificateEmailAddress string `json:"certificate_email_address,omitempty" path:"certificate_email_address,omitempty" url:"certificate_email_address,omitempty"`
	// Certificate private key.
	PrivateKey string `json:"private_key,omitempty" path:"private_key,omitempty" url:"private_key,omitempty"`
	// Certificate password.
	Password string `json:"password,omitempty" path:"password,omitempty" url:"password,omitempty"`
}

// Identifier returns the resource ID.
func (c Certificate) Identifier() interface{} {
	return c.Id
}

// CertificateCollection is a list of Certificate resources.
type CertificateCollection []Certificate

// CertificateListParams contains the request parameters for GET /certificates.
type CertificateListParams struct {
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// CertificateFindParams contains the request parameters for GET /certificates/{id}.
type CertificateFindParams struct {
	// Certificate ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// CertificateCreateParams contains the request parameters for POST /certificates.
type CertificateCreateParams struct {
	// Internal name of the SSL certificate. Required.
	Name string `url:"name" json:"name" path:"name"`
	// Domain for certificate.
	CertificateDomain string `url:"certificate_domain,omitempty" json:"certificate_domain,omitempty" path:"certificate_domain"`
	// Additional domains for the certificate CSR.
	CertificateExtraDomains []string `url:"certificate_extra_domains,omitempty" json:"certificate_extra_domains,omitempty" path:"certificate_extra_domains"`
	// Country.
	CertificateCountry string `url:"certificate_country,omitempty" json:"certificate_country,omitempty" path:"certificate_country"`
	// State or province.
	CertificateStateOrProvince string `url:"certificate_state_or_province,omitempty" json:"certificate_state_or_province,omitempty" path:"certificate_state_or_province"`
	// City or locale.
	CertificateCityOrLocale string `url:"certificate_city_or_locale,omitempty" json:"certificate_city_or_locale,omitempty" path:"certificate_city_or_locale"`
	// Company name.
	CertificateCompanyName string `url:"certificate_company_name,omitempty" json:"certificate_company_name,omitempty" path:"certificate_company_name"`
	// Department name or organization unit 1.
	CsrOu1 string `url:"csr_ou1,omitempty" json:"csr_ou1,omitempty" path:"csr_ou1"`
	// Department name or organization unit 2.
	CsrOu2 string `url:"csr_ou2,omitempty" json:"csr_ou2,omitempty" path:"csr_ou2"`
	// Department name or organization unit 3.
	CsrOu3 string `url:"csr_ou3,omitempty" json:"csr_ou3,omitempty" path:"csr_ou3"`
	// Email address for the certificate owner.
	CertificateEmailAddress string `url:"certificate_email_address,omitempty" json:"certificate_email_address,omitempty" path:"certificate_email_address"`
	// Any supported key type. Defaults to `RSA-4096`.
	KeyType string `url:"key_type,omitempty" json:"key_type,omitempty" path:"key_type"`
	// The certificate. PEM or PKCS7 formats accepted.
	Certificate string `url:"certificate,omitempty" json:"certificate,omitempty" path:"certificate"`
	// Certificate private key.
	PrivateKey string `url:"private_key,omitempty" json:"private_key,omitempty" path:"private_key"`
	// Certificate password.
	Password string `url:"password,omitempty" json:"password,omitempty" path:"password"`
	// Intermediate certificates. PEM or PKCS7 formats accepted.
	Intermediates string `url:"intermediates,omitempty" json:"intermediates,omitempty" path:"intermediates"`
}

// CertificateActivateParams contains the request parameters for POST /certificates/{id}/activate.
//
// Activate SSL Certificate
type CertificateActivateParams struct {
	// Certificate ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// Leave blank or set to `any` to replace any currently active certificate, if applicable. Set to `new_ips` to activate this certificate as an additional certificate on your Site by allocating a new set of Dedicated IPs (may require a Plan upgrade). Set to the ID of a currently active certificate to replace that certificate on its set of dedicated IPs. A certificate belonging to a Custom Domain activates only for that domain: leave blank, use `any`, or specify that domain's current certificate ID.
	ReplaceCert string `url:"replace_cert,omitempty" json:"replace_cert,omitempty" path:"replace_cert"`
}

// CertificateDeactivateParams contains the request parameters for POST /certificates/{id}/deactivate.
//
// Deactivate SSL Certificate
type CertificateDeactivateParams struct {
	// Certificate ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// CertificateReplaceWithFilesParams contains the request parameters for POST /certificates/{id}/replace_with_files_certificate.
//
// Replace this active custom SSL Certificate for the Primary Custom Domain with a Files.com SSL Certificate.  The Primary Custom Domain's CNAME record must point at the Site's subdomain.  The Files.com certificate is issued and activated in the background
type CertificateReplaceWithFilesParams struct {
	// Certificate ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// CertificateUpdateParams contains the request parameters for PATCH /certificates/{id}.
type CertificateUpdateParams struct {
	// Certificate ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// Internal certificate name.
	Name string `url:"name,omitempty" json:"name,omitempty" path:"name"`
	// Any intermediate certificates. PEM or PKCS7 formats accepted.
	Intermediates string `url:"intermediates,omitempty" json:"intermediates,omitempty" path:"intermediates"`
	// The certificate. PEM or PKCS7 formats accepted.
	Certificate string `url:"certificate,omitempty" json:"certificate,omitempty" path:"certificate"`
}

// CertificateDeleteParams contains the request parameters for DELETE /certificates/{id}.
type CertificateDeleteParams struct {
	// Certificate ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (c *Certificate) UnmarshalJSON(data []byte) error {
	type certificate Certificate
	var v certificate
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*c = Certificate(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (c *CertificateCollection) UnmarshalJSON(data []byte) error {
	type certificates CertificateCollection
	var v certificates
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*c = CertificateCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (c *CertificateCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*c))
	for i, v := range *c {
		ret[i] = v
	}

	return &ret
}
