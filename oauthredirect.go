package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type OauthRedirect struct {
	RedirectUri string `json:"redirect_uri,omitempty" path:"redirect_uri,omitempty" url:"redirect_uri,omitempty"`
}

// Identifier no path or id

type OauthRedirectCollection []OauthRedirect

func (o *OauthRedirect) UnmarshalJSON(data []byte) error {
	type oauthRedirect OauthRedirect
	var v oauthRedirect
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*o = OauthRedirect(v)
	return nil
}

func (o *OauthRedirectCollection) UnmarshalJSON(data []byte) error {
	type oauthRedirects OauthRedirectCollection
	var v oauthRedirects
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*o = OauthRedirectCollection(v)
	return nil
}

func (o *OauthRedirectCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*o))
	for i, v := range *o {
		ret[i] = v
	}

	return &ret
}
