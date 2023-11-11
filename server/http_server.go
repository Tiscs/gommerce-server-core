package server

import (
	"context"
	"net"
	"net/http"

	"github.com/choral-io/gommerce-server-core/config"
	"github.com/choral-io/gommerce-server-core/logging"
)

// HTTPServer is an implementation of Server for HTTP.
type HTTPServer struct {
	server *http.Server
	logger logging.Logger
	chDone chan error
}

var _ Server = (*HTTPServer)(nil)

// NewHTTPServer returns a new HTTPServer with the given config, logger, and handler.
func NewHTTPServer(cfg config.ServerHTTPConfig, logger logging.Logger, handler http.Handler) *HTTPServer {
	return &HTTPServer{
		server: &http.Server{Addr: cfg.GetAddr(), Handler: handler},
		logger: logger,
		chDone: make(chan error, 1),
	}
}

func (s *HTTPServer) Start(ctx context.Context) error {
	ln, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		return err
	}
	go func() {
		s.logger.Info(ctx, "serving http", "addr", s.server.Addr)
		if err := s.server.Serve(ln); err != nil && err != http.ErrServerClosed {
			s.logger.Error(ctx, "error while serving http", "error", err)
			s.chDone <- err
		}
		s.logger.Info(ctx, "http server stopped")
		close(s.chDone)
	}()
	return nil
}

func (s *HTTPServer) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *HTTPServer) Done() <-chan error {
	return s.chDone
}
