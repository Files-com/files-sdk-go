package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

type PublicKey struct {
	Id                int64      `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	Title             string     `json:"title,omitempty" path:"title,omitempty" url:"title,omitempty"`
	CreatedAt         *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	Fingerprint       string     `json:"fingerprint,omitempty" path:"fingerprint,omitempty" url:"fingerprint,omitempty"`
	FingerprintSha256 string     `json:"fingerprint_sha256,omitempty" path:"fingerprint_sha256,omitempty" url:"fingerprint_sha256,omitempty"`
	Username          string     `json:"username,omitempty" path:"username,omitempty" url:"username,omitempty"`
	UserId            int64      `json:"user_id,omitempty" path:"user_id,omitempty" url:"user_id,omitempty"`
	PublicKey         string     `json:"public_key,omitempty" path:"public_key,omitempty" url:"public_key,omitempty"`
}

func (p PublicKey) Identifier() interface{} {
	return p.Id
}

type PublicKeyCollection []PublicKey

type PublicKeyListParams struct {
	UserId int64  `url:"user_id,omitempty" required:"false" json:"user_id,omitempty" path:"user_id"`
	Action string `url:"action,omitempty" required:"false" json:"action,omitempty" path:"action"`
	ListParams
}

type PublicKeyFindParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
}

type PublicKeyCreateParams struct {
	UserId    int64  `url:"user_id,omitempty" required:"false" json:"user_id,omitempty" path:"user_id"`
	Title     string `url:"title,omitempty" required:"true" json:"title,omitempty" path:"title"`
	PublicKey string `url:"public_key,omitempty" required:"true" json:"public_key,omitempty" path:"public_key"`
}

type PublicKeyUpdateParams struct {
	Id    int64  `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
	Title string `url:"title,omitempty" required:"true" json:"title,omitempty" path:"title"`
}

type PublicKeyDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
}

func (p *PublicKey) UnmarshalJSON(data []byte) error {
	type publicKey PublicKey
	var v publicKey
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*p = PublicKey(v)
	return nil
}

func (p *PublicKeyCollection) UnmarshalJSON(data []byte) error {
	type publicKeys PublicKeyCollection
	var v publicKeys
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*p = PublicKeyCollection(v)
	return nil
}

func (p *PublicKeyCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*p))
	for i, v := range *p {
		ret[i] = v
	}

	return &ret
}
