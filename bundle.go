package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/lib"
)

type Bundle struct {
	Code                  string    `json:"code,omitempty"`
	Url                   string    `json:"url,omitempty"`
	Description           string    `json:"description,omitempty"`
	PasswordProtected     *bool     `json:"password_protected,omitempty"`
	RequireRegistration   *bool     `json:"require_registration,omitempty"`
	RequireShareRecipient *bool     `json:"require_share_recipient,omitempty"`
	ClickwrapBody         string    `json:"clickwrap_body,omitempty"`
	FormFieldSet          string    `json:"form_field_set,omitempty"`
	Id                    int64     `json:"id,omitempty"`
	CreatedAt             time.Time `json:"created_at,omitempty"`
	ExpiresAt             time.Time `json:"expires_at,omitempty"`
	MaxUses               int64     `json:"max_uses,omitempty"`
	Note                  string    `json:"note,omitempty"`
	UserId                int64     `json:"user_id,omitempty"`
	Username              string    `json:"username,omitempty"`
	ClickwrapId           int64     `json:"clickwrap_id,omitempty"`
	InboxId               int64     `json:"inbox_id,omitempty"`
	HasInbox              *bool     `json:"has_inbox,omitempty"`
	Paths                 []string  `json:"paths,omitempty"`
	Password              string    `json:"password,omitempty"`
	FormFieldSetId        int64     `json:"form_field_set_id,omitempty"`
}

type BundleCollection []Bundle

type BundleListParams struct {
	UserId     int64           `url:"user_id,omitempty" required:"false"`
	Cursor     string          `url:"cursor,omitempty" required:"false"`
	PerPage    int64           `url:"per_page,omitempty" required:"false"`
	SortBy     json.RawMessage `url:"sort_by,omitempty" required:"false"`
	Filter     json.RawMessage `url:"filter,omitempty" required:"false"`
	FilterGt   json.RawMessage `url:"filter_gt,omitempty" required:"false"`
	FilterGteq json.RawMessage `url:"filter_gteq,omitempty" required:"false"`
	FilterLike json.RawMessage `url:"filter_like,omitempty" required:"false"`
	FilterLt   json.RawMessage `url:"filter_lt,omitempty" required:"false"`
	FilterLteq json.RawMessage `url:"filter_lteq,omitempty" required:"false"`
	lib.ListParams
}

type BundleFindParams struct {
	Id int64 `url:"-,omitempty" required:"true"`
}

type BundleCreateParams struct {
	UserId                int64     `url:"user_id,omitempty" required:"false"`
	Paths                 []string  `url:"paths,omitempty" required:"true"`
	Password              string    `url:"password,omitempty" required:"false"`
	FormFieldSetId        int64     `url:"form_field_set_id,omitempty" required:"false"`
	ExpiresAt             time.Time `url:"expires_at,omitempty" required:"false"`
	MaxUses               int64     `url:"max_uses,omitempty" required:"false"`
	Description           string    `url:"description,omitempty" required:"false"`
	Note                  string    `url:"note,omitempty" required:"false"`
	Code                  string    `url:"code,omitempty" required:"false"`
	RequireRegistration   *bool     `url:"require_registration,omitempty" required:"false"`
	ClickwrapId           int64     `url:"clickwrap_id,omitempty" required:"false"`
	InboxId               int64     `url:"inbox_id,omitempty" required:"false"`
	RequireShareRecipient *bool     `url:"require_share_recipient,omitempty" required:"false"`
}

// Send email(s) with a link to bundle
type BundleShareParams struct {
	Id         int64    `url:"-,omitempty" required:"true"`
	To         []string `url:"to,omitempty" required:"false"`
	Note       string   `url:"note,omitempty" required:"false"`
	Recipients []string `url:"recipients,omitempty" required:"false"`
}

type BundleUpdateParams struct {
	Id                    int64     `url:"-,omitempty" required:"true"`
	Paths                 []string  `url:"paths,omitempty" required:"false"`
	Password              string    `url:"password,omitempty" required:"false"`
	FormFieldSetId        int64     `url:"form_field_set_id,omitempty" required:"false"`
	ClickwrapId           int64     `url:"clickwrap_id,omitempty" required:"false"`
	Code                  string    `url:"code,omitempty" required:"false"`
	Description           string    `url:"description,omitempty" required:"false"`
	ExpiresAt             time.Time `url:"expires_at,omitempty" required:"false"`
	InboxId               int64     `url:"inbox_id,omitempty" required:"false"`
	MaxUses               int64     `url:"max_uses,omitempty" required:"false"`
	Note                  string    `url:"note,omitempty" required:"false"`
	RequireRegistration   *bool     `url:"require_registration,omitempty" required:"false"`
	RequireShareRecipient *bool     `url:"require_share_recipient,omitempty" required:"false"`
}

type BundleDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"true"`
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

func (b *BundleCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*b))
	for i, v := range *b {
		ret[i] = v
	}

	return &ret
}
