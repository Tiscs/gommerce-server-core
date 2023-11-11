package logging

import (
	"context"
)

// Level is a logging level.
type Level int8

const (
	LevelDebug Level = -4
	LevelInfo  Level = 0
	LevelWarn  Level = 4
	LevelError Level = 8
	LevelPanic Level = 16
	LevelFatal Level = 20
)

// Logger used to log messages.
type Logger interface {
	// With returns a new Logger with args added to the logger's context.
	With(args ...any) Logger
	// Log logs a message with args at level.
	Log(ctx context.Context, level Level, message string, args ...any)
	// Debug logs a message at level Debug.
	Debug(ctx context.Context, message string, args ...any)
	// Info logs a message at level Info.
	Info(ctx context.Context, message string, args ...any)
	// Warn logs a message at level Warn.
	Warn(ctx context.Context, message string, args ...any)
	// Error logs a message at level Error.
	Error(ctx context.Context, message string, args ...any)
	// Panic logs a message at level Panic.
	Panic(ctx context.Context, message string, args ...any)
	// Fatal logs a message at level Fatal.
	Fatal(ctx context.Context, message string, args ...any)
}
