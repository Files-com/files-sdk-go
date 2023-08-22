package inbox_recipient_registration

import (
	files_sdk "github.com/Files-com/files-sdk-go/v2"
	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Create(params files_sdk.InboxRecipientRegistrationCreateParams, opts ...files_sdk.RequestResponseOption) (inboxRecipientRegistration files_sdk.InboxRecipientRegistration, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/inbox_recipient_registrations", Params: params, Entity: &inboxRecipientRegistration}, opts...)
	return
}

func Create(params files_sdk.InboxRecipientRegistrationCreateParams, opts ...files_sdk.RequestResponseOption) (inboxRecipientRegistration files_sdk.InboxRecipientRegistration, err error) {
	return (&Client{}).Create(params, opts...)
}
