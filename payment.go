package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// Payment is a Files.com API resource.
type Payment struct {
	Id                int64             `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	Amount            string            `json:"amount,omitempty" path:"amount,omitempty" url:"amount,omitempty"`
	Balance           string            `json:"balance,omitempty" path:"balance,omitempty" url:"balance,omitempty"`
	CreatedAt         *time.Time        `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	Currency          string            `json:"currency,omitempty" path:"currency,omitempty" url:"currency,omitempty"`
	DownloadUri       string            `json:"download_uri,omitempty" path:"download_uri,omitempty" url:"download_uri,omitempty"`
	InvoiceLineItems  []InvoiceLineItem `json:"invoice_line_items,omitempty" path:"invoice_line_items,omitempty" url:"invoice_line_items,omitempty"`
	Method            string            `json:"method,omitempty" path:"method,omitempty" url:"method,omitempty"`
	PaymentLineItems  []PaymentLineItem `json:"payment_line_items,omitempty" path:"payment_line_items,omitempty" url:"payment_line_items,omitempty"`
	PaymentReversedAt *time.Time        `json:"payment_reversed_at,omitempty" path:"payment_reversed_at,omitempty" url:"payment_reversed_at,omitempty"`
	PaymentType       string            `json:"payment_type,omitempty" path:"payment_type,omitempty" url:"payment_type,omitempty"`
	SiteName          string            `json:"site_name,omitempty" path:"site_name,omitempty" url:"site_name,omitempty"`
	Type              string            `json:"type,omitempty" path:"type,omitempty" url:"type,omitempty"`
}

// Identifier returns the resource ID.
func (p Payment) Identifier() interface{} {
	return p.Id
}

// PaymentCollection is a list of Payment resources.
type PaymentCollection []Payment

// PaymentListParams contains the request parameters for GET /payments.
type PaymentListParams struct {
	ListParams
}

// PaymentFindParams contains the request parameters for GET /payments/{id}.
type PaymentFindParams struct {
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (p *Payment) UnmarshalJSON(data []byte) error {
	type payment Payment
	var v payment
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*p = Payment(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (p *PaymentCollection) UnmarshalJSON(data []byte) error {
	type payments PaymentCollection
	var v payments
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*p = PaymentCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (p *PaymentCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*p))
	for i, v := range *p {
		ret[i] = v
	}

	return &ret
}
