package file

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/downloadurl"
	"github.com/Files-com/files-sdk-go/v3/internal/transferguard"
	"github.com/hashicorp/go-retryablehttp"
)

var errDownloadURLExpired = errors.New("download URL expired before it could be sent")
var errDownloadFrontend = errors.New("download response came from the Files.com web frontend")
var errDownloadAttemptLimit = errors.New("download range retry allowance exhausted")

func downloadResponseGuard(response *http.Response) error {
	for name := range response.Header {
		if strings.EqualFold(name, "X-Files-Frontend-App") {
			return errDownloadFrontend
		}
	}
	return nil
}

// A logical read shares one renewal across initial expiry, HTTP retries and
// partial-body continuation. Initial acquisition of a missing URL is not renewal.
type downloadReadBudget struct {
	mu        sync.Mutex
	renewed   bool
	limited   bool
	remaining int
}

func (b *downloadReadBudget) takeAttempt() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if !b.limited {
		return true
	}
	if b.remaining == 0 {
		return false
	}
	b.remaining--
	return true
}

func (b *downloadReadBudget) canAttempt() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return !b.limited || b.remaining > 0
}

func (b *downloadReadBudget) attemptsRemaining() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.remaining
}

func (b *downloadReadBudget) renew() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.renewed {
		return false
	}
	b.renewed = true
	return true
}

type downloadBudgetKey struct{}
type downloadRequestState struct {
	source          string
	mu              sync.Mutex
	expired         bool
	frontend        bool
	responseStarted bool
}

func withDownloadBudget(ctx context.Context) (context.Context, *downloadReadBudget) {
	if ctx == nil {
		ctx = context.Background()
	}
	if b, ok := ctx.Value(downloadBudgetKey{}).(*downloadReadBudget); ok {
		return ctx, b
	}
	b := &downloadReadBudget{}
	return context.WithValue(ctx, downloadBudgetKey{}, b), b
}

func knownDownloadExpired(uri string, within time.Duration) bool {
	u, err := downloadurl.New(uri)
	if err != nil || u.Time.IsZero() {
		return false
	}
	_, valid := u.Valid(within)
	return !valid
}

func (s *downloadRequestState) CheckRequest(req *http.Request) error {
	if knownDownloadExpired(s.source, 0) || knownDownloadExpired(req.URL.String(), 0) {
		s.mu.Lock()
		s.expired = true
		s.mu.Unlock()
		return errDownloadURLExpired
	}
	return nil
}

func (s *downloadRequestState) CheckResponse(response *http.Response) error {
	if err := downloadResponseGuard(response); err != nil {
		s.mu.Lock()
		s.frontend = true
		s.mu.Unlock()
		return err
	}
	return nil
}

// Copy fields, not retryablehttp's internal sync.Once values. Only this transfer
// sees the transport wrapper; the configured client, hooks and pool remain owned
// by their caller, including explicitly supplied HTTP clients.
func downloadTransferConfig(config files_sdk.Config, state *downloadRequestState) files_sdk.Config {
	config = config.Init()
	base := config.Client
	httpClient := http.Client{}
	if base.HTTPClient != nil {
		httpClient = *base.HTTPClient
	}
	httpClient.Transport = transferguard.Wrap(httpClient.Transport)
	check := base.CheckRetry
	if check == nil {
		check = retryablehttp.DefaultRetryPolicy
	}
	config.Client = &retryablehttp.Client{
		HTTPClient: &httpClient, Logger: base.Logger,
		RetryWaitMin: base.RetryWaitMin, RetryWaitMax: base.RetryWaitMax, RetryMax: base.RetryMax,
		RequestLogHook: base.RequestLogHook, ResponseLogHook: base.ResponseLogHook,
		Backoff: base.Backoff, ErrorHandler: base.ErrorHandler, PrepareRetry: base.PrepareRetry,
		CheckRetry: func(ctx context.Context, response *http.Response, err error) (bool, error) {
			retry, policyErr := check(ctx, response, err)
			if state.wasFrontend() {
				return false, errDownloadFrontend
			}
			if state.wasExpired() {
				return false, errDownloadURLExpired
			}
			if errors.Is(err, errDownloadURLExpired) {
				return false, errDownloadURLExpired
			}
			return retry, policyErr
		},
	}
	if base.ErrorHandler != nil {
		config.Client.ErrorHandler = func(response *http.Response, err error, attempts int) (*http.Response, error) {
			result, handlerErr := base.ErrorHandler(response, err, attempts)
			if state.wasFrontend() || state.wasExpired() {
				if result != nil && result.Body != nil {
					_ = result.Body.Close()
				}
				if state.wasExpired() && !state.wasFrontend() {
					return nil, errDownloadURLExpired
				}
				return nil, errDownloadFrontend
			}
			return result, handlerErr
		}
	}
	return config
}

func (s *downloadRequestState) wasExpired() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.expired
}

func (s *downloadRequestState) wasFrontend() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.frontend
}

func (s *downloadRequestState) startResponse() {
	s.mu.Lock()
	s.responseStarted = true
	s.mu.Unlock()
}

func (s *downloadRequestState) hasStartedResponse() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.responseStarted
}
