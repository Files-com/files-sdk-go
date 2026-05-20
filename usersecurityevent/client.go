package user_security_event

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

func (i *Iter) UserSecurityEvent() files_sdk.UserSecurityEvent {
	return i.Current().(files_sdk.UserSecurityEvent)
}

func (i *Iter) LoadResource(identifier interface{}, opts ...files_sdk.RequestResponseOption) (interface{}, error) {
	params := files_sdk.UserSecurityEventFindParams{}
	if id, ok := identifier.(int64); ok {
		params.Id = id
	}
	return i.Client.Find(params, opts...)
}

func (c *Client) List(params files_sdk.UserSecurityEventListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/user_security_events", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.UserSecurityEventCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

func List(params files_sdk.UserSecurityEventListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(params, opts...)
}

func (c *Client) Find(params files_sdk.UserSecurityEventFindParams, opts ...files_sdk.RequestResponseOption) (userSecurityEvent files_sdk.UserSecurityEvent, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/user_security_events/{id}", Params: params, Entity: &userSecurityEvent}, opts...)
	return
}

func Find(params files_sdk.UserSecurityEventFindParams, opts ...files_sdk.RequestResponseOption) (userSecurityEvent files_sdk.UserSecurityEvent, err error) {
	return (&Client{}).Find(params, opts...)
}
