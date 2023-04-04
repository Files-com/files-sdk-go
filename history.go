package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type History struct {
	Id          int64      `json:"id,omitempty" path:"id"`
	Path        string     `json:"path,omitempty" path:"path"`
	When        *time.Time `json:"when,omitempty" path:"when"`
	Destination string     `json:"destination,omitempty" path:"destination"`
	Display     string     `json:"display,omitempty" path:"display"`
	Ip          string     `json:"ip,omitempty" path:"ip"`
	Source      string     `json:"source,omitempty" path:"source"`
	Targets     []string   `json:"targets,omitempty" path:"targets"`
	UserId      int64      `json:"user_id,omitempty" path:"user_id"`
	Username    string     `json:"username,omitempty" path:"username"`
	Action      string     `json:"action,omitempty" path:"action"`
	FailureType string     `json:"failure_type,omitempty" path:"failure_type"`
	Interface   string     `json:"interface,omitempty" path:"interface"`
}

type HistoryCollection []History

type HistoryListForFileParams struct {
	StartAt *time.Time      `url:"start_at,omitempty" required:"false" json:"start_at,omitempty" path:"start_at"`
	EndAt   *time.Time      `url:"end_at,omitempty" required:"false" json:"end_at,omitempty" path:"end_at"`
	Display string          `url:"display,omitempty" required:"false" json:"display,omitempty" path:"display"`
	SortBy  json.RawMessage `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty" path:"sort_by"`
	Path    string          `url:"-,omitempty" required:"false" json:"-,omitempty" path:"path"`
	lib.ListParams
}

type HistoryListForFolderParams struct {
	StartAt *time.Time      `url:"start_at,omitempty" required:"false" json:"start_at,omitempty" path:"start_at"`
	EndAt   *time.Time      `url:"end_at,omitempty" required:"false" json:"end_at,omitempty" path:"end_at"`
	Display string          `url:"display,omitempty" required:"false" json:"display,omitempty" path:"display"`
	SortBy  json.RawMessage `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty" path:"sort_by"`
	Path    string          `url:"-,omitempty" required:"false" json:"-,omitempty" path:"path"`
	lib.ListParams
}

type HistoryListForUserParams struct {
	StartAt *time.Time      `url:"start_at,omitempty" required:"false" json:"start_at,omitempty" path:"start_at"`
	EndAt   *time.Time      `url:"end_at,omitempty" required:"false" json:"end_at,omitempty" path:"end_at"`
	Display string          `url:"display,omitempty" required:"false" json:"display,omitempty" path:"display"`
	SortBy  json.RawMessage `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty" path:"sort_by"`
	UserId  int64           `url:"-,omitempty" required:"false" json:"-,omitempty" path:"user_id"`
	lib.ListParams
}

type HistoryListLoginsParams struct {
	StartAt *time.Time      `url:"start_at,omitempty" required:"false" json:"start_at,omitempty" path:"start_at"`
	EndAt   *time.Time      `url:"end_at,omitempty" required:"false" json:"end_at,omitempty" path:"end_at"`
	Display string          `url:"display,omitempty" required:"false" json:"display,omitempty" path:"display"`
	SortBy  json.RawMessage `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty" path:"sort_by"`
	lib.ListParams
}

type HistoryListParams struct {
	StartAt      *time.Time      `url:"start_at,omitempty" required:"false" json:"start_at,omitempty" path:"start_at"`
	EndAt        *time.Time      `url:"end_at,omitempty" required:"false" json:"end_at,omitempty" path:"end_at"`
	Display      string          `url:"display,omitempty" required:"false" json:"display,omitempty" path:"display"`
	SortBy       json.RawMessage `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty" path:"sort_by"`
	Filter       json.RawMessage `url:"filter,omitempty" required:"false" json:"filter,omitempty" path:"filter"`
	FilterPrefix json.RawMessage `url:"filter_prefix,omitempty" required:"false" json:"filter_prefix,omitempty" path:"filter_prefix"`
	lib.ListParams
}

func (h *History) UnmarshalJSON(data []byte) error {
	type history History
	var v history
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*h = History(v)
	return nil
}

func (h *HistoryCollection) UnmarshalJSON(data []byte) error {
	type historys HistoryCollection
	var v historys
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*h = HistoryCollection(v)
	return nil
}

func (h *HistoryCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*h))
	for i, v := range *h {
		ret[i] = v
	}

	return &ret
}
