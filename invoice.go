package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Invoice struct {
	Id                int64           `json:"id,omitempty" path:"id"`
	Amount            string          `json:"amount,omitempty" path:"amount"`
	Balance           string          `json:"balance,omitempty" path:"balance"`
	CreatedAt         *time.Time      `json:"created_at,omitempty" path:"created_at"`
	Currency          string          `json:"currency,omitempty" path:"currency"`
	DownloadUri       string          `json:"download_uri,omitempty" path:"download_uri"`
	InvoiceLineItems  InvoiceLineItem `json:"invoice_line_items,omitempty" path:"invoice_line_items"`
	Method            string          `json:"method,omitempty" path:"method"`
	PaymentLineItems  PaymentLineItem `json:"payment_line_items,omitempty" path:"payment_line_items"`
	PaymentReversedAt *time.Time      `json:"payment_reversed_at,omitempty" path:"payment_reversed_at"`
	PaymentType       string          `json:"payment_type,omitempty" path:"payment_type"`
	SiteName          string          `json:"site_name,omitempty" path:"site_name"`
	Type              string          `json:"type,omitempty" path:"type"`
	UpdatedAt         *time.Time      `json:"updated_at,omitempty" path:"updated_at"`
}

type InvoiceCollection []Invoice

type InvoiceListParams struct {
	lib.ListParams
}

type InvoiceFindParams struct {
	Id int64 `url:"-,omitempty" required:"true" json:"-,omitempty" path:"id"`
}

func (i *Invoice) UnmarshalJSON(data []byte) error {
	type invoice Invoice
	var v invoice
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*i = Invoice(v)
	return nil
}

func (i *InvoiceCollection) UnmarshalJSON(data []byte) error {
	type invoices InvoiceCollection
	var v invoices
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*i = InvoiceCollection(v)
	return nil
}

func (i *InvoiceCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*i))
	for i, v := range *i {
		ret[i] = v
	}

	return &ret
}
