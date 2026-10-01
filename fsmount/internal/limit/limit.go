//go:build linux || windows

package sync

import (
	"context"
	"errors"
	"sync"

	"golang.org/x/sync/semaphore"
)

// ErrNoSlotsAvailable is returned when WithLimit cannot acquire a slot immediately.
var ErrNoSlotsAvailable = errors.New("no slots available")

// FuseOpType represents the type of operation for limiting purposes.
//
// e.g. upload vs download operations.
type FuseOpType uint8

const (
	_ FuseOpType = iota

	// FuseOpUpload represents upload operations.
	FuseOpUpload

	// FuseOpDownload represents download operations.
	FuseOpDownload

	// FuseOpOther represents other operations. e.g. list, find, delete, etc.
	FuseOpOther
)

type limiterKey struct{} // for reentrancy tagging

// FuseOpLimiter enforces multiple classes of limits, plus an optional global
// limit.
type FuseOpLimiter struct {
	class map[FuseOpType]*semaphore.Weighted
	// Optional global budget; nil to disable.
	global *semaphore.Weighted
}

// NewFuseOpLimiter creates a new MultiLimiter with the given per-class and
// global limits. A limit of 0 disables that limit.
func NewFuseOpLimiter(perClass map[FuseOpType]int64, global int64) *FuseOpLimiter {
	m := &FuseOpLimiter{class: make(map[FuseOpType]*semaphore.Weighted)}
	for k, v := range perClass {
		if v > 0 {
			m.class[k] = semaphore.NewWeighted(v)
		}
	}
	if global > 0 {
		m.global = semaphore.NewWeighted(global)
	}
	return m
}

// WithLimit applies the class (and global) limit. Reentrant-safe: if the
// context already holds the class token, it won't acquire again.
// This method blocks until slots are available.
func (m *FuseOpLimiter) WithLimit(ctx context.Context, cl FuseOpType, fn func(context.Context) error) error {
	ctx, release, err := m.Acquire(ctx, cl)
	if err != nil {
		return err
	}
	defer release()
	return fn(ctx)
}

// Acquire holds a slot until release is called. Streaming operations must keep
// it until their body is closed. The returned release function is idempotent.
func (m *FuseOpLimiter) Acquire(ctx context.Context, cl FuseOpType) (context.Context, func(), error) {
	acquire := func(s *semaphore.Weighted, n int64) error {
		if s == nil || n == 0 {
			return nil
		}
		if err := s.Acquire(ctx, n); err != nil {
			return err
		}
		return nil
	}
	return m.acquire(ctx, cl, acquire)
}

// TryWithLimit applies the class (and global) limit without blocking. Reentrant-safe: if the
// context already holds the class token, it won't acquire again.
// Returns ErrNoSlotsAvailable immediately if slots cannot be acquired.
func (m *FuseOpLimiter) TryWithLimit(ctx context.Context, cl FuseOpType, fn func(context.Context) error) error {
	acquire := func(s *semaphore.Weighted, n int64) error {
		if s == nil || n == 0 {
			return nil
		}
		if !s.TryAcquire(n) {
			return ErrNoSlotsAvailable
		}
		return nil
	}
	ctx, release, err := m.acquire(ctx, cl, acquire)
	if err != nil {
		return err
	}
	defer release()
	return fn(ctx)
}

func (m *FuseOpLimiter) acquire(ctx context.Context, cl FuseOpType, acquire func(*semaphore.Weighted, int64) error) (context.Context, func(), error) {
	held := getHeld(ctx)
	if held[cl] {
		return ctx, func() {}, nil // already holds this class; skip acquire
	}

	// Build the acquisition plan
	var toRelease []func()
	acquireWithRelease := func(s *semaphore.Weighted, n int64) error {
		if err := acquire(s, n); err != nil {
			return err
		}
		if s != nil && n > 0 {
			toRelease = append(toRelease, func() { s.Release(n) })
		}
		return nil
	}

	// Always acquire global first to avoid inversion.
	if err := acquireWithRelease(m.global, 1); err != nil {
		return nil, nil, err
	}
	if err := acquireWithRelease(m.class[cl], 1); err != nil {
		// release global if class acquire fails
		for i := len(toRelease) - 1; i >= 0; i-- {
			toRelease[i]()
		}
		return nil, nil, err
	}

	release := sync.OnceFunc(func() {
		for i := len(toRelease) - 1; i >= 0; i-- {
			toRelease[i]()
		}
	})
	return withHeld(ctx, cl), release, nil
}

type heldSet map[FuseOpType]bool

func getHeld(ctx context.Context) heldSet {
	if v := ctx.Value(limiterKey{}); v != nil {
		return v.(heldSet)
	}
	return heldSet{}
}

func withHeld(ctx context.Context, cl FuseOpType) context.Context {
	h := getHeld(ctx)
	// Copy on write to avoid mutating parent
	h2 := make(heldSet, len(h)+1)
	for k, v := range h {
		h2[k] = v
	}
	h2[cl] = true
	return context.WithValue(ctx, limiterKey{}, h2)
}
