package user_request

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

func (i *Iter) UserRequest() files_sdk.UserRequest {
	return i.Current().(files_sdk.UserRequest)
}

func (c *Client) List(ctx context.Context, params files_sdk.UserRequestListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/user_requests"
	i.ListParams = &params
	list := files_sdk.UserRequestCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.UserRequestListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}

func (c *Client) Find(ctx context.Context, params files_sdk.UserRequestFindParams) (files_sdk.UserRequest, error) {
	userRequest := files_sdk.UserRequest{}
	if params.Id == 0 {
		return userRequest, lib.CreateError(params, "Id")
	}
	path := "/user_requests/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "GET", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return userRequest, err
	}
	if res.StatusCode == 204 {
		return userRequest, nil
	}

	return userRequest, userRequest.UnmarshalJSON(*data)
}

func Find(ctx context.Context, params files_sdk.UserRequestFindParams) (files_sdk.UserRequest, error) {
	return (&Client{}).Find(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.UserRequestCreateParams) (files_sdk.UserRequest, error) {
	userRequest := files_sdk.UserRequest{}
	path := "/user_requests"
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return userRequest, err
	}
	if res.StatusCode == 204 {
		return userRequest, nil
	}

	return userRequest, userRequest.UnmarshalJSON(*data)
}

func Create(ctx context.Context, params files_sdk.UserRequestCreateParams) (files_sdk.UserRequest, error) {
	return (&Client{}).Create(ctx, params)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.UserRequestDeleteParams) error {
	userRequest := files_sdk.UserRequest{}
	if params.Id == 0 {
		return lib.CreateError(params, "Id")
	}
	path := "/user_requests/" + strconv.FormatInt(params.Id, 10) + ""
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

	return userRequest.UnmarshalJSON(*data)
}

func Delete(ctx context.Context, params files_sdk.UserRequestDeleteParams) error {
	return (&Client{}).Delete(ctx, params)
}
