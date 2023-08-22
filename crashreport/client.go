package crash_report

import (
	files_sdk "github.com/Files-com/files-sdk-go/v2"
	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Create(params files_sdk.CrashReportCreateParams, opts ...files_sdk.RequestResponseOption) (crashReport files_sdk.CrashReport, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/crash_reports", Params: params, Entity: &crashReport}, opts...)
	return
}

func Create(params files_sdk.CrashReportCreateParams, opts ...files_sdk.RequestResponseOption) (crashReport files_sdk.CrashReport, err error) {
	return (&Client{}).Create(params, opts...)
}
