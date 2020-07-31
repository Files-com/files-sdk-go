package files_sdk

import (
  lib "github.com/Files-com/files-sdk-go/lib"
  "encoding/json"
)

type Automation struct {
  Id int `json:"id,omitempty"`
  Automation string `json:"automation,omitempty"`
  Source string `json:"source,omitempty"`
  Destination string `json:"destination,omitempty"`
  DestinationReplaceFrom string `json:"destination_replace_from,omitempty"`
  DestinationReplaceTo string `json:"destination_replace_to,omitempty"`
  Interval string `json:"interval,omitempty"`
  NextProcessOn string `json:"next_process_on,omitempty"`
  Path string `json:"path,omitempty"`
  Realtime *bool `json:"realtime,omitempty"`
  UserId int `json:"user_id,omitempty"`
  UserIds []string `json:"user_ids,omitempty"`
  GroupIds []string `json:"group_ids,omitempty"`
}

type AutomationCollection []Automation

type AutomationListParams struct {
  Page int `url:"page,omitempty"`
  PerPage int `url:"per_page,omitempty"`
  Action string `url:"action,omitempty"`
  Cursor string `url:"cursor,omitempty"`
  SortBy json.RawMessage `url:"sort_by,omitempty"`
  Filter json.RawMessage `url:"filter,omitempty"`
  FilterGt json.RawMessage `url:"filter_gt,omitempty"`
  FilterGteq json.RawMessage `url:"filter_gteq,omitempty"`
  FilterLike json.RawMessage `url:"filter_like,omitempty"`
  FilterLt json.RawMessage `url:"filter_lt,omitempty"`
  FilterLteq json.RawMessage `url:"filter_lteq,omitempty"`
  Automation string `url:"automation,omitempty"`
  lib.ListParams
}

type AutomationFindParams struct {
  Id int `url:"-,omitempty"`
}

type AutomationCreateParams struct {
  Automation string `url:"automation,omitempty"`
  Source string `url:"source,omitempty"`
  Destination string `url:"destination,omitempty"`
  DestinationReplaceFrom string `url:"destination_replace_from,omitempty"`
  DestinationReplaceTo string `url:"destination_replace_to,omitempty"`
  Interval string `url:"interval,omitempty"`
  Path string `url:"path,omitempty"`
  UserIds string `url:"user_ids,omitempty"`
  GroupIds string `url:"group_ids,omitempty"`
}

type AutomationUpdateParams struct {
  Id int `url:"-,omitempty"`
  Automation string `url:"automation,omitempty"`
  Source string `url:"source,omitempty"`
  Destination string `url:"destination,omitempty"`
  DestinationReplaceFrom string `url:"destination_replace_from,omitempty"`
  DestinationReplaceTo string `url:"destination_replace_to,omitempty"`
  Interval string `url:"interval,omitempty"`
  Path string `url:"path,omitempty"`
  UserIds string `url:"user_ids,omitempty"`
  GroupIds string `url:"group_ids,omitempty"`
}

type AutomationDeleteParams struct {
  Id int `url:"-,omitempty"`
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

