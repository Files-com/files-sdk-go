package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// Sync is a Files.com API resource.
type Sync struct {
	// Sync ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Name for this sync job.
	Name string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	// Description for this sync job.
	Description string `json:"description,omitempty" path:"description,omitempty" url:"description,omitempty"`
	// Site ID this sync belongs to.
	SiteId int64 `json:"site_id,omitempty" path:"site_id,omitempty" url:"site_id,omitempty"`
	// Workspace ID this sync belongs to.
	WorkspaceId int64 `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	// User who created or owns this sync.
	UserId int64 `json:"user_id,omitempty" path:"user_id,omitempty" url:"user_id,omitempty"`
	// Absolute source path for the sync.
	SrcPath string `json:"src_path,omitempty" path:"src_path,omitempty" url:"src_path,omitempty"`
	// Absolute destination path for the sync.
	DestPath string `json:"dest_path,omitempty" path:"dest_path,omitempty" url:"dest_path,omitempty"`
	// Remote server ID for the source (if remote).
	SrcRemoteServerId int64 `json:"src_remote_server_id,omitempty" path:"src_remote_server_id,omitempty" url:"src_remote_server_id,omitempty"`
	// Remote server ID for the destination (if remote).
	DestRemoteServerId int64 `json:"dest_remote_server_id,omitempty" path:"dest_remote_server_id,omitempty" url:"dest_remote_server_id,omitempty"`
	// Source site ID if syncing from a child or partner site.
	SrcSiteId int64 `json:"src_site_id,omitempty" path:"src_site_id,omitempty" url:"src_site_id,omitempty"`
	// Destination site ID if syncing to a child or partner site.
	DestSiteId int64 `json:"dest_site_id,omitempty" path:"dest_site_id,omitempty" url:"dest_site_id,omitempty"`
	// Is this a two-way sync?
	TwoWay *bool `json:"two_way,omitempty" path:"two_way,omitempty" url:"two_way,omitempty"`
	// Keep files after copying?
	KeepAfterCopy *bool `json:"keep_after_copy,omitempty" path:"keep_after_copy,omitempty" url:"keep_after_copy,omitempty"`
	// Delete empty folders after sync?
	DeleteEmptyFolders *bool `json:"delete_empty_folders,omitempty" path:"delete_empty_folders,omitempty" url:"delete_empty_folders,omitempty"`
	// Is this sync disabled?
	Disabled *bool `json:"disabled,omitempty" path:"disabled,omitempty" url:"disabled,omitempty"`
	// Trigger type: daily, custom_schedule, or manual.
	Trigger string `json:"trigger,omitempty" path:"trigger,omitempty" url:"trigger,omitempty"`
	// Some MFT services request an empty file (known as a trigger file) to signal the sync is complete and they can begin further processing. If trigger_file is set, a zero-byte file will be sent at the end of the sync.
	TriggerFile string `json:"trigger_file,omitempty" path:"trigger_file,omitempty" url:"trigger_file,omitempty"`
	// If true, the trigger file will be sent at the end of a successful sync even when no files were transferred.
	AlwaysWriteTriggerFile *bool `json:"always_write_trigger_file,omitempty" path:"always_write_trigger_file,omitempty" url:"always_write_trigger_file,omitempty"`
	// Array of glob patterns to include.
	IncludePatterns []string `json:"include_patterns,omitempty" path:"include_patterns,omitempty" url:"include_patterns,omitempty"`
	// Array of glob patterns to exclude.
	ExcludePatterns []string `json:"exclude_patterns,omitempty" path:"exclude_patterns,omitempty" url:"exclude_patterns,omitempty"`
	// When this sync was created.
	CreatedAt *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	// When this sync was last updated.
	UpdatedAt *time.Time `json:"updated_at,omitempty" path:"updated_at,omitempty" url:"updated_at,omitempty"`
	// Frequency in minutes between syncs. If set, this value must be greater than or equal to the `remote_sync_interval` value for the site's plan. If left blank, the plan's `remote_sync_interval` will be used. This setting is only used if `trigger` is empty.
	SyncIntervalMinutes int64 `json:"sync_interval_minutes,omitempty" path:"sync_interval_minutes,omitempty" url:"sync_interval_minutes,omitempty"`
	// If trigger is `daily`, this specifies how often to run this sync. One of: `day`, `week`, `week_end`, `month`, `month_end`, `quarter`, `quarter_end`, `year`, `year_end`.
	Interval string `json:"interval,omitempty" path:"interval,omitempty" url:"interval,omitempty"`
	// If trigger type is `daily`, this specifies a day number to run in one of the supported intervals: `week`, `month`, `quarter`, `year`.
	RecurringDay int64 `json:"recurring_day,omitempty" path:"recurring_day,omitempty" url:"recurring_day,omitempty"`
	// If trigger type is `daily`, this specifies one or more day numbers to run in one of the supported intervals: `week`, `month`, `quarter`, `year`.
	RecurringDays []int64 `json:"recurring_days,omitempty" path:"recurring_days,omitempty" url:"recurring_days,omitempty"`
	// If trigger is `custom_schedule`, the reusable Schedule used instead of the sync's schedule fields.
	ScheduleId int64 `json:"schedule_id,omitempty" path:"schedule_id,omitempty" url:"schedule_id,omitempty"`
	// If trigger is `custom_schedule`, Custom schedule description for when the sync should be run. 0-based days of the week. 0 is Sunday, 1 is Monday, etc.
	ScheduleDaysOfWeek []int64 `json:"schedule_days_of_week,omitempty" path:"schedule_days_of_week,omitempty" url:"schedule_days_of_week,omitempty"`
	// Times of day to run in HH:MM format. For `custom_schedule`, run at these times on specified days of week. For `daily`, run at these times on the scheduled interval date.
	ScheduleTimesOfDay []string `json:"schedule_times_of_day,omitempty" path:"schedule_times_of_day,omitempty" url:"schedule_times_of_day,omitempty"`
	// Time zone for the schedule. If not set, times are interpreted as UTC.
	ScheduleTimeZone string `json:"schedule_time_zone,omitempty" path:"schedule_time_zone,omitempty" url:"schedule_time_zone,omitempty"`
	// Skip the sync if there is a formal, observed holiday for this region.
	HolidayRegion string `json:"holiday_region,omitempty" path:"holiday_region,omitempty" url:"holiday_region,omitempty"`
	// The latest run of this sync.
	LatestSyncRun SyncRun `json:"latest_sync_run,omitempty" path:"latest_sync_run,omitempty" url:"latest_sync_run,omitempty"`
}

// Identifier returns the resource ID.
func (s Sync) Identifier() interface{} {
	return s.Id
}

// SyncCollection is a list of Sync resources.
type SyncCollection []Sync

// SyncTriggerEnum is a string value for trigger.
// Enum lists the values documented by the API.
type SyncTriggerEnum string

// String returns the API parameter value.
func (u SyncTriggerEnum) String() string {
	return string(u)
}

// Enum returns the documented values keyed by their API strings.
func (u SyncTriggerEnum) Enum() map[string]SyncTriggerEnum {
	return map[string]SyncTriggerEnum{
		"daily":           SyncTriggerEnum("daily"),
		"custom_schedule": SyncTriggerEnum("custom_schedule"),
		"manual":          SyncTriggerEnum("manual"),
	}
}

// SyncListParams contains the request parameters for GET /syncs.
type SyncListParams struct {
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `site_id`, `workspace_id` or `name`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `workspace_id`, `disabled`, `src_remote_server_id` or `dest_remote_server_id`. Valid field combinations are `[ workspace_id, disabled ]`, `[ workspace_id, src_remote_server_id ]`, `[ workspace_id, dest_remote_server_id ]`, `[ disabled, src_remote_server_id ]`, `[ disabled, dest_remote_server_id ]`, `[ workspace_id, disabled, src_remote_server_id ]` or `[ workspace_id, disabled, dest_remote_server_id ]`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// SyncFindParams contains the request parameters for GET /syncs/{id}.
type SyncFindParams struct {
	// Sync ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// SyncCreateParams contains the request parameters for POST /syncs.
type SyncCreateParams struct {
	// Delete empty folders after sync?
	DeleteEmptyFolders *bool `url:"delete_empty_folders,omitempty" json:"delete_empty_folders,omitempty" path:"delete_empty_folders"`
	// Description for this sync job.
	Description string `url:"description,omitempty" json:"description,omitempty" path:"description"`
	// Absolute destination path for the sync.
	DestPath string `url:"dest_path,omitempty" json:"dest_path,omitempty" path:"dest_path"`
	// Remote server ID for the destination (if remote).
	DestRemoteServerId int64 `url:"dest_remote_server_id,omitempty" json:"dest_remote_server_id,omitempty" path:"dest_remote_server_id"`
	// Is this sync disabled?
	Disabled *bool `url:"disabled,omitempty" json:"disabled,omitempty" path:"disabled"`
	// Array of glob patterns to exclude.
	ExcludePatterns []string `url:"exclude_patterns,omitempty" json:"exclude_patterns,omitempty" path:"exclude_patterns"`
	// Skip the sync if there is a formal, observed holiday for this region.
	HolidayRegion string `url:"holiday_region,omitempty" json:"holiday_region,omitempty" path:"holiday_region"`
	// Array of glob patterns to include.
	IncludePatterns []string `url:"include_patterns,omitempty" json:"include_patterns,omitempty" path:"include_patterns"`
	// If trigger is `daily`, this specifies how often to run this sync. One of: `day`, `week`, `week_end`, `month`, `month_end`, `quarter`, `quarter_end`, `year`, `year_end`.
	Interval string `url:"interval,omitempty" json:"interval,omitempty" path:"interval"`
	// Keep files after copying?
	KeepAfterCopy *bool `url:"keep_after_copy,omitempty" json:"keep_after_copy,omitempty" path:"keep_after_copy"`
	// Name for this sync job.
	Name string `url:"name,omitempty" json:"name,omitempty" path:"name"`
	// If trigger type is `daily`, this specifies a day number to run in one of the supported intervals: `week`, `month`, `quarter`, `year`.
	RecurringDay int64 `url:"recurring_day,omitempty" json:"recurring_day,omitempty" path:"recurring_day"`
	// If trigger type is `daily`, this specifies one or more day numbers to run in one of the supported intervals: `week`, `month`, `quarter`, `year`.
	RecurringDays []int64 `url:"recurring_days,omitempty" json:"recurring_days,omitempty" path:"recurring_days"`
	// If trigger is `custom_schedule`, the reusable Schedule used instead of the sync's schedule fields.
	ScheduleId int64 `url:"schedule_id,omitempty" json:"schedule_id,omitempty" path:"schedule_id"`
	// If trigger is `custom_schedule`, Custom schedule description for when the sync should be run. 0-based days of the week. 0 is Sunday, 1 is Monday, etc.
	ScheduleDaysOfWeek []int64 `url:"schedule_days_of_week,omitempty" json:"schedule_days_of_week,omitempty" path:"schedule_days_of_week"`
	// Time zone for the schedule. If not set, times are interpreted as UTC.
	ScheduleTimeZone string `url:"schedule_time_zone,omitempty" json:"schedule_time_zone,omitempty" path:"schedule_time_zone"`
	// Times of day to run in HH:MM format. For `custom_schedule`, run at these times on specified days of week. For `daily`, run at these times on the scheduled interval date.
	ScheduleTimesOfDay []string `url:"schedule_times_of_day,omitempty" json:"schedule_times_of_day,omitempty" path:"schedule_times_of_day"`
	// Absolute source path for the sync.
	SrcPath string `url:"src_path,omitempty" json:"src_path,omitempty" path:"src_path"`
	// Remote server ID for the source (if remote).
	SrcRemoteServerId int64 `url:"src_remote_server_id,omitempty" json:"src_remote_server_id,omitempty" path:"src_remote_server_id"`
	// Frequency in minutes between syncs. If set, this value must be greater than or equal to the `remote_sync_interval` value for the site's plan. If left blank, the plan's `remote_sync_interval` will be used. This setting is only used if `trigger` is empty.
	SyncIntervalMinutes int64 `url:"sync_interval_minutes,omitempty" json:"sync_interval_minutes,omitempty" path:"sync_interval_minutes"`
	// Trigger type: daily, custom_schedule, or manual. Values: daily, custom_schedule, manual.
	Trigger SyncTriggerEnum `url:"trigger,omitempty" json:"trigger,omitempty" path:"trigger"`
	// Some MFT services request an empty file (known as a trigger file) to signal the sync is complete and they can begin further processing. If trigger_file is set, a zero-byte file will be sent at the end of the sync.
	TriggerFile string `url:"trigger_file,omitempty" json:"trigger_file,omitempty" path:"trigger_file"`
	// If true, the trigger file will be sent at the end of a successful sync even when no files were transferred.
	AlwaysWriteTriggerFile *bool `url:"always_write_trigger_file,omitempty" json:"always_write_trigger_file,omitempty" path:"always_write_trigger_file"`
	// Workspace ID this sync belongs to.
	WorkspaceId int64 `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
}

// SyncDryRunParams contains the request parameters for POST /syncs/{id}/dry_run.
//
// Dry Run Sync
type SyncDryRunParams struct {
	// Sync ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// SyncManualRunParams contains the request parameters for POST /syncs/{id}/manual_run.
//
// Manually Run Sync
type SyncManualRunParams struct {
	// Sync ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// SyncUpdateParams contains the request parameters for PATCH /syncs/{id}.
type SyncUpdateParams struct {
	// Sync ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// Delete empty folders after sync?
	DeleteEmptyFolders *bool `url:"delete_empty_folders,omitempty" json:"delete_empty_folders,omitempty" path:"delete_empty_folders"`
	// Description for this sync job.
	Description string `url:"description,omitempty" json:"description,omitempty" path:"description"`
	// Absolute destination path for the sync.
	DestPath string `url:"dest_path,omitempty" json:"dest_path,omitempty" path:"dest_path"`
	// Remote server ID for the destination (if remote).
	DestRemoteServerId int64 `url:"dest_remote_server_id,omitempty" json:"dest_remote_server_id,omitempty" path:"dest_remote_server_id"`
	// Is this sync disabled?
	Disabled *bool `url:"disabled,omitempty" json:"disabled,omitempty" path:"disabled"`
	// Array of glob patterns to exclude.
	ExcludePatterns []string `url:"exclude_patterns,omitempty" json:"exclude_patterns,omitempty" path:"exclude_patterns"`
	// Skip the sync if there is a formal, observed holiday for this region.
	HolidayRegion string `url:"holiday_region,omitempty" json:"holiday_region,omitempty" path:"holiday_region"`
	// Array of glob patterns to include.
	IncludePatterns []string `url:"include_patterns,omitempty" json:"include_patterns,omitempty" path:"include_patterns"`
	// If trigger is `daily`, this specifies how often to run this sync. One of: `day`, `week`, `week_end`, `month`, `month_end`, `quarter`, `quarter_end`, `year`, `year_end`.
	Interval string `url:"interval,omitempty" json:"interval,omitempty" path:"interval"`
	// Keep files after copying?
	KeepAfterCopy *bool `url:"keep_after_copy,omitempty" json:"keep_after_copy,omitempty" path:"keep_after_copy"`
	// Name for this sync job.
	Name string `url:"name,omitempty" json:"name,omitempty" path:"name"`
	// If trigger type is `daily`, this specifies a day number to run in one of the supported intervals: `week`, `month`, `quarter`, `year`.
	RecurringDay int64 `url:"recurring_day,omitempty" json:"recurring_day,omitempty" path:"recurring_day"`
	// If trigger type is `daily`, this specifies one or more day numbers to run in one of the supported intervals: `week`, `month`, `quarter`, `year`.
	RecurringDays []int64 `url:"recurring_days,omitempty" json:"recurring_days,omitempty" path:"recurring_days"`
	// If trigger is `custom_schedule`, the reusable Schedule used instead of the sync's schedule fields.
	ScheduleId int64 `url:"schedule_id,omitempty" json:"schedule_id,omitempty" path:"schedule_id"`
	// If trigger is `custom_schedule`, Custom schedule description for when the sync should be run. 0-based days of the week. 0 is Sunday, 1 is Monday, etc.
	ScheduleDaysOfWeek []int64 `url:"schedule_days_of_week,omitempty" json:"schedule_days_of_week,omitempty" path:"schedule_days_of_week"`
	// Time zone for the schedule. If not set, times are interpreted as UTC.
	ScheduleTimeZone string `url:"schedule_time_zone,omitempty" json:"schedule_time_zone,omitempty" path:"schedule_time_zone"`
	// Times of day to run in HH:MM format. For `custom_schedule`, run at these times on specified days of week. For `daily`, run at these times on the scheduled interval date.
	ScheduleTimesOfDay []string `url:"schedule_times_of_day,omitempty" json:"schedule_times_of_day,omitempty" path:"schedule_times_of_day"`
	// Absolute source path for the sync.
	SrcPath string `url:"src_path,omitempty" json:"src_path,omitempty" path:"src_path"`
	// Remote server ID for the source (if remote).
	SrcRemoteServerId int64 `url:"src_remote_server_id,omitempty" json:"src_remote_server_id,omitempty" path:"src_remote_server_id"`
	// Frequency in minutes between syncs. If set, this value must be greater than or equal to the `remote_sync_interval` value for the site's plan. If left blank, the plan's `remote_sync_interval` will be used. This setting is only used if `trigger` is empty.
	SyncIntervalMinutes int64 `url:"sync_interval_minutes,omitempty" json:"sync_interval_minutes,omitempty" path:"sync_interval_minutes"`
	// Trigger type: daily, custom_schedule, or manual. Values: daily, custom_schedule, manual.
	Trigger SyncTriggerEnum `url:"trigger,omitempty" json:"trigger,omitempty" path:"trigger"`
	// Some MFT services request an empty file (known as a trigger file) to signal the sync is complete and they can begin further processing. If trigger_file is set, a zero-byte file will be sent at the end of the sync.
	TriggerFile string `url:"trigger_file,omitempty" json:"trigger_file,omitempty" path:"trigger_file"`
	// If true, the trigger file will be sent at the end of a successful sync even when no files were transferred.
	AlwaysWriteTriggerFile *bool `url:"always_write_trigger_file,omitempty" json:"always_write_trigger_file,omitempty" path:"always_write_trigger_file"`
}

// SyncDeleteParams contains the request parameters for DELETE /syncs/{id}.
type SyncDeleteParams struct {
	// Sync ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (s *Sync) UnmarshalJSON(data []byte) error {
	type sync Sync
	var v sync
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*s = Sync(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (s *SyncCollection) UnmarshalJSON(data []byte) error {
	type syncs SyncCollection
	var v syncs
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*s = SyncCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (s *SyncCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*s))
	for i, v := range *s {
		ret[i] = v
	}

	return &ret
}
