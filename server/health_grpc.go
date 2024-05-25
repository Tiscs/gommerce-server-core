package server

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/redis/rueidis"
	"github.com/uptrace/bun"
	"google.golang.org/grpc"
	health "google.golang.org/grpc/health/grpc_health_v1"
)

type healthServiceServer struct {
	health.UnimplementedHealthServer

	sdb *sql.DB
	rdb rueidis.Client
}

// NewHealthServiceServer returns a new health service server.
func NewHealthServiceServer(bdb bun.IDB, rdb rueidis.Client) health.HealthServer {
	// this is a hack to get the underlying sql.DB from bun.IDB
	sdb := bdb.NewSelect().DB().DB
	return &healthServiceServer{
		sdb: sdb,
		rdb: rdb,
	}
}

// RegisterServerService implements ServerServiceRegister.
func (s *healthServiceServer) RegisterServerService(reg grpc.ServiceRegistrar) {
	reg.RegisterService(&health.Health_ServiceDesc, s)
}

func (s *healthServiceServer) check(ctx context.Context, _ *health.HealthCheckRequest) error {
	if s.sdb != nil {
		if err := s.sdb.PingContext(ctx); err != nil {
			slog.WarnContext(ctx, "db ping error", "error", err)
			return err
		}
	}
	if s.rdb != nil {
		if err := s.rdb.Do(context.Background(), s.rdb.B().Ping().Build()).Error(); err != nil {
			slog.WarnContext(ctx, "redis ping error", "error", err)
			return err
		}
	}
	return nil
}

// Check implements health.HealthServer.
// It checks the health of the server, and returns NOT_SERVING if the server is not healthy.
func (s *healthServiceServer) Check(ctx context.Context, req *health.HealthCheckRequest) (*health.HealthCheckResponse, error) {
	status := health.HealthCheckResponse_NOT_SERVING
	if s.check(ctx, req) == nil {
		status = health.HealthCheckResponse_SERVING
	}
	return &health.HealthCheckResponse{Status: status}, nil
}

func (s *healthServiceServer) reply(req *health.HealthCheckRequest, srv health.Health_WatchServer) error {
	status := health.HealthCheckResponse_NOT_SERVING
	if s.check(srv.Context(), req) == nil {
		status = health.HealthCheckResponse_SERVING
	}
	if err := srv.Send(&health.HealthCheckResponse{Status: status}); err != nil {
		slog.WarnContext(srv.Context(), "failed to send health check response", "error", err)
		return err
	}
	return nil
}

// Watch implements health.HealthServer.
// It checks the health of the server, and returns NOT_SERVING if the server is not healthy.
// It also sends a health check response every 30 seconds.
func (s *healthServiceServer) Watch(req *health.HealthCheckRequest, srv health.Health_WatchServer) error {
	if err := s.reply(req, srv); err != nil {
		return err
	}
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := s.reply(req, srv); err != nil {
				return err
			}
		case <-srv.Context().Done():
			return srv.Context().Err()
		}
	}
}
