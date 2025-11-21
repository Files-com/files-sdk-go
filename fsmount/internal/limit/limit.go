package sync

import (
	"context"

	"golang.org/x/sync/semaphore"
)

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
func (m *FuseOpLimiter) WithLimit(ctx context.Context, cl FuseOpType, fn func(context.Context) error) error {
	held := getHeld(ctx)
	if held[cl] {
		return fn(ctx) // already holds this class; skip acquire
	}

	// Build the acquisition plan
	var toRelease []func()
	acquire := func(s *semaphore.Weighted, n int64) error {
		if s == nil || n == 0 {
			return nil
		}
		if err := s.Acquire(ctx, n); err != nil {
			return err
		}
		toRelease = append(toRelease, func() { s.Release(n) })
		return nil
	}

	// Always acquire global first to avoid inversion.
	if err := acquire(m.global, 1); err != nil {
		return err
	}
	if err := acquire(m.class[cl], 1); err != nil {
		// release global if class acquire fails
		for i := len(toRelease) - 1; i >= 0; i-- {
			toRelease[i]()
		}
		return err
	}

	// Mark reentrancy and run
	ctx2 := withHeld(ctx, cl)
	err := fn(ctx2)

	// Release in LIFO
	for i := len(toRelease) - 1; i >= 0; i-- {
		toRelease[i]()
	}
	return err
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
