package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type ZipDownloadFiles struct {
	Files  []string `json:"files,omitempty" path:"files,omitempty" url:"files,omitempty"`
	Cursor string   `json:"cursor,omitempty" path:"cursor,omitempty" url:"cursor,omitempty"`
}

// Identifier no path or id

type ZipDownloadFilesCollection []ZipDownloadFiles

func (z *ZipDownloadFiles) UnmarshalJSON(data []byte) error {
	type zipDownloadFiles ZipDownloadFiles
	var v zipDownloadFiles
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*z = ZipDownloadFiles(v)
	return nil
}

func (z *ZipDownloadFilesCollection) UnmarshalJSON(data []byte) error {
	type zipDownloadFiless ZipDownloadFilesCollection
	var v zipDownloadFiless
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*z = ZipDownloadFilesCollection(v)
	return nil
}

func (z *ZipDownloadFilesCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*z))
	for i, v := range *z {
		ret[i] = v
	}

	return &ret
}
