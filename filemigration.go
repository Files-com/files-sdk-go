package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type FileMigration struct {
	Id         int64  `json:"id,omitempty" path:"id"`
	Path       string `json:"path,omitempty" path:"path"`
	DestPath   string `json:"dest_path,omitempty" path:"dest_path"`
	FilesMoved int64  `json:"files_moved,omitempty" path:"files_moved"`
	FilesTotal int64  `json:"files_total,omitempty" path:"files_total"`
	Operation  string `json:"operation,omitempty" path:"operation"`
	Region     string `json:"region,omitempty" path:"region"`
	Status     string `json:"status,omitempty" path:"status"`
	LogUrl     string `json:"log_url,omitempty" path:"log_url"`
}

type FileMigrationCollection []FileMigration

type FileMigrationFindParams struct {
	Id int64 `url:"-,omitempty" required:"true" json:"-,omitempty" path:"id"`
}

func (f *FileMigration) UnmarshalJSON(data []byte) error {
	type fileMigration FileMigration
	var v fileMigration
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*f = FileMigration(v)
	return nil
}

func (f *FileMigrationCollection) UnmarshalJSON(data []byte) error {
	type fileMigrations FileMigrationCollection
	var v fileMigrations
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
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
