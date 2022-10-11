package files_sdk

import (
	"encoding/json"
	"io"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Bundle struct {
	Code                            string          `json:"code,omitempty" path:"code"`
	Url                             string          `json:"url,omitempty" path:"url"`
	Description                     string          `json:"description,omitempty" path:"description"`
	PasswordProtected               *bool           `json:"password_protected,omitempty" path:"password_protected"`
	Permissions                     string          `json:"permissions,omitempty" path:"permissions"`
	PreviewOnly                     *bool           `json:"preview_only,omitempty" path:"preview_only"`
	RequireRegistration             *bool           `json:"require_registration,omitempty" path:"require_registration"`
	RequireShareRecipient           *bool           `json:"require_share_recipient,omitempty" path:"require_share_recipient"`
	ClickwrapBody                   string          `json:"clickwrap_body,omitempty" path:"clickwrap_body"`
	FormFieldSet                    FormFieldSet    `json:"form_field_set,omitempty" path:"form_field_set"`
	SkipName                        *bool           `json:"skip_name,omitempty" path:"skip_name"`
	SkipEmail                       *bool           `json:"skip_email,omitempty" path:"skip_email"`
	SkipCompany                     *bool           `json:"skip_company,omitempty" path:"skip_company"`
	Id                              int64           `json:"id,omitempty" path:"id"`
	CreatedAt                       *time.Time      `json:"created_at,omitempty" path:"created_at"`
	DontSeparateSubmissionsByFolder *bool           `json:"dont_separate_submissions_by_folder,omitempty" path:"dont_separate_submissions_by_folder"`
	ExpiresAt                       *time.Time      `json:"expires_at,omitempty" path:"expires_at"`
	MaxUses                         int64           `json:"max_uses,omitempty" path:"max_uses"`
	Note                            string          `json:"note,omitempty" path:"note"`
	PathTemplate                    string          `json:"path_template,omitempty" path:"path_template"`
	UserId                          int64           `json:"user_id,omitempty" path:"user_id"`
	Username                        string          `json:"username,omitempty" path:"username"`
	ClickwrapId                     int64           `json:"clickwrap_id,omitempty" path:"clickwrap_id"`
	InboxId                         int64           `json:"inbox_id,omitempty" path:"inbox_id"`
	WatermarkAttachment             Image           `json:"watermark_attachment,omitempty" path:"watermark_attachment"`
	WatermarkValue                  json.RawMessage `json:"watermark_value,omitempty" path:"watermark_value"`
	HasInbox                        *bool           `json:"has_inbox,omitempty" path:"has_inbox"`
	Paths                           []string        `json:"paths,omitempty" path:"paths"`
	Password                        string          `json:"password,omitempty" path:"password"`
	FormFieldSetId                  int64           `json:"form_field_set_id,omitempty" path:"form_field_set_id"`
	WatermarkAttachmentFile         io.Reader       `json:"watermark_attachment_file,omitempty" path:"watermark_attachment_file"`
	WatermarkAttachmentDelete       *bool           `json:"watermark_attachment_delete,omitempty" path:"watermark_attachment_delete"`
}

type BundleCollection []Bundle

type BundlePermissionsEnum string

func (u BundlePermissionsEnum) String() string {
	return string(u)
}

func (u BundlePermissionsEnum) Enum() map[string]BundlePermissionsEnum {
	return map[string]BundlePermissionsEnum{
		"read":         BundlePermissionsEnum("read"),
		"write":        BundlePermissionsEnum("write"),
		"read_write":   BundlePermissionsEnum("read_write"),
		"full":         BundlePermissionsEnum("full"),
		"none":         BundlePermissionsEnum("none"),
		"preview_only": BundlePermissionsEnum("preview_only"),
	}
}

type BundleListParams struct {
	UserId     int64           `url:"user_id,omitempty" required:"false" json:"user_id,omitempty" path:"user_id"`
	SortBy     json.RawMessage `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty" path:"sort_by"`
	Filter     json.RawMessage `url:"filter,omitempty" required:"false" json:"filter,omitempty" path:"filter"`
	FilterGt   json.RawMessage `url:"filter_gt,omitempty" required:"false" json:"filter_gt,omitempty" path:"filter_gt"`
	FilterGteq json.RawMessage `url:"filter_gteq,omitempty" required:"false" json:"filter_gteq,omitempty" path:"filter_gteq"`
	FilterLike json.RawMessage `url:"filter_like,omitempty" required:"false" json:"filter_like,omitempty" path:"filter_like"`
	FilterLt   json.RawMessage `url:"filter_lt,omitempty" required:"false" json:"filter_lt,omitempty" path:"filter_lt"`
	FilterLteq json.RawMessage `url:"filter_lteq,omitempty" required:"false" json:"filter_lteq,omitempty" path:"filter_lteq"`
	lib.ListParams
}

type BundleFindParams struct {
	Id int64 `url:"-,omitempty" required:"true" json:"-,omitempty" path:"id"`
}

type BundleCreateParams struct {
	UserId                          int64                 `url:"user_id,omitempty" required:"false" json:"user_id,omitempty" path:"user_id"`
	Paths                           []string              `url:"paths,omitempty" required:"true" json:"paths,omitempty" path:"paths"`
	Password                        string                `url:"password,omitempty" required:"false" json:"password,omitempty" path:"password"`
	FormFieldSetId                  int64                 `url:"form_field_set_id,omitempty" required:"false" json:"form_field_set_id,omitempty" path:"form_field_set_id"`
	DontSeparateSubmissionsByFolder *bool                 `url:"dont_separate_submissions_by_folder,omitempty" required:"false" json:"dont_separate_submissions_by_folder,omitempty" path:"dont_separate_submissions_by_folder"`
	ExpiresAt                       *time.Time            `url:"expires_at,omitempty" required:"false" json:"expires_at,omitempty" path:"expires_at"`
	MaxUses                         int64                 `url:"max_uses,omitempty" required:"false" json:"max_uses,omitempty" path:"max_uses"`
	Description                     string                `url:"description,omitempty" required:"false" json:"description,omitempty" path:"description"`
	Note                            string                `url:"note,omitempty" required:"false" json:"note,omitempty" path:"note"`
	Code                            string                `url:"code,omitempty" required:"false" json:"code,omitempty" path:"code"`
	PathTemplate                    string                `url:"path_template,omitempty" required:"false" json:"path_template,omitempty" path:"path_template"`
	Permissions                     BundlePermissionsEnum `url:"permissions,omitempty" required:"false" json:"permissions,omitempty" path:"permissions"`
	PreviewOnly                     *bool                 `url:"preview_only,omitempty" required:"false" json:"preview_only,omitempty" path:"preview_only"`
	RequireRegistration             *bool                 `url:"require_registration,omitempty" required:"false" json:"require_registration,omitempty" path:"require_registration"`
	ClickwrapId                     int64                 `url:"clickwrap_id,omitempty" required:"false" json:"clickwrap_id,omitempty" path:"clickwrap_id"`
	InboxId                         int64                 `url:"inbox_id,omitempty" required:"false" json:"inbox_id,omitempty" path:"inbox_id"`
	RequireShareRecipient           *bool                 `url:"require_share_recipient,omitempty" required:"false" json:"require_share_recipient,omitempty" path:"require_share_recipient"`
	SkipEmail                       *bool                 `url:"skip_email,omitempty" required:"false" json:"skip_email,omitempty" path:"skip_email"`
	SkipName                        *bool                 `url:"skip_name,omitempty" required:"false" json:"skip_name,omitempty" path:"skip_name"`
	SkipCompany                     *bool                 `url:"skip_company,omitempty" required:"false" json:"skip_company,omitempty" path:"skip_company"`
	WatermarkAttachmentFile         io.Writer             `url:"watermark_attachment_file,omitempty" required:"false" json:"watermark_attachment_file,omitempty" path:"watermark_attachment_file"`
}

// Send email(s) with a link to bundle
type BundleShareParams struct {
	Id         int64           `url:"-,omitempty" required:"true" json:"-,omitempty" path:"id"`
	To         []string        `url:"to,omitempty" required:"false" json:"to,omitempty" path:"to"`
	Note       string          `url:"note,omitempty" required:"false" json:"note,omitempty" path:"note"`
	Recipients json.RawMessage `url:"recipients,omitempty" required:"false" json:"recipients,omitempty" path:"recipients"`
}

type BundleUpdateParams struct {
	Id                              int64                 `url:"-,omitempty" required:"true" json:"-,omitempty" path:"id"`
	Paths                           []string              `url:"paths,omitempty" required:"false" json:"paths,omitempty" path:"paths"`
	Password                        string                `url:"password,omitempty" required:"false" json:"password,omitempty" path:"password"`
	FormFieldSetId                  int64                 `url:"form_field_set_id,omitempty" required:"false" json:"form_field_set_id,omitempty" path:"form_field_set_id"`
	ClickwrapId                     int64                 `url:"clickwrap_id,omitempty" required:"false" json:"clickwrap_id,omitempty" path:"clickwrap_id"`
	Code                            string                `url:"code,omitempty" required:"false" json:"code,omitempty" path:"code"`
	Description                     string                `url:"description,omitempty" required:"false" json:"description,omitempty" path:"description"`
	DontSeparateSubmissionsByFolder *bool                 `url:"dont_separate_submissions_by_folder,omitempty" required:"false" json:"dont_separate_submissions_by_folder,omitempty" path:"dont_separate_submissions_by_folder"`
	ExpiresAt                       *time.Time            `url:"expires_at,omitempty" required:"false" json:"expires_at,omitempty" path:"expires_at"`
	InboxId                         int64                 `url:"inbox_id,omitempty" required:"false" json:"inbox_id,omitempty" path:"inbox_id"`
	MaxUses                         int64                 `url:"max_uses,omitempty" required:"false" json:"max_uses,omitempty" path:"max_uses"`
	Note                            string                `url:"note,omitempty" required:"false" json:"note,omitempty" path:"note"`
	PathTemplate                    string                `url:"path_template,omitempty" required:"false" json:"path_template,omitempty" path:"path_template"`
	Permissions                     BundlePermissionsEnum `url:"permissions,omitempty" required:"false" json:"permissions,omitempty" path:"permissions"`
	PreviewOnly                     *bool                 `url:"preview_only,omitempty" required:"false" json:"preview_only,omitempty" path:"preview_only"`
	RequireRegistration             *bool                 `url:"require_registration,omitempty" required:"false" json:"require_registration,omitempty" path:"require_registration"`
	RequireShareRecipient           *bool                 `url:"require_share_recipient,omitempty" required:"false" json:"require_share_recipient,omitempty" path:"require_share_recipient"`
	SkipEmail                       *bool                 `url:"skip_email,omitempty" required:"false" json:"skip_email,omitempty" path:"skip_email"`
	SkipName                        *bool                 `url:"skip_name,omitempty" required:"false" json:"skip_name,omitempty" path:"skip_name"`
	SkipCompany                     *bool                 `url:"skip_company,omitempty" required:"false" json:"skip_company,omitempty" path:"skip_company"`
	WatermarkAttachmentDelete       *bool                 `url:"watermark_attachment_delete,omitempty" required:"false" json:"watermark_attachment_delete,omitempty" path:"watermark_attachment_delete"`
	WatermarkAttachmentFile         io.Writer             `url:"watermark_attachment_file,omitempty" required:"false" json:"watermark_attachment_file,omitempty" path:"watermark_attachment_file"`
}

type BundleDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"true" json:"-,omitempty" path:"id"`
}

func (b *Bundle) UnmarshalJSON(data []byte) error {
	type bundle Bundle
	var v bundle
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*b = Bundle(v)
	return nil
}

func (b *BundleCollection) UnmarshalJSON(data []byte) error {
	type bundles BundleCollection
	var v bundles
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
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
