package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// DnsRecord is a Files.com API resource.
type DnsRecord struct {
	Id     string `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	Domain string `json:"domain,omitempty" path:"domain,omitempty" url:"domain,omitempty"`
	Rrtype string `json:"rrtype,omitempty" path:"rrtype,omitempty" url:"rrtype,omitempty"`
	Value  string `json:"value,omitempty" path:"value,omitempty" url:"value,omitempty"`
}

// Identifier returns the resource ID.
func (d DnsRecord) Identifier() interface{} {
	return d.Id
}

// DnsRecordCollection is a list of DnsRecord resources.
type DnsRecordCollection []DnsRecord

// DnsRecordListParams contains the request parameters for GET /dns_records.
type DnsRecordListParams struct {
	ListParams
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (d *DnsRecord) UnmarshalJSON(data []byte) error {
	type dnsRecord DnsRecord
	var v dnsRecord
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*d = DnsRecord(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (d *DnsRecordCollection) UnmarshalJSON(data []byte) error {
	type dnsRecords DnsRecordCollection
	var v dnsRecords
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*d = DnsRecordCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (d *DnsRecordCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*d))
	for i, v := range *d {
		ret[i] = v
	}

	return &ret
}
