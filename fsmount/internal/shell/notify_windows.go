//go:build windows

package shell

import (
	"fmt"
	"sync"
	"syscall"
	"unsafe"
)

const (
	// SHCNE_UPDATEDIR notifies the shell that a directory has changed.
	SHCNE_UPDATEDIR = 0x00001000
	// SHCNF_PATHW indicates a Unicode path parameter.
	SHCNF_PATHW = 0x00000005
)

var (
	defaultNotifier     *Notifier
	defaultNotifierOnce sync.Once
	defaultNotifierErr  error
)

// Notifier provides methods to send shell change notifications to Windows Explorer.
type Notifier struct {
	dll  *syscall.DLL
	proc *syscall.Proc
}

// NewNotifier creates a new Notifier instance, loading shell32.dll and resolving the SHChangeNotify procedure.
// Returns an error if the DLL or procedure cannot be loaded.
func NewNotifier() (*Notifier, error) {
	dll, err := syscall.LoadDLL("shell32.dll")
	if err != nil {
		return nil, fmt.Errorf("failed to load shell32.dll: %w", err)
	}

	proc, err := dll.FindProc("SHChangeNotify")
	if err != nil {
		return nil, fmt.Errorf("failed to find SHChangeNotify procedure: %w", err)
	}

	return &Notifier{
		dll:  dll,
		proc: proc,
	}, nil
}

// notifyUpdatedDir sends a change notification to Windows Explorer for the specified directory path,
// causing Explorer to refresh any views showing that directory.
// This is a private method - external callers should use the NotifyUpdatedDir function instead.
func (n *Notifier) notifyUpdatedDir(path string) error {
	if path == "" {
		return nil
	}

	p, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return fmt.Errorf("failed to convert path to UTF16: %w", err)
	}

	n.proc.Call(
		uintptr(SHCNE_UPDATEDIR),
		uintptr(SHCNF_PATHW),
		uintptr(unsafe.Pointer(p)),
		uintptr(0),
	)
	return nil
}

// getDefaultNotifier returns a package-level singleton Notifier, initialized once via sync.Once.
// This avoids repeated DLL loading and procedure lookups for every notification.
func getDefaultNotifier() (*Notifier, error) {
	defaultNotifierOnce.Do(func() {
		defaultNotifier, defaultNotifierErr = NewNotifier()
	})
	return defaultNotifier, defaultNotifierErr
}

// NotifyUpdatedDir sends a directory update notification to Windows Explorer.
// It uses a singleton Notifier that is initialized once and reused for all calls.
func NotifyUpdatedDir(path string) error {
	notifier, err := getDefaultNotifier()
	if err != nil {
		return err
	}
	return notifier.notifyUpdatedDir(path)
}
