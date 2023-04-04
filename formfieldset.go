package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type FormFieldSet struct {
	Id          int64    `json:"id,omitempty" path:"id"`
	Title       string   `json:"title,omitempty" path:"title"`
	FormLayout  []int64  `json:"form_layout,omitempty" path:"form_layout"`
	FormFields  []string `json:"form_fields,omitempty" path:"form_fields"`
	SkipName    *bool    `json:"skip_name,omitempty" path:"skip_name"`
	SkipEmail   *bool    `json:"skip_email,omitempty" path:"skip_email"`
	SkipCompany *bool    `json:"skip_company,omitempty" path:"skip_company"`
	UserId      int64    `json:"user_id,omitempty" path:"user_id"`
}

type FormFieldSetCollection []FormFieldSet

type FormFieldSetListParams struct {
	UserId int64 `url:"user_id,omitempty" required:"false" json:"user_id,omitempty" path:"user_id"`
	lib.ListParams
}

type FormFieldSetFindParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
}

type FormFieldSetCreateParams struct {
	UserId      int64           `url:"user_id,omitempty" required:"" json:"user_id,omitempty" path:"user_id"`
	Title       string          `url:"title,omitempty" required:"" json:"title,omitempty" path:"title"`
	SkipEmail   *bool           `url:"skip_email,omitempty" required:"" json:"skip_email,omitempty" path:"skip_email"`
	SkipName    *bool           `url:"skip_name,omitempty" required:"" json:"skip_name,omitempty" path:"skip_name"`
	SkipCompany *bool           `url:"skip_company,omitempty" required:"" json:"skip_company,omitempty" path:"skip_company"`
	FormFields  json.RawMessage `url:"form_fields,omitempty" required:"" json:"form_fields,omitempty" path:"form_fields"`
}

type FormFieldSetUpdateParams struct {
	Id          int64           `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
	Title       string          `url:"title,omitempty" required:"" json:"title,omitempty" path:"title"`
	SkipEmail   *bool           `url:"skip_email,omitempty" required:"" json:"skip_email,omitempty" path:"skip_email"`
	SkipName    *bool           `url:"skip_name,omitempty" required:"" json:"skip_name,omitempty" path:"skip_name"`
	SkipCompany *bool           `url:"skip_company,omitempty" required:"" json:"skip_company,omitempty" path:"skip_company"`
	FormFields  json.RawMessage `url:"form_fields,omitempty" required:"" json:"form_fields,omitempty" path:"form_fields"`
}

type FormFieldSetDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
}

func (f *FormFieldSet) UnmarshalJSON(data []byte) error {
	type formFieldSet FormFieldSet
	var v formFieldSet
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*f = FormFieldSet(v)
	return nil
}

func (f *FormFieldSetCollection) UnmarshalJSON(data []byte) error {
	type formFieldSets FormFieldSetCollection
	var v formFieldSets
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*f = FormFieldSetCollection(v)
	return nil
}

func (f *FormFieldSetCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*f))
	for i, v := range *f {
		ret[i] = v
	}

	return &ret
}
