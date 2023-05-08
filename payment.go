package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Payment struct {
	Id                int64      `json:"id,omitempty" path:"id"`
	Amount            string     `json:"amount,omitempty" path:"amount"`
	Balance           string     `json:"balance,omitempty" path:"balance"`
	CreatedAt         *time.Time `json:"created_at,omitempty" path:"created_at"`
	Currency          string     `json:"currency,omitempty" path:"currency"`
	DownloadUri       string     `json:"download_uri,omitempty" path:"download_uri"`
	InvoiceLineItems  []string   `json:"invoice_line_items,omitempty" path:"invoice_line_items"`
	Method            string     `json:"method,omitempty" path:"method"`
	PaymentLineItems  []string   `json:"payment_line_items,omitempty" path:"payment_line_items"`
	PaymentReversedAt *time.Time `json:"payment_reversed_at,omitempty" path:"payment_reversed_at"`
	PaymentType       string     `json:"payment_type,omitempty" path:"payment_type"`
	SiteName          string     `json:"site_name,omitempty" path:"site_name"`
	Type              string     `json:"type,omitempty" path:"type"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty" path:"updated_at"`
}

func (p Payment) Identifier() interface{} {
	return p.Id
}

type PaymentCollection []Payment

type PaymentListParams struct {
	ListParams
}

type PaymentFindParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
}

func (p *Payment) UnmarshalJSON(data []byte) error {
	type payment Payment
	var v payment
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*p = Payment(v)
	return nil
}

func (p *PaymentCollection) UnmarshalJSON(data []byte) error {
	type payments PaymentCollection
	var v payments
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*p = PaymentCollection(v)
	return nil
}

func (p *PaymentCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*p))
	for i, v := range *p {
		ret[i] = v
	}

	return &ret
}
