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

func (c *Client) List(ctx context.Context, params files_sdk.BehaviorListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	path, err := lib.BuildPath("/behaviors", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.BehaviorCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.BehaviorListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}

func (c *Client) Find(ctx context.Context, params files_sdk.BehaviorFindParams) (behavior files_sdk.Behavior, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/behaviors/{id}", Params: params, Entity: &behavior})
	return
}

func Find(ctx context.Context, params files_sdk.BehaviorFindParams) (behavior files_sdk.Behavior, err error) {
	return (&Client{}).Find(ctx, params)
}

func (c *Client) ListFor(ctx context.Context, params files_sdk.BehaviorListForParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	path, err := lib.BuildPath("/behaviors/folders/{path}", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.BehaviorCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func ListFor(ctx context.Context, params files_sdk.BehaviorListForParams) (*Iter, error) {
	return (&Client{}).ListFor(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.BehaviorCreateParams) (behavior files_sdk.Behavior, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/behaviors", Params: params, Entity: &behavior})
	return
}

func Create(ctx context.Context, params files_sdk.BehaviorCreateParams) (behavior files_sdk.Behavior, err error) {
	return (&Client{}).Create(ctx, params)
}

func (c *Client) WebhookTest(ctx context.Context, params files_sdk.BehaviorWebhookTestParams) (err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/behaviors/webhook/test", Params: params, Entity: nil})
	return
}

func WebhookTest(ctx context.Context, params files_sdk.BehaviorWebhookTestParams) (err error) {
	return (&Client{}).WebhookTest(ctx, params)
}

func (c *Client) Update(ctx context.Context, params files_sdk.BehaviorUpdateParams) (behavior files_sdk.Behavior, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "PATCH", Path: "/behaviors/{id}", Params: params, Entity: &behavior})
	return
}

func Update(ctx context.Context, params files_sdk.BehaviorUpdateParams) (behavior files_sdk.Behavior, err error) {
	return (&Client{}).Update(ctx, params)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.BehaviorDeleteParams) (err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "DELETE", Path: "/behaviors/{id}", Params: params, Entity: nil})
	return
}

func Delete(ctx context.Context, params files_sdk.BehaviorDeleteParams) (err error) {
	return (&Client{}).Delete(ctx, params)
}
