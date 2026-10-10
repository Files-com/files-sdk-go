// Package partner_site_request provides the Files.com PartnerSiteRequest API client.
package partner_site_request

import (
	"iter"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	lib "github.com/Files-com/files-sdk-go/v3/lib"
	listquery "github.com/Files-com/files-sdk-go/v3/listquery"
)

// Client calls the PartnerSiteRequest API using its embedded Config.
type Client struct {
	files_sdk.Config
}

// Iter traverses a paginated API response. Range over All to get each resource
// with a nil error, followed by the error that stopped the listing, if any.
// Alternatively, call Next before reading the current resource, then check Err
// after Next returns false.
type Iter struct {
	*files_sdk.Iter
	*Client
}

// Reload returns a new iterator for the same listing, starting at the first page.
// See files_sdk.Iter.Reload for parameter and request-option handling.
func (i *Iter) Reload(opts ...files_sdk.RequestResponseOption) files_sdk.IterI {
	return &Iter{Iter: i.Iter.Reload(opts...).(*files_sdk.Iter), Client: i.Client}
}

// All returns an iterator over the listing's remaining resources, each with a
// nil error, followed by the error that stops the listing, if any. Pages are
// requested as the loop reaches them. See files_sdk.IterAll.
func (i *Iter) All() iter.Seq2[files_sdk.PartnerSiteRequest, error] {
	return files_sdk.IterAll[files_sdk.PartnerSiteRequest](i.Iter)
}

// PartnerSiteRequest returns the current resource. Call it only after Next returns true.
func (i *Iter) PartnerSiteRequest() files_sdk.PartnerSiteRequest {
	return i.Current().(files_sdk.PartnerSiteRequest)
}

// List returns an iterator for GET /partner_site_requests.
// Pages are requested as iteration reaches them: range over Iter.All, or call
// Next and then check Err.
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

// List returns a listing iterator using the default configuration.
// See Client.List for the operation and paging behavior.
func List(params files_sdk.PartnerSiteRequestListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(params, opts...)
}

// FindByPairingKey calls GET /partner_site_requests/find_by_pairing_key.
//
// API operation: Find partner site request by pairing key.
func (c *Client) FindByPairingKey(params files_sdk.PartnerSiteRequestFindByPairingKeyParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/partner_site_requests/find_by_pairing_key", Params: params, Entity: nil}, opts...)
	return
}

// FindByPairingKey calls Client.FindByPairingKey using the default configuration.
func FindByPairingKey(params files_sdk.PartnerSiteRequestFindByPairingKeyParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).FindByPairingKey(params, opts...)
}

// Create calls POST /partner_site_requests.
func (c *Client) Create(params files_sdk.PartnerSiteRequestCreateParams, opts ...files_sdk.RequestResponseOption) (partnerSiteRequest files_sdk.PartnerSiteRequest, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/partner_site_requests", Params: params, Entity: &partnerSiteRequest}, opts...)
	return
}

// Create calls Client.Create using the default configuration.
func Create(params files_sdk.PartnerSiteRequestCreateParams, opts ...files_sdk.RequestResponseOption) (partnerSiteRequest files_sdk.PartnerSiteRequest, err error) {
	return (&Client{}).Create(params, opts...)
}

// Reject calls POST /partner_site_requests/reject.
//
// API operation: Reject partner site request.
func (c *Client) Reject(params files_sdk.PartnerSiteRequestRejectParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/partner_site_requests/reject", Params: params, Entity: nil}, opts...)
	return
}

// Reject calls Client.Reject using the default configuration.
func Reject(params files_sdk.PartnerSiteRequestRejectParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Reject(params, opts...)
}

// Approve calls POST /partner_site_requests/approve.
//
// API operation: Approve partner site request.
func (c *Client) Approve(params files_sdk.PartnerSiteRequestApproveParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/partner_site_requests/approve", Params: params, Entity: nil}, opts...)
	return
}

// Approve calls Client.Approve using the default configuration.
func Approve(params files_sdk.PartnerSiteRequestApproveParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Approve(params, opts...)
}

// Delete calls DELETE /partner_site_requests/{id}.
func (c *Client) Delete(params files_sdk.PartnerSiteRequestDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "DELETE", Path: "/partner_site_requests/{id}", Params: params, Entity: nil}, opts...)
	return
}

// Delete calls Client.Delete using the default configuration.
func Delete(params files_sdk.PartnerSiteRequestDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(params, opts...)
}
