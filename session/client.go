package session

import (
	"context"

	files_sdk "github.com/Files-com/files-sdk-go"
	lib "github.com/Files-com/files-sdk-go/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Create(ctx context.Context, params files_sdk.SessionCreateParams) (files_sdk.Session, error) {
	session := files_sdk.Session{}
	path := "/sessions"
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return session, err
	}
	data, res, err := files_sdk.Call(ctx, "POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
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
	exportedParams, err := lib.ExportParams(lib.Interface())
	if err != nil {
		return session, err
	}
	data, res, err := files_sdk.Call(ctx, "DELETE", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
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
