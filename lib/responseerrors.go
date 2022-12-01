package lib

import (
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
)

type ResponseError struct {
	StatusCode int
	error
}

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

func NotStatus(status int) func(res *http.Response) error {
	return func(res *http.Response) error {
		if res.StatusCode != status {
			return errorFromBody(res)
		}

		return nil
	}
}

func IsStatus(status int) func(res *http.Response) error {
	return func(res *http.Response) error {
		if res.StatusCode == status {
			return errorFromBody(res)
		}

		return nil
	}
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
	body = make([]byte, int(math.Max(float64(res.ContentLength), float64(512))))
	_, err := io.ReadFull(res.Body, body)
	defer CloseBody(res)
	if err == nil || errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
		return ResponseError{StatusCode: res.StatusCode, error: fmt.Errorf(string(body))}
	} else {
		return ResponseError{StatusCode: res.StatusCode, error: err}
	}
}
