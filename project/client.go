package project

import (
	"strconv"

	files_sdk "github.com/Files-com/files-sdk-go"
	lib "github.com/Files-com/files-sdk-go/lib"
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

func (c *Client) List(params files_sdk.ProjectListParams) *Iter {
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	i := &Iter{Iter: &lib.Iter{}}
	path := "/projects"

	i.Query = func() (*[]interface{}, string, error) {
		data, res, err := files_sdk.Call("GET", c.Config, path, i.ExportParams())
		defaultValue := make([]interface{}, 0)
		if err != nil {
			return &defaultValue, "", err
		}
		list := files_sdk.ProjectCollection{}
		if err := list.UnmarshalJSON(*data); err != nil {
			return &defaultValue, "", err
		}

		ret := make([]interface{}, len(list))
		for i, v := range list {
			ret[i] = v
		}
		cursor := res.Header.Get("X-Files-Cursor")
		return &ret, cursor, nil
	}
	i.ListParams = &params
	return i
}

func List(params files_sdk.ProjectListParams) *Iter {
	return (&Client{}).List(params)
}

func (c *Client) Find(params files_sdk.ProjectFindParams) (files_sdk.Project, error) {
	project := files_sdk.Project{}
	path := "/projects/" + lib.QueryEscape(strconv.FormatInt(params.Id, 10)) + ""
	data, res, err := files_sdk.Call("GET", c.Config, path, lib.ExportParams(params))
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
	data, res, err := files_sdk.Call("POST", c.Config, path, lib.ExportParams(params))
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
	path := "/projects/" + lib.QueryEscape(strconv.FormatInt(params.Id, 10)) + ""
	data, res, err := files_sdk.Call("PATCH", c.Config, path, lib.ExportParams(params))
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
	path := "/projects/" + lib.QueryEscape(strconv.FormatInt(params.Id, 10)) + ""
	data, res, err := files_sdk.Call("DELETE", c.Config, path, lib.ExportParams(params))
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
