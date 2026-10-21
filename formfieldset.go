package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// FormFieldSet is a Files.com API resource.
type FormFieldSet struct {
	// Form field set id.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Title to be displayed.
	Title string `json:"title,omitempty" path:"title,omitempty" url:"title,omitempty"`
	// Reject answers over 1024 characters or more than 32 answered custom fields to ensure complete file metadata capture.
	RequireCompleteMetadataCapture *bool `json:"require_complete_metadata_capture,omitempty" path:"require_complete_metadata_capture,omitempty" url:"require_complete_metadata_capture,omitempty"`
	// Layout of the form.
	FormLayout []int64 `json:"form_layout,omitempty" path:"form_layout,omitempty" url:"form_layout,omitempty"`
	// Associated form field definitions; authenticated form field set responses include historical definitions, while form_layout identifies current fields.
	FormFields []FormField `json:"form_fields,omitempty" path:"form_fields,omitempty" url:"form_fields,omitempty"`
	// Any associated InboxRegistrations or BundleRegistrations can be saved without providing name.
	SkipName *bool `json:"skip_name,omitempty" path:"skip_name,omitempty" url:"skip_name,omitempty"`
	// Any associated InboxRegistrations or BundleRegistrations can be saved without providing email.
	SkipEmail *bool `json:"skip_email,omitempty" path:"skip_email,omitempty" url:"skip_email,omitempty"`
	// Any associated InboxRegistrations or BundleRegistrations can be saved without providing company.
	SkipCompany *bool `json:"skip_company,omitempty" path:"skip_company,omitempty" url:"skip_company,omitempty"`
	// Form Field Set is in use by an active Inbox / Bundle / Inbox Registration / Bundle Registration.
	InUse *bool `json:"in_use,omitempty" path:"in_use,omitempty" url:"in_use,omitempty"`
	// User ID. Provide a value of `0` to operate the current session's user.
	UserId int64 `json:"user_id,omitempty" path:"user_id,omitempty" url:"user_id,omitempty"`
	// Workspace ID.
	WorkspaceId int64 `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
}

// Identifier returns the resource ID.
func (f FormFieldSet) Identifier() interface{} {
	return f.Id
}

// FormFieldSetCollection is a list of FormFieldSet resources.
type FormFieldSetCollection []FormFieldSet

// FormFieldSetListParams contains the request parameters for GET /form_field_sets.
type FormFieldSetListParams struct {
	// User ID. Provide a value of `0` to operate the current session's user.
	UserId int64 `url:"user_id,omitempty" json:"user_id,omitempty" path:"user_id"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// FormFieldSetFindParams contains the request parameters for GET /form_field_sets/{id}.
type FormFieldSetFindParams struct {
	// Form Field Set ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// FormFieldSetCreateParams contains the request parameters for POST /form_field_sets.
type FormFieldSetCreateParams struct {
	// User ID. Provide a value of `0` to operate the current session's user.
	UserId int64 `url:"user_id,omitempty" json:"user_id,omitempty" path:"user_id"`
	// Title to be displayed.
	Title string `url:"title,omitempty" json:"title,omitempty" path:"title"`
	// Reject answers over 1024 characters or more than 32 answered custom fields to ensure complete file metadata capture.
	RequireCompleteMetadataCapture *bool `url:"require_complete_metadata_capture,omitempty" json:"require_complete_metadata_capture,omitempty" path:"require_complete_metadata_capture"`
	// Workspace ID.
	WorkspaceId int64 `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
	// Skip validating form email.
	SkipEmail *bool `url:"skip_email,omitempty" json:"skip_email,omitempty" path:"skip_email"`
	// Skip validating form name.
	SkipName *bool `url:"skip_name,omitempty" json:"skip_name,omitempty" path:"skip_name"`
	// Skip validating company.
	SkipCompany *bool `url:"skip_company,omitempty" json:"skip_company,omitempty" path:"skip_company"`
	// Each element is an object with these keys:
	// id (integer, int32): Id of existing Form Field.
	// label (string): Label of Field.
	// field_type (string): Type of field: text, text_area, dropdown, radio, or date.
	// help_text (string): Help text of field.
	// required (boolean): Is this a required field? (default true).
	// options_for_select (string): List of options for dropdown or radio.
	// default_option (string): Default option to be preselected in the dropdown or radio.
	FormFields []map[string]interface{} `url:"form_fields,omitempty" json:"form_fields,omitempty" path:"form_fields"`
}

// FormFieldSetUpdateParams contains the request parameters for PATCH /form_field_sets/{id}.
type FormFieldSetUpdateParams struct {
	// Form Field Set ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// Title to be displayed.
	Title string `url:"title,omitempty" json:"title,omitempty" path:"title"`
	// Reject answers over 1024 characters or more than 32 answered custom fields to ensure complete file metadata capture.
	RequireCompleteMetadataCapture *bool `url:"require_complete_metadata_capture,omitempty" json:"require_complete_metadata_capture,omitempty" path:"require_complete_metadata_capture"`
	// Workspace ID.
	WorkspaceId int64 `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
	// Skip validating form email.
	SkipEmail *bool `url:"skip_email,omitempty" json:"skip_email,omitempty" path:"skip_email"`
	// Skip validating form name.
	SkipName *bool `url:"skip_name,omitempty" json:"skip_name,omitempty" path:"skip_name"`
	// Skip validating company.
	SkipCompany *bool `url:"skip_company,omitempty" json:"skip_company,omitempty" path:"skip_company"`
	// Each element is an object with these keys:
	// id (integer, int32): Id of existing Form Field.
	// label (string): Label of Field.
	// field_type (string): Type of field: text, text_area, dropdown, radio, or date.
	// help_text (string): Help text of field.
	// required (boolean): Is this a required field? (default true).
	// options_for_select (string): List of options for dropdown or radio.
	// default_option (string): Default option to be preselected in the dropdown or radio.
	FormFields []map[string]interface{} `url:"form_fields,omitempty" json:"form_fields,omitempty" path:"form_fields"`
}

// FormFieldSetDeleteParams contains the request parameters for DELETE /form_field_sets/{id}.
type FormFieldSetDeleteParams struct {
	// Form Field Set ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (f *FormFieldSet) UnmarshalJSON(data []byte) error {
	type formFieldSet FormFieldSet
	var v formFieldSet
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*f = FormFieldSet(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (f *FormFieldSetCollection) UnmarshalJSON(data []byte) error {
	type formFieldSets FormFieldSetCollection
	var v formFieldSets
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*f = FormFieldSetCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (f *FormFieldSetCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*f))
	for i, v := range *f {
		ret[i] = v
	}

	return &ret
}
