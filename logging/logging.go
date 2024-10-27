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

func SetDefaultLogger(l Logger) {
	defaultLogger.Store(loggerWrapper{Logger: l})
}

func DefaultLogger() Logger {
	return defaultLogger.Load().(loggerWrapper).Logger
}
