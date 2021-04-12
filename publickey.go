package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/lib"
)

type PublicKey struct {
	Id          int64     `json:"id,omitempty"`
	Title       string    `json:"title,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	Fingerprint string    `json:"fingerprint,omitempty"`
	UserId      int64     `json:"user_id,omitempty"`
	PublicKey   string    `json:"public_key,omitempty"`
}

type PublicKeyCollection []PublicKey

type PublicKeyListParams struct {
	UserId  int64  `url:"user_id,omitempty" required:"false"`
	Cursor  string `url:"cursor,omitempty" required:"false"`
	PerPage int64  `url:"per_page,omitempty" required:"false"`
	lib.ListParams
}

type PublicKeyFindParams struct {
	Id int64 `url:"-,omitempty" required:"true"`
}

type PublicKeyCreateParams struct {
	UserId    int64  `url:"user_id,omitempty" required:"false"`
	Title     string `url:"title,omitempty" required:"true"`
	PublicKey string `url:"public_key,omitempty" required:"true"`
}

type PublicKeyUpdateParams struct {
	Id    int64  `url:"-,omitempty" required:"true"`
	Title string `url:"title,omitempty" required:"true"`
}

type PublicKeyDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"true"`
}

func (p *PublicKey) UnmarshalJSON(data []byte) error {
	type publicKey PublicKey
	var v publicKey
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*p = PublicKey(v)
	return nil
}

func (p *PublicKeyCollection) UnmarshalJSON(data []byte) error {
	type publicKeys []PublicKey
	var v publicKeys
	if err := json.Unmarshal(data, &v); err != nil {
		return err
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
