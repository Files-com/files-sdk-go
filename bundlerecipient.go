package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/lib"
)

type BundleRecipient struct {
	Company   string    `json:"company,omitempty"`
	Name      string    `json:"name,omitempty"`
	Note      string    `json:"note,omitempty"`
	Recipient string    `json:"recipient,omitempty"`
	SentAt    time.Time `json:"sent_at,omitempty"`
}

type BundleRecipientCollection []BundleRecipient

type BundleRecipientListParams struct {
	UserId   int64  `url:"user_id,omitempty"`
	Page     int    `url:"page,omitempty"`
	PerPage  int    `url:"per_page,omitempty"`
	Action   string `url:"action,omitempty"`
	Cursor   string `url:"cursor,omitempty"`
	BundleId int64  `url:"bundle_id,omitempty"`
	lib.ListParams
}

func (b *BundleRecipient) UnmarshalJSON(data []byte) error {
	type bundleRecipient BundleRecipient
	var v bundleRecipient
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*b = BundleRecipient(v)
	return nil
}

func (b *BundleRecipientCollection) UnmarshalJSON(data []byte) error {
	type bundleRecipients []BundleRecipient
	var v bundleRecipients
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*b = BundleRecipientCollection(v)
	return nil
}
