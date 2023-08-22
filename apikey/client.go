package api_key

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

func (i *Iter) ApiKey() files_sdk.ApiKey {
	return i.Current().(files_sdk.ApiKey)
}

func (i *Iter) LoadResource(identifier interface{}, opts ...files_sdk.RequestResponseOption) (interface{}, error) {
	params := files_sdk.ApiKeyFindParams{}
	if id, ok := identifier.(int64); ok {
		params.Id = id
	}
	return i.Client.Find(params, opts...)
}

func (c *Client) List(params files_sdk.ApiKeyListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/api_keys", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.ApiKeyCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

func List(params files_sdk.ApiKeyListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(params, opts...)
}

func (c *Client) Find(params files_sdk.ApiKeyFindParams, opts ...files_sdk.RequestResponseOption) (apiKey files_sdk.ApiKey, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/api_keys/{id}", Params: params, Entity: &apiKey}, opts...)
	return
}

func Find(params files_sdk.ApiKeyFindParams, opts ...files_sdk.RequestResponseOption) (apiKey files_sdk.ApiKey, err error) {
	return (&Client{}).Find(params, opts...)
}

func (c *Client) FindCurrent(opts ...files_sdk.RequestResponseOption) (apiKey files_sdk.ApiKey, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/api_key", Params: lib.Interface(), Entity: &apiKey}, opts...)
	return
}

func FindCurrent(opts ...files_sdk.RequestResponseOption) (apiKey files_sdk.ApiKey, err error) {
	return (&Client{}).FindCurrent(opts...)
}

func (c *Client) Create(params files_sdk.ApiKeyCreateParams, opts ...files_sdk.RequestResponseOption) (apiKey files_sdk.ApiKey, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/api_keys", Params: params, Entity: &apiKey}, opts...)
	return
}

func Create(params files_sdk.ApiKeyCreateParams, opts ...files_sdk.RequestResponseOption) (apiKey files_sdk.ApiKey, err error) {
	return (&Client{}).Create(params, opts...)
}

func (c *Client) Update(params files_sdk.ApiKeyUpdateParams, opts ...files_sdk.RequestResponseOption) (apiKey files_sdk.ApiKey, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/api_keys/{id}", Params: params, Entity: &apiKey}, opts...)
	return
}

func Update(params files_sdk.ApiKeyUpdateParams, opts ...files_sdk.RequestResponseOption) (apiKey files_sdk.ApiKey, err error) {
	return (&Client{}).Update(params, opts...)
}

func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (apiKey files_sdk.ApiKey, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/api_keys/{id}", Params: params, Entity: &apiKey}, opts...)
	return
}

func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (apiKey files_sdk.ApiKey, err error) {
	return (&Client{}).UpdateWithMap(params, opts...)
}

func (c *Client) UpdateCurrent(params files_sdk.ApiKeyUpdateCurrentParams, opts ...files_sdk.RequestResponseOption) (apiKey files_sdk.ApiKey, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/api_key", Params: params, Entity: &apiKey}, opts...)
	return
}

func UpdateCurrent(params files_sdk.ApiKeyUpdateCurrentParams, opts ...files_sdk.RequestResponseOption) (apiKey files_sdk.ApiKey, err error) {
	return (&Client{}).UpdateCurrent(params, opts...)
}

func (c *Client) UpdateCurrentWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (apiKey files_sdk.ApiKey, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/api_key", Params: params, Entity: &apiKey}, opts...)
	return
}

func UpdateCurrentWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (apiKey files_sdk.ApiKey, err error) {
	return (&Client{}).UpdateCurrentWithMap(params, opts...)
}

func (c *Client) Delete(params files_sdk.ApiKeyDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "DELETE", Path: "/api_keys/{id}", Params: params, Entity: nil}, opts...)
	return
}

func Delete(params files_sdk.ApiKeyDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(params, opts...)
}

func (c *Client) DeleteCurrent(opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "DELETE", Path: "/api_key", Params: lib.Interface(), Entity: nil}, opts...)
	return
}

func DeleteCurrent(opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).DeleteCurrent(opts...)
}
