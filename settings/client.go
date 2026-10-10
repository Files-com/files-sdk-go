// Package settings retrieves Files.com settings for desktop clients.
package settings

import (
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// Client retrieves settings using its embedded Config.
type Client struct {
	files_sdk.Config
}

// Get retrieves the settings returned by GET /settings.
func (c Client) Get(opts ...files_sdk.RequestResponseOption) (files_sdk.Settings, error) {
	settings := files_sdk.Settings{}
	err := files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/settings", Entity: &settings}, opts...)
	return settings, err
}

// Get retrieves settings using the default configuration.
func Get(opts ...files_sdk.RequestResponseOption) (files_sdk.Settings, error) {
	return (&Client{}).Get(opts...)
}
