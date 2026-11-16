//go:build windows

package privatefile

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
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
	assert.ErrorContains(t, err, "delete files or subfolders")
	assert.ErrorContains(t, err, "ask your IT administrator")

	child := filepath.Join(dir, "child")
	require.NoError(t, Mkdir(child))
	require.Error(t, CheckContainer(child), "a private directory below a modifiable one is still replaceable")
}

func TestCheckContainerAcceptsAnIndividualLocalAdministrator(t *testing.T) {
	administrator := localAccountSID(t, 500)
	dir := t.TempDir()
	require.NoError(t, EnsureDir(dir))
	acl, err := windows.ACLFromEntries([]windows.EXPLICIT_ACCESS{{
		AccessPermissions: fileDeleteChild | windows.DELETE,
		AccessMode:        windows.GRANT_ACCESS,
		Trustee:           windows.TRUSTEE{TrusteeForm: windows.TRUSTEE_IS_SID, TrusteeValue: windows.TrusteeValueFromSID(administrator)},
	}}, privateACLForTest(t))
	require.NoError(t, err)
	require.NoError(t, windows.SetNamedSecurityInfo(dir, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION|windows.PROTECTED_DACL_SECURITY_INFORMATION, nil, nil, acl, nil))
	require.NoError(t, CheckContainer(dir), "an explicit grant to a verified administrator is as trusted as a grant to Administrators")
}

func TestCheckContainerRejectsAnotherUsersModificationRights(t *testing.T) {
	guest := localAccountSID(t, 501)
	account, domain, _, err := guest.LookupAccount("")
	require.NoError(t, err)
	dir := t.TempDir()
	require.NoError(t, EnsureDir(dir))
	acl, err := windows.ACLFromEntries([]windows.EXPLICIT_ACCESS{{
		AccessPermissions: fileDeleteChild,
		AccessMode:        windows.GRANT_ACCESS,
		Trustee:           windows.TRUSTEE{TrusteeForm: windows.TRUSTEE_IS_SID, TrusteeValue: windows.TrusteeValueFromSID(guest)},
	}}, privateACLForTest(t))
	require.NoError(t, err)
	require.NoError(t, windows.SetNamedSecurityInfo(dir, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION|windows.PROTECTED_DACL_SECURITY_INFORMATION, nil, nil, acl, nil))
	err = CheckContainer(dir)
	require.ErrorContains(t, err, domain+`\`+account)
	assert.ErrorContains(t, err, "delete files or subfolders")
	assert.ErrorContains(t, err, "ask your IT administrator")
}

func localAccountSID(t *testing.T, rid uint32) *windows.SID {
	t.Helper()
	// The account database is named after the NetBIOS computer name, which
	// differs from the host name when that is longer than 15 characters.
	computer, err := windows.ComputerName()
	require.NoError(t, err)
	machine, _, kind, err := windows.LookupSID("", computer)
	require.NoError(t, err)
	require.Equal(t, uint32(windows.SidTypeDomain), kind)
	sid, err := windows.StringToSid(fmt.Sprintf("%s-%d", machine, rid))
	require.NoError(t, err)
	return sid
}

func TestCheckContainerHonorsDenyEntriesWithoutIgnoringRemainingAccess(t *testing.T) {
	trust, err := loadTrusted()
	require.NoError(t, err)
	for _, tc := range []struct {
		name    string
		entries string
		wantErr string
	}{
		{"denied delete-child access", "(D;;0x40;;;BU)(A;;0x40;;;BU)", ""},
		{"delete-folder access remains", "(D;;0x40;;;BU)(A;;0x10040;;;BU)", "delete or rename the folder"},
		{"deny to another group", "(D;;0x40;;;BA)(A;;0x40;;;BU)", "delete files or subfolders"},
		{"explicit allow before inherited deny", "(A;;0x40;;;BU)(D;ID;0x40;;;BU)", "delete files or subfolders"},
		{"denied account that cannot be looked up", "(D;;0x40;;;S-1-5-21-1-2-3-999999)(A;;0x40;;;S-1-5-21-1-2-3-999999)", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			require.NoError(t, EnsureDir(dir))
			setPrivateDACLWith(t, trust, dir, tc.entries)
			err = CheckContainer(dir)
			if tc.wantErr != "" {
				require.ErrorContains(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// setPrivateDACLWith gives path a protected DACL holding entries, an SDDL
// fragment kept in its on-disk order, before the private entries.
func setPrivateDACLWith(t *testing.T, trust trustedSIDs, path string, entries string) {
	t.Helper()
	sd, err := windows.SecurityDescriptorFromString(fmt.Sprintf("D:P%s(A;OICI;FA;;;%s)(A;OICI;FA;;;SY)(A;OICI;FA;;;BA)", entries, trust.user))
	require.NoError(t, err)
	acl, _, err := sd.DACL()
	require.NoError(t, err)
	require.NoError(t, windows.SetNamedSecurityInfo(path, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION|windows.PROTECTED_DACL_SECURITY_INFORMATION, nil, nil, acl, nil))
}

// OpenAI Codex grants its sandbox group read, write, execute and delete access
// (0x1301bf, without delete-child) to its project folder and everything below
// it. Started in the home folder, that covers the temporary folder and every
// folder between, with the right to create folders in them. A directory held
// open below them keeps them all in place on a local NTFS volume, so
// MkdirTempHeld accepts the temporary folder, while CheckContainer, whose
// caller holds nothing, still refuses it.
func TestMkdirTempHeldDoesNotDependOnTheFoldersAboveIt(t *testing.T) {
	trust, err := loadTrusted()
	require.NoError(t, err)
	home := t.TempDir()
	require.NoError(t, EnsureDir(home))
	setPrivateDACLWith(t, trust, home, "(A;OICI;0x1301bf;;;BU)")
	temp := filepath.Join(home, "AppData", "Local", "Temp")
	require.NoError(t, os.MkdirAll(temp, 0o700))
	require.ErrorContains(t, CheckContainer(temp), "delete or rename the folder")

	store, held, err := MkdirTempHeld(temp, "store-*")
	require.NoError(t, err)
	t.Cleanup(func() { held.Close() })
	assert.Equal(t, temp, filepath.Dir(store))
	requirePrivate(t, store)
	require.Error(t, os.Rename(store, store+"-moved"), "the storage is held from the moment it exists")
	for dir := temp; dir != home; dir = filepath.Dir(dir) {
		require.Error(t, os.Rename(dir, dir+"-moved"), "%s stays in place while the storage is held", dir)
	}
}

// A root accepted only while held may have been open to other accounts, so a
// directory found in it is used only if it is already private: one that is
// not could have been put there by another account, which might still hold it
// open.
func TestEnsureDirInCreatesOrAcceptsOnlyPrivateDirectories(t *testing.T) {
	rootPath := t.TempDir()
	root, err := os.Open(rootPath)
	require.NoError(t, err)
	t.Cleanup(func() { root.Close() })

	require.NoError(t, EnsureDirIn(root, "data"))
	requirePrivate(t, filepath.Join(rootPath, "data"))
	require.NoError(t, EnsureDirIn(root, "data"), "a private directory from an earlier run is reused")

	planted := filepath.Join(rootPath, "planted")
	require.NoError(t, os.Mkdir(planted, 0o700))
	grantEveryoneRead(t, planted)
	require.ErrorContains(t, EnsureDirIn(root, "planted"), "not private")
	assert.True(t, grantsEveryone(t, planted), "a refused directory is left as it was")

	require.NoError(t, os.WriteFile(filepath.Join(rootPath, "file"), nil, 0o600))
	require.Error(t, EnsureDirIn(root, "file"))
}

// A caller that holds the container itself open, as the disk cache does with
// its root, may use a folder another account can delete or rename, below
// folders others control, but not a folder whose entries others can remove.
func TestCheckContainerHeldAcceptsAFolderHeldOpen(t *testing.T) {
	trust, err := loadTrusted()
	require.NoError(t, err)
	private := t.TempDir()
	require.NoError(t, EnsureDir(private))
	held, err := CheckContainerHeld(private)
	require.NoError(t, err)
	assert.Nil(t, held, "a folder that passes CheckContainer needs nothing held")

	home := t.TempDir()
	require.NoError(t, EnsureDir(home))
	setPrivateDACLWith(t, trust, home, "(A;OICI;0x1301bf;;;BU)")
	cache := filepath.Join(home, "cache")
	require.NoError(t, os.Mkdir(cache, 0o700))
	require.Error(t, CheckContainer(cache))
	held, err = CheckContainerHeld(cache)
	require.NoError(t, err)
	require.NotNil(t, held)
	require.Error(t, os.Rename(cache, cache+"-moved"), "the held folder cannot be renamed")
	require.NoError(t, held.Close())
	require.NoError(t, os.Rename(cache, cache+"-moved"), "the folder is released once closed")

	shared := filepath.Join(home, "shared")
	require.NoError(t, Mkdir(shared))
	setPrivateDACLWith(t, trust, shared, "(A;;0x40;;;BU)")
	_, err = CheckContainerHeld(shared)
	require.ErrorContains(t, err, "delete files or subfolders")
}

// A held directory keeps in place only the directories it is actually inside.
// A folder above a junction on the way is not one of them, so delete access to
// it is refused even for a held directory.
func TestMkdirTempHeldRefusesDeleteAccessAboveALink(t *testing.T) {
	trust, err := loadTrusted()
	require.NoError(t, err)
	base := t.TempDir()
	require.NoError(t, EnsureDir(base))
	outer := filepath.Join(base, "outer")
	require.NoError(t, Mkdir(outer))
	target := filepath.Join(base, "target")
	require.NoError(t, Mkdir(target))
	link := filepath.Join(outer, "link")
	out, err := exec.Command("cmd", "/c", "mklink", "/J", link, target).CombinedOutput()
	require.NoError(t, err, string(out))
	setPrivateDACLWith(t, trust, outer, "(A;;0x1301bf;;;BU)")

	_, _, err = MkdirTempHeld(link, "store-*")
	require.Error(t, err)
	// Windows itself refuses to follow a junction a standard user created
	// (ERROR_UNTRUSTED_MOUNT_POINT), so only an administrator run reaches the
	// folder above it.
	const errorUntrustedMountPoint = windows.Errno(448)
	if !errors.Is(err, errorUntrustedMountPoint) {
		assert.ErrorContains(t, err, "delete or rename the folder")
		assert.ErrorContains(t, err, outer)
	}
	entries, err := os.ReadDir(target)
	require.NoError(t, err)
	assert.Empty(t, entries, "nothing is created behind a refused path")

	_, held, err := MkdirTempHeld(target, "store-*")
	require.NoError(t, err, "the same folder reached directly is not affected")
	held.Close()
}

func TestCheckContainerRejectsAnUnresolvedAccountWithActionableError(t *testing.T) {
	sid, err := windows.StringToSid("S-1-5-21-1-2-3-999999")
	require.NoError(t, err)
	dir := t.TempDir()
	require.NoError(t, EnsureDir(dir))
	acl, err := windows.ACLFromEntries([]windows.EXPLICIT_ACCESS{{
		AccessPermissions: fileDeleteChild,
		AccessMode:        windows.GRANT_ACCESS,
		Trustee:           windows.TRUSTEE{TrusteeForm: windows.TRUSTEE_IS_SID, TrusteeValue: windows.TrusteeValueFromSID(sid)},
	}}, privateACLForTest(t))
	require.NoError(t, err)
	require.NoError(t, windows.SetNamedSecurityInfo(dir, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION|windows.PROTECTED_DACL_SECURITY_INFORMATION, nil, nil, acl, nil))
	err = CheckContainer(dir)
	require.ErrorContains(t, err, "cannot verify Windows permissions")
	assert.ErrorContains(t, err, sid.String())
	assert.ErrorContains(t, err, dir)
	assert.ErrorContains(t, err, "the app is running as")
	assert.ErrorContains(t, err, "ask your IT administrator")
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
