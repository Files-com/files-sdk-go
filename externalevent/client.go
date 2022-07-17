package external_event

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

func (i *Iter) ExternalEvent() files_sdk.ExternalEvent {
	return i.Current().(files_sdk.ExternalEvent)
}

func (c *Client) List(ctx context.Context, params files_sdk.ExternalEventListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	path, err := lib.BuildPath("/external_events", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.ExternalEventCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.ExternalEventListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}

func (c *Client) Find(ctx context.Context, params files_sdk.ExternalEventFindParams) (externalEvent files_sdk.ExternalEvent, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/external_events/{id}", Params: params, Entity: &externalEvent})
	return
}

func Find(ctx context.Context, params files_sdk.ExternalEventFindParams) (externalEvent files_sdk.ExternalEvent, err error) {
	return (&Client{}).Find(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.ExternalEventCreateParams) (externalEvent files_sdk.ExternalEvent, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/external_events", Params: params, Entity: &externalEvent})
	return
}

func Create(ctx context.Context, params files_sdk.ExternalEventCreateParams) (externalEvent files_sdk.ExternalEvent, err error) {
	return (&Client{}).Create(ctx, params)
}
