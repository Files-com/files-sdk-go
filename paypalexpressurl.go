package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type PaypalExpressUrl struct {
	RedirectTo string `json:"redirect_to,omitempty" path:"redirect_to,omitempty" url:"redirect_to,omitempty"`
}

// Identifier no path or id

type PaypalExpressUrlCollection []PaypalExpressUrl

func (p *PaypalExpressUrl) UnmarshalJSON(data []byte) error {
	type paypalExpressUrl PaypalExpressUrl
	var v paypalExpressUrl
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*p = PaypalExpressUrl(v)
	return nil
}

func (p *PaypalExpressUrlCollection) UnmarshalJSON(data []byte) error {
	type paypalExpressUrls PaypalExpressUrlCollection
	var v paypalExpressUrls
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*p = PaypalExpressUrlCollection(v)
	return nil
}

func (p *PaypalExpressUrlCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*p))
	for i, v := range *p {
		ret[i] = v
	}

	return &ret
}
