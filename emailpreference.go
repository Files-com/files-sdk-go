package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type EmailPreference struct {
	Email               string   `json:"email,omitempty" path:"email,omitempty" url:"email,omitempty"`
	Notifications       []string `json:"notifications,omitempty" path:"notifications,omitempty" url:"notifications,omitempty"`
	BundleNotifications []string `json:"bundle_notifications,omitempty" path:"bundle_notifications,omitempty" url:"bundle_notifications,omitempty"`
	ReceiveAdminAlerts  *bool    `json:"receive_admin_alerts,omitempty" path:"receive_admin_alerts,omitempty" url:"receive_admin_alerts,omitempty"`
}

// Identifier no path or id

type EmailPreferenceCollection []EmailPreference

type UserParam struct {
	ReceiveAdminAlerts  string `url:"receive_admin_alerts,omitempty" json:"receive_admin_alerts,omitempty" path:"receive_admin_alerts"`
	Unsubscribed        string `url:"unsubscribed,omitempty" json:"unsubscribed,omitempty" path:"unsubscribed"`
	Notifications       string `url:"notifications,omitempty" json:"notifications,omitempty" path:"notifications"`
	BundleNotifications string `url:"bundle_notifications,omitempty" json:"bundle_notifications,omitempty" path:"bundle_notifications"`
	Unsubscribe         string `url:"unsubscribe,omitempty" json:"unsubscribe,omitempty" path:"unsubscribe"`
}

type EmailPreferenceGetParams struct {
	Token string `url:"-,omitempty" required:"false" json:"-,omitempty" path:"token"`
}

type EmailPreferenceUpdateParams struct {
	Token     string    `url:"-,omitempty" required:"false" json:"-,omitempty" path:"token"`
	UserParam UserParam `url:"user,omitempty" required:"false" json:"user,omitempty" path:"user"`
}

func (e *EmailPreference) UnmarshalJSON(data []byte) error {
	type emailPreference EmailPreference
	var v emailPreference
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*e = EmailPreference(v)
	return nil
}

func (e *EmailPreferenceCollection) UnmarshalJSON(data []byte) error {
	type emailPreferences EmailPreferenceCollection
	var v emailPreferences
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*e = EmailPreferenceCollection(v)
	return nil
}

func (e *EmailPreferenceCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*e))
	for i, v := range *e {
		ret[i] = v
	}

	return &ret
}
