package file

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

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

func pathSpec() []pathSpecTest {
	return []pathSpecTest{
		{
			name: "copy foo to dest",
			args: pathSpecArgs{
				src:  "src/foo",
				dest: "dest",
			},
			src: []pathSpecEntry{
				{dir: true, path: "src"},
				{dir: true, path: "src/foo"},
				{dir: false, path: "src/foo/baz.txt"},
			},
			dest: []pathSpecEntry{
				{dir: true, path: "dest", preexisting: true},
				{dir: true, path: "dest/foo"},
				{dir: false, path: "dest/foo/baz.txt"},
			},
		},
		{
			name: "copy contents of foo to dest/foo",
			args: pathSpecArgs{
				src:  "src/foo/",
				dest: "dest/foo",
			},
			src: []pathSpecEntry{
				{dir: true, path: "src"},
				{dir: true, path: "src/foo"},
				{dir: false, path: "src/foo/baz.txt"},
			},
			dest: []pathSpecEntry{
				{dir: true, path: "dest", preexisting: true},
				{dir: true, path: "dest/foo", preexisting: true},
				{dir: false, path: "dest/foo/baz.txt"},
			},
		},
		{
			name: "copy contents of foo to dest/bar",
			args: pathSpecArgs{
				src:  "src/foo/",
				dest: "dest/bar",
			},
			src: []pathSpecEntry{
				{dir: true, path: "src"},
				{dir: true, path: "src/foo"},
				{dir: false, path: "src/foo/baz.txt"},
			},
			dest: []pathSpecEntry{
				{dir: true, path: "dest", preexisting: true},
				{dir: true, path: "dest/bar", preexisting: true},
				{dir: false, path: "dest/bar/baz.txt"},
			},
		},
		{
			name: "copy baz.txt to dest",
			args: pathSpecArgs{
				src:  "src/foo/baz.txt",
				dest: "dest/baz.txt",
			},
			src: []pathSpecEntry{
				{dir: true, path: "src"},
				{dir: true, path: "src/foo"},
				{dir: false, path: "src/foo/baz.txt"},
			},
			dest: []pathSpecEntry{
				{dir: true, path: "dest", preexisting: true},
				{dir: false, path: "dest/baz.txt"},
			},
		},
		{
			name: "copy baz.txt to dest without name",
			args: pathSpecArgs{
				src:  "src/foo/baz.txt",
				dest: "dest/",
			},
			src: []pathSpecEntry{
				{dir: true, path: "src"},
				{dir: true, path: "src/foo"},
				{dir: false, path: "src/foo/baz.txt"},
			},
			dest: []pathSpecEntry{
				{dir: true, path: "dest", preexisting: true},
				{dir: false, path: "dest/baz.txt"},
			},
		},
		{
			name: "copy baz.txt to dest with rename",
			args: pathSpecArgs{
				src:  "src/foo/baz.txt",
				dest: "dest/taz.txt",
			},
			src: []pathSpecEntry{
				{dir: true, path: "src"},
				{dir: true, path: "src/foo"},
				{dir: false, path: "src/foo/baz.txt"},
			},
			dest: []pathSpecEntry{
				{dir: true, path: "dest", preexisting: true},
				{dir: false, path: "dest/taz.txt"},
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
		for _, tt := range pathSpec() {
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
					assert.NoError(t, err)
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
					assert.NoError(t, err)
				}
				cmd := exec.Command("rsync", "-av", strings.Join([]string{rsyncSrc, tt.args.src}, string(os.PathSeparator)), strings.Join([]string{rsyncDest, tt.args.dest}, string(os.PathSeparator)))
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
				for _, e := range tt.dest {
					fileInfo, err := os.Stat(filepath.Join(rsyncDest, e.path))
					assert.NoError(t, err, e.path)
					assert.Equal(t, e.dir, fileInfo.IsDir(), e.path)
				}

				assert.NoError(t, os.RemoveAll(rsyncSrc))
				assert.NoError(t, os.RemoveAll(rsyncDest))
			})
		}
	})
}
