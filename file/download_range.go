package file

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/internal/expiredurl"
	"github.com/Files-com/files-sdk-go/v3/lib"
)

var errDownloadWholeResponse = errors.New("download has no trustworthy range validator; a whole response is required")
var errDownloadRangeClose = errors.New("download response could not be closed")

type downloadRepresentation struct {
	mu           sync.Mutex
	initializing chan struct{}
	etag         string
	lastModified string
	total        int64
	whole        bool
}

// Dates are a best-effort guard, never a substitute for the strong validator
// required by multipart downloads and persistent resume.
func downloadLastModified(headers http.Header) string {
	values := headers.Values("Last-Modified")
	if len(values) != 1 {
		return ""
	}
	modified, err := http.ParseTime(values[0])
	if err != nil {
		return ""
	}
	return modified.UTC().Format(http.TimeFormat)
}

func downloadDateChanged(baseline string, headers http.Header) bool {
	bound, boundErr := http.ParseTime(baseline)
	modified, err := http.ParseTime(downloadLastModified(headers))
	return boundErr == nil && err == nil && modified.After(bound)
}

func strongDownloadETag(value string) string {
	value = strings.TrimSpace(value)
	if len(value) < 2 || value[0] != '"' || value[len(value)-1] != '"' {
		return ""
	}
	for _, b := range []byte(value[1 : len(value)-1]) {
		if b == '"' || b < 0x21 || b == 0x7f {
			return ""
		}
	}
	return value
}

func (s *downloadRepresentation) begin(ctx context.Context, allowSingle bool) (leader bool, etag, lastModified string, total int64, err error) {
	for {
		s.mu.Lock()
		if s.whole {
			lastModified = s.lastModified
			s.mu.Unlock()
			if allowSingle {
				return false, "", lastModified, 0, nil
			}
			return false, "", lastModified, 0, errDownloadWholeResponse
		}
		if s.etag != "" {
			etag, total = s.etag, s.total
			s.mu.Unlock()
			return
		}
		if s.initializing == nil {
			s.initializing = make(chan struct{})
			lastModified = s.lastModified
			s.mu.Unlock()
			return true, "", lastModified, 0, nil
		}
		ready := s.initializing
		s.mu.Unlock()
		select {
		case <-ready:
		case <-ctx.Done():
			return false, "", "", 0, ctx.Err()
		}
	}
}

func (s *downloadRepresentation) finish(etag, lastModified string, total int64, err error) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.lastModified == "" && (err == nil || errors.Is(err, errDownloadWholeResponse)) {
		s.lastModified = lastModified
	}
	if err == nil {
		if etag == "" && s.etag == "" {
			s.whole, s.total = true, total
		} else if s.etag != "" && etag != "" && s.etag != etag {
			err = downloadVersionChanged("concurrent response representation changed")
		} else if s.etag != "" && (etag == "" || s.total != total) {
			err = downloadRangeConflict("concurrent response validator or total was not preserved")
		} else if s.whole {
			err = errDownloadWholeResponse
		} else {
			s.etag, s.total = etag, total
		}
	}
	if errors.Is(err, errDownloadWholeResponse) {
		if s.etag != "" {
			err = errors.Join(downloadRangeConflict("range response lost its bound validator"), err)
		} else {
			s.whole = true
		}
	}
	close(s.initializing)
	s.initializing = nil
	return err
}

func (s *downloadRepresentation) snapshot() (etag string, total int64, whole bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.etag, s.total, s.whole
}

func downloadVersionChanged(reason string) error {
	return fmt.Errorf("%w: %s", files_sdk.ErrDownloadSourceChanged, reason)
}

func downloadRangeConflict(reason string) error {
	return downloadV2TerminalError{fmt.Errorf("%w: %s", errDownloadSizeConflict, reason)}
}

// This wrapper owns the same attempt allowance formerly used by the engine's
// full-part retries. Its next request starts after bytes already returned.
type resumableRangeReader struct {
	file             *File
	ctx              context.Context
	budget           *downloadReadBudget
	next, end, total int64
	etag             string
	lastModified     string
	body             io.ReadCloser
	pending          error
	terminal         error
	closed           bool
	requireVersion   bool
	noProof          bool
	trust            SizeTrust
}

func (f *File) openDownloadRange(ctx context.Context, off, end int64, attempts int) (io.ReadCloser, downloadV2Range, error) {
	return f.openRange(ctx, off, end, attempts, true)
}

func (f *File) openRange(ctx context.Context, off, end int64, attempts int, requireVersion bool) (io.ReadCloser, downloadV2Range, error) {
	if off < 0 || end < off {
		return nil, downloadV2Range{}, errors.New("invalid download byte range")
	}
	ctx, budget := withDownloadBudget(ctx)
	budget.mu.Lock()
	if !budget.limited {
		budget.limited, budget.remaining = true, attempts
	}
	budget.mu.Unlock()
	f.fileMutex.Lock()
	representation := f.representation
	f.fileMutex.Unlock()
	leader, etag, lastModified, total, err := representation.begin(ctx, !requireVersion)
	if err != nil {
		return nil, downloadV2Range{}, err
	}
	r := &resumableRangeReader{file: f, ctx: ctx, budget: budget, next: off, end: end, etag: etag, lastModified: lastModified, total: total, requireVersion: requireVersion}
	err = r.open()
	if leader {
		err = representation.finish(r.etag, r.lastModified, r.total, err)
	}
	if err != nil {
		return nil, downloadV2Range{}, errors.Join(err, r.closeBody())
	}
	return r, downloadV2Range{off: off, end: r.end, total: r.total, trust: r.trust}, nil
}

// Error parsers own their bodies and may replace Body after closing it. Keep
// the original close result observable even through that replacement.
type rangeResponseBody struct {
	io.ReadCloser
	closeOnce sync.Once
	closeErr  error
}

func (b *rangeResponseBody) Close() error {
	b.closeOnce.Do(func() {
		if b.ReadCloser != nil {
			b.closeErr = b.ReadCloser.Close()
		}
	})
	return b.closeErr
}

func discardRangeResponse(response *http.Response, body *rangeResponseBody, cause error) error {
	closeErr := body.Close()
	response.Body = http.NoBody
	if closeErr != nil {
		// A status may permit renewal/retry; failure to close that response
		// does not. Keep both causes and reuse the reader's terminal marker.
		return errors.Join(cause, rangeCloseFailure(closeErr))
	}
	return cause
}

func rangeCloseFailure(err error) error {
	if err == nil {
		return nil
	}
	// The distinct close cause also keeps mode/pause classifiers from
	// mistaking an EOF/cancellation returned by Close for ordinary completion.
	return downloadV2TerminalError{errors.Join(errDownloadRangeClose, err)}
}

func (r *resumableRangeReader) open() error {
	var last error
	for r.budget.canAttempt() {
		if err := r.ctx.Err(); err != nil {
			return err
		}
		r.file.fileMutex.Lock()
		file := *r.file.File
		r.file.fileMutex.Unlock()
		headers := &http.Header{}
		headers.Set("Range", fmt.Sprintf("bytes=%d-%d", r.next, r.end))
		if r.etag != "" {
			headers.Set("If-Match", r.etag)
		} else if r.file.checkDownloadModificationTime && r.lastModified != "" {
			headers.Set("If-Unmodified-Since", r.lastModified)
		}
		var maxConnections int
		var requestID string
		updated, err := (&Client{Config: r.file.Config}).Download(files_sdk.FileDownloadParams{File: file},
			files_sdk.WithContext(r.ctx), files_sdk.RequestHeadersOption(headers),
			files_sdk.ResponseOption(func(response *http.Response) error {
				body := &rangeResponseBody{ReadCloser: response.Body}
				response.Body = body
				reject := func(cause error) error { return discardRangeResponse(response, body, cause) }
				if r.file.checkDownloadModificationTime && r.etag == "" && (response.StatusCode == http.StatusOK || response.StatusCode == http.StatusPartialContent) {
					if downloadDateChanged(r.lastModified, response.Header) {
						return reject(downloadVersionChanged("download modification time changed"))
					}
					if r.lastModified == "" {
						r.lastModified = downloadLastModified(response.Header)
					}
				}
				if response.StatusCode == http.StatusOK {
					if r.etag == "" {
						return reject(errDownloadWholeResponse)
					}
					if etag := strongDownloadETag(response.Header.Get("ETag")); etag != "" && etag != r.etag {
						return reject(downloadVersionChanged("range representation changed"))
					}
					return reject(downloadRangeConflict("bound range was answered with a whole response"))
				}
				if response.StatusCode == http.StatusPreconditionFailed {
					return reject(downloadVersionChanged("range precondition failed"))
				}
				if err := expiredurl.ResponseErrors(response, lib.IsStatus(http.StatusForbidden), files_sdk.APIError(), lib.NotStatus(http.StatusPartialContent)); err != nil {
					return reject(err)
				}
				coverage, valid := parseContentRange(response.Header.Get("Content-Range"))
				etag := strongDownloadETag(response.Header.Get("ETag"))
				if r.etag != "" && etag != "" && etag != r.etag {
					return reject(downloadVersionChanged("range representation changed"))
				}
				// A managed stage can discard an initial response that identifies
				// no bytes at all. Positive but wrong coordinates, public reads,
				// and every response after binding remain strict conflicts.
				if r.requireVersion && r.etag == "" && (!valid || coverage.start < 0 || coverage.end < 0) {
					return reject(errDownloadWholeResponse)
				}
				expectedEnd := r.end
				if coverage.totalKnown {
					expectedEnd = min(expectedEnd, coverage.total-1)
				}
				if !valid || coverage.start != r.next || coverage.end != expectedEnd || coverage.totalKnown && coverage.total <= r.next {
					return reject(downloadRangeConflict("download range response did not cover the requested bytes and total"))
				}
				if r.etag == "" && (etag == "" || !coverage.totalKnown) && r.requireVersion {
					return reject(errDownloadWholeResponse)
				}
				if r.etag != "" && (etag != r.etag || !coverage.totalKnown || coverage.total != r.total) {
					return reject(downloadRangeConflict("bound range validator or total was not preserved"))
				}
				if response.ContentLength >= 0 && response.ContentLength != coverage.end-coverage.start+1 {
					return reject(downloadRangeConflict("download range content length did not match its coverage"))
				}
				r.etag, r.total, r.end, r.trust = etag, coverage.total, coverage.end, TrustedSizeValue
				r.noProof = etag == "" || !coverage.totalKnown
				if !coverage.totalKnown {
					r.etag, r.total, r.trust = "", -1, UntrustedSizeValue
				}
				maxConnections = parseMaxConnections(response)
				requestID = response.Header.Get("X-Files-Download-Request-Id")
				r.body = response.Body
				return nil
			}))
		r.file.fileMutex.Lock()
		*r.file.File = updated
		if err == nil {
			r.file.Size, r.file.SizeTrust = r.total, r.trust
			r.file.MaxConnections = maxConnections
			r.file.downloadRequestId = requestID
		}
		r.file.fileMutex.Unlock()
		if err == nil {
			if r.trust == UntrustedSizeValue && requestID != "" {
				r.body = ReaderCloserDownloadStatus{ReadWrapper: &ReadWrapper{ReadCloser: r.body}, file: r.file, expectedSize: r.end - r.next + 1, rangeRequest: true}
			}
			return nil
		}
		last = err
		// Inspect the joined result before looking for an expiry/status member.
		// A rejected response with a failed close must never become success on
		// a later URL or ordinary open retry.
		if downloadV2TerminalPartError(err) {
			return err
		}
		if downloadRequestExpired(err) && r.budget.renew() {
			r.file.fileMutex.Lock()
			r.file.DownloadUri = ""
			r.file.fileMutex.Unlock()
			continue
		}
		if errors.Is(err, errDownloadWholeResponse) || errors.Is(err, errDownloadFrontend) || errors.Is(err, errDownloadAttemptLimit) || downloadRequestExpired(err) || errors.Is(err, errDownloadURLExpired) {
			return err
		}
	}
	if last == nil {
		last = errDownloadAttemptLimit
	}
	return last
}

func (r *resumableRangeReader) closeBody() error {
	if r.body == nil {
		return nil
	}
	body := r.body
	r.body = nil
	return rangeCloseFailure(body.Close())
}

func (r *resumableRangeReader) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	if r.closed {
		return 0, io.ErrClosedPipe
	}
	for {
		if r.terminal != nil {
			return 0, r.terminal
		}
		if r.next > r.end {
			return 0, io.EOF
		}
		if r.pending != nil {
			readErr := r.pending
			r.pending = nil
			if err := r.closeBody(); err != nil {
				r.terminal = errors.Join(readErr, err)
				return 0, r.terminal
			}
			if r.noProof || !r.budget.canAttempt() || r.ctx.Err() != nil {
				r.terminal = errors.Join(readErr, r.ctx.Err())
				return 0, r.terminal
			}
			r.file.Config.LogPath(r.file.safeFile().Path, map[string]interface{}{"message": "resuming interrupted download range", "offset": r.next, "end": r.end, "attempts_remaining": r.budget.attemptsRemaining()})
			if err := r.open(); err != nil {
				r.terminal = err
				return 0, err
			}
		}
		want := min(int64(len(p)), r.end-r.next+1)
		n, err := r.body.Read(p[:want])
		r.next += int64(n)
		if r.next > r.end {
			// Confirm termination before a range is declared complete. Close also
			// retains this error when io.ReadFull consumes the final bytes.
			for err == nil {
				if r.ctx.Err() != nil {
					err = r.ctx.Err()
					break
				}
				var probe [1]byte
				more, probeErr := r.body.Read(probe[:])
				if more != 0 {
					err = errors.New("download range exceeded its declared length")
				} else {
					err = probeErr
				}
			}
			if err == io.EOF {
				err = nil
			}
			r.terminal = errors.Join(err, r.closeBody())
			if r.terminal != nil {
				// The final read is admitted only after termination and close
				// succeed. A paused prefix must not bypass that check later.
				r.next -= int64(n)
				return 0, r.terminal
			}
			return n, r.terminal
		}
		if err != nil {
			if errors.Is(err, io.EOF) {
				err = io.ErrUnexpectedEOF
			}
			r.pending = err
			if n == 0 {
				continue
			}
		}
		return n, nil
	}
}

func (r *resumableRangeReader) Close() error {
	if r.closed {
		return r.terminal
	}
	r.closed = true
	r.terminal = errors.Join(r.terminal, r.closeBody())
	return r.terminal
}
