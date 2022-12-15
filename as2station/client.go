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

func (c *Client) List(ctx context.Context, params files_sdk.As2StationListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	path, err := lib.BuildPath("/as2_stations", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.As2StationCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list, opts...)
	return i, nil
}

func List(ctx context.Context, params files_sdk.As2StationListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(ctx, params, opts...)
}

func (c *Client) Find(ctx context.Context, params files_sdk.As2StationFindParams, opts ...files_sdk.RequestResponseOption) (as2Station files_sdk.As2Station, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/as2_stations/{id}", Params: params, Entity: &as2Station}, opts...)
	return
}

func Find(ctx context.Context, params files_sdk.As2StationFindParams, opts ...files_sdk.RequestResponseOption) (as2Station files_sdk.As2Station, err error) {
	return (&Client{}).Find(ctx, params, opts...)
}

func (c *Client) Create(ctx context.Context, params files_sdk.As2StationCreateParams, opts ...files_sdk.RequestResponseOption) (as2Station files_sdk.As2Station, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/as2_stations", Params: params, Entity: &as2Station}, opts...)
	return
}

func Create(ctx context.Context, params files_sdk.As2StationCreateParams, opts ...files_sdk.RequestResponseOption) (as2Station files_sdk.As2Station, err error) {
	return (&Client{}).Create(ctx, params, opts...)
}

func (c *Client) Update(ctx context.Context, params files_sdk.As2StationUpdateParams, opts ...files_sdk.RequestResponseOption) (as2Station files_sdk.As2Station, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "PATCH", Path: "/as2_stations/{id}", Params: params, Entity: &as2Station}, opts...)
	return
}

func Update(ctx context.Context, params files_sdk.As2StationUpdateParams, opts ...files_sdk.RequestResponseOption) (as2Station files_sdk.As2Station, err error) {
	return (&Client{}).Update(ctx, params, opts...)
}

func (c *Client) UpdateWithMap(ctx context.Context, params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (as2Station files_sdk.As2Station, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "PATCH", Path: "/as2_stations/{id}", Params: params, Entity: &as2Station}, opts...)
	return
}

func UpdateWithMap(ctx context.Context, params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (as2Station files_sdk.As2Station, err error) {
	return (&Client{}).UpdateWithMap(ctx, params, opts...)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.As2StationDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "DELETE", Path: "/as2_stations/{id}", Params: params, Entity: nil}, opts...)
	return
}

func Delete(ctx context.Context, params files_sdk.As2StationDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(ctx, params, opts...)
}
