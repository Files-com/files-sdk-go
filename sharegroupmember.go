package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// ShareGroupMember is a Files.com API resource.
type ShareGroupMember struct {
	Name    string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	Company string `json:"company,omitempty" path:"company,omitempty" url:"company,omitempty"`
	Email   string `json:"email,omitempty" path:"email,omitempty" url:"email,omitempty"`
}

// Identifier no path or id

// ShareGroupMemberCollection is a list of ShareGroupMember resources.
type ShareGroupMemberCollection []ShareGroupMember

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (s *ShareGroupMember) UnmarshalJSON(data []byte) error {
	type shareGroupMember ShareGroupMember
	var v shareGroupMember
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*s = ShareGroupMember(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (s *ShareGroupMemberCollection) UnmarshalJSON(data []byte) error {
	type shareGroupMembers ShareGroupMemberCollection
	var v shareGroupMembers
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*s = ShareGroupMemberCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (s *ShareGroupMemberCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*s))
	for i, v := range *s {
		ret[i] = v
	}

	return &ret
}
