package logging

import (
	"log/slog"

	"go.uber.org/fx/fxevent"
	"go.uber.org/zap/zapcore"
)

// NewFxeventLogger returns a new fxevent.Logger that logs to logger.
func NewFxeventLogger(logger Logger, eventLevel Level, errorLevel Level) fxevent.Logger {
	// check if logger is *SlogLogger or *ZapLogger
	if s, ok := logger.(*SlogLogger); ok {
		l := &fxevent.SlogLogger{
			Logger: s.logger,
		}
		l.UseLogLevel(slog.Level(eventLevel))
		l.UseErrorLevel(slog.Level(errorLevel))
		return l
	}
	if z, ok := logger.(*ZapLogger); ok {
		l := &fxevent.ZapLogger{
			Logger: z.logger,
		}
		l.UseLogLevel(zapcore.Level(eventLevel / 4))
		l.UseErrorLevel(zapcore.Level(errorLevel / 4))
		return l
	}
	return &fxevent.ConsoleLogger{}
}
