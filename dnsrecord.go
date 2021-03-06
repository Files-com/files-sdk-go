package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/lib"
)

type DnsRecord struct {
	Id     string `json:"id,omitempty"`
	Domain string `json:"domain,omitempty"`
	Rrtype string `json:"rrtype,omitempty"`
	Value  string `json:"value,omitempty"`
}

type DnsRecordCollection []DnsRecord

type DnsRecordListParams struct {
	Cursor  string `url:"cursor,omitempty" required:"false"`
	PerPage int64  `url:"per_page,omitempty" required:"false"`
	lib.ListParams
}

func (d *DnsRecord) UnmarshalJSON(data []byte) error {
	type dnsRecord DnsRecord
	var v dnsRecord
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*d = DnsRecord(v)
	return nil
}

func (d *DnsRecordCollection) UnmarshalJSON(data []byte) error {
	type dnsRecords []DnsRecord
	var v dnsRecords
	if err := json.Unmarshal(data, &v); err != nil {
		return err
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
