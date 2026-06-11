package file

import (
	"io"
	"sync"
	"testing"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/lib"
)

func BenchmarkUploadV2PartPlannerKnownSize100GiB(b *testing.B) {
	part := uploadV2TestPart("https://uploads.example.com/upload/file")
	size := int64(100) * uploadV2GiB
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		iterator, ok := newUploadV2PartIterator(part, &size, 0, 0)
		if !ok {
			b.Fatal("planner unexpectedly unavailable")
		}
		var total int64
		for iterator != nil {
			var offset OffSet
			offset, iterator, _ = iterator()
			total += offset.len
		}
		if total != size {
			b.Fatalf("total = %d, want %d", total, size)
		}
	}
}

func BenchmarkUploadV2PartPlannerKnownSize5TiBS3(b *testing.B) {
	part := uploadV2TestPart("https://s3.amazonaws.com/bucket/key?partNumber=1")
	size := 5 * uploadV2TiB
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		iterator, ok := newUploadV2PartIterator(part, &size, 0, 0)
		if !ok {
			b.Fatal("planner unexpectedly unavailable")
		}
		var total int64
		for iterator != nil {
			var offset OffSet
			offset, iterator, _ = iterator()
			total += offset.len
		}
		if total != size {
			b.Fatalf("total = %d, want %d", total, size)
		}
	}
}

func BenchmarkUploadV2UnknownSizeGrowth(b *testing.B) {
	plan, ok, reason := newUploadV2PartPlan(uploadV2TargetDefault, nil)
	if !ok {
		b.Fatal(reason)
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var total int64
		for index := 0; index < 512; index++ {
			total += plan.partSizeForIndex(index)
		}
		if total == 0 {
			b.Fatal("unexpected zero total")
		}
	}
}

func BenchmarkUploadV2DecorateUploadURL(b *testing.B) {
	part := uploadV2TestPart("https://uploads.example.com/upload/file?partNumber=1&offset=old")
	engine := newUploadV2BenchmarkEngine(part)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		upload := part
		engine.decorateUploadURL(&upload, (i%10_000)+1, int64(i)*32*uploadV2MiB)
	}
}

func BenchmarkUploadV2ResumeValidationTenThousandParts(b *testing.B) {
	part := uploadV2TestPart("https://uploads.example.com/upload/file")
	engine := newUploadV2BenchmarkEngine(part)
	engine.u.Size = lib.Int64(10_000 * 32 * uploadV2MiB)
	engine.u.Parts = make(Parts, 10_000)
	for index := range engine.u.Parts {
		engine.u.Parts[index] = &Part{
			OffSet: OffSet{off: int64(index) * 32 * uploadV2MiB, len: 32 * uploadV2MiB},
			number: index + 1,
			error:  errUploadV2BenchmarkPending{},
		}
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if reason := engine.resumeResetReason(); reason != "" {
			b.Fatal(reason)
		}
	}
}

func BenchmarkUploadV2ProgressUnbatchedReaderAt64MiB(b *testing.B) {
	benchmarkUploadV2ProgressReaderAt64MiB(b, false)
}

func BenchmarkUploadV2ProgressBatchedReaderAt64MiB(b *testing.B) {
	benchmarkUploadV2ProgressReaderAt64MiB(b, true)
}

func benchmarkUploadV2ProgressReaderAt64MiB(b *testing.B, batched bool) {
	const partSize = 64 * uploadV2MiB
	buffer := make([]byte, 32*1024)
	readerAt := uploadV2BenchmarkReaderAt{}
	b.SetBytes(partSize)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var progress int64
		var calls int64
		var mu sync.Mutex
		progressCallback := func(delta int64) {
			mu.Lock()
			defer mu.Unlock()
			progress += delta
			calls++
		}
		onRead := progressCallback
		var batcher *uploadV2ProgressBatcher
		if batched {
			batcher = newUploadV2ProgressBatcher(progressCallback)
			onRead = batcher.Add
		}
		reader := &ProxyReaderAt{
			ReaderAt: readerAt,
			len:      partSize,
			onRead:   onRead,
		}
		n, err := io.CopyBuffer(uploadV2BenchmarkDiscard{}, reader, buffer)
		if err != nil {
			b.Fatal(err)
		}
		if batcher != nil {
			batcher.Flush()
		}
		mu.Lock()
		gotProgress := progress
		gotCalls := calls
		mu.Unlock()
		if n != partSize || gotProgress != partSize {
			b.Fatalf("read=%d progress=%d want=%d", n, gotProgress, partSize)
		}
		expectedUnbatchedCalls := partSize / int64(len(buffer))
		if batched && gotCalls >= expectedUnbatchedCalls {
			b.Fatalf("batched progress made %d callback calls", gotCalls)
		}
		if !batched && gotCalls != expectedUnbatchedCalls {
			b.Fatalf("unbatched progress made %d callback calls", gotCalls)
		}
	}
}

func BenchmarkUploadV2ProgressBatchedParallelReaders64MiB(b *testing.B) {
	const partSize = 64 * uploadV2MiB
	readerAt := uploadV2BenchmarkReaderAt{}
	b.SetBytes(partSize)
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		buffer := make([]byte, 32*1024)
		for pb.Next() {
			var progress int64
			var mu sync.Mutex
			batcher := newUploadV2ProgressBatcher(func(delta int64) {
				mu.Lock()
				defer mu.Unlock()
				progress += delta
			})
			reader := &ProxyReaderAt{
				ReaderAt: readerAt,
				len:      partSize,
				onRead:   batcher.Add,
			}
			n, err := io.CopyBuffer(uploadV2BenchmarkDiscard{}, reader, buffer)
			if err != nil {
				b.Fatal(err)
			}
			batcher.Flush()
			mu.Lock()
			gotProgress := progress
			mu.Unlock()
			if n != partSize || gotProgress != partSize {
				b.Fatalf("read=%d progress=%d want=%d", n, gotProgress, partSize)
			}
		}
	})
}

func BenchmarkUploadV2ProxyReaderAtNoTiming64MiB(b *testing.B) {
	benchmarkUploadV2ReadAtReader64MiB(b, func(readerAt io.ReaderAt, partSize int64) io.Reader {
		return &ProxyReaderAt{
			ReaderAt: readerAt,
			len:      partSize,
		}
	})
}

func BenchmarkUploadV2ProxyReaderAtTiming64MiB(b *testing.B) {
	benchmarkUploadV2ReadAtReader64MiB(b, func(readerAt io.ReaderAt, partSize int64) io.Reader {
		return &ProxyReaderAt{
			ReaderAt:          readerAt,
			len:               partSize,
			trackReadDuration: true,
		}
	})
}

func BenchmarkUploadV2ProxySectionReaderNoTiming64MiB(b *testing.B) {
	benchmarkUploadV2ReadAtReader64MiB(b, func(readerAt io.ReaderAt, partSize int64) io.Reader {
		return newProxySectionReader(readerAt, 0, partSize, nil, false)
	})
}

func BenchmarkUploadV2ProxySectionReaderTiming64MiB(b *testing.B) {
	benchmarkUploadV2ReadAtReader64MiB(b, func(readerAt io.ReaderAt, partSize int64) io.Reader {
		return newProxySectionReader(readerAt, 0, partSize, nil, true)
	})
}

func BenchmarkUploadV2SectionReader64MiB(b *testing.B) {
	benchmarkUploadV2ReadAtReader64MiB(b, func(readerAt io.ReaderAt, partSize int64) io.Reader {
		return io.NewSectionReader(readerAt, 0, partSize)
	})
}

func BenchmarkUploadV2SchedulerStatsHTTPCallSerial(b *testing.B) {
	var stats uploadV2SchedulerStats
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		stats.recordHTTPCall(20*time.Millisecond, 32*uploadV2MiB, nil)
	}
}

func BenchmarkUploadV2SchedulerStatsHTTPCallParallel(b *testing.B) {
	var stats uploadV2SchedulerStats
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			stats.recordHTTPCall(20*time.Millisecond, 32*uploadV2MiB, nil)
		}
	})
}

func BenchmarkUploadV2SchedulerStatsProducerPath(b *testing.B) {
	var stats uploadV2SchedulerStats
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		stats.recordAdaptiveWait(time.Nanosecond, true)
		stats.recordReaderBuild(32*uploadV2MiB, time.Nanosecond, nil)
		stats.recordPrepared(32*uploadV2MiB, 0)
		stats.recordReadyDepth(32, 0)
		stats.recordDispatchedPrepared()
	}
}

func BenchmarkUploadV2SchedulerStatsInstrumentedPartLifecycle(b *testing.B) {
	var stats uploadV2SchedulerStats
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		start := time.Now()
		stats.recordAdaptiveWait(time.Since(start), true)
		start = time.Now()
		stats.recordReaderBuild(32*uploadV2MiB, time.Since(start), nil)
		start = time.Now()
		stats.recordHTTPCall(time.Since(start), 32*uploadV2MiB, nil)
	}
}

func BenchmarkUploadV2SchedulerStatsAddLogAttrs(b *testing.B) {
	var stats uploadV2SchedulerStats
	stats.recordDirectScheduled()
	stats.recordPrepared(32*uploadV2MiB, 0)
	stats.recordDispatchedPrepared()
	stats.recordRunwayPartCapBlock()
	stats.recordReadyDepth(32, 0)
	stats.recordReaderBuild(32*uploadV2MiB, time.Millisecond, nil)
	stats.recordAdaptiveWait(2*time.Millisecond, true)
	stats.recordUploadURLRefresh(3*time.Millisecond, nil)
	stats.recordHTTPCall(4*time.Millisecond, 32*uploadV2MiB, nil)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		attrs := make(map[string]any, 32)
		stats.addLogAttrs(attrs)
	}
}

func benchmarkUploadV2ReadAtReader64MiB(b *testing.B, newReader func(io.ReaderAt, int64) io.Reader) {
	const partSize = 64 * uploadV2MiB
	readerAt := uploadV2BenchmarkReaderAt{}
	buffer := make([]byte, 32*1024)
	b.SetBytes(partSize)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		n, err := io.CopyBuffer(uploadV2BenchmarkDiscard{}, newReader(readerAt, partSize), buffer)
		if err != nil {
			b.Fatal(err)
		}
		if n != partSize {
			b.Fatalf("read=%d want=%d", n, partSize)
		}
	}
}

func newUploadV2BenchmarkEngine(part files_sdk.FileUploadPart) *uploadV2Engine {
	size := int64(100)
	plan, ok, reason := newUploadV2PartPlanForUpload(part, &size)
	if !ok {
		panic(reason)
	}
	return newUploadV2Engine(&uploadIO{FileUploadPart: part, Manager: lib.NewConstrainedWorkGroup(1)}, plan)
}

type errUploadV2BenchmarkPending struct{}

func (errUploadV2BenchmarkPending) Error() string {
	return "pending"
}

type uploadV2BenchmarkReaderAt struct{}

func (uploadV2BenchmarkReaderAt) ReadAt(p []byte, off int64) (int, error) {
	return len(p), nil
}

type uploadV2BenchmarkDiscard struct{}

func (uploadV2BenchmarkDiscard) Write(p []byte) (int, error) {
	return len(p), nil
}
