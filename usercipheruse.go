package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/lib"
)

type UserCipherUse struct {
	Id             int64     `json:"id,omitempty"`
	ProtocolCipher string    `json:"protocol_cipher,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
	Interface      string    `json:"interface,omitempty"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
	UserId         int64     `json:"user_id,omitempty"`
}

type UserCipherUseCollection []UserCipherUse

type UserCipherUseListParams struct {
	UserId  int64  `url:"user_id,omitempty" required:"false"`
	Cursor  string `url:"cursor,omitempty" required:"false"`
	PerPage int    `url:"per_page,omitempty" required:"false"`
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
