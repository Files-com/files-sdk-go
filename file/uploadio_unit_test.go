package file

import (
	"bytes"
	"errors"
	"io"
	"testing"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// readerWithUnsupportedReadAt wraps a reader but returns ErrUnsupported for ReadAt
type readerWithUnsupportedReadAt struct {
	io.Reader
}

func (r *readerWithUnsupportedReadAt) ReadAt(p []byte, off int64) (n int, err error) {
	return 0, errors.ErrUnsupported
}

func TestUploadIO_ReadAtUnsupportedFallback(t *testing.T) {
	testData := []byte("Hello, world! This is test data for ReadAt fallback functionality.")

	// Create a reader that implements ReadAt but returns ErrUnsupported
	reader := &readerWithUnsupportedReadAt{
		Reader: bytes.NewReader(testData),
	}

	uploadIO := &uploadIO{
		reader: reader,
		Size:   func() *int64 { s := int64(len(testData)); return &s }(),
		FileUploadPart: files_sdk.FileUploadPart{
			ParallelParts: func() *bool { b := false; return &b }(),
		},
	}

	// Test that ReaderAt detection works
	readerAt, ok := uploadIO.ReaderAt()
	assert.True(t, ok, "Should detect ReadAt interface")
	assert.Equal(t, reader, readerAt, "Should return the same reader")

	// Test buildReader with offset
	offset := OffSet{off: 0, len: 10}
	proxyReader, err := uploadIO.buildReader(offset)
	require.NoError(t, err, "buildReader should handle ErrUnsupported gracefully")

	// Should fall back to ProxyRead, not ProxyReaderAt
	_, isProxyReaderAt := proxyReader.(*ProxyReaderAt)
	assert.False(t, isProxyReaderAt, "Should fallback to ProxyRead when ReadAt unsupported")

	// Verify the reader works
	buf := make([]byte, 10)
	n, err := proxyReader.Read(buf)
	assert.NoError(t, err)
	assert.Equal(t, 10, n)
	assert.Equal(t, "Hello, wor", string(buf))
}

func TestUploadIO_ReadAtOtherError(t *testing.T) {
	testData := []byte("Hello, world!")

	// Create a reader that doesn't implement ReadAt (bytes.Buffer doesn't)
	reader := bytes.NewBuffer(testData)

	uploadIO := &uploadIO{
		reader: reader,
		Size:   func() *int64 { s := int64(len(testData)); return &s }(),
		FileUploadPart: files_sdk.FileUploadPart{
			ParallelParts: func() *bool { b := false; return &b }(),
		},
	}

	// This reader doesn't implement ReadAt, so ReaderAt should return false
	_, ok := uploadIO.ReaderAt()
	assert.False(t, ok, "Should not detect ReadAt interface for bytes.Buffer")
}
