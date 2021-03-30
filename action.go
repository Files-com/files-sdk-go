package files_sdk

import (
	"encoding/json"
	"time"
)

type Action struct {
	Id          int64     `json:"id,omitempty"`
	Path        string    `json:"path,omitempty"`
	When        time.Time `json:"when,omitempty"`
	Destination string    `json:"destination,omitempty"`
	Display     string    `json:"display,omitempty"`
	Ip          string    `json:"ip,omitempty"`
	Source      string    `json:"source,omitempty"`
	Targets     []string  `json:"targets,omitempty"`
	UserId      int64     `json:"user_id,omitempty"`
	Username    string    `json:"username,omitempty"`
	Action      string    `json:"action,omitempty"`
	FailureType string    `json:"failure_type,omitempty"`
	Interface   string    `json:"interface,omitempty"`
}

type ActionCollection []Action

func (a *Action) UnmarshalJSON(data []byte) error {
	type action Action
	var v action
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*a = Action(v)
	return nil
}

func (a *ActionCollection) UnmarshalJSON(data []byte) error {
	type actions []Action
	var v actions
	if err := json.Unmarshal(data, &v); err != nil {
		return err
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
