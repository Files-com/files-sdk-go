package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

type Export struct {
	Id           int64  `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	ExportStatus string `json:"export_status,omitempty" path:"export_status,omitempty" url:"export_status,omitempty"`
	ExportType   string `json:"export_type,omitempty" path:"export_type,omitempty" url:"export_type,omitempty"`
	DownloadUri  string `json:"download_uri,omitempty" path:"download_uri,omitempty" url:"download_uri,omitempty"`
}

func (e Export) Identifier() interface{} {
	return e.Id
}

type ExportCollection []Export

type ExportListParams struct {
	UserId       int64                  `url:"user_id,omitempty" json:"user_id,omitempty" path:"user_id"`
	SortBy       map[string]interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	Filter       Export                 `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	FilterPrefix map[string]interface{} `url:"filter_prefix,omitempty" json:"filter_prefix,omitempty" path:"filter_prefix"`
	ListParams
}

type ExportFindParams struct {
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

func (e *Export) UnmarshalJSON(data []byte) error {
	type export Export
	var v export
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*e = Export(v)
	return nil
}

func (e *ExportCollection) UnmarshalJSON(data []byte) error {
	type exports ExportCollection
	var v exports
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*e = ExportCollection(v)
	return nil
}

func (e *ExportCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*e))
	for i, v := range *e {
		ret[i] = v
	}

	return &ret
}
