package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Action struct {
	Id          int64      `json:"id,omitempty" path:"id"`
	Path        string     `json:"path,omitempty" path:"path"`
	When        *time.Time `json:"when,omitempty" path:"when"`
	Destination string     `json:"destination,omitempty" path:"destination"`
	Display     string     `json:"display,omitempty" path:"display"`
	Ip          string     `json:"ip,omitempty" path:"ip"`
	Source      string     `json:"source,omitempty" path:"source"`
	Targets     []string   `json:"targets,omitempty" path:"targets"`
	UserId      int64      `json:"user_id,omitempty" path:"user_id"`
	Username    string     `json:"username,omitempty" path:"username"`
	Action      string     `json:"action,omitempty" path:"action"`
	FailureType string     `json:"failure_type,omitempty" path:"failure_type"`
	Interface   string     `json:"interface,omitempty" path:"interface"`
}

type ActionCollection []Action

func (a *Action) UnmarshalJSON(data []byte) error {
	type action Action
	var v action
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*a = Action(v)
	return nil
}

func (a *ActionCollection) UnmarshalJSON(data []byte) error {
	type actions ActionCollection
	var v actions
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*a = ActionCollection(v)
	return nil
}

func (a *ActionCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*a))
	for i, v := range *a {
		ret[i] = v
	}

	return &ret
}
