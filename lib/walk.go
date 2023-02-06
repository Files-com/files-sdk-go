package lib

import (
	"context"
	"errors"
	"io/fs"

	"github.com/zenthangplus/goccm"
)

type Walk[T any] struct {
	fs.FS
	Queue[string]
	IterChan[T]
	goccm.ConcurrencyManager
	Root            string
	WalkFile        func(d fs.DirEntry, path string) (T, error)
	ListDirectories bool
}

type DirEntry struct {
	fs.DirEntry
	fs.FileInfo
	path string
	error
}

func (d DirEntry) Error() error {
	return d.error
}

func (d DirEntry) Path() string {
	return d.path
}

func (w *Walk[T]) Walk(ctx context.Context) *IterChan[T] {
	w.Queue.Init(1)
	it := (&IterChan[T]{}).Init()
	if w.Root == "" {
		w.Queue.Push(".")
	} else {
		w.Queue.Push(w.Root)
	}

	go func() {
		for {
			if ctx.Err() != nil {
				return
			}
			if dir := w.Queue.Pop(); dir != "" {
				w.Wait()
				go func() {
					err := w.walkDir(ctx, dir, it)
					if err != nil {
						it.SendError <- err
					}
					w.Done()
				}()
			}

			if it.Err() != nil {
				w.WaitAllDone()
				break
			}

			if w.Len() == 0 && w.RunningCount() == 0 {
				break
			}
		}
		it.Stop <- true
	}()
	return it
}

func (w *Walk[T]) walkDir(ctx context.Context, dir string, it *IterChan[T]) error {
	return fs.WalkDir(w.FS, dir, func(path string, d fs.DirEntry, err error) error {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if err != nil {
			if errors.Is(err, fs.ErrPermission) {
				it.Send <- DirEntry{d, nil, path, err}
				return fs.SkipDir
			}
			return err
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}

		if path == dir && d.IsDir() {
			if path == w.Root && w.ListDirectories {
				w.send(d, path, it)
			}
			return nil
		}

		if d.IsDir() && path != "." {
			if w.ListDirectories {
				w.send(d, path, it)
			}
			w.Queue.Push(path)
			return fs.SkipDir
		}

		if !d.Type().IsRegular() {
			return nil
		}

		w.send(d, path, it)

		return nil
	})
}

func (w *Walk[T]) send(d fs.DirEntry, path string, it *IterChan[T]) {
	toSend, err := w.WalkFile(d, path)
	if err != nil {
		it.SendError <- err
	} else {
		it.Send <- toSend
	}
}

func DirEntryWalkFile(d fs.DirEntry, path string) (DirEntry, error) {
	info, err := d.Info()
	return DirEntry{DirEntry: d, FileInfo: info, path: path, error: err}, err
}
