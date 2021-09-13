package status

import (
	"time"

	filesSDK "github.com/Files-com/files-sdk-go/v2"
)

type File struct {
	filesSDK.File
	*Job
	Status
	TransferBytes int64
	LocalPath     string
	RemotePath    string
	Id            string
	LastByte      time.Time
	Err           error
}

func (f *File) ToStatusFile() File {
	return *f
}

func (f *File) SetStatus(status Status, err error) {
	var setError bool
	f.Status, setError = SetStatus(f.Status, status, err)
	if setError {
		f.Err = err
	}
}

type Reporter func(File)
