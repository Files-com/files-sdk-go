package files_sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/hashicorp/go-retryablehttp"

	"moul.io/http2curl"
)

func Call(ctx context.Context, method string, config Config, resource string, params url.Values) (*[]byte, *http.Response, error) {
	defaultHeaders := &http.Header{}
	config.SetHeaders(defaultHeaders)
	opts := &CallParams{
		Method:  method,
		Config:  config,
		Uri:     config.RootPath() + resource,
		Params:  &params,
		Headers: defaultHeaders,
		Context: ctx,
	}
	request, err := buildRequest(opts)
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

type CallParams struct {
	Method  string
	Config  Config
	Uri     string
	Params  *url.Values
	BodyIo  io.ReadCloser
	Headers *http.Header
	context.Context
}

func CallRaw(params *CallParams) (*http.Response, error) {
	request, err := buildRequest(params)
	if err != nil {
		return &http.Response{}, err
	}
	if request.Body != nil {
		retryRequest := &retryablehttp.Request{Request: request}
		retryRequest.SetBody(func() (io.Reader, error) { return request.Body, nil })
		return params.Config.GetRawClient().Do(retryRequest)
	} else {
		return params.Config.GetHttpClient().Do(request)
	}
}

func buildRequest(opts *CallParams) (*http.Request, error) {
	if opts.Headers == nil {
		opts.Headers = &http.Header{}
	}
	if opts.Params != nil {
		removeDash(opts.Params)
	}
	var req *http.Request
	var err error
	if opts.Context != nil {
		req, err = http.NewRequestWithContext(opts.Context, opts.Method, opts.Uri, nil)
	} else {
		req, err = http.NewRequest(opts.Method, opts.Uri, nil)
	}

	if err != nil {
		return &http.Request{}, err
	}

	if opts.Headers.Get("Content-Length") != "" {
		c, _ := strconv.ParseInt(opts.Headers.Get("Content-Length"), 10, 64)
		req.ContentLength = c
	}

	switch opts.Method {
	case "GET", "HEAD", "DELETE":
		if opts.Params != nil {
			removeDash(opts.Params)
			req.URL.RawQuery = opts.Params.Encode()
		}
	default:
		if opts.BodyIo == nil {
			jsonBody, err := paramsToJson(opts.Params, opts.Headers)
			if err != nil {
				return &http.Request{}, err
			}
			req.Body = ioutil.NopCloser(jsonBody)
		} else {
			req.Body = opts.BodyIo
		}
	}

	req.Header = *opts.Headers
	if opts.Config.InDebug() {
		withoutBodyReq := *req
		withoutBodyReq.Body = nil
		command, err := http2curl.GetCurlCommand(&withoutBodyReq)
		if err != nil {
			panic(err)
		}
		opts.Config.Logger().Printf(" %v", command)
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
