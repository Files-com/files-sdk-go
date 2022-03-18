package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
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
	Cursor        string          `url:"cursor,omitempty" required:"false" json:"cursor,omitempty"`
	PerPage       int64           `url:"per_page,omitempty" required:"false" json:"per_page,omitempty"`
	SortBy        json.RawMessage `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty"`
	Filter        json.RawMessage `url:"filter,omitempty" required:"false" json:"filter,omitempty"`
	FilterGt      json.RawMessage `url:"filter_gt,omitempty" required:"false" json:"filter_gt,omitempty"`
	FilterGteq    json.RawMessage `url:"filter_gteq,omitempty" required:"false" json:"filter_gteq,omitempty"`
	FilterLike    json.RawMessage `url:"filter_like,omitempty" required:"false" json:"filter_like,omitempty"`
	FilterLt      json.RawMessage `url:"filter_lt,omitempty" required:"false" json:"filter_lt,omitempty"`
	FilterLteq    json.RawMessage `url:"filter_lteq,omitempty" required:"false" json:"filter_lteq,omitempty"`
	Path          string          `url:"path,omitempty" required:"false" json:"path,omitempty"`
	GroupId       string          `url:"group_id,omitempty" required:"false" json:"group_id,omitempty"`
	UserId        string          `url:"user_id,omitempty" required:"false" json:"user_id,omitempty"`
	IncludeGroups *bool           `url:"include_groups,omitempty" required:"false" json:"include_groups,omitempty"`
	lib.ListParams
}

type PermissionCreateParams struct {
	GroupId    int64  `url:"group_id,omitempty" required:"false" json:"group_id,omitempty"`
	Path       string `url:"path,omitempty" required:"false" json:"path,omitempty"`
	Permission string `url:"permission,omitempty" required:"false" json:"permission,omitempty"`
	Recursive  *bool  `url:"recursive,omitempty" required:"false" json:"recursive,omitempty"`
	UserId     int64  `url:"user_id,omitempty" required:"false" json:"user_id,omitempty"`
	Username   string `url:"username,omitempty" required:"false" json:"username,omitempty"`
}

type PermissionDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"true" json:"-,omitempty"`
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

func (p *PermissionCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*p))
	for i, v := range *p {
		ret[i] = v
	}

	return &ret
}
