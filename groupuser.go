package files_sdk

import (
  lib "github.com/Files-com/files-sdk-go/lib"
  "encoding/json"
)

type GroupUser struct {
  GroupName string `json:"group_name,omitempty"`
  GroupId int `json:"group_id,omitempty"`
  UserId int `json:"user_id,omitempty"`
  Admin *bool `json:"admin,omitempty"`
  Usernames []string `json:"usernames,omitempty"`
  Id int `json:"id,omitempty"`
}

type GroupUserCollection []GroupUser

type GroupUserListParams struct {
  UserId int `url:"user_id,omitempty"`
  Page int `url:"page,omitempty"`
  PerPage int `url:"per_page,omitempty"`
  Action string `url:"action,omitempty"`
  GroupId int `url:"group_id,omitempty"`
  lib.ListParams
}

type GroupUserUpdateParams struct {
  Id int `url:"-,omitempty"`
  GroupId int `url:"group_id,omitempty"`
  UserId int `url:"user_id,omitempty"`
  Admin *bool `url:"admin,omitempty"`
}

type GroupUserDeleteParams struct {
  Id int `url:"-,omitempty"`
  GroupId int `url:"group_id,omitempty"`
  UserId int `url:"user_id,omitempty"`
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

