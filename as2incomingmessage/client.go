package as2_incoming_message

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

func (i *Iter) As2IncomingMessage() files_sdk.As2IncomingMessage {
	return i.Current().(files_sdk.As2IncomingMessage)
}

func (c *Client) List(ctx context.Context, params files_sdk.As2IncomingMessageListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	path, err := lib.BuildPath("/as2_incoming_messages", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.As2IncomingMessageCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.As2IncomingMessageListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}
