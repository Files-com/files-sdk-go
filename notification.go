package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// Notification is a Files.com API resource.
type Notification struct {
	// Notification ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Folder path to notify on. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.
	Path string `json:"path,omitempty" path:"path,omitempty" url:"path,omitempty"`
	// ID of Group to receive notifications.
	GroupId int64 `json:"group_id,omitempty" path:"group_id,omitempty" url:"group_id,omitempty"`
	// Group name, if a Group ID is set.
	GroupName string `json:"group_name,omitempty" path:"group_name,omitempty" url:"group_name,omitempty"`
	// Group IDs when the notification requires multiple groups.
	GroupIds []int64 `json:"group_ids,omitempty" path:"group_ids,omitempty" url:"group_ids,omitempty"`
	// Group names when the notification requires multiple groups.
	GroupNames []string `json:"group_names,omitempty" path:"group_names,omitempty" url:"group_names,omitempty"`
	// If set, will only notify on actions made by a member of one of the specified groups.
	TriggeringGroupIds []int64 `json:"triggering_group_ids,omitempty" path:"triggering_group_ids,omitempty" url:"triggering_group_ids,omitempty"`
	// If set, will only notify on actions made one of the specified users.
	TriggeringUserIds []int64 `json:"triggering_user_ids,omitempty" path:"triggering_user_ids,omitempty" url:"triggering_user_ids,omitempty"`
	// Notify when actions are performed by a share recipient?
	TriggerByShareRecipients *bool `json:"trigger_by_share_recipients,omitempty" path:"trigger_by_share_recipients,omitempty" url:"trigger_by_share_recipients,omitempty"`
	// If true, will send notifications about a user's own activity to that user. If false, only activity performed by other users (or anonymous users) will be sent in notifications.
	NotifyUserActions *bool `json:"notify_user_actions,omitempty" path:"notify_user_actions,omitempty" url:"notify_user_actions,omitempty"`
	// Trigger on files copied to this path?
	NotifyOnCopy *bool `json:"notify_on_copy,omitempty" path:"notify_on_copy,omitempty" url:"notify_on_copy,omitempty"`
	// Trigger on files deleted in this path?
	NotifyOnDelete *bool `json:"notify_on_delete,omitempty" path:"notify_on_delete,omitempty" url:"notify_on_delete,omitempty"`
	// Trigger on files downloaded in this path?
	NotifyOnDownload *bool `json:"notify_on_download,omitempty" path:"notify_on_download,omitempty" url:"notify_on_download,omitempty"`
	// Trigger on files moved to this path?
	NotifyOnMove *bool `json:"notify_on_move,omitempty" path:"notify_on_move,omitempty" url:"notify_on_move,omitempty"`
	// Trigger on files created/uploaded/updated/changed in this path?
	NotifyOnUpload *bool `json:"notify_on_upload,omitempty" path:"notify_on_upload,omitempty" url:"notify_on_upload,omitempty"`
	// Apply notification recursively? This will enable notifications for each subfolder.
	Recursive *bool `json:"recursive,omitempty" path:"recursive,omitempty" url:"recursive,omitempty"`
	// The time interval that notifications are aggregated to.
	SendInterval string `json:"send_interval,omitempty" path:"send_interval,omitempty" url:"send_interval,omitempty"`
	// Custom subject line to use for notification emails.
	Subject string `json:"subject,omitempty" path:"subject,omitempty" url:"subject,omitempty"`
	// Custom message to include in notification emails.
	Message string `json:"message,omitempty" path:"message,omitempty" url:"message,omitempty"`
	// Array of filenames (possibly with wildcards) to scope trigger.
	TriggeringFilenames []string `json:"triggering_filenames,omitempty" path:"triggering_filenames,omitempty" url:"triggering_filenames,omitempty"`
	// Workspace ID. `0` means the default workspace.
	WorkspaceId int64 `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	// Is the user unsubscribed from this notification?
	Unsubscribed *bool `json:"unsubscribed,omitempty" path:"unsubscribed,omitempty" url:"unsubscribed,omitempty"`
	// The reason that the user unsubscribed.
	UnsubscribedReason string `json:"unsubscribed_reason,omitempty" path:"unsubscribed_reason,omitempty" url:"unsubscribed_reason,omitempty"`
	// Notification user ID.
	UserId int64 `json:"user_id,omitempty" path:"user_id,omitempty" url:"user_id,omitempty"`
	// Notification username.
	Username string `json:"username,omitempty" path:"username,omitempty" url:"username,omitempty"`
	// If true, it means that the recipient at this user's email address has manually unsubscribed from all emails, or had their email "hard bounce", which means that we are unable to send mail to this user's current email address. Notifications will resume if the user changes their email address.
	SuppressedEmail *bool `json:"suppressed_email,omitempty" path:"suppressed_email,omitempty" url:"suppressed_email,omitempty"`
}

// Identifier returns the resource ID.
func (n Notification) Identifier() interface{} {
	return n.Id
}

// NotificationCollection is a list of Notification resources.
type NotificationCollection []Notification

// NotificationListParams contains the request parameters for GET /notifications.
type NotificationListParams struct {
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id`, `path`, `user_id` or `group_id`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `path`, `user_id`, `workspace_id` or `group_id`. Valid field combinations are `[ workspace_id, path ]`, `[ workspace_id, user_id ]`, `[ workspace_id, group_id ]` or `[ workspace_id, user_id, path ]`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// If set, return records where the specified field is prefixed by the supplied value. Valid fields are `path`.
	FilterPrefix interface{} `url:"filter_prefix,omitempty" json:"filter_prefix,omitempty" path:"filter_prefix"`
	// Show notifications for this Path.
	Path string `url:"path,omitempty" json:"path,omitempty" path:"path"`
	// If `include_ancestors` is `true` and `path` is specified, include notifications for any parent paths. Ignored if `path` is not specified.
	IncludeAncestors *bool  `url:"include_ancestors,omitempty" json:"include_ancestors,omitempty" path:"include_ancestors"`
	GroupId          string `url:"group_id,omitempty" json:"group_id,omitempty" path:"group_id"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// NotificationFindParams contains the request parameters for GET /notifications/{id}.
type NotificationFindParams struct {
	// Notification ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// NotificationCreateParams contains the request parameters for POST /notifications.
type NotificationCreateParams struct {
	// The id of the user to notify. Provide `user_id`, `username` or `group_id`.
	UserId int64 `url:"user_id,omitempty" json:"user_id,omitempty" path:"user_id"`
	// If `true`, copying or moving resources into this path will trigger a notification, in addition to just uploads.
	NotifyOnCopy *bool `url:"notify_on_copy,omitempty" json:"notify_on_copy,omitempty" path:"notify_on_copy"`
	// Trigger on files deleted in this path?
	NotifyOnDelete *bool `url:"notify_on_delete,omitempty" json:"notify_on_delete,omitempty" path:"notify_on_delete"`
	// Trigger on files downloaded in this path?
	NotifyOnDownload *bool `url:"notify_on_download,omitempty" json:"notify_on_download,omitempty" path:"notify_on_download"`
	// Trigger on files moved to this path?
	NotifyOnMove *bool `url:"notify_on_move,omitempty" json:"notify_on_move,omitempty" path:"notify_on_move"`
	// Trigger on files created/uploaded/updated/changed in this path?
	NotifyOnUpload *bool `url:"notify_on_upload,omitempty" json:"notify_on_upload,omitempty" path:"notify_on_upload"`
	// If `true` actions initiated by the user will still result in a notification.
	NotifyUserActions *bool `url:"notify_user_actions,omitempty" json:"notify_user_actions,omitempty" path:"notify_user_actions"`
	// If `true`, enable notifications for each subfolder in this path.
	Recursive *bool `url:"recursive,omitempty" json:"recursive,omitempty" path:"recursive"`
	// The time interval that notifications are aggregated by. Can be `five_minutes`, `fifteen_minutes`, `hourly`, or `daily`.
	SendInterval string `url:"send_interval,omitempty" json:"send_interval,omitempty" path:"send_interval"`
	// Custom subject line to use for notification emails.
	Subject string `url:"subject,omitempty" json:"subject,omitempty" path:"subject"`
	// Custom message to include in notification emails.
	Message string `url:"message,omitempty" json:"message,omitempty" path:"message"`
	// Array of filenames (possibly with wildcards) to scope trigger.
	TriggeringFilenames []string `url:"triggering_filenames,omitempty" json:"triggering_filenames,omitempty" path:"triggering_filenames"`
	// If set, will only notify on actions made by a member of one of the specified groups.
	TriggeringGroupIds []int64 `url:"triggering_group_ids,omitempty" json:"triggering_group_ids,omitempty" path:"triggering_group_ids"`
	// If set, will only notify on actions made one of the specified users.
	TriggeringUserIds []int64 `url:"triggering_user_ids,omitempty" json:"triggering_user_ids,omitempty" path:"triggering_user_ids"`
	// Notify when actions are performed by a share recipient?
	TriggerByShareRecipients *bool `url:"trigger_by_share_recipients,omitempty" json:"trigger_by_share_recipients,omitempty" path:"trigger_by_share_recipients"`
	// Workspace ID. `0` means the default workspace.
	WorkspaceId int64 `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
	// The ID of the group to notify. Provide `user_id`, `username` or `group_id`.
	GroupId int64 `url:"group_id,omitempty" json:"group_id,omitempty" path:"group_id"`
	// Group IDs when the notification requires multiple groups. If sent as a string, it should be comma-delimited.
	GroupIds string `url:"group_ids,omitempty" json:"group_ids,omitempty" path:"group_ids"`
	// Path.
	Path string `url:"path,omitempty" json:"path,omitempty" path:"path"`
	// The username of the user to notify. Provide `user_id`, `username` or `group_id`.
	Username string `url:"username,omitempty" json:"username,omitempty" path:"username"`
}

// NotificationUpdateParams contains the request parameters for PATCH /notifications/{id}.
type NotificationUpdateParams struct {
	// Notification ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// If `true`, copying or moving resources into this path will trigger a notification, in addition to just uploads.
	NotifyOnCopy *bool `url:"notify_on_copy,omitempty" json:"notify_on_copy,omitempty" path:"notify_on_copy"`
	// Trigger on files deleted in this path?
	NotifyOnDelete *bool `url:"notify_on_delete,omitempty" json:"notify_on_delete,omitempty" path:"notify_on_delete"`
	// Trigger on files downloaded in this path?
	NotifyOnDownload *bool `url:"notify_on_download,omitempty" json:"notify_on_download,omitempty" path:"notify_on_download"`
	// Trigger on files moved to this path?
	NotifyOnMove *bool `url:"notify_on_move,omitempty" json:"notify_on_move,omitempty" path:"notify_on_move"`
	// Trigger on files created/uploaded/updated/changed in this path?
	NotifyOnUpload *bool `url:"notify_on_upload,omitempty" json:"notify_on_upload,omitempty" path:"notify_on_upload"`
	// If `true` actions initiated by the user will still result in a notification.
	NotifyUserActions *bool `url:"notify_user_actions,omitempty" json:"notify_user_actions,omitempty" path:"notify_user_actions"`
	// If `true`, enable notifications for each subfolder in this path.
	Recursive *bool `url:"recursive,omitempty" json:"recursive,omitempty" path:"recursive"`
	// The time interval that notifications are aggregated by. Can be `five_minutes`, `fifteen_minutes`, `hourly`, or `daily`.
	SendInterval string `url:"send_interval,omitempty" json:"send_interval,omitempty" path:"send_interval"`
	// Custom subject line to use for notification emails.
	Subject string `url:"subject,omitempty" json:"subject,omitempty" path:"subject"`
	// Custom message to include in notification emails.
	Message string `url:"message,omitempty" json:"message,omitempty" path:"message"`
	// Array of filenames (possibly with wildcards) to scope trigger.
	TriggeringFilenames []string `url:"triggering_filenames,omitempty" json:"triggering_filenames,omitempty" path:"triggering_filenames"`
	// If set, will only notify on actions made by a member of one of the specified groups.
	TriggeringGroupIds []int64 `url:"triggering_group_ids,omitempty" json:"triggering_group_ids,omitempty" path:"triggering_group_ids"`
	// If set, will only notify on actions made one of the specified users.
	TriggeringUserIds []int64 `url:"triggering_user_ids,omitempty" json:"triggering_user_ids,omitempty" path:"triggering_user_ids"`
	// Notify when actions are performed by a share recipient?
	TriggerByShareRecipients *bool `url:"trigger_by_share_recipients,omitempty" json:"trigger_by_share_recipients,omitempty" path:"trigger_by_share_recipients"`
	// Workspace ID. `0` means the default workspace.
	WorkspaceId int64 `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
}

// NotificationDeleteParams contains the request parameters for DELETE /notifications/{id}.
type NotificationDeleteParams struct {
	// Notification ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (n *Notification) UnmarshalJSON(data []byte) error {
	type notification Notification
	var v notification
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*n = Notification(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (n *NotificationCollection) UnmarshalJSON(data []byte) error {
	type notifications NotificationCollection
	var v notifications
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*n = NotificationCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (n *NotificationCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*n))
	for i, v := range *n {
		ret[i] = v
	}

	return &ret
}
