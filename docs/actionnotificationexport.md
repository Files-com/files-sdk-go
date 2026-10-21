# ActionNotificationExport

`files_sdk.ActionNotificationExport` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/actionnotificationexport`, named `action_notification_export`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	action_notification_export "github.com/Files-com/files-sdk-go/v3/actionnotificationexport"
)
```

## Operations

### Find

`GET /action_notification_exports/{id}`

Show Action Notification Export

```go
func (c *Client) Find(params files_sdk.ActionNotificationExportFindParams, opts ...files_sdk.RequestResponseOption) (actionNotificationExport files_sdk.ActionNotificationExport, err error)
func Find(params files_sdk.ActionNotificationExportFindParams, opts ...files_sdk.RequestResponseOption) (actionNotificationExport files_sdk.ActionNotificationExport, err error)
```

Parameters, `files_sdk.ActionNotificationExportFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Action Notification Export ID. |

Example:

```go
result, err := action_notification_export.Find(files_sdk.ActionNotificationExportFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Create

`POST /action_notification_exports`

Create Action Notification Export

```go
func (c *Client) Create(params files_sdk.ActionNotificationExportCreateParams, opts ...files_sdk.RequestResponseOption) (actionNotificationExport files_sdk.ActionNotificationExport, err error)
func Create(params files_sdk.ActionNotificationExportCreateParams, opts ...files_sdk.RequestResponseOption) (actionNotificationExport files_sdk.ActionNotificationExport, err error)
```

Parameters, `files_sdk.ActionNotificationExportCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `UserId` | `int64` | `user_id` | no | User ID. Provide a value of `0` to operate the current session's user. |
| `WorkspaceId` | `int64` | `workspace_id` | no | Workspace whose logs are exported. Set to `0` for the default workspace. A null value means a site-wide export. |
| `StartAt` | `*time.Time` | `start_at` | no | Start date/time of export range. |
| `EndAt` | `*time.Time` | `end_at` | no | End date/time of export range. |
| `QueryMessage` | `string` | `query_message` | no | Error message associated with the request, if any. |
| `QueryRequestMethod` | `string` | `query_request_method` | no | The HTTP request method used by the webhook. |
| `QueryRequestUrl` | `string` | `query_request_url` | no | The target webhook URL. |
| `QueryStatus` | `string` | `query_status` | no | The HTTP status returned from the server in response to the webhook request. |
| `QuerySuccess` | `*bool` | `query_success` | no | true if the webhook request succeeded (i.e. returned a 200 or 204 response status). false otherwise. |
| `QueryPath` | `string` | `query_path` | no | Return notifications that were triggered by actions on this specific path. |
| `QueryFolder` | `string` | `query_folder` | no | Return notifications that were triggered by actions in this folder. |

Example:

```go
result, err := action_notification_export.Create(files_sdk.ActionNotificationExportCreateParams{})
if err != nil {
	return err
}
fmt.Println(result)
```

## Model

`files_sdk.ActionNotificationExport` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | History Export ID |
| `WorkspaceId` | `int64` | `workspace_id` | Workspace whose logs are exported. Set to `0` for the default workspace. A null value means a site-wide export. |
| `ExportVersion` | `string` | `export_version` | Version of the underlying records for the export. |
| `StartAt` | `*time.Time` | `start_at` | Start date/time of export range. |
| `EndAt` | `*time.Time` | `end_at` | End date/time of export range. |
| `Status` | `string` | `status` | Status of export. Valid values: `building`, `ready`, or `failed` |
| `QueryPath` | `string` | `query_path` | Return notifications that were triggered by actions on this specific path. |
| `QueryFolder` | `string` | `query_folder` | Return notifications that were triggered by actions in this folder. |
| `QueryMessage` | `string` | `query_message` | Error message associated with the request, if any. |
| `QueryRequestMethod` | `string` | `query_request_method` | The HTTP request method used by the webhook. |
| `QueryRequestUrl` | `string` | `query_request_url` | The target webhook URL. |
| `QueryStatus` | `string` | `query_status` | The HTTP status returned from the server in response to the webhook request. |
| `QuerySuccess` | `*bool` | `query_success` | true if the webhook request succeeded (i.e. returned a 200 or 204 response status). false otherwise. |
| `ResultsUrl` | `string` | `results_url` | If `status` is `ready`, this will be a URL where all the results can be downloaded at once as a CSV. |
| `UserId` | `int64` | `user_id` | User ID. Provide a value of `0` to operate the current session's user. |
