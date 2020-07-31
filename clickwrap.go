package files_sdk

import (
	"encoding/json"
	lib "github.com/Files-com/files-sdk-go/lib"
)

type Clickwrap struct {
	Name           string `json:"name,omitempty"`
	Body           string `json:"body,omitempty"`
	UseWithUsers   string `json:"use_with_users,omitempty"`
	UseWithBundles string `json:"use_with_bundles,omitempty"`
	UseWithInboxes string `json:"use_with_inboxes,omitempty"`
	Id             int64  `json:"id,omitempty"`
}

type ClickwrapCollection []Clickwrap

type ClickwrapListParams struct {
	Page    int    `url:"page,omitempty"`
	PerPage int    `url:"per_page,omitempty"`
	Action  string `url:"action,omitempty"`
	lib.ListParams
}

type ClickwrapFindParams struct {
	Id int64 `url:"-,omitempty"`
}

type ClickwrapCreateParams struct {
	Name           string `url:"name,omitempty"`
	Body           string `url:"body,omitempty"`
	UseWithBundles string `url:"use_with_bundles,omitempty"`
	UseWithInboxes string `url:"use_with_inboxes,omitempty"`
	UseWithUsers   string `url:"use_with_users,omitempty"`
}

type ClickwrapUpdateParams struct {
	Id             int64  `url:"-,omitempty"`
	Name           string `url:"name,omitempty"`
	Body           string `url:"body,omitempty"`
	UseWithBundles string `url:"use_with_bundles,omitempty"`
	UseWithInboxes string `url:"use_with_inboxes,omitempty"`
	UseWithUsers   string `url:"use_with_users,omitempty"`
}

type ClickwrapDeleteParams struct {
	Id int64 `url:"-,omitempty"`
}

func (c *Clickwrap) UnmarshalJSON(data []byte) error {
	type clickwrap Clickwrap
	var v clickwrap
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*c = Clickwrap(v)
	return nil
}

func (c *ClickwrapCollection) UnmarshalJSON(data []byte) error {
	type clickwraps []Clickwrap
	var v clickwraps
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*c = ClickwrapCollection(v)
	return nil
}
