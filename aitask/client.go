package ai_task

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

func (i *Iter) AiTask() files_sdk.AiTask {
	return i.Current().(files_sdk.AiTask)
}

func (i *Iter) LoadResource(identifier interface{}, opts ...files_sdk.RequestResponseOption) (interface{}, error) {
	params := files_sdk.AiTaskFindParams{}
	if id, ok := identifier.(int64); ok {
		params.Id = id
	}
	return i.Client.Find(params, opts...)
}

func (c *Client) List(params files_sdk.AiTaskListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/ai_tasks", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.AiTaskCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

func List(params files_sdk.AiTaskListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(params, opts...)
}

func (c *Client) Find(params files_sdk.AiTaskFindParams, opts ...files_sdk.RequestResponseOption) (aiTask files_sdk.AiTask, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/ai_tasks/{id}", Params: params, Entity: &aiTask}, opts...)
	return
}

func Find(params files_sdk.AiTaskFindParams, opts ...files_sdk.RequestResponseOption) (aiTask files_sdk.AiTask, err error) {
	return (&Client{}).Find(params, opts...)
}

func (c *Client) Create(params files_sdk.AiTaskCreateParams, opts ...files_sdk.RequestResponseOption) (aiTask files_sdk.AiTask, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/ai_tasks", Params: params, Entity: &aiTask}, opts...)
	return
}

func Create(params files_sdk.AiTaskCreateParams, opts ...files_sdk.RequestResponseOption) (aiTask files_sdk.AiTask, err error) {
	return (&Client{}).Create(params, opts...)
}

func (c *Client) ManualRun(params files_sdk.AiTaskManualRunParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/ai_tasks/{id}/manual_run", Params: params, Entity: nil}, opts...)
	return
}

func ManualRun(params files_sdk.AiTaskManualRunParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).ManualRun(params, opts...)
}

func (c *Client) Update(params files_sdk.AiTaskUpdateParams, opts ...files_sdk.RequestResponseOption) (aiTask files_sdk.AiTask, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/ai_tasks/{id}", Params: params, Entity: &aiTask}, opts...)
	return
}

func Update(params files_sdk.AiTaskUpdateParams, opts ...files_sdk.RequestResponseOption) (aiTask files_sdk.AiTask, err error) {
	return (&Client{}).Update(params, opts...)
}

func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (aiTask files_sdk.AiTask, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/ai_tasks/{id}", Params: params, Entity: &aiTask}, opts...)
	return
}

func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (aiTask files_sdk.AiTask, err error) {
	return (&Client{}).UpdateWithMap(params, opts...)
}

func (c *Client) Delete(params files_sdk.AiTaskDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "DELETE", Path: "/ai_tasks/{id}", Params: params, Entity: nil}, opts...)
	return
}

func Delete(params files_sdk.AiTaskDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(params, opts...)
}
