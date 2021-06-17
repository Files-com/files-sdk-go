// Code generated for package ignore by go-bindata DO NOT EDIT. (@generated)
// sources:
// ignore/data/Linux.gitignore
// ignore/data/Windows.gitignore
// ignore/data/macOS.gitignore
package ignore

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _ignoreDataLinuxGitignore = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x44\x8f\xbd\x4e\x04\x21\x14\x85\x7b\x9e\xe2\x24\x76\x93\x38\x4f\xa1\x95\x96\xf6\xe6\x0e\x1c\x96\x1b\x19\x20\x17\xc6\x75\x1a\x9f\xdd\xe0\x6e\xb4\x84\x7b\xfe\xbe\xe5\xdb\xb9\x07\x0c\xee\xad\x9a\xd8\x89\xa8\x99\x1d\xd7\xa4\x3e\xc1\x4b\xc1\x46\x78\xa3\x0c\x06\x68\x84\xa0\x59\xf5\xec\x1d\x7d\x68\xce\x48\xd2\x21\x48\x52\x42\x26\x6a\x63\x41\x9d\xa2\xc0\xcc\xe9\x98\x61\x6e\x8d\x47\xe7\x7b\xd2\x10\x58\x96\x59\xf6\xf2\xf4\x8c\xa0\x46\x3f\xaa\x9d\x68\xc6\x48\x63\xf1\xec\x6e\xfd\xfb\x9e\xba\x57\x2d\xc7\x17\x86\x49\x4f\x88\x35\x07\xda\x7d\xd7\xae\x97\x34\x20\xad\x51\x0c\xb5\x40\xca\x89\x26\x36\x74\x68\x2d\xa8\x86\xa0\xfd\xc3\xad\x6f\xd3\xf9\xf8\x5b\xb9\x96\xd8\xef\x68\x62\xff\x44\xd7\xc4\xe9\xbe\x2d\x9f\x67\x68\x87\x71\xaf\x9f\x0c\xd8\x8e\x31\x9f\x37\xd0\x8d\x5a\x2e\x10\x3f\xd9\x19\xdc\xcc\x5b\xdc\x4f\x00\x00\x00\xff\xff\x9d\xb1\x6e\x51\x3c\x01\x00\x00")

func ignoreDataLinuxGitignoreBytes() ([]byte, error) {
	return bindataRead(
		_ignoreDataLinuxGitignore,
		"ignore/data/Linux.gitignore",
	)
}

func ignoreDataLinuxGitignore() (*asset, error) {
	bytes, err := ignoreDataLinuxGitignoreBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "ignore/data/Linux.gitignore", size: 316, mode: os.FileMode(420), modTime: time.Unix(1623959812, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _ignoreDataWindowsGitignore = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x4c\x8f\x4f\x4b\x03\x31\x10\xc5\xef\xf9\x14\x03\xf5\xb4\x87\x78\xf7\xd8\x3f\x42\x41\x3c\x14\x41\x44\x44\xb2\xc9\xd4\x1d\x36\x99\x84\xcc\x44\xed\xb7\x97\x6c\xad\xf6\xf2\xf2\x78\xf9\xbd\x81\xb7\x82\x67\xe2\x90\xbf\x04\x74\x6a\x69\x64\x47\x11\xbc\xf3\x13\xc2\x91\x22\x8a\x79\xea\xa9\xd8\x30\xfe\xbb\x3b\x64\x5f\x4f\x45\xdd\x18\xd1\xe0\xa4\x7f\xc4\xc5\xbf\x7f\x92\xa8\xeb\x89\x59\xc1\xb6\xa5\xb2\xdc\x32\x83\x15\x75\x7e\x0e\x2d\x95\xfe\x71\x9f\x63\xc0\x0a\x3e\xf3\x91\x3e\xce\xc4\xeb\x36\xbc\xa1\xcc\x9a\x8b\x25\xa6\x0e\x1d\xd0\x9f\x7c\x44\x58\x13\x43\x13\x0c\x90\x79\x41\x41\x26\x57\x51\xcc\xcd\x61\xb7\x79\xd9\x3c\xec\xec\x7a\xff\x78\xdb\x0b\x97\x35\x7b\x16\x75\x31\x62\xfd\xdd\x31\x58\xef\x46\x33\xd8\x24\x74\xd6\xef\xe5\x49\x8b\x96\xeb\xa6\x4c\xb9\xaa\x6f\xda\x3b\x91\x67\xf3\x13\x00\x00\xff\xff\xdc\x2a\x61\x5c\x22\x01\x00\x00")

func ignoreDataWindowsGitignoreBytes() ([]byte, error) {
	return bindataRead(
		_ignoreDataWindowsGitignore,
		"ignore/data/Windows.gitignore",
	)
}

func ignoreDataWindowsGitignore() (*asset, error) {
	bytes, err := ignoreDataWindowsGitignoreBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "ignore/data/Windows.gitignore", size: 290, mode: os.FileMode(420), modTime: time.Unix(1623959812, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _ignoreDataMacosGitignore = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x44\xd0\xc1\x8a\x1b\x31\x0c\x06\xe0\xbb\x60\xdf\x41\xb0\xb7\x42\xcd\xf6\x11\xb6\x0c\x29\x81\xd2\x96\x4e\xc8\xa9\x10\x9c\x99\x3f\xb5\x89\x6d\x19\x49\x33\x21\x6f\x5f\x26\xa4\xed\xcd\xf8\x17\x3f\x9f\xf4\xca\x5f\xd0\xa0\xb1\x50\x18\xc6\xd3\xe8\xa2\xa0\xf0\xde\x7b\xc1\x20\xcb\xb9\x80\xc2\xd7\xf1\xfb\x0a\xd5\x3c\x83\xe8\x95\xf7\x93\x34\xae\x8b\x39\xa3\xcd\x7c\xcb\x9e\xd8\x6f\xc2\xbf\x94\xb6\xe4\xe5\x65\x9b\x39\xa4\xa5\x9e\x5b\xcc\xc5\x28\x9c\x3e\x6c\x3f\xbb\x5c\x60\xec\x29\x3a\xd7\xfc\x3b\x39\xc7\xde\x11\x95\x73\x63\x4f\x60\x15\x71\x96\x0b\x47\x5e\xa5\x2c\x15\x14\x06\x99\x96\x8a\xe6\x3f\xb1\x66\xcb\xd2\xec\xe3\xf1\xd3\xdb\x1b\x85\x8b\x61\x45\x73\x9b\x29\x8c\x5d\xbc\x6c\x65\xcf\xe8\x80\xda\x45\xa3\xde\xf7\x8e\x6a\x14\x0e\x1a\x2d\xc1\x28\x1c\x1f\x9d\x1b\x2f\xe4\xa9\x19\x85\x49\x6a\x88\xdb\x8a\xc1\x73\x45\x8d\x53\xca\x0d\x61\x96\x26\xde\x15\x86\xe6\x9b\x79\xc8\x8a\xc9\x45\x33\x8c\xbb\x38\x9a\xe7\x58\xca\x9d\x27\x45\x74\xcc\x2c\x8d\x15\x55\x1c\xfc\xbe\xfb\xc1\x96\xe2\xff\xcb\x7d\xfe\xfb\x80\x5d\x5d\x3a\x7d\x83\xdf\x44\xaf\xfc\x20\xf1\x4e\xca\x0c\xa5\x7f\x5e\x7e\x82\x63\x9f\xb3\x5d\xe9\x4f\x00\x00\x00\xff\xff\x8e\xe0\xf8\x3a\x92\x01\x00\x00")

func ignoreDataMacosGitignoreBytes() ([]byte, error) {
	return bindataRead(
		_ignoreDataMacosGitignore,
		"ignore/data/macOS.gitignore",
	)
}

func ignoreDataMacosGitignore() (*asset, error) {
	bytes, err := ignoreDataMacosGitignoreBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "ignore/data/macOS.gitignore", size: 402, mode: os.FileMode(420), modTime: time.Unix(1623959812, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"ignore/data/Linux.gitignore":   ignoreDataLinuxGitignore,
	"ignore/data/Windows.gitignore": ignoreDataWindowsGitignore,
	"ignore/data/macOS.gitignore":   ignoreDataMacosGitignore,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"ignore": &bintree{nil, map[string]*bintree{
		"data": &bintree{nil, map[string]*bintree{
			"Linux.gitignore":   &bintree{ignoreDataLinuxGitignore, map[string]*bintree{}},
			"Windows.gitignore": &bintree{ignoreDataWindowsGitignore, map[string]*bintree{}},
			"macOS.gitignore":   &bintree{ignoreDataMacosGitignore, map[string]*bintree{}},
		}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
