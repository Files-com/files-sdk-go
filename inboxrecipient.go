package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type InboxRecipient struct {
	Company          string     `json:"company,omitempty" path:"company"`
	Name             string     `json:"name,omitempty" path:"name"`
	Note             string     `json:"note,omitempty" path:"note"`
	Recipient        string     `json:"recipient,omitempty" path:"recipient"`
	SentAt           *time.Time `json:"sent_at,omitempty" path:"sent_at"`
	UserId           int64      `json:"user_id,omitempty" path:"user_id"`
	InboxId          int64      `json:"inbox_id,omitempty" path:"inbox_id"`
	ShareAfterCreate *bool      `json:"share_after_create,omitempty" path:"share_after_create"`
}

type InboxRecipientCollection []InboxRecipient

type InboxRecipientListParams struct {
	UserId     int64           `url:"user_id,omitempty" required:"false" json:"user_id,omitempty" path:"user_id"`
	SortBy     json.RawMessage `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty" path:"sort_by"`
	Filter     json.RawMessage `url:"filter,omitempty" required:"false" json:"filter,omitempty" path:"filter"`
	FilterGt   json.RawMessage `url:"filter_gt,omitempty" required:"false" json:"filter_gt,omitempty" path:"filter_gt"`
	FilterGteq json.RawMessage `url:"filter_gteq,omitempty" required:"false" json:"filter_gteq,omitempty" path:"filter_gteq"`
	FilterLike json.RawMessage `url:"filter_like,omitempty" required:"false" json:"filter_like,omitempty" path:"filter_like"`
	FilterLt   json.RawMessage `url:"filter_lt,omitempty" required:"false" json:"filter_lt,omitempty" path:"filter_lt"`
	FilterLteq json.RawMessage `url:"filter_lteq,omitempty" required:"false" json:"filter_lteq,omitempty" path:"filter_lteq"`
	InboxId    int64           `url:"inbox_id,omitempty" required:"true" json:"inbox_id,omitempty" path:"inbox_id"`
	lib.ListParams
}

type InboxRecipientCreateParams struct {
	UserId           int64  `url:"user_id,omitempty" required:"false" json:"user_id,omitempty" path:"user_id"`
	InboxId          int64  `url:"inbox_id,omitempty" required:"true" json:"inbox_id,omitempty" path:"inbox_id"`
	Recipient        string `url:"recipient,omitempty" required:"true" json:"recipient,omitempty" path:"recipient"`
	Name             string `url:"name,omitempty" required:"false" json:"name,omitempty" path:"name"`
	Company          string `url:"company,omitempty" required:"false" json:"company,omitempty" path:"company"`
	Note             string `url:"note,omitempty" required:"false" json:"note,omitempty" path:"note"`
	ShareAfterCreate *bool  `url:"share_after_create,omitempty" required:"false" json:"share_after_create,omitempty" path:"share_after_create"`
}

func (i *InboxRecipient) UnmarshalJSON(data []byte) error {
	type inboxRecipient InboxRecipient
	var v inboxRecipient
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*i = InboxRecipient(v)
	return nil
}

func (i *InboxRecipientCollection) UnmarshalJSON(data []byte) error {
	type inboxRecipients InboxRecipientCollection
	var v inboxRecipients
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
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
