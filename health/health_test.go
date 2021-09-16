package health

import (
	"context"
	"reflect"
	"testing"

	health "google.golang.org/grpc/health/grpc_health_v1"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *Server
	}{
		{"ok", &Server{HealthServer: &health.UnimplementedHealthServer{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_Check(t *testing.T) {
	type fields struct {
		HealthServer health.HealthServer
	}
	type args struct {
		ctx context.Context
		in  *health.HealthCheckRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *health.HealthCheckResponse
		wantErr bool
	}{
		{"ok", fields{&health.UnimplementedHealthServer{}}, args{
			context.Background(), &health.HealthCheckRequest{Service: "foo"},
		}, &health.HealthCheckResponse{
			Status: health.HealthCheckResponse_SERVING,
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				HealthServer: tt.fields.HealthServer,
			}
			got, err := s.Check(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Server.Check() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server.Check() = %v, want %v", got, tt.want)
			}
		})
	}
}
