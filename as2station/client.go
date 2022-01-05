package as2_station

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

func (i *Iter) As2Station() files_sdk.As2Station {
	return i.Current().(files_sdk.As2Station)
}

func (c *Client) List(ctx context.Context, params files_sdk.As2StationListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/as2_stations"
	i.ListParams = &params
	list := files_sdk.As2StationCollection{}
	i.Query = listquery.Build(ctx, i, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.As2StationListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}

func (c *Client) Find(ctx context.Context, params files_sdk.As2StationFindParams) (files_sdk.As2Station, error) {
	as2Station := files_sdk.As2Station{}
	if params.Id == 0 {
		return as2Station, lib.CreateError(params, "Id")
	}
	path := "/as2_stations/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return as2Station, err
	}
	data, res, err := files_sdk.Call(ctx, "GET", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return as2Station, err
	}
	if res.StatusCode == 204 {
		return as2Station, nil
	}
	if err := as2Station.UnmarshalJSON(*data); err != nil {
		return as2Station, err
	}

	return as2Station, nil
}

func Find(ctx context.Context, params files_sdk.As2StationFindParams) (files_sdk.As2Station, error) {
	return (&Client{}).Find(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.As2StationCreateParams) (files_sdk.As2Station, error) {
	as2Station := files_sdk.As2Station{}
	path := "/as2_stations"
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return as2Station, err
	}
	data, res, err := files_sdk.Call(ctx, "POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return as2Station, err
	}
	if res.StatusCode == 204 {
		return as2Station, nil
	}
	if err := as2Station.UnmarshalJSON(*data); err != nil {
		return as2Station, err
	}

	return as2Station, nil
}

func Create(ctx context.Context, params files_sdk.As2StationCreateParams) (files_sdk.As2Station, error) {
	return (&Client{}).Create(ctx, params)
}

func (c *Client) Update(ctx context.Context, params files_sdk.As2StationUpdateParams) (files_sdk.As2Station, error) {
	as2Station := files_sdk.As2Station{}
	if params.Id == 0 {
		return as2Station, lib.CreateError(params, "Id")
	}
	path := "/as2_stations/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return as2Station, err
	}
	data, res, err := files_sdk.Call(ctx, "PATCH", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return as2Station, err
	}
	if res.StatusCode == 204 {
		return as2Station, nil
	}
	if err := as2Station.UnmarshalJSON(*data); err != nil {
		return as2Station, err
	}

	return as2Station, nil
}

func Update(ctx context.Context, params files_sdk.As2StationUpdateParams) (files_sdk.As2Station, error) {
	return (&Client{}).Update(ctx, params)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.As2StationDeleteParams) (files_sdk.As2Station, error) {
	as2Station := files_sdk.As2Station{}
	if params.Id == 0 {
		return as2Station, lib.CreateError(params, "Id")
	}
	path := "/as2_stations/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return as2Station, err
	}
	data, res, err := files_sdk.Call(ctx, "DELETE", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return as2Station, err
	}
	if res.StatusCode == 204 {
		return as2Station, nil
	}
	if err := as2Station.UnmarshalJSON(*data); err != nil {
		return as2Station, err
	}

	return as2Station, nil
}

func Delete(ctx context.Context, params files_sdk.As2StationDeleteParams) (files_sdk.As2Station, error) {
	return (&Client{}).Delete(ctx, params)
}
