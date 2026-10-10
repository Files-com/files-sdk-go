// Package partner_site provides the Files.com PartnerSite API client.
package partner_site

import (
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// Client calls the PartnerSite API using its embedded Config.
type Client struct {
	files_sdk.Config
}

// Delete calls DELETE /partner_sites/{id}.
func (c *Client) Delete(params files_sdk.PartnerSiteDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "DELETE", Path: "/partner_sites/{id}", Params: params, Entity: nil}, opts...)
	return
}

// Delete calls Client.Delete using the default configuration.
func Delete(params files_sdk.PartnerSiteDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(params, opts...)
}
