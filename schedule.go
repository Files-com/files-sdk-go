package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// Schedule is a Files.com API resource.
type Schedule struct {
	// Schedule ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Schedule name.
	Name string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	// 0-based weekdays used by the Schedule. 0 is Sunday.
	ScheduleDaysOfWeek []int64 `json:"schedule_days_of_week,omitempty" path:"schedule_days_of_week,omitempty" url:"schedule_days_of_week,omitempty"`
	// Times of day in HH:MM format (24-hour).
	ScheduleTimesOfDay []string `json:"schedule_times_of_day,omitempty" path:"schedule_times_of_day,omitempty" url:"schedule_times_of_day,omitempty"`
	// Time zone for scheduled times. If not set, times are interpreted as UTC.
	ScheduleTimeZone string `json:"schedule_time_zone,omitempty" path:"schedule_time_zone,omitempty" url:"schedule_time_zone,omitempty"`
	// Optional holiday region on which linked resources do not run.
	HolidayRegion string `json:"holiday_region,omitempty" path:"holiday_region,omitempty" url:"holiday_region,omitempty"`
	// Human-readable Schedule description.
	HumanReadableSchedule string `json:"human_readable_schedule,omitempty" path:"human_readable_schedule,omitempty" url:"human_readable_schedule,omitempty"`
	// Creation time.
	CreatedAt *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	// Last update time.
	UpdatedAt *time.Time `json:"updated_at,omitempty" path:"updated_at,omitempty" url:"updated_at,omitempty"`
}

// Identifier returns the resource ID.
func (s Schedule) Identifier() interface{} {
	return s.Id
}

// ScheduleCollection is a list of Schedule resources.
type ScheduleCollection []Schedule

// ScheduleListParams contains the request parameters for GET /schedules.
type ScheduleListParams struct {
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `name`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// ScheduleFindParams contains the request parameters for GET /schedules/{id}.
type ScheduleFindParams struct {
	// Schedule ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// ScheduleCreateParams contains the request parameters for POST /schedules.
type ScheduleCreateParams struct {
	// Schedule name. Required.
	Name string `url:"name" json:"name" path:"name"`
	// 0-based weekdays used by the Schedule. 0 is Sunday. Required.
	ScheduleDaysOfWeek []int64 `url:"schedule_days_of_week" json:"schedule_days_of_week" path:"schedule_days_of_week"`
	// Times of day in HH:MM format (24-hour). Required.
	ScheduleTimesOfDay []string `url:"schedule_times_of_day" json:"schedule_times_of_day" path:"schedule_times_of_day"`
	// Time zone for scheduled times. If not set, times are interpreted as UTC.
	ScheduleTimeZone string `url:"schedule_time_zone,omitempty" json:"schedule_time_zone,omitempty" path:"schedule_time_zone"`
	// Optional holiday region on which linked resources do not run.
	HolidayRegion string `url:"holiday_region,omitempty" json:"holiday_region,omitempty" path:"holiday_region"`
}

// ScheduleUpdateParams contains the request parameters for PATCH /schedules/{id}.
type ScheduleUpdateParams struct {
	// Schedule ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// Schedule name.
	Name string `url:"name,omitempty" json:"name,omitempty" path:"name"`
	// 0-based weekdays used by the Schedule. 0 is Sunday.
	ScheduleDaysOfWeek []int64 `url:"schedule_days_of_week,omitempty" json:"schedule_days_of_week,omitempty" path:"schedule_days_of_week"`
	// Times of day in HH:MM format (24-hour).
	ScheduleTimesOfDay []string `url:"schedule_times_of_day,omitempty" json:"schedule_times_of_day,omitempty" path:"schedule_times_of_day"`
	// Time zone for scheduled times. If not set, times are interpreted as UTC.
	ScheduleTimeZone string `url:"schedule_time_zone,omitempty" json:"schedule_time_zone,omitempty" path:"schedule_time_zone"`
	// Optional holiday region on which linked resources do not run.
	HolidayRegion string `url:"holiday_region,omitempty" json:"holiday_region,omitempty" path:"holiday_region"`
}

// ScheduleDeleteParams contains the request parameters for DELETE /schedules/{id}.
type ScheduleDeleteParams struct {
	// Schedule ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (s *Schedule) UnmarshalJSON(data []byte) error {
	type schedule Schedule
	var v schedule
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*s = Schedule(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (s *ScheduleCollection) UnmarshalJSON(data []byte) error {
	type schedules ScheduleCollection
	var v schedules
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*s = ScheduleCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (s *ScheduleCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*s))
	for i, v := range *s {
		ret[i] = v
	}

	return &ret
}
