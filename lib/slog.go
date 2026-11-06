package lib

import (
	"context"
	"fmt"
	"log/slog"
)

// SlogLogger sends SDK diagnostics to a caller-owned standard Go logger.
// Printf logs at Info, and Trace logs at Debug. Its handler controls filtering
// and routing, including changes made through a shared slog.LevelVar.
type SlogLogger struct {
	logger *slog.Logger
}

// NewSlogLogger adapts a non-nil logger for Config.Logger. Set it before Config.Init
// so the SDK's default retry client uses the same logger.
func NewSlogLogger(logger *slog.Logger) *SlogLogger {
	return &SlogLogger{logger: logger}
}

// DebugEnabled reports whether the handler enables detailed SDK diagnostics.
func (l *SlogLogger) DebugEnabled() bool {
	return l.logger.Enabled(context.Background(), slog.LevelDebug)
}

func (l *SlogLogger) Printf(format string, args ...any) {
	l.log(slog.LevelInfo, format, args...)
}

func (l *SlogLogger) Error(format string, args ...any) {
	l.log(slog.LevelError, format, args...)
}

func (l *SlogLogger) Warn(format string, args ...any) {
	l.log(slog.LevelWarn, format, args...)
}

func (l *SlogLogger) Info(format string, args ...any) {
	l.log(slog.LevelInfo, format, args...)
}

func (l *SlogLogger) Debug(format string, args ...any) {
	l.log(slog.LevelDebug, format, args...)
}

func (l *SlogLogger) Trace(format string, args ...any) {
	l.log(slog.LevelDebug, format, args...)
}

func (l *SlogLogger) log(level slog.Level, format string, args ...any) {
	ctx := context.Background()
	if l.logger.Enabled(ctx, level) {
		l.logger.Log(ctx, level, fmt.Sprintf(format, args...))
	}
}
