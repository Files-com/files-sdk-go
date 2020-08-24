package group

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

func (i *Iter) Group() files_sdk.Group {
	return i.Current().(files_sdk.Group)
}

func (c *Client) List(params files_sdk.GroupListParams) *Iter {
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	i := &Iter{Iter: &lib.Iter{}}
	path := "/groups"

	i.Query = func() (*[]interface{}, string, error) {
		data, res, err := files_sdk.Call("GET", c.Config, path, i.ExportParams())
		defaultValue := make([]interface{}, 0)
		if err != nil {
			return &defaultValue, "", err
		}
		list := files_sdk.GroupCollection{}
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

func List(params files_sdk.GroupListParams) *Iter {
	return (&Client{}).List(params)
}

func (c *Client) Find(params files_sdk.GroupFindParams) (files_sdk.Group, error) {
	group := files_sdk.Group{}
	path := "/groups/" + lib.QueryEscape(strconv.FormatInt(params.Id, 10)) + ""
	data, res, err := files_sdk.Call("GET", c.Config, path, lib.ExportParams(params))
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
	data, res, err := files_sdk.Call("POST", c.Config, path, lib.ExportParams(params))
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
	path := "/groups/" + lib.QueryEscape(strconv.FormatInt(params.Id, 10)) + ""
	data, res, err := files_sdk.Call("PATCH", c.Config, path, lib.ExportParams(params))
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
	path := "/groups/" + lib.QueryEscape(strconv.FormatInt(params.Id, 10)) + ""
	data, res, err := files_sdk.Call("DELETE", c.Config, path, lib.ExportParams(params))
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
