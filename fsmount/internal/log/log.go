// Package log provides logging interfaces and implementations for Files.com FUSE mount.
package log

type Logger interface {
	Debug(format string, v ...any)
	Error(format string, v ...any)
	Info(format string, v ...any)
	Trace(format string, v ...any)
	Warn(format string, v ...any)
}

type NoOpLogger struct{}

func (l *NoOpLogger) Debug(format string, v ...any) {}
func (l *NoOpLogger) Error(format string, v ...any) {}
func (l *NoOpLogger) Info(format string, v ...any)  {}
func (l *NoOpLogger) Trace(format string, v ...any) {}
func (l *NoOpLogger) Warn(format string, v ...any)  {}
