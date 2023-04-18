package manager

import (
	"sync"
	"sync/atomic"

	"github.com/zenthangplus/goccm"
)

func WithWaitGroup(manager goccm.ConcurrencyManager) *WaitGroup {
	return &WaitGroup{manager, sync.WaitGroup{}, 0}
}

type WaitGroup struct {
	manager   goccm.ConcurrencyManager
	waitGroup sync.WaitGroup
	count     int32
}

func (w *WaitGroup) Add() {
	w.manager.Wait()
	w.waitGroup.Add(1)
	atomic.AddInt32(&w.count, 1)
}

func (w *WaitGroup) Count() int {
	return int(atomic.LoadInt32(&w.count))
}

func (w *WaitGroup) Wait() {
	w.waitGroup.Wait()
}

func (w *WaitGroup) Done() {
	w.manager.Done()
	w.waitGroup.Done()
	atomic.AddInt32(&w.count, -1)
}
