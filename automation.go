package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/lib"
)

type Automation struct {
	Id                     int64           `json:"id,omitempty"`
	Automation             string          `json:"automation,omitempty"`
	Source                 string          `json:"source,omitempty"`
	Destination            string          `json:"destination,omitempty"`
	DestinationReplaceFrom string          `json:"destination_replace_from,omitempty"`
	DestinationReplaceTo   string          `json:"destination_replace_to,omitempty"`
	Interval               string          `json:"interval,omitempty"`
	NextProcessOn          string          `json:"next_process_on,omitempty"`
	Path                   string          `json:"path,omitempty"`
	Realtime               *bool           `json:"realtime,omitempty"`
	UserId                 int64           `json:"user_id,omitempty"`
	UserIds                []string        `json:"user_ids,omitempty"`
	GroupIds               []string        `json:"group_ids,omitempty"`
	Trigger                string          `json:"trigger,omitempty"`
	Schedule               json.RawMessage `json:"schedule,omitempty"`
}

type AutomationCollection []Automation

type AutomationListParams struct {
	Cursor     string          `url:"cursor,omitempty" required:"false"`
	PerPage    int             `url:"per_page,omitempty" required:"false"`
	SortBy     json.RawMessage `url:"sort_by,omitempty" required:"false"`
	Filter     json.RawMessage `url:"filter,omitempty" required:"false"`
	FilterGt   json.RawMessage `url:"filter_gt,omitempty" required:"false"`
	FilterGteq json.RawMessage `url:"filter_gteq,omitempty" required:"false"`
	FilterLike json.RawMessage `url:"filter_like,omitempty" required:"false"`
	FilterLt   json.RawMessage `url:"filter_lt,omitempty" required:"false"`
	FilterLteq json.RawMessage `url:"filter_lteq,omitempty" required:"false"`
	Automation string          `url:"automation,omitempty" required:"false"`
	lib.ListParams
}

type AutomationFindParams struct {
	Id int64 `url:"-,omitempty" required:"true"`
}

type AutomationCreateParams struct {
	Automation             string          `url:"automation,omitempty" required:"true"`
	Source                 string          `url:"source,omitempty" required:"false"`
	Destination            string          `url:"destination,omitempty" required:"false"`
	DestinationReplaceFrom string          `url:"destination_replace_from,omitempty" required:"false"`
	DestinationReplaceTo   string          `url:"destination_replace_to,omitempty" required:"false"`
	Interval               string          `url:"interval,omitempty" required:"false"`
	Path                   string          `url:"path,omitempty" required:"false"`
	UserIds                string          `url:"user_ids,omitempty" required:"false"`
	GroupIds               string          `url:"group_ids,omitempty" required:"false"`
	Schedule               json.RawMessage `url:"schedule,omitempty" required:"false"`
	Trigger                string          `url:"trigger,omitempty" required:"false"`
}

type AutomationUpdateParams struct {
	Id                     int64           `url:"-,omitempty" required:"true"`
	Automation             string          `url:"automation,omitempty" required:"true"`
	Source                 string          `url:"source,omitempty" required:"false"`
	Destination            string          `url:"destination,omitempty" required:"false"`
	DestinationReplaceFrom string          `url:"destination_replace_from,omitempty" required:"false"`
	DestinationReplaceTo   string          `url:"destination_replace_to,omitempty" required:"false"`
	Interval               string          `url:"interval,omitempty" required:"false"`
	Path                   string          `url:"path,omitempty" required:"false"`
	UserIds                string          `url:"user_ids,omitempty" required:"false"`
	GroupIds               string          `url:"group_ids,omitempty" required:"false"`
	Schedule               json.RawMessage `url:"schedule,omitempty" required:"false"`
	Trigger                string          `url:"trigger,omitempty" required:"false"`
}

type AutomationDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"true"`
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
