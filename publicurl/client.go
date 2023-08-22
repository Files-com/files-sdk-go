package public_url

import (
	files_sdk "github.com/Files-com/files-sdk-go/v2"
	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Create(params files_sdk.PublicUrlCreateParams, opts ...files_sdk.RequestResponseOption) (publicUrl files_sdk.PublicUrl, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/public_urls", Params: params, Entity: &publicUrl}, opts...)
	return
}

func Create(params files_sdk.PublicUrlCreateParams, opts ...files_sdk.RequestResponseOption) (publicUrl files_sdk.PublicUrl, err error) {
	return (&Client{}).Create(params, opts...)
}
