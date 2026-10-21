package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// Automation is a Files.com API resource.
type Automation struct {
	// Automation ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Workspace ID.
	WorkspaceId int64 `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	// Ordinarily, we will allow automation runs to run in parallel for non-scheduled automations. If this flag is `true` we will force automation runs to be serialized (run one at a time, one after another). This can resolve some issues with race conditions on remote systems at the cost of some performance.
	AlwaysSerializeJobs *bool `json:"always_serialize_jobs,omitempty" path:"always_serialize_jobs,omitempty" url:"always_serialize_jobs,omitempty"`
	// Ordinarily, files with identical size in the source and destination will be skipped from copy operations to prevent wasted transfer. If this flag is `true` we will overwrite the destination file always. Note that this may cause large amounts of wasted transfer usage. This setting has no effect unless `overwrite_files` is also set to `true`.
	AlwaysOverwriteSizeMatchingFiles *bool `json:"always_overwrite_size_matching_files,omitempty" path:"always_overwrite_size_matching_files,omitempty" url:"always_overwrite_size_matching_files,omitempty"`
	// Automation type.
	Automation string `json:"automation,omitempty" path:"automation,omitempty" url:"automation,omitempty"`
	// Indicates if the automation has been deleted.
	Deleted *bool `json:"deleted,omitempty" path:"deleted,omitempty" url:"deleted,omitempty"`
	// Description for the this Automation.
	Description string `json:"description,omitempty" path:"description,omitempty" url:"description,omitempty"`
	// Automation v2 graph definition.
	Definition interface{} `json:"definition,omitempty" path:"definition,omitempty" url:"definition,omitempty"`
	// If set, this string in the destination path will be replaced with the value in `destination_replace_to`.
	DestinationReplaceFrom string `json:"destination_replace_from,omitempty" path:"destination_replace_from,omitempty" url:"destination_replace_from,omitempty"`
	// If set, this string will replace the value `destination_replace_from` in the destination filename. You can use special patterns here.
	DestinationReplaceTo string `json:"destination_replace_to,omitempty" path:"destination_replace_to,omitempty" url:"destination_replace_to,omitempty"`
	// Destination Paths.
	Destinations []string `json:"destinations,omitempty" path:"destinations,omitempty" url:"destinations,omitempty"`
	// If true, this automation will not run.
	Disabled *bool `json:"disabled,omitempty" path:"disabled,omitempty" url:"disabled,omitempty"`
	// If set, this glob pattern will exclude files from the automation. Supports globs, except on remote mounts.
	ExcludePattern string `json:"exclude_pattern,omitempty" path:"exclude_pattern,omitempty" url:"exclude_pattern,omitempty"`
	// List of URLs to be imported and names to be used.
	ImportUrls []map[string]interface{} `json:"import_urls,omitempty" path:"import_urls,omitempty" url:"import_urls,omitempty"`
	// If trigger is `email`, this is the address that triggers the Automation.
	InboundEmailAddress string `json:"inbound_email_address,omitempty" path:"inbound_email_address,omitempty" url:"inbound_email_address,omitempty"`
	// Normally copy and move automations that use globs will implicitly preserve the source folder structure in the destination. If this flag is `true`, the source folder structure will be flattened in the destination. This is useful for copying or moving files from multiple folders into a single destination folder.
	FlattenDestinationStructure *bool `json:"flatten_destination_structure,omitempty" path:"flatten_destination_structure,omitempty" url:"flatten_destination_structure,omitempty"`
	// IDs of Groups for the Automation (i.e. who to Request File from).
	GroupIds []int64 `json:"group_ids,omitempty" path:"group_ids,omitempty" url:"group_ids,omitempty"`
	// If true, the Lock Folders behavior will be disregarded for automated actions.
	IgnoreLockedFolders *bool `json:"ignore_locked_folders,omitempty" path:"ignore_locked_folders,omitempty" url:"ignore_locked_folders,omitempty"`
	// If trigger is `daily`, this specifies how often to run this automation. One of: `day`, `week`, `week_end`, `month`, `month_end`, `quarter`, `quarter_end`, `year`, `year_end`.
	Interval string `json:"interval,omitempty" path:"interval,omitempty" url:"interval,omitempty"`
	// Time when automation was last modified. Does not change for name or description updates.
	LastModifiedAt *time.Time `json:"last_modified_at,omitempty" path:"last_modified_at,omitempty" url:"last_modified_at,omitempty"`
	// If `true`, use the legacy behavior for this automation, where it can operate on folders in addition to just files. This behavior no longer works and should not be used.
	LegacyFolderMatching *bool `json:"legacy_folder_matching,omitempty" path:"legacy_folder_matching,omitempty" url:"legacy_folder_matching,omitempty"`
	// Name for this automation.
	Name string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	// If true, existing files will be overwritten with new files on Move/Copy automations. Note: by default files will not be overwritten on Copy automations if they appear to be the same file size as the newly incoming file. Use the `always_overwrite_size_matching_files` option in conjunction with `overwrite_files` to override this behavior and overwrite files no matter what.
	OverwriteFiles *bool `json:"overwrite_files,omitempty" path:"overwrite_files,omitempty" url:"overwrite_files,omitempty"`
	// Path on which this Automation runs. Supports globs, except on remote mounts. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.
	Path string `json:"path,omitempty" path:"path,omitempty" url:"path,omitempty"`
	// Timezone to use when rendering timestamps in paths.
	PathTimeZone string `json:"path_time_zone,omitempty" path:"path_time_zone,omitempty" url:"path_time_zone,omitempty"`
	// Current Automation v2 definition version.
	Version int64 `json:"version,omitempty" path:"version,omitempty" url:"version,omitempty"`
	// If trigger type is `daily`, this specifies a day number to run in one of the supported intervals: `week`, `month`, `quarter`, `year`.
	RecurringDay int64 `json:"recurring_day,omitempty" path:"recurring_day,omitempty" url:"recurring_day,omitempty"`
	// If trigger type is `daily`, this specifies one or more day numbers to run in one of the supported intervals: `week`, `month`, `quarter`, `year`.
	RecurringDays []int64 `json:"recurring_days,omitempty" path:"recurring_days,omitempty" url:"recurring_days,omitempty"`
	// If trigger is `custom_schedule`, the reusable Schedule used instead of the automation's schedule fields.
	ScheduleId int64 `json:"schedule_id,omitempty" path:"schedule_id,omitempty" url:"schedule_id,omitempty"`
	// If the Automation fails, retry at this interval (in minutes). Acceptable values are 5 through 1440 (one day). Set to null to disable.
	RetryOnFailureIntervalInMinutes int64 `json:"retry_on_failure_interval_in_minutes,omitempty" path:"retry_on_failure_interval_in_minutes,omitempty" url:"retry_on_failure_interval_in_minutes,omitempty"`
	// If the Automation fails, retry at most this many times. Maximum allowed value: 10. Set to null to disable.
	RetryOnFailureNumberOfAttempts int64 `json:"retry_on_failure_number_of_attempts,omitempty" path:"retry_on_failure_number_of_attempts,omitempty" url:"retry_on_failure_number_of_attempts,omitempty"`
	// If trigger is `custom_schedule`, Custom schedule description for when the automation should be run in json format.
	Schedule interface{} `json:"schedule,omitempty" path:"schedule,omitempty" url:"schedule,omitempty"`
	// If trigger is `custom_schedule` or `daily` with times, Human readable schedule description for when the automation should be run.
	HumanReadableSchedule string `json:"human_readable_schedule,omitempty" path:"human_readable_schedule,omitempty" url:"human_readable_schedule,omitempty"`
	// If trigger is `custom_schedule`, Custom schedule description for when the automation should be run. 0 is Sunday, 1 is Monday, etc.
	ScheduleDaysOfWeek []int64 `json:"schedule_days_of_week,omitempty" path:"schedule_days_of_week,omitempty" url:"schedule_days_of_week,omitempty"`
	// Times of day to run in HH:MM format (24-hour). For `custom_schedule`, run at these times on specified days of week. For `daily`, run at these times on the scheduled interval date.
	ScheduleTimesOfDay []string `json:"schedule_times_of_day,omitempty" path:"schedule_times_of_day,omitempty" url:"schedule_times_of_day,omitempty"`
	// Time zone for the schedule. If not set, times are interpreted as UTC.
	ScheduleTimeZone string `json:"schedule_time_zone,omitempty" path:"schedule_time_zone,omitempty" url:"schedule_time_zone,omitempty"`
	// Source path/glob. See Automation docs for exact description, but this is used to filter for files in the `path` to find files to operate on. Supports globs, except on remote mounts.
	Source string `json:"source,omitempty" path:"source,omitempty" url:"source,omitempty"`
	// IDs of Syncs to run by this Automation.
	SyncIds []int64 `json:"sync_ids,omitempty" path:"sync_ids,omitempty" url:"sync_ids,omitempty"`
	// If trigger is `action`, this is the list of action types on which to trigger the automation. Valid actions are create, copy, move, archived_delete, update, read, destroy.
	TriggerActions []string `json:"trigger_actions,omitempty" path:"trigger_actions,omitempty" url:"trigger_actions,omitempty"`
	// How this automation is triggered to run.
	Trigger string `json:"trigger,omitempty" path:"trigger,omitempty" url:"trigger,omitempty"`
	// User ID of the Automation's creator.
	UserId int64 `json:"user_id,omitempty" path:"user_id,omitempty" url:"user_id,omitempty"`
	// IDs of Users for the Automation (i.e. who to Request File from).
	UserIds []int64 `json:"user_ids,omitempty" path:"user_ids,omitempty" url:"user_ids,omitempty"`
	// A Hash of attributes specific to the automation type.
	Value interface{} `json:"value,omitempty" path:"value,omitempty" url:"value,omitempty"`
	// If trigger is `webhook`, this is the URL of the webhook to trigger the Automation.
	WebhookUrl string `json:"webhook_url,omitempty" path:"webhook_url,omitempty" url:"webhook_url,omitempty"`
	// Skip the automation if there is a formal, observed holiday for this region.
	HolidayRegion string `json:"holiday_region,omitempty" path:"holiday_region,omitempty" url:"holiday_region,omitempty"`
}

// Identifier returns the resource ID.
func (a Automation) Identifier() interface{} {
	return a.Id
}

// AutomationCollection is a list of Automation resources.
type AutomationCollection []Automation

// AutomationTriggerEnum is a string value for trigger.
// Enum lists the values documented by the API.
type AutomationTriggerEnum string

// String returns the API parameter value.
func (u AutomationTriggerEnum) String() string {
	return string(u)
}

// Enum returns the documented values keyed by their API strings.
func (u AutomationTriggerEnum) Enum() map[string]AutomationTriggerEnum {
	return map[string]AutomationTriggerEnum{
		"manual":          AutomationTriggerEnum("manual"),
		"daily":           AutomationTriggerEnum("daily"),
		"custom_schedule": AutomationTriggerEnum("custom_schedule"),
		"webhook":         AutomationTriggerEnum("webhook"),
		"email":           AutomationTriggerEnum("email"),
		"action":          AutomationTriggerEnum("action"),
	}
}

// AutomationEnum is a string value for automation.
// Enum lists the values documented by the API.
type AutomationEnum string

// String returns the API parameter value.
func (u AutomationEnum) String() string {
	return string(u)
}

// Enum returns the documented values keyed by their API strings.
func (u AutomationEnum) Enum() map[string]AutomationEnum {
	return map[string]AutomationEnum{
		"create_folder": AutomationEnum("create_folder"),
		"delete_file":   AutomationEnum("delete_file"),
		"copy_file":     AutomationEnum("copy_file"),
		"move_file":     AutomationEnum("move_file"),
		"as2_send":      AutomationEnum("as2_send"),
		"run_sync":      AutomationEnum("run_sync"),
		"import_file":   AutomationEnum("import_file"),
		"v2":            AutomationEnum("v2"),
	}
}

// AutomationListParams contains the request parameters for GET /automations.
type AutomationListParams struct {
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id`, `name`, `automation`, `last_modified_at` or `disabled`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `disabled`, `last_modified_at`, `workspace_id` or `automation`. Valid field combinations are `[ disabled, last_modified_at ]`, `[ workspace_id, disabled ]`, `[ disabled, automation ]`, `[ workspace_id, last_modified_at ]`, `[ automation, last_modified_at ]`, `[ workspace_id, automation ]`, `[ workspace_id, disabled, last_modified_at ]`, `[ disabled, automation, last_modified_at ]`, `[ workspace_id, disabled, automation ]`, `[ workspace_id, automation, last_modified_at ]` or `[ workspace_id, disabled, automation, last_modified_at ]`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// If set, return records where the specified field is greater than the supplied value. Valid fields are `last_modified_at`.
	FilterGt interface{} `url:"filter_gt,omitempty" json:"filter_gt,omitempty" path:"filter_gt"`
	// If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `last_modified_at`.
	FilterGteq interface{} `url:"filter_gteq,omitempty" json:"filter_gteq,omitempty" path:"filter_gteq"`
	// If set, return records where the specified field is less than the supplied value. Valid fields are `last_modified_at`.
	FilterLt interface{} `url:"filter_lt,omitempty" json:"filter_lt,omitempty" path:"filter_lt"`
	// If set, return records where the specified field is less than or equal the supplied value. Valid fields are `last_modified_at`.
	FilterLteq interface{} `url:"filter_lteq,omitempty" json:"filter_lteq,omitempty" path:"filter_lteq"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// AutomationFindParams contains the request parameters for GET /automations/{id}.
type AutomationFindParams struct {
	// Automation ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// AutomationCreateParams contains the request parameters for POST /automations.
type AutomationCreateParams struct {
	// Source path/glob. See Automation docs for exact description, but this is used to filter for files in the `path` to find files to operate on. Supports globs, except on remote mounts.
	Source string `url:"source,omitempty" json:"source,omitempty" path:"source"`
	// A list of destination paths. Use a trailing slash for folder destinations and omit it for file destinations.
	Destinations []string `url:"destinations,omitempty" json:"destinations,omitempty" path:"destinations"`
	// If set, this string in the destination path will be replaced with the value in `destination_replace_to`.
	DestinationReplaceFrom string `url:"destination_replace_from,omitempty" json:"destination_replace_from,omitempty" path:"destination_replace_from"`
	// If set, this string will replace the value `destination_replace_from` in the destination filename. You can use special patterns here.
	DestinationReplaceTo string `url:"destination_replace_to,omitempty" json:"destination_replace_to,omitempty" path:"destination_replace_to"`
	// How often to run this automation? One of: `day`, `week`, `week_end`, `month`, `month_end`, `quarter`, `quarter_end`, `year`, `year_end`.
	Interval string `url:"interval,omitempty" json:"interval,omitempty" path:"interval"`
	// Path on which this Automation runs. Supports globs, except on remote mounts.
	Path string `url:"path,omitempty" json:"path,omitempty" path:"path"`
	// A list of sync IDs the automation is associated with. If sent as a string, it should be comma-delimited.
	SyncIds string `url:"sync_ids,omitempty" json:"sync_ids,omitempty" path:"sync_ids"`
	// A list of user IDs the automation is associated with. If sent as a string, it should be comma-delimited.
	UserIds string `url:"user_ids,omitempty" json:"user_ids,omitempty" path:"user_ids"`
	// A list of group IDs the automation is associated with. If sent as a string, it should be comma-delimited.
	GroupIds string `url:"group_ids,omitempty" json:"group_ids,omitempty" path:"group_ids"`
	// If trigger is `custom_schedule`, the reusable Schedule used instead of the automation's schedule fields.
	ScheduleId int64 `url:"schedule_id,omitempty" json:"schedule_id,omitempty" path:"schedule_id"`
	// If trigger is `custom_schedule` without `schedule_id`, a list of days of the week to run this automation. 0 is Sunday, 1 is Monday, etc.
	ScheduleDaysOfWeek []int64 `url:"schedule_days_of_week,omitempty" json:"schedule_days_of_week,omitempty" path:"schedule_days_of_week"`
	// Times of day to run in HH:MM format (24-hour). Required for `custom_schedule` without `schedule_id`. Optional for `daily` triggers; if not set, runs at midnight UTC.
	ScheduleTimesOfDay []string `url:"schedule_times_of_day,omitempty" json:"schedule_times_of_day,omitempty" path:"schedule_times_of_day"`
	// Time zone for schedule fields. Optional for `custom_schedule` without `schedule_id` and for `daily`. If not set, times are interpreted as UTC.
	ScheduleTimeZone string `url:"schedule_time_zone,omitempty" json:"schedule_time_zone,omitempty" path:"schedule_time_zone"`
	// Skip automation on holidays in this region. Optional for `custom_schedule` without `schedule_id` and for `daily`.
	HolidayRegion string `url:"holiday_region,omitempty" json:"holiday_region,omitempty" path:"holiday_region"`
	// Ordinarily, files with identical size in the source and destination will be skipped from copy operations to prevent wasted transfer. If this flag is `true` we will overwrite the destination file always. Note that this may cause large amounts of wasted transfer usage. This setting has no effect unless `overwrite_files` is also set to `true`.
	AlwaysOverwriteSizeMatchingFiles *bool `url:"always_overwrite_size_matching_files,omitempty" json:"always_overwrite_size_matching_files,omitempty" path:"always_overwrite_size_matching_files"`
	// Ordinarily, we will allow automation runs to run in parallel for non-scheduled automations. If this flag is `true` we will force automation runs to be serialized (run one at a time, one after another). This can resolve some issues with race conditions on remote systems at the cost of some performance.
	AlwaysSerializeJobs *bool `url:"always_serialize_jobs,omitempty" json:"always_serialize_jobs,omitempty" path:"always_serialize_jobs"`
	// Description for the this Automation.
	Description string `url:"description,omitempty" json:"description,omitempty" path:"description"`
	// Automation v2 graph definition.
	Definition interface{} `url:"definition,omitempty" json:"definition,omitempty" path:"definition"`
	// If true, this automation will not run.
	Disabled *bool `url:"disabled,omitempty" json:"disabled,omitempty" path:"disabled"`
	// If set, this glob pattern will exclude files from the automation. Supports globs, except on remote mounts.
	ExcludePattern string `url:"exclude_pattern,omitempty" json:"exclude_pattern,omitempty" path:"exclude_pattern"`
	// List of URLs to be imported and names to be used.
	ImportUrls []map[string]interface{} `url:"import_urls,omitempty" json:"import_urls,omitempty" path:"import_urls"`
	// Normally copy and move automations that use globs will implicitly preserve the source folder structure in the destination. If this flag is `true`, the source folder structure will be flattened in the destination. This is useful for copying or moving files from multiple folders into a single destination folder.
	FlattenDestinationStructure *bool `url:"flatten_destination_structure,omitempty" json:"flatten_destination_structure,omitempty" path:"flatten_destination_structure"`
	// If true, the Lock Folders behavior will be disregarded for automated actions.
	IgnoreLockedFolders *bool `url:"ignore_locked_folders,omitempty" json:"ignore_locked_folders,omitempty" path:"ignore_locked_folders"`
	// DEPRECATED: If `true`, use the legacy behavior for this automation, where it can operate on folders in addition to just files. This behavior no longer works and should not be used.
	LegacyFolderMatching *bool `url:"legacy_folder_matching,omitempty" json:"legacy_folder_matching,omitempty" path:"legacy_folder_matching"`
	// Name for this automation.
	Name string `url:"name,omitempty" json:"name,omitempty" path:"name"`
	// If true, existing files will be overwritten with new files on Move/Copy automations. Note: by default files will not be overwritten on Copy automations if they appear to be the same file size as the newly incoming file. Use the `always_overwrite_size_matching_files` option in conjunction with `overwrite_files` to override this behavior and overwrite files no matter what.
	OverwriteFiles *bool `url:"overwrite_files,omitempty" json:"overwrite_files,omitempty" path:"overwrite_files"`
	// Timezone to use when rendering timestamps in paths.
	PathTimeZone string `url:"path_time_zone,omitempty" json:"path_time_zone,omitempty" path:"path_time_zone"`
	// If the Automation fails, retry at this interval (in minutes). Acceptable values are 5 through 1440 (one day). Set to null to disable.
	RetryOnFailureIntervalInMinutes int64 `url:"retry_on_failure_interval_in_minutes,omitempty" json:"retry_on_failure_interval_in_minutes,omitempty" path:"retry_on_failure_interval_in_minutes"`
	// If the Automation fails, retry at most this many times. Maximum allowed value: 10. Set to null to disable.
	RetryOnFailureNumberOfAttempts int64 `url:"retry_on_failure_number_of_attempts,omitempty" json:"retry_on_failure_number_of_attempts,omitempty" path:"retry_on_failure_number_of_attempts"`
	// How this automation is triggered to run. Values: manual, daily, custom_schedule, webhook, email, action.
	Trigger AutomationTriggerEnum `url:"trigger,omitempty" json:"trigger,omitempty" path:"trigger"`
	// If trigger is `action`, this is the list of action types on which to trigger the automation. Valid actions are create, copy, move, archived_delete, update, read, destroy.
	TriggerActions []string `url:"trigger_actions,omitempty" json:"trigger_actions,omitempty" path:"trigger_actions"`
	// A Hash of attributes specific to the automation type.
	Value interface{} `url:"value,omitempty" json:"value,omitempty" path:"value"`
	// If trigger type is `daily`, this specifies a day number to run in one of the supported intervals: `week`, `month`, `quarter`, `year`.
	RecurringDay int64 `url:"recurring_day,omitempty" json:"recurring_day,omitempty" path:"recurring_day"`
	// If trigger type is `daily`, this specifies one or more day numbers to run in one of the supported intervals: `week`, `month`, `quarter`, `year`.
	RecurringDays []int64 `url:"recurring_days,omitempty" json:"recurring_days,omitempty" path:"recurring_days"`
	// Automation type. Values: create_folder, delete_file, copy_file, move_file, as2_send, run_sync, import_file, v2. Required.
	Automation AutomationEnum `url:"automation" json:"automation" path:"automation"`
	// Workspace ID.
	WorkspaceId int64 `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
}

// AutomationUpgradeParams contains the request parameters for POST /automations/{id}/upgrade.
//
// Upgrade a legacy Automation to Automation v2
type AutomationUpgradeParams struct {
	// Automation ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// AutomationManualRunParams contains the request parameters for POST /automations/{id}/manual_run.
//
// Manually Run Automation
type AutomationManualRunParams struct {
	// Automation ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// Initial items for a v2 manual trigger. Each item contains exactly one `file` path or `data` object.
	Items []map[string]interface{} `url:"items,omitempty" json:"items,omitempty" path:"items"`
}

// AutomationUpdateParams contains the request parameters for PATCH /automations/{id}.
type AutomationUpdateParams struct {
	// Automation ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// Source path/glob. See Automation docs for exact description, but this is used to filter for files in the `path` to find files to operate on. Supports globs, except on remote mounts.
	Source string `url:"source,omitempty" json:"source,omitempty" path:"source"`
	// A list of destination paths. Use a trailing slash for folder destinations and omit it for file destinations.
	Destinations []string `url:"destinations,omitempty" json:"destinations,omitempty" path:"destinations"`
	// If set, this string in the destination path will be replaced with the value in `destination_replace_to`.
	DestinationReplaceFrom string `url:"destination_replace_from,omitempty" json:"destination_replace_from,omitempty" path:"destination_replace_from"`
	// If set, this string will replace the value `destination_replace_from` in the destination filename. You can use special patterns here.
	DestinationReplaceTo string `url:"destination_replace_to,omitempty" json:"destination_replace_to,omitempty" path:"destination_replace_to"`
	// How often to run this automation? One of: `day`, `week`, `week_end`, `month`, `month_end`, `quarter`, `quarter_end`, `year`, `year_end`.
	Interval string `url:"interval,omitempty" json:"interval,omitempty" path:"interval"`
	// Path on which this Automation runs. Supports globs, except on remote mounts.
	Path string `url:"path,omitempty" json:"path,omitempty" path:"path"`
	// A list of sync IDs the automation is associated with. If sent as a string, it should be comma-delimited.
	SyncIds string `url:"sync_ids,omitempty" json:"sync_ids,omitempty" path:"sync_ids"`
	// A list of user IDs the automation is associated with. If sent as a string, it should be comma-delimited.
	UserIds string `url:"user_ids,omitempty" json:"user_ids,omitempty" path:"user_ids"`
	// A list of group IDs the automation is associated with. If sent as a string, it should be comma-delimited.
	GroupIds string `url:"group_ids,omitempty" json:"group_ids,omitempty" path:"group_ids"`
	// If trigger is `custom_schedule`, the reusable Schedule used instead of the automation's schedule fields.
	ScheduleId int64 `url:"schedule_id,omitempty" json:"schedule_id,omitempty" path:"schedule_id"`
	// If trigger is `custom_schedule` without `schedule_id`, a list of days of the week to run this automation. 0 is Sunday, 1 is Monday, etc.
	ScheduleDaysOfWeek []int64 `url:"schedule_days_of_week,omitempty" json:"schedule_days_of_week,omitempty" path:"schedule_days_of_week"`
	// Times of day to run in HH:MM format (24-hour). Required for `custom_schedule` without `schedule_id`. Optional for `daily` triggers; if not set, runs at midnight UTC.
	ScheduleTimesOfDay []string `url:"schedule_times_of_day,omitempty" json:"schedule_times_of_day,omitempty" path:"schedule_times_of_day"`
	// Time zone for schedule fields. Optional for `custom_schedule` without `schedule_id` and for `daily`. If not set, times are interpreted as UTC.
	ScheduleTimeZone string `url:"schedule_time_zone,omitempty" json:"schedule_time_zone,omitempty" path:"schedule_time_zone"`
	// Skip automation on holidays in this region. Optional for `custom_schedule` without `schedule_id` and for `daily`.
	HolidayRegion string `url:"holiday_region,omitempty" json:"holiday_region,omitempty" path:"holiday_region"`
	// Ordinarily, files with identical size in the source and destination will be skipped from copy operations to prevent wasted transfer. If this flag is `true` we will overwrite the destination file always. Note that this may cause large amounts of wasted transfer usage. This setting has no effect unless `overwrite_files` is also set to `true`.
	AlwaysOverwriteSizeMatchingFiles *bool `url:"always_overwrite_size_matching_files,omitempty" json:"always_overwrite_size_matching_files,omitempty" path:"always_overwrite_size_matching_files"`
	// Ordinarily, we will allow automation runs to run in parallel for non-scheduled automations. If this flag is `true` we will force automation runs to be serialized (run one at a time, one after another). This can resolve some issues with race conditions on remote systems at the cost of some performance.
	AlwaysSerializeJobs *bool `url:"always_serialize_jobs,omitempty" json:"always_serialize_jobs,omitempty" path:"always_serialize_jobs"`
	// Description for the this Automation.
	Description string `url:"description,omitempty" json:"description,omitempty" path:"description"`
	// Automation v2 graph definition.
	Definition interface{} `url:"definition,omitempty" json:"definition,omitempty" path:"definition"`
	// If true, this automation will not run.
	Disabled *bool `url:"disabled,omitempty" json:"disabled,omitempty" path:"disabled"`
	// If set, this glob pattern will exclude files from the automation. Supports globs, except on remote mounts.
	ExcludePattern string `url:"exclude_pattern,omitempty" json:"exclude_pattern,omitempty" path:"exclude_pattern"`
	// List of URLs to be imported and names to be used.
	ImportUrls []map[string]interface{} `url:"import_urls,omitempty" json:"import_urls,omitempty" path:"import_urls"`
	// Normally copy and move automations that use globs will implicitly preserve the source folder structure in the destination. If this flag is `true`, the source folder structure will be flattened in the destination. This is useful for copying or moving files from multiple folders into a single destination folder.
	FlattenDestinationStructure *bool `url:"flatten_destination_structure,omitempty" json:"flatten_destination_structure,omitempty" path:"flatten_destination_structure"`
	// If true, the Lock Folders behavior will be disregarded for automated actions.
	IgnoreLockedFolders *bool `url:"ignore_locked_folders,omitempty" json:"ignore_locked_folders,omitempty" path:"ignore_locked_folders"`
	// DEPRECATED: If `true`, use the legacy behavior for this automation, where it can operate on folders in addition to just files. This behavior no longer works and should not be used.
	LegacyFolderMatching *bool `url:"legacy_folder_matching,omitempty" json:"legacy_folder_matching,omitempty" path:"legacy_folder_matching"`
	// Name for this automation.
	Name string `url:"name,omitempty" json:"name,omitempty" path:"name"`
	// If true, existing files will be overwritten with new files on Move/Copy automations. Note: by default files will not be overwritten on Copy automations if they appear to be the same file size as the newly incoming file. Use the `always_overwrite_size_matching_files` option in conjunction with `overwrite_files` to override this behavior and overwrite files no matter what.
	OverwriteFiles *bool `url:"overwrite_files,omitempty" json:"overwrite_files,omitempty" path:"overwrite_files"`
	// Timezone to use when rendering timestamps in paths.
	PathTimeZone string `url:"path_time_zone,omitempty" json:"path_time_zone,omitempty" path:"path_time_zone"`
	// If the Automation fails, retry at this interval (in minutes). Acceptable values are 5 through 1440 (one day). Set to null to disable.
	RetryOnFailureIntervalInMinutes int64 `url:"retry_on_failure_interval_in_minutes,omitempty" json:"retry_on_failure_interval_in_minutes,omitempty" path:"retry_on_failure_interval_in_minutes"`
	// If the Automation fails, retry at most this many times. Maximum allowed value: 10. Set to null to disable.
	RetryOnFailureNumberOfAttempts int64 `url:"retry_on_failure_number_of_attempts,omitempty" json:"retry_on_failure_number_of_attempts,omitempty" path:"retry_on_failure_number_of_attempts"`
	// How this automation is triggered to run. Values: manual, daily, custom_schedule, webhook, email, action.
	Trigger AutomationTriggerEnum `url:"trigger,omitempty" json:"trigger,omitempty" path:"trigger"`
	// If trigger is `action`, this is the list of action types on which to trigger the automation. Valid actions are create, copy, move, archived_delete, update, read, destroy.
	TriggerActions []string `url:"trigger_actions,omitempty" json:"trigger_actions,omitempty" path:"trigger_actions"`
	// A Hash of attributes specific to the automation type.
	Value interface{} `url:"value,omitempty" json:"value,omitempty" path:"value"`
	// If trigger type is `daily`, this specifies a day number to run in one of the supported intervals: `week`, `month`, `quarter`, `year`.
	RecurringDay int64 `url:"recurring_day,omitempty" json:"recurring_day,omitempty" path:"recurring_day"`
	// If trigger type is `daily`, this specifies one or more day numbers to run in one of the supported intervals: `week`, `month`, `quarter`, `year`.
	RecurringDays []int64 `url:"recurring_days,omitempty" json:"recurring_days,omitempty" path:"recurring_days"`
	// Automation type. Values: create_folder, delete_file, copy_file, move_file, as2_send, run_sync, import_file, v2.
	Automation AutomationEnum `url:"automation,omitempty" json:"automation,omitempty" path:"automation"`
}

// AutomationDeleteParams contains the request parameters for DELETE /automations/{id}.
type AutomationDeleteParams struct {
	// Automation ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (a *Automation) UnmarshalJSON(data []byte) error {
	type automation Automation
	var v automation
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*a = Automation(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (a *AutomationCollection) UnmarshalJSON(data []byte) error {
	type automations AutomationCollection
	var v automations
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*a = AutomationCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (a *AutomationCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*a))
	for i, v := range *a {
		ret[i] = v
	}

	return &ret
}
