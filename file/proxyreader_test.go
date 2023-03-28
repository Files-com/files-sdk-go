package file

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProxyReader_Read_1(t *testing.T) {
	assert := assert.New(t)

	localFile, err := os.Open("../LICENSE")
	assert.NoError(err)
	reader := &ProxyReader{
		ReaderAt: localFile,
		off:      100,
		len:      101,
		onRead:   nil,
		read:     0,
	}

	b, err := io.ReadAll(reader)
	assert.NoError(err)
	assert.Equal(101, len(b))
	assert.Equal("granted, free of charge, to any person obtaining a copy\nof this software and associated documentation", string(b))
}

func TestProxyReader_Read_2(t *testing.T) {
	assert := assert.New(t)

	localFile, err := os.Open("../LICENSE")
	assert.NoError(err)
	reader := &ProxyReader{
		ReaderAt: localFile,
		off:      10,
		len:      1000,
		onRead:   nil,
		read:     0,
	}

	b, err := io.ReadAll(reader)
	assert.NoError(err)
	assert.Equal(1000, len(b))
}

func TestProxyReader_Read_OnRead(t *testing.T) {
	assert := assert.New(t)
	var read int64
	localFile, err := os.Open("../LICENSE")
	assert.NoError(err)
	reader := &ProxyReader{
		ReaderAt: localFile,
		off:      10,
		len:      1000,
		onRead: func(i int64) {
			read += i
		},
		read: 0,
	}

	b, err := io.ReadAll(reader)
	assert.NoError(err)
	assert.Equal(1000, len(b))
	assert.Equal(int64(1000), read)
	reader.Close()
	//	Can be read again
	b, err = io.ReadAll(reader)
	assert.NoError(err)
	assert.Equal(1000, len(b))
	assert.Equal(int64(1000), read)
}

func TestProxyReader_ReadWithOnRead(t *testing.T) {
	data := []byte("This is a test string for the ProxyReader implementation.")
	readerAt := bytes.NewReader(data)

	var bytesRead int64
	onRead := func(i int64) {
		bytesRead += i
	}

	proxyReader := &ProxyReader{
		ReaderAt: readerAt,
		off:      0,
		len:      int64(len(data)),
		onRead:   onRead,
	}

	buf := make([]byte, 10)
	n, err := proxyReader.Read(buf)

	if err != nil {
		t.Errorf("Error reading from ProxyReader: %v", err)
	}

	if n != 10 {
		t.Errorf("Expected to read 10 bytes, got %d bytes", n)
	}

	if bytesRead != int64(n) {
		t.Errorf("Expected bytesRead to be %d, got %d", n, bytesRead)
	}

	if string(buf) != "This is a " {
		t.Errorf("Unexpected read data: %s", buf)
	}
}

func TestProxyReader_SeekWithOnRead(t *testing.T) {
	data := []byte("This is a test string for the ProxyReader implementation.")
	readerAt := bytes.NewReader(data)

	var bytesRead int64
	onRead := func(i int64) {
		bytesRead += i
	}

	proxyReader := &ProxyReader{
		ReaderAt: readerAt,
		off:      0,
		len:      int64(len(data)),
		onRead:   onRead,
	}

	_, err := proxyReader.Seek(4, io.SeekStart)

	if err != nil {
		t.Errorf("Error seeking ProxyReader: %v", err)
	}

	if bytesRead != 4 {
		t.Errorf("Expected bytesRead to be 4, got %d", bytesRead)
	}

	buf := make([]byte, 4)
	n, err := proxyReader.Read(buf)

	if err != nil {
		t.Errorf("Error reading from ProxyReader: %v", err)
	}

	if n != 4 {
		t.Errorf("Expected to read 4 bytes, got %d bytes", n)
	}

	if bytesRead != 8 {
		t.Errorf("Expected bytesRead to be 8, got %d", bytesRead)
	}

	if string(buf) != " is " {
		t.Errorf("Unexpected read data: %s", buf)
	}
}

func TestProxyReader_CloseWithOnRead(t *testing.T) {
	data := []byte("This is a test string for the ProxyReader implementation.")
	readerAt := bytes.NewReader(data)

	var bytesRead int64
	onRead := func(i int64) {
		bytesRead += i
	}

	proxyReader := &ProxyReader{
		ReaderAt: readerAt,
		off:      0,
		len:      int64(len(data)),
		onRead:   onRead,
	}

	buf := make([]byte, 10)
	n, err := proxyReader.Read(buf)

	if err != nil {
		t.Errorf("Error reading from ProxyReader: %v", err)
	}

	if n != 10 {
		t.Errorf("Expected to read 10 bytes, got %d bytes", n)
	}

	if bytesRead != 10 {
		t.Errorf("Expected bytesRead to be 10, got %d", bytesRead)
	}

	err = proxyReader.Close()
	if err != nil {
		t.Errorf("Error closing ProxyReader: %v", err)
	}

	// Since the onRead callback is not called during Close(), bytesRead remains the same.
	if bytesRead != 10 {
		t.Errorf("Expected bytesRead to be 10, got %d", bytesRead)
	}
}
