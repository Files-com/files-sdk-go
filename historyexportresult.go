package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type HistoryExportResult struct {
	Id                  int64  `json:"id,omitempty" path:"id"`
	CreatedAt           int64  `json:"created_at,omitempty" path:"created_at"`
	CreatedAtIso8601    int64  `json:"created_at_iso8601,omitempty" path:"created_at_iso8601"`
	UserId              int64  `json:"user_id,omitempty" path:"user_id"`
	FileId              int64  `json:"file_id,omitempty" path:"file_id"`
	ParentId            int64  `json:"parent_id,omitempty" path:"parent_id"`
	Path                string `json:"path,omitempty" path:"path"`
	Folder              string `json:"folder,omitempty" path:"folder"`
	Src                 string `json:"src,omitempty" path:"src"`
	Destination         string `json:"destination,omitempty" path:"destination"`
	Ip                  string `json:"ip,omitempty" path:"ip"`
	Username            string `json:"username,omitempty" path:"username"`
	Action              string `json:"action,omitempty" path:"action"`
	FailureType         string `json:"failure_type,omitempty" path:"failure_type"`
	Interface           string `json:"interface,omitempty" path:"interface"`
	TargetId            int64  `json:"target_id,omitempty" path:"target_id"`
	TargetName          string `json:"target_name,omitempty" path:"target_name"`
	TargetPermission    string `json:"target_permission,omitempty" path:"target_permission"`
	TargetRecursive     *bool  `json:"target_recursive,omitempty" path:"target_recursive"`
	TargetExpiresAt     int64  `json:"target_expires_at,omitempty" path:"target_expires_at"`
	TargetPermissionSet string `json:"target_permission_set,omitempty" path:"target_permission_set"`
	TargetPlatform      string `json:"target_platform,omitempty" path:"target_platform"`
	TargetUsername      string `json:"target_username,omitempty" path:"target_username"`
	TargetUserId        int64  `json:"target_user_id,omitempty" path:"target_user_id"`
}

type HistoryExportResultCollection []HistoryExportResult

type HistoryExportResultListParams struct {
	UserId          int64 `url:"user_id,omitempty" required:"false" json:"user_id,omitempty" path:"user_id"`
	HistoryExportId int64 `url:"history_export_id,omitempty" required:"true" json:"history_export_id,omitempty" path:"history_export_id"`
	lib.ListParams
}

func (h *HistoryExportResult) UnmarshalJSON(data []byte) error {
	type historyExportResult HistoryExportResult
	var v historyExportResult
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*h = HistoryExportResult(v)
	return nil
}

func (h *HistoryExportResultCollection) UnmarshalJSON(data []byte) error {
	type historyExportResults HistoryExportResultCollection
	var v historyExportResults
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*h = HistoryExportResultCollection(v)
	return nil
}

func (h *HistoryExportResultCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*h))
	for i, v := range *h {
		ret[i] = v
	}

	return &ret
}
