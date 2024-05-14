package file

import (
	"io"
)

type ProxyReader interface {
	io.ReadCloser
	Len() int
	BytesRead() int
	Rewind() bool
}

type ProxyReaderAt struct {
	io.ReaderAt
	off    int64
	len    int64
	onRead func(i int64)
	read   int
	closed bool
	eof    bool
}

type ProxyRead struct {
	io.Reader
	len    int64
	onRead func(i int64)
	read   int
	closed bool
	eof    bool
}

func (x *ProxyReaderAt) Rewind() bool {
	x.onRead(int64(-x.read))
	x.read = 0
	return true
}

func (x *ProxyRead) Rewind() bool {
	return x.read == 0
}

func (x *ProxyReaderAt) Len() int {
	return int(x.len)
}

func (x *ProxyRead) Len() int {
	return int(x.len)
}

func (x *ProxyReaderAt) BytesRead() int {
	return x.read
}

func (x *ProxyRead) BytesRead() int {
	return x.read
}

func (x *ProxyReaderAt) Seek(offset int64, whence int) (int64, error) {
	x.onRead(-int64(x.read - int(offset))) // rewind progress
	x.read = int(offset)
	return offset, nil
}

func (x *ProxyReaderAt) Read(p []byte) (int, error) {
	if x.closed {
		x.onRead(-int64(x.read)) // rewind progress
		x.read = 0
		x.closed = false
	}

	if x.read == x.Len() {
		return 0, io.EOF
	}
	var n int
	var err error
	if len(p) > x.Len()-x.read {
		n, err = x.ReadAt(p[:min(x.Len()-x.read, len(p))], x.off+int64(x.read))
	} else {
		n, err = x.ReadAt(p, x.off+int64(x.read))
	}

	if err == io.EOF {
		x.eof = true
	}

	if err != nil {
		return n, err
	}

	x.read += n
	if x.onRead != nil {
		x.onRead(int64(n))
	}
	return n, nil
}

func (x *ProxyRead) Read(p []byte) (int, error) {
	if x.read == x.Len() || x.closed {
		return 0, io.EOF
	}

	n, err := x.Reader.Read(p[:min(x.Len()-x.read, len(p))])
	if err == io.EOF {
		x.eof = true
	}

	if err != nil {
		return n, err
	}

	x.read += n
	if x.onRead != nil {
		x.onRead(int64(n))
	}

	return n, err
}

func (x *ProxyReaderAt) Close() error {
	x.closed = true
	return nil
}

func (x *ProxyRead) Close() error {
	x.closed = true
	return nil
}
