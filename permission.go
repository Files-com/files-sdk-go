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
	Page          int             `url:"page,omitempty" required:"false"`
	PerPage       int             `url:"per_page,omitempty" required:"false"`
	Action        string          `url:"action,omitempty" required:"false"`
	Cursor        string          `url:"cursor,omitempty" required:"false"`
	SortBy        json.RawMessage `url:"sort_by,omitempty" required:"false"`
	Filter        json.RawMessage `url:"filter,omitempty" required:"false"`
	FilterGt      json.RawMessage `url:"filter_gt,omitempty" required:"false"`
	FilterGteq    json.RawMessage `url:"filter_gteq,omitempty" required:"false"`
	FilterLike    json.RawMessage `url:"filter_like,omitempty" required:"false"`
	FilterLt      json.RawMessage `url:"filter_lt,omitempty" required:"false"`
	FilterLteq    json.RawMessage `url:"filter_lteq,omitempty" required:"false"`
	Path          string          `url:"path,omitempty" required:"false"`
	GroupId       string          `url:"group_id,omitempty" required:"false"`
	UserId        string          `url:"user_id,omitempty" required:"false"`
	IncludeGroups *bool           `url:"include_groups,omitempty" required:"false"`
	lib.ListParams
}

type PermissionCreateParams struct {
	GroupId    int64  `url:"group_id,omitempty" required:"false"`
	Path       string `url:"path,omitempty" required:"false"`
	Permission string `url:"permission,omitempty" required:"false"`
	Recursive  *bool  `url:"recursive,omitempty" required:"false"`
	UserId     int64  `url:"user_id,omitempty" required:"false"`
	Username   string `url:"username,omitempty" required:"false"`
}

type PermissionDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"true"`
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
