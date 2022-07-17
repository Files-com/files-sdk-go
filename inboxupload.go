package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type InboxUpload struct {
	InboxRegistration InboxRegistration `json:"inbox_registration,omitempty" path:"inbox_registration"`
	Path              string            `json:"path,omitempty" path:"path"`
	CreatedAt         *time.Time        `json:"created_at,omitempty" path:"created_at"`
}

type InboxUploadCollection []InboxUpload

type InboxUploadListParams struct {
	SortBy              json.RawMessage `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty" path:"sort_by"`
	Filter              json.RawMessage `url:"filter,omitempty" required:"false" json:"filter,omitempty" path:"filter"`
	FilterGt            json.RawMessage `url:"filter_gt,omitempty" required:"false" json:"filter_gt,omitempty" path:"filter_gt"`
	FilterGteq          json.RawMessage `url:"filter_gteq,omitempty" required:"false" json:"filter_gteq,omitempty" path:"filter_gteq"`
	FilterLike          json.RawMessage `url:"filter_like,omitempty" required:"false" json:"filter_like,omitempty" path:"filter_like"`
	FilterLt            json.RawMessage `url:"filter_lt,omitempty" required:"false" json:"filter_lt,omitempty" path:"filter_lt"`
	FilterLteq          json.RawMessage `url:"filter_lteq,omitempty" required:"false" json:"filter_lteq,omitempty" path:"filter_lteq"`
	InboxRegistrationId int64           `url:"inbox_registration_id,omitempty" required:"false" json:"inbox_registration_id,omitempty" path:"inbox_registration_id"`
	InboxId             int64           `url:"inbox_id,omitempty" required:"false" json:"inbox_id,omitempty" path:"inbox_id"`
	lib.ListParams
}

func (i *InboxUpload) UnmarshalJSON(data []byte) error {
	type inboxUpload InboxUpload
	var v inboxUpload
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*i = InboxUpload(v)
	return nil
}

func (i *InboxUploadCollection) UnmarshalJSON(data []byte) error {
	type inboxUploads InboxUploadCollection
	var v inboxUploads
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
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
