package validator

import (
	"context"
	"log/slog"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func grpcStatusErr(ctx context.Context, err error) error {
	var errs []error
	switch err := err.(type) {
	case validationError:
		errs = []error{err}
	case validationMultiError:
		errs = err.AllErrors()
	default:
		errs = nil
	}
	if len(errs) == 0 {
		return err
	}
	st := status.New(codes.InvalidArgument, err.Error())
	br := &errdetails.BadRequest{}
	br.FieldViolations = make([]*errdetails.BadRequest_FieldViolation, len(errs))
	for i, err := range errs {
		if verr, ok := err.(validationError); ok {
			br.FieldViolations[i] = &errdetails.BadRequest_FieldViolation{
				Field:       verr.Field(),
				Description: verr.Reason(),
			}
		} else {
			br.FieldViolations[i] = &errdetails.BadRequest_FieldViolation{
				Field:       "",
				Description: err.Error(),
			}
		}
	}
	st, err = st.WithDetails(br)
	if err != nil {
		slog.ErrorContext(ctx, "failed to attach validation errors", "err", err)
		return status.Error(codes.Unknown, err.Error())
	}
	return st.Err()
}

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if err := Validate(ctx, req, grpcStatusErr); err != nil {
			return nil, err
		}
		if resp, err := handler(ctx, req); err != nil {
			return nil, grpcStatusErr(ctx, err)
		} else {
			return resp, nil
		}
	}
}

func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		wrapper := &recvWrapper{
			ServerStream: stream,
		}
		return handler(srv, wrapper)
	}
}

type recvWrapper struct {
	grpc.ServerStream
}

func (s *recvWrapper) RecvMsg(msg any) error {
	if err := s.ServerStream.RecvMsg(msg); err != nil {
		return err
	}
	if err := Validate(s.Context(), msg, grpcStatusErr); err != nil {
		return err
	}
	return nil
}
