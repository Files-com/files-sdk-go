package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// IpAddress is a Files.com API resource.
type IpAddress struct {
	// Unique label for list; used by Zapier and other integrations.
	Id string `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// The object that this public IP address list is associated with.
	AssociatedWith string `json:"associated_with,omitempty" path:"associated_with,omitempty" url:"associated_with,omitempty"`
	// Group ID.
	GroupId int64 `json:"group_id,omitempty" path:"group_id,omitempty" url:"group_id,omitempty"`
	// A list of IP addresses.
	IpAddresses []string `json:"ip_addresses,omitempty" path:"ip_addresses,omitempty" url:"ip_addresses,omitempty"`
}

// Identifier returns the resource ID.
func (i IpAddress) Identifier() interface{} {
	return i.Id
}

// IpAddressCollection is a list of IpAddress resources.
type IpAddressCollection []IpAddress

// IpAddressListParams contains the request parameters for GET /ip_addresses.
type IpAddressListParams struct {
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// IpAddressGetSmartfileReservedParams contains the request parameters for GET /ip_addresses/smartfile-reserved.
type IpAddressGetSmartfileReservedParams struct {
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// IpAddressGetExavaultReservedParams contains the request parameters for GET /ip_addresses/exavault-reserved.
type IpAddressGetExavaultReservedParams struct {
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// IpAddressGetReservedParams contains the request parameters for GET /ip_addresses/reserved.
type IpAddressGetReservedParams struct {
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (i *IpAddress) UnmarshalJSON(data []byte) error {
	type ipAddress IpAddress
	var v ipAddress
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*i = IpAddress(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (i *IpAddressCollection) UnmarshalJSON(data []byte) error {
	type ipAddresss IpAddressCollection
	var v ipAddresss
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*i = IpAddressCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (i *IpAddressCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*i))
	for i, v := range *i {
		ret[i] = v
	}

	return &ret
}
