package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

type Group struct {
	Id                int64  `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	Name              string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	AllowedIps        string `json:"allowed_ips,omitempty" path:"allowed_ips,omitempty" url:"allowed_ips,omitempty"`
	AdminIds          string `json:"admin_ids,omitempty" path:"admin_ids,omitempty" url:"admin_ids,omitempty"`
	Notes             string `json:"notes,omitempty" path:"notes,omitempty" url:"notes,omitempty"`
	UserIds           string `json:"user_ids,omitempty" path:"user_ids,omitempty" url:"user_ids,omitempty"`
	Usernames         string `json:"usernames,omitempty" path:"usernames,omitempty" url:"usernames,omitempty"`
	FtpPermission     *bool  `json:"ftp_permission,omitempty" path:"ftp_permission,omitempty" url:"ftp_permission,omitempty"`
	SftpPermission    *bool  `json:"sftp_permission,omitempty" path:"sftp_permission,omitempty" url:"sftp_permission,omitempty"`
	DavPermission     *bool  `json:"dav_permission,omitempty" path:"dav_permission,omitempty" url:"dav_permission,omitempty"`
	RestapiPermission *bool  `json:"restapi_permission,omitempty" path:"restapi_permission,omitempty" url:"restapi_permission,omitempty"`
}

func (g Group) Identifier() interface{} {
	return g.Id
}

type GroupCollection []Group

type GroupListParams struct {
	SortBy       map[string]interface{} `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty" path:"sort_by"`
	Filter       Group                  `url:"filter,omitempty" required:"false" json:"filter,omitempty" path:"filter"`
	FilterPrefix map[string]interface{} `url:"filter_prefix,omitempty" required:"false" json:"filter_prefix,omitempty" path:"filter_prefix"`
	Ids          string                 `url:"ids,omitempty" required:"false" json:"ids,omitempty" path:"ids"`
	ListParams
}

type GroupFindParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
}

type GroupCreateParams struct {
	Notes             string `url:"notes,omitempty" required:"false" json:"notes,omitempty" path:"notes"`
	UserIds           string `url:"user_ids,omitempty" required:"false" json:"user_ids,omitempty" path:"user_ids"`
	AdminIds          string `url:"admin_ids,omitempty" required:"false" json:"admin_ids,omitempty" path:"admin_ids"`
	FtpPermission     *bool  `url:"ftp_permission,omitempty" required:"false" json:"ftp_permission,omitempty" path:"ftp_permission"`
	SftpPermission    *bool  `url:"sftp_permission,omitempty" required:"false" json:"sftp_permission,omitempty" path:"sftp_permission"`
	DavPermission     *bool  `url:"dav_permission,omitempty" required:"false" json:"dav_permission,omitempty" path:"dav_permission"`
	RestapiPermission *bool  `url:"restapi_permission,omitempty" required:"false" json:"restapi_permission,omitempty" path:"restapi_permission"`
	AllowedIps        string `url:"allowed_ips,omitempty" required:"false" json:"allowed_ips,omitempty" path:"allowed_ips"`
	Name              string `url:"name,omitempty" required:"true" json:"name,omitempty" path:"name"`
}

type GroupUpdateParams struct {
	Id                int64  `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
	Notes             string `url:"notes,omitempty" required:"false" json:"notes,omitempty" path:"notes"`
	UserIds           string `url:"user_ids,omitempty" required:"false" json:"user_ids,omitempty" path:"user_ids"`
	AdminIds          string `url:"admin_ids,omitempty" required:"false" json:"admin_ids,omitempty" path:"admin_ids"`
	FtpPermission     *bool  `url:"ftp_permission,omitempty" required:"false" json:"ftp_permission,omitempty" path:"ftp_permission"`
	SftpPermission    *bool  `url:"sftp_permission,omitempty" required:"false" json:"sftp_permission,omitempty" path:"sftp_permission"`
	DavPermission     *bool  `url:"dav_permission,omitempty" required:"false" json:"dav_permission,omitempty" path:"dav_permission"`
	RestapiPermission *bool  `url:"restapi_permission,omitempty" required:"false" json:"restapi_permission,omitempty" path:"restapi_permission"`
	AllowedIps        string `url:"allowed_ips,omitempty" required:"false" json:"allowed_ips,omitempty" path:"allowed_ips"`
	Name              string `url:"name,omitempty" required:"false" json:"name,omitempty" path:"name"`
}

type GroupDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
}

func (g *Group) UnmarshalJSON(data []byte) error {
	type group Group
	var v group
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*g = Group(v)
	return nil
}

func (g *GroupCollection) UnmarshalJSON(data []byte) error {
	type groups GroupCollection
	var v groups
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*g = GroupCollection(v)
	return nil
}

func (g *GroupCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*g))
	for i, v := range *g {
		ret[i] = v
	}

	return &ret
}
