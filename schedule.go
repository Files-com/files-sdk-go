package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

type Schedule struct {
	Id                    int64      `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	Name                  string     `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	ScheduleDaysOfWeek    []int64    `json:"schedule_days_of_week,omitempty" path:"schedule_days_of_week,omitempty" url:"schedule_days_of_week,omitempty"`
	ScheduleTimesOfDay    []string   `json:"schedule_times_of_day,omitempty" path:"schedule_times_of_day,omitempty" url:"schedule_times_of_day,omitempty"`
	ScheduleTimeZone      string     `json:"schedule_time_zone,omitempty" path:"schedule_time_zone,omitempty" url:"schedule_time_zone,omitempty"`
	HolidayRegion         string     `json:"holiday_region,omitempty" path:"holiday_region,omitempty" url:"holiday_region,omitempty"`
	HumanReadableSchedule string     `json:"human_readable_schedule,omitempty" path:"human_readable_schedule,omitempty" url:"human_readable_schedule,omitempty"`
	CreatedAt             *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	UpdatedAt             *time.Time `json:"updated_at,omitempty" path:"updated_at,omitempty" url:"updated_at,omitempty"`
}

func (s Schedule) Identifier() interface{} {
	return s.Id
}

type ScheduleCollection []Schedule

type ScheduleListParams struct {
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	ListParams
}

type ScheduleFindParams struct {
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

type ScheduleCreateParams struct {
	Name               string   `url:"name" json:"name" path:"name"`
	ScheduleDaysOfWeek []int64  `url:"schedule_days_of_week" json:"schedule_days_of_week" path:"schedule_days_of_week"`
	ScheduleTimesOfDay []string `url:"schedule_times_of_day" json:"schedule_times_of_day" path:"schedule_times_of_day"`
	ScheduleTimeZone   string   `url:"schedule_time_zone,omitempty" json:"schedule_time_zone,omitempty" path:"schedule_time_zone"`
	HolidayRegion      string   `url:"holiday_region,omitempty" json:"holiday_region,omitempty" path:"holiday_region"`
}

type ScheduleUpdateParams struct {
	Id                 int64    `url:"-,omitempty" json:"-,omitempty" path:"id"`
	Name               string   `url:"name,omitempty" json:"name,omitempty" path:"name"`
	ScheduleDaysOfWeek []int64  `url:"schedule_days_of_week,omitempty" json:"schedule_days_of_week,omitempty" path:"schedule_days_of_week"`
	ScheduleTimesOfDay []string `url:"schedule_times_of_day,omitempty" json:"schedule_times_of_day,omitempty" path:"schedule_times_of_day"`
	ScheduleTimeZone   string   `url:"schedule_time_zone,omitempty" json:"schedule_time_zone,omitempty" path:"schedule_time_zone"`
	HolidayRegion      string   `url:"holiday_region,omitempty" json:"holiday_region,omitempty" path:"holiday_region"`
}

type ScheduleDeleteParams struct {
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

func (s *Schedule) UnmarshalJSON(data []byte) error {
	type schedule Schedule
	var v schedule
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*s = Schedule(v)
	return nil
}

func (s *ScheduleCollection) UnmarshalJSON(data []byte) error {
	type schedules ScheduleCollection
	var v schedules
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*s = ScheduleCollection(v)
	return nil
}

func (s *ScheduleCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*s))
	for i, v := range *s {
		ret[i] = v
	}

	return &ret
}
