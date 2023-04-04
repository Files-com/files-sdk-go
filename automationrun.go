package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type AutomationRun struct {
	Id                int64      `json:"id,omitempty" path:"id"`
	AutomationId      int64      `json:"automation_id,omitempty" path:"automation_id"`
	CompletedAt       *time.Time `json:"completed_at,omitempty" path:"completed_at"`
	CreatedAt         *time.Time `json:"created_at,omitempty" path:"created_at"`
	Status            string     `json:"status,omitempty" path:"status"`
	StatusMessagesUrl string     `json:"status_messages_url,omitempty" path:"status_messages_url"`
}

type AutomationRunCollection []AutomationRun

type AutomationRunListParams struct {
	UserId       int64           `url:"user_id,omitempty" required:"false" json:"user_id,omitempty" path:"user_id"`
	SortBy       json.RawMessage `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty" path:"sort_by"`
	Filter       json.RawMessage `url:"filter,omitempty" required:"false" json:"filter,omitempty" path:"filter"`
	AutomationId int64           `url:"automation_id,omitempty" required:"true" json:"automation_id,omitempty" path:"automation_id"`
	lib.ListParams
}

type AutomationRunFindParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
}

func (a *AutomationRun) UnmarshalJSON(data []byte) error {
	type automationRun AutomationRun
	var v automationRun
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*a = AutomationRun(v)
	return nil
}

func (a *AutomationRunCollection) UnmarshalJSON(data []byte) error {
	type automationRuns AutomationRunCollection
	var v automationRuns
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
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
