package lib

import (
	"context"
	"sync/atomic"
)

type IterChan[T any] struct {
	Send      chan T
	SendError chan error
	current   atomic.Value
	Error     atomic.Value
	Start     func(*IterChan[T])
	context.Context
	Stop context.CancelFunc
}

func (i *IterChan[T]) Init(ctx context.Context) *IterChan[T] {
	i.Send = make(chan T)
	i.SendError = make(chan error)
	i.Context, i.Stop = context.WithCancel(ctx)
	return i
}

func (i *IterChan[T]) Next() bool {
	for {
		select {
		case current := <-i.Send:
			i.current.Store(current)
			return true
		case err := <-i.SendError:
			// An error is not an iteration result: without a new value stored,
			// Resource() would return nil or repeat the previous value. Record
			// it for Err() and keep waiting for the next value or Done.
			i.Error.Store(err)
		case <-i.Done():
			return false
		}
	}
}

func (i *IterChan[T]) Current() interface{} {
	return i.current.Load()
}

func (i *IterChan[T]) Resource() T {
	return i.current.Load().(T)
}

func (i *IterChan[T]) Err() error {
	err := i.Error.Load()
	if err != nil {
		return err.(error)
	}
	return nil
}
