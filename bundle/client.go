package bundle

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

func (i *Iter) Bundle() files_sdk.Bundle {
	return i.Current().(files_sdk.Bundle)
}

func (c *Client) List(ctx context.Context, params files_sdk.BundleListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/bundles"
	i.ListParams = &params
	list := files_sdk.BundleCollection{}
	i.Query = listquery.Build(ctx, i, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.BundleListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}

func (c *Client) Find(ctx context.Context, params files_sdk.BundleFindParams) (files_sdk.Bundle, error) {
	bundle := files_sdk.Bundle{}
	if params.Id == 0 {
		return bundle, lib.CreateError(params, "Id")
	}
	path := "/bundles/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return bundle, err
	}
	data, res, err := files_sdk.Call(ctx, "GET", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return bundle, err
	}
	if res.StatusCode == 204 {
		return bundle, nil
	}
	if err := bundle.UnmarshalJSON(*data); err != nil {
		return bundle, err
	}

	return bundle, nil
}

func Find(ctx context.Context, params files_sdk.BundleFindParams) (files_sdk.Bundle, error) {
	return (&Client{}).Find(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.BundleCreateParams) (files_sdk.Bundle, error) {
	bundle := files_sdk.Bundle{}
	path := "/bundles"
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return bundle, err
	}
	data, res, err := files_sdk.Call(ctx, "POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return bundle, err
	}
	if res.StatusCode == 204 {
		return bundle, nil
	}
	if err := bundle.UnmarshalJSON(*data); err != nil {
		return bundle, err
	}

	return bundle, nil
}

func Create(ctx context.Context, params files_sdk.BundleCreateParams) (files_sdk.Bundle, error) {
	return (&Client{}).Create(ctx, params)
}

func (c *Client) Share(ctx context.Context, params files_sdk.BundleShareParams) (files_sdk.Bundle, error) {
	bundle := files_sdk.Bundle{}
	if params.Id == 0 {
		return bundle, lib.CreateError(params, "Id")
	}
	path := "/bundles/" + strconv.FormatInt(params.Id, 10) + "/share"
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return bundle, err
	}
	data, res, err := files_sdk.Call(ctx, "POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return bundle, err
	}
	if res.StatusCode == 204 {
		return bundle, nil
	}
	if err := bundle.UnmarshalJSON(*data); err != nil {
		return bundle, err
	}

	return bundle, nil
}

func Share(ctx context.Context, params files_sdk.BundleShareParams) (files_sdk.Bundle, error) {
	return (&Client{}).Share(ctx, params)
}

func (c *Client) Update(ctx context.Context, params files_sdk.BundleUpdateParams) (files_sdk.Bundle, error) {
	bundle := files_sdk.Bundle{}
	if params.Id == 0 {
		return bundle, lib.CreateError(params, "Id")
	}
	path := "/bundles/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return bundle, err
	}
	data, res, err := files_sdk.Call(ctx, "PATCH", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return bundle, err
	}
	if res.StatusCode == 204 {
		return bundle, nil
	}
	if err := bundle.UnmarshalJSON(*data); err != nil {
		return bundle, err
	}

	return bundle, nil
}

func Update(ctx context.Context, params files_sdk.BundleUpdateParams) (files_sdk.Bundle, error) {
	return (&Client{}).Update(ctx, params)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.BundleDeleteParams) (files_sdk.Bundle, error) {
	bundle := files_sdk.Bundle{}
	if params.Id == 0 {
		return bundle, lib.CreateError(params, "Id")
	}
	path := "/bundles/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return bundle, err
	}
	data, res, err := files_sdk.Call(ctx, "DELETE", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return bundle, err
	}
	if res.StatusCode == 204 {
		return bundle, nil
	}
	if err := bundle.UnmarshalJSON(*data); err != nil {
		return bundle, err
	}

	return bundle, nil
}

func Delete(ctx context.Context, params files_sdk.BundleDeleteParams) (files_sdk.Bundle, error) {
	return (&Client{}).Delete(ctx, params)
}
