package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type PublicUrl struct {
	HttpHeaders         string                 `json:"http_headers,omitempty" path:"http_headers,omitempty" url:"http_headers,omitempty"`
	Body                string                 `json:"body,omitempty" path:"body,omitempty" url:"body,omitempty"`
	DownloadUri         string                 `json:"download_uri,omitempty" path:"download_uri,omitempty" url:"download_uri,omitempty"`
	InternalDownloadUri string                 `json:"internal_download_uri,omitempty" path:"internal_download_uri,omitempty" url:"internal_download_uri,omitempty"`
	Error               string                 `json:"error,omitempty" path:"error,omitempty" url:"error,omitempty"`
	Redirect            string                 `json:"redirect,omitempty" path:"redirect,omitempty" url:"redirect,omitempty"`
	Status              int64                  `json:"status,omitempty" path:"status,omitempty" url:"status,omitempty"`
	MimeType            string                 `json:"mime_type,omitempty" path:"mime_type,omitempty" url:"mime_type,omitempty"`
	RemoteServerId      int64                  `json:"remote_server_id,omitempty" path:"remote_server_id,omitempty" url:"remote_server_id,omitempty"`
	Headers             map[string]interface{} `json:"headers,omitempty" path:"headers,omitempty" url:"headers,omitempty"`
	SocksIps            map[string]interface{} `json:"socks_ips,omitempty" path:"socks_ips,omitempty" url:"socks_ips,omitempty"`
	Hostname            string                 `json:"hostname,omitempty" path:"hostname,omitempty" url:"hostname,omitempty"`
	Path                string                 `json:"path,omitempty" path:"path,omitempty" url:"path,omitempty"`
}

func (p PublicUrl) Identifier() interface{} {
	return p.Path
}

type PublicUrlCollection []PublicUrl

type PublicUrlCreateParams struct {
	Hostname string `url:"hostname,omitempty" required:"true" json:"hostname,omitempty" path:"hostname"`
	Path     string `url:"path,omitempty" required:"true" json:"path,omitempty" path:"path"`
}

func (p *PublicUrl) UnmarshalJSON(data []byte) error {
	type publicUrl PublicUrl
	var v publicUrl
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*p = PublicUrl(v)
	return nil
}

func (p *PublicUrlCollection) UnmarshalJSON(data []byte) error {
	type publicUrls PublicUrlCollection
	var v publicUrls
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*p = PublicUrlCollection(v)
	return nil
}

func (p *PublicUrlCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*p))
	for i, v := range *p {
		ret[i] = v
	}

	return &ret
}
