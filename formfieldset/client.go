package form_field_set

import (
	"strconv"

	files_sdk "github.com/Files-com/files-sdk-go"
	lib "github.com/Files-com/files-sdk-go/lib"
)

type Client struct {
	files_sdk.Config
}

type Iter struct {
	*lib.Iter
}

func (i *Iter) FormFieldSet() files_sdk.FormFieldSet {
	return i.Current().(files_sdk.FormFieldSet)
}

func (c *Client) List(params files_sdk.FormFieldSetListParams) (*Iter, error) {
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	i := &Iter{Iter: &lib.Iter{}}
	path := "/form_field_sets"
	i.ListParams = &params
	exportParams, err := i.ExportParams()
	if err != nil {
		return i, err
	}
	i.Query = func() (*[]interface{}, string, error) {
		data, res, err := files_sdk.Call("GET", c.Config, path, exportParams)
		defaultValue := make([]interface{}, 0)
		if err != nil {
			return &defaultValue, "", err
		}
		list := files_sdk.FormFieldSetCollection{}
		if err := list.UnmarshalJSON(*data); err != nil {
			return &defaultValue, "", err
		}

		ret := make([]interface{}, len(list))
		for i, v := range list {
			ret[i] = v
		}
		cursor := res.Header.Get("X-Files-Cursor")
		return &ret, cursor, nil
	}
	return i, nil
}

func List(params files_sdk.FormFieldSetListParams) (*Iter, error) {
	return (&Client{}).List(params)
}

func (c *Client) Find(params files_sdk.FormFieldSetFindParams) (files_sdk.FormFieldSet, error) {
	formFieldSet := files_sdk.FormFieldSet{}
	if params.Id == 0 {
		return formFieldSet, lib.CreateError(params, "Id")
	}
	path := "/form_field_sets/" + lib.QueryEscape(strconv.FormatInt(params.Id, 10)) + ""
	exportedParms, err := lib.ExportParams(params)
	if err != nil {
		return formFieldSet, err
	}
	data, res, err := files_sdk.Call("GET", c.Config, path, exportedParms)
	if err != nil {
		return formFieldSet, err
	}
	if res.StatusCode == 204 {
		return formFieldSet, nil
	}
	if err := formFieldSet.UnmarshalJSON(*data); err != nil {
		return formFieldSet, err
	}

	return formFieldSet, nil
}

func Find(params files_sdk.FormFieldSetFindParams) (files_sdk.FormFieldSet, error) {
	return (&Client{}).Find(params)
}

func (c *Client) Create(params files_sdk.FormFieldSetCreateParams) (files_sdk.FormFieldSet, error) {
	formFieldSet := files_sdk.FormFieldSet{}
	path := "/form_field_sets"
	exportedParms, err := lib.ExportParams(params)
	if err != nil {
		return formFieldSet, err
	}
	data, res, err := files_sdk.Call("POST", c.Config, path, exportedParms)
	if err != nil {
		return formFieldSet, err
	}
	if res.StatusCode == 204 {
		return formFieldSet, nil
	}
	if err := formFieldSet.UnmarshalJSON(*data); err != nil {
		return formFieldSet, err
	}

	return formFieldSet, nil
}

func Create(params files_sdk.FormFieldSetCreateParams) (files_sdk.FormFieldSet, error) {
	return (&Client{}).Create(params)
}

func (c *Client) Update(params files_sdk.FormFieldSetUpdateParams) (files_sdk.FormFieldSet, error) {
	formFieldSet := files_sdk.FormFieldSet{}
	if params.Id == 0 {
		return formFieldSet, lib.CreateError(params, "Id")
	}
	path := "/form_field_sets/" + lib.QueryEscape(strconv.FormatInt(params.Id, 10)) + ""
	exportedParms, err := lib.ExportParams(params)
	if err != nil {
		return formFieldSet, err
	}
	data, res, err := files_sdk.Call("PATCH", c.Config, path, exportedParms)
	if err != nil {
		return formFieldSet, err
	}
	if res.StatusCode == 204 {
		return formFieldSet, nil
	}
	if err := formFieldSet.UnmarshalJSON(*data); err != nil {
		return formFieldSet, err
	}

	return formFieldSet, nil
}

func Update(params files_sdk.FormFieldSetUpdateParams) (files_sdk.FormFieldSet, error) {
	return (&Client{}).Update(params)
}

func (c *Client) Delete(params files_sdk.FormFieldSetDeleteParams) (files_sdk.FormFieldSet, error) {
	formFieldSet := files_sdk.FormFieldSet{}
	if params.Id == 0 {
		return formFieldSet, lib.CreateError(params, "Id")
	}
	path := "/form_field_sets/" + lib.QueryEscape(strconv.FormatInt(params.Id, 10)) + ""
	exportedParms, err := lib.ExportParams(params)
	if err != nil {
		return formFieldSet, err
	}
	data, res, err := files_sdk.Call("DELETE", c.Config, path, exportedParms)
	if err != nil {
		return formFieldSet, err
	}
	if res.StatusCode == 204 {
		return formFieldSet, nil
	}
	if err := formFieldSet.UnmarshalJSON(*data); err != nil {
		return formFieldSet, err
	}

	return formFieldSet, nil
}

func Delete(params files_sdk.FormFieldSetDeleteParams) (files_sdk.FormFieldSet, error) {
	return (&Client{}).Delete(params)
}
