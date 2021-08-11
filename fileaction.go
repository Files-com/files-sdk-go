package files_sdk

import (
	"encoding/json"
)

type FileAction struct {
	Status          string `json:"status,omitempty"`
	FileMigrationId int64  `json:"file_migration_id,omitempty"`
}

type FileActionCollection []FileAction

func (f *FileAction) UnmarshalJSON(data []byte) error {
	type fileAction FileAction
	var v fileAction
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*f = FileAction(v)
	return nil
}

func (f *FileActionCollection) UnmarshalJSON(data []byte) error {
	type fileActions []FileAction
	var v fileActions
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*f = FileActionCollection(v)
	return nil
}

func (f *FileActionCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*f))
	for i, v := range *f {
		ret[i] = v
	}

	return &ret
}
