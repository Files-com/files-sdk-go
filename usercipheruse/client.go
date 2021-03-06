package user_cipher_use

import (
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

func (i *Iter) UserCipherUse() files_sdk.UserCipherUse {
	return i.Current().(files_sdk.UserCipherUse)
}

func (c *Client) List(params files_sdk.UserCipherUseListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/user_cipher_uses"
	i.ListParams = &params
	list := files_sdk.UserCipherUseCollection{}
	i.Query = listquery.Build(i, c.Config, path, &list)
	return i, nil
}

func List(params files_sdk.UserCipherUseListParams) (*Iter, error) {
	return (&Client{}).List(params)
}
