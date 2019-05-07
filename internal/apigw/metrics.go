package apigw

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	metrics "github.com/tcfw/evntsrc/internal/metrics/protos"
	"google.golang.org/grpc"
)

func registerMetrics(ctx context.Context, mux *runtime.ServeMux, opts []grpc.DialOption) {
	if err := metrics.RegisterMetricsServiceHandlerFromEndpoint(ctx, mux, "metrics:443", opts); err != nil {
		panic(err)
	}
}
