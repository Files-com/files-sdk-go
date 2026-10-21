package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// ScheduledExport is a Files.com API resource.
type ScheduledExport struct {
	// Scheduled Export ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Name for this scheduled export.
	Name string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	// Export report type. Valid values: folder_size_audit, group_membership_audit, permission_audit, share_link_audit.
	ExportType string `json:"export_type,omitempty" path:"export_type,omitempty" url:"export_type,omitempty"`
	// Human-readable report name.
	ReportName string `json:"report_name,omitempty" path:"report_name,omitempty" url:"report_name,omitempty"`
	// Report-specific options. `permission_audit` supports `group_by` with `user` or `path`.
	ExportOptions interface{} `json:"export_options,omitempty" path:"export_options,omitempty" url:"export_options,omitempty"`
	// Site Admin user who receives the completed export e-mail.
	UserId int64 `json:"user_id,omitempty" path:"user_id,omitempty" url:"user_id,omitempty"`
	// If true, this scheduled export will not run.
	Disabled *bool `json:"disabled,omitempty" path:"disabled,omitempty" url:"disabled,omitempty"`
	// Schedule trigger type: `daily` or `custom_schedule`.
	Trigger string `json:"trigger,omitempty" path:"trigger,omitempty" url:"trigger,omitempty"`
	// If trigger is `daily`, this specifies how often to run the scheduled export.
	Interval string `json:"interval,omitempty" path:"interval,omitempty" url:"interval,omitempty"`
	// If trigger is `daily`, this selects the day number inside the chosen interval.
	RecurringDay int64 `json:"recurring_day,omitempty" path:"recurring_day,omitempty" url:"recurring_day,omitempty"`
	// If trigger is `daily`, this selects one or more day numbers inside a `week`, `month`, `quarter`, or `year` interval.
	RecurringDays []int64 `json:"recurring_days,omitempty" path:"recurring_days,omitempty" url:"recurring_days,omitempty"`
	// If trigger is `custom_schedule`, the reusable Schedule used instead of the scheduled export's schedule fields.
	ScheduleId int64 `json:"schedule_id,omitempty" path:"schedule_id,omitempty" url:"schedule_id,omitempty"`
	// If trigger is `custom_schedule`, the 0-based weekdays used by the schedule.
	ScheduleDaysOfWeek []int64 `json:"schedule_days_of_week,omitempty" path:"schedule_days_of_week,omitempty" url:"schedule_days_of_week,omitempty"`
	// Times of day in HH:MM format for the scheduled export schedule.
	ScheduleTimesOfDay []string `json:"schedule_times_of_day,omitempty" path:"schedule_times_of_day,omitempty" url:"schedule_times_of_day,omitempty"`
	// Time zone used by the scheduled export schedule.
	ScheduleTimeZone string `json:"schedule_time_zone,omitempty" path:"schedule_time_zone,omitempty" url:"schedule_time_zone,omitempty"`
	// Optional holiday region used by the scheduled export schedule.
	HolidayRegion string `json:"holiday_region,omitempty" path:"holiday_region,omitempty" url:"holiday_region,omitempty"`
	// Human-readable schedule description.
	HumanReadableSchedule string `json:"human_readable_schedule,omitempty" path:"human_readable_schedule,omitempty" url:"human_readable_schedule,omitempty"`
	// Most recent scheduled attempt time, including attempts that failed validation.
	LastRunAt *time.Time `json:"last_run_at,omitempty" path:"last_run_at,omitempty" url:"last_run_at,omitempty"`
	// Most recent Export ID created by this schedule.
	LastExportId int64 `json:"last_export_id,omitempty" path:"last_export_id,omitempty" url:"last_export_id,omitempty"`
	// Validation error from the most recent scheduled attempt. The schedule remains enabled and retries at its next scheduled time. Cleared when an export is successfully created; does not describe errors during export generation.
	LastError string `json:"last_error,omitempty" path:"last_error,omitempty" url:"last_error,omitempty"`
	// Creation time.
	CreatedAt *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	// Last update time.
	UpdatedAt *time.Time `json:"updated_at,omitempty" path:"updated_at,omitempty" url:"updated_at,omitempty"`
}

// Identifier returns the resource ID.
func (s ScheduledExport) Identifier() interface{} {
	return s.Id
}

// ScheduledExportCollection is a list of ScheduledExport resources.
type ScheduledExportCollection []ScheduledExport

// ScheduledExportTriggerEnum is a string value for trigger.
// Enum lists the values documented by the API.
type ScheduledExportTriggerEnum string

// String returns the API parameter value.
func (u ScheduledExportTriggerEnum) String() string {
	return string(u)
}

// Enum returns the documented values keyed by their API strings.
func (u ScheduledExportTriggerEnum) Enum() map[string]ScheduledExportTriggerEnum {
	return map[string]ScheduledExportTriggerEnum{
		"daily":           ScheduledExportTriggerEnum("daily"),
		"custom_schedule": ScheduledExportTriggerEnum("custom_schedule"),
	}
}

// ScheduledExportListParams contains the request parameters for GET /scheduled_exports.
type ScheduledExportListParams struct {
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `name`, `export_type` or `disabled`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `disabled` and `export_type`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// If set, return records where the specified field is prefixed by the supplied value. Valid fields are `export_type`.
	FilterPrefix interface{} `url:"filter_prefix,omitempty" json:"filter_prefix,omitempty" path:"filter_prefix"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// ScheduledExportFindParams contains the request parameters for GET /scheduled_exports/{id}.
type ScheduledExportFindParams struct {
	// Scheduled Export ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// ScheduledExportCreateParams contains the request parameters for POST /scheduled_exports.
type ScheduledExportCreateParams struct {
	// Name for this scheduled export. Required.
	Name string `url:"name" json:"name" path:"name"`
	// Export report type. Valid values: folder_size_audit, group_membership_audit, permission_audit, share_link_audit. Required.
	ExportType string `url:"export_type" json:"export_type" path:"export_type"`
	// Report-specific options. `permission_audit` supports `group_by` with `user` or `path`.
	ExportOptions interface{} `url:"export_options,omitempty" json:"export_options,omitempty" path:"export_options"`
	// Site Admin user who receives the completed export e-mail.
	UserId int64 `url:"user_id,omitempty" json:"user_id,omitempty" path:"user_id"`
	// If true, this scheduled export will not run.
	Disabled *bool `url:"disabled,omitempty" json:"disabled,omitempty" path:"disabled"`
	// Schedule trigger type: `daily` or `custom_schedule`. Values: daily, custom_schedule.
	Trigger ScheduledExportTriggerEnum `url:"trigger,omitempty" json:"trigger,omitempty" path:"trigger"`
	// If trigger is `daily`, this specifies how often to run the scheduled export.
	Interval string `url:"interval,omitempty" json:"interval,omitempty" path:"interval"`
	// If trigger is `daily`, this selects the day number inside the chosen interval.
	RecurringDay int64 `url:"recurring_day,omitempty" json:"recurring_day,omitempty" path:"recurring_day"`
	// If trigger is `daily`, this selects one or more day numbers inside a `week`, `month`, `quarter`, or `year` interval.
	RecurringDays []int64 `url:"recurring_days,omitempty" json:"recurring_days,omitempty" path:"recurring_days"`
	// If trigger is `custom_schedule`, the reusable Schedule used instead of the scheduled export's schedule fields.
	ScheduleId int64 `url:"schedule_id,omitempty" json:"schedule_id,omitempty" path:"schedule_id"`
	// If trigger is `custom_schedule`, the 0-based weekdays used by the schedule.
	ScheduleDaysOfWeek []int64 `url:"schedule_days_of_week,omitempty" json:"schedule_days_of_week,omitempty" path:"schedule_days_of_week"`
	// Times of day in HH:MM format for the scheduled export schedule.
	ScheduleTimesOfDay []string `url:"schedule_times_of_day,omitempty" json:"schedule_times_of_day,omitempty" path:"schedule_times_of_day"`
	// Time zone used by the scheduled export schedule.
	ScheduleTimeZone string `url:"schedule_time_zone,omitempty" json:"schedule_time_zone,omitempty" path:"schedule_time_zone"`
	// Optional holiday region used by the scheduled export schedule.
	HolidayRegion string `url:"holiday_region,omitempty" json:"holiday_region,omitempty" path:"holiday_region"`
}

// ScheduledExportUpdateParams contains the request parameters for PATCH /scheduled_exports/{id}.
type ScheduledExportUpdateParams struct {
	// Scheduled Export ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// Name for this scheduled export.
	Name string `url:"name,omitempty" json:"name,omitempty" path:"name"`
	// Export report type. Valid values: folder_size_audit, group_membership_audit, permission_audit, share_link_audit.
	ExportType string `url:"export_type,omitempty" json:"export_type,omitempty" path:"export_type"`
	// Report-specific options. `permission_audit` supports `group_by` with `user` or `path`.
	ExportOptions interface{} `url:"export_options,omitempty" json:"export_options,omitempty" path:"export_options"`
	// Site Admin user who receives the completed export e-mail.
	UserId int64 `url:"user_id,omitempty" json:"user_id,omitempty" path:"user_id"`
	// If true, this scheduled export will not run.
	Disabled *bool `url:"disabled,omitempty" json:"disabled,omitempty" path:"disabled"`
	// Schedule trigger type: `daily` or `custom_schedule`. Values: daily, custom_schedule.
	Trigger ScheduledExportTriggerEnum `url:"trigger,omitempty" json:"trigger,omitempty" path:"trigger"`
	// If trigger is `daily`, this specifies how often to run the scheduled export.
	Interval string `url:"interval,omitempty" json:"interval,omitempty" path:"interval"`
	// If trigger is `daily`, this selects the day number inside the chosen interval.
	RecurringDay int64 `url:"recurring_day,omitempty" json:"recurring_day,omitempty" path:"recurring_day"`
	// If trigger is `daily`, this selects one or more day numbers inside a `week`, `month`, `quarter`, or `year` interval.
	RecurringDays []int64 `url:"recurring_days,omitempty" json:"recurring_days,omitempty" path:"recurring_days"`
	// If trigger is `custom_schedule`, the reusable Schedule used instead of the scheduled export's schedule fields.
	ScheduleId int64 `url:"schedule_id,omitempty" json:"schedule_id,omitempty" path:"schedule_id"`
	// If trigger is `custom_schedule`, the 0-based weekdays used by the schedule.
	ScheduleDaysOfWeek []int64 `url:"schedule_days_of_week,omitempty" json:"schedule_days_of_week,omitempty" path:"schedule_days_of_week"`
	// Times of day in HH:MM format for the scheduled export schedule.
	ScheduleTimesOfDay []string `url:"schedule_times_of_day,omitempty" json:"schedule_times_of_day,omitempty" path:"schedule_times_of_day"`
	// Time zone used by the scheduled export schedule.
	ScheduleTimeZone string `url:"schedule_time_zone,omitempty" json:"schedule_time_zone,omitempty" path:"schedule_time_zone"`
	// Optional holiday region used by the scheduled export schedule.
	HolidayRegion string `url:"holiday_region,omitempty" json:"holiday_region,omitempty" path:"holiday_region"`
}

// ScheduledExportDeleteParams contains the request parameters for DELETE /scheduled_exports/{id}.
type ScheduledExportDeleteParams struct {
	// Scheduled Export ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (s *ScheduledExport) UnmarshalJSON(data []byte) error {
	type scheduledExport ScheduledExport
	var v scheduledExport
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*s = ScheduledExport(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (s *ScheduledExportCollection) UnmarshalJSON(data []byte) error {
	type scheduledExports ScheduledExportCollection
	var v scheduledExports
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*s = ScheduledExportCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (s *ScheduledExportCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*s))
	for i, v := range *s {
		ret[i] = v
	}

	return &ret
}
