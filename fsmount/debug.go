//go:build filescomfs_debug
// +build filescomfs_debug

package fsmount

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"net/http/pprof"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"

	_ "embed"

	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache"
	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache/disk"
	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache/mem"
)

const (
	pprofHostDefault = "localhost"
	pprofPortDefault = 6060

	dbgShutdownTimeout = 5 * time.Second
)

var pprofMu sync.Mutex

// cacheStoreWithStats extends cacheStore with statistics gathering methods.
// This interface is only available when building with the filescomfs_debug tag.
type cacheStoreWithStats interface {
	cacheStore
	Stats() *cache.Stats
}

// Compile-time verification that our cache implementations satisfy the extended interface
var _ cacheStoreWithStats = (*disk.DiskCache)(nil)
var _ cacheStoreWithStats = (*mem.MemoryCache)(nil)

//go:embed templates/debug.tmpl
var templateData []byte

func (reg *mountRegistry) startPprof() {
	pprofMu.Lock()
	defer pprofMu.Unlock()
	var pprofAddr string
	if reg.dbgSrv == nil {
		mux := reg.debugMux()
		pprofHost := os.Getenv("FILESCOMFS_DEBUG_PPROF_HOST")
		if pprofHost == "" {
			pprofHost = pprofHostDefault
		}
		envPort := os.Getenv("FILESCOMFS_DEBUG_PPROF_PORT")
		pprofPort := pprofPortDefault
		var err error
		if envPort != "" {
			pprofPort, err = strconv.Atoi(envPort)
			if err != nil {
				reg.log.Error("error parsing FILESCOMFS_DEBUG_PPROF_PORT environment variable: %v, defaulting to %d", err, pprofPortDefault)
			}
		}

		pprofAddr = fmt.Sprintf("%s:%d", pprofHost, pprofPort)
		reg.dbgSrv = &http.Server{Addr: pprofAddr, Handler: mux}
		go func() {
			if err := reg.dbgSrv.ListenAndServe(); err != nil {
				if !errors.Is(err, http.ErrServerClosed) {
					reg.log.Error("debug server in debug.go: %v", err)
				} else {
					reg.log.Info("debug server shut down successfully")
				}
			}
		}()
	}
	reg.log.Info("debug server listening on %s", pprofAddr)
}

func (reg *mountRegistry) stopPprof() {
	pprofMu.Lock()
	defer pprofMu.Unlock()
	// Shutdown the debug server if it was started
	if reg.dbgSrv != nil {
		reg.log.Info("Shutting down debug server")
		ctx, cancel := context.WithTimeout(context.Background(), dbgShutdownTimeout)
		defer cancel()
		if err := reg.dbgSrv.Shutdown(ctx); err != nil {
			reg.log.Error("error shutting down debug server: %v", err)
		} else {
			reg.log.Info("debug server shut down successfully")
		}
	}
	reg.dbgSrv = nil
}

// debugMux creates an *httpServeMux that exposes pprof handlers and handlers
// that expose internal file system state for use in debugging.
func (reg *mountRegistry) debugMux() *http.ServeMux {
	mux := http.NewServeMux()

	// root handler
	mux.HandleFunc("/", reg.handleDebugRoot)

	// ---- pprof endpoints ----
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	// ---- JSON endpoints ----
	mux.HandleFunc("/debug/state", reg.handleDebugState)
	mux.HandleFunc("/debug/handles", reg.handleDebugHandles)
	mux.HandleFunc("/debug/uploads", reg.handleDebugUploads)
	mux.HandleFunc("/debug/writers", reg.handleDebugWriters)
	mux.HandleFunc("/debug/nodes", reg.handleDebugNodes)
	mux.HandleFunc("/debug/locks", reg.handleDebugLocks)
	mux.HandleFunc("/debug/cache", reg.handleDebugCacheStats)

	return mux
}

// -------------------- handlers & helpers --------------------

type dbgHandle struct {
	ID        uint64    `json:"id"`
	Path      string    `json:"path"`
	ReadOnly  bool      `json:"readOnly"`
	BytesRead int64     `json:"bytesRead"`
	ReadAt    time.Time `json:"readAt"`
	// You can add more fields if your fileHandle exposes them safely
}

type dbgUpload struct {
	Path         string    `json:"path"`
	Ref          string    `json:"ref"`
	BytesWritten int64     `json:"bytesWritten"`
	LastActivity time.Time `json:"lastActivity"`
	Uploading    bool      `json:"uploading"`
	Finalizing   bool      `json:"finalizing"`
	Dirty        bool      `json:"dirty"`
}

type dbgSession struct {
	Path            string    `json:"path"`
	WorkingCopyPath string    `json:"workingCopyPath"`
	HandleCount     int       `json:"handleCount"`
	BaselineSize    int64     `json:"baselineSize"`
	CurrentSize     int64     `json:"currentSize"`
	Mtime           time.Time `json:"mtime"`
	Hydrated        bool      `json:"hydrated"`
	Dirty           bool      `json:"dirty"`
	Uploading       bool      `json:"uploading"`
	Finalizing      bool      `json:"finalizing"`
	LastUploadErr   string    `json:"lastUploadErr,omitempty"`
}

type dbgNode struct {
	Path        string       `json:"path"`
	Size        int64        `json:"size"`
	ModTime     time.Time    `json:"modTime"`
	DownloadURI bool         `json:"downloadUriCached"`
	HasSession  bool         `json:"hasSession"`
	HasUpload   bool         `json:"hasUpload"`
	Session     *dbgSession  `json:"session,omitempty"`
	Info        *dbgNodeInfo `json:"info"`
	Now         time.Time    `json:"now"`
	InfoExpires time.Time    `json:"infoExpires"`
	InfoExpired bool         `json:"infoExpired"`
}

type dbgNodeInfo struct {
	NodeType  string    `json:"nodeType"`
	Size      int64     `json:"size"`
	Created   time.Time `json:"created"`
	Modified  time.Time `json:"modified"`
	LockOwner string    `json:"lockOwner"`
}

func (reg *mountRegistry) list() []string {
	reg.hostsMu.Lock()
	defer reg.hostsMu.Unlock()
	mounts := make([]string, 0, len(reg.hosts))
	for mnt := range reg.hosts {
		mounts = append(mounts, mnt)
	}
	// sort the mount points for consistent ordering
	sort.Strings(mounts)
	return mounts
}

func (reg *mountRegistry) handleDebugRoot(w http.ResponseWriter, r *http.Request) {
	var tmplFile = "debug.tmpl"
	tmpl, err := template.New(tmplFile).Parse(string(templateData))
	if err != nil {
		writeJSON(w, map[string]string{"error": "error parsing template: " + err.Error()})
		return
	}

	mounts := reg.list()
	// simple index
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = tmpl.Execute(w, mounts)
}

// /debug/state — quick overview
func (reg *mountRegistry) handleDebugState(w http.ResponseWriter, r *http.Request) {
	type state struct {
		Now            time.Time `json:"now"`
		NumHandles     int       `json:"numHandles"`
		NumNodes       int       `json:"numNodes"`
		NumUploads     int       `json:"numUploads"`
		NumSessions    int       `json:"numSessions"`
		SampleUploads  []dbgUpload
		SampleSessions []dbgSession
	}
	now := time.Now()

	mnt := r.URL.Query().Get("mnt")

	if mnt == "" {
		writeJSON(w, map[string]string{"error": "missing 'mnt' query parameter"})
		return
	}

	host, ok := reg.get(mnt)
	if !ok {
		writeJSON(w, map[string]string{"error": "no such mount point"})
		return
	}
	fs := host.fs

	// Snapshot handles
	var handleCount int
	if fs.vfs != nil && fs.vfs.handles != nil {
		fs.vfs.handles.mu.Lock()
		handleCount = len(fs.vfs.handles.entries)
		fs.vfs.handles.mu.Unlock()
	}

	// Snapshot nodes
	nodes := fs.snapshotNodes()

	uploads := make([]dbgUpload, 0)
	sessions := make([]dbgSession, 0)
	for _, n := range nodes {
		if session := n.getWriteSession(); session != nil {
			snap := session.snapshot()
			session.mu.Lock()
			upload := session.upload
			handleCount := len(session.handles)
			workingCopyPath := session.workingCopyPath
			session.mu.Unlock()

			dbgSession := dbgSession{
				Path:            n.path,
				WorkingCopyPath: workingCopyPath,
				HandleCount:     handleCount,
				BaselineSize:    session.baselineSize,
				CurrentSize:     snap.currentSize,
				Mtime:           snap.mtime,
				Hydrated:        snap.hydrated,
				Dirty:           snap.dirty,
				Uploading:       snap.uploading,
				Finalizing:      snap.finalizing,
			}
			if snap.lastUploadErr != nil {
				dbgSession.LastUploadErr = snap.lastUploadErr.Error()
			}
			sessions = append(sessions, dbgSession)

			if upload != nil {
				ref, bytesWritten, lastActivity := upload.stats()
				uploads = append(uploads, dbgUpload{
					Path:         n.path,
					Ref:          ref,
					BytesWritten: bytesWritten,
					LastActivity: lastActivity,
					Uploading:    snap.uploading,
					Finalizing:   snap.finalizing,
					Dirty:        snap.dirty,
				})
			}
		}
	}

	// small samples (sorted by path for determinism)
	sort.Slice(uploads, func(i, j int) bool { return uploads[i].Path < uploads[j].Path })
	sort.Slice(sessions, func(i, j int) bool { return sessions[i].Path < sessions[j].Path })

	resp := state{
		Now:            now,
		NumHandles:     handleCount,
		NumNodes:       len(nodes),
		NumUploads:     len(uploads),
		NumSessions:    len(sessions),
		SampleUploads:  uploads,
		SampleSessions: sessions,
	}
	writeJSON(w, resp)
}

// /debug/handles — enumerate open file handles
func (reg *mountRegistry) handleDebugHandles(w http.ResponseWriter, r *http.Request) {
	out := make([]dbgHandle, 0)

	mnt := r.URL.Query().Get("mnt")

	if mnt == "" {
		writeJSON(w, map[string]string{"error": "missing 'mnt' query parameter"})
		return
	}

	host, ok := reg.get(mnt)
	if !ok {
		writeJSON(w, map[string]string{"error": "no such mount point"})
		return
	}
	fs := host.fs

	if fs.vfs != nil && fs.vfs.handles != nil {
		fs.vfs.handles.mu.Lock()
		for id, fh := range fs.vfs.handles.entries {
			if fh == nil || fh.node == nil {
				continue
			}
			item := dbgHandle{
				ID:        id,
				Path:      fh.node.path,
				ReadOnly:  fh.IsReadOnly(),
				BytesRead: fh.bytesRead.Load(),
				ReadAt:    fh.readAt,
			}
			out = append(out, item)
		}
		fs.vfs.handles.mu.Unlock()
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	writeJSON(w, out)
}

// /debug/uploads — nodes with active uploads
func (reg *mountRegistry) handleDebugUploads(w http.ResponseWriter, r *http.Request) {
	out := make([]dbgUpload, 0)

	mnt := r.URL.Query().Get("mnt")

	if mnt == "" {
		writeJSON(w, map[string]string{"error": "missing 'mnt' query parameter"})
		return
	}

	host, ok := reg.get(mnt)
	if !ok {
		writeJSON(w, map[string]string{"error": "no such mount point"})
		return
	}
	fs := host.fs

	for _, n := range fs.snapshotNodes() {
		session := n.getWriteSession()
		if session == nil {
			continue
		}
		snap := session.snapshot()
		session.mu.Lock()
		u := session.upload
		session.mu.Unlock()
		if u == nil {
			continue
		}
		ref, bytesWritten, lastActivity := u.stats()
		out = append(out, dbgUpload{
			Path:         n.path,
			Ref:          ref,
			BytesWritten: bytesWritten,
			LastActivity: lastActivity,
			Uploading:    snap.uploading,
			Finalizing:   snap.finalizing,
			Dirty:        snap.dirty,
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	writeJSON(w, out)
}

// /debug/writers — nodes with active write sessions
func (reg *mountRegistry) handleDebugWriters(w http.ResponseWriter, r *http.Request) {
	out := make([]dbgSession, 0)

	mnt := r.URL.Query().Get("mnt")

	if mnt == "" {
		writeJSON(w, map[string]string{"error": "missing 'mnt' query parameter"})
		return
	}

	host, ok := reg.get(mnt)
	if !ok {
		writeJSON(w, map[string]string{"error": "no such mount point"})
		return
	}
	fs := host.fs

	for _, n := range fs.snapshotNodes() {
		if session := n.getWriteSession(); session != nil {
			snap := session.snapshot()
			session.mu.Lock()
			handleCount := len(session.handles)
			workingCopyPath := session.workingCopyPath
			session.mu.Unlock()
			item := dbgSession{
				Path:            n.path,
				WorkingCopyPath: workingCopyPath,
				HandleCount:     handleCount,
				BaselineSize:    session.baselineSize,
				CurrentSize:     snap.currentSize,
				Mtime:           snap.mtime,
				Hydrated:        snap.hydrated,
				Dirty:           snap.dirty,
				Uploading:       snap.uploading,
				Finalizing:      snap.finalizing,
			}
			if snap.lastUploadErr != nil {
				item.LastUploadErr = snap.lastUploadErr.Error()
			}
			out = append(out, item)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	writeJSON(w, out)
}

// /debug/nodes — light metadata for all nodes
func (reg *mountRegistry) handleDebugNodes(w http.ResponseWriter, r *http.Request) {
	out := make([]dbgNode, 0)

	mnt := r.URL.Query().Get("mnt")

	if mnt == "" {
		writeJSON(w, map[string]string{"error": "missing 'mnt' query parameter"})
		return
	}

	host, ok := reg.get(mnt)
	if !ok {
		writeJSON(w, map[string]string{"error": "no such mount point"})
		return
	}
	fs := host.fs

	for _, n := range fs.snapshotNodes() {
		// status fields
		n.statusMu.Lock()
		path := n.path
		uriCached := n.downloadUri != ""
		size := n.info.size
		mod := n.info.modTime
		expires := n.infoExpires
		info := n.info
		n.statusMu.Unlock()

		hasSession := n.getWriteSession() != nil
		hasUpload := false
		var sessionInfo *dbgSession
		if session := n.getWriteSession(); session != nil {
			snap := session.snapshot()
			session.mu.Lock()
			hasUpload = session.upload != nil
			handleCount := len(session.handles)
			workingCopyPath := session.workingCopyPath
			session.mu.Unlock()
			sessionInfo = &dbgSession{
				Path:            path,
				WorkingCopyPath: workingCopyPath,
				HandleCount:     handleCount,
				BaselineSize:    session.baselineSize,
				CurrentSize:     snap.currentSize,
				Mtime:           snap.mtime,
				Hydrated:        snap.hydrated,
				Dirty:           snap.dirty,
				Uploading:       snap.uploading,
				Finalizing:      snap.finalizing,
			}
			if snap.lastUploadErr != nil {
				sessionInfo.LastUploadErr = snap.lastUploadErr.Error()
			}
		}

		out = append(out, dbgNode{
			Path:        path,
			Size:        size,
			ModTime:     mod,
			DownloadURI: uriCached,
			HasSession:  hasSession,
			HasUpload:   hasUpload,
			Session:     sessionInfo,
			Now:         time.Now(),
			InfoExpires: expires,
			InfoExpired: n.infoExpired(),
			Info: &dbgNodeInfo{
				NodeType: info.nodeType.String(),
				Size:     info.size,
				Created:  info.creationTime,
				Modified: info.modTime,
			},
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	writeJSON(w, out)
}

func (reg *mountRegistry) handleDebugCacheStats(w http.ResponseWriter, r *http.Request) {
	mnt := r.URL.Query().Get("mnt")

	if mnt == "" {
		writeJSON(w, map[string]string{"error": "missing 'mnt' query parameter"})
		return
	}

	host, ok := reg.get(mnt)
	if !ok {
		writeJSON(w, map[string]string{"error": "no such mount point"})
		return
	}
	fs := host.fs
	if fs.remote.cacheStore == nil {
		writeJSON(w, map[string]string{"error": "disk cache not enabled for this mount"})
		return
	}

	// Type-assert to the extended interface that includes Stats()
	statsCache, ok := fs.remote.cacheStore.(cacheStoreWithStats)
	if !ok {
		writeJSON(w, map[string]string{"error": "cache does not support stats"})
		return
	}

	writeJSON(w, statsCache.Stats())
}

// snapshotNodes returns a stable slice of *fsNode without holding the VFS lock.
func (fs *Filescomfs) snapshotNodes() []*fsNode {
	if fs.vfs == nil {
		return nil
	}
	fs.vfs.nodesMu.Lock()
	nodes := make([]*fsNode, 0, len(fs.vfs.nodes))
	for _, n := range fs.vfs.nodes {
		if n != nil {
			nodes = append(nodes, n)
		}
	}
	fs.vfs.nodesMu.Unlock()
	return nodes
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
}

// dbgLock is intentionally generic: Raw will include any exported fields
// from your internal lock entry type (whatever fs.lockMap stores).
type dbgLock struct {
	Path string      `json:"path"`
	Raw  interface{} `json:"raw,omitempty"`
}

// /debug/locks — enumerate current lockMap entries
func (reg *mountRegistry) handleDebugLocks(w http.ResponseWriter, r *http.Request) {
	out := make([]dbgLock, 0)

	mnt := r.URL.Query().Get("mnt")

	if mnt == "" {
		writeJSON(w, map[string]string{"error": "missing 'mnt' query parameter"})
		return
	}

	host, ok := reg.get(mnt)
	if !ok {
		writeJSON(w, map[string]string{"error": "no such mount point"})
		return
	}
	fs := host.fs

	// Snapshot under lockMapMutex, but don’t hold it while encoding.
	fs.remote.lockMapMutex.Lock()
	for p, li := range fs.remote.lockMap {
		out = append(out, dbgLock{
			Path: p,
			Raw:  li, // will marshal any exported fields on your lock entry type
		})
	}
	fs.remote.lockMapMutex.Unlock()

	// keep output stable/diff-friendly
	sort.Slice(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	writeJSON(w, out)
}
