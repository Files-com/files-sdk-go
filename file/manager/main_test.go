package manager

import "testing"

func TestAdaptiveUploadV2ConcurrentFilePartsOverrideCapsDefault(t *testing.T) {
	t.Cleanup(func() {
		SetAdaptiveUploadV2ConcurrentFileParts(0)
	})

	SetAdaptiveUploadV2ConcurrentFileParts(7)

	if got := EffectiveAdaptiveUploadV2ConcurrentFileParts(1024); got != 7 {
		t.Fatalf("EffectiveAdaptiveUploadV2ConcurrentFileParts(1024) = %d, want 7", got)
	}
	if got := EffectiveAdaptiveUploadV2ConcurrentFileParts(4); got != 4 {
		t.Fatalf("EffectiveAdaptiveUploadV2ConcurrentFileParts(4) = %d, want 4", got)
	}
}

func TestAdaptiveUploadV2ConcurrentFilePartsOverrideCanReset(t *testing.T) {
	t.Cleanup(func() {
		SetAdaptiveUploadV2ConcurrentFileParts(0)
	})

	SetAdaptiveUploadV2ConcurrentFileParts(7)
	SetAdaptiveUploadV2ConcurrentFileParts(0)

	if got := EffectiveAdaptiveUploadV2ConcurrentFileParts(1024); got != 1024 {
		t.Fatalf("EffectiveAdaptiveUploadV2ConcurrentFileParts(1024) = %d, want 1024", got)
	}
}

func TestAdaptiveDownloadV2ConcurrentFilePartsOverrideCapsDefault(t *testing.T) {
	t.Cleanup(func() {
		SetAdaptiveDownloadV2ConcurrentFileParts(0)
	})

	SetAdaptiveDownloadV2ConcurrentFileParts(7)

	if got := EffectiveAdaptiveDownloadV2ConcurrentFileParts(); got != 7 {
		t.Fatalf("EffectiveAdaptiveDownloadV2ConcurrentFileParts() = %d, want 7", got)
	}
}

func TestAdaptiveDownloadV2ConcurrentFilePartsOverrideCanReset(t *testing.T) {
	t.Cleanup(func() {
		SetAdaptiveDownloadV2ConcurrentFileParts(0)
	})

	SetAdaptiveDownloadV2ConcurrentFileParts(7)
	SetAdaptiveDownloadV2ConcurrentFileParts(0)

	if got := EffectiveAdaptiveDownloadV2ConcurrentFileParts(); got != AdaptiveDownloadV2ConcurrentFileParts {
		t.Fatalf("EffectiveAdaptiveDownloadV2ConcurrentFileParts() = %d, want %d", got, AdaptiveDownloadV2ConcurrentFileParts)
	}
}
