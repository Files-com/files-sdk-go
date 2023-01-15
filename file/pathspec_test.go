package file

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

type pathSpecEntry struct {
	dir         bool
	path        string
	preexisting bool
}
type pathSpecArgs struct {
	src  string
	dest string
}

type pathSpecTest struct {
	name string
	args pathSpecArgs
	dest []pathSpecEntry
	src  []pathSpecEntry
}

func pathSpec(srcPathSeparator string, destPathSeparator string) []pathSpecTest {
	return []pathSpecTest{
		{
			name: "copy foo to dest",
			args: pathSpecArgs{
				src:  join(srcPathSeparator, "src", "foo"),
				dest: "dest",
			},
			src: []pathSpecEntry{
				{dir: true, path: "src"},
				{dir: true, path: join(srcPathSeparator, "src", "foo")},
				{dir: false, path: join(srcPathSeparator, "src", "foo", "baz.txt")},
			},
			dest: []pathSpecEntry{
				{dir: true, path: "dest", preexisting: true},
				{dir: true, path: join(destPathSeparator, "dest", "foo")},
				{dir: false, path: join(destPathSeparator, "dest", "foo", "baz.txt")},
			},
		},
		{
			name: "copy contents of foo to dest/foo",
			args: pathSpecArgs{
				src:  join(srcPathSeparator, "src", "foo") + srcPathSeparator,
				dest: join(destPathSeparator, "dest", "foo"),
			},
			src: []pathSpecEntry{
				{dir: true, path: "src"},
				{dir: true, path: join(srcPathSeparator, "src", "foo")},
				{dir: false, path: join(srcPathSeparator, "src", "foo", "baz.txt")},
			},
			dest: []pathSpecEntry{
				{dir: true, path: "dest", preexisting: true},
				{dir: true, path: join(destPathSeparator, "dest", "foo"), preexisting: true},
				{dir: false, path: join(destPathSeparator, "dest", "foo", "baz.txt")},
			},
		},
		{
			name: "copy contents of foo to dest/bar",
			args: pathSpecArgs{
				src:  join(srcPathSeparator, "src", "foo") + srcPathSeparator,
				dest: join(destPathSeparator, "dest", "bar"),
			},
			src: []pathSpecEntry{
				{dir: true, path: "src"},
				{dir: true, path: join(srcPathSeparator, "src", "foo")},
				{dir: false, path: join(srcPathSeparator, "src", "foo", "baz.txt")},
			},
			dest: []pathSpecEntry{
				{dir: true, path: "dest", preexisting: true},
				{dir: true, path: join(destPathSeparator, "dest", "bar"), preexisting: true},
				{dir: false, path: join(destPathSeparator, "dest", "bar", "baz.txt")},
			},
		},
		{
			name: "copy baz.txt to dest",
			args: pathSpecArgs{
				src:  join(srcPathSeparator, "src", "foo", "baz.txt"),
				dest: join(destPathSeparator, "dest", "baz.txt"),
			},
			src: []pathSpecEntry{
				{dir: true, path: join(srcPathSeparator, "src")},
				{dir: true, path: join(srcPathSeparator, "src", "foo")},
				{dir: false, path: join(srcPathSeparator, "src", "foo", "baz.txt")},
			},
			dest: []pathSpecEntry{
				{dir: true, path: "dest", preexisting: true},
				{dir: false, path: join(destPathSeparator, "dest", "baz.txt")},
			},
		},
		{
			name: "copy baz.txt to dest without name",
			args: pathSpecArgs{
				src:  join(srcPathSeparator, "src", "foo", "baz.txt"),
				dest: "dest" + destPathSeparator,
			},
			src: []pathSpecEntry{
				{dir: true, path: "src"},
				{dir: true, path: join(srcPathSeparator, "src", "foo")},
				{dir: false, path: join(srcPathSeparator, "src", "foo", "baz.txt")},
			},
			dest: []pathSpecEntry{
				{dir: true, path: "dest", preexisting: true},
				{dir: false, path: join(destPathSeparator, "dest", "baz.txt")},
			},
		},
		{
			name: "copy baz.txt to dest with rename",
			args: pathSpecArgs{
				src:  join(srcPathSeparator, "src", "foo", "baz.txt"),
				dest: join(destPathSeparator, "dest", "taz.txt"),
			},
			src: []pathSpecEntry{
				{dir: true, path: "src"},
				{dir: true, path: join(srcPathSeparator, "src", "foo")},
				{dir: false, path: join(srcPathSeparator, "src", "foo", "baz.txt")},
			},
			dest: []pathSpecEntry{
				{dir: true, path: "dest", preexisting: true},
				{dir: false, path: join(destPathSeparator, "dest/taz.txt")},
			},
		},
		{
			name: "copy baz.txt to current working directory",
			args: pathSpecArgs{
				src:  join(srcPathSeparator, "src", "foo", "baz.txt"),
				dest: "",
			},
			src: []pathSpecEntry{
				{dir: true, path: "src"},
				{dir: true, path: join(srcPathSeparator, "src", "foo")},
				{dir: false, path: join(srcPathSeparator, "src", "foo", "baz.txt")},
			},
			dest: []pathSpecEntry{
				{dir: false, path: "baz.txt"},
			},
		},
	}
}

func TestRsync(t *testing.T) {
	_, err := exec.LookPath("rsync")
	if err != nil {
		t.Log(err)
		return
	}

	t.Run("rsync", func(t *testing.T) {
		mutex := &sync.Mutex{}
		for _, tt := range pathSpec(string(os.PathSeparator), string(os.PathSeparator)) {
			t.Run(tt.name, func(t *testing.T) {
				rsyncSrc, err := ioutil.TempDir("", "rsync-src")
				assert.NoError(t, err)
				rsyncDest, err := ioutil.TempDir("", "rsync-dest")
				assert.NoError(t, err)

				for _, e := range tt.src {
					if e.dir {
						err = os.MkdirAll(filepath.Join(rsyncSrc, e.path), 0750)
					} else {
						_, err = os.Create(filepath.Join(rsyncSrc, e.path))
					}
					require.NoError(t, err)
				}
				for _, e := range tt.dest {
					if !e.preexisting {
						continue
					}
					if e.dir {
						err = os.MkdirAll(filepath.Join(rsyncDest, e.path), 0750)
					} else {
						_, err = os.Create(filepath.Join(rsyncDest, e.path))
					}
					require.NoError(t, err)
				}

				source := join(string(os.PathSeparator), rsyncSrc, tt.args.src)
				destination := join(string(os.PathSeparator), rsyncDest, tt.args.dest)
				if tt.args.dest == "" {
					destination = ""
				}
				originalDir, err := os.Getwd()
				require.NoError(t, err)
				mutex.Lock()
				err = os.Chdir(rsyncDest)
				require.NoError(t, err)
				cmd := exec.Command("rsync", "-av", source, destination)
				var cleanedArgs []string

				for _, arg := range cmd.Args {
					arg = strings.Replace(arg, rsyncSrc, "", 1)
					arg = strings.Replace(arg, rsyncDest, "", 1)
					cleanedArgs = append(cleanedArgs, arg)
				}
				t.Log(cleanedArgs)
				stdout := bytes.NewBufferString("")
				stderr := bytes.NewBufferString("")
				cmd.Stderr = stderr
				cmd.Stdout = stdout
				err = cmd.Run()
				assert.NoError(t, err)
				if stderr.String() != "" {
					t.Log(stderr.String())
				}
				err = os.Chdir(originalDir)
				require.NoError(t, err)
				mutex.Unlock()
				for _, e := range tt.dest {
					fileInfo, err := os.Stat(filepath.Join(rsyncDest, e.path))
					require.NoError(t, err, e.path)
					assert.Equal(t, e.dir, fileInfo.IsDir(), e.path)
				}

				assert.NoError(t, os.RemoveAll(rsyncSrc))
				assert.NoError(t, os.RemoveAll(rsyncDest))
			})
		}
	})
}

func join(pathSeparator string, parts ...string) string {
	if parts[len(parts)-1] == pathSeparator {
		return strings.Join(parts, pathSeparator) + pathSeparator
	} else {
		return strings.Join(parts, pathSeparator)
	}
}
