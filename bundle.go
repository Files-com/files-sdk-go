package files_sdk

import (
	"encoding/json"
	lib "github.com/Files-com/files-sdk-go/lib"
	"time"
)

type Bundle struct {
	Code                  string    `json:"code,omitempty"`
	Url                   string    `json:"url,omitempty"`
	Description           string    `json:"description,omitempty"`
	PasswordProtected     *bool     `json:"password_protected,omitempty"`
	RequireRegistration   *bool     `json:"require_registration,omitempty"`
	RequireShareRecipient *bool     `json:"require_share_recipient,omitempty"`
	ClickwrapBody         string    `json:"clickwrap_body,omitempty"`
	Id                    int       `json:"id,omitempty"`
	CreatedAt             time.Time `json:"created_at,omitempty"`
	ExpiresAt             time.Time `json:"expires_at,omitempty"`
	MaxUses               int       `json:"max_uses,omitempty"`
	Note                  string    `json:"note,omitempty"`
	UserId                int       `json:"user_id,omitempty"`
	Username              string    `json:"username,omitempty"`
	ClickwrapId           int       `json:"clickwrap_id,omitempty"`
	InboxId               int       `json:"inbox_id,omitempty"`
	Paths                 []string  `json:"paths,omitempty"`
	Password              string    `json:"password,omitempty"`
}

type BundleCollection []Bundle

type BundleListParams struct {
	UserId     int             `url:"user_id,omitempty"`
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

type BundleFindParams struct {
	Id int `url:"-,omitempty"`
}

type BundleCreateParams struct {
	UserId                int      `url:"user_id,omitempty"`
	Paths                 []string `url:"paths,omitempty"`
	Password              string   `url:"password,omitempty"`
	ExpiresAt             string   `url:"expires_at,omitempty"`
	MaxUses               int      `url:"max_uses,omitempty"`
	Description           string   `url:"description,omitempty"`
	Note                  string   `url:"note,omitempty"`
	Code                  string   `url:"code,omitempty"`
	RequireRegistration   *bool    `url:"require_registration,omitempty"`
	ClickwrapId           int      `url:"clickwrap_id,omitempty"`
	InboxId               int      `url:"inbox_id,omitempty"`
	RequireShareRecipient *bool    `url:"require_share_recipient,omitempty"`
}

type BundleShareParams struct {
	Id   int      `url:"-,omitempty"`
	To   []string `url:"to,omitempty"`
	Note string   `url:"note,omitempty"`
}

type BundleUpdateParams struct {
	Id                    int    `url:"-,omitempty"`
	Password              string `url:"password,omitempty"`
	ClickwrapId           int    `url:"clickwrap_id,omitempty"`
	Code                  string `url:"code,omitempty"`
	Description           string `url:"description,omitempty"`
	ExpiresAt             string `url:"expires_at,omitempty"`
	InboxId               int    `url:"inbox_id,omitempty"`
	MaxUses               int    `url:"max_uses,omitempty"`
	Note                  string `url:"note,omitempty"`
	RequireRegistration   *bool  `url:"require_registration,omitempty"`
	RequireShareRecipient *bool  `url:"require_share_recipient,omitempty"`
}

type BundleDeleteParams struct {
	Id int `url:"-,omitempty"`
}

func (b *Bundle) UnmarshalJSON(data []byte) error {
	type bundle Bundle
	var v bundle
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*b = Bundle(v)
	return nil
}

func (b *BundleCollection) UnmarshalJSON(data []byte) error {
	type bundles []Bundle
	var v bundles
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*b = BundleCollection(v)
	return nil
}
