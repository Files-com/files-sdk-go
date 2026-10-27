# ScheduledExport

`files_sdk.ScheduledExport` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/scheduledexport`, named `scheduled_export`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	scheduled_export "github.com/Files-com/files-sdk-go/v3/scheduledexport"
)
```

## Operations

### List

`GET /scheduled_exports`

List Scheduled Exports

```go
func (c *Client) List(params files_sdk.ScheduledExportListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.ScheduledExportListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.ScheduledExport` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.ScheduledExportListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `name`, `export_type` or `disabled`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `disabled` and `export_type`. |
| `FilterPrefix` | `interface{}` | `filter_prefix` | no | If set, return records where the specified field is prefixed by the supplied value. Valid fields are `export_type`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := scheduled_export.List(files_sdk.ScheduledExportListParams{})
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

`GET /scheduled_exports/{id}`

Show Scheduled Export

```go
func (c *Client) Find(params files_sdk.ScheduledExportFindParams, opts ...files_sdk.RequestResponseOption) (scheduledExport files_sdk.ScheduledExport, err error)
func Find(params files_sdk.ScheduledExportFindParams, opts ...files_sdk.RequestResponseOption) (scheduledExport files_sdk.ScheduledExport, err error)
```

Parameters, `files_sdk.ScheduledExportFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Scheduled Export ID. |

Example:

```go
result, err := scheduled_export.Find(files_sdk.ScheduledExportFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Create

`POST /scheduled_exports`

Create Scheduled Export

```go
func (c *Client) Create(params files_sdk.ScheduledExportCreateParams, opts ...files_sdk.RequestResponseOption) (scheduledExport files_sdk.ScheduledExport, err error)
func Create(params files_sdk.ScheduledExportCreateParams, opts ...files_sdk.RequestResponseOption) (scheduledExport files_sdk.ScheduledExport, err error)
```

Parameters, `files_sdk.ScheduledExportCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Name` | `string` | `name` | yes | Name for this scheduled export. |
| `ExportType` | `string` | `export_type` | yes | Export report type. Valid values: folder_size_audit, group_membership_audit, permission_audit, share_link_audit |
| `ExportOptions` | `interface{}` | `export_options` | no | Options for the report selected by export_type. |
| `UserId` | `int64` | `user_id` | no | Site Admin user who receives the completed export e-mail. |
| `Disabled` | `*bool` | `disabled` | no | If true, this scheduled export will not run. |
| `Trigger` | `files_sdk.ScheduledExportTriggerEnum` | `trigger` | no | Schedule trigger type: `daily` or `custom_schedule`. Values: `daily`, `custom_schedule`. |
| `Interval` | `string` | `interval` | no | If trigger is `daily`, this specifies how often to run the scheduled export. |
| `RecurringDay` | `int64` | `recurring_day` | no | If trigger is `daily`, this selects the day number inside the chosen interval. |
| `RecurringDays` | `[]int64` | `recurring_days` | no | If trigger is `daily`, this selects one or more day numbers inside a `week`, `month`, `quarter`, or `year` interval. |
| `ScheduleId` | `int64` | `schedule_id` | no | If trigger is `custom_schedule`, the reusable Schedule used instead of the scheduled export's schedule fields. |
| `ScheduleDaysOfWeek` | `[]int64` | `schedule_days_of_week` | no | If trigger is `custom_schedule`, the 0-based weekdays used by the schedule. |
| `ScheduleTimesOfDay` | `[]string` | `schedule_times_of_day` | no | Times of day in HH:MM format for the scheduled export schedule. |
| `ScheduleTimeZone` | `string` | `schedule_time_zone` | no | Time zone used by the scheduled export schedule. |
| `HolidayRegion` | `string` | `holiday_region` | no | Optional holiday region used by the scheduled export schedule. |

Example:

```go
result, err := scheduled_export.Create(files_sdk.ScheduledExportCreateParams{
	Name:       "Monthly access review",
	ExportType: "permission_audit",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Update

`PATCH /scheduled_exports/{id}`

Update Scheduled Export

```go
func (c *Client) Update(params files_sdk.ScheduledExportUpdateParams, opts ...files_sdk.RequestResponseOption) (scheduledExport files_sdk.ScheduledExport, err error)
func Update(params files_sdk.ScheduledExportUpdateParams, opts ...files_sdk.RequestResponseOption) (scheduledExport files_sdk.ScheduledExport, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (scheduledExport files_sdk.ScheduledExport, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (scheduledExport files_sdk.ScheduledExport, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.ScheduledExportUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Scheduled Export ID. |
| `Name` | `string` | `name` | no | Name for this scheduled export. |
| `ExportType` | `string` | `export_type` | no | Export report type. Valid values: folder_size_audit, group_membership_audit, permission_audit, share_link_audit |
| `ExportOptions` | `interface{}` | `export_options` | no | Options for the report selected by export_type. |
| `UserId` | `int64` | `user_id` | no | Site Admin user who receives the completed export e-mail. |
| `Disabled` | `*bool` | `disabled` | no | If true, this scheduled export will not run. |
| `Trigger` | `files_sdk.ScheduledExportTriggerEnum` | `trigger` | no | Schedule trigger type: `daily` or `custom_schedule`. Values: `daily`, `custom_schedule`. |
| `Interval` | `string` | `interval` | no | If trigger is `daily`, this specifies how often to run the scheduled export. |
| `RecurringDay` | `int64` | `recurring_day` | no | If trigger is `daily`, this selects the day number inside the chosen interval. |
| `RecurringDays` | `[]int64` | `recurring_days` | no | If trigger is `daily`, this selects one or more day numbers inside a `week`, `month`, `quarter`, or `year` interval. |
| `ScheduleId` | `int64` | `schedule_id` | no | If trigger is `custom_schedule`, the reusable Schedule used instead of the scheduled export's schedule fields. |
| `ScheduleDaysOfWeek` | `[]int64` | `schedule_days_of_week` | no | If trigger is `custom_schedule`, the 0-based weekdays used by the schedule. |
| `ScheduleTimesOfDay` | `[]string` | `schedule_times_of_day` | no | Times of day in HH:MM format for the scheduled export schedule. |
| `ScheduleTimeZone` | `string` | `schedule_time_zone` | no | Time zone used by the scheduled export schedule. |
| `HolidayRegion` | `string` | `holiday_region` | no | Optional holiday region used by the scheduled export schedule. |

Example:

```go
result, err := scheduled_export.Update(files_sdk.ScheduledExportUpdateParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := scheduled_export.UpdateWithMap(map[string]interface{}{
	"id": 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /scheduled_exports/{id}`

Delete Scheduled Export

```go
func (c *Client) Delete(params files_sdk.ScheduledExportDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.ScheduledExportDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.ScheduledExportDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Scheduled Export ID. |

Example:

```go
err := scheduled_export.Delete(files_sdk.ScheduledExportDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.ScheduledExport` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Scheduled Export ID |
| `Name` | `string` | `name` | Name for this scheduled export. |
| `ExportType` | `string` | `export_type` | Export report type. Valid values: folder_size_audit, group_membership_audit, permission_audit, share_link_audit |
| `ReportName` | `string` | `report_name` | Human-readable report name. |
| `ExportOptions` | `interface{}` | `export_options` | Options for the report selected by export_type. |
| `UserId` | `int64` | `user_id` | Site Admin user who receives the completed export e-mail. |
| `Disabled` | `*bool` | `disabled` | If true, this scheduled export will not run. |
| `Trigger` | `string` | `trigger` | Schedule trigger type: `daily` or `custom_schedule`. |
| `Interval` | `string` | `interval` | If trigger is `daily`, this specifies how often to run the scheduled export. |
| `RecurringDay` | `int64` | `recurring_day` | If trigger is `daily`, this selects the day number inside the chosen interval. |
| `RecurringDays` | `[]int64` | `recurring_days` | If trigger is `daily`, this selects one or more day numbers inside a `week`, `month`, `quarter`, or `year` interval. |
| `ScheduleId` | `int64` | `schedule_id` | If trigger is `custom_schedule`, the reusable Schedule used instead of the scheduled export's schedule fields. |
| `ScheduleDaysOfWeek` | `[]int64` | `schedule_days_of_week` | If trigger is `custom_schedule`, the 0-based weekdays used by the schedule. |
| `ScheduleTimesOfDay` | `[]string` | `schedule_times_of_day` | Times of day in HH:MM format for the scheduled export schedule. |
| `ScheduleTimeZone` | `string` | `schedule_time_zone` | Time zone used by the scheduled export schedule. |
| `HolidayRegion` | `string` | `holiday_region` | Optional holiday region used by the scheduled export schedule. |
| `HumanReadableSchedule` | `string` | `human_readable_schedule` | Human-readable schedule description. |
| `LastRunAt` | `*time.Time` | `last_run_at` | Most recent scheduled attempt time, including attempts that failed validation. |
| `LastExportId` | `int64` | `last_export_id` | Most recent Export ID created by this schedule. |
| `LastError` | `string` | `last_error` | Validation error from the most recent scheduled attempt. The schedule remains enabled and retries at its next scheduled time. Cleared when an export is successfully created; does not describe errors during export generation. |
| `CreatedAt` | `*time.Time` | `created_at` | Creation time. |
| `UpdatedAt` | `*time.Time` | `updated_at` | Last update time. |
