package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type BundleRecipient struct {
	Company          string    `json:"company,omitempty"`
	Name             string    `json:"name,omitempty"`
	Note             string    `json:"note,omitempty"`
	Recipient        string    `json:"recipient,omitempty"`
	SentAt           time.Time `json:"sent_at,omitempty"`
	UserId           int64     `json:"user_id,omitempty"`
	BundleId         int64     `json:"bundle_id,omitempty"`
	ShareAfterCreate *bool     `json:"share_after_create,omitempty"`
}

type BundleRecipientCollection []BundleRecipient

type BundleRecipientListParams struct {
	UserId     int64           `url:"user_id,omitempty" required:"false"`
	Cursor     string          `url:"cursor,omitempty" required:"false"`
	PerPage    int64           `url:"per_page,omitempty" required:"false"`
	SortBy     json.RawMessage `url:"sort_by,omitempty" required:"false"`
	Filter     json.RawMessage `url:"filter,omitempty" required:"false"`
	FilterGt   json.RawMessage `url:"filter_gt,omitempty" required:"false"`
	FilterGteq json.RawMessage `url:"filter_gteq,omitempty" required:"false"`
	FilterLike json.RawMessage `url:"filter_like,omitempty" required:"false"`
	FilterLt   json.RawMessage `url:"filter_lt,omitempty" required:"false"`
	FilterLteq json.RawMessage `url:"filter_lteq,omitempty" required:"false"`
	BundleId   int64           `url:"bundle_id,omitempty" required:"true"`
	lib.ListParams
}

type BundleRecipientCreateParams struct {
	UserId           int64  `url:"user_id,omitempty" required:"false"`
	BundleId         int64  `url:"bundle_id,omitempty" required:"true"`
	Recipient        string `url:"recipient,omitempty" required:"true"`
	Name             string `url:"name,omitempty" required:"false"`
	Company          string `url:"company,omitempty" required:"false"`
	Note             string `url:"note,omitempty" required:"false"`
	ShareAfterCreate *bool  `url:"share_after_create,omitempty" required:"false"`
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

func (b *BundleRecipientCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*b))
	for i, v := range *b {
		ret[i] = v
	}

	return &ret
}
