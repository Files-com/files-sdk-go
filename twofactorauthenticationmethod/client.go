package two_factor_authentication_method

import (
	files_sdk "github.com/Files-com/files-sdk-go/v2"
	lib "github.com/Files-com/files-sdk-go/v2/lib"
	listquery "github.com/Files-com/files-sdk-go/v2/listquery"
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

func (i *Iter) TwoFactorAuthenticationMethod() files_sdk.TwoFactorAuthenticationMethod {
	return i.Current().(files_sdk.TwoFactorAuthenticationMethod)
}

func (c *Client) Get(params files_sdk.TwoFactorAuthenticationMethodGetParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/2fa", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.TwoFactorAuthenticationMethodCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

func Get(params files_sdk.TwoFactorAuthenticationMethodGetParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).Get(params, opts...)
}

func (c *Client) Create(params files_sdk.TwoFactorAuthenticationMethodCreateParams, opts ...files_sdk.RequestResponseOption) (twoFactorAuthenticationMethod files_sdk.TwoFactorAuthenticationMethod, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/2fa", Params: params, Entity: &twoFactorAuthenticationMethod}, opts...)
	return
}

func Create(params files_sdk.TwoFactorAuthenticationMethodCreateParams, opts ...files_sdk.RequestResponseOption) (twoFactorAuthenticationMethod files_sdk.TwoFactorAuthenticationMethod, err error) {
	return (&Client{}).Create(params, opts...)
}

func (c *Client) SendCode(params files_sdk.TwoFactorAuthenticationMethodSendCodeParams, opts ...files_sdk.RequestResponseOption) (u2fSignRequest files_sdk.U2fSignRequest, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/2fa/send_code", Params: params, Entity: &u2fSignRequest}, opts...)
	return
}

func SendCode(params files_sdk.TwoFactorAuthenticationMethodSendCodeParams, opts ...files_sdk.RequestResponseOption) (u2fSignRequest files_sdk.U2fSignRequest, err error) {
	return (&Client{}).SendCode(params, opts...)
}

func (c *Client) Update(params files_sdk.TwoFactorAuthenticationMethodUpdateParams, opts ...files_sdk.RequestResponseOption) (twoFactorAuthenticationMethod files_sdk.TwoFactorAuthenticationMethod, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/2fa/{id}", Params: params, Entity: &twoFactorAuthenticationMethod}, opts...)
	return
}

func Update(params files_sdk.TwoFactorAuthenticationMethodUpdateParams, opts ...files_sdk.RequestResponseOption) (twoFactorAuthenticationMethod files_sdk.TwoFactorAuthenticationMethod, err error) {
	return (&Client{}).Update(params, opts...)
}

func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (twoFactorAuthenticationMethod files_sdk.TwoFactorAuthenticationMethod, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/2fa/{id}", Params: params, Entity: &twoFactorAuthenticationMethod}, opts...)
	return
}

func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (twoFactorAuthenticationMethod files_sdk.TwoFactorAuthenticationMethod, err error) {
	return (&Client{}).UpdateWithMap(params, opts...)
}

func (c *Client) Delete(params files_sdk.TwoFactorAuthenticationMethodDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "DELETE", Path: "/2fa/{id}", Params: params, Entity: nil}, opts...)
	return
}

func Delete(params files_sdk.TwoFactorAuthenticationMethodDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(params, opts...)
}
