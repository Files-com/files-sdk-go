package user

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
	*files_sdk.Iter
	*Client
}

func (i *Iter) Reload(opts ...files_sdk.RequestResponseOption) files_sdk.IterI {
	return &Iter{Iter: i.Iter.Reload(opts...).(*files_sdk.Iter), Client: i.Client}
}

func (i *Iter) User() files_sdk.User {
	return i.Current().(files_sdk.User)
}

func (i *Iter) LoadResource(identifier interface{}, opts ...files_sdk.RequestResponseOption) (interface{}, error) {
	params := files_sdk.UserFindParams{}
	if id, ok := identifier.(int64); ok {
		params.Id = id
	}
	return i.Client.Find(context.Background(), params, opts...)
}

func (c *Client) List(ctx context.Context, params files_sdk.UserListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/users", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.UserCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list, opts...)
	return i, nil
}

func List(ctx context.Context, params files_sdk.UserListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(ctx, params, opts...)
}

func (c *Client) Find(ctx context.Context, params files_sdk.UserFindParams, opts ...files_sdk.RequestResponseOption) (user files_sdk.User, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/users/{id}", Params: params, Entity: &user}, opts...)
	return
}

func Find(ctx context.Context, params files_sdk.UserFindParams, opts ...files_sdk.RequestResponseOption) (user files_sdk.User, err error) {
	return (&Client{}).Find(ctx, params, opts...)
}

func (c *Client) Create(ctx context.Context, params files_sdk.UserCreateParams, opts ...files_sdk.RequestResponseOption) (user files_sdk.User, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/users", Params: params, Entity: &user}, opts...)
	return
}

func Create(ctx context.Context, params files_sdk.UserCreateParams, opts ...files_sdk.RequestResponseOption) (user files_sdk.User, err error) {
	return (&Client{}).Create(ctx, params, opts...)
}

func (c *Client) Unlock(ctx context.Context, params files_sdk.UserUnlockParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/users/{id}/unlock", Params: params, Entity: nil}, opts...)
	return
}

func Unlock(ctx context.Context, params files_sdk.UserUnlockParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Unlock(ctx, params, opts...)
}

func (c *Client) ResendWelcomeEmail(ctx context.Context, params files_sdk.UserResendWelcomeEmailParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/users/{id}/resend_welcome_email", Params: params, Entity: nil}, opts...)
	return
}

func ResendWelcomeEmail(ctx context.Context, params files_sdk.UserResendWelcomeEmailParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).ResendWelcomeEmail(ctx, params, opts...)
}

func (c *Client) User2faReset(ctx context.Context, params files_sdk.UserUser2faResetParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/users/{id}/2fa/reset", Params: params, Entity: nil}, opts...)
	return
}

func User2faReset(ctx context.Context, params files_sdk.UserUser2faResetParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).User2faReset(ctx, params, opts...)
}

func (c *Client) Update(ctx context.Context, params files_sdk.UserUpdateParams, opts ...files_sdk.RequestResponseOption) (user files_sdk.User, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "PATCH", Path: "/users/{id}", Params: params, Entity: &user}, opts...)
	return
}

func Update(ctx context.Context, params files_sdk.UserUpdateParams, opts ...files_sdk.RequestResponseOption) (user files_sdk.User, err error) {
	return (&Client{}).Update(ctx, params, opts...)
}

func (c *Client) UpdateWithMap(ctx context.Context, params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (user files_sdk.User, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "PATCH", Path: "/users/{id}", Params: params, Entity: &user}, opts...)
	return
}

func UpdateWithMap(ctx context.Context, params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (user files_sdk.User, err error) {
	return (&Client{}).UpdateWithMap(ctx, params, opts...)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.UserDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "DELETE", Path: "/users/{id}", Params: params, Entity: nil}, opts...)
	return
}

func Delete(ctx context.Context, params files_sdk.UserDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(ctx, params, opts...)
}
