package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type AutomationRun struct {
	AutomationId      int64  `json:"automation_id,omitempty"`
	Status            string `json:"status,omitempty"`
	StatusMessagesUrl string `json:"status_messages_url,omitempty"`
}

type AutomationRunCollection []AutomationRun

type AutomationRunListParams struct {
	UserId       int64  `url:"user_id,omitempty" required:"false"`
	Cursor       string `url:"cursor,omitempty" required:"false"`
	PerPage      int64  `url:"per_page,omitempty" required:"false"`
	AutomationId int64  `url:"automation_id,omitempty" required:"true"`
	lib.ListParams
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
