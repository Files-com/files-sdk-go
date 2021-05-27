package file

import (
	"io"
	"math"
)

type ProxyReader struct {
	io.ReaderAt
	off    int64
	len    int64
	onRead func(i int64)
	read   int
}

func (x *ProxyReader) Len() int {
	return int(x.len)
}

func (x *ProxyReader) Read(p []byte) (int, error) {
	if x.read == int(x.len) {
		return 0, io.EOF
	}
	buffLen := int(math.Min(float64(int(x.len)-x.read), float64(len(p))))
	buff := make([]byte, buffLen)
	n, err := x.ReadAt(buff, x.off+int64(x.read))
	if err != nil {
		return n, err
	}

	n = copy(p, buff)
	x.read += n
	if x.onRead != nil {
		x.onRead(int64(n))
	}
	return n, nil
}

func (x *ProxyReader) Close() error { return nil }
