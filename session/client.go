package session

import (
	files_sdk "github.com/Files-com/files-sdk-go/v2"
	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Create(params files_sdk.SessionCreateParams, opts ...files_sdk.RequestResponseOption) (session files_sdk.Session, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/sessions", Params: params, Entity: &session}, opts...)
	return
}

func Create(params files_sdk.SessionCreateParams, opts ...files_sdk.RequestResponseOption) (session files_sdk.Session, err error) {
	return (&Client{}).Create(params, opts...)
}

func (c *Client) ForgotReset(params files_sdk.SessionForgotResetParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/sessions/forgot/reset", Params: params, Entity: nil}, opts...)
	return
}

func ForgotReset(params files_sdk.SessionForgotResetParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).ForgotReset(params, opts...)
}

func (c *Client) ForgotValidate(params files_sdk.SessionForgotValidateParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/sessions/forgot/validate", Params: params, Entity: nil}, opts...)
	return
}

func ForgotValidate(params files_sdk.SessionForgotValidateParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).ForgotValidate(params, opts...)
}

func (c *Client) Forgot(params files_sdk.SessionForgotParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/sessions/forgot", Params: params, Entity: nil}, opts...)
	return
}

func Forgot(params files_sdk.SessionForgotParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Forgot(params, opts...)
}

func (c *Client) PairingKey(params files_sdk.SessionPairingKeyParams, opts ...files_sdk.RequestResponseOption) (pairedApiKey files_sdk.PairedApiKey, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/sessions/pairing_key/{key}", Params: params, Entity: &pairedApiKey}, opts...)
	return
}

func PairingKey(params files_sdk.SessionPairingKeyParams, opts ...files_sdk.RequestResponseOption) (pairedApiKey files_sdk.PairedApiKey, err error) {
	return (&Client{}).PairingKey(params, opts...)
}

func (c *Client) Oauth(params files_sdk.SessionOauthParams, opts ...files_sdk.RequestResponseOption) (oauthRedirect files_sdk.OauthRedirect, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/sessions/oauth", Params: params, Entity: &oauthRedirect}, opts...)
	return
}

func Oauth(params files_sdk.SessionOauthParams, opts ...files_sdk.RequestResponseOption) (oauthRedirect files_sdk.OauthRedirect, err error) {
	return (&Client{}).Oauth(params, opts...)
}

func (c *Client) Delete(opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "DELETE", Path: "/sessions", Params: lib.Interface(), Entity: nil}, opts...)
	return
}

func Delete(opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(opts...)
}
