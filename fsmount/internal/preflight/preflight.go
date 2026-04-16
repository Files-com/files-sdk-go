package preflight

import (
	"fmt"

	"github.com/winfsp/cgofuse/fuse"
)

func runFuseOptParse() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	// This forces cgofuse to load the platform FUSE implementation and
	// will panic if the fuse library cannot be loaded.
	// The panic is recovered via the defer and returned as an error.
	_, _ = fuse.OptParse([]string{}, "")
	return nil
}
