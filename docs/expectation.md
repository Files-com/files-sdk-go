# Expectation

`files_sdk.Expectation` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/expectation`, named `expectation`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/expectation"
)
```

## Operations

### List

`GET /expectations`

List Expectations

```go
func (c *Client) List(params files_sdk.ExpectationListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.ExpectationListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.Expectation` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.ExpectationListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id`, `name` or `disabled`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `disabled` and `workspace_id`. Valid field combinations are `[ workspace_id, disabled ]`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := expectation.List(files_sdk.ExpectationListParams{})
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

`GET /expectations/{id}`

Show Expectation

```go
func (c *Client) Find(params files_sdk.ExpectationFindParams, opts ...files_sdk.RequestResponseOption) (expectation files_sdk.Expectation, err error)
func Find(params files_sdk.ExpectationFindParams, opts ...files_sdk.RequestResponseOption) (expectation files_sdk.Expectation, err error)
```

Parameters, `files_sdk.ExpectationFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Expectation ID. |

Example:

```go
result, err := expectation.Find(files_sdk.ExpectationFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Create

`POST /expectations`

Create Expectation

```go
func (c *Client) Create(params files_sdk.ExpectationCreateParams, opts ...files_sdk.RequestResponseOption) (expectation files_sdk.Expectation, err error)
func Create(params files_sdk.ExpectationCreateParams, opts ...files_sdk.RequestResponseOption) (expectation files_sdk.Expectation, err error)
```

Parameters, `files_sdk.ExpectationCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Name` | `string` | `name` | no | Expectation name. |
| `Description` | `string` | `description` | no | Expectation description. |
| `Path` | `string` | `path` | no | Path scope for the expectation. Supports workspace-relative presentation. |
| `Source` | `string` | `source` | no | Source glob used to select candidate files. |
| `ExcludePattern` | `string` | `exclude_pattern` | no | Optional source exclusion glob. |
| `Disabled` | `*bool` | `disabled` | no | If true, the expectation is disabled. |
| `Trigger` | `files_sdk.ExpectationTriggerEnum` | `trigger` | no | How this expectation opens windows. Values: `manual`, `upload`, `daily`, `custom_schedule`. |
| `Interval` | `string` | `interval` | no | If trigger is `daily`, this specifies how often to run the expectation. |
| `RecurringDay` | `int64` | `recurring_day` | no | If trigger is `daily`, this selects the day number inside the chosen interval. |
| `RecurringDays` | `[]int64` | `recurring_days` | no | If trigger is `daily`, this selects one or more day numbers inside a `week`, `month`, `quarter`, or `year` interval. |
| `ScheduleId` | `int64` | `schedule_id` | no | If trigger is `custom_schedule`, the reusable Schedule used instead of the Expectation's schedule fields. |
| `ScheduleDaysOfWeek` | `[]int64` | `schedule_days_of_week` | no | If trigger is `custom_schedule`, the 0-based weekdays used by the schedule. |
| `ScheduleTimesOfDay` | `[]string` | `schedule_times_of_day` | no | Times of day in HH:MM format for the Expectation schedule. |
| `ScheduleTimeZone` | `string` | `schedule_time_zone` | no | Time zone used by the Expectation schedule. |
| `HolidayRegion` | `string` | `holiday_region` | no | Optional holiday region used by the Expectation schedule. |
| `LookbackInterval` | `int64` | `lookback_interval` | no | How many seconds before the due boundary the window starts. |
| `LateAcceptanceInterval` | `int64` | `late_acceptance_interval` | no | How many seconds a schedule-driven window may remain eligible to close as late. |
| `InactivityInterval` | `int64` | `inactivity_interval` | no | How many quiet seconds are required before final closure. |
| `MaxOpenInterval` | `int64` | `max_open_interval` | no | Hard-stop duration in seconds for unscheduled expectations. |
| `Criteria` | `interface{}` | `criteria` | no | Versioned success criteria definition for the expectation, including optional Files Transform Script content validation in criteria v2. |
| `WorkspaceId` | `int64` | `workspace_id` | no | Workspace ID. `0` means the default workspace. |

Example:

```go
result, err := expectation.Create(files_sdk.ExpectationCreateParams{})
if err != nil {
	return err
}
fmt.Println(result)
```

### TriggerEvaluation

`POST /expectations/{id}/trigger_evaluation`

Manually open an Expectation window

```go
func (c *Client) TriggerEvaluation(params files_sdk.ExpectationTriggerEvaluationParams, opts ...files_sdk.RequestResponseOption) (expectationEvaluation files_sdk.ExpectationEvaluation, err error)
func TriggerEvaluation(params files_sdk.ExpectationTriggerEvaluationParams, opts ...files_sdk.RequestResponseOption) (expectationEvaluation files_sdk.ExpectationEvaluation, err error)
```

Parameters, `files_sdk.ExpectationTriggerEvaluationParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Expectation ID. |

Example:

```go
result, err := expectation.TriggerEvaluation(files_sdk.ExpectationTriggerEvaluationParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Update

`PATCH /expectations/{id}`

Update Expectation

```go
func (c *Client) Update(params files_sdk.ExpectationUpdateParams, opts ...files_sdk.RequestResponseOption) (expectation files_sdk.Expectation, err error)
func Update(params files_sdk.ExpectationUpdateParams, opts ...files_sdk.RequestResponseOption) (expectation files_sdk.Expectation, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (expectation files_sdk.Expectation, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (expectation files_sdk.Expectation, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.ExpectationUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Expectation ID. |
| `Name` | `string` | `name` | no | Expectation name. |
| `Description` | `string` | `description` | no | Expectation description. |
| `Path` | `string` | `path` | no | Path scope for the expectation. Supports workspace-relative presentation. |
| `Source` | `string` | `source` | no | Source glob used to select candidate files. |
| `ExcludePattern` | `string` | `exclude_pattern` | no | Optional source exclusion glob. |
| `Disabled` | `*bool` | `disabled` | no | If true, the expectation is disabled. |
| `Trigger` | `files_sdk.ExpectationTriggerEnum` | `trigger` | no | How this expectation opens windows. Values: `manual`, `upload`, `daily`, `custom_schedule`. |
| `Interval` | `string` | `interval` | no | If trigger is `daily`, this specifies how often to run the expectation. |
| `RecurringDay` | `int64` | `recurring_day` | no | If trigger is `daily`, this selects the day number inside the chosen interval. |
| `RecurringDays` | `[]int64` | `recurring_days` | no | If trigger is `daily`, this selects one or more day numbers inside a `week`, `month`, `quarter`, or `year` interval. |
| `ScheduleId` | `int64` | `schedule_id` | no | If trigger is `custom_schedule`, the reusable Schedule used instead of the Expectation's schedule fields. |
| `ScheduleDaysOfWeek` | `[]int64` | `schedule_days_of_week` | no | If trigger is `custom_schedule`, the 0-based weekdays used by the schedule. |
| `ScheduleTimesOfDay` | `[]string` | `schedule_times_of_day` | no | Times of day in HH:MM format for the Expectation schedule. |
| `ScheduleTimeZone` | `string` | `schedule_time_zone` | no | Time zone used by the Expectation schedule. |
| `HolidayRegion` | `string` | `holiday_region` | no | Optional holiday region used by the Expectation schedule. |
| `LookbackInterval` | `int64` | `lookback_interval` | no | How many seconds before the due boundary the window starts. |
| `LateAcceptanceInterval` | `int64` | `late_acceptance_interval` | no | How many seconds a schedule-driven window may remain eligible to close as late. |
| `InactivityInterval` | `int64` | `inactivity_interval` | no | How many quiet seconds are required before final closure. |
| `MaxOpenInterval` | `int64` | `max_open_interval` | no | Hard-stop duration in seconds for unscheduled expectations. |
| `Criteria` | `interface{}` | `criteria` | no | Versioned success criteria definition for the expectation, including optional Files Transform Script content validation in criteria v2. |
| `WorkspaceId` | `int64` | `workspace_id` | no | Workspace ID. `0` means the default workspace. |

Example:

```go
result, err := expectation.Update(files_sdk.ExpectationUpdateParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := expectation.UpdateWithMap(map[string]interface{}{
	"id": 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /expectations/{id}`

Delete Expectation

```go
func (c *Client) Delete(params files_sdk.ExpectationDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.ExpectationDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.ExpectationDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Expectation ID. |

Example:

```go
err := expectation.Delete(files_sdk.ExpectationDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.Expectation` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Expectation ID |
| `WorkspaceId` | `int64` | `workspace_id` | Workspace ID. `0` means the default workspace. |
| `Name` | `string` | `name` | Expectation name. |
| `Description` | `string` | `description` | Expectation description. |
| `Path` | `string` | `path` | Path scope for the expectation. Supports workspace-relative presentation. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters. |
| `Source` | `string` | `source` | Source glob used to select candidate files. |
| `ExcludePattern` | `string` | `exclude_pattern` | Optional source exclusion glob. |
| `Disabled` | `*bool` | `disabled` | If true, the expectation is disabled. |
| `ExpectationsVersion` | `int64` | `expectations_version` | Criteria schema version for this expectation. |
| `Trigger` | `string` | `trigger` | How this expectation opens windows. |
| `Interval` | `string` | `interval` | If trigger is `daily`, this specifies how often to run the expectation. |
| `RecurringDay` | `int64` | `recurring_day` | If trigger is `daily`, this selects the day number inside the chosen interval. |
| `RecurringDays` | `[]int64` | `recurring_days` | If trigger is `daily`, this selects one or more day numbers inside a `week`, `month`, `quarter`, or `year` interval. |
| `ScheduleId` | `int64` | `schedule_id` | If trigger is `custom_schedule`, the reusable Schedule used instead of the Expectation's schedule fields. |
| `ScheduleDaysOfWeek` | `[]int64` | `schedule_days_of_week` | If trigger is `custom_schedule`, the 0-based weekdays used by the schedule. |
| `ScheduleTimesOfDay` | `[]string` | `schedule_times_of_day` | Times of day in HH:MM format for the Expectation schedule. |
| `ScheduleTimeZone` | `string` | `schedule_time_zone` | Time zone used by the Expectation schedule. |
| `HolidayRegion` | `string` | `holiday_region` | Optional holiday region used by the Expectation schedule. |
| `LookbackInterval` | `int64` | `lookback_interval` | How many seconds before the due boundary the window starts. |
| `LateAcceptanceInterval` | `int64` | `late_acceptance_interval` | How many seconds a schedule-driven window may remain eligible to close as late. |
| `InactivityInterval` | `int64` | `inactivity_interval` | How many quiet seconds are required before final closure. |
| `MaxOpenInterval` | `int64` | `max_open_interval` | Hard-stop duration in seconds for unscheduled expectations. |
| `Criteria` | `interface{}` | `criteria` | Versioned success criteria definition for the expectation. Criteria v2 supports optional FTS content validation. |
| `LastEvaluatedAt` | `*time.Time` | `last_evaluated_at` | Last time this expectation was evaluated. |
| `LastSuccessAt` | `*time.Time` | `last_success_at` | Last time this expectation closed successfully. |
| `LastFailureAt` | `*time.Time` | `last_failure_at` | Last time this expectation closed with a failure result. |
| `LastResult` | `string` | `last_result` | Most recent terminal result for this expectation. |
| `CreatedAt` | `*time.Time` | `created_at` | Creation time. |
| `UpdatedAt` | `*time.Time` | `updated_at` | Last update time. |
