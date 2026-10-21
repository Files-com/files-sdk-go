package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// PartnerChannel is a Files.com API resource.
type PartnerChannel struct {
	// The unique ID of the Partner Channel.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// ID of the Workspace associated with this Partner Channel.
	WorkspaceId int64 `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	// Channel directions. `two_way` enables both directions, `to_partner` enables outgoing downloads, and `from_partner` enables incoming uploads.
	Direction string `json:"direction,omitempty" path:"direction,omitempty" url:"direction,omitempty"`
	// Use the Channel folder directly for a one-way exchange. Defaults to false. Cannot be changed after creation. Folder name overrides must be blank when enabled, and the Channel must remain one-way.
	UseChannelRoot *bool `json:"use_channel_root,omitempty" path:"use_channel_root,omitempty" url:"use_channel_root,omitempty"`
	// ID of the Partner this Channel belongs to.
	PartnerId int64 `json:"partner_id,omitempty" path:"partner_id,omitempty" url:"partner_id,omitempty"`
	// ID of the Partner Channel Template that manages this Channel, if any.
	PartnerChannelTemplateId int64 `json:"partner_channel_template_id,omitempty" path:"partner_channel_template_id,omitempty" url:"partner_channel_template_id,omitempty"`
	// Channel path relative to the Partner root folder. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.
	Path string `json:"path,omitempty" path:"path,omitempty" url:"path,omitempty"`
	// Optional Channel-level to-Partner folder name override.
	ToPartnerFolderName string `json:"to_partner_folder_name,omitempty" path:"to_partner_folder_name,omitempty" url:"to_partner_folder_name,omitempty"`
	// Optional Channel-level from-Partner folder name override.
	FromPartnerFolderName string `json:"from_partner_folder_name,omitempty" path:"from_partner_folder_name,omitempty" url:"from_partner_folder_name,omitempty"`
	// Optional route path for files uploaded by the Partner.
	FromPartnerRoutePath string `json:"from_partner_route_path,omitempty" path:"from_partner_route_path,omitempty" url:"from_partner_route_path,omitempty"`
	// Optional route path for files delivered to the Partner.
	ToPartnerRoutePath string `json:"to_partner_route_path,omitempty" path:"to_partner_route_path,omitempty" url:"to_partner_route_path,omitempty"`
	// Managed folder paths inside the to-Partner folder.
	ToPartnerManagedFolderPaths []string `json:"to_partner_managed_folder_paths,omitempty" path:"to_partner_managed_folder_paths,omitempty" url:"to_partner_managed_folder_paths,omitempty"`
	// Managed folder paths inside the from-Partner folder.
	FromPartnerManagedFolderPaths []string `json:"from_partner_managed_folder_paths,omitempty" path:"from_partner_managed_folder_paths,omitempty" url:"from_partner_managed_folder_paths,omitempty"`
	// Resolved to-Partner subfolder name. Null when the direction is disabled or uses the Channel folder directly.
	EffectiveToPartnerFolderName string `json:"effective_to_partner_folder_name,omitempty" path:"effective_to_partner_folder_name,omitempty" url:"effective_to_partner_folder_name,omitempty"`
	// Resolved from-Partner subfolder name. Null when the direction is disabled or uses the Channel folder directly.
	EffectiveFromPartnerFolderName string `json:"effective_from_partner_folder_name,omitempty" path:"effective_from_partner_folder_name,omitempty" url:"effective_from_partner_folder_name,omitempty"`
	// Resolved Channel folder path.
	ChannelPath string `json:"channel_path,omitempty" path:"channel_path,omitempty" url:"channel_path,omitempty"`
	// Resolved to-Partner folder path.
	ToPartnerFolderPath string `json:"to_partner_folder_path,omitempty" path:"to_partner_folder_path,omitempty" url:"to_partner_folder_path,omitempty"`
	// Resolved from-Partner folder path.
	FromPartnerFolderPath string `json:"from_partner_folder_path,omitempty" path:"from_partner_folder_path,omitempty" url:"from_partner_folder_path,omitempty"`
}

// Identifier returns the resource ID.
func (p PartnerChannel) Identifier() interface{} {
	return p.Id
}

// PartnerChannelCollection is a list of PartnerChannel resources.
type PartnerChannelCollection []PartnerChannel

// PartnerChannelDirectionEnum is a string value for direction.
// Enum lists the values documented by the API.
type PartnerChannelDirectionEnum string

// String returns the API parameter value.
func (u PartnerChannelDirectionEnum) String() string {
	return string(u)
}

// Enum returns the documented values keyed by their API strings.
func (u PartnerChannelDirectionEnum) Enum() map[string]PartnerChannelDirectionEnum {
	return map[string]PartnerChannelDirectionEnum{
		"two_way":      PartnerChannelDirectionEnum("two_way"),
		"to_partner":   PartnerChannelDirectionEnum("to_partner"),
		"from_partner": PartnerChannelDirectionEnum("from_partner"),
	}
}

// PartnerChannelListParams contains the request parameters for GET /partner_channels.
type PartnerChannelListParams struct {
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id`, `path` or `partner_id`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `partner_id` and `workspace_id`. Valid field combinations are `[ workspace_id, partner_id ]`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// PartnerChannelFindParams contains the request parameters for GET /partner_channels/{id}.
type PartnerChannelFindParams struct {
	// Partner Channel ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// PartnerChannelCreateParams contains the request parameters for POST /partner_channels.
type PartnerChannelCreateParams struct {
	// Channel directions. `two_way` enables both directions, `to_partner` enables outgoing downloads, and `from_partner` enables incoming uploads. Values: two_way, to_partner, from_partner.
	Direction PartnerChannelDirectionEnum `url:"direction,omitempty" json:"direction,omitempty" path:"direction"`
	// Use the Channel folder directly for a one-way exchange. Defaults to false. Cannot be changed after creation. Folder name overrides must be blank when enabled, and the Channel must remain one-way.
	UseChannelRoot *bool `url:"use_channel_root,omitempty" json:"use_channel_root,omitempty" path:"use_channel_root"`
	// Optional Channel-level from-Partner folder name override.
	FromPartnerFolderName string `url:"from_partner_folder_name,omitempty" json:"from_partner_folder_name,omitempty" path:"from_partner_folder_name"`
	// Managed folder paths inside the from-Partner folder.
	FromPartnerManagedFolderPaths []string `url:"from_partner_managed_folder_paths,omitempty" json:"from_partner_managed_folder_paths,omitempty" path:"from_partner_managed_folder_paths"`
	// Optional route path for files uploaded by the Partner.
	FromPartnerRoutePath string `url:"from_partner_route_path,omitempty" json:"from_partner_route_path,omitempty" path:"from_partner_route_path"`
	// Optional Channel-level to-Partner folder name override.
	ToPartnerFolderName string `url:"to_partner_folder_name,omitempty" json:"to_partner_folder_name,omitempty" path:"to_partner_folder_name"`
	// Managed folder paths inside the to-Partner folder.
	ToPartnerManagedFolderPaths []string `url:"to_partner_managed_folder_paths,omitempty" json:"to_partner_managed_folder_paths,omitempty" path:"to_partner_managed_folder_paths"`
	// Optional route path for files delivered to the Partner.
	ToPartnerRoutePath string `url:"to_partner_route_path,omitempty" json:"to_partner_route_path,omitempty" path:"to_partner_route_path"`
	// ID of the Partner this Channel belongs to. Required.
	PartnerId int64 `url:"partner_id" json:"partner_id" path:"partner_id"`
	// Channel path relative to the Partner root folder. Required.
	Path string `url:"path" json:"path" path:"path"`
	// ID of the Workspace associated with this Partner Channel.
	WorkspaceId int64 `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
}

// PartnerChannelUpdateParams contains the request parameters for PATCH /partner_channels/{id}.
type PartnerChannelUpdateParams struct {
	// Partner Channel ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// Channel directions. `two_way` enables both directions, `to_partner` enables outgoing downloads, and `from_partner` enables incoming uploads. Values: two_way, to_partner, from_partner.
	Direction PartnerChannelDirectionEnum `url:"direction,omitempty" json:"direction,omitempty" path:"direction"`
	// Use the Channel folder directly for a one-way exchange. Defaults to false. Cannot be changed after creation. Folder name overrides must be blank when enabled, and the Channel must remain one-way.
	UseChannelRoot *bool `url:"use_channel_root,omitempty" json:"use_channel_root,omitempty" path:"use_channel_root"`
	// Optional Channel-level from-Partner folder name override.
	FromPartnerFolderName string `url:"from_partner_folder_name,omitempty" json:"from_partner_folder_name,omitempty" path:"from_partner_folder_name"`
	// Managed folder paths inside the from-Partner folder.
	FromPartnerManagedFolderPaths []string `url:"from_partner_managed_folder_paths,omitempty" json:"from_partner_managed_folder_paths,omitempty" path:"from_partner_managed_folder_paths"`
	// Optional route path for files uploaded by the Partner.
	FromPartnerRoutePath string `url:"from_partner_route_path,omitempty" json:"from_partner_route_path,omitempty" path:"from_partner_route_path"`
	// Optional Channel-level to-Partner folder name override.
	ToPartnerFolderName string `url:"to_partner_folder_name,omitempty" json:"to_partner_folder_name,omitempty" path:"to_partner_folder_name"`
	// Managed folder paths inside the to-Partner folder.
	ToPartnerManagedFolderPaths []string `url:"to_partner_managed_folder_paths,omitempty" json:"to_partner_managed_folder_paths,omitempty" path:"to_partner_managed_folder_paths"`
	// Optional route path for files delivered to the Partner.
	ToPartnerRoutePath string `url:"to_partner_route_path,omitempty" json:"to_partner_route_path,omitempty" path:"to_partner_route_path"`
	// Channel path relative to the Partner root folder.
	Path string `url:"path,omitempty" json:"path,omitempty" path:"path"`
}

// PartnerChannelDeleteParams contains the request parameters for DELETE /partner_channels/{id}.
type PartnerChannelDeleteParams struct {
	// Partner Channel ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (p *PartnerChannel) UnmarshalJSON(data []byte) error {
	type partnerChannel PartnerChannel
	var v partnerChannel
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*p = PartnerChannel(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (p *PartnerChannelCollection) UnmarshalJSON(data []byte) error {
	type partnerChannels PartnerChannelCollection
	var v partnerChannels
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*p = PartnerChannelCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (p *PartnerChannelCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*p))
	for i, v := range *p {
		ret[i] = v
	}

	return &ret
}
