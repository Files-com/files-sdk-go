package project

import (
	"strconv"

	files_sdk "github.com/Files-com/files-sdk-go"
	lib "github.com/Files-com/files-sdk-go/lib"
	listquery "github.com/Files-com/files-sdk-go/listquery"
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

func (c *Client) List(params files_sdk.ProjectListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/projects"
	i.ListParams = &params
	list := files_sdk.ProjectCollection{}
	i.Query = listquery.Build(i, c.Config, path, &list)
	return i, nil
}

func List(params files_sdk.ProjectListParams) (*Iter, error) {
	return (&Client{}).List(params)
}

func (c *Client) Find(params files_sdk.ProjectFindParams) (files_sdk.Project, error) {
	project := files_sdk.Project{}
	if params.Id == 0 {
		return project, lib.CreateError(params, "Id")
	}
	path := "/projects/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return project, err
	}
	data, res, err := files_sdk.Call("GET", c.Config, path, exportedParams)
	if err != nil {
		return project, err
	}
	if res.StatusCode == 204 {
		return project, nil
	}
	if err := project.UnmarshalJSON(*data); err != nil {
		return project, err
	}

	return project, nil
}

func Find(params files_sdk.ProjectFindParams) (files_sdk.Project, error) {
	return (&Client{}).Find(params)
}

func (c *Client) Create(params files_sdk.ProjectCreateParams) (files_sdk.Project, error) {
	project := files_sdk.Project{}
	path := "/projects"
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return project, err
	}
	data, res, err := files_sdk.Call("POST", c.Config, path, exportedParams)
	if err != nil {
		return project, err
	}
	if res.StatusCode == 204 {
		return project, nil
	}
	if err := project.UnmarshalJSON(*data); err != nil {
		return project, err
	}

	return project, nil
}

func Create(params files_sdk.ProjectCreateParams) (files_sdk.Project, error) {
	return (&Client{}).Create(params)
}

func (c *Client) Update(params files_sdk.ProjectUpdateParams) (files_sdk.Project, error) {
	project := files_sdk.Project{}
	if params.Id == 0 {
		return project, lib.CreateError(params, "Id")
	}
	path := "/projects/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return project, err
	}
	data, res, err := files_sdk.Call("PATCH", c.Config, path, exportedParams)
	if err != nil {
		return project, err
	}
	if res.StatusCode == 204 {
		return project, nil
	}
	if err := project.UnmarshalJSON(*data); err != nil {
		return project, err
	}

	return project, nil
}

func Update(params files_sdk.ProjectUpdateParams) (files_sdk.Project, error) {
	return (&Client{}).Update(params)
}

func (c *Client) Delete(params files_sdk.ProjectDeleteParams) (files_sdk.Project, error) {
	project := files_sdk.Project{}
	if params.Id == 0 {
		return project, lib.CreateError(params, "Id")
	}
	path := "/projects/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return project, err
	}
	data, res, err := files_sdk.Call("DELETE", c.Config, path, exportedParams)
	if err != nil {
		return project, err
	}
	if res.StatusCode == 204 {
		return project, nil
	}
	if err := project.UnmarshalJSON(*data); err != nil {
		return project, err
	}

	return project, nil
}

func Delete(params files_sdk.ProjectDeleteParams) (files_sdk.Project, error) {
	return (&Client{}).Delete(params)
}
