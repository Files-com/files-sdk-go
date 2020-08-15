package files_sdk

import (
	"encoding/json"
	lib "github.com/Files-com/files-sdk-go/lib"
)

type HistoryExportResult struct {
	Id                  int64  `json:"id,omitempty"`
	CreatedAt           int    `json:"created_at,omitempty"`
	UserId              int64  `json:"user_id,omitempty"`
	FileId              int64  `json:"file_id,omitempty"`
	ParentId            int64  `json:"parent_id,omitempty"`
	Path                string `json:"path,omitempty"`
	Folder              string `json:"folder,omitempty"`
	Src                 string `json:"src,omitempty"`
	Destination         string `json:"destination,omitempty"`
	Ip                  string `json:"ip,omitempty"`
	Username            string `json:"username,omitempty"`
	Action              string `json:"action,omitempty"`
	FailureType         string `json:"failure_type,omitempty"`
	Interface           string `json:"interface,omitempty"`
	TargetId            int64  `json:"target_id,omitempty"`
	TargetName          string `json:"target_name,omitempty"`
	TargetPermission    string `json:"target_permission,omitempty"`
	TargetRecursive     *bool  `json:"target_recursive,omitempty"`
	TargetExpiresAt     int    `json:"target_expires_at,omitempty"`
	TargetPermissionSet string `json:"target_permission_set,omitempty"`
	TargetPlatform      string `json:"target_platform,omitempty"`
	TargetUsername      string `json:"target_username,omitempty"`
	TargetUserId        int64  `json:"target_user_id,omitempty"`
}

type HistoryExportResultCollection []HistoryExportResult

type HistoryExportResultListParams struct {
	UserId          int64  `url:"user_id,omitempty"`
	Page            int    `url:"page,omitempty"`
	PerPage         int    `url:"per_page,omitempty"`
	Action          string `url:"action,omitempty"`
	Cursor          string `url:"cursor,omitempty"`
	HistoryExportId int64  `url:"history_export_id,omitempty"`
	lib.ListParams
}

func (h *HistoryExportResult) UnmarshalJSON(data []byte) error {
	type historyExportResult HistoryExportResult
	var v historyExportResult
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*h = HistoryExportResult(v)
	return nil
}

func (h *HistoryExportResultCollection) UnmarshalJSON(data []byte) error {
	type historyExportResults []HistoryExportResult
	var v historyExportResults
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*h = HistoryExportResultCollection(v)
	return nil
}
