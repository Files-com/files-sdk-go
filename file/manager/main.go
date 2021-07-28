package manager

import (
	"github.com/zenthangplus/goccm"
)

var (
	ConcurrentFiles     = 25
	ConcurrentFileParts = 100
)

type Manager struct {
	FilesManager     goccm.ConcurrencyManager
	FilePartsManager goccm.ConcurrencyManager
}

func New(ConcurrentFiles int, ConcurrentFileParts int) *Manager {
	return &Manager{
		FilesManager:     goccm.New(ConcurrentFiles),
		FilePartsManager: goccm.New(ConcurrentFileParts),
	}
}

func Default() *Manager {
	return New(ConcurrentFiles, ConcurrentFileParts)
}
