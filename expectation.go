package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// Expectation is a Files.com API resource.
type Expectation struct {
	// Expectation ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Workspace ID. `0` means the default workspace.
	WorkspaceId int64 `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	// Expectation name.
	Name string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	// Expectation description.
	Description string `json:"description,omitempty" path:"description,omitempty" url:"description,omitempty"`
	// Path scope for the expectation. Supports workspace-relative presentation. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.
	Path string `json:"path,omitempty" path:"path,omitempty" url:"path,omitempty"`
	// Source glob used to select candidate files.
	Source string `json:"source,omitempty" path:"source,omitempty" url:"source,omitempty"`
	// Optional source exclusion glob.
	ExcludePattern string `json:"exclude_pattern,omitempty" path:"exclude_pattern,omitempty" url:"exclude_pattern,omitempty"`
	// If true, the expectation is disabled.
	Disabled *bool `json:"disabled,omitempty" path:"disabled,omitempty" url:"disabled,omitempty"`
	// Criteria schema version for this expectation.
	ExpectationsVersion int64 `json:"expectations_version,omitempty" path:"expectations_version,omitempty" url:"expectations_version,omitempty"`
	// How this expectation opens windows.
	Trigger string `json:"trigger,omitempty" path:"trigger,omitempty" url:"trigger,omitempty"`
	// If trigger is `daily`, this specifies how often to run the expectation.
	Interval string `json:"interval,omitempty" path:"interval,omitempty" url:"interval,omitempty"`
	// If trigger is `daily`, this selects the day number inside the chosen interval.
	RecurringDay int64 `json:"recurring_day,omitempty" path:"recurring_day,omitempty" url:"recurring_day,omitempty"`
	// If trigger is `daily`, this selects one or more day numbers inside a `week`, `month`, `quarter`, or `year` interval.
	RecurringDays []int64 `json:"recurring_days,omitempty" path:"recurring_days,omitempty" url:"recurring_days,omitempty"`
	// If trigger is `custom_schedule`, the reusable Schedule used instead of the Expectation's schedule fields.
	ScheduleId int64 `json:"schedule_id,omitempty" path:"schedule_id,omitempty" url:"schedule_id,omitempty"`
	// If trigger is `custom_schedule`, the 0-based weekdays used by the schedule.
	ScheduleDaysOfWeek []int64 `json:"schedule_days_of_week,omitempty" path:"schedule_days_of_week,omitempty" url:"schedule_days_of_week,omitempty"`
	// Times of day in HH:MM format for the Expectation schedule.
	ScheduleTimesOfDay []string `json:"schedule_times_of_day,omitempty" path:"schedule_times_of_day,omitempty" url:"schedule_times_of_day,omitempty"`
	// Time zone used by the Expectation schedule.
	ScheduleTimeZone string `json:"schedule_time_zone,omitempty" path:"schedule_time_zone,omitempty" url:"schedule_time_zone,omitempty"`
	// Optional holiday region used by the Expectation schedule.
	HolidayRegion string `json:"holiday_region,omitempty" path:"holiday_region,omitempty" url:"holiday_region,omitempty"`
	// How many seconds before the due boundary the window starts.
	LookbackInterval int64 `json:"lookback_interval,omitempty" path:"lookback_interval,omitempty" url:"lookback_interval,omitempty"`
	// How many seconds a schedule-driven window may remain eligible to close as late.
	LateAcceptanceInterval int64 `json:"late_acceptance_interval,omitempty" path:"late_acceptance_interval,omitempty" url:"late_acceptance_interval,omitempty"`
	// How many quiet seconds are required before final closure.
	InactivityInterval int64 `json:"inactivity_interval,omitempty" path:"inactivity_interval,omitempty" url:"inactivity_interval,omitempty"`
	// Hard-stop duration in seconds for unscheduled expectations.
	MaxOpenInterval int64 `json:"max_open_interval,omitempty" path:"max_open_interval,omitempty" url:"max_open_interval,omitempty"`
	// Versioned success criteria definition for the expectation. Criteria v2 supports optional FTS content validation.
	Criteria interface{} `json:"criteria,omitempty" path:"criteria,omitempty" url:"criteria,omitempty"`
	// Last time this expectation was evaluated.
	LastEvaluatedAt *time.Time `json:"last_evaluated_at,omitempty" path:"last_evaluated_at,omitempty" url:"last_evaluated_at,omitempty"`
	// Last time this expectation closed successfully.
	LastSuccessAt *time.Time `json:"last_success_at,omitempty" path:"last_success_at,omitempty" url:"last_success_at,omitempty"`
	// Last time this expectation closed with a failure result.
	LastFailureAt *time.Time `json:"last_failure_at,omitempty" path:"last_failure_at,omitempty" url:"last_failure_at,omitempty"`
	// Most recent terminal result for this expectation.
	LastResult string `json:"last_result,omitempty" path:"last_result,omitempty" url:"last_result,omitempty"`
	// Creation time.
	CreatedAt *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	// Last update time.
	UpdatedAt *time.Time `json:"updated_at,omitempty" path:"updated_at,omitempty" url:"updated_at,omitempty"`
}

// Identifier returns the resource ID.
func (e Expectation) Identifier() interface{} {
	return e.Id
}

// ExpectationCollection is a list of Expectation resources.
type ExpectationCollection []Expectation

// ExpectationTriggerEnum is a string value for trigger.
// Enum lists the values documented by the API.
type ExpectationTriggerEnum string

// String returns the API parameter value.
func (u ExpectationTriggerEnum) String() string {
	return string(u)
}

// Enum returns the documented values keyed by their API strings.
func (u ExpectationTriggerEnum) Enum() map[string]ExpectationTriggerEnum {
	return map[string]ExpectationTriggerEnum{
		"manual":          ExpectationTriggerEnum("manual"),
		"upload":          ExpectationTriggerEnum("upload"),
		"daily":           ExpectationTriggerEnum("daily"),
		"custom_schedule": ExpectationTriggerEnum("custom_schedule"),
	}
}

// ExpectationListParams contains the request parameters for GET /expectations.
type ExpectationListParams struct {
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id`, `name` or `disabled`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `disabled` and `workspace_id`. Valid field combinations are `[ workspace_id, disabled ]`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// ExpectationFindParams contains the request parameters for GET /expectations/{id}.
type ExpectationFindParams struct {
	// Expectation ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// ExpectationCreateParams contains the request parameters for POST /expectations.
type ExpectationCreateParams struct {
	// Expectation name.
	Name string `url:"name,omitempty" json:"name,omitempty" path:"name"`
	// Expectation description.
	Description string `url:"description,omitempty" json:"description,omitempty" path:"description"`
	// Path scope for the expectation. Supports workspace-relative presentation.
	Path string `url:"path,omitempty" json:"path,omitempty" path:"path"`
	// Source glob used to select candidate files.
	Source string `url:"source,omitempty" json:"source,omitempty" path:"source"`
	// Optional source exclusion glob.
	ExcludePattern string `url:"exclude_pattern,omitempty" json:"exclude_pattern,omitempty" path:"exclude_pattern"`
	// If true, the expectation is disabled.
	Disabled *bool `url:"disabled,omitempty" json:"disabled,omitempty" path:"disabled"`
	// How this expectation opens windows. Values: manual, upload, daily, custom_schedule.
	Trigger ExpectationTriggerEnum `url:"trigger,omitempty" json:"trigger,omitempty" path:"trigger"`
	// If trigger is `daily`, this specifies how often to run the expectation.
	Interval string `url:"interval,omitempty" json:"interval,omitempty" path:"interval"`
	// If trigger is `daily`, this selects the day number inside the chosen interval.
	RecurringDay int64 `url:"recurring_day,omitempty" json:"recurring_day,omitempty" path:"recurring_day"`
	// If trigger is `daily`, this selects one or more day numbers inside a `week`, `month`, `quarter`, or `year` interval.
	RecurringDays []int64 `url:"recurring_days,omitempty" json:"recurring_days,omitempty" path:"recurring_days"`
	// If trigger is `custom_schedule`, the reusable Schedule used instead of the Expectation's schedule fields.
	ScheduleId int64 `url:"schedule_id,omitempty" json:"schedule_id,omitempty" path:"schedule_id"`
	// If trigger is `custom_schedule`, the 0-based weekdays used by the schedule.
	ScheduleDaysOfWeek []int64 `url:"schedule_days_of_week,omitempty" json:"schedule_days_of_week,omitempty" path:"schedule_days_of_week"`
	// Times of day in HH:MM format for the Expectation schedule.
	ScheduleTimesOfDay []string `url:"schedule_times_of_day,omitempty" json:"schedule_times_of_day,omitempty" path:"schedule_times_of_day"`
	// Time zone used by the Expectation schedule.
	ScheduleTimeZone string `url:"schedule_time_zone,omitempty" json:"schedule_time_zone,omitempty" path:"schedule_time_zone"`
	// Optional holiday region used by the Expectation schedule.
	HolidayRegion string `url:"holiday_region,omitempty" json:"holiday_region,omitempty" path:"holiday_region"`
	// How many seconds before the due boundary the window starts.
	LookbackInterval int64 `url:"lookback_interval,omitempty" json:"lookback_interval,omitempty" path:"lookback_interval"`
	// How many seconds a schedule-driven window may remain eligible to close as late.
	LateAcceptanceInterval int64 `url:"late_acceptance_interval,omitempty" json:"late_acceptance_interval,omitempty" path:"late_acceptance_interval"`
	// How many quiet seconds are required before final closure.
	InactivityInterval int64 `url:"inactivity_interval,omitempty" json:"inactivity_interval,omitempty" path:"inactivity_interval"`
	// Hard-stop duration in seconds for unscheduled expectations.
	MaxOpenInterval int64 `url:"max_open_interval,omitempty" json:"max_open_interval,omitempty" path:"max_open_interval"`
	// Versioned success criteria definition for the expectation, including optional Files Transform Script content validation in criteria v2.
	Criteria interface{} `url:"criteria,omitempty" json:"criteria,omitempty" path:"criteria"`
	// Workspace ID. `0` means the default workspace.
	WorkspaceId int64 `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
}

// ExpectationTriggerEvaluationParams contains the request parameters for POST /expectations/{id}/trigger_evaluation.
//
// Manually open an Expectation window
type ExpectationTriggerEvaluationParams struct {
	// Expectation ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// ExpectationUpdateParams contains the request parameters for PATCH /expectations/{id}.
type ExpectationUpdateParams struct {
	// Expectation ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// Expectation name.
	Name string `url:"name,omitempty" json:"name,omitempty" path:"name"`
	// Expectation description.
	Description string `url:"description,omitempty" json:"description,omitempty" path:"description"`
	// Path scope for the expectation. Supports workspace-relative presentation.
	Path string `url:"path,omitempty" json:"path,omitempty" path:"path"`
	// Source glob used to select candidate files.
	Source string `url:"source,omitempty" json:"source,omitempty" path:"source"`
	// Optional source exclusion glob.
	ExcludePattern string `url:"exclude_pattern,omitempty" json:"exclude_pattern,omitempty" path:"exclude_pattern"`
	// If true, the expectation is disabled.
	Disabled *bool `url:"disabled,omitempty" json:"disabled,omitempty" path:"disabled"`
	// How this expectation opens windows. Values: manual, upload, daily, custom_schedule.
	Trigger ExpectationTriggerEnum `url:"trigger,omitempty" json:"trigger,omitempty" path:"trigger"`
	// If trigger is `daily`, this specifies how often to run the expectation.
	Interval string `url:"interval,omitempty" json:"interval,omitempty" path:"interval"`
	// If trigger is `daily`, this selects the day number inside the chosen interval.
	RecurringDay int64 `url:"recurring_day,omitempty" json:"recurring_day,omitempty" path:"recurring_day"`
	// If trigger is `daily`, this selects one or more day numbers inside a `week`, `month`, `quarter`, or `year` interval.
	RecurringDays []int64 `url:"recurring_days,omitempty" json:"recurring_days,omitempty" path:"recurring_days"`
	// If trigger is `custom_schedule`, the reusable Schedule used instead of the Expectation's schedule fields.
	ScheduleId int64 `url:"schedule_id,omitempty" json:"schedule_id,omitempty" path:"schedule_id"`
	// If trigger is `custom_schedule`, the 0-based weekdays used by the schedule.
	ScheduleDaysOfWeek []int64 `url:"schedule_days_of_week,omitempty" json:"schedule_days_of_week,omitempty" path:"schedule_days_of_week"`
	// Times of day in HH:MM format for the Expectation schedule.
	ScheduleTimesOfDay []string `url:"schedule_times_of_day,omitempty" json:"schedule_times_of_day,omitempty" path:"schedule_times_of_day"`
	// Time zone used by the Expectation schedule.
	ScheduleTimeZone string `url:"schedule_time_zone,omitempty" json:"schedule_time_zone,omitempty" path:"schedule_time_zone"`
	// Optional holiday region used by the Expectation schedule.
	HolidayRegion string `url:"holiday_region,omitempty" json:"holiday_region,omitempty" path:"holiday_region"`
	// How many seconds before the due boundary the window starts.
	LookbackInterval int64 `url:"lookback_interval,omitempty" json:"lookback_interval,omitempty" path:"lookback_interval"`
	// How many seconds a schedule-driven window may remain eligible to close as late.
	LateAcceptanceInterval int64 `url:"late_acceptance_interval,omitempty" json:"late_acceptance_interval,omitempty" path:"late_acceptance_interval"`
	// How many quiet seconds are required before final closure.
	InactivityInterval int64 `url:"inactivity_interval,omitempty" json:"inactivity_interval,omitempty" path:"inactivity_interval"`
	// Hard-stop duration in seconds for unscheduled expectations.
	MaxOpenInterval int64 `url:"max_open_interval,omitempty" json:"max_open_interval,omitempty" path:"max_open_interval"`
	// Versioned success criteria definition for the expectation, including optional Files Transform Script content validation in criteria v2.
	Criteria interface{} `url:"criteria,omitempty" json:"criteria,omitempty" path:"criteria"`
	// Workspace ID. `0` means the default workspace.
	WorkspaceId int64 `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
}

// ExpectationDeleteParams contains the request parameters for DELETE /expectations/{id}.
type ExpectationDeleteParams struct {
	// Expectation ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (e *Expectation) UnmarshalJSON(data []byte) error {
	type expectation Expectation
	var v expectation
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*e = Expectation(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (e *ExpectationCollection) UnmarshalJSON(data []byte) error {
	type expectations ExpectationCollection
	var v expectations
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*e = ExpectationCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (e *ExpectationCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*e))
	for i, v := range *e {
		ret[i] = v
	}

	return &ret
}
