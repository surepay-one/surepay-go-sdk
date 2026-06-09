package core

// Logger is a minimal structured-log interface compatible with *log/slog.Logger.
type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type NoopLogger struct{}

func (NoopLogger) Debug(_ string, _ ...any) {}
func (NoopLogger) Info(_ string, _ ...any)  {}
func (NoopLogger) Warn(_ string, _ ...any)  {}
func (NoopLogger) Error(_ string, _ ...any) {}
