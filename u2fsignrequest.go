package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type U2fSignRequest struct {
	AppId                         string   `json:"app_id,omitempty" path:"app_id,omitempty" url:"app_id,omitempty"`
	Challenge                     string   `json:"challenge,omitempty" path:"challenge,omitempty" url:"challenge,omitempty"`
	SignRequest                   string   `json:"sign_request,omitempty" path:"sign_request,omitempty" url:"sign_request,omitempty"`
	WebauthnAuthenticationOptions []string `json:"webauthn_authentication_options,omitempty" path:"webauthn_authentication_options,omitempty" url:"webauthn_authentication_options,omitempty"`
}

// Identifier no path or id

type U2fSignRequestCollection []U2fSignRequest

func (u *U2fSignRequest) UnmarshalJSON(data []byte) error {
	type u2fSignRequest U2fSignRequest
	var v u2fSignRequest
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*u = U2fSignRequest(v)
	return nil
}

func (u *U2fSignRequestCollection) UnmarshalJSON(data []byte) error {
	type u2fSignRequests U2fSignRequestCollection
	var v u2fSignRequests
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*u = U2fSignRequestCollection(v)
	return nil
}

func (u *U2fSignRequestCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*u))
	for i, v := range *u {
		ret[i] = v
	}

	return &ret
}
