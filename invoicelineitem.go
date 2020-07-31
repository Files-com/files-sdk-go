package files_sdk

import (
	"encoding/json"
	"time"
)

type InvoiceLineItem struct {
	Amount         float32   `json:"amount,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
	Description    string    `json:"description,omitempty"`
	Type           string    `json:"type,omitempty"`
	ServiceEndAt   time.Time `json:"service_end_at,omitempty"`
	ServiceStartAt time.Time `json:"service_start_at,omitempty"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
	Plan           string    `json:"plan,omitempty"`
	Site           string    `json:"site,omitempty"`
}

type InvoiceLineItemCollection []InvoiceLineItem

func (i *InvoiceLineItem) UnmarshalJSON(data []byte) error {
	type invoiceLineItem InvoiceLineItem
	var v invoiceLineItem
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*i = InvoiceLineItem(v)
	return nil
}

func (i *InvoiceLineItemCollection) UnmarshalJSON(data []byte) error {
	type invoiceLineItems []InvoiceLineItem
	var v invoiceLineItems
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*i = InvoiceLineItemCollection(v)
	return nil
}
