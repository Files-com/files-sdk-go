package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// SiteSubdomainRedirect is a Files.com API resource.
type SiteSubdomainRedirect struct {
	Id        int64      `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	Subdomain string     `json:"subdomain,omitempty" path:"subdomain,omitempty" url:"subdomain,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" path:"updated_at,omitempty" url:"updated_at,omitempty"`
}

// Identifier returns the resource ID.
func (s SiteSubdomainRedirect) Identifier() interface{} {
	return s.Id
}

// SiteSubdomainRedirectCollection is a list of SiteSubdomainRedirect resources.
type SiteSubdomainRedirectCollection []SiteSubdomainRedirect

// SiteSubdomainRedirectListParams contains the request parameters for GET /site_subdomain_redirects.
type SiteSubdomainRedirectListParams struct {
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	ListParams
}

// SiteSubdomainRedirectFindParams contains the request parameters for GET /site_subdomain_redirects/{id}.
type SiteSubdomainRedirectFindParams struct {
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// SiteSubdomainRedirectDeleteParams contains the request parameters for DELETE /site_subdomain_redirects/{id}.
type SiteSubdomainRedirectDeleteParams struct {
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (s *SiteSubdomainRedirect) UnmarshalJSON(data []byte) error {
	type siteSubdomainRedirect SiteSubdomainRedirect
	var v siteSubdomainRedirect
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*s = SiteSubdomainRedirect(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (s *SiteSubdomainRedirectCollection) UnmarshalJSON(data []byte) error {
	type siteSubdomainRedirects SiteSubdomainRedirectCollection
	var v siteSubdomainRedirects
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*s = SiteSubdomainRedirectCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (s *SiteSubdomainRedirectCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*s))
	for i, v := range *s {
		ret[i] = v
	}

	return &ret
}
