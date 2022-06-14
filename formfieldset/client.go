package form_field_set

import (
	"context"
	"strconv"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	lib "github.com/Files-com/files-sdk-go/v2/lib"
	listquery "github.com/Files-com/files-sdk-go/v2/listquery"
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

func (c *Client) List(ctx context.Context, params files_sdk.FormFieldSetListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/form_field_sets"
	i.ListParams = &params
	list := files_sdk.FormFieldSetCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.FormFieldSetListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}

func (c *Client) Find(ctx context.Context, params files_sdk.FormFieldSetFindParams) (files_sdk.FormFieldSet, error) {
	formFieldSet := files_sdk.FormFieldSet{}
	if params.Id == 0 {
		return formFieldSet, lib.CreateError(params, "Id")
	}
	path := "/form_field_sets/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "GET", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return formFieldSet, err
	}
	if res.StatusCode == 204 {
		return formFieldSet, nil
	}

	return formFieldSet, formFieldSet.UnmarshalJSON(*data)
}

func Find(ctx context.Context, params files_sdk.FormFieldSetFindParams) (files_sdk.FormFieldSet, error) {
	return (&Client{}).Find(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.FormFieldSetCreateParams) (files_sdk.FormFieldSet, error) {
	formFieldSet := files_sdk.FormFieldSet{}
	path := "/form_field_sets"
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return formFieldSet, err
	}
	if res.StatusCode == 204 {
		return formFieldSet, nil
	}

	return formFieldSet, formFieldSet.UnmarshalJSON(*data)
}

func Create(ctx context.Context, params files_sdk.FormFieldSetCreateParams) (files_sdk.FormFieldSet, error) {
	return (&Client{}).Create(ctx, params)
}

func (c *Client) Update(ctx context.Context, params files_sdk.FormFieldSetUpdateParams) (files_sdk.FormFieldSet, error) {
	formFieldSet := files_sdk.FormFieldSet{}
	if params.Id == 0 {
		return formFieldSet, lib.CreateError(params, "Id")
	}
	path := "/form_field_sets/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "PATCH", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return formFieldSet, err
	}
	if res.StatusCode == 204 {
		return formFieldSet, nil
	}

	return formFieldSet, formFieldSet.UnmarshalJSON(*data)
}

func Update(ctx context.Context, params files_sdk.FormFieldSetUpdateParams) (files_sdk.FormFieldSet, error) {
	return (&Client{}).Update(ctx, params)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.FormFieldSetDeleteParams) error {
	formFieldSet := files_sdk.FormFieldSet{}
	if params.Id == 0 {
		return lib.CreateError(params, "Id")
	}
	path := "/form_field_sets/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "DELETE", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return err
	}
	if res.StatusCode == 204 {
		return nil
	}

	return formFieldSet.UnmarshalJSON(*data)
}

func Delete(ctx context.Context, params files_sdk.FormFieldSetDeleteParams) error {
	return (&Client{}).Delete(ctx, params)
}
