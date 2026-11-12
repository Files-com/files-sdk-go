package file

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"net/http"
	"strconv"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/lib"
)

// uploadFinalizationPollDelay is the wait between the pendingAnswers-th answer that Files.com is still finishing an
// upload and the next identical request: a second after the first such answer, two after the second and three after
// every later one. Each upload's finalization counts its own answers from one.
func uploadFinalizationPollDelay(pendingAnswers int) time.Duration {
	return time.Duration(min(pendingAnswers, 3)) * time.Second
}

// UploadCompletionUnknownError is the error of an upload whose request to finish it (action=end) Files.com accepted,
// answering that it was still finishing the upload (HTTP 202), when the client stopped before it learned the outcome:
// the upload's context ended, the same request sent again kept failing without an answer from Files.com, or an answer
// could not be read. Files.com may still complete the upload. Err is why the client stopped, so errors.Is finds a
// context's context.Canceled or context.DeadlineExceeded. Resuming the upload with the UploadResumable returned with
// this error sends the same request again when the UploadResumable kept the upload's parts, as it does for a file or
// an io.ReaderAt of known size.
type UploadCompletionUnknownError struct {
	// Path is the upload's destination.
	Path string
	// Ref is the upload's reference.
	Ref string
	// Err is why the client stopped waiting.
	Err error
}

func (e *UploadCompletionUnknownError) Error() string {
	return fmt.Sprintf("stopped waiting for Files.com to finish the upload of %v (ref %v), which Files.com may still complete: %v", e.Path, e.Ref, e.Err)
}

func (e *UploadCompletionUnknownError) Unwrap() error {
	return e.Err
}

// completeUpload sends the request that finishes the upload, action=end with async_finalize, and returns the File
// Files.com answers with. Files.com can instead answer that it is still finishing the upload (HTTP 202 with status
// pending, the upload's ref and an optional progress percentage). The same request is then sent again after each such
// answer, waiting one, two, then three seconds (uploadFinalizationPollDelay), until a File or an error ends it, with the
// Config the first request was sent with. sendCtx is the first request's context; waitCtx stops the waits and the
// requests after them.
func (u *uploadIO) completeUpload(sendCtx context.Context, waitCtx context.Context, providedMtime *time.Time, etags []files_sdk.EtagsParam, bytesWritten int64, path string, ref string) (files_sdk.File, error) {
	if providedMtime != nil && providedMtime.IsZero() {
		providedMtime = nil
	}
	params := files_sdk.FileCreateParams{
		ProvidedMtime:    providedMtime,
		EtagsParam:       etags,
		Action:           "end",
		Path:             path,
		Ref:              ref,
		Size:             bytesWritten,
		MkdirParents:     lib.Bool(true),
		ActionAttributes: maps.Clone(u.actionAttributes),
		AsyncFinalize:    lib.Bool(true),
	}
	unknown := func(err error) error {
		return &UploadCompletionUnknownError{Path: path, Ref: ref, Err: err}
	}

	// Every request is sent with the first one's credentials, workspace, endpoint, headers and HTTP client, even if
	// the Client's Config changes meanwhile, as an upload job's reporter can change it.
	config := uploadFinalizationConfig(u.Config)
	answer, err := sendUploadFinalization(sendCtx, config, params)
	if err != nil {
		if answer.status == http.StatusAccepted {
			return files_sdk.File{}, unknown(err)
		}
		// A first request that fails, or whose answer is lost, is handled as it always was.
		return answer.file, err
	}
	for pendingAnswers := 1; answer.status == http.StatusAccepted; pendingAnswers++ {
		percent, err := pendingUploadFinalization(answer.body, ref)
		if err != nil {
			return files_sdk.File{}, unknown(err)
		}
		logged := map[string]any{"timestamp": time.Now(), "event": "upload finalization pending"}
		if percent != nil {
			logged["progress"] = *percent
		}
		u.LogPath(path, logged)
		if u.finalizationProgress != nil {
			u.finalizationProgress(percent)
		}
		if err := waitForUploadFinalization(waitCtx, uploadFinalizationPollDelay(pendingAnswers)); err != nil {
			return files_sdk.File{}, unknown(err)
		}
		answer, err = sendUploadFinalization(waitCtx, config, params)
		if err != nil {
			if filesAnsweredUploadFinalization(answer.status) {
				// Files.com's error is the upload's outcome.
				return files_sdk.File{}, err
			}
			if stop := uploadContextStop(waitCtx); stop != nil {
				err = stop
			}
			return files_sdk.File{}, unknown(err)
		}
	}
	return answer.file, nil
}

// uploadFinalizationAnswer is the answer to a request that finishes an upload: its HTTP status, set before the body is
// decoded, and the File, or for HTTP 202 the body saying that Files.com is still finishing the upload.
type uploadFinalizationAnswer struct {
	status int
	file   files_sdk.File
	body   []byte
}

func (a *uploadFinalizationAnswer) UnmarshalJSON(data []byte) error {
	if a.status == http.StatusAccepted {
		a.body = append([]byte(nil), data...)
		return nil
	}
	return a.file.UnmarshalJSON(data)
}

// uploadFinalizationConfig copies config for the requests that finish one upload, so that what the first request sends
// does not change for the requests sent again: its own AdditionalHeaders, and the values a request would otherwise
// resolve when it is built. Those are FILES_API_KEY, which counts only when config has neither an API key nor a session
// ID, so the first request's choice between the two holds; the endpoint its override, environment and subdomain
// give; and the default retry client, over lib.DefaultClient, and user agent of a config without them.
func uploadFinalizationConfig(config files_sdk.Config) files_sdk.Config {
	config.AdditionalHeaders = maps.Clone(config.AdditionalHeaders)
	if config.APIKey == "" && config.SessionId == "" {
		config.APIKey = config.GetAPIKey()
	}
	config.EndpointOverride = config.Endpoint()
	initialized := config.Init()
	config.Client = initialized.Client
	config.UserAgent = initialized.UserAgent
	return config
}

// sendUploadFinalization sends params with config as file.Client.Create does.
func sendUploadFinalization(ctx context.Context, config files_sdk.Config, params files_sdk.FileCreateParams) (uploadFinalizationAnswer, error) {
	answer := &uploadFinalizationAnswer{}
	err := files_sdk.Resource(
		config,
		lib.Resource{Method: "POST", Path: "/files/{path}", Params: params, Entity: answer},
		files_sdk.WithContext(ctx),
		files_sdk.ResponseOption(func(res *http.Response) error {
			answer.status = res.StatusCode
			return nil
		}),
	)
	return *answer, err
}

// pendingUploadFinalization reads the body of an HTTP 202 answer, which must be a JSON object with the status pending
// and the upload's ref, and returns its progress: a JSON integer from 0 to 100, or nil when it has none or another
// value.
func pendingUploadFinalization(body []byte, ref string) (*int64, error) {
	var reply map[string]json.RawMessage
	if err := json.Unmarshal(body, &reply); err != nil || reply == nil {
		return nil, errors.New("Files.com answered HTTP 202 without a JSON object")
	}
	var replyStatus, replyRef string
	if err := json.Unmarshal(reply["status"], &replyStatus); err != nil || replyStatus != "pending" {
		return nil, fmt.Errorf("Files.com answered HTTP 202 with the status %s, not pending", reply["status"])
	}
	if err := json.Unmarshal(reply["ref"], &replyRef); err != nil || replyRef != ref {
		return nil, fmt.Errorf("Files.com answered HTTP 202 for the ref %s, not this upload's", reply["ref"])
	}
	return uploadFinalizationPercent(reply["progress"]), nil
}

func uploadFinalizationPercent(raw json.RawMessage) *int64 {
	percent, err := strconv.ParseInt(string(bytes.TrimSpace(raw)), 10, 64)
	if err != nil || percent < 0 || percent > 100 {
		return nil
	}
	return &percent
}

// filesAnsweredUploadFinalization reports whether status is Files.com's error answer, the outcome of an upload whose
// finalization it was waiting for. A server error or a rate limit (429), which the request's retries could not get
// past, is not an outcome.
func filesAnsweredUploadFinalization(status int) bool {
	return status >= 400 && status < 500 && status != http.StatusTooManyRequests
}

func waitForUploadFinalization(ctx context.Context, delay time.Duration) error {
	if stop := uploadContextStop(ctx); stop != nil {
		return stop
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return uploadContextStop(ctx)
	case <-timer.C:
		return nil
	}
}

// uploadContextStop returns ctx's error with its cause, such as ErrJobPaused, or nil while ctx has not ended.
func uploadContextStop(ctx context.Context) error {
	err := ctx.Err()
	if err == nil {
		return nil
	}
	if cause := context.Cause(ctx); cause != nil && !errors.Is(cause, err) {
		return fmt.Errorf("%w: %w", err, cause)
	}
	return err
}
