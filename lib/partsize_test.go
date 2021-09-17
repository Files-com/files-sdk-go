package lib

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PartSize(t *testing.T) {
	assert := assert.New(t)
	s := PartSizes
	testSize := int64(0)
	for _, b := range s {
		testSize += b
	}

	assert.Equal("4.90", fmt.Sprintf("%.2f", float64(testSize)/1024/1024/1024/1024))
	assert.Equal(len(s), 10_000)
	assert.LessOrEqual(testSize, int64(1024*1024*1024*1024*5))
}
