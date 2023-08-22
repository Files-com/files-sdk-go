package site

import (
	files_sdk "github.com/Files-com/files-sdk-go/v2"
	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Get(opts ...files_sdk.RequestResponseOption) (site files_sdk.Site, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/site", Params: lib.Interface(), Entity: &site}, opts...)
	return
}

func Get(opts ...files_sdk.RequestResponseOption) (site files_sdk.Site, err error) {
	return (&Client{}).Get(opts...)
}

func (c *Client) GetSwitchToPlan(opts ...files_sdk.RequestResponseOption) (plan files_sdk.Plan, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/site/switch_to_plan", Params: lib.Interface(), Entity: &plan}, opts...)
	return
}

func GetSwitchToPlan(opts ...files_sdk.RequestResponseOption) (plan files_sdk.Plan, err error) {
	return (&Client{}).GetSwitchToPlan(opts...)
}

func (c *Client) GetPlan(params files_sdk.SiteGetPlanParams, opts ...files_sdk.RequestResponseOption) (plan files_sdk.Plan, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/site/plan", Params: params, Entity: &plan}, opts...)
	return
}

func GetPlan(params files_sdk.SiteGetPlanParams, opts ...files_sdk.RequestResponseOption) (plan files_sdk.Plan, err error) {
	return (&Client{}).GetPlan(params, opts...)
}

func (c *Client) GetPaypalExpressInfo(params files_sdk.SiteGetPaypalExpressInfoParams, opts ...files_sdk.RequestResponseOption) (paypalExpressInfo files_sdk.PaypalExpressInfo, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/site/paypal_express_info", Params: params, Entity: &paypalExpressInfo}, opts...)
	return
}

func GetPaypalExpressInfo(params files_sdk.SiteGetPaypalExpressInfoParams, opts ...files_sdk.RequestResponseOption) (paypalExpressInfo files_sdk.PaypalExpressInfo, err error) {
	return (&Client{}).GetPaypalExpressInfo(params, opts...)
}

func (c *Client) GetPaypalExpress(params files_sdk.SiteGetPaypalExpressParams, opts ...files_sdk.RequestResponseOption) (paypalExpressUrl files_sdk.PaypalExpressUrl, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/site/paypal_express", Params: params, Entity: &paypalExpressUrl}, opts...)
	return
}

func GetPaypalExpress(params files_sdk.SiteGetPaypalExpressParams, opts ...files_sdk.RequestResponseOption) (paypalExpressUrl files_sdk.PaypalExpressUrl, err error) {
	return (&Client{}).GetPaypalExpress(params, opts...)
}

func (c *Client) GetUsage(opts ...files_sdk.RequestResponseOption) (usageSnapshot files_sdk.UsageSnapshot, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/site/usage", Params: lib.Interface(), Entity: &usageSnapshot}, opts...)
	return
}

func GetUsage(opts ...files_sdk.RequestResponseOption) (usageSnapshot files_sdk.UsageSnapshot, err error) {
	return (&Client{}).GetUsage(opts...)
}

func (c *Client) Create(params files_sdk.SiteCreateParams, opts ...files_sdk.RequestResponseOption) (site files_sdk.Site, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/site", Params: params, Entity: &site}, opts...)
	return
}

func Create(params files_sdk.SiteCreateParams, opts ...files_sdk.RequestResponseOption) (site files_sdk.Site, err error) {
	return (&Client{}).Create(params, opts...)
}

func (c *Client) Update(params files_sdk.SiteUpdateParams, opts ...files_sdk.RequestResponseOption) (site files_sdk.Site, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/site", Params: params, Entity: &site}, opts...)
	return
}

func Update(params files_sdk.SiteUpdateParams, opts ...files_sdk.RequestResponseOption) (site files_sdk.Site, err error) {
	return (&Client{}).Update(params, opts...)
}

func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (site files_sdk.Site, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/site", Params: params, Entity: &site}, opts...)
	return
}

func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (site files_sdk.Site, err error) {
	return (&Client{}).UpdateWithMap(params, opts...)
}

func (c *Client) Delete(params files_sdk.SiteDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "DELETE", Path: "/site", Params: params, Entity: nil}, opts...)
	return
}

func Delete(params files_sdk.SiteDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(params, opts...)
}
