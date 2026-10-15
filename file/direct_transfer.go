package file

import (
	"context"
	"net/http"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/lib"
)

// directTransferSessionContextKey carries the DirectTransferSession that
// Download uses for the requests of one Download V2 transfer.
type directTransferSessionContextKey struct{}

type directTransferBackpressure struct {
	retryAfter time.Duration
	seen       bool
}

type directTransferBackpressureContextKey struct{}

func directTransferBackpressureFromContext(ctx context.Context) *directTransferBackpressure {
	if ctx == nil {
		return nil
	}
	backpressure, _ := ctx.Value(directTransferBackpressureContextKey{}).(*directTransferBackpressure)
	return backpressure
}

func (b *directTransferBackpressure) record(retryAfter time.Duration) {
	if b == nil {
		return
	}
	b.seen = true
	if retryAfter > b.retryAfter {
		b.retryAfter = retryAfter
	}
}

func directDownloadSessionFromContext(ctx context.Context) *files_sdk.DirectTransferSession {
	if ctx == nil {
		return nil
	}
	session, _ := ctx.Value(directTransferSessionContextKey{}).(*files_sdk.DirectTransferSession)
	return session
}

// newDirectDownloadSession returns a session for direct download requests. A
// direct 429 is sent again after its Retry-After, up to downloadV2RetryAttempts
// requests, and recorded as back pressure for the Download V2 range that
// received it.
func newDirectDownloadSession(onDisable func(reason string, err error)) *files_sdk.DirectTransferSession {
	return files_sdk.NewDirectTransferSession(files_sdk.DirectTransferSessionOptions{
		RetryResponse: retryDirectDownloadBackpressure,
		OnDisable:     onDisable,
	})
}

func retryDirectDownloadBackpressure(request *http.Request, attempt int, err *files_sdk.DirectTransferResponseError) (time.Duration, bool) {
	if err.StatusCode != http.StatusTooManyRequests {
		return 0, false
	}
	directTransferBackpressureFromContext(request.Context()).record(err.RetryAfter)
	return err.RetryAfter, attempt < downloadV2RetryAttempts
}

type directTransferUploadStats interface {
	recordDirectAttempt()
	recordDirectSuccess()
	recordDirectFailure()
	recordDirectDisabled()
}

type directTransferUploadOptions struct {
	partNumber int64
	canReplay  func() bool
	observer   func(lib.UploadRetryAttempt)
	stats      directTransferUploadStats
}

func (u *uploadIO) callRawWithDirectTransfer(ctx context.Context, info files_sdk.DirectConnectionInfo, params *files_sdk.CallParams, options directTransferUploadOptions) (*http.Response, bool) {
	if !u.directTransferUploadAvailable(info) {
		return nil, false
	}

	options.stats.recordDirectAttempt()
	directURL, directClient, err := files_sdk.DirectTransferRetryableClient(ctx, u.Config, info)
	if err != nil {
		options.stats.recordDirectFailure()
		if u.disableDirectTransfersForUpload("direct_client_unavailable", err, options.partNumber) {
			options.stats.recordDirectDisabled()
		}
		return nil, false
	}

	directParams := *params
	directParams.Uri = directURL
	directParams.Headers = files_sdk.DirectTransferRequestHeaders(params.Headers)
	directParams.Client = lib.UploadRetryableHttpWithObserver(directClient, options.canReplay, options.observer)

	res, err := files_sdk.CallRaw(&directParams)
	if err == nil {
		err = lib.ResponseErrors(res, files_sdk.APIError(), lib.S3XMLError, lib.NonOkError)
		if err == nil {
			options.stats.recordDirectSuccess()
			return res, true
		}
	}

	lib.CloseBody(res)
	options.stats.recordDirectFailure()
	if u.disableDirectTransfersForUpload("direct_request_failed", err, options.partNumber) {
		options.stats.recordDirectDisabled()
	}
	return nil, false
}
