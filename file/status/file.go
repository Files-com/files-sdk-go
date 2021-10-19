package status

import (
	"sync"
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
	Mutex         *sync.RWMutex
}

type Reporter func(File)
