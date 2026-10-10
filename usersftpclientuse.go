package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// UserSftpClientUse is a Files.com API resource.
type UserSftpClientUse struct {
	Id         int64      `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	SftpClient string     `json:"sftp_client,omitempty" path:"sftp_client,omitempty" url:"sftp_client,omitempty"`
	CreatedAt  *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty" path:"updated_at,omitempty" url:"updated_at,omitempty"`
	UserId     int64      `json:"user_id,omitempty" path:"user_id,omitempty" url:"user_id,omitempty"`
}

// Identifier returns the resource ID.
func (u UserSftpClientUse) Identifier() interface{} {
	return u.Id
}

// UserSftpClientUseCollection is a list of UserSftpClientUse resources.
type UserSftpClientUseCollection []UserSftpClientUse

// UserSftpClientUseListParams contains the request parameters for GET /user_sftp_client_uses.
type UserSftpClientUseListParams struct {
	UserId int64       `url:"user_id,omitempty" json:"user_id,omitempty" path:"user_id"`
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	ListParams
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (u *UserSftpClientUse) UnmarshalJSON(data []byte) error {
	type userSftpClientUse UserSftpClientUse
	var v userSftpClientUse
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*u = UserSftpClientUse(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (u *UserSftpClientUseCollection) UnmarshalJSON(data []byte) error {
	type userSftpClientUses UserSftpClientUseCollection
	var v userSftpClientUses
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*u = UserSftpClientUseCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (u *UserSftpClientUseCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*u))
	for i, v := range *u {
		ret[i] = v
	}

	return &ret
}
