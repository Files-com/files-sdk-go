# AiTask

`files_sdk.AiTask` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/aitask`, named `ai_task`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	ai_task "github.com/Files-com/files-sdk-go/v3/aitask"
)
```

## Operations

### List

`GET /ai_tasks`

List Ai Tasks

```go
func (c *Client) List(params files_sdk.AiTaskListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.AiTaskListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.AiTask` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.AiTaskListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id`, `name`, `disabled` or `updated_at`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `disabled`, `trigger` or `workspace_id`. Valid field combinations are `[ workspace_id, disabled ]`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := ai_task.List(files_sdk.AiTaskListParams{})
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

`GET /ai_tasks/{id}`

Show Ai Task

```go
func (c *Client) Find(params files_sdk.AiTaskFindParams, opts ...files_sdk.RequestResponseOption) (aiTask files_sdk.AiTask, err error)
func Find(params files_sdk.AiTaskFindParams, opts ...files_sdk.RequestResponseOption) (aiTask files_sdk.AiTask, err error)
```

Parameters, `files_sdk.AiTaskFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Ai Task ID. |

Example:

```go
result, err := ai_task.Find(files_sdk.AiTaskFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Create

`POST /ai_tasks`

Create Ai Task

```go
func (c *Client) Create(params files_sdk.AiTaskCreateParams, opts ...files_sdk.RequestResponseOption) (aiTask files_sdk.AiTask, err error)
func Create(params files_sdk.AiTaskCreateParams, opts ...files_sdk.RequestResponseOption) (aiTask files_sdk.AiTask, err error)
```

Parameters, `files_sdk.AiTaskCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Description` | `string` | `description` | no | AI Task description. |
| `Disabled` | `*bool` | `disabled` | no | If true, this AI Task will not run. |
| `HolidayRegion` | `string` | `holiday_region` | no | Optional holiday region used by the AI Task schedule. |
| `Interval` | `string` | `interval` | no | If trigger is `daily`, this specifies how often to run the AI Task. |
| `Name` | `string` | `name` | yes | AI Task name. |
| `Path` | `string` | `path` | no | Path scope used for action-triggered AI Tasks. |
| `PermissionSet` | `files_sdk.AiTaskPermissionSetEnum` | `permission_set` | no | Permissions used by the internal API key for this AI Task. Valid values are `full` and `files_only`. Values: `full`, `files_only`. |
| `Prompt` | `string` | `prompt` | yes | Prompt sent when this AI Task is invoked. |
| `RecurringDay` | `int64` | `recurring_day` | no | If trigger is `daily`, this selects the day number inside the chosen interval. |
| `RecurringDays` | `[]int64` | `recurring_days` | no | If trigger is `daily`, this selects one or more day numbers inside a `week`, `month`, `quarter`, or `year` interval. |
| `ScheduleId` | `int64` | `schedule_id` | no | If trigger is `custom_schedule`, the reusable Schedule used instead of the AI Task's schedule fields. |
| `ScheduleDaysOfWeek` | `[]int64` | `schedule_days_of_week` | no | If trigger is `custom_schedule`, the 0-based weekdays used by the schedule. |
| `ScheduleTimeZone` | `string` | `schedule_time_zone` | no | Time zone used by the AI Task schedule. |
| `ScheduleTimesOfDay` | `[]string` | `schedule_times_of_day` | no | Times of day in HH:MM format for the AI Task schedule. |
| `Source` | `string` | `source` | no | Source glob used with `path` for action-triggered AI Tasks. |
| `Trigger` | `files_sdk.AiTaskTriggerEnum` | `trigger` | no | How this AI Task is triggered. Values: `manual`, `daily`, `custom_schedule`, `action`. |
| `TriggerActions` | `[]string` | `trigger_actions` | no | If trigger is `action`, the file action types that invoke this AI Task. Valid actions are create, copy, move, archived_delete, update, read, destroy. |
| `WorkspaceId` | `int64` | `workspace_id` | no | Workspace ID. `0` means the default workspace. |

Example:

```go
result, err := ai_task.Create(files_sdk.AiTaskCreateParams{
	Name:   "Summarize daily reports",
	Prompt: "Summarize the uploaded file and identify follow-up actions.",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### ManualRun

`POST /ai_tasks/{id}/manual_run`

Manually Run AI Task

```go
func (c *Client) ManualRun(params files_sdk.AiTaskManualRunParams, opts ...files_sdk.RequestResponseOption) (err error)
func ManualRun(params files_sdk.AiTaskManualRunParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.AiTaskManualRunParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Ai Task ID. |

Example:

```go
err := ai_task.ManualRun(files_sdk.AiTaskManualRunParams{
	Id: 1,
})
if err != nil {
	return err
}
```

### Update

`PATCH /ai_tasks/{id}`

Update Ai Task

```go
func (c *Client) Update(params files_sdk.AiTaskUpdateParams, opts ...files_sdk.RequestResponseOption) (aiTask files_sdk.AiTask, err error)
func Update(params files_sdk.AiTaskUpdateParams, opts ...files_sdk.RequestResponseOption) (aiTask files_sdk.AiTask, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (aiTask files_sdk.AiTask, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (aiTask files_sdk.AiTask, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.AiTaskUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Ai Task ID. |
| `Description` | `string` | `description` | no | AI Task description. |
| `Disabled` | `*bool` | `disabled` | no | If true, this AI Task will not run. |
| `HolidayRegion` | `string` | `holiday_region` | no | Optional holiday region used by the AI Task schedule. |
| `Interval` | `string` | `interval` | no | If trigger is `daily`, this specifies how often to run the AI Task. |
| `Name` | `string` | `name` | no | AI Task name. |
| `Path` | `string` | `path` | no | Path scope used for action-triggered AI Tasks. |
| `PermissionSet` | `files_sdk.AiTaskPermissionSetEnum` | `permission_set` | no | Permissions used by the internal API key for this AI Task. Valid values are `full` and `files_only`. Values: `full`, `files_only`. |
| `Prompt` | `string` | `prompt` | no | Prompt sent when this AI Task is invoked. |
| `RecurringDay` | `int64` | `recurring_day` | no | If trigger is `daily`, this selects the day number inside the chosen interval. |
| `RecurringDays` | `[]int64` | `recurring_days` | no | If trigger is `daily`, this selects one or more day numbers inside a `week`, `month`, `quarter`, or `year` interval. |
| `ScheduleId` | `int64` | `schedule_id` | no | If trigger is `custom_schedule`, the reusable Schedule used instead of the AI Task's schedule fields. |
| `ScheduleDaysOfWeek` | `[]int64` | `schedule_days_of_week` | no | If trigger is `custom_schedule`, the 0-based weekdays used by the schedule. |
| `ScheduleTimeZone` | `string` | `schedule_time_zone` | no | Time zone used by the AI Task schedule. |
| `ScheduleTimesOfDay` | `[]string` | `schedule_times_of_day` | no | Times of day in HH:MM format for the AI Task schedule. |
| `Source` | `string` | `source` | no | Source glob used with `path` for action-triggered AI Tasks. |
| `Trigger` | `files_sdk.AiTaskTriggerEnum` | `trigger` | no | How this AI Task is triggered. Values: `manual`, `daily`, `custom_schedule`, `action`. |
| `TriggerActions` | `[]string` | `trigger_actions` | no | If trigger is `action`, the file action types that invoke this AI Task. Valid actions are create, copy, move, archived_delete, update, read, destroy. |
| `WorkspaceId` | `int64` | `workspace_id` | no | Workspace ID. `0` means the default workspace. |

Example:

```go
result, err := ai_task.Update(files_sdk.AiTaskUpdateParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := ai_task.UpdateWithMap(map[string]interface{}{
	"id": 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /ai_tasks/{id}`

Delete Ai Task

```go
func (c *Client) Delete(params files_sdk.AiTaskDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.AiTaskDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.AiTaskDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Ai Task ID. |

Example:

```go
err := ai_task.Delete(files_sdk.AiTaskDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.AiTask` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | AI Task ID. |
| `WorkspaceId` | `int64` | `workspace_id` | Workspace ID. `0` means the default workspace. |
| `Name` | `string` | `name` | AI Task name. |
| `Description` | `string` | `description` | AI Task description. |
| `Prompt` | `string` | `prompt` | Prompt sent when this AI Task is invoked. |
| `PermissionSet` | `string` | `permission_set` | Permissions used by the internal API key for this AI Task. Valid values are `full` and `files_only`. |
| `Path` | `string` | `path` | Path scope used for action-triggered AI Tasks. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters. |
| `Source` | `string` | `source` | Source glob used with `path` for action-triggered AI Tasks. |
| `Disabled` | `*bool` | `disabled` | If true, this AI Task will not run. |
| `Trigger` | `string` | `trigger` | How this AI Task is triggered. |
| `TriggerActions` | `[]string` | `trigger_actions` | If trigger is `action`, the file action types that invoke this AI Task. Valid actions are create, copy, move, archived_delete, update, read, destroy. |
| `Interval` | `string` | `interval` | If trigger is `daily`, this specifies how often to run the AI Task. |
| `RecurringDay` | `int64` | `recurring_day` | If trigger is `daily`, this selects the day number inside the chosen interval. |
| `RecurringDays` | `[]int64` | `recurring_days` | If trigger is `daily`, this selects one or more day numbers inside a `week`, `month`, `quarter`, or `year` interval. |
| `ScheduleId` | `int64` | `schedule_id` | If trigger is `custom_schedule`, the reusable Schedule used instead of the AI Task's schedule fields. |
| `ScheduleDaysOfWeek` | `[]int64` | `schedule_days_of_week` | If trigger is `custom_schedule`, the 0-based weekdays used by the schedule. |
| `ScheduleTimesOfDay` | `[]string` | `schedule_times_of_day` | Times of day in HH:MM format for the AI Task schedule. |
| `ScheduleTimeZone` | `string` | `schedule_time_zone` | Time zone used by the AI Task schedule. |
| `HolidayRegion` | `string` | `holiday_region` | Optional holiday region used by the AI Task schedule. |
| `HumanReadableSchedule` | `string` | `human_readable_schedule` | Human-readable schedule description. |
| `LastRunAt` | `*time.Time` | `last_run_at` | Most recent successful invocation time. |
| `MasterAdminUserId` | `int64` | `master_admin_user_id` | Master User ID used for AI Task invocations. |
| `CreatedAt` | `*time.Time` | `created_at` | Creation time. |
| `UpdatedAt` | `*time.Time` | `updated_at` | Last update time. |
