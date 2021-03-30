package files_sdk

import (
	"encoding/json"
)

type PublicIpAddress struct {
	IpAddress  string `json:"ip_address,omitempty"`
	ServerName string `json:"server_name,omitempty"`
}

type PublicIpAddressCollection []PublicIpAddress

func (p *PublicIpAddress) UnmarshalJSON(data []byte) error {
	type publicIpAddress PublicIpAddress
	var v publicIpAddress
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*p = PublicIpAddress(v)
	return nil
}

func (p *PublicIpAddressCollection) UnmarshalJSON(data []byte) error {
	type publicIpAddresss []PublicIpAddress
	var v publicIpAddresss
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*p = PublicIpAddressCollection(v)
	return nil
}

func (p *PublicIpAddressCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*p))
	for i, v := range *p {
		ret[i] = v
	}

	return &ret
}
