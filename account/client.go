package account

import (
	files_sdk "github.com/Files-com/files-sdk-go/v2"
	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Get(opts ...files_sdk.RequestResponseOption) (account files_sdk.Account, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/account", Params: lib.Interface(), Entity: &account}, opts...)
	return
}

func Get(opts ...files_sdk.RequestResponseOption) (account files_sdk.Account, err error) {
	return (&Client{}).Get(opts...)
}

func (c *Client) Create(params files_sdk.AccountCreateParams, opts ...files_sdk.RequestResponseOption) (account files_sdk.Account, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/account", Params: params, Entity: &account}, opts...)
	return
}

func Create(params files_sdk.AccountCreateParams, opts ...files_sdk.RequestResponseOption) (account files_sdk.Account, err error) {
	return (&Client{}).Create(params, opts...)
}

func (c *Client) Update(params files_sdk.AccountUpdateParams, opts ...files_sdk.RequestResponseOption) (account files_sdk.Account, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/account", Params: params, Entity: &account}, opts...)
	return
}

func Update(params files_sdk.AccountUpdateParams, opts ...files_sdk.RequestResponseOption) (account files_sdk.Account, err error) {
	return (&Client{}).Update(params, opts...)
}

func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (account files_sdk.Account, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/account", Params: params, Entity: &account}, opts...)
	return
}

func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (account files_sdk.Account, err error) {
	return (&Client{}).UpdateWithMap(params, opts...)
}
