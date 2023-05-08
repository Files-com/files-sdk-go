package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type FormField struct {
	Id               int64    `json:"id,omitempty" path:"id"`
	Label            string   `json:"label,omitempty" path:"label"`
	Required         *bool    `json:"required,omitempty" path:"required"`
	HelpText         string   `json:"help_text,omitempty" path:"help_text"`
	FieldType        string   `json:"field_type,omitempty" path:"field_type"`
	OptionsForSelect []string `json:"options_for_select,omitempty" path:"options_for_select"`
	DefaultOption    string   `json:"default_option,omitempty" path:"default_option"`
	FormFieldSetId   int64    `json:"form_field_set_id,omitempty" path:"form_field_set_id"`
}

func (f FormField) Identifier() interface{} {
	return f.Id
}

type FormFieldCollection []FormField

func (f *FormField) UnmarshalJSON(data []byte) error {
	type formField FormField
	var v formField
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*f = FormField(v)
	return nil
}

func (f *FormFieldCollection) UnmarshalJSON(data []byte) error {
	type formFields FormFieldCollection
	var v formFields
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*f = FormFieldCollection(v)
	return nil
}

func (f *FormFieldCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*f))
	for i, v := range *f {
		ret[i] = v
	}

	return &ret
}
