package files_sdk

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/Files-com/files-sdk-go/v3/lib"
)

const (
	directTransferRequestFailed   = "direct_request_failed"
	directTransferResponseFailed  = "direct_response_failed"
	directTransferBodyReadFailed  = "direct_body_read_failed"
	directTransferBodyCloseFailed = "direct_body_close_failed"
)

// DirectTransferSession sends the requests of one transfer, such as the ranges
// of a download, to an Agent's direct endpoint until a direct request fails.
// From then on Send returns ErrDirectTransferUnavailable without sending, and
// the caller sends the rest of the transfer over its own route. A session is
// safe for concurrent requests. Call Close when the transfer ends.
type DirectTransferSession struct {
	options        DirectTransferSessionOptions
	cache          directTransferClientCache
	attempts       atomic.Int64
	successes      atomic.Int64
	disabledReason atomic.Pointer[string]
}

// DirectTransferSessionOptions customize how a DirectTransferSession handles
// direct failures. The zero value disables direct transfer at the first one.
type DirectTransferSessionOptions struct {
	// RetryResponse is called after each unsuccessful direct response. request
	// is the request as sent, after request options, and attempt counts the
	// direct requests sent for it so far. It returns how long to wait and
	// whether to send the request directly again. When RetryResponse is nil or
	// returns false, the response disables direct transfer for the session.
	RetryResponse func(request *http.Request, attempt int, err *DirectTransferResponseError) (time.Duration, bool)
	// OnDisable is called once, by the request or response body whose failure
	// disabled direct transfer, with DirectTransferSessionStats.DisabledReason
	// and the failure.
	OnDisable func(reason string, err error)
	// WrapTransport, when set, is called once, when the session builds its
	// direct client, with the SDK's direct transport, and returns the round
	// tripper the session sends direct requests with. The SDK transport
	// verifies the Agent's TLS identity against info's CA and server name,
	// dials only info's candidate addresses within bounded connect budgets, and
	// times out waiting for response headers. Wrap it to observe, log or adjust
	// requests and keep those protections; a round tripper that does not call
	// it takes them over.
	//
	// Requests reach the round tripper already addressed to the direct URI,
	// with Files.com credentials removed, and the session's client never
	// follows redirects. The direct URI's query carries the transfer token:
	// the session removes it from every *url.Error in a failure's chain,
	// including one from an http.Client the round tripper sends through, so
	// keep the request URL out of any other error the round tripper returns. Return responses with unread bodies: the session
	// tracks body failures and never sends a request again once its response
	// arrived. The round tripper is used by concurrent requests. Close closes
	// the idle connections of the SDK transport, and of the round tripper when
	// it has a CloseIdleConnections method; any other resources it holds belong
	// to the caller. A session with WrapTransport always uses its own client,
	// never a WithDirectTransferClientCache cache from the request context.
	WrapTransport func(direct http.RoundTripper) http.RoundTripper
}

// DirectTransferSessionStats summarizes the direct requests of a session.
type DirectTransferSessionStats struct {
	// Attempts counts the direct requests sent, including retries.
	Attempts int64
	// Successes counts the direct requests that returned a successful response.
	Successes int64
	// DisabledReason is empty while direct transfer is enabled. Otherwise it
	// names the first failure: direct_request_failed when no response arrived,
	// direct_response_failed for an unsuccessful response, or
	// direct_body_read_failed or direct_body_close_failed for a response body.
	DisabledReason string
}

// NewDirectTransferSession returns a session for one transfer.
func NewDirectTransferSession(options DirectTransferSessionOptions) *DirectTransferSession {
	return &DirectTransferSession{options: options}
}

// Send sends a copy of request to info's direct endpoint, as
// WrapDirectTransferOptions does, and the caller must close the body of a
// successful response. A failure while reading or closing that body disables
// direct transfer for the session; the response itself is never sent again.
// A body failure after the request context was canceled, or after its deadline
// passed, leaves the session enabled. To report the response as failed, cancel
// the context with another cause, using context.WithCancelCause.
//
// Send returns ErrDirectTransferUnavailable without sending when config
// disables direct transfers, info is not present, or the session is disabled.
// When the direct request fails, Send disables the session and returns an
// error that wraps ErrDirectTransferUnavailable; the caller can then send
// request over its own route. A request context that ends before the response
// arrives never disables the session, and Send returns the context's error.
// When response options fail, Send closes the response and returns an error
// that wraps ErrDirectTransferResponseStarted and the options' error. That
// request must not be sent again, even when the options' error also wraps
// ErrDirectTransferUnavailable.
func (s *DirectTransferSession) Send(config Config, info DirectConnectionInfo, request *http.Request, opts ...RequestResponseOption) (*http.Response, error) {
	if config.DisableDirectTransfers || !DirectConnectionInfoPresent(info) || s.disabledReason.Load() != nil {
		return nil, ErrDirectTransferUnavailable
	}

	directRequest, err := BuildRequest(request.Clone(request.Context()), opts...)
	if err != nil {
		return nil, err
	}
	ctx := directRequest.Context()
	for attempt := 1; ; attempt++ {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		s.attempts.Add(1)
		response, err := sendDirectTransferRequest(config, info, &s.cache, s.options.WrapTransport, directRequest)
		if err == nil {
			s.successes.Add(1)
			response.Body = &directTransferSessionBody{ReadCloser: response.Body, session: s, ctx: ctx}
			processedResponse, err := BuildResponse(response, opts...)
			if err != nil {
				lib.CloseBody(response)
				return nil, errors.Join(ErrDirectTransferResponseStarted, err)
			}
			return processedResponse, nil
		}
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}

		var responseErr *DirectTransferResponseError
		if !errors.As(err, &responseErr) {
			s.disable(directTransferRequestFailed, err)
			if errors.Is(err, ErrDirectTransferUnavailable) {
				return nil, err
			}
			return nil, fmt.Errorf("%w: %w", ErrDirectTransferUnavailable, err)
		}
		wait, retry := time.Duration(0), false
		if s.options.RetryResponse != nil {
			wait, retry = s.options.RetryResponse(directRequest, attempt, responseErr)
		}
		if !retry {
			s.disable(directTransferResponseFailed, err)
			return nil, err
		}
		if wait > 0 {
			timer := time.NewTimer(wait)
			select {
			case <-ctx.Done():
				timer.Stop()
				return nil, ctx.Err()
			case <-timer.C:
			}
		}
	}
}

// Stats returns the session's current counts and disabled reason.
func (s *DirectTransferSession) Stats() DirectTransferSessionStats {
	stats := DirectTransferSessionStats{
		Attempts:  s.attempts.Load(),
		Successes: s.successes.Load(),
	}
	if reason := s.disabledReason.Load(); reason != nil {
		stats.DisabledReason = *reason
	}
	return stats
}

// Close closes the idle connections of the session's direct client. Call it
// once the transfer's requests and response bodies are finished. Requests whose
// context carries a WithDirectTransferClientCache cache use that cache instead.
func (s *DirectTransferSession) Close() {
	s.cache.close()
}

func (s *DirectTransferSession) disable(reason string, err error) {
	if !s.disabledReason.CompareAndSwap(nil, &reason) {
		return
	}
	if s.options.OnDisable != nil {
		s.options.OnDisable(reason, err)
	}
}

type directTransferSessionBody struct {
	io.ReadCloser
	session *DirectTransferSession
	ctx     context.Context
}

func (b *directTransferSessionBody) Read(p []byte) (int, error) {
	n, err := b.ReadCloser.Read(p)
	if err != nil && !errors.Is(err, io.EOF) {
		b.fail(directTransferBodyReadFailed, err)
	}
	return n, err
}

func (b *directTransferSessionBody) Close() error {
	err := b.ReadCloser.Close()
	if err != nil {
		b.fail(directTransferBodyCloseFailed, err)
	}
	return err
}

// fail disables the session for a body failure unless the caller ended the
// request: its context was canceled, or its deadline passed, without another
// cause. The check uses the caller's context, so a transport's own timeouts
// still count as failures.
func (b *directTransferSessionBody) fail(reason string, err error) {
	if b.ctx.Err() != nil {
		cause := context.Cause(b.ctx)
		if errors.Is(cause, context.Canceled) || errors.Is(cause, context.DeadlineExceeded) {
			return
		}
	}
	b.session.disable(reason, err)
}
