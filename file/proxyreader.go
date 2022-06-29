package file

import (
	"io"
)

type ProxyReader struct {
	io.ReaderAt
	off    int64
	len    int64
	onRead func(i int64)
	read   int
	closed bool
}

func (x *ProxyReader) Len() int {
	return int(x.len)
}

func (x *ProxyReader) Seek(offset int64, whence int) (int64, error) {
	x.onRead(-int64(x.read - int(offset))) // rewind progress
	x.read = int(offset)
	return offset, nil
}

func (x *ProxyReader) Read(p []byte) (int, error) {
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
		n, err = x.ReadAt(p[0:x.Len()-x.read], x.off+int64(x.read))
	} else {
		n, err = x.ReadAt(p, x.off+int64(x.read))
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

func (x *ProxyReader) Close() error {
	x.closed = true
	return nil
}
