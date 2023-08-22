package certificate

import (
	files_sdk "github.com/Files-com/files-sdk-go/v2"
	lib "github.com/Files-com/files-sdk-go/v2/lib"
	listquery "github.com/Files-com/files-sdk-go/v2/listquery"
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

func (i *Iter) Certificate() files_sdk.Certificate {
	return i.Current().(files_sdk.Certificate)
}

func (i *Iter) LoadResource(identifier interface{}, opts ...files_sdk.RequestResponseOption) (interface{}, error) {
	params := files_sdk.CertificateFindParams{}
	if id, ok := identifier.(int64); ok {
		params.Id = id
	}
	return i.Client.Find(params, opts...)
}

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

func List(params files_sdk.CertificateListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(params, opts...)
}

func (c *Client) Find(params files_sdk.CertificateFindParams, opts ...files_sdk.RequestResponseOption) (certificate files_sdk.Certificate, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/certificates/{id}", Params: params, Entity: &certificate}, opts...)
	return
}

func Find(params files_sdk.CertificateFindParams, opts ...files_sdk.RequestResponseOption) (certificate files_sdk.Certificate, err error) {
	return (&Client{}).Find(params, opts...)
}

func (c *Client) Create(params files_sdk.CertificateCreateParams, opts ...files_sdk.RequestResponseOption) (certificate files_sdk.Certificate, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/certificates", Params: params, Entity: &certificate}, opts...)
	return
}

func Create(params files_sdk.CertificateCreateParams, opts ...files_sdk.RequestResponseOption) (certificate files_sdk.Certificate, err error) {
	return (&Client{}).Create(params, opts...)
}

func (c *Client) Deactivate(params files_sdk.CertificateDeactivateParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/certificates/{id}/deactivate", Params: params, Entity: nil}, opts...)
	return
}

func Deactivate(params files_sdk.CertificateDeactivateParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Deactivate(params, opts...)
}

func (c *Client) Activate(params files_sdk.CertificateActivateParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/certificates/{id}/activate", Params: params, Entity: nil}, opts...)
	return
}

func Activate(params files_sdk.CertificateActivateParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Activate(params, opts...)
}

func (c *Client) Update(params files_sdk.CertificateUpdateParams, opts ...files_sdk.RequestResponseOption) (certificate files_sdk.Certificate, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/certificates/{id}", Params: params, Entity: &certificate}, opts...)
	return
}

func Update(params files_sdk.CertificateUpdateParams, opts ...files_sdk.RequestResponseOption) (certificate files_sdk.Certificate, err error) {
	return (&Client{}).Update(params, opts...)
}

func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (certificate files_sdk.Certificate, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/certificates/{id}", Params: params, Entity: &certificate}, opts...)
	return
}

func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (certificate files_sdk.Certificate, err error) {
	return (&Client{}).UpdateWithMap(params, opts...)
}

func (c *Client) Delete(params files_sdk.CertificateDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "DELETE", Path: "/certificates/{id}", Params: params, Entity: nil}, opts...)
	return
}

func Delete(params files_sdk.CertificateDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(params, opts...)
}
