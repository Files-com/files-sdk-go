package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type EmailPreferenceBundleNotification struct {
	Id           int64  `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	BundleCode   string `json:"bundle_code,omitempty" path:"bundle_code,omitempty" url:"bundle_code,omitempty"`
	Unsubscribed *bool  `json:"unsubscribed,omitempty" path:"unsubscribed,omitempty" url:"unsubscribed,omitempty"`
}

func (e EmailPreferenceBundleNotification) Identifier() interface{} {
	return e.Id
}

type EmailPreferenceBundleNotificationCollection []EmailPreferenceBundleNotification

func (e *EmailPreferenceBundleNotification) UnmarshalJSON(data []byte) error {
	type emailPreferenceBundleNotification EmailPreferenceBundleNotification
	var v emailPreferenceBundleNotification
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*e = EmailPreferenceBundleNotification(v)
	return nil
}

func (e *EmailPreferenceBundleNotificationCollection) UnmarshalJSON(data []byte) error {
	type emailPreferenceBundleNotifications EmailPreferenceBundleNotificationCollection
	var v emailPreferenceBundleNotifications
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*e = EmailPreferenceBundleNotificationCollection(v)
	return nil
}

func (e *EmailPreferenceBundleNotificationCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*e))
	for i, v := range *e {
		ret[i] = v
	}

	return &ret
}
