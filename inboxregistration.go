package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type InboxRegistration struct {
	Code                           string                 `json:"code,omitempty" path:"code,omitempty" url:"code,omitempty"`
	Name                           string                 `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	Company                        string                 `json:"company,omitempty" path:"company,omitempty" url:"company,omitempty"`
	Email                          string                 `json:"email,omitempty" path:"email,omitempty" url:"email,omitempty"`
	ClickwrapBody                  string                 `json:"clickwrap_body,omitempty" path:"clickwrap_body,omitempty" url:"clickwrap_body,omitempty"`
	FormFieldSetId                 int64                  `json:"form_field_set_id,omitempty" path:"form_field_set_id,omitempty" url:"form_field_set_id,omitempty"`
	FormFieldData                  map[string]interface{} `json:"form_field_data,omitempty" path:"form_field_data,omitempty" url:"form_field_data,omitempty"`
	InboxId                        int64                  `json:"inbox_id,omitempty" path:"inbox_id,omitempty" url:"inbox_id,omitempty"`
	InboxRecipientId               int64                  `json:"inbox_recipient_id,omitempty" path:"inbox_recipient_id,omitempty" url:"inbox_recipient_id,omitempty"`
	InboxTitle                     string                 `json:"inbox_title,omitempty" path:"inbox_title,omitempty" url:"inbox_title,omitempty"`
	CreatedAt                      *time.Time             `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	InboxCode                      string                 `json:"inbox_code,omitempty" path:"inbox_code,omitempty" url:"inbox_code,omitempty"`
	InboxRecipientRegistrationCode string                 `json:"inbox_recipient_registration_code,omitempty" path:"inbox_recipient_registration_code,omitempty" url:"inbox_recipient_registration_code,omitempty"`
	Password                       string                 `json:"password,omitempty" path:"password,omitempty" url:"password,omitempty"`
}

// Identifier no path or id

type InboxRegistrationCollection []InboxRegistration

type InboxRegistrationListParams struct {
	Action           string `url:"action,omitempty" required:"false" json:"action,omitempty" path:"action"`
	FolderBehaviorId int64  `url:"folder_behavior_id,omitempty" required:"false" json:"folder_behavior_id,omitempty" path:"folder_behavior_id"`
	ListParams
}

type InboxRegistrationCreateParams struct {
	InboxCode                      string `url:"inbox_code,omitempty" required:"true" json:"inbox_code,omitempty" path:"inbox_code"`
	InboxRecipientRegistrationCode string `url:"inbox_recipient_registration_code,omitempty" required:"false" json:"inbox_recipient_registration_code,omitempty" path:"inbox_recipient_registration_code"`
	Password                       string `url:"password,omitempty" required:"false" json:"password,omitempty" path:"password"`
	Name                           string `url:"name,omitempty" required:"false" json:"name,omitempty" path:"name"`
	Company                        string `url:"company,omitempty" required:"false" json:"company,omitempty" path:"company"`
	Email                          string `url:"email,omitempty" required:"false" json:"email,omitempty" path:"email"`
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
