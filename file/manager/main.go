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
	FilesManager            ConcurrencyManager
	FilePartsManager        ConcurrencyManager
	DirectoryListingManager ConcurrencyManager
}

type ConcurrencyManager struct {
	goccm.ConcurrencyManager
	maxGoRoutines int
}

func (ConcurrencyManager) New(maxGoRoutines int) ConcurrencyManager {
	return ConcurrencyManager{goccm.New(maxGoRoutines), maxGoRoutines}
}

func (c ConcurrencyManager) Max() int {
	return c.maxGoRoutines
}

func New(files, fileParts, directoryListing int) *Manager {
	return &Manager{
		FilesManager:            ConcurrencyManager{}.New(files),
		FilePartsManager:        ConcurrencyManager{}.New(fileParts),
		DirectoryListingManager: ConcurrencyManager{}.New(directoryListing),
	}
}

func Default() *Manager {
	return New(ConcurrentFiles, ConcurrentFileParts, ConcurrentDirectoryList)
}

func Sync() *Manager {
	return New(1, 1, 1)
}

func Wait(ctx context.Context, manager ConcurrencyManager) bool {
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
		FilesManager:            ConcurrencyManager{}.New(int(math.Max(float64(maxConcurrentConnections/2), 1))),
		FilePartsManager:        ConcurrencyManager{}.New(maxConcurrentConnections),
		DirectoryListingManager: ConcurrencyManager{}.New(int(math.Max(float64(maxConcurrentDirectoryLists), 1))),
	}
}
