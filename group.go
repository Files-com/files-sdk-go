package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Group struct {
	Id        int64    `json:"id,omitempty"`
	Name      string   `json:"name,omitempty"`
	AdminIds  string   `json:"admin_ids,omitempty"`
	Notes     string   `json:"notes,omitempty"`
	UserIds   []int64  `json:"user_ids,omitempty"`
	Usernames []string `json:"usernames,omitempty"`
}

type GroupCollection []Group

type GroupListParams struct {
	Cursor     string          `url:"cursor,omitempty" required:"false" json:"cursor,omitempty"`
	PerPage    int64           `url:"per_page,omitempty" required:"false" json:"per_page,omitempty"`
	SortBy     json.RawMessage `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty"`
	Filter     json.RawMessage `url:"filter,omitempty" required:"false" json:"filter,omitempty"`
	FilterGt   json.RawMessage `url:"filter_gt,omitempty" required:"false" json:"filter_gt,omitempty"`
	FilterGteq json.RawMessage `url:"filter_gteq,omitempty" required:"false" json:"filter_gteq,omitempty"`
	FilterLike json.RawMessage `url:"filter_like,omitempty" required:"false" json:"filter_like,omitempty"`
	FilterLt   json.RawMessage `url:"filter_lt,omitempty" required:"false" json:"filter_lt,omitempty"`
	FilterLteq json.RawMessage `url:"filter_lteq,omitempty" required:"false" json:"filter_lteq,omitempty"`
	Ids        string          `url:"ids,omitempty" required:"false" json:"ids,omitempty"`
	lib.ListParams
}

type GroupFindParams struct {
	Id int64 `url:"-,omitempty" required:"true" json:"-,omitempty"`
}

type GroupCreateParams struct {
	Name     string `url:"name,omitempty" required:"false" json:"name,omitempty"`
	Notes    string `url:"notes,omitempty" required:"false" json:"notes,omitempty"`
	UserIds  string `url:"user_ids,omitempty" required:"false" json:"user_ids,omitempty"`
	AdminIds string `url:"admin_ids,omitempty" required:"false" json:"admin_ids,omitempty"`
}

type GroupUpdateParams struct {
	Id       int64  `url:"-,omitempty" required:"true" json:"-,omitempty"`
	Name     string `url:"name,omitempty" required:"false" json:"name,omitempty"`
	Notes    string `url:"notes,omitempty" required:"false" json:"notes,omitempty"`
	UserIds  string `url:"user_ids,omitempty" required:"false" json:"user_ids,omitempty"`
	AdminIds string `url:"admin_ids,omitempty" required:"false" json:"admin_ids,omitempty"`
}

type GroupDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"true" json:"-,omitempty"`
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

func (g *GroupCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*g))
	for i, v := range *g {
		ret[i] = v
	}

	return &ret
}
