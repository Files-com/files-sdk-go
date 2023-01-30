package status

import (
	"sync"
	"time"
)

type Signal struct {
	c      chan time.Time
	called bool
	when   time.Time
	mu     sync.RWMutex
	subs   []chan time.Time
}

func (s *Signal) Called() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.called
}

func (s *Signal) When() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.when
}

func (s *Signal) call(t time.Time) {
	s.mu.Lock()
	if s.called {
		panic("already called")
	}
	s.when = t
	s.called = true
	s.mu.Unlock()
	for _, ch := range s.subs {
		ch <- t
	}
}

func (s *Signal) TrySubscribe(try func(chan time.Time)) {
	s.mu.Lock()
	if s.called {
		s.mu.Unlock()
		return
	}
	c := make(chan time.Time)
	s.subs = append(s.subs, c)
	s.mu.Unlock()
	try(c)
}

func (s *Signal) Wait() {
	s.TrySubscribe(func(c chan time.Time) {
		<-c
	})
}

func (s *Signal) Subscribe() chan time.Time {
	s.mu.Lock()
	if s.called {
		panic("Can't Subscribe after called")
	}
	defer s.mu.Unlock()
	c := make(chan time.Time)
	s.subs = append(s.subs, c)
	return c
}

func (s *Signal) Clear() {
	s.called = false
}
