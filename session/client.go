package session

import (
	"context"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Create(ctx context.Context, params files_sdk.SessionCreateParams) (files_sdk.Session, error) {
	session := files_sdk.Session{}
	path := "/sessions"
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return session, err
	}
	if res.StatusCode == 204 {
		return session, nil
	}
	if err := session.UnmarshalJSON(*data); err != nil {
		return session, err
	}

	return session, nil
}

func Create(ctx context.Context, params files_sdk.SessionCreateParams) (files_sdk.Session, error) {
	return (&Client{}).Create(ctx, params)
}

func (c *Client) Delete(ctx context.Context) (files_sdk.Session, error) {
	session := files_sdk.Session{}
	path := "/sessions"
	exportedParams := lib.Params{Params: lib.Interface()}
	data, res, err := files_sdk.Call(ctx, "DELETE", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return session, err
	}
	if res.StatusCode == 204 {
		return session, nil
	}
	if err := session.UnmarshalJSON(*data); err != nil {
		return session, err
	}

	return session, nil
}

func Delete(ctx context.Context) (files_sdk.Session, error) {
	return (&Client{}).Delete(ctx)
}
