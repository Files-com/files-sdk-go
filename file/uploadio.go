package file

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/Files-com/files-sdk-go/v2/file/manager"
	"github.com/Files-com/files-sdk-go/v2/file/status"
	"github.com/Files-com/files-sdk-go/v2/lib"
)

type Progress func(int64)

type Len interface {
	Len() int
}

type UploadResumable struct {
	files_sdk.FileUploadPart
	Parts
	files_sdk.File
}

type uploadIO struct {
	ByteOffset
	Path     string
	reader   io.Reader
	readerAt io.ReaderAt
	Size     *int64
	Progress
	Manager       lib.ConcurrencyManagerWithSubWorker
	ProvidedMtime time.Time
	Parts
	files_sdk.FileUploadPart
	MkdirParents    *bool
	passedInContext context.Context

	onComplete    chan *Part
	partsFinished chan struct{}
	bytesWritten  int64
	*Client
	etags                      []files_sdk.EtagsParam
	file                       files_sdk.File
	RewindAllProgressOnFailure bool
}

func (u *uploadIO) Run(ctx context.Context) (UploadResumable, error) {
	if u.Path == "" {
		return u.UploadResumable(), errors.New("UploadWithDestinationPath is required")
	}

	if u.reader == nil && u.readerAt == nil {
		return u.UploadResumable(), errors.New("UploadWithReader or UploadWithReaderAt required")
	}

	u.onComplete = make(chan *Part)
	u.partsFinished = make(chan struct{})
	u.etags = make([]files_sdk.EtagsParam, 0)
	if u.Manager == nil {
		u.Manager = manager.Sync().FilePartsManager
	}
	if u.Progress == nil {
		u.Progress = func(i int64) {}
	}

	var defaultSize int64
	if u.Size != nil {
		defaultSize = *u.Size
	}

	var expires time.Time
	var err error
	if u.FileUploadPart.Expires != "" {
		expires, _ = time.Parse(time.RFC3339, u.FileUploadPart.Expires)
	}
	if !time.Now().Before(expires) || !lib.UnWrapBool(u.FileUploadPart.ParallelParts) {
		u.Parts = Parts{} // parts are invalidated
	}

	if u.MkdirParents == nil {
		u.MkdirParents = lib.Bool(true)
	}

	if expires.IsZero() || !time.Now().Before(expires) {
		if u.Manager.WaitWithContext(ctx) {
			u.FileUploadPart, err = u.startUpload(
				ctx,
				files_sdk.FileBeginUploadParams{
					Size:         defaultSize,
					Path:         u.Path,
					MkdirParents: u.MkdirParents,
				},
			)
			u.Manager.Done()
		} else {
			return u.UploadResumable(), ctx.Err()
		}
		if err != nil {
			return u.UploadResumable(), err
		}
	}

	partCtx, cancel := context.WithCancelCause(ctx)

	u.FileUploadPart.Path = u.Path
	go u.uploadParts(partCtx)
	return u.waitOnParts(ctx, cancel)
}

func (u *uploadIO) UploadResumable() UploadResumable {
	return UploadResumable{Parts: u.Parts, FileUploadPart: u.FileUploadPart, File: u.file}
}

func (u *uploadIO) rewindSuccessfulParts() {
	if !u.RewindAllProgressOnFailure {
		return
	}
	for _, part := range u.Parts {
		if part.Successful() {
			u.Progress(-part.bytes)
		}
	}
}

func (u *uploadIO) waitOnParts(ctx context.Context, cancelParts context.CancelCauseFunc) (UploadResumable, error) {
	var allErrors error
	for {
		select {
		case <-u.partsFinished:
			close(u.onComplete)
			if allErrors != nil {
				u.rewindSuccessfulParts()
				return u.UploadResumable(), allErrors
			}
			// Rate limit all outgoing connections
			if u.Manager.WaitWithContext(ctx) {
				var err error
				u.file, err = u.completeUpload(ctx, &u.ProvidedMtime, u.etags, u.bytesWritten, u.Path, u.FileUploadPart.Ref)
				u.Manager.Done()
				if err != nil {
					u.rewindSuccessfulParts()
				}
				return u.UploadResumable(), err
			} else {
				u.rewindSuccessfulParts()
				return u.UploadResumable(), ctx.Err()
			}
		case part := <-u.onComplete:
			if part.error == nil {
				u.etags = append(u.etags, part.EtagsParam)
				u.bytesWritten += part.bytes
			} else {
				u.Progress(-part.bytes)
				allErrors = errors.Join(allErrors, part.error)

				if strings.Contains(part.Error(), "File Upload Not Found") {
					cancelParts(part.error)
					u.etags = []files_sdk.EtagsParam{}
					u.Parts = Parts{}
					u.FileUploadPart = files_sdk.FileUploadPart{}
				}
			}
		}
	}
}

func (u *uploadIO) Reader() (io.Reader, bool) {
	return u.reader, u.reader != nil
}

func (u *uploadIO) ReaderAt() (io.ReaderAt, bool) {
	if u.readerAt != nil {
		return u.readerAt, true
	}
	readerAt, ok := u.reader.(io.ReaderAt)
	return readerAt, ok
}

func (u *uploadIO) partBuilder(offset OffSet, final bool, number int) *Part {
	part := &Part{OffSet: offset, number: number, final: final, ProxyReader: u.buildReader(offset)}
	// Parts are stored so a retry can pick up failed parts. Since io.Reader is a stream is better to just retry the whole file
	if _, readerAtOk := u.ReaderAt(); readerAtOk && u.Size != nil {
		u.Parts = append(u.Parts, part)
	}

	return part
}

func (u *uploadIO) buildReader(offset OffSet) ProxyReader {
	var readerAt io.ReaderAt
	var readerAtOk bool

	if readerAt, readerAtOk = u.ReaderAt(); !readerAtOk {
		u.ParallelParts = lib.Bool(false)
	}

	if u.Size == nil {
		if readerAtOk {
			sectionReader := io.NewSectionReader(readerAt, offset.off, offset.len)
			buf := new(bytes.Buffer)

			// Use io.CopyN to copy up to fiveMB into buf
			n, err := io.CopyN(buf, sectionReader, offset.len)
			if err != nil && err != io.EOF {
				// handle error
			}

			return &ProxyRead{
				Reader: buf,
				len:    n,
				onRead: u.Progress,
			}
		} else {
			buf := new(bytes.Buffer)

			// Use io.CopyN to copy up to fiveMB into buf
			reader, _ := u.Reader()
			n, err := io.CopyN(buf, reader, offset.len)
			if err != nil && err != io.EOF {
				// handle error
			}

			return &ProxyRead{
				Reader: buf,
				len:    n,
				onRead: u.Progress,
			}
		}
	}

	if readerAtOk {
		return &ProxyReaderAt{
			ReaderAt: readerAt,
			off:      offset.off,
			len:      offset.len,
			onRead:   u.Progress,
		}
	} else {
		reader, _ := u.Reader()
		return &ProxyRead{
			Reader: reader,
			len:    offset.len,
			onRead: u.Progress,
		}
	}
}

type PartRunnerReturn int

func (u *uploadIO) manageUpdatePart(ctx context.Context, part *Part, wait lib.ConcurrencyManager) bool {
	if wait.WaitWithContext(ctx) {
		if *u.FileUploadPart.ParallelParts && u.Size != nil {
			go func() {
				u.runUploadPart(ctx, part)
				wait.Done()
			}()
		} else {
			u.runUploadPart(ctx, part)
			wait.Done()
		}
		if part.ProxyReader.Len() != int(part.len) {
			return true
		}
	} else {
		return true
	}

	return false
}

func (u *uploadIO) runUploadPart(ctx context.Context, part *Part) {
	fileUploadPart := u.FileUploadPart
	fileUploadPart.PartNumber = int64(part.number)
	part.EtagsParam, part.error = u.createPart(ctx, part.ProxyReader, int64(part.ProxyReader.Len()), fileUploadPart, part.final)
	part.bytes = int64(part.ProxyReader.BytesRead())
	part.Touch()
	if part.error != nil {
		var pathErr *os.PathError
		if errors.As(part.error, &pathErr) {
			part.error = pathErr
		}
	}

	u.onComplete <- part
}

func (u *uploadIO) uploadParts(ctx context.Context) {
	wait := u.Manager.NewSubWorker()
	defer func() {
		wait.WaitAllDone()
		u.partsFinished <- struct{}{}
	}()

	for _, part := range u.Parts {
		if part.Successful() {
			u.Progress(part.bytes)
			u.onComplete <- part
		} else {
			part.Clear()
			part.ProxyReader = u.buildReader(part.OffSet)

			if u.manageUpdatePart(ctx, part, wait) {
				break
			}
		}
	}
	var (
		offset   OffSet
		iterator Iterator
		index    int
	)
	if len(u.Parts) > 0 {
		lastPart := u.Parts[len(u.Parts)-1]
		iterator = u.ByteOffset.Resume(u.Size, lastPart.off, lastPart.number-1)
		offset, iterator, index = iterator()
	} else {
		iterator = u.ByteOffset.BySize(u.Size)
	}

	for {
		if iterator == nil {
			break
		}

		offset, iterator, index = iterator()

		if u.manageUpdatePart(
			ctx,
			u.partBuilder(
				offset,
				iterator == nil && u.Size != nil,
				index+1,
			),
			wait,
		) {
			break
		}
	}
}

func (u *uploadIO) startUpload(ctx context.Context, beginUpload files_sdk.FileBeginUploadParams) (files_sdk.FileUploadPart, error) {
	uploads, err := u.BeginUpload(beginUpload, files_sdk.WithContext(ctx))
	if err != nil {
		return files_sdk.FileUploadPart{}, err
	}
	return uploads[0], err
}

func (u *uploadIO) completeUpload(ctx context.Context, providedMtime *time.Time, etags []files_sdk.EtagsParam, bytesWritten int64, path string, ref string) (files_sdk.File, error) {
	if providedMtime.IsZero() {
		providedMtime = nil
	}

	return u.Create(files_sdk.FileCreateParams{
		ProvidedMtime: providedMtime,
		EtagsParam:    etags,
		Action:        "end",
		Path:          path,
		Ref:           ref,
		Size:          bytesWritten,
		MkdirParents:  lib.Bool(true),
	}, files_sdk.WithContext(ctx))
}

func (u *uploadIO) createPart(ctx context.Context, reader io.ReadCloser, len int64, fileUploadPart files_sdk.FileUploadPart, lastPart bool) (files_sdk.EtagsParam, error) {
	partNumber := fileUploadPart.PartNumber
	var err error
	if partNumber != 1 && *fileUploadPart.ParallelParts { // Remote Mounts use the same url
		fileUploadPart, err = u.startUpload(
			ctx, files_sdk.FileBeginUploadParams{Path: fileUploadPart.Path, Ref: fileUploadPart.Ref, Part: fileUploadPart.PartNumber, MkdirParents: lib.Bool(true)},
		)
		fileUploadPart.PartNumber = partNumber
		if err != nil {
			return files_sdk.EtagsParam{}, err
		}
	}
	uri, err := url.Parse(fileUploadPart.UploadUri)
	if err == nil {
		q := uri.Query()
		if q.Get("partNumber") == "" {
			q.Add("part_number", strconv.FormatInt(partNumber, 10))
			uri.RawQuery = q.Encode()
			fileUploadPart.UploadUri = uri.String()
		}
	}

	headers := http.Header{}
	headers.Add("Content-Length", strconv.FormatInt(len, 10))
	res, err := files_sdk.CallRaw(
		&files_sdk.CallParams{
			Method:   fileUploadPart.HttpMethod,
			Config:   u.Config,
			Uri:      fileUploadPart.UploadUri,
			BodyIo:   reader,
			Headers:  &headers,
			Context:  ctx,
			StayOpen: !*fileUploadPart.ParallelParts && !lastPart, // Since Remote Mounts use the same url only close the connection on the last part.
		},
	)
	if err != nil {
		return files_sdk.EtagsParam{}, err
	}
	if err := lib.ResponseErrors(res, files_sdk.APIError(), lib.NonOkError); err != nil {
		return files_sdk.EtagsParam{}, err
	}
	etag := strings.Trim(res.Header.Get("Etag"), "\"")
	if etag == "" {
		// With remote mounts this has no value, but the code strip the value causing a validation error.
		etag = "null"
	}
	if res != nil {
		err = res.Body.Close()
		if err != nil {
			return files_sdk.EtagsParam{}, err
		}
	}
	return files_sdk.EtagsParam{
		Etag: etag,
		Part: strconv.FormatInt(fileUploadPart.PartNumber, 10),
	}, nil
}

func uploadProgress(uploadStatus *UploadStatus) func(bytesCount int64) {
	return func(bytesCount int64) {
		uploadStatus.incrementDownloadedBytes(bytesCount)
		uploadStatus.Job().UpdateStatus(status.Uploading, uploadStatus, nil)
	}
}
