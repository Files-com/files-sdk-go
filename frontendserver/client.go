package front_end_server

import (
	files_sdk "github.com/Files-com/files-sdk-go/v2"
	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Create(params files_sdk.FrontEndServerCreateParams, opts ...files_sdk.RequestResponseOption) (ip files_sdk.Ip, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/front_end_servers", Params: params, Entity: &ip}, opts...)
	return
}

func Create(params files_sdk.FrontEndServerCreateParams, opts ...files_sdk.RequestResponseOption) (ip files_sdk.Ip, err error) {
	return (&Client{}).Create(params, opts...)
}
