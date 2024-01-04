package data

import (
	"github.com/choral-io/gommerce-server-core/config"
	"github.com/redis/rueidis"
	"github.com/redis/rueidis/rueidisotel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// NewRedisClient creates a new redis client with tracing and metrics.
func NewRedisClient(cfg config.ServerRedisConfig, tp trace.TracerProvider, mp metric.MeterProvider) (rueidis.Client, error) {
	return rueidisotel.NewClient(rueidis.ClientOption{
		InitAddress: cfg.GetInitAddress(),
		SelectDB:    cfg.GetSelectDB(),
	}, rueidisotel.WithTracerProvider(tp), rueidisotel.WithMeterProvider(mp))
}
