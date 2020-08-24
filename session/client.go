package session

import (
	files_sdk "github.com/Files-com/files-sdk-go"
	lib "github.com/Files-com/files-sdk-go/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Create(params files_sdk.SessionCreateParams) (files_sdk.Session, error) {
	session := files_sdk.Session{}
	path := "/sessions"
	data, res, err := files_sdk.Call("POST", c.Config, path, lib.ExportParams(params))
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

func Create(params files_sdk.SessionCreateParams) (files_sdk.Session, error) {
	return (&Client{}).Create(params)
}

func (c *Client) Delete(params files_sdk.SessionDeleteParams) (files_sdk.Session, error) {
	session := files_sdk.Session{}
	path := "/sessions"
	data, res, err := files_sdk.Call("DELETE", c.Config, path, lib.ExportParams(params))
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

func Delete(params files_sdk.SessionDeleteParams) (files_sdk.Session, error) {
	return (&Client{}).Delete(params)
}
