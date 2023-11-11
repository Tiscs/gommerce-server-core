package logging

import (
	"log/slog"

	"github.com/choral-io/gommerce-server-core/config"
)

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
