package files_sdk

import (
	"encoding/json"
	"time"
)

type ActionNotificationExport struct {
	Id                 int64     `json:"id,omitempty"`
	ExportVersion      string    `json:"export_version,omitempty"`
	StartAt            time.Time `json:"start_at,omitempty"`
	EndAt              time.Time `json:"end_at,omitempty"`
	Status             string    `json:"status,omitempty"`
	QueryPath          string    `json:"query_path,omitempty"`
	QueryFolder        string    `json:"query_folder,omitempty"`
	QueryMessage       string    `json:"query_message,omitempty"`
	QueryRequestMethod string    `json:"query_request_method,omitempty"`
	QueryRequestUrl    string    `json:"query_request_url,omitempty"`
	QueryStatus        string    `json:"query_status,omitempty"`
	QuerySuccess       *bool     `json:"query_success,omitempty"`
	ResultsUrl         string    `json:"results_url,omitempty"`
	UserId             int64     `json:"user_id,omitempty"`
}

type ActionNotificationExportCollection []ActionNotificationExport

type ActionNotificationExportFindParams struct {
	Id int64 `url:"-,omitempty" required:"true"`
}

type ActionNotificationExportCreateParams struct {
	UserId             int64     `url:"user_id,omitempty" required:"false"`
	StartAt            time.Time `url:"start_at,omitempty" required:"false"`
	EndAt              time.Time `url:"end_at,omitempty" required:"false"`
	QueryMessage       string    `url:"query_message,omitempty" required:"false"`
	QueryRequestMethod string    `url:"query_request_method,omitempty" required:"false"`
	QueryRequestUrl    string    `url:"query_request_url,omitempty" required:"false"`
	QueryStatus        string    `url:"query_status,omitempty" required:"false"`
	QuerySuccess       *bool     `url:"query_success,omitempty" required:"false"`
	QueryPath          string    `url:"query_path,omitempty" required:"false"`
	QueryFolder        string    `url:"query_folder,omitempty" required:"false"`
}

func (a *ActionNotificationExport) UnmarshalJSON(data []byte) error {
	type actionNotificationExport ActionNotificationExport
	var v actionNotificationExport
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*a = ActionNotificationExport(v)
	return nil
}

func (a *ActionNotificationExportCollection) UnmarshalJSON(data []byte) error {
	type actionNotificationExports []ActionNotificationExport
	var v actionNotificationExports
	if err := json.Unmarshal(data, &v); err != nil {
		return err
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
