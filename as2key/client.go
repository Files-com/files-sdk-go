package as2_key

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

func (i *Iter) As2Key() files_sdk.As2Key {
	return i.Current().(files_sdk.As2Key)
}

func (c *Client) List(ctx context.Context, params files_sdk.As2KeyListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/as2_keys"
	i.ListParams = &params
	list := files_sdk.As2KeyCollection{}
	i.Query = listquery.Build(ctx, i, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.As2KeyListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}

func (c *Client) Find(ctx context.Context, params files_sdk.As2KeyFindParams) (files_sdk.As2Key, error) {
	as2Key := files_sdk.As2Key{}
	if params.Id == 0 {
		return as2Key, lib.CreateError(params, "Id")
	}
	path := "/as2_keys/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return as2Key, err
	}
	data, res, err := files_sdk.Call(ctx, "GET", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return as2Key, err
	}
	if res.StatusCode == 204 {
		return as2Key, nil
	}
	if err := as2Key.UnmarshalJSON(*data); err != nil {
		return as2Key, err
	}

	return as2Key, nil
}

func Find(ctx context.Context, params files_sdk.As2KeyFindParams) (files_sdk.As2Key, error) {
	return (&Client{}).Find(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.As2KeyCreateParams) (files_sdk.As2Key, error) {
	as2Key := files_sdk.As2Key{}
	path := "/as2_keys"
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return as2Key, err
	}
	data, res, err := files_sdk.Call(ctx, "POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return as2Key, err
	}
	if res.StatusCode == 204 {
		return as2Key, nil
	}
	if err := as2Key.UnmarshalJSON(*data); err != nil {
		return as2Key, err
	}

	return as2Key, nil
}

func Create(ctx context.Context, params files_sdk.As2KeyCreateParams) (files_sdk.As2Key, error) {
	return (&Client{}).Create(ctx, params)
}

func (c *Client) Update(ctx context.Context, params files_sdk.As2KeyUpdateParams) (files_sdk.As2Key, error) {
	as2Key := files_sdk.As2Key{}
	if params.Id == 0 {
		return as2Key, lib.CreateError(params, "Id")
	}
	path := "/as2_keys/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return as2Key, err
	}
	data, res, err := files_sdk.Call(ctx, "PATCH", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return as2Key, err
	}
	if res.StatusCode == 204 {
		return as2Key, nil
	}
	if err := as2Key.UnmarshalJSON(*data); err != nil {
		return as2Key, err
	}

	return as2Key, nil
}

func Update(ctx context.Context, params files_sdk.As2KeyUpdateParams) (files_sdk.As2Key, error) {
	return (&Client{}).Update(ctx, params)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.As2KeyDeleteParams) (files_sdk.As2Key, error) {
	as2Key := files_sdk.As2Key{}
	if params.Id == 0 {
		return as2Key, lib.CreateError(params, "Id")
	}
	path := "/as2_keys/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return as2Key, err
	}
	data, res, err := files_sdk.Call(ctx, "DELETE", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return as2Key, err
	}
	if res.StatusCode == 204 {
		return as2Key, nil
	}
	if err := as2Key.UnmarshalJSON(*data); err != nil {
		return as2Key, err
	}

	return as2Key, nil
}

func Delete(ctx context.Context, params files_sdk.As2KeyDeleteParams) (files_sdk.As2Key, error) {
	return (&Client{}).Delete(ctx, params)
}
