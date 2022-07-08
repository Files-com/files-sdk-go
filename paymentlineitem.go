package files_sdk

import (
	"encoding/json"
	"time"
)

type PaymentLineItem struct {
	Amount    string     `json:"amount,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	InvoiceId int64      `json:"invoice_id,omitempty"`
	PaymentId int64      `json:"payment_id,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type PaymentLineItemCollection []PaymentLineItem

func (p *PaymentLineItem) UnmarshalJSON(data []byte) error {
	type paymentLineItem PaymentLineItem
	var v paymentLineItem
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*p = PaymentLineItem(v)
	return nil
}

func (p *PaymentLineItemCollection) UnmarshalJSON(data []byte) error {
	type paymentLineItems []PaymentLineItem
	var v paymentLineItems
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*p = PaymentLineItemCollection(v)
	return nil
}

func (p *PaymentLineItemCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*p))
	for i, v := range *p {
		ret[i] = v
	}

	return &ret
}
