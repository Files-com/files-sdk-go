package files_sdk

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
	"github.com/hashicorp/go-retryablehttp"

	"moul.io/http2curl"
)

func Call(ctx context.Context, method string, config Config, resource string, params lib.Values) (*[]byte, *http.Response, error) {
	defaultHeaders := &http.Header{}
	config.SetHeaders(defaultHeaders)
	opts := &CallParams{
		Method:  method,
		Config:  config,
		Uri:     config.RootPath() + resource,
		Params:  params,
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
	Method   string
	Config   Config
	Uri      string
	Params   lib.Values
	BodyIo   io.ReadCloser
	Headers  *http.Header
	StayOpen bool
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
	var bodyIsJson bool
	if opts.Headers == nil {
		opts.Headers = &http.Header{}
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
			values, err := opts.Params.ToValues()
			if err != nil {
				return nil, err
			}
			req.URL.RawQuery = values.Encode()
		}
	default:
		if opts.BodyIo == nil {
			bodyIsJson = true
			jsonBody, err := opts.Params.ToJSON()
			if err != nil {
				return &http.Request{}, err
			}
			req.Body = ioutil.NopCloser(jsonBody)
			req.Header.Add("Content-Type", "application/json")
		} else {
			req.Body = opts.BodyIo
		}
	}

	req.Header = *opts.Headers
	if !opts.StayOpen {
		req.Header.Set("Connection", "close")
		req.Close = true
	}

	if opts.Config.InDebug() {
		defer debugLog(bodyIsJson, req, opts.Config, opts.Params)
	}
	return req, nil
}

func debugLog(bodyIsJson bool, req *http.Request, config Config, params lib.Values) {
	clonedReq := req.Clone(context.Background())
	clonedReq.Body = nil
	if bodyIsJson {
		jsonBody, err := params.ToJSON()
		if err != nil {
			panic(err)
		}
		clonedReq.Body = ioutil.NopCloser(jsonBody)
	}
	command, err := http2curl.GetCurlCommand(clonedReq)
	if err != nil {
		panic(err)
	}
	config.Logger().Printf(" %v", command)
}
