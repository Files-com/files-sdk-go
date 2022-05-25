package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Automation struct {
	Id                     int64           `json:"id,omitempty"`
	Automation             string          `json:"automation,omitempty"`
	Deleted                *bool           `json:"deleted,omitempty"`
	Disabled               *bool           `json:"disabled,omitempty"`
	Trigger                string          `json:"trigger,omitempty"`
	Interval               string          `json:"interval,omitempty"`
	LastModifiedAt         time.Time       `json:"last_modified_at,omitempty"`
	Name                   string          `json:"name,omitempty"`
	Schedule               json.RawMessage `json:"schedule,omitempty"`
	Source                 string          `json:"source,omitempty"`
	Destinations           string          `json:"destinations,omitempty"`
	DestinationReplaceFrom string          `json:"destination_replace_from,omitempty"`
	DestinationReplaceTo   string          `json:"destination_replace_to,omitempty"`
	Description            string          `json:"description,omitempty"`
	Path                   string          `json:"path,omitempty"`
	UserId                 int64           `json:"user_id,omitempty"`
	UserIds                []int64         `json:"user_ids,omitempty"`
	GroupIds               []int64         `json:"group_ids,omitempty"`
	WebhookUrl             string          `json:"webhook_url,omitempty"`
	TriggerActions         string          `json:"trigger_actions,omitempty"`
	Value                  json.RawMessage `json:"value,omitempty"`
	Destination            string          `json:"destination,omitempty"`
	ClonedFrom             int64           `json:"cloned_from,omitempty"`
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
	}
}

type AutomationListParams struct {
	Cursor      string          `url:"cursor,omitempty" required:"false" json:"cursor,omitempty"`
	PerPage     int64           `url:"per_page,omitempty" required:"false" json:"per_page,omitempty"`
	SortBy      json.RawMessage `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty"`
	Filter      json.RawMessage `url:"filter,omitempty" required:"false" json:"filter,omitempty"`
	FilterGt    json.RawMessage `url:"filter_gt,omitempty" required:"false" json:"filter_gt,omitempty"`
	FilterGteq  json.RawMessage `url:"filter_gteq,omitempty" required:"false" json:"filter_gteq,omitempty"`
	FilterLike  json.RawMessage `url:"filter_like,omitempty" required:"false" json:"filter_like,omitempty"`
	FilterLt    json.RawMessage `url:"filter_lt,omitempty" required:"false" json:"filter_lt,omitempty"`
	FilterLteq  json.RawMessage `url:"filter_lteq,omitempty" required:"false" json:"filter_lteq,omitempty"`
	WithDeleted *bool           `url:"with_deleted,omitempty" required:"false" json:"with_deleted,omitempty"`
	Automation  string          `url:"automation,omitempty" required:"false" json:"automation,omitempty"`
	lib.ListParams
}

type AutomationFindParams struct {
	Id int64 `url:"-,omitempty" required:"true" json:"-,omitempty"`
}

type AutomationCreateParams struct {
	Source                 string                `url:"source,omitempty" required:"false" json:"source,omitempty"`
	Destination            string                `url:"destination,omitempty" required:"false" json:"destination,omitempty"`
	Destinations           []string              `url:"destinations,omitempty" required:"false" json:"destinations,omitempty"`
	DestinationReplaceFrom string                `url:"destination_replace_from,omitempty" required:"false" json:"destination_replace_from,omitempty"`
	DestinationReplaceTo   string                `url:"destination_replace_to,omitempty" required:"false" json:"destination_replace_to,omitempty"`
	Interval               string                `url:"interval,omitempty" required:"false" json:"interval,omitempty"`
	Path                   string                `url:"path,omitempty" required:"false" json:"path,omitempty"`
	UserIds                string                `url:"user_ids,omitempty" required:"false" json:"user_ids,omitempty"`
	GroupIds               string                `url:"group_ids,omitempty" required:"false" json:"group_ids,omitempty"`
	Schedule               json.RawMessage       `url:"schedule,omitempty" required:"false" json:"schedule,omitempty"`
	Description            string                `url:"description,omitempty" required:"false" json:"description,omitempty"`
	Disabled               *bool                 `url:"disabled,omitempty" required:"false" json:"disabled,omitempty"`
	Name                   string                `url:"name,omitempty" required:"false" json:"name,omitempty"`
	Trigger                AutomationTriggerEnum `url:"trigger,omitempty" required:"false" json:"trigger,omitempty"`
	TriggerActions         []string              `url:"trigger_actions,omitempty" required:"false" json:"trigger_actions,omitempty"`
	Value                  json.RawMessage       `url:"value,omitempty" required:"false" json:"value,omitempty"`
	Automation             AutomationEnum        `url:"automation,omitempty" required:"true" json:"automation,omitempty"`
	ClonedFrom             int64                 `url:"cloned_from,omitempty" required:"false" json:"cloned_from,omitempty"`
}

type AutomationUpdateParams struct {
	Id                     int64                 `url:"-,omitempty" required:"true" json:"-,omitempty"`
	Source                 string                `url:"source,omitempty" required:"false" json:"source,omitempty"`
	Destination            string                `url:"destination,omitempty" required:"false" json:"destination,omitempty"`
	Destinations           []string              `url:"destinations,omitempty" required:"false" json:"destinations,omitempty"`
	DestinationReplaceFrom string                `url:"destination_replace_from,omitempty" required:"false" json:"destination_replace_from,omitempty"`
	DestinationReplaceTo   string                `url:"destination_replace_to,omitempty" required:"false" json:"destination_replace_to,omitempty"`
	Interval               string                `url:"interval,omitempty" required:"false" json:"interval,omitempty"`
	Path                   string                `url:"path,omitempty" required:"false" json:"path,omitempty"`
	UserIds                string                `url:"user_ids,omitempty" required:"false" json:"user_ids,omitempty"`
	GroupIds               string                `url:"group_ids,omitempty" required:"false" json:"group_ids,omitempty"`
	Schedule               json.RawMessage       `url:"schedule,omitempty" required:"false" json:"schedule,omitempty"`
	Description            string                `url:"description,omitempty" required:"false" json:"description,omitempty"`
	Disabled               *bool                 `url:"disabled,omitempty" required:"false" json:"disabled,omitempty"`
	Name                   string                `url:"name,omitempty" required:"false" json:"name,omitempty"`
	Trigger                AutomationTriggerEnum `url:"trigger,omitempty" required:"false" json:"trigger,omitempty"`
	TriggerActions         []string              `url:"trigger_actions,omitempty" required:"false" json:"trigger_actions,omitempty"`
	Value                  json.RawMessage       `url:"value,omitempty" required:"false" json:"value,omitempty"`
	Automation             AutomationEnum        `url:"automation,omitempty" required:"false" json:"automation,omitempty"`
}

type AutomationDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"true" json:"-,omitempty"`
}

func (a *Automation) UnmarshalJSON(data []byte) error {
	type automation Automation
	var v automation
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*a = Automation(v)
	return nil
}

func (a *AutomationCollection) UnmarshalJSON(data []byte) error {
	type automations []Automation
	var v automations
	if err := json.Unmarshal(data, &v); err != nil {
		return err
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
