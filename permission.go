package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// Permission is a Files.com API resource.
type Permission struct {
	// Permission ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Folder path. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.
	Path string `json:"path,omitempty" path:"path,omitempty" url:"path,omitempty"`
	// User ID.
	UserId int64 `json:"user_id,omitempty" path:"user_id,omitempty" url:"user_id,omitempty"`
	// Username (if applicable).
	Username string `json:"username,omitempty" path:"username,omitempty" url:"username,omitempty"`
	// Group ID.
	GroupId int64 `json:"group_id,omitempty" path:"group_id,omitempty" url:"group_id,omitempty"`
	// Group name (if applicable).
	GroupName string `json:"group_name,omitempty" path:"group_name,omitempty" url:"group_name,omitempty"`
	// Group IDs when this Permission requires membership in every listed group.
	GroupIds []int64 `json:"group_ids,omitempty" path:"group_ids,omitempty" url:"group_ids,omitempty"`
	// Group names when this permission requires multiple groups.
	GroupNames []string `json:"group_names,omitempty" path:"group_names,omitempty" url:"group_names,omitempty"`
	// Partner ID (if applicable).
	PartnerId int64 `json:"partner_id,omitempty" path:"partner_id,omitempty" url:"partner_id,omitempty"`
	// Partner name (if applicable).
	PartnerName string `json:"partner_name,omitempty" path:"partner_name,omitempty" url:"partner_name,omitempty"`
	// Permission type. See the table referenced in the documentation for an explanation of each permission.
	Permission string `json:"permission,omitempty" path:"permission,omitempty" url:"permission,omitempty"`
	// Whether this Permission applies to subfolders. Must be true for `admin` Permissions.
	Recursive *bool `json:"recursive,omitempty" path:"recursive,omitempty" url:"recursive,omitempty"`
	// Site ID.
	SiteId int64 `json:"site_id,omitempty" path:"site_id,omitempty" url:"site_id,omitempty"`
}

// Identifier returns the resource ID.
func (p Permission) Identifier() interface{} {
	return p.Id
}

// PermissionCollection is a list of Permission resources.
type PermissionCollection []Permission

// PermissionListParams contains the request parameters for GET /permissions.
type PermissionListParams struct {
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `site_id`, `group_id`, `path`, `user_id`, `partner_id` or `id`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `path`, `group_id`, `partner_id` or `user_id`. Valid field combinations are `[ group_id, path ]`, `[ partner_id, path ]`, `[ user_id, path ]`, `[ user_id, group_id ]`, `[ user_id, group_id, path ]`, `[ user_id, group_id, partner_id ]` or `[ user_id, group_id, partner_id, path ]`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// If set, return records where the specified field is prefixed by the supplied value. Valid fields are `path`.
	FilterPrefix interface{} `url:"filter_prefix,omitempty" json:"filter_prefix,omitempty" path:"filter_prefix"`
	// Permission path. If provided, will scope all permissions(including upward) to this path.
	Path string `url:"path,omitempty" json:"path,omitempty" path:"path"`
	// If searching by user or group, also include user's permissions that are inherited from its groups?
	IncludeGroups *bool  `url:"include_groups,omitempty" json:"include_groups,omitempty" path:"include_groups"`
	GroupId       string `url:"group_id,omitempty" json:"group_id,omitempty" path:"group_id"`
	PartnerId     string `url:"partner_id,omitempty" json:"partner_id,omitempty" path:"partner_id"`
	UserId        string `url:"user_id,omitempty" json:"user_id,omitempty" path:"user_id"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// PermissionCreateParams contains the request parameters for POST /permissions.
type PermissionCreateParams struct {
	// Folder path. Required.
	Path string `url:"path" json:"path" path:"path"`
	// Group ID. Provide `group_name` or `group_id`.
	GroupId int64 `url:"group_id,omitempty" json:"group_id,omitempty" path:"group_id"`
	// Group IDs when the Permission requires membership in every listed group. If sent as a string, it should be comma-delimited.
	GroupIds string `url:"group_ids,omitempty" json:"group_ids,omitempty" path:"group_ids"`
	// Permission type. Can be `admin`, `full`, `readonly`, `writeonly`, `list`, or `history`.
	Permission string `url:"permission,omitempty" json:"permission,omitempty" path:"permission"`
	// Apply to subfolders recursively? Must be true for `admin` Permissions.
	Recursive *bool `url:"recursive,omitempty" json:"recursive,omitempty" path:"recursive"`
	// Partner ID if this Permission belongs to a partner.
	PartnerId int64 `url:"partner_id,omitempty" json:"partner_id,omitempty" path:"partner_id"`
	// User ID. Provide `username` or `user_id`.
	UserId int64 `url:"user_id,omitempty" json:"user_id,omitempty" path:"user_id"`
	// User username. Provide `username` or `user_id`.
	Username string `url:"username,omitempty" json:"username,omitempty" path:"username"`
	// Group name. Provide `group_name` or `group_id`.
	GroupName string `url:"group_name,omitempty" json:"group_name,omitempty" path:"group_name"`
	// Site ID. If not provided, will default to current site. Used when creating a permission for a child site.
	SiteId int64 `url:"site_id,omitempty" json:"site_id,omitempty" path:"site_id"`
}

// PermissionDeleteParams contains the request parameters for DELETE /permissions/{id}.
type PermissionDeleteParams struct {
	// Permission ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (p *Permission) UnmarshalJSON(data []byte) error {
	type permission Permission
	var v permission
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*p = Permission(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (p *PermissionCollection) UnmarshalJSON(data []byte) error {
	type permissions PermissionCollection
	var v permissions
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*p = PermissionCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (p *PermissionCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*p))
	for i, v := range *p {
		ret[i] = v
	}

	return &ret
}
