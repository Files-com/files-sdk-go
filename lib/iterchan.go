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
	select {
	case current := <-i.Send:
		i.current.Store(current)
	case err := <-i.SendError:
		i.Error.Store(iterError{err})
	case <-i.Done():
		return false
	}

	return true
}

func (i *IterChan[T]) Current() interface{} {
	return i.current.Load()
}

func (i *IterChan[T]) Resource() T {
	return i.current.Load().(T)
}

func (i *IterChan[T]) Err() error {
	switch err := i.Error.Load().(type) {
	case iterError:
		return err.error
	case error:
		return err
	}
	return nil
}

// iterError holds each error Next stores in Error, so that errors of different types share the one concrete type
// atomic.Value requires. It is itself an error that unwraps to the error received, which Err returns.
type iterError struct {
	error
}

func (e iterError) Unwrap() error {
	return e.error
}
