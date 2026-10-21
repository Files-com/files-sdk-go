# Automation

`files_sdk.Automation` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/automation`, named `automation`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/automation"
)
```

## Operations

### List

`GET /automations`

List Automations

```go
func (c *Client) List(params files_sdk.AutomationListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.AutomationListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.Automation` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.AutomationListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id`, `name`, `automation`, `last_modified_at` or `disabled`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `disabled`, `last_modified_at`, `workspace_id` or `automation`. Valid field combinations are `[ disabled, last_modified_at ]`, `[ workspace_id, disabled ]`, `[ disabled, automation ]`, `[ workspace_id, last_modified_at ]`, `[ automation, last_modified_at ]`, `[ workspace_id, automation ]`, `[ workspace_id, disabled, last_modified_at ]`, `[ disabled, automation, last_modified_at ]`, `[ workspace_id, disabled, automation ]`, `[ workspace_id, automation, last_modified_at ]` or `[ workspace_id, disabled, automation, last_modified_at ]`. |
| `FilterGt` | `interface{}` | `filter_gt` | no | If set, return records where the specified field is greater than the supplied value. Valid fields are `last_modified_at`. |
| `FilterGteq` | `interface{}` | `filter_gteq` | no | If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `last_modified_at`. |
| `FilterLt` | `interface{}` | `filter_lt` | no | If set, return records where the specified field is less than the supplied value. Valid fields are `last_modified_at`. |
| `FilterLteq` | `interface{}` | `filter_lteq` | no | If set, return records where the specified field is less than or equal the supplied value. Valid fields are `last_modified_at`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := automation.List(files_sdk.AutomationListParams{})
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

`GET /automations/{id}`

Show Automation

```go
func (c *Client) Find(params files_sdk.AutomationFindParams, opts ...files_sdk.RequestResponseOption) (automation files_sdk.Automation, err error)
func Find(params files_sdk.AutomationFindParams, opts ...files_sdk.RequestResponseOption) (automation files_sdk.Automation, err error)
```

Parameters, `files_sdk.AutomationFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Automation ID. |

Example:

```go
result, err := automation.Find(files_sdk.AutomationFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### GetAuthoringSchema

`GET /automations/authoring_schema`

Show the Automation v2 authoring schema and active node catalog

```go
func (c *Client) GetAuthoringSchema(opts ...files_sdk.RequestResponseOption) (automationAuthoringSchema files_sdk.AutomationAuthoringSchema, err error)
func GetAuthoringSchema(opts ...files_sdk.RequestResponseOption) (automationAuthoringSchema files_sdk.AutomationAuthoringSchema, err error)
```

Example:

```go
result, err := automation.GetAuthoringSchema()
if err != nil {
	return err
}
fmt.Println(result)
```

### Create

`POST /automations`

Create Automation

```go
func (c *Client) Create(params files_sdk.AutomationCreateParams, opts ...files_sdk.RequestResponseOption) (automation files_sdk.Automation, err error)
func Create(params files_sdk.AutomationCreateParams, opts ...files_sdk.RequestResponseOption) (automation files_sdk.Automation, err error)
```

Parameters, `files_sdk.AutomationCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Source` | `string` | `source` | no | Source path/glob. See Automation docs for exact description, but this is used to filter for files in the `path` to find files to operate on. Supports globs, except on remote mounts. |
| `Destinations` | `[]string` | `destinations` | no | A list of destination paths. Use a trailing slash for folder destinations and omit it for file destinations. |
| `DestinationReplaceFrom` | `string` | `destination_replace_from` | no | If set, this string in the destination path will be replaced with the value in `destination_replace_to`. |
| `DestinationReplaceTo` | `string` | `destination_replace_to` | no | If set, this string will replace the value `destination_replace_from` in the destination filename. You can use special patterns here. |
| `Interval` | `string` | `interval` | no | How often to run this automation? One of: `day`, `week`, `week_end`, `month`, `month_end`, `quarter`, `quarter_end`, `year`, `year_end` |
| `Path` | `string` | `path` | no | Path on which this Automation runs. Supports globs, except on remote mounts. |
| `SyncIds` | `string` | `sync_ids` | no | A list of sync IDs the automation is associated with. If sent as a string, it should be comma-delimited. |
| `UserIds` | `string` | `user_ids` | no | A list of user IDs the automation is associated with. If sent as a string, it should be comma-delimited. |
| `GroupIds` | `string` | `group_ids` | no | A list of group IDs the automation is associated with. If sent as a string, it should be comma-delimited. |
| `ScheduleId` | `int64` | `schedule_id` | no | If trigger is `custom_schedule`, the reusable Schedule used instead of the automation's schedule fields. |
| `ScheduleDaysOfWeek` | `[]int64` | `schedule_days_of_week` | no | If trigger is `custom_schedule` without `schedule_id`, a list of days of the week to run this automation. 0 is Sunday, 1 is Monday, etc. |
| `ScheduleTimesOfDay` | `[]string` | `schedule_times_of_day` | no | Times of day to run in HH:MM format (24-hour). Required for `custom_schedule` without `schedule_id`. Optional for `daily` triggers; if not set, runs at midnight UTC. |
| `ScheduleTimeZone` | `string` | `schedule_time_zone` | no | Time zone for schedule fields. Optional for `custom_schedule` without `schedule_id` and for `daily`. If not set, times are interpreted as UTC. |
| `HolidayRegion` | `string` | `holiday_region` | no | Skip automation on holidays in this region. Optional for `custom_schedule` without `schedule_id` and for `daily`. |
| `AlwaysOverwriteSizeMatchingFiles` | `*bool` | `always_overwrite_size_matching_files` | no | Ordinarily, files with identical size in the source and destination will be skipped from copy operations to prevent wasted transfer. If this flag is `true` we will overwrite the destination file always. Note that this may cause large amounts of wasted transfer usage. This setting has no effect unless `overwrite_files` is also set to `true`. |
| `AlwaysSerializeJobs` | `*bool` | `always_serialize_jobs` | no | Ordinarily, we will allow automation runs to run in parallel for non-scheduled automations. If this flag is `true` we will force automation runs to be serialized (run one at a time, one after another). This can resolve some issues with race conditions on remote systems at the cost of some performance. |
| `Description` | `string` | `description` | no | Description for the this Automation. |
| `Definition` | `interface{}` | `definition` | no | Automation v2 graph definition. |
| `Disabled` | `*bool` | `disabled` | no | If true, this automation will not run. |
| `ExcludePattern` | `string` | `exclude_pattern` | no | If set, this glob pattern will exclude files from the automation. Supports globs, except on remote mounts. |
| `ImportUrls` | `[]map[string]interface{}` | `import_urls` | no | List of URLs to be imported and names to be used. |
| `FlattenDestinationStructure` | `*bool` | `flatten_destination_structure` | no | Normally copy and move automations that use globs will implicitly preserve the source folder structure in the destination. If this flag is `true`, the source folder structure will be flattened in the destination. This is useful for copying or moving files from multiple folders into a single destination folder. |
| `IgnoreLockedFolders` | `*bool` | `ignore_locked_folders` | no | If true, the Lock Folders behavior will be disregarded for automated actions. |
| `LegacyFolderMatching` | `*bool` | `legacy_folder_matching` | no | DEPRECATED: If `true`, use the legacy behavior for this automation, where it can operate on folders in addition to just files. This behavior no longer works and should not be used. |
| `Name` | `string` | `name` | no | Name for this automation. |
| `OverwriteFiles` | `*bool` | `overwrite_files` | no | If true, existing files will be overwritten with new files on Move/Copy automations. Note: by default files will not be overwritten on Copy automations if they appear to be the same file size as the newly incoming file. Use the `always_overwrite_size_matching_files` option in conjunction with `overwrite_files` to override this behavior and overwrite files no matter what. |
| `PathTimeZone` | `string` | `path_time_zone` | no | Timezone to use when rendering timestamps in paths. |
| `RetryOnFailureIntervalInMinutes` | `int64` | `retry_on_failure_interval_in_minutes` | no | If the Automation fails, retry at this interval (in minutes). Acceptable values are 5 through 1440 (one day). Set to null to disable. |
| `RetryOnFailureNumberOfAttempts` | `int64` | `retry_on_failure_number_of_attempts` | no | If the Automation fails, retry at most this many times. Maximum allowed value: 10. Set to null to disable. |
| `Trigger` | `files_sdk.AutomationTriggerEnum` | `trigger` | no | How this automation is triggered to run. Values: `manual`, `daily`, `custom_schedule`, `webhook`, `email`, `action`. |
| `TriggerActions` | `[]string` | `trigger_actions` | no | If trigger is `action`, this is the list of action types on which to trigger the automation. Valid actions are create, copy, move, archived_delete, update, read, destroy |
| `Value` | `interface{}` | `value` | no | A Hash of attributes specific to the automation type. |
| `RecurringDay` | `int64` | `recurring_day` | no | If trigger type is `daily`, this specifies a day number to run in one of the supported intervals: `week`, `month`, `quarter`, `year`. |
| `RecurringDays` | `[]int64` | `recurring_days` | no | If trigger type is `daily`, this specifies one or more day numbers to run in one of the supported intervals: `week`, `month`, `quarter`, `year`. |
| `Automation` | `files_sdk.AutomationEnum` | `automation` | yes | Automation type Values: `create_folder`, `delete_file`, `copy_file`, `move_file`, `as2_send`, `run_sync`, `import_file`, `v2`. |
| `WorkspaceId` | `int64` | `workspace_id` | no | Workspace ID |

Example:

```go
result, err := automation.Create(files_sdk.AutomationCreateParams{
	Automation: "create_folder",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Upgrade

`POST /automations/{id}/upgrade`

Upgrade a legacy Automation to Automation v2

```go
func (c *Client) Upgrade(params files_sdk.AutomationUpgradeParams, opts ...files_sdk.RequestResponseOption) (automation files_sdk.Automation, err error)
func Upgrade(params files_sdk.AutomationUpgradeParams, opts ...files_sdk.RequestResponseOption) (automation files_sdk.Automation, err error)
```

Parameters, `files_sdk.AutomationUpgradeParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Automation ID. |

Example:

```go
result, err := automation.Upgrade(files_sdk.AutomationUpgradeParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### ManualRun

`POST /automations/{id}/manual_run`

Manually Run Automation

```go
func (c *Client) ManualRun(params files_sdk.AutomationManualRunParams, opts ...files_sdk.RequestResponseOption) (err error)
func ManualRun(params files_sdk.AutomationManualRunParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.AutomationManualRunParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Automation ID. |
| `Items` | `[]map[string]interface{}` | `items` | no | Initial items for a v2 manual trigger. Each item contains exactly one `file` path or `data` object. |

Example:

```go
err := automation.ManualRun(files_sdk.AutomationManualRunParams{
	Id: 1,
})
if err != nil {
	return err
}
```

### Update

`PATCH /automations/{id}`

Update Automation

```go
func (c *Client) Update(params files_sdk.AutomationUpdateParams, opts ...files_sdk.RequestResponseOption) (automation files_sdk.Automation, err error)
func Update(params files_sdk.AutomationUpdateParams, opts ...files_sdk.RequestResponseOption) (automation files_sdk.Automation, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (automation files_sdk.Automation, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (automation files_sdk.Automation, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.AutomationUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Automation ID. |
| `Source` | `string` | `source` | no | Source path/glob. See Automation docs for exact description, but this is used to filter for files in the `path` to find files to operate on. Supports globs, except on remote mounts. |
| `Destinations` | `[]string` | `destinations` | no | A list of destination paths. Use a trailing slash for folder destinations and omit it for file destinations. |
| `DestinationReplaceFrom` | `string` | `destination_replace_from` | no | If set, this string in the destination path will be replaced with the value in `destination_replace_to`. |
| `DestinationReplaceTo` | `string` | `destination_replace_to` | no | If set, this string will replace the value `destination_replace_from` in the destination filename. You can use special patterns here. |
| `Interval` | `string` | `interval` | no | How often to run this automation? One of: `day`, `week`, `week_end`, `month`, `month_end`, `quarter`, `quarter_end`, `year`, `year_end` |
| `Path` | `string` | `path` | no | Path on which this Automation runs. Supports globs, except on remote mounts. |
| `SyncIds` | `string` | `sync_ids` | no | A list of sync IDs the automation is associated with. If sent as a string, it should be comma-delimited. |
| `UserIds` | `string` | `user_ids` | no | A list of user IDs the automation is associated with. If sent as a string, it should be comma-delimited. |
| `GroupIds` | `string` | `group_ids` | no | A list of group IDs the automation is associated with. If sent as a string, it should be comma-delimited. |
| `ScheduleId` | `int64` | `schedule_id` | no | If trigger is `custom_schedule`, the reusable Schedule used instead of the automation's schedule fields. |
| `ScheduleDaysOfWeek` | `[]int64` | `schedule_days_of_week` | no | If trigger is `custom_schedule` without `schedule_id`, a list of days of the week to run this automation. 0 is Sunday, 1 is Monday, etc. |
| `ScheduleTimesOfDay` | `[]string` | `schedule_times_of_day` | no | Times of day to run in HH:MM format (24-hour). Required for `custom_schedule` without `schedule_id`. Optional for `daily` triggers; if not set, runs at midnight UTC. |
| `ScheduleTimeZone` | `string` | `schedule_time_zone` | no | Time zone for schedule fields. Optional for `custom_schedule` without `schedule_id` and for `daily`. If not set, times are interpreted as UTC. |
| `HolidayRegion` | `string` | `holiday_region` | no | Skip automation on holidays in this region. Optional for `custom_schedule` without `schedule_id` and for `daily`. |
| `AlwaysOverwriteSizeMatchingFiles` | `*bool` | `always_overwrite_size_matching_files` | no | Ordinarily, files with identical size in the source and destination will be skipped from copy operations to prevent wasted transfer. If this flag is `true` we will overwrite the destination file always. Note that this may cause large amounts of wasted transfer usage. This setting has no effect unless `overwrite_files` is also set to `true`. |
| `AlwaysSerializeJobs` | `*bool` | `always_serialize_jobs` | no | Ordinarily, we will allow automation runs to run in parallel for non-scheduled automations. If this flag is `true` we will force automation runs to be serialized (run one at a time, one after another). This can resolve some issues with race conditions on remote systems at the cost of some performance. |
| `Description` | `string` | `description` | no | Description for the this Automation. |
| `Definition` | `interface{}` | `definition` | no | Automation v2 graph definition. |
| `Disabled` | `*bool` | `disabled` | no | If true, this automation will not run. |
| `ExcludePattern` | `string` | `exclude_pattern` | no | If set, this glob pattern will exclude files from the automation. Supports globs, except on remote mounts. |
| `ImportUrls` | `[]map[string]interface{}` | `import_urls` | no | List of URLs to be imported and names to be used. |
| `FlattenDestinationStructure` | `*bool` | `flatten_destination_structure` | no | Normally copy and move automations that use globs will implicitly preserve the source folder structure in the destination. If this flag is `true`, the source folder structure will be flattened in the destination. This is useful for copying or moving files from multiple folders into a single destination folder. |
| `IgnoreLockedFolders` | `*bool` | `ignore_locked_folders` | no | If true, the Lock Folders behavior will be disregarded for automated actions. |
| `LegacyFolderMatching` | `*bool` | `legacy_folder_matching` | no | DEPRECATED: If `true`, use the legacy behavior for this automation, where it can operate on folders in addition to just files. This behavior no longer works and should not be used. |
| `Name` | `string` | `name` | no | Name for this automation. |
| `OverwriteFiles` | `*bool` | `overwrite_files` | no | If true, existing files will be overwritten with new files on Move/Copy automations. Note: by default files will not be overwritten on Copy automations if they appear to be the same file size as the newly incoming file. Use the `always_overwrite_size_matching_files` option in conjunction with `overwrite_files` to override this behavior and overwrite files no matter what. |
| `PathTimeZone` | `string` | `path_time_zone` | no | Timezone to use when rendering timestamps in paths. |
| `RetryOnFailureIntervalInMinutes` | `int64` | `retry_on_failure_interval_in_minutes` | no | If the Automation fails, retry at this interval (in minutes). Acceptable values are 5 through 1440 (one day). Set to null to disable. |
| `RetryOnFailureNumberOfAttempts` | `int64` | `retry_on_failure_number_of_attempts` | no | If the Automation fails, retry at most this many times. Maximum allowed value: 10. Set to null to disable. |
| `Trigger` | `files_sdk.AutomationTriggerEnum` | `trigger` | no | How this automation is triggered to run. Values: `manual`, `daily`, `custom_schedule`, `webhook`, `email`, `action`. |
| `TriggerActions` | `[]string` | `trigger_actions` | no | If trigger is `action`, this is the list of action types on which to trigger the automation. Valid actions are create, copy, move, archived_delete, update, read, destroy |
| `Value` | `interface{}` | `value` | no | A Hash of attributes specific to the automation type. |
| `RecurringDay` | `int64` | `recurring_day` | no | If trigger type is `daily`, this specifies a day number to run in one of the supported intervals: `week`, `month`, `quarter`, `year`. |
| `RecurringDays` | `[]int64` | `recurring_days` | no | If trigger type is `daily`, this specifies one or more day numbers to run in one of the supported intervals: `week`, `month`, `quarter`, `year`. |
| `Automation` | `files_sdk.AutomationEnum` | `automation` | no | Automation type Values: `create_folder`, `delete_file`, `copy_file`, `move_file`, `as2_send`, `run_sync`, `import_file`, `v2`. |

Example:

```go
result, err := automation.Update(files_sdk.AutomationUpdateParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := automation.UpdateWithMap(map[string]interface{}{
	"id": 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /automations/{id}`

Delete Automation

```go
func (c *Client) Delete(params files_sdk.AutomationDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.AutomationDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.AutomationDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Automation ID. |

Example:

```go
err := automation.Delete(files_sdk.AutomationDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.Automation` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Automation ID |
| `WorkspaceId` | `int64` | `workspace_id` | Workspace ID |
| `AlwaysSerializeJobs` | `*bool` | `always_serialize_jobs` | Ordinarily, we will allow automation runs to run in parallel for non-scheduled automations. If this flag is `true` we will force automation runs to be serialized (run one at a time, one after another). This can resolve some issues with race conditions on remote systems at the cost of some performance. |
| `AlwaysOverwriteSizeMatchingFiles` | `*bool` | `always_overwrite_size_matching_files` | Ordinarily, files with identical size in the source and destination will be skipped from copy operations to prevent wasted transfer. If this flag is `true` we will overwrite the destination file always. Note that this may cause large amounts of wasted transfer usage. This setting has no effect unless `overwrite_files` is also set to `true`. |
| `Automation` | `string` | `automation` | Automation type |
| `Deleted` | `*bool` | `deleted` | Indicates if the automation has been deleted. |
| `Description` | `string` | `description` | Description for the this Automation. |
| `Definition` | `interface{}` | `definition` | Automation v2 graph definition. |
| `DestinationReplaceFrom` | `string` | `destination_replace_from` | If set, this string in the destination path will be replaced with the value in `destination_replace_to`. |
| `DestinationReplaceTo` | `string` | `destination_replace_to` | If set, this string will replace the value `destination_replace_from` in the destination filename. You can use special patterns here. |
| `Destinations` | `[]string` | `destinations` | Destination Paths |
| `Disabled` | `*bool` | `disabled` | If true, this automation will not run. |
| `ExcludePattern` | `string` | `exclude_pattern` | If set, this glob pattern will exclude files from the automation. Supports globs, except on remote mounts. |
| `ImportUrls` | `[]map[string]interface{}` | `import_urls` | List of URLs to be imported and names to be used. |
| `InboundEmailAddress` | `string` | `inbound_email_address` | If trigger is `email`, this is the address that triggers the Automation. |
| `FlattenDestinationStructure` | `*bool` | `flatten_destination_structure` | Normally copy and move automations that use globs will implicitly preserve the source folder structure in the destination. If this flag is `true`, the source folder structure will be flattened in the destination. This is useful for copying or moving files from multiple folders into a single destination folder. |
| `GroupIds` | `[]int64` | `group_ids` | IDs of Groups for the Automation (i.e. who to Request File from) |
| `IgnoreLockedFolders` | `*bool` | `ignore_locked_folders` | If true, the Lock Folders behavior will be disregarded for automated actions. |
| `Interval` | `string` | `interval` | If trigger is `daily`, this specifies how often to run this automation. One of: `day`, `week`, `week_end`, `month`, `month_end`, `quarter`, `quarter_end`, `year`, `year_end` |
| `LastModifiedAt` | `*time.Time` | `last_modified_at` | Time when automation was last modified. Does not change for name or description updates. |
| `LegacyFolderMatching` | `*bool` | `legacy_folder_matching` | If `true`, use the legacy behavior for this automation, where it can operate on folders in addition to just files. This behavior no longer works and should not be used. |
| `Name` | `string` | `name` | Name for this automation. |
| `OverwriteFiles` | `*bool` | `overwrite_files` | If true, existing files will be overwritten with new files on Move/Copy automations. Note: by default files will not be overwritten on Copy automations if they appear to be the same file size as the newly incoming file. Use the `always_overwrite_size_matching_files` option in conjunction with `overwrite_files` to override this behavior and overwrite files no matter what. |
| `Path` | `string` | `path` | Path on which this Automation runs. Supports globs, except on remote mounts. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters. |
| `PathTimeZone` | `string` | `path_time_zone` | Timezone to use when rendering timestamps in paths. |
| `Version` | `int64` | `version` | Current Automation v2 definition version. |
| `RecurringDay` | `int64` | `recurring_day` | If trigger type is `daily`, this specifies a day number to run in one of the supported intervals: `week`, `month`, `quarter`, `year`. |
| `RecurringDays` | `[]int64` | `recurring_days` | If trigger type is `daily`, this specifies one or more day numbers to run in one of the supported intervals: `week`, `month`, `quarter`, `year`. |
| `ScheduleId` | `int64` | `schedule_id` | If trigger is `custom_schedule`, the reusable Schedule used instead of the automation's schedule fields. |
| `RetryOnFailureIntervalInMinutes` | `int64` | `retry_on_failure_interval_in_minutes` | If the Automation fails, retry at this interval (in minutes). Acceptable values are 5 through 1440 (one day). Set to null to disable. |
| `RetryOnFailureNumberOfAttempts` | `int64` | `retry_on_failure_number_of_attempts` | If the Automation fails, retry at most this many times. Maximum allowed value: 10. Set to null to disable. |
| `Schedule` | `interface{}` | `schedule` | If trigger is `custom_schedule`, Custom schedule description for when the automation should be run in json format. |
| `HumanReadableSchedule` | `string` | `human_readable_schedule` | If trigger is `custom_schedule` or `daily` with times, Human readable schedule description for when the automation should be run. |
| `ScheduleDaysOfWeek` | `[]int64` | `schedule_days_of_week` | If trigger is `custom_schedule`, Custom schedule description for when the automation should be run. 0 is Sunday, 1 is Monday, etc. |
| `ScheduleTimesOfDay` | `[]string` | `schedule_times_of_day` | Times of day to run in HH:MM format (24-hour). For `custom_schedule`, run at these times on specified days of week. For `daily`, run at these times on the scheduled interval date. |
| `ScheduleTimeZone` | `string` | `schedule_time_zone` | Time zone for the schedule. If not set, times are interpreted as UTC. |
| `Source` | `string` | `source` | Source path/glob. See Automation docs for exact description, but this is used to filter for files in the `path` to find files to operate on. Supports globs, except on remote mounts. |
| `SyncIds` | `[]int64` | `sync_ids` | IDs of Syncs to run by this Automation. |
| `TriggerActions` | `[]string` | `trigger_actions` | If trigger is `action`, this is the list of action types on which to trigger the automation. Valid actions are create, copy, move, archived_delete, update, read, destroy |
| `Trigger` | `string` | `trigger` | How this automation is triggered to run. |
| `UserId` | `int64` | `user_id` | User ID of the Automation's creator. |
| `UserIds` | `[]int64` | `user_ids` | IDs of Users for the Automation (i.e. who to Request File from) |
| `Value` | `interface{}` | `value` | A Hash of attributes specific to the automation type. |
| `WebhookUrl` | `string` | `webhook_url` | If trigger is `webhook`, this is the URL of the webhook to trigger the Automation. |
| `HolidayRegion` | `string` | `holiday_region` | Skip the automation if there is a formal, observed holiday for this region. |
