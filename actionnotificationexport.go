package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type ActionNotificationExport struct {
	Id                 int64      `json:"id,omitempty" path:"id"`
	ExportVersion      string     `json:"export_version,omitempty" path:"export_version"`
	StartAt            *time.Time `json:"start_at,omitempty" path:"start_at"`
	EndAt              *time.Time `json:"end_at,omitempty" path:"end_at"`
	Status             string     `json:"status,omitempty" path:"status"`
	QueryPath          string     `json:"query_path,omitempty" path:"query_path"`
	QueryFolder        string     `json:"query_folder,omitempty" path:"query_folder"`
	QueryMessage       string     `json:"query_message,omitempty" path:"query_message"`
	QueryRequestMethod string     `json:"query_request_method,omitempty" path:"query_request_method"`
	QueryRequestUrl    string     `json:"query_request_url,omitempty" path:"query_request_url"`
	QueryStatus        string     `json:"query_status,omitempty" path:"query_status"`
	QuerySuccess       *bool      `json:"query_success,omitempty" path:"query_success"`
	ResultsUrl         string     `json:"results_url,omitempty" path:"results_url"`
	UserId             int64      `json:"user_id,omitempty" path:"user_id"`
}

func (a ActionNotificationExport) Identifier() interface{} {
	return a.Id
}

type ActionNotificationExportCollection []ActionNotificationExport

type ActionNotificationExportFindParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
}

type ActionNotificationExportCreateParams struct {
	UserId             int64      `url:"user_id,omitempty" required:"false" json:"user_id,omitempty" path:"user_id"`
	StartAt            *time.Time `url:"start_at,omitempty" required:"false" json:"start_at,omitempty" path:"start_at"`
	EndAt              *time.Time `url:"end_at,omitempty" required:"false" json:"end_at,omitempty" path:"end_at"`
	QueryMessage       string     `url:"query_message,omitempty" required:"false" json:"query_message,omitempty" path:"query_message"`
	QueryRequestMethod string     `url:"query_request_method,omitempty" required:"false" json:"query_request_method,omitempty" path:"query_request_method"`
	QueryRequestUrl    string     `url:"query_request_url,omitempty" required:"false" json:"query_request_url,omitempty" path:"query_request_url"`
	QueryStatus        string     `url:"query_status,omitempty" required:"false" json:"query_status,omitempty" path:"query_status"`
	QuerySuccess       *bool      `url:"query_success,omitempty" required:"false" json:"query_success,omitempty" path:"query_success"`
	QueryPath          string     `url:"query_path,omitempty" required:"false" json:"query_path,omitempty" path:"query_path"`
	QueryFolder        string     `url:"query_folder,omitempty" required:"false" json:"query_folder,omitempty" path:"query_folder"`
}

func (a *ActionNotificationExport) UnmarshalJSON(data []byte) error {
	type actionNotificationExport ActionNotificationExport
	var v actionNotificationExport
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*a = ActionNotificationExport(v)
	return nil
}

func (a *ActionNotificationExportCollection) UnmarshalJSON(data []byte) error {
	type actionNotificationExports ActionNotificationExportCollection
	var v actionNotificationExports
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*a = ActionNotificationExportCollection(v)
	return nil
}

func (a *ActionNotificationExportCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*a))
	for i, v := range *a {
		ret[i] = v
	}

	return &ret
}
