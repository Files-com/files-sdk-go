package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/lib"
)

type FormFieldSet struct {
	Id          int64    `json:"id,omitempty"`
	Title       string   `json:"title,omitempty"`
	FormLayout  int      `json:"form_layout,omitempty"`
	FormFields  []string `json:"form_fields,omitempty"`
	SkipName    *bool    `json:"skip_name,omitempty"`
	SkipEmail   *bool    `json:"skip_email,omitempty"`
	SkipCompany *bool    `json:"skip_company,omitempty"`
	UserId      int64    `json:"user_id,omitempty"`
}

type FormFieldSetCollection []FormFieldSet

type FormFieldSetListParams struct {
	UserId  int64  `url:"user_id,omitempty" required:"false"`
	Cursor  string `url:"cursor,omitempty" required:"false"`
	PerPage int    `url:"per_page,omitempty" required:"false"`
	lib.ListParams
}

type FormFieldSetFindParams struct {
	Id int64 `url:"-,omitempty" required:"true"`
}

type FormFieldSetCreateParams struct {
	UserId      int64    `url:"user_id,omitempty" required:""`
	Title       string   `url:"title,omitempty" required:""`
	SkipEmail   *bool    `url:"skip_email,omitempty" required:""`
	SkipName    *bool    `url:"skip_name,omitempty" required:""`
	SkipCompany *bool    `url:"skip_company,omitempty" required:""`
	FormFields  []string `url:"form_fields,omitempty" required:""`
}

type FormFieldSetUpdateParams struct {
	Id          int64    `url:"-,omitempty" required:"true"`
	Title       string   `url:"title,omitempty" required:""`
	SkipEmail   *bool    `url:"skip_email,omitempty" required:""`
	SkipName    *bool    `url:"skip_name,omitempty" required:""`
	SkipCompany *bool    `url:"skip_company,omitempty" required:""`
	FormFields  []string `url:"form_fields,omitempty" required:""`
}

type FormFieldSetDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"true"`
}

func (f *FormFieldSet) UnmarshalJSON(data []byte) error {
	type formFieldSet FormFieldSet
	var v formFieldSet
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*f = FormFieldSet(v)
	return nil
}

func (f *FormFieldSetCollection) UnmarshalJSON(data []byte) error {
	type formFieldSets []FormFieldSet
	var v formFieldSets
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*f = FormFieldSetCollection(v)
	return nil
}
