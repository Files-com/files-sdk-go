package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// Group is a Files.com API resource.
type Group struct {
	// Group ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Group name.
	Name string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	// A list of allowed IPs if applicable. Newline delimited.
	AllowedIps string `json:"allowed_ips,omitempty" path:"allowed_ips,omitempty" url:"allowed_ips,omitempty"`
	// Comma-delimited list of user IDs who are group administrators (separated by commas).
	AdminIds string `json:"admin_ids,omitempty" path:"admin_ids,omitempty" url:"admin_ids,omitempty"`
	// Notes about this group.
	Notes string `json:"notes,omitempty" path:"notes,omitempty" url:"notes,omitempty"`
	// Comma-delimited list of user IDs who belong to this group (separated by commas).
	UserIds string `json:"user_ids,omitempty" path:"user_ids,omitempty" url:"user_ids,omitempty"`
	// Comma-delimited list of usernames who belong to this group (separated by commas).
	Usernames string `json:"usernames,omitempty" path:"usernames,omitempty" url:"usernames,omitempty"`
	// AI Assistant Personality ID assigned to this Group, if any. Users in the Group inherit it unless a direct per-user or Partner assignment overrides it.
	AiAssistantPersonalityId int64 `json:"ai_assistant_personality_id,omitempty" path:"ai_assistant_personality_id,omitempty" url:"ai_assistant_personality_id,omitempty"`
	// If true, users in this group can use FTP to login. This will override a false value of `ftp_permission` on the user level.
	FtpPermission *bool `json:"ftp_permission,omitempty" path:"ftp_permission,omitempty" url:"ftp_permission,omitempty"`
	// If true, users in this group can use SFTP to login. This will override a false value of `sftp_permission` on the user level.
	SftpPermission *bool `json:"sftp_permission,omitempty" path:"sftp_permission,omitempty" url:"sftp_permission,omitempty"`
	// If true, users in this group can use WebDAV to login. This will override a false value of `dav_permission` on the user level.
	DavPermission *bool `json:"dav_permission,omitempty" path:"dav_permission,omitempty" url:"dav_permission,omitempty"`
	// If true, users in this group can use the REST API to login. This will override a false value of `restapi_permission` on the user level.
	RestapiPermission *bool `json:"restapi_permission,omitempty" path:"restapi_permission,omitempty" url:"restapi_permission,omitempty"`
	// If true, users in this group can access the S3-compatible endpoint. This will override a false value of `s3_compatible_endpoint_permission` on the user level. Defaults to false.
	S3CompatibleEndpointPermission *bool `json:"s3_compatible_endpoint_permission,omitempty" path:"s3_compatible_endpoint_permission,omitempty" url:"s3_compatible_endpoint_permission,omitempty"`
	// Desktop Configuration Profile ID assigned to this Group, if any. Users in the Group inherit it unless a direct per-user assignment overrides it.
	DesktopConfigurationProfileId int64 `json:"desktop_configuration_profile_id,omitempty" path:"desktop_configuration_profile_id,omitempty" url:"desktop_configuration_profile_id,omitempty"`
	// Integration Centric Profile ID assigned to this Group, if any. Users in the Group inherit it unless a direct per-user assignment overrides it.
	IntegrationCentricProfileId int64 `json:"integration_centric_profile_id,omitempty" path:"integration_centric_profile_id,omitempty" url:"integration_centric_profile_id,omitempty"`
	// Site ID.
	SiteId int64 `json:"site_id,omitempty" path:"site_id,omitempty" url:"site_id,omitempty"`
	// Workspace ID.
	WorkspaceId int64 `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
}

// Identifier returns the resource ID.
func (g Group) Identifier() interface{} {
	return g.Id
}

// GroupCollection is a list of Group resources.
type GroupCollection []Group

// GroupListParams contains the request parameters for GET /groups.
type GroupListParams struct {
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `site_id`, `workspace_id` or `name`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `name` and `workspace_id`. Valid field combinations are `[ workspace_id, name ]`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// If set, return records where the specified field is prefixed by the supplied value. Valid fields are `name`.
	FilterPrefix interface{} `url:"filter_prefix,omitempty" json:"filter_prefix,omitempty" path:"filter_prefix"`
	// Comma-separated list of group ids to include in results.
	Ids string `url:"ids,omitempty" json:"ids,omitempty" path:"ids"`
	// Include groups from the parent site.
	IncludeParentSiteGroups *bool `url:"include_parent_site_groups,omitempty" json:"include_parent_site_groups,omitempty" path:"include_parent_site_groups"`
	// When scoped to a Custom Workspace, include Default Workspace groups that have permissions in that Workspace. Available to Site Administrators, and to Workspace Administrators and Folder Admins in that Workspace when the site shares its Default Workspace directory with Workspaces.
	IncludeDefaultWorkspaceGroups *bool `url:"include_default_workspace_groups,omitempty" json:"include_default_workspace_groups,omitempty" path:"include_default_workspace_groups"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// GroupFindParams contains the request parameters for GET /groups/{id}.
type GroupFindParams struct {
	// Group ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// GroupCreateParams contains the request parameters for POST /groups.
type GroupCreateParams struct {
	// Group notes.
	Notes string `url:"notes,omitempty" json:"notes,omitempty" path:"notes"`
	// A list of user ids. If sent as a string, should be comma-delimited.
	UserIds string `url:"user_ids,omitempty" json:"user_ids,omitempty" path:"user_ids"`
	// A list of group admin user ids. If sent as a string, should be comma-delimited.
	AdminIds string `url:"admin_ids,omitempty" json:"admin_ids,omitempty" path:"admin_ids"`
	// AI Assistant Personality ID assigned to this Group, if any. Users in the Group inherit it unless a direct per-user or Partner assignment overrides it.
	AiAssistantPersonalityId int64 `url:"ai_assistant_personality_id,omitempty" json:"ai_assistant_personality_id,omitempty" path:"ai_assistant_personality_id"`
	// If true, users in this group can use FTP to login. This will override a false value of `ftp_permission` on the user level.
	FtpPermission *bool `url:"ftp_permission,omitempty" json:"ftp_permission,omitempty" path:"ftp_permission"`
	// If true, users in this group can use SFTP to login. This will override a false value of `sftp_permission` on the user level.
	SftpPermission *bool `url:"sftp_permission,omitempty" json:"sftp_permission,omitempty" path:"sftp_permission"`
	// If true, users in this group can use WebDAV to login. This will override a false value of `dav_permission` on the user level.
	DavPermission *bool `url:"dav_permission,omitempty" json:"dav_permission,omitempty" path:"dav_permission"`
	// If true, users in this group can use the REST API to login. This will override a false value of `restapi_permission` on the user level.
	RestapiPermission *bool `url:"restapi_permission,omitempty" json:"restapi_permission,omitempty" path:"restapi_permission"`
	// If true, users in this group can access the S3-compatible endpoint. This will override a false value of `s3_compatible_endpoint_permission` on the user level. Defaults to false.
	S3CompatibleEndpointPermission *bool `url:"s3_compatible_endpoint_permission,omitempty" json:"s3_compatible_endpoint_permission,omitempty" path:"s3_compatible_endpoint_permission"`
	// Desktop Configuration Profile ID assigned to this Group, if any. Users in the Group inherit it unless a direct per-user assignment overrides it.
	DesktopConfigurationProfileId int64 `url:"desktop_configuration_profile_id,omitempty" json:"desktop_configuration_profile_id,omitempty" path:"desktop_configuration_profile_id"`
	// Integration Centric Profile ID assigned to this Group, if any. Users in the Group inherit it unless a direct per-user assignment overrides it.
	IntegrationCentricProfileId int64 `url:"integration_centric_profile_id,omitempty" json:"integration_centric_profile_id,omitempty" path:"integration_centric_profile_id"`
	// A list of allowed IPs if applicable. Newline delimited.
	AllowedIps string `url:"allowed_ips,omitempty" json:"allowed_ips,omitempty" path:"allowed_ips"`
	// Group name. Required.
	Name string `url:"name" json:"name" path:"name"`
	// Workspace ID.
	WorkspaceId int64 `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
}

// GroupUpdateParams contains the request parameters for PATCH /groups/{id}.
type GroupUpdateParams struct {
	// Group ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// Group notes.
	Notes string `url:"notes,omitempty" json:"notes,omitempty" path:"notes"`
	// A list of user ids. If sent as a string, should be comma-delimited.
	UserIds string `url:"user_ids,omitempty" json:"user_ids,omitempty" path:"user_ids"`
	// A list of group admin user ids. If sent as a string, should be comma-delimited.
	AdminIds string `url:"admin_ids,omitempty" json:"admin_ids,omitempty" path:"admin_ids"`
	// AI Assistant Personality ID assigned to this Group, if any. Users in the Group inherit it unless a direct per-user or Partner assignment overrides it.
	AiAssistantPersonalityId int64 `url:"ai_assistant_personality_id,omitempty" json:"ai_assistant_personality_id,omitempty" path:"ai_assistant_personality_id"`
	// If true, users in this group can use FTP to login. This will override a false value of `ftp_permission` on the user level.
	FtpPermission *bool `url:"ftp_permission,omitempty" json:"ftp_permission,omitempty" path:"ftp_permission"`
	// If true, users in this group can use SFTP to login. This will override a false value of `sftp_permission` on the user level.
	SftpPermission *bool `url:"sftp_permission,omitempty" json:"sftp_permission,omitempty" path:"sftp_permission"`
	// If true, users in this group can use WebDAV to login. This will override a false value of `dav_permission` on the user level.
	DavPermission *bool `url:"dav_permission,omitempty" json:"dav_permission,omitempty" path:"dav_permission"`
	// If true, users in this group can use the REST API to login. This will override a false value of `restapi_permission` on the user level.
	RestapiPermission *bool `url:"restapi_permission,omitempty" json:"restapi_permission,omitempty" path:"restapi_permission"`
	// If true, users in this group can access the S3-compatible endpoint. This will override a false value of `s3_compatible_endpoint_permission` on the user level. Defaults to false.
	S3CompatibleEndpointPermission *bool `url:"s3_compatible_endpoint_permission,omitempty" json:"s3_compatible_endpoint_permission,omitempty" path:"s3_compatible_endpoint_permission"`
	// Desktop Configuration Profile ID assigned to this Group, if any. Users in the Group inherit it unless a direct per-user assignment overrides it.
	DesktopConfigurationProfileId int64 `url:"desktop_configuration_profile_id,omitempty" json:"desktop_configuration_profile_id,omitempty" path:"desktop_configuration_profile_id"`
	// Integration Centric Profile ID assigned to this Group, if any. Users in the Group inherit it unless a direct per-user assignment overrides it.
	IntegrationCentricProfileId int64 `url:"integration_centric_profile_id,omitempty" json:"integration_centric_profile_id,omitempty" path:"integration_centric_profile_id"`
	// A list of allowed IPs if applicable. Newline delimited.
	AllowedIps string `url:"allowed_ips,omitempty" json:"allowed_ips,omitempty" path:"allowed_ips"`
	// Group name.
	Name string `url:"name,omitempty" json:"name,omitempty" path:"name"`
}

// GroupDeleteParams contains the request parameters for DELETE /groups/{id}.
type GroupDeleteParams struct {
	// Group ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (g *Group) UnmarshalJSON(data []byte) error {
	type group Group
	var v group
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*g = Group(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (g *GroupCollection) UnmarshalJSON(data []byte) error {
	type groups GroupCollection
	var v groups
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*g = GroupCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (g *GroupCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*g))
	for i, v := range *g {
		ret[i] = v
	}

	return &ret
}
