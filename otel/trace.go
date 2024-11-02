package otel

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"

	"github.com/choral-io/gommerce-server-core/config"
)

// NewTracerProvider creates a new TracerProvider instance with the given config.
func NewTracerProvider(cfg config.TraceConfig, res *resource.Resource) (trace.TracerProvider, error) {
	ctx := context.Background()
	protocol := cfg.GetExporterConfig().GetProtocol()
	var exporter sdktrace.SpanExporter
	var err error
	if protocol == "otlp-grpc" {
		exporter, err = otlptracegrpc.New(ctx,
			otlptracegrpc.WithEndpoint(cfg.GetExporterConfig().GetEndpoint()),
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithTimeout(2*time.Second),
		)
	} else if protocol == "otlp-http" {
		exporter, err = otlptracehttp.New(ctx,
			otlptracehttp.WithEndpoint(cfg.GetExporterConfig().GetEndpoint()),
			otlptracehttp.WithInsecure(),
			otlptracehttp.WithTimeout(2*time.Second),
		)
	} else if protocol == "stdout" {
		exporter, err = stdout.New(stdout.WithPrettyPrint())
	} else if protocol == "noop" {
		exporter = tracetest.NewNoopExporter()
	} else {
		return nil, fmt.Errorf("invalid trace exporter protocol: %s", protocol)
	}
	if err != nil {
		return nil, err
	}
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(exporter),
	)
	return tracerProvider, nil
}
