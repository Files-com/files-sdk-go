package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
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
	Cursor  string          `url:"cursor,omitempty" required:"false" json:"cursor,omitempty"`
	PerPage int64           `url:"per_page,omitempty" required:"false" json:"per_page,omitempty"`
	SortBy  json.RawMessage `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty"`
	Mine    *bool           `url:"mine,omitempty" required:"false" json:"mine,omitempty"`
	Path    string          `url:"path,omitempty" required:"false" json:"path,omitempty"`
	lib.ListParams
}

type RequestGetFolderParams struct {
	Cursor  string          `url:"cursor,omitempty" required:"false" json:"cursor,omitempty"`
	PerPage int64           `url:"per_page,omitempty" required:"false" json:"per_page,omitempty"`
	SortBy  json.RawMessage `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty"`
	Mine    *bool           `url:"mine,omitempty" required:"false" json:"mine,omitempty"`
	Path    string          `url:"-,omitempty" required:"true" json:"-,omitempty"`
}

type RequestCreateParams struct {
	Path        string `url:"path,omitempty" required:"true" json:"path,omitempty"`
	Destination string `url:"destination,omitempty" required:"true" json:"destination,omitempty"`
	UserIds     string `url:"user_ids,omitempty" required:"false" json:"user_ids,omitempty"`
	GroupIds    string `url:"group_ids,omitempty" required:"false" json:"group_ids,omitempty"`
}

type RequestDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"true" json:"-,omitempty"`
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

func (r *RequestCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*r))
	for i, v := range *r {
		ret[i] = v
	}

	return &ret
}
