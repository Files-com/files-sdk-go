package manager

import (
	"sync"

	"github.com/zenthangplus/goccm"
)

func WithWaitGroup(manager goccm.ConcurrencyManager) *WaitGroup {
	return &WaitGroup{manager, sync.WaitGroup{}, 0, sync.Mutex{}}
}

type WaitGroup struct {
	manager   goccm.ConcurrencyManager
	waitGroup sync.WaitGroup
	count     int
	mutex     sync.Mutex
}

func (w *WaitGroup) Add() {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.manager.Wait()
	w.waitGroup.Add(1)
	w.count += 1
}

func (w *WaitGroup) Count() int {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	return w.count
}

func (w *WaitGroup) Wait() {
	w.waitGroup.Wait()
}

func (w *WaitGroup) Done() {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.manager.Done()
	w.waitGroup.Done()
	w.count -= 1
}
