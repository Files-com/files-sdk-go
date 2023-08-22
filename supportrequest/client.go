package support_request

import (
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

func (i *Iter) SupportRequest() files_sdk.SupportRequest {
	return i.Current().(files_sdk.SupportRequest)
}

func (c *Client) List(params files_sdk.SupportRequestListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/support_requests", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.SupportRequestCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

func List(params files_sdk.SupportRequestListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(params, opts...)
}

func (c *Client) Create(params files_sdk.SupportRequestCreateParams, opts ...files_sdk.RequestResponseOption) (supportRequest files_sdk.SupportRequest, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/support_requests", Params: params, Entity: &supportRequest}, opts...)
	return
}

func Create(params files_sdk.SupportRequestCreateParams, opts ...files_sdk.RequestResponseOption) (supportRequest files_sdk.SupportRequest, err error) {
	return (&Client{}).Create(params, opts...)
}

func (c *Client) Update(params files_sdk.SupportRequestUpdateParams, opts ...files_sdk.RequestResponseOption) (supportRequest files_sdk.SupportRequest, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/support_requests/{id}", Params: params, Entity: &supportRequest}, opts...)
	return
}

func Update(params files_sdk.SupportRequestUpdateParams, opts ...files_sdk.RequestResponseOption) (supportRequest files_sdk.SupportRequest, err error) {
	return (&Client{}).Update(params, opts...)
}

func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (supportRequest files_sdk.SupportRequest, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/support_requests/{id}", Params: params, Entity: &supportRequest}, opts...)
	return
}

func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (supportRequest files_sdk.SupportRequest, err error) {
	return (&Client{}).UpdateWithMap(params, opts...)
}
