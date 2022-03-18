package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type InboxRegistration struct {
	Code             string          `json:"code,omitempty"`
	Name             string          `json:"name,omitempty"`
	Company          string          `json:"company,omitempty"`
	Email            string          `json:"email,omitempty"`
	ClickwrapBody    string          `json:"clickwrap_body,omitempty"`
	FormFieldSetId   int64           `json:"form_field_set_id,omitempty"`
	FormFieldData    json.RawMessage `json:"form_field_data,omitempty"`
	InboxId          int64           `json:"inbox_id,omitempty"`
	InboxRecipientId int64           `json:"inbox_recipient_id,omitempty"`
	InboxTitle       string          `json:"inbox_title,omitempty"`
}

type InboxRegistrationCollection []InboxRegistration

type InboxRegistrationListParams struct {
	Cursor           string `url:"cursor,omitempty" required:"false" json:"cursor,omitempty"`
	PerPage          int64  `url:"per_page,omitempty" required:"false" json:"per_page,omitempty"`
	FolderBehaviorId int64  `url:"folder_behavior_id,omitempty" required:"false" json:"folder_behavior_id,omitempty"`
	lib.ListParams
}

func (i *InboxRegistration) UnmarshalJSON(data []byte) error {
	type inboxRegistration InboxRegistration
	var v inboxRegistration
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*i = InboxRegistration(v)
	return nil
}

func (i *InboxRegistrationCollection) UnmarshalJSON(data []byte) error {
	type inboxRegistrations []InboxRegistration
	var v inboxRegistrations
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*i = InboxRegistrationCollection(v)
	return nil
}

func (i *InboxRegistrationCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*i))
	for i, v := range *i {
		ret[i] = v
	}

	return &ret
}
