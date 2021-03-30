package bundle_recipient

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

func (i *Iter) BundleRecipient() files_sdk.BundleRecipient {
	return i.Current().(files_sdk.BundleRecipient)
}

func (c *Client) List(params files_sdk.BundleRecipientListParams) (*Iter, error) {
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	i := &Iter{Iter: &lib.Iter{}}
	path := "/bundle_recipients"
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
		list := files_sdk.BundleRecipientCollection{}
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

func List(params files_sdk.BundleRecipientListParams) (*Iter, error) {
	return (&Client{}).List(params)
}

func (c *Client) Create(params files_sdk.BundleRecipientCreateParams) (files_sdk.BundleRecipient, error) {
	bundleRecipient := files_sdk.BundleRecipient{}
	path := "/bundle_recipients"
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return bundleRecipient, err
	}
	data, res, err := files_sdk.Call("POST", c.Config, path, exportedParams)
	if err != nil {
		return bundleRecipient, err
	}
	if res.StatusCode == 204 {
		return bundleRecipient, nil
	}
	if err := bundleRecipient.UnmarshalJSON(*data); err != nil {
		return bundleRecipient, err
	}

	return bundleRecipient, nil
}

func Create(params files_sdk.BundleRecipientCreateParams) (files_sdk.BundleRecipient, error) {
	return (&Client{}).Create(params)
}
