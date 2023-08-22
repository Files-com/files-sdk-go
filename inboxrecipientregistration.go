package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type InboxRecipientRegistration struct {
	Code                  string `json:"code,omitempty" path:"code,omitempty" url:"code,omitempty"`
	InboxRegistrationCode string `json:"inbox_registration_code,omitempty" path:"inbox_registration_code,omitempty" url:"inbox_registration_code,omitempty"`
	Recipient             string `json:"recipient,omitempty" path:"recipient,omitempty" url:"recipient,omitempty"`
	Name                  string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	Company               string `json:"company,omitempty" path:"company,omitempty" url:"company,omitempty"`
	InboxRecipientCode    string `json:"inbox_recipient_code,omitempty" path:"inbox_recipient_code,omitempty" url:"inbox_recipient_code,omitempty"`
}

// Identifier no path or id

type InboxRecipientRegistrationCollection []InboxRecipientRegistration

type InboxRecipientRegistrationCreateParams struct {
	InboxRecipientCode string `url:"inbox_recipient_code,omitempty" required:"true" json:"inbox_recipient_code,omitempty" path:"inbox_recipient_code"`
}

func (i *InboxRecipientRegistration) UnmarshalJSON(data []byte) error {
	type inboxRecipientRegistration InboxRecipientRegistration
	var v inboxRecipientRegistration
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*i = InboxRecipientRegistration(v)
	return nil
}

func (i *InboxRecipientRegistrationCollection) UnmarshalJSON(data []byte) error {
	type inboxRecipientRegistrations InboxRecipientRegistrationCollection
	var v inboxRecipientRegistrations
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*i = InboxRecipientRegistrationCollection(v)
	return nil
}

func (i *InboxRecipientRegistrationCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*i))
	for i, v := range *i {
		ret[i] = v
	}

	return &ret
}
