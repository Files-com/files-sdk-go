package project

import (
	"context"
	"strconv"

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

func (c *Client) List(ctx context.Context, params files_sdk.ProjectListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/projects"
	i.ListParams = &params
	list := files_sdk.ProjectCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.ProjectListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}

func (c *Client) Find(ctx context.Context, params files_sdk.ProjectFindParams) (files_sdk.Project, error) {
	project := files_sdk.Project{}
	if params.Id == 0 {
		return project, lib.CreateError(params, "Id")
	}
	path := "/projects/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "GET", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return project, err
	}
	if res.StatusCode == 204 {
		return project, nil
	}

	return project, project.UnmarshalJSON(*data)
}

func Find(ctx context.Context, params files_sdk.ProjectFindParams) (files_sdk.Project, error) {
	return (&Client{}).Find(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.ProjectCreateParams) (files_sdk.Project, error) {
	project := files_sdk.Project{}
	path := "/projects"
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return project, err
	}
	if res.StatusCode == 204 {
		return project, nil
	}

	return project, project.UnmarshalJSON(*data)
}

func Create(ctx context.Context, params files_sdk.ProjectCreateParams) (files_sdk.Project, error) {
	return (&Client{}).Create(ctx, params)
}

func (c *Client) Update(ctx context.Context, params files_sdk.ProjectUpdateParams) (files_sdk.Project, error) {
	project := files_sdk.Project{}
	if params.Id == 0 {
		return project, lib.CreateError(params, "Id")
	}
	path := "/projects/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "PATCH", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return project, err
	}
	if res.StatusCode == 204 {
		return project, nil
	}

	return project, project.UnmarshalJSON(*data)
}

func Update(ctx context.Context, params files_sdk.ProjectUpdateParams) (files_sdk.Project, error) {
	return (&Client{}).Update(ctx, params)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.ProjectDeleteParams) error {
	project := files_sdk.Project{}
	if params.Id == 0 {
		return lib.CreateError(params, "Id")
	}
	path := "/projects/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "DELETE", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return err
	}
	if res.StatusCode == 204 {
		return nil
	}

	return project.UnmarshalJSON(*data)
}

func Delete(ctx context.Context, params files_sdk.ProjectDeleteParams) error {
	return (&Client{}).Delete(ctx, params)
}
