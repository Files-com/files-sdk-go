package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/lib"
)

type IpAddress struct {
	Id             string   `json:"id,omitempty"`
	AssociatedWith string   `json:"associated_with,omitempty"`
	GroupId        int64    `json:"group_id,omitempty"`
	IpAddresses    []string `json:"ip_addresses,omitempty"`
}

type IpAddressCollection []IpAddress

type IpAddressListParams struct {
	Page    int    `url:"page,omitempty" required:"false"`
	PerPage int    `url:"per_page,omitempty" required:"false"`
	Action  string `url:"action,omitempty" required:"false"`
	Cursor  string `url:"cursor,omitempty" required:"false"`
	lib.ListParams
}

type IpAddressGetReservedParams struct {
	Page    int    `url:"page,omitempty" required:"false"`
	PerPage int    `url:"per_page,omitempty" required:"false"`
	Action  string `url:"action,omitempty" required:"false"`
	Cursor  string `url:"cursor,omitempty" required:"false"`
}

func (i *IpAddress) UnmarshalJSON(data []byte) error {
	type ipAddress IpAddress
	var v ipAddress
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*i = IpAddress(v)
	return nil
}

func (i *IpAddressCollection) UnmarshalJSON(data []byte) error {
	type ipAddresss []IpAddress
	var v ipAddresss
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*i = IpAddressCollection(v)
	return nil
}
