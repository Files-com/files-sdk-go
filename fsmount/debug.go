//go:build filescomfs_debug
// +build filescomfs_debug

package fsmount

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"sort"
	"strconv"
	"time"
)

const (
	pprofHostDefault = "localhost"
	pprofPortDefault = 6060
)

func (fs *Filescomfs) startPprof() {
	mux := fs.debugMux()
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
			fs.log.Error("error parsing FILESCOMFS_DEBUG_PPROF_PORT environment variable: %v, defaulting to %d", err, pprofPortDefault)
		}
	}

	pprofAddr := fmt.Sprintf("%s:%d", pprofHost, pprofPort)
	fs.dbgSrv = &http.Server{Addr: pprofAddr, Handler: mux}
	go func() {
		if err := fs.dbgSrv.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				fs.log.Error("debug server in debug.go: %v", err)
			} else {
				fs.log.Info("debug server shut down successfully")
			}
		}
	}()
	fs.log.Info("debug server listening on %s", pprofAddr)
}

// debugMux creates an *httpServeMux that exposes pprof handlers and handlers
// that expose internal file system state for use in debugging.
func (fs *Filescomfs) debugMux() *http.ServeMux {
	mux := http.NewServeMux()

	// ---- pprof endpoints ----
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	// ---- JSON endpoints ----
	mux.HandleFunc("/debug/state", fs.handleDebugState)
	mux.HandleFunc("/debug/handles", fs.handleDebugHandles)
	mux.HandleFunc("/debug/uploads", fs.handleDebugUploads)
	mux.HandleFunc("/debug/writers", fs.handleDebugWriters)
	mux.HandleFunc("/debug/nodes", fs.handleDebugNodes)
	mux.HandleFunc("/debug/locks", fs.handleDebugLocks)

	// simple index
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `<html><body style="font-family: system-ui, sans-serif">
<h1>Files.com FUSE Debug</h1>
<ul>
  <li><a href="/debug/pprof/">/debug/pprof/</a></li>
  <li><a href="/debug/state">/debug/state</a></li>
  <li><a href="/debug/handles">/debug/handles</a></li>
  <li><a href="/debug/uploads">/debug/uploads</a></li>
  <li><a href="/debug/writers">/debug/writers</a></li>
  <li><a href="/debug/nodes">/debug/nodes</a></li>
  <li><a href="/debug/locks">/debug/locks</a></li>
</ul>
</body></html>`)
	})

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
	HasCancel    bool      `json:"hasCancel"`
	WriterOpen   bool      `json:"writerOpen"`
	Committed    bool      `json:"committed"`
}

type dbgWriter struct {
	Path      string `json:"path"`
	OwnerFH   uint64 `json:"ownerFh"`
	Offset    int64  `json:"offset"`
	Committed bool   `json:"committed"`
}

type dbgNode struct {
	Path        string       `json:"path"`
	Size        int64        `json:"size"`
	ModTime     time.Time    `json:"modTime"`
	DownloadURI bool         `json:"downloadUriCached"`
	HasWriter   bool         `json:"hasWriter"`
	HasUpload   bool         `json:"hasUpload"`
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

// /debug/state — quick overview
func (fs *Filescomfs) handleDebugState(w http.ResponseWriter, r *http.Request) {
	type state struct {
		Now           time.Time `json:"now"`
		NumHandles    int       `json:"numHandles"`
		NumNodes      int       `json:"numNodes"`
		NumUploads    int       `json:"numUploads"`
		NumWriters    int       `json:"numWriters"`
		SampleUploads []dbgUpload
		SampleWriters []dbgWriter
	}
	now := time.Now()

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
	writers := make([]dbgWriter, 0)
	for _, n := range nodes {
		// upload snapshot
		n.statusMu.Lock()
		u := n.upload
		path := n.path
		downloadURI := (n.downloadUri != "")
		size := n.info.size
		mod := n.info.modTime
		n.statusMu.Unlock()

		n.writeMu.Lock()
		w := n.writer
		owner := n.writerOwner
		var off int64
		if w != nil {
			off = w.Offset()
		}
		n.writeMu.Unlock()

		if u != nil {
			u.mu.Lock()
			uploads = append(uploads, dbgUpload{
				Path:         path,
				Ref:          u.ref,
				BytesWritten: u.bytesWritten,
				LastActivity: u.lastActivity,
				HasCancel:    u.cancel != nil,
				WriterOpen:   w != nil,
				Committed:    w != nil && off > 0,
			})
			u.mu.Unlock()
		}
		if w != nil {
			writers = append(writers, dbgWriter{
				Path:      path,
				OwnerFH:   owner,
				Offset:    off,
				Committed: off > 0,
			})
		}

		_ = downloadURI
		_ = size
		_ = mod
	}

	// small samples (sorted by path for determinism)
	sort.Slice(uploads, func(i, j int) bool { return uploads[i].Path < uploads[j].Path })
	sort.Slice(writers, func(i, j int) bool { return writers[i].Path < writers[j].Path })

	resp := state{
		Now:           now,
		NumHandles:    handleCount,
		NumNodes:      len(nodes),
		NumUploads:    len(uploads),
		NumWriters:    len(writers),
		SampleUploads: uploads,
		SampleWriters: writers,
	}
	writeJSON(w, resp)
}

// /debug/handles — enumerate open file handles
func (fs *Filescomfs) handleDebugHandles(w http.ResponseWriter, r *http.Request) {
	out := make([]dbgHandle, 0)

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
func (fs *Filescomfs) handleDebugUploads(w http.ResponseWriter, r *http.Request) {
	out := make([]dbgUpload, 0)
	for _, n := range fs.snapshotNodes() {
		n.statusMu.Lock()
		u := n.upload
		path := n.path
		n.statusMu.Unlock()
		if u == nil {
			continue
		}

		n.writeMu.Lock()
		wr := n.writer
		var off int64
		if wr != nil {
			off = wr.Offset()
		}
		n.writeMu.Unlock()

		u.mu.Lock()
		out = append(out, dbgUpload{
			Path:         path,
			Ref:          u.ref,
			BytesWritten: u.bytesWritten,
			LastActivity: u.lastActivity,
			HasCancel:    u.cancel != nil,
			WriterOpen:   wr != nil,
			Committed:    wr != nil && off > 0,
		})
		u.mu.Unlock()
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	writeJSON(w, out)
}

// /debug/writers — nodes with open writers
func (fs *Filescomfs) handleDebugWriters(w http.ResponseWriter, r *http.Request) {
	out := make([]dbgWriter, 0)
	for _, n := range fs.snapshotNodes() {
		n.writeMu.Lock()
		w := n.writer
		owner := n.writerOwner
		var off int64
		if w != nil {
			off = w.Offset()
		}
		path := n.path
		n.writeMu.Unlock()

		if w != nil {
			out = append(out, dbgWriter{
				Path:      path,
				OwnerFH:   owner,
				Offset:    off,
				Committed: off > 0,
			})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	writeJSON(w, out)
}

// /debug/nodes — light metadata for all nodes
func (fs *Filescomfs) handleDebugNodes(w http.ResponseWriter, r *http.Request) {
	out := make([]dbgNode, 0)
	for _, n := range fs.snapshotNodes() {
		// status fields
		n.statusMu.Lock()
		path := n.path
		uriCached := n.downloadUri != ""
		size := n.info.size
		mod := n.info.modTime
		hasUpload := n.upload != nil
		expires := n.infoExpires
		info := n.info
		n.statusMu.Unlock()

		// writer fields
		n.writeMu.Lock()
		hasWriter := (n.writer != nil)
		n.writeMu.Unlock()

		out = append(out, dbgNode{
			Path:        path,
			Size:        size,
			ModTime:     mod,
			DownloadURI: uriCached,
			HasWriter:   hasWriter,
			HasUpload:   hasUpload,
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
func (fs *Filescomfs) handleDebugLocks(w http.ResponseWriter, r *http.Request) {
	out := make([]dbgLock, 0)

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
