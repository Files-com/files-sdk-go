package manager

import (
	"fmt"
	"math"
	"net/http"
	"sync"

	"github.com/Files-com/files-sdk-go/v3/lib"
)

var (
	// ConcurrentFiles is the default number of files a transfer job may process at once.
	ConcurrentFiles = 50
	// ConcurrentFileParts is the default job-wide cap for concurrent upload/download part work.
	ConcurrentFileParts = 50
	// ConcurrentDirectoryList is the default number of local directory listing operations that may run at once.
	ConcurrentDirectoryList = 100
	// AdaptiveUploadV2ConcurrentFiles is the V2 upload job-level file concurrency cap.
	// Adaptive upload V2 still controls HTTP part concurrency dynamically; this value prevents
	// large multi-file jobs from being serialized while the per-destination controller learns.
	// The default is documented in docs/adaptive-upload-v2-file-concurrency-cap.md.
	AdaptiveUploadV2ConcurrentFiles = 128
	// AdaptiveUploadV2ConcurrentFileParts is the V2 upload HTTP part concurrency cap.
	// The adaptive manager treats this as a maximum, not a fixed target.
	AdaptiveUploadV2ConcurrentFileParts = 1024
)

var (
	sharedDefaultOnce    sync.Once
	sharedDefaultManager *Manager
)

type Manager struct {
	FilesManager            ConcurrencyManager
	FilePartsManager        ConcurrencyManager
	DirectoryListingManager ConcurrencyManager
}

type ConcurrencyManager struct {
	*lib.ConstrainedWorkGroup
	DownloadFilesAsSingleStream bool
}

func (ConcurrencyManager) New(maxGoRoutines int, downloadFilesAsSingleStream ...bool) ConcurrencyManager {
	if len(downloadFilesAsSingleStream) == 0 {
		downloadFilesAsSingleStream = append(downloadFilesAsSingleStream, false)
	}

	return ConcurrencyManager{ConstrainedWorkGroup: lib.NewConstrainedWorkGroup(maxGoRoutines), DownloadFilesAsSingleStream: downloadFilesAsSingleStream[0]}
}

func (c ConcurrencyManager) Max() int {
	return c.ConstrainedWorkGroup.Max()
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

func SharedDefault() *Manager {
	sharedDefaultOnce.Do(func() {
		sharedDefaultManager = Default()
	})
	return sharedDefaultManager
}

func Sync() *Manager {
	return New(1, 1, 1)
}

func Build(maxConcurrentConnections, maxConcurrentDirectoryLists int, downloadFilesAsSingleStream ...bool) *Manager {
	maxConcurrentConnections = int(math.Max(float64(maxConcurrentConnections), 1))
	return &Manager{
		FilesManager:            ConcurrencyManager{}.New(maxConcurrentConnections),
		FilePartsManager:        ConcurrencyManager{}.New(maxConcurrentConnections, downloadFilesAsSingleStream...),
		DirectoryListingManager: ConcurrencyManager{}.New(int(math.Max(float64(maxConcurrentDirectoryLists), 1))),
	}
}

func (m *Manager) CreateMatchingClient(client *http.Client) *http.Client {
	if fmt.Sprintf("%T", client.Transport) == "*recorder.Recorder" {
		// Can't modify VCR in Test mode.
		return client
	}

	switch t := client.Transport.(type) {
	case *lib.Transport:
		t.MaxConnsPerHost = m.FilePartsManager.Max()
		return client
	default:
		defaultTransport := lib.DefaultPooledTransport()
		defaultTransport.MaxConnsPerHost = m.FilePartsManager.Max()
		client.Transport = defaultTransport
		return client
	}
}
