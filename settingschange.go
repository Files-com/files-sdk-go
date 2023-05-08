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

// Identifier no path or id

type SettingsChangeCollection []SettingsChange

type SettingsChangeListParams struct {
	SortBy   json.RawMessage `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty" path:"sort_by"`
	ApiKeyId string          `url:"api_key_id,omitempty" required:"false" json:"api_key_id,omitempty" path:"api_key_id"`
	UserId   string          `url:"user_id,omitempty" required:"false" json:"user_id,omitempty" path:"user_id"`
	Filter   json.RawMessage `url:"filter,omitempty" required:"false" json:"filter,omitempty" path:"filter"`
	ListParams
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
