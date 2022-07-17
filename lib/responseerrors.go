package lib

import (
	"fmt"
	"net/http"
)

func ResponseErrors(res *http.Response, errorFunc ...func(res *http.Response) error) error {
	for _, errFunc := range errorFunc {
		if err := errFunc(res); err != nil {
			return err
		}
	}
	return nil
}

func NonOkError(res *http.Response) error {
	if res.StatusCode > 299 && res.StatusCode < 200 {
		return errorFromBody(res)
	}

	return nil
}

func IsJSON(res *http.Response) bool {
	return res.Header.Get("Content-type") == "application/json"
}

func NonJSONError(res *http.Response) error {
	if !IsJSON(res) {
		if res.StatusCode == 204 {
			return nil
		}
		return errorFromBody(res)
	}

	return nil
}

func CloseBody(res *http.Response) {
	if res != nil && res.Body != nil {
		res.Body.Close()
	}
}

func errorFromBody(res *http.Response) error {
	var body []byte
	if res.ContentLength == -1 {
		body = make([]byte, 512)
	} else {
		body = make([]byte, res.ContentLength)
	}
	_, err := res.Body.Read(body)
	if err == nil {
		return fmt.Errorf(string(body))
	} else {
		return err
	}
}
