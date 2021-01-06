package files_sdk

import (
	"encoding/json"
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
