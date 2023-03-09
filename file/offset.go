package file

import (
	"context"
	"math"
	"sync"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/Files-com/files-sdk-go/v2/lib"
)

type OffSet struct {
	off int64
	len int64
}

type Part struct {
	OffSet
	files_sdk.EtagsParam
	bytes    int64
	requests []time.Time
	error
	number     int64
	processing bool
	context.Context
	context.CancelFunc
	*sync.RWMutex
}

func (p *Part) Done() *Part {
	p.processing = false
	return p
}

func (p *Part) Start(ctx ...context.Context) *Part {
	p.Touch()
	p.processing = true
	if len(ctx) == 1 {
		p.WithContext(ctx[0])
	}
	return p
}

func (p *Part) WithContext(ctx context.Context) *Part {
	p.RWMutex = &sync.RWMutex{}
	p.Context, p.CancelFunc = context.WithCancel(ctx)
	return p
}

func (p *Part) Touch() {
	p.requests = append(p.requests, time.Now())
}

func (p *Part) Successful() bool {
	return p.bytes == p.len && p.error == nil
}

func (p *Part) Clear() {
	p.bytes = 0
	p.error = nil
}

func (p *Part) SetError(err error) {
	p.Lock()
	defer p.Unlock()
	p.error = err
}

func (p *Part) Err() error {
	p.RLock()
	defer p.RUnlock()
	if p.error != nil {
		return p.error
	}
	return p.Context.Err()
}

type Parts []*Part

func (p Parts) SuccessfulBytes() (b int64) {
	for _, part := range p {
		if part.Successful() {
			b += part.bytes
		}
	}

	return b
}

func byteOffsetSlice(size int64) []OffSet {
	partSizes := lib.PartSizes
	var blockSize int64
	var offsets []OffSet
	off := int64(0)
	blockSize, partSizes = partSizes[0], partSizes[1:]
	endRange := blockSize
	for {
		if off < size {
			endRange = int64(math.Min(float64(endRange), float64(size)))
			offsets = append(offsets, OffSet{off: off, len: endRange - off})
			off = endRange
			endRange = off + blockSize
			blockSize, partSizes = partSizes[0], partSizes[1:]
		} else {
			break
		}
	}
	return offsets
}

func byteChunkSlice(size int64, defaultChunkSize int64) []OffSet {
	var blockSize int64
	var offsets []OffSet
	off := int64(0)

	if size < defaultChunkSize {
		blockSize = size
	} else {
		blockSize = defaultChunkSize
	}
	endRange := blockSize
	for {
		if off < size {
			endRange = int64(math.Min(float64(endRange), float64(size)))
			offsets = append(offsets, OffSet{off: off, len: endRange - off})
			off = endRange
			endRange = off + blockSize
		} else {
			break
		}
	}
	return offsets
}
