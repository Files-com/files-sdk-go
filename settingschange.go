package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type SettingsChange struct {
	Changes            string    `json:"changes,omitempty"`
	CreatedAt          time.Time `json:"created_at,omitempty"`
	UserId             int64     `json:"user_id,omitempty"`
	UserIsFilesSupport *bool     `json:"user_is_files_support,omitempty"`
	Username           string    `json:"username,omitempty"`
}

type SettingsChangeCollection []SettingsChange

type SettingsChangeListParams struct {
	Cursor     string          `url:"cursor,omitempty" required:"false" json:"cursor,omitempty"`
	PerPage    int64           `url:"per_page,omitempty" required:"false" json:"per_page,omitempty"`
	SortBy     json.RawMessage `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty"`
	Filter     json.RawMessage `url:"filter,omitempty" required:"false" json:"filter,omitempty"`
	FilterGt   json.RawMessage `url:"filter_gt,omitempty" required:"false" json:"filter_gt,omitempty"`
	FilterGteq json.RawMessage `url:"filter_gteq,omitempty" required:"false" json:"filter_gteq,omitempty"`
	FilterLike json.RawMessage `url:"filter_like,omitempty" required:"false" json:"filter_like,omitempty"`
	FilterLt   json.RawMessage `url:"filter_lt,omitempty" required:"false" json:"filter_lt,omitempty"`
	FilterLteq json.RawMessage `url:"filter_lteq,omitempty" required:"false" json:"filter_lteq,omitempty"`
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

func (s *SettingsChangeCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*s))
	for i, v := range *s {
		ret[i] = v
	}

	return &ret
}
