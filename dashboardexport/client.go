package dashboard_export

import (
	files_sdk "github.com/Files-com/files-sdk-go/v2"
	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Find(params files_sdk.DashboardExportFindParams, opts ...files_sdk.RequestResponseOption) (dashboardExport files_sdk.DashboardExport, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/dashboard_exports/{id}", Params: params, Entity: &dashboardExport}, opts...)
	return
}

func Find(params files_sdk.DashboardExportFindParams, opts ...files_sdk.RequestResponseOption) (dashboardExport files_sdk.DashboardExport, err error) {
	return (&Client{}).Find(params, opts...)
}

func (c *Client) Create(params files_sdk.DashboardExportCreateParams, opts ...files_sdk.RequestResponseOption) (dashboardExport files_sdk.DashboardExport, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/dashboard_exports", Params: params, Entity: &dashboardExport}, opts...)
	return
}

func Create(params files_sdk.DashboardExportCreateParams, opts ...files_sdk.RequestResponseOption) (dashboardExport files_sdk.DashboardExport, err error) {
	return (&Client{}).Create(params, opts...)
}
