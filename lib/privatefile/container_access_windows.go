//go:build windows

package privatefile

import (
	"fmt"
	"runtime"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	authzDLL                       = windows.NewLazySystemDLL("authz.dll")
	authzInitializeResourceManager = authzDLL.NewProc("AuthzInitializeResourceManager")
	authzInitializeContextFromSid  = authzDLL.NewProc("AuthzInitializeContextFromSid")
	authzAccessCheck               = authzDLL.NewProc("AuthzAccessCheck")
	authzFreeContext               = authzDLL.NewProc("AuthzFreeContext")
	authzFreeResourceManager       = authzDLL.NewProc("AuthzFreeResourceManager")
	mapGenericMask                 = windows.NewLazySystemDLL("advapi32.dll").NewProc("MapGenericMask")
)

// FILE_ALL_ACCESS from WinNT.h includes every file-specific access right.
var fileGenericMapping = struct{ read, write, execute, all uint32 }{
	windows.FILE_GENERIC_READ,
	windows.FILE_GENERIC_WRITE,
	windows.FILE_GENERIC_EXECUTE,
	windows.STANDARD_RIGHTS_REQUIRED | windows.SYNCHRONIZE | 0x1ff,
}

type authzAccessRequest struct {
	desiredAccess        uint32
	principalSelfSID     *windows.SID
	objectTypeList       uintptr
	objectTypeListLength uint32
	optionalArguments    uintptr
}

type authzAccessReply struct {
	resultListLength      uint32
	grantedAccessMask     *uint32
	saclEvaluationResults *uint32
	errorCode             *uint32
}

// containerAccess evaluates the ordered allow and deny entries for this SID.
// Group expansion is deliberately disabled: a deny for some other group might
// not apply to every member of this SID, or in every kind of logon. Ignoring
// those denies can reject a safe folder, but cannot approve an unsafe one.
func containerAccess(sec security, sid *windows.SID) (uint32, error) {
	// Authz expects file-specific rights. Work on a copy so mapping generic
	// rights changes neither the caller's descriptor nor the folder's ACL.
	sd, err := sec.sd.ToAbsolute()
	if err != nil {
		return 0, err
	}
	acl, _, err := sd.DACL()
	if err != nil {
		return 0, err
	}
	if err := mapGenericMask.Find(); err != nil {
		return 0, err
	}
	for i := uint32(0); i < uint32(acl.AceCount); i++ {
		var ace *windows.ACCESS_ALLOWED_ACE
		if err := windows.GetAce(acl, i, &ace); err != nil {
			return 0, err
		}
		// Inherit-only entries apply to children, not to this folder.
		if ace.Header.AceFlags&windows.INHERIT_ONLY_ACE != 0 {
			continue
		}
		if ace.Header.AceType != windows.ACCESS_ALLOWED_ACE_TYPE && ace.Header.AceType != windows.ACCESS_DENIED_ACE_TYPE {
			return 0, fmt.Errorf("unsupported access entry type %d", ace.Header.AceType)
		}
		mapGenericMask.Call(uintptr(unsafe.Pointer(&ace.Mask)), uintptr(unsafe.Pointer(&fileGenericMapping)))
	}
	const (
		authzNoAudit         = 1
		authzSkipTokenGroups = 2
	)
	var manager windows.Handle
	if err := callAuthz(authzInitializeResourceManager, authzNoAudit, 0, 0, 0, 0, uintptr(unsafe.Pointer(&manager))); err != nil {
		return 0, err
	}
	defer callAuthz(authzFreeResourceManager, uintptr(manager))

	var context windows.Handle
	// The Identifier argument is an eight-byte LUID passed by value.
	if unsafe.Sizeof(uintptr(0)) == 4 {
		err = callAuthz(authzInitializeContextFromSid, authzSkipTokenGroups, uintptr(unsafe.Pointer(sid)), uintptr(manager), 0, 0, 0, 0, uintptr(unsafe.Pointer(&context)))
	} else {
		err = callAuthz(authzInitializeContextFromSid, authzSkipTokenGroups, uintptr(unsafe.Pointer(sid)), uintptr(manager), 0, 0, 0, uintptr(unsafe.Pointer(&context)))
	}
	if err != nil {
		return 0, err
	}
	runtime.KeepAlive(sid)
	defer callAuthz(authzFreeContext, uintptr(context))

	request := authzAccessRequest{desiredAccess: windows.MAXIMUM_ALLOWED}
	var granted, audit, result uint32
	reply := authzAccessReply{
		resultListLength:      1,
		grantedAccessMask:     &granted,
		saclEvaluationResults: &audit,
		errorCode:             &result,
	}
	if err := callAuthz(authzAccessCheck, 0, uintptr(context), uintptr(unsafe.Pointer(&request)), 0, uintptr(unsafe.Pointer(sd)), 0, 0, uintptr(unsafe.Pointer(&reply)), 0); err != nil {
		return 0, err
	}
	runtime.KeepAlive(sd)
	if result != 0 && result != uint32(windows.ERROR_ACCESS_DENIED) {
		return 0, windows.Errno(result)
	}
	return granted, nil
}

//go:uintptrescapes
func callAuthz(proc *windows.LazyProc, args ...uintptr) error {
	if err := proc.Find(); err != nil {
		return err
	}
	ok, _, err := proc.Call(args...)
	if ok != 0 {
		return nil
	}
	if err == windows.ERROR_SUCCESS {
		err = windows.ERROR_GEN_FAILURE
	}
	return fmt.Errorf("%s: %w", proc.Name, err)
}
