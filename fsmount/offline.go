//go:build linux || windows

package fsmount

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	path_lib "path"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"sync"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache/disk"
	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/shell"
)

// OfflineState says how far a file or folder on a mount is kept for offline use.
type OfflineState int

const (
	// OfflineNone is not kept: it opens from the storage.
	OfflineNone OfflineState = iota
	// OfflineDownloading is kept, but some of it isn't downloaded yet.
	OfflineDownloading
	// OfflineKept is kept and downloaded.
	OfflineKept
	// OfflinePartlyKept is a folder that isn't kept, with kept items inside.
	OfflinePartlyKept
)

func (s OfflineState) String() string {
	switch s {
	case OfflineDownloading:
		return "downloading"
	case OfflineKept:
		return "kept"
	case OfflinePartlyKept:
		return "partly kept"
	default:
		return "none"
	}
}

const (
	offlineDownloadWorkers = 2
	// At most this many of the keeper's listings run at once. They use the
	// mount's operation slots, which the mount's own callbacks need too.
	offlineListings = 2
	// A kept file that failed to download is tried again after this long.
	offlineRetryDelay = time.Minute
	// While the storage can't be reached, it is asked again in the background
	// at most this often, when the mount is used.
	offlineUnreachableWindow = 15 * time.Second
	offlineRootsFile         = "roots.json"
	offlineListingsDir       = "listings"
)

// With something kept, a call to the storage for a folder that has a saved
// listing, or made while the storage was just unreachable, gives up when the
// storage hasn't answered after this long. Storage clients retry a refused
// connection for minutes. A variable so tests can shorten it.
var offlineRemoteTimeout = 5 * time.Second

// savedEntry is what a saved listing records about a file or folder.
type savedEntry struct {
	Name        string    `json:"name"`
	Dir         bool      `json:"dir,omitempty"`
	Size        int64     `json:"size,omitempty"`
	ModTime     time.Time `json:"mtime,omitzero"`
	Created     time.Time `json:"ctime,omitzero"`
	Permissions string    `json:"permissions,omitempty"`
}

func savedEntryFor(name string, item files_sdk.File) savedEntry {
	return savedEntry{
		Name:        name,
		Dir:         item.IsDir(),
		Size:        item.Size,
		ModTime:     item.ModTime(),
		Created:     item.CreationTime(),
		Permissions: item.Permissions,
	}
}

func (e savedEntry) file(path string) files_sdk.File {
	item := files_sdk.File{
		Path:        path,
		DisplayName: e.Name,
		Type:        "file",
		Size:        e.Size,
		Permissions: e.Permissions,
	}
	if e.Dir {
		item.Type = "directory"
	}
	if !e.ModTime.IsZero() {
		mtime := e.ModTime
		item.Mtime = &mtime
	}
	if !e.Created.IsZero() {
		created := e.Created
		item.CreatedAt = &created
	}
	return item
}

// savedListing is a folder's listing saved so the folder shows when the
// storage can't be reached, including after a restart. A listing doesn't
// change once it is saved, because savedChildren reads it without k.mu: a
// change saves a new one.
type savedListing struct {
	Path    string       `json:"path"`
	Self    *savedEntry  `json:"self,omitempty"`
	Entries []savedEntry `json:"entries"`
}

// offlineKeeper keeps files and folders of a mount for offline use. A kept
// ("root") path and everything under it is downloaded into the disk cache and
// left there by the cache's limits. Kept folders and the folders above them
// keep their latest listing, which the mount shows when the storage can't be
// reached. Kept folders are brought up to date when they are listed: new and
// changed files are downloaded, removed ones are no longer kept.
type offlineKeeper struct {
	fs  *RemoteFs
	dir string

	mu               sync.Mutex
	roots            map[string]struct{}
	listings         map[string]*savedListing
	pending          map[string]struct{}
	walks            map[string]int
	unreachableUntil time.Time
	// wasUnreachable stays set from a failed call until a call succeeds.
	wasUnreachable bool
	checking       bool

	startOnce sync.Once
	// Holds a token for each of the keeper's listings that is running.
	listingSlots chan struct{}
	// Closed when the keeper closes.
	done chan struct{}

	queueMu  sync.Mutex
	queue    []string
	queued   map[string]struct{}
	queueCnd *sync.Cond
	closed   bool
	workers  sync.WaitGroup
}

func newOfflineKeeper(fs *RemoteFs) (*offlineKeeper, error) {
	dir, err := fs.cacheStore.OfflineDir()
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Join(dir, offlineListingsDir), 0o700); err != nil {
		return nil, err
	}
	k := &offlineKeeper{
		fs:           fs,
		dir:          dir,
		roots:        make(map[string]struct{}),
		listings:     make(map[string]*savedListing),
		pending:      make(map[string]struct{}),
		walks:        make(map[string]int),
		queued:       make(map[string]struct{}),
		listingSlots: make(chan struct{}, offlineListings),
		done:         make(chan struct{}),
	}
	k.queueCnd = sync.NewCond(&k.queueMu)
	data, err := os.ReadFile(filepath.Join(dir, offlineRootsFile))
	if err == nil {
		var roots []string
		if err := json.Unmarshal(data, &roots); err != nil {
			return nil, fmt.Errorf("unreadable kept paths: %w", err)
		}
		for _, root := range roots {
			k.roots[root] = struct{}{}
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	return k, nil
}

// start downloads what is kept but not downloaded yet, and brings kept folders
// up to date.
func (k *offlineKeeper) start() {
	k.startOnce.Do(k.startWork)
}

func (k *offlineKeeper) startWork() {
	for range offlineDownloadWorkers {
		k.workers.Add(1)
		go k.downloadLoop()
	}
	k.mu.Lock()
	roots := k.sortedRootsLocked()
	k.mu.Unlock()
	for _, root := range roots {
		k.beginWalk(root)
		go func() {
			defer k.endWalk(root)
			// A keep or rename cut short by the mount closing may not have
			// saved the folders above the kept path as they are now, which
			// it is reached through offline.
			for _, ancestor := range ancestors(root) {
				next := strings.SplitN(strings.TrimPrefix(strings.TrimPrefix(root, ancestor), "/"), "/", 2)[0]
				saved := func() bool { return k.savedListingHas(ancestor, next) }
				if !saved() {
					// Other kept paths in the folder may list it first.
					k.relistUnless(ancestor, saved)
				}
			}
			k.walk(root)
		}()
	}
}

func (k *offlineKeeper) isClosed() bool {
	select {
	case <-k.done:
		return true
	default:
		return false
	}
}

func (k *offlineKeeper) close() {
	k.queueMu.Lock()
	if !k.closed {
		k.closed = true
		close(k.done)
	}
	k.queueCnd.Broadcast()
	k.queueMu.Unlock()
	k.workers.Wait()
}

// keep keeps a path and everything under it for offline use.
func (k *offlineKeeper) keep(path string) error {
	path = cleanMountPath(path)
	k.mu.Lock()
	if _, covered := k.coveringRootLocked(path); !covered {
		for root := range k.roots {
			if isUnder(root, path) {
				delete(k.roots, root)
			}
		}
		k.roots[path] = struct{}{}
		if err := k.saveRootsLocked(); err != nil {
			k.mu.Unlock()
			return err
		}
	}
	k.walks[path]++
	k.mu.Unlock()
	k.notifyWithParents(path)

	go func() {
		defer k.endWalk(path)
		// The folders above save their listings, so the kept path can be
		// reached when the storage can't.
		for _, ancestor := range ancestors(path) {
			k.relist(ancestor)
		}
		k.walk(path)
	}()
	return nil
}

// release stops keeping a path and everything under it. A path kept because a
// folder above it is kept stops being kept on its own: the folder above stays
// kept for everything else in it.
func (k *offlineKeeper) release(path string) error {
	path = cleanMountPath(path)
	k.mu.Lock()
	if root, covered := k.coveringRootLocked(path); covered && root != path {
		// Keep everything else under the covering root: each folder on the way
		// down to path keeps its other entries.
		delete(k.roots, root)
		dir := root
		for _, part := range strings.Split(strings.Trim(strings.TrimPrefix(path, root), "/"), "/") {
			listing, ok := k.listingLocked(dir)
			if !ok {
				break
			}
			next := path_lib.Join(dir, part)
			for _, entry := range listing.Entries {
				if child := path_lib.Join(dir, entry.Name); child != next {
					k.roots[child] = struct{}{}
				}
			}
			dir = next
		}
	}
	for root := range k.roots {
		if root == path || isUnder(root, path) {
			delete(k.roots, root)
		}
	}
	err := k.saveRootsLocked()

	// Files under path are no longer kept, and folders under it drop their
	// saved listings, unless a folder above a remaining kept path needs them.
	var files []string
	var dirs []string
	k.collectLocked(path, &files, &dirs)
	if _, isDir := k.listingLocked(path); !isDir {
		files = append(files, path)
	}
	for pendingPath := range k.pending {
		if pendingPath == path || isUnder(pendingPath, path) {
			delete(k.pending, pendingPath)
		}
	}
	for _, dir := range dirs {
		if !k.neededLocked(dir) {
			k.forgetListingLocked(dir)
		}
	}
	// The cache learns what is kept while k.mu is held, so a keep that is
	// still listing can't keep these files again after this.
	if cacheErr := k.fs.cacheStore.SetKept(files, false); cacheErr != nil && err == nil {
		err = cacheErr
	}
	k.mu.Unlock()

	// Online Only frees the space right away: released files leave this
	// computer now, or when they are closed.
	for _, file := range files {
		k.fs.deleteCachedFile(file)
	}
	for _, file := range files {
		k.notify(file)
	}
	k.notifyWithParents(path)
	return err
}

// state says how far path is kept for offline use.
func (k *offlineKeeper) state(path string) OfflineState {
	path = cleanMountPath(path)
	k.mu.Lock()
	defer k.mu.Unlock()
	if _, covered := k.coveringRootLocked(path); !covered {
		for root := range k.roots {
			if isUnder(root, path) {
				return OfflinePartlyKept
			}
		}
		return OfflineNone
	}
	for walked := range k.walks {
		if walked == path || isUnder(path, walked) || isUnder(walked, path) {
			return OfflineDownloading
		}
	}
	if _, isDir := k.listingLocked(path); isDir {
		for pendingPath := range k.pending {
			if isUnder(pendingPath, path) {
				return OfflineDownloading
			}
		}
		return OfflineKept
	}
	if _, waiting := k.pending[path]; waiting {
		return OfflineDownloading
	}
	entry, ok := k.entryLocked(path)
	if ok && entry.Dir {
		// Its listing isn't saved yet, so it can't be shown offline.
		return OfflineDownloading
	}
	if ok && k.complete(path, entry) {
		return OfflineKept
	}
	return OfflineDownloading
}

// listed records a successful listing of path. Kept folders save it and keep
// their files, downloading the new and changed ones. Folders above a kept path
// save it so the kept path can be reached offline.
func (k *offlineKeeper) listed(path string, items []files_sdk.File) {
	path = cleanMountPath(path)
	k.mu.Lock()
	_, covered := k.coveringRootLocked(path)
	if !covered && !k.neededLocked(path) {
		k.mu.Unlock()
		return
	}
	listing := &savedListing{Path: path}
	if node, ok := k.fs.vfs.fetch(path); ok && path != "/" {
		info := node.currentInfo()
		self := savedEntry{
			Name:        path_lib.Base(path),
			Dir:         true,
			ModTime:     info.modTime,
			Created:     info.creationTime,
			Permissions: info.remotePermissions,
		}
		listing.Self = &self
	}
	names := make(map[string]struct{}, len(items))
	for _, item := range items {
		if k.fs.isLockFile(item.DisplayName) {
			continue
		}
		listing.Entries = append(listing.Entries, savedEntryFor(item.DisplayName, item))
		names[item.DisplayName] = struct{}{}
	}
	previous, hadPrevious := k.listingLocked(path)

	// What the storage no longer has is no longer kept.
	var removedFiles []string
	var removedDirs []string
	if hadPrevious && covered {
		for _, entry := range previous.Entries {
			if _, still := names[entry.Name]; still {
				continue
			}
			child := path_lib.Join(path, entry.Name)
			if entry.Dir {
				k.collectLocked(child, &removedFiles, &removedDirs)
				removedDirs = append(removedDirs, child)
			} else {
				removedFiles = append(removedFiles, child)
			}
			for pendingPath := range k.pending {
				if pendingPath == child || isUnder(pendingPath, child) {
					delete(k.pending, pendingPath)
				}
			}
		}
	}
	// A kept folder's listing is saved now, so it can be shown offline.
	if covered {
		delete(k.pending, path)
	}
	for _, dir := range removedDirs {
		k.forgetListingLocked(dir)
	}

	var keptFiles []string
	var incomplete []string
	var unlistedDirs []string
	for i, entry := range listing.Entries {
		child := path_lib.Join(path, entry.Name)
		if entry.Dir {
			// A subfolder that was never listed, such as a new one, waits
			// for its listing like a file waits for its download, and is
			// walked so what is in it is kept too. A walk in progress gets to
			// it on its own.
			if _, saved := k.listingLocked(child); covered && !saved {
				k.pending[child] = struct{}{}
				if !k.walkingLocked(child) {
					k.walks[child]++
					unlistedDirs = append(unlistedDirs, child)
				}
			}
			continue
		}
		if _, root := k.roots[child]; !covered && !root {
			continue
		}
		keptFiles = append(keptFiles, child)
		if k.complete(child, entry) {
			continue
		}
		// Until the new version is downloaded, the saved listing shows the
		// last one that is, so it still opens offline.
		if hadPrevious {
			for _, old := range previous.Entries {
				if old.Name == entry.Name && !old.Dir && k.complete(child, old) {
					listing.Entries[i] = old
					break
				}
			}
		}
		k.pending[child] = struct{}{}
		incomplete = append(incomplete, child)
	}
	if err := k.saveListingLocked(listing); err != nil {
		k.fs.log.Debug("offline: saving the listing of %v failed: %v", path, err)
	}
	// The cache learns what is kept while k.mu is held, so an Online Only that
	// comes in between can't be undone by this listing.
	if len(removedFiles) > 0 {
		_ = k.fs.cacheStore.SetKept(removedFiles, false)
	}
	if len(keptFiles) > 0 {
		if err := k.fs.cacheStore.SetKept(keptFiles, true); err != nil {
			k.fs.log.Debug("offline: keeping the files of %v failed: %v", path, err)
		}
	}
	k.mu.Unlock()

	for _, child := range incomplete {
		k.enqueue(child)
	}
	for _, dir := range unlistedDirs {
		go func() {
			defer k.endWalk(dir)
			k.walk(dir)
		}()
	}
}

// savedChildren shows a folder from its saved listing, when there is one.
func (k *offlineKeeper) savedChildren(path string) (map[string]struct{}, bool) {
	path = cleanMountPath(path)
	k.mu.Lock()
	listing, ok := k.listingLocked(path)
	k.mu.Unlock()
	if !ok {
		return nil, false
	}
	childPaths := make(map[string]struct{}, len(listing.Entries))
	for _, entry := range listing.Entries {
		child := path_lib.Join(path, entry.Name)
		childPaths[child] = struct{}{}
		k.fs.createNode(child, entry.file(k.fs.remotePath(child)))
	}
	return childPaths, true
}

// savedFolder returns a folder known from a saved listing.
func (k *offlineKeeper) savedFolder(path string) (files_sdk.File, bool) {
	path = cleanMountPath(path)
	k.mu.Lock()
	defer k.mu.Unlock()
	if path == "/" {
		if len(k.roots) == 0 {
			return files_sdk.File{}, false
		}
		return savedEntry{Name: "/", Dir: true}.file("/"), true
	}
	if entry, ok := k.entryLocked(path); ok && entry.Dir {
		return entry.file(k.fs.remotePath(path)), true
	}
	return files_sdk.File{}, false
}

// active reports whether anything is kept, which is when the mount answers from
// what it saved instead of waiting for storage that can't be reached.
func (k *offlineKeeper) active() bool {
	k.mu.Lock()
	defer k.mu.Unlock()
	return len(k.roots) > 0
}

// unreachable reports whether the storage didn't answer the last time it was
// asked. Saved folders are then shown without waiting for the storage, which
// is asked again in the background once offlineUnreachableWindow has passed.
func (k *offlineKeeper) unreachable() bool {
	k.mu.Lock()
	defer k.mu.Unlock()
	if !k.wasUnreachable {
		return false
	}
	if !k.checking && !time.Now().Before(k.unreachableUntil) {
		k.checking = true
		go k.checkReachable()
	}
	return true
}

// checkReachable lists the root to learn whether the storage answers again.
func (k *offlineKeeper) checkReachable() {
	ctx, cancel := newAnswerDeadline(offlineRemoteTimeout)
	defer cancel()
	_, err := k.fs.listDirContext(ctx, "/")
	k.mu.Lock()
	k.checking = false
	k.mu.Unlock()
	if offlineFallbackAllowed(err) {
		k.markUnreachable()
		return
	}
	k.markReachable()
}

func (k *offlineKeeper) markUnreachable() {
	k.mu.Lock()
	if !k.wasUnreachable {
		k.fs.log.Info("offline: the storage can't be reached, kept files are shown from this computer")
	}
	k.unreachableUntil = time.Now().Add(offlineUnreachableWindow)
	k.wasUnreachable = true
	k.mu.Unlock()
}

func (k *offlineKeeper) markReachable() {
	k.mu.Lock()
	if k.wasUnreachable {
		k.fs.log.Info("offline: the storage can be reached again")
	}
	k.unreachableUntil = time.Time{}
	k.wasUnreachable = false
	k.mu.Unlock()
}

// remoteContext limits how long a call to the storage for path may wait for
// an answer: a path with a saved listing, or a call made while the storage was
// just unreachable, gives up when the storage hasn't answered after
// offlineRemoteTimeout.
func (k *offlineKeeper) remoteContext(path string) (context.Context, context.CancelFunc) {
	path = cleanMountPath(path)
	k.mu.Lock()
	_, saved := k.listingLocked(path)
	if !saved {
		_, saved = k.entryLocked(path)
	}
	bounded := saved || k.wasUnreachable
	k.mu.Unlock()
	if !bounded {
		return context.Background(), func() {}
	}
	return newAnswerDeadline(offlineRemoteTimeout)
}

// answerDeadline is a context that gives up when the storage hasn't answered
// within a timeout, with context.DeadlineExceeded like a context with a
// timeout. A listing that started answering may take longer: the rest of its
// pages come from storage that can be reached. It is its own context, not one
// wrapped from the context package, so contexts made from it report its error.
type answerDeadline struct {
	done  chan struct{}
	timer *time.Timer
	mu    sync.Mutex
	err   error
}

func newAnswerDeadline(timeout time.Duration) (*answerDeadline, context.CancelFunc) {
	d := &answerDeadline{done: make(chan struct{})}
	d.timer = time.AfterFunc(timeout, func() { d.end(context.DeadlineExceeded) })
	return d, func() {
		d.timer.Stop()
		d.end(context.Canceled)
	}
}

func (d *answerDeadline) end(err error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.err == nil {
		d.err = err
		close(d.done)
	}
}

func (d *answerDeadline) Deadline() (time.Time, bool) { return time.Time{}, false }

func (d *answerDeadline) Done() <-chan struct{} { return d.done }

func (d *answerDeadline) Err() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.err
}

func (d *answerDeadline) Value(any) any { return nil }

// storageAnswered lifts the timeout of ctx, when it has one, once the storage
// answered.
func storageAnswered(ctx context.Context) {
	if d, ok := ctx.(*answerDeadline); ok {
		d.timer.Stop()
	}
}

// savedListingHas reports whether folder's saved listing has name in it.
func (k *offlineKeeper) savedListingHas(folder, name string) bool {
	k.mu.Lock()
	defer k.mu.Unlock()
	listing, ok := k.listingLocked(cleanMountPath(folder))
	if !ok {
		return false
	}
	// Saved listings are sorted by name.
	i := sort.Search(len(listing.Entries), func(i int) bool { return listing.Entries[i].Name >= name })
	return i < len(listing.Entries) && listing.Entries[i].Name == name
}

// renameInListingsLocked moves oldPath's entry from its folder's saved listing
// to newPath's folder's, under the new name. The folders then show the
// rename offline, also when the mount closes before they are listed again.
func (k *offlineKeeper) renameInListingsLocked(oldPath, newPath string) {
	oldName, newName := path_lib.Base(oldPath), path_lib.Base(newPath)
	from, ok := k.listingLocked(path_lib.Dir(oldPath))
	if !ok {
		return
	}
	var moved *savedEntry
	for _, entry := range from.Entries {
		if entry.Name == oldName {
			copied := entry
			moved = &copied
		}
	}
	if moved == nil {
		return
	}
	moved.Name = newName
	to := listingWithout(from, oldName)
	if path_lib.Dir(newPath) != path_lib.Dir(oldPath) {
		k.saveListingOrLogLocked(to)
		saved, ok := k.listingLocked(path_lib.Dir(newPath))
		if !ok {
			return
		}
		to = listingWithout(saved, newName)
	} else {
		to = listingWithout(to, newName)
	}
	to.Entries = append(to.Entries, *moved)
	k.saveListingOrLogLocked(to)
}

// listingWithout returns a copy of listing without name. Saved listings don't
// change: a change saves a new one.
func listingWithout(listing *savedListing, name string) *savedListing {
	changed := &savedListing{Path: listing.Path, Self: listing.Self}
	for _, entry := range listing.Entries {
		if entry.Name != name {
			changed.Entries = append(changed.Entries, entry)
		}
	}
	return changed
}

func (k *offlineKeeper) saveListingOrLogLocked(listing *savedListing) {
	if err := k.saveListingLocked(listing); err != nil {
		k.fs.log.Debug("offline: saving the listing of %v failed: %v", listing.Path, err)
	}
}

// hasSavedListing reports whether a folder has a saved listing.
func (k *offlineKeeper) hasSavedListing(path string) bool {
	k.mu.Lock()
	defer k.mu.Unlock()
	_, ok := k.listingLocked(cleanMountPath(path))
	return ok
}

// renamed moves what is kept under oldPath to newPath, after the mount renamed
// it. Nothing stays kept under the old name: its downloads and saved listings
// are dropped and its copies on this computer deleted. The folders above save
// their listings again, with the new name, and the walk of newPath keeps it
// again under the new name.
func (k *offlineKeeper) renamed(oldPath, newPath string) {
	oldPath, newPath = cleanMountPath(oldPath), cleanMountPath(newPath)
	k.mu.Lock()
	var movedRoots []string
	for root := range k.roots {
		if root == oldPath || isUnder(root, oldPath) {
			delete(k.roots, root)
			movedRoots = append(movedRoots, newPath+strings.TrimPrefix(root, oldPath))
		}
	}
	for _, root := range movedRoots {
		k.roots[root] = struct{}{}
	}
	if len(movedRoots) > 0 {
		_ = k.saveRootsLocked()
	}
	var oldFiles []string
	var oldDirs []string
	k.collectLocked(oldPath, &oldFiles, &oldDirs)
	if _, isDir := k.listingLocked(oldPath); !isDir {
		oldFiles = append(oldFiles, oldPath)
	}
	for pendingPath := range k.pending {
		if pendingPath == oldPath || isUnder(pendingPath, oldPath) {
			delete(k.pending, pendingPath)
		}
	}
	// Before the old folders' listings are forgotten, which can be the one
	// the entry moves from.
	k.renameInListingsLocked(oldPath, newPath)
	for _, dir := range oldDirs {
		k.forgetListingLocked(dir)
	}
	for _, ancestor := range ancestors(oldPath) {
		if !k.neededLocked(ancestor) {
			k.forgetListingLocked(ancestor)
		}
	}
	if err := k.fs.cacheStore.SetKept(oldFiles, false); err != nil {
		k.fs.log.Debug("offline: releasing the old name of %v failed: %v", newPath, err)
	}

	// The folders to list again: those above a kept new name, and a saved
	// folder that held the old or the new name.
	var relists []string
	if k.neededLocked(newPath) {
		relists = ancestors(newPath)
	} else if _, saved := k.listingLocked(path_lib.Dir(newPath)); saved {
		relists = []string{path_lib.Dir(newPath)}
	}
	if oldParent := path_lib.Dir(oldPath); oldParent != path_lib.Dir(newPath) {
		if _, saved := k.listingLocked(oldParent); saved {
			relists = append(relists, oldParent)
		}
	}
	var walkPaths []string
	if _, covered := k.coveringRootLocked(newPath); covered {
		walkPaths = append(walkPaths, newPath)
	} else {
		// What is kept inside a renamed folder that isn't kept itself is
		// walked on its own.
		walkPaths = append(walkPaths, movedRoots...)
	}
	for _, path := range walkPaths {
		k.walks[path]++
	}
	k.mu.Unlock()

	go func() {
		// Rename holds the old path's range lock while it calls this, and
		// deletes that path's own copy itself. The new name is walked after
		// the old copies are gone, so moving it back can't lose the copies
		// kept again.
		for _, file := range oldFiles {
			if file != oldPath {
				k.fs.deleteCachedFile(file)
			}
		}
		for _, folder := range relists {
			k.relist(folder)
		}
		for _, path := range walkPaths {
			k.walk(path)
			k.endWalk(path)
		}
	}()
}

// deleted stops keeping what the mount deleted.
func (k *offlineKeeper) deleted(path string) {
	path = cleanMountPath(path)
	k.mu.Lock()
	changed := false
	for root := range k.roots {
		if root == path || isUnder(root, path) {
			delete(k.roots, root)
			changed = true
		}
	}
	if changed {
		_ = k.saveRootsLocked()
	}
	for pendingPath := range k.pending {
		if pendingPath == path || isUnder(pendingPath, path) {
			delete(k.pending, pendingPath)
		}
	}
	k.mu.Unlock()
}

// walk lists a kept folder and the folders under it, which keeps their files
// and downloads the ones not downloaded yet. A kept file is downloaded.
func (k *offlineKeeper) walk(path string) {
	node, ok := k.fs.vfs.fetch(path)
	if !ok {
		// A fresh listing of the folder above finds it, also right after a
		// rename, before the renamed folder's listing counts as stale.
		k.relist(path_lib.Dir(path))
		if node, ok = k.fs.vfs.fetch(path); !ok {
			k.fs.log.Debug("offline: %v can't be found to keep it", path)
			return
		}
	}
	if info := node.currentInfo(); info.nodeType != nodeTypeDir {
		k.mu.Lock()
		// Online Only may have come while the walk waited for the storage.
		if _, covered := k.coveringRootLocked(path); !covered {
			k.mu.Unlock()
			return
		}
		entry := savedEntry{Name: path_lib.Base(path), Size: info.size, ModTime: info.modTime}
		complete := k.complete(path, entry)
		if !complete {
			k.pending[path] = struct{}{}
		}
		if err := k.fs.cacheStore.SetKept([]string{path}, true); err != nil {
			k.fs.log.Debug("offline: keeping %v failed: %v", path, err)
		}
		k.mu.Unlock()
		if !complete {
			k.enqueue(path)
		}
		return
	}
	k.relist(path)
	k.mu.Lock()
	listing, ok := k.listingLocked(path)
	var dirs []string
	if ok {
		for _, entry := range listing.Entries {
			if entry.Dir {
				dirs = append(dirs, path_lib.Join(path, entry.Name))
			}
		}
	}
	k.mu.Unlock()
	for _, dir := range dirs {
		k.walk(dir)
	}
}

// relist lists a folder from the storage, which saves its listing when it is
// kept or above a kept path.
func (k *offlineKeeper) relist(path string) {
	k.relistUnless(path, nil)
}

// relistUnless lists a folder like relist, unless done says it no longer
// needs to once a listing slot is free.
func (k *offlineKeeper) relistUnless(path string, done func() bool) {
	if k.unreachable() {
		return
	}
	// At most offlineListings of these run at once. A mount with many kept
	// paths used to list them all at once after a restart: the listings took
	// every operation slot, so the mount's first lookups failed and Windows
	// didn't mount the drive. Listings still waiting when the mount closes
	// don't run.
	select {
	case k.listingSlots <- struct{}{}:
		defer func() { <-k.listingSlots }()
	case <-k.done:
		return
	}
	if k.isClosed() || (done != nil && done()) {
		return
	}
	node, ok := k.fs.vfs.fetch(path)
	if !ok {
		var errc int
		if node, errc = k.fs.findDir(path); errc != 0 || node == nil {
			return
		}
	}
	k.fs.vfs.expireNodeInfo(path)
	node.expireChildPaths()
	if errc := k.fs.loadDir(node); errc != 0 {
		k.fs.log.Debug("offline: listing %v failed: errc=%d", path, errc)
	}
}

func (k *offlineKeeper) beginWalk(path string) {
	k.mu.Lock()
	k.walks[path]++
	k.mu.Unlock()
}

func (k *offlineKeeper) endWalk(path string) {
	k.mu.Lock()
	k.walks[path]--
	if k.walks[path] <= 0 {
		delete(k.walks, path)
	}
	k.mu.Unlock()
	k.notifyWithParents(path)
}

// walkingLocked reports whether a walk of path or a folder above it is running.
func (k *offlineKeeper) walkingLocked(path string) bool {
	for walked := range k.walks {
		if walked == path || isUnder(path, walked) {
			return true
		}
	}
	return false
}

func (k *offlineKeeper) enqueue(path string) {
	k.queueMu.Lock()
	defer k.queueMu.Unlock()
	if k.closed {
		return
	}
	if _, already := k.queued[path]; already {
		return
	}
	k.queued[path] = struct{}{}
	k.queue = append(k.queue, path)
	k.queueCnd.Signal()
}

func (k *offlineKeeper) downloadLoop() {
	defer k.workers.Done()
	for {
		k.queueMu.Lock()
		for len(k.queue) == 0 && !k.closed {
			k.queueCnd.Wait()
		}
		if k.closed {
			k.queueMu.Unlock()
			return
		}
		path := k.queue[0]
		k.queue = k.queue[1:]
		delete(k.queued, path)
		k.queueMu.Unlock()
		k.download(path)
	}
}

func (k *offlineKeeper) download(path string) {
	k.mu.Lock()
	_, waiting := k.pending[path]
	k.mu.Unlock()
	if !waiting {
		return
	}
	if k.unreachable() {
		k.retryLater(path)
		return
	}
	node, ok := k.fs.vfs.fetch(path)
	if !ok {
		if errc := k.fs.loadParent(path); errc == 0 {
			node, ok = k.fs.vfs.fetch(path)
		}
	}
	if !ok {
		k.retryLater(path)
		return
	}
	info := node.currentInfo()
	version := savedEntry{Name: path_lib.Base(path), Size: info.size, ModTime: info.modTime}
	err := k.fs.ensureFullyCached(path, version.Size)
	if err != nil {
		k.fs.log.Debug("offline: downloading %v failed: %v", path, err)
		k.retryLater(path)
		return
	}
	k.mu.Lock()
	delete(k.pending, path)
	k.savedVersionLocked(path, version)
	k.mu.Unlock()
	k.notifyWithParents(path)
}

// savedVersionLocked makes the saved listing show the version of a file that
// was just downloaded, in place of the older one it showed until then.
func (k *offlineKeeper) savedVersionLocked(path string, version savedEntry) {
	listing, ok := k.listingLocked(path_lib.Dir(path))
	if !ok || !k.complete(path, version) {
		return
	}
	for i, entry := range listing.Entries {
		if entry.Name != version.Name || entry.Dir {
			continue
		}
		if entry.Size != version.Size || !entry.ModTime.Equal(version.ModTime) {
			updated := &savedListing{Path: listing.Path, Self: listing.Self, Entries: slices.Clone(listing.Entries)}
			updated.Entries[i].Size = version.Size
			updated.Entries[i].ModTime = version.ModTime
			if err := k.saveListingLocked(updated); err != nil {
				k.fs.log.Debug("offline: saving the listing of %v failed: %v", path_lib.Dir(path), err)
			}
		}
		return
	}
}

func (k *offlineKeeper) retryLater(path string) {
	time.AfterFunc(offlineRetryDelay, func() {
		k.mu.Lock()
		_, waiting := k.pending[path]
		k.mu.Unlock()
		if waiting {
			k.enqueue(path)
		}
	})
}

// complete reports whether the cache holds all of the saved version of a file.
// It runs with k.mu held.
func (k *offlineKeeper) complete(path string, entry savedEntry) bool {
	if entry.Size <= 0 {
		return true
	}
	complete, err := k.fs.cacheStore.RangeEntryComplete(path, cacheEntryMetadata(path, entry.Size, entry.ModTime))
	return err == nil && complete
}

// coveringRootLocked returns the kept path that path is, or is under.
func (k *offlineKeeper) coveringRootLocked(path string) (string, bool) {
	for p := path; ; p = path_lib.Dir(p) {
		if _, ok := k.roots[p]; ok {
			return p, true
		}
		if p == "/" {
			return "", false
		}
	}
}

// neededLocked reports whether a folder's listing is needed to reach a kept path.
func (k *offlineKeeper) neededLocked(path string) bool {
	for root := range k.roots {
		if root == path || isUnder(root, path) || isUnder(path, root) {
			return true
		}
	}
	return false
}

// collectLocked gathers the files and folders under path from saved listings.
func (k *offlineKeeper) collectLocked(path string, files *[]string, dirs *[]string) {
	listing, ok := k.listingLocked(path)
	if !ok {
		return
	}
	*dirs = append(*dirs, path)
	for _, entry := range listing.Entries {
		child := path_lib.Join(path, entry.Name)
		if entry.Dir {
			k.collectLocked(child, files, dirs)
		} else {
			*files = append(*files, child)
		}
	}
}

func (k *offlineKeeper) entryLocked(path string) (savedEntry, bool) {
	if listing, ok := k.listingLocked(path); ok && listing.Self != nil {
		return *listing.Self, true
	}
	parent, ok := k.listingLocked(path_lib.Dir(path))
	if !ok {
		return savedEntry{}, false
	}
	name := path_lib.Base(path)
	for _, entry := range parent.Entries {
		if entry.Name == name {
			return entry, true
		}
	}
	return savedEntry{}, false
}

func (k *offlineKeeper) listingLocked(path string) (*savedListing, bool) {
	if listing, ok := k.listings[path]; ok {
		return listing, listing != nil
	}
	data, err := os.ReadFile(k.listingFile(path))
	if err != nil {
		k.listings[path] = nil
		return nil, false
	}
	var listing savedListing
	if err := json.Unmarshal(data, &listing); err != nil || listing.Path != path {
		k.listings[path] = nil
		return nil, false
	}
	k.listings[path] = &listing
	return &listing, true
}

func (k *offlineKeeper) saveListingLocked(listing *savedListing) error {
	sort.Slice(listing.Entries, func(i, j int) bool { return listing.Entries[i].Name < listing.Entries[j].Name })
	k.listings[listing.Path] = listing
	data, err := json.Marshal(listing)
	if err != nil {
		return err
	}
	return os.WriteFile(k.listingFile(listing.Path), data, 0o600)
}

func (k *offlineKeeper) forgetListingLocked(path string) {
	k.listings[path] = nil
	_ = os.Remove(k.listingFile(path))
}

func (k *offlineKeeper) listingFile(path string) string {
	sum := sha256.Sum256([]byte(path))
	return filepath.Join(k.dir, offlineListingsDir, hex.EncodeToString(sum[:])+".json")
}

func (k *offlineKeeper) sortedRootsLocked() []string {
	roots := make([]string, 0, len(k.roots))
	for root := range k.roots {
		roots = append(roots, root)
	}
	sort.Strings(roots)
	return roots
}

func (k *offlineKeeper) saveRootsLocked() error {
	data, err := json.Marshal(k.sortedRootsLocked())
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(k.dir, offlineRootsFile), data, 0o600)
}

// notifyWithParents tells Explorer that path and the folders above it may show
// a different overlay icon.
func (k *offlineKeeper) notifyWithParents(path string) {
	k.notify(path)
	for _, ancestor := range ancestors(path) {
		if ancestor != "/" {
			k.notify(ancestor)
		}
	}
}

func (k *offlineKeeper) notify(path string) {
	_ = shell.NotifyUpdatedItem(k.fs.localPath(path))
}

// listDirOrSaved lists a folder from the storage, or shows its saved listing
// when the storage can't be reached.
func (fs *RemoteFs) listDirOrSaved(path string) (map[string]struct{}, error) {
	k := fs.offline
	if k == nil || !k.active() {
		return fs.listDir(path)
	}
	if k.unreachable() {
		if childPaths, ok := k.savedChildren(path); ok {
			return childPaths, nil
		}
		// Folders that aren't saved fail right away while the storage can't
		// be reached, instead of waiting for it.
		return nil, errStorageUnreachable
	}
	ctx, cancel := k.remoteContext(path)
	defer cancel()
	childPaths, err := fs.listDirContext(ctx, path)
	if err == nil {
		k.markReachable()
		return childPaths, nil
	}
	if offlineFallbackAllowed(err) {
		k.markUnreachable()
		if saved, ok := k.savedChildren(path); ok {
			fs.log.Debug("offline: showing the saved listing of %v: %v", path, err)
			return saved, nil
		}
	}
	return childPaths, err
}

// deleteCachedFile deletes what is cached of path. An open file keeps its data
// until it is closed, so its readers aren't cut off.
func (fs *RemoteFs) deleteCachedFile(path string) {
	if fs.cacheStore.DeleteWhenClosed(path) {
		return
	}
	unlockRangeMutation := fs.lockRangeMutation(path)
	defer unlockRangeMutation()
	fs.cancelRangeDownloads(path)
	_ = fs.cacheStore.Delete(path)
}

// errStorageUnreachable is the error of a call not made because the storage
// can't be reached.
var errStorageUnreachable = errors.New("the storage can't be reached")

// offlineFallbackAllowed reports whether a failure may be answered from what
// is saved: the storage didn't answer, rather than answering with an error of
// its own, such as a missing path or access denied. An API error is the
// storage's answer, unless it says the service is unavailable.
func offlineFallbackAllowed(err error) bool {
	if class, _ := classifyMountError(err); class != "unknown" {
		return false
	}
	if _, answered := files_sdk.ResponseErrorTypeOf(err); answered {
		return files_sdk.IsServiceUnavailable(err)
	}
	return true
}

// KeepOffline keeps a file or folder of the mount for offline use, or stops
// keeping it. path is relative to the mount.
func (fs *RemoteFs) KeepOffline(path string, keep bool) error {
	if fs.offline == nil {
		return errors.New("keeping files offline needs the disk cache")
	}
	if keep {
		return fs.offline.keep(path)
	}
	return fs.offline.release(path)
}

// OfflineStates says how far each path of the mount is kept for offline use.
func (fs *RemoteFs) OfflineStates(paths []string) []OfflineState {
	states := make([]OfflineState, len(paths))
	if fs.offline == nil {
		return states
	}
	for i, path := range paths {
		states[i] = fs.offline.state(path)
	}
	return states
}

// HasKeptFiles reports whether the disk cache at diskCachePath keeps files
// offline, so a mount using it can start while its storage can't be reached.
func HasKeptFiles(diskCachePath string) bool {
	data, err := os.ReadFile(filepath.Join(disk.OfflineDirFor(diskCachePath), offlineRootsFile))
	if err != nil {
		return false
	}
	var roots []string
	return json.Unmarshal(data, &roots) == nil && len(roots) > 0
}

// cleanMountPath turns a path relative to the mount into the mount's form:
// forward slashes and a leading slash.
func cleanMountPath(path string) string {
	path = strings.ReplaceAll(path, "\\", "/")
	return path_lib.Clean("/" + path)
}

// isUnder reports whether path is inside the folder dir.
func isUnder(path, dir string) bool {
	if dir == "/" {
		return path != "/"
	}
	return strings.HasPrefix(path, dir+"/")
}

// ancestors returns the folders above path, from the mount root down.
func ancestors(path string) []string {
	var result []string
	for p := path_lib.Dir(path); ; p = path_lib.Dir(p) {
		result = append([]string{p}, result...)
		if p == "/" {
			break
		}
	}
	if path == "/" {
		return nil
	}
	return result
}
