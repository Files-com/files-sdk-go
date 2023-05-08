package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type BundleRecipient struct {
	Company          string     `json:"company,omitempty" path:"company"`
	Name             string     `json:"name,omitempty" path:"name"`
	Note             string     `json:"note,omitempty" path:"note"`
	Recipient        string     `json:"recipient,omitempty" path:"recipient"`
	SentAt           *time.Time `json:"sent_at,omitempty" path:"sent_at"`
	UserId           int64      `json:"user_id,omitempty" path:"user_id"`
	BundleId         int64      `json:"bundle_id,omitempty" path:"bundle_id"`
	ShareAfterCreate *bool      `json:"share_after_create,omitempty" path:"share_after_create"`
}

// Identifier no path or id

type BundleRecipientCollection []BundleRecipient

type BundleRecipientListParams struct {
	UserId   int64           `url:"user_id,omitempty" required:"false" json:"user_id,omitempty" path:"user_id"`
	SortBy   json.RawMessage `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty" path:"sort_by"`
	Filter   json.RawMessage `url:"filter,omitempty" required:"false" json:"filter,omitempty" path:"filter"`
	BundleId int64           `url:"bundle_id,omitempty" required:"true" json:"bundle_id,omitempty" path:"bundle_id"`
	ListParams
}

type BundleRecipientCreateParams struct {
	UserId           int64  `url:"user_id,omitempty" required:"false" json:"user_id,omitempty" path:"user_id"`
	BundleId         int64  `url:"bundle_id,omitempty" required:"true" json:"bundle_id,omitempty" path:"bundle_id"`
	Recipient        string `url:"recipient,omitempty" required:"true" json:"recipient,omitempty" path:"recipient"`
	Name             string `url:"name,omitempty" required:"false" json:"name,omitempty" path:"name"`
	Company          string `url:"company,omitempty" required:"false" json:"company,omitempty" path:"company"`
	Note             string `url:"note,omitempty" required:"false" json:"note,omitempty" path:"note"`
	ShareAfterCreate *bool  `url:"share_after_create,omitempty" required:"false" json:"share_after_create,omitempty" path:"share_after_create"`
}

func (b *BundleRecipient) UnmarshalJSON(data []byte) error {
	type bundleRecipient BundleRecipient
	var v bundleRecipient
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*b = BundleRecipient(v)
	return nil
}

func (b *BundleRecipientCollection) UnmarshalJSON(data []byte) error {
	type bundleRecipients BundleRecipientCollection
	var v bundleRecipients
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
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
