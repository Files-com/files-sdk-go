package click

import (
	files_sdk "github.com/Files-com/files-sdk-go/v2"
	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Create(params files_sdk.ClickCreateParams, opts ...files_sdk.RequestResponseOption) (click files_sdk.Click, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/clicks", Params: params, Entity: &click}, opts...)
	return
}

func Create(params files_sdk.ClickCreateParams, opts ...files_sdk.RequestResponseOption) (click files_sdk.Click, err error) {
	return (&Client{}).Create(params, opts...)
}
