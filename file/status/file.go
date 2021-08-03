package status

import (
	"context"

	filesSDK "github.com/Files-com/files-sdk-go"
)

type File struct {
	filesSDK.File
	*Job
	Status
	TransferBytes int64
	Cancel        context.CancelFunc
	LocalPath     string
	RemotePath    string
	Id            string
	Sync          bool
}

type Reporter func(File, error)
