package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/lib"
)

type Request struct {
	Id              int64  `json:"id,omitempty"`
	Path            string `json:"path,omitempty"`
	Source          string `json:"source,omitempty"`
	Destination     string `json:"destination,omitempty"`
	AutomationId    string `json:"automation_id,omitempty"`
	UserDisplayName string `json:"user_display_name,omitempty"`
	UserIds         string `json:"user_ids,omitempty"`
	GroupIds        string `json:"group_ids,omitempty"`
}

type RequestCollection []Request

type RequestListParams struct {
	Cursor  string          `url:"cursor,omitempty" required:"false"`
	PerPage int             `url:"per_page,omitempty" required:"false"`
	SortBy  json.RawMessage `url:"sort_by,omitempty" required:"false"`
	Mine    *bool           `url:"mine,omitempty" required:"false"`
	Path    string          `url:"path,omitempty" required:"false"`
	lib.ListParams
}

type RequestGetFolderParams struct {
	Cursor  string          `url:"cursor,omitempty" required:"false"`
	PerPage int             `url:"per_page,omitempty" required:"false"`
	SortBy  json.RawMessage `url:"sort_by,omitempty" required:"false"`
	Mine    *bool           `url:"mine,omitempty" required:"false"`
	Path    string          `url:"-,omitempty" required:"true"`
}

type RequestCreateParams struct {
	Path        string `url:"path,omitempty" required:"true"`
	Destination string `url:"destination,omitempty" required:"true"`
	UserIds     string `url:"user_ids,omitempty" required:"false"`
	GroupIds    string `url:"group_ids,omitempty" required:"false"`
}

type RequestDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"true"`
}

func (r *Request) UnmarshalJSON(data []byte) error {
	type request Request
	var v request
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*r = Request(v)
	return nil
}

func (r *RequestCollection) UnmarshalJSON(data []byte) error {
	type requests []Request
	var v requests
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*r = RequestCollection(v)
	return nil
}
