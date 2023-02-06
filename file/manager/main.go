package manager

import (
	"context"
	"math"

	"github.com/zenthangplus/goccm"
)

var (
	ConcurrentFiles         = 50
	ConcurrentFileParts     = 75
	ConcurrentDirectoryList = 100
)

type Manager struct {
	FilesManager            goccm.ConcurrencyManager
	FilePartsManager        goccm.ConcurrencyManager
	DirectoryListingManager goccm.ConcurrencyManager
}

func New(files, fileParts, directoryListing int) *Manager {
	return &Manager{
		FilesManager:            goccm.New(files),
		FilePartsManager:        goccm.New(fileParts),
		DirectoryListingManager: goccm.New(directoryListing),
	}
}

func Default() *Manager {
	return New(ConcurrentFiles, ConcurrentFileParts, ConcurrentDirectoryList)
}

func Sync() *Manager {
	return New(1, 1, 1)
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

func Build(maxConcurrentConnections, maxConcurrentDirectoryLists int) *Manager {
	maxConcurrentConnections = int(math.Max(float64(maxConcurrentConnections), 1))
	return &Manager{
		FilesManager:            goccm.New(int(math.Max(float64(maxConcurrentConnections/2), 1))),
		FilePartsManager:        goccm.New(maxConcurrentConnections),
		DirectoryListingManager: goccm.New(int(math.Max(float64(maxConcurrentDirectoryLists), 1))),
	}
}
