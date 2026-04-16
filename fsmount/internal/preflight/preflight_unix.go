//go:build darwin || linux

package preflight

// LoadFuse tests that the FUSE library can be loaded and returns any panic as an error.
func LoadFuse() error {
	return runFuseOptParse()
}
