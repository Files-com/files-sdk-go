package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type AutomationRun struct {
	Id                int64     `json:"id,omitempty"`
	AutomationId      int64     `json:"automation_id,omitempty"`
	CompletedAt       time.Time `json:"completed_at,omitempty"`
	CreatedAt         time.Time `json:"created_at,omitempty"`
	Status            string    `json:"status,omitempty"`
	StatusMessagesUrl string    `json:"status_messages_url,omitempty"`
}

type AutomationRunCollection []AutomationRun

type AutomationRunListParams struct {
	UserId       int64           `url:"user_id,omitempty" required:"false"`
	Cursor       string          `url:"cursor,omitempty" required:"false"`
	PerPage      int64           `url:"per_page,omitempty" required:"false"`
	SortBy       json.RawMessage `url:"sort_by,omitempty" required:"false"`
	Filter       json.RawMessage `url:"filter,omitempty" required:"false"`
	FilterGt     json.RawMessage `url:"filter_gt,omitempty" required:"false"`
	FilterGteq   json.RawMessage `url:"filter_gteq,omitempty" required:"false"`
	FilterLike   json.RawMessage `url:"filter_like,omitempty" required:"false"`
	FilterLt     json.RawMessage `url:"filter_lt,omitempty" required:"false"`
	FilterLteq   json.RawMessage `url:"filter_lteq,omitempty" required:"false"`
	AutomationId int64           `url:"automation_id,omitempty" required:"true"`
	lib.ListParams
}

type AutomationRunFindParams struct {
	Id int64 `url:"-,omitempty" required:"true"`
}

func (a *AutomationRun) UnmarshalJSON(data []byte) error {
	type automationRun AutomationRun
	var v automationRun
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*a = AutomationRun(v)
	return nil
}

func (a *AutomationRunCollection) UnmarshalJSON(data []byte) error {
	type automationRuns []AutomationRun
	var v automationRuns
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*a = AutomationRunCollection(v)
	return nil
}

func (a *AutomationRunCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*a))
	for i, v := range *a {
		ret[i] = v
	}

	return &ret
}
