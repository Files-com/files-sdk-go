package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// SftpHostKey is a Files.com API resource.
type SftpHostKey struct {
	Active            *bool  `json:"active,omitempty" path:"active,omitempty" url:"active,omitempty"`
	CustomDomainId    int64  `json:"custom_domain_id,omitempty" path:"custom_domain_id,omitempty" url:"custom_domain_id,omitempty"`
	Id                int64  `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	Name              string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	KeyType           string `json:"key_type,omitempty" path:"key_type,omitempty" url:"key_type,omitempty"`
	FingerprintMd5    string `json:"fingerprint_md5,omitempty" path:"fingerprint_md5,omitempty" url:"fingerprint_md5,omitempty"`
	FingerprintSha256 string `json:"fingerprint_sha256,omitempty" path:"fingerprint_sha256,omitempty" url:"fingerprint_sha256,omitempty"`
	PrivateKey        string `json:"private_key,omitempty" path:"private_key,omitempty" url:"private_key,omitempty"`
}

// Identifier returns the resource ID.
func (s SftpHostKey) Identifier() interface{} {
	return s.Id
}

// SftpHostKeyCollection is a list of SftpHostKey resources.
type SftpHostKeyCollection []SftpHostKey

// SftpHostKeyListParams contains the request parameters for GET /sftp_host_keys.
type SftpHostKeyListParams struct {
	ListParams
}

// SftpHostKeyFindParams contains the request parameters for GET /sftp_host_keys/{id}.
type SftpHostKeyFindParams struct {
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// SftpHostKeyCreateParams contains the request parameters for POST /sftp_host_keys.
type SftpHostKeyCreateParams struct {
	Active         *bool  `url:"active,omitempty" json:"active,omitempty" path:"active"`
	CustomDomainId int64  `url:"custom_domain_id,omitempty" json:"custom_domain_id,omitempty" path:"custom_domain_id"`
	Name           string `url:"name,omitempty" json:"name,omitempty" path:"name"`
	PrivateKey     string `url:"private_key,omitempty" json:"private_key,omitempty" path:"private_key"`
}

// SftpHostKeyUpdateParams contains the request parameters for PATCH /sftp_host_keys/{id}.
type SftpHostKeyUpdateParams struct {
	Id             int64  `url:"-,omitempty" json:"-,omitempty" path:"id"`
	Active         *bool  `url:"active,omitempty" json:"active,omitempty" path:"active"`
	CustomDomainId int64  `url:"custom_domain_id,omitempty" json:"custom_domain_id,omitempty" path:"custom_domain_id"`
	Name           string `url:"name,omitempty" json:"name,omitempty" path:"name"`
	PrivateKey     string `url:"private_key,omitempty" json:"private_key,omitempty" path:"private_key"`
}

// SftpHostKeyDeleteParams contains the request parameters for DELETE /sftp_host_keys/{id}.
type SftpHostKeyDeleteParams struct {
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (s *SftpHostKey) UnmarshalJSON(data []byte) error {
	type sftpHostKey SftpHostKey
	var v sftpHostKey
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*s = SftpHostKey(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (s *SftpHostKeyCollection) UnmarshalJSON(data []byte) error {
	type sftpHostKeys SftpHostKeyCollection
	var v sftpHostKeys
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*s = SftpHostKeyCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (s *SftpHostKeyCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*s))
	for i, v := range *s {
		ret[i] = v
	}

	return &ret
}
