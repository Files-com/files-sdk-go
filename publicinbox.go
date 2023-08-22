package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type PublicInbox struct {
	ColorLeft                 string       `json:"color_left,omitempty" path:"color_left,omitempty" url:"color_left,omitempty"`
	ColorLink                 string       `json:"color_link,omitempty" path:"color_link,omitempty" url:"color_link,omitempty"`
	ColorText                 string       `json:"color_text,omitempty" path:"color_text,omitempty" url:"color_text,omitempty"`
	ColorTop                  string       `json:"color_top,omitempty" path:"color_top,omitempty" url:"color_top,omitempty"`
	ColorTopText              string       `json:"color_top_text,omitempty" path:"color_top_text,omitempty" url:"color_top_text,omitempty"`
	Title                     string       `json:"title,omitempty" path:"title,omitempty" url:"title,omitempty"`
	Description               string       `json:"description,omitempty" path:"description,omitempty" url:"description,omitempty"`
	HelpText                  string       `json:"help_text,omitempty" path:"help_text,omitempty" url:"help_text,omitempty"`
	Key                       string       `json:"key,omitempty" path:"key,omitempty" url:"key,omitempty"`
	ShowOnLoginPage           *bool        `json:"show_on_login_page,omitempty" path:"show_on_login_page,omitempty" url:"show_on_login_page,omitempty"`
	HasPassword               *bool        `json:"has_password,omitempty" path:"has_password,omitempty" url:"has_password,omitempty"`
	RequireRegistration       *bool        `json:"require_registration,omitempty" path:"require_registration,omitempty" url:"require_registration,omitempty"`
	DontAllowFoldersInUploads *bool        `json:"dont_allow_folders_in_uploads,omitempty" path:"dont_allow_folders_in_uploads,omitempty" url:"dont_allow_folders_in_uploads,omitempty"`
	ClickwrapBody             string       `json:"clickwrap_body,omitempty" path:"clickwrap_body,omitempty" url:"clickwrap_body,omitempty"`
	FormFieldSet              FormFieldSet `json:"form_field_set,omitempty" path:"form_field_set,omitempty" url:"form_field_set,omitempty"`
	RequireLogout             *bool        `json:"require_logout,omitempty" path:"require_logout,omitempty" url:"require_logout,omitempty"`
}

// Identifier no path or id

type PublicInboxCollection []PublicInbox

type PublicInboxListParams struct {
	Action string `url:"action,omitempty" required:"false" json:"action,omitempty" path:"action"`
	ListParams
}

type PublicInboxGetKeyParams struct {
	Key           string `url:"-,omitempty" required:"false" json:"-,omitempty" path:"key"`
	RecipientCode string `url:"recipient_code,omitempty" required:"false" json:"recipient_code,omitempty" path:"recipient_code"`
}

func (p *PublicInbox) UnmarshalJSON(data []byte) error {
	type publicInbox PublicInbox
	var v publicInbox
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*p = PublicInbox(v)
	return nil
}

func (p *PublicInboxCollection) UnmarshalJSON(data []byte) error {
	type publicInboxs PublicInboxCollection
	var v publicInboxs
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*p = PublicInboxCollection(v)
	return nil
}

func (p *PublicInboxCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*p))
	for i, v := range *p {
		ret[i] = v
	}

	return &ret
}
