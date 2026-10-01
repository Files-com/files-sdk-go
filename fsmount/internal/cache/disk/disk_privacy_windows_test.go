//go:build windows

package disk_test

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"time"
	"unsafe"

	fscache "github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache"
	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache/disk"
	"golang.org/x/sys/windows"
)

// A cache root whose access control list hands Everyone read access down to
// everything created inside it, and a legacy cache file that additionally
// carries its own Everyone entry, must both end up granting no access beyond
// the owner, SYSTEM and Administrators once the cache opens. Windows checks a
// file's own list when a file is opened by a known path, so tightening the
// directories alone is not enough.
func TestNewDiskCacheProtectsLegacyStateUnderBroadWindowsACL(t *testing.T) {
	root := filepath.Join(t.TempDir(), "cache")
	if err := os.Mkdir(root, 0o755); err != nil {
		t.Fatal(err)
	}
	legacy, err := disk.NewDiskCache(root, disk.WithLruFlushInterval(time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	const path = "/confidential/payroll.txt"
	content := []byte("cached bytes that only the owner may read")
	if _, err := legacy.Write(path, content, 0); err != nil {
		t.Fatal(err)
	}
	if err := legacy.Commit(path, fscache.NewEntryMetadata(path, int64(len(content)), time.Unix(1, 0))); err != nil {
		t.Fatal(err)
	}
	if _, err := legacy.WritePartial("/confidential/in-progress.txt", []byte("partial"), 0); err != nil {
		t.Fatal(err)
	}
	legacy.StartMaintenance()
	legacy.StopMaintenance()

	// Reproduce a broad legacy configuration: inherited Everyone read on the
	// whole tree, plus an explicit Everyone entry on each file, as if it had
	// been granted directly.
	grantEveryoneRead(t, root, true)
	var files []string
	err = filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			files = append(files, p)
			grantEveryoneRead(t, p, false)
		} else if p != root {
			// Re-derive inheritance the way SetNamedSecurityInfo on the root does.
			grantEveryoneRead(t, p, true)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(files) < 4 {
		t.Fatalf("expected payload, partial, metadata and LRU files, found %d", len(files))
	}
	for _, f := range files {
		if !grantsEveryone(t, f) {
			t.Fatalf("test setup did not expose %s", f)
		}
	}

	reopened, err := disk.NewDiskCache(root, disk.WithLruFlushInterval(time.Hour))
	if err != nil {
		t.Fatalf("NewDiskCache on legacy state failed: %v", err)
	}
	for _, f := range files {
		if _, err := os.Lstat(f); errors.Is(err, fs.ErrNotExist) {
			// Opening the cache deletes partial data an earlier process left.
			continue
		}
		if grantsEveryone(t, f) {
			t.Errorf("%s still grants Everyone access after the cache opened", f)
		}
	}
	for _, dir := range []string{"data", "partial", "state"} {
		if grantsEveryone(t, filepath.Join(root, dir)) {
			t.Errorf("%s still grants Everyone access", dir)
		}
	}
	if !grantsEveryone(t, root) {
		t.Error("the caller's root must be left as the caller configured it")
	}
	buf := make([]byte, len(content))
	if n, err := reopened.Read(path, buf, 0); err != nil || n != len(content) || string(buf) != string(content) {
		t.Fatalf("payload unreadable by owner after protection: n=%d err=%v", n, err)
	}
	if err := reopened.Clear(); err != nil {
		t.Fatal(err)
	}
	if _, err := reopened.Write(path, content, 0); err != nil {
		t.Fatal(err)
	}
	if grantsEveryone(t, filepath.Join(root, "state", "metadata")) {
		t.Error("the metadata directory recreated by Clear grants Everyone access")
	}
}

// A legacy cache file that also has a name outside the cache shares its
// access control list with that name. Opening the cache must refuse it before
// the directory's list is changed, because changing a directory's list
// rewrites what its descendants inherit, so the caller's outside file keeps
// exactly the access it had.
func TestNewDiskCacheRefusesLegacyFileWithOutsideName(t *testing.T) {
	root := filepath.Join(t.TempDir(), "cache")
	if err := os.Mkdir(root, 0o755); err != nil {
		t.Fatal(err)
	}
	legacy, err := disk.NewDiskCache(root, disk.WithLruFlushInterval(time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	const path = "/confidential/payroll.txt"
	content := []byte("cached bytes shared with a caller's file")
	if _, err := legacy.Write(path, content, 0); err != nil {
		t.Fatal(err)
	}
	legacy.StartMaintenance()
	legacy.StopMaintenance()
	var payloadFile string
	err = filepath.WalkDir(filepath.Join(root, "data"), func(p string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			payloadFile = p
		}
		return err
	})
	if err != nil || payloadFile == "" {
		t.Fatalf("no payload file found: %v", err)
	}
	outside := filepath.Join(t.TempDir(), "callers-copy.txt")
	if err := os.Link(payloadFile, outside); err != nil {
		t.Fatal(err)
	}
	// The caller deliberately shares the outside name with Everyone.
	grantEveryoneRead(t, outside, false)
	before := sddlOf(t, outside)
	if !grantsEveryone(t, outside) {
		t.Fatal("test setup did not share the outside name")
	}

	if _, err := disk.NewDiskCache(root, disk.WithLruFlushInterval(time.Hour)); err == nil {
		t.Fatal("expected NewDiskCache to refuse a cache file that has a name outside the cache")
	}
	if after := sddlOf(t, outside); after != before {
		t.Fatalf("the caller's outside name changed access:\nbefore %s\nafter  %s", before, after)
	}
	if !grantsEveryone(t, outside) {
		t.Fatal("the caller's outside name lost the access the caller granted")
	}
	data, err := os.ReadFile(outside)
	if err != nil || string(data) != string(content) {
		t.Fatalf("the caller's outside file changed: %q, %v", data, err)
	}
}

func sddlOf(t *testing.T, path string) string {
	t.Helper()
	sd, err := windows.GetNamedSecurityInfo(path, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION|windows.OWNER_SECURITY_INFORMATION)
	if err != nil {
		t.Fatal(err)
	}
	return sd.String()
}

func everyoneSID(t *testing.T) *windows.SID {
	t.Helper()
	sid, err := windows.CreateWellKnownSid(windows.WinWorldSid)
	if err != nil {
		t.Fatal(err)
	}
	return sid
}

func daclOf(t *testing.T, path string) (*windows.ACL, error) {
	sd, err := windows.GetNamedSecurityInfo(path, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION)
	if err != nil {
		return nil, err
	}
	dacl, _, err := sd.DACL()
	return dacl, err
}

func grantEveryoneRead(t *testing.T, path string, inheritable bool) {
	t.Helper()
	existing, err := daclOf(t, path)
	if err != nil {
		t.Fatal(err)
	}
	inheritance := uint32(windows.NO_INHERITANCE)
	if inheritable {
		inheritance = windows.SUB_CONTAINERS_AND_OBJECTS_INHERIT
	}
	acl, err := windows.ACLFromEntries([]windows.EXPLICIT_ACCESS{{
		AccessPermissions: windows.GENERIC_READ,
		AccessMode:        windows.GRANT_ACCESS,
		Inheritance:       inheritance,
		Trustee:           windows.TRUSTEE{TrusteeForm: windows.TRUSTEE_IS_SID, TrusteeType: windows.TRUSTEE_IS_WELL_KNOWN_GROUP, TrusteeValue: windows.TrusteeValueFromSID(everyoneSID(t))},
	}}, existing)
	if err != nil {
		t.Fatal(err)
	}
	if err := windows.SetNamedSecurityInfo(path, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION|windows.UNPROTECTED_DACL_SECURITY_INFORMATION, nil, nil, acl, nil); err != nil {
		t.Fatal(err)
	}
}

func grantsEveryone(t *testing.T, path string) bool {
	t.Helper()
	dacl, err := daclOf(t, path)
	if err != nil {
		t.Fatal(err)
	}
	if dacl == nil {
		return true
	}
	everyone := everyoneSID(t)
	for i := uint32(0); i < uint32(dacl.AceCount); i++ {
		var ace *windows.ACCESS_ALLOWED_ACE
		if err := windows.GetAce(dacl, i, &ace); err != nil {
			t.Fatal(err)
		}
		if ace.Header.AceType == windows.ACCESS_ALLOWED_ACE_TYPE && (*windows.SID)(unsafe.Pointer(&ace.SidStart)).Equals(everyone) {
			return true
		}
	}
	return false
}
