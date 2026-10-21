package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// DesktopConfigurationProfile is a Files.com API resource.
type DesktopConfigurationProfile struct {
	// Desktop Configuration Profile ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Profile name.
	Name string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	// Workspace ID.
	WorkspaceId int64 `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	// Whether this profile applies to all users in the Workspace by default.
	UseForAllUsers *bool `json:"use_for_all_users,omitempty" path:"use_for_all_users,omitempty" url:"use_for_all_users,omitempty"`
	// Whether the desktop app should hide drive mounting, prevent new drive mounts, and unmount active drive mounts for users with this profile.
	DisableDriveMounting *bool `json:"disable_drive_mounting,omitempty" path:"disable_drive_mounting,omitempty" url:"disable_drive_mounting,omitempty"`
	// Mount point mappings for the desktop app. Keys must be a single uppercase Windows drive letter other than A, B, or C, and values are Files.com paths to mount there.
	MountMappings interface{} `json:"mount_mappings,omitempty" path:"mount_mappings,omitempty" url:"mount_mappings,omitempty"`
}

// Identifier returns the resource ID.
func (d DesktopConfigurationProfile) Identifier() interface{} {
	return d.Id
}

// DesktopConfigurationProfileCollection is a list of DesktopConfigurationProfile resources.
type DesktopConfigurationProfileCollection []DesktopConfigurationProfile

// DesktopConfigurationProfileListParams contains the request parameters for GET /desktop_configuration_profiles.
type DesktopConfigurationProfileListParams struct {
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id` and `name`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `workspace_id`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// DesktopConfigurationProfileFindParams contains the request parameters for GET /desktop_configuration_profiles/{id}.
type DesktopConfigurationProfileFindParams struct {
	// Desktop Configuration Profile ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// DesktopConfigurationProfileCreateParams contains the request parameters for POST /desktop_configuration_profiles.
type DesktopConfigurationProfileCreateParams struct {
	// Profile name. Required.
	Name string `url:"name" json:"name" path:"name"`
	// Mount point mappings for the desktop app. Keys must be a single uppercase Windows drive letter other than A, B, or C, and values are Files.com paths to mount there. Required.
	MountMappings interface{} `url:"mount_mappings" json:"mount_mappings" path:"mount_mappings"`
	// Workspace ID.
	WorkspaceId int64 `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
	// Whether this profile applies to all users in the Workspace by default.
	UseForAllUsers *bool `url:"use_for_all_users,omitempty" json:"use_for_all_users,omitempty" path:"use_for_all_users"`
	// Whether the desktop app should hide drive mounting, prevent new drive mounts, and unmount active drive mounts for users with this profile.
	DisableDriveMounting *bool `url:"disable_drive_mounting,omitempty" json:"disable_drive_mounting,omitempty" path:"disable_drive_mounting"`
}

// DesktopConfigurationProfileUpdateParams contains the request parameters for PATCH /desktop_configuration_profiles/{id}.
type DesktopConfigurationProfileUpdateParams struct {
	// Desktop Configuration Profile ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// Profile name.
	Name string `url:"name,omitempty" json:"name,omitempty" path:"name"`
	// Workspace ID.
	WorkspaceId int64 `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
	// Mount point mappings for the desktop app. Keys must be a single uppercase Windows drive letter other than A, B, or C, and values are Files.com paths to mount there.
	MountMappings interface{} `url:"mount_mappings,omitempty" json:"mount_mappings,omitempty" path:"mount_mappings"`
	// Whether this profile applies to all users in the Workspace by default.
	UseForAllUsers *bool `url:"use_for_all_users,omitempty" json:"use_for_all_users,omitempty" path:"use_for_all_users"`
	// Whether the desktop app should hide drive mounting, prevent new drive mounts, and unmount active drive mounts for users with this profile.
	DisableDriveMounting *bool `url:"disable_drive_mounting,omitempty" json:"disable_drive_mounting,omitempty" path:"disable_drive_mounting"`
}

// DesktopConfigurationProfileDeleteParams contains the request parameters for DELETE /desktop_configuration_profiles/{id}.
type DesktopConfigurationProfileDeleteParams struct {
	// Desktop Configuration Profile ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (d *DesktopConfigurationProfile) UnmarshalJSON(data []byte) error {
	type desktopConfigurationProfile DesktopConfigurationProfile
	var v desktopConfigurationProfile
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*d = DesktopConfigurationProfile(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (d *DesktopConfigurationProfileCollection) UnmarshalJSON(data []byte) error {
	type desktopConfigurationProfiles DesktopConfigurationProfileCollection
	var v desktopConfigurationProfiles
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*d = DesktopConfigurationProfileCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (d *DesktopConfigurationProfileCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*d))
	for i, v := range *d {
		ret[i] = v
	}

	return &ret
}
