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
	Recursive          *bool  `json:"recursive,omitempty"`
	SendInterval       string `json:"send_interval,omitempty"`
	Unsubscribed       *bool  `json:"unsubscribed,omitempty"`
	UnsubscribedReason string `json:"unsubscribed_reason,omitempty"`
	UserId             int64  `json:"user_id,omitempty"`
	Username           string `json:"username,omitempty"`
	SuppressedEmail    *bool  `json:"suppressed_email,omitempty"`
}

type NotificationCollection []Notification

type NotificationListParams struct {
	UserId           int64           `url:"user_id,omitempty" required:"false"`
	Cursor           string          `url:"cursor,omitempty" required:"false"`
	PerPage          int             `url:"per_page,omitempty" required:"false"`
	SortBy           json.RawMessage `url:"sort_by,omitempty" required:"false"`
	Filter           json.RawMessage `url:"filter,omitempty" required:"false"`
	FilterGt         json.RawMessage `url:"filter_gt,omitempty" required:"false"`
	FilterGteq       json.RawMessage `url:"filter_gteq,omitempty" required:"false"`
	FilterLike       json.RawMessage `url:"filter_like,omitempty" required:"false"`
	FilterLt         json.RawMessage `url:"filter_lt,omitempty" required:"false"`
	FilterLteq       json.RawMessage `url:"filter_lteq,omitempty" required:"false"`
	GroupId          int64           `url:"group_id,omitempty" required:"false"`
	Path             string          `url:"path,omitempty" required:"false"`
	IncludeAncestors *bool           `url:"include_ancestors,omitempty" required:"false"`
	lib.ListParams
}

type NotificationFindParams struct {
	Id int64 `url:"-,omitempty" required:"true"`
}

type NotificationCreateParams struct {
	UserId            int64  `url:"user_id,omitempty" required:"false"`
	NotifyOnCopy      *bool  `url:"notify_on_copy,omitempty" required:"false"`
	NotifyUserActions *bool  `url:"notify_user_actions,omitempty" required:"false"`
	Recursive         *bool  `url:"recursive,omitempty" required:"false"`
	SendInterval      string `url:"send_interval,omitempty" required:"false"`
	GroupId           int64  `url:"group_id,omitempty" required:"false"`
	Path              string `url:"path,omitempty" required:"false"`
	Username          string `url:"username,omitempty" required:"false"`
}

type NotificationUpdateParams struct {
	Id                int64  `url:"-,omitempty" required:"true"`
	NotifyOnCopy      *bool  `url:"notify_on_copy,omitempty" required:"false"`
	NotifyUserActions *bool  `url:"notify_user_actions,omitempty" required:"false"`
	Recursive         *bool  `url:"recursive,omitempty" required:"false"`
	SendInterval      string `url:"send_interval,omitempty" required:"false"`
}

type NotificationDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"true"`
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

func (n *NotificationCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*n))
	for i, v := range *n {
		ret[i] = v
	}

	return &ret
}
