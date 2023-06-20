package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Permission struct {
	Id         int64  `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	Path       string `json:"path,omitempty" path:"path,omitempty" url:"path,omitempty"`
	UserId     int64  `json:"user_id,omitempty" path:"user_id,omitempty" url:"user_id,omitempty"`
	Username   string `json:"username,omitempty" path:"username,omitempty" url:"username,omitempty"`
	GroupId    int64  `json:"group_id,omitempty" path:"group_id,omitempty" url:"group_id,omitempty"`
	GroupName  string `json:"group_name,omitempty" path:"group_name,omitempty" url:"group_name,omitempty"`
	Permission string `json:"permission,omitempty" path:"permission,omitempty" url:"permission,omitempty"`
	Recursive  *bool  `json:"recursive,omitempty" path:"recursive,omitempty" url:"recursive,omitempty"`
}

func (p Permission) Identifier() interface{} {
	return p.Id
}

type PermissionCollection []Permission

type PermissionListParams struct {
	SortBy        map[string]interface{} `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty" path:"sort_by"`
	Filter        Permission             `url:"filter,omitempty" required:"false" json:"filter,omitempty" path:"filter"`
	FilterPrefix  map[string]interface{} `url:"filter_prefix,omitempty" required:"false" json:"filter_prefix,omitempty" path:"filter_prefix"`
	Path          string                 `url:"path,omitempty" required:"false" json:"path,omitempty" path:"path"`
	GroupId       string                 `url:"group_id,omitempty" required:"false" json:"group_id,omitempty" path:"group_id"`
	UserId        string                 `url:"user_id,omitempty" required:"false" json:"user_id,omitempty" path:"user_id"`
	IncludeGroups *bool                  `url:"include_groups,omitempty" required:"false" json:"include_groups,omitempty" path:"include_groups"`
	ListParams
}

type PermissionCreateParams struct {
	GroupId    int64  `url:"group_id,omitempty" required:"false" json:"group_id,omitempty" path:"group_id"`
	Path       string `url:"path,omitempty" required:"false" json:"path,omitempty" path:"path"`
	Permission string `url:"permission,omitempty" required:"false" json:"permission,omitempty" path:"permission"`
	Recursive  *bool  `url:"recursive,omitempty" required:"false" json:"recursive,omitempty" path:"recursive"`
	UserId     int64  `url:"user_id,omitempty" required:"false" json:"user_id,omitempty" path:"user_id"`
	Username   string `url:"username,omitempty" required:"false" json:"username,omitempty" path:"username"`
}

type PermissionDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
}

func (p *Permission) UnmarshalJSON(data []byte) error {
	type permission Permission
	var v permission
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*p = Permission(v)
	return nil
}

func (p *PermissionCollection) UnmarshalJSON(data []byte) error {
	type permissions PermissionCollection
	var v permissions
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
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
