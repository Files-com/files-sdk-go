package lead

import (
	files_sdk "github.com/Files-com/files-sdk-go/v2"
	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Create(params files_sdk.LeadCreateParams, opts ...files_sdk.RequestResponseOption) (lead files_sdk.Lead, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/leads", Params: params, Entity: &lead}, opts...)
	return
}

func Create(params files_sdk.LeadCreateParams, opts ...files_sdk.RequestResponseOption) (lead files_sdk.Lead, err error) {
	return (&Client{}).Create(params, opts...)
}

func (c *Client) Update(params files_sdk.LeadUpdateParams, opts ...files_sdk.RequestResponseOption) (lead files_sdk.Lead, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/leads/{code}", Params: params, Entity: &lead}, opts...)
	return
}

func Update(params files_sdk.LeadUpdateParams, opts ...files_sdk.RequestResponseOption) (lead files_sdk.Lead, err error) {
	return (&Client{}).Update(params, opts...)
}

func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (lead files_sdk.Lead, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/leads/{code}", Params: params, Entity: &lead}, opts...)
	return
}

func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (lead files_sdk.Lead, err error) {
	return (&Client{}).UpdateWithMap(params, opts...)
}
