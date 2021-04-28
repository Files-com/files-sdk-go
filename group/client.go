package group

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

func (i *Iter) Group() files_sdk.Group {
	return i.Current().(files_sdk.Group)
}

func (c *Client) List(params files_sdk.GroupListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/groups"
	i.ListParams = &params
	list := files_sdk.GroupCollection{}
	i.Query = listquery.Build(i, c.Config, path, &list)
	return i, nil
}

func List(params files_sdk.GroupListParams) (*Iter, error) {
	return (&Client{}).List(params)
}

func (c *Client) Find(params files_sdk.GroupFindParams) (files_sdk.Group, error) {
	group := files_sdk.Group{}
	if params.Id == 0 {
		return group, lib.CreateError(params, "Id")
	}
	path := "/groups/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return group, err
	}
	data, res, err := files_sdk.Call("GET", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return group, err
	}
	if res.StatusCode == 204 {
		return group, nil
	}
	if err := group.UnmarshalJSON(*data); err != nil {
		return group, err
	}

	return group, nil
}

func Find(params files_sdk.GroupFindParams) (files_sdk.Group, error) {
	return (&Client{}).Find(params)
}

func (c *Client) Create(params files_sdk.GroupCreateParams) (files_sdk.Group, error) {
	group := files_sdk.Group{}
	path := "/groups"
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return group, err
	}
	data, res, err := files_sdk.Call("POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return group, err
	}
	if res.StatusCode == 204 {
		return group, nil
	}
	if err := group.UnmarshalJSON(*data); err != nil {
		return group, err
	}

	return group, nil
}

func Create(params files_sdk.GroupCreateParams) (files_sdk.Group, error) {
	return (&Client{}).Create(params)
}

func (c *Client) Update(params files_sdk.GroupUpdateParams) (files_sdk.Group, error) {
	group := files_sdk.Group{}
	if params.Id == 0 {
		return group, lib.CreateError(params, "Id")
	}
	path := "/groups/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return group, err
	}
	data, res, err := files_sdk.Call("PATCH", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return group, err
	}
	if res.StatusCode == 204 {
		return group, nil
	}
	if err := group.UnmarshalJSON(*data); err != nil {
		return group, err
	}

	return group, nil
}

func Update(params files_sdk.GroupUpdateParams) (files_sdk.Group, error) {
	return (&Client{}).Update(params)
}

func (c *Client) Delete(params files_sdk.GroupDeleteParams) (files_sdk.Group, error) {
	group := files_sdk.Group{}
	if params.Id == 0 {
		return group, lib.CreateError(params, "Id")
	}
	path := "/groups/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return group, err
	}
	data, res, err := files_sdk.Call("DELETE", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return group, err
	}
	if res.StatusCode == 204 {
		return group, nil
	}
	if err := group.UnmarshalJSON(*data); err != nil {
		return group, err
	}

	return group, nil
}

func Delete(params files_sdk.GroupDeleteParams) (files_sdk.Group, error) {
	return (&Client{}).Delete(params)
}
