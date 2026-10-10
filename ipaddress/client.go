// Package ip_address provides the Files.com IpAddress API client.
package ip_address

import (
	"iter"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	lib "github.com/Files-com/files-sdk-go/v3/lib"
	listquery "github.com/Files-com/files-sdk-go/v3/listquery"
)

// Client calls the IpAddress API using its embedded Config.
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
//
// This package's listings return different types, so All yields each resource
// unchanged, as an any value; use a type assertion or type switch to read it.
// List yields files_sdk.IpAddress values, and GetSmartfileReserved, GetExavaultReserved, and GetReserved yield files_sdk.PublicIpAddress values.
func (i *Iter) All() iter.Seq2[any, error] {
	return files_sdk.IterAll[any](i.Iter)
}

// IpAddress returns the current resource. Call it only after Next returns true.
func (i *Iter) IpAddress() files_sdk.IpAddress {
	return i.Current().(files_sdk.IpAddress)
}

// List returns an iterator for GET /ip_addresses.
// Pages are requested as iteration reaches them: range over Iter.All, or call
// Next and then check Err.
// Iter.All yields this listing's resources as files_sdk.IpAddress values.
//
// API operation: List IP Addresses associated with the current site.
func (c *Client) List(params files_sdk.IpAddressListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/ip_addresses", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.IpAddressCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

// List returns a listing iterator using the default configuration.
// See Client.List for the operation and paging behavior.
func List(params files_sdk.IpAddressListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(params, opts...)
}

// PublicIpAddress returns the current resource. Call it only after Next returns true.
func (i *Iter) PublicIpAddress() files_sdk.PublicIpAddress {
	return i.Current().(files_sdk.PublicIpAddress)
}

// GetSmartfileReserved returns an iterator for GET /ip_addresses/smartfile-reserved.
// Pages are requested as iteration reaches them: range over Iter.All, or call
// Next and then check Err.
// Iter.All yields this listing's resources as files_sdk.PublicIpAddress values.
//
// API operation: List all possible public SmartFile IP addresses.
func (c *Client) GetSmartfileReserved(params files_sdk.IpAddressGetSmartfileReservedParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/ip_addresses/smartfile-reserved", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.PublicIpAddressCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

// GetSmartfileReserved returns a listing iterator using the default configuration.
// See Client.GetSmartfileReserved for the operation and paging behavior.
func GetSmartfileReserved(params files_sdk.IpAddressGetSmartfileReservedParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).GetSmartfileReserved(params, opts...)
}

// GetExavaultReserved returns an iterator for GET /ip_addresses/exavault-reserved.
// Pages are requested as iteration reaches them: range over Iter.All, or call
// Next and then check Err.
// Iter.All yields this listing's resources as files_sdk.PublicIpAddress values.
//
// API operation: List all possible public ExaVault IP addresses.
func (c *Client) GetExavaultReserved(params files_sdk.IpAddressGetExavaultReservedParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/ip_addresses/exavault-reserved", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.PublicIpAddressCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

// GetExavaultReserved returns a listing iterator using the default configuration.
// See Client.GetExavaultReserved for the operation and paging behavior.
func GetExavaultReserved(params files_sdk.IpAddressGetExavaultReservedParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).GetExavaultReserved(params, opts...)
}

// GetReserved returns an iterator for GET /ip_addresses/reserved.
// Pages are requested as iteration reaches them: range over Iter.All, or call
// Next and then check Err.
// Iter.All yields this listing's resources as files_sdk.PublicIpAddress values.
//
// API operation: List all possible public IP addresses.
func (c *Client) GetReserved(params files_sdk.IpAddressGetReservedParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/ip_addresses/reserved", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.PublicIpAddressCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

// GetReserved returns a listing iterator using the default configuration.
// See Client.GetReserved for the operation and paging behavior.
func GetReserved(params files_sdk.IpAddressGetReservedParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).GetReserved(params, opts...)
}
