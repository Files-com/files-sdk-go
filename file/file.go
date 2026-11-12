package file

import (
	"encoding/json"
	"sync"
	"time"

	filesSDK "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file/status"
)

type JobFile struct {
	StatusName          string           `json:"status"`
	TransferBytes       int64            `json:"transferred_bytes"`
	Size                int64            `json:"size_bytes"`
	LocalPath           string           `json:"local_path"`
	RemotePath          string           `json:"remote_path"`
	CheckpointResumable UploadResumable  `json:"-"`
	TmpPath             string           `json:"-"`
	DownloadBinding     *DownloadBinding `json:"-"`
	EndedAt             time.Time        `json:"completed_at"`
	StartedAt           time.Time        `json:"started_at"`
	Err                 error            `json:"error"`
	Id                  string           `json:"-"`
	Attempts            int              `json:"attempts"`
	Mutex               *sync.RWMutex    `json:"-"`
	status.Status       `json:"-"`
	filesSDK.File       `json:"-"`
	*Job                `json:"-"`
	// Finalizing reports that Files.com answered the request finishing this upload by saying it is still finishing
	// it, and the upload is waiting for the outcome. Each such answer reports the file again to the job's
	// OnBytesChange(status.Uploading) reporters with the same TransferBytes. Only the Complete status means complete.
	// A file that ended Canceled or Errored while Finalizing stopped waiting with an UploadCompletionUnknownError:
	// Files.com may still complete it.
	Finalizing bool `json:"-"`
	// FinalizationProgress is the approximate percentage, from 0 to 100, of that work in Files.com's latest answer
	// while Finalizing, or nil when the answer gave none. 100 is still pending.
	FinalizationProgress *int64 `json:"-"`
	// FinalizingSince is when Files.com first answered that it was still finishing the upload, while Finalizing.
	FinalizingSince time.Time `json:"-"`
}

type MashableError struct {
	error
}

func (m MashableError) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.error)
}

func (m MashableError) Err() error {
	if m.error == nil {
		return nil
	}

	return m
}

func (m MashableError) Unwrap() error {
	return m.error
}

type Reporter func(JobFile)
