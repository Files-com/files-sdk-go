package status

import (
	"sync"
	"time"
)

type Signal struct {
	c      chan time.Time
	Called bool
	When   time.Time
	mu     sync.RWMutex
	subs   []chan time.Time
}

func (s *Signal) call(t time.Time) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.Called {
		panic("already called")
	}
	for _, ch := range s.subs {
		ch <- t
	}
	s.When = t
	s.Called = true
}

func (s *Signal) Subscribe() chan time.Time {
	s.mu.Lock()
	if s.Called {
		panic("Can't Subscribe after called")
	}
	defer s.mu.Unlock()
	c := make(chan time.Time)
	s.subs = append(s.subs, c)
	return c
}

func (s *Signal) Clear() {
	s.Called = false
}
