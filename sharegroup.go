package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// ShareGroup is a Files.com API resource.
type ShareGroup struct {
	Id      int64              `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	Name    string             `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	Notes   string             `json:"notes,omitempty" path:"notes,omitempty" url:"notes,omitempty"`
	UserId  int64              `json:"user_id,omitempty" path:"user_id,omitempty" url:"user_id,omitempty"`
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
	UserId int64 `url:"user_id,omitempty" json:"user_id,omitempty" path:"user_id"`
	ListParams
}

// ShareGroupFindParams contains the request parameters for GET /share_groups/{id}.
type ShareGroupFindParams struct {
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// ShareGroupCreateParams contains the request parameters for POST /share_groups.
type ShareGroupCreateParams struct {
	UserId  int64                    `url:"user_id,omitempty" json:"user_id,omitempty" path:"user_id"`
	Notes   string                   `url:"notes,omitempty" json:"notes,omitempty" path:"notes"`
	Name    string                   `url:"name" json:"name" path:"name"`
	Members []map[string]interface{} `url:"members" json:"members" path:"members"`
}

// ShareGroupUpdateParams contains the request parameters for PATCH /share_groups/{id}.
type ShareGroupUpdateParams struct {
	Id      int64                    `url:"-,omitempty" json:"-,omitempty" path:"id"`
	Notes   string                   `url:"notes,omitempty" json:"notes,omitempty" path:"notes"`
	Name    string                   `url:"name,omitempty" json:"name,omitempty" path:"name"`
	Members []map[string]interface{} `url:"members,omitempty" json:"members,omitempty" path:"members"`
}

// ShareGroupDeleteParams contains the request parameters for DELETE /share_groups/{id}.
type ShareGroupDeleteParams struct {
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
