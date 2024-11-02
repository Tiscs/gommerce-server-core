package logging

import (
	"log/slog"
	"sync/atomic"

	"github.com/choral-io/gommerce-server-core/config"
)

type loggerWrapper struct {
	Logger
}

var (
	defaultLogger atomic.Value
)

func init() {
	defaultLogger.Store(loggerWrapper{Logger: &SlogLogger{logger: slog.Default()}})
}

// NewLogger creates a new Logger instance with the given config.
func NewLogger(cfg config.LoggingConfig) (l Logger, err error) {
	zcfg := cfg.GetZapLogger()
	if zcfg != nil {
		if l, err = NewZapLogger(zcfg.GetPreset()); err != nil {
			return nil, err
		}
	}
	scfg := cfg.GetSlogLogger()
	if scfg != nil {
		if l, err = NewSlogLogger(scfg.GetHandler(), scfg.GetAddSource(), scfg.GetLeveler()); err != nil {
			return nil, err
		}
	}
	if l == nil {
		l = &SlogLogger{logger: slog.Default()}
	}
	return
}

// SetDefaultLogger sets the default logger.
func SetDefaultLogger(l Logger) {
	defaultLogger.Store(loggerWrapper{Logger: l})
}

// DefaultLogger returns the default logger.
// If no logger has been set, it returns the default slog logger.
func DefaultLogger() Logger {
	return defaultLogger.Load().(loggerWrapper).Logger
}
