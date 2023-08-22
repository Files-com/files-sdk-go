package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type BundleRecipientRegistration struct {
	Code                   string `json:"code,omitempty" path:"code,omitempty" url:"code,omitempty"`
	BundleRegistrationCode string `json:"bundle_registration_code,omitempty" path:"bundle_registration_code,omitempty" url:"bundle_registration_code,omitempty"`
	Recipient              string `json:"recipient,omitempty" path:"recipient,omitempty" url:"recipient,omitempty"`
	Name                   string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	Company                string `json:"company,omitempty" path:"company,omitempty" url:"company,omitempty"`
	InboxCode              string `json:"inbox_code,omitempty" path:"inbox_code,omitempty" url:"inbox_code,omitempty"`
	BundleRecipientCode    string `json:"bundle_recipient_code,omitempty" path:"bundle_recipient_code,omitempty" url:"bundle_recipient_code,omitempty"`
}

// Identifier no path or id

type BundleRecipientRegistrationCollection []BundleRecipientRegistration

type BundleRecipientRegistrationCreateParams struct {
	BundleRecipientCode string `url:"bundle_recipient_code,omitempty" required:"true" json:"bundle_recipient_code,omitempty" path:"bundle_recipient_code"`
}

func (b *BundleRecipientRegistration) UnmarshalJSON(data []byte) error {
	type bundleRecipientRegistration BundleRecipientRegistration
	var v bundleRecipientRegistration
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*b = BundleRecipientRegistration(v)
	return nil
}

func (b *BundleRecipientRegistrationCollection) UnmarshalJSON(data []byte) error {
	type bundleRecipientRegistrations BundleRecipientRegistrationCollection
	var v bundleRecipientRegistrations
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*b = BundleRecipientRegistrationCollection(v)
	return nil
}

func (b *BundleRecipientRegistrationCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*b))
	for i, v := range *b {
		ret[i] = v
	}

	return &ret
}
