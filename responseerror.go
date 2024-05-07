package files_sdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Files-com/files-sdk-go/v3/lib"
)

type ResponseError struct {
	Type         string          `json:"type,omitempty"`
	Title        string          `json:"title,omitempty"`
	ErrorMessage string          `json:"error,omitempty"`
	HttpCode     int             `json:"http-code,omitempty"`
	Errors       []ResponseError `json:"errors,omitempty"`
	Data         Data            `json:"data,omitempty"`
}

const (
	DestinationExists = "processing-failure/destination-exists"
)

func IsExist(err error) bool {
	var re ResponseError
	ok := errors.As(err, &re)
	return ok && re.Type == DestinationExists
}

func IsNotExist(err error) bool {
	var re ResponseError
	ok := errors.As(err, &re)
	return ok && strings.Split(re.Type, "/")[0] == "not-found"
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
	U2fSIgnRequests               []U2fSignRequests `json:"u2f_sign_requests,omitempty"`
	PartialSessionId              string            `json:"partial_session_id,omitempty"`
	TwoFactorAuthenticationMethod []string          `json:"two_factor_authentication_methods,omitempty"`
	Host                          string            `json:"host,omitempty"`
	// Download Request Status
	BytesTransferred int64     `json:"bytes_transferred,omitempty"`
	Status           string    `json:"status,omitempty"`
	StartedAt        time.Time `json:"started_at,omitempty"`
	CompletedAt      time.Time `json:"completed_at,omitempty"`
	TouchedAt        time.Time `json:"touched_at,omitempty"`
}

func (e ResponseError) Error() string {
	return fmt.Sprintf("%v - `%v`", e.Title, e.ErrorMessage)
}

func (e ResponseError) IsNil() bool {
	return e.ErrorMessage == ""
}

func (e ResponseError) Is(err error) bool {
	var responseError ResponseError
	return errors.As(err, &responseError)
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

func APIError(callbacks ...func(ResponseError) ResponseError) func(res *http.Response) error {
	return func(res *http.Response) error {
		if lib.IsNonOkStatus(res) && lib.IsHTML(res) && res.Header.Get("X-Request-Id") != "" && res.Header.Get("Server") == "nginx" {
			return fmt.Errorf("files.com Server error - request id: %v", res.Header.Get("X-Request-Id"))
		}

		if lib.IsNonOkStatus(res) && lib.IsJSON(res) {
			data, err := io.ReadAll(res.Body)
			if err != nil {
				return lib.NonOkError(res)
			}

			re := ResponseError{}

			err = re.UnmarshalJSON(data)
			if err != nil {
				return lib.NonOkError(res)
			}

			if re.IsNil() {
				return lib.NonOkError(res)
			}
			for _, callback := range callbacks {
				re = callback(re)
			}
			return re
		}
		return nil
	}
}
