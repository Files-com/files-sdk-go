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
	UserId   int64  `url:"user_id,omitempty" required:"false"`
	Page     int    `url:"page,omitempty" required:"false"`
	PerPage  int    `url:"per_page,omitempty" required:"false"`
	Action   string `url:"action,omitempty" required:"false"`
	Cursor   string `url:"cursor,omitempty" required:"false"`
	BundleId int64  `url:"bundle_id,omitempty" required:"true"`
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
