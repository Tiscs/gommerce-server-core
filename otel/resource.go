package otel

import (
	"context"

	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"

	"github.com/choral-io/gommerce-server-core/config"
)

// NewServerResource creates a new resource with the given config.
func NewServerResource(cfg config.ServerConfig) (*resource.Resource, error) {
	return resource.New(context.Background(), resource.WithAttributes(
		semconv.ServiceName(cfg.GetName()),
		semconv.ServiceVersion(cfg.GetVersion()),
		semconv.ServiceInstanceID(cfg.GetInstanceId()),
	))
}
