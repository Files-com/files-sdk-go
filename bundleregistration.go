package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type BundleRegistration struct {
	Code              string          `json:"code,omitempty"`
	Name              string          `json:"name,omitempty"`
	Company           string          `json:"company,omitempty"`
	Email             string          `json:"email,omitempty"`
	Ip                string          `json:"ip,omitempty"`
	InboxCode         string          `json:"inbox_code,omitempty"`
	ClickwrapBody     string          `json:"clickwrap_body,omitempty"`
	FormFieldSetId    int64           `json:"form_field_set_id,omitempty"`
	FormFieldData     json.RawMessage `json:"form_field_data,omitempty"`
	BundleCode        string          `json:"bundle_code,omitempty"`
	BundleId          int64           `json:"bundle_id,omitempty"`
	BundleRecipientId int64           `json:"bundle_recipient_id,omitempty"`
}

type BundleRegistrationCollection []BundleRegistration

type BundleRegistrationListParams struct {
	UserId   int64  `url:"user_id,omitempty" required:"false" json:"user_id,omitempty"`
	Cursor   string `url:"cursor,omitempty" required:"false" json:"cursor,omitempty"`
	PerPage  int64  `url:"per_page,omitempty" required:"false" json:"per_page,omitempty"`
	BundleId int64  `url:"bundle_id,omitempty" required:"false" json:"bundle_id,omitempty"`
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
