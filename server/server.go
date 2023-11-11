package server

import (
	"context"
)

// Server is implemented by servers.
type Server interface {
	// Start starts the server.
	Start(context.Context) error
	// Stop stops the server.
	Stop(context.Context) error
	// Done returns a channel that is closed when the server is stopped.
	// An error is sent to the channel if the server is stopped with an error.
	Done() <-chan error
}
