# EventDeliveryAttempt

`files_sdk.EventDeliveryAttempt` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/eventdeliveryattempt`, named `event_delivery_attempt`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	event_delivery_attempt "github.com/Files-com/files-sdk-go/v3/eventdeliveryattempt"
)
```

## Operations

### List

`GET /event_delivery_attempts`

List Event Delivery Attempts

```go
func (c *Client) List(params files_sdk.EventDeliveryAttemptListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.EventDeliveryAttemptListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.EventDeliveryAttempt` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.EventDeliveryAttemptListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `created_at`, `status`, `event_record_id`, `event_target_id` or `workspace_id`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `status`, `workspace_id`, `event_record_id` or `event_target_id`. Valid field combinations are `[ workspace_id, status ]`, `[ workspace_id, event_record_id ]` or `[ workspace_id, event_target_id ]`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := event_delivery_attempt.List(files_sdk.EventDeliveryAttemptListParams{})
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

`GET /event_delivery_attempts/{id}`

Show Event Delivery Attempt

```go
func (c *Client) Find(params files_sdk.EventDeliveryAttemptFindParams, opts ...files_sdk.RequestResponseOption) (eventDeliveryAttempt files_sdk.EventDeliveryAttempt, err error)
func Find(params files_sdk.EventDeliveryAttemptFindParams, opts ...files_sdk.RequestResponseOption) (eventDeliveryAttempt files_sdk.EventDeliveryAttempt, err error)
```

Parameters, `files_sdk.EventDeliveryAttemptFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Event Delivery Attempt ID. |

Example:

```go
result, err := event_delivery_attempt.Find(files_sdk.EventDeliveryAttemptFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

## Model

`files_sdk.EventDeliveryAttempt` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Event Delivery Attempt ID |
| `EventRecordId` | `int64` | `event_record_id` | Event Record ID |
| `EventSubscriptionId` | `int64` | `event_subscription_id` | Event Subscription ID |
| `EventTargetId` | `int64` | `event_target_id` | Event Target ID |
| `ActivityNotificationBatchId` | `int64` | `activity_notification_batch_id` | Shared summary ID when this event is delivered in a Slack or Teams subfolder summary. All events in a summary share its delivery outcome. |
| `WorkspaceId` | `int64` | `workspace_id` | Workspace ID. 0 means the default workspace or site-wide. |
| `Status` | `string` | `status` | Delivery status. |
| `AttemptNumber` | `int64` | `attempt_number` | Number of delivery attempts made. |
| `ResponseCode` | `int64` | `response_code` | HTTP response code, if applicable. |
| `ErrorMessage` | `string` | `error_message` | Delivery error message, if applicable. |
| `ResponseBody` | `string` | `response_body` | Delivery response body, if applicable. |
| `LatencyMs` | `int64` | `latency_ms` | Delivery latency in milliseconds. |
| `DeliveredAt` | `*time.Time` | `delivered_at` | Successful delivery date/time. |
| `LastAttemptedAt` | `*time.Time` | `last_attempted_at` | Most recent attempt date/time. |
| `NextAttemptAt` | `*time.Time` | `next_attempt_at` | Next scheduled attempt date/time. |
| `CreatedAt` | `*time.Time` | `created_at` | Delivery Attempt create date/time. |
