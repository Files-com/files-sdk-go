package files_sdk

import (
	"context"
	"io"
	"net/http"
)

type RequestResponseOption func(*http.Request, *http.Response) error

func RequestOption(opt func(req *http.Request) error) RequestResponseOption {
	return func(req *http.Request, res *http.Response) error {
		if req != nil {
			return opt(req)
		}
		return nil
	}
}

func ResponseOption(opt func(req *http.Response) error) RequestResponseOption {
	return func(req *http.Request, res *http.Response) error {
		if res != nil {
			return opt(res)
		}
		return nil
	}
}

func RequestHeadersOption(headers *http.Header) RequestResponseOption {
	return RequestOption(func(req *http.Request) error {
		for k, v := range *headers {
			req.Header.Set(k, v[0])
		}
		return nil
	})
}

func WithContext(ctx context.Context) RequestResponseOption {
	return RequestOption(func(req *http.Request) error {
		req.WithContext(ctx)
		return nil
	})
}

func ResponseBodyOption(opt func(io.ReadCloser) error) RequestResponseOption {
	return ResponseOption(func(res *http.Response) error {
		return opt(res.Body)
	})
}

func WrapRequestOptions(client HttpClient, request *http.Request, opts ...RequestResponseOption) (*http.Response, error) {
	for _, opt := range opts {
		err := opt(request, nil)
		if err != nil {
			return &http.Response{}, err
		}
	}

	resp, err := client.Do(request)
	if err != nil {
		return resp, err
	}

	for _, opt := range opts {
		err := opt(nil, resp)
		if err != nil {
			return resp, err
		}
	}

	return resp, nil
}
