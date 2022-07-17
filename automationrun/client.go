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
	*lib.Iter
}

func (i *Iter) AutomationRun() files_sdk.AutomationRun {
	return i.Current().(files_sdk.AutomationRun)
}

func (c *Client) List(ctx context.Context, params files_sdk.AutomationRunListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	path, err := lib.BuildPath("/automation_runs", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.AutomationRunCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.AutomationRunListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}

func (c *Client) Find(ctx context.Context, params files_sdk.AutomationRunFindParams) (automationRun files_sdk.AutomationRun, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/automation_runs/{id}", Params: params, Entity: &automationRun})
	return
}

func Find(ctx context.Context, params files_sdk.AutomationRunFindParams) (automationRun files_sdk.AutomationRun, err error) {
	return (&Client{}).Find(ctx, params)
}
