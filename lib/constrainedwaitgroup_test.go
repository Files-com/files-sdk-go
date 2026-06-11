package lib

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConstrainedWorkGroupWaitForADoneWithContextCancels(t *testing.T) {
	group := NewConstrainedWorkGroup(1)
	group.Wait()
	defer group.Done()

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan bool, 1)
	go func() {
		done <- group.WaitForADoneWithContext(ctx)
	}()

	cancel()
	select {
	case ok := <-done:
		assert.False(t, ok)
	case <-time.After(250 * time.Millisecond):
		t.Fatal("WaitForADoneWithContext did not return after cancellation")
	}
}

func TestConstrainedWorkGroupWaitForADoneWithContextWakesOnDone(t *testing.T) {
	group := NewConstrainedWorkGroup(1)
	group.Wait()

	done := make(chan bool, 1)
	waitingCtx, waiting := newObservedWaitContext()
	go func() {
		done <- group.WaitForADoneWithContext(waitingCtx)
	}()

	<-waiting
	group.Done()
	select {
	case ok := <-done:
		assert.True(t, ok)
	case <-time.After(250 * time.Millisecond):
		t.Fatal("WaitForADoneWithContext did not return after Done")
	}
}

type observedWaitContext struct {
	context.Context
	entered chan struct{}
	once    sync.Once
}

func newObservedWaitContext() (context.Context, <-chan struct{}) {
	ctx := &observedWaitContext{
		Context: context.Background(),
		entered: make(chan struct{}),
	}
	return ctx, ctx.entered
}

func (c *observedWaitContext) Done() <-chan struct{} {
	c.once.Do(func() {
		close(c.entered)
	})
	return c.Context.Done()
}

func TestSubWorkerWaitForADoneWithContextCancels(t *testing.T) {
	group := NewConstrainedWorkGroup(1)
	subWorker := group.NewSubWorker().(*SubWorker)
	subWorker.Wait()
	defer subWorker.Done()

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan bool, 1)
	go func() {
		done <- subWorker.WaitForADoneWithContext(ctx)
	}()

	cancel()
	select {
	case ok := <-done:
		assert.False(t, ok)
	case <-time.After(250 * time.Millisecond):
		t.Fatal("WaitForADoneWithContext did not return after cancellation")
	}
}
