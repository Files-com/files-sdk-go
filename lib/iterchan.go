package lib

import "sync/atomic"

type IterChan[T any] struct {
	Send      chan interface{}
	Stop      chan bool
	SendError chan error
	current   atomic.Value
	Error     atomic.Value
	Start     func(*IterChan[T])
}

func (i *IterChan[T]) Init() *IterChan[T] {
	i.Send = make(chan interface{})
	i.Stop = make(chan bool)
	i.SendError = make(chan error)

	return i
}

func (i *IterChan[T]) Next() bool {
	select {
	case current := <-i.Send:
		i.current.Store(current)
		return true
	case err := <-i.SendError:
		i.Error.Store(err)
		return false
	case <-i.Stop:
		return false
	}
}

func (i *IterChan[T]) Current() T {
	return i.current.Load().(T)
}

func (i *IterChan[T]) Err() error {
	err := i.Error.Load()
	if err != nil {
		return err.(error)
	}
	return nil
}
