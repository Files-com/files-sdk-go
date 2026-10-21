package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// HolidayCalendar is a Files.com API resource.
type HolidayCalendar struct {
	// Holiday Calendar ID. Set a scheduled resource's `holiday_region` to `custom_` followed by this ID to make it skip the days in this calendar.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Holiday Calendar name.
	Name string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	// Holiday rules for the calendar.
	Definition interface{} `json:"definition,omitempty" path:"definition,omitempty" url:"definition,omitempty"`
	// Creation time.
	CreatedAt *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	// Last update time.
	UpdatedAt *time.Time `json:"updated_at,omitempty" path:"updated_at,omitempty" url:"updated_at,omitempty"`
}

// Identifier returns the resource ID.
func (h HolidayCalendar) Identifier() interface{} {
	return h.Id
}

// HolidayCalendarCollection is a list of HolidayCalendar resources.
type HolidayCalendarCollection []HolidayCalendar

// HolidayCalendarListParams contains the request parameters for GET /holiday_calendars.
type HolidayCalendarListParams struct {
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are .
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// HolidayCalendarFindParams contains the request parameters for GET /holiday_calendars/{id}.
type HolidayCalendarFindParams struct {
	// Holiday Calendar ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// HolidayCalendarCreateParams contains the request parameters for POST /holiday_calendars.
type HolidayCalendarCreateParams struct {
	// Holiday rules for the calendar. Required.
	Definition interface{} `url:"definition" json:"definition" path:"definition"`
	// Holiday Calendar name. Required.
	Name string `url:"name" json:"name" path:"name"`
}

// HolidayCalendarUpdateParams contains the request parameters for PATCH /holiday_calendars/{id}.
type HolidayCalendarUpdateParams struct {
	// Holiday Calendar ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// Holiday rules for the calendar.
	Definition interface{} `url:"definition,omitempty" json:"definition,omitempty" path:"definition"`
	// Holiday Calendar name.
	Name string `url:"name,omitempty" json:"name,omitempty" path:"name"`
}

// HolidayCalendarDeleteParams contains the request parameters for DELETE /holiday_calendars/{id}.
type HolidayCalendarDeleteParams struct {
	// Holiday Calendar ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (h *HolidayCalendar) UnmarshalJSON(data []byte) error {
	type holidayCalendar HolidayCalendar
	var v holidayCalendar
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*h = HolidayCalendar(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (h *HolidayCalendarCollection) UnmarshalJSON(data []byte) error {
	type holidayCalendars HolidayCalendarCollection
	var v holidayCalendars
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*h = HolidayCalendarCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (h *HolidayCalendarCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*h))
	for i, v := range *h {
		ret[i] = v
	}

	return &ret
}
