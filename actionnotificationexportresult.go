package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type ActionNotificationExportResult struct {
	Id             int64  `json:"id,omitempty"`
	CreatedAt      int64  `json:"created_at,omitempty"`
	Status         int64  `json:"status,omitempty"`
	Message        string `json:"message,omitempty"`
	Success        *bool  `json:"success,omitempty"`
	RequestHeaders string `json:"request_headers,omitempty"`
	RequestMethod  string `json:"request_method,omitempty"`
	RequestUrl     string `json:"request_url,omitempty"`
	Path           string `json:"path,omitempty"`
	Folder         string `json:"folder,omitempty"`
}

type ActionNotificationExportResultCollection []ActionNotificationExportResult

type ActionNotificationExportResultListParams struct {
	UserId                     int64  `url:"user_id,omitempty" required:"false"`
	Cursor                     string `url:"cursor,omitempty" required:"false"`
	PerPage                    int64  `url:"per_page,omitempty" required:"false"`
	ActionNotificationExportId int64  `url:"action_notification_export_id,omitempty" required:"true"`
	lib.ListParams
}

func (a *ActionNotificationExportResult) UnmarshalJSON(data []byte) error {
	type actionNotificationExportResult ActionNotificationExportResult
	var v actionNotificationExportResult
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*a = ActionNotificationExportResult(v)
	return nil
}

func (a *ActionNotificationExportResultCollection) UnmarshalJSON(data []byte) error {
	type actionNotificationExportResults []ActionNotificationExportResult
	var v actionNotificationExportResults
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*a = ActionNotificationExportResultCollection(v)
	return nil
}

func (a *ActionNotificationExportResultCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*a))
	for i, v := range *a {
		ret[i] = v
	}

	return &ret
}
