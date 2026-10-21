package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// ShareGroup is a Files.com API resource.
type ShareGroup struct {
	// Share Group ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Name of the share group.
	Name string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	// Additional notes of the share group.
	Notes string `json:"notes,omitempty" path:"notes,omitempty" url:"notes,omitempty"`
	// Owner User ID.
	UserId int64 `json:"user_id,omitempty" path:"user_id,omitempty" url:"user_id,omitempty"`
	// A list of share group members.
	Members []ShareGroupMember `json:"members,omitempty" path:"members,omitempty" url:"members,omitempty"`
}

// Identifier returns the resource ID.
func (s ShareGroup) Identifier() interface{} {
	return s.Id
}

// ShareGroupCollection is a list of ShareGroup resources.
type ShareGroupCollection []ShareGroup

// ShareGroupListParams contains the request parameters for GET /share_groups.
type ShareGroupListParams struct {
	// User ID. Provide a value of `0` to operate the current session's user.
	UserId int64 `url:"user_id,omitempty" json:"user_id,omitempty" path:"user_id"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// ShareGroupFindParams contains the request parameters for GET /share_groups/{id}.
type ShareGroupFindParams struct {
	// Share Group ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// ShareGroupCreateParams contains the request parameters for POST /share_groups.
type ShareGroupCreateParams struct {
	// User ID. Provide a value of `0` to operate the current session's user.
	UserId int64 `url:"user_id,omitempty" json:"user_id,omitempty" path:"user_id"`
	// Additional notes of the share group.
	Notes string `url:"notes,omitempty" json:"notes,omitempty" path:"notes"`
	// Name of the share group. Required.
	Name string `url:"name" json:"name" path:"name"`
	// A list of share group members. Required.
	Members []map[string]interface{} `url:"members" json:"members" path:"members"`
}

// ShareGroupUpdateParams contains the request parameters for PATCH /share_groups/{id}.
type ShareGroupUpdateParams struct {
	// Share Group ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// Additional notes of the share group.
	Notes string `url:"notes,omitempty" json:"notes,omitempty" path:"notes"`
	// Name of the share group.
	Name string `url:"name,omitempty" json:"name,omitempty" path:"name"`
	// A list of share group members.
	Members []map[string]interface{} `url:"members,omitempty" json:"members,omitempty" path:"members"`
}

// ShareGroupDeleteParams contains the request parameters for DELETE /share_groups/{id}.
type ShareGroupDeleteParams struct {
	// Share Group ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (s *ShareGroup) UnmarshalJSON(data []byte) error {
	type shareGroup ShareGroup
	var v shareGroup
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*s = ShareGroup(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (s *ShareGroupCollection) UnmarshalJSON(data []byte) error {
	type shareGroups ShareGroupCollection
	var v shareGroups
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*s = ShareGroupCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (s *ShareGroupCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*s))
	for i, v := range *s {
		ret[i] = v
	}

	return &ret
}
