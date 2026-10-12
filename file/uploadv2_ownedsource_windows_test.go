package file

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"maps"
	"net/http"
	"os"
	"slices"
	"strconv"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Every part's own handle to the upload's file is closed before the upload reports how it ended, whether it succeeded,
// failed or was canceled while parts were running, and whether those parts were scheduled directly or from the ready
// runway. The upload's file here lets itself be deleted while it is open, but a part's own handle does not, as no
// handle Go opens does: deleting the source as soon as the upload returns fails if any part's handle is still open.
func TestUploadV2OwnedPartHandlesCloseBeforeTheResult(t *testing.T) {
	for _, runway := range []bool{false, true} {
		for _, outcome := range []string{"success", "failure", "cancel"} {
			t.Run(fmt.Sprintf("%s ready runway %t", outcome, runway), func(t *testing.T) {
				source, data := writeOwnedHandlesSource(t)
				file := openDeletableFile(t, source)
				server := (&MockAPIServer{T: t}).Do()
				t.Cleanup(server.Shutdown)
				serveOwnedHandlesBeginUpload(server, "source.bin", int64(len(data)))
				// The upload fails or is canceled once the first two parts are running.
				var mu sync.Mutex
				arrived := 0
				running := make(chan struct{})
				server.GetRouter().POST("/owned-handles-upload", func(c *gin.Context) {
					if _, err := io.Copy(io.Discard, c.Request.Body); err != nil {
						c.Status(http.StatusBadRequest)
						return
					}
					if outcome == "success" {
						c.Header("Etag", "etag-"+c.Query("part_number"))
						c.Status(http.StatusOK)
						return
					}
					mu.Lock()
					arrived++
					if arrived == 2 {
						close(running)
					}
					mu.Unlock()
					select {
					case <-running:
					case <-time.After(10 * time.Second):
						t.Error("two parts never ran at once")
					}
					if outcome == "failure" && c.Query("part_number") == "2" {
						c.Status(http.StatusBadRequest)
						return
					}
					<-c.Request.Context().Done()
				})
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				if outcome == "cancel" {
					go func() {
						<-running
						cancel()
					}()
				}
				var logs bytes.Buffer
				client := server.Client()
				client.Config.Logger = log.New(&logs, "", 0)
				opts := []UploadOption{
					UploadWithV2(),
					UploadWithContext(ctx),
					uploadWithOwnedFile(file),
					UploadWithSize(int64(len(data))),
					UploadWithDestinationPath("source.bin"),
					UploadWithManager(lib.NewConstrainedWorkGroup(4)),
				}
				if !runway {
					opts = append(opts, UploadWithV2ReadyRunway(0, 0))
				}

				_, err := client.UploadWithResume(opts...)

				assert.NoError(t, os.Remove(source))
				switch outcome {
				case "success":
					require.NoError(t, err)
				case "failure":
					require.ErrorContains(t, err, "status code: 400")
				case "cancel":
					require.ErrorIs(t, err, context.Canceled)
				}
				logText := logs.String()
				assert.Regexp(t, `scheduler_owned_part_handles: [1-9]`, logText)
				assert.Regexp(t, `scheduler_owned_part_handle_fallbacks: 0(,|\n)`, logText)
				if runway {
					assert.Regexp(t, `scheduler_dispatched_prepared_parts: [1-9]`, logText)
				} else {
					assert.Regexp(t, `scheduler_dispatched_prepared_parts: 0(,|\n)`, logText)
				}
			})
		}
	}
}

// A file a caller passes as the ReaderAt stays the caller's, even after the file the SDK opened for the upload: the
// upload reads it only through the caller's handle, so the caller closing it while the upload runs fails the upload
// rather than the upload reading on through handles of its own.
func TestUploadWithReaderAtReadsOnlyThroughTheCallersHandle(t *testing.T) {
	for _, test := range []struct {
		name  string
		owned bool
	}{
		{name: "caller's file"},
		{name: "caller's file after the SDK's", owned: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			source, data := writeOwnedHandlesSource(t)
			server := (&MockAPIServer{T: t}).Do()
			t.Cleanup(server.Shutdown)
			serveOwnedHandlesBeginUpload(server, "source.bin", int64(len(data)))
			caller, err := os.Open(source)
			require.NoError(t, err)
			arrived := make(chan struct{})
			var arrivedOnce sync.Once
			closed := make(chan struct{})
			go func() {
				<-arrived
				_ = caller.Close()
				close(closed)
			}()
			server.GetRouter().POST("/owned-handles-upload", func(c *gin.Context) {
				arrivedOnce.Do(func() { close(arrived) })
				<-closed
				_, _ = io.Copy(io.Discard, c.Request.Body)
				c.Header("Etag", "etag-"+c.Query("part_number"))
				c.Status(http.StatusOK)
			})
			opts := []UploadOption{
				UploadWithV2(),
				UploadWithSize(int64(len(data))),
				UploadWithDestinationPath("source.bin"),
				UploadWithManager(lib.NewConstrainedWorkGroup(4)),
			}
			if test.owned {
				owned, err := os.Open(source)
				require.NoError(t, err)
				t.Cleanup(func() { _ = owned.Close() })
				opts = append(opts, uploadWithOwnedFile(owned))
			}
			opts = append(opts, UploadWithReaderAt(caller))

			_, err = server.Client().UploadWithResume(opts...)

			require.ErrorIs(t, err, os.ErrClosed)
		})
	}
}

// A part reads through a handle of its own only when that handle is to the upload's file. Once the file's name names
// another file, as when the file is moved and another put in its place, the parts read the upload's file instead, and
// the upload sends the bytes of the file it opened.
func TestUploadV2OwnedPartHandlesReadTheUploadsFileAfterItsNameMoves(t *testing.T) {
	source, data := writeOwnedHandlesSource(t)
	file := openDeletableFile(t, source)
	require.NoError(t, os.Rename(source, source+".moved"))
	require.NoError(t, os.WriteFile(source, bytes.Repeat([]byte{0xa5}, len(data)), 0o600))
	server := (&MockAPIServer{T: t}).Do()
	t.Cleanup(server.Shutdown)
	serveOwnedHandlesBeginUpload(server, "source.bin", int64(len(data)))
	var mu sync.Mutex
	sent := map[int64][]byte{}
	server.GetRouter().POST("/owned-handles-upload", func(c *gin.Context) {
		offset, err := strconv.ParseInt(c.Query("part_offset"), 10, 64)
		body, readErr := io.ReadAll(c.Request.Body)
		if err != nil || readErr != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		mu.Lock()
		sent[offset] = body
		mu.Unlock()
		c.Header("Etag", "etag-"+c.Query("part_number"))
		c.Status(http.StatusOK)
	})
	var logs bytes.Buffer
	client := server.Client()
	client.Config.Logger = log.New(&logs, "", 0)

	_, err := client.UploadWithResume(
		UploadWithV2(),
		uploadWithOwnedFile(file),
		UploadWithSize(int64(len(data))),
		UploadWithDestinationPath("source.bin"),
		UploadWithManager(lib.NewConstrainedWorkGroup(4)),
	)

	require.NoError(t, err)
	mu.Lock()
	defer mu.Unlock()
	var uploaded []byte
	for _, offset := range slices.Sorted(maps.Keys(sent)) {
		uploaded = append(uploaded, sent[offset]...)
	}
	assert.True(t, bytes.Equal(data, uploaded), "the parts sent other bytes than the upload's file")
	assert.Regexp(t, `scheduler_owned_part_handles: 0(,|\n)`, logs.String())
	assert.Regexp(t, `scheduler_owned_part_handle_fallbacks: 3(,|\n)`, logs.String())
}

// openDeletableFile opens path for reading, as a file that does not keep it from being deleted or renamed.
func openDeletableFile(t *testing.T, path string) *os.File {
	t.Helper()
	name, err := syscall.UTF16PtrFromString(path)
	require.NoError(t, err)
	handle, err := syscall.CreateFile(name, syscall.GENERIC_READ, syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE|syscall.FILE_SHARE_DELETE, nil, syscall.OPEN_EXISTING, syscall.FILE_ATTRIBUTE_NORMAL, 0)
	require.NoError(t, err)
	file := os.NewFile(uintptr(handle), path)
	t.Cleanup(func() { _ = file.Close() })
	return file
}
