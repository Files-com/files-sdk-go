package clickwrap

import (
	"context"
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

func (i *Iter) Clickwrap() files_sdk.Clickwrap {
	return i.Current().(files_sdk.Clickwrap)
}

func (c *Client) List(ctx context.Context, params files_sdk.ClickwrapListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/clickwraps"
	i.ListParams = &params
	list := files_sdk.ClickwrapCollection{}
	i.Query = listquery.Build(ctx, i, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.ClickwrapListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}

func (c *Client) Find(ctx context.Context, params files_sdk.ClickwrapFindParams) (files_sdk.Clickwrap, error) {
	clickwrap := files_sdk.Clickwrap{}
	if params.Id == 0 {
		return clickwrap, lib.CreateError(params, "Id")
	}
	path := "/clickwraps/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return clickwrap, err
	}
	data, res, err := files_sdk.Call(ctx, "GET", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return clickwrap, err
	}
	if res.StatusCode == 204 {
		return clickwrap, nil
	}
	if err := clickwrap.UnmarshalJSON(*data); err != nil {
		return clickwrap, err
	}

	return clickwrap, nil
}

func Find(ctx context.Context, params files_sdk.ClickwrapFindParams) (files_sdk.Clickwrap, error) {
	return (&Client{}).Find(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.ClickwrapCreateParams) (files_sdk.Clickwrap, error) {
	clickwrap := files_sdk.Clickwrap{}
	path := "/clickwraps"
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return clickwrap, err
	}
	data, res, err := files_sdk.Call(ctx, "POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return clickwrap, err
	}
	if res.StatusCode == 204 {
		return clickwrap, nil
	}
	if err := clickwrap.UnmarshalJSON(*data); err != nil {
		return clickwrap, err
	}

	return clickwrap, nil
}

func Create(ctx context.Context, params files_sdk.ClickwrapCreateParams) (files_sdk.Clickwrap, error) {
	return (&Client{}).Create(ctx, params)
}

func (c *Client) Update(ctx context.Context, params files_sdk.ClickwrapUpdateParams) (files_sdk.Clickwrap, error) {
	clickwrap := files_sdk.Clickwrap{}
	if params.Id == 0 {
		return clickwrap, lib.CreateError(params, "Id")
	}
	path := "/clickwraps/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return clickwrap, err
	}
	data, res, err := files_sdk.Call(ctx, "PATCH", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return clickwrap, err
	}
	if res.StatusCode == 204 {
		return clickwrap, nil
	}
	if err := clickwrap.UnmarshalJSON(*data); err != nil {
		return clickwrap, err
	}

	return clickwrap, nil
}

func Update(ctx context.Context, params files_sdk.ClickwrapUpdateParams) (files_sdk.Clickwrap, error) {
	return (&Client{}).Update(ctx, params)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.ClickwrapDeleteParams) (files_sdk.Clickwrap, error) {
	clickwrap := files_sdk.Clickwrap{}
	if params.Id == 0 {
		return clickwrap, lib.CreateError(params, "Id")
	}
	path := "/clickwraps/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return clickwrap, err
	}
	data, res, err := files_sdk.Call(ctx, "DELETE", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return clickwrap, err
	}
	if res.StatusCode == 204 {
		return clickwrap, nil
	}
	if err := clickwrap.UnmarshalJSON(*data); err != nil {
		return clickwrap, err
	}

	return clickwrap, nil
}

func Delete(ctx context.Context, params files_sdk.ClickwrapDeleteParams) (files_sdk.Clickwrap, error) {
	return (&Client{}).Delete(ctx, params)
}
