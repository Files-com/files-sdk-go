package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type InvoiceLineItem struct {
	Amount         string     `json:"amount,omitempty" path:"amount,omitempty" url:"amount,omitempty"`
	CreatedAt      *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	Description    string     `json:"description,omitempty" path:"description,omitempty" url:"description,omitempty"`
	Type           string     `json:"type,omitempty" path:"type,omitempty" url:"type,omitempty"`
	ServiceEndAt   *time.Time `json:"service_end_at,omitempty" path:"service_end_at,omitempty" url:"service_end_at,omitempty"`
	ServiceStartAt *time.Time `json:"service_start_at,omitempty" path:"service_start_at,omitempty" url:"service_start_at,omitempty"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty" path:"updated_at,omitempty" url:"updated_at,omitempty"`
	Plan           string     `json:"plan,omitempty" path:"plan,omitempty" url:"plan,omitempty"`
	Site           string     `json:"site,omitempty" path:"site,omitempty" url:"site,omitempty"`
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
