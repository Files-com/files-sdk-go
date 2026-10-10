package files_sdk

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"

	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/hashicorp/go-retryablehttp"
	"moul.io/http2curl/v2"
)

// Resource sends an API resource request, decodes its response into the resource
// entity, and closes the response body.
func Resource(config Config, resource lib.Resource, opt ...RequestResponseOption) error {
	out, err := resource.Out()
	if err != nil {
		return err
	}
	data, res, err := Call(resource.Method, config, out.Path, out.Values, opt...)
	defer lib.CloseBody(res)
	if err != nil {
		return err
	}
	if res.StatusCode == 204 {
		return err
	}

	return out.Entity.UnmarshalJSON(*data)
}

// Call sends an API request and returns its response bytes, HTTP response, and
// any transport, response-reading, or API error. The caller must close the response body.
func Call(method string, config Config, resource string, params lib.Values, opts ...RequestResponseOption) (*[]byte, *http.Response, error) {
	defaultHeaders := &http.Header{}
	config.SetHeaders(defaultHeaders)
	callParams := &CallParams{
		Method:  method,
		Config:  config,
		Uri:     config.RootPath() + resource,
		Params:  params,
		Headers: defaultHeaders,
	}
	request, err := buildRequest(callParams)
	if err != nil {
		return nil, &http.Response{}, err
	}
	// In debug mode each response is logged as it arrives, including the ones the client retries and discards; the
	// final one is logged here only if the client did not log it.
	debug := config.InDebug()
	var logged *http.Response
	if debug && config.Client != nil && config.Client.HTTPClient != nil {
		config.Client = withResponseLogHook(config.Client, func(res *http.Response) {
			debugLogResponse(res, config)
			logged = res
		})
	}
	response, err := WrapRequestOptions(config, request, opts...)
	if err != nil {
		return nil, response, err
	}
	if debug && response != logged {
		debugLogResponse(response, config)
	}
	data, res, err := ParseResponse(response, resource)
	var responseError ResponseError
	ok := errors.As(err, &responseError)
	if ok {
		err = responseError
	}
	return data, res, err
}

// ParseResponse checks res for API errors and reads its body. It does not close
// the body; the caller must close it, including when this function returns an error.
func ParseResponse(res *http.Response, resource string) (*[]byte, *http.Response, error) {
	defaultValue := make([]byte, 0)
	if res.StatusCode == 204 {
		return &defaultValue, res, nil
	}
	nonOkError := lib.NonOkErrorCustom(func(err error) error {
		return fmt.Errorf("%v - %v", resource, err)
	})

	if err := lib.ResponseErrors(res, APIError(), nonOkError); err != nil {
		return &defaultValue, res, err
	}
	data, err := io.ReadAll(res.Body)
	return &data, res, err
}

// CallParams describes a raw HTTP request, including its body, context, and retry client.
type CallParams struct {
	Method string
	Config Config
	Uri    string
	Params lib.Values
	BodyIo io.ReadCloser
	// RetryableBody builds a fresh body reader for retries.
	// Leave nil unless this request needs retryable streaming body support.
	RetryableBody retryablehttp.ReaderFunc
	// Client overrides the HTTP client for this raw request.
	// Leave nil to use the configured SDK client.
	Client  *retryablehttp.Client
	Headers *http.Header
	context.Context
}

// CallRaw sends a request without decoding the response. Client, when provided,
// overrides the configured retry client. The caller must close the response body.
// Returned URL errors redact the request URL in either debug mode; DEBUG diagnostics
// may contain the full URL. Retry entries above DEBUG redact URLs and credentials.
func CallRaw(params *CallParams) (response *http.Response, err error) {
	defer func() {
		var urlErr *url.Error
		if errors.As(err, &urlErr) {
			if params.Config.InDebug() {
				params.Config.debugf("Transfer request failed: %v", err)
			}
			// Keep transport error types and cancellation checks without exposing the signed URL.
			err = &url.Error{Op: urlErr.Op, URL: "[redacted]", Err: urlErr.Err}
		}
	}()
	if params.Headers == nil {
		params.Headers = &http.Header{}
	}
	params.Config.SetUserAgentHeader(params.Headers)

	request, err := buildRequest(params)
	if err != nil {
		return &http.Response{}, err
	}
	retryRequest := &retryablehttp.Request{Request: request}
	if params.RetryableBody != nil {
		if err := retryRequest.SetBody(params.RetryableBody); err != nil {
			return &http.Response{}, err
		}
	} else {
		retryRequest.Body = request.Body
	}
	client := params.Config.Client
	if params.Client != nil {
		client = params.Client
	}
	if hasAuthHeaders(request.Header) {
		// Go copies custom headers to redirect targets, including other origins. Re-scope Files auth on each
		// redirect, as Config.Do does.
		client = params.Config.redirectSafeCopy(client)
	}
	// retryablehttp includes complete request URLs in its transport logs. This request's copy of the client keeps
	// them to DEBUG entries without changing the shared client's logger.
	client = &retryablehttp.Client{
		HTTPClient:      client.HTTPClient,
		Logger:          params.Config.retryLogFor(client.Logger, request),
		RetryWaitMin:    client.RetryWaitMin,
		RetryWaitMax:    client.RetryWaitMax,
		RetryMax:        client.RetryMax,
		RequestLogHook:  client.RequestLogHook,
		ResponseLogHook: client.ResponseLogHook,
		CheckRetry:      client.CheckRetry,
		Backoff:         client.Backoff,
		ErrorHandler:    client.ErrorHandler,
		PrepareRetry:    client.PrepareRetry,
	}
	return client.Do(retryRequest)
}

func buildRequest(opts *CallParams) (*http.Request, error) {
	var bodyIsJson bool
	if opts.Headers == nil {
		opts.Headers = &http.Header{}
	}

	paramsInQuery := opts.Method == "GET" || opts.Method == "HEAD" || opts.Method == "DELETE"
	var requestBody io.Reader
	var multipartContentType string
	var err error
	if !paramsInQuery && opts.BodyIo == nil {
		// Parameters that set a file are sent as multipart/form-data; all others as JSON.
		requestBody, multipartContentType, err = lib.MultipartBody(opts.Params)
		if err != nil {
			return &http.Request{}, err
		}
		if requestBody == nil {
			bodyIsJson = true
			requestBody, err = opts.Params.ToJSON()
			if err != nil {
				return &http.Request{}, err
			}
		}
	}

	var req *http.Request
	if opts.Context != nil {
		req, err = http.NewRequestWithContext(opts.Context, opts.Method, opts.Uri, requestBody)
	} else {
		req, err = http.NewRequest(opts.Method, opts.Uri, requestBody)
	}

	if err != nil {
		return &http.Request{}, err
	}
	req.Header = *opts.Headers
	if req.Header.Get("Content-Length") != "" {
		c, _ := strconv.ParseInt(req.Header.Get("Content-Length"), 10, 64)
		req.ContentLength = c
	}

	if paramsInQuery {
		if opts.Params != nil {
			values, err := opts.Params.ToValues()
			if err != nil {
				return nil, err
			}
			req.URL.RawQuery = values.Encode()
		}
	} else if bodyIsJson {
		req.Header.Add("Content-Type", "application/json")
	} else if multipartContentType != "" {
		req.Header.Set("Content-Type", multipartContentType)
	} else if req.ContentLength != 0 {
		req.Body = opts.BodyIo
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
		clonedReq.Body = io.NopCloser(jsonBody)
	}
	command, err := http2curl.GetCurlCommand(clonedReq)
	if err != nil {
		panic(err)
	}
	config.debugf(" %v", command)
}

// withResponseLogHook returns a copy of client for a single call that passes each response it receives, including
// the ones it retries, to logResponse before the client's own ResponseLogHook, if any. The copy shares the client's
// HTTP client, policies and hooks, so the client itself is not changed.
func withResponseLogHook(client *retryablehttp.Client, logResponse func(*http.Response)) *retryablehttp.Client {
	supplied := client.ResponseLogHook
	hook := func(logger retryablehttp.Logger, res *http.Response) {
		logResponse(res)
		if supplied != nil {
			supplied(logger, res)
		}
	}
	return &retryablehttp.Client{
		HTTPClient:      client.HTTPClient,
		Logger:          client.Logger,
		RetryWaitMin:    client.RetryWaitMin,
		RetryWaitMax:    client.RetryWaitMax,
		RetryMax:        client.RetryMax,
		RequestLogHook:  client.RequestLogHook,
		ResponseLogHook: hook,
		CheckRetry:      client.CheckRetry,
		Backoff:         client.Backoff,
		ErrorHandler:    client.ErrorHandler,
		PrepareRetry:    client.PrepareRetry,
	}
}

// debugLogResponse logs the response status line, headers and body, the counterpart of debugLog for the request.
// The body is read here and put back with the same bytes and, if reading failed, the same error, so the caller
// reads what the original body gave; closing it still closes the original response body.
func debugLogResponse(res *http.Response, config Config) {
	if res == nil || res.Body == nil {
		return
	}
	head, err := httputil.DumpResponse(res, false)
	if err != nil {
		return
	}
	original := res.Body
	body, readErr := io.ReadAll(original)
	res.Body = struct {
		io.Reader
		io.Closer
	}{io.MultiReader(bytes.NewReader(body), readResult{readErr}), original}
	config.debugf(" %s%s", head, body)
	if readErr != nil {
		config.debugf(" reading the response body failed: %v", readErr)
	}
}

// readResult ends a put-back body the way the original body ended: with its read error, or io.EOF when it was
// read in full.
type readResult struct {
	err error
}

func (r readResult) Read([]byte) (int, error) {
	if r.err == nil {
		return 0, io.EOF
	}
	return 0, r.err
}
