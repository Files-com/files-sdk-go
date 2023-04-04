package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Automation struct {
	Id                     int64           `json:"id,omitempty" path:"id"`
	Automation             string          `json:"automation,omitempty" path:"automation"`
	Deleted                *bool           `json:"deleted,omitempty" path:"deleted"`
	Disabled               *bool           `json:"disabled,omitempty" path:"disabled"`
	Trigger                string          `json:"trigger,omitempty" path:"trigger"`
	Interval               string          `json:"interval,omitempty" path:"interval"`
	LastModifiedAt         *time.Time      `json:"last_modified_at,omitempty" path:"last_modified_at"`
	Name                   string          `json:"name,omitempty" path:"name"`
	Schedule               json.RawMessage `json:"schedule,omitempty" path:"schedule"`
	Source                 string          `json:"source,omitempty" path:"source"`
	Destinations           []string        `json:"destinations,omitempty" path:"destinations"`
	DestinationReplaceFrom string          `json:"destination_replace_from,omitempty" path:"destination_replace_from"`
	DestinationReplaceTo   string          `json:"destination_replace_to,omitempty" path:"destination_replace_to"`
	Description            string          `json:"description,omitempty" path:"description"`
	RecurringDay           int64           `json:"recurring_day,omitempty" path:"recurring_day"`
	Path                   string          `json:"path,omitempty" path:"path"`
	UserId                 int64           `json:"user_id,omitempty" path:"user_id"`
	SyncIds                []int64         `json:"sync_ids,omitempty" path:"sync_ids"`
	UserIds                []int64         `json:"user_ids,omitempty" path:"user_ids"`
	GroupIds               []int64         `json:"group_ids,omitempty" path:"group_ids"`
	WebhookUrl             string          `json:"webhook_url,omitempty" path:"webhook_url"`
	TriggerActions         []string        `json:"trigger_actions,omitempty" path:"trigger_actions"`
	Value                  json.RawMessage `json:"value,omitempty" path:"value"`
	Destination            string          `json:"destination,omitempty" path:"destination"`
}

type AutomationCollection []Automation

type AutomationTriggerEnum string

func (u AutomationTriggerEnum) String() string {
	return string(u)
}

func (u AutomationTriggerEnum) Enum() map[string]AutomationTriggerEnum {
	return map[string]AutomationTriggerEnum{
		"realtime":        AutomationTriggerEnum("realtime"),
		"daily":           AutomationTriggerEnum("daily"),
		"custom_schedule": AutomationTriggerEnum("custom_schedule"),
		"webhook":         AutomationTriggerEnum("webhook"),
		"email":           AutomationTriggerEnum("email"),
		"action":          AutomationTriggerEnum("action"),
	}
}

type AutomationEnum string

func (u AutomationEnum) String() string {
	return string(u)
}

func (u AutomationEnum) Enum() map[string]AutomationEnum {
	return map[string]AutomationEnum{
		"create_folder":    AutomationEnum("create_folder"),
		"request_file":     AutomationEnum("request_file"),
		"request_move":     AutomationEnum("request_move"),
		"copy_newest_file": AutomationEnum("copy_newest_file"),
		"delete_file":      AutomationEnum("delete_file"),
		"copy_file":        AutomationEnum("copy_file"),
		"move_file":        AutomationEnum("move_file"),
		"as2_send":         AutomationEnum("as2_send"),
		"run_sync":         AutomationEnum("run_sync"),
	}
}

type AutomationListParams struct {
	SortBy      json.RawMessage `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty" path:"sort_by"`
	Automation  string          `url:"automation,omitempty" required:"false" json:"automation,omitempty" path:"automation"`
	Filter      json.RawMessage `url:"filter,omitempty" required:"false" json:"filter,omitempty" path:"filter"`
	FilterGt    json.RawMessage `url:"filter_gt,omitempty" required:"false" json:"filter_gt,omitempty" path:"filter_gt"`
	FilterGteq  json.RawMessage `url:"filter_gteq,omitempty" required:"false" json:"filter_gteq,omitempty" path:"filter_gteq"`
	FilterLt    json.RawMessage `url:"filter_lt,omitempty" required:"false" json:"filter_lt,omitempty" path:"filter_lt"`
	FilterLteq  json.RawMessage `url:"filter_lteq,omitempty" required:"false" json:"filter_lteq,omitempty" path:"filter_lteq"`
	WithDeleted *bool           `url:"with_deleted,omitempty" required:"false" json:"with_deleted,omitempty" path:"with_deleted"`
	lib.ListParams
}

type AutomationFindParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
}

type AutomationCreateParams struct {
	Source                 string                `url:"source,omitempty" required:"false" json:"source,omitempty" path:"source"`
	Destination            string                `url:"destination,omitempty" required:"false" json:"destination,omitempty" path:"destination"`
	Destinations           []string              `url:"destinations,omitempty" required:"false" json:"destinations,omitempty" path:"destinations"`
	DestinationReplaceFrom string                `url:"destination_replace_from,omitempty" required:"false" json:"destination_replace_from,omitempty" path:"destination_replace_from"`
	DestinationReplaceTo   string                `url:"destination_replace_to,omitempty" required:"false" json:"destination_replace_to,omitempty" path:"destination_replace_to"`
	Interval               string                `url:"interval,omitempty" required:"false" json:"interval,omitempty" path:"interval"`
	Path                   string                `url:"path,omitempty" required:"false" json:"path,omitempty" path:"path"`
	SyncIds                string                `url:"sync_ids,omitempty" required:"false" json:"sync_ids,omitempty" path:"sync_ids"`
	UserIds                string                `url:"user_ids,omitempty" required:"false" json:"user_ids,omitempty" path:"user_ids"`
	GroupIds               string                `url:"group_ids,omitempty" required:"false" json:"group_ids,omitempty" path:"group_ids"`
	Schedule               json.RawMessage       `url:"schedule,omitempty" required:"false" json:"schedule,omitempty" path:"schedule"`
	Description            string                `url:"description,omitempty" required:"false" json:"description,omitempty" path:"description"`
	Disabled               *bool                 `url:"disabled,omitempty" required:"false" json:"disabled,omitempty" path:"disabled"`
	Name                   string                `url:"name,omitempty" required:"false" json:"name,omitempty" path:"name"`
	Trigger                AutomationTriggerEnum `url:"trigger,omitempty" required:"false" json:"trigger,omitempty" path:"trigger"`
	TriggerActions         []string              `url:"trigger_actions,omitempty" required:"false" json:"trigger_actions,omitempty" path:"trigger_actions"`
	Value                  json.RawMessage       `url:"value,omitempty" required:"false" json:"value,omitempty" path:"value"`
	RecurringDay           int64                 `url:"recurring_day,omitempty" required:"false" json:"recurring_day,omitempty" path:"recurring_day"`
	Automation             AutomationEnum        `url:"automation,omitempty" required:"true" json:"automation,omitempty" path:"automation"`
}

type AutomationUpdateParams struct {
	Id                     int64                 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
	Source                 string                `url:"source,omitempty" required:"false" json:"source,omitempty" path:"source"`
	Destination            string                `url:"destination,omitempty" required:"false" json:"destination,omitempty" path:"destination"`
	Destinations           []string              `url:"destinations,omitempty" required:"false" json:"destinations,omitempty" path:"destinations"`
	DestinationReplaceFrom string                `url:"destination_replace_from,omitempty" required:"false" json:"destination_replace_from,omitempty" path:"destination_replace_from"`
	DestinationReplaceTo   string                `url:"destination_replace_to,omitempty" required:"false" json:"destination_replace_to,omitempty" path:"destination_replace_to"`
	Interval               string                `url:"interval,omitempty" required:"false" json:"interval,omitempty" path:"interval"`
	Path                   string                `url:"path,omitempty" required:"false" json:"path,omitempty" path:"path"`
	SyncIds                string                `url:"sync_ids,omitempty" required:"false" json:"sync_ids,omitempty" path:"sync_ids"`
	UserIds                string                `url:"user_ids,omitempty" required:"false" json:"user_ids,omitempty" path:"user_ids"`
	GroupIds               string                `url:"group_ids,omitempty" required:"false" json:"group_ids,omitempty" path:"group_ids"`
	Schedule               json.RawMessage       `url:"schedule,omitempty" required:"false" json:"schedule,omitempty" path:"schedule"`
	Description            string                `url:"description,omitempty" required:"false" json:"description,omitempty" path:"description"`
	Disabled               *bool                 `url:"disabled,omitempty" required:"false" json:"disabled,omitempty" path:"disabled"`
	Name                   string                `url:"name,omitempty" required:"false" json:"name,omitempty" path:"name"`
	Trigger                AutomationTriggerEnum `url:"trigger,omitempty" required:"false" json:"trigger,omitempty" path:"trigger"`
	TriggerActions         []string              `url:"trigger_actions,omitempty" required:"false" json:"trigger_actions,omitempty" path:"trigger_actions"`
	Value                  json.RawMessage       `url:"value,omitempty" required:"false" json:"value,omitempty" path:"value"`
	RecurringDay           int64                 `url:"recurring_day,omitempty" required:"false" json:"recurring_day,omitempty" path:"recurring_day"`
	Automation             AutomationEnum        `url:"automation,omitempty" required:"false" json:"automation,omitempty" path:"automation"`
}

type AutomationDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
}

func (a *Automation) UnmarshalJSON(data []byte) error {
	type automation Automation
	var v automation
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*a = Automation(v)
	return nil
}

func (a *AutomationCollection) UnmarshalJSON(data []byte) error {
	type automations AutomationCollection
	var v automations
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*a = AutomationCollection(v)
	return nil
}

func (a *AutomationCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*a))
	for i, v := range *a {
		ret[i] = v
	}

	return &ret
}
