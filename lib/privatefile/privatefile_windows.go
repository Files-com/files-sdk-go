//go:build windows

package privatefile

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"unsafe"

	"golang.org/x/sys/windows"
)

// On Windows, privacy is a discretionary access control list (DACL). A private
// entry has a protected DACL, so nothing is inherited from the directory above
// it, that grants access only to the owning account, SYSTEM and Administrators.
// Directories created this way pass the same list down to what is created
// inside them, so files that inherit are private as well.

// trustedSIDs are the accounts that may own or modify the directories above
// private storage without making it replaceable: the current user, and the
// system principals that own the standard directory tree.
type trustedSIDs struct {
	user             *windows.SID
	system           *windows.SID
	administrators   *windows.SID
	trustedInstaller *windows.SID
	creatorOwner     *windows.SID
	ownerRights      *windows.SID
}

var (
	trustedOnce sync.Once
	trusted     trustedSIDs
	trustedErr  error

	privateACLOnce sync.Once
	privateACL     *windows.ACL
	privateACLErr  error

	inheritOnlyOnce sync.Once
	inheritOnlyACL  *windows.ACL
	inheritOnlyErr  error
)

func loadTrusted() (trustedSIDs, error) {
	trustedOnce.Do(func() {
		user, err := windows.GetCurrentProcessToken().GetTokenUser()
		if err != nil {
			trustedErr = fmt.Errorf("privatefile: cannot determine the current user: %w", err)
			return
		}
		system, err := windows.CreateWellKnownSid(windows.WinLocalSystemSid)
		if err != nil {
			trustedErr = err
			return
		}
		administrators, err := windows.CreateWellKnownSid(windows.WinBuiltinAdministratorsSid)
		if err != nil {
			trustedErr = err
			return
		}
		creatorOwner, err := windows.CreateWellKnownSid(windows.WinCreatorOwnerSid)
		if err != nil {
			trustedErr = err
			return
		}
		ownerRights, err := windows.CreateWellKnownSid(windows.WinCreatorOwnerRightsSid)
		if err != nil {
			trustedErr = err
			return
		}
		trustedInstaller, err := windows.StringToSid("S-1-5-80-956008885-3418522649-1831038044-1853292631-2271478464")
		if err != nil {
			trustedErr = err
			return
		}
		trusted = trustedSIDs{
			user:             user.User.Sid,
			system:           system,
			administrators:   administrators,
			trustedInstaller: trustedInstaller,
			creatorOwner:     creatorOwner,
			ownerRights:      ownerRights,
		}
	})
	return trusted, trustedErr
}

// isTrustedOwner reports whether sid may own a directory above private
// storage or an entry inside it.
func (t trustedSIDs) isTrustedOwner(sid *windows.SID) bool {
	return sid.Equals(t.user) || sid.Equals(t.system) || sid.Equals(t.administrators) || sid.Equals(t.trustedInstaller)
}

// mayModifyContainer reports whether an allow entry for sid on a directory
// above private storage is acceptable. CREATOR OWNER and OWNER RIGHTS entries
// describe the creator of an entry, which for private storage is this user.
func (t trustedSIDs) mayModifyContainer(sid *windows.SID) bool {
	return t.isTrustedOwner(sid) || sid.Equals(t.creatorOwner) || sid.Equals(t.ownerRights)
}

// privateDACL grants full control to the current user, SYSTEM and
// Administrators, and passes the same entries down to created children.
func loadPrivateACL() (*windows.ACL, error) {
	privateACLOnce.Do(func() {
		t, err := loadTrusted()
		if err != nil {
			privateACLErr = err
			return
		}
		var entries []windows.EXPLICIT_ACCESS
		for _, sid := range []*windows.SID{t.user, t.system, t.administrators} {
			entries = append(entries, windows.EXPLICIT_ACCESS{
				AccessPermissions: windows.GENERIC_ALL,
				AccessMode:        windows.GRANT_ACCESS,
				Inheritance:       windows.SUB_CONTAINERS_AND_OBJECTS_INHERIT,
				Trustee: windows.TRUSTEE{
					TrusteeForm:  windows.TRUSTEE_IS_SID,
					TrusteeType:  windows.TRUSTEE_IS_UNKNOWN,
					TrusteeValue: windows.TrusteeValueFromSID(sid),
				},
			})
		}
		privateACL, privateACLErr = windows.ACLFromEntries(entries, nil)
	})
	return privateACL, privateACLErr
}

// loadInheritOnlyACL is an empty list: applied unprotected, it leaves an entry
// with exactly the entries it inherits from its directory.
func loadInheritOnlyACL() (*windows.ACL, error) {
	inheritOnlyOnce.Do(func() {
		sd, err := windows.SecurityDescriptorFromString("D:")
		if err != nil {
			inheritOnlyErr = err
			return
		}
		inheritOnlyACL, _, inheritOnlyErr = sd.DACL()
	})
	return inheritOnlyACL, inheritOnlyErr
}

// privateSecurityDescriptor is the protected private DACL in the form the
// creation APIs take.
func privateSecurityDescriptor() (*windows.SECURITY_DESCRIPTOR, error) {
	acl, err := loadPrivateACL()
	if err != nil {
		return nil, err
	}
	sd, err := windows.NewSecurityDescriptor()
	if err != nil {
		return nil, err
	}
	if err := sd.SetDACL(acl, true, false); err != nil {
		return nil, err
	}
	if err := sd.SetControl(windows.SE_DACL_PROTECTED, windows.SE_DACL_PROTECTED); err != nil {
		return nil, err
	}
	return sd, nil
}

// Mkdir creates a directory only the current user can access. It fails if
// path already exists.
func Mkdir(path string) error {
	sd, err := privateSecurityDescriptor()
	if err != nil {
		return &fs.PathError{Op: "mkdir", Path: path, Err: err}
	}
	name, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return &fs.PathError{Op: "mkdir", Path: path, Err: err}
	}
	sa := windows.SecurityAttributes{Length: uint32(unsafe.Sizeof(windows.SecurityAttributes{})), SecurityDescriptor: sd}
	if err := windows.CreateDirectory(name, &sa); err != nil {
		return &fs.PathError{Op: "mkdir", Path: path, Err: err}
	}
	return nil
}

// EnsureDir makes path a private directory this process owns: it creates a
// missing one, and replaces the DACL of an existing one with the protected
// private DACL. Windows then re-derives the inherited entries of everything
// below it; explicit entries below it are handled by EnsureEntryPrivate.
//
// An existing path must be a real directory (not a reparse point) owned by the
// current user or the system. Anything else is refused rather than used.
func EnsureDir(path string) error {
	err := Mkdir(path)
	if err == nil || !errors.Is(err, fs.ErrExist) {
		return err
	}
	if _, err := lstatDirectory(path); err != nil {
		return err
	}
	handle, err := openForSecurity(path, windows.READ_CONTROL|windows.WRITE_DAC)
	if err != nil {
		return &fs.PathError{Op: "ensuredir", Path: path, Err: err}
	}
	defer windows.CloseHandle(handle)
	sec, err := readSecurity(handle)
	if err != nil {
		return &fs.PathError{Op: "ensuredir", Path: path, Err: err}
	}
	t, err := loadTrusted()
	if err != nil {
		return err
	}
	if !t.isTrustedOwner(sec.owner) {
		return &fs.PathError{Op: "ensuredir", Path: path, Err: fmt.Errorf("owned by %s, not by this user or the system", sec.owner)}
	}
	acl, err := loadPrivateACL()
	if err != nil {
		return err
	}
	if err := windows.SetSecurityInfo(handle, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION|windows.PROTECTED_DACL_SECURITY_INFORMATION, nil, nil, acl, nil); err != nil {
		return &fs.PathError{Op: "ensuredir", Path: path, Err: err}
	}
	return nil
}

// CheckEntry reports whether one entry inside a directory this package is
// about to make private may be made private along with it. info must be the
// Lstat of path. Only entries this package could have created qualify: an
// ordinary file or directory owned by this user or the system, with a single
// name. A reparse point is refused: it does not belong in a private tree this
// package created, and changing its target's DACL could affect unrelated data.
// A file with several names is refused because its DACL is shared with a name
// somewhere this process does not control, and a directory's inheritable
// entries would reach it too. Nothing is modified.
func CheckEntry(path string, info fs.FileInfo) error {
	if isLink(info) {
		return &fs.PathError{Op: "setacl", Path: path, Err: errors.New("unexpected reparse point inside a private directory")}
	}
	handle, err := openForSecurity(path, windows.READ_CONTROL)
	if err != nil {
		return &fs.PathError{Op: "setacl", Path: path, Err: err}
	}
	defer windows.CloseHandle(handle)
	return checkEntryHandle(path, info, handle)
}

// checkEntryHandle is CheckEntry for an entry already opened for reading its
// security and attributes.
func checkEntryHandle(path string, info fs.FileInfo, handle windows.Handle) error {
	sec, err := readSecurity(handle)
	if err != nil {
		return &fs.PathError{Op: "setacl", Path: path, Err: err}
	}
	t, err := loadTrusted()
	if err != nil {
		return err
	}
	if !t.isTrustedOwner(sec.owner) {
		return &fs.PathError{Op: "setacl", Path: path, Err: fmt.Errorf("owned by %s, not by this user or the system", sec.owner)}
	}
	if !info.IsDir() {
		var attributes windows.ByHandleFileInformation
		if err := windows.GetFileInformationByHandle(handle, &attributes); err != nil {
			return &fs.PathError{Op: "setacl", Path: path, Err: err}
		}
		if attributes.NumberOfLinks > 1 {
			return &fs.PathError{Op: "setacl", Path: path, Err: fmt.Errorf("file has %d names; a private file must have only this one", attributes.NumberOfLinks)}
		}
	}
	return nil
}

// EnsureEntryPrivate makes one entry inside a directory that EnsureDir made
// private carry only the entries it inherits from that directory, after
// CheckEntry accepts it. info must be the Lstat of path. A protected DACL or
// an explicit entry on the file would otherwise keep granting whatever access
// it granted before, whatever the directory above it says: Windows checks the
// file's own DACL when a file is opened by a known path.
func EnsureEntryPrivate(path string, info fs.FileInfo) error {
	if isLink(info) {
		return &fs.PathError{Op: "setacl", Path: path, Err: errors.New("unexpected reparse point inside a private directory")}
	}
	handle, err := openForSecurity(path, windows.READ_CONTROL|windows.WRITE_DAC)
	if err != nil {
		return &fs.PathError{Op: "setacl", Path: path, Err: err}
	}
	defer windows.CloseHandle(handle)
	if err := checkEntryHandle(path, info, handle); err != nil {
		return err
	}
	sec, err := readSecurity(handle)
	if err != nil {
		return &fs.PathError{Op: "setacl", Path: path, Err: err}
	}
	if !sec.protected && !sec.hasExplicitEntries() {
		return nil
	}
	inherit, err := loadInheritOnlyACL()
	if err != nil {
		return err
	}
	if err := windows.SetSecurityInfo(handle, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION|windows.UNPROTECTED_DACL_SECURITY_INFORMATION, nil, nil, inherit, nil); err != nil {
		return &fs.PathError{Op: "setacl", Path: path, Err: err}
	}
	return nil
}

// checkLinkComponent refuses a reparse point on the way to a container unless
// this user or the system owns it: anyone else who owns a link can point it
// at a tree they control.
func checkLinkComponent(path string) error {
	handle, err := openForSecurity(path, windows.READ_CONTROL)
	if err != nil {
		return &fs.PathError{Op: "checkcontainer", Path: path, Err: err}
	}
	defer windows.CloseHandle(handle)
	sec, err := readSecurity(handle)
	if err != nil {
		return &fs.PathError{Op: "checkcontainer", Path: path, Err: err}
	}
	t, err := loadTrusted()
	if err != nil {
		return err
	}
	return t.checkContainerOwner(path, sec.owner)
}

// fileDeleteChild is the directory right to delete or rename any entry, which
// x/sys/windows does not name.
const fileDeleteChild = 0x40

// containerModifyRights are the rights that let an account replace or remove
// an existing entry of a directory, or grant itself the ability to: deleting
// any child, deleting or renaming the directory itself, or changing its
// security. Adding new entries (FILE_ADD_FILE, FILE_ADD_SUBDIRECTORY, generic
// write) is not among them: a new name beside private storage cannot replace
// it, and the standard drive root grants Authenticated Users exactly that.
const containerModifyRights = fileDeleteChild | windows.DELETE | windows.WRITE_DAC | windows.WRITE_OWNER | windows.GENERIC_ALL

// checkContainerComponent refuses a directory an untrusted account owns or
// can replace entries in. Verified local administrators are trusted.
//
// A directory below path that the caller holds open (heldBelow) keeps path in
// place, so path is not examined; see heldDirKeepsParents. A directory the
// caller holds open itself (heldHere) cannot be deleted or renamed by anyone
// while it is open, so DELETE on it is accepted, but the rights to remove its
// entries or change its security are still refused.
func checkContainerComponent(path string, hold containerHold) error {
	if hold == heldBelow {
		return nil
	}
	modifyRights := uint32(containerModifyRights)
	if hold == heldHere {
		modifyRights &^= windows.DELETE
	}
	handle, err := openForSecurity(path, windows.READ_CONTROL)
	if err != nil {
		return &fs.PathError{Op: "checkcontainer", Path: path, Err: err}
	}
	defer windows.CloseHandle(handle)
	sec, err := readSecurity(handle)
	if err != nil {
		return &fs.PathError{Op: "checkcontainer", Path: path, Err: err}
	}
	t, err := loadTrusted()
	if err != nil {
		return err
	}
	if err := t.checkContainerOwner(path, sec.owner); err != nil {
		return err
	}
	if sec.dacl == nil {
		return t.containerError(path, errors.New("folder has no access restrictions; any account could replace private storage below it"))
	}
	checked := make(map[string]bool)
	for _, ace := range sec.entries() {
		if ace.header.AceFlags&windows.INHERIT_ONLY_ACE != 0 {
			continue
		}
		if ace.header.AceType != windows.ACCESS_ALLOWED_ACE_TYPE && ace.header.AceType != windows.ACCESS_DENIED_ACE_TYPE {
			if uint32(ace.mask)&modifyRights != 0 {
				return t.containerError(path, fmt.Errorf("cannot verify Windows permissions: unsupported access entry type %d", ace.header.AceType))
			}
			continue
		}
		if ace.header.AceType != windows.ACCESS_ALLOWED_ACE_TYPE {
			continue
		}
		if uint32(ace.mask)&modifyRights == 0 || t.mayModifyContainer(ace.sid) {
			continue
		}
		if checked[ace.sid.String()] {
			continue
		}
		checked[ace.sid.String()] = true
		// Deny entries are evaluated before the account is looked up, so they
		// also apply to an account whose name cannot be resolved.
		rights, err := containerAccess(sec, ace.sid)
		if err != nil {
			return t.containerError(path, fmt.Errorf("cannot verify Windows permissions for %s: %w", ace.sid, err))
		}
		if rights&modifyRights == 0 {
			continue
		}
		account, administrator, err := t.containerAccount(ace.sid)
		if err != nil {
			return t.containerError(path, fmt.Errorf("cannot verify Windows permissions for %s: %w", account, err))
		}
		if administrator {
			continue
		}
		return t.containerError(path, fmt.Errorf("folder lets %s %s; private storage below it could be replaced by that account", account, containerRightsNames(rights&modifyRights)))
	}
	return nil
}

// heldDirKeepsParents reports whether f, a directory held open, keeps every
// directory above it in place. NTFS does not rename or move a directory while
// anything inside it is open, whatever the rights of the account trying, and
// f holds no share-delete access, so f itself cannot be deleted or renamed
// either. That is relied on only for local NTFS volumes, where it was
// verified; on a network share or another file system the directories above
// are checked as if nothing were held.
func heldDirKeepsParents(f *os.File) bool {
	h := windows.Handle(f.Fd())
	path := make([]uint16, windows.MAX_PATH)
	for {
		// Flags 0: FILE_NAME_NORMALIZED and VOLUME_NAME_DOS.
		n, err := windows.GetFinalPathNameByHandle(h, &path[0], uint32(len(path)), 0)
		if err != nil {
			return false
		}
		if n < uint32(len(path)) {
			path = path[:n]
			break
		}
		path = make([]uint16, n)
	}
	if strings.HasPrefix(windows.UTF16ToString(path), `\\?\UNC\`) {
		return false
	}
	var fileSystem [windows.MAX_PATH + 1]uint16
	if err := windows.GetVolumeInformationByHandle(h, nil, 0, nil, nil, nil, &fileSystem[0], uint32(len(fileSystem))); err != nil {
		return false
	}
	return windows.UTF16ToString(fileSystem[:]) == "NTFS"
}

// containerAccount trusts individual local administrators on the same basis
// as the Administrators group. A failed lookup never establishes trust.
func (t trustedSIDs) containerAccount(sid *windows.SID) (string, bool, error) {
	account, domain, kind, err := sid.LookupAccount("")
	if err != nil {
		return sid.String(), false, fmt.Errorf("cannot resolve the account name: %w", err)
	}
	name := accountLabel(sid, account, domain)
	if kind != windows.SidTypeUser {
		return name, false, nil
	}
	accountUser, err := user.LookupId(sid.String())
	if err != nil {
		return name, false, fmt.Errorf("cannot look up the Windows account: %w", err)
	}
	groups, err := accountUser.GroupIds()
	if err != nil {
		return name, false, fmt.Errorf("cannot verify local administrator membership: %w", err)
	}
	return name, slices.Contains(groups, t.administrators.String()), nil
}

func (t trustedSIDs) checkContainerOwner(path string, owner *windows.SID) error {
	if t.isTrustedOwner(owner) {
		return nil
	}
	account, administrator, err := t.containerAccount(owner)
	if err != nil {
		return t.containerError(path, fmt.Errorf("cannot verify the folder owner %s: %w", account, err))
	}
	if administrator {
		return nil
	}
	return t.containerError(path, fmt.Errorf("folder is owned by %s; that account could replace private storage below it", account))
}

func (t trustedSIDs) containerError(path string, err error) error {
	account, domain, _, lookupErr := t.user.LookupAccount("")
	appAccount := t.user.String()
	if lookupErr == nil {
		appAccount = accountLabel(t.user, account, domain)
	}
	return &fs.PathError{Op: "checkcontainer", Path: path, Err: fmt.Errorf("%w; the app is running as %s; ask your IT administrator to verify these accounts and remove any unintended access from this folder or its parent folders", err, appAccount)}
}

func accountLabel(sid *windows.SID, account string, domain string) string {
	if domain != "" {
		account = domain + `\` + account
	}
	return fmt.Sprintf("%s (%s)", account, sid)
}

func containerRightsNames(rights uint32) string {
	var names []string
	for _, right := range []struct {
		mask uint32
		name string
	}{
		{fileDeleteChild, "delete files or subfolders"},
		{windows.DELETE, "delete or rename the folder"},
		{windows.WRITE_DAC, "change folder permissions"},
		{windows.WRITE_OWNER, "change the folder owner"},
	} {
		if rights&right.mask != 0 {
			names = append(names, right.name)
		}
	}
	return strings.Join(names, ", ")
}

// OpenFile opens name inside root as a file only the current user can read.
// A file it creates has the protected private DACL from the moment it exists;
// a file that already existed receives it before OpenFile returns, so no data
// written through the returned file is ever exposed. perm is accepted for
// symmetry with the POSIX implementation; Windows takes no mode.
//
// The file is opened relative to a handle of its directory obtained through
// root, so the confinement root provides is kept, and the final component is
// never followed as a reparse point.
func OpenFile(root *os.Root, name string, flag int, perm fs.FileMode) (*os.File, error) {
	if perm&^0o700 != 0 {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fmt.Errorf("private file mode %04o grants access beyond the owner", perm)}
	}
	var access uint32 = windows.SYNCHRONIZE | windows.READ_CONTROL | windows.WRITE_DAC | windows.FILE_READ_ATTRIBUTES | windows.FILE_READ_EA
	switch flag & (os.O_RDONLY | os.O_WRONLY | os.O_RDWR) {
	case os.O_WRONLY:
		access |= windows.FILE_GENERIC_WRITE
	case os.O_RDWR:
		access |= windows.FILE_GENERIC_READ | windows.FILE_GENERIC_WRITE
	default:
		access |= windows.FILE_GENERIC_READ
	}
	if flag&os.O_CREATE != 0 {
		access |= windows.FILE_GENERIC_WRITE
	}
	var disposition uint32
	switch {
	case flag&(os.O_CREATE|os.O_EXCL) == os.O_CREATE|os.O_EXCL:
		disposition = windows.FILE_CREATE
	case flag&os.O_CREATE != 0:
		disposition = windows.FILE_OPEN_IF
	default:
		disposition = windows.FILE_OPEN
	}
	sd, err := privateSecurityDescriptor()
	if err != nil {
		return nil, &fs.PathError{Op: "open", Path: name, Err: err}
	}
	handle, opened, err := openRelative(root, name, access, disposition, sd)
	if err != nil {
		return nil, err
	}
	if opened {
		acl, err := loadPrivateACL()
		if err == nil {
			err = windows.SetSecurityInfo(handle, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION|windows.PROTECTED_DACL_SECURITY_INFORMATION, nil, nil, acl, nil)
		}
		if err != nil {
			windows.CloseHandle(handle)
			return nil, &fs.PathError{Op: "open", Path: filepath.Join(root.Name(), name), Err: err}
		}
	}
	if flag&os.O_TRUNC != 0 {
		if err := windows.Ftruncate(handle, 0); err != nil {
			windows.CloseHandle(handle)
			return nil, &fs.PathError{Op: "open", Path: filepath.Join(root.Name(), name), Err: err}
		}
	}
	return os.NewFile(uintptr(handle), filepath.Join(root.Name(), name)), nil
}

// openRelative opens the final element of name relative to a handle of its
// directory, which is opened through root so the root's confinement applies.
// It reports whether an existing file was opened rather than created.
func openRelative(root *os.Root, name string, access uint32, disposition uint32, sd *windows.SECURITY_DESCRIPTOR) (windows.Handle, bool, error) {
	dir, base := filepath.Split(name)
	if dir == "" {
		dir = "."
	}
	dirFile, err := root.Open(dir)
	if err != nil {
		return 0, false, err
	}
	defer dirFile.Close()
	objectName, err := windows.NewNTUnicodeString(base)
	if err != nil {
		return 0, false, &fs.PathError{Op: "open", Path: filepath.Join(root.Name(), name), Err: err}
	}
	oa := windows.OBJECT_ATTRIBUTES{
		RootDirectory:      windows.Handle(dirFile.Fd()),
		ObjectName:         objectName,
		Attributes:         windows.OBJ_CASE_INSENSITIVE,
		SecurityDescriptor: sd,
	}
	oa.Length = uint32(unsafe.Sizeof(oa))
	var handle windows.Handle
	var iosb windows.IO_STATUS_BLOCK
	err = windows.NtCreateFile(
		&handle,
		access,
		&oa,
		&iosb,
		nil,
		windows.FILE_ATTRIBUTE_NORMAL,
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE|windows.FILE_SHARE_DELETE,
		disposition,
		windows.FILE_SYNCHRONOUS_IO_NONALERT|windows.FILE_NON_DIRECTORY_FILE|windows.FILE_OPEN_REPARSE_POINT,
		0,
		0,
	)
	if err != nil {
		var errno error = err
		if status, ok := err.(windows.NTStatus); ok {
			errno = status.Errno()
		}
		return 0, false, &fs.PathError{Op: "open", Path: filepath.Join(root.Name(), name), Err: errno}
	}
	const fileOpened = 1 // FILE_OPENED: an existing file was opened
	return handle, iosb.Information == fileOpened, nil
}

// openContainer opens dir, following links, to create directories inside it
// through the handle. Other processes may still delete or rename dir: what
// keeps it in place is a directory then created inside it and held open.
func openContainer(dir string) (*os.File, error) {
	name, err := windows.UTF16PtrFromString(dir)
	if err != nil {
		return nil, &fs.PathError{Op: "open", Path: dir, Err: err}
	}
	handle, err := windows.CreateFile(name, windows.FILE_GENERIC_READ, windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE|windows.FILE_SHARE_DELETE, nil, windows.OPEN_EXISTING, windows.FILE_FLAG_BACKUP_SEMANTICS, 0)
	if err != nil {
		return nil, &fs.PathError{Op: "open", Path: dir, Err: err}
	}
	return os.NewFile(uintptr(handle), dir), nil
}

// createDirIn creates or opens the directory name inside the directory open
// as parent, through the handle rather than a path, and returns it open. A
// directory it creates has the protected private DACL from the moment it
// exists. The handle does not share delete access, so while it is open nobody
// can delete or rename the directory. created reports that it did not exist.
func createDirIn(parent *os.File, name string, disposition uint32) (handle windows.Handle, created bool, err error) {
	path := filepath.Join(parent.Name(), name)
	sd, err := privateSecurityDescriptor()
	if err != nil {
		return 0, false, &fs.PathError{Op: "mkdir", Path: path, Err: err}
	}
	objectName, err := windows.NewNTUnicodeString(name)
	if err != nil {
		return 0, false, &fs.PathError{Op: "mkdir", Path: path, Err: err}
	}
	oa := windows.OBJECT_ATTRIBUTES{
		RootDirectory:      windows.Handle(parent.Fd()),
		ObjectName:         objectName,
		Attributes:         windows.OBJ_CASE_INSENSITIVE,
		SecurityDescriptor: sd,
	}
	oa.Length = uint32(unsafe.Sizeof(oa))
	var iosb windows.IO_STATUS_BLOCK
	err = windows.NtCreateFile(
		&handle,
		windows.FILE_GENERIC_READ,
		&oa,
		&iosb,
		nil,
		0,
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE,
		disposition,
		windows.FILE_SYNCHRONOUS_IO_NONALERT|windows.FILE_DIRECTORY_FILE|windows.FILE_OPEN_REPARSE_POINT,
		0,
		0,
	)
	if err != nil {
		var errno error = err
		if status, ok := err.(windows.NTStatus); ok {
			errno = status.Errno()
		}
		return 0, false, &fs.PathError{Op: "mkdir", Path: path, Err: errno}
	}
	const fileCreated = 2 // FILE_CREATED: the directory did not exist
	return handle, iosb.Information == fileCreated, nil
}

// mkdirHeld creates name inside container, private and held open, in one step.
func mkdirHeld(container *os.File, name string) (*os.File, error) {
	handle, _, err := createDirIn(container, name, windows.FILE_CREATE)
	if err != nil {
		return nil, err
	}
	return os.NewFile(uintptr(handle), filepath.Join(container.Name(), name)), nil
}

// ensureDirIn is EnsureDirIn: name is created private inside root, or opened
// there and accepted only when it is a real directory that is already private.
func ensureDirIn(root *os.File, name string) error {
	handle, created, err := createDirIn(root, name, windows.FILE_OPEN_IF)
	if err != nil {
		return err
	}
	defer windows.CloseHandle(handle)
	if created {
		return nil
	}
	path := filepath.Join(root.Name(), name)
	var info windows.ByHandleFileInformation
	if err := windows.GetFileInformationByHandle(handle, &info); err != nil {
		return &fs.PathError{Op: "ensuredir", Path: path, Err: err}
	}
	if info.FileAttributes&windows.FILE_ATTRIBUTE_REPARSE_POINT != 0 {
		return &fs.PathError{Op: "ensuredir", Path: path, Err: errors.New("must be a directory, but it is a reparse point")}
	}
	sec, err := readSecurity(handle)
	if err != nil {
		return &fs.PathError{Op: "ensuredir", Path: path, Err: err}
	}
	t, err := loadTrusted()
	if err != nil {
		return err
	}
	if !t.isTrustedOwner(sec.owner) {
		return &fs.PathError{Op: "ensuredir", Path: path, Err: fmt.Errorf("owned by %s, not by this user or the system", sec.owner)}
	}
	if sec.dacl == nil || !sec.protected {
		return &fs.PathError{Op: "ensuredir", Path: path, Err: errors.New("existing directory is not private, so another account may have put it here")}
	}
	for _, ace := range sec.entries() {
		if ace.header.AceType == windows.ACCESS_ALLOWED_ACE_TYPE && !t.isTrustedOwner(ace.sid) {
			return &fs.PathError{Op: "ensuredir", Path: path, Err: fmt.Errorf("existing directory grants access to %s, so another account may have put it here", ace.sid)}
		}
	}
	return nil
}

// openForSecurity opens path for reading or changing its security and reading
// its attributes, without following a reparse point at path, so the entry
// named is the one examined.
func openForSecurity(path string, access uint32) (windows.Handle, error) {
	name, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return 0, err
	}
	return windows.CreateFile(name, access|windows.FILE_READ_ATTRIBUTES, windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE|windows.FILE_SHARE_DELETE, nil, windows.OPEN_EXISTING, windows.FILE_FLAG_BACKUP_SEMANTICS|windows.FILE_FLAG_OPEN_REPARSE_POINT, 0)
}

// security is the owner and DACL of an entry.
type security struct {
	sd        *windows.SECURITY_DESCRIPTOR
	owner     *windows.SID
	dacl      *windows.ACL // nil means no DACL: every account has full access
	protected bool
}

type aclEntry struct {
	header windows.ACE_HEADER
	mask   windows.ACCESS_MASK
	sid    *windows.SID
}

func readSecurity(handle windows.Handle) (security, error) {
	sd, err := windows.GetSecurityInfo(handle, windows.SE_FILE_OBJECT, windows.OWNER_SECURITY_INFORMATION|windows.DACL_SECURITY_INFORMATION)
	if err != nil {
		return security{}, err
	}
	owner, _, err := sd.Owner()
	if err != nil {
		return security{}, err
	}
	dacl, _, err := sd.DACL()
	if err != nil {
		return security{}, err
	}
	control, _, err := sd.Control()
	if err != nil {
		return security{}, err
	}
	return security{sd: sd, owner: owner, dacl: dacl, protected: control&windows.SE_DACL_PROTECTED != 0}, nil
}

// entries lists the access control entries of the DACL.
func (s security) entries() []aclEntry {
	if s.dacl == nil {
		return nil
	}
	var entries []aclEntry
	for i := uint32(0); i < uint32(s.dacl.AceCount); i++ {
		var ace *windows.ACCESS_ALLOWED_ACE
		if err := windows.GetAce(s.dacl, i, &ace); err != nil {
			continue
		}
		entries = append(entries, aclEntry{header: ace.Header, mask: ace.Mask, sid: (*windows.SID)(unsafe.Pointer(&ace.SidStart))})
	}
	return entries
}

// hasExplicitEntries reports whether any entry was set on the object itself
// rather than inherited from its directory.
func (s security) hasExplicitEntries() bool {
	for _, ace := range s.entries() {
		if ace.header.AceFlags&windows.INHERITED_ACE == 0 {
			return true
		}
	}
	return false
}

// publicationPlatform on Windows is the DACL the published file receives:
// the one the replaced file had, or, when there was none, only what the
// destination directory hands down.
type publicationPlatform struct {
	existing *security // nil: inherit from the destination directory
}

func publicationFor(root *os.Root, name string, _ string, _ fs.FileMode) (publicationPlatform, error) {
	existing, err := root.Open(name)
	if errors.Is(err, fs.ErrNotExist) {
		return publicationPlatform{}, nil
	}
	if err != nil {
		return publicationPlatform{}, err
	}
	defer existing.Close()
	info, err := existing.Stat()
	if err != nil {
		return publicationPlatform{}, err
	}
	if !info.Mode().IsRegular() {
		return publicationPlatform{}, nil
	}
	sec, err := readSecurity(windows.Handle(existing.Fd()))
	if err != nil {
		return publicationPlatform{}, &fs.PathError{Op: "getacl", Path: filepath.Join(root.Name(), name), Err: err}
	}
	return publicationPlatform{existing: &sec}, nil
}

// applyToPublished sets the published DACL on the file now at name, once a
// handle to it is checked to be the published file: the replaced file's DACL,
// protected or not as it was, or an empty unprotected DACL so the file has
// exactly what the destination directory hands down.
func (p publicationPlatform) applyToPublished(root *os.Root, name string, published fs.FileInfo) error {
	handle, _, err := openRelative(root, name, windows.SYNCHRONIZE|windows.READ_CONTROL|windows.WRITE_DAC|windows.FILE_READ_ATTRIBUTES, windows.FILE_OPEN, nil)
	if err != nil {
		return err
	}
	f := os.NewFile(uintptr(handle), filepath.Join(root.Name(), name))
	defer f.Close()
	got, err := f.Stat()
	if err != nil {
		return err
	}
	if published == nil || !os.SameFile(published, got) {
		return &fs.PathError{Op: "setacl", Path: filepath.Join(root.Name(), name), Err: errPublishedReplaced}
	}
	var acl *windows.ACL
	var flags windows.SECURITY_INFORMATION = windows.DACL_SECURITY_INFORMATION
	if p.existing != nil {
		acl = p.existing.dacl
		if p.existing.protected {
			flags |= windows.PROTECTED_DACL_SECURITY_INFORMATION
		} else {
			flags |= windows.UNPROTECTED_DACL_SECURITY_INFORMATION
		}
	} else {
		acl, err = loadInheritOnlyACL()
		if err != nil {
			return err
		}
		flags |= windows.UNPROTECTED_DACL_SECURITY_INFORMATION
	}
	if err := windows.SetSecurityInfo(windows.Handle(f.Fd()), windows.SE_FILE_OBJECT, flags, nil, nil, acl, nil); err != nil {
		return &fs.PathError{Op: "setacl", Path: filepath.Join(root.Name(), name), Err: err}
	}
	return nil
}
