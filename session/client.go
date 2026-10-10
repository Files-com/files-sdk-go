// Package session provides the Files.com Session API client.
package session

import (
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// Client calls the Session API using its embedded Config.
type Client struct {
	files_sdk.Config
}

// Create calls POST /sessions.
//
// API operation: Create user session (log in).
func (c *Client) Create(params files_sdk.SessionCreateParams, opts ...files_sdk.RequestResponseOption) (session files_sdk.Session, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/sessions", Params: params, Entity: &session}, opts...)
	return
}

// Create calls Client.Create using the default configuration.
func Create(params files_sdk.SessionCreateParams, opts ...files_sdk.RequestResponseOption) (session files_sdk.Session, err error) {
	return (&Client{}).Create(params, opts...)
}

// Delete calls DELETE /sessions.
//
// API operation: Delete user session (log out).
func (c *Client) Delete(opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "DELETE", Path: "/sessions", Entity: nil}, opts...)
	return
}

// Delete calls Client.Delete using the default configuration.
func Delete(opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(opts...)
}
