package as2_station

import (
	"context"

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
	path, err := lib.BuildPath("/as2_stations", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.As2StationCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.As2StationListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}

func (c *Client) Find(ctx context.Context, params files_sdk.As2StationFindParams) (as2Station files_sdk.As2Station, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/as2_stations/{id}", Params: params, Entity: &as2Station})
	return
}

func Find(ctx context.Context, params files_sdk.As2StationFindParams) (as2Station files_sdk.As2Station, err error) {
	return (&Client{}).Find(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.As2StationCreateParams) (as2Station files_sdk.As2Station, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/as2_stations", Params: params, Entity: &as2Station})
	return
}

func Create(ctx context.Context, params files_sdk.As2StationCreateParams) (as2Station files_sdk.As2Station, err error) {
	return (&Client{}).Create(ctx, params)
}

func (c *Client) Update(ctx context.Context, params files_sdk.As2StationUpdateParams) (as2Station files_sdk.As2Station, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "PATCH", Path: "/as2_stations/{id}", Params: params, Entity: &as2Station})
	return
}

func Update(ctx context.Context, params files_sdk.As2StationUpdateParams) (as2Station files_sdk.As2Station, err error) {
	return (&Client{}).Update(ctx, params)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.As2StationDeleteParams) (err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "DELETE", Path: "/as2_stations/{id}", Params: params, Entity: nil})
	return
}

func Delete(ctx context.Context, params files_sdk.As2StationDeleteParams) (err error) {
	return (&Client{}).Delete(ctx, params)
}
