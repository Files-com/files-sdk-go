//go:build linux || windows

package disk

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache"
)

func BenchmarkDiskCacheRangePublication(b *testing.B) {
	const size = 64 * 1024 * 1024
	for _, sparse := range []bool{false, true} {
		b.Run(fmt.Sprintf("sparse=%t", sparse), func(b *testing.B) {
			chunk := bytes.Repeat([]byte("x"), 128*1024)
			b.SetBytes(size)
			for iteration := 0; iteration < b.N; iteration++ {
				b.StopTimer()
				dc, err := NewDiskCache(b.TempDir())
				if err != nil {
					b.Fatal(err)
				}
				b.Cleanup(dc.StopMaintenance)
				meta := cache.NewEntryMetadata("/benchmark", size, time.Now())
				b.StartTimer()
				for offset := int64(0); offset < size; offset += int64(len(chunk)) {
					var err error
					if sparse {
						_, err = dc.WriteRange(meta.Path, meta, chunk, offset)
					} else {
						_, err = dc.WriteCompleteRange(meta.Path, meta, chunk, offset)
					}
					if err != nil {
						b.Fatal(err)
					}
				}
				b.StopTimer()
				complete, err := dc.RangeEntryComplete(meta.Path, meta)
				if err != nil || !complete {
					b.Fatalf("complete=%t err=%v", complete, err)
				}
				for offset := int64(0); offset < size; offset += int64(len(chunk)) {
					buffer := make([]byte, len(chunk))
					n, err := dc.ReadRange(meta.Path, meta, buffer, offset)
					if err != nil || n != len(buffer) || !bytes.Equal(buffer, chunk) {
						b.Fatalf("bad bytes at %d", offset)
					}
				}
			}
		})
	}
}

func TestDiskCacheEvictionDuringCompleteWriteDoesNotValidateHoles(t *testing.T) {
	dc, err := NewDiskCache(t.TempDir(), WithCapacityBytes(8))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(dc.StopMaintenance)
	meta := cache.NewEntryMetadata("/download", 8, time.Now())
	if _, err := dc.WriteCompleteRange(meta.Path, meta, []byte("AAAA"), 0); err != nil {
		t.Fatal(err)
	}
	other := cache.NewEntryMetadata("/other", 8, meta.ModTime)
	if _, err := dc.WriteCompleteRange(other.Path, other, []byte("12345678"), 0); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(dc.entryPath(meta.Path)); !os.IsNotExist(err) {
		t.Fatalf("fixture did not evict first entry: %v", err)
	}
	if _, err := dc.WriteCompleteRange(meta.Path, meta, []byte("BBBB"), 4); err != nil {
		t.Fatal(err)
	}
	complete, err := dc.RangeEntryComplete(meta.Path, meta)
	if err != nil {
		t.Fatal(err)
	}
	data := make([]byte, 8)
	n, err := dc.ReadRange(meta.Path, meta, data, 0)
	if complete || n != 0 || err != nil {
		t.Fatalf("evicted prefix must be a miss: complete=%t n=%d data=%q err=%v", complete, n, data, err)
	}
	if n, err := dc.ReadRange(meta.Path, meta, data[:4], 4); n != 4 || err != nil || string(data[:4]) != "BBBB" {
		t.Fatalf("new suffix was not retained: n=%d data=%q err=%v", n, data[:4], err)
	}
}
