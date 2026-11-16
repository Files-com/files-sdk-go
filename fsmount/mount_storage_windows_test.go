//go:build windows

package fsmount

import (
	"os"
	"path/filepath"
	"testing"
	"unsafe"

	"golang.org/x/sys/windows"
)

// While a mount runs, its open storage keeps the folders above it in place, so
// an account that may delete or rename one of them cannot move the storage
// away and leave its path to be recreated elsewhere.
func TestMountStorageKeepsTheFoldersAboveItInPlace(t *testing.T) {
	parent := filepath.Join(t.TempDir(), "Temp")
	if err := os.Mkdir(parent, 0o700); err != nil {
		t.Fatal(err)
	}
	storage, err := newMountStorage(parent, DefaultFileSystemName)
	if err != nil {
		t.Fatalf("newMountStorage failed: %v", err)
	}
	moved := parent + "-moved"
	if err := os.Rename(parent, moved); err == nil {
		t.Fatalf("%s was moved while the mount's storage inside it was open", parent)
	}
	if err := storage.remove(); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(parent, moved); err != nil {
		t.Fatalf("%s should move once the mount's storage is removed: %v", parent, err)
	}
}

// requireOwnerOnlyDir asserts that dir has a protected DACL that grants access
// only to the current user, SYSTEM and Administrators.
func requireOwnerOnlyDir(t *testing.T, dir string) {
	t.Helper()
	sd, err := windows.GetNamedSecurityInfo(dir, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION)
	if err != nil {
		t.Fatal(err)
	}
	control, _, err := sd.Control()
	if err != nil {
		t.Fatal(err)
	}
	if control&windows.SE_DACL_PROTECTED == 0 {
		t.Fatalf("%s inherits access from its directory", dir)
	}
	dacl, _, err := sd.DACL()
	if err != nil || dacl == nil {
		t.Fatalf("%s has no DACL: %v", dir, err)
	}
	user, err := windows.GetCurrentProcessToken().GetTokenUser()
	if err != nil {
		t.Fatal(err)
	}
	system, _ := windows.CreateWellKnownSid(windows.WinLocalSystemSid)
	admins, _ := windows.CreateWellKnownSid(windows.WinBuiltinAdministratorsSid)
	for i := uint32(0); i < uint32(dacl.AceCount); i++ {
		var ace *windows.ACCESS_ALLOWED_ACE
		if err := windows.GetAce(dacl, i, &ace); err != nil {
			t.Fatal(err)
		}
		if ace.Header.AceType != windows.ACCESS_ALLOWED_ACE_TYPE {
			continue
		}
		sid := (*windows.SID)(unsafe.Pointer(&ace.SidStart))
		if !sid.Equals(user.User.Sid) && !sid.Equals(system) && !sid.Equals(admins) {
			t.Fatalf("%s grants access to %s", dir, sid)
		}
	}
}
