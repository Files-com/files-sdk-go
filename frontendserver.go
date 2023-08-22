package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type FrontEndServer struct {
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
	Name                  string         `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	Hostname              string         `json:"hostname,omitempty" path:"hostname,omitempty" url:"hostname,omitempty"`
	Zone                  string         `json:"zone,omitempty" path:"zone,omitempty" url:"zone,omitempty"`
	Ips                   []string       `json:"ips,omitempty" path:"ips,omitempty" url:"ips,omitempty"`
	PrimaryIp             string         `json:"primary_ip,omitempty" path:"primary_ip,omitempty" url:"primary_ip,omitempty"`
	PrimaryIpPublic       string         `json:"primary_ip_public,omitempty" path:"primary_ip_public,omitempty" url:"primary_ip_public,omitempty"`
	SooIp                 string         `json:"soo_ip,omitempty" path:"soo_ip,omitempty" url:"soo_ip,omitempty"`
	SooIpPublic           string         `json:"soo_ip_public,omitempty" path:"soo_ip_public,omitempty" url:"soo_ip_public,omitempty"`
	ExavaultIp            string         `json:"exavault_ip,omitempty" path:"exavault_ip,omitempty" url:"exavault_ip,omitempty"`
	ExavaultIpPublic      string         `json:"exavault_ip_public,omitempty" path:"exavault_ip_public,omitempty" url:"exavault_ip_public,omitempty"`
	ExavaultSooIp         string         `json:"exavault_soo_ip,omitempty" path:"exavault_soo_ip,omitempty" url:"exavault_soo_ip,omitempty"`
	ExavaultSooIpPublic   string         `json:"exavault_soo_ip_public,omitempty" path:"exavault_soo_ip_public,omitempty" url:"exavault_soo_ip_public,omitempty"`
}

// Identifier no path or id

type FrontEndServerCollection []FrontEndServer

type IpsParam struct {
	PrivateIp string `url:"private_ip,omitempty" json:"private_ip,omitempty" path:"private_ip"`
	PublicIp  string `url:"public_ip,omitempty" json:"public_ip,omitempty" path:"public_ip"`
}

type FrontEndServerCreateParams struct {
	Name                string     `url:"name,omitempty" required:"true" json:"name,omitempty" path:"name"`
	Hostname            string     `url:"hostname,omitempty" required:"false" json:"hostname,omitempty" path:"hostname"`
	Zone                string     `url:"zone,omitempty" required:"false" json:"zone,omitempty" path:"zone"`
	Ips                 []string   `url:"ips,omitempty" required:"false" json:"ips,omitempty" path:"ips"`
	IpsParam            []IpsParam `url:"ips,omitempty" required:"false" json:"ips,omitempty" path:"ips"`
	PrimaryIp           string     `url:"primary_ip,omitempty" required:"false" json:"primary_ip,omitempty" path:"primary_ip"`
	PrimaryIpPublic     string     `url:"primary_ip_public,omitempty" required:"false" json:"primary_ip_public,omitempty" path:"primary_ip_public"`
	SooIp               string     `url:"soo_ip,omitempty" required:"false" json:"soo_ip,omitempty" path:"soo_ip"`
	SooIpPublic         string     `url:"soo_ip_public,omitempty" required:"false" json:"soo_ip_public,omitempty" path:"soo_ip_public"`
	ExavaultIp          string     `url:"exavault_ip,omitempty" required:"false" json:"exavault_ip,omitempty" path:"exavault_ip"`
	ExavaultIpPublic    string     `url:"exavault_ip_public,omitempty" required:"false" json:"exavault_ip_public,omitempty" path:"exavault_ip_public"`
	ExavaultSooIp       string     `url:"exavault_soo_ip,omitempty" required:"false" json:"exavault_soo_ip,omitempty" path:"exavault_soo_ip"`
	ExavaultSooIpPublic string     `url:"exavault_soo_ip_public,omitempty" required:"false" json:"exavault_soo_ip_public,omitempty" path:"exavault_soo_ip_public"`
}

func (f *FrontEndServer) UnmarshalJSON(data []byte) error {
	type frontEndServer FrontEndServer
	var v frontEndServer
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*f = FrontEndServer(v)
	return nil
}

func (f *FrontEndServerCollection) UnmarshalJSON(data []byte) error {
	type frontEndServers FrontEndServerCollection
	var v frontEndServers
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*f = FrontEndServerCollection(v)
	return nil
}

func (f *FrontEndServerCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*f))
	for i, v := range *f {
		ret[i] = v
	}

	return &ret
}
