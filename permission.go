package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/lib"
)

type Permission struct {
	Id         int64  `json:"id,omitempty"`
	Path       string `json:"path,omitempty"`
	UserId     int64  `json:"user_id,omitempty"`
	Username   string `json:"username,omitempty"`
	GroupId    int64  `json:"group_id,omitempty"`
	GroupName  string `json:"group_name,omitempty"`
	Permission string `json:"permission,omitempty"`
	Recursive  *bool  `json:"recursive,omitempty"`
}

type PermissionCollection []Permission

type PermissionListParams struct {
	Page          int             `url:"page,omitempty"`
	PerPage       int             `url:"per_page,omitempty"`
	Action        string          `url:"action,omitempty"`
	Cursor        string          `url:"cursor,omitempty"`
	SortBy        json.RawMessage `url:"sort_by,omitempty"`
	Filter        json.RawMessage `url:"filter,omitempty"`
	FilterGt      json.RawMessage `url:"filter_gt,omitempty"`
	FilterGteq    json.RawMessage `url:"filter_gteq,omitempty"`
	FilterLike    json.RawMessage `url:"filter_like,omitempty"`
	FilterLt      json.RawMessage `url:"filter_lt,omitempty"`
	FilterLteq    json.RawMessage `url:"filter_lteq,omitempty"`
	Path          string          `url:"path,omitempty"`
	GroupId       string          `url:"group_id,omitempty"`
	UserId        string          `url:"user_id,omitempty"`
	IncludeGroups *bool           `url:"include_groups,omitempty"`
	lib.ListParams
}

type PermissionCreateParams struct {
	GroupId    int64  `url:"group_id,omitempty"`
	Path       string `url:"path,omitempty"`
	Permission string `url:"permission,omitempty"`
	Recursive  *bool  `url:"recursive,omitempty"`
	UserId     int64  `url:"user_id,omitempty"`
	Username   string `url:"username,omitempty"`
}

type PermissionDeleteParams struct {
	Id int64 `url:"-,omitempty"`
}

func (p *Permission) UnmarshalJSON(data []byte) error {
	type permission Permission
	var v permission
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*p = Permission(v)
	return nil
}

func (p *PermissionCollection) UnmarshalJSON(data []byte) error {
	type permissions []Permission
	var v permissions
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*p = PermissionCollection(v)
	return nil
}
