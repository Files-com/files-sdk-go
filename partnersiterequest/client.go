package partner_site_request

import (
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	lib "github.com/Files-com/files-sdk-go/v3/lib"
	listquery "github.com/Files-com/files-sdk-go/v3/listquery"
)

type Client struct {
	files_sdk.Config
}

type Iter struct {
	*files_sdk.Iter
	*Client
}

func (i *Iter) Reload(opts ...files_sdk.RequestResponseOption) files_sdk.IterI {
	return &Iter{Iter: i.Iter.Reload(opts...).(*files_sdk.Iter), Client: i.Client}
}

func (i *Iter) PartnerSiteRequest() files_sdk.PartnerSiteRequest {
	return i.Current().(files_sdk.PartnerSiteRequest)
}

func (c *Client) List(params files_sdk.PartnerSiteRequestListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/partner_site_requests", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.PartnerSiteRequestCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

func List(params files_sdk.PartnerSiteRequestListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(params, opts...)
}

func (c *Client) FindByPairingKey(params files_sdk.PartnerSiteRequestFindByPairingKeyParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/partner_site_requests/find_by_pairing_key", Params: params, Entity: nil}, opts...)
	return
}

func FindByPairingKey(params files_sdk.PartnerSiteRequestFindByPairingKeyParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).FindByPairingKey(params, opts...)
}

func (c *Client) Create(params files_sdk.PartnerSiteRequestCreateParams, opts ...files_sdk.RequestResponseOption) (partnerSiteRequest files_sdk.PartnerSiteRequest, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/partner_site_requests", Params: params, Entity: &partnerSiteRequest}, opts...)
	return
}

func Create(params files_sdk.PartnerSiteRequestCreateParams, opts ...files_sdk.RequestResponseOption) (partnerSiteRequest files_sdk.PartnerSiteRequest, err error) {
	return (&Client{}).Create(params, opts...)
}

func (c *Client) Reject(params files_sdk.PartnerSiteRequestRejectParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/partner_site_requests/{id}/reject", Params: params, Entity: nil}, opts...)
	return
}

func Reject(params files_sdk.PartnerSiteRequestRejectParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Reject(params, opts...)
}

func (c *Client) Approve(params files_sdk.PartnerSiteRequestApproveParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/partner_site_requests/{id}/approve", Params: params, Entity: nil}, opts...)
	return
}

func Approve(params files_sdk.PartnerSiteRequestApproveParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Approve(params, opts...)
}

func (c *Client) Delete(params files_sdk.PartnerSiteRequestDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "DELETE", Path: "/partner_site_requests/{id}", Params: params, Entity: nil}, opts...)
	return
}

func Delete(params files_sdk.PartnerSiteRequestDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(params, opts...)
}
