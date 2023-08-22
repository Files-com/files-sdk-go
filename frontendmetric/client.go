package frontend_metric

import (
	files_sdk "github.com/Files-com/files-sdk-go/v2"
	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Create(params files_sdk.FrontendMetricCreateParams, opts ...files_sdk.RequestResponseOption) (frontendMetric files_sdk.FrontendMetric, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/frontend_metrics", Params: params, Entity: &frontendMetric}, opts...)
	return
}

func Create(params files_sdk.FrontendMetricCreateParams, opts ...files_sdk.RequestResponseOption) (frontendMetric files_sdk.FrontendMetric, err error) {
	return (&Client{}).Create(params, opts...)
}
