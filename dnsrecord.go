package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type DnsRecord struct {
	Id     string `json:"id,omitempty" path:"id"`
	Domain string `json:"domain,omitempty" path:"domain"`
	Rrtype string `json:"rrtype,omitempty" path:"rrtype"`
	Value  string `json:"value,omitempty" path:"value"`
}

type DnsRecordCollection []DnsRecord

type DnsRecordListParams struct {
	lib.ListParams
}

func (d *DnsRecord) UnmarshalJSON(data []byte) error {
	type dnsRecord DnsRecord
	var v dnsRecord
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*d = DnsRecord(v)
	return nil
}

func (d *DnsRecordCollection) UnmarshalJSON(data []byte) error {
	type dnsRecords DnsRecordCollection
	var v dnsRecords
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*d = DnsRecordCollection(v)
	return nil
}

func (d *DnsRecordCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*d))
	for i, v := range *d {
		ret[i] = v
	}

	return &ret
}
