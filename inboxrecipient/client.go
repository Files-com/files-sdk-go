package inbox_recipient

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

func (i *Iter) InboxRecipient() files_sdk.InboxRecipient {
	return i.Current().(files_sdk.InboxRecipient)
}

func (c *Client) List(ctx context.Context, params files_sdk.InboxRecipientListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/inbox_recipients"
	i.ListParams = &params
	list := files_sdk.InboxRecipientCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.InboxRecipientListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.InboxRecipientCreateParams) (files_sdk.InboxRecipient, error) {
	inboxRecipient := files_sdk.InboxRecipient{}
	path := "/inbox_recipients"
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
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

func Create(ctx context.Context, params files_sdk.InboxRecipientCreateParams) (files_sdk.InboxRecipient, error) {
	return (&Client{}).Create(ctx, params)
}
