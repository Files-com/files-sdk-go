package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

type HolidayCalendar struct {
	Id         int64       `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	Name       string      `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	Definition interface{} `json:"definition,omitempty" path:"definition,omitempty" url:"definition,omitempty"`
	CreatedAt  *time.Time  `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	UpdatedAt  *time.Time  `json:"updated_at,omitempty" path:"updated_at,omitempty" url:"updated_at,omitempty"`
}

func (h HolidayCalendar) Identifier() interface{} {
	return h.Id
}

type HolidayCalendarCollection []HolidayCalendar

type HolidayCalendarListParams struct {
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	ListParams
}

type HolidayCalendarFindParams struct {
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

type HolidayCalendarCreateParams struct {
	Name string `url:"name" json:"name" path:"name"`
}

type HolidayCalendarUpdateParams struct {
	Id   int64  `url:"-,omitempty" json:"-,omitempty" path:"id"`
	Name string `url:"name,omitempty" json:"name,omitempty" path:"name"`
}

type HolidayCalendarDeleteParams struct {
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

func (h *HolidayCalendar) UnmarshalJSON(data []byte) error {
	type holidayCalendar HolidayCalendar
	var v holidayCalendar
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*h = HolidayCalendar(v)
	return nil
}

func (h *HolidayCalendarCollection) UnmarshalJSON(data []byte) error {
	type holidayCalendars HolidayCalendarCollection
	var v holidayCalendars
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*h = HolidayCalendarCollection(v)
	return nil
}

func (h *HolidayCalendarCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*h))
	for i, v := range *h {
		ret[i] = v
	}

	return &ret
}
