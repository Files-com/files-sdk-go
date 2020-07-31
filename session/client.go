package session

import (
  lib "github.com/Files-com/files-sdk-go/lib"
  files_sdk "github.com/Files-com/files-sdk-go"
)

type Client struct {
	files_sdk.Config
}


func (c *Client) Create (params files_sdk.SessionCreateParams) (files_sdk.Session, error) {
  session := files_sdk.Session{}
	  path := "/sessions"
	data, _, err := files_sdk.Call("POST", c.Config, path, lib.ExportParams(params))
	if err != nil {
	  return session, err
	}
	if err := session.UnmarshalJSON(*data); err != nil {
	return session, err
	}

	return  session, nil
}

func Create (params files_sdk.SessionCreateParams) (files_sdk.Session, error) {
  client := Client{}
  return client.Create (params)
}

func (c *Client) Delete () (files_sdk.Session, error) {
  session := files_sdk.Session{}
	  path := "/sessions"
	data, _, err := files_sdk.Call("DELETE", c.Config, path, lib.ExportParams(lib.Interface()))
	if err != nil {
	  return session, err
	}
	if err := session.UnmarshalJSON(*data); err != nil {
	return session, err
	}

	return  session, nil
}

func Delete () (files_sdk.Session, error) {
  client := Client{}
  return client.Delete ()
}
