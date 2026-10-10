package lib

import "sync"

// signalCalls makes Call's check and close one step, so concurrent calls on one generation, including calls through
// copies of a Signal that share C, close it once. The package-wide lock is held only for the channel check and close.
var signalCalls sync.Mutex

type Signal struct {
	C chan bool
}

func (s *Signal) Init() *Signal {
	s.C = make(chan bool)
	return s
}

func (s *Signal) Called() bool {
	select {
	case <-s.C:
		return true
	default:
		return false
	}
}

func (s *Signal) Call() {
	signalCalls.Lock()
	defer signalCalls.Unlock()
	select {
	case <-s.C:
	default:
		close(s.C)
	}
}

func (s *Signal) Clear() {
	s.Init()
}
