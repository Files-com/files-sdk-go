package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// BundleNotification is a Files.com API resource.
type BundleNotification struct {
	BundleId             int64 `json:"bundle_id,omitempty" path:"bundle_id,omitempty" url:"bundle_id,omitempty"`
	Id                   int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	NotifyOnRegistration *bool `json:"notify_on_registration,omitempty" path:"notify_on_registration,omitempty" url:"notify_on_registration,omitempty"`
	NotifyOnUpload       *bool `json:"notify_on_upload,omitempty" path:"notify_on_upload,omitempty" url:"notify_on_upload,omitempty"`
	NotifyCurrentUser    *bool `json:"notify_current_user,omitempty" path:"notify_current_user,omitempty" url:"notify_current_user,omitempty"`
	NotifyUserId         int64 `json:"notify_user_id,omitempty" path:"notify_user_id,omitempty" url:"notify_user_id,omitempty"`
	WorkspaceId          int64 `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	UserId               int64 `json:"user_id,omitempty" path:"user_id,omitempty" url:"user_id,omitempty"`
}

// Identifier returns the resource ID.
func (b BundleNotification) Identifier() interface{} {
	return b.Id
}

// BundleNotificationCollection is a list of BundleNotification resources.
type BundleNotificationCollection []BundleNotification

// BundleNotificationListParams contains the request parameters for GET /bundle_notifications.
type BundleNotificationListParams struct {
	UserId   int64       `url:"user_id,omitempty" json:"user_id,omitempty" path:"user_id"`
	SortBy   interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	Filter   interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	BundleId int64       `url:"bundle_id,omitempty" json:"bundle_id,omitempty" path:"bundle_id"`
	ListParams
}

// BundleNotificationFindParams contains the request parameters for GET /bundle_notifications/{id}.
type BundleNotificationFindParams struct {
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// BundleNotificationCreateParams contains the request parameters for POST /bundle_notifications.
type BundleNotificationCreateParams struct {
	UserId               int64 `url:"user_id,omitempty" json:"user_id,omitempty" path:"user_id"`
	BundleId             int64 `url:"bundle_id" json:"bundle_id" path:"bundle_id"`
	NotifyUserId         int64 `url:"notify_user_id,omitempty" json:"notify_user_id,omitempty" path:"notify_user_id"`
	NotifyOnRegistration *bool `url:"notify_on_registration,omitempty" json:"notify_on_registration,omitempty" path:"notify_on_registration"`
	NotifyOnUpload       *bool `url:"notify_on_upload,omitempty" json:"notify_on_upload,omitempty" path:"notify_on_upload"`
}

// BundleNotificationUpdateParams contains the request parameters for PATCH /bundle_notifications/{id}.
type BundleNotificationUpdateParams struct {
	Id                   int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	NotifyOnRegistration *bool `url:"notify_on_registration,omitempty" json:"notify_on_registration,omitempty" path:"notify_on_registration"`
	NotifyOnUpload       *bool `url:"notify_on_upload,omitempty" json:"notify_on_upload,omitempty" path:"notify_on_upload"`
}

// BundleNotificationDeleteParams contains the request parameters for DELETE /bundle_notifications/{id}.
type BundleNotificationDeleteParams struct {
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (b *BundleNotification) UnmarshalJSON(data []byte) error {
	type bundleNotification BundleNotification
	var v bundleNotification
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*b = BundleNotification(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (b *BundleNotificationCollection) UnmarshalJSON(data []byte) error {
	type bundleNotifications BundleNotificationCollection
	var v bundleNotifications
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*b = BundleNotificationCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (b *BundleNotificationCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*b))
	for i, v := range *b {
		ret[i] = v
	}

	return &ret
}
