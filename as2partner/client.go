package as2_partner

import (
	"context"
	"strconv"

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
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/as2_partners"
	i.ListParams = &params
	list := files_sdk.As2PartnerCollection{}
	i.Query = listquery.Build(ctx, i, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.As2PartnerListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}

func (c *Client) Find(ctx context.Context, params files_sdk.As2PartnerFindParams) (files_sdk.As2Partner, error) {
	as2Partner := files_sdk.As2Partner{}
	if params.Id == 0 {
		return as2Partner, lib.CreateError(params, "Id")
	}
	path := "/as2_partners/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return as2Partner, err
	}
	data, res, err := files_sdk.Call(ctx, "GET", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return as2Partner, err
	}
	if res.StatusCode == 204 {
		return as2Partner, nil
	}
	if err := as2Partner.UnmarshalJSON(*data); err != nil {
		return as2Partner, err
	}

	return as2Partner, nil
}

func Find(ctx context.Context, params files_sdk.As2PartnerFindParams) (files_sdk.As2Partner, error) {
	return (&Client{}).Find(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.As2PartnerCreateParams) (files_sdk.As2Partner, error) {
	as2Partner := files_sdk.As2Partner{}
	path := "/as2_partners"
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return as2Partner, err
	}
	data, res, err := files_sdk.Call(ctx, "POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return as2Partner, err
	}
	if res.StatusCode == 204 {
		return as2Partner, nil
	}
	if err := as2Partner.UnmarshalJSON(*data); err != nil {
		return as2Partner, err
	}

	return as2Partner, nil
}

func Create(ctx context.Context, params files_sdk.As2PartnerCreateParams) (files_sdk.As2Partner, error) {
	return (&Client{}).Create(ctx, params)
}

func (c *Client) Update(ctx context.Context, params files_sdk.As2PartnerUpdateParams) (files_sdk.As2Partner, error) {
	as2Partner := files_sdk.As2Partner{}
	if params.Id == 0 {
		return as2Partner, lib.CreateError(params, "Id")
	}
	path := "/as2_partners/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return as2Partner, err
	}
	data, res, err := files_sdk.Call(ctx, "PATCH", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return as2Partner, err
	}
	if res.StatusCode == 204 {
		return as2Partner, nil
	}
	if err := as2Partner.UnmarshalJSON(*data); err != nil {
		return as2Partner, err
	}

	return as2Partner, nil
}

func Update(ctx context.Context, params files_sdk.As2PartnerUpdateParams) (files_sdk.As2Partner, error) {
	return (&Client{}).Update(ctx, params)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.As2PartnerDeleteParams) (files_sdk.As2Partner, error) {
	as2Partner := files_sdk.As2Partner{}
	if params.Id == 0 {
		return as2Partner, lib.CreateError(params, "Id")
	}
	path := "/as2_partners/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return as2Partner, err
	}
	data, res, err := files_sdk.Call(ctx, "DELETE", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return as2Partner, err
	}
	if res.StatusCode == 204 {
		return as2Partner, nil
	}
	if err := as2Partner.UnmarshalJSON(*data); err != nil {
		return as2Partner, err
	}

	return as2Partner, nil
}

func Delete(ctx context.Context, params files_sdk.As2PartnerDeleteParams) (files_sdk.As2Partner, error) {
	return (&Client{}).Delete(ctx, params)
}
