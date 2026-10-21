package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// GroupUser is a Files.com API resource.
type GroupUser struct {
	// Group name.
	GroupName string `json:"group_name,omitempty" path:"group_name,omitempty" url:"group_name,omitempty"`
	// Group ID.
	GroupId int64 `json:"group_id,omitempty" path:"group_id,omitempty" url:"group_id,omitempty"`
	// User ID.
	UserId int64 `json:"user_id,omitempty" path:"user_id,omitempty" url:"user_id,omitempty"`
	// Is this user an administrator of this group?
	Admin *bool `json:"admin,omitempty" path:"admin,omitempty" url:"admin,omitempty"`
	// Username of the user.
	Username string `json:"username,omitempty" path:"username,omitempty" url:"username,omitempty"`
	// Group User ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
}

// Identifier returns the resource ID.
func (g GroupUser) Identifier() interface{} {
	return g.Id
}

// GroupUserCollection is a list of GroupUser resources.
type GroupUserCollection []GroupUser

// GroupUserListParams contains the request parameters for GET /group_users.
type GroupUserListParams struct {
	// Group ID. If provided, returns memberships of this group. Requires a Site Administrator, a Read-only Administrator, a Workspace Administrator for the group's workspace, or a Group Administrator of this group.
	GroupId int64 `url:"group_id,omitempty" json:"group_id,omitempty" path:"group_id"`
	// User ID. If provided, will return group_users of this user.
	UserId int64 `url:"user_id,omitempty" json:"user_id,omitempty" path:"user_id"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// GroupUserCreateParams contains the request parameters for POST /group_users.
type GroupUserCreateParams struct {
	// Group ID to add user to. Required.
	GroupId int64 `url:"group_id" json:"group_id" path:"group_id"`
	// User ID to add to group. Required.
	UserId int64 `url:"user_id" json:"user_id" path:"user_id"`
	// Is the user a group administrator?
	Admin *bool `url:"admin,omitempty" json:"admin,omitempty" path:"admin"`
}

// GroupUserUpdateParams contains the request parameters for PATCH /group_users/{id}.
type GroupUserUpdateParams struct {
	// Group User ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// Group ID to add user to. Required.
	GroupId int64 `url:"group_id" json:"group_id" path:"group_id"`
	// User ID to add to group. Required.
	UserId int64 `url:"user_id" json:"user_id" path:"user_id"`
	// Is the user a group administrator?
	Admin *bool `url:"admin,omitempty" json:"admin,omitempty" path:"admin"`
}

// GroupUserDeleteParams contains the request parameters for DELETE /group_users/{id}.
type GroupUserDeleteParams struct {
	// Group User ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// Group ID from which to remove user. Required.
	GroupId int64 `url:"group_id" json:"group_id" path:"group_id"`
	// User ID to remove from group. Required.
	UserId int64 `url:"user_id" json:"user_id" path:"user_id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (g *GroupUser) UnmarshalJSON(data []byte) error {
	type groupUser GroupUser
	var v groupUser
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*g = GroupUser(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (g *GroupUserCollection) UnmarshalJSON(data []byte) error {
	type groupUsers GroupUserCollection
	var v groupUsers
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*g = GroupUserCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (g *GroupUserCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*g))
	for i, v := range *g {
		ret[i] = v
	}

	return &ret
}
