package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// AiTask is a Files.com API resource.
type AiTask struct {
	// AI Task ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Workspace ID. `0` means the default workspace.
	WorkspaceId int64 `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	// AI Task name.
	Name string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	// AI Task description.
	Description string `json:"description,omitempty" path:"description,omitempty" url:"description,omitempty"`
	// Prompt sent when this AI Task is invoked.
	Prompt string `json:"prompt,omitempty" path:"prompt,omitempty" url:"prompt,omitempty"`
	// Permissions used by the internal API key for this AI Task. Valid values are `full` and `files_only`.
	PermissionSet string `json:"permission_set,omitempty" path:"permission_set,omitempty" url:"permission_set,omitempty"`
	// Path scope used for action-triggered AI Tasks. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.
	Path string `json:"path,omitempty" path:"path,omitempty" url:"path,omitempty"`
	// Source glob used with `path` for action-triggered AI Tasks.
	Source string `json:"source,omitempty" path:"source,omitempty" url:"source,omitempty"`
	// If true, this AI Task will not run.
	Disabled *bool `json:"disabled,omitempty" path:"disabled,omitempty" url:"disabled,omitempty"`
	// How this AI Task is triggered.
	Trigger string `json:"trigger,omitempty" path:"trigger,omitempty" url:"trigger,omitempty"`
	// If trigger is `action`, the file action types that invoke this AI Task. Valid actions are create, copy, move, archived_delete, update, read, destroy.
	TriggerActions []string `json:"trigger_actions,omitempty" path:"trigger_actions,omitempty" url:"trigger_actions,omitempty"`
	// If trigger is `daily`, this specifies how often to run the AI Task.
	Interval string `json:"interval,omitempty" path:"interval,omitempty" url:"interval,omitempty"`
	// If trigger is `daily`, this selects the day number inside the chosen interval.
	RecurringDay int64 `json:"recurring_day,omitempty" path:"recurring_day,omitempty" url:"recurring_day,omitempty"`
	// If trigger is `daily`, this selects one or more day numbers inside a `week`, `month`, `quarter`, or `year` interval.
	RecurringDays []int64 `json:"recurring_days,omitempty" path:"recurring_days,omitempty" url:"recurring_days,omitempty"`
	// If trigger is `custom_schedule`, the reusable Schedule used instead of the AI Task's schedule fields.
	ScheduleId int64 `json:"schedule_id,omitempty" path:"schedule_id,omitempty" url:"schedule_id,omitempty"`
	// If trigger is `custom_schedule`, the 0-based weekdays used by the schedule.
	ScheduleDaysOfWeek []int64 `json:"schedule_days_of_week,omitempty" path:"schedule_days_of_week,omitempty" url:"schedule_days_of_week,omitempty"`
	// Times of day in HH:MM format for the AI Task schedule.
	ScheduleTimesOfDay []string `json:"schedule_times_of_day,omitempty" path:"schedule_times_of_day,omitempty" url:"schedule_times_of_day,omitempty"`
	// Time zone used by the AI Task schedule.
	ScheduleTimeZone string `json:"schedule_time_zone,omitempty" path:"schedule_time_zone,omitempty" url:"schedule_time_zone,omitempty"`
	// Optional holiday region used by the AI Task schedule.
	HolidayRegion string `json:"holiday_region,omitempty" path:"holiday_region,omitempty" url:"holiday_region,omitempty"`
	// Human-readable schedule description.
	HumanReadableSchedule string `json:"human_readable_schedule,omitempty" path:"human_readable_schedule,omitempty" url:"human_readable_schedule,omitempty"`
	// Most recent successful invocation time.
	LastRunAt *time.Time `json:"last_run_at,omitempty" path:"last_run_at,omitempty" url:"last_run_at,omitempty"`
	// Master User ID used for AI Task invocations.
	MasterAdminUserId int64 `json:"master_admin_user_id,omitempty" path:"master_admin_user_id,omitempty" url:"master_admin_user_id,omitempty"`
	// Creation time.
	CreatedAt *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	// Last update time.
	UpdatedAt *time.Time `json:"updated_at,omitempty" path:"updated_at,omitempty" url:"updated_at,omitempty"`
}

// Identifier returns the resource ID.
func (a AiTask) Identifier() interface{} {
	return a.Id
}

// AiTaskCollection is a list of AiTask resources.
type AiTaskCollection []AiTask

// AiTaskPermissionSetEnum is a string value for permission_set.
// Enum lists the values documented by the API.
type AiTaskPermissionSetEnum string

// String returns the API parameter value.
func (u AiTaskPermissionSetEnum) String() string {
	return string(u)
}

// Enum returns the documented values keyed by their API strings.
func (u AiTaskPermissionSetEnum) Enum() map[string]AiTaskPermissionSetEnum {
	return map[string]AiTaskPermissionSetEnum{
		"full":       AiTaskPermissionSetEnum("full"),
		"files_only": AiTaskPermissionSetEnum("files_only"),
	}
}

// AiTaskTriggerEnum is a string value for trigger.
// Enum lists the values documented by the API.
type AiTaskTriggerEnum string

// String returns the API parameter value.
func (u AiTaskTriggerEnum) String() string {
	return string(u)
}

// Enum returns the documented values keyed by their API strings.
func (u AiTaskTriggerEnum) Enum() map[string]AiTaskTriggerEnum {
	return map[string]AiTaskTriggerEnum{
		"manual":          AiTaskTriggerEnum("manual"),
		"daily":           AiTaskTriggerEnum("daily"),
		"custom_schedule": AiTaskTriggerEnum("custom_schedule"),
		"action":          AiTaskTriggerEnum("action"),
	}
}

// AiTaskListParams contains the request parameters for GET /ai_tasks.
type AiTaskListParams struct {
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id`, `name`, `disabled` or `updated_at`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `disabled`, `trigger` or `workspace_id`. Valid field combinations are `[ workspace_id, disabled ]`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// AiTaskFindParams contains the request parameters for GET /ai_tasks/{id}.
type AiTaskFindParams struct {
	// Ai Task ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// AiTaskCreateParams contains the request parameters for POST /ai_tasks.
type AiTaskCreateParams struct {
	// AI Task description.
	Description string `url:"description,omitempty" json:"description,omitempty" path:"description"`
	// If true, this AI Task will not run.
	Disabled *bool `url:"disabled,omitempty" json:"disabled,omitempty" path:"disabled"`
	// Optional holiday region used by the AI Task schedule.
	HolidayRegion string `url:"holiday_region,omitempty" json:"holiday_region,omitempty" path:"holiday_region"`
	// If trigger is `daily`, this specifies how often to run the AI Task.
	Interval string `url:"interval,omitempty" json:"interval,omitempty" path:"interval"`
	// AI Task name. Required.
	Name string `url:"name" json:"name" path:"name"`
	// Path scope used for action-triggered AI Tasks.
	Path string `url:"path,omitempty" json:"path,omitempty" path:"path"`
	// Permissions used by the internal API key for this AI Task. Valid values are `full` and `files_only`. Values: full, files_only.
	PermissionSet AiTaskPermissionSetEnum `url:"permission_set,omitempty" json:"permission_set,omitempty" path:"permission_set"`
	// Prompt sent when this AI Task is invoked. Required.
	Prompt string `url:"prompt" json:"prompt" path:"prompt"`
	// If trigger is `daily`, this selects the day number inside the chosen interval.
	RecurringDay int64 `url:"recurring_day,omitempty" json:"recurring_day,omitempty" path:"recurring_day"`
	// If trigger is `daily`, this selects one or more day numbers inside a `week`, `month`, `quarter`, or `year` interval.
	RecurringDays []int64 `url:"recurring_days,omitempty" json:"recurring_days,omitempty" path:"recurring_days"`
	// If trigger is `custom_schedule`, the reusable Schedule used instead of the AI Task's schedule fields.
	ScheduleId int64 `url:"schedule_id,omitempty" json:"schedule_id,omitempty" path:"schedule_id"`
	// If trigger is `custom_schedule`, the 0-based weekdays used by the schedule.
	ScheduleDaysOfWeek []int64 `url:"schedule_days_of_week,omitempty" json:"schedule_days_of_week,omitempty" path:"schedule_days_of_week"`
	// Time zone used by the AI Task schedule.
	ScheduleTimeZone string `url:"schedule_time_zone,omitempty" json:"schedule_time_zone,omitempty" path:"schedule_time_zone"`
	// Times of day in HH:MM format for the AI Task schedule.
	ScheduleTimesOfDay []string `url:"schedule_times_of_day,omitempty" json:"schedule_times_of_day,omitempty" path:"schedule_times_of_day"`
	// Source glob used with `path` for action-triggered AI Tasks.
	Source string `url:"source,omitempty" json:"source,omitempty" path:"source"`
	// How this AI Task is triggered. Values: manual, daily, custom_schedule, action.
	Trigger AiTaskTriggerEnum `url:"trigger,omitempty" json:"trigger,omitempty" path:"trigger"`
	// If trigger is `action`, the file action types that invoke this AI Task. Valid actions are create, copy, move, archived_delete, update, read, destroy.
	TriggerActions []string `url:"trigger_actions,omitempty" json:"trigger_actions,omitempty" path:"trigger_actions"`
	// Workspace ID. `0` means the default workspace.
	WorkspaceId int64 `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
}

// AiTaskManualRunParams contains the request parameters for POST /ai_tasks/{id}/manual_run.
//
// Manually Run AI Task
type AiTaskManualRunParams struct {
	// Ai Task ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// AiTaskUpdateParams contains the request parameters for PATCH /ai_tasks/{id}.
type AiTaskUpdateParams struct {
	// Ai Task ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// AI Task description.
	Description string `url:"description,omitempty" json:"description,omitempty" path:"description"`
	// If true, this AI Task will not run.
	Disabled *bool `url:"disabled,omitempty" json:"disabled,omitempty" path:"disabled"`
	// Optional holiday region used by the AI Task schedule.
	HolidayRegion string `url:"holiday_region,omitempty" json:"holiday_region,omitempty" path:"holiday_region"`
	// If trigger is `daily`, this specifies how often to run the AI Task.
	Interval string `url:"interval,omitempty" json:"interval,omitempty" path:"interval"`
	// AI Task name.
	Name string `url:"name,omitempty" json:"name,omitempty" path:"name"`
	// Path scope used for action-triggered AI Tasks.
	Path string `url:"path,omitempty" json:"path,omitempty" path:"path"`
	// Permissions used by the internal API key for this AI Task. Valid values are `full` and `files_only`. Values: full, files_only.
	PermissionSet AiTaskPermissionSetEnum `url:"permission_set,omitempty" json:"permission_set,omitempty" path:"permission_set"`
	// Prompt sent when this AI Task is invoked.
	Prompt string `url:"prompt,omitempty" json:"prompt,omitempty" path:"prompt"`
	// If trigger is `daily`, this selects the day number inside the chosen interval.
	RecurringDay int64 `url:"recurring_day,omitempty" json:"recurring_day,omitempty" path:"recurring_day"`
	// If trigger is `daily`, this selects one or more day numbers inside a `week`, `month`, `quarter`, or `year` interval.
	RecurringDays []int64 `url:"recurring_days,omitempty" json:"recurring_days,omitempty" path:"recurring_days"`
	// If trigger is `custom_schedule`, the reusable Schedule used instead of the AI Task's schedule fields.
	ScheduleId int64 `url:"schedule_id,omitempty" json:"schedule_id,omitempty" path:"schedule_id"`
	// If trigger is `custom_schedule`, the 0-based weekdays used by the schedule.
	ScheduleDaysOfWeek []int64 `url:"schedule_days_of_week,omitempty" json:"schedule_days_of_week,omitempty" path:"schedule_days_of_week"`
	// Time zone used by the AI Task schedule.
	ScheduleTimeZone string `url:"schedule_time_zone,omitempty" json:"schedule_time_zone,omitempty" path:"schedule_time_zone"`
	// Times of day in HH:MM format for the AI Task schedule.
	ScheduleTimesOfDay []string `url:"schedule_times_of_day,omitempty" json:"schedule_times_of_day,omitempty" path:"schedule_times_of_day"`
	// Source glob used with `path` for action-triggered AI Tasks.
	Source string `url:"source,omitempty" json:"source,omitempty" path:"source"`
	// How this AI Task is triggered. Values: manual, daily, custom_schedule, action.
	Trigger AiTaskTriggerEnum `url:"trigger,omitempty" json:"trigger,omitempty" path:"trigger"`
	// If trigger is `action`, the file action types that invoke this AI Task. Valid actions are create, copy, move, archived_delete, update, read, destroy.
	TriggerActions []string `url:"trigger_actions,omitempty" json:"trigger_actions,omitempty" path:"trigger_actions"`
	// Workspace ID. `0` means the default workspace.
	WorkspaceId int64 `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
}

// AiTaskDeleteParams contains the request parameters for DELETE /ai_tasks/{id}.
type AiTaskDeleteParams struct {
	// Ai Task ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (a *AiTask) UnmarshalJSON(data []byte) error {
	type aiTask AiTask
	var v aiTask
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*a = AiTask(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (a *AiTaskCollection) UnmarshalJSON(data []byte) error {
	type aiTasks AiTaskCollection
	var v aiTasks
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*a = AiTaskCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (a *AiTaskCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*a))
	for i, v := range *a {
		ret[i] = v
	}

	return &ret
}
