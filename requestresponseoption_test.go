package files_sdk

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithContext_NoRequest(t *testing.T) {
	// Test that WithContext stores the context when there's no request
	ctx := context.WithValue(context.Background(), "test-key", "test-value")
	opt := &requestResponseOption{}

	err := WithContext(ctx)(opt)

	require.NoError(t, err)
	assert.NotNil(t, opt.Context)
	assert.Equal(t, "test-value", opt.Context.Value("test-key"))
}

func TestWithContext_WithRequest(t *testing.T) {
	// Test that WithContext sets the context on the request
	ctx := context.WithValue(context.Background(), "test-key", "test-value")
	req, err := http.NewRequest("GET", "http://example.com", nil)
	require.NoError(t, err)

	opt := &requestResponseOption{Request: req}

	err = WithContext(ctx)(opt)

	require.NoError(t, err)
	assert.NotNil(t, opt.Request.Context())
	assert.Equal(t, "test-value", opt.Request.Context().Value("test-key"))
}

func TestWithContext_NilContext(t *testing.T) {
	// Test that WithContext handles nil context gracefully
	req, err := http.NewRequest("GET", "http://example.com", nil)
	require.NoError(t, err)

	opt := &requestResponseOption{Request: req}

	err = WithContext(nil)(opt)

	require.NoError(t, err)
	// The original request context should remain unchanged
	assert.NotNil(t, opt.Request.Context())
}

func TestWithContext_NilContextNoRequest(t *testing.T) {
	// Test that WithContext stores nil context when there's no request
	opt := &requestResponseOption{}

	err := WithContext(nil)(opt)

	require.NoError(t, err)
	assert.Nil(t, opt.Context)
}

func TestContextOption_RetrievesContext(t *testing.T) {
	// Test that ContextOption can retrieve the context from options
	ctx := context.WithValue(context.Background(), "test-key", "test-value")
	opts := []RequestResponseOption{WithContext(ctx)}

	retrievedCtx := ContextOption(opts)

	assert.NotNil(t, retrievedCtx)
	assert.Equal(t, "test-value", retrievedCtx.Value("test-key"))
}

func TestContextOption_NoContext(t *testing.T) {
	// Test that ContextOption returns context.Background() when no context is provided
	opts := []RequestResponseOption{}

	retrievedCtx := ContextOption(opts)

	assert.NotNil(t, retrievedCtx)
	assert.Equal(t, context.Background(), retrievedCtx)
}

func TestContextOption_MultipleOptions(t *testing.T) {
	// Test that ContextOption works with multiple options
	ctx := context.WithValue(context.Background(), "test-key", "test-value")
	opts := []RequestResponseOption{
		RequestOption(func(req *http.Request) error {
			req.Header.Set("X-Test", "test")
			return nil
		}),
		WithContext(ctx),
	}

	retrievedCtx := ContextOption(opts)

	assert.NotNil(t, retrievedCtx)
	assert.Equal(t, "test-value", retrievedCtx.Value("test-key"))
}

func TestContextOption_ContextOverride(t *testing.T) {
	// Test that later context options override earlier ones
	ctx1 := context.WithValue(context.Background(), "key1", "value1")
	ctx2 := context.WithValue(context.Background(), "key2", "value2")

	opts := []RequestResponseOption{
		WithContext(ctx1),
		WithContext(ctx2),
	}

	retrievedCtx := ContextOption(opts)

	assert.NotNil(t, retrievedCtx)
	assert.Nil(t, retrievedCtx.Value("key1"))
	assert.Equal(t, "value2", retrievedCtx.Value("key2"))
}

func TestWithContext_RequestContainsContext(t *testing.T) {
	// Test that the request properly contains the context after applying the option
	ctx := context.WithValue(context.Background(), "test-key", "test-value")

	req, err := http.NewRequest("GET", "http://example.com", nil)
	require.NoError(t, err)

	// Apply the WithContext option using BuildRequest
	modifiedReq, err := BuildRequest(req, WithContext(ctx))
	require.NoError(t, err)

	// Verify the request contains the context
	assert.NotNil(t, modifiedReq.Context())
	assert.Equal(t, "test-value", modifiedReq.Context().Value("test-key"))
}

func TestBuildRequest(t *testing.T) {
	// Test BuildRequest applies multiple options correctly
	ctx := context.WithValue(context.Background(), "test-key", "test-value")
	req, err := http.NewRequest("GET", "http://example.com", nil)
	require.NoError(t, err)

	opts := []RequestResponseOption{
		WithContext(ctx),
		RequestOption(func(r *http.Request) error {
			r.Header.Set("X-Custom-Header", "custom-value")
			return nil
		}),
	}

	modifiedReq, err := BuildRequest(req, opts...)
	require.NoError(t, err)
	assert.NotNil(t, modifiedReq)
	assert.Equal(t, "test-value", modifiedReq.Context().Value("test-key"))
	assert.Equal(t, "custom-value", modifiedReq.Header.Get("X-Custom-Header"))
}

func TestBuildResponse(t *testing.T) {
	// Test BuildResponse applies response options correctly
	resp := &http.Response{
		StatusCode: 200,
		Header:     make(http.Header),
	}

	opts := []RequestResponseOption{
		ResponseOption(func(r *http.Response) error {
			r.Header.Set("X-Response-Header", "response-value")
			return nil
		}),
	}

	modifiedResp, err := BuildResponse(resp, opts...)
	require.NoError(t, err)
	assert.NotNil(t, modifiedResp)
	assert.Equal(t, 200, modifiedResp.StatusCode)
	assert.Equal(t, "response-value", modifiedResp.Header.Get("X-Response-Header"))
}

func TestBuildRequest_Error(t *testing.T) {
	// Test BuildRequest returns error when option fails
	req, err := http.NewRequest("GET", "http://example.com", nil)
	require.NoError(t, err)

	opts := []RequestResponseOption{
		RequestOption(func(r *http.Request) error {
			return assert.AnError
		}),
	}

	_, err = BuildRequest(req, opts...)
	assert.Error(t, err)
}

func TestBuildResponse_Error(t *testing.T) {
	// Test BuildResponse returns error when option fails
	resp := &http.Response{StatusCode: 200}

	opts := []RequestResponseOption{
		ResponseOption(func(r *http.Response) error {
			return assert.AnError
		}),
	}

	_, err := BuildResponse(resp, opts...)
	assert.Error(t, err)
}
