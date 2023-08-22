package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type EmailPreferenceNotification struct {
	Id           int64  `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	Path         string `json:"path,omitempty" path:"path,omitempty" url:"path,omitempty"`
	SendInterval string `json:"send_interval,omitempty" path:"send_interval,omitempty" url:"send_interval,omitempty"`
	Unsubscribed *bool  `json:"unsubscribed,omitempty" path:"unsubscribed,omitempty" url:"unsubscribed,omitempty"`
}

func (e EmailPreferenceNotification) Identifier() interface{} {
	return e.Id
}

type EmailPreferenceNotificationCollection []EmailPreferenceNotification

func (e *EmailPreferenceNotification) UnmarshalJSON(data []byte) error {
	type emailPreferenceNotification EmailPreferenceNotification
	var v emailPreferenceNotification
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*e = EmailPreferenceNotification(v)
	return nil
}

func (e *EmailPreferenceNotificationCollection) UnmarshalJSON(data []byte) error {
	type emailPreferenceNotifications EmailPreferenceNotificationCollection
	var v emailPreferenceNotifications
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*e = EmailPreferenceNotificationCollection(v)
	return nil
}

func (e *EmailPreferenceNotificationCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*e))
	for i, v := range *e {
		ret[i] = v
	}

	return &ret
}
