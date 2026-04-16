//go:build windows

package preflight

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"syscall"

	"golang.org/x/sys/windows/registry"
)

type winfspDiagnostics struct {
	registryLookupSucceeded bool
	registryLookupError     error
	installDir              string
	installDirError         error
	dllName                 string
	dllPath                 string
	dllExists               bool
	dllExistenceKnown       bool
	dllExistsError          error
	explicitLoadError       error
}

// LoadFuse tests that the FUSE library can be loaded and returns actionable WinFSP diagnostics on failure.
func LoadFuse() error {
	panicErr := runFuseOptParse()
	if panicErr == nil {
		return nil
	}
	diag := collectWinFSPDiagnostics()

	message := fmt.Sprintf(
		"failed to initialize FUSE; WinFSP preflight diagnostics: registry_lookup=%s install_dir=%q dll_path=%q dll_exists=%s explicit_load_error=%s fuse_optparse_panic=%q",
		diag.registryLookupStatus(),
		diag.installDirValue(),
		diag.dllPathValue(),
		diag.dllExistsStatus(),
		diag.explicitLoadStatus(),
		panicErr.Error(),
	)
	if diag.explicitLoadError != nil {
		return fmt.Errorf("%s: %w", message, diag.explicitLoadError)
	}
	return fmt.Errorf("%s", message)
}

func collectWinFSPDiagnostics() winfspDiagnostics {
	diag := winfspDiagnostics{
		dllName: expectedWinFSPDLLName(),
	}

	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `Software\WinFsp`, registry.QUERY_VALUE|registry.WOW64_32KEY)
	if err != nil {
		diag.registryLookupError = fmt.Errorf("open HKLM\\Software\\WinFsp in 32-bit registry view: %w", err)
		return diag
	}
	defer key.Close()
	diag.registryLookupSucceeded = true

	installDir, _, err := key.GetStringValue("InstallDir")
	if err != nil {
		diag.installDirError = fmt.Errorf("read InstallDir from HKLM\\Software\\WinFsp: %w", err)
		return diag
	}
	diag.installDir = installDir
	diag.dllPath = filepath.Join(installDir, "bin", diag.dllName)

	st, err := os.Stat(diag.dllPath)
	switch {
	case err == nil:
		diag.dllExistenceKnown = true
		diag.dllExists = !st.IsDir()
	case os.IsNotExist(err):
		diag.dllExistenceKnown = true
		diag.dllExists = false
	default:
		diag.dllExistsError = fmt.Errorf("stat %q: %w", diag.dllPath, err)
	}

	dll, err := syscall.LoadDLL(diag.dllPath)
	if err != nil {
		diag.explicitLoadError = err
		return diag
	}
	_ = dll.Release()
	return diag
}

func expectedWinFSPDLLName() string {
	switch runtime.GOARCH {
	case "arm64":
		return "winfsp-a64.dll"
	case "amd64":
		return "winfsp-x64.dll"
	case "386":
		return "winfsp-x86.dll"
	default:
		return "winfsp.dll"
	}
}

func (d winfspDiagnostics) registryLookupStatus() string {
	switch {
	case d.registryLookupSucceeded:
		return "ok"
	case d.registryLookupError != nil:
		return d.registryLookupError.Error()
	default:
		return "not attempted"
	}
}

func (d winfspDiagnostics) installDirValue() string {
	switch {
	case d.installDir != "":
		return d.installDir
	case d.installDirError != nil:
		return "unavailable (" + d.installDirError.Error() + ")"
	default:
		return ""
	}
}

func (d winfspDiagnostics) dllPathValue() string {
	if d.dllPath != "" {
		return d.dllPath
	}
	if d.installDir == "" && d.dllName != "" {
		return filepath.Join("<InstallDir>", "bin", d.dllName)
	}
	return ""
}

func (d winfspDiagnostics) dllExistsStatus() string {
	switch {
	case d.dllExistenceKnown:
		if d.dllExists {
			return "true"
		}
		return "false"
	case d.dllExistsError != nil:
		return "unknown (" + d.dllExistsError.Error() + ")"
	default:
		return "unknown"
	}
}

func (d winfspDiagnostics) explicitLoadStatus() string {
	if d.explicitLoadError == nil {
		if d.dllPath == "" {
			return "skipped"
		}
		return "ok"
	}
	return formatWindowsLoaderError(d.explicitLoadError)
}

func formatWindowsLoaderError(err error) string {
	var errno syscall.Errno
	if !errorAsErrno(err, &errno) {
		return err.Error()
	}
	return fmt.Sprintf("%s (errno=%d/0x%x)", err.Error(), uint32(errno), uint32(errno))
}

func errorAsErrno(err error, target *syscall.Errno) bool {
	return errors.As(err, target)
}
