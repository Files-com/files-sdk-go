package files_sdk

import (
	"context"
	"io"
	"net/http"
	"strconv"
)

type requestResponseOption struct {
	*http.Request
	*http.Response
	context.Context
}

// RequestResponseOption customizes an individual SDK request or response.
// Use WithContext for cancellation and deadlines. Use Config.APIKey,
// Config.SessionId, and Config.WorkspaceId for authentication and workspace
// scoping. Request mutation options support custom HTTP behavior; some
// higher-level helpers do not pass them through.
type RequestResponseOption func(*requestResponseOption) error

// RequestOption calls call with the HTTP request before it is sent.
// An error from call stops the request.
func RequestOption(call func(req *http.Request) error) RequestResponseOption {
	return func(opt *requestResponseOption) error {
		if opt.Request != nil {
			return call(opt.Request)
		}
		return nil
	}
}

// ResponseOption calls call with the HTTP response before the SDK processes it.
// An error from call stops response processing.
func ResponseOption(call func(req *http.Response) error) RequestResponseOption {
	return func(opt *requestResponseOption) error {
		if opt.Response != nil {
			return call(opt.Response)
		}
		return nil
	}
}

// RequestHeadersOption sets request headers from headers. Only the first value
// of each header is used; each supplied header must have at least one value.
func RequestHeadersOption(headers *http.Header) RequestResponseOption {
	return RequestOption(func(req *http.Request) error {
		for k, v := range *headers {
			req.Header.Set(k, v[0])
		}
		return nil
	})
}

// WithWorkspaceId sets the workspace header for a request.
//
// Deprecated: use Config.WorkspaceId; request/header mutations, including API-key headers, are not guaranteed through every helper.
func WithWorkspaceId(workspaceId int64) RequestResponseOption {
	return RequestOption(func(req *http.Request) error {
		req.Header.Set(workspaceIdHeader, strconv.FormatInt(workspaceId, 10))
		return nil
	})
}

// WithoutWorkspaceId removes the workspace header from a request.
//
// Deprecated: use Config without WorkspaceId; request/header mutations, including API-key headers, are not guaranteed through every helper.
func WithoutWorkspaceId() RequestResponseOption {
	return RequestOption(func(req *http.Request) error {
		req.Header.Del(workspaceIdHeader)
		return nil
	})
}

// WithContext sets the context used for request cancellation and deadlines.
// A non-nil ctx replaces any context set by earlier options.
func WithContext(ctx context.Context) RequestResponseOption {
	return func(opt *requestResponseOption) error {
		if opt.Request != nil && ctx != nil {
			opt.Request = opt.Request.WithContext(ctx)
		} else {
			opt.Context = ctx
		}
		return nil
	}
}

// ResponseBodyOption calls opt with the response body before the SDK processes it.
// Reading the body consumes those bytes for subsequent response processing.
func ResponseBodyOption(opt func(io.ReadCloser) error) RequestResponseOption {
	return ResponseOption(func(res *http.Response) error {
		return opt(res.Body)
	})
}

// WrapRequestOptions applies request options, sends the request with config,
// and applies response options. The caller must close a returned response body.
func WrapRequestOptions(config Config, request *http.Request, opts ...RequestResponseOption) (*http.Response, error) {
	// Apply request options
	modifiedRequest, err := BuildRequest(request, opts...)
	if err != nil {
		return &http.Response{}, err
	}

	// Execute the request
	response, err := config.Do(modifiedRequest)
	if err != nil {
		return response, err
	}

	return BuildResponse(response, opts...)
}

// ContextOption extracts the context from opts, or returns context.Background
// when no context is set. Errors from options are ignored.
func ContextOption(opts []RequestResponseOption) context.Context {
	params := &requestResponseOption{}
	for _, opt := range opts {
		opt(params)
	}
	if params.Context == nil {
		params.Context = context.Background()
	}
	return params.Context
}

// BuildRequest applies opts in order and returns the resulting request.
// It stops at the first option error and does not send the request.
func BuildRequest(request *http.Request, opts ...RequestResponseOption) (*http.Request, error) {
	optionRequestResponse := &requestResponseOption{Request: request}

	for _, opt := range opts {
		if err := opt(optionRequestResponse); err != nil {
			return nil, err
		}
	}

	return optionRequestResponse.Request, nil
}

// BuildResponse applies opts in order and returns the resulting response.
// It stops at the first option error.
func BuildResponse(response *http.Response, opts ...RequestResponseOption) (*http.Response, error) {
	optionRequestResponse := &requestResponseOption{Response: response}

	for _, opt := range opts {
		if err := opt(optionRequestResponse); err != nil {
			return nil, err
		}
	}

	return optionRequestResponse.Response, nil
}
