package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

type Invoice struct {
	Id                int64                    `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	Amount            string                   `json:"amount,omitempty" path:"amount,omitempty" url:"amount,omitempty"`
	Balance           string                   `json:"balance,omitempty" path:"balance,omitempty" url:"balance,omitempty"`
	CreatedAt         *time.Time               `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	Currency          string                   `json:"currency,omitempty" path:"currency,omitempty" url:"currency,omitempty"`
	DownloadUri       string                   `json:"download_uri,omitempty" path:"download_uri,omitempty" url:"download_uri,omitempty"`
	InvoiceLineItems  []map[string]interface{} `json:"invoice_line_items,omitempty" path:"invoice_line_items,omitempty" url:"invoice_line_items,omitempty"`
	Method            string                   `json:"method,omitempty" path:"method,omitempty" url:"method,omitempty"`
	PaymentLineItems  []map[string]interface{} `json:"payment_line_items,omitempty" path:"payment_line_items,omitempty" url:"payment_line_items,omitempty"`
	PaymentReversedAt *time.Time               `json:"payment_reversed_at,omitempty" path:"payment_reversed_at,omitempty" url:"payment_reversed_at,omitempty"`
	PaymentType       string                   `json:"payment_type,omitempty" path:"payment_type,omitempty" url:"payment_type,omitempty"`
	SiteName          string                   `json:"site_name,omitempty" path:"site_name,omitempty" url:"site_name,omitempty"`
	Type              string                   `json:"type,omitempty" path:"type,omitempty" url:"type,omitempty"`
}

func (i Invoice) Identifier() interface{} {
	return i.Id
}

type InvoiceCollection []Invoice

type InvoiceListParams struct {
	ListParams
}

type InvoiceFindParams struct {
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
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
