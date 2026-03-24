package expectation_incident

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

func (i *Iter) ExpectationIncident() files_sdk.ExpectationIncident {
	return i.Current().(files_sdk.ExpectationIncident)
}

func (i *Iter) LoadResource(identifier interface{}, opts ...files_sdk.RequestResponseOption) (interface{}, error) {
	params := files_sdk.ExpectationIncidentFindParams{}
	if id, ok := identifier.(int64); ok {
		params.Id = id
	}
	return i.Client.Find(params, opts...)
}

func (c *Client) List(params files_sdk.ExpectationIncidentListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/expectation_incidents", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.ExpectationIncidentCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

func List(params files_sdk.ExpectationIncidentListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(params, opts...)
}

func (c *Client) Find(params files_sdk.ExpectationIncidentFindParams, opts ...files_sdk.RequestResponseOption) (expectationIncident files_sdk.ExpectationIncident, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/expectation_incidents/{id}", Params: params, Entity: &expectationIncident}, opts...)
	return
}

func Find(params files_sdk.ExpectationIncidentFindParams, opts ...files_sdk.RequestResponseOption) (expectationIncident files_sdk.ExpectationIncident, err error) {
	return (&Client{}).Find(params, opts...)
}

func (c *Client) Resolve(params files_sdk.ExpectationIncidentResolveParams, opts ...files_sdk.RequestResponseOption) (expectationIncident files_sdk.ExpectationIncident, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/expectation_incidents/{id}/resolve", Params: params, Entity: &expectationIncident}, opts...)
	return
}

func Resolve(params files_sdk.ExpectationIncidentResolveParams, opts ...files_sdk.RequestResponseOption) (expectationIncident files_sdk.ExpectationIncident, err error) {
	return (&Client{}).Resolve(params, opts...)
}

func (c *Client) Snooze(params files_sdk.ExpectationIncidentSnoozeParams, opts ...files_sdk.RequestResponseOption) (expectationIncident files_sdk.ExpectationIncident, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/expectation_incidents/{id}/snooze", Params: params, Entity: &expectationIncident}, opts...)
	return
}

func Snooze(params files_sdk.ExpectationIncidentSnoozeParams, opts ...files_sdk.RequestResponseOption) (expectationIncident files_sdk.ExpectationIncident, err error) {
	return (&Client{}).Snooze(params, opts...)
}

func (c *Client) Acknowledge(params files_sdk.ExpectationIncidentAcknowledgeParams, opts ...files_sdk.RequestResponseOption) (expectationIncident files_sdk.ExpectationIncident, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/expectation_incidents/{id}/acknowledge", Params: params, Entity: &expectationIncident}, opts...)
	return
}

func Acknowledge(params files_sdk.ExpectationIncidentAcknowledgeParams, opts ...files_sdk.RequestResponseOption) (expectationIncident files_sdk.ExpectationIncident, err error) {
	return (&Client{}).Acknowledge(params, opts...)
}
