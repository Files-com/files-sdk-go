// Package transferguard composes SDK-owned transfer checks without giving a
// request ownership of its borrowed connection pool.
package transferguard

import (
	"context"
	"net/http"
)

// Guard belongs to one logical request, not to the transport/cache. A rejected
// response is closed by this package before any higher response observer sees it.
type Guard interface {
	CheckRequest(*http.Request) error
	CheckResponse(*http.Response) error
}

type contextKey struct{}

func WithContext(ctx context.Context, guard Guard) context.Context {
	return context.WithValue(ctx, contextKey{}, guard)
}

type transport struct{ base http.RoundTripper }

// Wrap is stateless and may be shared by a cached client. It deliberately does
// not forward CloseIdleConnections: the pool's existing owner retains cleanup.
func Wrap(base http.RoundTripper) http.RoundTripper {
	if base == nil {
		base = http.DefaultTransport
	}
	return transport{base: base}
}

func (t transport) RoundTrip(request *http.Request) (*http.Response, error) {
	guard, present := request.Context().Value(contextKey{}).(Guard)
	if !present {
		return t.base.RoundTrip(request)
	}
	if err := guard.CheckRequest(request); err != nil {
		return nil, err
	}
	response, err := t.base.RoundTrip(request)
	if response != nil {
		if rejected := guard.CheckResponse(response); rejected != nil {
			if response.Body != nil {
				_ = response.Body.Close()
			}
			return nil, rejected
		}
	}
	return response, err
}
