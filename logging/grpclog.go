package logging

import (
	"context"
	"io"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func errorToCode(err error) codes.Code {
	switch err {
	case context.DeadlineExceeded:
		return codes.DeadlineExceeded
	case context.Canceled:
		return codes.Canceled
	case io.ErrUnexpectedEOF:
		return codes.Internal
	}
	return status.Code(err)
}

// GRPCLogger provides a grpc middleware that logs grpc calls.
type GRPCLogger struct {
	logger Logger
	opts   []logging.Option
}

// NewGRPCLogger creates a new GRPCLogger with logger.
func NewGRPCLogger(logger Logger) *GRPCLogger {
	return &GRPCLogger{
		logger: logger,
		opts: []logging.Option{
			logging.WithLogOnEvents(logging.StartCall, logging.FinishCall), // logging.PayloadReceived, logging.PayloadSent
			logging.WithCodes(errorToCode),
		},
	}
}

// Log implements logging.Logger.
func (l *GRPCLogger) Log(ctx context.Context, level logging.Level, msg string, fields ...any) {
	l.logger.Log(ctx, Level(level), msg, fields...)
}

// UnaryServerInterceptor returns a grpc.UnaryServerInterceptor that logs grpc calls.
func (l *GRPCLogger) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return logging.UnaryServerInterceptor(l, l.opts...)
}

// StreamServerInterceptor returns a grpc.StreamServerInterceptor that logs grpc calls.
func (l *GRPCLogger) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return logging.StreamServerInterceptor(l, l.opts...)
}

// UnaryClientInterceptor returns a grpc.UnaryClientInterceptor that logs grpc calls.
func (l *GRPCLogger) UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return logging.UnaryClientInterceptor(l, l.opts...)
}

// StreamClientInterceptor returns a grpc.StreamClientInterceptor that logs grpc calls.
func (l *GRPCLogger) StreamClientInterceptor() grpc.StreamClientInterceptor {
	return logging.StreamClientInterceptor(l, l.opts...)
}
