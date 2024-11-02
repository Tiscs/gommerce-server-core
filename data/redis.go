package data

import (
	"strings"

	"github.com/redis/rueidis"
	"github.com/redis/rueidis/rueidisotel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"

	"github.com/choral-io/gommerce-server-core/config"
)

// NewRedisClient creates a new redis client with tracing and metrics.
func NewRedisClient(cfg config.ServerRedisConfig, tp trace.TracerProvider, mp metric.MeterProvider) (rueidis.Client, error) {
	return rueidisotel.NewClient(rueidis.ClientOption{
		InitAddress: processInitAddress(cfg.GetInitAddr()),
		SelectDB:    cfg.GetSelectDB(),
	}, rueidisotel.WithTracerProvider(tp), rueidisotel.WithMeterProvider(mp))
}

func processInitAddress(url string) []string {
	addr := strings.Split(url, ",")
	var j int
	for _, s := range addr {
		u := strings.TrimSpace(s)
		if len(u) > 0 {
			addr[j] = u
			j++
		}
	}
	return addr[:j]
}
