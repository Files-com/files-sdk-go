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
	maxGoRoutines               int
	DownloadFilesAsSingleStream bool
}

func (ConcurrencyManager) New(maxGoRoutines int, downloadFilesAsSingleStream ...bool) ConcurrencyManager {
	if len(downloadFilesAsSingleStream) == 0 {
		downloadFilesAsSingleStream = append(downloadFilesAsSingleStream, false)
	}
	return ConcurrencyManager{ConcurrencyManager: goccm.New(maxGoRoutines), maxGoRoutines: maxGoRoutines, DownloadFilesAsSingleStream: downloadFilesAsSingleStream[0]}
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

func Build(maxConcurrentConnections, maxConcurrentDirectoryLists int, downloadFilesAsSingleStream ...bool) *Manager {
	maxConcurrentConnections = int(math.Max(float64(maxConcurrentConnections), 1))
	return &Manager{
		FilesManager:            ConcurrencyManager{}.New(maxConcurrentConnections),
		FilePartsManager:        ConcurrencyManager{}.New(maxConcurrentConnections, downloadFilesAsSingleStream...),
		DirectoryListingManager: ConcurrencyManager{}.New(int(math.Max(float64(maxConcurrentDirectoryLists), 1))),
	}
}
