package files_sdk

import (
	"encoding/json"
	"fmt"
)

type ResponseError struct {
	Type         string          `json:"type"`
	Title        string          `json:"title"`
	ErrorMessage string          `json:"error"`
	HttpCode     int             `json:"http-code"`
	Errors       []ResponseError `json:"errors"`
	Data         Data            `json:"data"`
}

type SignRequest struct {
	Version   string `json:"version"`
	KeyHandle string `json:"keyHandle"`
}

type U2fSignRequests struct {
	AppId       string      `json:"app_id"`
	Challenge   string      `json:"challenge"`
	SignRequest SignRequest `json:"sign_request"`
}

type Data struct {
	U2fSIgnRequests               []U2fSignRequests `json:"u2f_sign_requests"`
	PartialSessionId              string            `json:"partial_session_id"`
	TwoFactorAuthenticationMethod []string          `json:"two_factor_authentication_methods"`
}

func (e ResponseError) Error() string {
	return fmt.Sprintf("%v - http-code: %v", e.ErrorMessage, e.HttpCode)
}

func (e ResponseError) IsNil() bool {
	return e.ErrorMessage == ""
}

func (e *ResponseError) UnmarshalJSON(data []byte) error {
	type re ResponseError
	var v re
	err := json.Unmarshal(data, &v)

	if err != nil {
		if er, ok := err.(*json.UnmarshalTypeError); ok && er.Value != "array" {
			return err
		} else {
			GlobalConfig.GetLogger().Printf(err.Error())
		}
	}

	*e = ResponseError(v)
	return nil
}
