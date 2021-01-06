package files_sdk

import (
	"encoding/json"
)

type BundleRegistration struct {
	Code           string `json:"code,omitempty"`
	Name           string `json:"name,omitempty"`
	Company        string `json:"company,omitempty"`
	Email          string `json:"email,omitempty"`
	InboxCode      string `json:"inbox_code,omitempty"`
	FormFieldSetId int64  `json:"form_field_set_id,omitempty"`
	FormFieldData  string `json:"form_field_data,omitempty"`
}

type BundleRegistrationCollection []BundleRegistration

func (b *BundleRegistration) UnmarshalJSON(data []byte) error {
	type bundleRegistration BundleRegistration
	var v bundleRegistration
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*b = BundleRegistration(v)
	return nil
}

func (b *BundleRegistrationCollection) UnmarshalJSON(data []byte) error {
	type bundleRegistrations []BundleRegistration
	var v bundleRegistrations
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*b = BundleRegistrationCollection(v)
	return nil
}
