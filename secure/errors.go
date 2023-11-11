package secure

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Status represents an error status.
type Status string

const (
	// ErrInvalidToken is returned when the token is invalid.
	ErrInvalidToken = Status("token is invalid")
	// ErrExpiredToken is returned when the token is expired.
	ErrExpiredToken = Status("token is expired")
	// ErrInvalidTokenType is returned when the token type is invalid.
	ErrInvalidTokenType = Status("token type is invalid")
	// ErrInvalidTokenSignature is returned when the token signature is invalid.
	ErrInvalidTokenSignature = Status("token signature is invalid")
	// ErrUnsupportedSigningMethod is returned when the signing method is not supported.
	ErrUnsupportedSigningMethod = Status("unsupported signing method")
	// ErrUnsupportedTokenType is returned when the token type is not supported.
	ErrUnsupportedTokenType = Status("unsupported token type")
	// ErrUnsupportedOperation is returned when the operation is not supported.
	ErrUnsupportedOperation = Status("unsupported operation")
	// ErrInvalidAuthExprOutput is returned when the authorization expression does not return a boolean.
	ErrInvalidAuthExprOutput = Status("authorization expression must return a boolean")
	// ErrUnauthenticated is returned when the user is not authenticated.
	ErrUnauthenticated = Status("unauthenticated")
	// ErrPermissionDenied is returned when the user does not have permission to perform the operation.
	ErrPermissionDenied = Status("permission denied")
)

func (s Status) Error() string {
	return string(s)
}

// GRPCStatus returns the gRPC status for the error.
// Implements the GRPCStatus() method, see status.FromError(error).
func (s Status) GRPCStatus() *status.Status {
	code := codes.Unknown
	switch s {
	case
		ErrInvalidToken,
		ErrExpiredToken,
		ErrInvalidTokenType,
		ErrInvalidTokenSignature,
		ErrUnauthenticated:
		code = codes.Unauthenticated
	case
		ErrPermissionDenied:
		code = codes.PermissionDenied
	}
	return status.New(code, s.Error())
}
