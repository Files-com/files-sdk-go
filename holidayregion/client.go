package holiday_region

import (
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	lib "github.com/Files-com/files-sdk-go/v3/lib"
	listquery "github.com/Files-com/files-sdk-go/v3/listquery"
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

func (i *Iter) HolidayRegion() files_sdk.HolidayRegion {
	return i.Current().(files_sdk.HolidayRegion)
}

func (c *Client) GetSupported(params files_sdk.HolidayRegionGetSupportedParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/holiday_regions/supported", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.HolidayRegionCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

func GetSupported(params files_sdk.HolidayRegionGetSupportedParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).GetSupported(params, opts...)
}
