package file

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	goFs "io/fs"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/folder"
	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/samber/lo"
)

type SizeTrust int

const (
	NullSizeTrust SizeTrust = iota
	UntrustedSizeValue
	TrustedSizeValue
)

type FS struct {
	files_sdk.Config
	context.Context
	Root       string
	cache      *sync.Map
	cacheDir   *sync.Map
	useCache   bool
	cacheMutex *lib.KeyedMutex
}

func (f *FS) Init(config files_sdk.Config, cache bool) *FS {
	f.Config = config
	f.ClearCache()
	f.useCache = cache
	return f
}

func (f *FS) WithContext(ctx context.Context) interface{} {
	return &FS{Context: ctx, Config: f.Config, Root: f.Root, cache: f.cache, useCache: f.useCache, cacheDir: f.cacheDir, cacheMutex: f.cacheMutex}
}

func (f *FS) ClearCache() {
	f.cache = &sync.Map{}
	f.cacheDir = &sync.Map{}
	m := lib.NewKeyedMutex()
	f.cacheMutex = &m
}

type File struct {
	*files_sdk.File
	*FS
	ReadCloser        lib.AtomicValue[io.ReadCloser]
	downloadRequestId string
	MaxConnections    int
	stat              bool
	fileMutex         *sync.Mutex
	readDirState      *readDirState
	SizeTrust
	serverBytesSent int64
	representation  *downloadRepresentation
}

type ReadDirFile struct {
	*File
	state *readDirState
}

type readDirState struct {
	mutex    sync.Mutex
	iterator *folder.Iter
	done     bool
	err      error
}

func (f *File) safeFile() files_sdk.File {
	f.fileMutex.Lock()
	defer f.fileMutex.Unlock()
	return *f.File
}

func (f *File) Init() *File {
	f.fileMutex = &sync.Mutex{}
	f.readDirState = &readDirState{}
	f.SizeTrust = NullSizeTrust
	f.representation = &downloadRepresentation{}
	return f
}

func (f *File) Name() string {
	return f.safeFile().DisplayName
}

func (f *File) IsDir() bool {
	return f.safeFile().Type == "directory"
}

func (f *File) Type() goFs.FileMode {
	if f.IsDir() {
		return goFs.ModeDir
	}
	return 0
}

func (f *File) Info() (goFs.FileInfo, error) {
	return f.Stat()
}

type Info struct {
	files_sdk.File
	sizeTrust SizeTrust
}

func (i Info) Name() string {
	return i.File.DisplayName
}

func (i Info) Size() int64 {
	return i.File.Size
}

type UntrustedSize interface {
	UntrustedSize() bool
	SizeTrust() SizeTrust
	goFs.FileInfo
}

func (i Info) UntrustedSize() bool {
	return i.sizeTrust == UntrustedSizeValue || i.sizeTrust == NullSizeTrust
}

func (i Info) SizeTrust() SizeTrust {
	return i.sizeTrust
}

type PossibleSize interface {
	PossibleSize() int64
}

func (i Info) PossibleSize() int64 {
	return i.File.Size
}

func (i Info) Mode() goFs.FileMode {
	mode := goFs.ModePerm
	if i.IsDir() {
		mode |= goFs.ModeDir
	}
	return mode
}

func (i Info) IsDir() bool {
	return i.File.Type == "directory"
}

func (i Info) Sys() interface{} {
	return i.File
}

func (i Info) RemoteMount() bool {
	if i.Crc32 != "" { // Detect if is Files.com native file.
		return false
	}

	return true
}

func (f *File) Stat() (goFs.FileInfo, error) {
	f.fileMutex.Lock()
	defer f.fileMutex.Unlock()
	return Info{File: *f.File, sizeTrust: f.SizeTrust}, nil
}

func (f *File) Read(b []byte) (n int, err error) {
	f.ReadCloser.CompareAndUpdate(nil, func() io.ReadCloser {
		f.fileMutex.Lock()
		defer f.fileMutex.Unlock()

		var readCloser io.ReadCloser
		ctx, budget := withDownloadBudget(f.Context)
		readCloser, err = f.readCloserInitContext(ctx)
		if downloadRequestExpired(err) && budget.renew() {
			f.Config.LogPath(f.File.Path, map[string]interface{}{"message": "downloadRequestExpired", "error": err})
			f.File.DownloadUri = "" // force a new query
			*f.File, err = (&Client{Config: f.Config}).DownloadUri(files_sdk.FileDownloadParams{File: *f.File}, files_sdk.WithContext(ctx))
			if err == nil {
				readCloser, err = f.readCloserInitContext(ctx)
			}
		}

		if err != nil {
			readCloser = nil
		}
		if err != nil && !downloadSourceChanged(err) && f.File.DownloadUri != "" && f.downloadRequestId != "" {
			// The download request status explains most failures; after a source
			// change the typed error is the contract and the old URL's status is
			// no longer authoritative.
			status, statusErr := (&Client{Config: f.Config}).DownloadRequestStatus(f.File.DownloadUri, f.downloadRequestId, files_sdk.WithContext(f.Context))
			if statusErr != nil {
				err = statusErr
			}
			if !status.IsNil() {
				err = status
			}
		}
		return readCloser
	})

	if err != nil {
		return
	}

	return f.ReadCloser.Load().Read(b)
}

func parseSize(response *http.Response) (size int64, sizeTrust SizeTrust) {
	if response.StatusCode == http.StatusPartialContent {
		if contentRange, ok := parseContentRange(response.Header.Get("Content-Range")); ok && contentRange.totalKnown {
			return contentRange.total, TrustedSizeValue
		}
	} else if response.ContentLength > -1 {
		sizeTrust = TrustedSizeValue
		size = response.ContentLength

		return
	}

	// The response states no size, like a chunked stream from a source that
	// cannot know its length. Its size is -1, as net/http reports an unknown
	// ContentLength, so it stays distinct from a stated empty size. A full
	// response is complete when its body ends without an error. A download
	// request status, when the response advertises one, is checked on Close
	// and can fail the download or confirm its size; it is not required.
	return -1, UntrustedSizeValue
}

func parseMaxConnections(response *http.Response) int {
	maxConnections, _ := strconv.Atoi(response.Header.Get("X-Files-Max-Connections"))
	return maxConnections
}

func (f *File) readCloserInit() (readCloser io.ReadCloser, err error) {
	return f.readCloserInitContext(f.Context)
}

func (f *File) readCloserInitContext(ctx context.Context) (readCloser io.ReadCloser, err error) {
	var size int64
	var sizeTrust SizeTrust
	var fileInfo files_sdk.File
	fileInfo, err = (&Client{Config: f.Config}).Download(
		files_sdk.FileDownloadParams{File: *f.File},
		files_sdk.WithContext(ctx),
		files_sdk.ResponseOption(func(response *http.Response) error {
			f.MaxConnections = parseMaxConnections(response)
			f.downloadRequestId = response.Header.Get("X-Files-Download-Request-Id")
			if err := lib.ResponseErrors(response, files_sdk.APIError(), lib.NotStatus(http.StatusOK)); err != nil {
				return &goFs.PathError{Path: f.File.Path, Err: err, Op: "read"}
			}
			size, sizeTrust = parseSize(response)
			etag := strongDownloadETag(response.Header.Get("ETag"))
			f.representation.mu.Lock()
			bound, total := f.representation.etag, f.representation.total
			if bound != "" && (bound != etag || total != size || sizeTrust != TrustedSizeValue) {
				f.representation.mu.Unlock()
				return downloadVersionChanged("whole response representation changed")
			}
			if etag != "" && sizeTrust == TrustedSizeValue {
				f.representation.etag, f.representation.total, f.representation.whole = etag, size, false
			} else {
				f.representation.whole = true
			}
			f.representation.mu.Unlock()
			readCloser = &ReadWrapper{ReadCloser: response.Body}
			return nil
		}),
	)
	// Download may refresh the download URL from the API, whose metadata can
	// report a stale size. The size the transfer response reports is the one the
	// downloaded bytes are validated against, so it must win.
	*f.File = fileInfo
	if readCloser != nil {
		f.Size, f.SizeTrust = size, sizeTrust
	}
	return readCloser, err
}

type ReaderRange interface {
	ReaderRange(off int64, end int64) (io.ReadCloser, error)
	goFs.File
}

type ReadAtLeastWrapper struct {
	io.ReadCloser
	io.Reader
}

func (r ReadAtLeastWrapper) Close() error {
	return r.ReadCloser.Close()
}

func (f ReadAtLeastWrapper) Read(b []byte) (n int, err error) {
	return f.Reader.Read(b)
}

func (f *File) ReaderRange(off int64, end int64) (r io.ReadCloser, err error) {
	ctx, _ := withDownloadBudget(f.Context)
	attempts := f.Config.Init().RetryMax + 1
	if attempts < 1 {
		attempts = 1
	}
	r, _, err = f.openRange(ctx, off, end, attempts, false)
	if err != nil {
		return nil, &goFs.PathError{Path: f.safeFile().Path, Op: "ReaderRange", Err: err}
	}
	return r, nil
}

type ReadWrapper struct {
	io.ReadCloser
	read int
}

func (r *ReadWrapper) Read(p []byte) (n int, err error) {
	n, err = r.ReadCloser.Read(p)
	r.read += n
	return
}

type ReaderCloserDownloadStatus struct {
	*ReadWrapper
	file         *File
	expectedSize int64
	rangeRequest bool
	UntrustedSizeRangeRequestSize
}

type UntrustedSizeRangeRequestSize struct {
	ExpectedSize int64
	SentSize     int64
	ReceivedSize int64
	Status       string
}

func (u UntrustedSizeRangeRequestSize) VerifyReceived() error {
	if u.Status == "started" { // Race condition where server does not record download status. Trust what we asked for and got is correct.
		if u.ReceivedSize != u.ExpectedSize {
			errors.Join(UntrustedSizeRangeRequestSizeExpectedReceived, fmt.Errorf("expected %v bytes %v received", u.ExpectedSize, u.ReceivedSize))
		}
	} else if u.ReceivedSize != u.SentSize {
		return errors.Join(UntrustedSizeRangeRequestSizeSentReceived, fmt.Errorf("expected %v bytes sent %v received", u.SentSize, u.ReceivedSize))
	}
	return nil
}

func (u UntrustedSizeRangeRequestSize) Log() map[string]interface{} {
	return map[string]interface{}{
		"expected_size":  u.ExpectedSize,
		"sent_size":      u.SentSize,
		"received_size":  u.ReceivedSize,
		"VerifyReceived": u.VerifyReceived(),
		"Mismatch":       u.Mismatch(),
		"Status":         u.Status,
	}
}

var UntrustedSizeRangeRequestSizeExpectedReceived = fmt.Errorf("received size did not match server expected size")
var UntrustedSizeRangeRequestSizeSentReceived = fmt.Errorf("received size did not match server send size")

func (u UntrustedSizeRangeRequestSize) Mismatch() error {
	if u.Status == "started" {
		return nil
	}
	if u.ExpectedSize > u.SentSize {
		return UntrustedSizeRangeRequestSizeSentLessThanExpected
	}
	if u.ExpectedSize < u.SentSize {
		return UntrustedSizeRangeRequestSizeSentMoreThanExpected
	}
	return nil
}

var UntrustedSizeRangeRequestSizeSentMoreThanExpected = fmt.Errorf("server send more than expected")

var UntrustedSizeRangeRequestSizeSentLessThanExpected = fmt.Errorf("server send less than expected")

func (r ReaderCloserDownloadStatus) Close() error {
	if r.ReadCloser == nil {
		return nil
	}
	err := r.ReadCloser.Close()
	defer func() { r.ReadCloser = nil }()
	if err != nil {
		return err
	}

	if r.file.downloadRequestId == "" {
		return nil
	}

	info, err := r.file.Info()
	if err != nil {
		return err
	}

	if untrustedInfo, ok := info.(UntrustedSize); ok && (untrustedInfo.UntrustedSize() || untrustedInfo.SizeTrust() == NullSizeTrust) {
		r.file.fileMutex.Lock()

		status, err := (&Client{Config: r.file.Config}).DownloadRequestStatus(r.file.DownloadUri, r.file.downloadRequestId, files_sdk.WithContext(r.file.Context))
		r.file.fileMutex.Unlock()
		if err != nil {
			return err
		}
		if !status.IsNil() && (status.Data.Status == "failed" || status.Data.Status == "error") {
			return status
		}
		r.UntrustedSizeRangeRequestSize = UntrustedSizeRangeRequestSize{
			r.expectedSize,
			status.Data.BytesTransferred,
			int64(r.ReadWrapper.read),
			status.Data.Status,
		}

		if err := r.UntrustedSizeRangeRequestSize.VerifyReceived(); err != nil {
			r.file.Config.LogPath(info.Name(), r.UntrustedSizeRangeRequestSize.Log())
			return err
		}

		// The true size can only be known after the server determines that the full file has been sent without any errors.
		if r.rangeRequest {
			if err := r.UntrustedSizeRangeRequestSize.Mismatch(); err != nil {
				r.file.Config.LogPath(info.Name(), r.UntrustedSizeRangeRequestSize.Log())
				return err
			}

			if r.file.SizeTrust == UntrustedSizeValue {
				r.file.serverBytesSent += status.Data.BytesTransferred
			}
		} else if status.Data.Status != "started" {
			r.file.SizeTrust = TrustedSizeValue
			r.file.Size = status.Data.BytesTransferred
		}

		if dataBytes, err := json.Marshal(status.Data); err == nil {
			dataMap := make(map[string]interface{})
			if err = json.Unmarshal(dataBytes, &dataMap); err == nil {
				r.file.Config.LogPath(info.Name(), lo.Assign(dataMap, map[string]interface{}{"message": "download request server status"}))
			}
		}
	}
	return nil
}

func (f *File) ReadAt(p []byte, off int64) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	reader, err := f.ReaderRange(off, off+int64(len(p))-1)
	if err != nil {
		return 0, err
	}
	n, err = io.ReadFull(reader, p)
	return n, errors.Join(err, reader.Close())
}

func downloadRequestExpired(err error) bool {
	if err == nil {
		return false
	}
	if files_sdk.IsExpired(err) {
		return true
	}
	var directErr *files_sdk.DirectTransferResponseError
	if errors.As(err, &directErr) && directErr.StatusCode == http.StatusForbidden {
		return true
	}
	var responseErr lib.ResponseError
	ok := errors.As(err, &responseErr)
	return ok && responseErr.StatusCode == http.StatusForbidden
}

func (f *File) downloadURI() (err error) {
	f.fileMutex.Lock()
	*f.File, err = (&Client{Config: f.Config}).DownloadUri(files_sdk.FileDownloadParams{File: *f.File}, files_sdk.WithContext(f.Context))
	f.fileMutex.Unlock()
	return
}

func (f *File) Close() error {
	return f.closeRead(true)
}

// closeRead closes the current reader and retains legacy read-after-Close behavior.
// Completion checks apply to legacy Close and EOF-confirmed native reads.
func (f *File) closeRead(verifyCompletion bool) error {
	f.fileMutex.Lock()
	f.fileMutex.Unlock()
	defer func() { f.ReadCloser.Store(nil) }()
	closer := f.ReadCloser.Load()
	if closer == nil {
		return nil
	}
	if !verifyCompletion {
		return closer.Close()
	}
	reader, ok := closer.(*ReadWrapper)
	if !ok {
		reader = &ReadWrapper{ReadCloser: closer}
	}
	return ReaderCloserDownloadStatus{ReadWrapper: reader, file: f}.Close()
}

// restartDownload discards the current download lifecycle: the URL, its
// request status and any size learned from its responses. The next read obtains
// a fresh download request, which is required after download_source_changed
// because every range of the old URL is rejected.
func (f *File) restartDownload() {
	f.fileMutex.Lock()
	defer f.fileMutex.Unlock()
	f.File.DownloadUri = ""
	f.downloadRequestId = ""
	f.serverBytesSent = 0
	f.SizeTrust = NullSizeTrust
	f.ReadCloser.Store(nil)
	// Context copies can still own the old representation. A retry starts a
	// new generation instead of changing the proof beneath those readers.
	f.representation = &downloadRepresentation{}
}

func (f *File) WithContext(ctx context.Context) goFs.File {
	f.fileMutex.Lock()
	newF := *f
	f.fileMutex.Unlock()
	fs := *newF.FS
	newF.FS = fs.WithContext(ctx).(*FS)
	return &newF
}

func (f *FS) Open(name string) (goFs.File, error) {
	if name == "." {
		name = ""
	}
	result, ok := f.cache.Load(lib.NormalizeForComparison(name))
	if ok {
		file := result.(*File)
		if file.IsDir() {
			return &ReadDirFile{File: file, state: &readDirState{}}, nil
		}
		return file, nil
	}
	path := lib.UrlJoinNoEscape(f.Root, name)
	file, err := f.openRemoteFile(path)
	if err != nil {
		return &File{}, &goFs.PathError{Path: path, Err: err, Op: "open"}
	}
	if f.useCache {
		f.cache.Store(lib.NormalizeForComparison(path), file)
	}
	if file.IsDir() {
		return &ReadDirFile{File: file, state: &readDirState{}}, nil
	} else {
		return file, nil
	}
}

// openRemoteFile initializes a fresh handle for an already resolved API path.
// Callers own Root resolution, caching and directory wrapping.
func (f *FS) openRemoteFile(path string) (*File, error) {
	fileInfo := files_sdk.File{Type: "directory"}
	if path != "" {
		var err error
		fileInfo, err = (&Client{Config: f.Config}).Find(files_sdk.FileFindParams{Path: path}, files_sdk.WithContext(f.Context))
		if err != nil {
			return nil, err
		}
	}
	return (&File{File: &fileInfo, FS: f}).Init(), nil
}

type DirEntryError struct {
	DirEntries []goFs.DirEntry
	error
}

func (f *FS) ReadDir(name string) ([]goFs.DirEntry, error) {
	if name == "." {
		name = ""
	}
	cacheName := lib.NormalizeForComparison(name)
	if f.useCache {
		f.cacheMutex.Lock(cacheName)
		defer f.cacheMutex.Unlock(cacheName)

		dirs, ok := f.cacheDir.Load(cacheName)
		if ok {
			dirEntryError := dirs.(DirEntryError)
			return dirEntryError.DirEntries, dirEntryError.error
		}
	}

	dirs, err := (&ReadDirFile{File: (&File{File: &files_sdk.File{Path: name}, FS: f}).Init()}).ReadDir(0)
	if f.useCache && (err == nil || errors.Is(err, files_sdk.ResponseError{})) {
		f.cacheDir.Store(cacheName, DirEntryError{dirs, err})
	}
	return dirs, err
}

func (f ReadDirFile) ReadDir(n int) ([]goFs.DirEntry, error) {
	state := f.state
	if state == nil {
		state = f.File.readDirState
	}
	state.mutex.Lock()
	defer state.mutex.Unlock()
	files := make([]goFs.DirEntry, 0)
	if state.done {
		if state.err != nil {
			return files, state.err
		}
		if n > 0 {
			return files, io.EOF
		}
		return files, nil
	}
	if f.Context != nil && f.Context.Err() != nil {
		return files, &goFs.PathError{Path: f.Path, Err: f.Context.Err(), Op: "readdir"}
	}
	if state.iterator == nil {
		var err error
		folderClient := folder.Client{Config: f.Config}
		state.iterator, err = folderClient.ListFor(files_sdk.FolderListForParams{Path: f.Path}, files_sdk.WithContext(f.Context))
		if err != nil {
			state.done = true
			state.err = &goFs.PathError{Path: f.Path, Err: err, Op: "readdir"}
			return files, state.err
		}
	}
	for n <= 0 || len(files) < n {
		if !state.iterator.Next() {
			state.done = true
			if err := state.iterator.Err(); err != nil {
				state.err = &goFs.PathError{Path: f.Path, Err: err, Op: "readdir"}
				return files, state.err
			}
			if n > 0 && len(files) == 0 {
				return files, io.EOF
			}
			return files, nil
		}
		fi := state.iterator.File()
		child, err := listingEntryIsChild(f.Path, fi)
		if err != nil {
			state.done = true
			state.err = &goFs.PathError{Path: f.Path, Err: err, Op: "readdir"}
			return files, state.err
		}
		if child {
			file := (&File{File: &fi, FS: f.FS}).Init()
			if f.useCache {
				f.cache.Store(lib.NormalizeForComparison(fi.Path), file)
			}
			files = append(files, file)
		}
	}
	return files, nil
}

// listingEntryIsChild classifies a listing entry by its server path relative
// to the listed folder. A direct child is listed. A deeper descendant is
// skipped: the API can return nested files, and they are listed with their own
// folder. A path outside the folder or with malformed elements, or a display
// name that is not a single element, is an error so the bad response is
// reported instead of being dropped or followed.
func listingEntryIsChild(folder string, fi files_sdk.File) (bool, error) {
	folderSegments := remoteSegments(folder)
	segments := remoteSegments(fi.Path)
	for _, segment := range segments {
		if !validRemoteElement(segment) {
			return false, fmt.Errorf("listing entry %q has a malformed path", fi.Path)
		}
	}
	if len(segments) <= len(folderSegments) {
		return false, fmt.Errorf("listing entry %q is not below the folder %q", fi.Path, folder)
	}
	for i, segment := range folderSegments {
		if lib.NormalizeForComparison(segment) != lib.NormalizeForComparison(segments[i]) {
			return false, fmt.Errorf("listing entry %q is not below the folder %q", fi.Path, folder)
		}
	}
	if len(segments) > len(folderSegments)+1 {
		return false, nil
	}
	// fs.WalkDir joins the display name onto the folder path to reach the
	// entry, so a name with separators or parent references would make the
	// walk list a folder other than the entry's own.
	if !validRemoteElement(fi.DisplayName) || strings.Contains(fi.DisplayName, "/") {
		return false, fmt.Errorf("listing entry %q has an invalid name %q", fi.Path, fi.DisplayName)
	}
	return true, nil
}

// validRemoteElement reports whether a slash-separated remote path element can
// name a file or folder.
func validRemoteElement(element string) bool {
	return element != "" && element != "." && element != ".." && !strings.Contains(element, "\x00")
}

func (f *FS) MkdirAll(dir string, _ goFs.FileMode) error {
	var parentPath string
	for _, dirPath := range strings.Split(dir, "/") {
		if dirPath == "" {
			break
		}
		folderClient := folder.Client{Config: f.Config}
		_, err := folderClient.Create(files_sdk.FolderCreateParams{Path: lib.UrlJoinNoEscape(parentPath, dirPath)}, files_sdk.WithContext(f.Context))
		if err != nil && !errors.Is(err, files_sdk.ErrDestinationExists) {
			return err
		}

		parentPath = lib.UrlJoinNoEscape(parentPath, dirPath)
	}
	return nil
}

func (f *FS) PathSeparator() string {
	return "/"
}

// MkdirTemp is unsupported because the API cannot guarantee an exclusively new remote directory.
// It returns errors.ErrUnsupported in a PathError without creating anything or sending a request.
func (f *FS) MkdirTemp(dir, pattern string) (string, error) {
	return "", &goFs.PathError{
		Op:   "mkdirtemp",
		Path: f.PathJoin(dir, pattern),
		Err:  fmt.Errorf("atomic exclusive remote directory allocation is unavailable: %w", errors.ErrUnsupported),
	}
}

// WritableFile holds everything written to it in memory, and Close uploads it to path.
//
// Deprecated: FS.Create no longer returns a WritableFile; the writer it returns keeps what is written in a temporary
// file on disk instead. WritableFile is unchanged for code that still uses it.
type WritableFile struct {
	*Client
	*FS
	path string
	*bytes.Buffer
}

func (w WritableFile) init() WritableFile {
	w.Buffer = bytes.NewBuffer([]byte{})
	return w
}

func (w WritableFile) Write(p []byte) (int, error) {
	return w.Buffer.Write(p)
}

func (w WritableFile) Close() (err error) {
	return w.Client.Upload(
		UploadWithContext(w.Context),
		UploadWithReader(bytes.NewReader(w.Buffer.Bytes())),
		UploadWithDestinationPath(w.path),
		UploadWithSize(int64(w.Buffer.Len())),
	)
}

// Create returns a writer for the file at path under Root, where Open finds it. The writer is an io.WriteCloser, no
// longer a WritableFile. It keeps what is written in a temporary file in os.TempDir, using disk space rather than memory,
// and sends nothing until Close uploads that file, as Client.Upload does, and removes it. A nil Close means the
// upload, its finalization and removing the temporary file all succeeded, as the client saw them. An error does
// not mean nothing was written: the remote file may already have been created or replaced, as when only removing
// the temporary file failed, or the outcome may be unknown, as when the request that finalizes the upload gets no
// answer. If a write fails, Close uploads nothing and returns that error; Close uploads at most once, and a later
// Close returns the same result. The writer uses the FS's context, or context.Background when it has none: if that
// context has ended before Close, Close uploads nothing. Always call Close, even after a write or copy fails or the
// context has ended: only Close removes the temporary file. Create fails if it cannot create the temporary file.
func (f *FS) Create(path string) (io.WriteCloser, error) {
	path = lib.UrlJoinNoEscape(f.Root, path)
	staging, err := os.CreateTemp("", "files-sdk-create-*")
	if err != nil {
		return nil, &goFs.PathError{Op: "create", Path: path, Err: err}
	}
	ctx := f.Context
	if ctx == nil {
		ctx = context.Background()
	}
	return &stagedRemoteFile{fs: f, client: &Client{Config: f.Config}, ctx: ctx, path: path, staging: staging}, nil
}

// stagedRemoteFile is the writer Create returns. Only it opens its temporary file, which holds what was written. A write
// that fails, or writes only part of what it was given, fails the file: Close then uploads nothing.
type stagedRemoteFile struct {
	fs      *FS
	client  *Client
	ctx     context.Context
	path    string
	staging *os.File

	mutex    sync.Mutex
	size     int64
	writeErr error
	closed   bool

	closeOnce sync.Once
	closeErr  error
}

// Write adds p to the temporary file. After Close, or once a write has failed, it adds nothing and returns an error:
// fs.ErrClosed, or the failed write's.
func (s *stagedRemoteFile) Write(p []byte) (int, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.closed {
		return 0, &goFs.PathError{Op: "write", Path: s.path, Err: goFs.ErrClosed}
	}
	if s.writeErr != nil {
		return 0, s.writeErr
	}
	n, err := s.staging.Write(p)
	s.size += int64(n)
	if err != nil {
		s.writeErr = err
	}
	return n, err
}

// Close uploads the temporary file, unless a write failed or the context has ended, then removes it. Ending the context
// during the upload stops the parts still being sent; once every part has been sent, the upload is completed regardless.
// After an upload is attempted, the files and listings cached by the FS, which every FS WithContext made from it shares,
// are dropped, so the next Open or ReadDir asks the API again; reads running alongside Close, and other clients' writes,
// can still leave them stale. Close runs once, and a later Close returns its error without uploading again. The error is
// the failed write's, the context's or the upload's, joined with any failure to close or remove the temporary file.
func (s *stagedRemoteFile) Close() error {
	s.closeOnce.Do(func() { s.closeErr = s.upload() })
	return s.closeErr
}

func (s *stagedRemoteFile) upload() error {
	s.mutex.Lock()
	s.closed = true
	err, size := s.writeErr, s.size
	s.mutex.Unlock()

	if err == nil {
		err = s.ctx.Err()
	}
	if err == nil {
		err = s.client.Upload(
			UploadWithContext(s.ctx),
			UploadWithReaderAt(s.staging),
			UploadWithDestinationPath(s.path),
			UploadWithSize(size),
		)
		s.fs.dropCaches()
	}

	// net/http can still be finishing a read of a part's body when that part's request has ended (see
	// ProxySectionReader). Where an open file cannot be removed, as on Windows, such a read can make the removal fail;
	// Close reports that failure.
	if cleanupErr := errors.Join(s.staging.Close(), os.Remove(s.staging.Name())); cleanupErr != nil {
		return errors.Join(err, cleanupErr)
	}
	return err
}

// dropCaches empties, in place, the caches this FS shares with every FS WithContext made from it, so none of them finds
// a file or listing it cached before.
func (f *FS) dropCaches() {
	if f.cache != nil {
		f.cache.Clear()
	}
	if f.cacheDir != nil {
		f.cacheDir.Clear()
	}
}

func (f *FS) RemoveAll(path string) error {
	return (&Client{Config: f.Config}).Delete(files_sdk.FileDeleteParams{Path: path, Recursive: lib.Bool(true)}, files_sdk.WithContext(f.Context))
}

func (f *FS) Remove(path string) error {
	return (&Client{Config: f.Config}).Delete(files_sdk.FileDeleteParams{Path: path}, files_sdk.WithContext(f.Context))
}

func (f *FS) PathJoin(s ...string) string {
	return lib.UrlJoinNoEscape(s...)
}

func (f *FS) RelPath(parent, child string) (string, error) {
	path := strings.Replace(child, parent, "", 1)
	if path == "" {
		return ".", nil
	}
	path = strings.TrimSuffix(path, f.PathSeparator())
	path = strings.TrimPrefix(path, f.PathSeparator())
	return path, nil
}

func (f *FS) SplitPath(path string) (string, string) {
	if path == "" {
		return "", ""
	}

	parts := strings.Split(path, f.PathSeparator())

	return f.PathJoin(parts[:int(math.Min(float64(len(parts)-2), float64(len(parts))))]...), parts[len(parts)-1]
}

func (f *FS) TempDir() string {
	return "tmp"
}

func (f *FS) Chtimes(name string, _atime time.Time, mtime time.Time) error {
	_, err := (&Client{Config: f.Config}).Update(files_sdk.FileUpdateParams{Path: name, ProvidedMtime: &mtime}, files_sdk.WithContext(f.Context))
	if err != nil {
		return &goFs.PathError{Path: name, Err: err, Op: "chtimes"}
	}
	return nil
}
