package user_request

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

func (i *Iter) UserRequest() files_sdk.UserRequest {
	return i.Current().(files_sdk.UserRequest)
}

func (c *Client) List(params files_sdk.UserRequestListParams) (*Iter, error) {
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	i := &Iter{Iter: &lib.Iter{}}
	path := "/user_requests"
	i.ListParams = &params
	exportParams, err := i.ExportParams()
	if err != nil {
		return i, err
	}
	i.Query = func() (*[]interface{}, string, error) {
		data, res, err := files_sdk.Call("GET", c.Config, path, exportParams)
		defaultValue := make([]interface{}, 0)
		if err != nil {
			return &defaultValue, "", err
		}
		list := files_sdk.UserRequestCollection{}
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
	return i, nil
}

func List(params files_sdk.UserRequestListParams) (*Iter, error) {
	return (&Client{}).List(params)
}

func (c *Client) Find(params files_sdk.UserRequestFindParams) (files_sdk.UserRequest, error) {
	userRequest := files_sdk.UserRequest{}
	if params.Id == 0 {
		return userRequest, lib.CreateError(params, "Id")
	}
	path := "/user_requests/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return userRequest, err
	}
	data, res, err := files_sdk.Call("GET", c.Config, path, exportedParams)
	if err != nil {
		return userRequest, err
	}
	if res.StatusCode == 204 {
		return userRequest, nil
	}
	if err := userRequest.UnmarshalJSON(*data); err != nil {
		return userRequest, err
	}

	return userRequest, nil
}

func Find(params files_sdk.UserRequestFindParams) (files_sdk.UserRequest, error) {
	return (&Client{}).Find(params)
}

func (c *Client) Create(params files_sdk.UserRequestCreateParams) (files_sdk.UserRequest, error) {
	userRequest := files_sdk.UserRequest{}
	path := "/user_requests"
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return userRequest, err
	}
	data, res, err := files_sdk.Call("POST", c.Config, path, exportedParams)
	if err != nil {
		return userRequest, err
	}
	if res.StatusCode == 204 {
		return userRequest, nil
	}
	if err := userRequest.UnmarshalJSON(*data); err != nil {
		return userRequest, err
	}

	return userRequest, nil
}

func Create(params files_sdk.UserRequestCreateParams) (files_sdk.UserRequest, error) {
	return (&Client{}).Create(params)
}

func (c *Client) Delete(params files_sdk.UserRequestDeleteParams) (files_sdk.UserRequest, error) {
	userRequest := files_sdk.UserRequest{}
	if params.Id == 0 {
		return userRequest, lib.CreateError(params, "Id")
	}
	path := "/user_requests/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return userRequest, err
	}
	data, res, err := files_sdk.Call("DELETE", c.Config, path, exportedParams)
	if err != nil {
		return userRequest, err
	}
	if res.StatusCode == 204 {
		return userRequest, nil
	}
	if err := userRequest.UnmarshalJSON(*data); err != nil {
		return userRequest, err
	}

	return userRequest, nil
}

func Delete(params files_sdk.UserRequestDeleteParams) (files_sdk.UserRequest, error) {
	return (&Client{}).Delete(params)
}
