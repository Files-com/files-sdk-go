package file

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/Files-com/files-sdk-go/v2/file/manager"
	"github.com/Files-com/files-sdk-go/v2/file/status"
	"github.com/Files-com/files-sdk-go/v2/lib"
	"github.com/zenthangplus/goccm"
)

type UploadIOParams struct {
	Path     string
	Reader   io.ReaderAt
	Size     int64
	Progress func(int64)
	Manager  goccm.ConcurrencyManager
	Parts
	files_sdk.FileUploadPart
}

func (c *Client) UploadIO(parentCtx context.Context, params UploadIOParams) (files_sdk.File, files_sdk.FileUploadPart, Parts, error) {
	var workingParts Parts
	var allParts Parts

	if params.Manager == nil {
		params.Manager = manager.Sync().FilePartsManager
	}
	if params.Progress == nil {
		params.Progress = func(i int64) {}
	}
	onComplete := make(chan *Part)
	onError := make(chan error)
	bytesWritten := int64(0)
	etags := make([]files_sdk.EtagsParam, 0)
	fileUploadPart := params.FileUploadPart
	var expires time.Time
	var err error
	if params.FileUploadPart.Expires != "" {
		expires, _ = time.Parse(time.RFC3339, params.FileUploadPart.Expires)
	}
	if !time.Now().Before(expires) || !lib.UnWrapBool(params.FileUploadPart.ParallelParts) {
		params.Parts = Parts{} // parts are invalidated
	}
	if expires.IsZero() || !time.Now().Before(expires) {
		fileUploadPart, err = c.startUpload(
			parentCtx,
			files_sdk.FileBeginUploadParams{
				Size:         params.Size,
				Path:         params.Path,
				MkdirParents: lib.Bool(true),
			},
		)
	}
	if err != nil {
		return files_sdk.File{}, fileUploadPart, workingParts, err
	}
	fileUploadPart.Path = params.Path
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()
	run := func(part *Part, fileUploadPart files_sdk.FileUploadPart, lastPart bool) {
		proxyReader := &ProxyReader{
			ReaderAt: params.Reader,
			off:      part.off,
			len:      part.len,
			onRead:   params.Progress,
		}
		fileUploadPart.PartNumber = part.number
		part.EtagsParam, part.bytes, part.error = c.createPart(ctx, proxyReader, part.len, fileUploadPart, lastPart)
		part.Touch()
		if part.error != nil {
			if *fileUploadPart.ParallelParts {
				params.Manager.Done()
			}
			onError <- part.error
			return
		}
		if *fileUploadPart.ParallelParts {
			params.Manager.Done()
		}
		onComplete <- part
	}
	if len(params.Parts) == 0 {
		for i, offset := range byteOffsetSlice(params.Size) {
			part := &Part{OffSet: offset, number: int64(i) + 1}
			workingParts = append(workingParts, part)
		}
		allParts = workingParts
	} else {
		for _, part := range params.Parts {
			if part.Successful() {
				bytesWritten = +part.bytes
				params.Progress(part.bytes)
				etags = append(etags, part.EtagsParam)
			} else {
				part.Clear()
				workingParts = append(workingParts, part)
			}
			allParts = append(allParts, part)
		}
	}

	go func() {
		for i, part := range workingParts {
			if *fileUploadPart.ParallelParts {
				params.Manager.Wait()
				go run(part, fileUploadPart, len(workingParts) == i+1)
			} else {
				run(part, fileUploadPart, len(workingParts) == i+1)
			}
		}
	}()

	for range workingParts {
		select {
		case <-ctx.Done():
			return files_sdk.File{}, fileUploadPart, allParts, ctx.Err()
		case err = <-onError:
			cancel()
			return files_sdk.File{}, fileUploadPart, allParts, err
		case part := <-onComplete:
			etags = append(etags, part.EtagsParam)
			bytesWritten += part.bytes
		}
	}

	f, err := c.completeUpload(ctx, etags, bytesWritten, fileUploadPart.Path, fileUploadPart.Ref)
	return f, fileUploadPart, allParts, err
}

func (c *Client) startUpload(ctx context.Context, beginUpload files_sdk.FileBeginUploadParams) (files_sdk.FileUploadPart, error) {
	uploads, err := c.BeginUpload(ctx, beginUpload)
	if err != nil {
		return files_sdk.FileUploadPart{}, err
	}
	return uploads[0], err
}

func (c *Client) completeUpload(ctx context.Context, etags []files_sdk.EtagsParam, bytesWritten int64, path string, ref string) (files_sdk.File, error) {
	return c.Create(ctx, files_sdk.FileCreateParams{
		ProvidedMtime: time.Now(),
		EtagsParam:    etags,
		Action:        "end",
		Path:          path,
		Ref:           ref,
		Size:          bytesWritten,
	})
}

type OffSet struct {
	off int64
	len int64
}

type Part struct {
	OffSet
	files_sdk.EtagsParam
	bytes    int64
	requests []time.Time
	error
	number int64
}

func (p *Part) Touch() {
	p.requests = append(p.requests, time.Now())
}

func (p *Part) Successful() bool {
	return p.bytes == p.len && p.error == nil
}

func (p *Part) Clear() {
	p.bytes = 0
	p.error = nil
}

type Parts []*Part

func (p Parts) SuccessfulBytes() (b int64) {
	for _, part := range p {
		if part.Successful() {
			b += part.bytes
		}
	}

	return b
}

func byteOffsetSlice(size int64) []OffSet {
	partSizes := lib.PartSizes
	var blockSize int64
	var offsets []OffSet
	off := int64(0)
	blockSize, partSizes = partSizes[0], partSizes[1:]
	endRange := blockSize
	for {
		if off < size {
			endRange = int64(math.Min(float64(endRange), float64(size)))
			offsets = append(offsets, OffSet{off: off, len: endRange - off})
			off = endRange
			endRange = off + blockSize
			blockSize, partSizes = partSizes[0], partSizes[1:]
		} else {
			break
		}
	}
	return offsets
}

func (c *Client) createPart(ctx context.Context, reader io.ReadCloser, len int64, fileUploadPart files_sdk.FileUploadPart, lastPart bool) (files_sdk.EtagsParam, int64, error) {
	partNumber := fileUploadPart.PartNumber
	var err error
	if partNumber != 1 && *fileUploadPart.ParallelParts { // Remote Mounts use the same url
		fileUploadPart, err = c.startUpload(
			ctx, files_sdk.FileBeginUploadParams{Path: fileUploadPart.Path, Ref: fileUploadPart.Ref, Part: fileUploadPart.PartNumber, MkdirParents: lib.Bool(true)},
		)
		fileUploadPart.PartNumber = partNumber
		if err != nil {
			return files_sdk.EtagsParam{}, int64(0), err
		}
	}

	headers := http.Header{}
	headers.Add("Content-Length", strconv.FormatInt(len, 10))
	res, err := files_sdk.CallRaw(
		&files_sdk.CallParams{
			Method:   fileUploadPart.HttpMethod,
			Config:   c.Config,
			Uri:      fileUploadPart.UploadUri,
			BodyIo:   reader,
			Headers:  &headers,
			Context:  ctx,
			StayOpen: !*fileUploadPart.ParallelParts && !lastPart, // Since Remote Mounts use the same url only close the connection on the last part.
		},
	)
	defer func() {
		reader.Close()
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return files_sdk.EtagsParam{}, len, err
	}

	if res.StatusCode != 200 {
		out, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return files_sdk.EtagsParam{}, len, err
		}
		return files_sdk.EtagsParam{}, len, fmt.Errorf(string(out))
	}
	etag := strings.Trim(res.Header.Get("Etag"), "\"")
	if etag == "" {
		// With remote mounts this has no value, but the code strip the value causing a validation error.
		etag = "null"
	}
	return files_sdk.EtagsParam{
		Etag: etag,
		Part: strconv.FormatInt(fileUploadPart.PartNumber, 10),
	}, len, nil
}

func uploadProgress(uploadStatus *UploadStatus) func(bytesCount int64) {
	return func(bytesCount int64) {
		uploadStatus.incrementDownloadedBytes(bytesCount)
		uploadStatus.Job().UpdateStatus(status.Uploading, uploadStatus, nil)
	}
}
