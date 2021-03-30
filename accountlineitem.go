package files_sdk

import (
	"encoding/json"
	"time"
)

type AccountLineItem struct {
	Id                int64     `json:"id,omitempty"`
	Amount            float32   `json:"amount,omitempty"`
	Balance           float32   `json:"balance,omitempty"`
	CreatedAt         time.Time `json:"created_at,omitempty"`
	Currency          string    `json:"currency,omitempty"`
	DownloadUri       string    `json:"download_uri,omitempty"`
	InvoiceLineItems  []string  `json:"invoice_line_items,omitempty"`
	Method            string    `json:"method,omitempty"`
	PaymentLineItems  []string  `json:"payment_line_items,omitempty"`
	PaymentReversedAt time.Time `json:"payment_reversed_at,omitempty"`
	PaymentType       string    `json:"payment_type,omitempty"`
	SiteName          string    `json:"site_name,omitempty"`
	Type              string    `json:"type,omitempty"`
	UpdatedAt         time.Time `json:"updated_at,omitempty"`
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
