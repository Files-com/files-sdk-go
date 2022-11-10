package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type IpAddress struct {
	Id             string   `json:"id,omitempty" path:"id"`
	AssociatedWith string   `json:"associated_with,omitempty" path:"associated_with"`
	GroupId        int64    `json:"group_id,omitempty" path:"group_id"`
	IpAddresses    []string `json:"ip_addresses,omitempty" path:"ip_addresses"`
}

type IpAddressCollection []IpAddress

type IpAddressListParams struct {
	lib.ListParams
}

type IpAddressGetExavaultReservedParams struct {
	lib.ListParams
}

type IpAddressGetReservedParams struct {
	lib.ListParams
}

func (i *IpAddress) UnmarshalJSON(data []byte) error {
	type ipAddress IpAddress
	var v ipAddress
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*i = IpAddress(v)
	return nil
}

func (i *IpAddressCollection) UnmarshalJSON(data []byte) error {
	type ipAddresss IpAddressCollection
	var v ipAddresss
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*i = IpAddressCollection(v)
	return nil
}

func (i *IpAddressCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*i))
	for i, v := range *i {
		ret[i] = v
	}

	return &ret
}
