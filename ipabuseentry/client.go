package ip_abuse_entry

import (
	files_sdk "github.com/Files-com/files-sdk-go/v2"
	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Create(params files_sdk.IpAbuseEntryCreateParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/ip_abuse_entries", Params: params, Entity: nil}, opts...)
	return
}

func Create(params files_sdk.IpAbuseEntryCreateParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Create(params, opts...)
}
