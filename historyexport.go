package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type HistoryExport struct {
	Id                       int64      `json:"id,omitempty" path:"id"`
	HistoryVersion           string     `json:"history_version,omitempty" path:"history_version"`
	StartAt                  *time.Time `json:"start_at,omitempty" path:"start_at"`
	EndAt                    *time.Time `json:"end_at,omitempty" path:"end_at"`
	Status                   string     `json:"status,omitempty" path:"status"`
	QueryAction              string     `json:"query_action,omitempty" path:"query_action"`
	QueryInterface           string     `json:"query_interface,omitempty" path:"query_interface"`
	QueryUserId              string     `json:"query_user_id,omitempty" path:"query_user_id"`
	QueryFileId              string     `json:"query_file_id,omitempty" path:"query_file_id"`
	QueryParentId            string     `json:"query_parent_id,omitempty" path:"query_parent_id"`
	QueryPath                string     `json:"query_path,omitempty" path:"query_path"`
	QueryFolder              string     `json:"query_folder,omitempty" path:"query_folder"`
	QuerySrc                 string     `json:"query_src,omitempty" path:"query_src"`
	QueryDestination         string     `json:"query_destination,omitempty" path:"query_destination"`
	QueryIp                  string     `json:"query_ip,omitempty" path:"query_ip"`
	QueryUsername            string     `json:"query_username,omitempty" path:"query_username"`
	QueryFailureType         string     `json:"query_failure_type,omitempty" path:"query_failure_type"`
	QueryTargetId            string     `json:"query_target_id,omitempty" path:"query_target_id"`
	QueryTargetName          string     `json:"query_target_name,omitempty" path:"query_target_name"`
	QueryTargetPermission    string     `json:"query_target_permission,omitempty" path:"query_target_permission"`
	QueryTargetUserId        string     `json:"query_target_user_id,omitempty" path:"query_target_user_id"`
	QueryTargetUsername      string     `json:"query_target_username,omitempty" path:"query_target_username"`
	QueryTargetPlatform      string     `json:"query_target_platform,omitempty" path:"query_target_platform"`
	QueryTargetPermissionSet string     `json:"query_target_permission_set,omitempty" path:"query_target_permission_set"`
	ResultsUrl               string     `json:"results_url,omitempty" path:"results_url"`
	UserId                   int64      `json:"user_id,omitempty" path:"user_id"`
}

func (h HistoryExport) Identifier() interface{} {
	return h.Id
}

type HistoryExportCollection []HistoryExport

type HistoryExportFindParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
}

type HistoryExportCreateParams struct {
	UserId                   int64      `url:"user_id,omitempty" required:"false" json:"user_id,omitempty" path:"user_id"`
	StartAt                  *time.Time `url:"start_at,omitempty" required:"false" json:"start_at,omitempty" path:"start_at"`
	EndAt                    *time.Time `url:"end_at,omitempty" required:"false" json:"end_at,omitempty" path:"end_at"`
	QueryAction              string     `url:"query_action,omitempty" required:"false" json:"query_action,omitempty" path:"query_action"`
	QueryInterface           string     `url:"query_interface,omitempty" required:"false" json:"query_interface,omitempty" path:"query_interface"`
	QueryUserId              string     `url:"query_user_id,omitempty" required:"false" json:"query_user_id,omitempty" path:"query_user_id"`
	QueryFileId              string     `url:"query_file_id,omitempty" required:"false" json:"query_file_id,omitempty" path:"query_file_id"`
	QueryParentId            string     `url:"query_parent_id,omitempty" required:"false" json:"query_parent_id,omitempty" path:"query_parent_id"`
	QueryPath                string     `url:"query_path,omitempty" required:"false" json:"query_path,omitempty" path:"query_path"`
	QueryFolder              string     `url:"query_folder,omitempty" required:"false" json:"query_folder,omitempty" path:"query_folder"`
	QuerySrc                 string     `url:"query_src,omitempty" required:"false" json:"query_src,omitempty" path:"query_src"`
	QueryDestination         string     `url:"query_destination,omitempty" required:"false" json:"query_destination,omitempty" path:"query_destination"`
	QueryIp                  string     `url:"query_ip,omitempty" required:"false" json:"query_ip,omitempty" path:"query_ip"`
	QueryUsername            string     `url:"query_username,omitempty" required:"false" json:"query_username,omitempty" path:"query_username"`
	QueryFailureType         string     `url:"query_failure_type,omitempty" required:"false" json:"query_failure_type,omitempty" path:"query_failure_type"`
	QueryTargetId            string     `url:"query_target_id,omitempty" required:"false" json:"query_target_id,omitempty" path:"query_target_id"`
	QueryTargetName          string     `url:"query_target_name,omitempty" required:"false" json:"query_target_name,omitempty" path:"query_target_name"`
	QueryTargetPermission    string     `url:"query_target_permission,omitempty" required:"false" json:"query_target_permission,omitempty" path:"query_target_permission"`
	QueryTargetUserId        string     `url:"query_target_user_id,omitempty" required:"false" json:"query_target_user_id,omitempty" path:"query_target_user_id"`
	QueryTargetUsername      string     `url:"query_target_username,omitempty" required:"false" json:"query_target_username,omitempty" path:"query_target_username"`
	QueryTargetPlatform      string     `url:"query_target_platform,omitempty" required:"false" json:"query_target_platform,omitempty" path:"query_target_platform"`
	QueryTargetPermissionSet string     `url:"query_target_permission_set,omitempty" required:"false" json:"query_target_permission_set,omitempty" path:"query_target_permission_set"`
}

func (h *HistoryExport) UnmarshalJSON(data []byte) error {
	type historyExport HistoryExport
	var v historyExport
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*h = HistoryExport(v)
	return nil
}

func (h *HistoryExportCollection) UnmarshalJSON(data []byte) error {
	type historyExports HistoryExportCollection
	var v historyExports
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*h = HistoryExportCollection(v)
	return nil
}

func (h *HistoryExportCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*h))
	for i, v := range *h {
		ret[i] = v
	}

	return &ret
}
