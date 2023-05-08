package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type GroupUser struct {
	GroupName string   `json:"group_name,omitempty" path:"group_name"`
	GroupId   int64    `json:"group_id,omitempty" path:"group_id"`
	UserId    int64    `json:"user_id,omitempty" path:"user_id"`
	Admin     *bool    `json:"admin,omitempty" path:"admin"`
	Usernames []string `json:"usernames,omitempty" path:"usernames"`
	Id        int64    `json:"id,omitempty" path:"id"`
}

func (g GroupUser) Identifier() interface{} {
	return g.Id
}

type GroupUserCollection []GroupUser

type GroupUserListParams struct {
	UserId  int64 `url:"user_id,omitempty" required:"false" json:"user_id,omitempty" path:"user_id"`
	GroupId int64 `url:"group_id,omitempty" required:"false" json:"group_id,omitempty" path:"group_id"`
	ListParams
}

type GroupUserCreateParams struct {
	GroupId int64 `url:"group_id,omitempty" required:"true" json:"group_id,omitempty" path:"group_id"`
	UserId  int64 `url:"user_id,omitempty" required:"true" json:"user_id,omitempty" path:"user_id"`
	Admin   *bool `url:"admin,omitempty" required:"false" json:"admin,omitempty" path:"admin"`
}

type GroupUserUpdateParams struct {
	Id      int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
	GroupId int64 `url:"group_id,omitempty" required:"true" json:"group_id,omitempty" path:"group_id"`
	UserId  int64 `url:"user_id,omitempty" required:"true" json:"user_id,omitempty" path:"user_id"`
	Admin   *bool `url:"admin,omitempty" required:"false" json:"admin,omitempty" path:"admin"`
}

type GroupUserDeleteParams struct {
	Id      int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
	GroupId int64 `url:"group_id,omitempty" required:"true" json:"group_id,omitempty" path:"group_id"`
	UserId  int64 `url:"user_id,omitempty" required:"true" json:"user_id,omitempty" path:"user_id"`
}

func (g *GroupUser) UnmarshalJSON(data []byte) error {
	type groupUser GroupUser
	var v groupUser
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*g = GroupUser(v)
	return nil
}

func (g *GroupUserCollection) UnmarshalJSON(data []byte) error {
	type groupUsers GroupUserCollection
	var v groupUsers
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*g = GroupUserCollection(v)
	return nil
}

func (g *GroupUserCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*g))
	for i, v := range *g {
		ret[i] = v
	}

	return &ret
}
