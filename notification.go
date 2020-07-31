package files_sdk

import (
	"encoding/json"
	lib "github.com/Files-com/files-sdk-go/lib"
)

type Notification struct {
	Id                 int64  `json:"id,omitempty"`
	Path               string `json:"path,omitempty"`
	GroupId            int64  `json:"group_id,omitempty"`
	GroupName          string `json:"group_name,omitempty"`
	NotifyUserActions  *bool  `json:"notify_user_actions,omitempty"`
	NotifyOnCopy       *bool  `json:"notify_on_copy,omitempty"`
	SendInterval       string `json:"send_interval,omitempty"`
	Unsubscribed       *bool  `json:"unsubscribed,omitempty"`
	UnsubscribedReason string `json:"unsubscribed_reason,omitempty"`
	UserId             int64  `json:"user_id,omitempty"`
	Username           string `json:"username,omitempty"`
	SuppressedEmail    *bool  `json:"suppressed_email,omitempty"`
}

type NotificationCollection []Notification

type NotificationListParams struct {
	UserId           int64           `url:"user_id,omitempty"`
	Page             int             `url:"page,omitempty"`
	PerPage          int             `url:"per_page,omitempty"`
	Action           string          `url:"action,omitempty"`
	Cursor           string          `url:"cursor,omitempty"`
	SortBy           json.RawMessage `url:"sort_by,omitempty"`
	Filter           json.RawMessage `url:"filter,omitempty"`
	FilterGt         json.RawMessage `url:"filter_gt,omitempty"`
	FilterGteq       json.RawMessage `url:"filter_gteq,omitempty"`
	FilterLike       json.RawMessage `url:"filter_like,omitempty"`
	FilterLt         json.RawMessage `url:"filter_lt,omitempty"`
	FilterLteq       json.RawMessage `url:"filter_lteq,omitempty"`
	GroupId          int64           `url:"group_id,omitempty"`
	Path             string          `url:"path,omitempty"`
	IncludeAncestors *bool           `url:"include_ancestors,omitempty"`
	lib.ListParams
}

type NotificationFindParams struct {
	Id int64 `url:"-,omitempty"`
}

type NotificationCreateParams struct {
	UserId            int64  `url:"user_id,omitempty"`
	NotifyOnCopy      *bool  `url:"notify_on_copy,omitempty"`
	NotifyUserActions *bool  `url:"notify_user_actions,omitempty"`
	SendInterval      string `url:"send_interval,omitempty"`
	GroupId           int64  `url:"group_id,omitempty"`
	Path              string `url:"path,omitempty"`
	Username          string `url:"username,omitempty"`
}

type NotificationUpdateParams struct {
	Id                int64  `url:"-,omitempty"`
	NotifyOnCopy      *bool  `url:"notify_on_copy,omitempty"`
	NotifyUserActions *bool  `url:"notify_user_actions,omitempty"`
	SendInterval      string `url:"send_interval,omitempty"`
}

type NotificationDeleteParams struct {
	Id int64 `url:"-,omitempty"`
}

func (n *Notification) UnmarshalJSON(data []byte) error {
	type notification Notification
	var v notification
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*n = Notification(v)
	return nil
}

func (n *NotificationCollection) UnmarshalJSON(data []byte) error {
	type notifications []Notification
	var v notifications
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*n = NotificationCollection(v)
	return nil
}
