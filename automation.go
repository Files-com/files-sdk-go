package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Automation struct {
	Id                     int64                  `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	Automation             string                 `json:"automation,omitempty" path:"automation,omitempty" url:"automation,omitempty"`
	Deleted                *bool                  `json:"deleted,omitempty" path:"deleted,omitempty" url:"deleted,omitempty"`
	Disabled               *bool                  `json:"disabled,omitempty" path:"disabled,omitempty" url:"disabled,omitempty"`
	Trigger                string                 `json:"trigger,omitempty" path:"trigger,omitempty" url:"trigger,omitempty"`
	Interval               string                 `json:"interval,omitempty" path:"interval,omitempty" url:"interval,omitempty"`
	LastModifiedAt         *time.Time             `json:"last_modified_at,omitempty" path:"last_modified_at,omitempty" url:"last_modified_at,omitempty"`
	Name                   string                 `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	Schedule               map[string]interface{} `json:"schedule,omitempty" path:"schedule,omitempty" url:"schedule,omitempty"`
	Source                 string                 `json:"source,omitempty" path:"source,omitempty" url:"source,omitempty"`
	Destinations           []string               `json:"destinations,omitempty" path:"destinations,omitempty" url:"destinations,omitempty"`
	DestinationReplaceFrom string                 `json:"destination_replace_from,omitempty" path:"destination_replace_from,omitempty" url:"destination_replace_from,omitempty"`
	DestinationReplaceTo   string                 `json:"destination_replace_to,omitempty" path:"destination_replace_to,omitempty" url:"destination_replace_to,omitempty"`
	Description            string                 `json:"description,omitempty" path:"description,omitempty" url:"description,omitempty"`
	RecurringDay           int64                  `json:"recurring_day,omitempty" path:"recurring_day,omitempty" url:"recurring_day,omitempty"`
	Path                   string                 `json:"path,omitempty" path:"path,omitempty" url:"path,omitempty"`
	UserId                 int64                  `json:"user_id,omitempty" path:"user_id,omitempty" url:"user_id,omitempty"`
	SyncIds                []int64                `json:"sync_ids,omitempty" path:"sync_ids,omitempty" url:"sync_ids,omitempty"`
	UserIds                []int64                `json:"user_ids,omitempty" path:"user_ids,omitempty" url:"user_ids,omitempty"`
	GroupIds               []int64                `json:"group_ids,omitempty" path:"group_ids,omitempty" url:"group_ids,omitempty"`
	WebhookUrl             string                 `json:"webhook_url,omitempty" path:"webhook_url,omitempty" url:"webhook_url,omitempty"`
	TriggerActions         []string               `json:"trigger_actions,omitempty" path:"trigger_actions,omitempty" url:"trigger_actions,omitempty"`
	Value                  map[string]interface{} `json:"value,omitempty" path:"value,omitempty" url:"value,omitempty"`
	Destination            string                 `json:"destination,omitempty" path:"destination,omitempty" url:"destination,omitempty"`
	ClonedFrom             int64                  `json:"cloned_from,omitempty" path:"cloned_from,omitempty" url:"cloned_from,omitempty"`
}

func (a Automation) Identifier() interface{} {
	return a.Id
}

type AutomationCollection []Automation

type AutomationTriggerEnum string

func (u AutomationTriggerEnum) String() string {
	return string(u)
}

func (u AutomationTriggerEnum) Enum() map[string]AutomationTriggerEnum {
	return map[string]AutomationTriggerEnum{
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
		"create_folder": AutomationEnum("create_folder"),
		"delete_file":   AutomationEnum("delete_file"),
		"copy_file":     AutomationEnum("copy_file"),
		"move_file":     AutomationEnum("move_file"),
		"as2_send":      AutomationEnum("as2_send"),
		"run_sync":      AutomationEnum("run_sync"),
	}
}

type AutomationListParams struct {
	Action      string                 `url:"action,omitempty" required:"false" json:"action,omitempty" path:"action"`
	SortBy      map[string]interface{} `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty" path:"sort_by"`
	Filter      Automation             `url:"filter,omitempty" required:"false" json:"filter,omitempty" path:"filter"`
	FilterGt    map[string]interface{} `url:"filter_gt,omitempty" required:"false" json:"filter_gt,omitempty" path:"filter_gt"`
	FilterGteq  map[string]interface{} `url:"filter_gteq,omitempty" required:"false" json:"filter_gteq,omitempty" path:"filter_gteq"`
	FilterLt    map[string]interface{} `url:"filter_lt,omitempty" required:"false" json:"filter_lt,omitempty" path:"filter_lt"`
	FilterLteq  map[string]interface{} `url:"filter_lteq,omitempty" required:"false" json:"filter_lteq,omitempty" path:"filter_lteq"`
	WithDeleted *bool                  `url:"with_deleted,omitempty" required:"false" json:"with_deleted,omitempty" path:"with_deleted"`
	ListParams
}

type AutomationFindParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
}

type AutomationCreateParams struct {
	Source                 string                 `url:"source,omitempty" required:"false" json:"source,omitempty" path:"source"`
	Destination            string                 `url:"destination,omitempty" required:"false" json:"destination,omitempty" path:"destination"`
	Destinations           []string               `url:"destinations,omitempty" required:"false" json:"destinations,omitempty" path:"destinations"`
	DestinationReplaceFrom string                 `url:"destination_replace_from,omitempty" required:"false" json:"destination_replace_from,omitempty" path:"destination_replace_from"`
	DestinationReplaceTo   string                 `url:"destination_replace_to,omitempty" required:"false" json:"destination_replace_to,omitempty" path:"destination_replace_to"`
	Interval               string                 `url:"interval,omitempty" required:"false" json:"interval,omitempty" path:"interval"`
	Path                   string                 `url:"path,omitempty" required:"false" json:"path,omitempty" path:"path"`
	SyncIds                string                 `url:"sync_ids,omitempty" required:"false" json:"sync_ids,omitempty" path:"sync_ids"`
	UserIds                string                 `url:"user_ids,omitempty" required:"false" json:"user_ids,omitempty" path:"user_ids"`
	GroupIds               string                 `url:"group_ids,omitempty" required:"false" json:"group_ids,omitempty" path:"group_ids"`
	Schedule               map[string]interface{} `url:"schedule,omitempty" required:"false" json:"schedule,omitempty" path:"schedule"`
	Description            string                 `url:"description,omitempty" required:"false" json:"description,omitempty" path:"description"`
	Disabled               *bool                  `url:"disabled,omitempty" required:"false" json:"disabled,omitempty" path:"disabled"`
	Name                   string                 `url:"name,omitempty" required:"false" json:"name,omitempty" path:"name"`
	Trigger                AutomationTriggerEnum  `url:"trigger,omitempty" required:"false" json:"trigger,omitempty" path:"trigger"`
	TriggerActions         []string               `url:"trigger_actions,omitempty" required:"false" json:"trigger_actions,omitempty" path:"trigger_actions"`
	Value                  map[string]interface{} `url:"value,omitempty" required:"false" json:"value,omitempty" path:"value"`
	RecurringDay           int64                  `url:"recurring_day,omitempty" required:"false" json:"recurring_day,omitempty" path:"recurring_day"`
	Automation             AutomationEnum         `url:"automation,omitempty" required:"true" json:"automation,omitempty" path:"automation"`
	ClonedFrom             int64                  `url:"cloned_from,omitempty" required:"false" json:"cloned_from,omitempty" path:"cloned_from"`
}

// Manually run automation
type AutomationManualRunParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
}

type AutomationUpdateParams struct {
	Id                     int64                  `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
	Source                 string                 `url:"source,omitempty" required:"false" json:"source,omitempty" path:"source"`
	Destination            string                 `url:"destination,omitempty" required:"false" json:"destination,omitempty" path:"destination"`
	Destinations           []string               `url:"destinations,omitempty" required:"false" json:"destinations,omitempty" path:"destinations"`
	DestinationReplaceFrom string                 `url:"destination_replace_from,omitempty" required:"false" json:"destination_replace_from,omitempty" path:"destination_replace_from"`
	DestinationReplaceTo   string                 `url:"destination_replace_to,omitempty" required:"false" json:"destination_replace_to,omitempty" path:"destination_replace_to"`
	Interval               string                 `url:"interval,omitempty" required:"false" json:"interval,omitempty" path:"interval"`
	Path                   string                 `url:"path,omitempty" required:"false" json:"path,omitempty" path:"path"`
	SyncIds                string                 `url:"sync_ids,omitempty" required:"false" json:"sync_ids,omitempty" path:"sync_ids"`
	UserIds                string                 `url:"user_ids,omitempty" required:"false" json:"user_ids,omitempty" path:"user_ids"`
	GroupIds               string                 `url:"group_ids,omitempty" required:"false" json:"group_ids,omitempty" path:"group_ids"`
	Schedule               map[string]interface{} `url:"schedule,omitempty" required:"false" json:"schedule,omitempty" path:"schedule"`
	Description            string                 `url:"description,omitempty" required:"false" json:"description,omitempty" path:"description"`
	Disabled               *bool                  `url:"disabled,omitempty" required:"false" json:"disabled,omitempty" path:"disabled"`
	Name                   string                 `url:"name,omitempty" required:"false" json:"name,omitempty" path:"name"`
	Trigger                AutomationTriggerEnum  `url:"trigger,omitempty" required:"false" json:"trigger,omitempty" path:"trigger"`
	TriggerActions         []string               `url:"trigger_actions,omitempty" required:"false" json:"trigger_actions,omitempty" path:"trigger_actions"`
	Value                  map[string]interface{} `url:"value,omitempty" required:"false" json:"value,omitempty" path:"value"`
	RecurringDay           int64                  `url:"recurring_day,omitempty" required:"false" json:"recurring_day,omitempty" path:"recurring_day"`
	Automation             AutomationEnum         `url:"automation,omitempty" required:"false" json:"automation,omitempty" path:"automation"`
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
