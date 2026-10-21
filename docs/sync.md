# Sync

`files_sdk.Sync` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/sync`, named `sync`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/sync"
)
```

## Operations

### List

`GET /syncs`

List Syncs

```go
func (c *Client) List(params files_sdk.SyncListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.SyncListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.Sync` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.SyncListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `site_id`, `workspace_id` or `name`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `workspace_id`, `disabled`, `src_remote_server_id` or `dest_remote_server_id`. Valid field combinations are `[ workspace_id, disabled ]`, `[ workspace_id, src_remote_server_id ]`, `[ workspace_id, dest_remote_server_id ]`, `[ disabled, src_remote_server_id ]`, `[ disabled, dest_remote_server_id ]`, `[ workspace_id, disabled, src_remote_server_id ]` or `[ workspace_id, disabled, dest_remote_server_id ]`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := sync.List(files_sdk.SyncListParams{})
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

`GET /syncs/{id}`

Show Sync

```go
func (c *Client) Find(params files_sdk.SyncFindParams, opts ...files_sdk.RequestResponseOption) (sync files_sdk.Sync, err error)
func Find(params files_sdk.SyncFindParams, opts ...files_sdk.RequestResponseOption) (sync files_sdk.Sync, err error)
```

Parameters, `files_sdk.SyncFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Sync ID. |

Example:

```go
result, err := sync.Find(files_sdk.SyncFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Create

`POST /syncs`

Create Sync

```go
func (c *Client) Create(params files_sdk.SyncCreateParams, opts ...files_sdk.RequestResponseOption) (sync files_sdk.Sync, err error)
func Create(params files_sdk.SyncCreateParams, opts ...files_sdk.RequestResponseOption) (sync files_sdk.Sync, err error)
```

Parameters, `files_sdk.SyncCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `DeleteEmptyFolders` | `*bool` | `delete_empty_folders` | no | Delete empty folders after sync? |
| `Description` | `string` | `description` | no | Description for this sync job |
| `DestPath` | `string` | `dest_path` | no | Absolute destination path for the sync |
| `DestRemoteServerId` | `int64` | `dest_remote_server_id` | no | Remote server ID for the destination (if remote) |
| `Disabled` | `*bool` | `disabled` | no | Is this sync disabled? |
| `ExcludePatterns` | `[]string` | `exclude_patterns` | no | Array of glob patterns to exclude |
| `HolidayRegion` | `string` | `holiday_region` | no | Skip the sync if there is a formal, observed holiday for this region. |
| `IncludePatterns` | `[]string` | `include_patterns` | no | Array of glob patterns to include |
| `Interval` | `string` | `interval` | no | If trigger is `daily`, this specifies how often to run this sync. One of: `day`, `week`, `week_end`, `month`, `month_end`, `quarter`, `quarter_end`, `year`, `year_end` |
| `KeepAfterCopy` | `*bool` | `keep_after_copy` | no | Keep files after copying? |
| `Name` | `string` | `name` | no | Name for this sync job |
| `RecurringDay` | `int64` | `recurring_day` | no | If trigger type is `daily`, this specifies a day number to run in one of the supported intervals: `week`, `month`, `quarter`, `year`. |
| `RecurringDays` | `[]int64` | `recurring_days` | no | If trigger type is `daily`, this specifies one or more day numbers to run in one of the supported intervals: `week`, `month`, `quarter`, `year`. |
| `ScheduleId` | `int64` | `schedule_id` | no | If trigger is `custom_schedule`, the reusable Schedule used instead of the sync's schedule fields. |
| `ScheduleDaysOfWeek` | `[]int64` | `schedule_days_of_week` | no | If trigger is `custom_schedule`, Custom schedule description for when the sync should be run. 0-based days of the week. 0 is Sunday, 1 is Monday, etc. |
| `ScheduleTimeZone` | `string` | `schedule_time_zone` | no | Time zone for the schedule. If not set, times are interpreted as UTC. |
| `ScheduleTimesOfDay` | `[]string` | `schedule_times_of_day` | no | Times of day to run in HH:MM format. For `custom_schedule`, run at these times on specified days of week. For `daily`, run at these times on the scheduled interval date. |
| `SrcPath` | `string` | `src_path` | no | Absolute source path for the sync |
| `SrcRemoteServerId` | `int64` | `src_remote_server_id` | no | Remote server ID for the source (if remote) |
| `SyncIntervalMinutes` | `int64` | `sync_interval_minutes` | no | Frequency in minutes between syncs. If set, this value must be greater than or equal to the `remote_sync_interval` value for the site's plan. If left blank, the plan's `remote_sync_interval` will be used. This setting is only used if `trigger` is empty. |
| `Trigger` | `files_sdk.SyncTriggerEnum` | `trigger` | no | Trigger type: daily, custom_schedule, or manual Values: `daily`, `custom_schedule`, `manual`. |
| `TriggerFile` | `string` | `trigger_file` | no | Some MFT services request an empty file (known as a trigger file) to signal the sync is complete and they can begin further processing. If trigger_file is set, a zero-byte file will be sent at the end of the sync. |
| `AlwaysWriteTriggerFile` | `*bool` | `always_write_trigger_file` | no | If true, the trigger file will be sent at the end of a successful sync even when no files were transferred. |
| `WorkspaceId` | `int64` | `workspace_id` | no | Workspace ID this sync belongs to |

Example:

```go
result, err := sync.Create(files_sdk.SyncCreateParams{})
if err != nil {
	return err
}
fmt.Println(result)
```

### DryRun

`POST /syncs/{id}/dry_run`

Dry Run Sync

```go
func (c *Client) DryRun(params files_sdk.SyncDryRunParams, opts ...files_sdk.RequestResponseOption) (err error)
func DryRun(params files_sdk.SyncDryRunParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.SyncDryRunParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Sync ID. |

Example:

```go
err := sync.DryRun(files_sdk.SyncDryRunParams{
	Id: 1,
})
if err != nil {
	return err
}
```

### ManualRun

`POST /syncs/{id}/manual_run`

Manually Run Sync

```go
func (c *Client) ManualRun(params files_sdk.SyncManualRunParams, opts ...files_sdk.RequestResponseOption) (err error)
func ManualRun(params files_sdk.SyncManualRunParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.SyncManualRunParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Sync ID. |

Example:

```go
err := sync.ManualRun(files_sdk.SyncManualRunParams{
	Id: 1,
})
if err != nil {
	return err
}
```

### Update

`PATCH /syncs/{id}`

Update Sync

```go
func (c *Client) Update(params files_sdk.SyncUpdateParams, opts ...files_sdk.RequestResponseOption) (sync files_sdk.Sync, err error)
func Update(params files_sdk.SyncUpdateParams, opts ...files_sdk.RequestResponseOption) (sync files_sdk.Sync, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (sync files_sdk.Sync, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (sync files_sdk.Sync, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.SyncUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Sync ID. |
| `DeleteEmptyFolders` | `*bool` | `delete_empty_folders` | no | Delete empty folders after sync? |
| `Description` | `string` | `description` | no | Description for this sync job |
| `DestPath` | `string` | `dest_path` | no | Absolute destination path for the sync |
| `DestRemoteServerId` | `int64` | `dest_remote_server_id` | no | Remote server ID for the destination (if remote) |
| `Disabled` | `*bool` | `disabled` | no | Is this sync disabled? |
| `ExcludePatterns` | `[]string` | `exclude_patterns` | no | Array of glob patterns to exclude |
| `HolidayRegion` | `string` | `holiday_region` | no | Skip the sync if there is a formal, observed holiday for this region. |
| `IncludePatterns` | `[]string` | `include_patterns` | no | Array of glob patterns to include |
| `Interval` | `string` | `interval` | no | If trigger is `daily`, this specifies how often to run this sync. One of: `day`, `week`, `week_end`, `month`, `month_end`, `quarter`, `quarter_end`, `year`, `year_end` |
| `KeepAfterCopy` | `*bool` | `keep_after_copy` | no | Keep files after copying? |
| `Name` | `string` | `name` | no | Name for this sync job |
| `RecurringDay` | `int64` | `recurring_day` | no | If trigger type is `daily`, this specifies a day number to run in one of the supported intervals: `week`, `month`, `quarter`, `year`. |
| `RecurringDays` | `[]int64` | `recurring_days` | no | If trigger type is `daily`, this specifies one or more day numbers to run in one of the supported intervals: `week`, `month`, `quarter`, `year`. |
| `ScheduleId` | `int64` | `schedule_id` | no | If trigger is `custom_schedule`, the reusable Schedule used instead of the sync's schedule fields. |
| `ScheduleDaysOfWeek` | `[]int64` | `schedule_days_of_week` | no | If trigger is `custom_schedule`, Custom schedule description for when the sync should be run. 0-based days of the week. 0 is Sunday, 1 is Monday, etc. |
| `ScheduleTimeZone` | `string` | `schedule_time_zone` | no | Time zone for the schedule. If not set, times are interpreted as UTC. |
| `ScheduleTimesOfDay` | `[]string` | `schedule_times_of_day` | no | Times of day to run in HH:MM format. For `custom_schedule`, run at these times on specified days of week. For `daily`, run at these times on the scheduled interval date. |
| `SrcPath` | `string` | `src_path` | no | Absolute source path for the sync |
| `SrcRemoteServerId` | `int64` | `src_remote_server_id` | no | Remote server ID for the source (if remote) |
| `SyncIntervalMinutes` | `int64` | `sync_interval_minutes` | no | Frequency in minutes between syncs. If set, this value must be greater than or equal to the `remote_sync_interval` value for the site's plan. If left blank, the plan's `remote_sync_interval` will be used. This setting is only used if `trigger` is empty. |
| `Trigger` | `files_sdk.SyncTriggerEnum` | `trigger` | no | Trigger type: daily, custom_schedule, or manual Values: `daily`, `custom_schedule`, `manual`. |
| `TriggerFile` | `string` | `trigger_file` | no | Some MFT services request an empty file (known as a trigger file) to signal the sync is complete and they can begin further processing. If trigger_file is set, a zero-byte file will be sent at the end of the sync. |
| `AlwaysWriteTriggerFile` | `*bool` | `always_write_trigger_file` | no | If true, the trigger file will be sent at the end of a successful sync even when no files were transferred. |

Example:

```go
result, err := sync.Update(files_sdk.SyncUpdateParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := sync.UpdateWithMap(map[string]interface{}{
	"id": 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /syncs/{id}`

Delete Sync

```go
func (c *Client) Delete(params files_sdk.SyncDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.SyncDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.SyncDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Sync ID. |

Example:

```go
err := sync.Delete(files_sdk.SyncDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.Sync` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Sync ID |
| `Name` | `string` | `name` | Name for this sync job |
| `Description` | `string` | `description` | Description for this sync job |
| `SiteId` | `int64` | `site_id` | Site ID this sync belongs to |
| `WorkspaceId` | `int64` | `workspace_id` | Workspace ID this sync belongs to |
| `UserId` | `int64` | `user_id` | User who created or owns this sync |
| `SrcPath` | `string` | `src_path` | Absolute source path for the sync |
| `DestPath` | `string` | `dest_path` | Absolute destination path for the sync |
| `SrcRemoteServerId` | `int64` | `src_remote_server_id` | Remote server ID for the source (if remote) |
| `DestRemoteServerId` | `int64` | `dest_remote_server_id` | Remote server ID for the destination (if remote) |
| `SrcSiteId` | `int64` | `src_site_id` | Source site ID if syncing from a child or partner site |
| `DestSiteId` | `int64` | `dest_site_id` | Destination site ID if syncing to a child or partner site |
| `TwoWay` | `*bool` | `two_way` | Is this a two-way sync? |
| `KeepAfterCopy` | `*bool` | `keep_after_copy` | Keep files after copying? |
| `DeleteEmptyFolders` | `*bool` | `delete_empty_folders` | Delete empty folders after sync? |
| `Disabled` | `*bool` | `disabled` | Is this sync disabled? |
| `Trigger` | `string` | `trigger` | Trigger type: daily, custom_schedule, or manual |
| `TriggerFile` | `string` | `trigger_file` | Some MFT services request an empty file (known as a trigger file) to signal the sync is complete and they can begin further processing. If trigger_file is set, a zero-byte file will be sent at the end of the sync. |
| `AlwaysWriteTriggerFile` | `*bool` | `always_write_trigger_file` | If true, the trigger file will be sent at the end of a successful sync even when no files were transferred. |
| `IncludePatterns` | `[]string` | `include_patterns` | Array of glob patterns to include |
| `ExcludePatterns` | `[]string` | `exclude_patterns` | Array of glob patterns to exclude |
| `CreatedAt` | `*time.Time` | `created_at` | When this sync was created |
| `UpdatedAt` | `*time.Time` | `updated_at` | When this sync was last updated |
| `SyncIntervalMinutes` | `int64` | `sync_interval_minutes` | Frequency in minutes between syncs. If set, this value must be greater than or equal to the `remote_sync_interval` value for the site's plan. If left blank, the plan's `remote_sync_interval` will be used. This setting is only used if `trigger` is empty. |
| `Interval` | `string` | `interval` | If trigger is `daily`, this specifies how often to run this sync. One of: `day`, `week`, `week_end`, `month`, `month_end`, `quarter`, `quarter_end`, `year`, `year_end` |
| `RecurringDay` | `int64` | `recurring_day` | If trigger type is `daily`, this specifies a day number to run in one of the supported intervals: `week`, `month`, `quarter`, `year`. |
| `RecurringDays` | `[]int64` | `recurring_days` | If trigger type is `daily`, this specifies one or more day numbers to run in one of the supported intervals: `week`, `month`, `quarter`, `year`. |
| `ScheduleId` | `int64` | `schedule_id` | If trigger is `custom_schedule`, the reusable Schedule used instead of the sync's schedule fields. |
| `ScheduleDaysOfWeek` | `[]int64` | `schedule_days_of_week` | If trigger is `custom_schedule`, Custom schedule description for when the sync should be run. 0-based days of the week. 0 is Sunday, 1 is Monday, etc. |
| `ScheduleTimesOfDay` | `[]string` | `schedule_times_of_day` | Times of day to run in HH:MM format. For `custom_schedule`, run at these times on specified days of week. For `daily`, run at these times on the scheduled interval date. |
| `ScheduleTimeZone` | `string` | `schedule_time_zone` | Time zone for the schedule. If not set, times are interpreted as UTC. |
| `HolidayRegion` | `string` | `holiday_region` | Skip the sync if there is a formal, observed holiday for this region. |
| `LatestSyncRun` | `files_sdk.SyncRun` | `latest_sync_run` | The latest run of this sync |
