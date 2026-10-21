# EventSubscription

`files_sdk.EventSubscription` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/eventsubscription`, named `event_subscription`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	event_subscription "github.com/Files-com/files-sdk-go/v3/eventsubscription"
)
```

## Operations

### List

`GET /event_subscriptions`

List Event Subscriptions

```go
func (c *Client) List(params files_sdk.EventSubscriptionListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.EventSubscriptionListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.EventSubscription` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.EventSubscriptionListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `name`, `enabled`, `event_channel_id` or `workspace_id`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `enabled`, `event_channel_id` or `workspace_id`. Valid field combinations are `[ enabled, event_channel_id ]`, `[ workspace_id, enabled ]`, `[ workspace_id, event_channel_id ]` or `[ workspace_id, enabled, event_channel_id ]`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := event_subscription.List(files_sdk.EventSubscriptionListParams{})
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

`GET /event_subscriptions/{id}`

Show Event Subscription

```go
func (c *Client) Find(params files_sdk.EventSubscriptionFindParams, opts ...files_sdk.RequestResponseOption) (eventSubscription files_sdk.EventSubscription, err error)
func Find(params files_sdk.EventSubscriptionFindParams, opts ...files_sdk.RequestResponseOption) (eventSubscription files_sdk.EventSubscription, err error)
```

Parameters, `files_sdk.EventSubscriptionFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Event Subscription ID. |

Example:

```go
result, err := event_subscription.Find(files_sdk.EventSubscriptionFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Create

`POST /event_subscriptions`

Create Event Subscription

```go
func (c *Client) Create(params files_sdk.EventSubscriptionCreateParams, opts ...files_sdk.RequestResponseOption) (eventSubscription files_sdk.EventSubscription, err error)
func Create(params files_sdk.EventSubscriptionCreateParams, opts ...files_sdk.RequestResponseOption) (eventSubscription files_sdk.EventSubscription, err error)
```

Parameters, `files_sdk.EventSubscriptionCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `EventChannelId` | `int64` | `event_channel_id` | no | Event Channel ID |
| `WorkspaceId` | `int64` | `workspace_id` | no | Workspace ID. 0 means the default workspace or site-wide. |
| `ApplyToAllWorkspaces` | `*bool` | `apply_to_all_workspaces` | no | If true, this default-workspace subscription applies to events from all workspaces. |
| `Name` | `string` | `name` | yes | Event Subscription name. |
| `Subject` | `string` | `subject` | no | Custom subject line to use for notification emails. |
| `Message` | `string` | `message` | no | Custom message to include in notification emails. |
| `MessageOnly` | `*bool` | `message_only` | no | If true, notification email bodies contain only the custom message, omitting event details and the review button. Requires a custom message, defaults to false, and does not affect non-email targets. |
| `Enabled` | `*bool` | `enabled` | no | Whether this Event Subscription can dispatch events. |
| `EventTypes` | `[]string` | `event_types` | no | Event type strings matched by this subscription. Blank means all event types. Valid values: `automation_run.canceled.v1`, `automation_run.failure.v1`, `automation_run.failure_will_retry.v1`, `automation_run.partial_failure.v1`, `automation_run.partial_failure_will_retry.v1`, `automation_run.skipped.v1`, `automation_run.success.v1`, `expectation_evaluation.invalid.v1`, `expectation_evaluation.late.v1`, `expectation_evaluation.missing.v1`, `expectation_evaluation.success.v1`, `expectation_incident.acknowledged.v1`, `expectation_incident.open.v1`, `expectation_incident.resolved.v1`, `expectation_incident.snoozed.v1`, `external_event.client_log.failure.v1`, `external_event.client_log.partial_failure.v1`, `external_event.client_log.skipped.v1`, `external_event.client_log.success.v1`, `pending_work_event.failure.v1`, `pending_work_event.partial_failure.v1`, `pending_work_event.skipped.v1`, `pending_work_event.success.v1`, `siem_http_destination_event.failure.v1`, `siem_http_destination_event.partial_failure.v1`, `siem_http_destination_event.skipped.v1`, `siem_http_destination_event.success.v1`, `sso_event.ldap_login.failure.v1`, `sso_event.ldap_login.partial_failure.v1`, `sso_event.ldap_login.skipped.v1`, `sso_event.ldap_login.success.v1`, `sso_event.ldap_sync.failure.v1`, `sso_event.ldap_sync.partial_failure.v1`, `sso_event.ldap_sync.skipped.v1`, `sso_event.ldap_sync.success.v1`, `sso_event.saml_login.failure.v1`, `sso_event.saml_login.partial_failure.v1`, `sso_event.saml_login.skipped.v1`, `sso_event.saml_login.success.v1`, `sync_run.failure.v1`, `sync_run.partial_failure.v1`, `sync_run.skipped.v1`, `sync_run.success.v1`, `user_security_event.lockout.v1` |
| `Filter` | `interface{}` | `filter` | no | Conditions on event fields that an event must meet to be delivered, keyed by dot-notation field path. All conditions must match. |
| `DeliveryPolicy` | `interface{}` | `delivery_policy` | no | Event Subscription delivery policy. |
| `EventTargetIds` | `[]int64` | `event_target_ids` | no | Event Target IDs this subscription sends to. |

Example:

```go
result, err := event_subscription.Create(files_sdk.EventSubscriptionCreateParams{
	Name: "example",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Update

`PATCH /event_subscriptions/{id}`

Update Event Subscription

```go
func (c *Client) Update(params files_sdk.EventSubscriptionUpdateParams, opts ...files_sdk.RequestResponseOption) (eventSubscription files_sdk.EventSubscription, err error)
func Update(params files_sdk.EventSubscriptionUpdateParams, opts ...files_sdk.RequestResponseOption) (eventSubscription files_sdk.EventSubscription, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (eventSubscription files_sdk.EventSubscription, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (eventSubscription files_sdk.EventSubscription, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.EventSubscriptionUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Event Subscription ID. |
| `EventChannelId` | `int64` | `event_channel_id` | no | Event Channel ID |
| `WorkspaceId` | `int64` | `workspace_id` | no | Workspace ID. 0 means the default workspace or site-wide. |
| `ApplyToAllWorkspaces` | `*bool` | `apply_to_all_workspaces` | no | If true, this default-workspace subscription applies to events from all workspaces. |
| `Name` | `string` | `name` | no | Event Subscription name. |
| `Subject` | `string` | `subject` | no | Custom subject line to use for notification emails. |
| `Message` | `string` | `message` | no | Custom message to include in notification emails. |
| `MessageOnly` | `*bool` | `message_only` | no | If true, notification email bodies contain only the custom message, omitting event details and the review button. Requires a custom message, defaults to false, and does not affect non-email targets. |
| `Enabled` | `*bool` | `enabled` | no | Whether this Event Subscription can dispatch events. |
| `EventTypes` | `[]string` | `event_types` | no | Event type strings matched by this subscription. Blank means all event types. Valid values: `automation_run.canceled.v1`, `automation_run.failure.v1`, `automation_run.failure_will_retry.v1`, `automation_run.partial_failure.v1`, `automation_run.partial_failure_will_retry.v1`, `automation_run.skipped.v1`, `automation_run.success.v1`, `expectation_evaluation.invalid.v1`, `expectation_evaluation.late.v1`, `expectation_evaluation.missing.v1`, `expectation_evaluation.success.v1`, `expectation_incident.acknowledged.v1`, `expectation_incident.open.v1`, `expectation_incident.resolved.v1`, `expectation_incident.snoozed.v1`, `external_event.client_log.failure.v1`, `external_event.client_log.partial_failure.v1`, `external_event.client_log.skipped.v1`, `external_event.client_log.success.v1`, `pending_work_event.failure.v1`, `pending_work_event.partial_failure.v1`, `pending_work_event.skipped.v1`, `pending_work_event.success.v1`, `siem_http_destination_event.failure.v1`, `siem_http_destination_event.partial_failure.v1`, `siem_http_destination_event.skipped.v1`, `siem_http_destination_event.success.v1`, `sso_event.ldap_login.failure.v1`, `sso_event.ldap_login.partial_failure.v1`, `sso_event.ldap_login.skipped.v1`, `sso_event.ldap_login.success.v1`, `sso_event.ldap_sync.failure.v1`, `sso_event.ldap_sync.partial_failure.v1`, `sso_event.ldap_sync.skipped.v1`, `sso_event.ldap_sync.success.v1`, `sso_event.saml_login.failure.v1`, `sso_event.saml_login.partial_failure.v1`, `sso_event.saml_login.skipped.v1`, `sso_event.saml_login.success.v1`, `sync_run.failure.v1`, `sync_run.partial_failure.v1`, `sync_run.skipped.v1`, `sync_run.success.v1`, `user_security_event.lockout.v1` |
| `Filter` | `interface{}` | `filter` | no | Conditions on event fields that an event must meet to be delivered, keyed by dot-notation field path. All conditions must match. |
| `DeliveryPolicy` | `interface{}` | `delivery_policy` | no | Event Subscription delivery policy. |
| `EventTargetIds` | `[]int64` | `event_target_ids` | no | Event Target IDs this subscription sends to. |

Example:

```go
result, err := event_subscription.Update(files_sdk.EventSubscriptionUpdateParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := event_subscription.UpdateWithMap(map[string]interface{}{
	"id": 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /event_subscriptions/{id}`

Delete Event Subscription

```go
func (c *Client) Delete(params files_sdk.EventSubscriptionDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.EventSubscriptionDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.EventSubscriptionDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Event Subscription ID. |

Example:

```go
err := event_subscription.Delete(files_sdk.EventSubscriptionDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.EventSubscription` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Event Subscription ID |
| `EventChannelId` | `int64` | `event_channel_id` | Event Channel ID |
| `WorkspaceId` | `int64` | `workspace_id` | Workspace ID. 0 means the default workspace or site-wide. |
| `ApplyToAllWorkspaces` | `*bool` | `apply_to_all_workspaces` | If true, this default-workspace subscription applies to events from all workspaces. |
| `Name` | `string` | `name` | Event Subscription name. |
| `Subject` | `string` | `subject` | Custom subject line to use for notification emails. |
| `Message` | `string` | `message` | Custom message to include in notification emails. |
| `MessageOnly` | `*bool` | `message_only` | If true, notification email bodies contain only the custom message, omitting event details and the review button. Requires a custom message, defaults to false, and does not affect non-email targets. |
| `Enabled` | `*bool` | `enabled` | Whether this Event Subscription can dispatch events. |
| `EventTypes` | `[]string` | `event_types` | Event type strings matched by this subscription. Blank means all event types. Valid values: `automation_run.canceled.v1`, `automation_run.failure.v1`, `automation_run.failure_will_retry.v1`, `automation_run.partial_failure.v1`, `automation_run.partial_failure_will_retry.v1`, `automation_run.skipped.v1`, `automation_run.success.v1`, `expectation_evaluation.invalid.v1`, `expectation_evaluation.late.v1`, `expectation_evaluation.missing.v1`, `expectation_evaluation.success.v1`, `expectation_incident.acknowledged.v1`, `expectation_incident.open.v1`, `expectation_incident.resolved.v1`, `expectation_incident.snoozed.v1`, `external_event.client_log.failure.v1`, `external_event.client_log.partial_failure.v1`, `external_event.client_log.skipped.v1`, `external_event.client_log.success.v1`, `pending_work_event.failure.v1`, `pending_work_event.partial_failure.v1`, `pending_work_event.skipped.v1`, `pending_work_event.success.v1`, `siem_http_destination_event.failure.v1`, `siem_http_destination_event.partial_failure.v1`, `siem_http_destination_event.skipped.v1`, `siem_http_destination_event.success.v1`, `sso_event.ldap_login.failure.v1`, `sso_event.ldap_login.partial_failure.v1`, `sso_event.ldap_login.skipped.v1`, `sso_event.ldap_login.success.v1`, `sso_event.ldap_sync.failure.v1`, `sso_event.ldap_sync.partial_failure.v1`, `sso_event.ldap_sync.skipped.v1`, `sso_event.ldap_sync.success.v1`, `sso_event.saml_login.failure.v1`, `sso_event.saml_login.partial_failure.v1`, `sso_event.saml_login.skipped.v1`, `sso_event.saml_login.success.v1`, `sync_run.failure.v1`, `sync_run.partial_failure.v1`, `sync_run.skipped.v1`, `sync_run.success.v1`, `user_security_event.lockout.v1` |
| `Filter` | `interface{}` | `filter` | Conditions on event fields that an event must meet to be delivered, keyed by dot-notation field path. All conditions must match. |
| `DeliveryPolicy` | `interface{}` | `delivery_policy` | Event Subscription delivery policy. |
| `EventTargetIds` | `[]int64` | `event_target_ids` | Event Target IDs this subscription sends to. |
| `CreatedAt` | `*time.Time` | `created_at` | Event Subscription create date/time. |
| `UpdatedAt` | `*time.Time` | `updated_at` | Event Subscription update date/time. |
