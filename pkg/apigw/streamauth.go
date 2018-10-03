package apigw

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	streamauth "github.com/tcfw/evntsrc/pkg/streamauth/protos"
	"google.golang.org/grpc"
)

const streamauthEndpoint = "streamauth:443"

func registerStreamAuth(ctx context.Context, mux *runtime.ServeMux, opts []grpc.DialOption) {
	err := streamauth.RegisterStreamAuthServiceHandlerFromEndpoint(ctx, mux, streamauthEndpoint, opts)
	if err != nil {
		panic(err)
	}
}
