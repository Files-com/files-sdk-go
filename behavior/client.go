package behavior

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

func (i *Iter) Behavior() files_sdk.Behavior {
	return i.Current().(files_sdk.Behavior)
}

func (c *Client) List(ctx context.Context, params files_sdk.BehaviorListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	path, err := lib.BuildPath("/behaviors", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.BehaviorCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list, opts...)
	return i, nil
}

func List(ctx context.Context, params files_sdk.BehaviorListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(ctx, params, opts...)
}

func (c *Client) Find(ctx context.Context, params files_sdk.BehaviorFindParams, opts ...files_sdk.RequestResponseOption) (behavior files_sdk.Behavior, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/behaviors/{id}", Params: params, Entity: &behavior}, opts...)
	return
}

func Find(ctx context.Context, params files_sdk.BehaviorFindParams, opts ...files_sdk.RequestResponseOption) (behavior files_sdk.Behavior, err error) {
	return (&Client{}).Find(ctx, params, opts...)
}

func (c *Client) ListFor(ctx context.Context, params files_sdk.BehaviorListForParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	path, err := lib.BuildPath("/behaviors/folders/{path}", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.BehaviorCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list, opts...)
	return i, nil
}

func ListFor(ctx context.Context, params files_sdk.BehaviorListForParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).ListFor(ctx, params, opts...)
}

func (c *Client) Create(ctx context.Context, params files_sdk.BehaviorCreateParams, opts ...files_sdk.RequestResponseOption) (behavior files_sdk.Behavior, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/behaviors", Params: params, Entity: &behavior}, opts...)
	return
}

func Create(ctx context.Context, params files_sdk.BehaviorCreateParams, opts ...files_sdk.RequestResponseOption) (behavior files_sdk.Behavior, err error) {
	return (&Client{}).Create(ctx, params, opts...)
}

func (c *Client) WebhookTest(ctx context.Context, params files_sdk.BehaviorWebhookTestParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/behaviors/webhook/test", Params: params, Entity: nil}, opts...)
	return
}

func WebhookTest(ctx context.Context, params files_sdk.BehaviorWebhookTestParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).WebhookTest(ctx, params, opts...)
}

func (c *Client) Update(ctx context.Context, params files_sdk.BehaviorUpdateParams, opts ...files_sdk.RequestResponseOption) (behavior files_sdk.Behavior, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "PATCH", Path: "/behaviors/{id}", Params: params, Entity: &behavior}, opts...)
	return
}

func Update(ctx context.Context, params files_sdk.BehaviorUpdateParams, opts ...files_sdk.RequestResponseOption) (behavior files_sdk.Behavior, err error) {
	return (&Client{}).Update(ctx, params, opts...)
}

func (c *Client) UpdateWithMap(ctx context.Context, params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (behavior files_sdk.Behavior, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "PATCH", Path: "/behaviors/{id}", Params: params, Entity: &behavior}, opts...)
	return
}

func UpdateWithMap(ctx context.Context, params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (behavior files_sdk.Behavior, err error) {
	return (&Client{}).UpdateWithMap(ctx, params, opts...)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.BehaviorDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "DELETE", Path: "/behaviors/{id}", Params: params, Entity: nil}, opts...)
	return
}

func Delete(ctx context.Context, params files_sdk.BehaviorDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(ctx, params, opts...)
}
