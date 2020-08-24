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
	UserId  int64  `url:"user_id,omitempty"`
	Page    int    `url:"page,omitempty"`
	PerPage int    `url:"per_page,omitempty"`
	Action  string `url:"action,omitempty"`
	Cursor  string `url:"cursor,omitempty"`
	lib.ListParams
}

type PublicKeyFindParams struct {
	Id int64 `url:"-,omitempty"`
}

type PublicKeyCreateParams struct {
	UserId    int64  `url:"user_id,omitempty"`
	Title     string `url:"title,omitempty"`
	PublicKey string `url:"public_key,omitempty"`
}

type PublicKeyUpdateParams struct {
	Id    int64  `url:"-,omitempty"`
	Title string `url:"title,omitempty"`
}

type PublicKeyDeleteParams struct {
	Id int64 `url:"-,omitempty"`
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
