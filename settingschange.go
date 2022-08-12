package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type SettingsChange struct {
	Changes            []string   `json:"changes,omitempty" path:"changes"`
	CreatedAt          *time.Time `json:"created_at,omitempty" path:"created_at"`
	UserId             int64      `json:"user_id,omitempty" path:"user_id"`
	UserIsFilesSupport *bool      `json:"user_is_files_support,omitempty" path:"user_is_files_support"`
	Username           string     `json:"username,omitempty" path:"username"`
}

type SettingsChangeCollection []SettingsChange

type SettingsChangeListParams struct {
	SortBy     json.RawMessage `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty" path:"sort_by"`
	Filter     json.RawMessage `url:"filter,omitempty" required:"false" json:"filter,omitempty" path:"filter"`
	FilterGt   json.RawMessage `url:"filter_gt,omitempty" required:"false" json:"filter_gt,omitempty" path:"filter_gt"`
	FilterGteq json.RawMessage `url:"filter_gteq,omitempty" required:"false" json:"filter_gteq,omitempty" path:"filter_gteq"`
	FilterLike json.RawMessage `url:"filter_like,omitempty" required:"false" json:"filter_like,omitempty" path:"filter_like"`
	FilterLt   json.RawMessage `url:"filter_lt,omitempty" required:"false" json:"filter_lt,omitempty" path:"filter_lt"`
	FilterLteq json.RawMessage `url:"filter_lteq,omitempty" required:"false" json:"filter_lteq,omitempty" path:"filter_lteq"`
	lib.ListParams
}

func (s *SettingsChange) UnmarshalJSON(data []byte) error {
	type settingsChange SettingsChange
	var v settingsChange
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*s = SettingsChange(v)
	return nil
}

func (s *SettingsChangeCollection) UnmarshalJSON(data []byte) error {
	type settingsChanges SettingsChangeCollection
	var v settingsChanges
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
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
