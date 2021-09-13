package manager

import (
	"context"

	"github.com/zenthangplus/goccm"
)

var (
	ConcurrentFiles     = 50
	ConcurrentFileParts = 75
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

func Wait(ctx context.Context, manager goccm.ConcurrencyManager) bool {
	if ctx.Err() != nil {
		return false
	}
	manager.Wait()
	if ctx.Err() != nil {
		return false
	}
	return true
}
