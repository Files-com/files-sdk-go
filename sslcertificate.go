package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type SslCertificate struct {
	Name                string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	Certificate         string `json:"certificate,omitempty" path:"certificate,omitempty" url:"certificate,omitempty"`
	PrivateKey          string `json:"private_key,omitempty" path:"private_key,omitempty" url:"private_key,omitempty"`
	Key                 string `json:"key,omitempty" path:"key,omitempty" url:"key,omitempty"`
	Intermediates       string `json:"intermediates,omitempty" path:"intermediates,omitempty" url:"intermediates,omitempty"`
	Id                  string `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	DomainHstsHeader    string `json:"domain_hsts_header,omitempty" path:"domain_hsts_header,omitempty" url:"domain_hsts_header,omitempty"`
	FtpsEnabled         string `json:"ftps_enabled,omitempty" path:"ftps_enabled,omitempty" url:"ftps_enabled,omitempty"`
	HttpsEnabled        string `json:"https_enabled,omitempty" path:"https_enabled,omitempty" url:"https_enabled,omitempty"`
	SftpInsecureCiphers string `json:"sftp_insecure_ciphers,omitempty" path:"sftp_insecure_ciphers,omitempty" url:"sftp_insecure_ciphers,omitempty"`
	TlsDisabled         string `json:"tls_disabled,omitempty" path:"tls_disabled,omitempty" url:"tls_disabled,omitempty"`
	Subdomain           string `json:"subdomain,omitempty" path:"subdomain,omitempty" url:"subdomain,omitempty"`
}

func (s SslCertificate) Identifier() interface{} {
	return s.Id
}

type SslCertificateCollection []SslCertificate

func (s *SslCertificate) UnmarshalJSON(data []byte) error {
	type sslCertificate SslCertificate
	var v sslCertificate
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*s = SslCertificate(v)
	return nil
}

func (s *SslCertificateCollection) UnmarshalJSON(data []byte) error {
	type sslCertificates SslCertificateCollection
	var v sslCertificates
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*s = SslCertificateCollection(v)
	return nil
}

func (s *SslCertificateCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*s))
	for i, v := range *s {
		ret[i] = v
	}

	return &ret
}
