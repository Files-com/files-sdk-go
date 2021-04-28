package sso_strategy

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

func (i *Iter) SsoStrategy() files_sdk.SsoStrategy {
	return i.Current().(files_sdk.SsoStrategy)
}

func (c *Client) List(params files_sdk.SsoStrategyListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/sso_strategies"
	i.ListParams = &params
	list := files_sdk.SsoStrategyCollection{}
	i.Query = listquery.Build(i, c.Config, path, &list)
	return i, nil
}

func List(params files_sdk.SsoStrategyListParams) (*Iter, error) {
	return (&Client{}).List(params)
}

func (c *Client) Find(params files_sdk.SsoStrategyFindParams) (files_sdk.SsoStrategy, error) {
	ssoStrategy := files_sdk.SsoStrategy{}
	if params.Id == 0 {
		return ssoStrategy, lib.CreateError(params, "Id")
	}
	path := "/sso_strategies/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return ssoStrategy, err
	}
	data, res, err := files_sdk.Call("GET", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return ssoStrategy, err
	}
	if res.StatusCode == 204 {
		return ssoStrategy, nil
	}
	if err := ssoStrategy.UnmarshalJSON(*data); err != nil {
		return ssoStrategy, err
	}

	return ssoStrategy, nil
}

func Find(params files_sdk.SsoStrategyFindParams) (files_sdk.SsoStrategy, error) {
	return (&Client{}).Find(params)
}
