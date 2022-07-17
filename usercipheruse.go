package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type UserCipherUse struct {
	Id             int64      `json:"id,omitempty" path:"id"`
	ProtocolCipher string     `json:"protocol_cipher,omitempty" path:"protocol_cipher"`
	CreatedAt      *time.Time `json:"created_at,omitempty" path:"created_at"`
	Interface      string     `json:"interface,omitempty" path:"interface"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty" path:"updated_at"`
	UserId         int64      `json:"user_id,omitempty" path:"user_id"`
}

type UserCipherUseCollection []UserCipherUse

type UserCipherUseListParams struct {
	UserId int64 `url:"user_id,omitempty" required:"false" json:"user_id,omitempty" path:"user_id"`
	lib.ListParams
}

func (u *UserCipherUse) UnmarshalJSON(data []byte) error {
	type userCipherUse UserCipherUse
	var v userCipherUse
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*u = UserCipherUse(v)
	return nil
}

func (u *UserCipherUseCollection) UnmarshalJSON(data []byte) error {
	type userCipherUses UserCipherUseCollection
	var v userCipherUses
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*u = UserCipherUseCollection(v)
	return nil
}

func (u *UserCipherUseCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*u))
	for i, v := range *u {
		ret[i] = v
	}

	return &ret
}
