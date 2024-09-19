package siem_http_destination

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

func (i *Iter) SiemHttpDestination() files_sdk.SiemHttpDestination {
	return i.Current().(files_sdk.SiemHttpDestination)
}

func (i *Iter) LoadResource(identifier interface{}, opts ...files_sdk.RequestResponseOption) (interface{}, error) {
	params := files_sdk.SiemHttpDestinationFindParams{}
	if id, ok := identifier.(int64); ok {
		params.Id = id
	}
	return i.Client.Find(params, opts...)
}

func (c *Client) List(params files_sdk.SiemHttpDestinationListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/siem_http_destinations", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.SiemHttpDestinationCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

func List(params files_sdk.SiemHttpDestinationListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(params, opts...)
}

func (c *Client) Find(params files_sdk.SiemHttpDestinationFindParams, opts ...files_sdk.RequestResponseOption) (siemHttpDestination files_sdk.SiemHttpDestination, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/siem_http_destinations/{id}", Params: params, Entity: &siemHttpDestination}, opts...)
	return
}

func Find(params files_sdk.SiemHttpDestinationFindParams, opts ...files_sdk.RequestResponseOption) (siemHttpDestination files_sdk.SiemHttpDestination, err error) {
	return (&Client{}).Find(params, opts...)
}

func (c *Client) Create(params files_sdk.SiemHttpDestinationCreateParams, opts ...files_sdk.RequestResponseOption) (siemHttpDestination files_sdk.SiemHttpDestination, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/siem_http_destinations", Params: params, Entity: &siemHttpDestination}, opts...)
	return
}

func Create(params files_sdk.SiemHttpDestinationCreateParams, opts ...files_sdk.RequestResponseOption) (siemHttpDestination files_sdk.SiemHttpDestination, err error) {
	return (&Client{}).Create(params, opts...)
}

func (c *Client) SendTestEntry(params files_sdk.SiemHttpDestinationSendTestEntryParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/siem_http_destinations/send_test_entry", Params: params, Entity: nil}, opts...)
	return
}

func SendTestEntry(params files_sdk.SiemHttpDestinationSendTestEntryParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).SendTestEntry(params, opts...)
}

func (c *Client) Update(params files_sdk.SiemHttpDestinationUpdateParams, opts ...files_sdk.RequestResponseOption) (siemHttpDestination files_sdk.SiemHttpDestination, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/siem_http_destinations/{id}", Params: params, Entity: &siemHttpDestination}, opts...)
	return
}

func Update(params files_sdk.SiemHttpDestinationUpdateParams, opts ...files_sdk.RequestResponseOption) (siemHttpDestination files_sdk.SiemHttpDestination, err error) {
	return (&Client{}).Update(params, opts...)
}

func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (siemHttpDestination files_sdk.SiemHttpDestination, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/siem_http_destinations/{id}", Params: params, Entity: &siemHttpDestination}, opts...)
	return
}

func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (siemHttpDestination files_sdk.SiemHttpDestination, err error) {
	return (&Client{}).UpdateWithMap(params, opts...)
}

func (c *Client) Delete(params files_sdk.SiemHttpDestinationDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "DELETE", Path: "/siem_http_destinations/{id}", Params: params, Entity: nil}, opts...)
	return
}

func Delete(params files_sdk.SiemHttpDestinationDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(params, opts...)
}
