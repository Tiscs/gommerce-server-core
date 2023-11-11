package logging

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func extractTracingAttrs(ctx context.Context) []slog.Attr {
	if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
		return []slog.Attr{
			slog.String("trace.id", span.TraceID().String()),
			slog.String("span.id", span.SpanID().String()),
		}
	}
	return nil
}

func extractTracingFields(ctx context.Context) []zap.Field {
	if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
		return []zap.Field{
			zap.String("trace.id", span.TraceID().String()),
			zap.String("span.id", span.SpanID().String()),
		}
	}
	return nil
}
