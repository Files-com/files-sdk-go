package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// InvoiceLineItem is a Files.com API resource.
type InvoiceLineItem struct {
	// Invoice Line item Id.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Invoice line item amount.
	Amount string `json:"amount,omitempty" path:"amount,omitempty" url:"amount,omitempty"`
	// Invoice line item created at date/time.
	CreatedAt *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	// Invoice line item description.
	Description string `json:"description,omitempty" path:"description,omitempty" url:"description,omitempty"`
	// Invoice line item type.
	Type string `json:"type,omitempty" path:"type,omitempty" url:"type,omitempty"`
	// Invoice line item service end date/time.
	ServiceEndAt *time.Time `json:"service_end_at,omitempty" path:"service_end_at,omitempty" url:"service_end_at,omitempty"`
	// Invoice line item service start date/time.
	ServiceStartAt *time.Time `json:"service_start_at,omitempty" path:"service_start_at,omitempty" url:"service_start_at,omitempty"`
	// Plan name.
	Plan string `json:"plan,omitempty" path:"plan,omitempty" url:"plan,omitempty"`
	// Site name.
	Site string `json:"site,omitempty" path:"site,omitempty" url:"site,omitempty"`
	// Prepaid bytes purchased for this invoice line item.
	PrepaidBytes int64 `json:"prepaid_bytes,omitempty" path:"prepaid_bytes,omitempty" url:"prepaid_bytes,omitempty"`
	// When the prepaid bytes expire.
	PrepaidBytesExpireAt *time.Time `json:"prepaid_bytes_expire_at,omitempty" path:"prepaid_bytes_expire_at,omitempty" url:"prepaid_bytes_expire_at,omitempty"`
	// Total prepaid bytes used for this invoice line item.
	PrepaidBytesUsed int64 `json:"prepaid_bytes_used,omitempty" path:"prepaid_bytes_used,omitempty" url:"prepaid_bytes_used,omitempty"`
	// Available prepaid bytes for this invoice line item.
	PrepaidBytesAvailable int64 `json:"prepaid_bytes_available,omitempty" path:"prepaid_bytes_available,omitempty" url:"prepaid_bytes_available,omitempty"`
}

// Identifier returns the resource ID.
func (i InvoiceLineItem) Identifier() interface{} {
	return i.Id
}

// InvoiceLineItemCollection is a list of InvoiceLineItem resources.
type InvoiceLineItemCollection []InvoiceLineItem

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (i *InvoiceLineItem) UnmarshalJSON(data []byte) error {
	type invoiceLineItem InvoiceLineItem
	var v invoiceLineItem
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*i = InvoiceLineItem(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (i *InvoiceLineItemCollection) UnmarshalJSON(data []byte) error {
	type invoiceLineItems InvoiceLineItemCollection
	var v invoiceLineItems
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*i = InvoiceLineItemCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (i *InvoiceLineItemCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*i))
	for i, v := range *i {
		ret[i] = v
	}

	return &ret
}
