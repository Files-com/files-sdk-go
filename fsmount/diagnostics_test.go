package fsmount

import (
	"testing"
	"time"

	"github.com/winfsp/cgofuse/fuse"
)

func TestShouldLogMountCallback(t *testing.T) {
	tests := []struct {
		name     string
		op       string
		errc     int
		duration time.Duration
		want     bool
	}{
		{
			name:     "fast success is quiet",
			op:       "Getattr",
			errc:     0,
			duration: time.Millisecond,
			want:     false,
		},
		{
			name:     "slow success is logged",
			op:       "Getattr",
			errc:     0,
			duration: mountDiagnosticsSlowThreshold,
			want:     true,
		},
		{
			name:     "fast getattr enoent is quiet",
			op:       "Getattr",
			errc:     -fuse.ENOENT,
			duration: time.Millisecond,
			want:     false,
		},
		{
			name:     "slow getattr enoent is logged",
			op:       "Getattr",
			errc:     -fuse.ENOENT,
			duration: mountDiagnosticsSlowThreshold,
			want:     true,
		},
		{
			name:     "fast create enoent is logged",
			op:       "Create",
			errc:     -fuse.ENOENT,
			duration: time.Millisecond,
			want:     true,
		},
		{
			name:     "fast permission error is logged",
			op:       "Open",
			errc:     -fuse.EACCES,
			duration: time.Millisecond,
			want:     true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := shouldLogMountCallback(test.op, test.errc, test.duration)
			if got != test.want {
				t.Fatalf("shouldLogMountCallback(%q, %d, %s) = %t, want %t", test.op, test.errc, test.duration, got, test.want)
			}
		})
	}
}
