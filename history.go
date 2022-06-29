package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type History struct {
	Id          int64      `json:"id,omitempty"`
	Path        string     `json:"path,omitempty"`
	When        *time.Time `json:"when,omitempty"`
	Destination string     `json:"destination,omitempty"`
	Display     string     `json:"display,omitempty"`
	Ip          string     `json:"ip,omitempty"`
	Source      string     `json:"source,omitempty"`
	Targets     []string   `json:"targets,omitempty"`
	UserId      int64      `json:"user_id,omitempty"`
	Username    string     `json:"username,omitempty"`
	Action      string     `json:"action,omitempty"`
	FailureType string     `json:"failure_type,omitempty"`
	Interface   string     `json:"interface,omitempty"`
}

type HistoryCollection []History

type HistoryListForFileParams struct {
	StartAt *time.Time      `url:"start_at,omitempty" required:"false" json:"start_at,omitempty"`
	EndAt   *time.Time      `url:"end_at,omitempty" required:"false" json:"end_at,omitempty"`
	Display string          `url:"display,omitempty" required:"false" json:"display,omitempty"`
	Cursor  string          `url:"cursor,omitempty" required:"false" json:"cursor,omitempty"`
	PerPage int64           `url:"per_page,omitempty" required:"false" json:"per_page,omitempty"`
	SortBy  json.RawMessage `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty"`
	Path    string          `url:"-,omitempty" required:"true" json:"-,omitempty"`
	lib.ListParams
}

type HistoryListForFolderParams struct {
	StartAt *time.Time      `url:"start_at,omitempty" required:"false" json:"start_at,omitempty"`
	EndAt   *time.Time      `url:"end_at,omitempty" required:"false" json:"end_at,omitempty"`
	Display string          `url:"display,omitempty" required:"false" json:"display,omitempty"`
	Cursor  string          `url:"cursor,omitempty" required:"false" json:"cursor,omitempty"`
	PerPage int64           `url:"per_page,omitempty" required:"false" json:"per_page,omitempty"`
	SortBy  json.RawMessage `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty"`
	Path    string          `url:"-,omitempty" required:"true" json:"-,omitempty"`
	lib.ListParams
}

type HistoryListForUserParams struct {
	StartAt *time.Time      `url:"start_at,omitempty" required:"false" json:"start_at,omitempty"`
	EndAt   *time.Time      `url:"end_at,omitempty" required:"false" json:"end_at,omitempty"`
	Display string          `url:"display,omitempty" required:"false" json:"display,omitempty"`
	Cursor  string          `url:"cursor,omitempty" required:"false" json:"cursor,omitempty"`
	PerPage int64           `url:"per_page,omitempty" required:"false" json:"per_page,omitempty"`
	SortBy  json.RawMessage `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty"`
	UserId  int64           `url:"-,omitempty" required:"true" json:"-,omitempty"`
	lib.ListParams
}

type HistoryListLoginsParams struct {
	StartAt *time.Time      `url:"start_at,omitempty" required:"false" json:"start_at,omitempty"`
	EndAt   *time.Time      `url:"end_at,omitempty" required:"false" json:"end_at,omitempty"`
	Display string          `url:"display,omitempty" required:"false" json:"display,omitempty"`
	Cursor  string          `url:"cursor,omitempty" required:"false" json:"cursor,omitempty"`
	PerPage int64           `url:"per_page,omitempty" required:"false" json:"per_page,omitempty"`
	SortBy  json.RawMessage `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty"`
	lib.ListParams
}

type HistoryListParams struct {
	StartAt    *time.Time      `url:"start_at,omitempty" required:"false" json:"start_at,omitempty"`
	EndAt      *time.Time      `url:"end_at,omitempty" required:"false" json:"end_at,omitempty"`
	Display    string          `url:"display,omitempty" required:"false" json:"display,omitempty"`
	Cursor     string          `url:"cursor,omitempty" required:"false" json:"cursor,omitempty"`
	PerPage    int64           `url:"per_page,omitempty" required:"false" json:"per_page,omitempty"`
	SortBy     json.RawMessage `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty"`
	Filter     json.RawMessage `url:"filter,omitempty" required:"false" json:"filter,omitempty"`
	FilterGt   json.RawMessage `url:"filter_gt,omitempty" required:"false" json:"filter_gt,omitempty"`
	FilterGteq json.RawMessage `url:"filter_gteq,omitempty" required:"false" json:"filter_gteq,omitempty"`
	FilterLike json.RawMessage `url:"filter_like,omitempty" required:"false" json:"filter_like,omitempty"`
	FilterLt   json.RawMessage `url:"filter_lt,omitempty" required:"false" json:"filter_lt,omitempty"`
	FilterLteq json.RawMessage `url:"filter_lteq,omitempty" required:"false" json:"filter_lteq,omitempty"`
	lib.ListParams
}

func (h *History) UnmarshalJSON(data []byte) error {
	type history History
	var v history
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*h = History(v)
	return nil
}

func (h *HistoryCollection) UnmarshalJSON(data []byte) error {
	type historys []History
	var v historys
	if err := json.Unmarshal(data, &v); err != nil {
		return err
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
