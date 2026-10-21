# EventRecord

`files_sdk.EventRecord` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/eventrecord`, named `event_record`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	event_record "github.com/Files-com/files-sdk-go/v3/eventrecord"
)
```

## Operations

### List

`GET /event_records`

List Event Records

```go
func (c *Client) List(params files_sdk.EventRecordListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.EventRecordListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.EventRecord` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.EventRecordListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `event_type`, `created_at` or `workspace_id`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `created_at`, `event_type` or `workspace_id`. Valid field combinations are `[ event_type, created_at ]`, `[ workspace_id, created_at ]`, `[ workspace_id, event_type ]` or `[ workspace_id, event_type, created_at ]`. |
| `FilterGt` | `interface{}` | `filter_gt` | no | If set, return records where the specified field is greater than the supplied value. Valid fields are `created_at`. |
| `FilterGteq` | `interface{}` | `filter_gteq` | no | If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `created_at`. |
| `FilterPrefix` | `interface{}` | `filter_prefix` | no | If set, return records where the specified field is prefixed by the supplied value. Valid fields are `event_type`. |
| `FilterLt` | `interface{}` | `filter_lt` | no | If set, return records where the specified field is less than the supplied value. Valid fields are `created_at`. |
| `FilterLteq` | `interface{}` | `filter_lteq` | no | If set, return records where the specified field is less than or equal the supplied value. Valid fields are `created_at`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := event_record.List(files_sdk.EventRecordListParams{})
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

`GET /event_records/{id}`

Show Event Record

```go
func (c *Client) Find(params files_sdk.EventRecordFindParams, opts ...files_sdk.RequestResponseOption) (eventRecord files_sdk.EventRecord, err error)
func Find(params files_sdk.EventRecordFindParams, opts ...files_sdk.RequestResponseOption) (eventRecord files_sdk.EventRecord, err error)
```

Parameters, `files_sdk.EventRecordFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Event Record ID. |

Example:

```go
result, err := event_record.Find(files_sdk.EventRecordFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

## Model

`files_sdk.EventRecord` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Event Record ID |
| `WorkspaceId` | `int64` | `workspace_id` | Workspace ID. 0 means the default workspace or site-wide. |
| `EventUuid` | `string` | `event_uuid` | Stable event UUID. |
| `EventType` | `string` | `event_type` | Versioned event type string. Valid values: `automation_run.canceled.v1`, `automation_run.failure.v1`, `automation_run.failure_will_retry.v1`, `automation_run.partial_failure.v1`, `automation_run.partial_failure_will_retry.v1`, `automation_run.skipped.v1`, `automation_run.success.v1`, `expectation_evaluation.invalid.v1`, `expectation_evaluation.late.v1`, `expectation_evaluation.missing.v1`, `expectation_evaluation.success.v1`, `expectation_incident.acknowledged.v1`, `expectation_incident.open.v1`, `expectation_incident.resolved.v1`, `expectation_incident.snoozed.v1`, `external_event.client_log.failure.v1`, `external_event.client_log.partial_failure.v1`, `external_event.client_log.skipped.v1`, `external_event.client_log.success.v1`, `pending_work_event.failure.v1`, `pending_work_event.partial_failure.v1`, `pending_work_event.skipped.v1`, `pending_work_event.success.v1`, `siem_http_destination_event.failure.v1`, `siem_http_destination_event.partial_failure.v1`, `siem_http_destination_event.skipped.v1`, `siem_http_destination_event.success.v1`, `sso_event.ldap_login.failure.v1`, `sso_event.ldap_login.partial_failure.v1`, `sso_event.ldap_login.skipped.v1`, `sso_event.ldap_login.success.v1`, `sso_event.ldap_sync.failure.v1`, `sso_event.ldap_sync.partial_failure.v1`, `sso_event.ldap_sync.skipped.v1`, `sso_event.ldap_sync.success.v1`, `sso_event.saml_login.failure.v1`, `sso_event.saml_login.partial_failure.v1`, `sso_event.saml_login.skipped.v1`, `sso_event.saml_login.success.v1`, `sync_run.failure.v1`, `sync_run.partial_failure.v1`, `sync_run.skipped.v1`, `sync_run.success.v1`, `user_security_event.lockout.v1` |
| `Severity` | `string` | `severity` | Event severity. |
| `SourceType` | `string` | `source_type` | Source record type. |
| `SourceId` | `int64` | `source_id` | Source record ID. |
| `OccurredAt` | `*time.Time` | `occurred_at` | Event occurrence date/time. |
| `HumanTitle` | `string` | `human_title` | Human-readable event title. |
| `HumanSummary` | `string` | `human_summary` | Human-readable event summary. |
| `HumanFields` | `[]map[string]interface{}` | `human_fields` | Human-readable event detail fields. |
| `Actor` | `interface{}` | `actor` | Actor associated with the event. |
| `Resources` | `[]map[string]interface{}` | `resources` | Resources associated with the event. |
| `Payload` | `interface{}` | `payload` | Event payload. |
| `CreatedAt` | `*time.Time` | `created_at` | Event Record create date/time. |
