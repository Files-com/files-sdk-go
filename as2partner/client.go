package as2_partner

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

func (i *Iter) As2Partner() files_sdk.As2Partner {
	return i.Current().(files_sdk.As2Partner)
}

func (c *Client) List(ctx context.Context, params files_sdk.As2PartnerListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	path, err := lib.BuildPath("/as2_partners", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.As2PartnerCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.As2PartnerListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}

func (c *Client) Find(ctx context.Context, params files_sdk.As2PartnerFindParams) (as2Partner files_sdk.As2Partner, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/as2_partners/{id}", Params: params, Entity: &as2Partner})
	return
}

func Find(ctx context.Context, params files_sdk.As2PartnerFindParams) (as2Partner files_sdk.As2Partner, err error) {
	return (&Client{}).Find(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.As2PartnerCreateParams) (as2Partner files_sdk.As2Partner, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/as2_partners", Params: params, Entity: &as2Partner})
	return
}

func Create(ctx context.Context, params files_sdk.As2PartnerCreateParams) (as2Partner files_sdk.As2Partner, err error) {
	return (&Client{}).Create(ctx, params)
}

func (c *Client) Update(ctx context.Context, params files_sdk.As2PartnerUpdateParams) (as2Partner files_sdk.As2Partner, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "PATCH", Path: "/as2_partners/{id}", Params: params, Entity: &as2Partner})
	return
}

func Update(ctx context.Context, params files_sdk.As2PartnerUpdateParams) (as2Partner files_sdk.As2Partner, err error) {
	return (&Client{}).Update(ctx, params)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.As2PartnerDeleteParams) (err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "DELETE", Path: "/as2_partners/{id}", Params: params, Entity: nil})
	return
}

func Delete(ctx context.Context, params files_sdk.As2PartnerDeleteParams) (err error) {
	return (&Client{}).Delete(ctx, params)
}
