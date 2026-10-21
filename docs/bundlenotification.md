# BundleNotification

`files_sdk.BundleNotification` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/bundlenotification`, named `bundle_notification`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	bundle_notification "github.com/Files-com/files-sdk-go/v3/bundlenotification"
)
```

## Operations

### List

`GET /bundle_notifications`

List Share Link Notifications

```go
func (c *Client) List(params files_sdk.BundleNotificationListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.BundleNotificationListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.BundleNotification` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.BundleNotificationListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `UserId` | `int64` | `user_id` | no | User ID. Provide a value of `0` to operate the current session's user. |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id` and `bundle_id`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `bundle_id`. |
| `BundleId` | `int64` | `bundle_id` | no | Bundle ID |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := bundle_notification.List(files_sdk.BundleNotificationListParams{})
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

`GET /bundle_notifications/{id}`

Show Share Link Notification

```go
func (c *Client) Find(params files_sdk.BundleNotificationFindParams, opts ...files_sdk.RequestResponseOption) (bundleNotification files_sdk.BundleNotification, err error)
func Find(params files_sdk.BundleNotificationFindParams, opts ...files_sdk.RequestResponseOption) (bundleNotification files_sdk.BundleNotification, err error)
```

Parameters, `files_sdk.BundleNotificationFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Bundle Notification ID. |

Example:

```go
result, err := bundle_notification.Find(files_sdk.BundleNotificationFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Create

`POST /bundle_notifications`

Create Share Link Notification

```go
func (c *Client) Create(params files_sdk.BundleNotificationCreateParams, opts ...files_sdk.RequestResponseOption) (bundleNotification files_sdk.BundleNotification, err error)
func Create(params files_sdk.BundleNotificationCreateParams, opts ...files_sdk.RequestResponseOption) (bundleNotification files_sdk.BundleNotification, err error)
```

Parameters, `files_sdk.BundleNotificationCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `UserId` | `int64` | `user_id` | no | User ID. Provide a value of `0` to operate the current session's user. |
| `BundleId` | `int64` | `bundle_id` | yes | Bundle ID to notify on |
| `NotifyUserId` | `int64` | `notify_user_id` | no | The id of the user to notify. |
| `NotifyOnRegistration` | `*bool` | `notify_on_registration` | no | Triggers bundle notification when a registration action occurs for it. |
| `NotifyOnUpload` | `*bool` | `notify_on_upload` | no | Triggers bundle notification when a upload action occurs for it. |

Example:

```go
result, err := bundle_notification.Create(files_sdk.BundleNotificationCreateParams{
	BundleId: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Update

`PATCH /bundle_notifications/{id}`

Update Share Link Notification

```go
func (c *Client) Update(params files_sdk.BundleNotificationUpdateParams, opts ...files_sdk.RequestResponseOption) (bundleNotification files_sdk.BundleNotification, err error)
func Update(params files_sdk.BundleNotificationUpdateParams, opts ...files_sdk.RequestResponseOption) (bundleNotification files_sdk.BundleNotification, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (bundleNotification files_sdk.BundleNotification, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (bundleNotification files_sdk.BundleNotification, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.BundleNotificationUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Bundle Notification ID. |
| `NotifyOnRegistration` | `*bool` | `notify_on_registration` | no | Triggers bundle notification when a registration action occurs for it. |
| `NotifyOnUpload` | `*bool` | `notify_on_upload` | no | Triggers bundle notification when a upload action occurs for it. |

Example:

```go
result, err := bundle_notification.Update(files_sdk.BundleNotificationUpdateParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := bundle_notification.UpdateWithMap(map[string]interface{}{
	"id": 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /bundle_notifications/{id}`

Delete Share Link Notification

```go
func (c *Client) Delete(params files_sdk.BundleNotificationDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.BundleNotificationDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.BundleNotificationDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Bundle Notification ID. |

Example:

```go
err := bundle_notification.Delete(files_sdk.BundleNotificationDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.BundleNotification` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `BundleId` | `int64` | `bundle_id` | Bundle ID to notify on |
| `Id` | `int64` | `id` | Bundle Notification ID |
| `NotifyOnRegistration` | `*bool` | `notify_on_registration` | Triggers bundle notification when a registration action occurs for it. |
| `NotifyOnUpload` | `*bool` | `notify_on_upload` | Triggers bundle notification when a upload action occurs for it. |
| `NotifyCurrentUser` | `*bool` | `notify_current_user` | Is the current user the user to notify? |
| `NotifyUserId` | `int64` | `notify_user_id` | The id of the user to notify. |
| `WorkspaceId` | `int64` | `workspace_id` | Workspace ID. `0` means the default workspace. |
| `UserId` | `int64` | `user_id` | User ID. Provide a value of `0` to operate the current session's user. |
