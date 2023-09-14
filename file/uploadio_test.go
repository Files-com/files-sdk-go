//go:build LIVE_TESTS
// +build LIVE_TESTS

package file

import (
	"bytes"
	"context"
	"crypto/sha256"
	"io"
	"net/http"
	"sort"
	"strconv"
	"testing"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/snabb/httpreaderat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_UploadIO(t *testing.T) {
	client := &Client{}
	assert := assert.New(t)

	req, _ := http.NewRequest("GET", "https://download.virtualbox.org/virtualbox/6.1.26/VirtualBox-6.1.26-145957-OSX.dmg", nil)

	f, err := httpreaderat.New(nil, req, nil)
	if err != nil {
		panic(err)
	}

	progressCounter := int64(0)
	progress := func(i int64) {
		progressCounter += i
	}
	var r UploadResumable
	r, err = client.UploadWithResume(
		UploadWithDestinationPath("VirtualBox.dmg"),
		UploadWithReaderAt(f),
		UploadWithSize(f.Size()),
		UploadWithManager(lib.NewConstrainedWorkGroup(10)),
		UploadWithProgress(progress),
	)
	require.NoError(t, err)
	assert.Equal(f.Size(), r.Size)
	assert.Equal(23, len(r.Parts))
	assert.Equal(progressCounter, r.Size)
	var buf bytes.Buffer
	file, err := client.Download(context.Background(), files_sdk.FileDownloadParams{Path: r.File.Path}, files_sdk.ResponseBodyOption(func(closer io.ReadCloser) error {
		_, err := io.Copy(&buf, closer)
		return err
	}))
	assert.NoError(err)
	assert.Equal(r.File.Size, int64(buf.Len()))
	assert.Equal(f.Size(), int64(buf.Len()))
	assert.Equal(f.Size(), file.Size)
	assert.Equal(r.File.Size, file.Size)
	remoteSHA := sha256.Sum256(buf.Bytes())
	response, err := client.GetHttpClient().Do(req)
	b, err := io.ReadAll(response.Body)
	localShow := sha256.Sum256(b)

	assert.Equal(localShow, remoteSHA)
}

func TestClient_UploadIO_Cancel_Restart(t *testing.T) {
	client := &Client{}
	assert := assert.New(t)

	req, _ := http.NewRequest("GET", "https://download.virtualbox.org/virtualbox/6.1.26/VirtualBox-6.1.26-145957-OSX.dmg", nil)

	f, err := httpreaderat.New(nil, req, nil)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	progressCounter := int64(0)
	progress := func(i int64) {
		progressCounter += i
	}
	var r UploadResumable
	manager := lib.NewConstrainedWorkGroup(4)
	r, err = client.UploadWithResume(
		UploadWithContext(ctx),
		UploadWithDestinationPath("VirtualBox.dmg"),
		UploadWithReaderAt(f),
		UploadWithSize(f.Size()),
		UploadWithManager(manager),
		UploadWithProgress(progress),
	)
	assert.Equal(0, manager.RunningCount())
	assert.Error(err, "context canceled")
	assert.Equal(int64(0), r.File.Size)
	assert.GreaterOrEqual(len(r.Parts), 2)
	var successful int
	var alreadyRan Parts
	for _, part := range r.Parts {
		if part.Successful() {
			alreadyRan = append(alreadyRan, part)
			successful += 1
			assert.Equal(1, len(part.requests))
		} else {
			assert.LessOrEqual(len(part.requests), 1)
		}
	}
	assert.InDelta(5, successful, 10)
	assert.Equal(r.Parts.SuccessfulBytes(), progressCounter)
	var uploadedBytes int64
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	progressCounter = 0
	r, err = client.UploadWithResume(
		UploadWithContext(ctx),
		UploadWithDestinationPath("VirtualBox.dmg"),
		UploadWithReaderAt(f),
		UploadWithSize(f.Size()),
		UploadWithManager(manager),
		UploadWithProgress(progress),
		UploadWithResume(r),
	)
	require.NoError(t, err)
	assert.Equal(r.File.Size, progressCounter)
	assert.Equal(r.File.Size, r.Parts.SuccessfulBytes())
	successful = 0
	for _, part := range r.Parts {
		for _, rpart := range alreadyRan {
			if rpart == part {
				assert.Equal(1, len(part.requests))
			}
		}
		assert.Equal(true, part.Successful())
		assert.LessOrEqual(len(part.requests), 2)
		uploadedBytes += part.bytes
		assert.NoError(part.error)
	}

	sort.SliceStable(r.Parts, func(i, j int) bool {
		x, _ := strconv.ParseInt(r.Parts[i].Part, 10, 64)
		y, _ := strconv.ParseInt(r.Parts[j].Part, 10, 64)
		return x < y
	})

	assert.Equal(r.File.Size, uploadedBytes)
	var buf bytes.Buffer
	file, err := client.Download(context.Background(), files_sdk.FileDownloadParams{Path: r.File.Path}, files_sdk.ResponseBodyOption(func(closer io.ReadCloser) error {
		_, err := io.Copy(&buf, closer)
		return err
	}))
	assert.NoError(err)
	assert.Equal(r.File.Size, int64(buf.Len()))
	assert.Equal(f.Size(), int64(buf.Len()))
	assert.Equal(f.Size(), file.Size)
	assert.Equal(r.File.Size, file.Size)
	remoteSHA := sha256.Sum256(buf.Bytes())
	response, err := client.GetHttpClient().Do(req)
	b, err := io.ReadAll(response.Body)
	localShow := sha256.Sum256(b)

	assert.Equal(localShow, remoteSHA)
}

func TestClient_UploadIO_Cancel_Restart_Expired(t *testing.T) {
	client := &Client{}
	assert := assert.New(t)

	req, _ := http.NewRequest("GET", "https://download.virtualbox.org/virtualbox/6.1.26/VirtualBox-6.1.26-145957-OSX.dmg", nil)

	f, err := httpreaderat.New(nil, req, nil)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	assert.NoError(err)
	progressCounter := int64(0)
	progress := func(i int64) {
		progressCounter += i
	}
	var r UploadResumable
	manager := lib.NewConstrainedWorkGroup(4)
	r, err = client.UploadWithResume(
		UploadWithContext(ctx),
		UploadWithDestinationPath("VirtualBox.dmg"),
		UploadWithReaderAt(f),
		UploadWithSize(f.Size()),
		UploadWithManager(manager),
		UploadWithProgress(progress),
	)

	require.Error(t, err)
	var successful int
	var alreadyRan Parts
	for _, part := range r.Parts {
		if part.Successful() {
			alreadyRan = append(alreadyRan, part)
			successful += 1
			assert.Equal(1, len(part.requests))
		} else {
			assert.LessOrEqual(len(part.requests), 1)
		}
	}
	assert.InDelta(5, successful, 10)
	assert.InDelta(r.Parts.SuccessfulBytes(), progressCounter, float64(lib.BasePart))
	var uploadedBytes int64
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	r.FileUploadPart.Expires = time.Now().AddDate(0, -1, 0).Format(time.RFC3339)
	progressCounter = 0
	r, err = client.UploadWithResume(
		UploadWithContext(ctx),
		UploadWithDestinationPath("VirtualBox.dmg"),
		UploadWithReaderAt(f),
		UploadWithSize(f.Size()),
		UploadWithManager(manager),
		UploadWithProgress(progress),
		UploadWithResume(r),
	)
	assert.NoError(err)
	assert.InDelta(r.File.Size, progressCounter, 50_000)
	successful = 0
	for _, part := range r.Parts {
		for _, rpart := range alreadyRan {
			if rpart == part {
				assert.Equal(1, len(part.requests))
			}
		}
		assert.Equal(true, part.Successful())
		assert.LessOrEqual(len(part.requests), 2)
		uploadedBytes += part.bytes
		assert.NoError(part.error)
	}

	sort.SliceStable(r.Parts, func(i, j int) bool {
		x, _ := strconv.ParseInt(r.Parts[i].Part, 10, 64)
		y, _ := strconv.ParseInt(r.Parts[j].Part, 10, 64)
		return x < y
	})

	assert.Equal(r.File.Size, uploadedBytes)
	var buf bytes.Buffer

	file, err := client.Download(context.Background(), files_sdk.FileDownloadParams{Path: r.File.Path}, files_sdk.ResponseBodyOption(func(closer io.ReadCloser) error {
		_, err := io.Copy(&buf, closer)
		return err
	}))
	assert.NoError(err)
	assert.Equal(f.Size(), int64(buf.Len()))
	assert.Equal(f.Size(), file.Size)
	remoteSHA := sha256.Sum256(buf.Bytes())
	response, err := client.GetHttpClient().Do(req)
	b, err := io.ReadAll(response.Body)
	localShow := sha256.Sum256(b)

	assert.Equal(localShow, remoteSHA)
}
