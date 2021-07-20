package files_sdk

import (
	"encoding/json"
)

type FileMigration struct {
	Id         int64  `json:"id,omitempty"`
	Path       string `json:"path,omitempty"`
	DestPath   string `json:"dest_path,omitempty"`
	FilesMoved int64  `json:"files_moved,omitempty"`
	FilesTotal int64  `json:"files_total,omitempty"`
	Operation  string `json:"operation,omitempty"`
	Region     string `json:"region,omitempty"`
	Status     string `json:"status,omitempty"`
}

type FileMigrationCollection []FileMigration

type FileMigrationFindParams struct {
	Id int64 `url:"-,omitempty" required:"true"`
}

func (f *FileMigration) UnmarshalJSON(data []byte) error {
	type fileMigration FileMigration
	var v fileMigration
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*f = FileMigration(v)
	return nil
}

func (f *FileMigrationCollection) UnmarshalJSON(data []byte) error {
	type fileMigrations []FileMigration
	var v fileMigrations
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*f = FileMigrationCollection(v)
	return nil
}

func (f *FileMigrationCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*f))
	for i, v := range *f {
		ret[i] = v
	}

	return &ret
}
