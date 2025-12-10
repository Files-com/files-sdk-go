package sync

import "sync"

// Pool is a generic wrapper around sync.Pool for type-safe pooling
type Pool[T any] struct {
	pool sync.Pool
}

// NewPool creates a new type-safe pool with the given constructor function
func NewPool[T any](newFn func() T) *Pool[T] {
	return &Pool[T]{
		pool: sync.Pool{
			New: func() any {
				return newFn()
			},
		},
	}
}

// Get retrieves an item from the pool
func (p *Pool[T]) Get() T {
	return p.pool.Get().(T)
}

// Put returns an item to the pool
func (p *Pool[T]) Put(x T) {
	p.pool.Put(x)
}
