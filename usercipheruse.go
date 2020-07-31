package files_sdk

import (
	"encoding/json"
	lib "github.com/Files-com/files-sdk-go/lib"
	"time"
)

type UserCipherUse struct {
	Id             int       `json:"id,omitempty"`
	ProtocolCipher string    `json:"protocol_cipher,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
	Interface      string    `json:"interface,omitempty"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
	UserId         int       `json:"user_id,omitempty"`
}

type UserCipherUseCollection []UserCipherUse

type UserCipherUseListParams struct {
	UserId  int    `url:"user_id,omitempty"`
	Page    int    `url:"page,omitempty"`
	PerPage int    `url:"per_page,omitempty"`
	Action  string `url:"action,omitempty"`
	lib.ListParams
}

func (u *UserCipherUse) UnmarshalJSON(data []byte) error {
	type userCipherUse UserCipherUse
	var v userCipherUse
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*u = UserCipherUse(v)
	return nil
}

func (u *UserCipherUseCollection) UnmarshalJSON(data []byte) error {
	type userCipherUses []UserCipherUse
	var v userCipherUses
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*u = UserCipherUseCollection(v)
	return nil
}
