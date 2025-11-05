// Package flags provides utilities for working with FUSE file open flags.
package flags

import (
	"fmt"
	"os"
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

// IsCreateExclusive checks if the flag is set to create and exclusive.
func (f FuseFlags) IsCreateExclusive() bool {
	return f.IsCreate() && f.IsExclusive()
}

// Without returns a new FuseFlags instance with the specified flags removed.
func (f FuseFlags) Without(remove int) FuseFlags {
	return FuseFlags(int(f) &^ remove)
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

// AsOsFlags translates POSIX/FUSE-style open(2) flags carried in FuseFlags
// into the platform-appropriate flags defined by Go’s os package.
//
// Why this is cross-platform:
//   - Access mode: FUSE follows POSIX semantics where the lower two bits encode
//     the access mode (O_RDONLY=0, O_WRONLY=1, O_RDWR=2). Masking with 0x3 and
//     mapping to os.O_RDONLY / os.O_WRONLY / os.O_RDWR allows the os package to turn
//     these into the correct native flags for Linux, macOS, and Windows.
//   - Behavior/creation bits: FUSE flags like O_CREAT, O_EXCL, O_TRUNC, and
//     O_APPEND have direct counterparts in the os package (os.O_CREATE,
//     os.O_EXCL, os.O_TRUNC, os.O_APPEND). OR-ing these through lets the os
//     package handle the per-platform syscall details (including WinFsp on
//     Windows for cgofuse).
//
// How it works:
//  1. Mask the lower two bits (0x3) to select exactly one of R/O, W/O, or R/W.
//  2. OR in any creation/behavior bits that are present.
//  3. Return the composite mask, suitable for passing to os.OpenFile.
//
// Notes:
//   - O_EVTONLY (macOS) is intentionally ignored here because it’s about event
//     notifications rather than the access/creation mode used by os.OpenFile.
//   - Add additional mappings here if more FUSE flags need to flow through.
func (f FuseFlags) AsOsFlags() int {
	osf := 0

	// Access mode (mask is 0|1|2 for R/W/RW)
	switch f & 0x3 { // O_RDONLY=0, O_WRONLY=1, O_RDWR=2 in POSIX
	case 0:
		osf |= os.O_RDONLY
	case 1:
		osf |= os.O_WRONLY
	case 2:
		osf |= os.O_RDWR
	}

	// Creation / behavior bits
	if int(f)&fuse.O_CREAT != 0 {
		osf |= os.O_CREATE
	}
	if int(f)&fuse.O_EXCL != 0 {
		osf |= os.O_EXCL
	}
	if int(f)&fuse.O_TRUNC != 0 {
		osf |= os.O_TRUNC
	}
	if int(f)&fuse.O_APPEND != 0 {
		osf |= os.O_APPEND
	}

	return osf
}
