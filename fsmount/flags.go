package fsmount

import (
	"fmt"
	"strings"

	"github.com/winfsp/cgofuse/fuse"
)

const (
	// O_EVTONLY is a flag used to indicate that the file is opened for event notifications only.
	O_EVTONLY = 0x8000
)

type FuseFlags int

// NewFuseFlags initializes a FuseFlags instance with the given integer flag value.
func NewFuseFlags(flags int) FuseFlags {
	return FuseFlags(flags)
}

// IsEventOnly checks if the flag is open for event notifications only.
func (f FuseFlags) IsEventOnly() bool {
	return int(f)&O_EVTONLY != 0
}

// IsReadOnly checks if the flag is set to read-only.
func (f FuseFlags) IsReadOnly() bool {
	return int(f) == 0
}

// IsWriteOnly checks if the flag is set to write-only.
func (f FuseFlags) IsWriteOnly() bool {
	return int(f)&fuse.O_WRONLY != 0
}

// IsReadWrite checks if the flag is set to read-write.
func (f FuseFlags) IsReadWrite() bool {
	return int(f)&fuse.O_RDWR != 0
}

// IsCreate checks if the flag is set to create.
func (f FuseFlags) IsCreate() bool {
	return int(f)&fuse.O_CREAT != 0
}

// IsExclusive checks if the flag is set to exclusive.
func (f FuseFlags) IsExclusive() bool {
	return int(f)&fuse.O_EXCL != 0
}

// IsTruncate checks if the flag is set to truncate.
func (f FuseFlags) IsTruncate() bool {
	return int(f)&fuse.O_TRUNC != 0
}

// IsAppend checks if the flag is set to append.
func (f FuseFlags) IsAppend() bool {
	return int(f)&fuse.O_APPEND != 0
}

func (f FuseFlags) IsCreateExclusive() bool {
	return f.IsCreate() && f.IsExclusive()
}

// String returns a string representation of the FuseFlags.
func (f FuseFlags) String() string {
	flags := []string{}
	if f.IsEventOnly() {
		flags = append(flags, "O_EVTONLY")
	}
	if f.IsReadOnly() {
		flags = append(flags, "O_RDONLY")
	}
	if f.IsWriteOnly() {
		flags = append(flags, "O_WRONLY")
	}
	if f.IsReadWrite() {
		flags = append(flags, "O_RDWR")
	}
	if f.IsCreate() {
		flags = append(flags, "O_CREAT")
	}
	if f.IsExclusive() {
		flags = append(flags, "O_EXCL")
	}
	if f.IsTruncate() {
		flags = append(flags, "O_TRUNC")
	}
	if f.IsAppend() {
		flags = append(flags, "O_APPEND")
	}
	if len(flags) == 0 {
		return "NONE"
	}
	return fmt.Sprintf("FuseFlags{%s}", strings.Join(flags, "|"))
}
