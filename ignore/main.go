package ignore

import (
	"fmt"
	"runtime"
	"strings"

	ignore "github.com/sabhiram/go-gitignore"
)

func New() (*ignore.GitIgnore, error) {
	os := runtime.GOOS
	switch os {
	case "windows":
		data, err := Asset("ignore/data/Windows.gitignore")
		if err != nil {
			return &ignore.GitIgnore{}, err
		}
		return ignore.CompileIgnoreLines(strings.Split(string(data), "\n")...), nil
	case "darwin":
		data, err := Asset("ignore/data/macOS.gitignore")
		if err != nil {
			return &ignore.GitIgnore{}, err
		}
		return ignore.CompileIgnoreLines(strings.Split(string(data), "\n")...), nil
	case "linux":
		data, err := Asset("ignore/data/Linux.gitignore")
		if err != nil {
			return &ignore.GitIgnore{}, err
		}
		return ignore.CompileIgnoreLines(strings.Split(string(data), "\n")...), nil
	default:
		return &ignore.GitIgnore{}, fmt.Errorf("unknown os %s", os)
	}
}
