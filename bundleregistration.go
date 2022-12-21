package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type BundleRegistration struct {
	Code              string          `json:"code,omitempty" path:"code"`
	Name              string          `json:"name,omitempty" path:"name"`
	Company           string          `json:"company,omitempty" path:"company"`
	Email             string          `json:"email,omitempty" path:"email"`
	Ip                string          `json:"ip,omitempty" path:"ip"`
	InboxCode         string          `json:"inbox_code,omitempty" path:"inbox_code"`
	ClickwrapBody     string          `json:"clickwrap_body,omitempty" path:"clickwrap_body"`
	FormFieldSetId    int64           `json:"form_field_set_id,omitempty" path:"form_field_set_id"`
	FormFieldData     json.RawMessage `json:"form_field_data,omitempty" path:"form_field_data"`
	BundleCode        string          `json:"bundle_code,omitempty" path:"bundle_code"`
	BundleId          int64           `json:"bundle_id,omitempty" path:"bundle_id"`
	BundleRecipientId int64           `json:"bundle_recipient_id,omitempty" path:"bundle_recipient_id"`
	CreatedAt         *time.Time      `json:"created_at,omitempty" path:"created_at"`
}

type BundleRegistrationCollection []BundleRegistration

type BundleRegistrationListParams struct {
	UserId   int64 `url:"user_id,omitempty" required:"false" json:"user_id,omitempty" path:"user_id"`
	BundleId int64 `url:"bundle_id,omitempty" required:"false" json:"bundle_id,omitempty" path:"bundle_id"`
	lib.ListParams
}

func (b *BundleRegistration) UnmarshalJSON(data []byte) error {
	type bundleRegistration BundleRegistration
	var v bundleRegistration
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*b = BundleRegistration(v)
	return nil
}

func (b *BundleRegistrationCollection) UnmarshalJSON(data []byte) error {
	type bundleRegistrations BundleRegistrationCollection
	var v bundleRegistrations
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
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
