package files_sdk

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"moul.io/http2curl"
)

func Call(method string, config Config, resource string, params url.Values) (*[]byte, *http.Response, error) {
	defaultHeaders := &http.Header{}
	config.SetHeaders(defaultHeaders)
	request, err := buildRequest(method, config, config.RootPath()+resource, &params, nil, defaultHeaders)
	if err != nil {
		return nil, &http.Response{}, err
	}
	response, err := config.GetHttpClient().Do(request)
	if err != nil {
		return nil, response, err
	}
	data, res, err := ParseResponse(response)
	responseError, ok := err.(ResponseError)
	if ok {
		err = responseError
	}
	return data, res, err
}

func ParseResponse(res *http.Response) (*[]byte, *http.Response, error) {
	defaultValue := make([]byte, 0)
	if res.StatusCode == 204 {
		return &defaultValue, res, nil
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return &defaultValue, res, err
	}
	re := ResponseError{}
	err = re.UnmarshalJSON(data)
	if err != nil {
		return &data, res, err
	}
	if !re.IsNil() {
		return &data, res, re
	}
	return &data, res, err

}

func CallRaw(method string, config Config, uri string, params *url.Values, body *[]byte, headers *http.Header) (*http.Response, error) {
	request, err := buildRequest(method, config, uri, params, body, headers)
	if err != nil {
		return &http.Response{}, err
	}
	return config.GetHttpClient().Do(request)
}

func buildRequest(method string, config Config, uri string, params *url.Values, body *[]byte, headers *http.Header) (*http.Request, error) {
	if headers == nil {
		headers = &http.Header{}
	}
	if params != nil {
		removeDash(params)
	}

	req, err := http.NewRequest(method, uri, nil)

	if err != nil {
		return &http.Request{}, err
	}

	if headers.Get("Content-Length") != "" {
		c, _ := strconv.ParseInt(headers.Get("Content-Length"), 10, 64)
		req.ContentLength = c
	}

	switch method {
	case "GET", "HEAD", "DELETE":
		if params != nil {
			removeDash(params)
			req.URL.RawQuery = params.Encode()
		}
	default:
		if body == nil {
			jsonBody, err := paramsToJson(params, headers)
			if err != nil {
				return &http.Request{}, err
			}
			req.Body = ioutil.NopCloser(jsonBody)
		} else {
			req.Body = ioutil.NopCloser(bytes.NewReader(*body))
		}
	}

	req.Header = *headers
	if config.InDebug() {
		command, err := http2curl.GetCurlCommand(req)
		if err != nil {
			panic(err)
		}
		config.Logger().Printf(" %v", command)
	}

	return req, nil
}

func paramsToJson(params *url.Values, headers *http.Header) (*bytes.Buffer, error) {
	bodyParams := make(map[string]string)
	for key, value := range *params {
		bodyParams[key] = value[0]
	}
	bodyBytes, err := json.Marshal(bodyParams)
	if err != nil {
		return nil, err
	}
	body := bytes.NewBuffer(bodyBytes)

	headers.Add("Content-Type", "application/json")
	return body, nil
}

func removeDash(params *url.Values) {
	for key := range *params {
		if string(key[0]) == "-" {
			params.Del(key)
		}
	}
}
