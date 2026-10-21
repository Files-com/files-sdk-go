package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// InboxRegistration is a Files.com API resource.
type InboxRegistration struct {
	// Registration cookie code.
	Code string `json:"code,omitempty" path:"code,omitempty" url:"code,omitempty"`
	// Registrant name.
	Name string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	// Registrant company name.
	Company string `json:"company,omitempty" path:"company,omitempty" url:"company,omitempty"`
	// Registrant email address.
	Email string `json:"email,omitempty" path:"email,omitempty" url:"email,omitempty"`
	// Registrant IP Address.
	Ip string `json:"ip,omitempty" path:"ip,omitempty" url:"ip,omitempty"`
	// Clickwrap text that was shown to the registrant.
	ClickwrapBody string `json:"clickwrap_body,omitempty" path:"clickwrap_body,omitempty" url:"clickwrap_body,omitempty"`
	// Id of associated form field set.
	FormFieldSetId int64 `json:"form_field_set_id,omitempty" path:"form_field_set_id,omitempty" url:"form_field_set_id,omitempty"`
	// Data for form field set with form field ids as keys and user data as values.
	FormFieldData interface{} `json:"form_field_data,omitempty" path:"form_field_data,omitempty" url:"form_field_data,omitempty"`
	// Id of associated inbox.
	InboxId int64 `json:"inbox_id,omitempty" path:"inbox_id,omitempty" url:"inbox_id,omitempty"`
	// Id of associated inbox recipient.
	InboxRecipientId int64 `json:"inbox_recipient_id,omitempty" path:"inbox_recipient_id,omitempty" url:"inbox_recipient_id,omitempty"`
	// Title of associated inbox.
	InboxTitle string `json:"inbox_title,omitempty" path:"inbox_title,omitempty" url:"inbox_title,omitempty"`
	// Registration creation date/time.
	CreatedAt *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
}

// Identifier no path or id

// InboxRegistrationCollection is a list of InboxRegistration resources.
type InboxRegistrationCollection []InboxRegistration

// InboxRegistrationListParams contains the request parameters for GET /inbox_registrations.
type InboxRegistrationListParams struct {
	// ID of the associated Inbox. This is required if the user is not a site admin.
	FolderBehaviorId int64 `url:"folder_behavior_id,omitempty" json:"folder_behavior_id,omitempty" path:"folder_behavior_id"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (i *InboxRegistration) UnmarshalJSON(data []byte) error {
	type inboxRegistration InboxRegistration
	var v inboxRegistration
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*i = InboxRegistration(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (i *InboxRegistrationCollection) UnmarshalJSON(data []byte) error {
	type inboxRegistrations InboxRegistrationCollection
	var v inboxRegistrations
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*i = InboxRegistrationCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (i *InboxRegistrationCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*i))
	for i, v := range *i {
		ret[i] = v
	}

	return &ret
}
