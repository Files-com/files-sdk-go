package files_sdk

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type ResponseError struct {
	Type         string          `json:"type"`
	Title        string          `json:"title"`
	ErrorMessage string          `json:"error"`
	HttpCode     int             `json:"http-code"`
	Errors       []ResponseError `json:"errors"`
	Data         Data            `json:"data"`
}

const (
	DestinationExists = "processing-failure/destination-exists"
)

func IsDestinationExistsError(err error) bool {
	re, ok := err.(ResponseError)
	return ok && re.Type == DestinationExists
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
	return fmt.Sprintf("%v - `%v`", e.Title, e.ErrorMessage)
}

func (e ResponseError) IsNil() bool {
	return e.ErrorMessage == ""
}

func (e *ResponseError) UnmarshalJSON(data []byte) error {
	type re ResponseError
	var v re
	err := json.Unmarshal(data, &v)

	if err != nil {
		jsonError, ok := err.(*json.UnmarshalTypeError)

		if ok && jsonError.Field == "" {
			if jsonError.Value == "string" {
				var str string
				json.Unmarshal(data, &str)
				v.ErrorMessage = str
			} else if jsonError.Value != "array" {
				return err
			}
		} else if ok && jsonError.Field == "http-code" {
			tmp := make(map[string]interface{})
			json.Unmarshal(data, &tmp)
			intVar, _ := strconv.Atoi(tmp["http-code"].(string))
			v.HttpCode = intVar
		} else {
			return err
		}

		jsonSyntaxErr, ok := err.(*json.SyntaxError)
		if ok && jsonSyntaxErr.Error() == "invalid character '<' looking for beginning of value" {
			return fmt.Errorf(string(data))
		}
	}

	*e = ResponseError(v)
	return nil
}
