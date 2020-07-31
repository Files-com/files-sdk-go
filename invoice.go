package files_sdk

import (
  lib "github.com/Files-com/files-sdk-go/lib"
  "encoding/json"
  "time"
)

type Invoice struct {
  Id int `json:"id,omitempty"`
  Amount float32 `json:"amount,omitempty"`
  Balance float32 `json:"balance,omitempty"`
  CreatedAt time.Time `json:"created_at,omitempty"`
  Currency string `json:"currency,omitempty"`
  DownloadUri string `json:"download_uri,omitempty"`
  InvoiceLineItems []string `json:"invoice_line_items,omitempty"`
  Method string `json:"method,omitempty"`
  PaymentLineItems []string `json:"payment_line_items,omitempty"`
  PaymentReversedAt time.Time `json:"payment_reversed_at,omitempty"`
  PaymentType string `json:"payment_type,omitempty"`
  SiteName string `json:"site_name,omitempty"`
  Type string `json:"type,omitempty"`
  UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type InvoiceCollection []Invoice

type InvoiceListParams struct {
  Page int `url:"page,omitempty"`
  PerPage int `url:"per_page,omitempty"`
  Action string `url:"action,omitempty"`
  lib.ListParams
}

type InvoiceFindParams struct {
  Id int `url:"-,omitempty"`
}


func (i *Invoice) UnmarshalJSON(data []byte) error {
	type invoice Invoice
	var v invoice
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*i = Invoice(v)
	return nil
}

func (i *InvoiceCollection) UnmarshalJSON(data []byte) error {
	type invoices []Invoice
	var v invoices
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*i = InvoiceCollection(v)
	return nil
}

