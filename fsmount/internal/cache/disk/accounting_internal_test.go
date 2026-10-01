//go:build linux || windows

package disk

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache"
	cachelog "github.com/Files-com/files-sdk-go/v3/fsmount/internal/log"
)

type blockingMaintenanceLogger struct {
	cachelog.NoOpLogger
	started chan struct{}
	release chan struct{}
	once    sync.Once
}

func (l *blockingMaintenanceLogger) Debug(format string, _ ...any) {
	if !strings.Contains(format, "performing maintenance") {
		return
	}
	l.once.Do(func() { close(l.started) })
	<-l.release
}

func TestMaintenancePreservesConcurrentRangeAccounting(t *testing.T) {
	cacheStore, err := NewDiskCache(t.TempDir())
	if err != nil {
		t.Fatalf("NewDiskCache failed: %v", err)
	}
	// Maintenance skips a cycle while the startup scan runs.
	cacheStore.waitForScan()

	path := "/concurrent-maintenance.bin"
	meta := cache.NewEntryMetadata(path, 20, time.Date(2026, time.August, 14, 12, 0, 0, 0, time.UTC))
	if _, err := cacheStore.WriteRange(path, meta, bytes.Repeat([]byte("a"), 10), 0); err != nil {
		t.Fatalf("initial WriteRange failed: %v", err)
	}

	orphan := filepath.Join(cacheStore.dataDir, "orphan.bin")
	if err := os.WriteFile(orphan, []byte("orphan"), 0o644); err != nil {
		t.Fatalf("WriteFile orphan failed: %v", err)
	}

	logger := &blockingMaintenanceLogger{started: make(chan struct{}), release: make(chan struct{})}
	cacheStore.log = logger
	cacheStore.delMu.Lock()
	delMuLocked := true
	var releaseOnce sync.Once
	releaseMaintenance := func() {
		releaseOnce.Do(func() { close(logger.release) })
	}
	defer func() {
		releaseMaintenance()
		if delMuLocked {
			cacheStore.delMu.Unlock()
		}
	}()

	maintenanceDone := make(chan struct{})
	go func() {
		cacheStore.runMaintenanceOnce(context.Background())
		close(maintenanceDone)
	}()

	select {
	case <-logger.started:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for maintenance snapshot")
	}

	writeDone := make(chan error, 1)
	go func() {
		_, writeErr := cacheStore.WriteRange(path, meta, bytes.Repeat([]byte("b"), 10), 10)
		writeDone <- writeErr
	}()
	releaseMaintenance()

	select {
	case err := <-writeDone:
		if err != nil {
			t.Fatalf("concurrent WriteRange failed: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("concurrent WriteRange did not pass the completed maintenance snapshot")
	}

	cacheStore.delMu.Unlock()
	delMuLocked = false
	select {
	case <-maintenanceDone:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for maintenance completion")
	}

	stats := cacheStore.Stats()
	if got := stats.SizeBytes.Load(); got != meta.Size {
		t.Fatalf("SizeBytes after maintenance = %d, want %d", got, meta.Size)
	}
	if got := stats.FileCount.Load(); got != 1 {
		t.Fatalf("FileCount after maintenance = %d, want 1", got)
	}
}
