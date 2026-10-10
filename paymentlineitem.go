package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// PaymentLineItem is a Files.com API resource.
type PaymentLineItem struct {
	Amount    string     `json:"amount,omitempty" path:"amount,omitempty" url:"amount,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	InvoiceId int64      `json:"invoice_id,omitempty" path:"invoice_id,omitempty" url:"invoice_id,omitempty"`
	PaymentId int64      `json:"payment_id,omitempty" path:"payment_id,omitempty" url:"payment_id,omitempty"`
}

// Identifier no path or id

// PaymentLineItemCollection is a list of PaymentLineItem resources.
type PaymentLineItemCollection []PaymentLineItem

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (p *PaymentLineItem) UnmarshalJSON(data []byte) error {
	type paymentLineItem PaymentLineItem
	var v paymentLineItem
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*p = PaymentLineItem(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (p *PaymentLineItemCollection) UnmarshalJSON(data []byte) error {
	type paymentLineItems PaymentLineItemCollection
	var v paymentLineItems
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*p = PaymentLineItemCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (p *PaymentLineItemCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*p))
	for i, v := range *p {
		ret[i] = v
	}

	return &ret
}
