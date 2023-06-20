package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type AccountLineItem struct {
	Id                int64      `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	Amount            string     `json:"amount,omitempty" path:"amount,omitempty" url:"amount,omitempty"`
	Balance           string     `json:"balance,omitempty" path:"balance,omitempty" url:"balance,omitempty"`
	CreatedAt         *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	Currency          string     `json:"currency,omitempty" path:"currency,omitempty" url:"currency,omitempty"`
	DownloadUri       string     `json:"download_uri,omitempty" path:"download_uri,omitempty" url:"download_uri,omitempty"`
	InvoiceLineItems  []string   `json:"invoice_line_items,omitempty" path:"invoice_line_items,omitempty" url:"invoice_line_items,omitempty"`
	Method            string     `json:"method,omitempty" path:"method,omitempty" url:"method,omitempty"`
	PaymentLineItems  []string   `json:"payment_line_items,omitempty" path:"payment_line_items,omitempty" url:"payment_line_items,omitempty"`
	PaymentReversedAt *time.Time `json:"payment_reversed_at,omitempty" path:"payment_reversed_at,omitempty" url:"payment_reversed_at,omitempty"`
	PaymentType       string     `json:"payment_type,omitempty" path:"payment_type,omitempty" url:"payment_type,omitempty"`
	SiteName          string     `json:"site_name,omitempty" path:"site_name,omitempty" url:"site_name,omitempty"`
	Type              string     `json:"type,omitempty" path:"type,omitempty" url:"type,omitempty"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty" path:"updated_at,omitempty" url:"updated_at,omitempty"`
}

func (a AccountLineItem) Identifier() interface{} {
	return a.Id
}

type AccountLineItemCollection []AccountLineItem

func (a *AccountLineItem) UnmarshalJSON(data []byte) error {
	type accountLineItem AccountLineItem
	var v accountLineItem
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*a = AccountLineItem(v)
	return nil
}

func (a *AccountLineItemCollection) UnmarshalJSON(data []byte) error {
	type accountLineItems AccountLineItemCollection
	var v accountLineItems
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*a = AccountLineItemCollection(v)
	return nil
}

func (a *AccountLineItemCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*a))
	for i, v := range *a {
		ret[i] = v
	}

	return &ret
}
