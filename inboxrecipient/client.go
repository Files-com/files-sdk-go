package inbox_recipient

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

func (i *Iter) InboxRecipient() files_sdk.InboxRecipient {
	return i.Current().(files_sdk.InboxRecipient)
}

func (c *Client) List(params files_sdk.InboxRecipientListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/inbox_recipients"
	i.ListParams = &params
	list := files_sdk.InboxRecipientCollection{}
	i.Query = listquery.Build(i, c.Config, path, &list)
	return i, nil
}

func List(params files_sdk.InboxRecipientListParams) (*Iter, error) {
	return (&Client{}).List(params)
}

func (c *Client) Create(params files_sdk.InboxRecipientCreateParams) (files_sdk.InboxRecipient, error) {
	inboxRecipient := files_sdk.InboxRecipient{}
	path := "/inbox_recipients"
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return inboxRecipient, err
	}
	data, res, err := files_sdk.Call("POST", c.Config, path, exportedParams)
	if err != nil {
		return inboxRecipient, err
	}
	if res.StatusCode == 204 {
		return inboxRecipient, nil
	}
	if err := inboxRecipient.UnmarshalJSON(*data); err != nil {
		return inboxRecipient, err
	}

	return inboxRecipient, nil
}

func Create(params files_sdk.InboxRecipientCreateParams) (files_sdk.InboxRecipient, error) {
	return (&Client{}).Create(params)
}
