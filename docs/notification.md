# Notification

`files_sdk.Notification` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/notification`, named `notification`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/notification"
)
```

## Operations

### List

`GET /notifications`

List Notifications

```go
func (c *Client) List(params files_sdk.NotificationListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.NotificationListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.Notification` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.NotificationListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id`, `path`, `user_id` or `group_id`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `path`, `user_id`, `workspace_id` or `group_id`. Valid field combinations are `[ workspace_id, path ]`, `[ workspace_id, user_id ]`, `[ workspace_id, group_id ]` or `[ workspace_id, user_id, path ]`. |
| `FilterPrefix` | `interface{}` | `filter_prefix` | no | If set, return records where the specified field is prefixed by the supplied value. Valid fields are `path`. |
| `Path` | `string` | `path` | no | Show notifications for this Path. |
| `IncludeAncestors` | `*bool` | `include_ancestors` | no | If `include_ancestors` is `true` and `path` is specified, include notifications for any parent paths. Ignored if `path` is not specified. |
| `GroupId` | `string` | `group_id` | no |  |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := notification.List(files_sdk.NotificationListParams{})
if err != nil {
	return err
}
for item, err := range it.All() {
	if err != nil {
		return err
	}
	fmt.Println(item)
}
```

### Find

`GET /notifications/{id}`

Show Notification

```go
func (c *Client) Find(params files_sdk.NotificationFindParams, opts ...files_sdk.RequestResponseOption) (notification files_sdk.Notification, err error)
func Find(params files_sdk.NotificationFindParams, opts ...files_sdk.RequestResponseOption) (notification files_sdk.Notification, err error)
```

Parameters, `files_sdk.NotificationFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Notification ID. |

Example:

```go
result, err := notification.Find(files_sdk.NotificationFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Create

`POST /notifications`

Create Notification

```go
func (c *Client) Create(params files_sdk.NotificationCreateParams, opts ...files_sdk.RequestResponseOption) (notification files_sdk.Notification, err error)
func Create(params files_sdk.NotificationCreateParams, opts ...files_sdk.RequestResponseOption) (notification files_sdk.Notification, err error)
```

Parameters, `files_sdk.NotificationCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `UserId` | `int64` | `user_id` | no | The id of the user to notify. Provide `user_id`, `username` or `group_id`. |
| `NotifyOnCopy` | `*bool` | `notify_on_copy` | no | If `true`, copying or moving resources into this path will trigger a notification, in addition to just uploads. |
| `NotifyOnDelete` | `*bool` | `notify_on_delete` | no | Trigger on files deleted in this path? |
| `NotifyOnDownload` | `*bool` | `notify_on_download` | no | Trigger on files downloaded in this path? |
| `NotifyOnMove` | `*bool` | `notify_on_move` | no | Trigger on files moved to this path? |
| `NotifyOnUpload` | `*bool` | `notify_on_upload` | no | Trigger on files created/uploaded/updated/changed in this path? |
| `NotifyUserActions` | `*bool` | `notify_user_actions` | no | If `true` actions initiated by the user will still result in a notification |
| `Recursive` | `*bool` | `recursive` | no | If `true`, enable notifications for each subfolder in this path |
| `SendInterval` | `string` | `send_interval` | no | The time interval that notifications are aggregated by. Can be `five_minutes`, `fifteen_minutes`, `hourly`, or `daily`. |
| `Subject` | `string` | `subject` | no | Custom subject line to use for notification emails |
| `Message` | `string` | `message` | no | Custom message to include in notification emails |
| `TriggeringFilenames` | `[]string` | `triggering_filenames` | no | Array of filenames (possibly with wildcards) to scope trigger |
| `TriggeringGroupIds` | `[]int64` | `triggering_group_ids` | no | If set, will only notify on actions made by a member of one of the specified groups |
| `TriggeringUserIds` | `[]int64` | `triggering_user_ids` | no | If set, will only notify on actions made one of the specified users |
| `TriggerByShareRecipients` | `*bool` | `trigger_by_share_recipients` | no | Notify when actions are performed by a share recipient? |
| `WorkspaceId` | `int64` | `workspace_id` | no | Workspace ID. `0` means the default workspace. |
| `GroupId` | `int64` | `group_id` | no | The ID of the group to notify. Provide `user_id`, `username` or `group_id`. |
| `GroupIds` | `string` | `group_ids` | no | Group IDs when the notification requires multiple groups. If sent as a string, it should be comma-delimited. |
| `Path` | `string` | `path` | no | Path |
| `Username` | `string` | `username` | no | The username of the user to notify. Provide `user_id`, `username` or `group_id`. |

Example:

```go
result, err := notification.Create(files_sdk.NotificationCreateParams{})
if err != nil {
	return err
}
fmt.Println(result)
```

### Update

`PATCH /notifications/{id}`

Update Notification

```go
func (c *Client) Update(params files_sdk.NotificationUpdateParams, opts ...files_sdk.RequestResponseOption) (notification files_sdk.Notification, err error)
func Update(params files_sdk.NotificationUpdateParams, opts ...files_sdk.RequestResponseOption) (notification files_sdk.Notification, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (notification files_sdk.Notification, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (notification files_sdk.Notification, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.NotificationUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Notification ID. |
| `NotifyOnCopy` | `*bool` | `notify_on_copy` | no | If `true`, copying or moving resources into this path will trigger a notification, in addition to just uploads. |
| `NotifyOnDelete` | `*bool` | `notify_on_delete` | no | Trigger on files deleted in this path? |
| `NotifyOnDownload` | `*bool` | `notify_on_download` | no | Trigger on files downloaded in this path? |
| `NotifyOnMove` | `*bool` | `notify_on_move` | no | Trigger on files moved to this path? |
| `NotifyOnUpload` | `*bool` | `notify_on_upload` | no | Trigger on files created/uploaded/updated/changed in this path? |
| `NotifyUserActions` | `*bool` | `notify_user_actions` | no | If `true` actions initiated by the user will still result in a notification |
| `Recursive` | `*bool` | `recursive` | no | If `true`, enable notifications for each subfolder in this path |
| `SendInterval` | `string` | `send_interval` | no | The time interval that notifications are aggregated by. Can be `five_minutes`, `fifteen_minutes`, `hourly`, or `daily`. |
| `Subject` | `string` | `subject` | no | Custom subject line to use for notification emails |
| `Message` | `string` | `message` | no | Custom message to include in notification emails |
| `TriggeringFilenames` | `[]string` | `triggering_filenames` | no | Array of filenames (possibly with wildcards) to scope trigger |
| `TriggeringGroupIds` | `[]int64` | `triggering_group_ids` | no | If set, will only notify on actions made by a member of one of the specified groups |
| `TriggeringUserIds` | `[]int64` | `triggering_user_ids` | no | If set, will only notify on actions made one of the specified users |
| `TriggerByShareRecipients` | `*bool` | `trigger_by_share_recipients` | no | Notify when actions are performed by a share recipient? |
| `WorkspaceId` | `int64` | `workspace_id` | no | Workspace ID. `0` means the default workspace. |

Example:

```go
result, err := notification.Update(files_sdk.NotificationUpdateParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := notification.UpdateWithMap(map[string]interface{}{
	"id": 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /notifications/{id}`

Delete Notification

```go
func (c *Client) Delete(params files_sdk.NotificationDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.NotificationDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.NotificationDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Notification ID. |

Example:

```go
err := notification.Delete(files_sdk.NotificationDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.Notification` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Notification ID |
| `Path` | `string` | `path` | Folder path to notify on. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters. |
| `GroupId` | `int64` | `group_id` | ID of Group to receive notifications |
| `GroupName` | `string` | `group_name` | Group name, if a Group ID is set |
| `GroupIds` | `[]int64` | `group_ids` | Group IDs when the notification requires multiple groups |
| `GroupNames` | `[]string` | `group_names` | Group names when the notification requires multiple groups |
| `TriggeringGroupIds` | `[]int64` | `triggering_group_ids` | If set, will only notify on actions made by a member of one of the specified groups |
| `TriggeringUserIds` | `[]int64` | `triggering_user_ids` | If set, will only notify on actions made one of the specified users |
| `TriggerByShareRecipients` | `*bool` | `trigger_by_share_recipients` | Notify when actions are performed by a share recipient? |
| `NotifyUserActions` | `*bool` | `notify_user_actions` | If true, will send notifications about a user's own activity to that user. If false, only activity performed by other users (or anonymous users) will be sent in notifications. |
| `NotifyOnCopy` | `*bool` | `notify_on_copy` | Trigger on files copied to this path? |
| `NotifyOnDelete` | `*bool` | `notify_on_delete` | Trigger on files deleted in this path? |
| `NotifyOnDownload` | `*bool` | `notify_on_download` | Trigger on files downloaded in this path? |
| `NotifyOnMove` | `*bool` | `notify_on_move` | Trigger on files moved to this path? |
| `NotifyOnUpload` | `*bool` | `notify_on_upload` | Trigger on files created/uploaded/updated/changed in this path? |
| `Recursive` | `*bool` | `recursive` | Apply notification recursively? This will enable notifications for each subfolder. |
| `SendInterval` | `string` | `send_interval` | The time interval that notifications are aggregated to |
| `Subject` | `string` | `subject` | Custom subject line to use for notification emails |
| `Message` | `string` | `message` | Custom message to include in notification emails |
| `TriggeringFilenames` | `[]string` | `triggering_filenames` | Array of filenames (possibly with wildcards) to scope trigger |
| `WorkspaceId` | `int64` | `workspace_id` | Workspace ID. `0` means the default workspace. |
| `Unsubscribed` | `*bool` | `unsubscribed` | Is the user unsubscribed from this notification? |
| `UnsubscribedReason` | `string` | `unsubscribed_reason` | The reason that the user unsubscribed |
| `UserId` | `int64` | `user_id` | Notification user ID |
| `Username` | `string` | `username` | Notification username |
| `SuppressedEmail` | `*bool` | `suppressed_email` | If true, it means that the recipient at this user's email address has manually unsubscribed from all emails, or had their email "hard bounce", which means that we are unable to send mail to this user's current email address. Notifications will resume if the user changes their email address. |
