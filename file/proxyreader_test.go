package file

import (
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
