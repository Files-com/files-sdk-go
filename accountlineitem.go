package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// AccountLineItem is a Files.com API resource.
type AccountLineItem struct {
	// Line item Id.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Line item amount.
	Amount string `json:"amount,omitempty" path:"amount,omitempty" url:"amount,omitempty"`
	// Line item balance.
	Balance string `json:"balance,omitempty" path:"balance,omitempty" url:"balance,omitempty"`
	// Line item created at.
	CreatedAt *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	// Line item currency.
	Currency string `json:"currency,omitempty" path:"currency,omitempty" url:"currency,omitempty"`
	// Line item download uri.
	DownloadUri string `json:"download_uri,omitempty" path:"download_uri,omitempty" url:"download_uri,omitempty"`
	// Associated invoice line items.
	InvoiceLineItems []InvoiceLineItem `json:"invoice_line_items,omitempty" path:"invoice_line_items,omitempty" url:"invoice_line_items,omitempty"`
	// Line item payment method.
	Method string `json:"method,omitempty" path:"method,omitempty" url:"method,omitempty"`
	// Associated payment line items.
	PaymentLineItems []PaymentLineItem `json:"payment_line_items,omitempty" path:"payment_line_items,omitempty" url:"payment_line_items,omitempty"`
	// Date/time payment was reversed if applicable.
	PaymentReversedAt *time.Time `json:"payment_reversed_at,omitempty" path:"payment_reversed_at,omitempty" url:"payment_reversed_at,omitempty"`
	// Type of payment if applicable.
	PaymentType string `json:"payment_type,omitempty" path:"payment_type,omitempty" url:"payment_type,omitempty"`
	// Site name this line item is for.
	SiteName string `json:"site_name,omitempty" path:"site_name,omitempty" url:"site_name,omitempty"`
	// Type of line item, either payment or invoice.
	Type string `json:"type,omitempty" path:"type,omitempty" url:"type,omitempty"`
}

// Identifier returns the resource ID.
func (a AccountLineItem) Identifier() interface{} {
	return a.Id
}

// AccountLineItemCollection is a list of AccountLineItem resources.
type AccountLineItemCollection []AccountLineItem

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (a *AccountLineItem) UnmarshalJSON(data []byte) error {
	type accountLineItem AccountLineItem
	var v accountLineItem
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*a = AccountLineItem(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (a *AccountLineItemCollection) UnmarshalJSON(data []byte) error {
	type accountLineItems AccountLineItemCollection
	var v accountLineItems
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*a = AccountLineItemCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (a *AccountLineItemCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*a))
	for i, v := range *a {
		ret[i] = v
	}

	return &ret
}
