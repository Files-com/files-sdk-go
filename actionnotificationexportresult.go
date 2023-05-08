package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type ActionNotificationExportResult struct {
	Id             int64  `json:"id,omitempty" path:"id"`
	CreatedAt      int64  `json:"created_at,omitempty" path:"created_at"`
	Status         int64  `json:"status,omitempty" path:"status"`
	Message        string `json:"message,omitempty" path:"message"`
	Success        *bool  `json:"success,omitempty" path:"success"`
	RequestHeaders string `json:"request_headers,omitempty" path:"request_headers"`
	RequestMethod  string `json:"request_method,omitempty" path:"request_method"`
	RequestUrl     string `json:"request_url,omitempty" path:"request_url"`
	Path           string `json:"path,omitempty" path:"path"`
	Folder         string `json:"folder,omitempty" path:"folder"`
}

func (a ActionNotificationExportResult) Identifier() interface{} {
	return a.Id
}

type ActionNotificationExportResultCollection []ActionNotificationExportResult

type ActionNotificationExportResultListParams struct {
	UserId                     int64 `url:"user_id,omitempty" required:"false" json:"user_id,omitempty" path:"user_id"`
	ActionNotificationExportId int64 `url:"action_notification_export_id,omitempty" required:"true" json:"action_notification_export_id,omitempty" path:"action_notification_export_id"`
	ListParams
}

func (a *ActionNotificationExportResult) UnmarshalJSON(data []byte) error {
	type actionNotificationExportResult ActionNotificationExportResult
	var v actionNotificationExportResult
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*a = ActionNotificationExportResult(v)
	return nil
}

func (a *ActionNotificationExportResultCollection) UnmarshalJSON(data []byte) error {
	type actionNotificationExportResults ActionNotificationExportResultCollection
	var v actionNotificationExportResults
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
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
