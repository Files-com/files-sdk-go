package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Ip struct {
	Ip                    string         `json:"ip,omitempty" path:"ip,omitempty" url:"ip,omitempty"`
	ExternalIp            string         `json:"external_ip,omitempty" path:"external_ip,omitempty" url:"external_ip,omitempty"`
	Assigned              string         `json:"assigned,omitempty" path:"assigned,omitempty" url:"assigned,omitempty"`
	Site                  SslCertificate `json:"site,omitempty" path:"site,omitempty" url:"site,omitempty"`
	FtpEnabled            string         `json:"ftp_enabled,omitempty" path:"ftp_enabled,omitempty" url:"ftp_enabled,omitempty"`
	SftpEnabled           string         `json:"sftp_enabled,omitempty" path:"sftp_enabled,omitempty" url:"sftp_enabled,omitempty"`
	SftpHostKeyType       string         `json:"sftp_host_key_type,omitempty" path:"sftp_host_key_type,omitempty" url:"sftp_host_key_type,omitempty"`
	SftpHostKeyPrivateKey string         `json:"sftp_host_key_private_key,omitempty" path:"sftp_host_key_private_key,omitempty" url:"sftp_host_key_private_key,omitempty"`
	SiteId                string         `json:"site_id,omitempty" path:"site_id,omitempty" url:"site_id,omitempty"`
	MotdText              string         `json:"motd_text,omitempty" path:"motd_text,omitempty" url:"motd_text,omitempty"`
	MotdUseForFtp         *bool          `json:"motd_use_for_ftp,omitempty" path:"motd_use_for_ftp,omitempty" url:"motd_use_for_ftp,omitempty"`
	MotdUseForSftp        *bool          `json:"motd_use_for_sftp,omitempty" path:"motd_use_for_sftp,omitempty" url:"motd_use_for_sftp,omitempty"`
	PairType              string         `json:"pair_type,omitempty" path:"pair_type,omitempty" url:"pair_type,omitempty"`
}

// Identifier no path or id

type IpCollection []Ip

func (i *Ip) UnmarshalJSON(data []byte) error {
	type ip Ip
	var v ip
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*i = Ip(v)
	return nil
}

func (i *IpCollection) UnmarshalJSON(data []byte) error {
	type ips IpCollection
	var v ips
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*i = IpCollection(v)
	return nil
}

func (i *IpCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*i))
	for i, v := range *i {
		ret[i] = v
	}

	return &ret
}
