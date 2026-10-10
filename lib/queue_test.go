package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueueIsUsableAfterClear(t *testing.T) {
	q := (&Queue[string]{}).Init(2)
	q.Push("old")
	q.Clear()

	assert.Equal(t, 0, q.Len())
	q.Push("new")
	assert.Equal(t, "new", q.Pop())
}
