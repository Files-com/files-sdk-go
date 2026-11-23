// Package certificate is the Files.com API client package for the files_sdk.Certificate model.
package certificate

import (
	"iter"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	lib "github.com/Files-com/files-sdk-go/v3/lib"
	listquery "github.com/Files-com/files-sdk-go/v3/listquery"
)

// Client calls the Certificate API using its embedded Config.
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
func (i *Iter) All() iter.Seq2[files_sdk.Certificate, error] {
	return files_sdk.IterAll[files_sdk.Certificate](i.Iter)
}

// Certificate returns the current resource. Call it only after Next returns true.
func (i *Iter) Certificate() files_sdk.Certificate {
	return i.Current().(files_sdk.Certificate)
}

// LoadResource fetches a resource by its ID. identifier must have type
// int64.
func (i *Iter) LoadResource(identifier interface{}, opts ...files_sdk.RequestResponseOption) (interface{}, error) {
	params := files_sdk.CertificateFindParams{}
	if id, ok := identifier.(int64); ok {
		params.Id = id
	}
	return i.Client.Find(params, opts...)
}

// List returns an iterator for GET /certificates.
// Pages are requested as iteration reaches them: range over Iter.All, or call
// Next and then check Err.
func (c *Client) List(params files_sdk.CertificateListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/certificates", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.CertificateCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

// List returns a listing iterator using the default configuration.
// See Client.List for the operation and paging behavior.
func List(params files_sdk.CertificateListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(params, opts...)
}

// Find calls GET /certificates/{id}.
func (c *Client) Find(params files_sdk.CertificateFindParams, opts ...files_sdk.RequestResponseOption) (certificate files_sdk.Certificate, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/certificates/{id}", Params: params, Entity: &certificate}, opts...)
	return
}

// Find calls Client.Find using the default configuration.
func Find(params files_sdk.CertificateFindParams, opts ...files_sdk.RequestResponseOption) (certificate files_sdk.Certificate, err error) {
	return (&Client{}).Find(params, opts...)
}

// Create calls POST /certificates.
//
// API operation: Create CSR or import existing SSL Certificate.
func (c *Client) Create(params files_sdk.CertificateCreateParams, opts ...files_sdk.RequestResponseOption) (certificate files_sdk.Certificate, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/certificates", Params: params, Entity: &certificate}, opts...)
	return
}

// Create calls Client.Create using the default configuration.
func Create(params files_sdk.CertificateCreateParams, opts ...files_sdk.RequestResponseOption) (certificate files_sdk.Certificate, err error) {
	return (&Client{}).Create(params, opts...)
}

// Activate calls POST /certificates/{id}/activate.
//
// API operation: Activate SSL Certificate.
func (c *Client) Activate(params files_sdk.CertificateActivateParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/certificates/{id}/activate", Params: params, Entity: nil}, opts...)
	return
}

// Activate calls Client.Activate using the default configuration.
func Activate(params files_sdk.CertificateActivateParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Activate(params, opts...)
}

// Deactivate calls POST /certificates/{id}/deactivate.
//
// API operation: Deactivate SSL Certificate.
func (c *Client) Deactivate(params files_sdk.CertificateDeactivateParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/certificates/{id}/deactivate", Params: params, Entity: nil}, opts...)
	return
}

// Deactivate calls Client.Deactivate using the default configuration.
func Deactivate(params files_sdk.CertificateDeactivateParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Deactivate(params, opts...)
}

// ReplaceWithFiles calls POST /certificates/{id}/replace_with_files_certificate.
//
// API operation: Replace this active custom SSL Certificate for the Primary Custom Domain with a Files.com SSL Certificate.  The Primary Custom Domain's CNAME record must point at the Site's subdomain.  The Files.com certificate is issued and activated in the background.
func (c *Client) ReplaceWithFiles(params files_sdk.CertificateReplaceWithFilesParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/certificates/{id}/replace_with_files_certificate", Params: params, Entity: nil}, opts...)
	return
}

// ReplaceWithFiles calls Client.ReplaceWithFiles using the default configuration.
func ReplaceWithFiles(params files_sdk.CertificateReplaceWithFilesParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).ReplaceWithFiles(params, opts...)
}

// Update calls PATCH /certificates/{id}.
func (c *Client) Update(params files_sdk.CertificateUpdateParams, opts ...files_sdk.RequestResponseOption) (certificate files_sdk.Certificate, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/certificates/{id}", Params: params, Entity: &certificate}, opts...)
	return
}

// Update calls Client.Update using the default configuration.
func Update(params files_sdk.CertificateUpdateParams, opts ...files_sdk.RequestResponseOption) (certificate files_sdk.Certificate, err error) {
	return (&Client{}).Update(params, opts...)
}

// UpdateWithMap calls PATCH /certificates/{id} using API parameter names as map keys.
// Include any path parameters in the map. Unlike optional struct fields, explicit
// zero values in the map are included in the request.
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (certificate files_sdk.Certificate, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/certificates/{id}", Params: params, Entity: &certificate}, opts...)
	return
}

// UpdateWithMap calls Client.UpdateWithMap using the default configuration.
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (certificate files_sdk.Certificate, err error) {
	return (&Client{}).UpdateWithMap(params, opts...)
}

// Delete calls DELETE /certificates/{id}.
func (c *Client) Delete(params files_sdk.CertificateDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "DELETE", Path: "/certificates/{id}", Params: params, Entity: nil}, opts...)
	return
}

// Delete calls Client.Delete using the default configuration.
func Delete(params files_sdk.CertificateDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(params, opts...)
}
