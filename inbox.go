package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Inbox struct {
	Title                                    string       `json:"title,omitempty" path:"title,omitempty" url:"title,omitempty"`
	Description                              string       `json:"description,omitempty" path:"description,omitempty" url:"description,omitempty"`
	HelpText                                 string       `json:"help_text,omitempty" path:"help_text,omitempty" url:"help_text,omitempty"`
	Key                                      string       `json:"key,omitempty" path:"key,omitempty" url:"key,omitempty"`
	ShowOnLoginPage                          *bool        `json:"show_on_login_page,omitempty" path:"show_on_login_page,omitempty" url:"show_on_login_page,omitempty"`
	HasPassword                              *bool        `json:"has_password,omitempty" path:"has_password,omitempty" url:"has_password,omitempty"`
	RequireRegistration                      *bool        `json:"require_registration,omitempty" path:"require_registration,omitempty" url:"require_registration,omitempty"`
	DontAllowFoldersInUploads                *bool        `json:"dont_allow_folders_in_uploads,omitempty" path:"dont_allow_folders_in_uploads,omitempty" url:"dont_allow_folders_in_uploads,omitempty"`
	ClickwrapBody                            string       `json:"clickwrap_body,omitempty" path:"clickwrap_body,omitempty" url:"clickwrap_body,omitempty"`
	FormFieldSet                             FormFieldSet `json:"form_field_set,omitempty" path:"form_field_set,omitempty" url:"form_field_set,omitempty"`
	NotifySendersOnSuccessfulUploadsViaEmail *bool        `json:"notify_senders_on_successful_uploads_via_email,omitempty" path:"notify_senders_on_successful_uploads_via_email,omitempty" url:"notify_senders_on_successful_uploads_via_email,omitempty"`
	AllowWhitelisting                        *bool        `json:"allow_whitelisting,omitempty" path:"allow_whitelisting,omitempty" url:"allow_whitelisting,omitempty"`
	Whitelist                                []string     `json:"whitelist,omitempty" path:"whitelist,omitempty" url:"whitelist,omitempty"`
	InboundEmailAddress                      string       `json:"inbound_email_address,omitempty" path:"inbound_email_address,omitempty" url:"inbound_email_address,omitempty"`
}

// Identifier no path or id

type InboxCollection []Inbox

type InboxListParams struct {
	Action string                 `url:"action,omitempty" required:"false" json:"action,omitempty" path:"action"`
	SortBy map[string]interface{} `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty" path:"sort_by"`
	Filter Inbox                  `url:"filter,omitempty" required:"false" json:"filter,omitempty" path:"filter"`
	ListParams
}

func (i *Inbox) UnmarshalJSON(data []byte) error {
	type inbox Inbox
	var v inbox
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*i = Inbox(v)
	return nil
}

func (i *InboxCollection) UnmarshalJSON(data []byte) error {
	type inboxs InboxCollection
	var v inboxs
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*i = InboxCollection(v)
	return nil
}

func (i *InboxCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*i))
	for i, v := range *i {
		ret[i] = v
	}

	return &ret
}
