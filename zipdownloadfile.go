package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type ZipDownloadFile struct {
	Files  []string `json:"files,omitempty" path:"files,omitempty" url:"files,omitempty"`
	Cursor string   `json:"cursor,omitempty" path:"cursor,omitempty" url:"cursor,omitempty"`
	Code   string   `json:"code,omitempty" path:"code,omitempty" url:"code,omitempty"`
	Limit  int64    `json:"limit,omitempty" path:"limit,omitempty" url:"limit,omitempty"`
	SiteId int64    `json:"site_id,omitempty" path:"site_id,omitempty" url:"site_id,omitempty"`
}

// Identifier no path or id

type ZipDownloadFileCollection []ZipDownloadFile

type ZipDownloadFileCreateParams struct {
	Code   string `url:"code,omitempty" required:"true" json:"code,omitempty" path:"code"`
	Limit  int64  `url:"limit,omitempty" required:"false" json:"limit,omitempty" path:"limit"`
	SiteId int64  `url:"site_id,omitempty" required:"false" json:"site_id,omitempty" path:"site_id"`
	ListParams
}

func (z *ZipDownloadFile) UnmarshalJSON(data []byte) error {
	type zipDownloadFile ZipDownloadFile
	var v zipDownloadFile
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*z = ZipDownloadFile(v)
	return nil
}

func (z *ZipDownloadFileCollection) UnmarshalJSON(data []byte) error {
	type zipDownloadFiles ZipDownloadFileCollection
	var v zipDownloadFiles
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*z = ZipDownloadFileCollection(v)
	return nil
}

func (z *ZipDownloadFileCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*z))
	for i, v := range *z {
		ret[i] = v
	}

	return &ret
}
