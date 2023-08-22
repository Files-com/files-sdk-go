package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type BundleNotification struct {
	BundleId             int64 `json:"bundle_id,omitempty" path:"bundle_id,omitempty" url:"bundle_id,omitempty"`
	Id                   int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	NotifyOnRegistration *bool `json:"notify_on_registration,omitempty" path:"notify_on_registration,omitempty" url:"notify_on_registration,omitempty"`
	NotifyOnUpload       *bool `json:"notify_on_upload,omitempty" path:"notify_on_upload,omitempty" url:"notify_on_upload,omitempty"`
	UserId               int64 `json:"user_id,omitempty" path:"user_id,omitempty" url:"user_id,omitempty"`
}

func (b BundleNotification) Identifier() interface{} {
	return b.Id
}

type BundleNotificationCollection []BundleNotification

type BundleNotificationListParams struct {
	UserId int64                  `url:"user_id,omitempty" required:"false" json:"user_id,omitempty" path:"user_id"`
	SortBy map[string]interface{} `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty" path:"sort_by"`
	Filter BundleNotification     `url:"filter,omitempty" required:"false" json:"filter,omitempty" path:"filter"`
	ListParams
}

type BundleNotificationFindParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
}

type BundleNotificationCreateParams struct {
	UserId               int64 `url:"user_id,omitempty" required:"false" json:"user_id,omitempty" path:"user_id"`
	NotifyOnRegistration *bool `url:"notify_on_registration,omitempty" required:"false" json:"notify_on_registration,omitempty" path:"notify_on_registration"`
	NotifyOnUpload       *bool `url:"notify_on_upload,omitempty" required:"false" json:"notify_on_upload,omitempty" path:"notify_on_upload"`
	BundleId             int64 `url:"bundle_id,omitempty" required:"true" json:"bundle_id,omitempty" path:"bundle_id"`
}

type BundleNotificationUpdateParams struct {
	Id                   int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
	NotifyOnRegistration *bool `url:"notify_on_registration,omitempty" required:"false" json:"notify_on_registration,omitempty" path:"notify_on_registration"`
	NotifyOnUpload       *bool `url:"notify_on_upload,omitempty" required:"false" json:"notify_on_upload,omitempty" path:"notify_on_upload"`
}

type BundleNotificationDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
}

func (b *BundleNotification) UnmarshalJSON(data []byte) error {
	type bundleNotification BundleNotification
	var v bundleNotification
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*b = BundleNotification(v)
	return nil
}

func (b *BundleNotificationCollection) UnmarshalJSON(data []byte) error {
	type bundleNotifications BundleNotificationCollection
	var v bundleNotifications
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*b = BundleNotificationCollection(v)
	return nil
}

func (b *BundleNotificationCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*b))
	for i, v := range *b {
		ret[i] = v
	}

	return &ret
}
