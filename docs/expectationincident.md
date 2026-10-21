# ExpectationIncident

`files_sdk.ExpectationIncident` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/expectationincident`, named `expectation_incident`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	expectation_incident "github.com/Files-com/files-sdk-go/v3/expectationincident"
)
```

## Operations

### List

`GET /expectation_incidents`

List Expectation Incidents

```go
func (c *Client) List(params files_sdk.ExpectationIncidentListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.ExpectationIncidentListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.ExpectationIncident` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.ExpectationIncidentListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id`, `created_at` or `expectation_id`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `expectation_id` and `workspace_id`. Valid field combinations are `[ workspace_id, expectation_id ]`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := expectation_incident.List(files_sdk.ExpectationIncidentListParams{})
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

`GET /expectation_incidents/{id}`

Show Expectation Incident

```go
func (c *Client) Find(params files_sdk.ExpectationIncidentFindParams, opts ...files_sdk.RequestResponseOption) (expectationIncident files_sdk.ExpectationIncident, err error)
func Find(params files_sdk.ExpectationIncidentFindParams, opts ...files_sdk.RequestResponseOption) (expectationIncident files_sdk.ExpectationIncident, err error)
```

Parameters, `files_sdk.ExpectationIncidentFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Expectation Incident ID. |

Example:

```go
result, err := expectation_incident.Find(files_sdk.ExpectationIncidentFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Resolve

`POST /expectation_incidents/{id}/resolve`

Resolve an expectation incident

```go
func (c *Client) Resolve(params files_sdk.ExpectationIncidentResolveParams, opts ...files_sdk.RequestResponseOption) (expectationIncident files_sdk.ExpectationIncident, err error)
func Resolve(params files_sdk.ExpectationIncidentResolveParams, opts ...files_sdk.RequestResponseOption) (expectationIncident files_sdk.ExpectationIncident, err error)
```

Parameters, `files_sdk.ExpectationIncidentResolveParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Expectation Incident ID. |

Example:

```go
result, err := expectation_incident.Resolve(files_sdk.ExpectationIncidentResolveParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Snooze

`POST /expectation_incidents/{id}/snooze`

Snooze an expectation incident until a specified time

```go
func (c *Client) Snooze(params files_sdk.ExpectationIncidentSnoozeParams, opts ...files_sdk.RequestResponseOption) (expectationIncident files_sdk.ExpectationIncident, err error)
func Snooze(params files_sdk.ExpectationIncidentSnoozeParams, opts ...files_sdk.RequestResponseOption) (expectationIncident files_sdk.ExpectationIncident, err error)
```

Parameters, `files_sdk.ExpectationIncidentSnoozeParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Expectation Incident ID. |
| `SnoozedUntil` | `*time.Time` | `snoozed_until` | yes | Time until which the incident should remain snoozed. |

Example:

```go
// Set this required value before the call.
var snoozedUntil *time.Time
result, err := expectation_incident.Snooze(files_sdk.ExpectationIncidentSnoozeParams{
	Id:           1,
	SnoozedUntil: snoozedUntil,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Acknowledge

`POST /expectation_incidents/{id}/acknowledge`

Acknowledge an expectation incident

```go
func (c *Client) Acknowledge(params files_sdk.ExpectationIncidentAcknowledgeParams, opts ...files_sdk.RequestResponseOption) (expectationIncident files_sdk.ExpectationIncident, err error)
func Acknowledge(params files_sdk.ExpectationIncidentAcknowledgeParams, opts ...files_sdk.RequestResponseOption) (expectationIncident files_sdk.ExpectationIncident, err error)
```

Parameters, `files_sdk.ExpectationIncidentAcknowledgeParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Expectation Incident ID. |

Example:

```go
result, err := expectation_incident.Acknowledge(files_sdk.ExpectationIncidentAcknowledgeParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

## Model

`files_sdk.ExpectationIncident` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Expectation Incident ID |
| `WorkspaceId` | `int64` | `workspace_id` | Workspace ID. `0` means the default workspace. |
| `ExpectationId` | `int64` | `expectation_id` | Expectation ID. |
| `Status` | `string` | `status` | Incident status. |
| `OpenedAt` | `*time.Time` | `opened_at` | When the incident was opened. |
| `LastFailedAt` | `*time.Time` | `last_failed_at` | When the most recent failing evaluation contributing to the incident occurred. |
| `AcknowledgedAt` | `*time.Time` | `acknowledged_at` | When the incident was acknowledged. |
| `SnoozedUntil` | `*time.Time` | `snoozed_until` | When the current snooze expires. |
| `ResolvedAt` | `*time.Time` | `resolved_at` | When the incident was resolved. |
| `OpenedByEvaluationId` | `int64` | `opened_by_evaluation_id` | Evaluation that first opened the incident. |
| `LastEvaluationId` | `int64` | `last_evaluation_id` | Most recent evaluation linked to the incident. |
| `ResolvedByEvaluationId` | `int64` | `resolved_by_evaluation_id` | Evaluation that resolved the incident. |
| `Summary` | `interface{}` | `summary` | Compact incident summary payload. |
| `CreatedAt` | `*time.Time` | `created_at` | Creation time. |
| `UpdatedAt` | `*time.Time` | `updated_at` | Last update time. |
