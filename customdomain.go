package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// CustomDomain is a Files.com API resource.
type CustomDomain struct {
	// Custom Domain ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Workspace ID (0 for the default workspace).
	WorkspaceId int64 `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	// Allow all workspaces to use this default-workspace Custom Domain.
	AvailableToAllWorkspaces *bool `json:"available_to_all_workspaces,omitempty" path:"available_to_all_workspaces,omitempty" url:"available_to_all_workspaces,omitempty"`
	// Eligible public IP addresses for Remote Server outbound connections. Empty when this Custom Domain is not eligible for outbound selection.
	OutboundIpAddresses []string `json:"outbound_ip_addresses,omitempty" path:"outbound_ip_addresses,omitempty" url:"outbound_ip_addresses,omitempty"`
	// Customer-owned domain name.
	Domain string `json:"domain,omitempty" path:"domain,omitempty" url:"domain,omitempty"`
	// Where this custom domain routes. Can be `site_alias`, `public_hosting`, `s3_endpoint`, or `unassigned` (not routing traffic). Set to `unassigned` automatically when a bound `public_hosting` folder behavior is deleted, and can be set manually via the API for any reason.
	Destination string `json:"destination,omitempty" path:"destination,omitempty" url:"destination,omitempty"`
	// Current DNS verification status.
	DnsStatus string `json:"dns_status,omitempty" path:"dns_status,omitempty" url:"dns_status,omitempty"`
	// Current SSL certificate ID.
	SslCertificateId int64 `json:"ssl_certificate_id,omitempty" path:"ssl_certificate_id,omitempty" url:"ssl_certificate_id,omitempty"`
	// Is this domain's SSL certificate automatically managed and renewed by Files.com?
	BrickManaged *bool `json:"brick_managed,omitempty" path:"brick_managed,omitempty" url:"brick_managed,omitempty"`
	// Public Hosting behavior ID when this domain routes to a specific Public Hosting behavior. Preserved as historical context when `destination` becomes `unassigned`.
	FolderBehaviorId int64 `json:"folder_behavior_id,omitempty" path:"folder_behavior_id,omitempty" url:"folder_behavior_id,omitempty"`
	// Dedicated public IP addresses allocated to this Custom Domain.
	IpAddresses []string `json:"ip_addresses,omitempty" path:"ip_addresses,omitempty" url:"ip_addresses,omitempty"`
	// When this Custom Domain was created.
	CreatedAt *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	// When this Custom Domain was last updated.
	UpdatedAt *time.Time `json:"updated_at,omitempty" path:"updated_at,omitempty" url:"updated_at,omitempty"`
}

// Identifier returns the resource ID.
func (c CustomDomain) Identifier() interface{} {
	return c.Id
}

// CustomDomainCollection is a list of CustomDomain resources.
type CustomDomainCollection []CustomDomain

// CustomDomainDestinationEnum is a string value for destination.
// Enum lists the values documented by the API.
type CustomDomainDestinationEnum string

// String returns the API parameter value.
func (u CustomDomainDestinationEnum) String() string {
	return string(u)
}

// Enum returns the documented values keyed by their API strings.
func (u CustomDomainDestinationEnum) Enum() map[string]CustomDomainDestinationEnum {
	return map[string]CustomDomainDestinationEnum{
		"site_alias":     CustomDomainDestinationEnum("site_alias"),
		"public_hosting": CustomDomainDestinationEnum("public_hosting"),
		"s3_endpoint":    CustomDomainDestinationEnum("s3_endpoint"),
		"unassigned":     CustomDomainDestinationEnum("unassigned"),
	}
}

// CustomDomainListParams contains the request parameters for GET /custom_domains.
type CustomDomainListParams struct {
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id`, `id` or `available_to_all_workspaces`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `workspace_id`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// CustomDomainFindParams contains the request parameters for GET /custom_domains/{id}.
type CustomDomainFindParams struct {
	// Custom Domain ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// CustomDomainCreateAllocateIpParams contains the request parameters for POST /custom_domains/{id}/allocate_ips.
type CustomDomainCreateAllocateIpParams struct {
	// Custom Domain ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// Number of dedicated IP addresses to allocate. Values: 1, 2. Required.
	Count int64 `url:"count" json:"count" path:"count"`
}

// CustomDomainCreateParams contains the request parameters for POST /custom_domains.
type CustomDomainCreateParams struct {
	// Allow all workspaces to use this default-workspace Custom Domain.
	AvailableToAllWorkspaces *bool `url:"available_to_all_workspaces,omitempty" json:"available_to_all_workspaces,omitempty" path:"available_to_all_workspaces"`
	// Workspace ID (0 for the default workspace).
	WorkspaceId int64 `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
	// Where this custom domain routes. Can be `site_alias`, `public_hosting`, `s3_endpoint`, or `unassigned` (not routing traffic). Set to `unassigned` automatically when a bound `public_hosting` folder behavior is deleted, and can be set manually via the API for any reason. Values: site_alias, public_hosting, s3_endpoint, unassigned.
	Destination CustomDomainDestinationEnum `url:"destination,omitempty" json:"destination,omitempty" path:"destination"`
	// Public Hosting behavior ID when this domain routes to a specific Public Hosting behavior. Preserved as historical context when `destination` becomes `unassigned`.
	FolderBehaviorId int64 `url:"folder_behavior_id,omitempty" json:"folder_behavior_id,omitempty" path:"folder_behavior_id"`
	// Current SSL certificate ID.
	SslCertificateId int64 `url:"ssl_certificate_id,omitempty" json:"ssl_certificate_id,omitempty" path:"ssl_certificate_id"`
	// Customer-owned domain name. Required.
	Domain string `url:"domain" json:"domain" path:"domain"`
}

// CustomDomainUpdateParams contains the request parameters for PATCH /custom_domains/{id}.
type CustomDomainUpdateParams struct {
	// Custom Domain ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// Allow all workspaces to use this default-workspace Custom Domain.
	AvailableToAllWorkspaces *bool `url:"available_to_all_workspaces,omitempty" json:"available_to_all_workspaces,omitempty" path:"available_to_all_workspaces"`
	// Workspace ID (0 for the default workspace).
	WorkspaceId int64 `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
	// Where this custom domain routes. Can be `site_alias`, `public_hosting`, `s3_endpoint`, or `unassigned` (not routing traffic). Set to `unassigned` automatically when a bound `public_hosting` folder behavior is deleted, and can be set manually via the API for any reason. Values: site_alias, public_hosting, s3_endpoint, unassigned.
	Destination CustomDomainDestinationEnum `url:"destination,omitempty" json:"destination,omitempty" path:"destination"`
	// Public Hosting behavior ID when this domain routes to a specific Public Hosting behavior. Preserved as historical context when `destination` becomes `unassigned`.
	FolderBehaviorId int64 `url:"folder_behavior_id,omitempty" json:"folder_behavior_id,omitempty" path:"folder_behavior_id"`
	// Current SSL certificate ID.
	SslCertificateId int64 `url:"ssl_certificate_id,omitempty" json:"ssl_certificate_id,omitempty" path:"ssl_certificate_id"`
	// Customer-owned domain name.
	Domain string `url:"domain,omitempty" json:"domain,omitempty" path:"domain"`
}

// CustomDomainDeleteParams contains the request parameters for DELETE /custom_domains/{id}.
type CustomDomainDeleteParams struct {
	// Custom Domain ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (c *CustomDomain) UnmarshalJSON(data []byte) error {
	type customDomain CustomDomain
	var v customDomain
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*c = CustomDomain(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (c *CustomDomainCollection) UnmarshalJSON(data []byte) error {
	type customDomains CustomDomainCollection
	var v customDomains
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*c = CustomDomainCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (c *CustomDomainCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*c))
	for i, v := range *c {
		ret[i] = v
	}

	return &ret
}
