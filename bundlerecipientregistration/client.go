package bundle_recipient_registration

import (
	files_sdk "github.com/Files-com/files-sdk-go/v2"
	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Create(params files_sdk.BundleRecipientRegistrationCreateParams, opts ...files_sdk.RequestResponseOption) (bundleRecipientRegistration files_sdk.BundleRecipientRegistration, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/bundle_recipient_registrations", Params: params, Entity: &bundleRecipientRegistration}, opts...)
	return
}

func Create(params files_sdk.BundleRecipientRegistrationCreateParams, opts ...files_sdk.RequestResponseOption) (bundleRecipientRegistration files_sdk.BundleRecipientRegistration, err error) {
	return (&Client{}).Create(params, opts...)
}
