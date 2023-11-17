package logging

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

// GrpcLogger provides a grpc middleware that logs grpc calls.
type GrpcLogger struct {
	logger Logger
	opts   []logging.Option
}

// NewGrpcLogger creates a new GrpcLogger with logger.
func NewGrpcLogger(logger Logger) *GrpcLogger {
	return &GrpcLogger{
		logger: logger,
		opts: []logging.Option{
			logging.WithLogOnEvents(logging.StartCall, logging.FinishCall), // logging.PayloadReceived, logging.PayloadSent
		},
	}
}

// Log implements logging.Logger.
func (l *GrpcLogger) Log(ctx context.Context, level logging.Level, msg string, fields ...any) {
	l.logger.Log(ctx, Level(level), msg, fields...)
}

// UnaryServerInterceptor returns a grpc.UnaryServerInterceptor that logs grpc calls.
func (l *GrpcLogger) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return logging.UnaryServerInterceptor(l, l.opts...)
}

// StreamServerInterceptor returns a grpc.StreamServerInterceptor that logs grpc calls.
func (l *GrpcLogger) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return logging.StreamServerInterceptor(l, l.opts...)
}

// UnaryClientInterceptor returns a grpc.UnaryClientInterceptor that logs grpc calls.
func (l *GrpcLogger) UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return logging.UnaryClientInterceptor(l, l.opts...)
}

// StreamClientInterceptor returns a grpc.StreamClientInterceptor that logs grpc calls.
func (l *GrpcLogger) StreamClientInterceptor() grpc.StreamClientInterceptor {
	return logging.StreamClientInterceptor(l, l.opts...)
}
