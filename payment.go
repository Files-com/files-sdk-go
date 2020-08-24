package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/lib"
)

type Payment struct {
	Id                int64     `json:"id,omitempty"`
	Amount            float32   `json:"amount,omitempty"`
	Balance           float32   `json:"balance,omitempty"`
	CreatedAt         time.Time `json:"created_at,omitempty"`
	Currency          string    `json:"currency,omitempty"`
	DownloadUri       string    `json:"download_uri,omitempty"`
	InvoiceLineItems  []string  `json:"invoice_line_items,omitempty"`
	Method            string    `json:"method,omitempty"`
	PaymentLineItems  []string  `json:"payment_line_items,omitempty"`
	PaymentReversedAt time.Time `json:"payment_reversed_at,omitempty"`
	PaymentType       string    `json:"payment_type,omitempty"`
	SiteName          string    `json:"site_name,omitempty"`
	Type              string    `json:"type,omitempty"`
	UpdatedAt         time.Time `json:"updated_at,omitempty"`
}

type PaymentCollection []Payment

type PaymentListParams struct {
	Page    int    `url:"page,omitempty"`
	PerPage int    `url:"per_page,omitempty"`
	Action  string `url:"action,omitempty"`
	Cursor  string `url:"cursor,omitempty"`
	lib.ListParams
}

type PaymentFindParams struct {
	Id int64 `url:"-,omitempty"`
}

func (p *Payment) UnmarshalJSON(data []byte) error {
	type payment Payment
	var v payment
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*p = Payment(v)
	return nil
}

func (p *PaymentCollection) UnmarshalJSON(data []byte) error {
	type payments []Payment
	var v payments
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*p = PaymentCollection(v)
	return nil
}
