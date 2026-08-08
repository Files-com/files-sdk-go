package lib

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIterChan_ErrorBeforeFirstValue(t *testing.T) {
	it := (&IterChan[int]{}).Init(context.Background())
	sendErr := errors.New("stat failed")

	go func() {
		it.SendError <- sendErr
		it.Send <- 1
		it.Stop()
	}()

	var values []int
	for it.Next() {
		values = append(values, it.Resource())
	}

	assert.Equal(t, []int{1}, values)
	assert.Equal(t, sendErr, it.Err())
}

func TestIterChan_ErrorBetweenValues(t *testing.T) {
	it := (&IterChan[int]{}).Init(context.Background())
	sendErr := errors.New("stat failed")

	go func() {
		it.Send <- 1
		it.SendError <- sendErr
		it.Send <- 2
		it.Stop()
	}()

	var values []int
	for it.Next() {
		values = append(values, it.Resource())
	}

	assert.Equal(t, []int{1, 2}, values)
	assert.Equal(t, sendErr, it.Err())
}

func TestIterChan_ErrorOnly(t *testing.T) {
	it := (&IterChan[int]{}).Init(context.Background())
	sendErr := errors.New("stat failed")

	go func() {
		it.SendError <- sendErr
		it.Stop()
	}()

	for it.Next() {
		t.Fatalf("Next returned true with no value sent")
	}

	assert.Equal(t, sendErr, it.Err())
}
