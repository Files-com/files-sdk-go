package file

import (
	"os"
	"unsafe"

	"golang.org/x/sys/windows"
)

// ownedFileIdentity is a file's FILE_ID_INFO: the serial number of its volume and its 128-bit file ID, which together
// identify the file uniquely on the computer. os.SameFile compares the 64-bit file index of BY_HANDLE_FILE_INFORMATION
// instead, which Windows does not guarantee unique on ReFS.
type ownedFileIdentity struct {
	volumeSerialNumber uint64
	fileID             [16]byte
}

// ownedFileIdentityOf returns the identity of file by its handle, if its file system gives one; an identity of zeros
// identifies nothing.
func ownedFileIdentityOf(file *os.File) (ownedFileIdentity, bool) {
	conn, err := file.SyscallConn()
	if err != nil {
		return ownedFileIdentity{}, false
	}
	var id ownedFileIdentity
	var infoErr error
	if err := conn.Control(func(handle uintptr) {
		infoErr = windows.GetFileInformationByHandleEx(windows.Handle(handle), windows.FileIdInfo, (*byte)(unsafe.Pointer(&id)), uint32(unsafe.Sizeof(id)))
	}); err != nil || infoErr != nil || id.volumeSerialNumber == 0 || id.fileID == [16]byte{} {
		return ownedFileIdentity{}, false
	}
	return id, true
}

func sameOwnedFile(a, b ownedFileIdentity) bool {
	return a == b
}
