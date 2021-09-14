package priority

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

func (i *Iter) Priority() files_sdk.Priority {
	return i.Current().(files_sdk.Priority)
}

func (c *Client) List(ctx context.Context, params files_sdk.PriorityListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/priorities"
	i.ListParams = &params
	list := files_sdk.PriorityCollection{}
	i.Query = listquery.Build(ctx, i, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.PriorityListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}
