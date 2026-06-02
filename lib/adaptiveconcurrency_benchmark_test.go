package lib

import (
	"context"
	"testing"
	"time"
)

func BenchmarkAdaptiveConcurrencyManagerAcquireRelease(b *testing.B) {
	manager := NewAdaptiveConcurrencyManager(256)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		manager.Wait()
		manager.DoneWithSample(AdaptiveConcurrencySample{
			Success:  true,
			Bytes:    32 * 1024 * 1024,
			Duration: time.Millisecond,
		})
	}
}

func BenchmarkAdaptiveConcurrencyManagerWaitWithContext(b *testing.B) {
	manager := NewAdaptiveConcurrencyManager(256)
	ctx := context.Background()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if !manager.WaitWithContext(ctx) {
			b.Fatal("unexpected canceled acquire")
		}
		manager.Done()
	}
}

func BenchmarkAdaptiveConcurrencyManagerBackPressureSample(b *testing.B) {
	manager := NewAdaptiveConcurrencyManager(256)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		manager.mu.Lock()
		manager.pauseUntil = time.Time{}
		manager.mu.Unlock()
		manager.Wait()
		manager.DoneWithSample(AdaptiveConcurrencySample{
			Success:      true,
			Bytes:        32 * 1024 * 1024,
			Duration:     time.Millisecond,
			StatusCode:   503,
			BackPressure: true,
		})
	}
}
