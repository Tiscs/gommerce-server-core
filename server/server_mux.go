package server

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

// ServerMuxRoute represents a route for server mux of gRPC-Gateway server.
type ServerMuxRoute struct {
	// Methods is a list of HTTP methods, which are allowed for the route.
	Methods []string
	// Pattern is a URL pattern, which is used for matching the route.
	Pattern string
	// Handler is a handler function, which is called when the route is matched.
	Handler runtime.HandlerFunc
}
