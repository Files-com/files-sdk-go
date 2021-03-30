package form_field_set

import (
	"strconv"

	files_sdk "github.com/Files-com/files-sdk-go"
	lib "github.com/Files-com/files-sdk-go/lib"
	listquery "github.com/Files-com/files-sdk-go/listquery"
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
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/form_field_sets"
	i.ListParams = &params
	list := files_sdk.FormFieldSetCollection{}
	i.Query = listquery.Build(i, c.Config, path, &list)
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
	path := "/form_field_sets/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return formFieldSet, err
	}
	data, res, err := files_sdk.Call("GET", c.Config, path, exportedParams)
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
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return formFieldSet, err
	}
	data, res, err := files_sdk.Call("POST", c.Config, path, exportedParams)
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
	path := "/form_field_sets/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return formFieldSet, err
	}
	data, res, err := files_sdk.Call("PATCH", c.Config, path, exportedParams)
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
	path := "/form_field_sets/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return formFieldSet, err
	}
	data, res, err := files_sdk.Call("DELETE", c.Config, path, exportedParams)
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
