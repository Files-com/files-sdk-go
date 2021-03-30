package inbox_recipient

import (
	files_sdk "github.com/Files-com/files-sdk-go"
	lib "github.com/Files-com/files-sdk-go/lib"
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
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	i := &Iter{Iter: &lib.Iter{}}
	path := "/inbox_recipients"
	i.ListParams = &params
	exportParams, err := i.ExportParams()
	if err != nil {
		return i, err
	}
	i.Query = func() (*[]interface{}, string, error) {
		data, res, err := files_sdk.Call("GET", c.Config, path, exportParams)
		defaultValue := make([]interface{}, 0)
		if err != nil {
			return &defaultValue, "", err
		}
		list := files_sdk.InboxRecipientCollection{}
		if err := list.UnmarshalJSON(*data); err != nil {
			return &defaultValue, "", err
		}

		ret := make([]interface{}, len(list))
		for i, v := range list {
			ret[i] = v
		}
		cursor := res.Header.Get("X-Files-Cursor")
		return &ret, cursor, nil
	}
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
