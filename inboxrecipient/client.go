package inbox_recipient

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

func (i *Iter) InboxRecipient() files_sdk.InboxRecipient {
	return i.Current().(files_sdk.InboxRecipient)
}

func (c *Client) List(params files_sdk.InboxRecipientListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/inbox_recipients", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.InboxRecipientCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

func List(params files_sdk.InboxRecipientListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(params, opts...)
}

func (c *Client) Create(params files_sdk.InboxRecipientCreateParams, opts ...files_sdk.RequestResponseOption) (inboxRecipient files_sdk.InboxRecipient, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/inbox_recipients", Params: params, Entity: &inboxRecipient}, opts...)
	return
}

func Create(params files_sdk.InboxRecipientCreateParams, opts ...files_sdk.RequestResponseOption) (inboxRecipient files_sdk.InboxRecipient, err error) {
	return (&Client{}).Create(params, opts...)
}
