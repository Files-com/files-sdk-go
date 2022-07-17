package automation

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

func (i *Iter) Automation() files_sdk.Automation {
	return i.Current().(files_sdk.Automation)
}

func (c *Client) List(ctx context.Context, params files_sdk.AutomationListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	path, err := lib.BuildPath("/automations", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.AutomationCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.AutomationListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}

func (c *Client) Find(ctx context.Context, params files_sdk.AutomationFindParams) (automation files_sdk.Automation, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/automations/{id}", Params: params, Entity: &automation})
	return
}

func Find(ctx context.Context, params files_sdk.AutomationFindParams) (automation files_sdk.Automation, err error) {
	return (&Client{}).Find(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.AutomationCreateParams) (automation files_sdk.Automation, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/automations", Params: params, Entity: &automation})
	return
}

func Create(ctx context.Context, params files_sdk.AutomationCreateParams) (automation files_sdk.Automation, err error) {
	return (&Client{}).Create(ctx, params)
}

func (c *Client) Update(ctx context.Context, params files_sdk.AutomationUpdateParams) (automation files_sdk.Automation, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "PATCH", Path: "/automations/{id}", Params: params, Entity: &automation})
	return
}

func Update(ctx context.Context, params files_sdk.AutomationUpdateParams) (automation files_sdk.Automation, err error) {
	return (&Client{}).Update(ctx, params)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.AutomationDeleteParams) (err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "DELETE", Path: "/automations/{id}", Params: params, Entity: nil})
	return
}

func Delete(ctx context.Context, params files_sdk.AutomationDeleteParams) (err error) {
	return (&Client{}).Delete(ctx, params)
}
