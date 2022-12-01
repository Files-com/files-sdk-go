package form_field_set

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

func (i *Iter) FormFieldSet() files_sdk.FormFieldSet {
	return i.Current().(files_sdk.FormFieldSet)
}

func (c *Client) List(ctx context.Context, params files_sdk.FormFieldSetListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	path, err := lib.BuildPath("/form_field_sets", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.FormFieldSetCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list, opts...)
	return i, nil
}

func List(ctx context.Context, params files_sdk.FormFieldSetListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(ctx, params, opts...)
}

func (c *Client) Find(ctx context.Context, params files_sdk.FormFieldSetFindParams, opts ...files_sdk.RequestResponseOption) (formFieldSet files_sdk.FormFieldSet, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/form_field_sets/{id}", Params: params, Entity: &formFieldSet}, opts...)
	return
}

func Find(ctx context.Context, params files_sdk.FormFieldSetFindParams, opts ...files_sdk.RequestResponseOption) (formFieldSet files_sdk.FormFieldSet, err error) {
	return (&Client{}).Find(ctx, params, opts...)
}

func (c *Client) Create(ctx context.Context, params files_sdk.FormFieldSetCreateParams, opts ...files_sdk.RequestResponseOption) (formFieldSet files_sdk.FormFieldSet, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/form_field_sets", Params: params, Entity: &formFieldSet}, opts...)
	return
}

func Create(ctx context.Context, params files_sdk.FormFieldSetCreateParams, opts ...files_sdk.RequestResponseOption) (formFieldSet files_sdk.FormFieldSet, err error) {
	return (&Client{}).Create(ctx, params, opts...)
}

func (c *Client) Update(ctx context.Context, params files_sdk.FormFieldSetUpdateParams, opts ...files_sdk.RequestResponseOption) (formFieldSet files_sdk.FormFieldSet, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "PATCH", Path: "/form_field_sets/{id}", Params: params, Entity: &formFieldSet}, opts...)
	return
}

func Update(ctx context.Context, params files_sdk.FormFieldSetUpdateParams, opts ...files_sdk.RequestResponseOption) (formFieldSet files_sdk.FormFieldSet, err error) {
	return (&Client{}).Update(ctx, params, opts...)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.FormFieldSetDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "DELETE", Path: "/form_field_sets/{id}", Params: params, Entity: nil}, opts...)
	return
}

func Delete(ctx context.Context, params files_sdk.FormFieldSetDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(ctx, params, opts...)
}
