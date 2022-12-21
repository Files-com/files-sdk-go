package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type InboxRegistration struct {
	Code             string          `json:"code,omitempty" path:"code"`
	Name             string          `json:"name,omitempty" path:"name"`
	Company          string          `json:"company,omitempty" path:"company"`
	Email            string          `json:"email,omitempty" path:"email"`
	ClickwrapBody    string          `json:"clickwrap_body,omitempty" path:"clickwrap_body"`
	FormFieldSetId   int64           `json:"form_field_set_id,omitempty" path:"form_field_set_id"`
	FormFieldData    json.RawMessage `json:"form_field_data,omitempty" path:"form_field_data"`
	InboxId          int64           `json:"inbox_id,omitempty" path:"inbox_id"`
	InboxRecipientId int64           `json:"inbox_recipient_id,omitempty" path:"inbox_recipient_id"`
	InboxTitle       string          `json:"inbox_title,omitempty" path:"inbox_title"`
	CreatedAt        *time.Time      `json:"created_at,omitempty" path:"created_at"`
}

type InboxRegistrationCollection []InboxRegistration

type InboxRegistrationListParams struct {
	FolderBehaviorId int64 `url:"folder_behavior_id,omitempty" required:"false" json:"folder_behavior_id,omitempty" path:"folder_behavior_id"`
	lib.ListParams
}

func (i *InboxRegistration) UnmarshalJSON(data []byte) error {
	type inboxRegistration InboxRegistration
	var v inboxRegistration
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*i = InboxRegistration(v)
	return nil
}

func (i *InboxRegistrationCollection) UnmarshalJSON(data []byte) error {
	type inboxRegistrations InboxRegistrationCollection
	var v inboxRegistrations
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
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
