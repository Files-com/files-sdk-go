package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// FormField is a Files.com API resource.
type FormField struct {
	// Form field id.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Label to be displayed.
	Label string `json:"label,omitempty" path:"label,omitempty" url:"label,omitempty"`
	// Is this a required field?
	Required *bool `json:"required,omitempty" path:"required,omitempty" url:"required,omitempty"`
	// Help text to be displayed.
	HelpText string `json:"help_text,omitempty" path:"help_text,omitempty" url:"help_text,omitempty"`
	// Type of Field.
	FieldType string `json:"field_type,omitempty" path:"field_type,omitempty" url:"field_type,omitempty"`
	// Options to display for radio and dropdown.
	OptionsForSelect []string `json:"options_for_select,omitempty" path:"options_for_select,omitempty" url:"options_for_select,omitempty"`
	// Default option for radio and dropdown.
	DefaultOption string `json:"default_option,omitempty" path:"default_option,omitempty" url:"default_option,omitempty"`
	// Form field set id.
	FormFieldSetId int64 `json:"form_field_set_id,omitempty" path:"form_field_set_id,omitempty" url:"form_field_set_id,omitempty"`
}

// Identifier returns the resource ID.
func (f FormField) Identifier() interface{} {
	return f.Id
}

// FormFieldCollection is a list of FormField resources.
type FormFieldCollection []FormField

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (f *FormField) UnmarshalJSON(data []byte) error {
	type formField FormField
	var v formField
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*f = FormField(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (f *FormFieldCollection) UnmarshalJSON(data []byte) error {
	type formFields FormFieldCollection
	var v formFields
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*f = FormFieldCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (f *FormFieldCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*f))
	for i, v := range *f {
		ret[i] = v
	}

	return &ret
}
