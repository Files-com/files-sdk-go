package user

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

func (i *Iter) User() files_sdk.User {
	return i.Current().(files_sdk.User)
}

func (c *Client) List(ctx context.Context, params files_sdk.UserListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/users"
	i.ListParams = &params
	list := files_sdk.UserCollection{}
	i.Query = listquery.Build(ctx, i, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.UserListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}

func (c *Client) Find(ctx context.Context, params files_sdk.UserFindParams) (files_sdk.User, error) {
	user := files_sdk.User{}
	if params.Id == 0 {
		return user, lib.CreateError(params, "Id")
	}
	path := "/users/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return user, err
	}
	data, res, err := files_sdk.Call(ctx, "GET", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return user, err
	}
	if res.StatusCode == 204 {
		return user, nil
	}
	if err := user.UnmarshalJSON(*data); err != nil {
		return user, err
	}

	return user, nil
}

func Find(ctx context.Context, params files_sdk.UserFindParams) (files_sdk.User, error) {
	return (&Client{}).Find(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.UserCreateParams) (files_sdk.User, error) {
	user := files_sdk.User{}
	path := "/users"
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return user, err
	}
	data, res, err := files_sdk.Call(ctx, "POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return user, err
	}
	if res.StatusCode == 204 {
		return user, nil
	}
	if err := user.UnmarshalJSON(*data); err != nil {
		return user, err
	}

	return user, nil
}

func Create(ctx context.Context, params files_sdk.UserCreateParams) (files_sdk.User, error) {
	return (&Client{}).Create(ctx, params)
}

func (c *Client) Unlock(ctx context.Context, params files_sdk.UserUnlockParams) (files_sdk.User, error) {
	user := files_sdk.User{}
	if params.Id == 0 {
		return user, lib.CreateError(params, "Id")
	}
	path := "/users/" + strconv.FormatInt(params.Id, 10) + "/unlock"
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return user, err
	}
	data, res, err := files_sdk.Call(ctx, "POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return user, err
	}
	if res.StatusCode == 204 {
		return user, nil
	}
	if err := user.UnmarshalJSON(*data); err != nil {
		return user, err
	}

	return user, nil
}

func Unlock(ctx context.Context, params files_sdk.UserUnlockParams) (files_sdk.User, error) {
	return (&Client{}).Unlock(ctx, params)
}

func (c *Client) ResendWelcomeEmail(ctx context.Context, params files_sdk.UserResendWelcomeEmailParams) (files_sdk.User, error) {
	user := files_sdk.User{}
	if params.Id == 0 {
		return user, lib.CreateError(params, "Id")
	}
	path := "/users/" + strconv.FormatInt(params.Id, 10) + "/resend_welcome_email"
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return user, err
	}
	data, res, err := files_sdk.Call(ctx, "POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return user, err
	}
	if res.StatusCode == 204 {
		return user, nil
	}
	if err := user.UnmarshalJSON(*data); err != nil {
		return user, err
	}

	return user, nil
}

func ResendWelcomeEmail(ctx context.Context, params files_sdk.UserResendWelcomeEmailParams) (files_sdk.User, error) {
	return (&Client{}).ResendWelcomeEmail(ctx, params)
}

func (c *Client) User2faReset(ctx context.Context, params files_sdk.UserUser2faResetParams) (files_sdk.User, error) {
	user := files_sdk.User{}
	if params.Id == 0 {
		return user, lib.CreateError(params, "Id")
	}
	path := "/users/" + strconv.FormatInt(params.Id, 10) + "/2fa/reset"
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return user, err
	}
	data, res, err := files_sdk.Call(ctx, "POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return user, err
	}
	if res.StatusCode == 204 {
		return user, nil
	}
	if err := user.UnmarshalJSON(*data); err != nil {
		return user, err
	}

	return user, nil
}

func User2faReset(ctx context.Context, params files_sdk.UserUser2faResetParams) (files_sdk.User, error) {
	return (&Client{}).User2faReset(ctx, params)
}

func (c *Client) Update(ctx context.Context, params files_sdk.UserUpdateParams) (files_sdk.User, error) {
	user := files_sdk.User{}
	if params.Id == 0 {
		return user, lib.CreateError(params, "Id")
	}
	path := "/users/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return user, err
	}
	data, res, err := files_sdk.Call(ctx, "PATCH", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return user, err
	}
	if res.StatusCode == 204 {
		return user, nil
	}
	if err := user.UnmarshalJSON(*data); err != nil {
		return user, err
	}

	return user, nil
}

func Update(ctx context.Context, params files_sdk.UserUpdateParams) (files_sdk.User, error) {
	return (&Client{}).Update(ctx, params)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.UserDeleteParams) (files_sdk.User, error) {
	user := files_sdk.User{}
	if params.Id == 0 {
		return user, lib.CreateError(params, "Id")
	}
	path := "/users/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return user, err
	}
	data, res, err := files_sdk.Call(ctx, "DELETE", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return user, err
	}
	if res.StatusCode == 204 {
		return user, nil
	}
	if err := user.UnmarshalJSON(*data); err != nil {
		return user, err
	}

	return user, nil
}

func Delete(ctx context.Context, params files_sdk.UserDeleteParams) (files_sdk.User, error) {
	return (&Client{}).Delete(ctx, params)
}
