package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/lib"
)

type BundleRegistration struct {
	Code           string `json:"code,omitempty"`
	Name           string `json:"name,omitempty"`
	Company        string `json:"company,omitempty"`
	Email          string `json:"email,omitempty"`
	InboxCode      string `json:"inbox_code,omitempty"`
	ClickwrapBody  string `json:"clickwrap_body,omitempty"`
	FormFieldSetId int64  `json:"form_field_set_id,omitempty"`
	FormFieldData  string `json:"form_field_data,omitempty"`
}

type BundleRegistrationCollection []BundleRegistration

type BundleRegistrationListParams struct {
	UserId   int64  `url:"user_id,omitempty" required:"false"`
	Cursor   string `url:"cursor,omitempty" required:"false"`
	PerPage  int64  `url:"per_page,omitempty" required:"false"`
	BundleId int64  `url:"bundle_id,omitempty" required:"true"`
	lib.ListParams
}

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

func (b *BundleRegistrationCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*b))
	for i, v := range *b {
		ret[i] = v
	}

	return &ret
}
