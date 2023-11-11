package logging

import (
	"context"
	"errors"
	"log/slog"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const badKey = "!BADKEY"

// ZapLogger provides a logger that uses zap.
type ZapLogger struct {
	logger *zap.Logger
}

var _ Logger = (*ZapLogger)(nil)

func NewZapLogger(preset string) (*ZapLogger, error) {
	var f func(options ...zap.Option) (*zap.Logger, error)
	if preset == "development" {
		f = zap.NewDevelopment
	} else if preset == "production" {
		f = zap.NewProduction
	} else {
		return nil, errors.New("unknown logging preset")
	}
	if l, err := f(zap.AddCallerSkip(2)); err != nil {
		return nil, err
	} else {
		zap.ReplaceGlobals(l)
	}
	return &ZapLogger{logger: zap.L()}, nil
}

func argsToFields(args []any) []zap.Field {
	fields := []zap.Field{}
	for len(args) > 0 {
		switch x := args[0].(type) {
		case string:
			if len(args) == 1 {
				fields = append(fields, zap.String(badKey, x))
				args = nil
				continue
			}
			fields = append(fields, zap.Any(x, args[1]))
			args = args[2:]
		case zap.Field:
			fields = append(fields, x)
			args = args[1:]
		case slog.Attr:
			fields = append(fields, zap.Any(x.Key, x.Value))
			args = args[1:]
		default:
			fields = append(fields, zap.Any(badKey, x))
			args = args[1:]
		}
	}
	return fields
}

func (l *ZapLogger) With(args ...any) Logger {
	return &ZapLogger{logger: l.logger.With(argsToFields(args)...)}
}

func (l *ZapLogger) log(ctx context.Context, level Level, message string, args ...any) {
	fields := argsToFields(args)
	if fs := extractTracingFields(ctx); len(fs) > 0 {
		fields = append(fields, fs...)
	}
	l.logger.Log(zapcore.Level(level/4), message, fields...)
}

func (l *ZapLogger) Log(ctx context.Context, level Level, message string, args ...any) {
	l.log(ctx, level, message, args...)
}

func (l *ZapLogger) Debug(ctx context.Context, message string, args ...any) {
	l.log(ctx, LevelDebug, message, args...)
}

func (l *ZapLogger) Info(ctx context.Context, message string, args ...any) {
	l.log(ctx, LevelInfo, message, args...)
}

func (l *ZapLogger) Warn(ctx context.Context, message string, args ...any) {
	l.log(ctx, LevelWarn, message, args...)
}

func (l *ZapLogger) Error(ctx context.Context, message string, args ...any) {
	l.log(ctx, LevelError, message, args...)
}

func (l *ZapLogger) Panic(ctx context.Context, message string, args ...any) {
	l.log(ctx, LevelPanic, message, args...)
}

func (l *ZapLogger) Fatal(ctx context.Context, message string, args ...any) {
	l.log(ctx, LevelFatal, message, args...)
}
