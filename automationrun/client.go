package automation_run

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
	*files_sdk.Iter
	*Client
}

func (i *Iter) Reload(opts ...files_sdk.RequestResponseOption) files_sdk.IterI {
	return &Iter{Iter: i.Iter.Reload(opts...).(*files_sdk.Iter), Client: i.Client}
}

func (i *Iter) AutomationRun() files_sdk.AutomationRun {
	return i.Current().(files_sdk.AutomationRun)
}

func (i *Iter) LoadResource(identifier interface{}, opts ...files_sdk.RequestResponseOption) (interface{}, error) {
	params := files_sdk.AutomationRunFindParams{}
	if id, ok := identifier.(int64); ok {
		params.Id = id
	}
	return i.Client.Find(context.Background(), params, opts...)
}

func (c *Client) List(ctx context.Context, params files_sdk.AutomationRunListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/automation_runs", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.AutomationRunCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list, opts...)
	return i, nil
}

func List(ctx context.Context, params files_sdk.AutomationRunListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(ctx, params, opts...)
}

func (c *Client) Find(ctx context.Context, params files_sdk.AutomationRunFindParams, opts ...files_sdk.RequestResponseOption) (automationRun files_sdk.AutomationRun, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/automation_runs/{id}", Params: params, Entity: &automationRun}, opts...)
	return
}

func Find(ctx context.Context, params files_sdk.AutomationRunFindParams, opts ...files_sdk.RequestResponseOption) (automationRun files_sdk.AutomationRun, err error) {
	return (&Client{}).Find(ctx, params, opts...)
}
