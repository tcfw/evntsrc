package metrics

import (
	"context"
	"sync"

	pb "github.com/tcfw/evntsrc/internal/metrics/protos"
)

//Server core struct
type Server struct {
	mu sync.Mutex
}

//NewServer creates a new struct to interface the streams server
func NewServer() *Server {
	return &Server{}
}

//MetricsCount fetches metrics relating to event count from timeseries services
func (s *Server) MetricsCount(ctx context.Context, req *pb.MetricsCountRequest) (*pb.MetricsCountResponse, error) {
	return nil, nil
}
