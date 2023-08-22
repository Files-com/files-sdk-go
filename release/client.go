package release

import (
	files_sdk "github.com/Files-com/files-sdk-go/v2"
	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) GetLatest(params files_sdk.ReleaseGetLatestParams, opts ...files_sdk.RequestResponseOption) (release files_sdk.Release, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/releases/latest", Params: params, Entity: &release}, opts...)
	return
}

func GetLatest(params files_sdk.ReleaseGetLatestParams, opts ...files_sdk.RequestResponseOption) (release files_sdk.Release, err error) {
	return (&Client{}).GetLatest(params, opts...)
}
