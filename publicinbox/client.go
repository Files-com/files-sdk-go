package public_inbox

import (
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

func (i *Iter) PublicInbox() files_sdk.PublicInbox {
	return i.Current().(files_sdk.PublicInbox)
}

func (c *Client) List(params files_sdk.PublicInboxListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/public_inboxes", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.PublicInboxCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

func List(params files_sdk.PublicInboxListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(params, opts...)
}

func (c *Client) GetKey(params files_sdk.PublicInboxGetKeyParams, opts ...files_sdk.RequestResponseOption) (publicInbox files_sdk.PublicInbox, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/public_inboxes/key/{key}", Params: params, Entity: &publicInbox}, opts...)
	return
}

func GetKey(params files_sdk.PublicInboxGetKeyParams, opts ...files_sdk.RequestResponseOption) (publicInbox files_sdk.PublicInbox, err error) {
	return (&Client{}).GetKey(params, opts...)
}
