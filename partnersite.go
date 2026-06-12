package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

type PartnerSite struct {
	HostPartnerId    int64  `json:"host_partner_id,omitempty" path:"host_partner_id,omitempty" url:"host_partner_id,omitempty"`
	HostPartnerName  string `json:"host_partner_name,omitempty" path:"host_partner_name,omitempty" url:"host_partner_name,omitempty"`
	GuestPartnerId   int64  `json:"guest_partner_id,omitempty" path:"guest_partner_id,omitempty" url:"guest_partner_id,omitempty"`
	GuestPartnerName string `json:"guest_partner_name,omitempty" path:"guest_partner_name,omitempty" url:"guest_partner_name,omitempty"`
	HostSiteId       int64  `json:"host_site_id,omitempty" path:"host_site_id,omitempty" url:"host_site_id,omitempty"`
	HostSiteName     string `json:"host_site_name,omitempty" path:"host_site_name,omitempty" url:"host_site_name,omitempty"`
	GuestSiteId      int64  `json:"guest_site_id,omitempty" path:"guest_site_id,omitempty" url:"guest_site_id,omitempty"`
	GuestSiteName    string `json:"guest_site_name,omitempty" path:"guest_site_name,omitempty" url:"guest_site_name,omitempty"`
	WorkspaceId      int64  `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
}

// Identifier no path or id

type PartnerSiteCollection []PartnerSite

type PartnerSiteListParams struct {
	ListParams
}

func (p *PartnerSite) UnmarshalJSON(data []byte) error {
	type partnerSite PartnerSite
	var v partnerSite
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*p = PartnerSite(v)
	return nil
}

func (p *PartnerSiteCollection) UnmarshalJSON(data []byte) error {
	type partnerSites PartnerSiteCollection
	var v partnerSites
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*p = PartnerSiteCollection(v)
	return nil
}

func (p *PartnerSiteCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*p))
	for i, v := range *p {
		ret[i] = v
	}

	return &ret
}
