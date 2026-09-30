//go:build !windows

package privatefile

import (
	"io/fs"
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func perm(t *testing.T, path string) fs.FileMode {
	t.Helper()
	info, err := os.Lstat(path)
	require.NoError(t, err)
	return info.Mode().Perm()
}

// Directory modes here are asserted with a permissive umask, so the bits the
// package sets are observable rather than removed by the environment.
func withUmask(t *testing.T, mask int) {
	t.Helper()
	previous := syscall.Umask(mask)
	t.Cleanup(func() { syscall.Umask(previous) })
}

func TestMkdirTempCreatesOwnerOnlyDirectory(t *testing.T) {
	withUmask(t, 0)
	dir, err := MkdirTemp(t.TempDir(), "Files.com-v6-*")
	require.NoError(t, err)
	assert.Equal(t, fs.FileMode(0o700), perm(t, dir))
	assert.Contains(t, filepath.Base(dir), "Files.com-v6-")
}

func TestEnsureDirTightensOwnedDirectoryAndRefusesLinks(t *testing.T) {
	withUmask(t, 0)
	root := t.TempDir()
	shared := filepath.Join(root, "data")
	require.NoError(t, os.Mkdir(shared, 0o755))
	require.NoError(t, EnsureDir(shared))
	assert.Equal(t, fs.FileMode(0o700), perm(t, shared), "an existing directory we own loses group and other access")

	created := filepath.Join(root, "state")
	require.NoError(t, EnsureDir(created))
	assert.Equal(t, fs.FileMode(0o700), perm(t, created))

	target := filepath.Join(root, "elsewhere")
	require.NoError(t, os.Mkdir(target, 0o755))
	link := filepath.Join(root, "partial")
	require.NoError(t, os.Symlink(target, link))
	require.Error(t, EnsureDir(link), "a symbolic link where a private directory should be is refused")
	assert.Equal(t, fs.FileMode(0o755), perm(t, target), "the link target is left alone")

	file := filepath.Join(root, "file")
	require.NoError(t, os.WriteFile(file, nil, 0o644))
	require.Error(t, EnsureDir(file))
}

func TestEnsureEntryPrivateRepairsOwnedEntriesOnly(t *testing.T) {
	withUmask(t, 0)
	root := t.TempDir()
	file := filepath.Join(root, "payload")
	require.NoError(t, os.WriteFile(file, []byte("x"), 0o644))
	info, err := os.Lstat(file)
	require.NoError(t, err)
	require.NoError(t, EnsureEntryPrivate(file, info))
	assert.Equal(t, fs.FileMode(0o600), perm(t, file))

	dir := filepath.Join(root, "shard")
	require.NoError(t, os.Mkdir(dir, 0o755))
	info, err = os.Lstat(dir)
	require.NoError(t, err)
	require.NoError(t, EnsureEntryPrivate(dir, info))
	assert.Equal(t, fs.FileMode(0o700), perm(t, dir))

	outside := filepath.Join(t.TempDir(), "outside")
	require.NoError(t, os.WriteFile(outside, []byte("y"), 0o644))
	link := filepath.Join(root, "link")
	require.NoError(t, os.Symlink(outside, link))
	info, err = os.Lstat(link)
	require.NoError(t, err)
	require.Error(t, EnsureEntryPrivate(link, info), "a symbolic link is refused rather than followed")
	assert.Equal(t, fs.FileMode(0o644), perm(t, outside), "the link target is not changed")

	hardlink := filepath.Join(root, "second-name")
	require.NoError(t, os.Link(outside, hardlink))
	info, err = os.Lstat(hardlink)
	require.NoError(t, err)
	require.Error(t, EnsureEntryPrivate(hardlink, info), "a file with another name elsewhere is refused")
	assert.Equal(t, fs.FileMode(0o644), perm(t, outside))
}

func TestCheckContainerAcceptsTrustedLocationsAndRejectsReplaceableOnes(t *testing.T) {
	require.NoError(t, CheckContainer(os.TempDir()), "the system temporary directory is sticky and trusted")

	root := t.TempDir()
	// With no umask, directories are created with the modes named below.
	// t.TempDir itself creates its directory with 0777 less the umask, which
	// is group-writable under the common 0002 umask, so it is tightened.
	withUmask(t, 0)
	require.NoError(t, os.Chmod(root, 0o700))
	require.NoError(t, CheckContainer(root), "a directory only this user can modify is a trusted location")
	open := filepath.Join(root, "open")
	require.NoError(t, os.Mkdir(open, 0o777))
	err := CheckContainer(open)
	require.Error(t, err, "a world-writable directory without the sticky bit is refused")
	assert.ErrorContains(t, err, "writable by every user")

	below := filepath.Join(open, "below")
	require.NoError(t, os.Mkdir(below, 0o700))
	require.Error(t, CheckContainer(below), "a private directory below a replaceable one is still replaceable")

	sticky := filepath.Join(root, "sticky")
	require.NoError(t, os.Mkdir(sticky, 0o777|fs.ModeSticky))
	require.NoError(t, os.Chmod(sticky, 0o777|fs.ModeSticky))
	require.NoError(t, CheckContainer(sticky), "a sticky world-writable directory only lets owners remove their entries")

	groupShared := filepath.Join(root, "group")
	require.NoError(t, os.Mkdir(groupShared, 0o775))
	err = CheckContainer(groupShared)
	require.Error(t, err, "a group-writable directory is refused even when this user owns it")
	assert.ErrorContains(t, err, "writable by a group")

	private := filepath.Join(root, "private")
	require.NoError(t, os.Mkdir(private, 0o700))
	require.NoError(t, os.Mkdir(filepath.Join(private, "below"), 0o700))
	link := filepath.Join(root, "link")
	require.NoError(t, os.Symlink(private, link))
	require.NoError(t, CheckContainer(link), "a link the user owns to a trusted location is accepted")

	// A link that resolves to a safe tree is still replaceable when the
	// directory holding the link is not: everyone could swap the link.
	plantedLink := filepath.Join(open, "location")
	require.NoError(t, os.Symlink(private, plantedLink))
	err = CheckContainer(plantedLink)
	require.Error(t, err, "a link inside a world-writable directory is refused even when it points at a safe tree")
	assert.ErrorContains(t, err, "writable by every user")

	// The same holds for a link that only reaches the unsafe directory
	// through another link, and for a relative link target.
	hop := filepath.Join(root, "hop")
	require.NoError(t, os.Symlink(filepath.Join("open", "location"), hop))
	require.Error(t, CheckContainer(hop), "a chain of links passing through a world-writable directory is refused")
	require.Error(t, CheckContainer(filepath.Join(plantedLink, "below")), "a path below such a link is refused too")
}

func TestOpenFileCreatesPrivateFilesAndTightensExistingOnes(t *testing.T) {
	withUmask(t, 0)
	dir := t.TempDir()
	root, err := os.OpenRoot(dir)
	require.NoError(t, err)
	defer root.Close()

	f, err := OpenFile(root, "fresh", os.O_RDWR|os.O_CREATE|os.O_EXCL, 0o700)
	require.NoError(t, err)
	defer f.Close()
	assert.Equal(t, fs.FileMode(0o700), perm(t, filepath.Join(dir, "fresh")))

	require.NoError(t, os.WriteFile(filepath.Join(dir, "legacy"), []byte("old"), 0o644))
	require.NoError(t, os.Chmod(filepath.Join(dir, "legacy"), 0o644))
	g, err := OpenFile(root, "legacy", os.O_WRONLY|os.O_CREATE, 0o700)
	require.NoError(t, err)
	defer g.Close()
	assert.Equal(t, fs.FileMode(0o600), perm(t, filepath.Join(dir, "legacy")), "the existing file is private before any write")

	_, err = OpenFile(root, "bad", os.O_RDWR|os.O_CREATE, 0o644)
	require.Error(t, err, "a private mode must not grant group or other access")
}

func TestPublicationForIntersectsFreshExistingAndIntended(t *testing.T) {
	withUmask(t, 0o027)
	dir := t.TempDir()
	root, err := os.OpenRoot(dir)
	require.NoError(t, err)
	defer root.Close()

	publish := func(name string, intended fs.FileMode) fs.FileMode {
		t.Helper()
		p, err := PublicationFor(root, name, ".probe", intended)
		require.NoError(t, err)
		_, err = root.Lstat(".probe")
		require.ErrorIs(t, err, fs.ErrNotExist, "the probe is removed")
		return p.platform.mode
	}
	t.Run("apply only to the published file", func(t *testing.T) {
		require.NoError(t, os.WriteFile(filepath.Join(dir, "stage"), []byte("x"), 0o600))
		require.NoError(t, os.Chmod(filepath.Join(dir, "stage"), 0o600))
		info, err := os.Lstat(filepath.Join(dir, "stage"))
		require.NoError(t, err)
		p, err := PublicationFor(root, "out", ".probe", fs.ModePerm)
		require.NoError(t, err)
		require.NoError(t, root.Rename("stage", "out"))
		require.NoError(t, p.ApplyToPublished(root, "out", info))
		assert.Equal(t, fs.FileMode(0o640), perm(t, filepath.Join(dir, "out")))

		require.NoError(t, os.WriteFile(filepath.Join(dir, "other"), []byte("y"), 0o600))
		require.NoError(t, os.Chmod(filepath.Join(dir, "other"), 0o600))
		err = p.ApplyToPublished(root, "other", info)
		require.ErrorIs(t, err, errPublishedReplaced, "another file at the name is left alone")
		assert.Equal(t, fs.FileMode(0o600), perm(t, filepath.Join(dir, "other")))
	})
	assert.Equal(t, fs.FileMode(0o640), publish("missing", fs.ModePerm), "a fresh file gets 0666 before the umask")
	require.NoError(t, os.WriteFile(filepath.Join(dir, "restricted"), []byte("x"), 0o600))
	require.NoError(t, os.Chmod(filepath.Join(dir, "restricted"), 0o600))
	assert.Equal(t, fs.FileMode(0o600), publish("restricted", fs.ModePerm), "an existing restrictive file keeps its restriction")
	require.NoError(t, os.WriteFile(filepath.Join(dir, "wide"), []byte("x"), 0o666))
	require.NoError(t, os.Chmod(filepath.Join(dir, "wide"), 0o666))
	assert.Equal(t, fs.FileMode(0o640), publish("wide", fs.ModePerm), "an existing wide file does not widen the fresh mode")
	assert.Equal(t, fs.FileMode(0o600), publish("missing", 0o600), "a stage the caller prepared as 0600 stays 0600")
	assert.Equal(t, fs.FileMode(0o600), publish("missing", 0o700), "a stage the caller prepared as 0700 gains no group or other access")
	assert.Equal(t, fs.FileMode(0o640), publish("missing", 0o640))
}
