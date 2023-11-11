package secure

import (
	"context"
	"errors"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/metadata"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// authorizer is implemented by servers that require authorization.
type authorizer interface {
	// Authorize authorizes the identity in the given context to perform the given procedure.
	Authorize(ctx context.Context, procedure string) error
}

// ServerAuthorizer provides server-side grpc interceptors for authorization.
type ServerAuthorizer struct {
	stores map[string]TokenStore
}

// NewServerAuthorizer returns a new ServerAuthorizer with the given token stores.
// The key of the map stores is the authentication schema.
func NewServerAuthorizer(stores map[string]TokenStore) *ServerAuthorizer {
	auth := &ServerAuthorizer{stores: make(map[string]TokenStore, len(stores))}
	for schema, store := range stores {
		auth.stores[strings.ToLower(schema)] = store
	}
	return auth
}

// resolveIdentity resolves the identity from the given authorization header value.
func (auth *ServerAuthorizer) resolveIdentity(ahv string) (*Identity, error) {
	if splits := strings.SplitN(ahv, " ", 2); len(splits) == 2 {
		schema := strings.ToLower(splits[0])
		if store, ok := auth.stores[schema]; ok {
			if t, err := store.Verify(splits[1]); err == nil {
				return &Identity{schema: schema, token: t}, nil
			} else if errors.Is(err, jwt.ErrTokenExpired) {
				return nil, ErrExpiredToken
			} else {
				return nil, err
			}
		}
	}
	return nil, nil
}

// UnaryServerInterceptor returns a grpc.UnaryServerInterceptor that authorizes the identity in the context.
func (auth *ServerAuthorizer) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if user, err := auth.resolveIdentity(metadata.ExtractIncoming(ctx).Get(AUTH_HEADER_KEY)); err != nil {
			return nil, err
		} else if user != nil {
			ctx = context.WithValue(ctx, identityKey{}, user)
		}
		if authorizer, ok := info.Server.(authorizer); ok {
			if err := authorizer.Authorize(ctx, info.FullMethod); err != nil {
				return nil, err
			}
		}
		return handler(ctx, req)
	}
}

// StreamServerInterceptor returns a grpc.StreamServerInterceptor that authorizes the identity in the context.
func (auth *ServerAuthorizer) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if user, err := auth.resolveIdentity(metadata.ExtractIncoming(ss.Context()).Get(AUTH_HEADER_KEY)); err != nil {
			return err
		} else if user != nil {
			ss = &middleware.WrappedServerStream{ServerStream: ss, WrappedContext: context.WithValue(ss.Context(), identityKey{}, user)}
		}
		if authorizer, ok := srv.(authorizer); ok {
			if err := authorizer.Authorize(ss.Context(), info.FullMethod); err != nil {
				return err
			}
		}
		return handler(srv, ss)
	}
}

// ClientAuthorizer provides client-side grpc interceptors for authorization.
type ClientAuthorizer struct {
	cred credentials.PerRPCCredentials
}

// NewClientAuthorizer returns a new ClientAuthorizer with the given credentials.
func NewClientAuthorizer(cred credentials.PerRPCCredentials) *ClientAuthorizer {
	return &ClientAuthorizer{cred: cred}
}

// UnaryClientInterceptor returns a grpc.UnaryClientInterceptor that authorizes the client connection with the given credentials.
func (auth *ClientAuthorizer) UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, conn *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		opts = append(opts, grpc.PerRPCCredentials(auth.cred))
		return invoker(ctx, method, req, reply, conn, opts...)
	}
}

// StreamClientInterceptor returns a grpc.StreamClientInterceptor that authorizes the client connection with the given credentials.
func (auth *ClientAuthorizer) StreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, conn *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		opts = append(opts, grpc.PerRPCCredentials(auth.cred))
		return streamer(ctx, desc, conn, method, opts...)
	}
}
