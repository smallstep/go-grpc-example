package health

import (
	"context"

	health "google.golang.org/grpc/health/grpc_health_v1"
)

// Server implements the GRPC Health Checking Protocol.
type Server struct {
	health.HealthServer
}

// New returns a new health server.
func New() *Server {
	return &Server{
		HealthServer: &health.UnimplementedHealthServer{},
	}
}

// Check implements the GRPC Health Checking Protocol and returns always serving.
func (s *Server) Check(ctx context.Context, in *health.HealthCheckRequest) (*health.HealthCheckResponse, error) {
	return &health.HealthCheckResponse{
		Status: health.HealthCheckResponse_SERVING,
	}, nil
}
