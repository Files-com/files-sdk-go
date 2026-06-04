package lib

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAdaptiveConcurrencyManagerStartsAtCapAndGrows(t *testing.T) {
	manager := NewAdaptiveConcurrencyManager(4)
	assert.Equal(t, 4, manager.Target())
	assert.Equal(t, 4, manager.Max())
}

func TestAdaptiveConcurrencyManagerSupportsCustomInitialTarget(t *testing.T) {
	manager := NewAdaptiveConcurrencyManagerWithInitial(256, 64)
	assert.Equal(t, 64, manager.Target())
	assert.Equal(t, 256, manager.Max())

	capped := NewAdaptiveConcurrencyManagerWithInitial(4, 64)
	assert.Equal(t, 4, capped.Target())
	assert.Equal(t, 4, capped.Max())

	minimum := NewAdaptiveConcurrencyManagerWithInitial(20, 0)
	assert.Equal(t, 1, minimum.Target())
	assert.Equal(t, 20, minimum.Max())
}

func TestAdaptiveConcurrencyManagerSupportsTunedConfig(t *testing.T) {
	manager := NewAdaptiveConcurrencyManagerWithConfig(AdaptiveConcurrencyConfig{
		MaxConcurrency:            50,
		InitialTarget:             24,
		MinTarget:                 8,
		GrowEvery:                 2,
		GrowStep:                  3,
		FailureShrinkPercent:      35,
		BackPressureShrinkPercent: 10,
	})
	assert.Equal(t, 24, manager.Target())

	for i := 0; i < 2; i++ {
		manager.Wait()
		manager.DoneWithSample(AdaptiveConcurrencySample{Success: true})
	}
	assert.Equal(t, 27, manager.Target())

	manager.Wait()
	manager.DoneWithSample(AdaptiveConcurrencySample{Success: true, BackPressure: true, StatusCode: 504})
	assert.Equal(t, 24, manager.Target())
	snapshot := manager.Snapshot()
	assert.Equal(t, 1, snapshot.BackPressureTotal)
	assert.Equal(t, 1, snapshot.ShrinkTotal)
}

func TestAdaptiveConcurrencyManagerBacksOffWhenGrowthDoesNotImproveThroughput(t *testing.T) {
	manager := NewAdaptiveConcurrencyManagerWithConfig(AdaptiveConcurrencyConfig{
		MaxConcurrency:            100,
		InitialTarget:             50,
		MinTarget:                 8,
		GrowEvery:                 2,
		GrowStep:                  50,
		FailureShrinkPercent:      35,
		BackPressureShrinkPercent: 10,
		ThroughputWindow:          2,
		ThroughputMinGainPercent:  10,
		ThroughputShrinkPercent:   15,
		ThroughputHoldWindows:     1,
	})

	for i := 0; i < 2; i++ {
		manager.Wait()
		time.Sleep(5 * time.Millisecond)
		manager.DoneWithSample(AdaptiveConcurrencySample{
			Success:  true,
			Bytes:    100,
			Duration: 10 * time.Millisecond,
		})
	}
	assert.Equal(t, 100, manager.Target())

	for i := 0; i < 2; i++ {
		manager.Wait()
		time.Sleep(10 * time.Millisecond)
		manager.DoneWithSample(AdaptiveConcurrencySample{
			Success:  true,
			Bytes:    100,
			Duration: 20 * time.Millisecond,
		})
	}

	assert.Equal(t, 85, manager.Target())
	snapshot := manager.Snapshot()
	assert.Equal(t, 1, snapshot.ThroughputBackoffTotal)
	assert.Equal(t, 0, snapshot.LatencyBackoffTotal)
	assert.Greater(t, snapshot.LastThroughputBytesPerSecond, float64(0))
	assert.Greater(t, snapshot.BestThroughputBytesPerSecond, float64(0))
}

func TestAdaptiveConcurrencyManagerRequiresRepeatedProbeMissesWhenConfigured(t *testing.T) {
	manager := NewAdaptiveConcurrencyManagerWithConfig(AdaptiveConcurrencyConfig{
		MaxConcurrency:            100,
		InitialTarget:             50,
		MinTarget:                 8,
		GrowEvery:                 2,
		GrowStep:                  50,
		FailureShrinkPercent:      35,
		BackPressureShrinkPercent: 10,
		ThroughputWindow:          2,
		ThroughputMinGainPercent:  10,
		ThroughputShrinkPercent:   15,
		ThroughputProbeMinWindows: 2,
	})

	for i := 0; i < 2; i++ {
		manager.Wait()
		time.Sleep(5 * time.Millisecond)
		manager.DoneWithSample(AdaptiveConcurrencySample{
			Success:  true,
			Bytes:    100,
			Duration: 10 * time.Millisecond,
		})
	}
	assert.Equal(t, 100, manager.Target())

	for i := 0; i < 2; i++ {
		manager.Wait()
		time.Sleep(10 * time.Millisecond)
		manager.DoneWithSample(AdaptiveConcurrencySample{
			Success:  true,
			Bytes:    100,
			Duration: 20 * time.Millisecond,
		})
	}
	assert.Equal(t, 100, manager.Target())

	for i := 0; i < 2; i++ {
		manager.Wait()
		time.Sleep(10 * time.Millisecond)
		manager.DoneWithSample(AdaptiveConcurrencySample{
			Success:  true,
			Bytes:    100,
			Duration: 20 * time.Millisecond,
		})
	}
	assert.Equal(t, 85, manager.Target())
	snapshot := manager.Snapshot()
	assert.Equal(t, 2, snapshot.ThroughputProbeMissTotal)
	assert.Equal(t, 1, snapshot.ThroughputBackoffTotal)
}

func TestAdaptiveConcurrencyManagerRejectsWeakGainAboveProbePlateau(t *testing.T) {
	manager := NewAdaptiveConcurrencyManagerWithConfig(AdaptiveConcurrencyConfig{
		MaxConcurrency:                         200,
		InitialTarget:                          150,
		MinTarget:                              8,
		GrowEvery:                              2,
		GrowStep:                               20,
		FailureShrinkPercent:                   35,
		BackPressureShrinkPercent:              10,
		ThroughputWindow:                       2,
		ThroughputMinGainPercent:               1,
		ThroughputShrinkPercent:                10,
		ThroughputProbePlateauTarget:           150,
		ThroughputProbeMinGainPerTargetPercent: 0.15,
	})

	for i := 0; i < 2; i++ {
		manager.Wait()
		manager.DoneWithSample(AdaptiveConcurrencySample{
			Success:  true,
			Bytes:    1_000_000,
			Duration: 100 * time.Millisecond,
		})
	}
	assert.Equal(t, 170, manager.Target())

	for i := 0; i < 2; i++ {
		manager.Wait()
		manager.DoneWithSample(AdaptiveConcurrencySample{
			Success:  true,
			Bytes:    1_025_000,
			Duration: 100 * time.Millisecond,
		})
	}

	snapshot := manager.Snapshot()
	assert.Equal(t, 153, snapshot.Target)
	assert.Equal(t, 1, snapshot.ThroughputBackoffTotal)
	assert.Equal(t, 1, snapshot.ThroughputProbeEfficiencyMissTotal)
	assert.Equal(t, 20, snapshot.LastThroughputProbeTargetDelta)
	assert.Greater(t, snapshot.LastThroughputProbeGainPercent, float64(1))
	assert.Less(t, snapshot.LastThroughputProbeGainPercent, float64(3))
	assert.Greater(t, snapshot.LastThroughputProbeGainPerTargetPercent, float64(0))
	assert.Less(t, snapshot.LastThroughputProbeGainPerTargetPercent, 0.15)
}

func TestAdaptiveConcurrencyManagerAcceptsStrongGainAboveProbePlateau(t *testing.T) {
	manager := NewAdaptiveConcurrencyManagerWithConfig(AdaptiveConcurrencyConfig{
		MaxConcurrency:                         200,
		InitialTarget:                          150,
		MinTarget:                              8,
		GrowEvery:                              2,
		GrowStep:                               20,
		FailureShrinkPercent:                   35,
		BackPressureShrinkPercent:              10,
		ThroughputWindow:                       2,
		ThroughputMinGainPercent:               1,
		ThroughputShrinkPercent:                10,
		ThroughputProbePlateauTarget:           150,
		ThroughputProbeMinGainPerTargetPercent: 0.15,
	})

	for i := 0; i < 2; i++ {
		manager.Wait()
		manager.DoneWithSample(AdaptiveConcurrencySample{
			Success:  true,
			Bytes:    1_000_000,
			Duration: 100 * time.Millisecond,
		})
	}
	assert.Equal(t, 170, manager.Target())

	for i := 0; i < 2; i++ {
		manager.Wait()
		manager.DoneWithSample(AdaptiveConcurrencySample{
			Success:  true,
			Bytes:    1_070_000,
			Duration: 100 * time.Millisecond,
		})
	}

	snapshot := manager.Snapshot()
	assert.Equal(t, 190, snapshot.Target)
	assert.Equal(t, 0, snapshot.ThroughputBackoffTotal)
	assert.Equal(t, 0, snapshot.ThroughputProbeEfficiencyMissTotal)
	assert.Equal(t, 20, snapshot.LastThroughputProbeTargetDelta)
	assert.Greater(t, snapshot.LastThroughputProbeGainPercent, float64(3))
	assert.Greater(t, snapshot.LastThroughputProbeGainPerTargetPercent, 0.15)
}

func TestAdaptiveConcurrencyManagerAcceptsNeutralGainThroughProbePlateau(t *testing.T) {
	manager := NewAdaptiveConcurrencyManagerWithConfig(AdaptiveConcurrencyConfig{
		MaxConcurrency:                         200,
		InitialTarget:                          150,
		MinTarget:                              8,
		GrowthCeiling:                          150,
		GrowthCeilingProbeSuccesses:            1,
		GrowthCeilingProbeRate:                 1,
		GrowEvery:                              2,
		GrowStep:                               20,
		FailureShrinkPercent:                   35,
		BackPressureShrinkPercent:              10,
		ThroughputWindow:                       2,
		ThroughputMinGainPercent:               1,
		ThroughputShrinkPercent:                10,
		ThroughputProbePlateauTarget:           200,
		ThroughputProbeMinGainPerTargetPercent: 0.15,
		ThroughputProbeLossTolerancePercent:    2,
	})

	for i := 0; i < 2; i++ {
		manager.Wait()
		manager.DoneWithSample(AdaptiveConcurrencySample{
			Success:  true,
			Bytes:    1_000_000,
			Duration: 100 * time.Millisecond,
		})
	}
	assert.Equal(t, 170, manager.Target())

	for i := 0; i < 2; i++ {
		manager.Wait()
		manager.DoneWithSample(AdaptiveConcurrencySample{
			Success:  true,
			Bytes:    990_000,
			Duration: 100 * time.Millisecond,
		})
	}

	snapshot := manager.Snapshot()
	assert.Equal(t, 0, snapshot.ThroughputBackoffTotal)
	assert.Equal(t, 0, snapshot.ThroughputProbeEfficiencyMissTotal)
	assert.GreaterOrEqual(t, snapshot.Target, 170)
}

func TestAdaptiveConcurrencyManagerRequiresProbeGainThroughPlateau(t *testing.T) {
	manager := NewAdaptiveConcurrencyManagerWithConfig(AdaptiveConcurrencyConfig{
		MaxConcurrency:                         220,
		InitialTarget:                          150,
		MinTarget:                              8,
		GrowEvery:                              2,
		GrowStep:                               20,
		FailureShrinkPercent:                   35,
		BackPressureShrinkPercent:              10,
		ThroughputWindow:                       2,
		ThroughputMinGainPercent:               10,
		ThroughputShrinkPercent:                10,
		ThroughputProbeFloor:                   150,
		ThroughputProbeFloorRate:               1,
		ThroughputProbePlateauTarget:           200,
		ThroughputProbeMinGainPerTargetPercent: 0.15,
	})

	for i := 0; i < 2; i++ {
		manager.Wait()
		manager.DoneWithSample(AdaptiveConcurrencySample{
			Success:  true,
			Bytes:    1_000_000,
			Duration: 100 * time.Millisecond,
		})
	}
	assert.Equal(t, 170, manager.Target())

	for i := 0; i < 2; i++ {
		manager.Wait()
		manager.DoneWithSample(AdaptiveConcurrencySample{
			Success:  true,
			Bytes:    900_000,
			Duration: 100 * time.Millisecond,
		})
	}
	assert.Equal(t, 153, manager.Target())

	for i := 0; i < 2; i++ {
		manager.Wait()
		manager.DoneWithSample(AdaptiveConcurrencySample{
			Success:  true,
			Bytes:    900_000,
			Duration: 100 * time.Millisecond,
		})
	}

	snapshot := manager.Snapshot()
	assert.Equal(t, 155, snapshot.Target)
	assert.Equal(t, 2, snapshot.ThroughputBackoffTotal)
	assert.Equal(t, 2, snapshot.ThroughputProbeMissTotal)
}

func TestAdaptiveConcurrencyManagerUnlocksGrowthCeilingAfterSustainedThroughput(t *testing.T) {
	manager := NewAdaptiveConcurrencyManagerWithConfig(AdaptiveConcurrencyConfig{
		MaxConcurrency:           220,
		InitialTarget:            150,
		MinTarget:                8,
		GrowthCeiling:            150,
		GrowthCeilingProbeBytes:  200,
		GrowthCeilingProbeRate:   1,
		GrowEvery:                2,
		GrowStep:                 20,
		ThroughputWindow:         2,
		ThroughputMinGainPercent: 1,
		ThroughputShrinkPercent:  10,
	})

	for i := 0; i < 2; i++ {
		manager.Wait()
		manager.DoneWithSample(AdaptiveConcurrencySample{
			Success:  true,
			Bytes:    100,
			Duration: 10 * time.Millisecond,
		})
	}

	snapshot := manager.Snapshot()
	assert.True(t, snapshot.GrowthCeilingUnlocked)
	assert.Equal(t, 170, snapshot.Target)
	assert.Equal(t, 150, snapshot.GrowthCeiling)
}

func TestAdaptiveConcurrencyManagerUnlocksGrowthCeilingAfterManySuccessfulParts(t *testing.T) {
	manager := NewAdaptiveConcurrencyManagerWithConfig(AdaptiveConcurrencyConfig{
		MaxConcurrency:              220,
		InitialTarget:               150,
		MinTarget:                   8,
		GrowthCeiling:               150,
		GrowthCeilingProbeBytes:     10_000,
		GrowthCeilingProbeSuccesses: 4,
		GrowthCeilingProbeRate:      1,
		GrowEvery:                   2,
		GrowStep:                    20,
		ThroughputWindow:            2,
		ThroughputMinGainPercent:    1,
		ThroughputShrinkPercent:     10,
	})

	for i := 0; i < 4; i++ {
		manager.Wait()
		manager.DoneWithSample(AdaptiveConcurrencySample{
			Success:  true,
			Bytes:    100,
			Duration: 10 * time.Millisecond,
		})
	}

	snapshot := manager.Snapshot()
	assert.True(t, snapshot.GrowthCeilingUnlocked)
	assert.Equal(t, 170, snapshot.Target)
	assert.Equal(t, 150, snapshot.GrowthCeiling)
	assert.Equal(t, 4, snapshot.GrowthCeilingProbeSuccessThreshold)
}

func TestAdaptiveConcurrencyManagerDoesNotTreatDisabledBytesGateAsReady(t *testing.T) {
	manager := NewAdaptiveConcurrencyManagerWithConfig(AdaptiveConcurrencyConfig{
		MaxConcurrency:              220,
		InitialTarget:               150,
		MinTarget:                   8,
		GrowthCeiling:               150,
		GrowthCeilingProbeSuccesses: 4,
		GrowthCeilingProbeRate:      1,
		GrowEvery:                   2,
		GrowStep:                    20,
		ThroughputWindow:            2,
		ThroughputMinGainPercent:    1,
		ThroughputShrinkPercent:     10,
	})

	for i := 0; i < 2; i++ {
		manager.Wait()
		manager.DoneWithSample(AdaptiveConcurrencySample{
			Success:  true,
			Bytes:    100,
			Duration: 10 * time.Millisecond,
		})
	}

	snapshot := manager.Snapshot()
	assert.False(t, snapshot.GrowthCeilingUnlocked)
	assert.Equal(t, 150, snapshot.Target)
}

func TestAdaptiveConcurrencyManagerKeepsGrowthCeilingWithoutThroughput(t *testing.T) {
	manager := NewAdaptiveConcurrencyManagerWithConfig(AdaptiveConcurrencyConfig{
		MaxConcurrency:              220,
		InitialTarget:               150,
		MinTarget:                   8,
		GrowthCeiling:               150,
		GrowthCeilingProbeSuccesses: 2,
		GrowthCeilingProbeRate:      10_000_000,
		GrowEvery:                   2,
		GrowStep:                    20,
		ThroughputWindow:            2,
		ThroughputMinGainPercent:    1,
		ThroughputShrinkPercent:     10,
	})

	for i := 0; i < 4; i++ {
		manager.Wait()
		manager.DoneWithSample(AdaptiveConcurrencySample{
			Success:  true,
			Bytes:    100,
			Duration: time.Second,
		})
	}

	snapshot := manager.Snapshot()
	assert.False(t, snapshot.GrowthCeilingUnlocked)
	assert.Equal(t, 150, snapshot.Target)
}

func TestAdaptiveConcurrencyManagerUsesProbeStepGainAbovePlateau(t *testing.T) {
	manager := NewAdaptiveConcurrencyManagerWithConfig(AdaptiveConcurrencyConfig{
		MaxConcurrency:                         220,
		InitialTarget:                          150,
		MinTarget:                              8,
		GrowEvery:                              2,
		GrowStep:                               20,
		FailureShrinkPercent:                   35,
		BackPressureShrinkPercent:              10,
		ThroughputWindow:                       2,
		ThroughputMinGainPercent:               1,
		ThroughputShrinkPercent:                10,
		ThroughputProbePlateauTarget:           150,
		ThroughputProbeMinGainPerTargetPercent: 0.15,
	})

	for i := 0; i < 2; i++ {
		manager.Wait()
		manager.DoneWithSample(AdaptiveConcurrencySample{
			Success:  true,
			Bytes:    1_000_000,
			Duration: 100 * time.Millisecond,
		})
	}
	assert.Equal(t, 170, manager.Target())

	for i := 0; i < 2; i++ {
		manager.Wait()
		manager.DoneWithSample(AdaptiveConcurrencySample{
			Success:  true,
			Bytes:    1_070_000,
			Duration: 100 * time.Millisecond,
		})
	}
	assert.Equal(t, 190, manager.Target())

	for i := 0; i < 2; i++ {
		manager.Wait()
		manager.DoneWithSample(AdaptiveConcurrencySample{
			Success:  true,
			Bytes:    1_108_000,
			Duration: 100 * time.Millisecond,
		})
	}

	snapshot := manager.Snapshot()
	assert.Equal(t, 210, snapshot.Target)
	assert.Equal(t, 0, snapshot.ThroughputBackoffTotal)
	assert.Equal(t, 0, snapshot.ThroughputProbeEfficiencyMissTotal)
	assert.Equal(t, 20, snapshot.LastThroughputProbeTargetDelta)
	assert.Greater(t, snapshot.LastThroughputProbeGainPercent, float64(3))
	assert.Less(t, snapshot.LastThroughputProbeGainPercent, float64(6))
	assert.Greater(t, snapshot.LastThroughputProbeGainPerTargetPercent, 0.15)
}

func TestAdaptiveConcurrencyManagerSkipsThroughputProbeBelowFloorAfterFastRate(t *testing.T) {
	manager := NewAdaptiveConcurrencyManagerWithConfig(AdaptiveConcurrencyConfig{
		MaxConcurrency:            200,
		InitialTarget:             50,
		MinTarget:                 8,
		GrowEvery:                 2,
		GrowStep:                  50,
		FailureShrinkPercent:      35,
		BackPressureShrinkPercent: 10,
		ThroughputWindow:          2,
		ThroughputMinGainPercent:  10,
		ThroughputShrinkPercent:   15,
		ThroughputProbeFloor:      150,
		ThroughputProbeFloorRate:  1,
	})

	for i := 0; i < 2; i++ {
		manager.Wait()
		manager.DoneWithSample(AdaptiveConcurrencySample{
			Success:  true,
			Bytes:    100,
			Duration: 10 * time.Millisecond,
		})
	}
	assert.Equal(t, 100, manager.Target())

	for i := 0; i < 2; i++ {
		manager.Wait()
		manager.DoneWithSample(AdaptiveConcurrencySample{
			Success:  true,
			Bytes:    100,
			Duration: 20 * time.Millisecond,
		})
	}

	snapshot := manager.Snapshot()
	assert.Equal(t, 150, snapshot.Target)
	assert.Equal(t, 150, snapshot.PeakTarget)
	assert.Equal(t, 0, snapshot.ThroughputBackoffTotal)
}

func TestAdaptiveConcurrencyManagerDoesNotThroughputBackoffBelowProbeFloorAfterFastRate(t *testing.T) {
	manager := NewAdaptiveConcurrencyManagerWithConfig(AdaptiveConcurrencyConfig{
		MaxConcurrency:            200,
		InitialTarget:             150,
		MinTarget:                 8,
		GrowEvery:                 2,
		GrowStep:                  20,
		FailureShrinkPercent:      35,
		BackPressureShrinkPercent: 10,
		ThroughputWindow:          2,
		ThroughputMinGainPercent:  10,
		ThroughputShrinkPercent:   15,
		ThroughputProbeFloor:      150,
		ThroughputProbeFloorRate:  1,
	})

	for i := 0; i < 2; i++ {
		manager.Wait()
		manager.DoneWithSample(AdaptiveConcurrencySample{
			Success:  true,
			Bytes:    100,
			Duration: 10 * time.Millisecond,
		})
	}
	assert.Equal(t, 170, manager.Target())

	for i := 0; i < 2; i++ {
		manager.Wait()
		manager.DoneWithSample(AdaptiveConcurrencySample{
			Success:  true,
			Bytes:    100,
			Duration: 20 * time.Millisecond,
		})
	}

	snapshot := manager.Snapshot()
	assert.Equal(t, 150, snapshot.Target)
	assert.Equal(t, 170, snapshot.PeakTarget)
	assert.Equal(t, 1, snapshot.ThroughputBackoffTotal)
}

func TestAdaptiveConcurrencyManagerStillProbesBelowFloorOnSlowRate(t *testing.T) {
	manager := NewAdaptiveConcurrencyManagerWithConfig(AdaptiveConcurrencyConfig{
		MaxConcurrency:            200,
		InitialTarget:             50,
		MinTarget:                 8,
		GrowEvery:                 2,
		GrowStep:                  50,
		FailureShrinkPercent:      35,
		BackPressureShrinkPercent: 10,
		ThroughputWindow:          2,
		ThroughputMinGainPercent:  10,
		ThroughputShrinkPercent:   15,
		ThroughputProbeFloor:      150,
		ThroughputProbeFloorRate:  1 << 60,
	})

	for i := 0; i < 2; i++ {
		manager.Wait()
		manager.DoneWithSample(AdaptiveConcurrencySample{
			Success:  true,
			Bytes:    100,
			Duration: 10 * time.Millisecond,
		})
	}
	assert.Equal(t, 100, manager.Target())

	for i := 0; i < 2; i++ {
		manager.Wait()
		manager.DoneWithSample(AdaptiveConcurrencySample{
			Success:  true,
			Bytes:    100,
			Duration: 20 * time.Millisecond,
		})
	}

	snapshot := manager.Snapshot()
	assert.Equal(t, 85, snapshot.Target)
	assert.Equal(t, 100, snapshot.PeakTarget)
	assert.Equal(t, 1, snapshot.ThroughputBackoffTotal)
}

func TestAdaptiveConcurrencyManagerThroughputWindowUsesPartDurationForBurstCompletions(t *testing.T) {
	manager := NewAdaptiveConcurrencyManagerWithConfig(AdaptiveConcurrencyConfig{
		MaxConcurrency:            100,
		InitialTarget:             50,
		MinTarget:                 8,
		GrowEvery:                 100,
		GrowStep:                  1,
		FailureShrinkPercent:      35,
		BackPressureShrinkPercent: 10,
		ThroughputWindow:          2,
		ThroughputMinGainPercent:  10,
		ThroughputShrinkPercent:   15,
	})

	for i := 0; i < 2; i++ {
		manager.Wait()
		manager.DoneWithSample(AdaptiveConcurrencySample{
			Success:  true,
			Bytes:    100,
			Duration: 100 * time.Millisecond,
		})
	}

	snapshot := manager.Snapshot()
	assert.Greater(t, snapshot.BestThroughputBytesPerSecond, float64(0))
	assert.Less(t, snapshot.BestThroughputBytesPerSecond, float64(5_000))
}

func TestAdaptiveConcurrencyManagerBacksOffOnLatencyQueuePressure(t *testing.T) {
	manager := NewAdaptiveConcurrencyManagerWithConfig(AdaptiveConcurrencyConfig{
		MaxConcurrency:            100,
		InitialTarget:             50,
		MinTarget:                 8,
		GrowEvery:                 100,
		GrowStep:                  1,
		FailureShrinkPercent:      35,
		BackPressureShrinkPercent: 10,
		ThroughputWindow:          2,
		ThroughputMinGainPercent:  3,
		ThroughputShrinkPercent:   15,
		ThroughputHoldWindows:     1,
		LatencyShrinkPercent:      10,
		LatencyQueueHigh:          12,
	})

	for i := 0; i < 2; i++ {
		manager.Wait()
		time.Sleep(5 * time.Millisecond)
		manager.DoneWithSample(AdaptiveConcurrencySample{
			Success:  true,
			Bytes:    100,
			Duration: 10 * time.Millisecond,
		})
	}

	for i := 0; i < 2; i++ {
		manager.Wait()
		time.Sleep(15 * time.Millisecond)
		manager.DoneWithSample(AdaptiveConcurrencySample{
			Success:  true,
			Bytes:    100,
			Duration: 40 * time.Millisecond,
		})
	}

	assert.Equal(t, 45, manager.Target())
	snapshot := manager.Snapshot()
	assert.Equal(t, 1, snapshot.LatencyBackoffTotal)
	assert.Greater(t, snapshot.LastQueueEstimate, float64(12))
	assert.Greater(t, snapshot.MinDurationPerByte, float64(0))
	assert.Greater(t, snapshot.LastDurationPerByte, snapshot.MinDurationPerByte)
}

func TestAdaptiveConcurrencyManagerSuppressesGrowthOnLatencyQueuePressure(t *testing.T) {
	manager := NewAdaptiveConcurrencyManagerWithConfig(AdaptiveConcurrencyConfig{
		MaxConcurrency:            100,
		InitialTarget:             50,
		MinTarget:                 8,
		GrowEvery:                 2,
		GrowStep:                  10,
		FailureShrinkPercent:      35,
		BackPressureShrinkPercent: 10,
		ThroughputWindow:          2,
		ThroughputMinGainPercent:  3,
		LatencyGrowthQueueHigh:    12,
	})

	for i := 0; i < 2; i++ {
		manager.Wait()
		manager.DoneWithSample(AdaptiveConcurrencySample{
			Success:  true,
			Bytes:    100,
			Duration: 10 * time.Millisecond,
		})
	}
	assert.Equal(t, 60, manager.Target())

	for i := 0; i < 2; i++ {
		manager.Wait()
		manager.DoneWithSample(AdaptiveConcurrencySample{
			Success:  true,
			Bytes:    100,
			Duration: 40 * time.Millisecond,
		})
	}

	snapshot := manager.Snapshot()
	assert.Equal(t, 60, snapshot.Target)
	assert.Equal(t, 1, snapshot.GrowTotal)
	assert.Equal(t, 1, snapshot.LatencyGrowthSuppressionTotal)
	assert.Greater(t, snapshot.LastQueueEstimate, float64(12))
}

func TestAdaptiveConcurrencyManagerSkipsLatencyBackoffBelowProbeFloorAfterFastRate(t *testing.T) {
	manager := NewAdaptiveConcurrencyManagerWithConfig(AdaptiveConcurrencyConfig{
		MaxConcurrency:            200,
		InitialTarget:             100,
		MinTarget:                 8,
		GrowEvery:                 100,
		GrowStep:                  1,
		FailureShrinkPercent:      35,
		BackPressureShrinkPercent: 10,
		ThroughputWindow:          2,
		ThroughputMinGainPercent:  3,
		ThroughputShrinkPercent:   15,
		ThroughputProbeFloor:      150,
		ThroughputProbeFloorRate:  1,
		LatencyShrinkPercent:      10,
		LatencyQueueHigh:          12,
	})

	for i := 0; i < 2; i++ {
		manager.Wait()
		manager.DoneWithSample(AdaptiveConcurrencySample{
			Success:  true,
			Bytes:    100,
			Duration: 10 * time.Millisecond,
		})
	}

	for i := 0; i < 2; i++ {
		manager.Wait()
		manager.DoneWithSample(AdaptiveConcurrencySample{
			Success:  true,
			Bytes:    100,
			Duration: 40 * time.Millisecond,
		})
	}

	snapshot := manager.Snapshot()
	assert.Equal(t, 100, snapshot.Target)
	assert.Equal(t, 0, snapshot.LatencyBackoffTotal)
	assert.Greater(t, snapshot.LastQueueEstimate, float64(12))
}

func TestAdaptiveConcurrencyManagerKeepsLatencyBackoffBelowProbeFloorOnSlowRate(t *testing.T) {
	manager := NewAdaptiveConcurrencyManagerWithConfig(AdaptiveConcurrencyConfig{
		MaxConcurrency:            200,
		InitialTarget:             100,
		MinTarget:                 8,
		GrowEvery:                 100,
		GrowStep:                  1,
		FailureShrinkPercent:      35,
		BackPressureShrinkPercent: 10,
		ThroughputWindow:          2,
		ThroughputMinGainPercent:  3,
		ThroughputShrinkPercent:   15,
		ThroughputProbeFloor:      150,
		ThroughputProbeFloorRate:  1 << 60,
		LatencyShrinkPercent:      10,
		LatencyQueueHigh:          12,
	})

	for i := 0; i < 2; i++ {
		manager.Wait()
		manager.DoneWithSample(AdaptiveConcurrencySample{
			Success:  true,
			Bytes:    100,
			Duration: 10 * time.Millisecond,
		})
	}

	for i := 0; i < 2; i++ {
		manager.Wait()
		manager.DoneWithSample(AdaptiveConcurrencySample{
			Success:  true,
			Bytes:    100,
			Duration: 40 * time.Millisecond,
		})
	}

	snapshot := manager.Snapshot()
	assert.Equal(t, 90, snapshot.Target)
	assert.Equal(t, 1, snapshot.LatencyBackoffTotal)
	assert.Greater(t, snapshot.LastQueueEstimate, float64(12))
}

func TestAdaptiveConcurrencyManagerDoesNotBackOffBelowLatencyFloor(t *testing.T) {
	manager := NewAdaptiveConcurrencyManagerWithConfig(AdaptiveConcurrencyConfig{
		MaxConcurrency:            100,
		InitialTarget:             60,
		MinTarget:                 8,
		GrowEvery:                 100,
		GrowStep:                  1,
		FailureShrinkPercent:      35,
		BackPressureShrinkPercent: 10,
		ThroughputWindow:          2,
		ThroughputMinGainPercent:  3,
		ThroughputShrinkPercent:   15,
		LatencyFloor:              50,
		LatencyShrinkPercent:      50,
		LatencyQueueHigh:          12,
	})

	for i := 0; i < 2; i++ {
		manager.Wait()
		time.Sleep(5 * time.Millisecond)
		manager.DoneWithSample(AdaptiveConcurrencySample{
			Success:  true,
			Bytes:    100,
			Duration: 10 * time.Millisecond,
		})
	}

	for i := 0; i < 4; i++ {
		manager.Wait()
		time.Sleep(15 * time.Millisecond)
		manager.DoneWithSample(AdaptiveConcurrencySample{
			Success:  true,
			Bytes:    100,
			Duration: 40 * time.Millisecond,
		})
	}

	assert.Equal(t, 50, manager.Target())
	snapshot := manager.Snapshot()
	assert.Equal(t, 1, snapshot.LatencyBackoffTotal)
	assert.Greater(t, snapshot.LastQueueEstimate, float64(12))
}

func TestAdaptiveConcurrencyManagerGrowsAndShrinksFromSamples(t *testing.T) {
	manager := NewAdaptiveConcurrencyManager(20)
	assert.Equal(t, 8, manager.Target())

	for i := 0; i < 8; i++ {
		manager.Wait()
		manager.DoneWithSample(AdaptiveConcurrencySample{Success: true})
	}
	assert.Equal(t, 9, manager.Target())

	manager.Wait()
	manager.DoneWithSample(AdaptiveConcurrencySample{Success: false})
	assert.Equal(t, 4, manager.Target())
	snapshot := manager.Snapshot()
	assert.Equal(t, 8, snapshot.SuccessTotal)
	assert.Equal(t, 1, snapshot.FailureTotal)
	assert.Equal(t, 1, snapshot.GrowTotal)
	assert.Equal(t, 1, snapshot.ShrinkTotal)
}

func TestAdaptiveConcurrencyManagerBackPressurePausesAndTracksRetryAfter(t *testing.T) {
	manager := NewAdaptiveConcurrencyManager(20)
	assert.Equal(t, 8, manager.Target())

	manager.Wait()
	manager.DoneWithSample(AdaptiveConcurrencySample{
		Success:      false,
		BackPressure: true,
		RetryAfter:   25 * time.Millisecond,
		StatusCode:   503,
	})

	assert.Equal(t, 4, manager.Target())
	snapshot := manager.Snapshot()
	assert.Equal(t, 1, snapshot.BackPressureTotal)
	assert.Equal(t, 1, snapshot.RetryAfterTotal)
}

func TestAdaptiveConcurrencyManagerSnapshotIncludesBytesAndDuration(t *testing.T) {
	manager := NewAdaptiveConcurrencyManager(4)
	manager.Wait()
	manager.DoneWithSample(AdaptiveConcurrencySample{
		Success:  true,
		Bytes:    1024,
		Duration: 10 * time.Millisecond,
	})

	snapshot := manager.Snapshot()
	assert.Equal(t, int64(1024), snapshot.BytesTotal)
	assert.Equal(t, 10*time.Millisecond, snapshot.AverageDuration)
	assert.Equal(t, 4, snapshot.PeakTarget)
	assert.Equal(t, 1, snapshot.PeakRunning)
}

func TestAdaptiveConcurrencyManagerShrinksAfterSuccessfulPartWithBackPressure(t *testing.T) {
	manager := NewAdaptiveConcurrencyManager(20)
	manager.Wait()
	manager.DoneWithSample(AdaptiveConcurrencySample{
		Success:      true,
		BackPressure: true,
		RetryAfter:   25 * time.Millisecond,
		StatusCode:   503,
	})

	assert.Equal(t, 4, manager.Target())
	snapshot := manager.Snapshot()
	assert.Equal(t, 1, snapshot.SuccessTotal)
	assert.Equal(t, 0, snapshot.FailureTotal)
	assert.Equal(t, 1, snapshot.BackPressureTotal)
	assert.Equal(t, 1, snapshot.RetryAfterTotal)
	assert.Equal(t, 1, snapshot.ShrinkTotal)
}

func TestAdaptiveConcurrencyManagerWaitWithContextCancelsWhileAtCapacity(t *testing.T) {
	manager := NewAdaptiveConcurrencyManager(1)
	manager.Wait()
	defer manager.Done()

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan bool, 1)
	go func() {
		done <- manager.WaitWithContext(ctx)
	}()

	cancel()
	select {
	case acquired := <-done:
		assert.False(t, acquired)
	case <-time.After(250 * time.Millisecond):
		t.Fatal("WaitWithContext did not return after cancellation")
	}
}

func TestAdaptiveConcurrencyManagerWaitWithContextCancelsDuringRetryAfterPause(t *testing.T) {
	manager := NewAdaptiveConcurrencyManager(2)
	manager.Wait()
	manager.DoneWithSample(AdaptiveConcurrencySample{
		Success:      false,
		BackPressure: true,
		RetryAfter:   time.Hour,
		StatusCode:   503,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Millisecond)
	defer cancel()
	start := time.Now()
	acquired := manager.WaitWithContext(ctx)

	assert.False(t, acquired)
	assert.Less(t, time.Since(start), 250*time.Millisecond)
}

func TestAdaptiveConcurrencyManagerClearsExpiredPauseOnAcquire(t *testing.T) {
	manager := NewAdaptiveConcurrencyManager(2)
	manager.pauseUntil = time.Now().Add(-time.Second)

	assert.True(t, manager.WaitWithContext(context.Background()))
	manager.Done()

	assert.True(t, manager.pauseUntil.IsZero())
}

func TestAdaptiveConcurrencyManagerWaitsForActivePauseToExpire(t *testing.T) {
	manager := NewAdaptiveConcurrencyManager(2)
	manager.Wait()
	manager.DoneWithSample(AdaptiveConcurrencySample{
		Success:      false,
		BackPressure: true,
		RetryAfter:   25 * time.Millisecond,
		StatusCode:   503,
	})

	start := time.Now()
	assert.True(t, manager.WaitWithContext(context.Background()))
	elapsed := time.Since(start)
	manager.Done()

	assert.GreaterOrEqual(t, elapsed, 10*time.Millisecond)
	assert.Less(t, elapsed, 250*time.Millisecond)
	assert.True(t, manager.pauseUntil.IsZero())
}
