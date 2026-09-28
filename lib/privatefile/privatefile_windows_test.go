//go:build windows

package privatefile

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/windows"
)

// These tests read the effective discretionary access control lists Windows
// stores, not the mode bits Go reports, which mean nothing on Windows.

func securityOf(t *testing.T, path string) security {
	t.Helper()
	handle, err := openForSecurity(path, windows.READ_CONTROL)
	require.NoError(t, err)
	defer windows.CloseHandle(handle)
	sec, err := readSecurity(handle)
	require.NoError(t, err)
	return sec
}

// requirePrivate asserts that path has a protected DACL granting access only
// to the current user, SYSTEM and Administrators.
func requirePrivate(t *testing.T, path string) {
	t.Helper()
	sec := securityOf(t, path)
	trust, err := loadTrusted()
	require.NoError(t, err)
	assert.True(t, sec.protected, "%s must not inherit access from its directory", path)
	require.NotNil(t, sec.dacl, "%s has no DACL", path)
	for _, ace := range sec.entries() {
		if ace.header.AceType != windows.ACCESS_ALLOWED_ACE_TYPE {
			continue
		}
		assert.True(t, trust.isTrustedOwner(ace.sid), "%s grants access to %s", path, ace.sid)
	}
}

func grantEveryoneRead(t *testing.T, path string) {
	t.Helper()
	everyone, err := windows.CreateWellKnownSid(windows.WinWorldSid)
	require.NoError(t, err)
	existing := securityOf(t, path)
	acl, err := windows.ACLFromEntries([]windows.EXPLICIT_ACCESS{{
		AccessPermissions: windows.GENERIC_READ,
		AccessMode:        windows.GRANT_ACCESS,
		Inheritance:       windows.SUB_CONTAINERS_AND_OBJECTS_INHERIT,
		Trustee:           windows.TRUSTEE{TrusteeForm: windows.TRUSTEE_IS_SID, TrusteeType: windows.TRUSTEE_IS_WELL_KNOWN_GROUP, TrusteeValue: windows.TrusteeValueFromSID(everyone)},
	}}, existing.dacl)
	require.NoError(t, err)
	require.NoError(t, windows.SetNamedSecurityInfo(path, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION|windows.UNPROTECTED_DACL_SECURITY_INFORMATION, nil, nil, acl, nil))
}

func grantsEveryone(t *testing.T, path string) bool {
	t.Helper()
	everyone, err := windows.CreateWellKnownSid(windows.WinWorldSid)
	require.NoError(t, err)
	for _, ace := range securityOf(t, path).entries() {
		if ace.header.AceType == windows.ACCESS_ALLOWED_ACE_TYPE && ace.sid.Equals(everyone) {
			return true
		}
	}
	return false
}

func TestMkdirAndMkdirTempCreateProtectedPrivateDirectories(t *testing.T) {
	parent := t.TempDir()
	grantEveryoneRead(t, parent)
	dir := filepath.Join(parent, "private")
	require.NoError(t, Mkdir(dir))
	requirePrivate(t, dir)
	require.ErrorIs(t, Mkdir(dir), os.ErrExist)

	temp, err := MkdirTemp(parent, "Files.com-v6-*")
	require.NoError(t, err)
	requirePrivate(t, temp)

	// A file created inside with ordinary Go calls inherits only private entries.
	inside := filepath.Join(temp, "scratch.tmp")
	require.NoError(t, os.WriteFile(inside, []byte("x"), 0o644))
	assert.False(t, grantsEveryone(t, inside))
}

func TestEnsureDirProtectsExistingDirectoryAndItsInheritedEntries(t *testing.T) {
	parent := t.TempDir()
	dir := filepath.Join(parent, "data")
	require.NoError(t, os.Mkdir(dir, 0o755))
	grantEveryoneRead(t, dir)
	legacy := filepath.Join(dir, "legacy.bin")
	require.NoError(t, os.WriteFile(legacy, []byte("legacy"), 0o644))
	require.True(t, grantsEveryone(t, legacy), "the legacy file inherited the broad entry")

	require.NoError(t, EnsureDir(dir))
	requirePrivate(t, dir)
	assert.False(t, grantsEveryone(t, legacy), "tightening the directory re-derives what the file inherits")
}

func TestEnsureEntryPrivateStripsExplicitEntries(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, EnsureDir(dir))
	file := filepath.Join(dir, "payload.bin")
	require.NoError(t, os.WriteFile(file, []byte("payload"), 0o644))
	grantEveryoneRead(t, file)
	require.True(t, grantsEveryone(t, file))

	info, err := os.Lstat(file)
	require.NoError(t, err)
	require.NoError(t, EnsureEntryPrivate(file, info))
	assert.False(t, grantsEveryone(t, file), "an explicit entry on the file itself is removed")
	assert.False(t, securityOf(t, file).protected)
}

func TestEnsureEntryPrivateRefusesFileWithAnotherName(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, EnsureDir(dir))
	file := filepath.Join(dir, "payload.bin")
	require.NoError(t, os.WriteFile(file, []byte("payload"), 0o644))
	outside := filepath.Join(t.TempDir(), "callers-copy.bin")
	require.NoError(t, os.Link(file, outside))
	grantEveryoneRead(t, outside)
	before := securityOf(t, outside).sd.String()

	info, err := os.Lstat(file)
	require.NoError(t, err)
	require.Error(t, CheckEntry(file, info), "a file with another name cannot be made private on its own")
	require.Error(t, EnsureEntryPrivate(file, info))
	assert.Equal(t, before, securityOf(t, outside).sd.String(), "the outside name keeps exactly the access it had")
	assert.True(t, grantsEveryone(t, outside))
}

func TestOpenFileCreatesAndTightensPrivateFiles(t *testing.T) {
	dir := t.TempDir()
	grantEveryoneRead(t, dir)
	root, err := os.OpenRoot(dir)
	require.NoError(t, err)
	defer root.Close()

	f, err := OpenFile(root, "fresh.download", os.O_RDWR|os.O_CREATE|os.O_EXCL, 0o700)
	require.NoError(t, err)
	_, err = f.Write([]byte("staged"))
	require.NoError(t, err)
	require.NoError(t, f.Close())
	requirePrivate(t, filepath.Join(dir, "fresh.download"))

	legacy := filepath.Join(dir, "legacy.download")
	require.NoError(t, os.WriteFile(legacy, []byte("partial"), 0o644))
	require.True(t, grantsEveryone(t, legacy))
	g, err := OpenFile(root, "legacy.download", os.O_WRONLY|os.O_CREATE, 0o700)
	require.NoError(t, err)
	requirePrivate(t, legacy)
	require.NoError(t, g.Close())
	data, err := os.ReadFile(legacy)
	require.NoError(t, err)
	assert.Equal(t, "partial", string(data), "opening an existing stage keeps its bytes")

	_, err = OpenFile(root, "fresh.download", os.O_RDWR|os.O_CREATE|os.O_EXCL, 0o700)
	require.ErrorIs(t, err, os.ErrExist)
}

func TestCheckContainerRejectsDirectoriesOtherAccountsCanModify(t *testing.T) {
	cacheDir, err := os.UserCacheDir()
	require.NoError(t, err)
	require.NoError(t, CheckContainer(cacheDir), "the user's local application data is a trusted location")

	dir := t.TempDir()
	require.NoError(t, EnsureDir(dir))
	require.NoError(t, CheckContainer(dir))

	users, err := windows.CreateWellKnownSid(windows.WinBuiltinUsersSid)
	require.NoError(t, err)
	grant := func(rights windows.ACCESS_MASK) {
		t.Helper()
		acl, err := windows.ACLFromEntries([]windows.EXPLICIT_ACCESS{{
			AccessPermissions: rights,
			AccessMode:        windows.GRANT_ACCESS,
			Inheritance:       windows.NO_INHERITANCE,
			Trustee:           windows.TRUSTEE{TrusteeForm: windows.TRUSTEE_IS_SID, TrusteeType: windows.TRUSTEE_IS_WELL_KNOWN_GROUP, TrusteeValue: windows.TrusteeValueFromSID(users)},
		}}, privateACLForTest(t))
		require.NoError(t, err)
		require.NoError(t, windows.SetNamedSecurityInfo(dir, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION|windows.PROTECTED_DACL_SECURITY_INFORMATION, nil, nil, acl, nil))
	}
	// Adding new entries, as the standard drive root allows Authenticated
	// Users to do, cannot replace an existing protected child.
	grant(windows.FILE_WRITE_DATA | windows.FILE_APPEND_DATA | windows.GENERIC_WRITE | windows.FILE_GENERIC_READ)
	require.NoError(t, CheckContainer(dir), "a directory others may add entries to is still acceptable")

	grant(windows.FILE_GENERIC_READ | fileDeleteChild)
	err = CheckContainer(dir)
	require.Error(t, err, "a directory whose entries Users may delete is refused")
	assert.ErrorContains(t, err, "modify its entries")

	child := filepath.Join(dir, "child")
	require.NoError(t, Mkdir(child))
	require.Error(t, CheckContainer(child), "a private directory below a modifiable one is still replaceable")
}

func privateACLForTest(t *testing.T) *windows.ACL {
	t.Helper()
	acl, err := loadPrivateACL()
	require.NoError(t, err)
	return acl
}

func TestPublicationRestoresExistingDACLOrInherits(t *testing.T) {
	dir := t.TempDir()
	root, err := os.OpenRoot(dir)
	require.NoError(t, err)
	defer root.Close()

	// Fresh output: after publication the file has only what the directory hands down.
	f, err := OpenFile(root, "stage", os.O_RDWR|os.O_CREATE|os.O_EXCL, 0o700)
	require.NoError(t, err)
	staged, err := f.Stat()
	require.NoError(t, err)
	require.NoError(t, f.Close())
	pub, err := PublicationFor(root, "fresh.txt", ".probe", os.ModePerm)
	require.NoError(t, err)
	require.NoError(t, root.Rename("stage", "fresh.txt"))
	require.NoError(t, pub.ApplyToPublished(root, "fresh.txt", staged))
	fresh := securityOf(t, filepath.Join(dir, "fresh.txt"))
	assert.False(t, fresh.protected)
	assert.False(t, fresh.hasExplicitEntries(), "a fresh file carries only inherited entries")

	// Existing restricted output keeps its own protected DACL.
	existing := filepath.Join(dir, "restricted.txt")
	require.NoError(t, os.WriteFile(existing, []byte("old"), 0o600))
	privateACL, err := loadPrivateACL()
	require.NoError(t, err)
	require.NoError(t, windows.SetNamedSecurityInfo(existing, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION|windows.PROTECTED_DACL_SECURITY_INFORMATION, nil, nil, privateACL, nil))
	grantEveryoneRead(t, dir) // the directory would hand down Everyone read
	g, err := OpenFile(root, "stage2", os.O_RDWR|os.O_CREATE|os.O_EXCL, 0o700)
	require.NoError(t, err)
	_, err = g.Write([]byte("new"))
	require.NoError(t, err)
	staged2, err := g.Stat()
	require.NoError(t, err)
	require.NoError(t, g.Close())
	pub, err = PublicationFor(root, "restricted.txt", ".probe", os.ModePerm)
	require.NoError(t, err)
	require.NoError(t, root.Rename("stage2", "restricted.txt"))
	require.NoError(t, pub.ApplyToPublished(root, "restricted.txt", staged2))
	require.ErrorIs(t, pub.ApplyToPublished(root, "fresh.txt", staged2), errPublishedReplaced, "the published access is never applied to another file")
	requirePrivate(t, existing)
	assert.False(t, grantsEveryone(t, existing), "replacing a restricted file must not widen it")
	data, err := os.ReadFile(existing)
	require.NoError(t, err)
	assert.Equal(t, "new", string(data))
}
