package status

import (
	"sync"
	"time"

	"encoding/json"

	filesSDK "github.com/Files-com/files-sdk-go/v2"
)

type File struct {
	StatusName    string        `json:"status"`
	TransferBytes int64         `json:"transferred_bytes"`
	Size          int64         `json:"size_bytes"`
	LocalPath     string        `json:"local_path"`
	RemotePath    string        `json:"remote_path"`
	EndedAt       time.Time     `json:"transferred_at"`
	Err           error         `json:"error"`
	Id            string        `json:"-"`
	Attempts      int           `json:"attempts"`
	Mutex         *sync.RWMutex `json:"-"`
	LastByte      time.Time     `json:"-"`
	Status        `json:"-"`
	filesSDK.File `json:"-"`
	*Job          `json:"-"`
}

type MashableError struct {
	error
}

func (me MashableError) MarshalJSON() ([]byte, error) {
	return json.Marshal(me.Error())
}

func (me MashableError) Err() error {
	if me.error == nil {
		return nil
	}

	return me
}

type Reporter func(File)
