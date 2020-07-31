package files_sdk

import (
	"encoding/json"
	lib "github.com/Files-com/files-sdk-go/lib"
	"time"
)

type History struct {
	Id          int       `json:"id,omitempty"`
	Path        string    `json:"path,omitempty"`
	When        time.Time `json:"when,omitempty"`
	Destination string    `json:"destination,omitempty"`
	Display     string    `json:"display,omitempty"`
	Ip          string    `json:"ip,omitempty"`
	Source      string    `json:"source,omitempty"`
	Targets     []string  `json:"targets,omitempty"`
	UserId      int       `json:"user_id,omitempty"`
	Username    string    `json:"username,omitempty"`
	Action      string    `json:"action,omitempty"`
	FailureType string    `json:"failure_type,omitempty"`
	Interface   string    `json:"interface,omitempty"`
}

type HistoryCollection []History

type HistoryListForFileParams struct {
	StartAt string          `url:"start_at,omitempty"`
	EndAt   string          `url:"end_at,omitempty"`
	Display string          `url:"display,omitempty"`
	Page    int             `url:"page,omitempty"`
	PerPage int             `url:"per_page,omitempty"`
	Action  string          `url:"action,omitempty"`
	Cursor  string          `url:"cursor,omitempty"`
	SortBy  json.RawMessage `url:"sort_by,omitempty"`
	Path    string          `url:"-,omitempty"`
	lib.ListParams
}

type HistoryListForFolderParams struct {
	StartAt string          `url:"start_at,omitempty"`
	EndAt   string          `url:"end_at,omitempty"`
	Display string          `url:"display,omitempty"`
	Page    int             `url:"page,omitempty"`
	PerPage int             `url:"per_page,omitempty"`
	Action  string          `url:"action,omitempty"`
	Cursor  string          `url:"cursor,omitempty"`
	SortBy  json.RawMessage `url:"sort_by,omitempty"`
	Path    string          `url:"-,omitempty"`
	lib.ListParams
}

type HistoryListForUserParams struct {
	StartAt string          `url:"start_at,omitempty"`
	EndAt   string          `url:"end_at,omitempty"`
	Display string          `url:"display,omitempty"`
	Page    int             `url:"page,omitempty"`
	PerPage int             `url:"per_page,omitempty"`
	Action  string          `url:"action,omitempty"`
	Cursor  string          `url:"cursor,omitempty"`
	SortBy  json.RawMessage `url:"sort_by,omitempty"`
	UserId  int             `url:"-,omitempty"`
	lib.ListParams
}

type HistoryListLoginsParams struct {
	StartAt string          `url:"start_at,omitempty"`
	EndAt   string          `url:"end_at,omitempty"`
	Display string          `url:"display,omitempty"`
	Page    int             `url:"page,omitempty"`
	PerPage int             `url:"per_page,omitempty"`
	Action  string          `url:"action,omitempty"`
	Cursor  string          `url:"cursor,omitempty"`
	SortBy  json.RawMessage `url:"sort_by,omitempty"`
	lib.ListParams
}

type HistoryListParams struct {
	StartAt    string          `url:"start_at,omitempty"`
	EndAt      string          `url:"end_at,omitempty"`
	Display    string          `url:"display,omitempty"`
	Page       int             `url:"page,omitempty"`
	PerPage    int             `url:"per_page,omitempty"`
	Action     string          `url:"action,omitempty"`
	Cursor     string          `url:"cursor,omitempty"`
	SortBy     json.RawMessage `url:"sort_by,omitempty"`
	Filter     json.RawMessage `url:"filter,omitempty"`
	FilterGt   json.RawMessage `url:"filter_gt,omitempty"`
	FilterGteq json.RawMessage `url:"filter_gteq,omitempty"`
	FilterLike json.RawMessage `url:"filter_like,omitempty"`
	FilterLt   json.RawMessage `url:"filter_lt,omitempty"`
	FilterLteq json.RawMessage `url:"filter_lteq,omitempty"`
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
