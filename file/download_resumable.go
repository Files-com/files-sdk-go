package file

import "errors"

// ErrJobPaused is the cancel cause set when a job is paused rather than fully canceled.
// Use errors.Is(context.Cause(ctx), ErrJobPaused) to distinguish pause from cancel.
var ErrJobPaused = errors.New("job paused")

type JobDownloadCheckpoint struct {
	CompletedPaths []string
}
