package ignore

import (
	"fmt"
	"runtime"
	"strings"

	_ "embed"

	ignore "github.com/sabhiram/go-gitignore"
)

//go:embed data/Windows.gitignore
var WindowsGitignore []byte

//go:embed data/macOS.gitignore
var macOSGitignore []byte

//go:embed data/Linux.gitignore
var LinuxGitignore []byte

func New() (*ignore.GitIgnore, error) {
	os := runtime.GOOS
	switch os {
	case "windows":
		return ignore.CompileIgnoreLines(strings.Split(string(WindowsGitignore), "\n")...), nil
	case "darwin":
		return ignore.CompileIgnoreLines(strings.Split(string(macOSGitignore), "\n")...), nil
	case "linux":
		return ignore.CompileIgnoreLines(strings.Split(string(LinuxGitignore), "\n")...), nil
	default:
		return &ignore.GitIgnore{}, fmt.Errorf("unknown os %s", os)
	}
}
