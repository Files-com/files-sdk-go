package cache

import (
	"context"
	"io"
	"sync"
)

// ReadyGate manages blocking cache consumers until the requested data is available.
type ReadyGate struct {
	mu        sync.Mutex
	cond      *sync.Cond
	available int64 // bytes [0, available) present contiguously
	total     int64 // -1 unknown
	err       error
	done      bool
	waiters   int
	Cancel    context.CancelFunc
}

func NewReadyGate() *ReadyGate {
	gate := &ReadyGate{total: -1}
	gate.cond = sync.NewCond(&gate.mu)
	return gate
}

func (gate *ReadyGate) Add() {
	gate.mu.Lock()
	gate.waiters++
	gate.mu.Unlock()
}

func (gate *ReadyGate) Done() {
	gate.mu.Lock()
	if gate.waiters > 0 {
		gate.waiters--
	}
	gate.mu.Unlock()
}

func (gate *ReadyGate) SetAvailable(x int64) {
	gate.mu.Lock()
	if x > gate.available {
		gate.available = x
	}
	gate.cond.Broadcast()
	gate.mu.Unlock()
}

func (gate *ReadyGate) SetCancel(cancel context.CancelFunc) {
	gate.mu.Lock()
	gate.Cancel = cancel
	gate.mu.Unlock()
}

func (gate *ReadyGate) Finish(err error, total int64) {
	gate.mu.Lock()
	gate.err, gate.done = err, true
	if total >= 0 {
		gate.total = total
	}
	gate.cond.Broadcast()
	gate.mu.Unlock()
}

func (gate *ReadyGate) WaitFor(end int64) error {
	gate.mu.Lock()
	defer gate.mu.Unlock()
	for {
		if gate.available >= end {
			return nil
		}
		if gate.done {
			if gate.err != nil {
				return gate.err
			}
			return io.EOF // smaller than requested
		}
		gate.cond.Wait()
	}
}

func (gate *ReadyGate) Available() int64 {
	gate.mu.Lock()
	defer gate.mu.Unlock()
	return gate.available
}
