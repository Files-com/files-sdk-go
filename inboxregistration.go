package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/lib"
)

type InboxRegistration struct {
	Code           string `json:"code,omitempty"`
	Name           string `json:"name,omitempty"`
	Company        string `json:"company,omitempty"`
	Email          string `json:"email,omitempty"`
	FormFieldSetId int64  `json:"form_field_set_id,omitempty"`
	FormFieldData  string `json:"form_field_data,omitempty"`
}

type InboxRegistrationCollection []InboxRegistration

type InboxRegistrationListParams struct {
	Cursor           string `url:"cursor,omitempty" required:"false"`
	PerPage          int64  `url:"per_page,omitempty" required:"false"`
	FolderBehaviorId int64  `url:"folder_behavior_id,omitempty" required:"true"`
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
