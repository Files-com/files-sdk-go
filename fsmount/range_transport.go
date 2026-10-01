//go:build linux || windows

package fsmount

import (
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache"
	"github.com/Files-com/files-sdk-go/v3/lib"
)

var (
	errRangeValidatorUnavailable = errors.New("remote range has no strong ETag")
	errRangeDownloadUnsupported  = errors.New("remote backend does not support ranged downloads")
	errRangeNotSatisfiable       = errors.New("remote range is not satisfiable")
	errMalformedRangeResponse    = errors.New("malformed remote range response")
	errShortRangeResponse        = errors.New("remote range response ended early")
	errLongRangeResponse         = errors.New("remote range response exceeded its declared length")
)

// remoteRangeResponse describes the bytes returned for one remote range request.
// TotalSize is -1 when the response does not supply it.
type remoteRangeResponse struct {
	ETag      string
	File      files_sdk.File
	Returned  cache.ByteRange
	TotalSize int64
	Partial   bool
	Body      io.ReadCloser
}

func (b *sdkRemoteBackend) downloadRange(params files_sdk.FileDownloadParams, requested cache.ByteRange, opts ...files_sdk.RequestResponseOption) (remoteRangeResponse, error) {
	if err := requested.Validate(); err != nil {
		return remoteRangeResponse{}, err
	}
	if requested.Start == requested.End {
		return remoteRangeResponse{}, fmt.Errorf("%w: empty requested range [%d, %d)", cache.ErrInvalidByteRange, requested.Start, requested.End)
	}

	for attempt := 0; attempt < 2; attempt++ {
		result := remoteRangeResponse{TotalSize: -1}
		headers := &http.Header{}
		headers.Set("Range", fmt.Sprintf("bytes=%d-%d", requested.Start, requested.End-1))
		rangeOpts := append([]files_sdk.RequestResponseOption(nil), opts...)
		rangeOpts = append(
			rangeOpts,
			files_sdk.RequestHeadersOption(headers),
			files_sdk.ResponseOption(func(response *http.Response) error {
				return captureRangeResponse(&result, params.File.Size, requested, response)
			}),
		)

		file, err := b.fileClient.Download(params, rangeOpts...)
		result.File = file
		if err == nil {
			return result, nil
		}
		if result.Body != nil {
			_ = result.Body.Close()
		}
		if attempt == 0 && (rangeDownloadExpired(err) || errors.Is(err, errRangeNotSatisfiable)) {
			params.File = file
			params.File.DownloadUri = ""
			continue
		}
		return remoteRangeResponse{}, err
	}
	return remoteRangeResponse{}, errors.New("remote range download retry exhausted")
}

func (b *providerRemoteBackend) downloadRange(_ files_sdk.FileDownloadParams, _ cache.ByteRange, _ ...files_sdk.RequestResponseOption) (remoteRangeResponse, error) {
	return remoteRangeResponse{}, errRangeDownloadUnsupported
}

func captureRangeResponse(result *remoteRangeResponse, knownSize int64, requested cache.ByteRange, response *http.Response) error {
	result.ETag = strongETag(response.Header.Get("ETag"))
	switch response.StatusCode {
	case http.StatusPartialContent:
		returned, totalSize, err := parseContentRange(response.Header.Get("Content-Range"))
		if err != nil {
			lib.CloseBody(response)
			return err
		}
		if returned != requested {
			lib.CloseBody(response)
			return fmt.Errorf("%w: returned range [%d, %d), requested [%d, %d)", errMalformedRangeResponse, returned.Start, returned.End, requested.Start, requested.End)
		}
		expected := returned.End - returned.Start
		if response.ContentLength >= 0 && response.ContentLength != expected {
			lib.CloseBody(response)
			return fmt.Errorf("%w: content length %d, expected %d", errMalformedRangeResponse, response.ContentLength, expected)
		}
		result.Returned = returned
		result.TotalSize = totalSize
		result.Partial = true
		result.Body = newExactLengthReadCloser(response.Body, expected)
		return nil
	case http.StatusOK:
		totalSize := response.ContentLength
		if totalSize < 0 {
			totalSize = knownSize
		}
		if totalSize < 0 {
			lib.CloseBody(response)
			return fmt.Errorf("%w: complete response did not identify its length", errMalformedRangeResponse)
		}
		result.Returned = cache.ByteRange{Start: 0, End: totalSize}
		result.TotalSize = totalSize
		result.Partial = false
		result.Body = newExactLengthReadCloser(response.Body, totalSize)
		return nil
	case http.StatusPreconditionFailed:
		lib.CloseBody(response)
		return errRangeVersionChanged
	case http.StatusRequestedRangeNotSatisfiable:
		lib.CloseBody(response)
		return fmt.Errorf("%w: %s", errRangeNotSatisfiable, response.Header.Get("Content-Range"))
	default:
		if err := files_sdk.APIError()(response); err != nil {
			return err
		}
		if err := lib.NonOkError(response); err != nil {
			return err
		}
		lib.CloseBody(response)
		return fmt.Errorf("%w: unexpected HTTP status %d", errMalformedRangeResponse, response.StatusCode)
	}
}

// Weak or malformed validators cannot prove that separate responses contain
// the same bytes. RFC 9110 entity tags are quoted opaque values.
func strongETag(value string) string {
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

func parseContentRange(value string) (cache.ByteRange, int64, error) {
	unit, remainder, ok := strings.Cut(strings.TrimSpace(value), " ")
	if !ok || !strings.EqualFold(unit, "bytes") {
		return cache.ByteRange{}, -1, fmt.Errorf("%w: invalid Content-Range %q", errMalformedRangeResponse, value)
	}
	rangeValue, totalValue, ok := strings.Cut(remainder, "/")
	if !ok {
		return cache.ByteRange{}, -1, fmt.Errorf("%w: invalid Content-Range %q", errMalformedRangeResponse, value)
	}
	startValue, endValue, ok := strings.Cut(rangeValue, "-")
	if !ok {
		return cache.ByteRange{}, -1, fmt.Errorf("%w: invalid Content-Range %q", errMalformedRangeResponse, value)
	}
	start, startErr := strconv.ParseInt(strings.TrimSpace(startValue), 10, 64)
	endInclusive, endErr := strconv.ParseInt(strings.TrimSpace(endValue), 10, 64)
	if startErr != nil || endErr != nil || start < 0 || endInclusive < start || endInclusive == math.MaxInt64 {
		return cache.ByteRange{}, -1, fmt.Errorf("%w: invalid Content-Range %q", errMalformedRangeResponse, value)
	}
	returned := cache.ByteRange{Start: start, End: endInclusive + 1}

	totalSize := int64(-1)
	totalValue = strings.TrimSpace(totalValue)
	if totalValue != "*" {
		parsedTotal, err := strconv.ParseInt(totalValue, 10, 64)
		if err != nil || parsedTotal < returned.End {
			return cache.ByteRange{}, -1, fmt.Errorf("%w: invalid Content-Range %q", errMalformedRangeResponse, value)
		}
		totalSize = parsedTotal
	}
	return returned, totalSize, nil
}

func rangeDownloadExpired(err error) bool {
	if files_sdk.IsExpired(err) {
		return true
	}
	var responseErr lib.ResponseError
	return errors.As(err, &responseErr) && responseErr.StatusCode == http.StatusForbidden
}

type exactLengthReadCloser struct {
	body     io.ReadCloser
	expected int64
	read     int64
	done     bool
}

func newExactLengthReadCloser(body io.ReadCloser, expected int64) io.ReadCloser {
	return &exactLengthReadCloser{body: body, expected: expected}
}

func (r *exactLengthReadCloser) Read(buffer []byte) (int, error) {
	if r.done {
		return 0, io.EOF
	}
	remaining := r.expected - r.read
	if remaining > 0 {
		if int64(len(buffer)) > remaining {
			buffer = buffer[:remaining]
		}
		n, err := r.body.Read(buffer)
		r.read += int64(n)
		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
				if r.read != r.expected {
					return n, fmt.Errorf("%w: received %d bytes, expected %d", errShortRangeResponse, r.read, r.expected)
				}
				r.done = true
			}
			return n, err
		}
		return n, nil
	}

	probe := []byte{0}
	n, err := r.body.Read(probe)
	if n > 0 {
		r.done = true
		return 0, fmt.Errorf("%w: expected %d bytes", errLongRangeResponse, r.expected)
	}
	if errors.Is(err, io.EOF) {
		r.done = true
	}
	return 0, err
}

func (r *exactLengthReadCloser) Close() error {
	return r.body.Close()
}
