package status

import (
	filesSDK "github.com/Files-com/files-sdk-go"
)

type Report interface {
	TransferBytes() int64
	File() filesSDK.File
	Cancel()
	Destination() string
	IStatus
	Job() Job
	Id() string
}
