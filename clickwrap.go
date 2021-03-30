package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/lib"
)

type Clickwrap struct {
	Id             int64  `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	Body           string `json:"body,omitempty"`
	UseWithUsers   string `json:"use_with_users,omitempty"`
	UseWithBundles string `json:"use_with_bundles,omitempty"`
	UseWithInboxes string `json:"use_with_inboxes,omitempty"`
}

type ClickwrapCollection []Clickwrap

type ClickwrapListParams struct {
	Cursor  string `url:"cursor,omitempty" required:"false"`
	PerPage int    `url:"per_page,omitempty" required:"false"`
	lib.ListParams
}

type ClickwrapFindParams struct {
	Id int64 `url:"-,omitempty" required:"true"`
}

type ClickwrapCreateParams struct {
	Name           string `url:"name,omitempty" required:"false"`
	Body           string `url:"body,omitempty" required:"false"`
	UseWithBundles string `url:"use_with_bundles,omitempty" required:"false"`
	UseWithInboxes string `url:"use_with_inboxes,omitempty" required:"false"`
	UseWithUsers   string `url:"use_with_users,omitempty" required:"false"`
}

type ClickwrapUpdateParams struct {
	Id             int64  `url:"-,omitempty" required:"true"`
	Name           string `url:"name,omitempty" required:"false"`
	Body           string `url:"body,omitempty" required:"false"`
	UseWithBundles string `url:"use_with_bundles,omitempty" required:"false"`
	UseWithInboxes string `url:"use_with_inboxes,omitempty" required:"false"`
	UseWithUsers   string `url:"use_with_users,omitempty" required:"false"`
}

type ClickwrapDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"true"`
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
