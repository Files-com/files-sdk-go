package group_user

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

func (i *Iter) GroupUser() files_sdk.GroupUser {
	return i.Current().(files_sdk.GroupUser)
}

func (c *Client) List(params files_sdk.GroupUserListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/group_users"
	i.ListParams = &params
	list := files_sdk.GroupUserCollection{}
	i.Query = listquery.Build(i, c.Config, path, &list)
	return i, nil
}

func List(params files_sdk.GroupUserListParams) (*Iter, error) {
	return (&Client{}).List(params)
}

func (c *Client) Create(params files_sdk.GroupUserCreateParams) (files_sdk.GroupUser, error) {
	groupUser := files_sdk.GroupUser{}
	path := "/group_users"
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return groupUser, err
	}
	data, res, err := files_sdk.Call("POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return groupUser, err
	}
	if res.StatusCode == 204 {
		return groupUser, nil
	}
	if err := groupUser.UnmarshalJSON(*data); err != nil {
		return groupUser, err
	}

	return groupUser, nil
}

func Create(params files_sdk.GroupUserCreateParams) (files_sdk.GroupUser, error) {
	return (&Client{}).Create(params)
}

func (c *Client) Update(params files_sdk.GroupUserUpdateParams) (files_sdk.GroupUser, error) {
	groupUser := files_sdk.GroupUser{}
	if params.Id == 0 {
		return groupUser, lib.CreateError(params, "Id")
	}
	path := "/group_users/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return groupUser, err
	}
	data, res, err := files_sdk.Call("PATCH", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return groupUser, err
	}
	if res.StatusCode == 204 {
		return groupUser, nil
	}
	if err := groupUser.UnmarshalJSON(*data); err != nil {
		return groupUser, err
	}

	return groupUser, nil
}

func Update(params files_sdk.GroupUserUpdateParams) (files_sdk.GroupUser, error) {
	return (&Client{}).Update(params)
}

func (c *Client) Delete(params files_sdk.GroupUserDeleteParams) (files_sdk.GroupUser, error) {
	groupUser := files_sdk.GroupUser{}
	if params.Id == 0 {
		return groupUser, lib.CreateError(params, "Id")
	}
	path := "/group_users/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return groupUser, err
	}
	data, res, err := files_sdk.Call("DELETE", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return groupUser, err
	}
	if res.StatusCode == 204 {
		return groupUser, nil
	}
	if err := groupUser.UnmarshalJSON(*data); err != nil {
		return groupUser, err
	}

	return groupUser, nil
}

func Delete(params files_sdk.GroupUserDeleteParams) (files_sdk.GroupUser, error) {
	return (&Client{}).Delete(params)
}
