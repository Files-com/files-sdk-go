package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

type AutomationRun struct {
	Id                   int64      `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	AutomationId         int64      `json:"automation_id,omitempty" path:"automation_id,omitempty" url:"automation_id,omitempty"`
	CompletedAt          *time.Time `json:"completed_at,omitempty" path:"completed_at,omitempty" url:"completed_at,omitempty"`
	CreatedAt            *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	Runtime              string     `json:"runtime,omitempty" path:"runtime,omitempty" url:"runtime,omitempty"`
	Status               string     `json:"status,omitempty" path:"status,omitempty" url:"status,omitempty"`
	SuccessfulOperations int64      `json:"successful_operations,omitempty" path:"successful_operations,omitempty" url:"successful_operations,omitempty"`
	FailedOperations     int64      `json:"failed_operations,omitempty" path:"failed_operations,omitempty" url:"failed_operations,omitempty"`
	StatusMessagesUrl    string     `json:"status_messages_url,omitempty" path:"status_messages_url,omitempty" url:"status_messages_url,omitempty"`
}

func (a AutomationRun) Identifier() interface{} {
	return a.Id
}

type AutomationRunCollection []AutomationRun

type AutomationRunListParams struct {
	SortBy       map[string]interface{} `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty" path:"sort_by"`
	Filter       AutomationRun          `url:"filter,omitempty" required:"false" json:"filter,omitempty" path:"filter"`
	AutomationId int64                  `url:"automation_id,omitempty" required:"true" json:"automation_id,omitempty" path:"automation_id"`
	ListParams
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
