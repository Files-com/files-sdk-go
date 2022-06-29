package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type BundleRecipient struct {
	Company          string     `json:"company,omitempty"`
	Name             string     `json:"name,omitempty"`
	Note             string     `json:"note,omitempty"`
	Recipient        string     `json:"recipient,omitempty"`
	SentAt           *time.Time `json:"sent_at,omitempty"`
	UserId           int64      `json:"user_id,omitempty"`
	BundleId         int64      `json:"bundle_id,omitempty"`
	ShareAfterCreate *bool      `json:"share_after_create,omitempty"`
}

type BundleRecipientCollection []BundleRecipient

type BundleRecipientListParams struct {
	UserId     int64           `url:"user_id,omitempty" required:"false" json:"user_id,omitempty"`
	Cursor     string          `url:"cursor,omitempty" required:"false" json:"cursor,omitempty"`
	PerPage    int64           `url:"per_page,omitempty" required:"false" json:"per_page,omitempty"`
	SortBy     json.RawMessage `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty"`
	Filter     json.RawMessage `url:"filter,omitempty" required:"false" json:"filter,omitempty"`
	FilterGt   json.RawMessage `url:"filter_gt,omitempty" required:"false" json:"filter_gt,omitempty"`
	FilterGteq json.RawMessage `url:"filter_gteq,omitempty" required:"false" json:"filter_gteq,omitempty"`
	FilterLike json.RawMessage `url:"filter_like,omitempty" required:"false" json:"filter_like,omitempty"`
	FilterLt   json.RawMessage `url:"filter_lt,omitempty" required:"false" json:"filter_lt,omitempty"`
	FilterLteq json.RawMessage `url:"filter_lteq,omitempty" required:"false" json:"filter_lteq,omitempty"`
	BundleId   int64           `url:"bundle_id,omitempty" required:"true" json:"bundle_id,omitempty"`
	lib.ListParams
}

type BundleRecipientCreateParams struct {
	UserId           int64  `url:"user_id,omitempty" required:"false" json:"user_id,omitempty"`
	BundleId         int64  `url:"bundle_id,omitempty" required:"true" json:"bundle_id,omitempty"`
	Recipient        string `url:"recipient,omitempty" required:"true" json:"recipient,omitempty"`
	Name             string `url:"name,omitempty" required:"false" json:"name,omitempty"`
	Company          string `url:"company,omitempty" required:"false" json:"company,omitempty"`
	Note             string `url:"note,omitempty" required:"false" json:"note,omitempty"`
	ShareAfterCreate *bool  `url:"share_after_create,omitempty" required:"false" json:"share_after_create,omitempty"`
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
