package expectation

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

func (i *Iter) Expectation() files_sdk.Expectation {
	return i.Current().(files_sdk.Expectation)
}

func (i *Iter) LoadResource(identifier interface{}, opts ...files_sdk.RequestResponseOption) (interface{}, error) {
	params := files_sdk.ExpectationFindParams{}
	if id, ok := identifier.(int64); ok {
		params.Id = id
	}
	return i.Client.Find(params, opts...)
}

func (c *Client) List(params files_sdk.ExpectationListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/expectations", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.ExpectationCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

func List(params files_sdk.ExpectationListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(params, opts...)
}

func (c *Client) Find(params files_sdk.ExpectationFindParams, opts ...files_sdk.RequestResponseOption) (expectation files_sdk.Expectation, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/expectations/{id}", Params: params, Entity: &expectation}, opts...)
	return
}

func Find(params files_sdk.ExpectationFindParams, opts ...files_sdk.RequestResponseOption) (expectation files_sdk.Expectation, err error) {
	return (&Client{}).Find(params, opts...)
}

func (c *Client) Create(params files_sdk.ExpectationCreateParams, opts ...files_sdk.RequestResponseOption) (expectation files_sdk.Expectation, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/expectations", Params: params, Entity: &expectation}, opts...)
	return
}

func Create(params files_sdk.ExpectationCreateParams, opts ...files_sdk.RequestResponseOption) (expectation files_sdk.Expectation, err error) {
	return (&Client{}).Create(params, opts...)
}

func (c *Client) Trigger(params files_sdk.ExpectationTriggerParams, opts ...files_sdk.RequestResponseOption) (expectationEvaluation files_sdk.ExpectationEvaluation, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/expectations/{id}/trigger", Params: params, Entity: &expectationEvaluation}, opts...)
	return
}

func Trigger(params files_sdk.ExpectationTriggerParams, opts ...files_sdk.RequestResponseOption) (expectationEvaluation files_sdk.ExpectationEvaluation, err error) {
	return (&Client{}).Trigger(params, opts...)
}

func (c *Client) Update(params files_sdk.ExpectationUpdateParams, opts ...files_sdk.RequestResponseOption) (expectation files_sdk.Expectation, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/expectations/{id}", Params: params, Entity: &expectation}, opts...)
	return
}

func Update(params files_sdk.ExpectationUpdateParams, opts ...files_sdk.RequestResponseOption) (expectation files_sdk.Expectation, err error) {
	return (&Client{}).Update(params, opts...)
}

func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (expectation files_sdk.Expectation, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/expectations/{id}", Params: params, Entity: &expectation}, opts...)
	return
}

func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (expectation files_sdk.Expectation, err error) {
	return (&Client{}).UpdateWithMap(params, opts...)
}

func (c *Client) Delete(params files_sdk.ExpectationDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "DELETE", Path: "/expectations/{id}", Params: params, Entity: nil}, opts...)
	return
}

func Delete(params files_sdk.ExpectationDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(params, opts...)
}
