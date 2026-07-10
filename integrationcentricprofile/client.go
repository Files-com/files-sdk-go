package integration_centric_profile

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

func (i *Iter) IntegrationCentricProfile() files_sdk.IntegrationCentricProfile {
	return i.Current().(files_sdk.IntegrationCentricProfile)
}

func (i *Iter) LoadResource(identifier interface{}, opts ...files_sdk.RequestResponseOption) (interface{}, error) {
	params := files_sdk.IntegrationCentricProfileFindParams{}
	if id, ok := identifier.(int64); ok {
		params.Id = id
	}
	return i.Client.Find(params, opts...)
}

func (c *Client) List(params files_sdk.IntegrationCentricProfileListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/integration_centric_profiles", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.IntegrationCentricProfileCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

func List(params files_sdk.IntegrationCentricProfileListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(params, opts...)
}

func (c *Client) Find(params files_sdk.IntegrationCentricProfileFindParams, opts ...files_sdk.RequestResponseOption) (integrationCentricProfile files_sdk.IntegrationCentricProfile, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/integration_centric_profiles/{id}", Params: params, Entity: &integrationCentricProfile}, opts...)
	return
}

func Find(params files_sdk.IntegrationCentricProfileFindParams, opts ...files_sdk.RequestResponseOption) (integrationCentricProfile files_sdk.IntegrationCentricProfile, err error) {
	return (&Client{}).Find(params, opts...)
}

func (c *Client) Create(params files_sdk.IntegrationCentricProfileCreateParams, opts ...files_sdk.RequestResponseOption) (integrationCentricProfile files_sdk.IntegrationCentricProfile, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/integration_centric_profiles", Params: params, Entity: &integrationCentricProfile}, opts...)
	return
}

func Create(params files_sdk.IntegrationCentricProfileCreateParams, opts ...files_sdk.RequestResponseOption) (integrationCentricProfile files_sdk.IntegrationCentricProfile, err error) {
	return (&Client{}).Create(params, opts...)
}

func (c *Client) Update(params files_sdk.IntegrationCentricProfileUpdateParams, opts ...files_sdk.RequestResponseOption) (integrationCentricProfile files_sdk.IntegrationCentricProfile, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/integration_centric_profiles/{id}", Params: params, Entity: &integrationCentricProfile}, opts...)
	return
}

func Update(params files_sdk.IntegrationCentricProfileUpdateParams, opts ...files_sdk.RequestResponseOption) (integrationCentricProfile files_sdk.IntegrationCentricProfile, err error) {
	return (&Client{}).Update(params, opts...)
}

func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (integrationCentricProfile files_sdk.IntegrationCentricProfile, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/integration_centric_profiles/{id}", Params: params, Entity: &integrationCentricProfile}, opts...)
	return
}

func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (integrationCentricProfile files_sdk.IntegrationCentricProfile, err error) {
	return (&Client{}).UpdateWithMap(params, opts...)
}

func (c *Client) Delete(params files_sdk.IntegrationCentricProfileDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "DELETE", Path: "/integration_centric_profiles/{id}", Params: params, Entity: nil}, opts...)
	return
}

func Delete(params files_sdk.IntegrationCentricProfileDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(params, opts...)
}
