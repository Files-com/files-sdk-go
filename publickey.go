package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type PublicKey struct {
	Id          int64      `json:"id,omitempty" path:"id"`
	Title       string     `json:"title,omitempty" path:"title"`
	CreatedAt   *time.Time `json:"created_at,omitempty" path:"created_at"`
	Fingerprint string     `json:"fingerprint,omitempty" path:"fingerprint"`
	UserId      int64      `json:"user_id,omitempty" path:"user_id"`
	PublicKey   string     `json:"public_key,omitempty" path:"public_key"`
}

type PublicKeyCollection []PublicKey

type PublicKeyListParams struct {
	UserId int64 `url:"user_id,omitempty" required:"false" json:"user_id,omitempty" path:"user_id"`
	lib.ListParams
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
