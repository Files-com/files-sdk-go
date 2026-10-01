//go:build windows

package disk

import (
	"errors"
	"fmt"
	"os"

	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache"
	"golang.org/x/sys/windows"
)

// markSparseFile enables sparse allocation for an open Windows file.
func markSparseFile(file *os.File) error {
	var bytesReturned uint32
	err := windows.DeviceIoControl(
		windows.Handle(file.Fd()),
		windows.FSCTL_SET_SPARSE,
		nil,
		0,
		nil,
		0,
		&bytesReturned,
		nil,
	)
	if err != nil {
		if errors.Is(err, windows.ERROR_INVALID_FUNCTION) || errors.Is(err, windows.ERROR_NOT_SUPPORTED) {
			return fmt.Errorf("%w: %v", cache.ErrSparseFilesUnsupported, err)
		}
		return fmt.Errorf("marking cache file %s as sparse failed: %w", file.Name(), err)
	}
	return nil
}
