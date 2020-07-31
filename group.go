package files_sdk

import (
	"encoding/json"
	lib "github.com/Files-com/files-sdk-go/lib"
)

type Group struct {
	Id        int      `json:"id,omitempty"`
	Name      string   `json:"name,omitempty"`
	AdminIds  []string `json:"admin_ids,omitempty"`
	Notes     string   `json:"notes,omitempty"`
	UserIds   []string `json:"user_ids,omitempty"`
	Usernames []string `json:"usernames,omitempty"`
}

type GroupCollection []Group

type GroupListParams struct {
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
	Ids        string          `url:"ids,omitempty"`
	lib.ListParams
}

type GroupFindParams struct {
	Id int `url:"-,omitempty"`
}

type GroupCreateParams struct {
	Name     string `url:"name,omitempty"`
	Notes    string `url:"notes,omitempty"`
	UserIds  string `url:"user_ids,omitempty"`
	AdminIds string `url:"admin_ids,omitempty"`
}

type GroupUpdateParams struct {
	Id       int    `url:"-,omitempty"`
	Name     string `url:"name,omitempty"`
	Notes    string `url:"notes,omitempty"`
	UserIds  string `url:"user_ids,omitempty"`
	AdminIds string `url:"admin_ids,omitempty"`
}

type GroupDeleteParams struct {
	Id int `url:"-,omitempty"`
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
