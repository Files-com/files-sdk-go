package files_sdk

import (
	"encoding/json"
)

type FormField struct {
	Id               int64  `json:"id,omitempty"`
	Label            string `json:"label,omitempty"`
	Required         *bool  `json:"required,omitempty"`
	HelpText         string `json:"help_text,omitempty"`
	FieldType        string `json:"field_type,omitempty"`
	OptionsForSelect string `json:"options_for_select,omitempty"`
	DefaultOption    string `json:"default_option,omitempty"`
	FormFieldSetId   int64  `json:"form_field_set_id,omitempty"`
}

type FormFieldCollection []FormField

func (f *FormField) UnmarshalJSON(data []byte) error {
	type formField FormField
	var v formField
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*f = FormField(v)
	return nil
}

func (f *FormFieldCollection) UnmarshalJSON(data []byte) error {
	type formFields []FormField
	var v formFields
	if err := json.Unmarshal(data, &v); err != nil {
		return err
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
