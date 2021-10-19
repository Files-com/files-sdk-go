//go:build LIVE_TESTS
// +build LIVE_TESTS

package file

import (
	"bytes"
	"context"
	"crypto/sha256"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"testing"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/Files-com/files-sdk-go/v2/lib"
	"github.com/stretchr/testify/assert"
	"github.com/zenthangplus/goccm"

	"github.com/snabb/httpreaderat"
)

func TestClient_UploadIO_Cancel_Restart(t *testing.T) {
	client := &Client{}
	assert := assert.New(t)

	req, _ := http.NewRequest("GET", "https://download.virtualbox.org/virtualbox/6.1.26/VirtualBox-6.1.26-145957-OSX.dmg", nil)

	f, err := httpreaderat.New(nil, req, nil)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	progressCounter := int64(0)
	progress := func(i int64) {
		progressCounter += i
	}
	params := UploadIOParams{
		Path:     "VirtualBox.dmg",
		Reader:   f,
		Size:     f.Size(),
		Manager:  goccm.New(4),
		Progress: progress,
	}
	var fi files_sdk.File
	var parts Parts
	var fileUploadPart files_sdk.FileUploadPart
	cancelLater := func() {
		time.Sleep(10 * time.Second)
		cancel()
	}
	go cancelLater()
	fi, fileUploadPart, parts, err = client.UploadIO(ctx, params)
	assert.Error(err, "context canceled")
	assert.Equal(int64(0), fi.Size)
	assert.Equal(23, len(parts))
	var successful int
	var alreadyRan Parts
	for _, part := range parts {
		if part.Successful() {
			alreadyRan = append(alreadyRan, part)
			successful += 1
			assert.Equal(1, len(part.requests))
		} else {
			assert.LessOrEqual(len(part.requests), 1)
		}
	}
	assert.InDelta(5, successful, 10)
	assert.InDelta(parts.SuccessfulBytes(), progressCounter, 2*1024*1024)
	params.Parts = parts
	var uploadedBytes int64
	ctx, _ = context.WithCancel(context.Background())
	params.FileUploadPart = fileUploadPart
	var newFileUploadPart files_sdk.FileUploadPart
	progressCounter = 0

	fi, newFileUploadPart, parts, err = client.UploadIO(ctx, params)
	assert.NoError(err)
	assert.Equal(fileUploadPart, newFileUploadPart)
	assert.InDelta(params.Size, progressCounter, float64(lib.BasePart))
	assert.NotEqual(params.Size, fi.Size, "Returned size will not always match")
	assert.Equal(params.Size, parts.SuccessfulBytes())
	successful = 0
	for _, part := range parts {
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

	sort.SliceStable(parts, func(i, j int) bool {
		x, _ := strconv.ParseInt(parts[i].Part, 10, 64)
		y, _ := strconv.ParseInt(parts[j].Part, 10, 64)
		return x < y
	})

	assert.Equal(params.Size, uploadedBytes)
	var buf bytes.Buffer
	file, err := client.Download(context.Background(), files_sdk.FileDownloadParams{Writer: &buf, Path: params.Path})
	assert.NoError(err)
	assert.Equal(params.Size, int64(buf.Len()))
	assert.Equal(f.Size(), int64(buf.Len()))
	assert.Equal(f.Size(), file.Size)
	assert.Equal(params.Size, file.Size)
	remoteSHA := sha256.Sum256(buf.Bytes())
	response, err := client.GetHttpClient().Do(req)
	b, err := ioutil.ReadAll(response.Body)
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

	ctx, cancel := context.WithCancel(context.Background())
	assert.NoError(err)
	progressCounter := int64(0)
	progress := func(i int64) {
		progressCounter += i
	}
	params := UploadIOParams{
		Path:     "VirtualBox.dmg",
		Reader:   f,
		Size:     f.Size(),
		Manager:  goccm.New(4),
		Progress: progress,
	}

	var parts Parts
	var fileUploadPart files_sdk.FileUploadPart
	cancelLater := func() {
		time.Sleep(10 * time.Second)
		cancel()
	}
	go cancelLater()
	_, fileUploadPart, parts, _ = client.UploadIO(ctx, params)
	assert.Equal(23, len(parts))
	var successful int
	var alreadyRan Parts
	for _, part := range parts {
		if part.Successful() {
			alreadyRan = append(alreadyRan, part)
			successful += 1
			assert.Equal(1, len(part.requests))
		} else {
			assert.LessOrEqual(len(part.requests), 1)
		}
	}
	assert.InDelta(5, successful, 10)
	assert.InDelta(parts.SuccessfulBytes(), progressCounter, float64(lib.BasePart))
	params.Parts = parts
	var uploadedBytes int64
	ctx, _ = context.WithCancel(context.Background())
	params.FileUploadPart = fileUploadPart
	params.FileUploadPart.Expires = time.Now().AddDate(0, -1, 0).Format(time.RFC3339)
	var newFileUploadPart files_sdk.FileUploadPart
	progressCounter = 0
	_, newFileUploadPart, parts, err = client.UploadIO(ctx, params)
	assert.NoError(err)
	assert.NotEqual(fileUploadPart, newFileUploadPart)
	assert.InDelta(params.Size, progressCounter, 50_000)
	successful = 0
	for _, part := range parts {
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

	sort.SliceStable(parts, func(i, j int) bool {
		x, _ := strconv.ParseInt(parts[i].Part, 10, 64)
		y, _ := strconv.ParseInt(parts[j].Part, 10, 64)
		return x < y
	})

	assert.Equal(params.Size, uploadedBytes)
	var buf bytes.Buffer

	file, err := client.Download(context.Background(), files_sdk.FileDownloadParams{Writer: &buf, Path: params.Path})
	assert.NoError(err)
	assert.Equal(f.Size(), int64(buf.Len()))
	assert.Equal(f.Size(), file.Size)
	remoteSHA := sha256.Sum256(buf.Bytes())
	response, err := client.GetHttpClient().Do(req)
	b, err := ioutil.ReadAll(response.Body)
	localShow := sha256.Sum256(b)

	assert.Equal(localShow, remoteSHA)
}
