# ActionNotificationExportResult

`files_sdk.ActionNotificationExportResult` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/actionnotificationexportresult`, named `action_notification_export_result`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	action_notification_export_result "github.com/Files-com/files-sdk-go/v3/actionnotificationexportresult"
)
```

## Operations

### List

`GET /action_notification_export_results`

List Action Notification Export Results

```go
func (c *Client) List(params files_sdk.ActionNotificationExportResultListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.ActionNotificationExportResultListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.ActionNotificationExportResult` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.ActionNotificationExportResultListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `UserId` | `int64` | `user_id` | no | User ID. Provide a value of `0` to operate the current session's user. |
| `ActionNotificationExportId` | `int64` | `action_notification_export_id` | yes | ID of the associated action notification export. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := action_notification_export_result.List(files_sdk.ActionNotificationExportResultListParams{
	ActionNotificationExportId: 1,
})
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

## Model

`files_sdk.ActionNotificationExportResult` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Notification ID |
| `CreatedAt` | `int64` | `created_at` | When the notification was sent. |
| `Status` | `int64` | `status` | HTTP status code returned in the webhook response. |
| `Message` | `string` | `message` | A message indicating the overall status of the webhook notification. |
| `Success` | `*bool` | `success` | `true` if the webhook succeeded by receiving a 200 or 204 response. |
| `RequestHeaders` | `string` | `request_headers` | A JSON-encoded string with headers that were sent with the webhook. |
| `RequestMethod` | `string` | `request_method` | The HTTP verb used to perform the webhook. |
| `RequestUrl` | `string` | `request_url` | The webhook request URL. |
| `Path` | `string` | `path` | The path to the actual file that triggered this notification. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters. |
| `Folder` | `string` | `folder` | The folder associated with the triggering action for this notification. |
