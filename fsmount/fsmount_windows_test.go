package fsmount

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidMountPoint(t *testing.T) {
	tests := []struct {
		mountPoint string
		valid      bool
	}{
		{"", true},
		{"C:", true},
		{"Z:", true},
		{"A:", true},
		{"1:", false},
		{"C", false},
		{"CC:", false},
		{"C::", false},
		{"C:/", false},
		{"C:\\", false},
		{"C:extra", false},
	}

	for _, test := range tests {
		err := validateMountPoint(test.mountPoint)
		if (err == nil) != test.valid {
			t.Errorf("validateMountPoint(%q) = %v; want valid=%v", test.mountPoint, err, test.valid)
		}
	}
}

func TestOpenLocalFileAllowsRenameWhileOpen(t *testing.T) {
	dir := t.TempDir()
	oldPath := filepath.Join(dir, "illustrator.tmp")
	newPath := filepath.Join(dir, "~ai-rename.tmp")

	f, err := openLocalFile(oldPath, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		t.Fatalf("openLocalFile failed: %v", err)
	}
	defer f.Close()

	if _, err := f.Write([]byte("test")); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	if err := os.Rename(oldPath, newPath); err != nil {
		t.Fatalf("Rename while file is open failed: %v", err)
	}
}
