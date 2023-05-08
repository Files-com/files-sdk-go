package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type InvoiceLineItem struct {
	Amount         string     `json:"amount,omitempty" path:"amount"`
	CreatedAt      *time.Time `json:"created_at,omitempty" path:"created_at"`
	Description    string     `json:"description,omitempty" path:"description"`
	Type           string     `json:"type,omitempty" path:"type"`
	ServiceEndAt   *time.Time `json:"service_end_at,omitempty" path:"service_end_at"`
	ServiceStartAt *time.Time `json:"service_start_at,omitempty" path:"service_start_at"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty" path:"updated_at"`
	Plan           string     `json:"plan,omitempty" path:"plan"`
	Site           string     `json:"site,omitempty" path:"site"`
}

// Identifier no path or id

type InvoiceLineItemCollection []InvoiceLineItem

func (i *InvoiceLineItem) UnmarshalJSON(data []byte) error {
	type invoiceLineItem InvoiceLineItem
	var v invoiceLineItem
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*i = InvoiceLineItem(v)
	return nil
}

func (i *InvoiceLineItemCollection) UnmarshalJSON(data []byte) error {
	type invoiceLineItems InvoiceLineItemCollection
	var v invoiceLineItems
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*i = InvoiceLineItemCollection(v)
	return nil
}

func (i *InvoiceLineItemCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*i))
	for i, v := range *i {
		ret[i] = v
	}

	return &ret
}
