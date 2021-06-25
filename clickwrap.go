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

type ClickwrapUseWithBundlesEnum string

func (u ClickwrapUseWithBundlesEnum) String() string {
	return string(u)
}

const (
	NoneUseWithBundles      ClickwrapUseWithBundlesEnum = "none"
	AvailableUseWithBundles ClickwrapUseWithBundlesEnum = "available"
	RequireUseWithBundles   ClickwrapUseWithBundlesEnum = "require"
)

func (u ClickwrapUseWithBundlesEnum) Enum() map[string]ClickwrapUseWithBundlesEnum {
	return map[string]ClickwrapUseWithBundlesEnum{
		"none":      NoneUseWithBundles,
		"available": AvailableUseWithBundles,
		"require":   RequireUseWithBundles,
	}
}

type ClickwrapUseWithInboxesEnum string

func (u ClickwrapUseWithInboxesEnum) String() string {
	return string(u)
}

const (
	NoneUseWithInboxes      ClickwrapUseWithInboxesEnum = "none"
	AvailableUseWithInboxes ClickwrapUseWithInboxesEnum = "available"
	RequireUseWithInboxes   ClickwrapUseWithInboxesEnum = "require"
)

func (u ClickwrapUseWithInboxesEnum) Enum() map[string]ClickwrapUseWithInboxesEnum {
	return map[string]ClickwrapUseWithInboxesEnum{
		"none":      NoneUseWithInboxes,
		"available": AvailableUseWithInboxes,
		"require":   RequireUseWithInboxes,
	}
}

type ClickwrapUseWithUsersEnum string

func (u ClickwrapUseWithUsersEnum) String() string {
	return string(u)
}

const (
	NoneUseWithUsers    ClickwrapUseWithUsersEnum = "none"
	RequireUseWithUsers ClickwrapUseWithUsersEnum = "require"
)

func (u ClickwrapUseWithUsersEnum) Enum() map[string]ClickwrapUseWithUsersEnum {
	return map[string]ClickwrapUseWithUsersEnum{
		"none":    NoneUseWithUsers,
		"require": RequireUseWithUsers,
	}
}

type ClickwrapListParams struct {
	Cursor  string `url:"cursor,omitempty" required:"false"`
	PerPage int64  `url:"per_page,omitempty" required:"false"`
	lib.ListParams
}

type ClickwrapFindParams struct {
	Id int64 `url:"-,omitempty" required:"true"`
}

type ClickwrapCreateParams struct {
	Name           string                      `url:"name,omitempty" required:"false"`
	Body           string                      `url:"body,omitempty" required:"false"`
	UseWithBundles ClickwrapUseWithBundlesEnum `url:"use_with_bundles,omitempty" required:"false"`
	UseWithInboxes ClickwrapUseWithInboxesEnum `url:"use_with_inboxes,omitempty" required:"false"`
	UseWithUsers   ClickwrapUseWithUsersEnum   `url:"use_with_users,omitempty" required:"false"`
}

type ClickwrapUpdateParams struct {
	Id             int64                       `url:"-,omitempty" required:"true"`
	Name           string                      `url:"name,omitempty" required:"false"`
	Body           string                      `url:"body,omitempty" required:"false"`
	UseWithBundles ClickwrapUseWithBundlesEnum `url:"use_with_bundles,omitempty" required:"false"`
	UseWithInboxes ClickwrapUseWithInboxesEnum `url:"use_with_inboxes,omitempty" required:"false"`
	UseWithUsers   ClickwrapUseWithUsersEnum   `url:"use_with_users,omitempty" required:"false"`
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

func (c *ClickwrapCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*c))
	for i, v := range *c {
		ret[i] = v
	}

	return &ret
}
