package files_sdk

import (
	"encoding/json"
	lib "github.com/Files-com/files-sdk-go/lib"
	"time"
)

type As2Key struct {
	Id                 int64     `json:"id,omitempty"`
	As2PartnershipName string    `json:"as2_partnership_name,omitempty"`
	CreatedAt          time.Time `json:"created_at,omitempty"`
	Fingerprint        string    `json:"fingerprint,omitempty"`
	UserId             int64     `json:"user_id,omitempty"`
	PublicKey          string    `json:"public_key,omitempty"`
}

type As2KeyCollection []As2Key

type As2KeyListParams struct {
	UserId  int64  `url:"user_id,omitempty"`
	Page    int    `url:"page,omitempty"`
	PerPage int    `url:"per_page,omitempty"`
	Action  string `url:"action,omitempty"`
	lib.ListParams
}

type As2KeyFindParams struct {
	Id int64 `url:"-,omitempty"`
}

type As2KeyCreateParams struct {
	UserId             int64  `url:"user_id,omitempty"`
	As2PartnershipName string `url:"as2_partnership_name,omitempty"`
	PublicKey          string `url:"public_key,omitempty"`
}

type As2KeyUpdateParams struct {
	Id                 int64  `url:"-,omitempty"`
	As2PartnershipName string `url:"as2_partnership_name,omitempty"`
}

type As2KeyDeleteParams struct {
	Id int64 `url:"-,omitempty"`
}

func (a *As2Key) UnmarshalJSON(data []byte) error {
	type as2Key As2Key
	var v as2Key
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*a = As2Key(v)
	return nil
}

func (a *As2KeyCollection) UnmarshalJSON(data []byte) error {
	type as2Keys []As2Key
	var v as2Keys
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*a = As2KeyCollection(v)
	return nil
}
