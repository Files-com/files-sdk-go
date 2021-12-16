package lib

import "sync/atomic"

type IterChan struct {
	Send      chan interface{}
	Stop      chan bool
	SendError chan error
	current   atomic.Value
	Error     atomic.Value
	Start     func(*IterChan)
}

func (i IterChan) Init() *IterChan {
	i.Send = make(chan interface{})
	i.Stop = make(chan bool)
	i.SendError = make(chan error)

	return &i
}

func (i *IterChan) Next() bool {
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

func (i *IterChan) Current() interface{} {
	return i.current.Load()
}

func (i *IterChan) Err() error {
	err := i.Error.Load()
	if err != nil {
		return err.(error)
	}
	return nil
}
