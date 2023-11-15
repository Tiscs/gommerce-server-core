package logging

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"runtime"
	"time"
)

// SlogLogger is a Logger implementation that uses slog.
type SlogLogger struct {
	logger *slog.Logger
}

var _ Logger = (*SlogLogger)(nil)

func NewSlogLogger(handler string, addSource bool, level slog.Leveler) (*SlogLogger, error) {
	options := &slog.HandlerOptions{
		AddSource: addSource,
		Level:     level,
	}
	if handler == "text" {
		slog.SetDefault(slog.New(&wrappedSlogHandler{
			innerHandler: slog.NewTextHandler(os.Stderr, options),
			extractAttrs: []func(context.Context) []slog.Attr{extractTracingAttrs},
		}))
	} else if handler == "json" {
		slog.SetDefault(slog.New(&wrappedSlogHandler{
			innerHandler: slog.NewJSONHandler(os.Stderr, options),
			extractAttrs: []func(context.Context) []slog.Attr{extractTracingAttrs},
		}))
	} else {
		return nil, errors.New("unknown logging handler")
	}
	return &SlogLogger{logger: slog.Default()}, nil
}

func (l *SlogLogger) With(args ...any) Logger {
	return &SlogLogger{logger: l.logger.With(args...)}
}

func (l *SlogLogger) log(ctx context.Context, level Level, message string, args ...any) {
	if !l.logger.Enabled(ctx, slog.Level(level)) {
		return
	}
	var pc uintptr
	var pcs [1]uintptr
	runtime.Callers(3, pcs[:])
	pc = pcs[0]
	r := slog.NewRecord(time.Now(), slog.Level(level), message, pc)
	r.Add(args...)
	if ctx == nil {
		ctx = context.Background()
	}
	_ = l.logger.Handler().Handle(ctx, r)
}

func (l *SlogLogger) Log(ctx context.Context, level Level, message string, args ...any) {
	l.log(ctx, level, message, args...)
}

func (l *SlogLogger) Debug(ctx context.Context, message string, args ...any) {
	l.log(ctx, LevelDebug, message, args...)
}

func (l *SlogLogger) Info(ctx context.Context, message string, args ...any) {
	l.log(ctx, LevelInfo, message, args...)
}

func (l *SlogLogger) Warn(ctx context.Context, message string, args ...any) {
	l.log(ctx, LevelWarn, message, args...)
}

func (l *SlogLogger) Error(ctx context.Context, message string, args ...any) {
	l.log(ctx, LevelError, message, args...)
}

func (l *SlogLogger) Panic(ctx context.Context, message string, args ...any) {
	l.log(ctx, LevelPanic, message, args...)
	panic(message)
}

func (l *SlogLogger) Fatal(ctx context.Context, message string, args ...any) {
	l.log(ctx, LevelFatal, message, args...)
	os.Exit(1)
}

type wrappedSlogHandler struct {
	innerHandler slog.Handler
	extractAttrs []func(context.Context) []slog.Attr
}

var _ slog.Handler = (*wrappedSlogHandler)(nil)

func (h *wrappedSlogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.innerHandler.Enabled(ctx, level)
}

func (h *wrappedSlogHandler) Handle(ctx context.Context, record slog.Record) error {
	for _, extract := range h.extractAttrs {
		if attrs := extract(ctx); attrs != nil {
			record.AddAttrs(attrs...)
		}
	}
	return h.innerHandler.Handle(ctx, record)
}

func (h *wrappedSlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &wrappedSlogHandler{innerHandler: h.innerHandler.WithAttrs(attrs)}
}

func (h *wrappedSlogHandler) WithGroup(name string) slog.Handler {
	return &wrappedSlogHandler{innerHandler: h.innerHandler.WithGroup(name)}
}
