package file

import (
	"testing"

	"github.com/Files-com/files-sdk-go/v3/file/manager"
)

func TestUploadV2MaxConcurrencyUsesAdaptiveUploadOverride(t *testing.T) {
	t.Cleanup(func() {
		manager.SetAdaptiveUploadV2ConcurrentFileParts(0)
	})

	manager.SetAdaptiveUploadV2ConcurrentFileParts(7)

	if got := (&uploadIO{}).uploadV2MaxConcurrency(); got != 7 {
		t.Fatalf("uploadV2MaxConcurrency() = %d, want 7", got)
	}
}

func TestUploadV2MaxConcurrencyDoesNotRaiseTargetDefault(t *testing.T) {
	t.Cleanup(func() {
		manager.SetAdaptiveUploadV2ConcurrentFileParts(0)
	})

	manager.SetAdaptiveUploadV2ConcurrentFileParts(1024)

	if got := (&uploadIO{}).uploadV2MaxConcurrency(); got != AdaptiveTransferDefaultMaxConcurrency {
		t.Fatalf("uploadV2MaxConcurrency() = %d, want %d", got, AdaptiveTransferDefaultMaxConcurrency)
	}
}

func TestDownloadV2MaxConcurrencyUsesAdaptiveDownloadOverride(t *testing.T) {
	t.Cleanup(func() {
		manager.SetAdaptiveDownloadV2ConcurrentFileParts(0)
	})

	manager.SetAdaptiveDownloadV2ConcurrentFileParts(7)

	if got := downloadV2MaxConcurrency(&Job{}, DownloaderParams{}); got != 7 {
		t.Fatalf("downloadV2MaxConcurrency() = %d, want 7", got)
	}
}
