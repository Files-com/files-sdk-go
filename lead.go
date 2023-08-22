package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Lead struct {
	Id                 int64  `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	Code               string `json:"code,omitempty" path:"code,omitempty" url:"code,omitempty"`
	Name               string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	Address            string `json:"address,omitempty" path:"address,omitempty" url:"address,omitempty"`
	Address2           string `json:"address_2,omitempty" path:"address_2,omitempty" url:"address_2,omitempty"`
	City               string `json:"city,omitempty" path:"city,omitempty" url:"city,omitempty"`
	CompanyName        string `json:"company_name,omitempty" path:"company_name,omitempty" url:"company_name,omitempty"`
	ContactName        string `json:"contact_name,omitempty" path:"contact_name,omitempty" url:"contact_name,omitempty"`
	Country            string `json:"country,omitempty" path:"country,omitempty" url:"country,omitempty"`
	Currency           string `json:"currency,omitempty" path:"currency,omitempty" url:"currency,omitempty"`
	Email              string `json:"email,omitempty" path:"email,omitempty" url:"email,omitempty"`
	Language           string `json:"language,omitempty" path:"language,omitempty" url:"language,omitempty"`
	PhoneNumber        string `json:"phone_number,omitempty" path:"phone_number,omitempty" url:"phone_number,omitempty"`
	State              string `json:"state,omitempty" path:"state,omitempty" url:"state,omitempty"`
	Zip                string `json:"zip,omitempty" path:"zip,omitempty" url:"zip,omitempty"`
	LeadLevel          string `json:"lead_level,omitempty" path:"lead_level,omitempty" url:"lead_level,omitempty"`
	ClickCookieCode    string `json:"click_cookie_code,omitempty" path:"click_cookie_code,omitempty" url:"click_cookie_code,omitempty"`
	FormName           string `json:"form_name,omitempty" path:"form_name,omitempty" url:"form_name,omitempty"`
	OpportunityComment string `json:"opportunity_comment,omitempty" path:"opportunity_comment,omitempty" url:"opportunity_comment,omitempty"`
	OpportunityType    string `json:"opportunity_type,omitempty" path:"opportunity_type,omitempty" url:"opportunity_type,omitempty"`
	Gclid              string `json:"gclid,omitempty" path:"gclid,omitempty" url:"gclid,omitempty"`
	OriginalBrand      string `json:"original_brand,omitempty" path:"original_brand,omitempty" url:"original_brand,omitempty"`
	UtmCampaign        string `json:"utm_campaign,omitempty" path:"utm_campaign,omitempty" url:"utm_campaign,omitempty"`
	UtmContent         string `json:"utm_content,omitempty" path:"utm_content,omitempty" url:"utm_content,omitempty"`
	UtmDomain          string `json:"utm_domain,omitempty" path:"utm_domain,omitempty" url:"utm_domain,omitempty"`
	UtmMedium          string `json:"utm_medium,omitempty" path:"utm_medium,omitempty" url:"utm_medium,omitempty"`
	UtmSource          string `json:"utm_source,omitempty" path:"utm_source,omitempty" url:"utm_source,omitempty"`
	UtmTerm            string `json:"utm_term,omitempty" path:"utm_term,omitempty" url:"utm_term,omitempty"`
	TimeZone           string `json:"time_zone,omitempty" path:"time_zone,omitempty" url:"time_zone,omitempty"`
	TimeZoneOffset     int64  `json:"time_zone_offset,omitempty" path:"time_zone_offset,omitempty" url:"time_zone_offset,omitempty"`
}

func (l Lead) Identifier() interface{} {
	return l.Id
}

type LeadCollection []Lead

type LeadCreateParams struct {
	Name               string `url:"name,omitempty" required:"false" json:"name,omitempty" path:"name"`
	Address            string `url:"address,omitempty" required:"false" json:"address,omitempty" path:"address"`
	Address2           string `url:"address_2,omitempty" required:"false" json:"address_2,omitempty" path:"address_2"`
	City               string `url:"city,omitempty" required:"false" json:"city,omitempty" path:"city"`
	ContactName        string `url:"contact_name,omitempty" required:"false" json:"contact_name,omitempty" path:"contact_name"`
	Currency           string `url:"currency,omitempty" required:"false" json:"currency,omitempty" path:"currency"`
	Email              string `url:"email,omitempty" required:"false" json:"email,omitempty" path:"email"`
	Language           string `url:"language,omitempty" required:"false" json:"language,omitempty" path:"language"`
	PhoneNumber        string `url:"phone_number,omitempty" required:"false" json:"phone_number,omitempty" path:"phone_number"`
	State              string `url:"state,omitempty" required:"false" json:"state,omitempty" path:"state"`
	Zip                string `url:"zip,omitempty" required:"false" json:"zip,omitempty" path:"zip"`
	ClickCookieCode    string `url:"click_cookie_code,omitempty" required:"false" json:"click_cookie_code,omitempty" path:"click_cookie_code"`
	FormName           string `url:"form_name,omitempty" required:"false" json:"form_name,omitempty" path:"form_name"`
	OpportunityComment string `url:"opportunity_comment,omitempty" required:"false" json:"opportunity_comment,omitempty" path:"opportunity_comment"`
	OpportunityType    string `url:"opportunity_type,omitempty" required:"false" json:"opportunity_type,omitempty" path:"opportunity_type"`
	Gclid              string `url:"gclid,omitempty" required:"false" json:"gclid,omitempty" path:"gclid"`
	OriginalBrand      string `url:"original_brand,omitempty" required:"false" json:"original_brand,omitempty" path:"original_brand"`
	UtmCampaign        string `url:"utm_campaign,omitempty" required:"false" json:"utm_campaign,omitempty" path:"utm_campaign"`
	UtmContent         string `url:"utm_content,omitempty" required:"false" json:"utm_content,omitempty" path:"utm_content"`
	UtmDomain          string `url:"utm_domain,omitempty" required:"false" json:"utm_domain,omitempty" path:"utm_domain"`
	UtmMedium          string `url:"utm_medium,omitempty" required:"false" json:"utm_medium,omitempty" path:"utm_medium"`
	UtmSource          string `url:"utm_source,omitempty" required:"false" json:"utm_source,omitempty" path:"utm_source"`
	UtmTerm            string `url:"utm_term,omitempty" required:"false" json:"utm_term,omitempty" path:"utm_term"`
	TimeZone           string `url:"time_zone,omitempty" required:"false" json:"time_zone,omitempty" path:"time_zone"`
	TimeZoneOffset     int64  `url:"time_zone_offset,omitempty" required:"false" json:"time_zone_offset,omitempty" path:"time_zone_offset"`
}

type LeadUpdateParams struct {
	Code               string `url:"-,omitempty" required:"false" json:"-,omitempty" path:"code"`
	Name               string `url:"name,omitempty" required:"false" json:"name,omitempty" path:"name"`
	Address            string `url:"address,omitempty" required:"false" json:"address,omitempty" path:"address"`
	Address2           string `url:"address_2,omitempty" required:"false" json:"address_2,omitempty" path:"address_2"`
	City               string `url:"city,omitempty" required:"false" json:"city,omitempty" path:"city"`
	ContactName        string `url:"contact_name,omitempty" required:"false" json:"contact_name,omitempty" path:"contact_name"`
	Currency           string `url:"currency,omitempty" required:"false" json:"currency,omitempty" path:"currency"`
	Email              string `url:"email,omitempty" required:"false" json:"email,omitempty" path:"email"`
	Language           string `url:"language,omitempty" required:"false" json:"language,omitempty" path:"language"`
	PhoneNumber        string `url:"phone_number,omitempty" required:"false" json:"phone_number,omitempty" path:"phone_number"`
	State              string `url:"state,omitempty" required:"false" json:"state,omitempty" path:"state"`
	Zip                string `url:"zip,omitempty" required:"false" json:"zip,omitempty" path:"zip"`
	ClickCookieCode    string `url:"click_cookie_code,omitempty" required:"false" json:"click_cookie_code,omitempty" path:"click_cookie_code"`
	FormName           string `url:"form_name,omitempty" required:"false" json:"form_name,omitempty" path:"form_name"`
	OpportunityComment string `url:"opportunity_comment,omitempty" required:"false" json:"opportunity_comment,omitempty" path:"opportunity_comment"`
	OpportunityType    string `url:"opportunity_type,omitempty" required:"false" json:"opportunity_type,omitempty" path:"opportunity_type"`
	Gclid              string `url:"gclid,omitempty" required:"false" json:"gclid,omitempty" path:"gclid"`
	OriginalBrand      string `url:"original_brand,omitempty" required:"false" json:"original_brand,omitempty" path:"original_brand"`
	UtmCampaign        string `url:"utm_campaign,omitempty" required:"false" json:"utm_campaign,omitempty" path:"utm_campaign"`
	UtmContent         string `url:"utm_content,omitempty" required:"false" json:"utm_content,omitempty" path:"utm_content"`
	UtmDomain          string `url:"utm_domain,omitempty" required:"false" json:"utm_domain,omitempty" path:"utm_domain"`
	UtmMedium          string `url:"utm_medium,omitempty" required:"false" json:"utm_medium,omitempty" path:"utm_medium"`
	UtmSource          string `url:"utm_source,omitempty" required:"false" json:"utm_source,omitempty" path:"utm_source"`
	UtmTerm            string `url:"utm_term,omitempty" required:"false" json:"utm_term,omitempty" path:"utm_term"`
	TimeZone           string `url:"time_zone,omitempty" required:"false" json:"time_zone,omitempty" path:"time_zone"`
	TimeZoneOffset     int64  `url:"time_zone_offset,omitempty" required:"false" json:"time_zone_offset,omitempty" path:"time_zone_offset"`
}

func (l *Lead) UnmarshalJSON(data []byte) error {
	type lead Lead
	var v lead
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*l = Lead(v)
	return nil
}

func (l *LeadCollection) UnmarshalJSON(data []byte) error {
	type leads LeadCollection
	var v leads
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*l = LeadCollection(v)
	return nil
}

func (l *LeadCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*l))
	for i, v := range *l {
		ret[i] = v
	}

	return &ret
}
