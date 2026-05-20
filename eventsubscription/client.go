package event_subscription

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

func (i *Iter) EventSubscription() files_sdk.EventSubscription {
	return i.Current().(files_sdk.EventSubscription)
}

func (i *Iter) LoadResource(identifier interface{}, opts ...files_sdk.RequestResponseOption) (interface{}, error) {
	params := files_sdk.EventSubscriptionFindParams{}
	if id, ok := identifier.(int64); ok {
		params.Id = id
	}
	return i.Client.Find(params, opts...)
}

func (c *Client) List(params files_sdk.EventSubscriptionListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/event_subscriptions", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.EventSubscriptionCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

func List(params files_sdk.EventSubscriptionListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(params, opts...)
}

func (c *Client) Find(params files_sdk.EventSubscriptionFindParams, opts ...files_sdk.RequestResponseOption) (eventSubscription files_sdk.EventSubscription, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/event_subscriptions/{id}", Params: params, Entity: &eventSubscription}, opts...)
	return
}

func Find(params files_sdk.EventSubscriptionFindParams, opts ...files_sdk.RequestResponseOption) (eventSubscription files_sdk.EventSubscription, err error) {
	return (&Client{}).Find(params, opts...)
}

func (c *Client) Create(params files_sdk.EventSubscriptionCreateParams, opts ...files_sdk.RequestResponseOption) (eventSubscription files_sdk.EventSubscription, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/event_subscriptions", Params: params, Entity: &eventSubscription}, opts...)
	return
}

func Create(params files_sdk.EventSubscriptionCreateParams, opts ...files_sdk.RequestResponseOption) (eventSubscription files_sdk.EventSubscription, err error) {
	return (&Client{}).Create(params, opts...)
}

func (c *Client) Update(params files_sdk.EventSubscriptionUpdateParams, opts ...files_sdk.RequestResponseOption) (eventSubscription files_sdk.EventSubscription, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/event_subscriptions/{id}", Params: params, Entity: &eventSubscription}, opts...)
	return
}

func Update(params files_sdk.EventSubscriptionUpdateParams, opts ...files_sdk.RequestResponseOption) (eventSubscription files_sdk.EventSubscription, err error) {
	return (&Client{}).Update(params, opts...)
}

func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (eventSubscription files_sdk.EventSubscription, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/event_subscriptions/{id}", Params: params, Entity: &eventSubscription}, opts...)
	return
}

func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (eventSubscription files_sdk.EventSubscription, err error) {
	return (&Client{}).UpdateWithMap(params, opts...)
}

func (c *Client) Delete(params files_sdk.EventSubscriptionDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "DELETE", Path: "/event_subscriptions/{id}", Params: params, Entity: nil}, opts...)
	return
}

func Delete(params files_sdk.EventSubscriptionDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(params, opts...)
}
