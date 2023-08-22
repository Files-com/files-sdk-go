package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type BundleRecipient struct {
	Company          string     `json:"company,omitempty" path:"company,omitempty" url:"company,omitempty"`
	Name             string     `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	Note             string     `json:"note,omitempty" path:"note,omitempty" url:"note,omitempty"`
	Recipient        string     `json:"recipient,omitempty" path:"recipient,omitempty" url:"recipient,omitempty"`
	SentAt           *time.Time `json:"sent_at,omitempty" path:"sent_at,omitempty" url:"sent_at,omitempty"`
	UserId           int64      `json:"user_id,omitempty" path:"user_id,omitempty" url:"user_id,omitempty"`
	BundleId         int64      `json:"bundle_id,omitempty" path:"bundle_id,omitempty" url:"bundle_id,omitempty"`
	Method           string     `json:"method,omitempty" path:"method,omitempty" url:"method,omitempty"`
	ShareAfterCreate *bool      `json:"share_after_create,omitempty" path:"share_after_create,omitempty" url:"share_after_create,omitempty"`
}

// Identifier no path or id

type BundleRecipientCollection []BundleRecipient

type BundleRecipientMethodEnum string

func (u BundleRecipientMethodEnum) String() string {
	return string(u)
}

func (u BundleRecipientMethodEnum) Enum() map[string]BundleRecipientMethodEnum {
	return map[string]BundleRecipientMethodEnum{
		"email": BundleRecipientMethodEnum("email"),
	}
}

type BundleRecipientListParams struct {
	UserId   int64                  `url:"user_id,omitempty" required:"false" json:"user_id,omitempty" path:"user_id"`
	Action   string                 `url:"action,omitempty" required:"false" json:"action,omitempty" path:"action"`
	SortBy   map[string]interface{} `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty" path:"sort_by"`
	Filter   BundleRecipient        `url:"filter,omitempty" required:"false" json:"filter,omitempty" path:"filter"`
	BundleId int64                  `url:"bundle_id,omitempty" required:"true" json:"bundle_id,omitempty" path:"bundle_id"`
	ListParams
}

type BundleRecipientCreateParams struct {
	UserId           int64                     `url:"user_id,omitempty" required:"false" json:"user_id,omitempty" path:"user_id"`
	BundleId         int64                     `url:"bundle_id,omitempty" required:"true" json:"bundle_id,omitempty" path:"bundle_id"`
	Recipient        string                    `url:"recipient,omitempty" required:"true" json:"recipient,omitempty" path:"recipient"`
	Name             string                    `url:"name,omitempty" required:"false" json:"name,omitempty" path:"name"`
	Company          string                    `url:"company,omitempty" required:"false" json:"company,omitempty" path:"company"`
	Note             string                    `url:"note,omitempty" required:"false" json:"note,omitempty" path:"note"`
	Method           BundleRecipientMethodEnum `url:"method,omitempty" required:"false" json:"method,omitempty" path:"method"`
	ShareAfterCreate *bool                     `url:"share_after_create,omitempty" required:"false" json:"share_after_create,omitempty" path:"share_after_create"`
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
