package bundle

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

func (i *Iter) Bundle() files_sdk.Bundle {
	return i.Current().(files_sdk.Bundle)
}

func (c *Client) List(ctx context.Context, params files_sdk.BundleListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	path, err := lib.BuildPath("/bundles", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.BundleCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list, opts...)
	return i, nil
}

func List(ctx context.Context, params files_sdk.BundleListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(ctx, params, opts...)
}

func (c *Client) Find(ctx context.Context, params files_sdk.BundleFindParams, opts ...files_sdk.RequestResponseOption) (bundle files_sdk.Bundle, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/bundles/{id}", Params: params, Entity: &bundle}, opts...)
	return
}

func Find(ctx context.Context, params files_sdk.BundleFindParams, opts ...files_sdk.RequestResponseOption) (bundle files_sdk.Bundle, err error) {
	return (&Client{}).Find(ctx, params, opts...)
}

func (c *Client) Create(ctx context.Context, params files_sdk.BundleCreateParams, opts ...files_sdk.RequestResponseOption) (bundle files_sdk.Bundle, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/bundles", Params: params, Entity: &bundle}, opts...)
	return
}

func Create(ctx context.Context, params files_sdk.BundleCreateParams, opts ...files_sdk.RequestResponseOption) (bundle files_sdk.Bundle, err error) {
	return (&Client{}).Create(ctx, params, opts...)
}

func (c *Client) Share(ctx context.Context, params files_sdk.BundleShareParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/bundles/{id}/share", Params: params, Entity: nil}, opts...)
	return
}

func Share(ctx context.Context, params files_sdk.BundleShareParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Share(ctx, params, opts...)
}

func (c *Client) Update(ctx context.Context, params files_sdk.BundleUpdateParams, opts ...files_sdk.RequestResponseOption) (bundle files_sdk.Bundle, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "PATCH", Path: "/bundles/{id}", Params: params, Entity: &bundle}, opts...)
	return
}

func Update(ctx context.Context, params files_sdk.BundleUpdateParams, opts ...files_sdk.RequestResponseOption) (bundle files_sdk.Bundle, err error) {
	return (&Client{}).Update(ctx, params, opts...)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.BundleDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "DELETE", Path: "/bundles/{id}", Params: params, Entity: nil}, opts...)
	return
}

func Delete(ctx context.Context, params files_sdk.BundleDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(ctx, params, opts...)
}
