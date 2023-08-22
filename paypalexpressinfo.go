package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type PaypalExpressInfo struct {
	BillingEmail       string `json:"billing_email,omitempty" path:"billing_email,omitempty" url:"billing_email,omitempty"`
	BillingCompanyName string `json:"billing_company_name,omitempty" path:"billing_company_name,omitempty" url:"billing_company_name,omitempty"`
	BillingAddress     string `json:"billing_address,omitempty" path:"billing_address,omitempty" url:"billing_address,omitempty"`
	BillingAddress2    string `json:"billing_address_2,omitempty" path:"billing_address_2,omitempty" url:"billing_address_2,omitempty"`
	BillingCity        string `json:"billing_city,omitempty" path:"billing_city,omitempty" url:"billing_city,omitempty"`
	BillingState       string `json:"billing_state,omitempty" path:"billing_state,omitempty" url:"billing_state,omitempty"`
	BillingCountry     string `json:"billing_country,omitempty" path:"billing_country,omitempty" url:"billing_country,omitempty"`
	BillingZip         string `json:"billing_zip,omitempty" path:"billing_zip,omitempty" url:"billing_zip,omitempty"`
	BillingName        string `json:"billing_name,omitempty" path:"billing_name,omitempty" url:"billing_name,omitempty"`
	BillingPhone       string `json:"billing_phone,omitempty" path:"billing_phone,omitempty" url:"billing_phone,omitempty"`
	PaypalPayerId      int64  `json:"paypal_payer_id,omitempty" path:"paypal_payer_id,omitempty" url:"paypal_payer_id,omitempty"`
}

// Identifier no path or id

type PaypalExpressInfoCollection []PaypalExpressInfo

func (p *PaypalExpressInfo) UnmarshalJSON(data []byte) error {
	type paypalExpressInfo PaypalExpressInfo
	var v paypalExpressInfo
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*p = PaypalExpressInfo(v)
	return nil
}

func (p *PaypalExpressInfoCollection) UnmarshalJSON(data []byte) error {
	type paypalExpressInfos PaypalExpressInfoCollection
	var v paypalExpressInfos
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*p = PaypalExpressInfoCollection(v)
	return nil
}

func (p *PaypalExpressInfoCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*p))
	for i, v := range *p {
		ret[i] = v
	}

	return &ret
}
