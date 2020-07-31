package site

import (
  lib "github.com/Files-com/files-sdk-go/lib"
  files_sdk "github.com/Files-com/files-sdk-go"
)

type Client struct {
	files_sdk.Config
}


func (c *Client) Get () (files_sdk.Site, error) {
  site := files_sdk.Site{}
	  path := "/site"
	data, _, err := files_sdk.Call("GET", c.Config, path, lib.ExportParams(lib.Interface()))
	if err != nil {
	  return site, err
	}
	if err := site.UnmarshalJSON(*data); err != nil {
	return site, err
	}

	return  site, nil
}

func Get () (files_sdk.Site, error) {
  client := Client{}
  return client.Get ()
}

func (c *Client) GetUsage () (files_sdk.Site, error) {
  site := files_sdk.Site{}
	  path := "/site/usage"
	data, _, err := files_sdk.Call("GET", c.Config, path, lib.ExportParams(lib.Interface()))
	if err != nil {
	  return site, err
	}
	if err := site.UnmarshalJSON(*data); err != nil {
	return site, err
	}

	return  site, nil
}

func GetUsage () (files_sdk.Site, error) {
  client := Client{}
  return client.GetUsage ()
}

func (c *Client) Update (params files_sdk.SiteUpdateParams) (files_sdk.Site, error) {
  site := files_sdk.Site{}
	  path := "/site"
	data, _, err := files_sdk.Call("PATCH", c.Config, path, lib.ExportParams(params))
	if err != nil {
	  return site, err
	}
	if err := site.UnmarshalJSON(*data); err != nil {
	return site, err
	}

	return  site, nil
}

func Update (params files_sdk.SiteUpdateParams) (files_sdk.Site, error) {
  client := Client{}
  return client.Update (params)
}
