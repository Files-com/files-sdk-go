package event_target

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

func (i *Iter) EventTarget() files_sdk.EventTarget {
	return i.Current().(files_sdk.EventTarget)
}

func (i *Iter) LoadResource(identifier interface{}, opts ...files_sdk.RequestResponseOption) (interface{}, error) {
	params := files_sdk.EventTargetFindParams{}
	if id, ok := identifier.(int64); ok {
		params.Id = id
	}
	return i.Client.Find(params, opts...)
}

func (c *Client) List(params files_sdk.EventTargetListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/event_targets", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.EventTargetCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

func List(params files_sdk.EventTargetListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(params, opts...)
}

func (c *Client) Find(params files_sdk.EventTargetFindParams, opts ...files_sdk.RequestResponseOption) (eventTarget files_sdk.EventTarget, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/event_targets/{id}", Params: params, Entity: &eventTarget}, opts...)
	return
}

func Find(params files_sdk.EventTargetFindParams, opts ...files_sdk.RequestResponseOption) (eventTarget files_sdk.EventTarget, err error) {
	return (&Client{}).Find(params, opts...)
}

func (c *Client) Create(params files_sdk.EventTargetCreateParams, opts ...files_sdk.RequestResponseOption) (eventTarget files_sdk.EventTarget, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/event_targets", Params: params, Entity: &eventTarget}, opts...)
	return
}

func Create(params files_sdk.EventTargetCreateParams, opts ...files_sdk.RequestResponseOption) (eventTarget files_sdk.EventTarget, err error) {
	return (&Client{}).Create(params, opts...)
}

func (c *Client) Update(params files_sdk.EventTargetUpdateParams, opts ...files_sdk.RequestResponseOption) (eventTarget files_sdk.EventTarget, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/event_targets/{id}", Params: params, Entity: &eventTarget}, opts...)
	return
}

func Update(params files_sdk.EventTargetUpdateParams, opts ...files_sdk.RequestResponseOption) (eventTarget files_sdk.EventTarget, err error) {
	return (&Client{}).Update(params, opts...)
}

func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (eventTarget files_sdk.EventTarget, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/event_targets/{id}", Params: params, Entity: &eventTarget}, opts...)
	return
}

func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (eventTarget files_sdk.EventTarget, err error) {
	return (&Client{}).UpdateWithMap(params, opts...)
}

func (c *Client) Delete(params files_sdk.EventTargetDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "DELETE", Path: "/event_targets/{id}", Params: params, Entity: nil}, opts...)
	return
}

func Delete(params files_sdk.EventTargetDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(params, opts...)
}
