package files_sdk

import (
  lib "github.com/Files-com/files-sdk-go/lib"
  "encoding/json"
  "time"
)

type PublicKey struct {
  Id int `json:"id,omitempty"`
  Title string `json:"title,omitempty"`
  CreatedAt time.Time `json:"created_at,omitempty"`
  Fingerprint string `json:"fingerprint,omitempty"`
  UserId int `json:"user_id,omitempty"`
  PublicKey string `json:"public_key,omitempty"`
}

type PublicKeyCollection []PublicKey

type PublicKeyListParams struct {
  UserId int `url:"user_id,omitempty"`
  Page int `url:"page,omitempty"`
  PerPage int `url:"per_page,omitempty"`
  Action string `url:"action,omitempty"`
  lib.ListParams
}

type PublicKeyFindParams struct {
  Id int `url:"-,omitempty"`
}

type PublicKeyCreateParams struct {
  UserId int `url:"user_id,omitempty"`
  Title string `url:"title,omitempty"`
  PublicKey string `url:"public_key,omitempty"`
}

type PublicKeyUpdateParams struct {
  Id int `url:"-,omitempty"`
  Title string `url:"title,omitempty"`
}

type PublicKeyDeleteParams struct {
  Id int `url:"-,omitempty"`
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

