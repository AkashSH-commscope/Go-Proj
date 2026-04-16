package grpcserver

import (
	"context"
	telemetryv1 "example.com/edge-telemetry-bridge/gen/telemetry/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)
type Server struct {
	telemetryv1.UnimplementedTelemetryServiceServer
}
func (s *Server) ListRecentReadings(ctx context.Context, req *telemetryv1.ListRecentReadingsRequest) (*telemetryv1.ListRecentReadingsResponse, error) {

	return nil, status.Errorf(codes.Unimplemented, "method ListRecentReadings not implemented")
}
func (s *Server) SubscribeReadings(req *telemetryv1.SubscribeReadingsRequest, stream telemetryv1.TelemetryService_SubscribeReadingsServer) error {
	<-stream.Context().Done()
	return nil
}