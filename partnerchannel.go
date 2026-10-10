package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// PartnerChannel is a Files.com API resource.
type PartnerChannel struct {
	Id                             int64    `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	WorkspaceId                    int64    `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	Direction                      string   `json:"direction,omitempty" path:"direction,omitempty" url:"direction,omitempty"`
	UseChannelRoot                 *bool    `json:"use_channel_root,omitempty" path:"use_channel_root,omitempty" url:"use_channel_root,omitempty"`
	PartnerId                      int64    `json:"partner_id,omitempty" path:"partner_id,omitempty" url:"partner_id,omitempty"`
	PartnerChannelTemplateId       int64    `json:"partner_channel_template_id,omitempty" path:"partner_channel_template_id,omitempty" url:"partner_channel_template_id,omitempty"`
	Path                           string   `json:"path,omitempty" path:"path,omitempty" url:"path,omitempty"`
	ToPartnerFolderName            string   `json:"to_partner_folder_name,omitempty" path:"to_partner_folder_name,omitempty" url:"to_partner_folder_name,omitempty"`
	FromPartnerFolderName          string   `json:"from_partner_folder_name,omitempty" path:"from_partner_folder_name,omitempty" url:"from_partner_folder_name,omitempty"`
	FromPartnerRoutePath           string   `json:"from_partner_route_path,omitempty" path:"from_partner_route_path,omitempty" url:"from_partner_route_path,omitempty"`
	ToPartnerRoutePath             string   `json:"to_partner_route_path,omitempty" path:"to_partner_route_path,omitempty" url:"to_partner_route_path,omitempty"`
	ToPartnerManagedFolderPaths    []string `json:"to_partner_managed_folder_paths,omitempty" path:"to_partner_managed_folder_paths,omitempty" url:"to_partner_managed_folder_paths,omitempty"`
	FromPartnerManagedFolderPaths  []string `json:"from_partner_managed_folder_paths,omitempty" path:"from_partner_managed_folder_paths,omitempty" url:"from_partner_managed_folder_paths,omitempty"`
	EffectiveToPartnerFolderName   string   `json:"effective_to_partner_folder_name,omitempty" path:"effective_to_partner_folder_name,omitempty" url:"effective_to_partner_folder_name,omitempty"`
	EffectiveFromPartnerFolderName string   `json:"effective_from_partner_folder_name,omitempty" path:"effective_from_partner_folder_name,omitempty" url:"effective_from_partner_folder_name,omitempty"`
	ChannelPath                    string   `json:"channel_path,omitempty" path:"channel_path,omitempty" url:"channel_path,omitempty"`
	ToPartnerFolderPath            string   `json:"to_partner_folder_path,omitempty" path:"to_partner_folder_path,omitempty" url:"to_partner_folder_path,omitempty"`
	FromPartnerFolderPath          string   `json:"from_partner_folder_path,omitempty" path:"from_partner_folder_path,omitempty" url:"from_partner_folder_path,omitempty"`
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
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	ListParams
}

// PartnerChannelFindParams contains the request parameters for GET /partner_channels/{id}.
type PartnerChannelFindParams struct {
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// PartnerChannelCreateParams contains the request parameters for POST /partner_channels.
type PartnerChannelCreateParams struct {
	Direction                     PartnerChannelDirectionEnum `url:"direction,omitempty" json:"direction,omitempty" path:"direction"`
	UseChannelRoot                *bool                       `url:"use_channel_root,omitempty" json:"use_channel_root,omitempty" path:"use_channel_root"`
	FromPartnerFolderName         string                      `url:"from_partner_folder_name,omitempty" json:"from_partner_folder_name,omitempty" path:"from_partner_folder_name"`
	FromPartnerManagedFolderPaths []string                    `url:"from_partner_managed_folder_paths,omitempty" json:"from_partner_managed_folder_paths,omitempty" path:"from_partner_managed_folder_paths"`
	FromPartnerRoutePath          string                      `url:"from_partner_route_path,omitempty" json:"from_partner_route_path,omitempty" path:"from_partner_route_path"`
	ToPartnerFolderName           string                      `url:"to_partner_folder_name,omitempty" json:"to_partner_folder_name,omitempty" path:"to_partner_folder_name"`
	ToPartnerManagedFolderPaths   []string                    `url:"to_partner_managed_folder_paths,omitempty" json:"to_partner_managed_folder_paths,omitempty" path:"to_partner_managed_folder_paths"`
	ToPartnerRoutePath            string                      `url:"to_partner_route_path,omitempty" json:"to_partner_route_path,omitempty" path:"to_partner_route_path"`
	PartnerId                     int64                       `url:"partner_id" json:"partner_id" path:"partner_id"`
	Path                          string                      `url:"path" json:"path" path:"path"`
	WorkspaceId                   int64                       `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
}

// PartnerChannelUpdateParams contains the request parameters for PATCH /partner_channels/{id}.
type PartnerChannelUpdateParams struct {
	Id                            int64                       `url:"-,omitempty" json:"-,omitempty" path:"id"`
	Direction                     PartnerChannelDirectionEnum `url:"direction,omitempty" json:"direction,omitempty" path:"direction"`
	UseChannelRoot                *bool                       `url:"use_channel_root,omitempty" json:"use_channel_root,omitempty" path:"use_channel_root"`
	FromPartnerFolderName         string                      `url:"from_partner_folder_name,omitempty" json:"from_partner_folder_name,omitempty" path:"from_partner_folder_name"`
	FromPartnerManagedFolderPaths []string                    `url:"from_partner_managed_folder_paths,omitempty" json:"from_partner_managed_folder_paths,omitempty" path:"from_partner_managed_folder_paths"`
	FromPartnerRoutePath          string                      `url:"from_partner_route_path,omitempty" json:"from_partner_route_path,omitempty" path:"from_partner_route_path"`
	ToPartnerFolderName           string                      `url:"to_partner_folder_name,omitempty" json:"to_partner_folder_name,omitempty" path:"to_partner_folder_name"`
	ToPartnerManagedFolderPaths   []string                    `url:"to_partner_managed_folder_paths,omitempty" json:"to_partner_managed_folder_paths,omitempty" path:"to_partner_managed_folder_paths"`
	ToPartnerRoutePath            string                      `url:"to_partner_route_path,omitempty" json:"to_partner_route_path,omitempty" path:"to_partner_route_path"`
	Path                          string                      `url:"path,omitempty" json:"path,omitempty" path:"path"`
}

// PartnerChannelDeleteParams contains the request parameters for DELETE /partner_channels/{id}.
type PartnerChannelDeleteParams struct {
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
