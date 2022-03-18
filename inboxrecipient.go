package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type InboxRecipient struct {
	Company          string    `json:"company,omitempty"`
	Name             string    `json:"name,omitempty"`
	Note             string    `json:"note,omitempty"`
	Recipient        string    `json:"recipient,omitempty"`
	SentAt           time.Time `json:"sent_at,omitempty"`
	UserId           int64     `json:"user_id,omitempty"`
	InboxId          int64     `json:"inbox_id,omitempty"`
	ShareAfterCreate *bool     `json:"share_after_create,omitempty"`
}

type InboxRecipientCollection []InboxRecipient

type InboxRecipientListParams struct {
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
	InboxId    int64           `url:"inbox_id,omitempty" required:"true" json:"inbox_id,omitempty"`
	lib.ListParams
}

type InboxRecipientCreateParams struct {
	UserId           int64  `url:"user_id,omitempty" required:"false" json:"user_id,omitempty"`
	InboxId          int64  `url:"inbox_id,omitempty" required:"true" json:"inbox_id,omitempty"`
	Recipient        string `url:"recipient,omitempty" required:"true" json:"recipient,omitempty"`
	Name             string `url:"name,omitempty" required:"false" json:"name,omitempty"`
	Company          string `url:"company,omitempty" required:"false" json:"company,omitempty"`
	Note             string `url:"note,omitempty" required:"false" json:"note,omitempty"`
	ShareAfterCreate *bool  `url:"share_after_create,omitempty" required:"false" json:"share_after_create,omitempty"`
}

func (i *InboxRecipient) UnmarshalJSON(data []byte) error {
	type inboxRecipient InboxRecipient
	var v inboxRecipient
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*i = InboxRecipient(v)
	return nil
}

func (i *InboxRecipientCollection) UnmarshalJSON(data []byte) error {
	type inboxRecipients []InboxRecipient
	var v inboxRecipients
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*i = InboxRecipientCollection(v)
	return nil
}

func (i *InboxRecipientCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*i))
	for i, v := range *i {
		ret[i] = v
	}

	return &ret
}
