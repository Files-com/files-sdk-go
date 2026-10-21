package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// IntegrationCentricProfile is a Files.com API resource.
type IntegrationCentricProfile struct {
	// Integration Centric Profile ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Profile name.
	Name string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	// Workspace ID.
	WorkspaceId int64 `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	// Whether this profile applies to all users in the Workspace by default.
	UseForAllUsers *bool `json:"use_for_all_users,omitempty" path:"use_for_all_users,omitempty" url:"use_for_all_users,omitempty"`
	// Remote Server integrations the user is expected to add and connect. Each entry requires `server_type` and may include a display `name`.
	ExpectedRemoteServers []map[string]interface{} `json:"expected_remote_servers,omitempty" path:"expected_remote_servers,omitempty" url:"expected_remote_servers,omitempty"`
}

// Identifier returns the resource ID.
func (i IntegrationCentricProfile) Identifier() interface{} {
	return i.Id
}

// IntegrationCentricProfileCollection is a list of IntegrationCentricProfile resources.
type IntegrationCentricProfileCollection []IntegrationCentricProfile

// IntegrationCentricProfileListParams contains the request parameters for GET /integration_centric_profiles.
type IntegrationCentricProfileListParams struct {
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id` and `name`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `workspace_id`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// IntegrationCentricProfileFindParams contains the request parameters for GET /integration_centric_profiles/{id}.
type IntegrationCentricProfileFindParams struct {
	// Integration Centric Profile ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// IntegrationCentricProfileCreateParams contains the request parameters for POST /integration_centric_profiles.
type IntegrationCentricProfileCreateParams struct {
	// Profile name. Required.
	Name string `url:"name" json:"name" path:"name"`
	// Remote Server integrations the user is expected to add and connect. Each entry requires `server_type` and may include a display `name`. Required.
	ExpectedRemoteServers []map[string]interface{} `url:"expected_remote_servers" json:"expected_remote_servers" path:"expected_remote_servers"`
	// Workspace ID.
	WorkspaceId int64 `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
	// Whether this profile applies to all users in the Workspace by default.
	UseForAllUsers *bool `url:"use_for_all_users,omitempty" json:"use_for_all_users,omitempty" path:"use_for_all_users"`
}

// IntegrationCentricProfileUpdateParams contains the request parameters for PATCH /integration_centric_profiles/{id}.
type IntegrationCentricProfileUpdateParams struct {
	// Integration Centric Profile ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// Profile name.
	Name string `url:"name,omitempty" json:"name,omitempty" path:"name"`
	// Workspace ID.
	WorkspaceId int64 `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
	// Remote Server integrations the user is expected to add and connect. Each entry requires `server_type` and may include a display `name`.
	ExpectedRemoteServers []map[string]interface{} `url:"expected_remote_servers,omitempty" json:"expected_remote_servers,omitempty" path:"expected_remote_servers"`
	// Whether this profile applies to all users in the Workspace by default.
	UseForAllUsers *bool `url:"use_for_all_users,omitempty" json:"use_for_all_users,omitempty" path:"use_for_all_users"`
}

// IntegrationCentricProfileDeleteParams contains the request parameters for DELETE /integration_centric_profiles/{id}.
type IntegrationCentricProfileDeleteParams struct {
	// Integration Centric Profile ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (i *IntegrationCentricProfile) UnmarshalJSON(data []byte) error {
	type integrationCentricProfile IntegrationCentricProfile
	var v integrationCentricProfile
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*i = IntegrationCentricProfile(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (i *IntegrationCentricProfileCollection) UnmarshalJSON(data []byte) error {
	type integrationCentricProfiles IntegrationCentricProfileCollection
	var v integrationCentricProfiles
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*i = IntegrationCentricProfileCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (i *IntegrationCentricProfileCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*i))
	for i, v := range *i {
		ret[i] = v
	}

	return &ret
}
