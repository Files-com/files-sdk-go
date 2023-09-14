package external_event

import (
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	lib "github.com/Files-com/files-sdk-go/v3/lib"
	listquery "github.com/Files-com/files-sdk-go/v3/listquery"
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

func (i *Iter) ExternalEvent() files_sdk.ExternalEvent {
	return i.Current().(files_sdk.ExternalEvent)
}

func (i *Iter) LoadResource(identifier interface{}, opts ...files_sdk.RequestResponseOption) (interface{}, error) {
	params := files_sdk.ExternalEventFindParams{}
	if id, ok := identifier.(int64); ok {
		params.Id = id
	}
	return i.Client.Find(params, opts...)
}

func (c *Client) List(params files_sdk.ExternalEventListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/external_events", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.ExternalEventCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

func List(params files_sdk.ExternalEventListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(params, opts...)
}

func (c *Client) Find(params files_sdk.ExternalEventFindParams, opts ...files_sdk.RequestResponseOption) (externalEvent files_sdk.ExternalEvent, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/external_events/{id}", Params: params, Entity: &externalEvent}, opts...)
	return
}

func Find(params files_sdk.ExternalEventFindParams, opts ...files_sdk.RequestResponseOption) (externalEvent files_sdk.ExternalEvent, err error) {
	return (&Client{}).Find(params, opts...)
}

func (c *Client) Create(params files_sdk.ExternalEventCreateParams, opts ...files_sdk.RequestResponseOption) (externalEvent files_sdk.ExternalEvent, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/external_events", Params: params, Entity: &externalEvent}, opts...)
	return
}

func Create(params files_sdk.ExternalEventCreateParams, opts ...files_sdk.RequestResponseOption) (externalEvent files_sdk.ExternalEvent, err error) {
	return (&Client{}).Create(params, opts...)
}
