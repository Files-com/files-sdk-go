package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Notification struct {
	Id                       int64    `json:"id,omitempty" path:"id"`
	Path                     string   `json:"path,omitempty" path:"path"`
	GroupId                  int64    `json:"group_id,omitempty" path:"group_id"`
	GroupName                string   `json:"group_name,omitempty" path:"group_name"`
	TriggeringGroupIds       []int64  `json:"triggering_group_ids,omitempty" path:"triggering_group_ids"`
	TriggeringUserIds        []int64  `json:"triggering_user_ids,omitempty" path:"triggering_user_ids"`
	TriggerByShareRecipients *bool    `json:"trigger_by_share_recipients,omitempty" path:"trigger_by_share_recipients"`
	NotifyUserActions        *bool    `json:"notify_user_actions,omitempty" path:"notify_user_actions"`
	NotifyOnCopy             *bool    `json:"notify_on_copy,omitempty" path:"notify_on_copy"`
	NotifyOnDelete           *bool    `json:"notify_on_delete,omitempty" path:"notify_on_delete"`
	NotifyOnDownload         *bool    `json:"notify_on_download,omitempty" path:"notify_on_download"`
	NotifyOnMove             *bool    `json:"notify_on_move,omitempty" path:"notify_on_move"`
	NotifyOnUpload           *bool    `json:"notify_on_upload,omitempty" path:"notify_on_upload"`
	Recursive                *bool    `json:"recursive,omitempty" path:"recursive"`
	SendInterval             string   `json:"send_interval,omitempty" path:"send_interval"`
	Message                  string   `json:"message,omitempty" path:"message"`
	TriggeringFilenames      []string `json:"triggering_filenames,omitempty" path:"triggering_filenames"`
	Unsubscribed             *bool    `json:"unsubscribed,omitempty" path:"unsubscribed"`
	UnsubscribedReason       string   `json:"unsubscribed_reason,omitempty" path:"unsubscribed_reason"`
	UserId                   int64    `json:"user_id,omitempty" path:"user_id"`
	Username                 string   `json:"username,omitempty" path:"username"`
	SuppressedEmail          *bool    `json:"suppressed_email,omitempty" path:"suppressed_email"`
}

func (n Notification) Identifier() interface{} {
	return n.Id
}

type NotificationCollection []Notification

type NotificationListParams struct {
	UserId           int64           `url:"user_id,omitempty" required:"false" json:"user_id,omitempty" path:"user_id"`
	SortBy           json.RawMessage `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty" path:"sort_by"`
	GroupId          string          `url:"group_id,omitempty" required:"false" json:"group_id,omitempty" path:"group_id"`
	Filter           json.RawMessage `url:"filter,omitempty" required:"false" json:"filter,omitempty" path:"filter"`
	FilterPrefix     json.RawMessage `url:"filter_prefix,omitempty" required:"false" json:"filter_prefix,omitempty" path:"filter_prefix"`
	Path             string          `url:"path,omitempty" required:"false" json:"path,omitempty" path:"path"`
	IncludeAncestors *bool           `url:"include_ancestors,omitempty" required:"false" json:"include_ancestors,omitempty" path:"include_ancestors"`
	ListParams
}

type NotificationFindParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
}

type NotificationCreateParams struct {
	UserId                   int64    `url:"user_id,omitempty" required:"false" json:"user_id,omitempty" path:"user_id"`
	NotifyOnCopy             *bool    `url:"notify_on_copy,omitempty" required:"false" json:"notify_on_copy,omitempty" path:"notify_on_copy"`
	NotifyOnDelete           *bool    `url:"notify_on_delete,omitempty" required:"false" json:"notify_on_delete,omitempty" path:"notify_on_delete"`
	NotifyOnDownload         *bool    `url:"notify_on_download,omitempty" required:"false" json:"notify_on_download,omitempty" path:"notify_on_download"`
	NotifyOnMove             *bool    `url:"notify_on_move,omitempty" required:"false" json:"notify_on_move,omitempty" path:"notify_on_move"`
	NotifyOnUpload           *bool    `url:"notify_on_upload,omitempty" required:"false" json:"notify_on_upload,omitempty" path:"notify_on_upload"`
	NotifyUserActions        *bool    `url:"notify_user_actions,omitempty" required:"false" json:"notify_user_actions,omitempty" path:"notify_user_actions"`
	Recursive                *bool    `url:"recursive,omitempty" required:"false" json:"recursive,omitempty" path:"recursive"`
	SendInterval             string   `url:"send_interval,omitempty" required:"false" json:"send_interval,omitempty" path:"send_interval"`
	Message                  string   `url:"message,omitempty" required:"false" json:"message,omitempty" path:"message"`
	TriggeringFilenames      []string `url:"triggering_filenames,omitempty" required:"false" json:"triggering_filenames,omitempty" path:"triggering_filenames"`
	TriggeringGroupIds       []int64  `url:"triggering_group_ids,omitempty" required:"false" json:"triggering_group_ids,omitempty" path:"triggering_group_ids"`
	TriggeringUserIds        []int64  `url:"triggering_user_ids,omitempty" required:"false" json:"triggering_user_ids,omitempty" path:"triggering_user_ids"`
	TriggerByShareRecipients *bool    `url:"trigger_by_share_recipients,omitempty" required:"false" json:"trigger_by_share_recipients,omitempty" path:"trigger_by_share_recipients"`
	GroupId                  int64    `url:"group_id,omitempty" required:"false" json:"group_id,omitempty" path:"group_id"`
	Path                     string   `url:"path,omitempty" required:"false" json:"path,omitempty" path:"path"`
	Username                 string   `url:"username,omitempty" required:"false" json:"username,omitempty" path:"username"`
}

type NotificationUpdateParams struct {
	Id                       int64    `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
	NotifyOnCopy             *bool    `url:"notify_on_copy,omitempty" required:"false" json:"notify_on_copy,omitempty" path:"notify_on_copy"`
	NotifyOnDelete           *bool    `url:"notify_on_delete,omitempty" required:"false" json:"notify_on_delete,omitempty" path:"notify_on_delete"`
	NotifyOnDownload         *bool    `url:"notify_on_download,omitempty" required:"false" json:"notify_on_download,omitempty" path:"notify_on_download"`
	NotifyOnMove             *bool    `url:"notify_on_move,omitempty" required:"false" json:"notify_on_move,omitempty" path:"notify_on_move"`
	NotifyOnUpload           *bool    `url:"notify_on_upload,omitempty" required:"false" json:"notify_on_upload,omitempty" path:"notify_on_upload"`
	NotifyUserActions        *bool    `url:"notify_user_actions,omitempty" required:"false" json:"notify_user_actions,omitempty" path:"notify_user_actions"`
	Recursive                *bool    `url:"recursive,omitempty" required:"false" json:"recursive,omitempty" path:"recursive"`
	SendInterval             string   `url:"send_interval,omitempty" required:"false" json:"send_interval,omitempty" path:"send_interval"`
	Message                  string   `url:"message,omitempty" required:"false" json:"message,omitempty" path:"message"`
	TriggeringFilenames      []string `url:"triggering_filenames,omitempty" required:"false" json:"triggering_filenames,omitempty" path:"triggering_filenames"`
	TriggeringGroupIds       []int64  `url:"triggering_group_ids,omitempty" required:"false" json:"triggering_group_ids,omitempty" path:"triggering_group_ids"`
	TriggeringUserIds        []int64  `url:"triggering_user_ids,omitempty" required:"false" json:"triggering_user_ids,omitempty" path:"triggering_user_ids"`
	TriggerByShareRecipients *bool    `url:"trigger_by_share_recipients,omitempty" required:"false" json:"trigger_by_share_recipients,omitempty" path:"trigger_by_share_recipients"`
}

type NotificationDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
}

func (n *Notification) UnmarshalJSON(data []byte) error {
	type notification Notification
	var v notification
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*n = Notification(v)
	return nil
}

func (n *NotificationCollection) UnmarshalJSON(data []byte) error {
	type notifications NotificationCollection
	var v notifications
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*n = NotificationCollection(v)
	return nil
}

func (n *NotificationCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*n))
	for i, v := range *n {
		ret[i] = v
	}

	return &ret
}
