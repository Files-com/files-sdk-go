package sftp_host_key

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
	*lib.Iter
}

func (i *Iter) SftpHostKey() files_sdk.SftpHostKey {
	return i.Current().(files_sdk.SftpHostKey)
}

func (c *Client) List(ctx context.Context, params files_sdk.SftpHostKeyListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	path, err := lib.BuildPath("/sftp_host_keys", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.SftpHostKeyCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.SftpHostKeyListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}

func (c *Client) Find(ctx context.Context, params files_sdk.SftpHostKeyFindParams) (sftpHostKey files_sdk.SftpHostKey, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/sftp_host_keys/{id}", Params: params, Entity: &sftpHostKey})
	return
}

func Find(ctx context.Context, params files_sdk.SftpHostKeyFindParams) (sftpHostKey files_sdk.SftpHostKey, err error) {
	return (&Client{}).Find(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.SftpHostKeyCreateParams) (sftpHostKey files_sdk.SftpHostKey, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/sftp_host_keys", Params: params, Entity: &sftpHostKey})
	return
}

func Create(ctx context.Context, params files_sdk.SftpHostKeyCreateParams) (sftpHostKey files_sdk.SftpHostKey, err error) {
	return (&Client{}).Create(ctx, params)
}

func (c *Client) Update(ctx context.Context, params files_sdk.SftpHostKeyUpdateParams) (sftpHostKey files_sdk.SftpHostKey, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "PATCH", Path: "/sftp_host_keys/{id}", Params: params, Entity: &sftpHostKey})
	return
}

func Update(ctx context.Context, params files_sdk.SftpHostKeyUpdateParams) (sftpHostKey files_sdk.SftpHostKey, err error) {
	return (&Client{}).Update(ctx, params)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.SftpHostKeyDeleteParams) (err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "DELETE", Path: "/sftp_host_keys/{id}", Params: params, Entity: nil})
	return
}

func Delete(ctx context.Context, params files_sdk.SftpHostKeyDeleteParams) (err error) {
	return (&Client{}).Delete(ctx, params)
}