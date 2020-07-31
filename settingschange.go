package files_sdk

import (
	"encoding/json"
	lib "github.com/Files-com/files-sdk-go/lib"
	"time"
)

type SettingsChange struct {
	ChangeDetails json.RawMessage `json:"change_details,omitempty"`
	CreatedAt     time.Time       `json:"created_at,omitempty"`
	UserId        int64           `json:"user_id,omitempty"`
}

type SettingsChangeCollection []SettingsChange

type SettingsChangeListParams struct {
	Page       int             `url:"page,omitempty"`
	PerPage    int             `url:"per_page,omitempty"`
	Action     string          `url:"action,omitempty"`
	Cursor     string          `url:"cursor,omitempty"`
	SortBy     json.RawMessage `url:"sort_by,omitempty"`
	Filter     json.RawMessage `url:"filter,omitempty"`
	FilterGt   json.RawMessage `url:"filter_gt,omitempty"`
	FilterGteq json.RawMessage `url:"filter_gteq,omitempty"`
	FilterLike json.RawMessage `url:"filter_like,omitempty"`
	FilterLt   json.RawMessage `url:"filter_lt,omitempty"`
	FilterLteq json.RawMessage `url:"filter_lteq,omitempty"`
	lib.ListParams
}

func (s *SettingsChange) UnmarshalJSON(data []byte) error {
	type settingsChange SettingsChange
	var v settingsChange
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*s = SettingsChange(v)
	return nil
}

func (s *SettingsChangeCollection) UnmarshalJSON(data []byte) error {
	type settingsChanges []SettingsChange
	var v settingsChanges
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*s = SettingsChangeCollection(v)
	return nil
}
