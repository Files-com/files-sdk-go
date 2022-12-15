package project

import (
	"context"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	lib "github.com/Files-com/files-sdk-go/v2/lib"
	listquery "github.com/Files-com/files-sdk-go/v2/listquery"
)

type Client struct {
	files_sdk.Config
}

type Iter struct {
	*lib.Iter
}

func (i *Iter) Project() files_sdk.Project {
	return i.Current().(files_sdk.Project)
}

func (c *Client) List(ctx context.Context, params files_sdk.ProjectListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	path, err := lib.BuildPath("/projects", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.ProjectCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list, opts...)
	return i, nil
}

func List(ctx context.Context, params files_sdk.ProjectListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(ctx, params, opts...)
}

func (c *Client) Find(ctx context.Context, params files_sdk.ProjectFindParams, opts ...files_sdk.RequestResponseOption) (project files_sdk.Project, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/projects/{id}", Params: params, Entity: &project}, opts...)
	return
}

func Find(ctx context.Context, params files_sdk.ProjectFindParams, opts ...files_sdk.RequestResponseOption) (project files_sdk.Project, err error) {
	return (&Client{}).Find(ctx, params, opts...)
}

func (c *Client) Create(ctx context.Context, params files_sdk.ProjectCreateParams, opts ...files_sdk.RequestResponseOption) (project files_sdk.Project, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/projects", Params: params, Entity: &project}, opts...)
	return
}

func Create(ctx context.Context, params files_sdk.ProjectCreateParams, opts ...files_sdk.RequestResponseOption) (project files_sdk.Project, err error) {
	return (&Client{}).Create(ctx, params, opts...)
}

func (c *Client) Update(ctx context.Context, params files_sdk.ProjectUpdateParams, opts ...files_sdk.RequestResponseOption) (project files_sdk.Project, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "PATCH", Path: "/projects/{id}", Params: params, Entity: &project}, opts...)
	return
}

func Update(ctx context.Context, params files_sdk.ProjectUpdateParams, opts ...files_sdk.RequestResponseOption) (project files_sdk.Project, err error) {
	return (&Client{}).Update(ctx, params, opts...)
}

func (c *Client) UpdateWithMap(ctx context.Context, params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (project files_sdk.Project, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "PATCH", Path: "/projects/{id}", Params: params, Entity: &project}, opts...)
	return
}

func UpdateWithMap(ctx context.Context, params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (project files_sdk.Project, err error) {
	return (&Client{}).UpdateWithMap(ctx, params, opts...)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.ProjectDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "DELETE", Path: "/projects/{id}", Params: params, Entity: nil}, opts...)
	return
}

func Delete(ctx context.Context, params files_sdk.ProjectDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(ctx, params, opts...)
}
