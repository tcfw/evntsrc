package apigw

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	streams "github.com/tcfw/evntsrc/internal/streams/protos"
	"google.golang.org/grpc"
)

func registerStreams(ctx context.Context, mux *runtime.ServeMux, opts []grpc.DialOption) {
	if err := streams.RegisterStreamsServiceHandlerFromEndpoint(ctx, mux, "streams:443", opts); err != nil {
		panic(err)
	}
}
