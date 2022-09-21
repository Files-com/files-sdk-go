package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type SftpHostKey struct {
	Id                int64  `json:"id,omitempty" path:"id"`
	Name              string `json:"name,omitempty" path:"name"`
	FingerprintMd5    string `json:"fingerprint_md5,omitempty" path:"fingerprint_md5"`
	FingerprintSha256 string `json:"fingerprint_sha256,omitempty" path:"fingerprint_sha256"`
	PrivateKey        string `json:"private_key,omitempty" path:"private_key"`
}

type SftpHostKeyCollection []SftpHostKey

type SftpHostKeyListParams struct {
	lib.ListParams
}

type SftpHostKeyFindParams struct {
	Id int64 `url:"-,omitempty" required:"true" json:"-,omitempty" path:"id"`
}

type SftpHostKeyCreateParams struct {
	Name       string `url:"name,omitempty" required:"false" json:"name,omitempty" path:"name"`
	PrivateKey string `url:"private_key,omitempty" required:"false" json:"private_key,omitempty" path:"private_key"`
}

type SftpHostKeyUpdateParams struct {
	Id         int64  `url:"-,omitempty" required:"true" json:"-,omitempty" path:"id"`
	Name       string `url:"name,omitempty" required:"false" json:"name,omitempty" path:"name"`
	PrivateKey string `url:"private_key,omitempty" required:"false" json:"private_key,omitempty" path:"private_key"`
}

type SftpHostKeyDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"true" json:"-,omitempty" path:"id"`
}

func (s *SftpHostKey) UnmarshalJSON(data []byte) error {
	type sftpHostKey SftpHostKey
	var v sftpHostKey
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*s = SftpHostKey(v)
	return nil
}

func (s *SftpHostKeyCollection) UnmarshalJSON(data []byte) error {
	type sftpHostKeys SftpHostKeyCollection
	var v sftpHostKeys
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*s = SftpHostKeyCollection(v)
	return nil
}

func (s *SftpHostKeyCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*s))
	for i, v := range *s {
		ret[i] = v
	}

	return &ret
}
