package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// PartnerSiteRequest is a Files.com API resource.
type PartnerSiteRequest struct {
	// Partner Site Request ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Host Partner ID.
	HostPartnerId int64 `json:"host_partner_id,omitempty" path:"host_partner_id,omitempty" url:"host_partner_id,omitempty"`
	// Guest Site URL.
	GuestSiteUrl string `json:"guest_site_url,omitempty" path:"guest_site_url,omitempty" url:"guest_site_url,omitempty"`
	// Request status (pending, approved, rejected).
	Status string `json:"status,omitempty" path:"status,omitempty" url:"status,omitempty"`
	// Host Site Name.
	HostSiteName string `json:"host_site_name,omitempty" path:"host_site_name,omitempty" url:"host_site_name,omitempty"`
	// Pairing key used to approve this request on the Guest Site.
	PairingKey string `json:"pairing_key,omitempty" path:"pairing_key,omitempty" url:"pairing_key,omitempty"`
	// Request creation date/time.
	CreatedAt *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	// Request last updated date/time.
	UpdatedAt *time.Time `json:"updated_at,omitempty" path:"updated_at,omitempty" url:"updated_at,omitempty"`
}

// Identifier returns the resource ID.
func (p PartnerSiteRequest) Identifier() interface{} {
	return p.Id
}

// PartnerSiteRequestCollection is a list of PartnerSiteRequest resources.
type PartnerSiteRequestCollection []PartnerSiteRequest

// PartnerSiteRequestListParams contains the request parameters for GET /partner_site_requests.
type PartnerSiteRequestListParams struct {
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `host_partner_id`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `host_partner_id`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// PartnerSiteRequestFindByPairingKeyParams contains the request parameters for GET /partner_site_requests/find_by_pairing_key.
type PartnerSiteRequestFindByPairingKeyParams struct {
	// Pairing key for the partner site request. Required.
	PairingKey string `url:"pairing_key" json:"pairing_key" path:"pairing_key"`
}

// PartnerSiteRequestCreateParams contains the request parameters for POST /partner_site_requests.
type PartnerSiteRequestCreateParams struct {
	// Host Partner ID to link with. Required.
	HostPartnerId int64 `url:"host_partner_id" json:"host_partner_id" path:"host_partner_id"`
	// Guest Site URL to link to. Required.
	GuestSiteUrl string `url:"guest_site_url" json:"guest_site_url" path:"guest_site_url"`
}

// PartnerSiteRequestRejectParams contains the request parameters for POST /partner_site_requests/reject.
type PartnerSiteRequestRejectParams struct {
	// Pairing key for the partner site request. Required.
	PairingKey string `url:"pairing_key" json:"pairing_key" path:"pairing_key"`
}

// PartnerSiteRequestApproveParams contains the request parameters for POST /partner_site_requests/approve.
type PartnerSiteRequestApproveParams struct {
	// Pairing key for the partner site request. Required.
	PairingKey string `url:"pairing_key" json:"pairing_key" path:"pairing_key"`
	// ID of an existing Partner on this site, with the host role, that represents the requesting organization. The connection binds to that Partner and makes it host_and_guest. When omitted, a guest Partner named after the host site is created.
	PartnerId int64 `url:"partner_id,omitempty" json:"partner_id,omitempty" path:"partner_id"`
}

// PartnerSiteRequestDeleteParams contains the request parameters for DELETE /partner_site_requests/{id}.
type PartnerSiteRequestDeleteParams struct {
	// Partner Site Request ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (p *PartnerSiteRequest) UnmarshalJSON(data []byte) error {
	type partnerSiteRequest PartnerSiteRequest
	var v partnerSiteRequest
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*p = PartnerSiteRequest(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (p *PartnerSiteRequestCollection) UnmarshalJSON(data []byte) error {
	type partnerSiteRequests PartnerSiteRequestCollection
	var v partnerSiteRequests
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*p = PartnerSiteRequestCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (p *PartnerSiteRequestCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*p))
	for i, v := range *p {
		ret[i] = v
	}

	return &ret
}
