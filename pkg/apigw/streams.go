package apigw

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	streams "github.com/tcfw/evntsrc/pkg/streams/protos"
	"google.golang.org/grpc"
)

func registerStreams(ctx context.Context, mux *runtime.ServeMux, opts []grpc.DialOption) {
	err := streams.RegisterStreamsServiceHandlerFromEndpoint(ctx, mux, "streams:443", opts)
	if err != nil {
		panic(err)
	}
}
