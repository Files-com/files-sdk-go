package file

import (
	"testing"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
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

func TestUploadV2MaxConcurrencyUsesDirectProfile(t *testing.T) {
	upload := &uploadIO{FileUploadPart: files_sdk.FileUploadPart{
		UploadUri:            "://unused",
		DirectConnectionInfo: testDirectConnectionInfo("/uploads?jwt=direct-token"),
	}}

	if got := upload.uploadV2MaxConcurrency(); got != AdaptiveTransferDirectMaxConcurrency {
		t.Fatalf("uploadV2MaxConcurrency() = %d, want %d", got, AdaptiveTransferDirectMaxConcurrency)
	}
}

func TestDownloadV2MaxConcurrencyUsesAdaptiveDownloadOverride(t *testing.T) {
	t.Cleanup(func() {
		manager.SetAdaptiveDownloadV2ConcurrentFileParts(0)
	})

	manager.SetAdaptiveDownloadV2ConcurrentFileParts(7)

	if got := downloadV2MaxConcurrency(&Job{}, DownloaderParams{}, downloadV2TargetDefault); got != 7 {
		t.Fatalf("downloadV2MaxConcurrency() = %d, want 7", got)
	}
}

func TestDownloadV2MaxConcurrencyDoesNotRaiseTargetDefault(t *testing.T) {
	t.Cleanup(func() {
		manager.SetAdaptiveDownloadV2ConcurrentFileParts(0)
	})

	manager.SetAdaptiveDownloadV2ConcurrentFileParts(1024)

	if got := downloadV2MaxConcurrency(&Job{}, DownloaderParams{}, downloadV2TargetDefault); got != AdaptiveTransferDefaultMaxConcurrency {
		t.Fatalf("downloadV2MaxConcurrency() = %d, want %d", got, AdaptiveTransferDefaultMaxConcurrency)
	}
}

func TestDownloadV2MaxConcurrencyUsesDirectProfile(t *testing.T) {
	if got := downloadV2MaxConcurrency(&Job{}, DownloaderParams{}, downloadV2TargetDirect); got != AdaptiveTransferDirectMaxConcurrency {
		t.Fatalf("downloadV2MaxConcurrency() = %d, want %d", got, AdaptiveTransferDirectMaxConcurrency)
	}
}
