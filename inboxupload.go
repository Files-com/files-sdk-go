package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/lib"
)

type InboxUpload struct {
	InboxRegistration string    `json:"inbox_registration,omitempty"`
	Path              string    `json:"path,omitempty"`
	CreatedAt         time.Time `json:"created_at,omitempty"`
}

type InboxUploadCollection []InboxUpload

type InboxUploadListParams struct {
	Cursor              string          `url:"cursor,omitempty" required:"false"`
	PerPage             int             `url:"per_page,omitempty" required:"false"`
	SortBy              json.RawMessage `url:"sort_by,omitempty" required:"false"`
	Filter              json.RawMessage `url:"filter,omitempty" required:"false"`
	FilterGt            json.RawMessage `url:"filter_gt,omitempty" required:"false"`
	FilterGteq          json.RawMessage `url:"filter_gteq,omitempty" required:"false"`
	FilterLike          json.RawMessage `url:"filter_like,omitempty" required:"false"`
	FilterLt            json.RawMessage `url:"filter_lt,omitempty" required:"false"`
	FilterLteq          json.RawMessage `url:"filter_lteq,omitempty" required:"false"`
	InboxRegistrationId int64           `url:"inbox_registration_id,omitempty" required:"false"`
	InboxId             int64           `url:"inbox_id,omitempty" required:"false"`
	lib.ListParams
}

func (i *InboxUpload) UnmarshalJSON(data []byte) error {
	type inboxUpload InboxUpload
	var v inboxUpload
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*i = InboxUpload(v)
	return nil
}

func (i *InboxUploadCollection) UnmarshalJSON(data []byte) error {
	type inboxUploads []InboxUpload
	var v inboxUploads
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*i = InboxUploadCollection(v)
	return nil
}

func (i *InboxUploadCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*i))
	for i, v := range *i {
		ret[i] = v
	}

	return &ret
}
