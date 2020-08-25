package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/lib"
)

type GroupUser struct {
	GroupName string   `json:"group_name,omitempty"`
	GroupId   int64    `json:"group_id,omitempty"`
	UserId    int64    `json:"user_id,omitempty"`
	Admin     *bool    `json:"admin,omitempty"`
	Usernames []string `json:"usernames,omitempty"`
	Id        int64    `json:"id,omitempty"`
}

type GroupUserCollection []GroupUser

type GroupUserListParams struct {
	UserId  int64  `url:"user_id,omitempty" required:"false"`
	Page    int    `url:"page,omitempty" required:"false"`
	PerPage int    `url:"per_page,omitempty" required:"false"`
	Action  string `url:"action,omitempty" required:"false"`
	Cursor  string `url:"cursor,omitempty" required:"false"`
	GroupId int64  `url:"group_id,omitempty" required:"false"`
	lib.ListParams
}

type GroupUserUpdateParams struct {
	Id      int64 `url:"-,omitempty" required:"true"`
	GroupId int64 `url:"group_id,omitempty" required:"true"`
	UserId  int64 `url:"user_id,omitempty" required:"true"`
	Admin   *bool `url:"admin,omitempty" required:"false"`
}

type GroupUserDeleteParams struct {
	Id      int64 `url:"-,omitempty" required:"true"`
	GroupId int64 `url:"group_id,omitempty" required:"true"`
	UserId  int64 `url:"user_id,omitempty" required:"true"`
}

func (g *GroupUser) UnmarshalJSON(data []byte) error {
	type groupUser GroupUser
	var v groupUser
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*g = GroupUser(v)
	return nil
}

func (g *GroupUserCollection) UnmarshalJSON(data []byte) error {
	type groupUsers []GroupUser
	var v groupUsers
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*g = GroupUserCollection(v)
	return nil
}
