package files_sdk

import (
	"encoding/json"
	"time"
)

type HistoryExport struct {
	Id                       int64     `json:"id,omitempty"`
	StartAt                  time.Time `json:"start_at,omitempty"`
	EndAt                    time.Time `json:"end_at,omitempty"`
	Status                   string    `json:"status,omitempty"`
	QueryAction              string    `json:"query_action,omitempty"`
	QueryInterface           string    `json:"query_interface,omitempty"`
	QueryUserId              string    `json:"query_user_id,omitempty"`
	QueryFileId              string    `json:"query_file_id,omitempty"`
	QueryParentId            string    `json:"query_parent_id,omitempty"`
	QueryPath                string    `json:"query_path,omitempty"`
	QueryFolder              string    `json:"query_folder,omitempty"`
	QuerySrc                 string    `json:"query_src,omitempty"`
	QueryDestination         string    `json:"query_destination,omitempty"`
	QueryIp                  string    `json:"query_ip,omitempty"`
	QueryUsername            string    `json:"query_username,omitempty"`
	QueryFailureType         string    `json:"query_failure_type,omitempty"`
	QueryTargetId            string    `json:"query_target_id,omitempty"`
	QueryTargetName          string    `json:"query_target_name,omitempty"`
	QueryTargetPermission    string    `json:"query_target_permission,omitempty"`
	QueryTargetUserId        string    `json:"query_target_user_id,omitempty"`
	QueryTargetUsername      string    `json:"query_target_username,omitempty"`
	QueryTargetPlatform      string    `json:"query_target_platform,omitempty"`
	QueryTargetPermissionSet string    `json:"query_target_permission_set,omitempty"`
	ResultsUrl               string    `json:"results_url,omitempty"`
	UserId                   int64     `json:"user_id,omitempty"`
}

type HistoryExportCollection []HistoryExport

type HistoryExportFindParams struct {
	Id int64 `url:"-,omitempty"`
}

type HistoryExportCreateParams struct {
	UserId                   int64     `url:"user_id,omitempty"`
	StartAt                  time.Time `url:"start_at,omitempty"`
	EndAt                    time.Time `url:"end_at,omitempty"`
	QueryAction              string    `url:"query_action,omitempty"`
	QueryInterface           string    `url:"query_interface,omitempty"`
	QueryUserId              string    `url:"query_user_id,omitempty"`
	QueryFileId              string    `url:"query_file_id,omitempty"`
	QueryParentId            string    `url:"query_parent_id,omitempty"`
	QueryPath                string    `url:"query_path,omitempty"`
	QueryFolder              string    `url:"query_folder,omitempty"`
	QuerySrc                 string    `url:"query_src,omitempty"`
	QueryDestination         string    `url:"query_destination,omitempty"`
	QueryIp                  string    `url:"query_ip,omitempty"`
	QueryUsername            string    `url:"query_username,omitempty"`
	QueryFailureType         string    `url:"query_failure_type,omitempty"`
	QueryTargetId            string    `url:"query_target_id,omitempty"`
	QueryTargetName          string    `url:"query_target_name,omitempty"`
	QueryTargetPermission    string    `url:"query_target_permission,omitempty"`
	QueryTargetUserId        string    `url:"query_target_user_id,omitempty"`
	QueryTargetUsername      string    `url:"query_target_username,omitempty"`
	QueryTargetPlatform      string    `url:"query_target_platform,omitempty"`
	QueryTargetPermissionSet string    `url:"query_target_permission_set,omitempty"`
}

func (h *HistoryExport) UnmarshalJSON(data []byte) error {
	type historyExport HistoryExport
	var v historyExport
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*h = HistoryExport(v)
	return nil
}

func (h *HistoryExportCollection) UnmarshalJSON(data []byte) error {
	type historyExports []HistoryExport
	var v historyExports
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*h = HistoryExportCollection(v)
	return nil
}
