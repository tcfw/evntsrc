package metrics

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/prometheus/common/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	promApi "github.com/prometheus/client_golang/api/prometheus/v1"

	promCli "github.com/prometheus/client_golang/api"
	pb "github.com/tcfw/evntsrc/internal/metrics/protos"
	streams "github.com/tcfw/evntsrc/internal/streams/protos"
	"github.com/tcfw/evntsrc/internal/tracing"
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
	if err := s.canAccess(ctx, req.Stream); err != nil {
		return nil, err
	}

	promClient, err := promCli.NewClient(promCli.Config{Address: os.Getenv("PROM_HOST"), RoundTripper: promCli.DefaultRoundTripper})
	if err != nil {
		return nil, err
	}

	api := promApi.NewAPI(promClient)

	interval := time.Now()
	resolution := 30 * time.Second

	switch req.Interval {
	case pb.MetricsCountRequest_min10:
		interval = interval.Add(-10 * time.Minute)
		break
	case pb.MetricsCountRequest_min30:
		interval = interval.Add(-30 * time.Minute)
		break
	case pb.MetricsCountRequest_hour:
		interval = interval.Add(-time.Hour)
		break
	case pb.MetricsCountRequest_hour12:
		interval = interval.Add(-12 * time.Hour)
		break
	case pb.MetricsCountRequest_day:
		interval = interval.Add(-24 * time.Hour)
		resolution = 30 * time.Minute
		break
	case pb.MetricsCountRequest_week:
		interval = interval.Add(-24 * 7 * time.Hour)
		resolution = 4 * time.Hour
		break
	case pb.MetricsCountRequest_month:
		interval = interval.Add(-24 * 31 * time.Hour)
		resolution = 24 * time.Hour
		break
	}

	modelVal, err := api.QueryRange(ctx, fmt.Sprintf(`sum(increase(storer_store_request_count{stream="%d"}[2m]))`, req.Stream), promApi.Range{Start: interval, End: time.Now(), Step: resolution})
	if err != nil {
		return nil, err
	}

	matrix := modelVal.(model.Matrix)
	metrics := matrix[0]

	vals := []*pb.MetricCount{}

	for _, sample := range metrics.Values {
		vals = append(vals, &pb.MetricCount{
			Count: float32(sample.Value),
			Timestamp: &pb.MetricCount_Timestamp{
				Seconds: sample.Timestamp.Time().Unix(),
			},
		})
	}

	return &pb.MetricsCountResponse{Metrics: vals}, nil
}

func (s *Server) canAccess(ctx context.Context, stream int32) error {
	md, _ := metadata.FromIncomingContext(ctx)
	ctxOg := metadata.NewOutgoingContext(ctx, md)
	opts := tracing.GRPCClientOptions()

	conn, err := grpc.DialContext(ctxOg, "streams:443", opts...)
	if err != nil {
		return err
	}

	streamsClient := streams.NewStreamsServiceClient(conn)

	_, err = streamsClient.Get(ctxOg, &streams.GetRequest{ID: stream})
	return err
}
