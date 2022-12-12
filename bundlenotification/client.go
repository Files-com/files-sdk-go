package bundle_notification

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

func (i *Iter) BundleNotification() files_sdk.BundleNotification {
	return i.Current().(files_sdk.BundleNotification)
}

func (c *Client) List(ctx context.Context, params files_sdk.BundleNotificationListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	path, err := lib.BuildPath("/bundle_notifications", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.BundleNotificationCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list, opts...)
	return i, nil
}

func List(ctx context.Context, params files_sdk.BundleNotificationListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(ctx, params, opts...)
}

func (c *Client) Find(ctx context.Context, params files_sdk.BundleNotificationFindParams, opts ...files_sdk.RequestResponseOption) (bundleNotification files_sdk.BundleNotification, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/bundle_notifications/{id}", Params: params, Entity: &bundleNotification}, opts...)
	return
}

func Find(ctx context.Context, params files_sdk.BundleNotificationFindParams, opts ...files_sdk.RequestResponseOption) (bundleNotification files_sdk.BundleNotification, err error) {
	return (&Client{}).Find(ctx, params, opts...)
}

func (c *Client) Create(ctx context.Context, params files_sdk.BundleNotificationCreateParams, opts ...files_sdk.RequestResponseOption) (bundleNotification files_sdk.BundleNotification, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/bundle_notifications", Params: params, Entity: &bundleNotification}, opts...)
	return
}

func Create(ctx context.Context, params files_sdk.BundleNotificationCreateParams, opts ...files_sdk.RequestResponseOption) (bundleNotification files_sdk.BundleNotification, err error) {
	return (&Client{}).Create(ctx, params, opts...)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.BundleNotificationDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "DELETE", Path: "/bundle_notifications/{id}", Params: params, Entity: nil}, opts...)
	return
}

func Delete(ctx context.Context, params files_sdk.BundleNotificationDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(ctx, params, opts...)
}
