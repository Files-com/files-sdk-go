package files_sdk

import (
	"encoding/json"
	"io"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Bundle struct {
	Code                      string          `json:"code,omitempty"`
	Url                       string          `json:"url,omitempty"`
	Description               string          `json:"description,omitempty"`
	PasswordProtected         *bool           `json:"password_protected,omitempty"`
	PreviewOnly               *bool           `json:"preview_only,omitempty"`
	RequireRegistration       *bool           `json:"require_registration,omitempty"`
	RequireShareRecipient     *bool           `json:"require_share_recipient,omitempty"`
	ClickwrapBody             string          `json:"clickwrap_body,omitempty"`
	FormFieldSet              FormFieldSet    `json:"form_field_set,omitempty"`
	Id                        int64           `json:"id,omitempty"`
	CreatedAt                 time.Time       `json:"created_at,omitempty"`
	ExpiresAt                 time.Time       `json:"expires_at,omitempty"`
	MaxUses                   int64           `json:"max_uses,omitempty"`
	Note                      string          `json:"note,omitempty"`
	UserId                    int64           `json:"user_id,omitempty"`
	Username                  string          `json:"username,omitempty"`
	ClickwrapId               int64           `json:"clickwrap_id,omitempty"`
	InboxId                   int64           `json:"inbox_id,omitempty"`
	WatermarkAttachment       Image           `json:"watermark_attachment,omitempty"`
	WatermarkValue            json.RawMessage `json:"watermark_value,omitempty"`
	HasInbox                  *bool           `json:"has_inbox,omitempty"`
	Paths                     []string        `json:"paths,omitempty"`
	Password                  string          `json:"password,omitempty"`
	FormFieldSetId            int64           `json:"form_field_set_id,omitempty"`
	WatermarkAttachmentFile   io.Reader       `json:"watermark_attachment_file,omitempty"`
	WatermarkAttachmentDelete *bool           `json:"watermark_attachment_delete,omitempty"`
}

type BundleCollection []Bundle

type BundleListParams struct {
	UserId     int64           `url:"user_id,omitempty" required:"false" json:"user_id,omitempty"`
	Cursor     string          `url:"cursor,omitempty" required:"false" json:"cursor,omitempty"`
	PerPage    int64           `url:"per_page,omitempty" required:"false" json:"per_page,omitempty"`
	SortBy     json.RawMessage `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty"`
	Filter     json.RawMessage `url:"filter,omitempty" required:"false" json:"filter,omitempty"`
	FilterGt   json.RawMessage `url:"filter_gt,omitempty" required:"false" json:"filter_gt,omitempty"`
	FilterGteq json.RawMessage `url:"filter_gteq,omitempty" required:"false" json:"filter_gteq,omitempty"`
	FilterLike json.RawMessage `url:"filter_like,omitempty" required:"false" json:"filter_like,omitempty"`
	FilterLt   json.RawMessage `url:"filter_lt,omitempty" required:"false" json:"filter_lt,omitempty"`
	FilterLteq json.RawMessage `url:"filter_lteq,omitempty" required:"false" json:"filter_lteq,omitempty"`
	lib.ListParams
}

type BundleFindParams struct {
	Id int64 `url:"-,omitempty" required:"true" json:"-,omitempty"`
}

type BundleCreateParams struct {
	UserId                  int64     `url:"user_id,omitempty" required:"false" json:"user_id,omitempty"`
	Paths                   []string  `url:"paths,omitempty" required:"true" json:"paths,omitempty"`
	Password                string    `url:"password,omitempty" required:"false" json:"password,omitempty"`
	FormFieldSetId          int64     `url:"form_field_set_id,omitempty" required:"false" json:"form_field_set_id,omitempty"`
	ExpiresAt               time.Time `url:"expires_at,omitempty" required:"false" json:"expires_at,omitempty"`
	MaxUses                 int64     `url:"max_uses,omitempty" required:"false" json:"max_uses,omitempty"`
	Description             string    `url:"description,omitempty" required:"false" json:"description,omitempty"`
	Note                    string    `url:"note,omitempty" required:"false" json:"note,omitempty"`
	Code                    string    `url:"code,omitempty" required:"false" json:"code,omitempty"`
	PreviewOnly             *bool     `url:"preview_only,omitempty" required:"false" json:"preview_only,omitempty"`
	RequireRegistration     *bool     `url:"require_registration,omitempty" required:"false" json:"require_registration,omitempty"`
	ClickwrapId             int64     `url:"clickwrap_id,omitempty" required:"false" json:"clickwrap_id,omitempty"`
	InboxId                 int64     `url:"inbox_id,omitempty" required:"false" json:"inbox_id,omitempty"`
	RequireShareRecipient   *bool     `url:"require_share_recipient,omitempty" required:"false" json:"require_share_recipient,omitempty"`
	WatermarkAttachmentFile io.Writer `url:"watermark_attachment_file,omitempty" required:"false" json:"watermark_attachment_file,omitempty"`
}

// Send email(s) with a link to bundle
type BundleShareParams struct {
	Id         int64    `url:"-,omitempty" required:"true" json:"-,omitempty"`
	To         []string `url:"to,omitempty" required:"false" json:"to,omitempty"`
	Note       string   `url:"note,omitempty" required:"false" json:"note,omitempty"`
	Recipients []string `url:"recipients,omitempty" required:"false" json:"recipients,omitempty"`
}

type BundleUpdateParams struct {
	Id                        int64     `url:"-,omitempty" required:"true" json:"-,omitempty"`
	Paths                     []string  `url:"paths,omitempty" required:"false" json:"paths,omitempty"`
	Password                  string    `url:"password,omitempty" required:"false" json:"password,omitempty"`
	FormFieldSetId            int64     `url:"form_field_set_id,omitempty" required:"false" json:"form_field_set_id,omitempty"`
	ClickwrapId               int64     `url:"clickwrap_id,omitempty" required:"false" json:"clickwrap_id,omitempty"`
	Code                      string    `url:"code,omitempty" required:"false" json:"code,omitempty"`
	Description               string    `url:"description,omitempty" required:"false" json:"description,omitempty"`
	ExpiresAt                 time.Time `url:"expires_at,omitempty" required:"false" json:"expires_at,omitempty"`
	InboxId                   int64     `url:"inbox_id,omitempty" required:"false" json:"inbox_id,omitempty"`
	MaxUses                   int64     `url:"max_uses,omitempty" required:"false" json:"max_uses,omitempty"`
	Note                      string    `url:"note,omitempty" required:"false" json:"note,omitempty"`
	PreviewOnly               *bool     `url:"preview_only,omitempty" required:"false" json:"preview_only,omitempty"`
	RequireRegistration       *bool     `url:"require_registration,omitempty" required:"false" json:"require_registration,omitempty"`
	RequireShareRecipient     *bool     `url:"require_share_recipient,omitempty" required:"false" json:"require_share_recipient,omitempty"`
	WatermarkAttachmentDelete *bool     `url:"watermark_attachment_delete,omitempty" required:"false" json:"watermark_attachment_delete,omitempty"`
	WatermarkAttachmentFile   io.Writer `url:"watermark_attachment_file,omitempty" required:"false" json:"watermark_attachment_file,omitempty"`
}

type BundleDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"true" json:"-,omitempty"`
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
