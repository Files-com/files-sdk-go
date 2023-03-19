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
	InboxId          int64      `json:"inbox_id,omitempty" path:"inbox_id"`
	ShareAfterCreate *bool      `json:"share_after_create,omitempty" path:"share_after_create"`
}

type InboxRecipientCollection []InboxRecipient

type InboxRecipientListParams struct {
	SortBy  json.RawMessage `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty" path:"sort_by"`
	Filter  json.RawMessage `url:"filter,omitempty" required:"false" json:"filter,omitempty" path:"filter"`
	InboxId int64           `url:"inbox_id,omitempty" required:"true" json:"inbox_id,omitempty" path:"inbox_id"`
	lib.ListParams
}

type InboxRecipientCreateParams struct {
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
