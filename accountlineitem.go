package files_sdk

import (
	"encoding/json"
	"time"
)

type AccountLineItem struct {
	Id                int64           `json:"id,omitempty"`
	Amount            string          `json:"amount,omitempty"`
	Balance           string          `json:"balance,omitempty"`
	CreatedAt         *time.Time      `json:"created_at,omitempty"`
	Currency          string          `json:"currency,omitempty"`
	DownloadUri       string          `json:"download_uri,omitempty"`
	InvoiceLineItems  InvoiceLineItem `json:"invoice_line_items,omitempty"`
	Method            string          `json:"method,omitempty"`
	PaymentLineItems  PaymentLineItem `json:"payment_line_items,omitempty"`
	PaymentReversedAt *time.Time      `json:"payment_reversed_at,omitempty"`
	PaymentType       string          `json:"payment_type,omitempty"`
	SiteName          string          `json:"site_name,omitempty"`
	Type              string          `json:"type,omitempty"`
	UpdatedAt         *time.Time      `json:"updated_at,omitempty"`
}

type AccountLineItemCollection []AccountLineItem

func (a *AccountLineItem) UnmarshalJSON(data []byte) error {
	type accountLineItem AccountLineItem
	var v accountLineItem
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*a = AccountLineItem(v)
	return nil
}

func (a *AccountLineItemCollection) UnmarshalJSON(data []byte) error {
	type accountLineItems []AccountLineItem
	var v accountLineItems
	if err := json.Unmarshal(data, &v); err != nil {
		return err
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
