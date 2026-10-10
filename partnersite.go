package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// PartnerSite is a Files.com API resource.
type PartnerSite struct {
}

// Identifier no path or id

// PartnerSiteCollection is a list of PartnerSite resources.
type PartnerSiteCollection []PartnerSite

// PartnerSiteDeleteParams contains the request parameters for DELETE /partner_sites/{id}.
type PartnerSiteDeleteParams struct {
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (p *PartnerSite) UnmarshalJSON(data []byte) error {
	type partnerSite PartnerSite
	var v partnerSite
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*p = PartnerSite(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (p *PartnerSiteCollection) UnmarshalJSON(data []byte) error {
	type partnerSites PartnerSiteCollection
	var v partnerSites
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*p = PartnerSiteCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (p *PartnerSiteCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*p))
	for i, v := range *p {
		ret[i] = v
	}

	return &ret
}
