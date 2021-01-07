package lib

import "io"

type ProgressWriter struct {
	Writer          io.Writer
	ProgressWatcher func(int64)
}

func (w ProgressWriter) Write(p []byte) (n int, err error) {
	n, err = w.Writer.Write(p)

	w.ProgressWatcher(int64(n))
	return n, err
}
