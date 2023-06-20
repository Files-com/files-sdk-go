package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type FileMigration struct {
	Id         int64  `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	Path       string `json:"path,omitempty" path:"path,omitempty" url:"path,omitempty"`
	DestPath   string `json:"dest_path,omitempty" path:"dest_path,omitempty" url:"dest_path,omitempty"`
	FilesMoved int64  `json:"files_moved,omitempty" path:"files_moved,omitempty" url:"files_moved,omitempty"`
	FilesTotal int64  `json:"files_total,omitempty" path:"files_total,omitempty" url:"files_total,omitempty"`
	Operation  string `json:"operation,omitempty" path:"operation,omitempty" url:"operation,omitempty"`
	Region     string `json:"region,omitempty" path:"region,omitempty" url:"region,omitempty"`
	Status     string `json:"status,omitempty" path:"status,omitempty" url:"status,omitempty"`
	LogUrl     string `json:"log_url,omitempty" path:"log_url,omitempty" url:"log_url,omitempty"`
}

func (f FileMigration) Identifier() interface{} {
	return f.Id
}

type FileMigrationCollection []FileMigration

type FileMigrationFindParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
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
