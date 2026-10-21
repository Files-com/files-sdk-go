package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// PartnerChannelTemplate is a Files.com API resource.
type PartnerChannelTemplate struct {
	// The unique ID of the Partner Channel Template.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// ID of the Workspace associated with this Partner Channel Template.
	WorkspaceId int64 `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	// Channel directions. `two_way` enables both directions, `to_partner` enables outgoing downloads, and `from_partner` enables incoming uploads.
	Direction string `json:"direction,omitempty" path:"direction,omitempty" url:"direction,omitempty"`
	// Use the Channel folder directly for a one-way exchange. Defaults to false. Cannot be changed after creation. Folder name overrides must be blank when enabled, and the Template must remain one-way.
	UseChannelRoot *bool `json:"use_channel_root,omitempty" path:"use_channel_root,omitempty" url:"use_channel_root,omitempty"`
	// The name of the Partner Channel Template.
	Name string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	// Channel path relative to the Partner root folder. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.
	Path string `json:"path,omitempty" path:"path,omitempty" url:"path,omitempty"`
	// Optional Channel-level to-Partner folder name override.
	ToPartnerFolderName string `json:"to_partner_folder_name,omitempty" path:"to_partner_folder_name,omitempty" url:"to_partner_folder_name,omitempty"`
	// Optional Channel-level from-Partner folder name override.
	FromPartnerFolderName string `json:"from_partner_folder_name,omitempty" path:"from_partner_folder_name,omitempty" url:"from_partner_folder_name,omitempty"`
	// Optional route path pattern for files uploaded by the Partner. Supports {{partner_name}}.
	FromPartnerRoutePathPattern string `json:"from_partner_route_path_pattern,omitempty" path:"from_partner_route_path_pattern,omitempty" url:"from_partner_route_path_pattern,omitempty"`
	// Optional route path pattern for files delivered to the Partner. Supports {{partner_name}}.
	ToPartnerRoutePathPattern string `json:"to_partner_route_path_pattern,omitempty" path:"to_partner_route_path_pattern,omitempty" url:"to_partner_route_path_pattern,omitempty"`
	// Managed folder paths inside the to-Partner folder.
	ToPartnerManagedFolderPaths []string `json:"to_partner_managed_folder_paths,omitempty" path:"to_partner_managed_folder_paths,omitempty" url:"to_partner_managed_folder_paths,omitempty"`
	// Managed folder paths inside the from-Partner folder.
	FromPartnerManagedFolderPaths []string `json:"from_partner_managed_folder_paths,omitempty" path:"from_partner_managed_folder_paths,omitempty" url:"from_partner_managed_folder_paths,omitempty"`
	// Resolved to-Partner subfolder name. Null when the direction is disabled or uses the Channel folder directly.
	EffectiveToPartnerFolderName string `json:"effective_to_partner_folder_name,omitempty" path:"effective_to_partner_folder_name,omitempty" url:"effective_to_partner_folder_name,omitempty"`
	// Resolved from-Partner subfolder name. Null when the direction is disabled or uses the Channel folder directly.
	EffectiveFromPartnerFolderName string `json:"effective_from_partner_folder_name,omitempty" path:"effective_from_partner_folder_name,omitempty" url:"effective_from_partner_folder_name,omitempty"`
}

// Identifier returns the resource ID.
func (p PartnerChannelTemplate) Identifier() interface{} {
	return p.Id
}

// PartnerChannelTemplateCollection is a list of PartnerChannelTemplate resources.
type PartnerChannelTemplateCollection []PartnerChannelTemplate

// PartnerChannelTemplateDirectionEnum is a string value for direction.
// Enum lists the values documented by the API.
type PartnerChannelTemplateDirectionEnum string

// String returns the API parameter value.
func (u PartnerChannelTemplateDirectionEnum) String() string {
	return string(u)
}

// Enum returns the documented values keyed by their API strings.
func (u PartnerChannelTemplateDirectionEnum) Enum() map[string]PartnerChannelTemplateDirectionEnum {
	return map[string]PartnerChannelTemplateDirectionEnum{
		"two_way":      PartnerChannelTemplateDirectionEnum("two_way"),
		"to_partner":   PartnerChannelTemplateDirectionEnum("to_partner"),
		"from_partner": PartnerChannelTemplateDirectionEnum("from_partner"),
	}
}

// PartnerChannelTemplateListParams contains the request parameters for GET /partner_channel_templates.
type PartnerChannelTemplateListParams struct {
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id` and `name`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `workspace_id`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// PartnerChannelTemplateFindParams contains the request parameters for GET /partner_channel_templates/{id}.
type PartnerChannelTemplateFindParams struct {
	// Partner Channel Template ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// PartnerChannelTemplateCreateParams contains the request parameters for POST /partner_channel_templates.
type PartnerChannelTemplateCreateParams struct {
	// Channel directions. `two_way` enables both directions, `to_partner` enables outgoing downloads, and `from_partner` enables incoming uploads. Values: two_way, to_partner, from_partner.
	Direction PartnerChannelTemplateDirectionEnum `url:"direction,omitempty" json:"direction,omitempty" path:"direction"`
	// Use the Channel folder directly for a one-way exchange. Defaults to false. Cannot be changed after creation. Folder name overrides must be blank when enabled, and the Template must remain one-way.
	UseChannelRoot *bool `url:"use_channel_root,omitempty" json:"use_channel_root,omitempty" path:"use_channel_root"`
	// Optional Channel-level from-Partner folder name override.
	FromPartnerFolderName string `url:"from_partner_folder_name,omitempty" json:"from_partner_folder_name,omitempty" path:"from_partner_folder_name"`
	// Managed folder paths inside the from-Partner folder.
	FromPartnerManagedFolderPaths []string `url:"from_partner_managed_folder_paths,omitempty" json:"from_partner_managed_folder_paths,omitempty" path:"from_partner_managed_folder_paths"`
	// Optional route path pattern for files uploaded by the Partner. Supports {{partner_name}}.
	FromPartnerRoutePathPattern string `url:"from_partner_route_path_pattern,omitempty" json:"from_partner_route_path_pattern,omitempty" path:"from_partner_route_path_pattern"`
	// Optional Channel-level to-Partner folder name override.
	ToPartnerFolderName string `url:"to_partner_folder_name,omitempty" json:"to_partner_folder_name,omitempty" path:"to_partner_folder_name"`
	// Managed folder paths inside the to-Partner folder.
	ToPartnerManagedFolderPaths []string `url:"to_partner_managed_folder_paths,omitempty" json:"to_partner_managed_folder_paths,omitempty" path:"to_partner_managed_folder_paths"`
	// Optional route path pattern for files delivered to the Partner. Supports {{partner_name}}.
	ToPartnerRoutePathPattern string `url:"to_partner_route_path_pattern,omitempty" json:"to_partner_route_path_pattern,omitempty" path:"to_partner_route_path_pattern"`
	// The name of the Partner Channel Template. Required.
	Name string `url:"name" json:"name" path:"name"`
	// Channel path relative to the Partner root folder. Required.
	Path string `url:"path" json:"path" path:"path"`
	// ID of the Workspace associated with this Partner Channel Template.
	WorkspaceId int64 `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
}

// PartnerChannelTemplateUpdateParams contains the request parameters for PATCH /partner_channel_templates/{id}.
type PartnerChannelTemplateUpdateParams struct {
	// Partner Channel Template ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// Channel directions. `two_way` enables both directions, `to_partner` enables outgoing downloads, and `from_partner` enables incoming uploads. Values: two_way, to_partner, from_partner.
	Direction PartnerChannelTemplateDirectionEnum `url:"direction,omitempty" json:"direction,omitempty" path:"direction"`
	// Use the Channel folder directly for a one-way exchange. Defaults to false. Cannot be changed after creation. Folder name overrides must be blank when enabled, and the Template must remain one-way.
	UseChannelRoot *bool `url:"use_channel_root,omitempty" json:"use_channel_root,omitempty" path:"use_channel_root"`
	// Optional Channel-level from-Partner folder name override.
	FromPartnerFolderName string `url:"from_partner_folder_name,omitempty" json:"from_partner_folder_name,omitempty" path:"from_partner_folder_name"`
	// Managed folder paths inside the from-Partner folder.
	FromPartnerManagedFolderPaths []string `url:"from_partner_managed_folder_paths,omitempty" json:"from_partner_managed_folder_paths,omitempty" path:"from_partner_managed_folder_paths"`
	// Optional route path pattern for files uploaded by the Partner. Supports {{partner_name}}.
	FromPartnerRoutePathPattern string `url:"from_partner_route_path_pattern,omitempty" json:"from_partner_route_path_pattern,omitempty" path:"from_partner_route_path_pattern"`
	// Optional Channel-level to-Partner folder name override.
	ToPartnerFolderName string `url:"to_partner_folder_name,omitempty" json:"to_partner_folder_name,omitempty" path:"to_partner_folder_name"`
	// Managed folder paths inside the to-Partner folder.
	ToPartnerManagedFolderPaths []string `url:"to_partner_managed_folder_paths,omitempty" json:"to_partner_managed_folder_paths,omitempty" path:"to_partner_managed_folder_paths"`
	// Optional route path pattern for files delivered to the Partner. Supports {{partner_name}}.
	ToPartnerRoutePathPattern string `url:"to_partner_route_path_pattern,omitempty" json:"to_partner_route_path_pattern,omitempty" path:"to_partner_route_path_pattern"`
	// The name of the Partner Channel Template.
	Name string `url:"name,omitempty" json:"name,omitempty" path:"name"`
	// Channel path relative to the Partner root folder.
	Path string `url:"path,omitempty" json:"path,omitempty" path:"path"`
}

// PartnerChannelTemplateDeleteParams contains the request parameters for DELETE /partner_channel_templates/{id}.
type PartnerChannelTemplateDeleteParams struct {
	// Partner Channel Template ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (p *PartnerChannelTemplate) UnmarshalJSON(data []byte) error {
	type partnerChannelTemplate PartnerChannelTemplate
	var v partnerChannelTemplate
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*p = PartnerChannelTemplate(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (p *PartnerChannelTemplateCollection) UnmarshalJSON(data []byte) error {
	type partnerChannelTemplates PartnerChannelTemplateCollection
	var v partnerChannelTemplates
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*p = PartnerChannelTemplateCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (p *PartnerChannelTemplateCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*p))
	for i, v := range *p {
		ret[i] = v
	}

	return &ret
}
