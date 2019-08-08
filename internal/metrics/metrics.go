package metrics

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/opentracing/opentracing-go"

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

//EventsCount fetches metrics relating to event count from timeseries services
func (s *Server) EventsCount(ctx context.Context, req *pb.MetricsRequest) (*pb.MetricsResponse, error) {
	return s.processByQuery(ctx, req, fmt.Sprintf(`sum(increase(storer_store_request_count{stream="%d"}[2m]))`, req.Stream))
}

//EventsSize fetches metrics relating to event byte size from timeseries services
func (s *Server) EventsSize(ctx context.Context, req *pb.MetricsRequest) (*pb.MetricsResponse, error) {
	return s.processByQuery(ctx, req, fmt.Sprintf(`sum(increase(event_subscribe_byte_count{stream="%d"}[2m])+increase(event_publish_byte_count{stream="%d"}[2m]))`, req.Stream, req.Stream))
}

//Connections fetches metrics relating to socket count from timeseries services
func (s *Server) Connections(ctx context.Context, req *pb.MetricsRequest) (*pb.MetricsResponse, error) {
	return s.processByQuery(ctx, req, fmt.Sprintf(`sum(ws_conn{stream="%d"})`, req.Stream))
}

func (s *Server) processByQuery(ctx context.Context, req *pb.MetricsRequest, query string) (*pb.MetricsResponse, error) {
	if err := s.canAccess(ctx, req.Stream); err != nil {
		return nil, err
	}

	promClient, err := promCli.NewClient(promCli.Config{Address: os.Getenv("PROM_HOST"), RoundTripper: promCli.DefaultRoundTripper})
	if err != nil {
		return nil, err
	}

	api := promApi.NewAPI(promClient)

	interval, resolution := apiInterval(&req.Interval)

	span := tracing.StartChildSpan(opentracing.SpanFromContext(ctx), "Prometheus")
	modelVal, _, err := api.QueryRange(ctx, query, promApi.Range{Start: interval, End: time.Now(), Step: resolution})
	if err != nil {
		return nil, err
	}

	matrix := modelVal.(model.Matrix)
	metrics := matrix[0]
	span.Finish()

	vals := []*pb.MetricCount{}

	for _, sample := range metrics.Values {
		vals = append(vals, &pb.MetricCount{
			Count: float32(sample.Value),
			Timestamp: &pb.MetricCount_Timestamp{
				Seconds: sample.Timestamp.Time().Unix(),
			},
		})
	}

	return &pb.MetricsResponse{Metrics: vals}, nil
}

func apiInterval(reqInterval *pb.Interval) (time.Time, time.Duration) {
	interval := time.Now()
	resolution := 30 * time.Second

	switch *reqInterval {
	case pb.Interval_min10:
		interval = interval.Add(-10 * time.Minute)
		break
	case pb.Interval_min30:
		interval = interval.Add(-30 * time.Minute)
		break
	case pb.Interval_hour:
		interval = interval.Add(-time.Hour)
		break
	case pb.Interval_hour12:
		interval = interval.Add(-12 * time.Hour)
		break
	case pb.Interval_day:
		interval = interval.Add(-24 * time.Hour)
		resolution = 30 * time.Minute
		break
	case pb.Interval_week:
		interval = interval.Add(-24 * 7 * time.Hour)
		resolution = 4 * time.Hour
		break
	case pb.Interval_month:
		interval = interval.Add(-24 * 31 * time.Hour)
		resolution = 24 * time.Hour
		break
	}

	return interval, resolution
}

var (
	streamsConn *grpc.ClientConn
)

func (s *Server) canAccess(ctx context.Context, stream int32) error {
	md, _ := metadata.FromIncomingContext(ctx)
	ctxOg := metadata.NewOutgoingContext(ctx, md)
	opts := tracing.GRPCClientOptions()

	if streamsConn == nil {
		conn, err := grpc.Dial("streams:443", opts...)
		if err != nil {
			return err
		}
		streamsConn = conn
	}

	streamsClient := streams.NewStreamsServiceClient(streamsConn)

	_, err := streamsClient.Get(ctxOg, &streams.GetRequest{ID: stream})
	return err
}
