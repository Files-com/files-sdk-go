//go:build linux

package fsmount

import (
	"os"
	"testing"
)

func requireOwnerOnlyDir(t *testing.T, dir string) {
	t.Helper()
	info, err := os.Lstat(dir)
	if err != nil {
		t.Fatal(err)
	}
	if !info.IsDir() || info.Mode().Perm()&^0o700 != 0 {
		t.Fatalf("%s must be a directory readable only by its owner, mode %04o", dir, info.Mode().Perm())
	}
}
