package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/lib"
)

type Group struct {
	Id        int64    `json:"id,omitempty"`
	Name      string   `json:"name,omitempty"`
	AdminIds  []string `json:"admin_ids,omitempty"`
	Notes     string   `json:"notes,omitempty"`
	UserIds   []string `json:"user_ids,omitempty"`
	Usernames []string `json:"usernames,omitempty"`
}

type GroupCollection []Group

type GroupListParams struct {
	Cursor     string          `url:"cursor,omitempty" required:"false"`
	PerPage    int             `url:"per_page,omitempty" required:"false"`
	SortBy     json.RawMessage `url:"sort_by,omitempty" required:"false"`
	Filter     json.RawMessage `url:"filter,omitempty" required:"false"`
	FilterGt   json.RawMessage `url:"filter_gt,omitempty" required:"false"`
	FilterGteq json.RawMessage `url:"filter_gteq,omitempty" required:"false"`
	FilterLike json.RawMessage `url:"filter_like,omitempty" required:"false"`
	FilterLt   json.RawMessage `url:"filter_lt,omitempty" required:"false"`
	FilterLteq json.RawMessage `url:"filter_lteq,omitempty" required:"false"`
	Ids        string          `url:"ids,omitempty" required:"false"`
	lib.ListParams
}

type GroupFindParams struct {
	Id int64 `url:"-,omitempty" required:"true"`
}

type GroupCreateParams struct {
	Name     string `url:"name,omitempty" required:"false"`
	Notes    string `url:"notes,omitempty" required:"false"`
	UserIds  string `url:"user_ids,omitempty" required:"false"`
	AdminIds string `url:"admin_ids,omitempty" required:"false"`
}

type GroupUpdateParams struct {
	Id       int64  `url:"-,omitempty" required:"true"`
	Name     string `url:"name,omitempty" required:"false"`
	Notes    string `url:"notes,omitempty" required:"false"`
	UserIds  string `url:"user_ids,omitempty" required:"false"`
	AdminIds string `url:"admin_ids,omitempty" required:"false"`
}

type GroupDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"true"`
}

func (g *Group) UnmarshalJSON(data []byte) error {
	type group Group
	var v group
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*g = Group(v)
	return nil
}

func (g *GroupCollection) UnmarshalJSON(data []byte) error {
	type groups []Group
	var v groups
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*g = GroupCollection(v)
	return nil
}
