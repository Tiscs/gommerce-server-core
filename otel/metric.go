package otel

import (
	"context"
	"fmt"
	"time"

	"github.com/choral-io/gommerce-server-core/config"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
)

// NewMeterProvider creates a new MeterProvider instance with the given config.
func NewMeterProvider(cfg config.MetricConfig, res *resource.Resource) (metric.MeterProvider, error) {
	ctx := context.Background()
	protocol := cfg.GetExporterConfig().GetProtocol()
	var exporter sdkmetric.Exporter
	var err error
	if protocol == "otlp-grpc" {
		exporter, err = otlpmetricgrpc.New(ctx,
			otlpmetricgrpc.WithEndpoint(cfg.GetExporterConfig().GetEndpoint()),
			otlpmetricgrpc.WithInsecure(),
			otlpmetricgrpc.WithTimeout(2*time.Second),
		)
	} else if protocol == "otlp-http" {
		exporter, err = otlpmetrichttp.New(ctx,
			otlpmetrichttp.WithEndpoint(cfg.GetExporterConfig().GetEndpoint()),
			otlpmetrichttp.WithInsecure(),
			otlpmetrichttp.WithTimeout(2*time.Second),
		)
	} else if protocol == "stdout" {
		exporter, err = stdout.New(stdout.WithPrettyPrint())
	} else if protocol == "noop" {
		exporter = nil
	} else {
		return nil, fmt.Errorf("invalid trace exporter protocol: %s", protocol)
	}
	if err != nil {
		return nil, err
	}
	var reader sdkmetric.Reader = nil
	if exporter != nil {
		reader = sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(time.Second*5))
	}
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(reader),
	)
	return meterProvider, nil
}
